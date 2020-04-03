package emulators

import (
	"fmt"
	"strconv"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newEmulatedDeviceForFarmGroup returns an AndroidDevice configuration for a farm
// group at the given index. If STF is being used for this farm, we also add
// STF annotations.
// TODO: The STF values should be retrieved from a common utility function.
func newEmulatedDeviceForFarmGroup(logger logr.Logger, farm *androidv1alpha1.AndroidFarm, idx int32, group *androidv1alpha1.DeviceGroup, checksum string) *androidv1alpha1.AndroidDevice {
	annotations := make(map[string]string)
	if !farm.STFDisabled() {
		annotations[androidv1alpha1.STFProviderAnnotation] = fmt.Sprintf("%s-provider-%s.%s.svc", farm.STFNamePrefix(), group.Name, farm.STFConfig().GetNamespace())
	}
	annotations[androidv1alpha1.DeviceConfigSHAAnnotation] = checksum
	return &androidv1alpha1.AndroidDevice{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-%s", group.Name, strconv.Itoa(int(idx))),
			Namespace:   group.GetNamespace(),
			Labels:      util.DeviceFarmLabels(farm, group),
			Annotations: annotations,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         farm.APIVersion,
					Kind:               farm.Kind,
					Name:               farm.Name,
					UID:                farm.UID,
					Controller:         util.BoolPointer(true),
					BlockOwnerDeletion: util.BoolPointer(true),
				},
			},
		},
		Spec: androidv1alpha1.AndroidDeviceSpec{
			DeviceConfig: group.Emulators.DeviceConfig,
			ConfigRef:    group.Emulators.ConfigRef,
			Hostname:     group.GetHostname(logger, idx),
			Subdomain:    group.GetSubdomain(),
		},
	}
}
