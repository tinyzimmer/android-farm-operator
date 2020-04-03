package rethinkdb

import (
	"context"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

// checkRethinkDBIsReady returns nil if the rethinkdb cluster is ready. If
// it isn't a requeue error is returned.
func (r *RethinkDBReconciler) checkRethinkDBIsReady(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	nn := types.NamespacedName{
		Name:      instance.RethinkDBName(),
		Namespace: instance.STFConfig().GetNamespace(),
	}
	ss := &appsv1.StatefulSet{}
	if err := r.client.Get(context.TODO(), nn, ss); err != nil {
		return err
	}

	if ss.Status.ReadyReplicas != *instance.STFConfig().RethinkDBReplicas() {
		return errors.NewRequeueError("Requeing until RethinkDB is ready", 5)
	}

	return nil
}

// checkRethinkDBProxyIsReady returns nil if the rethinkdb proxy is ready. If
// it isn't a requeue error is returned.
func (r *RethinkDBReconciler) checkRethinkDBProxyIsReady(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	nn := types.NamespacedName{
		Name:      instance.RethinkDBProxyName(),
		Namespace: instance.STFConfig().GetNamespace(),
	}
	ss := &appsv1.StatefulSet{}
	if err := r.client.Get(context.TODO(), nn, ss); err != nil {
		return err
	}

	if ss.Status.ReadyReplicas != *instance.STFConfig().RethinkDBProxyReplicas() {
		return errors.NewRequeueError("Requeing until RethinkDB is ready", 5)
	}

	return nil
}
