# Variables
APP_NAME := Appcrons
CMD_DIR := ./cmd
BIN_DIR := ./bin
TEST_DIR := ./tests

# Ensure bin directory exists
$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

# Default target
.PHONY: all
all: install run

# Run the development server
.PHONY: run
run:
	@echo "Starting development server..."
	@ @GO_ENV=development go run $(CMD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests with GO_ENV=testing..."
	@GO_ENV=testing go test -v $(TEST_DIR)/...


# Install dependencies
.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf $(BIN_DIR)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt $(CMD_DIR) $(TEST_DIR)

# Build the application
.PHONY: build
build: $(BIN_DIR)
	@echo "Building application..."
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

# Run application in the background
.PHONY: start
start: $(BIN_DIR)
	@echo "Starting application in the background..."
	@nohup $(BIN_DIR)/$(APP_NAME) &

# Stop the application
.PHONY: stop
stop:
	@echo "Stopping application..."
	@pkill -f "$(BIN_DIR)/$(APP_NAME)"
