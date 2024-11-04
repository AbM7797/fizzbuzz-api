.PHONY: dev
dev:
	wrangler dev --env local

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.23.1 -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm .

.PHONY: deploy-staging
deploy-staging:
	wrangler publish --env staging

.PHONY: deploy-prod
deploy-prod:
	wrangler publish --env production