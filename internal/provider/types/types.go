package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VellumProviderModel describes the provider data model.
type VellumProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}