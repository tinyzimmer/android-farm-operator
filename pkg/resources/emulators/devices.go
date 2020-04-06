package emulators

import (
	"fmt"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newEmulatedDeviceForFarmGroup returns an AndroidDevice configuration for a farm
// group at the given index. If STF is being used for this farm, we also add
// STF annotations.
func newEmulatedDeviceForFarmGroup(logger logr.Logger, farm *androidv1alpha1.AndroidFarm, idx int32, group *androidv1alpha1.DeviceGroup, checksum string) *androidv1alpha1.AndroidDevice {
	annotations := make(map[string]string)
	if !farm.STFDisabled() {
		annotations[androidv1alpha1.STFProviderAnnotation] = fmt.Sprintf("%s.%s.svc", group.GetProviderName(), farm.STFConfig().GetNamespace())
	}
	annotations[androidv1alpha1.DeviceConfigSHAAnnotation] = checksum
	return &androidv1alpha1.AndroidDevice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s", group.Name, util.DeviceIntToString(int(idx))),
			Namespace:       group.GetNamespace(),
			Labels:          util.DeviceFarmLabels(farm, group),
			Annotations:     annotations,
			OwnerReferences: farm.OwnerReferences(),
		},
		Spec: androidv1alpha1.AndroidDeviceSpec{
			DeviceConfig: group.Emulators.DeviceConfig,
			ConfigRef:    group.Emulators.ConfigRef,
			Hostname:     group.GetHostname(logger, util.DeviceIntToString(int(idx))),
			Subdomain:    group.GetSubdomain(),
		},
	}
}
