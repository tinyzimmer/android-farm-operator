package rethinkdb

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
)

// reconcileMigrationJob ensures that an stf database migration job has run
// in an AndroidFarm's rethinkdb cluster.
func (r *RethinkDBReconciler) reconcileMigrationJob(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	job := newJobForCR(instance)
	return util.ReconcileJob(reqLogger, r.client, job, true)
}

// newJobForCR returns a new stf migration job definition for an AndroidFarm.
func newJobForCR(cr *androidv1alpha1.AndroidFarm) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-migrate", cr.RethinkDBName()),
			Namespace:       cr.STFConfig().GetNamespace(),
			Labels:          cr.RethinkDBLabels("migrat"),
			OwnerReferences: cr.OwnerReferences(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.RethinkDBLabels("migrate"),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.STFConfig().GetServiceAccount(),
					RestartPolicy:      "OnFailure",
					SecurityContext:    cr.STFConfig().PodSecurityContext(),
					Containers: []corev1.Container{
						{
							Name:            "stf-migrate",
							Image:           cr.STFConfig().OpenSTFImage(),
							Command:         []string{"/bin/bash", "-c"},
							Args:            []string{"while ! getent hosts ${PROXY_NAME} ; do sleep 3 ; done && stf migrate"},
							SecurityContext: cr.STFConfig().ContainerSecurityContext(),
							Env: append(cr.RethinkDBEnvVars(), corev1.EnvVar{
								Name:  "RETHINKDB_PORT_28015_TCP",
								Value: stfutil.RethinkDBProxyEndpoint(cr),
							}),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    *resource.NewMilliQuantity(100, resource.DecimalSI),
									"memory": *resource.NewQuantity(256*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}
