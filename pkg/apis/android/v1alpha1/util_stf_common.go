package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// STFConfig returns the STF configuration for this AndroidFarm instance.
func (a *AndroidFarm) STFConfig() *STFConfig {
	return a.Spec.STFConfig
}

// STFDisabled returns true if there is no STF configuration provided, which
// means we skip deploying those resources.
func (a *AndroidFarm) STFDisabled() bool {
	return a.Spec.STFConfig == nil
}

// OwnerReferences returns an owner reference slice with this AndroidFarm
// instance as the owner.
func (a *AndroidFarm) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion:         a.APIVersion,
			Kind:               a.Kind,
			Name:               a.GetName(),
			UID:                a.GetUID(),
			Controller:         &trueVal,
			BlockOwnerDeletion: &trueVal,
		},
	}
}

// STFNamePrefix returns the string to prefix STF resource names with.
func (a *AndroidFarm) STFNamePrefix() string {
	return fmt.Sprintf("%s-stf", a.GetName())
}

// STFComponentLabels returns the labels to apply to STF resources for this
// AndroidFarm.
func (a *AndroidFarm) STFComponentLabels(component string) map[string]string {
	return map[string]string{
		DeviceFarmLabel: a.GetName(),
		"app":           "stf",
		"component":     component,
	}
}

func (a *AndroidFarm) STFComponentName(name string) string {
	return fmt.Sprintf("%s-%s", a.STFNamePrefix(), name)
}

// GetNamespace returns the namespace where STF resources should be deployed for
// this AndroidFarm.
func (s *STFConfig) GetNamespace() string {
	if s.Namespace != "" {
		return s.Namespace
	}
	return "default"
}

// GetServiceAccount returns the service account (if any) we should use for STF
// deployments in this farm.
func (s *STFConfig) GetServiceAccount() string {
	return s.ServiceAccount
}

// OpenSTFImage returns the OpenSTF docker image to use for deployments in this
// farm.
func (s *STFConfig) OpenSTFImage() string {
	if s.STFImage != nil {
		if s.STFImage.Image != "" {
			return s.STFImage.Image
		}
	}
	return defaultSTFImage
}

// PullPolicy returns the image pull policy for STF resources.
func (s *STFConfig) PullPolicy() corev1.PullPolicy {
	if s.STFImage != nil {
		if s.STFImage.ImagePullPolicy != "" {
			return s.STFImage.ImagePullPolicy
		}
	}
	return corev1.PullIfNotPresent
}

// PullSecrets returns the pull secrets (if any) to be used for pulling STF
// images.
func (s *STFConfig) PullSecrets() []corev1.LocalObjectReference {
	secrets := make([]corev1.LocalObjectReference, 0)
	if s.STFImage != nil {
		if s.STFImage.ImagePullSecrets != nil {
			secrets = append(secrets, s.STFImage.ImagePullSecrets...)
		}
	}
	if s.ADB != nil {
		if s.ADB.ImagePullSecrets != nil {
			secrets = append(secrets, s.ADB.ImagePullSecrets...)
		}
	}
	return secrets
}

// GetAppHostname returns the external hostname (or IP) that is used when
// configuring STF deployments.
func (s *STFConfig) GetAppHostname() string {
	return s.AppHostname
}

func (s *STFConfig) GetAppExternalURL() string {
	return fmt.Sprintf("%s://%s", s.GetHTTPScheme(), s.GetAppHostname())
}

func (s *STFConfig) GetAppExternalWebsocketURL() string {
	return fmt.Sprintf("%s://%s", s.GetWebsocketScheme(), s.GetAppHostname())
}

// PodSecurityContext returns the pod security context to use for STF deployments
// in this farm.
func (s *STFConfig) PodSecurityContext() *corev1.PodSecurityContext {
	if s.PrivilegedDeployments {
		return &corev1.PodSecurityContext{
			RunAsNonRoot: &falseVal,
		}
	}
	return nil
}

// ContainerSecurityContext returns the container security context to use for
// STF deployments in this farm.
func (s *STFConfig) ContainerSecurityContext() *corev1.SecurityContext {
	if s.PrivilegedDeployments {
		return &corev1.SecurityContext{
			Privileged: &trueVal,
		}
	}
	return nil
}

// GetSTFSecretKey returns the key in our STFSecret where we can find the
// secret to give to STF deployments.
func (s *STFConfig) GetSTFSecretKey() string {
	if s.STFSecretKey != "" {
		return s.STFSecretKey
	}
	return defaultSTFSecretKey
}
