package emulators

import (
	"context"
	"strings"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
	rdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// runGC runs garbage collection for a given android farm. This function is invoked at
// the end of every reconcile event.
func runGC(reqLogger logr.Logger, c client.Client, farm *androidv1alpha1.AndroidFarm) error {

	// make a map of device groups to their count
	deviceGroups := make(map[string]int32)
	for _, group := range farm.DeviceGroups() {
		if group.IsEmulatedGroup() {
			deviceGroups[group.Name] = group.GetCount()
		}
	}

	// fetch all devices for this farm
	devices := &androidv1alpha1.AndroidDeviceList{}
	if err := c.List(context.TODO(), devices, client.InNamespace(metav1.NamespaceAll), farm.MatchingLabels()); err != nil {
		return err
	}

	// Iterate and delete devices that either don't reference a group, or reference
	// a group that no longer exists. Finally, check if their device index is larger
	// than the number of devices supposed to be in the group.
	for _, device := range devices.Items {
		var delete bool = false
		if group, ok := device.Labels[androidv1alpha1.DeviceGroupLabel]; !ok {
			reqLogger.Info("Deleting device with no group label", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
			delete = true
		} else if count, ok := deviceGroups[group]; !ok {
			reqLogger.Info("Deleting device with reference to group that no longer exists", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
			delete = true
		} else {
			devidx, err := getDeviceIdx(device.Name)
			if err != nil {
				return err
			}
			if int32(devidx) > count-1 {
				reqLogger.Info("Deleting device as the group is being scaled down", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
				delete = true
			}
		}
		// If we are deleting, remove the device from STF and delete it's resources
		if delete {
			if err := removeFromRethinkDB(reqLogger, farm, &device); err != nil {
				return err
			}
			if err := c.Delete(context.TODO(), &device); err != nil {
				return err
			}
		}
	}

	return nil
}

// removeFromRethinkDB will remove the given device from the rethinkdb instance
// of its farm. This is not a crucial function as its only purpose is to keep
// the OpenSTF UI tidy.
func removeFromRethinkDB(reqLogger logr.Logger, farm *androidv1alpha1.AndroidFarm, device *androidv1alpha1.AndroidDevice) error {
	if device.GetAnnotations() == nil {
		// TODO - I mean we could. Some more utility functions can be defined to
		// deal with values like this, instead of annotations.
		reqLogger.Info("Device has no annotations, can't determine serial", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
		return nil
	}
	serial, ok := device.Annotations[androidv1alpha1.ProviderSerialAnnotation]
	if !ok {
		reqLogger.Info("Device has no provider serial annotations, can't determine stf name", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
		return nil
	}
	// connect to the master rethinkdb instance
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: strings.TrimPrefix(stfutil.RethinkDBProxyEndpoint(farm), "tcp://"),
	})
	if err != nil {
		return err
	}
	// remove any device from the devices table that matches the serial of the removed
	// device.
	_, err = rdb.DB("stf").Table("devices").Filter(func(uu rdb.Term) rdb.Term {
		return uu.Field("serial").Eq(serial)
	}).Delete().Run(session)
	return err
}
