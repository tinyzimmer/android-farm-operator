package stf

import (
	"bytes"
	"fmt"
	"strconv"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"

	// traefikv1 "github.com/containous/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"github.com/go-logr/logr"
	cm "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *STFReconciler) reconcileIngress(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	if instance.STFConfig().GetIssuerReference() != nil {
		cert := newCertificateForCR(instance)
		if err := util.ReconcileCertificate(reqLogger, r.client, cert, true); err != nil {
			return err
		}
	}
	if instance.STFConfig().UseIngressRouteCRD() {
		// WIP
		return nil
	}
	return reconcileIngressDeployment(reqLogger, r.client, instance)
}

func reconcileIngressDeployment(reqLogger logr.Logger, c client.Client, instance *androidv1alpha1.AndroidFarm) error {
	var static bytes.Buffer
	if err := traefikStaticConfigTmpl.Execute(&static, map[string]interface{}{
		"UseSSL":    instance.STFConfig().TLSEnabled() && !instance.STFConfig().UseExternalTLS(),
		"Services":  calculateServiceDefinitions(instance),
		"AccessLog": instance.STFConfig().TraefikAccessLogsEnabled(),
		"Backtick":  "`",
	}); err != nil {
		return err
	}
	var dynamic bytes.Buffer
	if err := traefikDynamicConfigTmpl.Execute(&dynamic, map[string]interface{}{
		"UseSSL":             instance.STFConfig().TLSEnabled() && !instance.STFConfig().UseExternalTLS(),
		"UseSelfSigned":      instance.TLSSecret() == "",
		"Services":           calculateServiceDefinitions(instance),
		"Proxies":            caclculateProxyDefinitions(instance),
		"DashboardEnabled":   instance.STFConfig().TraefikDashboardEnabled(),
		"DashboardWhitelist": instance.STFConfig().TraefikDashboardWhitelistString(),
		"DashboardRule":      fmt.Sprintf("Host(`%s`)", instance.STFConfig().TraefikDashboardHost()),
		"Backtick":           "`",
	}); err != nil {
		return err
	}
	builder := builders.NewDeploymentBuilder(reqLogger, instance, "traefik").
		WithImage(instance.STFConfig().TraefikImage()).
		WithReplicas(instance.STFConfig().TraefikReplicas()).
		WithFile("config.toml", static.String()).
		WithFile("routes/stf.toml", dynamic.String()).
		WithCommand([]string{"traefik"}).
		WithArgs([]string{
			"--configfile", "/etc/configmap/config.toml",
		}).
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithServiceAnnotations(instance.STFConfig().TraefikServiceAnnotations()).
		WithContainerSecurityContext(&corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
				Add:  []corev1.Capability{"NET_BIND_SERVICE"},
			},
		}).
		WithPort("stf-internal", 8880)

	if instance.STFConfig().TLSEnabled() {

		if instance.STFConfig().UseExternalTLS() {
			builder = builder.WithPort("web", 8088).
				WithServiceMap(instance.STFConfig().TraefikServiceType(), []map[string]int32{{"web": 80}})
		} else {
			builder = builder.WithPort("websecure", 8443).
				WithServiceMap(instance.STFConfig().TraefikServiceType(), []map[string]int32{{"websecure": 443}})
		}
		if !instance.STFConfig().UseExternalTLS() && instance.TLSSecret() != "" {
			builder = builder.WithVolumes(
				[]corev1.Volume{
					{
						Name: "tls",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: instance.TLSSecret(),
							},
						},
					},
				},
				[]corev1.VolumeMount{
					{
						Name:      "tls",
						MountPath: "/etc/traefik/ssl",
					},
				},
			)
		}
	} else {
		builder = builder.WithPort("web", 8088).
			WithServiceMap(instance.STFConfig().TraefikServiceType(), []map[string]int32{{"web": 80}})

	}

	for _, group := range instance.DeviceGroups() {
		if !group.UseClusterLocalADB() {
			// Provider ADB Ports
			for _, svc := range calculateProviderSvcDefinitions(instance, group, true) {
				for svcName, svcAttrs := range svc {
					port, _ := strconv.Atoi(svcAttrs.Port)
					if len(svcName) > 15 {
						svcName = svcName[len(svcName)-15:]
					}
					builder = builder.WithPort(svcName, int32(port))
				}
			}
		}
	}

	return builder.Reconcile(c)
}

func newCertificateForCR(instance *androidv1alpha1.AndroidFarm) *cm.Certificate {
	dnsNames := []string{instance.STFConfig().GetAppHostname()}
	if instance.STFConfig().TraefikDashboardEnabled() {
		dnsNames = append(dnsNames, instance.STFConfig().TraefikDashboardHost())
	}
	return &cm.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:            instance.STFConfig().GetAppHostname(),
			Namespace:       instance.STFConfig().GetNamespace(),
			Labels:          instance.STFComponentLabels("ingress"),
			OwnerReferences: instance.OwnerReferences(),
		},
		Spec: cm.CertificateSpec{
			DNSNames:   dnsNames,
			SecretName: instance.TLSSecret(),
			KeySize:    4096,
			IssuerRef:  *instance.STFConfig().GetIssuerReference(),
		},
	}
}

// func buildIngressRoute(instance *androidv1alpha1.AndroidFarm) *traefikv1.IngressRoute {
// 	ingressRoute := &traefikv1.IngressRoute{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:            fmt.Sprintf("%s-ingress", instance.STFNamePrefix()),
// 			Namespace:       instance.STFConfig().GetNamespace(),
// 			Labels:          instance.STFComponentLabels("ingress"),
// 			OwnerReferences: instance.OwnerReferences(),
// 		},
// 		Spec: traefikv1.IngressRouteSpec{
// 			EntryPoints: []string{instance.STFConfig().GetHTTPScheme()},
// 			Routes:      make([]traefikv1.Route, 0),
// 		},
// 	}
// 	if instance.STFConfig().TLSEnabled() {
// 		ingressRoute.Spec.TLS = &traefikv1.TLS{
// 			SecretName: instance.TLSSecret(),
// 		}
// 	}
//
// 	return ingressRoute
// }
