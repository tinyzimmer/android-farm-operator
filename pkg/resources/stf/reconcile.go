package stf

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type STFReconciler struct {
	resources.FarmReconciler

	client client.Client
	scheme *runtime.Scheme
}

var _ resources.FarmReconciler = &STFReconciler{}

func New(c client.Client, s *runtime.Scheme) resources.FarmReconciler {
	return &STFReconciler{client: c, scheme: s}
}

func (r *STFReconciler) Reconcile(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	// Don't do anything if there is no STF config for the instance
	if instance.STFDisabled() {
		return nil
	}

	reconcileFuncs := []struct {
		msg   string
		rfunc util.FarmReconcileFunc
	}{
		{"Reconciling storage service for OpenSTF", r.reconcileStorage},
		{"Reconciling apk storage service for OpenSTF", r.reconcileAPKStorage},
		{"Reconciling image storage service for OpenSTF", r.reconcileImageStorage},
		{"Reconciling triproxy app service for OpenSTF", r.reconcileTriproxyApp},
		{"Reconciling triproxy device service for OpenSTF", r.reconcileTriproxyDev},
		{"Reconciling processor for OpenSTF", r.reconcileProcessor},
		{"Reconciling reaper for OpenSTF", r.reconcileReaper},
		{"Reconciling websocket service for OpenSTF", r.reconcileWebsocket},
		{"Reconciling API service for OpenSTF", r.reconcileAPI},
		{"Reconciling App service for OpenSTF", r.reconcileApp},
		{"Reconciling Auth service for OpenSTF", r.reconcileAuth},
		{"Reconciling Device Group Providers for OpenSTF", r.reconcileProviders},
		{"Reconciling Traefik service for OpenSTF", r.reconcileIngress}, // Use traefik for rethinkdb-proxy?
	}

	for _, task := range reconcileFuncs {
		reqLogger.Info(task.msg)
		if err := task.rfunc(reqLogger, instance); err != nil {
			return err
		}
	}

	return nil
}
