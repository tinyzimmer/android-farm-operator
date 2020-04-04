package rethinkdb

import (
	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// rethinkDBStartScript is the start script used on rethinkdb cluster nodes
var rethinkDBStartScript = `
set -exo pipefail

ORDINAL=$(echo "${POD_NAME}" | rev | cut -d "-" -f1 | rev)

if [[ "${ORDINAL}" != "0" ]]; then
	while ! getent hosts ${SERVICE_NAME}.${POD_NAMESPACE} ; do sleep 3 ; done
	ENDPOINT="${SERVICE_NAME}-0.${SERVICE_NAME}.${POD_NAMESPACE}.svc.${CLUSTER_SUFFIX}:29015"
	echo "Join to ${SERVICE_NAME} on ${ENDPOINT}"
	exec rethinkdb \
		--bind all \
    --directory /data/rethinkdb_data \
		--join ${ENDPOINT} \
		--server-name ${POD_NAME} \
		--server-tag ${POD_NAME} \
		--server-tag ${NODE_NAME} \
		--canonical-address ${POD_IP}
else
	echo "Start single/master instance"
	exec rethinkdb \
		--bind all \
    --directory /data/rethinkdb_data \
		--server-name ${POD_NAME} \
		--server-tag ${POD_NAME} \
		--server-tag ${NODE_NAME} \
		--canonical-address ${POD_IP}
fi
`

func (r *RethinkDBReconciler) reconcileStatefulSet(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	if err := util.ReconcileStatefulSet(reqLogger, r.client, newStatefulSetForCR(instance)); err != nil {
		return err
	}
	return nil
}

// newStatefulSetForCR returns a new rethinkdb statefulset configuration for the given
// AndroidFarm instance.
func newStatefulSetForCR(cr *androidv1alpha1.AndroidFarm) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.RethinkDBName(),
			Namespace:       cr.STFConfig().GetNamespace(),
			Labels:          cr.RethinkDBLabels("rethinkdb"),
			OwnerReferences: cr.OwnerReferences(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: cr.STFConfig().RethinkDBReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: cr.RethinkDBLabels("rethinkdb"),
			},
			ServiceName:          cr.RethinkDBName(),
			VolumeClaimTemplates: cr.RethinkDBVolumeClaims(),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.RethinkDBLabels("rethinkdb"),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.STFConfig().GetServiceAccount(),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: util.Int64Ptr(1000),
					},
					Volumes: cr.RethinkDBVolumes(),
					Containers: []corev1.Container{
						{
							Name:            "rethinkdb",
							ImagePullPolicy: cr.STFConfig().RethinkDBPullPolicy(),
							Image:           cr.STFConfig().RethinkDBImage(),
							Env: append(cr.RethinkDBEnvVars(), corev1.EnvVar{
								Name:  "CLUSTER_SUFFIX",
								Value: util.GetClusterSuffix(),
							}),
							VolumeMounts: cr.RethinkDBVolumeMounts(),
							Command:      []string{"/bin/bash", "-c"},
							Args:         []string{rethinkDBStartScript},
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
							Resources: cr.STFConfig().RethinkDBResourceRequirements(),
						},
					},
				},
			},
		},
	}
}
