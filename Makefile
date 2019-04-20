all: deps-init deps build

ifeq ($(OS),$(filter $(OS),Linux Darwin))
    LD_FLAGS = -ldflags="-s -w"
endif

.PHONY: deps-init
deps-init:
	rm go.mod go.sum
	@GO111MODULE=on go mod init
	@GO111MODULE=on go mod tidy

.PHONY: deps
deps:
@GO111MODULE=on go mod download

.PHONY: build
build:
	GO111MODULE=on go build $(LD_FLAGS) -o gotr
