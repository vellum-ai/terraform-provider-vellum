package documentindex

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vellum-ai/terraform-provider-vellum/internal/vellum"
	vellumclient "github.com/vellum-ai/terraform-provider-vellum/internal/vellum/client"
)

// DataSourceModel represents the document index data source model.
// This is expressed in terraform with the following:
//
//	data "vellum_document_index" "reference" {
//	  name = "reference"
//	}
type DataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Environment types.String `tfsdk:"environment"`
	Status      types.String `tfsdk:"status"`
	Created     types.String `tfsdk:"created"`
}

// DataSource represents the document index data source. This type is extendable;
// simply override any of the exported methods in the data_source_hooks.go file
// to customize the data source.
type DataSource struct {
	base *baseDataSource
}

// Compile-time assertion that ensures the provider satisfies the datasource.DataSource
// interface.
var _ datasource.DataSource = (*DataSource)(nil)

// NewDataSource returns a new document index data source.
func NewDataSource() datasource.DataSource {
	return &DataSource{
		base: newBaseDataSource(),
	}
}

// baseDataSource implements the base functionality of the data source.
//
// Do NOT edit this type directly; use data_source_hooks.go instead.
type baseDataSource struct {
	Vellum *vellumclient.Client
}

var _ datasource.DataSource = (*DataSource)(nil)

func newBaseDataSource() *baseDataSource {
	return &baseDataSource{}
}

func (b *baseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document_index"
}

func (b *baseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Document Index data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The Document Index's ID",
				MarkdownDescription: "The Document Index's ID",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "A name that uniquely identifies this index within its workspace",
				MarkdownDescription: "A name that uniquely identifies this index within its workspace",
			},
			"label": schema.StringAttribute{
				Computed:            true,
				Description:         "A human-readable label for the document index",
				MarkdownDescription: "A human-readable label for the document index",
			},
			"environment": schema.StringAttribute{
				Computed:            true,
				Description:         "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
				MarkdownDescription: "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The current status of the document index\n\n* `ACTIVE` - Active\n* `ARCHIVED` - Archived",
				MarkdownDescription: "The current status of the document index\n\n* `ACTIVE` - Active\n* `ARCHIVED` - Archived",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (b *baseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vellumclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *vellum.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	b.Vellum = client
}

func (b *baseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model *DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := b.Vellum.DocumentIndexes.Retrieve(
		ctx,
		b.modelToRetrieveRequest(model),
	)
	if err != nil {
		resp.Diagnostics.AddError("error getting document index information", err.Error())
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			b.retrieveResponseToModel(response),
		)...,
	)
}

func (b *baseDataSource) modelToRetrieveRequest(model *DataSourceModel) string {
	return model.Name.ValueString()
}

func (b *baseDataSource) retrieveResponseToModel(response *vellum.DocumentIndexRead) *DataSourceModel {
	return &DataSourceModel{
		Id:          types.StringValue(response.Id),
		Name:        types.StringValue(response.Name),
		Created:     types.StringValue(response.Created.Format(time.RFC3339)),
		Environment: types.StringValue(string(*response.Environment)),
		Label:       types.StringValue(response.Label),
		Status:      types.StringValue(string(*response.Status)),
	}
}
