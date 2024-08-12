// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ml_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
	vellumclient "terraform-provider-vellum/internal/sdk/client"
)

var _ resource.ResourceWithConfigure = &MLModelResource{}
var _ resource.ResourceWithImportState = &MLModelResource{}

type MLModelResource struct {
	client *vellumclient.Client
}

func Resource() resource.Resource {
	return &MLModelResource{}
}

type TfHuggingFaceTokenizerConfig struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

type TfTikTokenTokenizerConfig struct {
	Name types.String `tfsdk:"name"`
}

type TfMlModelTokenizerConfig struct {
	Type        types.String                 `tfsdk:"type"`
	HuggingFace TfHuggingFaceTokenizerConfig `tfsdk:"hugging_face"`
	Tiktoken    TfTikTokenTokenizerConfig    `tfsdk:"tiktoken"`
}

type TfMLModelExecConfig struct {
	ModelIdentifier        types.String             `tfsdk:"model_identifier"`
	BaseUrl                types.String             `tfsdk:"base_url"`
	Features               types.List               `tfsdk:"features"`
	Metadata               types.Map                `tfsdk:"metadata"`
	TokenizerConfig        TfMlModelTokenizerConfig `tfsdk:"tokenizer_config"`
	ForceSystemCredentials types.Bool               `tfsdk:"force_system_credentials"`
}

type TfOpenApiNumberProperty struct {
	Minimum          types.Float64 `tfsdk:"minimum"`
	Maximum          types.Float64 `tfsdk:"maximum"`
	Format           types.String  `tfsdk:"format"`
	ExclusiveMinimum types.Bool    `tfsdk:"exclusive_minimum"`
	ExclusiveMaximum types.Bool    `tfsdk:"exclusive_maximum"`
	Default          types.Float64 `tfsdk:"default"`
	Title            types.String  `tfsdk:"title"`
	Description      types.String  `tfsdk:"description"`
}

type TfOpenApiIntegerProperty struct {
	Minimum          types.Int64  `tfsdk:"minimum"`
	Maximum          types.Int64  `tfsdk:"maximum"`
	ExclusiveMinimum types.Bool   `tfsdk:"exclusive_minimum"`
	ExclusiveMaximum types.Bool   `tfsdk:"exclusive_maximum"`
	Default          types.Int64  `tfsdk:"default"`
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
}

type TfOpenApiArrayProperty struct {
	MinItems    types.Int64       `tfsdk:"min_items"`
	MaxItems    types.Int64       `tfsdk:"max_items"`
	UniqueItems types.Bool        `tfsdk:"unique_items"`
	Items       TfOpenApiProperty `tfsdk:"items"`
	PrefixItems types.List        `tfsdk:"prefix_items"`
	Contains    TfOpenApiProperty `tfsdk:"contains"`
	MinContains types.Int64       `tfsdk:"min_contains"`
	MaxContains types.Int64       `tfsdk:"max_contains"`
	Default     types.List        `tfsdk:"default"`
	Title       types.String      `tfsdk:"title"`
	Description types.String      `tfsdk:"description"`
}

type TfOpenApiObjectProperty struct {
	Properties           types.Map         `tfsdk:"properties"`
	Required             types.List        `tfsdk:"required"`
	MinProperties        types.Int64       `tfsdk:"min_properties"`
	MaxProperties        types.Int64       `tfsdk:"max_properties"`
	PropertyNames        TfOpenApiProperty `tfsdk:"property_names"`
	AdditionalProperties TfOpenApiProperty `tfsdk:"additional_properties"`
	PatternProperties    types.Map         `tfsdk:"pattern_properties"`
	Default              types.Map         `tfsdk:"default"`
	Title                types.String      `tfsdk:"title"`
	Description          types.String      `tfsdk:"description"`
}

type TfOpenApiStringProperty struct {
	MinLength   types.Int64  `tfsdk:"min_length"`
	MaxLength   types.Int64  `tfsdk:"max_length"`
	Pattern     types.String `tfsdk:"pattern"`
	Format      types.String `tfsdk:"format"`
	Default     types.String `tfsdk:"default"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

type TfOpenApiBooleanProperty struct {
	Default     types.Bool   `tfsdk:"default"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

type TfOpenApiOneOfProperty struct {
	OneOf       types.List   `tfsdk:"oneOf"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

type TfOpenApiConstProperty struct {
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Const       types.String `tfsdk:"const"`
}

type TfOpenApiProperty struct {
	Type    types.String
	Array   TfOpenApiArrayProperty
	Object  TfOpenApiObjectProperty
	Integer TfOpenApiIntegerProperty
	Number  TfOpenApiNumberProperty
	String  TfOpenApiStringProperty
	Boolean TfOpenApiBooleanProperty
	OneOf   TfOpenApiOneOfProperty
	Const   TfOpenApiConstProperty
}

type TfMlModelParameterConfig struct {
	Temperature      types.Number `tfsdk:"temperature"`
	MaxTokens        types.Int64  `tfsdk:"max_tokens"`
	Stop             types.List   `tfsdk:"stop"`
	TopP             types.Number `tfsdk:"top_p"`
	TopK             types.Int64  `tfsdk:"top_k"`
	FrequencyPenalty types.Number `tfsdk:"frequency_penalty"`
	PresencePenalty  types.Number `tfsdk:"presence_penalty"`
	LogitBias        types.Map    `tfsdk:"logit_bias"`
	CustomParameters types.Map    `tfsdk:"custom_parameters"`
}

type TfMlModelDisplayConfigLabelled struct {
	Label                  types.String  `tfsdk:"label"`
	Description            types.String  `tfsdk:"description"`
	Tags                   types.String  `tfsdk:"tags"`
	DefaultDisplayPriority types.Float64 `tfsdk:"default_display_priority"`
}

type TfMLModelResourceModel struct {
	Id              types.String                   `tfsdk:"id"`
	Name            types.String                   `tfsdk:"name"`
	Visibility      types.String                   `tfsdk:"visibility"`
	HostedBy        types.String                   `tfsdk:"hosted_by"`
	DevelopedBy     types.String                   `tfsdk:"developed_by"`
	Family          types.String                   `tfsdk:"family"`
	ExecConfig      TfMLModelExecConfig            `tfsdk:"exec_config"`
	ParameterConfig TfMlModelParameterConfig       `tfsdk:"parameter_config"`
	DisplayConfig   TfMlModelDisplayConfigLabelled `tfsdk:"display_config"`
}

func (r *MLModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ml_model"
}

func (r *MLModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ML Model resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ML Model's ID",
				MarkdownDescription: "The ML Model's ID",
			},

			"name": schema.StringAttribute{
				Required:            true,
				Description:         "A name that uniquely identifies this ML Model",
				MarkdownDescription: "A name that uniquely identifies this ML Model",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 150),
				},
			},
			"visibility": schema.StringAttribute{
				Description:         "The visibility of the ML Model.",
				MarkdownDescription: "The visibility of the ML Model.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DEFAULT",
						"PUBLIC",
						"PRIVATE",
						"DISABLED",
					),
				},
			},
			"hosted_by": schema.StringAttribute{
				Description:         "The organization hosting the ML Model.",
				MarkdownDescription: "The organization hosting the ML Model.",
				Required:            true,
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
				Description:         "The organization that developed the ML Model.",
				MarkdownDescription: "The organization that developed the ML Model.",
				Required:            true,
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
				Description:         "The family of the ML Model.",
				MarkdownDescription: "The family of the ML Model.",
				Required:            true,
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
			"exec_config": schema.ObjectAttribute{
				Description:         "The execution configuration of the ML Model.",
				MarkdownDescription: "The execution configuration of the ML Model.",
				Required:            true,
				AttributeTypes: map[string]attr.Type{
					"model_identifier": schema.StringAttribute{
						Description:         "The model identifier",
						MarkdownDescription: "The model identifier",
						Required:            true,
					}.GetType(),
					"base_url": schema.StringAttribute{
						Description:         "The base URL",
						MarkdownDescription: "The base URL",
						Required:            true,
					}.GetType(),
					"features": schema.ListAttribute{
						Description:         "The features",
						MarkdownDescription: "The features",
						Required:            true,
						ElementType:         schema.StringAttribute{}.GetType(),
					}.GetType(),
					"metadata": schema.MapAttribute{
						Description: "Arbitrary JSON object",
						Required:    true,
						ElementType: types.StringType,
					}.GetType(),
				},
			},
		},
	}
}

func (r *MLModelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vellumclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MLModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var mlModelPlan *TfMLModelResourceModel

	diags := req.Plan.Get(ctx, &mlModelPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mlModelRequest, d := NewVellumMLModelCreateRequest(ctx, mlModelPlan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	mlModel, err := r.client.MLModels.Create(ctx, mlModelRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ML Model, got error: %s", err))
		return
	}

	mlModelModel, diagnostic := NewTfMLModelModel(ctx, mlModelPlan, mlModel)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &mlModelModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MLModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var mlModelState TfMLModelResourceModel
	var err error
	resp.Diagnostics.Append(req.State.Get(ctx, &mlModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mlModel, err := r.client.MLModels.Retrieve(ctx, mlModelState.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ML Model, got error: %s", err))
		return
	}

	mlModelModel, diagnostic := NewTfMLModelModel(ctx, &mlModelState, mlModel)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &mlModelModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MLModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var mlModelPlan *TfMLModelResourceModel
	var mlModelState *TfMLModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &mlModelPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &mlModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := mlModelState.Id.ValueString()

	var visibility *vellum.VisibilityEnum
	if mlModelPlan.Visibility.ValueString() != "" {
		s, _ := vellum.NewVisibilityEnumFromString(mlModelPlan.Visibility.ValueString())
		visibility = &s
	}

	mlModel, err := r.client.MLModels.PartialUpdate(ctx,
		id,
		&vellum.PatchedMlModelUpdateRequest{
			Visibility: visibility,
		})

	if err != nil {
		resp.Diagnostics.AddError("error during ML Model update", err.Error())
		return
	}

	mlModelModel, diagnostic := NewTfMLModelModel(ctx, mlModelPlan, mlModel)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &mlModelModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MLModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var mlModelState *TfMLModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &mlModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := mlModelState.Id.ValueString()
	visibility := vellum.VisibilityEnum("DISABLED")

	_, err := r.client.MLModels.PartialUpdate(ctx,
		id,
		&vellum.PatchedMlModelUpdateRequest{
			Visibility: &visibility,
		})

	if err != nil {
		resp.Diagnostics.AddError("error when disabling the ML Model resource", err.Error())
		return
	}
}

func (r *MLModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
