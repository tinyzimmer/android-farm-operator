package v1alpha1

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NamespacedName returns the namespaced name object for this device
func (a *AndroidDevice) NamespacedName() types.NamespacedName {
	return types.NamespacedName{Name: a.Name, Namespace: a.Namespace}
}

// IsFarmedDevice returns true if this device instance is managed by an
// AndroidFarm.
func (a *AndroidDevice) IsFarmedDevice() bool {
	if a.Labels != nil {
		if _, ok := a.Labels[DeviceFarmLabel]; ok {
			return true
		}
	}
	return false
}

// ConfigChecksum returns the checksum of the configuration currently present
// on the device pods.
func (a *AndroidDevice) ConfigChecksum() string {
	if sum, ok := a.Annotations[DeviceConfigSHAAnnotation]; ok {
		return sum
	}
	return ""
}

// GetConfig returns the desired configuration state of this device instance.
// The configref is looked up if provided, and then any overrides are merged on
// top of it.
func (a *AndroidDevice) GetConfig(c client.Client) (*AndroidDeviceConfig, error) {
	found := &AndroidDeviceConfig{}
	if a.Spec.ConfigRef != nil {
		if err := c.Get(
			context.TODO(),
			types.NamespacedName{Name: a.Spec.ConfigRef.Name, Namespace: metav1.NamespaceAll},
			found,
		); err != nil {
			return nil, err
		}
	}
	if a.Spec.DeviceConfig != nil {
		merged, err := a.Spec.DeviceConfig.MergeInto(found.Spec)
		if err != nil {
			return nil, err
		}
		found.Spec = *merged
	}
	return found, nil
}

// GetFarm returns the parent AndroidFarm for the current device.
// TODO: This could be grabbed from the owner reference as well.
func (a *AndroidDevice) GetFarm(c client.Client) (*AndroidFarm, error) {
	if a.Labels == nil {
		return nil, errors.New("Labels are empty for this device, unable to determine group")
	}
	farm, ok := a.Labels[DeviceFarmLabel]
	if !ok {
		return nil, errors.New("No farm referenced for this device from labels")
	}
	namespacedName := types.NamespacedName{Name: farm, Namespace: metav1.NamespaceAll}
	found := &AndroidFarm{}
	if err := c.Get(context.TODO(), namespacedName, found); err != nil {
		return nil, err
	}
	return found, nil
}

// GetDeviceGroup returns the device group for this device.
func (a *AndroidDevice) GetDeviceGroup(c client.Client) (*DeviceGroup, error) {
	group, ok := a.Labels[DeviceGroupLabel]
	if !ok {
		return nil, errors.New("No farm group referenced for device from labels")
	}
	farm, err := a.GetFarm(c)
	if err != nil {
		return nil, err
	}
	for _, farmGroup := range farm.DeviceGroups() {
		if group == farmGroup.Name {
			return farmGroup, nil
		}
	}
	return nil, errors.New("Failed to locate device group for device")
}
