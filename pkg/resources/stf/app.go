package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	"github.com/go-logr/logr"
)

var appStartScript = `
for i in $(env | grep STF | grep -v APP_SECRET | cut -d '=' -f1) ; do unset $i ; done
stf app \
  --port 3000 \
  --auth-url %s \
  --websocket-url %s/
`

func (r *STFReconciler) reconcileApp(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	startCmd := fmt.Sprintf(appStartScript,
		instance.STFConfig().GetAuthURL(),
		instance.STFConfig().GetAppExternalWebsocketURL(),
	)
	return builders.NewDeploymentBuilder(reqLogger, instance, "app").
		WithResourceRequirements(instance.STFConfig().AppResourceRequirements()).
		WithPort("app", 3000).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{startCmd}).
		WithSecret("STF_APP_SECRET", instance.STFConfig().GetSTFSecretKey()).
		WithRethinkDB().
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}
