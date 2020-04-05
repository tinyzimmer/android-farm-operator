package v1alpha1

import (
	"encoding/json"
	"fmt"
)

// UseIngressRouteCRD returns true if we are provisioning IngressRoutes for an
// external traefik deployment.
func (s *STFConfig) UseIngressRouteCRD() bool {
	if s.Traefik != nil {
		if s.Traefik.UseIngressRoute {
			return true
		}
	}
	return false
}

// TraefikImage returns the docker image to use for the traefik deployment.
func (s *STFConfig) TraefikImage() string {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.Version != "" {
				return fmt.Sprintf("traefik:%s", s.Traefik.Deployment.Version)
			}
		}
	}
	return fmt.Sprintf("traefik:%s", defaultTraefikVersion)
}

// TraefikReplicas returns the number of replicas to run in the traefik deployment.
func (s *STFConfig) TraefikReplicas() int32 {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.Replicas != 0 {
				return s.Traefik.Deployment.Replicas
			}
		}
	}
	return defaultTraefikReplicas
}

// TraefikServiceType returns the type of service to create in front of the
// traefik deployment.
func (s *STFConfig) TraefikServiceType() string {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.ServiceType != "" {
				return s.Traefik.Deployment.ServiceType
			}
		}
	}
	return defaultTraefikServiceType
}

// TraefikAccessLogsEnabled returns true if the access log should be enabled in
// traefik
func (s *STFConfig) TraefikAccessLogsEnabled() bool {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			return s.Traefik.Deployment.AccessLogs
		}
	}
	return false
}

// TraefikDashboardEnabled returns true if the dashboard and api should be
// enabled in the traefik deployment.
func (s *STFConfig) TraefikDashboardEnabled() bool {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			return s.Traefik.Deployment.Dashboard != nil
		}
	}
	return false
}

// TraefikDashboardHost returns the hostname for which requests should be
// routed to the traefik dashboard.
func (s *STFConfig) TraefikDashboardHost() string {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.Dashboard != nil {
				return s.Traefik.Deployment.Dashboard.Host
			}
		}
	}
	return ""
}

// TraefikDashboardWhitelistString returns the json representation of the
// IPs to be whitelisted for the traefik dashboard. If the user-provided list
// can't be serialized, we use the default allow-all rule.
func (s *STFConfig) TraefikDashboardWhitelistString() string {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.Dashboard != nil {
				if s.Traefik.Deployment.Dashboard.IPWhitelist != nil {
					out, err := json.Marshal(s.Traefik.Deployment.Dashboard.IPWhitelist)
					if err == nil {
						return string(out)
					}
				}
			}
		}
	}
	return defaultTraefikWhitelist
}

func (s *STFConfig) TraefikServiceNames() []string {
	if s.Traefik != nil && s.Traefik.Deployment != nil && s.Traefik.Deployment.ServiceNames != nil {
		return s.Traefik.Deployment.ServiceNames
	}
	return []string{}
}

func (s *STFConfig) TraefikServiceAnnotations() map[string]string {
	if s.Traefik != nil && s.Traefik.Deployment != nil && s.Traefik.Deployment.Annotations != nil {
		return s.Traefik.Deployment.Annotations
	}
	return map[string]string{}
}
