# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

## Prerequisites

Create a new file called `.terraformrc` in your home directory (`~`). This allows you to run terraform against the locally built provider:

```
provider_installation {
  dev_overrides {
    "hashicorp.com/ai/vellum" = "/Users/<Username>/go/bin"
  }

  direct {}
}
```

Be sure to replace `<Username>` with your home username.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or ar testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named data source page
