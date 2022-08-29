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
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	cryptorand "crypto/rand"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// containerSecurityContext prepares SecurityContext for containers based on the provided parameter.
func containerSecurityContext(containerSecurityContext v1alpha1.ContainerSecurityContext) *corev1.SecurityContext {
	var securityContext *corev1.SecurityContext
	if containerSecurityContext.Enabled {
		securityContext = &corev1.SecurityContext{
			RunAsUser:              containerSecurityContext.RunAsUser,
			RunAsNonRoot:           containerSecurityContext.RunAsNonRootUser,
			ReadOnlyRootFilesystem: containerSecurityContext.ReadOnlyRootFilesystem,
		}
	} else {
		securityContext = &corev1.SecurityContext{}
	}

	return securityContext
}

// getContainerSecurityContext prepares SecurityContext for containers based on the provided parameter.
func podSecurityContext(podSecurityContext v1alpha1.PodSecurityContext) *corev1.PodSecurityContext {
	var securityContext *corev1.PodSecurityContext
	if podSecurityContext.Enabled {
		securityContext = &corev1.PodSecurityContext{
			FSGroup: podSecurityContext.FsGroup,
		}
	} else {
		securityContext = &corev1.PodSecurityContext{}
	}

	return securityContext
}

// imageString prepares image string from the provided Image struct.
func imageString(globalRegistry string, image v1alpha1.Image) string {
	var imageStr string
	if globalRegistry != "" {
		imageStr = fmt.Sprintf("%s/%s:%s", globalRegistry, image.Repository, image.Tag)
	} else if image.Registry != "" {
		imageStr = fmt.Sprintf("%s/%s:%s", image.Registry, image.Repository, image.Tag)
	} else {
		imageStr = fmt.Sprintf("%s:%s", image.Repository, image.Tag)
	}
	return imageStr
}

// imagePullSecrets prepares imagePullSecrets string slice from the provided Image struct.
func imagePullSecrets(globalPullSecrets []string, image v1alpha1.Image) []corev1.LocalObjectReference {
	imagePullSecrets := []corev1.LocalObjectReference{}
	globalImagePullSecrets := []corev1.LocalObjectReference{}
	for _, secret := range globalPullSecrets {
		globalImagePullSecrets = append(globalImagePullSecrets, corev1.LocalObjectReference{
			Name: secret,
		})
	}

	for _, secret := range image.PullSecrets {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{
			Name: secret,
		})
	}

	return mergeImagePullSecrets(globalImagePullSecrets, imagePullSecrets)
}

// containerEnvFrom prepares EnvFrom resource for Agent and Controllers' container.
func containerEnvFrom(controllerSpec v1alpha1.CommonSpec) []corev1.EnvFromSource {
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

// containerProbes prepares livenessProbe and readinessProbe based on the provided parameters.
func containerProbes(spec v1alpha1.CommonSpec) (*corev1.Probe, *corev1.Probe) {
	var livenessProbe *corev1.Probe
	var readinessProbe *corev1.Probe
	if spec.LivenessProbe.Enabled {
		livenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/v1/status/liveness",
					Port:   intstr.FromString("grpc"),
					Scheme: corev1.URISchemeHTTP,
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
					Path:   "/v1/status/readiness",
					Port:   intstr.FromString("grpc"),
					Scheme: corev1.URISchemeHTTP,
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

// agentEnv prepares env resources for Agents' container.
func agentEnv(instance *v1alpha1.Aperture, agentGroup string) []corev1.EnvVar {
	spec := instance.Spec.Agent

	envs := []corev1.EnvVar{
		{
			Name: "NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: v1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		},
	}

	if agentGroup == "" && spec.AgentGroup != "" {
		agentGroup = spec.AgentGroup
	}

	if agentGroup != "" {
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
			Value: agentGroup,
		})
	}

	if instance.Spec.Sidecar.Enabled {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: v1Version,
					FieldPath:  "metadata.name",
				},
			},
		})
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
			Value: "false",
		})
	} else {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: v1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		})
		envs = append(envs, corev1.EnvVar{
			Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
			Value: "true",
		})
	}

	if instance.Spec.FluxNinjaPlugin.Enabled {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_AGENT_FLUXNINJA_PLUGIN_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretName(instance.GetName(), "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret.Agent),
					},
					Key:      secretDataKey(&instance.Spec.FluxNinjaPlugin.APIKeySecret.Agent.SecretKeyRef),
					Optional: pointer.BoolPtr(false),
				},
			},
		})
	}

	return mergeEnvVars(envs, spec.ExtraEnvVars)
}

// agentVolumeMounts prepares volumeMounts for Agents' container.
func agentVolumeMounts(agentSpec v1alpha1.AgentSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "aperture-agent-config",
			MountPath: "/etc/aperture/aperture-agent/config",
		},
	}

	return mergeVolumeMounts(volumeMounts, agentSpec.ExtraVolumeMounts)
}

// agentVolumes prepares volumes for Agent.
func agentVolumes(agentSpec v1alpha1.AgentSpec) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "aperture-agent-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32Ptr(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: agentServiceName,
					},
				},
			},
		},
	}

	return mergeVolumes(volumes, agentSpec.ExtraVolumes)
}

// controllerEnv prepares env resources for Controller' container.
func controllerEnv(instance *v1alpha1.Aperture) []corev1.EnvVar {
	spec := instance.Spec.Controller

	envs := []corev1.EnvVar{
		{
			Name: "APERTURE_CONTROLLER_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: v1Version,
					FieldPath:  "spec.nodeName",
				},
			},
		},
	}

	if instance.Spec.FluxNinjaPlugin.Enabled {
		envs = append(envs, corev1.EnvVar{
			Name: "APERTURE_CONTROLLER_FLUXNINJA_PLUGIN_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretName(instance.GetName(), "controller", &instance.Spec.FluxNinjaPlugin.APIKeySecret.Controller),
					},
					Key:      secretDataKey(&instance.Spec.FluxNinjaPlugin.APIKeySecret.Controller.SecretKeyRef),
					Optional: pointer.BoolPtr(false),
				},
			},
		})
	}

	return mergeEnvVars(envs, spec.ExtraEnvVars)
}

// getVolumes prepares volumeMounts for Controllers' container.
func controllerVolumeMounts(controllerSpec v1alpha1.CommonSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "aperture-controller-config",
			MountPath: "/etc/aperture/aperture-controller/config",
		},
		{
			Name:      "etc-aperture-policies",
			MountPath: "/etc/aperture/aperture-controller/policies",
			ReadOnly:  true,
		},
		{
			Name:      "etc-aperture-classification",
			MountPath: "/etc/aperture/aperture-controller/classifiers",
			ReadOnly:  true,
		},
		{
			Name:      "webhook-cert",
			MountPath: "/etc/aperture/aperture-controller/certs",
			ReadOnly:  true,
		},
	}

	return mergeVolumeMounts(volumeMounts, controllerSpec.ExtraVolumeMounts)
}

// controllerVolumes prepares volumes for Controller.
func controllerVolumes(instance *v1alpha1.Aperture) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "aperture-controller-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32Ptr(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: controllerServiceName,
					},
				},
			},
		},
		{
			Name: "etc-aperture-policies",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32Ptr(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "policies",
					},
					Optional: pointer.BoolPtr(true),
				},
			},
		},
		{
			Name: "etc-aperture-classification",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: pointer.Int32Ptr(420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "classification",
					},
					Optional: pointer.BoolPtr(true),
				},
			},
		},
		{
			Name: "webhook-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: pointer.Int32Ptr(420),
					SecretName:  fmt.Sprintf("%s-controller-cert", instance.GetName()),
				},
			},
		},
	}

	return mergeVolumes(volumes, instance.Spec.Controller.ExtraVolumes)
}

// commonLabels prepares common labels used by all resources.
func commonLabels(instance *v1alpha1.Aperture, component string) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       appName,
		"app.kubernetes.io/instance":   instance.GetName(),
		"app.kubernetes.io/managed-by": operatorName,
		"app.kubernetes.io/component":  component,
	}

	if instance.Spec.Labels != nil {
		if err := mergo.Map(&labels, instance.Spec.Labels, mergo.WithOverride); err != nil {
			return labels
		}
	}

	return labels
}

// selectorLabels prepares the labels used for Selector.
func selectorLabels(instance, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       appName,
		"app.kubernetes.io/instance":   instance,
		"app.kubernetes.io/managed-by": operatorName,
		"app.kubernetes.io/component":  component,
	}
}

// getAnnotationsWithOwnerRef prepares the map for Annotation with reference to the creator instance.
func getAnnotationsWithOwnerRef(instance *v1alpha1.Aperture) map[string]string {
	annotations := instance.Spec.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["fluxninja.com/primary-resource"] = fmt.Sprintf("%s/%s", instance.GetNamespace(), instance.GetName())
	annotations["fluxninja.com/primary-resource-type"] = fmt.Sprintf("%s.%s",
		instance.GetObjectKind().GroupVersionKind().GroupKind().Kind, instance.GetObjectKind().GroupVersionKind().GroupKind().Group)

	return annotations
}

// checkEtcdEndpoints generates endpoints list based on the release name if that is not provided else returns the provided values.
func checkEtcdEndpoints(etcd v1alpha1.EtcdSpec, name, namespace string) v1alpha1.EtcdSpec {
	endpoints := []string{}
	if etcd.Endpoints != nil {
		for _, endpoint := range etcd.Endpoints {
			if endpoint != "" {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	if len(endpoints) == 0 {
		endpoints = append(endpoints, fmt.Sprintf("http://%s-etcd.%s:2379", name, namespace))
	}

	etcd.Endpoints = endpoints
	return etcd
}

// checkPrometheusAddress generates prometheus address based on the release name if that is not provided else returns the provided value.
func checkPrometheusAddress(address, name, namespace string) string {
	if address == "" {
		address = fmt.Sprintf("http://%s-prometheus-server.%s:80", name, namespace)
	}
	return strings.TrimRight(address, "/")
}

// secretName fetches name for ApiKey secret from config or generates the name if not present in config.
func secretName(instance, component string, spec *v1alpha1.APIKeySecret) string {
	name := spec.SecretKeyRef.Name
	if name != "" {
		return name
	}

	return fmt.Sprintf("%s-%s-apikey", instance, component)
}

// secretDataKey fetches Key for ApiKey secret from config or generates the Key if not present in config.
func secretDataKey(spec *v1alpha1.SecretKeyRef) string {
	key := spec.Key
	if key != "" {
		return key
	}

	return secretKey
}

// checkCertificate checks if existing certificates are available.
func checkCertificate() bool {
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

// generateCertificate generates certificate and stores it in the desired location.
func generateCertificate(dnsPrefix, namespace string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
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

// writeFile writes data in the file at the given path.
func writeFile(filepath string, sCert *bytes.Buffer) error {
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
func CheckAndGenerateCertForOperator() error {
	if checkCertificate() {
		return nil
	}

	namespace := os.Getenv("APERTURE_OPERATOR_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	serviceName := os.Getenv("APERTURE_OPERATOR_SERVICE_NAME")
	if serviceName == "" {
		return fmt.Errorf("the value for environment variable 'APERTURE_OPERATOR_SERVICE_NAME' is not configured")
	}

	serverCertPEM, serverPrivKeyPEM, caPEM, err := generateCertificate(serviceName, namespace)
	if err != nil {
		return err
	}

	err = os.MkdirAll(os.Getenv("APERTURE_OPERATOR_CERT_DIR"), 0o777)
	if err != nil {
		return err
	}

	err = writeFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME")), serverCertPEM)
	if err != nil {
		return err
	}

	err = writeFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_KEY_NAME")), serverPrivKeyPEM)
	if err != nil {
		return err
	}

	err = writeFile(fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName), caPEM)
	if err != nil {
		return err
	}

	return nil
}

// mergeEnvVars merges common and provided extra Environment variables of Kubernetes container.
func mergeEnvVars(common, extra []corev1.EnvVar) []corev1.EnvVar {
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

// mergeEnvFromSources merges common and provided extra Environment From of Kubernetes container.
func mergeEnvFromSources(common, extra []corev1.EnvFromSource) []corev1.EnvFromSource {
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

// mergeVolumeMounts merges common and provided extra Volume mounts of Kubernetes container.
func mergeVolumeMounts(common, extra []corev1.VolumeMount) []corev1.VolumeMount {
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

// mergeVolumes merges common and provided extra Volume of Kubernetes Pod.
func mergeVolumes(common, extra []corev1.Volume) []corev1.Volume {
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

// mergeContainers merges common and provided Container/Init Container of Kubernetes container.
func mergeContainers(common, extra []corev1.Container) []corev1.Container {
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

// mergeImagePullSecrets merges common and provided Image Pull Secrets of Kubernetes.
func mergeImagePullSecrets(common, extra []corev1.LocalObjectReference) []corev1.LocalObjectReference {
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

// updateAperture updates the Aperture resource in Kubernetes.
func updateResource(client client.Client, ctx context.Context, instance client.Object) error {
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
