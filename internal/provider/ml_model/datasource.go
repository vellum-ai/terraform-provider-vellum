package ml_model

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
	return &MLModelDataSource{}
}

type MLModelDataSource struct {
	client *vellumclient.Client
}

var _ datasource.DataSource = &MLModelDataSource{}
var _ datasource.DataSourceWithConfigure = &MLModelDataSource{}

type TfMLModelDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Visibility  types.String `tfsdk:"visibility"`
	HostedBy    types.String `tfsdk:"hosted_by"`
	DevelopedBy types.String `tfsdk:"developed_by"`
	Family      types.String `tfsdk:"family"`
}

func (d *MLModelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ml_model"
}

func (d *MLModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ML Model data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "The ML Model's ID",
				MarkdownDescription: "The ML Model's ID",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "A name that uniquely identifies this ML Model",
				MarkdownDescription: "A name that uniquely identifies this ML Model",
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				Description:         "The visibility of the ML Model.",
				MarkdownDescription: "The visibility of the ML Model.",
			},
			"hosted_by": schema.StringAttribute{
				Computed:            true,
				Description:         "The organization hosting the ML Model.",
				MarkdownDescription: "The organization hosting the ML Model.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ANTHROPIC",
						"AWS_BEDROCK",
						"AZURE_OPENAI",
						"COHERE",
						"CUSTOM",
						"FIREWORKS_AI",
						"GOOGLE",
						"GOOGLE_VERTEX_AI",
						"GROQ",
						"HUGGINGFACE",
						"IBM_WATSONX",
						"MOSAICML",
						"MYSTIC",
						"OPENAI",
						"OPENPIPE",
						"PYQ",
						"REPLICATE",
					),
				},
			},
			"developed_by": schema.StringAttribute{
				Computed:            true,
				Description:         "The organization that developed the ML Model.",
				MarkdownDescription: "The organization that developed the ML Model.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"01_AI",
						"AMAZON",
						"ANTHROPIC",
						"COHERE",
						"ELUTHERAI",
						"FIREWORKS_AI",
						"GOOGLE",
						"HUGGINGFACE",
						"IBM",
						"META",
						"MISTRAL_AI",
						"MOSAICML",
						"NOUS_RESEARCH",
						"OPENAI",
						"OPENCHAT",
						"OPENPIPE",
						"TII",
						"WIZARDLM",
					),
				},
			},
			"family": schema.StringAttribute{
				Computed:            true,
				Description:         "The family of the ML Model.",
				MarkdownDescription: "The family of the ML Model.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CAPYBARA",
						"CHAT_GPT",
						"CLAUDE",
						"COHERE",
						"FALCON",
						"GEMINI",
						"GRANITE",
						"GPT3",
						"FIREWORKS",
						"LLAMA2",
						"LLAMA3",
						"MISTRAL",
						"MPT",
						"OPENCHAT",
						"PALM",
						"SOLAR",
						"TITAN",
						"WIZARD",
						"YI",
						"ZEPHYR",
					),
				},
			},
		},
	}
}

func (d *MLModelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MLModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var mLModelModel *TfMLModelDataSourceModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &mLModelModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mlModelRetrieveParameter := mLModelModel.Name.ValueString()
	if mlModelRetrieveParameter == "" {
		mlModelRetrieveParameter = mLModelModel.Id.ValueString()
	}
	if mlModelRetrieveParameter == "" {
		resp.Diagnostics.AddError("failed to read ML Model", "either `id` or `name` must be set")
		return
	}

	MLModel, err := d.client.MLModels.Retrieve(ctx, mlModelRetrieveParameter)
	if err != nil {
		resp.Diagnostics.AddError("error getting ML Model information", err.Error())
		return
	}

	MLModelModel, diagnostic := NewTfMLModelDataSourceModel(ctx, MLModel)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &MLModelModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
