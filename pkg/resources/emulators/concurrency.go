package emulators

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getDeviceIdx returns the ordinal number of a device in a device group
func getDeviceIdx(devName string) (int64, error) {
	spl := strings.Split(devName, "-")
	// Make sure it will fit into an int32
	return strconv.ParseInt(spl[len(spl)-1], 10, 32)
}

// groupIndexReadyToCreate takes a device group, a management policy, and a device
// index and determines if it is ready to be created. Devices are iterated in numerical
// order up until this device index. If a device does not exist or is not finished booting,
// it is marked as pending. After all devices are iterated, if the pending count is greater
// than the allowed concurrency, return a requeue.
func groupIndexReadyToCreate(c client.Client, group *androidv1alpha1.DeviceGroup, policy *androidv1alpha1.DeviceManagementPolicy, devidx int32) error {
	pending := int32(0)
	for i := int32(0); i < devidx; i++ {
		// Look up the device at this index
		devName := fmt.Sprintf("%s-%s", group.Name, util.DeviceIntToString(int(i)))
		nn := types.NamespacedName{Name: devName, Namespace: group.GetNamespace()}
		found := &androidv1alpha1.AndroidDevice{}
		if err := c.Get(context.TODO(), nn, found); err != nil {
			if client.IgnoreNotFound(err) == nil {
				// Could not find the device, assume it is pending creation
				pending++
				continue
			}
			// return all other api errors
			return err
		}
		// check if boot-completed annotation exists (this is provided by the AndroidDevice controller)
		if found.GetAnnotations() != nil {
			if bootCompleted, ok := found.Annotations[androidv1alpha1.BootCompletedAnnotation]; ok {
				if bootCompleted == "true" {
					// This device is ready
					continue
				}
			}
		}
		// The device is still pending
		pending++
	}
	if pending >= policy.GetConcurrency() {
		return errors.NewRequeueError("Waiting to create device due to concurrency policy", 5)
	}
	return nil
}

// groupIndexReadyToUpdateFunc returns a function that will return true if the
// provided device index is ready to be updated. Like the create func above,
// devices are iterated up until the provided index. See logs below for the checks
// made on each device.
// TODO: Both of these functions desperately need unit tests
func groupIndexReadyToUpdateFunc(reqLogger logr.Logger, c client.Client, farm *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup, devidx int32) func(string) bool {
	return func(newChecksum string) bool {
		// If there is no management policy for this group, return true immediately
		policy := farm.GetDeviceManagementPolicy(group.Name)
		if policy == nil {
			return true
		}

		// Iterate the devices
		pending := int32(0)
		for i := int32(0); i < devidx; i++ {
			devName := fmt.Sprintf("%s-%s", group.Name, util.DeviceIntToString(int(i)))
			nn := types.NamespacedName{Name: devName, Namespace: group.GetNamespace()}
			found := &androidv1alpha1.AndroidDevice{}
			if err := c.Get(context.TODO(), nn, found); err != nil {
				if client.IgnoreNotFound(err) == nil {
					pending++
					continue
				}
				reqLogger.Error(err, "Error looking up device in farm, not allowing update")
				return false
			}
			if found.GetAnnotations() == nil {
				reqLogger.Info("Found device with no annotations, marking it as pending")
				pending++
				continue
			} else if _, ok := found.Annotations[androidv1alpha1.DeviceConfigSHAAnnotation]; !ok {
				reqLogger.Info("Found device with no config checksum, marking it as pending")
				pending++
				continue
			}
			currentChecksum := found.Annotations[androidv1alpha1.DeviceConfigSHAAnnotation]
			if newChecksum != currentChecksum {
				reqLogger.Info("Existing device's config checksum does not match the new one, marking as pending")
				pending++
				continue
			}
			bootCompleted, ok := found.Annotations[androidv1alpha1.BootCompletedAnnotation]
			if !ok {
				reqLogger.Info("Existing device has no boot status annotation, marking as pending")
				pending++
				continue
			}
			if bootCompleted != "true" {
				reqLogger.Info("Device is still booting, marking as pending")
				pending++
				continue
			}
		}
		// Return whether or not the number of pending devices is less than the maximum
		// allowed concurrency.
		return pending < policy.GetConcurrency()
	}
}
