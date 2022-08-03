package webhooks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	nethttp "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
)

const (
	webhooksKey         = "webhooks"
	defaultListenerAddr = ":8086"
)

// Module is a fx module that provides annotated http server using the default listener
// and kubernetes registry. Currently, the module is setting up a separate server with TLS enabled and
// it cannot be done on main server yet because OTEL collector does not support TLS.
func Module() fx.Option {
	return fx.Options(
		tlsconfig.Constructor{Name: webhooksKey, Key: webhooksKey + ".tls"}.Annotate(),
		fx.Provide(listener.Constructor{
			Name: webhooksKey,
			Key:  webhooksKey,
			DefaultConfig: listener.ListenerConfig{
				Addr: defaultListenerAddr,
			},
		}.ProvideAnnotated()),
		nethttp.ServerConstructor{Name: webhooksKey, ListenerName: webhooksKey, TLSConfigName: webhooksKey}.Annotate(),
		fx.Provide(fx.Annotate(
			ProvideK8sRegistry,
			fx.ParamTags(config.NameTag(webhooksKey)),
		)),
	)
}

// ProvideK8sRegistry provides k8s registry which allows serving k8s webhooks
// using given router.
func ProvideK8sRegistry(mux *mux.Router) *K8sRegistry {
	return &K8sRegistry{
		mux:      mux,
		handlers: make(map[string]*handler),
	}
}

// K8sRegistry allows registering k8s webhooks as validators.
type K8sRegistry struct {
	mux      *mux.Router
	handlers map[string]*handler
}

// RegisterValidator adds a validator to be handled on a given HTTP path
//
// This function should be only called before Start phase.
func (r *K8sRegistry) RegisterValidator(path string, validator Validator) {
	oldHandler, exists := r.handlers[path]
	if !exists {
		newHandler := &handler{
			validators: []Validator{validator},
			path:       path,
		}
		r.mux.Handle(path, newHandler)
		// remember the handler so we can add more validators on the same path later
		r.handlers[path] = newHandler
	} else {
		oldHandler.validators = append(oldHandler.validators, validator)
	}
}

// FIXME make me configurable? (note: this should match the timeout on
// ValidatingWebhookConfiguration).
const requestTimeout = 5 * time.Second

// Validator is an interface for k8s validating webhooks.
type Validator interface {
	// ValidateObject validates an object given in admission.Request
	//
	// returns:
	// * true, "", nil if object passes validation
	// * false, msg, nil if object fails validation
	// * false, "", err if validator itself fails
	//
	// Note: if validator is "not interested" in a given object, it should
	// return true.
	ValidateObject(
		ctx context.Context,
		req *admissionv1.AdmissionRequest,
	) (ok bool, msg string, err error)
}

// handler combines all the hooks to be served from concrete HTTP path.
type handler struct {
	path       string // for logging
	validators []Validator
}

// ServeHTTP implements http.Handler interface and invokes the handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Trace().Str("path", h.path).Msg("Webhooks HandleFunc() start")

	if req.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", http.StatusUnsupportedMediaType)
		return
	}
	if req.Body == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var aReq admissionv1.AdmissionReview
	err = json.Unmarshal(body, &aReq)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	reqCtx, cancel := context.WithTimeout(req.Context(), requestTimeout)
	defer cancel()

	aResp := admissionv1.AdmissionReview{
		TypeMeta: aReq.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     aReq.Request.UID,
			Allowed: true,
		},
	}

	var ok bool
	var msg string
	for _, validator := range h.validators {
		ok, msg, err = validator.ValidateObject(reqCtx, aReq.Request)
		if err != nil {
			log.Error().Err(err).Str(h.path, "path").Msg("Validator failed to validate")
			http.Error(
				w, "internal error occurred in validator", http.StatusInternalServerError,
			)
			return
		}

		if !ok {
			aResp.Response.Allowed = false
			aResp.Response.Result = &metav1.Status{
				Message: msg,
			}
			break
		}
	}

	respBody, err := json.Marshal(aResp)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling webhook response body")
		http.Error(w, "", http.StatusInternalServerError)
	}
	_, err = w.Write(respBody)
	if err != nil {
		log.Error().Err(err).Msg("Error writing webhook response body")
		http.Error(w, "", http.StatusInternalServerError)
	}

	if aResp.Response.Allowed {
		log.Trace().Str("path", h.path).Msg("Webhook accepted validation request")
	} else {
		log.Info().Str("path", h.path).Str("msg", aResp.Response.Result.Message).Msg(
			"Validating webhook rejected validation request",
		)
	}
}
