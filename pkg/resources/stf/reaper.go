package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
)

func (r *STFReconciler) reconcileReaper(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "reaper").
		WithResourceRequirements(instance.STFConfig().ReaperResourceRequirements()).
		WithArgs([]string{
			"reaper", "dev",
			"--connect-push", fmt.Sprintf("%s:7270", stfutil.TriproxyEndpoint(instance, "dev")),
			"--connect-sub", fmt.Sprintf("%s:7150", stfutil.TriproxyEndpoint(instance, "app")),
			"--heartbeat-timeout", "30000",
		}).
		WithRethinkDB().
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}
