WHAT := todo

PWD ?= $(shell pwd)

VERSION   ?= $(shell git describe --tags)
REVISION  ?= $(shell git rev-parse HEAD)
BRANCH    ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDUSER ?= $(shell id -un)
BUILDTIME ?= $(shell date '+%Y%m%d-%H:%M:%S')

.PHONY: build build-darwin-amd64 build-linux-amd64 build-windows-amd64

build: 
	for target in $(WHAT); do \
		go build -ldflags "-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Version=${VERSION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Revision=${REVISION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Branch=${BRANCH} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildUser=${BUILDUSER} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target ./cmd/$$target; \
	done

build-darwin-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags "-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Version=${VERSION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Revision=${REVISION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Branch=${BRANCH} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildUser=${BUILDUSER} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/todo-${VERSION}-darwin-amd64/$$target ./cmd/$$target; \
	done

build-linux-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags "-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Version=${VERSION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Revision=${REVISION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Branch=${BRANCH} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildUser=${BUILDUSER} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/todo-${VERSION}-linux-amd64/$$target ./cmd/$$target; \
	done

build-windows-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -a -installsuffix cgo -ldflags "-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Version=${VERSION} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.Revision=${REVISION} \
			-X github.com/yelimot/fullstack-todo-app-backendpkg/version.Branch=${BRANCH} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildUser=${BUILDUSER} \
			-X github.com/yelimot/fullstack-todo-app-backend/pkg/version.BuildDate=${BUILDTIME}" \
			-o ./bin/todo-${VERSION}-windows-amd64/$$target.exe ./cmd/$$target; \
	done