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

// type TfMLModelExecConfigMetadata := types.MapType {
// 	ElementType: types.Object,
// }

type TfMLModelExecConfig struct {
	ModelIdentifier types.String `tfsdk:"model_identifier"`
	BaseUrl         types.String `tfsdk:"base_url"`
	Metadata        types.Map `tfsdk:"metadata"`
	Features        types.List   `tfsdk:"features"`
}

type TfMLModelResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Visibility  types.String `tfsdk:"visibility"`
	HostedBy    types.String `tfsdk:"hosted_by"`
	DevelopedBy types.String `tfsdk:"developed_by"`
	Family      types.String `tfsdk:"family"`
	ExecConfig  TfMLModelExecConfig `tfsdk:"exec_config"`
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
					"metadata": schema.MapAttribute{
						Description:         "The metadata",
						MarkdownDescription: "The metadata",
						Required:            true,
					}.GetType(),
					"features": schema.ListAttribute{
						Description:         "The features",
						MarkdownDescription: "The features",
						Required:            true,
						ElementType:         schema.StringAttribute{}.GetType(),
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
