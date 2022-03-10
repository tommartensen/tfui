.PHONY: test build

APP_NAME := tfui
PLATFORMS := darwin linux
ARCHS := amd64 arm64

build: ## Build binary for your OS and architecture
	go build -o ./build/${APP_NAME} ./cmd

clean: ## Clean up the build folder
	@-rm -r ./build

format: ## Auto-format the code to conform with common Go style
	go fmt ./...

lint: ## Run the linter to enforce best practices
	golangci-lint run

test: ## Run all tests
	go test ./...

release: clean go-release go-checksums ## Cross-compile binaries and create checksums

docker-build: ## Build docker container
	docker build . -t tommartensen/${APP_NAME}:0.0.1

docker-run: ## Run docker container
	docker run -v plans:/plans -p 8080:8080 tommartensen/${APP_NAME}:0.0.1

go-vendor: ## Get dependencies
	go mod vendor

go-release: go-vendor ## Cross-compile with native Go methods
	@for GOOS in ${PLATFORMS}; do \
	  for GOARCH in ${ARCHS}; do \
	  	export GOOS=$$GOOS; \
			export GOARCH=$$GOARCH; \
			go build -v -o ./build/${APP_NAME}-$$GOOS-$$GOARCH ./cmd; \
	  done \
	done

go-checksums: ## Create checksum file
	@-rm -f ./build/SHA256SUMS
	cd ./build/ && sha256sum * > SHA256SUMS

helm-deploy:
	helm upgrade --install --wait --timeout 300s ${APP_NAME} deploy/chart/

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
