include ./Makefile.includes.mk 

PWD := $(shell pwd)
OS := $(shell uname)
ARCH := $(shell uname -m)

BIN ?= $(PWD)/.bin
DST ?= $(PWD)/.dst

BUN_VERSION := 1.2.19
GORELEASER_VERSION := 2.11.1

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
	# create table
	go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir $(PWD)/pkg/db/migrations -seq $(TABLE)

.PHONY: generate-db-files
generate-db-files:
	# generate db files
	go run github.com/sqlc-dev/sqlc/cmd/sqlc 

.PHONY: install-tools
install-tools:

BUN_URL := https://github.com/oven-sh/bun/releases/download/bun-v$(BUN_VERSION)/bun-$(OS)-$(ARCH).zip
$(eval $(call create-install-tool-from-zip,bun,$(BUN_URL),1))

GORELEASER_ARCH := $(ARCH)
ifeq ($(GORELEASER_ARCH),aarch64)
	GORELEASER_ARCH := arm64
endif
GORELEASER_URL := https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_$(OS)_$(GORELEASER_ARCH).tar.gz
$(eval $(call create-install-tool-from-tar-gz,goreleaser,$(GORELEASER_URL),0))

$(BIN):
	# make bin folder
	mkdir -p $(BIN)

$(DST):
	# make dst folder
	mkdir -p $(DST)
