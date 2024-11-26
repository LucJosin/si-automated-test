COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

# Run tests
test:
	go test ./... -v

# Run tests with coverage and generate HTML output
test-coverage:
	go test ./... -cover -coverprofile=$(COVERAGE_OUT)
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Test coverage report generated: $(COVERAGE_HTML)"

# Clean up coverage files
clean:
	@rm -rf $(COVERAGE_OUT) $(COVERAGE_HTML)
	@echo "Removed $(COVERAGE_OUT) and $(COVERAGE_HTML)"