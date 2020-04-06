package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
)

func RethinkDBProxyEndpoint(a *androidv1alpha1.AndroidFarm) string {
	return fmt.Sprintf("tcp://%s:28015", RethinkDBProxyURI(a))
}

func RethinkDBAdminEndpoint(a *androidv1alpha1.AndroidFarm) string {
	return fmt.Sprintf("%s:8080", RethinkDBProxyURI(a))
}

func RethinkDBProxyURI(a *androidv1alpha1.AndroidFarm) string {
	return fmt.Sprintf("%s.%s.svc.%s",
		a.RethinkDBProxyName(),
		a.STFConfig().GetNamespace(),
		util.GetClusterSuffix(),
	)
}

func RethinkDBProxyIndexURI(a *androidv1alpha1.AndroidFarm, idx int32) string {
	return fmt.Sprintf("%s-%d.%s", a.RethinkDBProxyName(), idx, RethinkDBProxyURI(a))
}

func TriproxyEndpoint(a *androidv1alpha1.AndroidFarm, component string) string {
	return fmt.Sprintf(
		"tcp://%s-triproxy-%s.%s.svc.%s",
		a.STFNamePrefix(),
		component,
		a.STFConfig().GetNamespace(),
		util.GetClusterSuffix(),
	)
}

func GetGroupADBAdvertiseURL(farm *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) string {
	if hostname := group.ProviderHostnameOverride(); hostname != "" {
		return hostname
	}
	if !group.UseClusterLocalADB() {
		return farm.STFConfig().GetAppHostname()
	}
	return fmt.Sprintf("%s-%s-traefik.%s.svc.%s", farm.STFNamePrefix(), group.GetProviderName(), farm.STFConfig().GetNamespace(), util.GetClusterSuffix())
}
