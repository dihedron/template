# Add custom targets below...

#
# compile is the default target; it builds the
# application for the default platform (linux/amd64)
#
.DEFAULT_GOAL := compile

.PHONY: compile
compile: go-dev ## build for the default linux/amd64 platform

.PHONY: snapshot
snapshot: go-snapshot ## build a snapshot version for the supported platforms

.PHONY: release
release: go-release ## build a release version (requires a valid tag)

.PHONY: clean
clean: ## clean the binary directory
	@rm -rf dist

.PHONY: test-json
test-json: compile ## test against a JSON file
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.json --log-enabled

.PHONY: test-yaml
test-yaml: compile ## test against a YAML file
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.yaml --log-enabled
