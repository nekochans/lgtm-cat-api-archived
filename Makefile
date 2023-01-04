.PHONY: sqlc-generate
sqlc-generate:
	docker run --rm -v `pwd`:/src -w /src kjconroy/sqlc generate

.PHONY: lint
lint:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v

.PHONY: format
format:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v --fix

.PHONY: run-normal-build
run-normal-build:
	AIR_CONF=air.normal.toml docker compose up dev --build -d

.PHONY: run-normal
run-normal:
	AIR_CONF=air.normal.toml docker compose up dev -d

.PHONY: run-debug-build
run-debug-build:
	AIR_CONF=air.debug.toml docker compose up dev -d

.PHONY: run-debug
run-debug:
	AIR_CONF=air.debug.toml docker compose up dev --build -d

.PHONY: test-build
test-build:
	docker compose up test --build

.PHONY: test
test:
	docker compose up test
.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build
build: deps
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/lambda ./cmd/lambda/main.go

.PHONY: zip
zip: build
	zip -j bin/lambda.zip bin/lambda

.PHONY: deploy-stg
deploy-stg: zip
	aws lambda update-function-code \
		--region ap-northeast-1 \
		--function-name stg-lgtm-cat-api \
		--zip-file fileb://bin/lambda.zip \
		--profile lgtm-cat

.PHONY: deploy-prod
deploy-prod: zip
	aws lambda update-function-code \
		--region ap-northeast-1 \
		--function-name prod-lgtm-cat-api \
		--zip-file fileb://bin/lambda.zip \
		--profile lgtm-cat
