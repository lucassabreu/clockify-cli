export GO111MODULE=on

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## clean all buildable files
	rm -rf dist

install-deps: ## install golang dependencies
	go mod download

build: dist

dist: install-deps dist/darwin dist/linux dist/windows ## build all cli versions (default)

dist/darwin:
	mkdir -p dist/darwin
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/clockify-cli

dist/linux:
	mkdir -p dist/linux
	GOOS=linux GOARCH=amd64 go build -o dist/linux/clockify-cli

dist/windows:
	mkdir -p dist/windows
	GOOS=windows GOARCH=amd64 go build -o dist/windows/clockify-cli

go-install: ## install dev version
	go install

goreleaser-test: tag=Unreleased
goreleaser-test: release

ifeq ($(tag),Unreleased)
SNAPSHOT=1
endif
tag=
release: ## releases a tagged version
	sed "/^## \[$(tag)/, /^## \[/!d" CHANGELOG.md | tail -n +2 | head -n -2 > /tmp/rn.md
	curl -sL https://git.io/goreleaser | bash /dev/stdin --release-notes /tmp/rn.md \
		--rm-dist $(if $(SNAPSHOT),--snapshot --skip-publish,)
ifneq ($(SNAPSHOT),1)
	curl -X POST -d {} "https://api.netlify.com/build_hooks/5eef4f99028bddbb4093e4c6?trigger_branch=$(tag)&trigger_title=Releasing $(tag)"
endif
