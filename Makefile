# Include toolbox tasks
include ./.toolbox.mk

# Run go fmt against code
fmt:
	go fmt ./...
	gofmt -s -w .

# Run go vet against code
vet:
	go vet ./...

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tidy fmt vet
	go test ./...  -coverprofile=coverage.out
	go tool cover -func=coverage.out

release: tb.semver tb.goreleaser
	@version=$$($(TB_SEMVER)); \
	git tag -s $$version -m"Release $$version"; \
	git push origin $$version
	PATH=$(TB_LOCALBIN):$${PATH} $(TB_GORELEASER) --clean --parallelism 2

test-release: tb.goreleaser 
	PATH=$(TB_LOCALBIN):$${PATH} $(TB_GORELEASER) --skip=publish --snapshot --clean --parallelism 2
