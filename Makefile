GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --always --tags)
BUILD_NODE_PAR = -ldflags "-X github.com/ontio/mixbee/common/config.Version=$(VERSION)" #-race

ARCH=$(shell uname -m)
DBUILD=docker build
DRUN=docker run
DOCKER_NS ?= ontio
DOCKER_TAG=$(ARCH)-$(VERSION)
ONT_CFG_IN_DOCKER=config.json
WALLET_FILE=wallet.dat

SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)
TOOLS=./tools
ABI=$(TOOLS)/abi
NATIVE_ABI_SCRIPT=./cmd/abi/native_abi_script

mixbee: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o mixbee main.go
 
sigsvr: $(SRC_FILES) abi 
	$(GC)  $(BUILD_NODE_PAR) -o sigsvr sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr $(TOOLS)

abi: 
	@if [ ! -d $(ABI) ];then mkdir -p $(ABI) ;fi
	@cp $(NATIVE_ABI_SCRIPT)/*.json $(ABI)

tools: sigsvr abi

all: mixbee tools

mixbee-cross: mixbee-windows mixbee-linux mixbee-darwin

mixbee-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o mixbee-windows-amd64.exe main.go

mixbee-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o mixbee-linux-amd64 main.go

mixbee-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o mixbee-darwin-amd64 main.go

tools-cross: tools-windows tools-linux tools-darwin

tools-windows: abi 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-windows-amd64.exe sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-windows-amd64.exe $(TOOLS)

tools-linux: abi 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-linux-amd64 sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-linux-amd64 $(TOOLS)

tools-darwin: abi 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-darwin-amd64 sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-darwin-amd64 $(TOOLS)

all-cross: mixbee-cross tools-cross abi

format:
	$(GOFMT) -w main.go

$(WALLET_FILE):
	@if [ ! -e $(WALLET_FILE) ]; then $(error Please create wallet file first) ; fi

docker/payload: docker/build/bin/mixbee docker/Dockerfile $(ONT_CFG_IN_DOCKER) $(WALLET_FILE)
	@echo "Building mixbee payload"
	@mkdir -p $@
	@cp docker/Dockerfile $@
	@cp docker/build/bin/mixbee $@
	@cp -f $(ONT_CFG_IN_DOCKER) $@/config.json
	@cp -f $(WALLET_FILE) $@
	@tar czf $@/config.tgz -C $@ config.json $(WALLET_FILE)
	@touch $@

docker/build/bin/%: Makefile
	@echo "Building mixbee in docker"
	@mkdir -p docker/build/bin docker/build/pkg
	@$(DRUN) --rm \
		-v $(abspath docker/build/bin):/go/bin \
		-v $(abspath docker/build/pkg):/go/pkg \
		-v $(GOPATH)/src:/go/src \
		-w /go/src/github.com/ontio/mixbee \
		golang:1.9.5-stretch \
		$(GC)  $(BUILD_NODE_PAR) -o docker/build/bin/mixbee main.go
	@touch $@

docker: Makefile docker/payload docker/Dockerfile 
	@echo "Building mixbee docker"
	@$(DBUILD) -t $(DOCKER_NS)/mixbee docker/payload
	@docker tag $(DOCKER_NS)/mixbee $(DOCKER_NS)/mixbee:$(DOCKER_TAG)
	@touch $@

clean:
	rm -rf *.8 *.o *.out *.6 *exe
	rm -rf mixbee tools docker/payload docker/build

