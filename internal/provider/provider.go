// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-vellum/internal/provider/document_index"
	"terraform-provider-vellum/internal/provider/ml_model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellumclient "terraform-provider-vellum/internal/sdk/client"
)

// Ensure VellumProvider satisfies various provider interfaces.
var _ provider.Provider = &VellumProvider{}
var _ provider.ProviderWithFunctions = &VellumProvider{}

// VellumProvider defines the provider implementation.
type VellumProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// VellumProviderModel describes the provider data model.
type VellumProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseUrl types.String `tfsdk:"base_url"`
}

func (p *VellumProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vellum"
	resp.Version = p.version
}

func (p *VellumProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key to authenticate with the Vellum API",
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL to use with the Vellum API",
				Optional:            true,
			},
		},
	}
}

func (p *VellumProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data VellumProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl := os.Getenv("VELLUM_BASE_URL")
	if baseUrl == "" {
		baseUrl = data.BaseUrl.ValueString()
	}

	client := vellumclient.NewClient(
		vellumclient.WithApiKeyAndBaseUrl(
			os.Getenv("VELLUM_API_KEY"),
			baseUrl,
		),
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *VellumProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		document_index.Resource,
		ml_model.Resource,
	}
}

func (p *VellumProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		document_index.DataSource,
		ml_model.DataSource,
	}
}

func (p *VellumProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VellumProvider{
			version: version,
		}
	}
}
