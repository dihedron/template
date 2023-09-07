default: build

.PHONY: build
build:
	goreleaser build --snapshot --single-target --rm-dist

.PHONY: clean
clean:
	@rm -rf dist/ 

.PHONY: test-json
test-json: build
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.json --log-enabled

.PHONY: test-yaml
test-yaml: build
	dist/template_linux_amd64_v1/template --template=_test/outer.tpl --template=_test/inner.tpl --input=@_test/input.yaml --log-enabled
