package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func resourceLimits(cpu, memMB int64) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			"cpu":    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
			"memory": *resource.NewQuantity(memMB*1024*1024, resource.BinarySI),
		},
	}
}

var defaultRequirements = resourceLimits(100, 256)

func (s *STFConfig) ADBResourceRequirements() corev1.ResourceRequirements {
	if s.ADB != nil {
		if s.ADB.Resources != nil {
			return *s.ADB.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) APIResourceRequirements() corev1.ResourceRequirements {
	if s.API != nil {
		if s.API.Resources != nil {
			return *s.API.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) AppResourceRequirements() corev1.ResourceRequirements {
	if s.App != nil {
		if s.App.Resources != nil {
			return *s.App.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) AuthResourceRequirements() corev1.ResourceRequirements {
	if s.Auth != nil {
		if s.Auth.Resources != nil {
			return *s.Auth.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) ProcessorResourceRequirements() corev1.ResourceRequirements {
	if s.Processor != nil {
		if s.Processor.Resources != nil {
			return *s.Processor.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) ProviderResourceRequirements() corev1.ResourceRequirements {
	if s.Provider != nil {
		if s.Provider.Resources != nil {
			return *s.Provider.Resources
		}
	}
	return resourceLimits(500, 1024)
}

func (s *STFConfig) ReaperResourceRequirements() corev1.ResourceRequirements {
	if s.Reaper != nil {
		if s.Reaper.Resources != nil {
			return *s.Reaper.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) StorageResourceRequirements() corev1.ResourceRequirements {
	if s.Storage != nil {
		if s.Storage.Resources != nil {
			return *s.Storage.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) APKStorageResourceRequirements() corev1.ResourceRequirements {
	return s.StorageResourceRequirements()
}

func (s *STFConfig) ImgStorageResourceRequirements() corev1.ResourceRequirements {
	return s.StorageResourceRequirements()
}

func (s *STFConfig) TriproxyAppResourceRequirements() corev1.ResourceRequirements {
	if s.TriproxyApp != nil {
		if s.TriproxyApp.Resources != nil {
			return *s.TriproxyApp.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) TriproxyDevResourceRequirements() corev1.ResourceRequirements {
	if s.TriproxyDev != nil {
		if s.TriproxyDev.Resources != nil {
			return *s.TriproxyDev.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) WebsocketResourceRequirements() corev1.ResourceRequirements {
	if s.Websocket != nil {
		if s.Websocket.Resources != nil {
			return *s.Websocket.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) RethinkDBResourceRequirements() corev1.ResourceRequirements {
	if s.RethinkDB != nil {
		if s.RethinkDB.Resources != nil {
			return *s.RethinkDB.Resources
		}
	}
	return resourceLimits(250, 1024)
}

func (s *STFConfig) RethinkDBProxyResourceRequirements() corev1.ResourceRequirements {
	if s.RethinkDBProxy != nil {
		if s.RethinkDBProxy.Resources != nil {
			return *s.RethinkDBProxy.Resources
		}
	}
	return defaultRequirements
}

func (s *STFConfig) TraefikResourceRequirements() corev1.ResourceRequirements {
	if s.Traefik != nil {
		if s.Traefik.Deployment != nil {
			if s.Traefik.Deployment.Resources != nil {
				return *s.Traefik.Deployment.Resources
			}
		}
	}
	return defaultRequirements
}
