aperture_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

generate-api:
	@echo Generating API
	@cd api && $(MAKE) generate

go-generate:
	@echo Generating go code
	@./scripts/go_generate.sh

go-mod-tidy:
	@echo Download go.mod dependencies
	@GOPRIVATE=github.com/fluxninja,github.com/aperture-control
	@go mod tidy

go-test:
	@echo Running go tests
	@{ \
		export KUBEBUILDER_ASSETS=$(shell make operator-setup_envtest -s); \
		gotestsum --format=pkgname; \
	}

go-lint:
	@echo Linting go code
	@./scripts/golangci_lint.sh --fix

go-build:
	@echo Building go code
	@./scripts/go_build.sh

go-build-plugins:
	@echo Building go plugins
	@./scripts/go_build_plugins.sh

install-asdf-tools:
	@echo Installing Asdf tools
	@./scripts/manage_tools.sh setup

install-go-tools:
	@echo Installing tools from tools.go
	@./scripts/install_go_tools.sh

install-python-tools:
	@echo Installing tools from tools.py
	@pip install -r requirements.txt

go-generate-swagger:
	@echo Generating swagger code
	@echo Generating swagger specs from go code
	@./scripts/go_generate_swagger.sh

generate-docs: generate-blueprints generate-helm-readme generate-doc-assets
	@echo Generating docs

generate-config-markdown: go-generate-swagger generate-api
	@cd ./docs && $(MAKE) generate-config-markdown

generate-helm-readme:
	@echo Generating helm readme
	@cd ./manifests/charts && $(MAKE) generate-helm-readme

helm-lint:
	@echo helm lint
	@cd ./manifests/charts && $(MAKE) helm-lint

generate-blueprints: generate-config-markdown
	@cd ./blueprints && $(MAKE) generate-blueprints

generate-doc-assets:
	@cd ./docs && $(MAKE) generate-jsonnet
	@cd ./docs && $(MAKE) generate-mermaid

coverage_profile:
	gotestsum --format=testname -- -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage_in_browser: profile.coverprofile
	go tool cover -html profile.coverprofile

all: install-asdf-tools install-go-tools generate-api go-generate go-mod-tidy go-lint go-build go-build-plugins go-test generate-docs generate-helm-readme generate-blueprints helm-lint
	@echo "Done"

.PHONY: install-asdf-tools install-go-tools generate-api go-generate go-generate-swagger go-mod-tidy generate-config-markdown generate-doc-assets generate-docs go-test go-lint go-build go-build-plugins coverage_profile show_coverage_in_browser generate-helm-readme helm-lint generate-blueprints

#####################################
###### OPERATOR section starts ######
#####################################

# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
IMAGE_TAG_BASE ?= fluxninja/aperture-operator

# USE_IMAGE_DIGESTS defines if images are resolved via tags or digests
# You can enable this value if you would like to use SHA Based Digests
# To enable set flag to true
USE_IMAGE_DIGESTS ?= false
ifeq ($(USE_IMAGE_DIGESTS), true)
	BUNDLE_GEN_FLAGS += --use-image-digests
endif

# Image URL to use all building/pushing image targets
IMG ?= aperture-operator:latest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.23

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: operator-all
operator-all: operator-build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: operator-help
operator-help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: operator-manifests
operator-manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd:ignoreUnexportedFields=true,allowDangerousTypes=true webhook paths="./operator/..." output:crd:artifacts:config=operator/config/crd/bases output:rbac:artifacts:config=operator/config/rbac output:webhook:artifacts:config=operator/config/webhook
	./operator/hack/create_policy_sample.sh

.PHONY: operator-generate
operator-generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="operator/hack/boilerplate.go.txt" paths="./pkg/..."
	$(CONTROLLER_GEN) object:headerFile="operator/hack/boilerplate.go.txt" paths="./operator/..."

.PHONY: operator-fmt
operator-fmt: ## Run go fmt against code.
	go fmt ./operator/...

.PHONY: operator-vet
operator-vet: ## Run go vet against code.
	go vet ./operator/...

.PHONY: operator-test
operator-test: operator-manifests operator-generate operator-fmt operator-vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./operator/... -coverprofile operator/cover.out

.PHONY: operator-setup_envtest
operator-setup_envtest: envtest ## Run tests.
	echo "$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)"

##@ Build

.PHONY: operator-build
operator-build: operator-generate operator-fmt operator-vet ## Build manager binary.
	go build -o bin/manager operator/main.go

.PHONY: operator-run
operator-run: operator-manifests operator-generate operator-fmt operator-vet ## Run a controller from your host.
	go run ./operator/main.go

.PHONY: operator-docker-build
operator-docker-build: operator-test ## Build docker image with the manager.
	docker build -t ${IMG} ./ -f operator/Dockerfile

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: operator-install
operator-install: operator-manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build operator/config/crd | kubectl apply -f -

.PHONY: operator-uninstall
operator-uninstall: operator-manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build operator/config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: operator-deploy
operator-deploy: operator-manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd operator/config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build operator/config/default | kubectl create -f -

.PHONY: operator-undeploy
operator-undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build operator/config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v4.5.7
CONTROLLER_TOOLS_VERSION ?= v0.9.2

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

#####################################
###### OPERATOR section ends ########
#####################################
