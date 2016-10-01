TARGETS = linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-386 windows-amd64
COMMAND_NAME = httpstat
PACKAGE_NAME = github.com/davecheney/$(COMMAND_NAME)
LDFLAGS = -ldflags=-X=main.version=$(VERSION)
OBJECTS = $(patsubst $(COMMAND_NAME)-windows-amd64%,$(COMMAND_NAME)-windows-amd64%.exe, $(patsubst $(COMMAND_NAME)-windows-386%,$(COMMAND_NAME)-windows-386%.exe, $(patsubst %,$(COMMAND_NAME)-%-v$(VERSION), $(TARGETS)))) 

release: check-env $(OBJECTS) ## Build release binaries (requires VERSION)

clean: check-env ## Remove release binaries
	rm $(OBJECTS)

$(OBJECTS): $(wildcard *.go)
	env GOOS=`echo $@ | cut -d'-' -f2` GOARCH=`echo $@ | cut -d'-' -f3 | cut -d'.' -f 1` go build -o $@ $(LDFLAGS) $(PACKAGE_NAME)

.PHONY: help check-env

check-env:
ifndef VERSION
	$(error VERSION is undefined)
endif

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
