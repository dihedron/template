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

#
# The following targets are used to create a new release of the application;
# they extract the latest tag (in the format vX.Y.Z), increment the major,
# minor or patch version number, and create a new tag with the new version.
#
.PHONY: new-major-release
new-major-release: ## create a new major release (e.g. v1.2.3 -> v2.0.0)
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0 || echo "v0.0.0"))
	$(eval NEW_VERSION=$(shell echo $(OLD_VERSION) | awk -F. '{print $$1"."$$2"."$$3}' | awk -F. '{print "v"$$1+1".0.0"}' | sed 's/vv/v/g'))
	@echo "New major release: $(OLD_VERSION) -> $(NEW_VERSION)"
	@git tag -a $(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	@git push origin tag $(NEW_VERSION)

.PHONY: new-minor-release
new-minor-release: ## create a new minor release (e.g. v1.2.3 -> v1.3.0)
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0 || echo "v0.0.0"))
	$(eval NEW_VERSION=$(shell echo $(OLD_VERSION) | awk -F. '{print $$1"."$$2"."$$3}' | awk -F. '{print "v"$$1"."$$2+1".0"}' | sed 's/vv/v/g'))
	@echo "New minor release: $(OLD_VERSION) -> $(NEW_VERSION)"
	@git tag -a $(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	@git push origin tag $(NEW_VERSION)

.PHONY: new-revision-release
new-revision-release: ## create a new revision release (e.g. v1.2.3 -> v1.2.4)
	$(eval OLD_VERSION=$(shell git describe --tags --abbrev=0 || echo "v0.0.0"))
	$(eval NEW_VERSION=$(shell echo $(OLD_VERSION) | awk -F. '{print $$1"."$$2"."$$3}' | awk -F. '{print "v"$$1"."$$2"."$$3+1}' | sed 's/vv/v/g'))
	@echo "New revision release: $(OLD_VERSION) -> $(NEW_VERSION)"
	@git tag -a $(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	@git push origin tag $(NEW_VERSION)

.PHONY: test-json
test-json: compile ## test against a JSON file
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.json

.PHONY: test-yaml
test-yaml: compile ## test against a YAML file
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.yaml
