# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

VERSION ?= $(shell bin/get_version.sh)

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