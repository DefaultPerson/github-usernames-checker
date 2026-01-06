BINARY_NAME=Github_Username_Checker

PLATFORMS := linux darwin windows
ARCHITECTURES := amd64 arm64

GO := go
GOFLAGS := -ldflags "-s -w"

BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/main user_exists_checker.go

run:
	go run user_exists_checker.go

$(shell mkdir -p $(BUILD_DIR))

all: build-all

build-all: $(addsuffix .build,$(foreach platform,$(PLATFORMS),$(foreach arch,$(ARCHITECTURES),$(platform)-$(arch))))

%.build:
	@echo "Building for $*"
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(word 1,$(subst -, ,$*))-$(word 2,$(subst -, ,$*))

clean:
	rm -rf $(BUILD_DIR)

# PHONY targets
.PHONY: all build-all clean $(addsuffix .build,$(foreach platform,$(PLATFORMS),$(foreach arch,$(ARCHITECTURES),$(platform)-$(arch))))
