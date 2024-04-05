package terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// --- Override the following methods to customize the provider ---

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	p.base.Metadata(ctx, req, resp)
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	p.base.Schema(ctx, req, resp)
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.base.Configure(ctx, req, resp)
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return p.base.Resources(ctx)
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.base.DataSources(ctx)
}
