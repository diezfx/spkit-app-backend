GOLANG_CI_VERSION=v1.54.0

generate:
	go generate ./...
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./internal/ent/schema --target gen/ent

fmt:
	gofumpt -l -w .

test:
	go test ./... -v

lint: 
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin $(GOLANG_CI_VERSION) 
	./.bin/golangci-lint run -c ./.golangci.yml


docker/build:
	docker build . -t ghcr.io/diezfx/split-backend-app:latest -f "deployment/Dockerfile" --build-arg="APP_NAME=split-app-backend"


