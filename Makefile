CMD_DIR := ./cmd
TEST_DIR := ./tests

# Runs the development server
.PHONY: run
run:
	@echo "Starting development server..."
	@GO_ENV=development go run $(CMD_DIR)

# Runs tests in the development environment
.PHONY: test
test:
	@echo "Running tests with GO_ENV=testing..."
	@GO_ENV=testing go test -v $(TEST_DIR)/...

# Runs tests in the staging environment
.PHONY: stage
stage:
	@echo "Running tests with GO_ENV=staging..."
	@GO_ENV=staging go test -v $(TEST_DIR)/...

# Installs the packages
.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

