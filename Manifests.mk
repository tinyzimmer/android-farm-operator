## Cluster manifests

define KIND_CLUSTER_MANIFEST
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
  extraMounts:
  - hostPath: /dev/kvm
    containerPath: /dev/kvm
  - hostPath: /dev/bus/usb
    containerPath: /dev/bus/usb
- role: worker
  extraMounts:
  - hostPath: /dev/kvm
    containerPath: /dev/kvm
endef

define METALLB_CONFIG
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: metallb-system
  name: config
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - 172.17.255.1-172.17.255.250
endef

export KIND_CLUSTER_MANIFEST
export METALLB_CONFIG

##
