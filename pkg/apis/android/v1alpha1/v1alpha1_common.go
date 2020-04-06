package v1alpha1

// Labels used for selecting pods and devices based off their inheritance
const (
	// DeviceConfigLabel is the selector matching devices to configurations
	// that they derive their spec from.
	DeviceConfigLabel = "deviceConfig"
	// DeviceFarmLabel is the selector matching devices to the farm they belong to.
	DeviceFarmLabel = "deviceFarm"
	// DeviceGroupLabel is the selector matching devices to the device group they
	// belong to.
	DeviceGroupLabel = "deviceGroup"
)

// Annotations used for internal operations on resources
const (
	// CreationSpecAnnotation contains the serialized creation spec of a resource
	// to be compared against desired state.
	CreationSpecAnnotation = "android.stf.io/creation-spec"
	// STFProviderAnnotation contains a reference to the stf-provider instance
	// that manages a device.
	STFProviderAnnotation = "android.stf.io/stf-provider"
	// ADBConnectedAnnotation is used to signal that a device is connected to
	// its ADB server.
	ADBConnectedAnnotation = "android.stf.io/adb-connected"
	// BootCompletedAnnotation is used to signal that a device has completed its
	// boot process.
	BootCompletedAnnotation = "android.stf.io/boot-completed"
	// ConfigMapSHAAnnotation is used to store the checksum of the configmap data
	// used when a deployment was created. If a deployment is reconciled with a new
	// checksum, it means its configuration has changed and pods need to be cycled.
	ConfigMapSHAAnnotation = "android.stf.io/configmap-checksum"
	// DeviceConfigSHAAnnotation is used to store the checksum of the configuration
	// used to provision a given device instance.
	DeviceConfigSHAAnnotation = "android.stf.io/device-config-checksum"
	// ProviderSerialAnnotation contains the name of a device as known by its
	// stf provider.
	ProviderSerialAnnotation = "android.stf.io/stf-serial"
)

// Defaults and other static vars
var (
	// defaultSTFImage is the default STF image to use for OpenSTF deployments.
	// Latest is less buggy and hopefully will be tagged soon
	defaultSTFImage = "openstf/stf:latest"
	// defaultRDBVersion is the default RethinkDB version to use for the StatefulSets.
	defaultRDBVersion = "2.4"
	// defaultTraefikVersion is the default version of traefik to run for OpenSTF
	// deployments
	defaultTraefikVersion = "v2.2.0"
	// defaultTraefikServiceType is the default service type to use for traefik
	// services.
	defaultTraefikServiceType = "LoadBalancer"
	// defaultTraefikWhitelist is the default whitelist to apply to the traefik
	// dashboard if ssl is enabled on the instance.
	defaultTraefikWhitelist = `["0.0.0.0/0"]`
	// The default service type to use for provider traefik instances.
	defaultTraefikProviderServiceType = "ClusterIP"
	// defaultRDBProxyReplicas is the default number of rethinkdb proxies to run in
	// front of the StatefulSets.
	defaultRDBProxyReplicas int32 = 1
	// defaultRDBReplicas is the default number of rethinkdb replicas to run in
	// the StatefulSets.
	defaultRDBReplicas int32 = 1
	// defaultRDBShards is the default number of shards to use per table in the
	// rethinkdb database.
	defaultRDBShards int32 = 1
	// defaultTraefikReplicas is the default number of traefik replicas to run in
	// OpenSTF deployments.
	defaultTraefikReplicas int32 = 1
	// defaultSTFSecretKey is the default key where the stf secret is found inside
	// a provided kubernetes secret.
	defaultSTFSecretKey = "stf-secret"
	// defaultOAuthClientIDKey is the default key where the oauth client ID is found
	// inside a provided kubernetes secret
	defaultOAuthClientIDKey = "client-id"
	// defaultOauthClientSecretKey is the default key where the oauth client secret
	// is found inside a provided kubernetes secret.
	defaultOAuthClientSecretKey = "client-secret"
	// defaultRunUser is the default user to run stf containers as
	defaultRunUser int64 = 1000

	// predefined bools to easily grab pointers to
	trueVal  = true
	falseVal = false
)
