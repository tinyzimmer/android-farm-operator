package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
)

var apiStartScript = `
for i in $(env | grep STF | grep -v API | cut -d '=' -f1) ; do unset $i ; done

stf api || \
	echo Detected older OpenSTF version && \
	unset STF_API_CONNECT_SUB_DEV && \
	unset STF_API_CONNECT_PUSH_DEV && \
	stf api
`

func (r *STFReconciler) reconcileAPI(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	startCmd := fmt.Sprintf(apiStartScript,
		fmt.Sprintf("%s:7150", stfutil.TriproxyEndpoint(instance, "app")),
		fmt.Sprintf("%s:7170", stfutil.TriproxyEndpoint(instance, "app")),
		fmt.Sprintf("%s:7250", stfutil.TriproxyEndpoint(instance, "dev")),
		fmt.Sprintf("%s:7270", stfutil.TriproxyEndpoint(instance, "dev")),
	)
	return builders.NewDeploymentBuilder(reqLogger, instance, "api").
		WithResourceRequirements(instance.STFConfig().APIResourceRequirements()).
		WithPort("api", 3000).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{startCmd}).
		WithRethinkDB().
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithSecret("STF_API_SECRET", instance.STFConfig().GetSTFSecretKey()).
		WithEnvVar("STF_API_PORT", "3000").
		WithEnvVar("STF_API_CONNECT_SUB", fmt.Sprintf("%s:7150", stfutil.TriproxyEndpoint(instance, "app"))).
		WithEnvVar("STF_API_CONNECT_PUSH", fmt.Sprintf("%s:7170", stfutil.TriproxyEndpoint(instance, "app"))).
		WithEnvVar("STF_API_CONNECT_SUB_DEV", fmt.Sprintf("%s:7250", stfutil.TriproxyEndpoint(instance, "dev"))).
		WithEnvVar("STF_API_CONNECT_PUSH_DEV", fmt.Sprintf("%s:7270", stfutil.TriproxyEndpoint(instance, "dev"))).
		WithWait().
		Reconcile(r.client)
}
