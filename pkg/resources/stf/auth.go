package stf

import (
	"strings"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	"github.com/go-logr/logr"
)

func (r *STFReconciler) reconcileAuth(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	builder := builders.NewDeploymentBuilder(reqLogger, instance, "auth").
		WithResourceRequirements(instance.STFConfig().AuthResourceRequirements()).
		WithPort("stf-auth", 3000).
		WithRethinkDB().
		WithService("ClusterIP").
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithWait()

	if instance.STFConfig().UseMockAuth() {
		reqLogger.Info("Using mock authentication service")
		return builder.
			WithSecret("STF_AUTH_MOCK_SECRET", instance.STFConfig().GetSTFSecretKey()).
			WithArgs([]string{
				"auth-mock",
				"--port", "3000",
				"--app-url", instance.STFConfig().GetAppExternalURL(),
			}).
			Reconcile(r.client)
	}

	return builder.
		WithSecret("STF_AUTH_OAUTH2_SECRET", instance.STFConfig().GetSTFSecretKey()).
		WithSecret("STF_AUTH_OAUTH2_OAUTH_CLIENT_ID", instance.STFConfig().GetOAuthClientIDKey()).
		WithSecret("STF_AUTH_OAUTH2_OAUTH_CLIENT_SECRET", instance.STFConfig().GetOAuthClientSecretKey()).
		WithArgs([]string{
			"auth-oauth2",
			"--port", "3000",
			"--app-url", instance.STFConfig().GetAppExternalURL(),
			"--oauth-authorization-url", instance.STFConfig().Auth.OAuth.AuthorizationURL,
			"--oauth-token-url", instance.STFConfig().Auth.OAuth.TokenURL,
			"--oauth-userinfo-url", instance.STFConfig().Auth.OAuth.UserInfoURL,
			"--oauth-scope", strings.Join(instance.STFConfig().Auth.OAuth.Scopes, " "),
			"--oauth-callback-url", instance.STFConfig().Auth.OAuth.CallbackURL,
		}).
		Reconcile(r.client)
}
