APP_VERSION ?= 0.0.1


-include Makevars.mk
-include Manifests.mk


###
# Building
###

.PHONY: build
build: build-operator

# Build all images
build-all: build-operator build-emulator build-adbmon build-goredir

# Same as build
build-operator: generate
	docker build . \
		-f build/operatorDockerfile \
		-t ${OPERATOR_IMAGE} \
		--build-arg GIT_COMMIT=`git rev-parse HEAD`

# Build an emulator image
.PHONY: build-emulator
build-emulator:
	mkdir -p build/apps
	docker build . \
		-f build/emulatorDockerfile \
		-t ${EMULATOR_IMAGE} \
		--build-arg PLATFORM=${EMULATOR_PLATFORM} \
		--build-arg CLI_TOOLS_VERSION=${EMULATOR_CLI_TOOLS}

# Build the adbmon image
build-adbmon:
	docker build . \
		-f build/adbmonDockerfile \
		-t ${ADB_IMAGE} \
		--build-arg ANDROID_TOOLS_VERSION=${ANDROID_TOOLS_VERSION}

# Build the goredir image
build-goredir:
	docker build . \
		-f build/goredirDockerfile \
		-t ${REDIR_IMAGE}

###
# Push images
###

push: build-operator push-operator

push-operator: build-operator
	docker push ${OPERATOR_IMAGE}

push-emulator: build-emulator
	docker push ${EMULATOR_IMAGE}

push-adbmon: build-adbmon
	docker push ${ADB_IMAGE}

push-goredir: build-goredir
	docker push ${REDIR_IMAGE}

push-all: push-operator push-emulator push-adbmon push-goredir

###
# Codegen
###

# Ensures a local copy of the operator-sdk
${OPERATOR_SDK}:
	mkdir -p _bin
	curl -JL -o ${OPERATOR_SDK} ${OPERATOR_SDK_URL}
	chmod +x ${OPERATOR_SDK}

# Generates deep copy code
generate: ${OPERATOR_SDK}
	GOROOT=${GOROOT} ${OPERATOR_SDK} generate k8s --verbose

# Generates CRD manifest
manifests: ${OPERATOR_SDK}
	${OPERATOR_SDK} generate crds --verbose

###
# Linting
###

${GOLANGCI_LINT}:
	mkdir -p $(dir ${GOLANGCI_LINT})
	cd $(dir ${GOLANGCI_LINT}) && curl -JL ${GOLANGCI_DOWNLOAD_URL} | tar xzf -
	chmod +x $(dir ${GOLANGCI_LINT})golangci-lint-${GOLANGCI_VERSION}-$(shell uname | tr A-Z a-z)-amd64/golangci-lint
	ln -s golangci-lint-${GOLANGCI_VERSION}-$(shell uname | tr A-Z a-z)-amd64/golangci-lint ${GOLANGCI_LINT}

# Lint files
lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} run -v --timeout 300s

# Tests
test:
	echo 'no tests yet, but needed'

###
# Kind helpers for local testing
###

# Ensures a repo-local installation of kind
${KIND}:
	mkdir -p $(dir ${KIND})
	curl -JL -o ${KIND} ${KIND_DOWNLOAD_URL}
	chmod +x ${KIND}

# Make a local test cluster and load a pre-baked emulator image into it
test-cluster: ${KIND}
	echo "$$KIND_CLUSTER_MANIFEST" | ${KIND} create cluster --config - --image kindest/node:${KUBERNETES_VERSION}
	$(MAKE) test-ingress

# Loads the operator image into the local kind cluster
load: load-operator

load-operator: ${KIND} build
	${KIND} load docker-image ${OPERATOR_IMAGE}

# Loads the emulator into the kind cluster
load-emulator: ${KIND}
	${KIND} load docker-image ${EMULATOR_IMAGE}

# Load adbmon image
load-adbmon: ${KIND}
	${KIND} load docker-image ${ADB_IMAGE}

load-goredir:
	${KIND} load docker-image ${REDIR_IMAGE}

load-all: load-operator load-emulator load-adbmon load-goredir

# Deploys metallb load balancer to the kind cluster
test-ingress:
	kubectl --context kind-kind apply -f https://raw.githubusercontent.com/google/metallb/${METALLB_VERSION}/manifests/namespace.yaml
	kubectl --context kind-kind apply -f https://raw.githubusercontent.com/google/metallb/${METALLB_VERSION}/manifests/metallb.yaml
	kubectl --context kind-kind create secret generic -n metallb-system memberlist --from-literal=secretkey="`openssl rand -base64 128`" || echo
	echo "$$METALLB_CONFIG" | kubectl --context kind-kind apply -f -

# # Deploys a stand alone traefik into the cluster
# traefik:
# 	helm repo add traefik https://containous.github.io/traefik-helm-chart
# 	helm repo update
# 	helm upgrade --install --namespace ingress traefik traefik/traefik

# Builds and deploys the operator into a local kind cluster, requires helm.
.PHONY: deploy
deploy:
	helm install --kube-context kind-kind android-farm-operator deploy/charts/android-farm-operator ${HELM_ARGS} --wait

example-farm:
	kubectl apply --context kind-kind -f deploy/examples/example-config.yaml -f deploy/examples/example-farm.yaml

## Doc generation

${REFDOCS_CLONE}:
	mkdir -p $(dir ${REFDOCS})
	git clone https://github.com/ahmetb/gen-crd-api-reference-docs "${REFDOCS_CLONE}"

${REFDOCS}: ${REFDOCS_CLONE}
	cd "${REFDOCS_CLONE}" && go build .
	mv "${REFDOCS_CLONE}/gen-crd-api-reference-docs" "${REFDOCS}"

api-docs: ${REFDOCS}
	go mod vendor
	bash hack/update-api-docs.sh

# just do everything
all-of-it: lint generate manifests api-docs build-all test-cluster load-all deploy example-farm
