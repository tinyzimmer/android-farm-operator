package stf

import (
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
)

var websocketStartScript = `
for i in $(env | grep STF | grep -v WEBSOCKET_SECRET | cut -d '=' -f1) ; do unset $i ; done
stf websocket \
  --port 3000 \
  --storage-url %s \
  --connect-sub %s \
  --connect-push %s
`

func (r *STFReconciler) reconcileWebsocket(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	startCmd := fmt.Sprintf(websocketStartScript,
		instance.STFConfig().GetAppExternalURL(),
		fmt.Sprintf("%s:7150", stfutil.TriproxyEndpoint(instance, "app")),
		fmt.Sprintf("%s:7170", stfutil.TriproxyEndpoint(instance, "app")),
	)
	return builders.NewDeploymentBuilder(reqLogger, instance, "websocket").
		WithResourceRequirements(instance.STFConfig().WebsocketResourceRequirements()).
		WithPort("websocket", 3000).
		WithCommand([]string{"/bin/bash", "-c"}).
		WithArgs([]string{startCmd}).
		WithSecret("STF_WEBSOCKET_SECRET", instance.STFConfig().GetSTFSecretKey()).
		WithRethinkDB().
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait().
		Reconcile(r.client)
}
