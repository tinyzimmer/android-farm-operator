package stf

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	corev1 "k8s.io/api/core/v1"
)

var providerStartScriptTmpl = template.Must(template.New("provider-start").Parse(`
timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/127.0.0.1/5037; do sleep 1; done'
echo "ADB is available, launching {{ .ProviderName }}"
unset "${!STF@}" ; stf provider \
	--name {{ .ProviderName }} \
	--adb-host 127.0.0.1 \
  --connect-push "tcp://{{ .TriproxyDev }}:7270" \
  --connect-sub "tcp://{{ .TriproxyDev }}:7250" \
	--storage-url "{{ .StorageURL }}" \
	--heartbeat-interval 10000 \
  --allow-remote \
  --screen-ws-url-pattern "{{ .AppWebsocketURL }}/d/{{ .ProviderName }}/<%= serial %>/<%= publicPort %>/" \
  --min-port {{ .MinPort }} --max-port {{ .MaxPort }} \
	{{ if .NoCleanup }}--no-cleanup{{ end }} \
	--public-ip "{{ .ProviderHostname }}"
`))

func (r *STFReconciler) reconcileProviders(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	for _, group := range instance.DeviceGroups() {

		if err := r.reconcileGroupProvider(reqLogger, instance, group); err != nil {
			return err
		}

		if err := r.reconcileGroupProviderTraefik(reqLogger, instance, group); err != nil {
			return err
		}

	}

	return nil
}

func (r *STFReconciler) reconcileGroupProvider(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) error {
	startScript, err := getGroupStartScript(instance, group)
	if err != nil {
		return err
	}

	name := group.GetProviderName()

	builder := builders.NewDeploymentBuilder(reqLogger, instance, name).
		WithResourceRequirements(instance.STFConfig().ProviderResourceRequirements()).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{startScript}).
		WithService("ClusterIP").
		WithRethinkDB().
		WithPodSecurityContext(instance.STFConfig().ADBPodSecurityContext(group)).
		WithContainerSecurityContext(instance.STFConfig().ADBContainerSecurityContext(group)).
		WithSidecar(instance.STFConfig().ADBSidecarContainer(name, group)).
		WithWait()

	var maxPort int32
	if group.IsEmulatedGroup() {
		maxPort = group.GetCount()
	} else if group.IsUSBGroup() {
		maxPort = group.MaxUSBDevices()
		builder = builder.
			WithNodeSelector(group.ProviderNodeSelector()).
			WithVolumes([]corev1.Volume{
				{
					Name: "usb",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/dev/bus/usb",
						},
					},
				}}, nil)
	}

	for i := group.GetProviderStartPort(); i <= getProviderMaxPort(group.GetProviderStartPort(), maxPort); i++ {
		builder = builder.WithPort(fmt.Sprintf("provider-%d", i), i)
	}

	return builder.Reconcile(r.client)
}

func (r *STFReconciler) reconcileGroupProviderTraefik(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) error {
	providerDefs := calculateProviderSvcDefinitions(instance, group, false)

	var static bytes.Buffer
	if err := providerTraefikStaticConfigTmpl.Execute(&static, map[string]interface{}{
		"Services":  providerDefs,
		"AccessLog": group.ProviderTraefikAccessLogsEnabled(),
		"Backtick":  "`",
	}); err != nil {
		return err
	}
	var dynamic bytes.Buffer
	if err := providerTraefikDynamicConfigTmpl.Execute(&dynamic, map[string]interface{}{
		"Services": providerDefs,
		"Backtick": "`",
	}); err != nil {
		return err
	}

	builder := builders.NewDeploymentBuilder(reqLogger, instance, fmt.Sprintf("%s-traefik", group.GetProviderName())).
		WithImage(instance.STFConfig().TraefikImage()).
		WithReplicas(group.ProviderTraefikReplicas()).
		WithResourceRequirements(group.ProviderTraefikResourceRequirements()).
		WithFile("config.toml", static.String()).
		WithFile("routes/stf.toml", dynamic.String()).
		WithCommand([]string{"traefik"}).
		WithArgs([]string{
			"--configfile", "/etc/configmap/config.toml",
		}).
		WithService(group.ProviderTraefikServiceType()).
		WithServiceAnnotations(group.ProviderTraefikServiceAnnotations()).
		WithPort("web", 8088).
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext())

	for _, svc := range providerDefs {
		for svcName, svcAttrs := range svc {
			port, _ := strconv.Atoi(svcAttrs.Port)
			if len(svcName) > 15 {
				svcName = svcName[len(svcName)-15:]
			}
			builder = builder.WithPort(svcName, int32(port))
		}
	}

	return builder.Reconcile(r.client)
}

func getGroupStartScript(instance *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) (string, error) {
	var buf bytes.Buffer
	var deviceCount int32
	if group.IsEmulatedGroup() {
		deviceCount = group.GetCount()
	} else {
		deviceCount = group.MaxUSBDevices()
	}
	if err := providerStartScriptTmpl.Execute(&buf, map[string]interface{}{
		"ProviderName":     group.GetProviderName(),
		"TriproxyDev":      fmt.Sprintf("%s-triproxy-dev", instance.STFNamePrefix()),
		"HTTPScheme":       instance.STFConfig().GetHTTPScheme(),
		"WebsocketScheme":  instance.STFConfig().GetWebsocketScheme(),
		"ProviderHostname": stfutil.GetGroupADBAdvertiseURL(instance, group),
		"MinPort":          strconv.Itoa(int(group.GetProviderStartPort())),
		"MaxPort":          strconv.Itoa(int(getProviderMaxPort(group.GetProviderStartPort(), deviceCount))),
		"StorageURL":       instance.InternalStorageURL(),
		"AppWebsocketURL":  instance.STFConfig().GetAppExternalWebsocketURL(),
		"NoCleanup":        group.ProviderNoCleanup(),
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getProviderMaxPort(startPort, deviceCount int32) int32 {
	return (deviceCount * 4) + startPort
}
