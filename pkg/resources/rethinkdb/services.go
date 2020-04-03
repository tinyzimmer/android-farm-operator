package rethinkdb

import (
	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// reconcileService ensures the rethinkdb statefulset's headless and admin
// services
func (r *RethinkDBReconciler) reconcileServices(reqLogger logr.Logger, cr *androidv1alpha1.AndroidFarm) error {
	if err := util.ReconcileService(reqLogger, r.client, newServiceForCR(cr)); err != nil {
		return err
	}
	if err := util.ReconcileService(reqLogger, r.client, newProxyServiceForCR(cr)); err != nil {
		return err
	}
	return nil
}

// newServiceForCR returns a new headless service for the rehtinkdb cluster
func newServiceForCR(cr *androidv1alpha1.AndroidFarm) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.RethinkDBName(),
			Namespace:       cr.STFConfig().GetNamespace(),
			Labels:          cr.RethinkDBLabels("rethinkdb"),
			OwnerReferences: cr.OwnerReferences(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       "driver-port",
					Port:       int32(28015),
					TargetPort: intstr.Parse("28015"),
				},
				{
					Name:       "cluster-port",
					Port:       int32(29015),
					TargetPort: intstr.Parse("29015"),
				},
			},
			Selector: cr.RethinkDBLabels("rethinkdb"),
		},
	}
}

// newProxyServiceForCR returns a headless service for the rehtinkdb proxies
func newProxyServiceForCR(cr *androidv1alpha1.AndroidFarm) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.RethinkDBProxyName(),
			Namespace:       cr.STFConfig().GetNamespace(),
			Labels:          cr.RethinkDBLabels("proxy"),
			OwnerReferences: cr.OwnerReferences(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       "admin-port",
					Port:       int32(8080),
					TargetPort: intstr.Parse("8080"),
				},
				{
					Name:       "driver-port",
					Port:       int32(28015),
					TargetPort: intstr.Parse("28015"),
				},
				{
					Name:       "cluster-port",
					Port:       int32(29015),
					TargetPort: intstr.Parse("29015"),
				},
			},
			Selector: cr.RethinkDBLabels("proxy"),
		},
	}
}
