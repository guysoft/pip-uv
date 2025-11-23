.PHONY: build install clean

BINARY_NAME=pip
BUILD_DIR=bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) src/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	rm -rf $(BUILD_DIR)

# Example helper to install to a specific venv (usage: make install VENV_PATH=../my-venv)
install: build
	@if [ -z "$(VENV_PATH)" ]; then \
		echo "Error: VENV_PATH is not set. Usage: make install VENV_PATH=/path/to/venv"; \
		exit 1; \
	fi
	@if [ ! -d "$(VENV_PATH)/bin" ]; then \
		echo "Error: $(VENV_PATH)/bin does not exist. Is this a valid venv?"; \
		exit 1; \
	fi
	cp $(BUILD_DIR)/$(BINARY_NAME) $(VENV_PATH)/bin/pip
	@echo "Installed shim to $(VENV_PATH)/bin/pip"

