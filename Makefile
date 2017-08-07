.DEFAULT_GOAL := build

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=notify-irc
LD_FLAGS += -s -w

release: clean fetch
	$(GOPATH)/bin/goreleaser --skip-publish

publish: clean fetch
	$(GOPATH)/bin/goreleaser

snapshot: clean fetch
	$(GOPATH)/bin/goreleaser --snapshot --skip-validate --skip-publish

update-deps: fetch
	$(GOPATH)/bin/govendor add +external
	$(GOPATH)/bin/govendor remove +unused
	$(GOPATH)/bin/govendor update +external

fetch:
	test -f $(GOPATH)/bin/govendor || go get -u -v github.com/kardianos/govendor
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser
	$(GOPATH)/bin/govendor sync

clean:
	/bin/rm -rfv "dist/" "${BINARY}"

build: fetch
	go build -ldflags "${LD_FLAGS}" -i -v -o ${BINARY}
