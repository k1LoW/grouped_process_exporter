PKG = github.com/k1LoW/grouped_process_exporter
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	SED = gsed
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	SED = sed
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test sec

test:
	go test ./... -coverprofile=coverage.out -covermode=count

sec:
	gosec ./...

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

prerelease:
	git pull origin master --tag
	ghch -w -N ${VER}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS
	git commit -m'Bump up version number'
	git tag ${VER}

release:
	git push origin master --tag
	goreleaser --rm-dist

.PHONY: default test
