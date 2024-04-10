package terraform

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vellum-ai/terraform-provider-vellum/internal/terraform/documentindex"
	vellumclient "github.com/vellum-ai/terraform-provider-vellum/internal/vellum/client"
	"github.com/vellum-ai/terraform-provider-vellum/internal/vellum/option"
)

// ProviderModel represents the vellum terraform provider model.
// This is expressed in terraform with the following:
//
//	provider "vellum" {
//	  api_key = "YOUR_API_KEY"
//	}
type ProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// Provider represents the Vellum terraform provider. This type is extendable;
// simply override any of the exported methods in the hooks.go file to customize
// the provider.
type Provider struct {
	base *baseProvider

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Compile-time assertion that ensures the provider satisfies the provider.Provider
// interface.
var _ provider.Provider = (*Provider)(nil)

// NewProvider returns a new Vellum terraform provider.
func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			base:    newBaseProvider(version),
			version: version,
		}
	}
}

// baseProvider implements the base functionality of the provider.
//
// Do NOT edit this type directly; use hooks.go instead.
type baseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

var _ provider.Provider = (*baseProvider)(nil)

func newBaseProvider(version string) *baseProvider {
	return &baseProvider{
		version: version,
	}
}

func (b *baseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vellum"
	resp.Version = b.version
}

func (b *baseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key to authenticate with the Vellum API",
				Optional:            true,
			},
		},
	}
}

func (b *baseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("VELLUM_API_KEY")
	if !model.ApiKey.IsNull() {
		apiKey = model.ApiKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"An API key is required to use the vellum provider",
			"You must set a VELLUM_API_KEY or specify an api_key in the provider constructor",
		)
		return
	}

	client := vellumclient.NewClient(
		option.WithApiKey(
			apiKey,
		),
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (b *baseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		documentindex.NewResource,
	}
}

func (b *baseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		documentindex.NewDataSource,
	}
}
