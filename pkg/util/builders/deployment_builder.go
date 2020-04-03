package builders

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeploymentBuilder struct {
	cr                       *androidv1alpha1.AndroidFarm
	component                string
	image                    string
	replicas                 int32
	cmd                      []string
	args                     []string
	volumes                  []corev1.Volume
	volumeMounts             []corev1.VolumeMount
	ports                    []map[string]int32
	resources                corev1.ResourceRequirements
	envVars                  []corev1.EnvVar
	service                  bool
	serviceType              corev1.ServiceType
	headless                 bool
	files                    []map[string]string
	podSecurityContext       *corev1.PodSecurityContext
	containerSecurityContext *corev1.SecurityContext
	sidecars                 []corev1.Container
	nodeSelector             map[string]string
	wait                     bool
	logger                   logr.Logger
}

func NewDeploymentBuilder(logger logr.Logger, cr *androidv1alpha1.AndroidFarm, component string) *DeploymentBuilder {
	return &DeploymentBuilder{
		cr:        cr,
		logger:    logger,
		component: component,
		cmd:       []string{"stf"},
		envVars: []corev1.EnvVar{
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
				Name: "POD_NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.namespace",
					},
				},
			},
			{
				Name:  "CLUSTER_SUFFIX",
				Value: util.GetClusterSuffix(),
			},
		},
		resources: corev1.ResourceRequirements{
			Limits:   corev1.ResourceList{},
			Requests: corev1.ResourceList{},
		},
	}
}

func (d *DeploymentBuilder) WithReplicas(i int32) *DeploymentBuilder {
	d.replicas = i
	return d
}

func (d *DeploymentBuilder) WithCommand(s []string) *DeploymentBuilder {
	d.cmd = s
	return d
}

func (d *DeploymentBuilder) WithArgs(s []string) *DeploymentBuilder {
	d.args = s
	return d
}

func (d *DeploymentBuilder) WithImage(img string) *DeploymentBuilder {
	d.image = img
	return d
}

func (d *DeploymentBuilder) WithVolumes(vols []corev1.Volume, mounts []corev1.VolumeMount) *DeploymentBuilder {
	if d.volumes == nil {
		d.volumes = make([]corev1.Volume, 0)
		d.volumeMounts = make([]corev1.VolumeMount, 0)
	}
	d.volumes = append(d.volumes, vols...)
	if mounts != nil {
		d.volumeMounts = append(d.volumeMounts, mounts...)
	}
	return d
}

func (d *DeploymentBuilder) WithPort(name string, port int32) *DeploymentBuilder {
	if d.ports == nil {
		d.ports = make([]map[string]int32, 0)
	}
	d.ports = append(d.ports, map[string]int32{name: port})
	return d
}

func (d *DeploymentBuilder) WithResourceRequirements(reqs corev1.ResourceRequirements) *DeploymentBuilder {
	d.resources = reqs
	return d
}

func (d *DeploymentBuilder) WithService(t string) *DeploymentBuilder {
	d.service = true
	d.serviceType = corev1.ServiceType(t)
	d.envVars = append(d.envVars, corev1.EnvVar{
		Name:  "SERVICE_NAME",
		Value: fmt.Sprintf("%s-%s", d.cr.STFNamePrefix(), d.component),
	})
	return d
}

func (d *DeploymentBuilder) WithHeadlessService() *DeploymentBuilder {
	d.headless = true
	return d.WithService("")
}

func (d *DeploymentBuilder) WithEnvVar(name, value string) *DeploymentBuilder {
	d.envVars = append(d.envVars, corev1.EnvVar{
		Name:  name,
		Value: value,
	})
	return d
}

func (d *DeploymentBuilder) WithEnvVarFromSecret(name, key string) *DeploymentBuilder {
	d.envVars = append(d.envVars, corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: d.cr.STFConfig().Secret,
				},
				Key: key,
			},
		},
	})
	return d
}

func (d *DeploymentBuilder) WithRethinkDB() *DeploymentBuilder {
	return d.WithEnvVar("RETHINKDB_PORT_28015_TCP", stfutil.RethinkDBProxyEndpoint(d.cr))
	// return d.WithEnvVar("RETHINKDB_PORT_28015_TCP", stfutil.RethinkDBEndpoint(d.cr, master))
}

func (d *DeploymentBuilder) WithSecret(name, key string) *DeploymentBuilder {
	return d.WithEnvVarFromSecret(name, key)
}

func (d *DeploymentBuilder) WithFile(path, content string) *DeploymentBuilder {
	if d.files == nil {
		d.files = make([]map[string]string, 0)
	}
	d.files = append(d.files, map[string]string{path: content})
	return d
}

func (d *DeploymentBuilder) WithSidecar(container corev1.Container) *DeploymentBuilder {
	if d.sidecars == nil {
		d.sidecars = make([]corev1.Container, 0)
	}
	container.Env = d.envVars
	d.sidecars = append(d.sidecars, container)
	return d
}

func (d *DeploymentBuilder) WithWait() *DeploymentBuilder {
	d.wait = true
	return d
}

func (d *DeploymentBuilder) WithPodSecurityContext(ctx *corev1.PodSecurityContext) *DeploymentBuilder {
	d.podSecurityContext = ctx
	return d
}

func (d *DeploymentBuilder) WithContainerSecurityContext(ctx *corev1.SecurityContext) *DeploymentBuilder {
	d.containerSecurityContext = ctx
	return d
}

func (d *DeploymentBuilder) WithNodeSelector(selector map[string]string) *DeploymentBuilder {
	d.nodeSelector = selector
	return d
}

func (d *DeploymentBuilder) getReplicas() *int32 {
	if d.replicas == 0 {
		d.replicas = 1
	}
	return &d.replicas
}

func (d *DeploymentBuilder) getCommand() []string {
	return d.cmd
}

func (d *DeploymentBuilder) getArgs() []string {
	if d.args == nil {
		return []string{}
	}
	return d.args
}

func (d *DeploymentBuilder) getVolumes() []corev1.Volume {
	if d.volumes == nil {
		return []corev1.Volume{}
	}
	return d.volumes
}

func (d *DeploymentBuilder) getVolumeMounts() []corev1.VolumeMount {
	if d.volumeMounts == nil {
		return []corev1.VolumeMount{}
	}
	return d.volumeMounts
}

func (d *DeploymentBuilder) getContainerPorts() []corev1.ContainerPort {
	ports := make([]corev1.ContainerPort, 0)
	if d.ports == nil {
		return ports
	}
	for _, portDef := range d.ports {
		for name, port := range portDef {
			ports = append(ports, corev1.ContainerPort{
				Name:          name,
				ContainerPort: port,
			})
		}
	}
	return ports
}

func (d *DeploymentBuilder) getImage() string {
	if d.image == "" {
		return d.cr.STFConfig().OpenSTFImage()
	}
	return d.image
}

func (d *DeploymentBuilder) buildDeployment() *appsv1.Deployment {
	containers := []corev1.Container{
		{
			Name:            fmt.Sprintf("stf-%s", d.component),
			ImagePullPolicy: d.cr.STFConfig().PullPolicy(),
			Image:           d.getImage(),
			Env:             d.envVars,
			VolumeMounts:    d.getVolumeMounts(),
			Command:         d.getCommand(),
			Args:            d.getArgs(),
			Ports:           d.getContainerPorts(),
			Resources:       d.resources,
			SecurityContext: d.containerSecurityContext,
		},
	}
	if d.sidecars != nil {
		containers = append(containers, d.sidecars...)
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s", d.cr.STFNamePrefix(), d.component),
			Namespace:       d.cr.STFConfig().GetNamespace(),
			Labels:          d.cr.STFComponentLabels(d.component),
			OwnerReferences: d.cr.OwnerReferences(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: d.getReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: d.cr.STFComponentLabels(d.component),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.cr.STFComponentLabels(d.component),
				},
				Spec: corev1.PodSpec{
					NodeSelector:       d.nodeSelector,
					ServiceAccountName: d.cr.STFConfig().GetServiceAccount(),
					SecurityContext:    d.podSecurityContext,
					Volumes:            d.getVolumes(),
					ImagePullSecrets:   d.cr.STFConfig().PullSecrets(),
					Containers:         containers,
				},
			},
		},
	}
}

func (d *DeploymentBuilder) buildService() *corev1.Service {
	ports := make([]corev1.ServicePort, 0)
	for _, portDef := range d.ports {
		for name, port := range portDef {
			ports = append(ports, corev1.ServicePort{
				Name:       name,
				Port:       port,
				TargetPort: intstr.FromInt(int(port)),
			})
		}
	}
	if d.sidecars != nil {
		for _, sidecar := range d.sidecars {
			if sidecar.Ports != nil {
				for _, port := range sidecar.Ports {
					ports = append(ports, corev1.ServicePort{
						Name:       port.Name,
						Port:       port.ContainerPort,
						TargetPort: intstr.FromInt(int(port.ContainerPort)),
					})
				}
			}
		}
	}
	var clusterIP string
	if d.headless {
		clusterIP = "None"
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s", d.cr.STFNamePrefix(), d.component),
			Namespace:       d.cr.STFConfig().GetNamespace(),
			Labels:          d.cr.STFComponentLabels(d.component),
			OwnerReferences: d.cr.OwnerReferences(),
		},
		Spec: corev1.ServiceSpec{
			Type:      d.serviceType,
			ClusterIP: clusterIP,
			Ports:     ports,
			Selector:  d.cr.STFComponentLabels(d.component),
		},
	}
}

func (d *DeploymentBuilder) buildConfigMap() *corev1.ConfigMap {
	data := make(map[string]string)
	for _, fileDef := range d.files {
		for fpath, fcontent := range fileDef {
			data[filepath.Base(fpath)] = fcontent
		}
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%s", d.cr.STFNamePrefix(), d.component),
			Namespace:       d.cr.STFConfig().GetNamespace(),
			Labels:          d.cr.STFComponentLabels(d.component),
			OwnerReferences: d.cr.OwnerReferences(),
		},
		Data: data,
	}
}

func (d *DeploymentBuilder) appendConfigMapVolumes() {
	items := make([]corev1.KeyToPath, 0)
	for _, fileDef := range d.files {
		for fpath := range fileDef {
			items = append(items, corev1.KeyToPath{Key: filepath.Base(fpath), Path: fpath})
		}
	}
	d.WithVolumes(
		[]corev1.Volume{
			{
				Name: "configmap",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s-%s", d.cr.STFNamePrefix(), d.component),
						},
						Items: items,
					},
				},
			},
		},
		[]corev1.VolumeMount{
			{
				Name:      "configmap",
				MountPath: "/etc/configmap",
			},
		},
	)
}

func (d *DeploymentBuilder) Reconcile(c client.Client) error {
	var cm *corev1.ConfigMap
	if d.files != nil {
		d.appendConfigMapVolumes()
		cm = d.buildConfigMap()
		if err := util.ReconcileConfigMap(d.logger, c, cm); err != nil {
			return err
		}
	}

	if d.service {
		if err := util.ReconcileService(d.logger, c, d.buildService()); err != nil {
			return err
		}
	}

	deployment := d.buildDeployment()
	if cm != nil {
		sha, err := computeConfigmapSHA(cm)
		if err != nil {
			return err
		}
		deployment.Spec.Template.SetAnnotations(map[string]string{
			androidv1alpha1.ConfigMapSHAAnnotation: sha,
		})
	}
	if err := util.ReconcileDeployment(d.logger, c, deployment, d.wait); err != nil {
		return err
	}

	return nil
}

func computeConfigmapSHA(cm *corev1.ConfigMap) (string, error) {
	out, err := json.Marshal(cm.Data)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := h.Write(out); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
