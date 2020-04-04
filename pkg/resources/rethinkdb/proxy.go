package rethinkdb

import (
	"fmt"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// proxyStartScript is the start script used for the proxy instances
var proxyStartScript = `
set -exo pipefail

while ! getent hosts ${WAIT_SERVERS} ; do sleep 3 ; done

exec rethinkdb proxy \
	${JOIN_SERVERS} \
	--bind all \
	--canonical-address ${POD_IP}
`

// reconcileProxy ensures the proxy servers for the rethinkdb statefulset. We use
// a stateful set so we can provider static server definitions to the traefik instance.
// This in turn allows for sticky sessions with the backend servers when using
// the rethinkdb admin interface.
func (r *RethinkDBReconciler) reconcileProxy(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	clusterSuffix := util.GetClusterSuffix()
	name := instance.RethinkDBName()
	namespace := instance.STFConfig().GetNamespace()
	joinStr := ""
	waitStr := ""
	for i := int32(0); i < *instance.STFConfig().RethinkDBReplicas(); i++ {
		addr := fmt.Sprintf("%s-%d.%s.%s.svc.%s", name, i, name, namespace, clusterSuffix)
		joinStr = joinStr + fmt.Sprintf(" --join %s:29015 ", addr)
		waitStr = waitStr + fmt.Sprintf(" %s ", addr)
	}
	ss := newProxyStatefulSetForCR(instance, joinStr, waitStr)
	if err := util.ReconcileStatefulSet(reqLogger, r.client, ss); err != nil {
		return err
	}
	return nil
}

// newProxyStatefulSetForCR returns a new proxy statefulset for an AndroidFarm instance.
func newProxyStatefulSetForCR(cr *androidv1alpha1.AndroidFarm, joinStr, waitStr string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.RethinkDBProxyName(),
			Namespace:       cr.STFConfig().GetNamespace(),
			Labels:          cr.RethinkDBLabels("proxy"),
			OwnerReferences: cr.OwnerReferences(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: cr.STFConfig().RethinkDBProxyReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: cr.RethinkDBLabels("proxy"),
			},
			ServiceName: cr.RethinkDBProxyName(),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.RethinkDBLabels("proxy"),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.STFConfig().GetServiceAccount(),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: util.BoolPointer(false),
					},
					Containers: []corev1.Container{
						{
							Name:            "rethinkdb-proxy",
							ImagePullPolicy: cr.STFConfig().RethinkDBPullPolicy(),
							Image:           cr.STFConfig().RethinkDBImage(),
							Command:         []string{"/bin/bash", "-c"},
							Args:            []string{proxyStartScript},
							Env: append(cr.RethinkDBEnvVars(), []corev1.EnvVar{
								{
									Name:  "JOIN_SERVERS",
									Value: joinStr,
								},
								{
									Name:  "WAIT_SERVERS",
									Value: waitStr,
								},
								{
									Name:  "CLUSTER_SUFFIX",
									Value: util.GetClusterSuffix(),
								},
							}...),
							Ports: []corev1.ContainerPort{
								{
									Name:          "admin-port",
									ContainerPort: 8080,
								},
								{
									Name:          "driver-port",
									ContainerPort: 28015,
								},
								{
									Name:          "cluster-port",
									ContainerPort: 29015,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.Parse("driver-port"),
									},
								},
							},
							Resources: cr.STFConfig().RethinkDBProxyResourceRequirements(),
						},
					},
				},
			},
		},
	}
}
