## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi

# Helper functions
STRIP_V = $(patsubst v%,%,$(1))

## Tool Binaries
TB_GORELEASER ?= $(TB_LOCALBIN)/goreleaser
TB_SEMVER ?= $(TB_LOCALBIN)/semver

## Tool Versions
# renovate: packageName=github.com/goreleaser/goreleaser/v2
TB_GORELEASER_VERSION ?= v2.16.0
TB_GORELEASER_VERSION_NUM ?= $(call STRIP_V,$(TB_GORELEASER_VERSION))
# renovate: packageName=github.com/bakito/semver
TB_SEMVER_VERSION ?= v1.1.10
TB_SEMVER_VERSION_NUM ?= $(call STRIP_V,$(TB_SEMVER_VERSION))

## Tool Installer
.PHONY: tb.goreleaser
tb.goreleaser: ## Download goreleaser locally if necessary.
	@test -s $(TB_GORELEASER) && $(TB_GORELEASER) --version | grep -q $(TB_GORELEASER_VERSION_NUM) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/goreleaser/goreleaser/v2@$(TB_GORELEASER_VERSION)
.PHONY: tb.semver
tb.semver: ## Download semver locally if necessary.
	@test -s $(TB_SEMVER) && $(TB_SEMVER) -version | grep -q $(TB_SEMVER_VERSION_NUM) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/bakito/semver@$(TB_SEMVER_VERSION)

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_GORELEASER) \
		$(TB_SEMVER)

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile --renovate -f $(TB_LOCALDIR)/Makefile \
		github.com/goreleaser/goreleaser/v2?--version \
		github.com/bakito/semver?-version
## toolbox - end
