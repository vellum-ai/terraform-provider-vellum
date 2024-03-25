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

Then, ensure you have the following tools installed:

- [Go 1.21+](https://golang.org/doc/install) installed and configured.
- [Terraform v1.5+](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) installed locally.

Finally, run `go install` from the repository's root which builds the provider's binaries into the `~/go/bin` directory.

## Running

Within the `/examples/provider` directory, you should run:

```bash
terraform plan
```

This will run the example against local terraform state, specifying all the resources that will be built, modified, or deleted as a result.
