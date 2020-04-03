package androidfarm

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources/emulators"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources/rethinkdb"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources/stf"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_androidfarm")

var farmFinalizer = "finalizer.androidfarms.android.stf.io"

// Add creates a new AndroidFarm Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAndroidFarm{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("androidfarm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AndroidFarm
	err = c.Watch(&source.Kind{Type: &androidv1alpha1.AndroidFarm{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch devices and requeue the parent farm
	err = c.Watch(&source.Kind{Type: &androidv1alpha1.AndroidDevice{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch configmaps and requeue the parent farm
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch services and requeue the parent farm
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch deployments and requeue the parent farm
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch statefulsets and requeue the parent farm
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch device configs and requeue farms that reference them
	err = c.Watch(
		&source.Kind{Type: &androidv1alpha1.AndroidDeviceConfig{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
				reqs, err := getAffectedFarms(mgr.GetClient(), a.Meta.GetName())
				if err != nil {
					fmt.Println("Error requeuing farms:", err)
				}
				return reqs
			}),
		})
	if err != nil {
		return err
	}

	return nil
}

// getAffectedFarms lists all farms, and returns the ones that reference
// the given device config.
func getAffectedFarms(c client.Client, config string) ([]reconcile.Request, error) {
	reqs := make([]reconcile.Request, 0)
	farms := &androidv1alpha1.AndroidFarmList{}
	if err := c.List(context.TODO(), farms, client.InNamespace(metav1.NamespaceAll)); err != nil {
		return reqs, err
	}
	for _, farm := range farms.Items {
		for _, group := range farm.DeviceGroups() {
			if group.IsEmulatedGroup() {
				if group.Emulators.ConfigRef != nil && group.Emulators.ConfigRef.Name == config {
					reqs = append(reqs, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      farm.Name,
							Namespace: farm.Namespace,
						},
					})
				}
			}
		}
	}
	return reqs, nil
}

// blank assignment to verify that ReconcileAndroidFarm implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAndroidFarm{}

// ReconcileAndroidFarm reconciles a AndroidFarm object
type ReconcileAndroidFarm struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AndroidFarm object and makes changes based on the state read
// and what is in the AndroidFarm.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAndroidFarm) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AndroidFarm")

	// Fetch the AndroidFarm instance
	instance := &androidv1alpha1.AndroidFarm{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// check if marked for deletion and run finalizers
	if util.IsMarkedForDeletion(instance) {
		return r.runFinalizers(reqLogger, instance)
	}

	// populate reconcilers for the instance
	reconcilers := []resources.FarmReconciler{
		rethinkdb.New(r.client, r.scheme),
		stf.New(r.client, r.scheme),
		emulators.NewForFarm(r.client, r.scheme),
	}

	// run each reconciler
	for _, r := range reconcilers {
		if err := r.Reconcile(reqLogger, instance); err != nil {
			if requeue, ok := errors.IsRequeueError(err); ok {
				reqLogger.Info(err.Error())
				return reconcile.Result{
					Requeue:      true,
					RequeueAfter: requeue.Duration(),
				}, nil
			}
			return reconcile.Result{}, err
		}
	}

	// ensure finalizer for cleanup
	if !contains(instance.GetFinalizers(), farmFinalizer) {
		instance.SetFinalizers(append(instance.GetFinalizers(), farmFinalizer))
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileAndroidFarm) runFinalizers(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) (reconcile.Result, error) {
	for _, group := range instance.Spec.DeviceGroups {
		if err := r.client.DeleteAllOf(
			context.TODO(),
			&androidv1alpha1.AndroidDevice{},
			client.InNamespace(group.GetNamespace()),
			client.MatchingLabels(util.DeviceFarmLabels(instance, group)),
		); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return reconcile.Result{}, err
			}
			reqLogger.Info("No matching pods found for device group", "DeviceGroup", group.Name, "Namespace", group.GetNamespace())
			continue
		}
		reqLogger.Info("Sent delete for all pods in device group", "DeviceGroup", group.Name, "Namespace", group.GetNamespace())
	}

	// remove the finalizer
	instance.SetFinalizers(remove(instance.GetFinalizers(), farmFinalizer))
	if err := r.client.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func remove(in []string, rm string) []string {
	out := make([]string, 0)
	for _, x := range in {
		if x != rm {
			out = append(out, x)
		}
	}
	return out
}
