Android Farm Operator CRD Reference
-----------------------------------

### Packages:

-   [android.stf.io/v1alpha1](#android.stf.io%2fv1alpha1)

Types

-   [ADBConfig](#%23android.stf.io%2fv1alpha1.ADBConfig)
-   [APIConfig](#%23android.stf.io%2fv1alpha1.APIConfig)
-   [AndroidDevice](#%23android.stf.io%2fv1alpha1.AndroidDevice)
-   [AndroidDeviceConfig](#%23android.stf.io%2fv1alpha1.AndroidDeviceConfig)
-   [AndroidDeviceConfigSpec](#%23android.stf.io%2fv1alpha1.AndroidDeviceConfigSpec)
-   [AndroidDeviceSpec](#%23android.stf.io%2fv1alpha1.AndroidDeviceSpec)
-   [AndroidFarm](#%23android.stf.io%2fv1alpha1.AndroidFarm)
-   [AndroidFarmSpec](#%23android.stf.io%2fv1alpha1.AndroidFarmSpec)
-   [AppConfig](#%23android.stf.io%2fv1alpha1.AppConfig)
-   [AuthConfig](#%23android.stf.io%2fv1alpha1.AuthConfig)
-   [DeviceGroup](#%23android.stf.io%2fv1alpha1.DeviceGroup)
-   [DeviceManagementPolicy](#%23android.stf.io%2fv1alpha1.DeviceManagementPolicy)
-   [EmulatorConfig](#%23android.stf.io%2fv1alpha1.EmulatorConfig)
-   [GlobalProviderConfig](#%23android.stf.io%2fv1alpha1.GlobalProviderConfig)
-   [HostUSBConfig](#%23android.stf.io%2fv1alpha1.HostUSBConfig)
-   [PodManagementPolicy](#%23android.stf.io%2fv1alpha1.PodManagementPolicy)
-   [ProcessorConfig](#%23android.stf.io%2fv1alpha1.ProcessorConfig)
-   [ProviderConfig](#%23android.stf.io%2fv1alpha1.ProviderConfig)
-   [ReaperConfig](#%23android.stf.io%2fv1alpha1.ReaperConfig)
-   [RethinkDBConfig](#%23android.stf.io%2fv1alpha1.RethinkDBConfig)
-   [RethinkDBProxyConfig](#%23android.stf.io%2fv1alpha1.RethinkDBProxyConfig)
-   [STFConfig](#%23android.stf.io%2fv1alpha1.STFConfig)
-   [STFImage](#%23android.stf.io%2fv1alpha1.STFImage)
-   [STFOAuth](#%23android.stf.io%2fv1alpha1.STFOAuth)
-   [StorageConfig](#%23android.stf.io%2fv1alpha1.StorageConfig)
-   [TCPRedirConfig](#%23android.stf.io%2fv1alpha1.TCPRedirConfig)
-   [TLSConfig](#%23android.stf.io%2fv1alpha1.TLSConfig)
-   [TraefikConfig](#%23android.stf.io%2fv1alpha1.TraefikConfig)
-   [TraefikDashboard](#%23android.stf.io%2fv1alpha1.TraefikDashboard)
-   [TraefikDeployment](#%23android.stf.io%2fv1alpha1.TraefikDeployment)
-   [TriproxyAppConfig](#%23android.stf.io%2fv1alpha1.TriproxyAppConfig)
-   [TriproxyDevConfig](#%23android.stf.io%2fv1alpha1.TriproxyDevConfig)
-   [Volume](#%23android.stf.io%2fv1alpha1.Volume)
-   [WebsocketConfig](#%23android.stf.io%2fv1alpha1.WebsocketConfig)

android.stf.io/v1alpha1
-----------------------

Package v1alpha1 contains API Schema definitions for the android
v1alpha1 API group

Resource Types:

### ADBConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

ADBConfig represents configuration options for the adb containers

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>image</code> <em>string</em></td>
<td><p>Image for the adb servers. Defaults to <code>quay.io/tinyzimmer/adbmon</code>. Source in this repository.</p></td>
</tr>
<tr class="even">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to attach to deployments using this image.</p></td>
</tr>
<tr class="odd">
<td><code>imagePullSecrets</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">[]Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Any pull secrets required for downloading the image.</p></td>
</tr>
<tr class="even">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the provider adb sidecars.</p></td>
</tr>
</tbody>
</table>

### APIConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

APIConfig represents configuration options for the api servers

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the stf api servers</p></td>
</tr>
<tr class="even">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of api server replicas to run</p></td>
</tr>
</tbody>
</table>

### AndroidDevice

AndroidDevice is the Schema for the androiddevices API

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta">Kubernetes meta/v1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceSpec">AndroidDeviceSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>configRef</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>A reference to an AndroidDeviceConfig to use for the emulators in this group.</p></td>
</tr>
<tr class="even">
<td><code>deviceConfig</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceConfigSpec">AndroidDeviceConfigSpec</a></em></td>
<td><p>Any overrides to the config represented by the ConfigRef. Any values supplied here will be merged into the found AndroidDeviceConfig, with fields in this object taking precedence over existing ones in the AndroidDeviceConfig.</p></td>
</tr>
<tr class="odd">
<td><code>hostname</code> <em>string</em></td>
<td><p>A hostname to apply to the device (used by AndroidFarm controller)</p></td>
</tr>
<tr class="even">
<td><code>subdomain</code> <em>string</em></td>
<td><p>A subdomain to apply to the device (used by AndroidFarm controller)</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceStatus">AndroidDeviceStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### AndroidDeviceConfig

AndroidDeviceConfig is the Schema for the androiddeviceconfigs API

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta">Kubernetes meta/v1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceConfigSpec">AndroidDeviceConfigSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>dockerImage</code> <em>string</em></td>
<td><p>The docker image to use for emulator devices</p></td>
</tr>
<tr class="even">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to use for emulator pods</p></td>
</tr>
<tr class="odd">
<td><code>imagePullSecrets</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">[]Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Pull secrets required for the docker image.</p></td>
</tr>
<tr class="even">
<td><code>adbPort</code> <em>int32</em></td>
<td><p>The ADB port that the emulator listens on. Defaults to 5555. A sidecar will be spawned within emulator pods that redirects external traffic to this port.</p></td>
</tr>
<tr class="odd">
<td><code>command</code> <em>[]string</em></td>
<td><p>An optional command to run when starting an emulator image</p></td>
</tr>
<tr class="even">
<td><code>args</code> <em>[]string</em></td>
<td><p>Any arguments to pass to the above command.</p></td>
</tr>
<tr class="odd">
<td><code>extraPorts</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#containerport-v1-core">[]Kubernetes core/v1.ContainerPort</a></em></td>
<td><p>Extra port mappings to apply to the emulator pods.</p></td>
</tr>
<tr class="even">
<td><code>extraEnvVars</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvar-v1-core">[]Kubernetes core/v1.EnvVar</a></em></td>
<td><p>Extra environment variables to supply to the emulator pods.</p></td>
</tr>
<tr class="odd">
<td><code>kvmEnabled</code> <em>bool</em></td>
<td><p>Whether to mount the kvm device to the pods, will require that the operator can launch privileged pods.</p></td>
</tr>
<tr class="even">
<td><code>volumes</code> <em><a href="#android.stf.io/v1alpha1.Volume">[]Volume</a></em></td>
<td><p>A list of volume configurations to apply to the emulator pods.</p></td>
</tr>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints to place on the emulators.</p></td>
</tr>
<tr class="even">
<td><code>startupJobTemplates</code> <em>[]string</em></td>
<td><p>A list of AndroidJobTemplates to execute against new instances. TODO: Very very very beta</p></td>
</tr>
<tr class="odd">
<td><code>tcpRedir</code> <em><a href="#android.stf.io/v1alpha1.TCPRedirConfig">TCPRedirConfig</a></em></td>
<td><p>Configuration for the tcp redirection side car</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceConfigStatus">AndroidDeviceConfigStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### AndroidDeviceConfigSpec

(*Appears on:*
[AndroidDeviceConfig](#android.stf.io/v1alpha1.AndroidDeviceConfig),
[AndroidDeviceSpec](#android.stf.io/v1alpha1.AndroidDeviceSpec),
[EmulatorConfig](#android.stf.io/v1alpha1.EmulatorConfig))

AndroidDeviceConfigSpec defines the desired state of AndroidDeviceConfig

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>dockerImage</code> <em>string</em></td>
<td><p>The docker image to use for emulator devices</p></td>
</tr>
<tr class="even">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to use for emulator pods</p></td>
</tr>
<tr class="odd">
<td><code>imagePullSecrets</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">[]Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Pull secrets required for the docker image.</p></td>
</tr>
<tr class="even">
<td><code>adbPort</code> <em>int32</em></td>
<td><p>The ADB port that the emulator listens on. Defaults to 5555. A sidecar will be spawned within emulator pods that redirects external traffic to this port.</p></td>
</tr>
<tr class="odd">
<td><code>command</code> <em>[]string</em></td>
<td><p>An optional command to run when starting an emulator image</p></td>
</tr>
<tr class="even">
<td><code>args</code> <em>[]string</em></td>
<td><p>Any arguments to pass to the above command.</p></td>
</tr>
<tr class="odd">
<td><code>extraPorts</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#containerport-v1-core">[]Kubernetes core/v1.ContainerPort</a></em></td>
<td><p>Extra port mappings to apply to the emulator pods.</p></td>
</tr>
<tr class="even">
<td><code>extraEnvVars</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvar-v1-core">[]Kubernetes core/v1.EnvVar</a></em></td>
<td><p>Extra environment variables to supply to the emulator pods.</p></td>
</tr>
<tr class="odd">
<td><code>kvmEnabled</code> <em>bool</em></td>
<td><p>Whether to mount the kvm device to the pods, will require that the operator can launch privileged pods.</p></td>
</tr>
<tr class="even">
<td><code>volumes</code> <em><a href="#android.stf.io/v1alpha1.Volume">[]Volume</a></em></td>
<td><p>A list of volume configurations to apply to the emulator pods.</p></td>
</tr>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints to place on the emulators.</p></td>
</tr>
<tr class="even">
<td><code>startupJobTemplates</code> <em>[]string</em></td>
<td><p>A list of AndroidJobTemplates to execute against new instances. TODO: Very very very beta</p></td>
</tr>
<tr class="odd">
<td><code>tcpRedir</code> <em><a href="#android.stf.io/v1alpha1.TCPRedirConfig">TCPRedirConfig</a></em></td>
<td><p>Configuration for the tcp redirection side car</p></td>
</tr>
</tbody>
</table>

### AndroidDeviceSpec

(*Appears on:* [AndroidDevice](#android.stf.io/v1alpha1.AndroidDevice))

AndroidDeviceSpec defines the desired state of AndroidDevice

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>configRef</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>A reference to an AndroidDeviceConfig to use for the emulators in this group.</p></td>
</tr>
<tr class="even">
<td><code>deviceConfig</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceConfigSpec">AndroidDeviceConfigSpec</a></em></td>
<td><p>Any overrides to the config represented by the ConfigRef. Any values supplied here will be merged into the found AndroidDeviceConfig, with fields in this object taking precedence over existing ones in the AndroidDeviceConfig.</p></td>
</tr>
<tr class="odd">
<td><code>hostname</code> <em>string</em></td>
<td><p>A hostname to apply to the device (used by AndroidFarm controller)</p></td>
</tr>
<tr class="even">
<td><code>subdomain</code> <em>string</em></td>
<td><p>A subdomain to apply to the device (used by AndroidFarm controller)</p></td>
</tr>
</tbody>
</table>

### AndroidFarm

AndroidFarm is the Schema for the androidfarms API

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta">Kubernetes meta/v1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#android.stf.io/v1alpha1.AndroidFarmSpec">AndroidFarmSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>deviceGroups</code> <em><a href="#android.stf.io/v1alpha1.*github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1.DeviceGroup">[]*github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1.DeviceGroup</a></em></td>
<td><p>A list of device groups and their configurations to run on the cluster</p></td>
</tr>
<tr class="even">
<td><code>deviceManagementPolicy</code> <em><a href="#android.stf.io/v1alpha1.DeviceManagementPolicy">DeviceManagementPolicy</a></em></td>
<td><p>A device management policy to apply globally unless overridden on the group level</p></td>
</tr>
<tr class="odd">
<td><code>stfConfig</code> <em><a href="#android.stf.io/v1alpha1.STFConfig">STFConfig</a></em></td>
<td><p>The configuration for the OpenSTF Deployment</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#android.stf.io/v1alpha1.AndroidFarmStatus">AndroidFarmStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### AndroidFarmSpec

(*Appears on:* [AndroidFarm](#android.stf.io/v1alpha1.AndroidFarm))

AndroidFarmSpec defines the desired state of AndroidFarm

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>deviceGroups</code> <em><a href="#android.stf.io/v1alpha1.*github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1.DeviceGroup">[]*github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1.DeviceGroup</a></em></td>
<td><p>A list of device groups and their configurations to run on the cluster</p></td>
</tr>
<tr class="even">
<td><code>deviceManagementPolicy</code> <em><a href="#android.stf.io/v1alpha1.DeviceManagementPolicy">DeviceManagementPolicy</a></em></td>
<td><p>A device management policy to apply globally unless overridden on the group level</p></td>
</tr>
<tr class="odd">
<td><code>stfConfig</code> <em><a href="#android.stf.io/v1alpha1.STFConfig">STFConfig</a></em></td>
<td><p>The configuration for the OpenSTF Deployment</p></td>
</tr>
</tbody>
</table>

### AppConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

AppConfig represents configuration options for the app deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
<tr class="even">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of app replicas to run</p></td>
</tr>
</tbody>
</table>

### AuthConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

STFAuth represents the authentication configuration for OpenSTF.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>mock</code> <em>bool</em></td>
<td><p>Use the stf mock authentication adapter.</p></td>
</tr>
<tr class="even">
<td><code>oauth</code> <em><a href="#android.stf.io/v1alpha1.STFOAuth">STFOAuth</a></em></td>
<td><p>Use OAuth with the provided parameters for authentication.</p></td>
</tr>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Auth deployment resource requirements.</p></td>
</tr>
</tbody>
</table>

### DeviceGroup

DeviceGroup represents a collection of android devices that share a
common configuration.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>name</code> <em>string</em></td>
<td><p>A name for the device group, this field is required. The name is prefixed to all resources created for the group.</p></td>
</tr>
<tr class="even">
<td><code>provider</code> <em><a href="#android.stf.io/v1alpha1.ProviderConfig">ProviderConfig</a></em></td>
<td><p>STF provider configurations for this device group.</p></td>
</tr>
<tr class="odd">
<td><code>emulators</code> <em><a href="#android.stf.io/v1alpha1.EmulatorConfig">EmulatorConfig</a></em></td>
<td><p>A configuration for emulated devices running in pods on the kubernetes cluster</p></td>
</tr>
<tr class="even">
<td><code>hostUSB</code> <em><a href="#android.stf.io/v1alpha1.HostUSBConfig">HostUSBConfig</a></em></td>
<td><p>A configuration for connecting host usb devices to the AndroidFarm.</p></td>
</tr>
<tr class="odd">
<td><code>omitFromSTF</code> <em>bool</em></td>
<td><p>TODO: implement</p></td>
</tr>
</tbody>
</table>

### DeviceManagementPolicy

(*Appears on:*
[AndroidFarmSpec](#android.stf.io/v1alpha1.AndroidFarmSpec),
[EmulatorConfig](#android.stf.io/v1alpha1.EmulatorConfig))

DeviceManagementPolicy represents a policy for managing concurrency
during the creation and updating of emulator pods.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>podManagementPolicy</code> <em><a href="#android.stf.io/v1alpha1.PodManagementPolicy">PodManagementPolicy</a></em></td>
<td><p>The type of policy to enforce, currently only OrderedReady.</p></td>
</tr>
<tr class="even">
<td><code>concurrency</code> <em>int32</em></td>
<td><p>The maximum number of devices that can be booting at any point in time.</p></td>
</tr>
</tbody>
</table>

### EmulatorConfig

(*Appears on:* [DeviceGroup](#android.stf.io/v1alpha1.DeviceGroup))

EmulatorConfig is a configuration for virtual emulators running in pods
on the kubernetes cluster.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>namespace</code> <em>string</em></td>
<td><p>The namespace to run the device group, defaults to the default namespace.</p></td>
</tr>
<tr class="even">
<td><code>count</code> <em>int32</em></td>
<td><p>The number of devices to run in the group. Defaults to no devices.</p></td>
</tr>
<tr class="odd">
<td><code>hostnameTemplate</code> <em>string</em></td>
<td><p>A go-template to use for configuring the hostname of the devices. Currently only {{ .Index }} is passed to thte template, but more will come. A headless service is put in front of device groups to make the individual pods accessible by their hostname/subdomain</p></td>
</tr>
<tr class="even">
<td><code>subdomain</code> <em>string</em></td>
<td><p>A subdomain to use for the pods in the device group. This also becomes the name of the headless service.</p></td>
</tr>
<tr class="odd">
<td><code>deviceManagementPolicy</code> <em><a href="#android.stf.io/v1alpha1.DeviceManagementPolicy">DeviceManagementPolicy</a></em></td>
<td><p>A policy for managing concurrency during provisioning/updates of android emulators.</p></td>
</tr>
<tr class="even">
<td><code>configRef</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>A reference to an AndroidDeviceConfig to use for the emulators in this group.</p></td>
</tr>
<tr class="odd">
<td><code>deviceConfig</code> <em><a href="#android.stf.io/v1alpha1.AndroidDeviceConfigSpec">AndroidDeviceConfigSpec</a></em></td>
<td><p>Any overrides to the config represented by the ConfigRef. Any values supplied here will be merged into the found AndroidDeviceConfig, with fields in this object taking precedence over existing ones in the AndroidDeviceConfig.</p></td>
</tr>
</tbody>
</table>

### GlobalProviderConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

GlobalProviderConfig represents global configuration options for the
provider deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

### HostUSBConfig

(*Appears on:* [DeviceGroup](#android.stf.io/v1alpha1.DeviceGroup))

HostUSBConfig is a configuration for connecting devices attached
physically to the kubernetes hosts.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>nodeName</code> <em>string</em></td>
<td><p>The node to launch an ADB server on for binding devices to STF.</p></td>
</tr>
<tr class="even">
<td><code>maxDevices</code> <em>int32</em></td>
<td><p>Specify the maximum number of devices expected to run on this host. This is required because for lack of a better solution we dynamically allocate provider service ports at the moment and we need to determine how many to do for a usb farm.</p></td>
</tr>
</tbody>
</table>

PodManagementPolicy (`string` alias)

(*Appears on:*
[DeviceManagementPolicy](#android.stf.io/v1alpha1.DeviceManagementPolicy))

### ProcessorConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

ProcessorConfig represents configuration options for the processor
deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

### ProviderConfig

(*Appears on:* [DeviceGroup](#android.stf.io/v1alpha1.DeviceGroup))

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>startPort</code> <em>int32</em></td>
<td><p>The starting port to use for provider services. Defaults to 15000. If specifying multiple device groups, you should set this explicitly for each group and ensure they are not able to overlap. (~4 ports per device)</p></td>
</tr>
<tr class="even">
<td><code>clusterLocalADB</code> <em>bool</em></td>
<td><p>When set to true, the provider will advertise it’s cluster local service address for adb connections. The default behavior is to advertise the external app hostname.</p></td>
</tr>
<tr class="odd">
<td><code>persistDeviceState</code> <em>bool</em></td>
<td><p>Set to true to persist device state (apps, accounts, caches) between user sessions.</p></td>
</tr>
</tbody>
</table>

### ReaperConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

ReaperConfig represents configuration options for the reaper deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

### RethinkDBConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

RethinkDBConfig represents configurations for the RethinkDB StatefulSet.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>version</code> <em>string</em></td>
<td><p>The version of rethinkdb to run. Defaults to 2.4.</p></td>
</tr>
<tr class="even">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to use for the rethinkdb image.</p></td>
</tr>
<tr class="odd">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of rethinkdb replicas per shard to make.</p></td>
</tr>
<tr class="even">
<td><code>shards</code> <em>int32</em></td>
<td><p>The number of shards to use for each table in the stf database.</p></td>
</tr>
<tr class="odd">
<td><code>pvcSpec</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaimspec-v1-core">Kubernetes core/v1.PersistentVolumeClaimSpec</a></em></td>
<td><p>A PVCSpec to use for RethinkDB persistence.</p></td>
</tr>
<tr class="even">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints for the rethinkdb instances</p></td>
</tr>
</tbody>
</table>

### RethinkDBProxyConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

RethinkDBProxyConfig represents configuration options for the rethinkdb
proxy deployment.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of proxy instances to run</p></td>
</tr>
<tr class="even">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints for the proxy instances</p></td>
</tr>
</tbody>
</table>

### STFConfig

(*Appears on:*
[AndroidFarmSpec](#android.stf.io/v1alpha1.AndroidFarmSpec))

STFConfig represents configuration options for the OpenSTF deployment in
this AndroidFarm.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>appHostname</code> <em>string</em></td>
<td><p>The external hostname to use when configuring OpenSTF services. The OpenSTF deployment must be accessible at this address (or IP).</p></td>
</tr>
<tr class="even">
<td><code>secret</code> <em>string</em></td>
<td><p>The name of the kubernetes secret containing secrets for configuring OpenSTF.</p></td>
</tr>
<tr class="odd">
<td><code>stfSecretKey</code> <em>string</em></td>
<td><p>The key in the above secret where the OpenSTF secret is. Defaults to ‘stf-secret’.</p></td>
</tr>
<tr class="even">
<td><code>serviceAccount</code> <em>string</em></td>
<td><p>A kubernetes service account to attach to OpenSTF deployments. This can be required if you are launching privileged pods that need to be validated against a PodSecurityPolicy.</p></td>
</tr>
<tr class="odd">
<td><code>privilegedDeployments</code> <em>bool</em></td>
<td><p>Use privileged security contexts for OpenSTF deployments. This is required if you are using an image that runs as root.</p></td>
</tr>
<tr class="even">
<td><code>namespace</code> <em>string</em></td>
<td><p>The namespace to provision the STF deployments in. Defaults to the default namespace.</p></td>
</tr>
<tr class="odd">
<td><code>stfImage</code> <em><a href="#android.stf.io/v1alpha1.STFImage">STFImage</a></em></td>
<td><p>The docker image configuration to use for the STF services.</p></td>
</tr>
<tr class="even">
<td><code>adb</code> <em><a href="#android.stf.io/v1alpha1.ADBConfig">ADBConfig</a></em></td>
<td><p>ADB extra configuration options</p></td>
</tr>
<tr class="odd">
<td><code>api</code> <em><a href="#android.stf.io/v1alpha1.APIConfig">APIConfig</a></em></td>
<td><p>API extra configuration options</p></td>
</tr>
<tr class="even">
<td><code>app</code> <em><a href="#android.stf.io/v1alpha1.AppConfig">AppConfig</a></em></td>
<td><p>App extra configuration options</p></td>
</tr>
<tr class="odd">
<td><code>auth</code> <em><a href="#android.stf.io/v1alpha1.AuthConfig">AuthConfig</a></em></td>
<td><p>Authentication configuration options</p></td>
</tr>
<tr class="even">
<td><code>processor</code> <em><a href="#android.stf.io/v1alpha1.ProcessorConfig">ProcessorConfig</a></em></td>
<td><p>Processor configuration options</p></td>
</tr>
<tr class="odd">
<td><code>provider</code> <em><a href="#android.stf.io/v1alpha1.GlobalProviderConfig">GlobalProviderConfig</a></em></td>
<td><p>Provider configuration options</p></td>
</tr>
<tr class="even">
<td><code>reaper</code> <em><a href="#android.stf.io/v1alpha1.ReaperConfig">ReaperConfig</a></em></td>
<td><p>Reaper configuration options</p></td>
</tr>
<tr class="odd">
<td><code>triproxyApp</code> <em><a href="#android.stf.io/v1alpha1.TriproxyAppConfig">TriproxyAppConfig</a></em></td>
<td><p>Triproxy App configuration options</p></td>
</tr>
<tr class="even">
<td><code>triproxyDev</code> <em><a href="#android.stf.io/v1alpha1.TriproxyDevConfig">TriproxyDevConfig</a></em></td>
<td><p>Triproxy Dev configuration options</p></td>
</tr>
<tr class="odd">
<td><code>websocket</code> <em><a href="#android.stf.io/v1alpha1.WebsocketConfig">WebsocketConfig</a></em></td>
<td><p>Websocket configuration options</p></td>
</tr>
<tr class="even">
<td><code>traefik</code> <em><a href="#android.stf.io/v1alpha1.TraefikConfig">TraefikConfig</a></em></td>
<td><p>A configuration for the traefik deployment/routes put in front of the STF deployments.</p></td>
</tr>
<tr class="odd">
<td><code>rethinkdb</code> <em><a href="#android.stf.io/v1alpha1.RethinkDBConfig">RethinkDBConfig</a></em></td>
<td><p>A configuration for the RethinKDB statefulset.</p></td>
</tr>
<tr class="even">
<td><code>rethinkdbProxy</code> <em><a href="#android.stf.io/v1alpha1.RethinkDBProxyConfig">RethinkDBProxyConfig</a></em></td>
<td><p>A configuration for the rethinkdb proxy deployment</p></td>
</tr>
<tr class="odd">
<td><code>storage</code> <em><a href="#android.stf.io/v1alpha1.StorageConfig">StorageConfig</a></em></td>
<td><p>A configuration for the OpenSTF storage services.</p></td>
</tr>
</tbody>
</table>

### STFImage

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

STFImage is the configuration for the docker image used in STF
deployments.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>image</code> <em>string</em></td>
<td><p>Image is the repository to download the image from. Defaults to openstf/stf:latest.</p></td>
</tr>
<tr class="even">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to attach to deployments using this image.</p></td>
</tr>
<tr class="odd">
<td><code>imagePullSecrets</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">[]Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Any pull secrets required for downloading the image.</p></td>
</tr>
</tbody>
</table>

### STFOAuth

(*Appears on:* [AuthConfig](#android.stf.io/v1alpha1.AuthConfig))

STFOauth represents an OAuth configuration to use for the STF oauth
adapter.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>authorizationURL</code> <em>string</em></td>
<td><p>The Authorization URL for the OAuth service</p></td>
</tr>
<tr class="even">
<td><code>tokenURL</code> <em>string</em></td>
<td><p>The Token URL for the OAuth service.</p></td>
</tr>
<tr class="odd">
<td><code>userInfoURL</code> <em>string</em></td>
<td><p>The User Info URL for the OAuth service.</p></td>
</tr>
<tr class="even">
<td><code>scopes</code> <em>[]string</em></td>
<td><p>The scopes needed for OAuth.</p></td>
</tr>
<tr class="odd">
<td><code>callbackURL</code> <em>string</em></td>
<td><p>The OAuth callback URL. TODO : This doesn’t need to be required and can default to: http(s):///auth/oauth/callback</p></td>
</tr>
<tr class="even">
<td><code>clientIDKey</code> <em>string</em></td>
<td><p>The key in the STF secret that contains the client id. Defaults to ‘client-id’.</p></td>
</tr>
<tr class="odd">
<td><code>clientSecretKey</code> <em>string</em></td>
<td><p>The key in the STF secret that contains the client secret key. Defaults to ‘client-secret’.</p></td>
</tr>
</tbody>
</table>

### StorageConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

STFStorageConfig represents configurations for the OpenSTF storage
service.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of stf-storage replicas to run</p></td>
</tr>
<tr class="even">
<td><code>pvcSpec</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaimspec-v1-core">Kubernetes core/v1.PersistentVolumeClaimSpec</a></em></td>
<td><p>A PVC spec to use for storage persistence. If specifying more than one replica, only one PVC will be created and it should allow ReadWriteMany.</p></td>
</tr>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Storage deployments resource requirements.</p></td>
</tr>
</tbody>
</table>

### TCPRedirConfig

(*Appears on:*
[AndroidDeviceConfigSpec](#android.stf.io/v1alpha1.AndroidDeviceConfigSpec))

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>enabled</code> <em>bool</em></td>
<td><p>Whether to run a sidecar with emulator pods that redirects TCP traffic on the adb port to the emulator adb server listening on the loopback interface. This is required for the image used in this repository, but if you are using an image that exposes ADB on all interfaces itself, this is not required.</p></td>
</tr>
<tr class="even">
<td><code>image</code> <em>string</em></td>
<td><p>Image is the repository to download the image from. Defaults to quay.io/tinyzimmer/goredir whose source is in this repository.</p></td>
</tr>
<tr class="odd">
<td><code>imagePullPolicy</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">Kubernetes core/v1.PullPolicy</a></em></td>
<td><p>The pull policy to attach to deployments using this image.</p></td>
</tr>
<tr class="even">
<td><code>imagePullSecrets</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">[]Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Any pull secrets required for downloading the image.</p></td>
</tr>
</tbody>
</table>

### TLSConfig

(*Appears on:* [TraefikConfig](#android.stf.io/v1alpha1.TraefikConfig))

SSLConfig represents the SSL configuration for the STF deployments.
Specify an empty object to configure SSL with traefik’s default
self-signed certificate. Should only be done for testing.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>tlsSecret</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#localobjectreference-v1-core">Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>Specifies a preexisting TLS secret in the cluster to use. It must follow the standard format with a tls.crt and tls.key.</p></td>
</tr>
<tr class="even">
<td><code>issuerRef</code> <em><a href="https://godoc.org/github.com/jetstack/cert-manager/pkg/apis/meta/v1#ObjectReference">github.com/jetstack/cert-manager/pkg/apis/meta/v1.ObjectReference</a></em></td>
<td><p>(WIP) A cert-manager issuer reference to use to provision a TLS secret.</p></td>
</tr>
<tr class="odd">
<td><code>external</code> <em>bool</em></td>
<td><p>Specify that SSL is managed externally. OpenSTF will be configured to know it is being served over HTTPS, but you will be responsible for terminating TLS before sending traffic to the traefik instance.</p></td>
</tr>
</tbody>
</table>

### TraefikConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

TraefikConfig represents configurations for the traefik deployment
placed in front of the OpenSTF services.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>useIngressRoute</code> <em>bool</em></td>
<td><p>(WIP) - use IngressRoute CRDs for an existing traefik deployment instead of creating a standalone traefik service.</p></td>
</tr>
<tr class="even">
<td><code>deployment</code> <em><a href="#android.stf.io/v1alpha1.TraefikDeployment">TraefikDeployment</a></em></td>
<td><p>Configuration options for the traefik deployment.</p></td>
</tr>
<tr class="odd">
<td><code>tls</code> <em><a href="#android.stf.io/v1alpha1.TLSConfig">TLSConfig</a></em></td>
<td><p>TLS configurations for traefik</p></td>
</tr>
</tbody>
</table>

### TraefikDashboard

(*Appears on:*
[TraefikDeployment](#android.stf.io/v1alpha1.TraefikDeployment))

TraefikDashboard represents configuration options for the Traefik
dashboard.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>host</code> <em>string</em></td>
<td><p>The hostname that should route to the traefik dashboard.</p></td>
</tr>
<tr class="even">
<td><code>ipWhitelist</code> <em>[]string</em></td>
<td><p>A list of IP addresses to whitelist for dashboard access.</p></td>
</tr>
</tbody>
</table>

### TraefikDeployment

(*Appears on:* [TraefikConfig](#android.stf.io/v1alpha1.TraefikConfig))

TraefikDeployment represents configuration options for the traefik
deployment.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>replicas</code> <em>int32</em></td>
<td><p>The number of traefik instances to run.</p></td>
</tr>
<tr class="even">
<td><code>version</code> <em>string</em></td>
<td><p>The version of traefik to run, only &gt;2.0 supported. Defaults to 2.2.0.</p></td>
</tr>
<tr class="odd">
<td><code>serviceType</code> <em>string</em></td>
<td><p>The type of service to create for Traefik. Defaults to LoadBalancer. If using external SSL from a pre-existing ingress controller, you’ll want to set this to ClusterIP.</p></td>
</tr>
<tr class="even">
<td><code>accessLogs</code> <em>bool</em></td>
<td><p>Set to true if you wish for traefik to produce access logs</p></td>
</tr>
<tr class="odd">
<td><code>dashboard</code> <em><a href="#android.stf.io/v1alpha1.TraefikDashboard">TraefikDashboard</a></em></td>
<td><p>A configuration for the traefik dashboard</p></td>
</tr>
<tr class="even">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints for the traefik deployment</p></td>
</tr>
</tbody>
</table>

### TriproxyAppConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

TriproxyAppConfig represents configuration options for the triproxy app
deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

### TriproxyDevConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

TriproxyDevConfig represents configuration options for the triproxy dev
deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

### Volume

(*Appears on:*
[AndroidDeviceConfigSpec](#android.stf.io/v1alpha1.AndroidDeviceConfigSpec))

Volume represents a volume configuration for the emulator.

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>volumePrefix</code> <em>string</em></td>
<td><p>A prefix to apply to PVCs created for devices using this configuration.</p></td>
</tr>
<tr class="even">
<td><code>mountPoint</code> <em>string</em></td>
<td><p>Where to mount the volume in emulator pods.</p></td>
</tr>
<tr class="odd">
<td><code>pvcSpec</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaimspec-v1-core">Kubernetes core/v1.PersistentVolumeClaimSpec</a></em></td>
<td><p>A PVC spec to use for creating the emulator volumes.</p></td>
</tr>
</tbody>
</table>

### WebsocketConfig

(*Appears on:* [STFConfig](#android.stf.io/v1alpha1.STFConfig))

WebsocketConfig represents configuration options for the websocket
deployments

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>resources</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>The resource restraints for the app deployment</p></td>
</tr>
</tbody>
</table>

------------------------------------------------------------------------

*Generated with `gen-crd-api-reference-docs` on git commit `c445a32`.*
