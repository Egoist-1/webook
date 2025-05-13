.PHONY: docker
docker:
	@rm -f webook || true
	@docker rmi -f egoist/webook:v0.01 || true
	@set GOOS=linux set GOARCH=arm
	@go build -o webook .
	@docker build -t egoist/webook:v0.0.1 .
mock:
	@go generate ./...
	@go mod tidy