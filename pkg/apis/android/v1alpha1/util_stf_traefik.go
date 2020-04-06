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

// TraefikServiceAnnotations returns annotations to apply to the traefik service
// for an stf cluster.
func (s *STFConfig) TraefikServiceAnnotations() map[string]string {
	if s.Traefik != nil && s.Traefik.Deployment != nil && s.Traefik.Deployment.Annotations != nil {
		return s.Traefik.Deployment.Annotations
	}
	return map[string]string{}
}

// ProviderTraefikReplicas returns the number of replicas to run in the traefik
// deployment for this device group's provider.
func (g *DeviceGroup) ProviderTraefikReplicas() int32 {
	if g.Provider != nil {
		if g.Provider.Traefik != nil {
			if g.Provider.Traefik.Replicas != 0 {
				return g.Provider.Traefik.Replicas
			}
		}
	}
	return defaultTraefikReplicas
}

// ProviderTraefikServiceType returns the type of service to use in the traefik
// deployment for this device group's provider.
func (g *DeviceGroup) ProviderTraefikServiceType() string {
	if g.Provider != nil {
		if g.Provider.Traefik != nil {
			if g.Provider.Traefik.ServiceType != "" {
				return g.Provider.Traefik.ServiceType
			}
		}
	}
	return defaultTraefikProviderServiceType
}

// ProviderTraefikDashboardEnabled returns true if we should enable the traefik
// dashboard on this provider deployment.
func (g *DeviceGroup) ProviderTraefikDashboardEnabled() bool {
	if g.Provider != nil {
		if g.Provider.Traefik != nil {
			if g.Provider.Traefik.Dashboard != nil {
				return true
			}
		}
	}
	return false
}

// ProviderTraefikAccessLogsEnabled returns true if we should enable the access
// logs on this provider's traefik instance.
func (g *DeviceGroup) ProviderTraefikAccessLogsEnabled() bool {
	if g.Provider != nil {
		if g.Provider.Traefik != nil {
			return g.Provider.Traefik.AccessLogs
		}
	}
	return false
}

// ProviderTraefikServiceAnnotations returns the annotations to apply to the
// traefik instance for this device group.
func (g *DeviceGroup) ProviderTraefikServiceAnnotations() map[string]string {
	if g.Provider != nil {
		if g.Provider.Traefik != nil {
			if g.Provider.Traefik.Annotations != nil {
				return g.Provider.Traefik.Annotations
			}
		}
	}
	return map[string]string{}
}

// UseClusterLocalADB returns true if we should only set up adb entrypoints on
// the provider traefik instances, instead of the main one.
func (g *DeviceGroup) UseClusterLocalADB() bool {
	return g.Provider != nil && g.Provider.ClusterLocalADB
}

func (g *DeviceGroup) ProviderHostnameOverride() string {
	if g.Provider != nil {
		if g.Provider.HostnameOverride != "" {
			return g.Provider.HostnameOverride
		}
	}
	return ""
}
