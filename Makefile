.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build
build: deps
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/lambda main.go

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
