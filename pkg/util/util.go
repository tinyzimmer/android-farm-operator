package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type FarmReconcileFunc func(logr.Logger, *androidv1alpha1.AndroidFarm) error

func BoolPointer(b bool) *bool {
	return &b
}

func IsMarkedForDeletion(cr *androidv1alpha1.AndroidFarm) bool {
	return cr.GetDeletionTimestamp() != nil
}

func SetCreationSpecAnnotation(meta *metav1.ObjectMeta, obj runtime.Object) error {
	annotations := meta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	spec, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	annotations[androidv1alpha1.CreationSpecAnnotation] = string(spec)
	meta.SetAnnotations(annotations)
	return nil
}

func CreationSpecsEqual(m1 metav1.ObjectMeta, m2 metav1.ObjectMeta) bool {
	m1ann := m1.GetAnnotations()
	m2ann := m2.GetAnnotations()
	spec1, ok := m1ann[androidv1alpha1.CreationSpecAnnotation]
	if !ok {
		return false
	}
	spec2, ok := m2ann[androidv1alpha1.CreationSpecAnnotation]
	if !ok {
		return false
	}
	return spec1 == spec2
}

func GetClusterSuffix() string {
	resolvconf, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		return ""
	}
	re := regexp.MustCompile("search.*")
	match := re.FindString(string(resolvconf))
	if strings.TrimSpace(match) == "" {
		return ""
	}
	fields := strings.Fields(match)
	return fields[len(fields)-1]
}

func GetPodADBPort(pod corev1.Pod) (int32, error) {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.Name == "adb" {
				return port.ContainerPort, nil
			}
		}
	}
	return 0, fmt.Errorf("Could not determine ADB port for device: %s/%s", pod.Namespace, pod.Name)
}
