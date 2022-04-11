test: ## Runs tests
	go test ./...
run:  ## Builds & Runs the application
	go build . && ./centurion
follow-logs: ## Follows the logs in the file
	tail -f logs
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
