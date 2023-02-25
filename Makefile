all: deps-init deps build

ifeq ($(OS),$(filter $(OS),Linux Darwin))
    LD_FLAGS = -ldflags="-s -w"
endif

.PHONY: deps-init
deps-init:
	rm go.mod go.sum
	go mod init
	go mod tidy

.PHONY: deps
deps:
	go mod download

.PHONY: build
build:
	go build $(LD_FLAGS) -o gotr
