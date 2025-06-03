vet:
	go vet -trimpath ./...
	staticcheck ./...

test: vet
	go test -trimpath -race ./...

release:
	bump_version minor main.go
	git push origin --tags
