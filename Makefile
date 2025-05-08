.PHONY: docs
docs:
	@go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs; \
	tfplugindocs
