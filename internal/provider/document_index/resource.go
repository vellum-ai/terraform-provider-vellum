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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	vellum "terraform-provider-vellum/internal/sdk"
	vellumclient "terraform-provider-vellum/internal/sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DocumentIndexResource{}
var _ resource.ResourceWithImportState = &DocumentIndexResource{}

func NewDocumentIndexResource() resource.Resource {
	return &DocumentIndexResource{}
}

// DocumentIndexResource defines the resource implementation.
type DocumentIndexResource struct {
	client *vellumclient.Client
}

// DocumentIndexResourceModel describes the resource data model.
type DocumentIndexResourceModel struct {
	CopyDocumentsFromIndexId types.String `tfsdk:"copy_documents_from_index_id"`
	Created                  types.String `tfsdk:"created"`
	Environment              types.String `tfsdk:"environment"`
	Id                       types.String `tfsdk:"id"`
	Label                    types.String `tfsdk:"label"`
	Name                     types.String `tfsdk:"name"`
	Status                   types.String `tfsdk:"status"`
}

func (r *DocumentIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document_index"
}

func (r *DocumentIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Document Index resource",

		Attributes: map[string]schema.Attribute{
			"copy_documents_from_index_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Optionally specify the id of a document index from which you'd like to copy and re-index its documents into this newly created index",
				MarkdownDescription: "Optionally specify the id of a document index from which you'd like to copy and re-index its documents into this newly created index",
			},
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
				Description:         "Either the Document Index's ID or its unique name",
				MarkdownDescription: "Either the Document Index's ID or its unique name",
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
	var data DocumentIndexResourceModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Replace this with data.indexing_config, make indexing_config optional in vellum backend.
	DefaultIndexingConfig := map[string]interface{}{
		"chunking": map[string]interface{}{
			"chunker_name": "sentence-chunker",
			"chunker_config": map[string]interface{}{
				"character_limit":   1000,
				"min_overlap_ratio": 0.5,
			},
		},
		"vectorizer": map[string]interface{}{
			"model_name": "hkunlp/instructor-xl",
			"config": map[string]interface{}{
				"instruction_domain":             "",
				"instruction_document_text_type": "plain_text",
				"instruction_query_text_type":    "plain_text",
			},
		},
	}

	// If applicable, this is a great opportunity to initialize any necessary provider client data and make a call using it.
	httpResp, err := r.client.DocumentIndexes.Create(ctx, &vellum.DocumentIndexCreateRequest{
		Label:          data.Label.ValueString(),
		Name:           data.Name.ValueString(),
		IndexingConfig: DefaultIndexingConfig,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create document index, got error: %s", err))
		return
	}

	data.Id = types.StringValue(httpResp.Id)
	data.Created = types.StringValue(httpResp.Created.String())
	data.Environment = types.StringValue(string(*httpResp.Environment))
	data.CopyDocumentsFromIndexId = types.StringNull()
	data.Status = types.StringValue(string(*httpResp.Status))
	tflog.Trace(ctx, fmt.Sprintf("created a document index: %s", data.Id))

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DocumentIndexResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read document index, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DocumentIndexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update document index, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DocumentIndexResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete document index, got error: %s", err))
	//     return
	// }
}

func (r *DocumentIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
