// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ml_model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

type TfOpenApiNumberProperty struct {
	Minimum          types.Number `tfsdk:"minimum"`
	Maximum          types.Number `tfsdk:"maximum"`
	Format           types.String `tfsdk:"format"`
	ExclusiveMinimum types.Bool   `tfsdk:"exclusive_minimum"`
	ExclusiveMaximum types.Bool   `tfsdk:"exclusive_maximum"`
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
}

type TfOpenApiIntegerProperty struct {
	Minimum          types.Int64  `tfsdk:"minimum"`
	Maximum          types.Int64  `tfsdk:"maximum"`
	ExclusiveMinimum types.Bool   `tfsdk:"exclusive_minimum"`
	ExclusiveMaximum types.Bool   `tfsdk:"exclusive_maximum"`
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
}

type TfOpenApiStringProperty struct {
	MinLength   types.Int64  `tfsdk:"min_length"`
	MaxLength   types.Int64  `tfsdk:"max_length"`
	Pattern     types.String `tfsdk:"pattern"`
	Format      types.String `tfsdk:"format"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

type TfOpenApiBooleanProperty struct {
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

type TfOpenApiArrayProperty struct {
	MinItems    types.Int64          `tfsdk:"min_items"`
	MaxItems    types.Int64          `tfsdk:"max_items"`
	UniqueItems types.Bool           `tfsdk:"unique_items"`
	Items       *TfOpenApiProperty   `tfsdk:"items"`
	PrefixItems []*TfOpenApiProperty `tfsdk:"prefix_items"`
	Contains    *TfOpenApiProperty   `tfsdk:"contains"`
	MinContains types.Int64          `tfsdk:"min_contains"`
	MaxContains types.Int64          `tfsdk:"max_contains"`
	Title       types.String         `tfsdk:"title"`
	Description types.String         `tfsdk:"description"`
}

type TfOpenApiObjectProperty struct {
	Properties           map[string]*TfOpenApiProperty `tfsdk:"properties"`
	Required             types.List                    `tfsdk:"required"`
	MinProperties        types.Int64                   `tfsdk:"min_properties"`
	MaxProperties        types.Int64                   `tfsdk:"max_properties"`
	PropertyNames        *TfOpenApiProperty            `tfsdk:"property_names"`
	AdditionalProperties *TfOpenApiProperty            `tfsdk:"additional_properties"`
	PatternProperties    map[string]*TfOpenApiProperty `tfsdk:"pattern_properties"`
	Title                types.String                  `tfsdk:"title"`
	Description          types.String                  `tfsdk:"description"`
}

type TfOpenApiProperty struct {
	Integer *TfOpenApiIntegerProperty `tfsdk:"integer"`
	Number  *TfOpenApiNumberProperty  `tfsdk:"number"`
	String  *TfOpenApiStringProperty  `tfsdk:"string"`
	Boolean *TfOpenApiBooleanProperty `tfsdk:"boolean"`
	// TODO: Add object and array
	// 	Note: Terraform plugin doesn't support recursive schemas
}

type TfMLModelParameterConfig struct {
	Temperature      *TfOpenApiNumberProperty      `tfsdk:"temperature"`
	MaxTokens        *TfOpenApiIntegerProperty     `tfsdk:"max_tokens"`
	Stop             *TfOpenApiArrayProperty       `tfsdk:"stop"`
	TopP             *TfOpenApiNumberProperty      `tfsdk:"top_p"`
	TopK             *TfOpenApiIntegerProperty     `tfsdk:"top_k"`
	FrequencyPenalty *TfOpenApiNumberProperty      `tfsdk:"frequency_penalty"`
	PresencePenalty  *TfOpenApiNumberProperty      `tfsdk:"presence_penalty"`
	LogitBias        *TfOpenApiObjectProperty      `tfsdk:"logit_bias"`
	CustomParameters map[string]*TfOpenApiProperty `tfsdk:"custom_parameters"`
}

type TfMLModelDisplayConfig struct {
	Label                  types.String `tfsdk:"label"`
	Description            types.String `tfsdk:"description"`
	Tags                   types.List   `tfsdk:"tags"`
	DefaultDisplayPriority types.Number `tfsdk:"default_display_priority"`
}

type TfHuggingFaceTokenizerConfig struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

type TiktokenTokenizerConfig struct {
	Name types.String `tfsdk:"name"`
}

type TfMlModelTokenizerConfig struct {
	Type        types.String                  `tfsdk:"type"`
	HuggingFace *TfHuggingFaceTokenizerConfig `tfsdk:"hugging_face"`
	Tiktoken    *TiktokenTokenizerConfig      `tfsdk:"tiktoken"`
}

type TfMLModelExecConfig struct {
	ModelIdentifier        types.String              `tfsdk:"model_identifier"`
	BaseUrl                types.String              `tfsdk:"base_url"`
	Features               types.List                `tfsdk:"features"`
	Metadata               types.Map                 `tfsdk:"metadata"`
	ForceSystemCredentials types.Bool                `tfsdk:"force_system_credentials"`
	TokenizerConfig        *TfMlModelTokenizerConfig `tfsdk:"tokenizer_config"`
}

type TfMLModelResourceModel struct {
	Id              types.String              `tfsdk:"id"`
	Name            types.String              `tfsdk:"name"`
	Visibility      types.String              `tfsdk:"visibility"`
	HostedBy        types.String              `tfsdk:"hosted_by"`
	DevelopedBy     types.String              `tfsdk:"developed_by"`
	Family          types.String              `tfsdk:"family"`
	ExecConfig      TfMLModelExecConfig       `tfsdk:"exec_config"`
	ParameterConfig *TfMLModelParameterConfig `tfsdk:"parameter_config"`
	DisplayConfig   *TfMLModelDisplayConfig   `tfsdk:"display_config"`
}

var TfOpenApiNumberPropertySchema = map[string]schema.Attribute{
	"minimum": schema.NumberAttribute{
		Description:         "Minimum value",
		MarkdownDescription: "Minimum value",
		Optional:            true,
	},
	"maximum": schema.NumberAttribute{
		Description:         "Maximum value",
		MarkdownDescription: "Maximum value",
		Optional:            true,
	},
	"format": schema.StringAttribute{
		Description:         "Format",
		MarkdownDescription: "Format",
		Optional:            true,
	},
	"exclusive_minimum": schema.BoolAttribute{
		Description:         "Exclusive Minimum",
		MarkdownDescription: "Exclusive Minimum",
		Optional:            true,
	},
	"exclusive_maximum": schema.BoolAttribute{
		Description:         "Exclusive Maximum",
		MarkdownDescription: "Exclusive Maximum",
		Optional:            true,
	},
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiIntegerPropertySchema = map[string]schema.Attribute{
	"minimum": schema.Int64Attribute{
		Description:         "Minimum value",
		MarkdownDescription: "Minimum value",
		Optional:            true,
	},
	"maximum": schema.Int64Attribute{
		Description:         "Maximum value",
		MarkdownDescription: "Maximum value",
		Optional:            true,
	},
	"exclusive_minimum": schema.BoolAttribute{
		Description:         "Exclusive Minimum",
		MarkdownDescription: "Exclusive Minimum",
		Optional:            true,
	},
	"exclusive_maximum": schema.BoolAttribute{
		Description:         "Exclusive Maximum",
		MarkdownDescription: "Exclusive Maximum",
		Optional:            true,
	},
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiStringPropertySchema = map[string]schema.Attribute{
	"min_length": schema.Int64Attribute{
		Description:         "Min length",
		MarkdownDescription: "Min length",
		Optional:            true,
	},
	"max_length": schema.Int64Attribute{
		Description:         "Max length",
		MarkdownDescription: "Max length",
		Optional:            true,
	},
	"pattern": schema.StringAttribute{
		Description:         "Pattern",
		MarkdownDescription: "Pattern",
		Optional:            true,
	},
	"format": schema.StringAttribute{
		Description:         "Format",
		MarkdownDescription: "Format",
		Optional:            true,
	},
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiBooleanPropertySchema = map[string]schema.Attribute{
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiArrayPropertySchema = map[string]schema.Attribute{
	"min_items": schema.Int64Attribute{
		Description:         "Min items",
		MarkdownDescription: "Min items",
		Optional:            true,
	},
	"max_items": schema.Int64Attribute{
		Description:         "Max items",
		MarkdownDescription: "Max items",
		Optional:            true,
	},
	"unique_items": schema.BoolAttribute{
		Description:         "Unique items",
		MarkdownDescription: "Unique items",
		Optional:            true,
	},
	"items": schema.SingleNestedAttribute{
		Description:         "Items",
		MarkdownDescription: "Items",
		Optional:            true,
		Attributes:          TfOpenApiPropertySchema,
	},
	"prefix_items": schema.ListNestedAttribute{
		Description:         "Prefix items",
		MarkdownDescription: "Prefix items",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: TfOpenApiPropertySchema,
		},
	},
	"contains": schema.SingleNestedAttribute{
		Description:         "Contains",
		MarkdownDescription: "Contains",
		Optional:            true,
		Attributes:          TfOpenApiPropertySchema,
	},
	"min_contains": schema.Int64Attribute{
		Description:         "Min contains",
		MarkdownDescription: "Min contains",
		Optional:            true,
	},
	"max_contains": schema.Int64Attribute{
		Description:         "Max contains",
		MarkdownDescription: "Max contains",
		Optional:            true,
	},
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiObjectPropertySchema = map[string]schema.Attribute{
	"properties": schema.MapNestedAttribute{
		Description:         "Properties",
		MarkdownDescription: "Properties",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: TfOpenApiPropertySchema,
		},
	},
	"required": schema.ListAttribute{
		Description:         "Required",
		MarkdownDescription: "Required",
		Optional:            true,
		ElementType:         types.StringType,
	},
	"min_properties": schema.Int64Attribute{
		Description:         "Min properties",
		MarkdownDescription: "Min properties",
		Optional:            true,
	},
	"max_properties": schema.Int64Attribute{
		Description:         "Max properties",
		MarkdownDescription: "Max properties",
		Optional:            true,
	},
	"property_names": schema.SingleNestedAttribute{
		Description:         "Property names",
		MarkdownDescription: "Property names",
		Optional:            true,
		Attributes:          TfOpenApiPropertySchema,
	},
	"additional_properties": schema.SingleNestedAttribute{
		Description:         "Additional properties",
		MarkdownDescription: "Additional properties",
		Optional:            true,
		Attributes:          TfOpenApiPropertySchema,
	},
	"pattern_properties": schema.MapNestedAttribute{
		Description:         "Pattern properties",
		MarkdownDescription: "Pattern properties",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: TfOpenApiPropertySchema,
		},
	},
	"title": schema.StringAttribute{
		Description:         "Title",
		MarkdownDescription: "Title",
		Optional:            true,
	},
	"description": schema.StringAttribute{
		Description:         "Description",
		MarkdownDescription: "Description",
		Optional:            true,
	},
}

var TfOpenApiPropertySchema = map[string]schema.Attribute{
	"number": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfOpenApiNumberPropertySchema,
	},
	"string": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfOpenApiStringPropertySchema,
	},
	"integer": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfOpenApiIntegerPropertySchema,
	},
	"boolean": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfOpenApiBooleanPropertySchema,
	},
	// TODO: Add object and array
	// 	Note: Terraform plugin doesn't support recursive schemas
}

var TfHuggingFaceTokenizerConfigSchema = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required: true,
	},
	"path": schema.StringAttribute{
		Optional: true,
	},
}

var TfTikTokenTokenizerConfigSchema = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required: true,
	},
}

var TfMlModelTokenizerConfigSchema = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Required: true,
	},
	"hugging_face": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfHuggingFaceTokenizerConfigSchema,
	},
	"tiktoken": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: TfTikTokenTokenizerConfigSchema,
	},
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
			"exec_config": schema.SingleNestedAttribute{
				Description:         "The execution configuration of the ML Model.",
				MarkdownDescription: "The execution configuration of the ML Model.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"model_identifier": schema.StringAttribute{
						Description:         "The model identifier",
						MarkdownDescription: "The model identifier",
						Required:            true,
					},
					"base_url": schema.StringAttribute{
						Description:         "The base URL",
						MarkdownDescription: "The base URL",
						Required:            true,
					},
					"features": schema.ListAttribute{
						Description:         "The features",
						MarkdownDescription: "The features",
						Required:            true,
						ElementType:         schema.StringAttribute{}.GetType(),
					},
					"metadata": schema.MapAttribute{
						Description: "Arbitrary JSON object",
						Required:    true,
						ElementType: types.StringType,
					},
					"force_system_credentials": schema.BoolAttribute{
						Description: "",
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"tokenizer_config": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes:  TfMlModelTokenizerConfigSchema,
					},
				},
			},
			"parameter_config": schema.SingleNestedAttribute{
				Description:         "Configuration for the ML Model's parameters.",
				MarkdownDescription: "Configuration for the ML Model's parameters.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"temperature": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiNumberPropertySchema,
					},
					"max_tokens": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiIntegerPropertySchema,
					},
					"stop": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiArrayPropertySchema,
					},
					"top_p": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiNumberPropertySchema,
					},
					"top_k": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiIntegerPropertySchema,
					},
					"frequency_penalty": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiNumberPropertySchema,
					},
					"presence_penalty": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiNumberPropertySchema,
					},
					"logit_bias": schema.SingleNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Attributes:          TfOpenApiObjectPropertySchema,
					},
					"custom_parameters": schema.MapNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: TfOpenApiPropertySchema,
						},
					},
				},
			},
			"display_config": schema.SingleNestedAttribute{
				Description:         "Configuration for how to display the ML Model.",
				MarkdownDescription: "Configuration for how to display the ML Model.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"label": schema.StringAttribute{
						Required: true,
					},
					"description": schema.StringAttribute{
						Required: true,
					},
					"tags": schema.ListAttribute{
						Optional: true,
						ElementType: schema.StringAttribute{
							Validators: []validator.String{
								stringvalidator.OneOf(
									"TEXT",
									"CHAT",
									"OPEN_SOURCE",
									"FINETUNED",
									"NEW",
									"ALPHA",
									"BETA",
									"DEPRECATED",
								),
							},
						}.GetType(),
					},
					"default_display_priority": schema.NumberAttribute{
						Optional: true,
					},
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

	mlModel, err := r.client.MlModels.Create(ctx, mlModelRequest)
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

	mlModel, err := r.client.MlModels.Retrieve(ctx, mlModelState.Id.ValueString())
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

	mlModel, err := r.client.MlModels.PartialUpdate(ctx,
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

	_, err := r.client.MlModels.PartialUpdate(ctx,
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
