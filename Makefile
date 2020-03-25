PROJECT_NAME := "dougkirkley/kube-deployer"
PKG := "github.com/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test lint

all: build

format: ## format the files
	@gofmt -w ${GO_FILES}
	
lint: ## Lint the files 
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests  
	@go test -short ${PKG_LIST}

dep: ## Get the dependencies  
	@go get -v ./...

build: dep ## Build the binary file 
	@go build -i -v $(PKG)

clean: ## Remove previous build 
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen 
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 
