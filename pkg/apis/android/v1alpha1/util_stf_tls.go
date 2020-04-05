package v1alpha1

import (
	"fmt"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
)

// TLSSecret returns the secret to be used for SSL configurations, or an empty
// string if SSL is disabled in the OpenSTF farm.
func (a *AndroidFarm) TLSSecret() string {
	if a.STFConfig().Traefik != nil {
		if a.STFConfig().Traefik.TLS != nil {
			if a.STFConfig().Traefik.TLS.TLSSecret != nil {
				return a.STFConfig().Traefik.TLS.TLSSecret.Name
			}
			if a.STFConfig().Traefik.TLS.IssuerRef != nil {
				return fmt.Sprintf("%s-ingress-tls", a.GetName())
			}
		}
	}
	return ""
}

func (s *STFConfig) GetIssuerReference() *cmmeta.ObjectReference {
	if s.Traefik != nil && s.Traefik.TLS != nil && s.Traefik.TLS.IssuerRef != nil {
		return s.Traefik.TLS.IssuerRef
	}
	return nil
}

// TLSEnabled returns true if TLS is being used in front of the OpenSTF
// farm.
func (s *STFConfig) TLSEnabled() bool {
	if s.Traefik != nil {
		return s.Traefik.TLS != nil
	}
	return false
}

// UseExternalTLS returns true if TLS is being managed externally for the
// OpenSTF farm.
func (s *STFConfig) UseExternalTLS() bool {
	if s.Traefik != nil && s.Traefik.TLS != nil && s.Traefik.TLS.External {
		return true
	}
	return false
}

// GetHTTPScheme returns the external HTTP scheme to be broadcasted by the
// OpenSTF deployments.
func (s *STFConfig) GetHTTPScheme() string {
	if s.TLSEnabled() {
		return "https"
	}
	return "http"
}

// GetWebsocketScheme returns the external websocket scheme to be broadcasted by
// the provider servers.
func (s *STFConfig) GetWebsocketScheme() string {
	if s.TLSEnabled() {
		return "wss"
	}
	return "ws"
}
