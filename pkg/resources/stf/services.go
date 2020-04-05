package stf

import (
	"fmt"
	"strconv"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
)

type svcDef struct {
	Rule        string
	Endpoints   []string
	Port        string
	TCPPort     string
	Priority    string
	Middlewares []string
	IsProvider  bool
}

func caclculateProxyDefinitions(cr *androidv1alpha1.AndroidFarm) []map[string]svcDef {
	hostname := cr.InternalProxyHost()
	return []map[string]svcDef{
		// Storage Route
		{"stf-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/`))", hostname),
			Endpoints: []string{cr.STFComponentName("storage")},
			Priority:  "9",
			Port:      "3000",
		}},
		// APK Storage Route
		{"stf-apk-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/apk/`))", hostname),
			Endpoints: []string{cr.STFComponentName("apk-storage")},
			Priority:  "10",
			Port:      "3000",
		}},
		// Image Storage Route
		{"stf-img-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/image/`))", hostname),
			Endpoints: []string{cr.STFComponentName("img-storage")},
			Priority:  "10",
			Port:      "3000",
		}},
	}
}

func calculateServiceDefinitions(cr *androidv1alpha1.AndroidFarm) []map[string]svcDef {
	// app hostname
	hostname := cr.STFConfig().GetAppHostname()
	// RethinkDB proxy routes so we can use sticky sessions
	rdbProxies := make([]string, 0)
	for i := int32(0); i < *cr.STFConfig().RethinkDBProxyReplicas(); i++ {
		rdbProxies = append(rdbProxies, stfutil.RethinkDBProxyIndexURI(cr, i))
	}
	svcs := []map[string]svcDef{
		// App Route
		{"stf-app": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-app", cr.STFNamePrefix())},
			Priority:  "5",
			Port:      "3000",
		}},
		// API Route
		{"stf-api": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/api/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-api", cr.STFNamePrefix())},
			Priority:  "10",
			Port:      "3000",
		}},
		// Auth Route
		{"stf-auth": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/auth/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-auth", cr.STFNamePrefix())},
			Priority:  "10",
			Port:      "3000",
		}},
		/// Websocket Route
		{"stf-websocket": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/socket.io/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-websocket", cr.STFNamePrefix())},
			Priority:  "10",
			Port:      "3000",
		}},
		// Storage Route
		{"stf-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-storage", cr.STFNamePrefix())},
			Priority:  "9",
			Port:      "3000",
		}},
		// APK Storage Route
		{"stf-apk-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/apk/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-apk-storage", cr.STFNamePrefix())},
			Priority:  "10",
			Port:      "3000",
		}},
		// Image Storage Route
		{"stf-img-storage": {
			Rule:      fmt.Sprintf("(Host(`%s`) && PathPrefix(`/s/image/`))", hostname),
			Endpoints: []string{fmt.Sprintf("%s-img-storage", cr.STFNamePrefix())},
			Priority:  "10",
			Port:      "3000",
		}},
		// RethinkDB Admin UI
		{"rethinkdb-admin": {
			Rule:        fmt.Sprintf("(Host(`%s`) && PathPrefix(`/rethinkdb`))", hostname),
			Endpoints:   rdbProxies,
			Priority:    "10",
			Port:        "8080",
			Middlewares: []string{"strip-rethinkdb"},
		}},
	}

	for _, group := range cr.DeviceGroups() {
		if group.UseClusterLocalADB() {
			providerName := group.GetProviderName()
			svcs = append(svcs, map[string]svcDef{
				providerName: svcDef{
					Rule:       fmt.Sprintf("(Host(`%s`) && PathPrefix(`/d/%s/`))", cr.STFConfig().GetAppHostname(), providerName),
					Endpoints:  []string{fmt.Sprintf("%s-%s-traefik", cr.STFNamePrefix(), providerName)},
					Priority:   "10",
					Port:       "8088",
					IsProvider: false,
				},
			})
		} else {
			svcs = append(svcs, calculateProviderSvcDefinitions(cr, group, true)...)
		}
	}

	return svcs
}

func calculateProviderSvcDefinitions(cr *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup, toTraefik bool) []map[string]svcDef {
	svcs := make([]map[string]svcDef, 0)
	var maxPort int32
	minPort := group.GetProviderStartPort()
	if group.IsEmulatedGroup() {
		maxPort = getProviderMaxPort(minPort, group.GetCount())
	} else {
		maxPort = getProviderMaxPort(minPort, group.MaxUSBDevices())
	}
	providerName := group.GetProviderName()
	for i := minPort; i <= maxPort; i++ {
		svcName := fmt.Sprintf("%s-%d", providerName, i)
		var endpoint string
		if toTraefik {
			endpoint = fmt.Sprintf("%s-%s-traefik", cr.STFNamePrefix(), providerName)
		} else {
			endpoint = fmt.Sprintf("%s-%s", cr.STFNamePrefix(), providerName)
		}
		svcs = append(svcs, map[string]svcDef{svcName: svcDef{
			Rule:       fmt.Sprintf("(Host(`%s`) && PathPrefix(`/d/%s/{serial:[^/]+}/%d/`))", cr.STFConfig().GetAppHostname(), providerName, i),
			Endpoints:  []string{endpoint},
			Priority:   "10",
			Port:       strconv.Itoa(int(i)),
			IsProvider: true,
		}})
	}
	return svcs
}
