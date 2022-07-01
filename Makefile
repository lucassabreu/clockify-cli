export GO111MODULE=on
MAIN_PKG=./cmd/clockify-cli

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## clean all buildable files
	rm -rf dist

deps-install: ## install golang dependencies
	go mod download

deps-upgrade: ## upgrade go dependencies
	go get -u -v $(MAIN_PKG)
	go mod tidy

build: dist

dist: deps-install dist/darwin dist/linux dist/windows ## build all cli versions (default)

dist-internal:
	mkdir -p dist/$(goos)
	GOOS=$(goos) GOARCH=$(goarch) go build -o dist/darwin/clockify-cli $(MAIN_PKG)

dist/darwin:
	make dist-internal goos=darwin goarch=amd64

dist/linux:
	make dist-internal goos=linux goarch=amd64

dist/windows:
	make dist-internal goos=windows goarch=amd64

go-install: deps-install ## install dev version
	go install $(MAIN_PKG)

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
	curl -X POST -d '{"trigger_branch":"$(tag)","trigger_title":"Releasing $(tag)"}' https://api.netlify.com/build_hooks/5eef4f99028bddbb4093e4c6 -v
endif

site/themes/hugo-theme-learn/.git:
	git submodule update --init

site-build: site/themes/hugo-theme-learn/.git ## generates command documents and builds the site
	./scripts/site-build

site-serve: site-build ## builds the site, and serves it locally
	cd site && hugo serve
