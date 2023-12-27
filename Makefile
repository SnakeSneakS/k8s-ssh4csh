ENVTEST_K8S_VERSION ?= 1.28.0

# gobin
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# container tool
CONTAINER_TOOL ?= docker

# shell tool
SHELL ?= /usr/bin/env bash -o pipefail
.SHELLFLAGS ?= -ec

# default
.PHONY: all
all: help

##@ General
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet envtest #manifests generate ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out

DEV_POD_TEMPLATE_FILE?=./manifest/dev/pod.template.yaml
DEV_POD_OUT_FILE?=./manifest/dev/pod.yaml
TMP_FILE?=./manifest/dev/tmp.pod.yaml
KEY_INPUT_CODE_DIR?=<INPUT_CODE_DIR>
KEY_INPUT_SSH_PUB_KEY?=<INPUT_SSH_PUB_KEY>
.PHONY: dev-pod-deploy
dev-pod-deploy: ## Deploy golang development environment
	test -f $(SSH_KEY_FILE) \
	&& echo "ssh key for dev already exists in $(SSH_KEY_FILE)" \
	|| make dev-ssh-keygen 
	cat $(DEV_POD_TEMPLATE_FILE) > $(TMP_FILE) && \
	cat $(TMP_FILE) | sed "s@$(KEY_INPUT_SSH_PUB_KEY)@`cat $(SSH_KEY_FILE).pub`@g" > $(TMP_FILE) && \
	cat $(TMP_FILE) | sed "s@$(KEY_INPUT_CODE_DIR)@`pwd`@g" > $(DEV_POD_OUT_FILE) && \
	rm $(TMP_FILE)
	kubectl kustomize ./manifest/dev | kubectl apply -f -


SUBCMD?=-h
DEV_PODNAME?=dev-ssh-sidecar
DEV_CONTAINER_NAME?=dev-golang
.PHONY: dev-pod-run
dev-pod-run: ## Set $SUBCMD to change sub command. Run golang code in the deployed golang development environmen_ (e.g. SUBCMD="echo-containers" make dev-pod-run )
#&&を使うとKUBERNETES_SERVICE_HOSTなどの環境変数が消えるためNG 
	kubectl exec $(DEV_PODNAME) -c $(DEV_CONTAINER_NAME) -- go mod download;
	kubectl exec $(DEV_PODNAME) -c $(DEV_CONTAINER_NAME) -it -- go run main.go $(SUBCMD) 
# どうやったら最後のgo run main.go にkill signalを送ることができるのだろう?


SSH_KEY_FILE?=./out/dev-sshkey
.PHONY: dev-ssh-keygen
dev-ssh-keygen: ## Generate ssh key for development environment
	ssh-keygen -t ed25519 -C "ssh key for development usage" -f $(SSH_KEY_FILE)
	echo "ssh key is generated in $(SSH_KEY_FILE) & $(SSH_KEY_FILE).pub"

DEV_SSH_CONTAINER_IP?=
DEV_SSH_CONTAINER_PORT?=22
DEV_SVC_NAME?=dev-ssh-sidecar
.PHONY: dev-ssh-connect 
dev-ssh-connect: ## Connect to dev container through ssh
	@echo $(DEV_SSH_CONTAINER_IP)
ifndef DEV_SSH_CONTAINER_IP
	@echo 'DEV_SSH_CONTAINER_IP is required, e.g.) DEV_SSH_CONTAINER_IP=127.0.0.1 make dev-ssh-connect'
	@echo 'you can find ip address by running \`kubectl get svc $(DEV_PODNAME)\`'
	exit 1
endif 
	ssh -i $(SSH_KEY_FILE) $(DEV_CONTAINER_NAME)@$(DEV_SSH_CONTAINER_IP) 

##@ Build Dependencies 

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY:tmp
tmp: ## tmp command for try makefile
	sleep 1 && echo 1
	sleep 1 && echo 2
	sleep 1 && echo 3
	sleep 1 && echo 4
	sleep 1 && echo 5
