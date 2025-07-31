define create-install-tool-from-tar-gz
install-tools: install-tool-$(1)
.PHONY: install-tool-$(1)
install-tool-$(1): $$(BIN)/$(1)
$$(BIN)/$(1): | $$(BIN)
	# clean temp paths
	rm -rf $$(BIN)/.extract $$(BIN)/.archive.tar.gz && mkdir -p $$(BIN)/.extract
	# download $(1) archive
	curl -o $$(BIN)/.archive.tar.gz -fsSL $(2)
	# extract $(1)
	bsdtar xvzf $$(BIN)/.archive.tar.gz --strip-components $(3) -C $$(BIN)/.extract
	# move $(1)
	mv $$(BIN)/.extract/$(1) $$(BIN)/$(1)
	# clean temp paths
	rm -rf $$(BIN)/.extract $$(BIN)/.archive.tar.gz 
endef

define create-install-tool-from-zip
install-tools: install-tool-$(1)
.PHONY: install-tool-$(1)
install-tool-$(1): $$(BIN)/$(1)
$$(BIN)/$(1): | $$(BIN)
	# clean temp paths
	rm -rf $$(BIN)/.extract $$(BIN)/.archive.zip && mkdir -p $$(BIN)/.extract
	# download $(1) archive
	curl -o $$(BIN)/.archive.zip -fsSL $(2)
	# extract $(1)
	bsdtar xvf $$(BIN)/.archive.zip --strip-components $(3) -C $$(BIN)/.extract
	# move $(1)
	mv $$(BIN)/.extract/$(1) $$(BIN)/$(1)
	# clean temp paths
	rm -rf $$(BIN)/.extract $$(BIN)/.archive.zip 
endef