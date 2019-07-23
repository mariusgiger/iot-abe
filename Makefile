PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# some build parameters
ifeq (${VERSION},)
  VERSION=$(shell git describe --tags 2>/dev/null)
endif
ifeq (${GITHASH},)
  GITHASH=$(shell git log -1 --format='%H')
endif
ifeq (${BUILDTIME},)
  BUILDTIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
endif

.PHONY: all install build test cover coverhtml lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}
	golangci-lint run

loc: ## displays the lines of code of this project
	gocloc --not-match-d="vendor|node_modules|lib|analysis" --exclude-ext="json" .

version: ## Show version
	@echo "  Version: ${VERSION}"
	@echo "  GitHash: ${GITHASH}"
	@echo "BuildTime: ${BUILDTIME}"

setup: ## installs development dependencies
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/derekparker/delve/cmd/dlv
	# solc install
	brew update
	brew upgrade
	brew tap ethereum/ethereum
	brew install solidity 
	# abigen install
	cd $GOPATH/src/github.com/ethereum/go-ethereum
	go install ./cmd/abigen
	# TODO ganache install
	dep ensure

install: ## Install dependencies
	dep ensure
	$(MAKE) fix-eth-bug 

install-ci: ## Install dependencies for cicd
	dep ensure -v -vendor-only

fix-eth-bug: ## Workaround for Ethereum client bug using dep
	@go get -v -u "github.com/ethereum/go-ethereum/crypto/secp256k1"
	@cp -r "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

contract-build: ## builds the smart contract
	solc --bin --abi --overwrite -o output contract/contracts/AccessControl.sol 

contract-abigen: ## generates go-bindings for the smart contract
	abigen --abi output/AccessControl.abi --bin output/AccessControl.bin --pkg contract --type AccessControl --out ./pkg/contract/contract.go

contract-test: ## runs smart contract tests
	$(MAKE) -C ./contract all

clean: ## cleans the repo
	rm -rf output/
	rm -rf vendor/
	rm -f cover.out
	docker rm -v $(docker ps --filter status=exited -q 2>/dev/null)
	docker rmi $(docker images --filter dangling=true -q 2>/dev/null)

build: ## Builds the cli
	CGO_ENABLED=1 go build -tags netgo -o output/iot-abe \
-ldflags "-X github.com/mariusgiger/iot-abe/cmd.Version=${VERSION} -X github.com/mariusgiger/iot-abe/cmd.GitHash=${GITHASH} -X github.com/mariusgiger/iot-abe/cmd.BuildTime=${BUILDTIME}" .

docker-build: ## Builds a docker image iot-abe
	docker build -t "mariusgiger/iot-abe:builder"  --target builder . && \
    docker build -t "mariusgiger/iot-abe:latest" .

docker-build-arm: ## Builds a docker image for iot-abe for ARM architecture
	docker build -f Dockerfile.arm -t "mariusgiger/iot-abe:arm-builder"  --target builder . && \
    docker build -f Dockerfile.arm  -t "mariusgiger/iot-abe:arm-latest" .

docker-publish-arm:
	docker push mariusgiger/iot-abe:arm-latest

docker-run-rpi: ## Runs iot-abe server on a raspberry
	docker run --name iot-abe --rm --privileged --device=/dev/vchiq -p8080:8080 mariusgiger/iot-abe:arm-latest server

# "-count=1" is used to avoid test result caching
test-unit: ## Do the unit tests
	go test -v -count=1 ${PKG_LIST}

# '-run NOTHING' is used to ignore all unit tests
test-perf: ## Do the performance tests
	go test -run NOTHING -bench=. -benchmem ${PKG_LIST}

test-integration: ## Do the integration tests
	go test -count=1 -run Remote ./...

cover: ## Generates a global code coverage report
	go test -coverprofile=cover.out ${PKG_LIST}
	go tool cover -func=cover.out

cover-html: ## Generates a global code coverage report in HTML
	go test -coverprofile=cover.out ${PKG_LIST}
	go tool cover -html=cover.out

analyze: ## Opens analysis env
	$(shell cd ./analysis && source bin/activate && jupyter notebook)

rpi-connect: ## Connects to a RPI
	ssh pi@192.168.1.54

help: ## Display this help screen
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
