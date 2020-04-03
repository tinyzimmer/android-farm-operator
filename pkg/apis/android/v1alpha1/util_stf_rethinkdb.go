package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RethinkDBName returns the name to use for the rethinkdb resources in this
// AndroidFarm.
func (a *AndroidFarm) RethinkDBName() string {
	return fmt.Sprintf("%s-rethinkdb", a.GetName())
}

func (a *AndroidFarm) RethinkDBProxyName() string {
	return fmt.Sprintf("%s-proxy", a.RethinkDBName())
}

// RethinkDBAdminName returns the name to use for the rethinkdb admin resources
func (a *AndroidFarm) RethinkDBAdminName() string {
	return fmt.Sprintf("%s-admin", a.RethinkDBName())
}

// RethinkDBEnvVars returns the environment variables to supply to the
// rethinkdb pods.
func (a *AndroidFarm) RethinkDBEnvVars() []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name: "POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
		{
			Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name:  "SERVICE_NAME",
			Value: a.RethinkDBName(),
		},
		{
			Name:  "PROXY_NAME",
			Value: a.RethinkDBProxyName(),
		},
	}
	return envVars
}

// RethinkDBLabels returns the labels to apply to rethinkdb resources.
func (a *AndroidFarm) RethinkDBLabels(component string) map[string]string {
	return map[string]string{
		DeviceFarmLabel: a.GetName(),
		"db":            component,
	}
}

// RethinkDBVolumeClaims returns the PVC templates to supply to the RethinkDB
// StatefulSet.
func (a *AndroidFarm) RethinkDBVolumeClaims() []corev1.PersistentVolumeClaim {
	if a.STFConfig() != nil {
		if a.STFConfig().RethinkDB != nil {
			if a.STFConfig().RethinkDB.PVCSpec != nil {
				return []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: fmt.Sprintf("%s-rethinkdb-storage", a.GetName()),
						},
						Spec: *a.STFConfig().RethinkDB.PVCSpec,
					},
				}
			}
		}
	}
	return []corev1.PersistentVolumeClaim{}
}

// RethinkDBVolumeMounts returns the volume mounts to apply to the pods in the
// RethinkDB StatefulSet.
func (a *AndroidFarm) RethinkDBVolumeMounts() []corev1.VolumeMount {
	if a.STFConfig() != nil {
		if a.STFConfig().RethinkDB != nil {
			if a.STFConfig().RethinkDB.PVCSpec != nil {
				return []corev1.VolumeMount{
					{
						Name:      fmt.Sprintf("%s-rethinkdb-storage", a.GetName()),
						MountPath: "/data/rethinkdb_data",
					},
				}
			}
		}
	}
	return []corev1.VolumeMount{}
}

func (s *STFConfig) RethinkDBProxyReplicas() *int32 {
	if s.RethinkDBProxy != nil {
		if s.RethinkDBProxy.Replicas > 0 {
			return &s.RethinkDBProxy.Replicas
		}
	}
	return &defaultRDBProxyReplicas
}

// RethinkDBReplicas returns the number of replicas per shard to run in the rethinkdb
// StatefulSet.
func (s *STFConfig) RethinkDBReplicas() *int32 {
	if s.RethinkDB != nil {
		if s.RethinkDB.Replicas > 0 {
			return &s.RethinkDB.Replicas
		}
	}
	return &defaultRDBReplicas
}

// RethinkDBShards returns the number of shards to use in the rethinkdb
// StatefulSet tables.
func (s *STFConfig) RethinkDBShards() *int32 {
	if s.RethinkDB != nil {
		if s.RethinkDB.Shards > 0 {
			return &s.RethinkDB.Shards
		}
	}
	return &defaultRDBShards
}

// RethinkDBPullPolicy returns the pull policy to use in the rethinkdb StatefulSet.
func (s *STFConfig) RethinkDBPullPolicy() corev1.PullPolicy {
	if s.RethinkDB != nil {
		if s.RethinkDB.ImagePullPolicy != "" {
			return s.RethinkDB.ImagePullPolicy
		}
	}
	return corev1.PullAlways
}

// RethinkDBImage returns the image to use in the rethinkdb StatefulSet.
func (s *STFConfig) RethinkDBImage() string {
	if s.RethinkDB != nil {
		if s.RethinkDB.Version != "" {
			return fmt.Sprintf("rethinkdb:%s", s.RethinkDB.Version)
		}
	}
	return fmt.Sprintf("rethinkdb:%s", defaultRDBVersion)
}
