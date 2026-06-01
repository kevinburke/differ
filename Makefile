vet:
	go vet -trimpath ./...
	staticcheck ./...

test: vet
	go test -trimpath -race ./...

version ?= minor

.PHONY: release
release: test
	go run github.com/kevinburke/bump_version@latest --tag-prefix=v $(version) version.go
