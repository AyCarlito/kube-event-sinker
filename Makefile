# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

VERSION ?= $(shell bin/get_version.sh)

IMAGE_REGISTRY ?=
IMAGE_TAG_BASE ?= aycarlito/kube-event-sinker
IMG ?= $(IMAGE_TAG_BASE):$(VERSION)

ifneq ($(IMAGE_REGISTRY),)
IMG := $(IMAGE_REGISTRY)/$(IMG)
endif

CHART_REGISTRY ?=

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

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: fmt
fmt: goimports ## Run goimports against code.
	find . -name '*.go' -exec sed -i ' /^import/,/)/ { /^$$/ d } ' {} + ; \
	$(LOCALBIN)/goimports -local github.com/AyCarlito/kube-event-sinker -w -e -format-only .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: run
run: ## Run the application.
	go run . $(CMD) $(FLAGS)

.PHONY: docker-run
docker-run: ## Run the docker image.
	docker run --network host \
	--user $(shell id -u):$(shell id -g) \
	-v ~/.kube/config:/.kube/config \
	$(IMG) $(CMD) $(FLAGS) 

.PHONY: helm-install
helm-install: ## Install the helm chart.
	helm upgrade --install kube-event-sinker -n kube-event-sinker helm/kube-event-sinker-$(VERSION).tgz

.PHONY: helm-uninstall
helm-uninstall: ## Uninstall the helm chart.
	helm uninstall kube-event-sinker -n kube-event-sinker

##@ Build
clean:
	go clean -modcache

clean-all:
	go clean -cache

pre:
	go mod tidy

.PHONY: build
build: pre fmt vet ## Build binary.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/kube-event-sinker

.PHONY: docker-build 
docker-build: ## Build docker image.
	docker build --platform linux/amd64 -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image.
	docker push ${IMG}

.PHONY: helm
helm: ## Build helm chart.
	sed -i "s/tag:.*/tag: $(VERSION)/" ./helm/values.yaml
	sed -i "s/appVersion:.*/appVersion: \"$(VERSION)\"/" ./helm/Chart.yaml

	helm lint ./helm
	helm package ./helm --destination ./helm --version "$(VERSION)" --app-version "$(VERSION)"

.PHONY: helm-push
helm-push: ## Push helm chart.
	helm push helm/kube-event-sinker-$(VERSION).tgz ${CHART_REGISTRY}

##@ Release
.PHONY: generate-latest-tag
generate-latest-tag: ## Generates the latest tag.
	./bin/bump_tag.sh

.PHONY: create-release-branch 
create-release-branch: generate-latest-tag ## Creates a release branch.
	./bin/release.sh

.PHONY: create-release-notes
create-release-notes:  ## Creates release notes.
	./bin/generate_release_notes.sh

##@ Build Dependencies
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOIMPORTS ?= $(LOCALBIN)/goimports

## Tool Binaries Versions
GOIMPORTS_VERSION ?= v0.16.0

.PHONY: goimports
goimports: ## Download goimports locally if necessary.
	test -s $(GOIMPORTS)/ || GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)