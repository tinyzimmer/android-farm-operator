package emulators

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/android"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EmulatorDeviceReconciler represents a reconciler for AndroidDevices.
type EmulatorDeviceReconciler struct {
	resources.FarmReconciler

	client client.Client
	scheme *runtime.Scheme
}

// NewForDevice returns a reconciler for an AndroidDevice object
func NewForDevice(c client.Client, s *runtime.Scheme) resources.DeviceReconciler {
	return &EmulatorDeviceReconciler{client: c, scheme: s}
}

// Reconcile reconciles an AndroidDevice in the cluster with its desired state.
func (r *EmulatorDeviceReconciler) Reconcile(reqLogger logr.Logger, instance *androidv1alpha1.AndroidDevice) error {
	// TODO : Use a finalizer for this
	if instance.GetDeletionTimestamp() != nil && instance.IsFarmedDevice() {
		farm, err := instance.GetFarm(r.client)
		if err != nil {
			reqLogger.Error(err, "Could not find farm for deleted device, won't be able to clean up rethinkdb")
			return nil
		}
		if err := removeFromRethinkDB(reqLogger, farm, instance); err != nil {
			reqLogger.Error(err, "Failed to remove device from rethinkdb, leaving some garbage behind")
			return nil
		}
	}

	reqLogger.Info("Reconciling pod for android device", "DeviceName", instance.Name, "DeviceNamespace", instance.Namespace)

	config, err := instance.GetConfig(r.client)
	if err != nil {
		return err
	}

	pod := newPodForDevice(instance, config)

	if len(config.Spec.Volumes) > 0 {
		for _, vol := range config.Spec.Volumes {
			// retrieve an existing pvc or create a new one
			pvc, err := reconcilePVCForPod(reqLogger, r.client, pod, vol, util.DeviceLabels(instance))
			if err != nil {
				return err
			}
			// attach the pvc to the pod
			// TODO : This could be done with getters from newPodForDevice
			pod = appendPVCToPod(pod, pvc, vol)
		}
	}

	if created, err := util.ReconcilePod(reqLogger, r.client, pod); err != nil {
		return err
	} else if created {
		return errors.NewRequeueError("Requeueing to check pod boot progress", 3)
	}

	// Fetch the created pod
	found := &corev1.Pod{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		return errors.NewRequeueError("Could not find pod for new android device", 3)
	}

	// Get our ADB Port
	adbPort, err := util.GetPodADBPort(*found)
	if err != nil {
		return err
	}

	// connect to the device and check boot status
	if found.Status.PodIP == "" {
		return errors.NewRequeueError("The device has not yet been assigned a private IP address", 3)
	}
	reqLogger.Info(fmt.Sprintf("Connecting to android device %s on %s:%d", found.Name, found.Status.PodIP, adbPort))
	sess, err := android.NewSession(reqLogger, found.Status.PodIP, adbPort)
	if err != nil {
		return err
	}
	defer sess.Close()
	if complete, err := sess.BootCompleted(); err != nil {
		if strings.Contains(err.Error(), "device offline") {
			if err := resetDeviceAnnotations(r.client, instance); err != nil {
				return err
			}
			return errors.NewRequeueError("ADB needs some time to catch up...", 3)
		}
		reqLogger.Error(err, "Unhandled ADB error while checking boot status")
		return errors.NewRequeueError("ADB needs some time to catch up...", 3)
	} else if !complete {
		if err := resetDeviceAnnotations(r.client, instance); err != nil {
			return err
		}
		return errors.NewRequeueError("Device is still booting", 3)
	}

	// mark device as booted
	if instance.GetAnnotations() == nil {
		instance.Annotations = make(map[string]string)
	}
	if booted, ok := instance.Annotations[androidv1alpha1.BootCompletedAnnotation]; !ok || booted != "true" {
		instance.Annotations[androidv1alpha1.BootCompletedAnnotation] = "true"
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return err
		}
		return errors.NewRequeueError("Marking device as finished booting and requeueing", 1)
	}

	// Check if we are binding this device to an ADB server
	if err := reconcileSTFBinding(reqLogger, r.client, instance, found); err != nil {
		return err
	}

	return nil
}

// resetDeviceAnnotations sets all boot/adb status annotations to false for a device,
// and then updates the remote state if necessary.
func resetDeviceAnnotations(c client.Client, device *androidv1alpha1.AndroidDevice) error {
	annotations := device.GetAnnotations()
	if annotations == nil {
		return nil
	}
	var changed bool
	if adbConnected, ok := annotations[androidv1alpha1.ADBConnectedAnnotation]; ok {
		if adbConnected == "true" {
			annotations[androidv1alpha1.ADBConnectedAnnotation] = "false"
			changed = true
		}
	}
	if bootCompleted, ok := annotations[androidv1alpha1.BootCompletedAnnotation]; ok {
		if bootCompleted == "true" {
			annotations[androidv1alpha1.BootCompletedAnnotation] = "false"
			changed = true
		}
	}
	if changed {
		device.SetAnnotations(annotations)
		return c.Update(context.TODO(), device)
	}
	return nil
}
