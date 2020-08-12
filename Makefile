# Task runner

.PHONY: help build

.DEFAULT_GOAL := help
.ONESHELL:

SHELL := /bin/bash

APP_VERSION  := $(shell git show-ref --dereference --tags | grep ^`git rev-parse HEAD` | sed -e 's,.* refs/tags/,,' -e 's/\^{}//')
GIT_HASH     := $(shell git rev-parse --short HEAD)

ANSI_TITLE        := '\e[1;32m'
ANSI_CMD          := '\e[0;32m'
ANSI_TITLE        := '\e[0;33m'
ANSI_SUBTITLE     := '\e[0;37m'
ANSI_WARNING      := '\e[1;31m'
ANSI_OFF          := '\e[0m'

TIMESTAMP := $(shell date "+%s")

help: ## Show this menu
	@echo -e $(ANSI_TITLE)gopro$(ANSI_OFF)$(ANSI_SUBTITLE)" - Libraries, Daemons and Controllers for GoPro on Linux"$(ANSI_OFF)
	@echo -e "\nUsage: $ make \$${COMMAND} \n"
	@echo -e "Variables use the \$${VARIABLE} syntax, and are supplied as environment variables before the command. For example, \n"
	@echo -e "  \$$ VARIABLE="foo" make help\n"
	@echo -e $(ANSI_TITLE)Commands:$(ANSI_OFF)
	@grep -E '^[.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: verify-requirements
verify-requirements-%: ## Verifies the system contains the required binaries
	@which $* > /dev/null || echo -e "\n  You need '$*'.\n"
	@which $* > /dev/null || exit 1

clean: ## Clean up the directory
	rm -rf dist

app: verify-requirements-go ## Builds the application itself
	go build -o dist/fs/usr/bin/goprod cmd/goprod/main.go
	go build -o dist/fs/usr/bin/goproctl cmd/goproctl/main.go

deb: verify-requirements-fpm ## Generate the debian package
	mkdir -p dist

	# Find the package name
	[ ! -z $(APP_VERSION) ] && export V=$(APP_VERSION) || export V="0.0.0-$(GIT_HASH)"

	fpm --input-type dir \
		--output-type deb \
		--name "gopro" \
		--version $$V \
		--package dist/gopro.$$V.deb \
		dist/fs/=/

deb.install: ## Install the debian package
	# Find the package name
	[ ! -z $(APP_VERSION) ] && export V=$(APP_VERSION) || export V="0.0.0-$(GIT_HASH)"

	dpkg -i dist/gopro.$$V.deb
