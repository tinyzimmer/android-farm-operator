package resources

import (
	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
)

// FarmReconciler represents an interface for ensuring resources for an AndroidFarm
type FarmReconciler interface {
	Reconcile(logr.Logger, *androidv1alpha1.AndroidFarm) error
}

// DeviceReconciler represents an interface for ensuring resources for a single device emulator.
type DeviceReconciler interface {
	Reconcile(logr.Logger, *androidv1alpha1.AndroidDevice) error
}
