include ./Makefile.includes.mk 

PWD := $(shell pwd)
OS := $(shell uname)
ARCH := $(shell uname -m)

BIN ?= $(PWD)/.bin
DST ?= $(PWD)/.dst

BUN_VERSION := 1.2.19
GORELEASER_VERSION := 2.11.1
SQLC_VERSION := 1.29.0

.PHONY: default
default: list-targets

.PHONY: list-targets
list-targets:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: clean
clean:
	# remove dst folder
	rm -rf $(DST)

.PHONY: create-migration
create-migration:
	# create migration with NAME build setting
	go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir $(PWD)/pkg/db/migrations -seq $(NAME)

.PHONY: generate-db-files
generate-db-files:
	# delete generated db files
	rm -rf pkg/db/sqlc
	# generate db files
	sqlc generate --file $(PWD)/.sqlc.yaml

.PHONY: install-tools
install-tools:

BUN_ARCH := $(ARCH)
ifeq ($(BUN_ARCH),x86_64)
	BUN_ARCH := x64
endif
BUN_URL := https://github.com/oven-sh/bun/releases/download/bun-v$(BUN_VERSION)/bun-$(OS)-$(BUN_ARCH).zip
$(eval $(call create-install-tool-from-zip,bun,$(BUN_URL),1))

GORELEASER_ARCH := $(ARCH)
ifeq ($(GORELEASER_ARCH),aarch64)
	GORELEASER_ARCH := arm64
endif
GORELEASER_URL := https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_$(OS)_$(GORELEASER_ARCH).tar.gz
$(eval $(call create-install-tool-from-tar-gz,goreleaser,$(GORELEASER_URL),0))

SQLC_ARCH := $(ARCH)
ifeq ($(SQLC_ARCH),x86_64)
	SQLC_ARCH := amd64
endif
SQLC_URL := https://github.com/sqlc-dev/sqlc/releases/download/v$(SQLC_VERSION)/sqlc_$(SQLC_VERSION)_$(OS)_$(SQLC_ARCH).tar.gz
$(eval $(call create-install-tool-from-tar-gz,sqlc,$(SQLC_URL),0))

$(BIN):
	# make bin folder
	mkdir -p $(BIN)

$(DST):
	# make dst folder
	mkdir -p $(DST)
