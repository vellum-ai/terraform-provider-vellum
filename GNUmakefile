default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


# Used to generate the provider spec from the OpenAPI spec.
generate-provider-spec:
	tfplugingen-openapi generate \
		--config internal/specs/generator_config.yml \
		--output internal/specs/provider_code_spec.json \
		internal/specs/openapi.yml

# Used to generate the provider code from the provider spec.
generate-provider-code:
	tfplugingen-framework generate all \
		--input internal/specs/provider_code_spec.json \
		--output internal/provider

# Used to generate both the provider spec and provider code.
generate-all:
	make generate-provider-spec
	make generate-provider-code

# Used to generate the boilerplate for a new resource that'll then be hand-written.
scaffold-resource:
	mkdir -p "./internal/provider/resource_$(name)" && \
	tfplugingen-framework scaffold resource --name $(name) \
		--output-dir internal/provider \
 		--output-file=resource_$(name)/$(name)_resource.go \
		--force