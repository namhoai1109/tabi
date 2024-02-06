depends: ## Install & build dependencies
	go get ./...
	go build ./...
	go mod tidy

mod.clean:
	go clean -cache
	go clean -modcache

mod: ## Update dependencies
	go mod tidy && go mod vendor
