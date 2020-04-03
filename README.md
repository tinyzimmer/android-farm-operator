# android-farm-operator

This repository contains a Kubernetes operator for managing android emulators, physical devices, and OpenSTF deployments on bare-metal clusters.
It is written in `go` and based off the [`operator-sdk`](https://github.com/operator-framework/operator-sdk) framework.

The main CRD provided by the operator is the `AndroidFarm`. The controller will manage the following resources for you:

 - RethinkDB and proxies (replicated/sharded if desired)
 - An [OpenSTF](https://openstf.io/) cluster
 - Providers and ADB servers for host usb devices
 - Emulator farms within the cluster
   - Defined by the CRD `AndroidDeviceConfig`
   - Manages emulator pods, providers, and ADB servers for each emulator "device group".
 - Traefik HTTP/TCP Proxies for `rethinkdb-admin`, `openstf`, and `provider-adb` interfaces.
   - TLS or cleartext HTTP configurations available

There is also some heavy WIP functionality for defining `AndroidJobTemplates` and `AndroidJobs` to run commands/inputs across multiple devices at once using custom resources.
More docs on that and a stable implementation may come later.

## Getting Started

### Quickstart

```bash
$> git clone https://github.com/tinyzimmer/android-farm-operator
$> cd android-farm-operator
# Install using the helm chart
$> helm install android-farm-operator deploy/charts/android-farm-operator
```

Take a look in `deploy/examples` to see some basic configuration options. There are more outlined in the full reference. Edit `example-config.yaml` and `example-farm.yaml` to your desires (skip the `example-config` if just doing host USB devices) and then apply the manifests.

```bash
$> kubectl apply -f deploy/examples/example-config.yaml  # If launching emulated devices, requires KVM
$> kubectl apply -f deploy/examples/example-farm.yaml    # Deploys rethinkdb, stf resources, and optional
                                                         # emulator pods.
```

Please refer to the example manifests or the full [CRD reference]((doc/crds.md)) for explanations
on configuration options.

If you don't have a bare metal kubernetes cluster available to you but you still want to play with the operator, I've tested most things using [`kind`](https://github.com/kubernetes-sigs/kind). So as long as you have docker available you should be able to deploy the operator and all of it's resources. Some considerations per host OS though:

- **Linux distro** - Everything should work, as long as you have kvm enabled
- **macOS** - USB devices won't work on regular docker-for-mac. But if you use `docker-machine` with the virtualbox driver and enable usb support in your machine, it should work. Emulators won't work tho coz no kvm.
- **Windows** - nah

I have helpers in the Makefile to get you going quickly.
You'll need to have at least `go`, `docker`, and `helm3` installed.

```bash
$> make all-of-it      # Just do everything (with no previous cache, takes some time)

$> make build-emulator # Build a local copy of the emulator docker image (or you can pull it
                       # with docker pull quay.io/tinyzimmer/android-emulator:android-29-slim.
                       # It's about 2 GB.)
$> make test-cluster   # Creates a kind cluster with kvm/usb passthrough and metallb load balancing

$> make build-operator load-operator # Builds the operator image and loads it into kind
$> make deploy                       # Deploys the operator and crds

$> make example-farm                 # Deploys the example farm (defaults assume you edit /etc/hosts
                                     # to point stf.kind.local and stf-traefik.kind.local to the metallb
                                     # IP address)
```

On a standard linux installation of docker, (probably possible with some port-forwarding on mac) you can now edit your `/etc/hosts` so that `stf.kind.local` and `traefik.kind.local` resolve to 172.17.255.1. And then in your browser use the following endpoints:

- https://stf.kind.local - The main STF UI
- https://stf.kind.local/rethinkdb/ - The rethinkdb admin interface
- https://traefik.kind.local - The traefik dashboard

### Extra resources

 - [Helm Configurations](doc/helm.md)
 - [`CustomResourceDefinitions` Reference](doc/crds.md)



#### TODO

  - [ ] Unit tests
  - [ ] More docs
  - [ ] CI
