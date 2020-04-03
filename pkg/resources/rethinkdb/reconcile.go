package rethinkdb

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	rethinkdbutil "github.com/tinyzimmer/android-farm-operator/pkg/util/rethinkdb"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RethinkDBReconciler represents a reconciler for RethinkDB clusters
type RethinkDBReconciler struct {
	resources.FarmReconciler

	client client.Client
	scheme *runtime.Scheme
}

// Blank assignment to ensure RethinkDBReconciler implements FarmReconciler
var _ resources.FarmReconciler = &RethinkDBReconciler{}

// New returns a new rethinkdb reconciler
func New(c client.Client, s *runtime.Scheme) resources.FarmReconciler {
	return &RethinkDBReconciler{client: c, scheme: s}
}

// Reconcile will reconcile the desired state of an AndroidFarm's rethinkdb
// cluster with the resources running on the Kubernetes cluster.
func (r *RethinkDBReconciler) Reconcile(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	// Don't do anything if there is no STF config for the instance
	if instance.STFDisabled() {
		return nil
	}

	// define our reconcile functions
	reconcileFuncs := []struct {
		msg   string
		rfunc util.FarmReconcileFunc
	}{
		{"Reconciling headless and admin services for RethinkDB", r.reconcileServices},
		{"Reconciling StatefulSet for RethinkDB", r.reconcileStatefulSet},
		{"Checking if RethinkDB is ready", r.checkRethinkDBIsReady},
		{"Reconciling RethinkDB Proxy deployment", r.reconcileProxy},
		{"Checking if RethinkDB proxy is ready", r.checkRethinkDBProxyIsReady},
		{"Ensuring STF DB migration job", r.reconcileMigrationJob},
		{"Ensuring table replication across RethinkDB nodes", rethinkdbutil.EnsureRethinkDBReplicas},
	}

	// Run each function and return any error (can be requeue errors)
	for _, task := range reconcileFuncs {
		reqLogger.Info(task.msg)
		if err := task.rfunc(reqLogger, instance); err != nil {
			return err
		}
	}

	return nil
}
