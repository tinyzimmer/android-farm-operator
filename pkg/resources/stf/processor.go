package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
)

func (r *STFReconciler) reconcileProcessor(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "processor").
		WithResourceRequirements(instance.STFConfig().ProcessorResourceRequirements()).
		WithArgs([]string{
			"processor", "stf-processor",
			"--connect-app-dealer", fmt.Sprintf("%s:7160", stfutil.TriproxyEndpoint(instance, "app")),
			"--connect-dev-dealer", fmt.Sprintf("%s:7260", stfutil.TriproxyEndpoint(instance, "dev")),
		}).
		WithRethinkDB().
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}
