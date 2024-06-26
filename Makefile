default: all

.PHONY: all
all: dep test build

GIT_REVISION=$(shell git describe --always --dirty)
GIT_BRANCH=$(shell git rev-parse --symbolic-full-name --abbrev-ref HEAD)

LDFLAGS=-ldflags "-s -X main.gitRevision=$(GIT_REVISION) -X main.gitBranch=$(GIT_BRANCH)"

.PHONY: clean
clean:
	[ -d dist ] || mkdir dist
	rm dist/* || true

.PHONY: dep
dep:
	go mod tidy

.PHONY: checkdep
checkdep:
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build: clean dep
	go build $(LDFLAGS) -o dist/ ./

.PHONY: gox
gox: clean dep
	gox --osarch="linux/amd64 windows/amd64" -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}" $(LDFLAGS) ./