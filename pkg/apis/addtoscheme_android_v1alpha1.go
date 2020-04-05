package apis

import (
	cm "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	"github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, cm.SchemeBuilder.AddToScheme)
}
