CMD_DIR := ./cmd
TEST_DIR := ./tests

.PHONY: run
run:
	@echo "Starting development server..."
	@GO_ENV=development go run $(CMD_DIR)

.PHONY: test
test:
	@echo "Running tests with GO_ENV=testing..."
	@GO_ENV=testing go test -v $(TEST_DIR)/...


.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

