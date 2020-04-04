package v1alpha1

import (
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodManagementPolicy string

const (
	OrderedReady PodManagementPolicy = "OrderedReady"
)

// AndroidFarmSpec defines the desired state of AndroidFarm
type AndroidFarmSpec struct {
	// A list of device groups and their configurations to run on the cluster
	DeviceGroups []*DeviceGroup `json:"deviceGroups"`
	// A device management policy to apply globally unless overridden
	// on the group level
	DeviceManagementPolicy *DeviceManagementPolicy `json:"deviceManagementPolicy,omitempty"`
	// The configuration for the OpenSTF Deployment
	STFConfig *STFConfig `json:"stfConfig,omitempty"`
}

// DeviceGroup represents a collection of android devices that share a common
// configuration.
type DeviceGroup struct {
	// A name for the device group, this field is required. The name is prefixed
	// to all resources created for the group.
	Name string `json:"name"`
	// STF provider configurations for this device group.
	Provider *ProviderConfig `json:"provider,omitempty"`
	// A configuration for emulated devices running in pods on the kubernetes cluster
	Emulators *EmulatorConfig `json:"emulators,omitempty"`
	// A configuration for connecting host usb devices to the AndroidFarm.
	HostUSB *HostUSBConfig `json:"hostUSB,omitempty"`
	// TODO: implement
	OmitFromSTF bool `json:"omitFromSTF,omitempty"`
}

type ProviderConfig struct {
	// The starting port to use for provider services. Defaults to 15000. If specifying
	// multiple device groups, you should set this explicitly for each group and
	// ensure they are not able to overlap. (~4 ports per device)
	StartPort int32 `json:"startPort,omitempty"`
	// When set to true, the provider will advertise it's cluster local service
	// address for adb connections. The default behavior is to advertise the external
	// app hostname.
	ClusterLocalADB bool `json:"clusterLocalADB,omitempty"`
	// Set to true to persist device state (apps, accounts, caches) between user
	// sessions.
	PersistDeviceState bool `json:"persistDeviceState,omitempty"`
}

// EmulatorConfig is a configuration for virtual emulators running in pods on
// the kubernetes cluster.
type EmulatorConfig struct {
	// The namespace to run the device group, defaults to the default namespace.
	Namespace string `json:"namespace,omitempty"`
	// The number of devices to run in the group. Defaults to no devices.
	Count int32 `json:"count,omitempty"`
	// A go-template to use for configuring the hostname of the devices. Currently
	// only {{ .Index }} is passed to thte template, but more will come. A headless
	// service is put in front of device groups to make the individual pods accessible
	// by their hostname/subdomain
	HostnameTemplate string `json:"hostnameTemplate,omitempty"`
	// A subdomain to use for the pods in the device group. This also becomes the
	// name of the headless service.
	Subdomain string `json:"subdomain,omitempty"`
	// A policy for managing concurrency during provisioning/updates of android
	// emulators.
	DeviceManagementPolicy *DeviceManagementPolicy `json:"deviceManagementPolicy,omitempty"`
	// A reference to an AndroidDeviceConfig to use for the emulators in this group.
	ConfigRef *corev1.LocalObjectReference `json:"configRef,omitempty"`
	// Any overrides to the config represented by the ConfigRef. Any values supplied here
	// will be merged into the found AndroidDeviceConfig, with fields in this object
	// taking precedence over existing ones in the AndroidDeviceConfig.
	DeviceConfig *AndroidDeviceConfigSpec `json:"deviceConfig,omitempty"`
}

// HostUSBConfig is a configuration for connecting devices attached physically
// to the kubernetes hosts.
type HostUSBConfig struct {
	// The node to launch an ADB server on for binding devices to STF.
	NodeName string `json:"nodeName,omitempty"`
	// Specify the maximum number of devices expected to run on this host. This
	// is required because for lack of a better solution we dynamically allocate provider
	// service ports at the moment and we need to determine how many to do for a usb farm.
	MaxDevices int32 `json:"maxDevices,omitempty"`
}

// DeviceManagementPolicy represents a policy for managing concurrency during
// the creation and updating of emulator pods.
type DeviceManagementPolicy struct {
	// The type of policy to enforce, currently only OrderedReady.
	PodManagementPolicy PodManagementPolicy `json:"podManagementPolicy,omitempty"`
	// The maximum number of devices that can be booting at any point in time.
	Concurrency int32 `json:"concurrency,omitempty"`
}

// STFConfig represents configuration options for the OpenSTF deployment in this
// AndroidFarm.
type STFConfig struct {
	// The external hostname to use when configuring OpenSTF services. The OpenSTF
	// deployment must be accessible at this address (or IP).
	AppHostname string `json:"appHostname"`
	// The name of the kubernetes secret containing secrets for configuring OpenSTF.
	Secret string `json:"secret"`
	// The key in the above secret where the OpenSTF secret is. Defaults to 'stf-secret'.
	STFSecretKey string `json:"stfSecretKey,omitempty"`
	// A kubernetes service account to attach to OpenSTF deployments. This can be
	// required if you are launching privileged pods that need to be validated against
	// a PodSecurityPolicy.
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// The namespace to provision the STF deployments in. Defaults to the default
	// namespace.
	Namespace string `json:"namespace,omitempty"`
	// The docker image configuration to use for the STF services.
	STFImage *STFImage `json:"stfImage,omitempty"`
	// ADB extra configuration options
	ADB *ADBConfig `json:"adb,omitempty"`
	// API extra configuration options
	API *APIConfig `json:"api,omitempty"`
	// App extra configuration options
	App *AppConfig `json:"app,omitempty"`
	// Authentication configuration options
	Auth *AuthConfig `json:"auth,omitempty"`
	// Processor configuration options
	Processor *ProcessorConfig `json:"processor,omitempty"`
	// Provider configuration options
	Provider *GlobalProviderConfig `json:"provider,omitempty"`
	// Reaper configuration options
	Reaper *ReaperConfig `json:"reaper,omitempty"`
	// Triproxy App configuration options
	TriproxyApp *TriproxyAppConfig `json:"triproxyApp,omitempty"`
	// Triproxy Dev configuration options
	TriproxyDev *TriproxyDevConfig `json:"triproxyDev,omitempty"`
	// Websocket configuration options
	Websocket *WebsocketConfig `json:"websocket,omitempty"`
	// A configuration for the traefik deployment/routes put in front of the
	// STF deployments.
	Traefik *TraefikConfig `json:"traefik,omitempty"`
	// A configuration for the RethinKDB statefulset.
	RethinkDB *RethinkDBConfig `json:"rethinkdb,omitempty"`
	// A configuration for the rethinkdb proxy deployment
	RethinkDBProxy *RethinkDBProxyConfig `json:"rethinkdbProxy,omitempty"`
	// A configuration for the OpenSTF storage services.
	Storage *StorageConfig `json:"storage,omitempty"`
}

// ADBConfig represents configuration options for the adb containers
type ADBConfig struct {
	// Image for the adb servers. Defaults to `quay.io/tinyzimmer/adbmon`. Source
	// in this repository.
	Image string `json:"image,omitempty"`
	// The pull policy to attach to deployments using this image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Any pull secrets required for downloading the image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// The resource restraints for the provider adb sidecars.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// APIConfig represents configuration options for the api servers
type APIConfig struct {
	// The resource restraints for the stf api servers
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// The number of api server replicas to run
	Replicas int32 `json:"replicas,omitempty"`
}

// AppConfig represents configuration options for the app deployments
type AppConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// The number of app replicas to run
	Replicas int32 `json:"replicas,omitempty"`
}

// ProcessorConfig represents configuration options for the processor deployments
type ProcessorConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// GlobalProviderConfig represents global configuration options for the provider deployments
type GlobalProviderConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// ReaperConfig represents configuration options for the reaper deployments
type ReaperConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// TriproxyAppConfig represents configuration options for the triproxy app deployments
type TriproxyAppConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// TriproxyDevConfig represents configuration options for the triproxy dev deployments
type TriproxyDevConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// WebsocketConfig represents configuration options for the websocket deployments
type WebsocketConfig struct {
	// The resource restraints for the app deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// STFImage is the configuration for the docker image used in STF deployments.
type STFImage struct {
	// Image is the repository to download the image from.
	// Defaults to openstf/stf:latest.
	Image string `json:"image,omitempty"`
	// The pull policy to attach to deployments using this image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Any pull secrets required for downloading the image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// RethinkDBConfig represents configurations for the RethinkDB StatefulSet.
type RethinkDBConfig struct {
	// The version of rethinkdb to run. Defaults to 2.4.
	Version string `json:"version,omitempty"`
	// The pull policy to use for the rethinkdb image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// The number of rethinkdb replicas per shard to make.
	Replicas int32 `json:"replicas,omitempty"`
	// The number of shards to use for each table in the stf database.
	Shards int32 `json:"shards,omitempty"`
	// A PVCSpec to use for RethinkDB persistence.
	PVCSpec *corev1.PersistentVolumeClaimSpec `json:"pvcSpec,omitempty"`
	// Resource restraints for the rethinkdb instances
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// RethinkDBProxyConfig represents configuration options for the rethinkdb proxy
// deployment.
type RethinkDBProxyConfig struct {
	// The number of proxy instances to run
	Replicas int32 `json:"replicas,omitempty"`
	// Resource restraints for the proxy instances
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// SSLConfig represents the SSL configuration for the STF deployments. Specify
// an empty object to configure SSL with traefik's default self-signed certificate.
// Should only be done for testing.
type TLSConfig struct {
	// Specifies a preexisting TLS secret in the cluster to use. It must follow
	// the standard format with a tls.crt and tls.key.
	TLSSecret *corev1.LocalObjectReference `json:"tlsSecret,omitempty"`
	// (WIP) A cert-manager issuer reference to use to provision a TLS secret.
	IssuerRef *cmmeta.ObjectReference `json:"issuerRef,omitempty"`
	// Specify that SSL is managed externally. OpenSTF will be configured to know
	// it is being served over HTTPS, but you will be responsible for terminating
	// TLS before sending traffic to the traefik instance. When using this option,
	// traefik will listen for requests on port 80, and you can set up an ingress to
	// `<farm_name>-stf-traefik`.
	External bool `json:"external,omitempty"`
}

// STFStorageConfig represents configurations for the OpenSTF storage service.
type StorageConfig struct {
	// The number of stf-storage replicas to run
	Replicas int32 `json:"replicas,omitempty"`
	// A PVC spec to use for storage persistence. If specifying more than one
	// replica, only one PVC will be created and it should allow ReadWriteMany.
	PVCSpec *corev1.PersistentVolumeClaimSpec `json:"pvcSpec,omitempty"`
	// Storage deployments resource requirements.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// STFAuth represents the authentication configuration for OpenSTF.
type AuthConfig struct {
	// Use the stf mock authentication adapter.
	Mock bool `json:"mock,omitempty"`
	// Use OAuth with the provided parameters for authentication.
	OAuth *STFOAuth `json:"oauth,omitempty"`
	// Auth deployment resource requirements.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// STFOauth represents an OAuth configuration to use for the STF oauth adapter.
type STFOAuth struct {
	// The Authorization URL for the OAuth service
	AuthorizationURL string `json:"authorizationURL"`
	// The Token URL for the OAuth service.
	TokenURL string `json:"tokenURL"`
	// The User Info URL for the OAuth service.
	UserInfoURL string `json:"userInfoURL"`
	// The scopes needed for OAuth.
	Scopes []string `json:"scopes"`
	// The OAuth callback URL.
	// TODO : This doesn't need to be required and can default to:
	//        http(s)://<app_hostname>/auth/oauth/callback
	CallbackURL string `json:"callbackURL"`
	// The key in the STF secret that contains the client id. Defaults to 'client-id'.
	ClientIDKey string `json:"clientIDKey,omitempty"`
	// The key in the STF secret that contains the client secret key. Defaults to 'client-secret'.
	ClientSecretKey string `json:"clientSecretKey,omitempty"`
}

// TraefikConfig represents configurations for the traefik deployment placed
// in front of the OpenSTF services.
type TraefikConfig struct {
	// (WIP) - use IngressRoute CRDs for an existing traefik deployment
	// instead of creating a standalone traefik service.
	UseIngressRoute bool `json:"useIngressRoute,omitempty"`
	// Configuration options for the traefik deployment.
	Deployment *TraefikDeployment `json:"deployment,omitempty"`
	// TLS configurations for traefik
	TLS *TLSConfig `json:"tls,omitempty"`
}

// TraefikDeployment represents configuration options for the traefik deployment.
type TraefikDeployment struct {
	// The number of traefik instances to run.
	Replicas int32 `json:"replicas,omitempty"`
	// The version of traefik to run, only >2.0 supported. Defaults to 2.2.0.
	Version string `json:"version,omitempty"`
	// The type of service to create for Traefik. Defaults to LoadBalancer.
	// If using external SSL from a pre-existing ingress controller, you'll want to
	// set this to ClusterIP.
	ServiceType string `json:"serviceType,omitempty"`
	// Set to true if you wish for traefik to produce access logs
	AccessLogs bool `json:"accessLogs,omitempty"`
	// A configuration for the traefik dashboard
	Dashboard *TraefikDashboard `json:"dashboard,omitempty"`
	// Resource restraints for the traefik deployment
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// TraefikDashboard represents configuration options for the Traefik dashboard.
type TraefikDashboard struct {
	// The hostname that should route to the traefik dashboard.
	Host string `json:"host"`
	// A list of IP addresses to whitelist for dashboard access.
	IPWhitelist []string `json:"ipWhitelist,omitempty"`
}

// AndroidFarmStatus defines the observed state of AndroidFarm
type AndroidFarmStatus struct {
	State string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidFarm is the Schema for the androidfarms API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=androidfarms,scope=Cluster
type AndroidFarm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AndroidFarmSpec   `json:"spec,omitempty"`
	Status AndroidFarmStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidFarmList contains a list of AndroidFarm
type AndroidFarmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AndroidFarm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AndroidFarm{}, &AndroidFarmList{})
}
