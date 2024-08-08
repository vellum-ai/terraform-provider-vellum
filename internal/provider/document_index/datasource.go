package document_index

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellumclient "terraform-provider-vellum/internal/sdk/client"
)

func DataSource() datasource.DataSource {
	return &DocumentIndexDataSource{}
}

type DocumentIndexDataSource struct {
	client *vellumclient.Client
}

var _ datasource.DataSource = &DocumentIndexDataSource{}
var _ datasource.DataSourceWithConfigure = &DocumentIndexDataSource{}

type TfDocumentIndexDataSourceModel struct {
	Created     types.String `tfsdk:"created"`
	Environment types.String `tfsdk:"environment"`
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
}

func (d *DocumentIndexDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document_index"
}

func (d *DocumentIndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Document Index data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "The Document Index's ID",
				MarkdownDescription: "The Document Index's ID",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "A name that uniquely identifies this index within its workspace",
				MarkdownDescription: "A name that uniquely identifies this index within its workspace",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"environment": schema.StringAttribute{
				Computed:            true,
				Description:         "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
				MarkdownDescription: "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
			},
			"label": schema.StringAttribute{
				Computed:            true,
				Description:         "A human-readable label for the document index",
				MarkdownDescription: "A human-readable label for the document index",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The current status of the document index\n\n* `ACTIVE` - Active\n* `ARCHIVED` - Archived",
				MarkdownDescription: "The current status of the document index\n\n* `ACTIVE` - Active\n* `ARCHIVED` - Archived",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ACTIVE",
						"ARCHIVED",
					),
				},
			},
		},
	}
}

func (d *DocumentIndexDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vellumclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DocumentIndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var documentIndexModel *TfDocumentIndexDataSourceModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &documentIndexModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentIndexRetrieveParameter := documentIndexModel.Name.ValueString()
	if documentIndexRetrieveParameter == "" {
		documentIndexRetrieveParameter = documentIndexModel.Id.ValueString()
	}
	if documentIndexRetrieveParameter == "" {
		resp.Diagnostics.AddError("failed to read Document Index", "either `id` or `name` must be set")
		return
	}

	documentIndex, err := d.client.DocumentIndexes.Retrieve(ctx, documentIndexRetrieveParameter)
	if err != nil {
		resp.Diagnostics.AddError("error getting Document Index information", err.Error())
		return
	}

	documentIndexModel, diagnostic := NewTfDocumentIndexDataSourceModel(ctx, documentIndex)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &documentIndexModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
