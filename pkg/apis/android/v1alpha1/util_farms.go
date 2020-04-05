package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MatchingLabels returns the select for finding devices/pods belonging to this
// farm instance.
func (a *AndroidFarm) MatchingLabels() client.MatchingLabels {
	return client.MatchingLabels{DeviceFarmLabel: a.Name}
}

// DeviceGroups returns the device groups for this farm instance.
func (a *AndroidFarm) DeviceGroups() []*DeviceGroup {
	if a.Spec.DeviceGroups == nil {
		return []*DeviceGroup{}
	}
	return a.Spec.DeviceGroups
}

// GetDeviceManagementPolicy returns the device management policy for a device
// group. If one is provided on the group level, it is returned immediately,
// otherwise any global policy on the AndroidFarm is returned.
func (a *AndroidFarm) GetDeviceManagementPolicy(group string) *DeviceManagementPolicy {
	for _, devGroup := range a.DeviceGroups() {
		if devGroup.Name == group {
			if devGroup.Emulators != nil {
				if devGroup.Emulators.DeviceManagementPolicy != nil {
					return devGroup.Emulators.DeviceManagementPolicy
				}
			}
		}
	}
	if a.Spec.DeviceManagementPolicy != nil {
		return a.Spec.DeviceManagementPolicy
	}
	return nil
}

func (a *AndroidFarm) GetGroupADBAdvertiseURL(group *DeviceGroup) string {
	if !group.UseClusterLocalADB() {
		return a.STFConfig().GetAppHostname()
	}
	return fmt.Sprintf("%s.%s.svc", group.GetProviderName(), a.STFConfig().GetNamespace())
}

// GetPodManagementPolicy returns the pod management policy for the device
// mnanagement policy. Currently there is only "OrderedReady".
func (d *DeviceManagementPolicy) GetPodManagementPolicy() PodManagementPolicy {
	if d.PodManagementPolicy == "" {
		return GroupedOrderedReady
	}
	return d.PodManagementPolicy
}

// GetConcurrency returns the maximum number of devices that can be in the
// booting state for this policy instance.
func (d *DeviceManagementPolicy) GetConcurrency() int32 {
	if d.Concurrency == 0 {
		return 1
	}
	return d.Concurrency
}

func (g *DeviceGroup) UseClusterLocalADB() bool {
	return g.Provider != nil && g.Provider.ClusterLocalADB
}

// MatchingLabels returns the selector for finding devices/pods in this
// device group.
func (g *DeviceGroup) MatchingLabels() client.MatchingLabels {
	return client.MatchingLabels{DeviceGroupLabel: g.Name}
}

// IsEmulatedGroup returns true if the device group is a group of emulated devices
// on the kubernetes cluster
func (g *DeviceGroup) IsEmulatedGroup() bool {
	return g.Emulators != nil && g.HostUSB == nil
}

// IsUSBGroup returns true if the device group is for USB devices on nodes
func (g *DeviceGroup) IsUSBGroup() bool {
	return g.HostUSB != nil && g.Emulators == nil
}

// GetNamespace returns the namespace that devices in this group should be
// provisioned in.
func (g *DeviceGroup) GetNamespace() string {
	if g.Emulators == nil {
		return ""
	}
	if g.Emulators.Namespace != "" {
		return g.Emulators.Namespace
	}
	return "default"
}

// GetCount returns the number of devices that should be in this device group.
func (f *DeviceGroup) GetCount() int32 {
	if f.Emulators == nil {
		return 0
	}
	return f.Emulators.Count
}

// GetConfig returns the desired configuration for this device group. If a config
// reference is provided, it is looked up first. Then any overrides in the group
// itself are merged on top of it.
func (f *DeviceGroup) GetConfig(c client.Client) (*AndroidDeviceConfig, error) {
	found := &AndroidDeviceConfig{}
	if f.Emulators.ConfigRef != nil {
		if err := c.Get(
			context.TODO(),
			types.NamespacedName{Name: f.Emulators.ConfigRef.Name, Namespace: metav1.NamespaceAll},
			found,
		); err != nil {
			return nil, err
		}
	}
	if f.Emulators.DeviceConfig != nil {
		merged, err := f.Emulators.DeviceConfig.MergeInto(found.Spec)
		if err != nil {
			return nil, err
		}
		found.Spec = *merged
	}
	return found, nil
}

// GetHostname returns the device hostname for a device in this group at the given
// index.
func (f *DeviceGroup) GetHostname(logger logr.Logger, idx int32) string {
	if f.Emulators.HostnameTemplate == "" {
		return ""
	}
	t, err := template.New("hostname").Parse(f.Emulators.HostnameTemplate)
	if err != nil {
		logger.Error(err, "Could not parse hostname template, falling back to k8s default")
		return ""
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]string{
		"Index": strconv.Itoa(int(idx)),
	})
	if err != nil {
		logger.Error(err, "Failed to execute hostname template, falling back to k8s default")
		return ""
	}
	return buf.String()
}

// GetSubdomain returns the subdomain (and headless service name) to be used for
// devices in this group.
func (f *DeviceGroup) GetSubdomain() string {
	if f.Emulators.Subdomain != "" {
		return f.Emulators.Subdomain
	}
	return f.Name
}
