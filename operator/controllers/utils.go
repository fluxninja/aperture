/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// ContainerSecurityContext prepares SecurityContext for containers based on the provided parameter.
func ContainerSecurityContext(containerSecurityContext common.ContainerSecurityContext) *corev1.SecurityContext {
	var securityContext *corev1.SecurityContext
	if containerSecurityContext.Enabled {
		securityContext = &corev1.SecurityContext{
			RunAsUser:              pointer.Int64(containerSecurityContext.RunAsUser),
			RunAsNonRoot:           pointer.Bool(containerSecurityContext.RunAsNonRootUser),
			ReadOnlyRootFilesystem: pointer.Bool(containerSecurityContext.ReadOnlyRootFilesystem),
		}
	} else {
		securityContext = &corev1.SecurityContext{}
	}

	return securityContext
}

// PodSecurityContext prepares SecurityContext for Pods based on the provided parameter.
func PodSecurityContext(podSecurityContext common.PodSecurityContext) *corev1.PodSecurityContext {
	var securityContext *corev1.PodSecurityContext
	if podSecurityContext.Enabled {
		securityContext = &corev1.PodSecurityContext{
			FSGroup: pointer.Int64(podSecurityContext.FsGroup),
		}
	} else {
		securityContext = &corev1.PodSecurityContext{}
	}

	return securityContext
}

// ImageString prepares image string from the provided Image struct.
func ImageString(image common.Image, repository string) string {
	var imageStr string
	if image.Registry != "" {
		imageStr = fmt.Sprintf("%s/%s:%s", image.Registry, repository, image.Tag)
	} else {
		imageStr = fmt.Sprintf("%s:%s", repository, image.Tag)
	}
	return imageStr
}

// ImagePullSecrets prepares ImagePullSecrets string slice from the provided Image struct.
func ImagePullSecrets(image common.Image) []corev1.LocalObjectReference {
	imagePullSecrets := []corev1.LocalObjectReference{}

	for _, secret := range image.PullSecrets {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{
			Name: secret,
		})
	}

	return imagePullSecrets
}

// ContainerEnvFrom prepares EnvFrom resource for Agent and Controllers' container.
func ContainerEnvFrom(controllerSpec common.CommonSpec) []corev1.EnvFromSource {
	envFrom := []corev1.EnvFromSource{}
	if controllerSpec.ExtraEnvVarsCM != "" {
		envFrom = append(envFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: controllerSpec.ExtraEnvVarsCM,
				},
			},
		})
	}

	if controllerSpec.ExtraEnvVarsSecret != "" {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: controllerSpec.ExtraEnvVarsSecret,
				},
			},
		})
	}

	return envFrom
}

// ContainerProbes prepares livenessProbe and readinessProbe based on the provided parameters.
func ContainerProbes(spec common.CommonSpec, scheme corev1.URIScheme) (*corev1.Probe, *corev1.Probe) {
	var livenessProbe *corev1.Probe
	var readinessProbe *corev1.Probe
	if spec.LivenessProbe.Enabled {
		livenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/v1/status/system/liveness",
					Port:   intstr.FromString(Server),
					Scheme: scheme,
				},
			},
			InitialDelaySeconds: spec.LivenessProbe.InitialDelaySeconds,
			TimeoutSeconds:      spec.LivenessProbe.TimeoutSeconds,
			PeriodSeconds:       spec.LivenessProbe.PeriodSeconds,
			FailureThreshold:    spec.LivenessProbe.FailureThreshold,
			SuccessThreshold:    spec.LivenessProbe.SuccessThreshold,
		}
	} else if spec.CustomLivenessProbe != nil {
		livenessProbe = spec.CustomLivenessProbe
	}

	if spec.ReadinessProbe.Enabled {
		readinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/v1/status/system/readiness",
					Port:   intstr.FromString(Server),
					Scheme: scheme,
				},
			},
			InitialDelaySeconds: spec.ReadinessProbe.InitialDelaySeconds,
			TimeoutSeconds:      spec.ReadinessProbe.TimeoutSeconds,
			PeriodSeconds:       spec.ReadinessProbe.PeriodSeconds,
			FailureThreshold:    spec.ReadinessProbe.FailureThreshold,
			SuccessThreshold:    spec.ReadinessProbe.SuccessThreshold,
		}
	} else if spec.CustomReadinessProbe != nil {
		readinessProbe = spec.CustomReadinessProbe
	}

	return livenessProbe, readinessProbe
}

// AgentEnv prepares env resources for Agents' container.
func AgentEnv(instance *agentv1alpha1.Agent, agentGroup string) []corev1.EnvVar {
	spec := instance.Spec

	envs := []corev1.EnvVar{
		{
			Name: "NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: V1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		},
	}

	if agentGroup == "" && spec.ConfigSpec.AgentInfo.AgentGroup != "" {
		agentGroup = spec.ConfigSpec.AgentInfo.AgentGroup
	}

	if agentGroup != "" {
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
			Value: agentGroup,
		})
	}

	if spec.Sidecar.Enabled {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: V1Version,
					FieldPath:  "metadata.name",
				},
			},
		})
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_ENABLED",
			Value: "false",
		})
	} else {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: V1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		})
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_ENABLED",
			Value: "true",
		})
	}

	if instance.Spec.Secrets.FluxNinjaExtension.Create || instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef.Name != "" {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_FLUXNINJA_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: SecretName(instance.GetName(), "agent", &instance.Spec.Secrets.FluxNinjaExtension),
					},
					Key:      SecretDataKey(&instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef),
					Optional: pointer.Bool(false),
				},
			},
		})
	}

	return MergeEnvVars(envs, spec.ExtraEnvVars)
}

// AgentVolumeMounts prepares volumeMounts for Agents' container.
func AgentVolumeMounts(agentSpec agentv1alpha1.AgentSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "aperture-agent-config",
			MountPath: "/etc/aperture/aperture-agent/config",
		},
	}

	if len(agentSpec.ConfigSpec.AgentFunctions.Endpoints) > 0 &&
		agentSpec.ControllerClientCertConfig.ConfigMapName != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      agentSpec.ControllerClientCertConfig.ConfigMapName,
			MountPath: AgentControllerClientCertPath,
		})
	}

	return MergeVolumeMounts(volumeMounts, agentSpec.ExtraVolumeMounts)
}

// AgentVolumes prepares volumes for Agent.
func AgentVolumes(agentSpec agentv1alpha1.AgentSpec) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "aperture-agent-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: AgentServiceName,
					},
				},
			},
		},
	}

	if agentSpec.ControllerClientCertConfig.ConfigMapName != "" {
		volumes = append(volumes, corev1.Volume{
			Name: agentSpec.ControllerClientCertConfig.ConfigMapName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: agentSpec.ControllerClientCertConfig.ConfigMapName,
					},
				},
			},
		})
	}

	return MergeVolumes(volumes, agentSpec.ExtraVolumes)
}

// ControllerEnv prepares env resources for Controller' container.
func ControllerEnv(instance *controllerv1alpha1.Controller) []corev1.EnvVar {
	spec := instance.Spec

	envs := []corev1.EnvVar{
		{
			Name: "APERTURE_CONTROLLER_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: V1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		},
		{
			Name:  "APERTURE_CONTROLLER_NAMESPACE",
			Value: instance.GetNamespace(),
		},
	}

	if spec.Secrets.FluxNinjaExtension.Create || instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef.Name != "" {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_CONTROLLER_FLUXNINJA_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: SecretName(instance.GetName(), "controller", &instance.Spec.Secrets.FluxNinjaExtension),
					},
					Key:      SecretDataKey(&instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef),
					Optional: pointer.Bool(false),
				},
			},
		})
	}

	return MergeEnvVars(envs, spec.ExtraEnvVars)
}

// ControllerVolumeMounts prepares volumeMounts for Controllers' container.
func ControllerVolumeMounts(controllerSpec common.CommonSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "aperture-controller-config",
			MountPath: "/etc/aperture/aperture-controller/config",
		},
		{
			Name:      "server-cert",
			MountPath: "/etc/aperture/aperture-controller/certs",
			ReadOnly:  true,
		},
	}

	return MergeVolumeMounts(volumeMounts, controllerSpec.ExtraVolumeMounts)
}

// ControllerVolumes prepares volumes for Controller.
func ControllerVolumes(instance *controllerv1alpha1.Controller) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "aperture-controller-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: ControllerServiceName,
					},
				},
			},
		},
		{
			Name: "server-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: pointer.Int32(420),
					SecretName:  fmt.Sprintf("%s-controller-cert", instance.GetName()),
				},
			},
		},
	}

	return MergeVolumes(volumes, instance.Spec.ExtraVolumes)
}

// CommonLabels prepares common labels used by all resources.
func CommonLabels(commonLabels map[string]string, instanceName, component string) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       AppName,
		"app.kubernetes.io/instance":   instanceName,
		"app.kubernetes.io/managed-by": OperatorName,
		"app.kubernetes.io/component":  component,
	}

	if commonLabels != nil {
		if err := mergo.Map(&labels, commonLabels, mergo.WithOverride); err != nil {
			return labels
		}
	}

	return labels
}

// SelectorLabels prepares the labels used for Selector.
func SelectorLabels(instance, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       AppName,
		"app.kubernetes.io/instance":   instance,
		"app.kubernetes.io/managed-by": OperatorName,
		"app.kubernetes.io/component":  component,
	}
}

// ControllerAnnotationsWithOwnerRef prepares the map for Annotation with reference to the creator instance.
func ControllerAnnotationsWithOwnerRef(instance *controllerv1alpha1.Controller) map[string]string {
	annotations := instance.Spec.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["fluxninja.com/primary-resource"] = fmt.Sprintf("%s/%s", instance.GetNamespace(), instance.GetName())
	annotations["fluxninja.com/primary-resource-type"] = fmt.Sprintf("%s.%s",
		instance.GetObjectKind().GroupVersionKind().GroupKind().Kind, instance.GetObjectKind().GroupVersionKind().GroupKind().Group)

	return annotations
}

// AgentAnnotationsWithOwnerRef prepares the map for Annotation with reference to the creator instance.
func AgentAnnotationsWithOwnerRef(instance *agentv1alpha1.Agent) map[string]string {
	annotations := instance.Spec.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["fluxninja.com/primary-resource"] = fmt.Sprintf("%s/%s", instance.GetNamespace(), instance.GetName())
	annotations["fluxninja.com/primary-resource-type"] = fmt.Sprintf("%s.%s",
		instance.GetObjectKind().GroupVersionKind().GroupKind().Kind, instance.GetObjectKind().GroupVersionKind().GroupKind().Group)

	return annotations
}

// SecretName fetches name for ApiKey secret from config or generates the name if not present in config.
func SecretName(instance, component string, spec *common.APIKeySecret) string {
	name := spec.SecretKeyRef.Name
	if name != "" {
		return name
	}

	return fmt.Sprintf("%s-%s-apikey", instance, component)
}

// SecretDataKey fetches Key for ApiKey secret from config or generates the Key if not present in config.
func SecretDataKey(spec *common.SecretKeyRef) string {
	key := spec.Key
	if key != "" {
		return key
	}

	return SecretKey
}

// CheckCertificate checks if existing certificates are available.
func CheckCertificate() bool {
	certDir := strings.TrimRight(os.Getenv("APERTURE_OPERATOR_CERT_DIR"), "/")
	if certDir == "" {
		certDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
		os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
	}

	certName := os.Getenv("APERTURE_OPERATOR_CERT_NAME")
	if certName == "" {
		certName = "tls.crt"
		os.Setenv("APERTURE_OPERATOR_CERT_NAME", certName)
	}

	keyName := os.Getenv("APERTURE_OPERATOR_KEY_NAME")
	if keyName == "" {
		keyName = "tls.key"
		os.Setenv("APERTURE_OPERATOR_KEY_NAME", keyName)
	}

	_, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME")),
		fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_KEY_NAME")))

	return err == nil
}

// GenerateCertificate generates certificate and stores it in the desired location.
func GenerateCertificate(dnsPrefix, namespace string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	var caPEM, serverCertPEM, serverPrivKeyPEM *bytes.Buffer

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization: []string{"fluxninja.com"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// CA private key
	caPrivKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// Self signed CA certificate
	caBytes, err := x509.CreateCertificate(cryptorand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// PEM encode CA cert
	caPEM = new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	dnsNames := []string{
		dnsPrefix,
		fmt.Sprintf("%s.%s", dnsPrefix, namespace),
		fmt.Sprintf("%s.%s.svc", dnsPrefix, namespace),
		fmt.Sprintf("%s.%s.svc.cluster.local", dnsPrefix, namespace),
	}

	commonName := fmt.Sprintf("%s.%s.svc", dnsPrefix, namespace)

	// server cert config
	cert := &x509.Certificate{
		DNSNames:     dnsNames,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"fluxninja.com"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// server private key
	serverPrivKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// sign the server cert
	serverCertBytes, err := x509.CreateCertificate(cryptorand.Reader, cert, ca, &serverPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// PEM encode the  server cert and key
	serverCertPEM = new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})

	serverPrivKeyPEM = new(bytes.Buffer)
	_ = pem.Encode(serverPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey),
	})

	return serverCertPEM, serverPrivKeyPEM, caPEM, nil
}

// WriteFile writes data in the file at the given path.
func WriteFile(filepath string, sCert *bytes.Buffer) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(sCert.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// CheckAndGenerateCertForOperator checks if existing certificates are present and creates new if not present.
func CheckAndGenerateCertForOperator(config *rest.Config) error {
	if CheckCertificate() {
		return nil
	}

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}

	namespace := os.Getenv("APERTURE_OPERATOR_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	serviceName := os.Getenv("APERTURE_OPERATOR_SERVICE_NAME")
	if serviceName == "" {
		return fmt.Errorf("the value for environment variable 'APERTURE_OPERATOR_SERVICE_NAME' is not configured")
	}

	serverCertPEM, serverPrivKeyPEM, caPEM, err := GenerateCertificate(serviceName, namespace)
	if err != nil {
		return err
	}

	var updateSecret bool
	secretName := fmt.Sprintf("%s-cert", serviceName)
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels:    CommonLabels(map[string]string{}, serviceName, AppName),
		},
	}

	err = k8sClient.Get(context.Background(), types.NamespacedName{Name: secretName, Namespace: namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			updateSecret = true
		} else {
			return err
		}
	} else {
		certBytes, ok := secret.Data[OperatorCertName]
		if !ok {
			updateSecret = true
		} else {
			block, _ := pem.Decode(certBytes)
			var cert *x509.Certificate
			cert, err = x509.ParseCertificate(block.Bytes)
			if err == nil && cert.NotAfter.After(time.Now()) {
				serverCertPEM = bytes.NewBuffer(secret.Data[OperatorCertName])
				serverPrivKeyPEM = bytes.NewBuffer(secret.Data[OperatorCertKeyName])
				caPEM = bytes.NewBuffer(secret.Data[OperatorCAName])
			} else {
				updateSecret = true
			}
		}
	}

	if updateSecret {
		secret.Data = map[string][]byte{
			OperatorCertName:    serverCertPEM.Bytes(),
			OperatorCertKeyName: serverPrivKeyPEM.Bytes(),
			OperatorCAName:      caPEM.Bytes(),
		}

		_, err = controllerutil.CreateOrUpdate(context.Background(), k8sClient, secret, SecretMutate(secret, secret.Data, secret.OwnerReferences))
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(os.Getenv("APERTURE_OPERATOR_CERT_DIR"), 0o777)
	if err != nil {
		return err
	}

	err = WriteFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME")), serverCertPEM)
	if err != nil {
		return err
	}

	err = WriteFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_KEY_NAME")), serverPrivKeyPEM)
	if err != nil {
		return err
	}

	err = WriteFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName), caPEM)
	if err != nil {
		return err
	}

	return nil
}

// GetOrGenerateCertificate returns the TLS/SSL certificates of the Controller.
func GetOrGenerateCertificate(client client.Client, instance *controllerv1alpha1.Controller) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	secretName := fmt.Sprintf("%s-controller-cert", instance.GetName())

	generateCert := func() (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
		// generate certificates
		serverCertPEM, serverPrivKeyPEM, caPEM, err := GenerateCertificate(ControllerServiceName, instance.GetNamespace())
		if err != nil {
			return nil, nil, nil, err
		}

		return serverCertPEM, serverPrivKeyPEM, caPEM, nil
	}

	secret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: instance.GetNamespace()}, secret)
	if err != nil {
		return generateCert()
	}

	existingCert, ok := secret.Data[ControllerCertName]
	if !ok {
		return generateCert()
	}

	existingKey, ok := secret.Data[ControllerCertKeyName]
	if !ok {
		return generateCert()
	}

	block, _ := pem.Decode(existingCert)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return generateCert()
	}

	// regenerate certificate if it is expired
	if time.Now().After(cert.NotAfter) {
		return generateCert()
	}

	var existingClientCert []byte
	cm := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-controller-client-cert", instance.GetName()), Namespace: instance.GetNamespace()}, cm)
	if err != nil {
		// Checking existing ValidatingWebhookConfiguration for backward compatibility
		vwc := &admissionregistrationv1.ValidatingWebhookConfiguration{}
		err := client.Get(context.TODO(), types.NamespacedName{Name: ControllerServiceName}, vwc)
		if err != nil || len(vwc.Webhooks) == 0 {
			// No ValidatingWebhookConfiguration found, generate new certificate
			return generateCert()
		}

		existingClientCert = vwc.Webhooks[0].ClientConfig.CABundle
	} else {
		existingClientCert = []byte(cm.Data[ControllerClientCertKey])
	}

	return bytes.NewBuffer(existingCert), bytes.NewBuffer(existingKey), bytes.NewBuffer(existingClientCert), nil
}

// MergeEnvVars merges common and provided extra Environment variables of Kubernetes container.
func MergeEnvVars(common, extra []corev1.EnvVar) []corev1.EnvVar {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		keys[env.Name] = true
	}

	for _, env := range common {
		if _, ok := keys[env.Name]; !ok {
			extra = append(extra, env)
			keys[env.Name] = true
		}
	}

	return extra
}

// MergeEnvFromSources merges common and provided extra Environment From of Kubernetes container.
func MergeEnvFromSources(common, extra []corev1.EnvFromSource) []corev1.EnvFromSource {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		if env.ConfigMapRef != nil {
			keys[fmt.Sprintf("configmap_%s", env.ConfigMapRef.LocalObjectReference.Name)] = true
		}

		if env.SecretRef != nil {
			keys[fmt.Sprintf("secret_%s", env.SecretRef.LocalObjectReference.Name)] = true
		}
	}

	for _, env := range common {
		envFrom := corev1.EnvFromSource{}
		if env.ConfigMapRef != nil {
			if _, ok := keys[fmt.Sprintf("configmap_%s", env.ConfigMapRef.LocalObjectReference.Name)]; !ok {
				envFrom.ConfigMapRef = env.ConfigMapRef
				keys[fmt.Sprintf("configmap_%s", env.ConfigMapRef.LocalObjectReference.Name)] = true
			}
		}

		if env.SecretRef != nil {
			if _, ok := keys[fmt.Sprintf("secret_%s", env.SecretRef.LocalObjectReference.Name)]; !ok {
				envFrom.SecretRef = env.SecretRef
				keys[fmt.Sprintf("secret_%s", env.SecretRef.LocalObjectReference.Name)] = true
			}
		}

		if envFrom.ConfigMapRef != nil || envFrom.SecretRef != nil {
			extra = append(extra, envFrom)
		}
	}
	return extra
}

// MergeVolumeMounts merges common and provided extra Volume mounts of Kubernetes container.
func MergeVolumeMounts(common, extra []corev1.VolumeMount) []corev1.VolumeMount {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		keys[env.Name] = true
	}

	for _, env := range common {
		if _, ok := keys[env.Name]; !ok {
			extra = append(extra, env)
			keys[env.Name] = true
		}
	}
	return extra
}

// MergeVolumes merges common and provided extra Volume of Kubernetes Pod.
func MergeVolumes(common, extra []corev1.Volume) []corev1.Volume {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		keys[env.Name] = true
	}

	for _, env := range common {
		if _, ok := keys[env.Name]; !ok {
			extra = append(extra, env)
			keys[env.Name] = true
		}
	}
	return extra
}

// MergeContainers merges common and provided Container/Init Container of Kubernetes container.
func MergeContainers(common, extra []corev1.Container) []corev1.Container {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		keys[env.Name] = true
	}

	for _, env := range common {
		if _, ok := keys[env.Name]; !ok {
			extra = append(extra, env)
			keys[env.Name] = true
		}
	}
	return extra
}

// MergeImagePullSecrets merges common and provided Image Pull Secrets of Kubernetes.
func MergeImagePullSecrets(common, extra []corev1.LocalObjectReference) []corev1.LocalObjectReference {
	if extra == nil {
		return common
	}

	keys := map[string]bool{}
	for _, env := range extra {
		keys[env.Name] = true
	}

	for _, env := range common {
		if _, ok := keys[env.Name]; !ok {
			extra = append(extra, env)
			keys[env.Name] = true
		}
	}
	return extra
}

// UpdateResource updates the Aperture resource in Kubernetes.
func UpdateResource(client client.Client, ctx context.Context, instance client.Object) error {
	attempt := 5
	for attempt > 0 {
		attempt -= 1
		if err := client.Update(ctx, instance); err != nil {
			if errors.IsConflict(err) {
				namespacesName := types.NamespacedName{
					Namespace: instance.GetNamespace(),
					Name:      instance.GetName(),
				}
				if err = client.Get(ctx, namespacesName, instance); err != nil {
					return err
				}
				continue
			}
			return err
		}
	}

	return nil
}

// GetPort parses port value from the Address string.
func GetPort(addr string) (int32, error) {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	// read 32 bit integer from string
	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(port), nil
}

// GetControllerClientCert returns the controller client certificate from the controller configmap.
func GetControllerClientCert(endpoints []string, client_ client.Client, ctx context.Context) []byte {
	var localControllerCert []byte
	for _, endpoint := range endpoints {
		controllerNS := localControllerNamespaceFromEndpoint(endpoint)
		if controllerNS == "" {
			continue
		}

		var configMaps corev1.ConfigMapList
		err := client_.List(ctx, &configMaps, &client.ListOptions{Namespace: controllerNS})
		if err != nil {
			continue
		}

		for _, cm := range configMaps.Items {
			if !strings.HasSuffix(cm.Name, "-controller-client-cert") {
				continue
			}
			localControllerCert = append(localControllerCert, cm.Data[ControllerClientCertKey]...)
		}
	}

	return localControllerCert
}

// localControllerNamespaceFromEndpoint returns the namespace of the local controller.
func localControllerNamespaceFromEndpoint(endpoint string) string {
	addr, port, ok := strings.Cut(endpoint, ":")
	if !ok {
		return ""
	}

	if port != "8080" {
		return ""
	}

	subdomains := strings.Split(addr, ".")
	if len(subdomains) < 2 {
		return ""
	}

	if !strings.Contains(subdomains[0], "aperture-controller") {
		return ""
	}

	tail := strings.Join(subdomains[2:], ".") + "."
	if !strings.HasPrefix("svc.cluster.local.", tail) { //nolint:gocritic
		return ""
	}

	return subdomains[1]
}
