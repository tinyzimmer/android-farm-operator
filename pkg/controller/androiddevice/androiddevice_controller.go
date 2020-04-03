package androiddevice

import (
	"context"
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources/emulators"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
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

var log = logf.Log.WithName("controller_androiddevice")

// Add creates a new AndroidDevice Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAndroidDevice{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("androiddevice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AndroidDevice
	err = c.Watch(&source.Kind{Type: &androidv1alpha1.AndroidDevice{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner AndroidDevice
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidDevice{},
	})
	if err != nil {
		return err
	}

	// Watch created volumes and requeue for owner
	// TOOD: Currently we leave volumes in tact - should this behavior change?
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &androidv1alpha1.AndroidFarm{},
	})
	if err != nil {
		return err
	}

	// Watch device configs and requeue non-farmed devices that use them.
	err = c.Watch(
		&source.Kind{Type: &androidv1alpha1.AndroidDeviceConfig{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
				reqs, err := getAffectedDevices(mgr.GetClient(), a.Meta.GetName())
				if err != nil {
					fmt.Println("Error requeuing devices:", err)
				}
				return reqs
			}),
		})
	if err != nil {
		return err
	}

	return nil
}

// getAffectedDevices returns the non-farmed devices that are affected by a change
// to an AndroidDeviceConfig.
func getAffectedDevices(c client.Client, config string) ([]reconcile.Request, error) {
	reqs := make([]reconcile.Request, 0)
	devices := &androidv1alpha1.AndroidDeviceList{}
	if err := c.List(context.TODO(), devices, client.InNamespace(metav1.NamespaceAll), client.MatchingLabels{androidv1alpha1.DeviceConfigLabel: config}); err != nil {
		return reqs, err
	}
	for _, dev := range devices.Items {
		if dev.IsFarmedDevice() {
			continue
		}
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      dev.Name,
				Namespace: dev.Namespace,
			},
		})
	}
	return reqs, nil
}

// blank assignment to verify that ReconcileAndroidDevice implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAndroidDevice{}

// ReconcileAndroidDevice reconciles a AndroidDevice object
type ReconcileAndroidDevice struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AndroidDevice object and makes changes based on the state read
// and what is in the AndroidDevice.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAndroidDevice) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AndroidDevice")

	// Fetch the AndroidDevice instance
	instance := &androidv1alpha1.AndroidDevice{}
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

	reconciler := emulators.NewForDevice(r.client, r.scheme)
	if err := reconciler.Reconcile(reqLogger, instance); err != nil {
		if requeue, ok := errors.IsRequeueError(err); ok {
			reqLogger.Info(err.Error())
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: requeue.Duration(),
			}, nil
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
