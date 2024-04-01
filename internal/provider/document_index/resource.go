// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package document_index

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
	vellumclient "terraform-provider-vellum/internal/sdk/client"
)

var _ resource.ResourceWithConfigure = &DocumentIndexResource{}
var _ resource.ResourceWithImportState = &DocumentIndexResource{}

type DocumentIndexResource struct {
	client *vellumclient.Client
}

func Resource() resource.Resource {
	return &DocumentIndexResource{}
}

type TfDocumentIndexResourceModel struct {
	Created     types.String `tfsdk:"created"`
	Environment types.String `tfsdk:"environment"`
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
}

func (r *DocumentIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document_index"
}

func (r *DocumentIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Document Index resource",

		Attributes: map[string]schema.Attribute{
			"created": schema.StringAttribute{
				Computed: true,
			},
			"environment": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
				MarkdownDescription: "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DEVELOPMENT",
						"STAGING",
						"PRODUCTION",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The Document Index's ID",
				MarkdownDescription: "The Document Index's ID",
			},
			"label": schema.StringAttribute{
				Required:            true,
				Description:         "A human-readable label for the document index",
				MarkdownDescription: "A human-readable label for the document index",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 150),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "A name that uniquely identifies this index within its workspace",
				MarkdownDescription: "A name that uniquely identifies this index within its workspace",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 150),
				},
			},
			"status": schema.StringAttribute{
				Optional:            true,
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

func (r *DocumentIndexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DocumentIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var documentIndexPlan *TfDocumentIndexResourceModel

	diags := req.Plan.Get(ctx, &documentIndexPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentIndexRequest, d := NewVellumDocumentIndexCreateRequest(ctx, documentIndexPlan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentIndex, err := r.client.DocumentIndexes.Create(ctx, documentIndexRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create document index, got error: %s", err))
		return
	}

	documentIndexModel, diagnostic := NewTfDocumentIndexModel(ctx, documentIndexPlan, documentIndex)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &documentIndexModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DocumentIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var documentIndexState TfDocumentIndexResourceModel
	var err error
	resp.Diagnostics.Append(req.State.Get(ctx, &documentIndexState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentIndex, err := r.client.DocumentIndexes.Retrieve(ctx, documentIndexState.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read document index, got error: %s", err))
		return
	}

	documentIndexModel, diagnostic := NewTfDocumentIndexModel(ctx, &documentIndexState, documentIndex)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &documentIndexModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DocumentIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var documentIndexPlan *TfDocumentIndexResourceModel
	var documentIndexState *TfDocumentIndexResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &documentIndexPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &documentIndexState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := documentIndexState.Id.ValueString()
	label := documentIndexPlan.Label.ValueString()

	var status *vellum.EntityStatus
	if documentIndexPlan.Status.ValueString() != "" {
		s, _ := vellum.NewEntityStatusFromString(documentIndexPlan.Status.ValueString())
		status = &s
	}

	var environment *vellum.EnvironmentEnum;
	if documentIndexPlan.Environment.ValueString() != "" {
		env, _ := vellum.NewEnvironmentEnumFromString(documentIndexPlan.Environment.ValueString())
		environment = &env
	}

	documentIndex, err := r.client.DocumentIndexes.PartialUpdate(ctx,
		id,
		&vellum.PatchedDocumentIndexUpdateRequest{
			Label:       &label,
			Status:      status,
			Environment: environment,
		})

	if err != nil {
		resp.Diagnostics.AddError("error during document index update", err.Error())
		return
	}

	documentIndexModel, diagnostic := NewTfDocumentIndexModel(ctx, documentIndexPlan, documentIndex)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &documentIndexModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DocumentIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var documentIndexState *TfDocumentIndexResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &documentIndexState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DocumentIndexes.Destroy(
		ctx,
		documentIndexState.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error when destroying the document index resource", err.Error())
		return
	}
}

func (r *DocumentIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
