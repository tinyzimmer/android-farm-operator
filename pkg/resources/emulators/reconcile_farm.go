package emulators

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/resources"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EmulatorFarmReconciler represents a reconciler for AndroidFarm CRs
type EmulatorFarmReconciler struct {
	resources.FarmReconciler

	client client.Client
	scheme *runtime.Scheme
}

var _ resources.FarmReconciler = &EmulatorFarmReconciler{}
var _ resources.DeviceReconciler = &EmulatorDeviceReconciler{}

// NewForFarm returns a new reconciler for an AndroidFarm
func NewForFarm(c client.Client, s *runtime.Scheme) resources.FarmReconciler {
	return &EmulatorFarmReconciler{client: c, scheme: s}
}

// Reconcile will reconcile the desired state of the devices for an AndroidFarm
// with what is running in the cluster.
func (r *EmulatorFarmReconciler) Reconcile(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	// Iterate device farms and reconcile devices
	for _, group := range instance.DeviceGroups() {
		reqLogger.Info("Reconciling device group")
		// If it's an emulated device group, reconcile the emulator devices
		if group.IsEmulatedGroup() {
			logger := reqLogger.WithValues("Group", group.Name, "Namespace", group.GetNamespace())
			logger.Info("Device group has an emulated device configuration, reconciling pods")
			if err := r.ReconcileEmulatedDeviceGroup(logger, instance, group); err != nil {
				return err
			}
		}
	}

	// Run garbage collection on the farm
	return runGC(reqLogger, r.client, instance)
}

func (r *EmulatorFarmReconciler) ReconcileEmulatedDeviceGroup(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) error {
	if group.GetCount() == 0 {
		reqLogger.Info("Device group has 0 devices, skipping")
		return nil
	}

	// Get the config
	config, err := group.GetConfig(r.client)
	if err != nil {
		return err
	}

	// get the config's checksum
	checksum, err := config.Checksum()
	if err != nil {
		return err
	}
	reqLogger.Info("Calculated config checksum for device group", "Checksum", checksum)

	// Create a headless service for DNS resolution
	svc := serviceForDeviceGroup(instance, group, config)
	if err := util.ReconcileService(reqLogger, r.client, svc); err != nil {
		return err
	}

	// Create devices for the group
	for i := int32(0); i < group.GetCount(); i++ {
		// check if we are enforcing concurrency
		if policy := instance.GetDeviceManagementPolicy(group.Name); policy != nil {
			if err := groupIndexReadyToCreate(r.client, group, policy, i); err != nil {
				return err
			}
		}
		// Define a new Device object
		reqLogger.Info("Reconciling emulator device for device farm", "Group", group.Name, "PodNumber", i)
		device := newEmulatedDeviceForFarmGroup(reqLogger, instance, i, group, checksum)
		if err := util.ReconcileDevice(reqLogger, r.client, device, groupIndexReadyToUpdateFunc(reqLogger, r.client, instance, group, i)); err != nil {
			return err
		}
	}
	return nil
}
