.PHONY: lint test coverage

lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint -E gofmt run ./...

test:
	@type gosec > /dev/null || go get github.com/securego/gosec/cmd/gosec
	gosec ./...
	go test ./...

coverage:
	go test -covermode=count -coverprofile=coverage.txt ./...
	@type goveralls > /dev/null || go get -u github.com/mattn/goveralls
	goveralls -coverprofile=coverage.txt -service=travis-ci
