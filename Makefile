PROJECT = tubeguardian
CMD_DIR = ./cmd/$(PROJECT)
BIN_DIR = ./bin
CONFIG_DIR = ./configs
LOGS_DIR = ./logs

GOOS_LIST = linux darwin windows
GOARCH_LIST = amd64 arm64

PKG_NAME = $(PROJECT)-v1.0.0
PKG_FILE = $(PKG_NAME).tar.gz

.PHONY: all build clean package run

all: build

# Build all OS/ARCH combinations
build:
	@mkdir -p $(BIN_DIR)
	@for os in $(GOOS_LIST); do \
		for arch in $(GOARCH_LIST); do \
			echo "🚀 Building $(PROJECT)-$$os-$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -o $(BIN_DIR)/$(PROJECT)-$$os-$$arch $(CMD_DIR); \
		done \
	done

# Package everything
package: build
	@echo "📦 Creating package $(PKG_FILE)..."
	@mkdir -p release_temp
	@cp -r $(BIN_DIR) release_temp/
	@cp -r $(CONFIG_DIR) release_temp/
	@tar -czf $(PKG_FILE) -C release_temp .
	@rm -rf release_temp
	@echo "✅ Package created: $(PKG_FILE)"

# Run locally
run:
	go run $(CMD_DIR)

# Clean
clean:
	rm -rf $(BIN_DIR)
	rm -f $(PKG_FILE)
	rm -rf $(LOGS_DIR)
