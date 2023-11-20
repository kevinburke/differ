STATICCHECK := $(GOPATH)/bin/staticcheck

$(STATICCHECK):
	go get -u honnef.co/go/tools/cmd/staticcheck

vet: $(STATICCHECK)
	go vet ./...
	$(STATICCHECK) ./...

test: vet
	go test -trimpath -race ./...

release:
	bump_version minor main.go
	git push origin --tags
