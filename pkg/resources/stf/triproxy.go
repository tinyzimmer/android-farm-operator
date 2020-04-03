package stf

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	"github.com/go-logr/logr"
)

var triproxyAppStartScript = `
unset "${!STF@}" ; stf triproxy app \
  --bind-pub "tcp://*:7150" \
  --bind-dealer "tcp://*:7160" \
  --bind-pull "tcp://*:7170"
`

var triproxyDevStartScript = `
unset "${!STF@}" ; stf triproxy dev \
  --bind-pub "tcp://*:7250" \
  --bind-dealer "tcp://*:7260" \
  --bind-pull "tcp://*:7270"
`

func (r *STFReconciler) reconcileTriproxyApp(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "triproxy-app").
		WithPort("bind-pub", 7150).
		WithPort("bind-dealer", 7160).
		WithPort("bind-pull", 7170).
		WithResourceRequirements(instance.STFConfig().TriproxyAppResourceRequirements()).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{triproxyAppStartScript}).
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}

func (r *STFReconciler) reconcileTriproxyDev(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "triproxy-dev").
		WithPort("bind-pub", 7250).
		WithPort("bind-dealer", 7260).
		WithPort("bind-pull", 7270).
		WithResourceRequirements(instance.STFConfig().TriproxyDevResourceRequirements()).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{triproxyDevStartScript}).
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}
