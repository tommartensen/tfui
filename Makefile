.PHONY: test

build: ## Build binary
	go build -o tfui ./cmd

format: ## Auto-format the code to conform with common Go style
	go fmt ./...

lint: ## Run the linter to enforce best practices
	golangci-lint run

test: ## Run all tests
	go test ./...

docker-build: ## Build docker container
	docker build . -t tommartensen/tfui:0.0.1

docker-run: ## Run docker container
	docker run -v plans:/plans -p 8080:8080 tommartensen/tfui:0.0.1

helm-deploy:
	helm upgrade --install --wait --timeout 300s tfui deploy/chart/

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
