package documentindex

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vellum-ai/terraform-provider-vellum/internal/vellum"
	vellumclient "github.com/vellum-ai/terraform-provider-vellum/internal/vellum/client"
)

// ResourceModel represents the document index resource model.
// This is expressed in terraform with the following:
//
//	resource "vellum_document_index" "managed" {
//	  label = "Managed Index"
//	  name  = "managed-index"
//	}
type ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Environment types.String `tfsdk:"environment"`
	Status      types.String `tfsdk:"status"`
	Created     types.String `tfsdk:"created"`
}

// Resource represents the document index resource. This type is extendable;
// simply override any of the exported methods in the resource_hooks.go file
// to customize the resource.
type Resource struct {
	base *baseResource
}

// Compile-time assertion that ensures the provider satisfies the resource.Resource
// interface.
var _ resource.Resource = (*Resource)(nil)

// NewResource returns a new document index resource.
func NewResource() resource.Resource {
	return &Resource{
		base: newBaseResource(),
	}
}

// baseResource implements the base functionality of the resource.
//
// Do NOT edit this type directly; use resource_hooks.go instead.
type baseResource struct {
	Vellum *vellumclient.Client
}

var _ resource.Resource = (*baseResource)(nil)

func newBaseResource() *baseResource {
	return &baseResource{}
}

func (b *baseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document_index"
}

func (b *baseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Document Index resource",

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
				Required:            true,
				Description:         "A human-readable label for the document index",
				MarkdownDescription: "A human-readable label for the document index",
			},
			"environment": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
				MarkdownDescription: "The environment this document index is used in\n\n* `DEVELOPMENT` - Development\n* `STAGING` - Staging\n* `PRODUCTION` - Production",
			},
			"status": schema.StringAttribute{
				Optional:            true,
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

func (b *baseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (b *baseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := b.modelToCreateRequest(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create document index error",
			fmt.Sprintf("Unable to create document index request: %v", err),
		)
		return
	}

	response, err := b.Vellum.DocumentIndexes.Create(
		ctx,
		request,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create document index error",
			fmt.Sprintf("Unable to create document index, got error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			b.createResponseToModel(response),
		)...,
	)
}

func (b *baseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model *ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := b.Vellum.DocumentIndexes.Retrieve(
		ctx,
		b.modelToRetrieveRequest(model),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read document index error",
			fmt.Sprintf("Unable to read document index, got error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			b.retrieveResponseToModel(response),
		)...,
	)
}

func (b *baseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := b.modelToPartialUpdateRequest(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update document index error",
			fmt.Sprintf("Unable to create document index request: %v", err),
		)
		return
	}

	response, err := b.Vellum.DocumentIndexes.PartialUpdate(ctx,
		model.Id.ValueString(),
		request,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update document index error",
			fmt.Sprintf("Unable to update document index, got error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			b.updateResponseToModel(response),
		)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *baseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model *ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := b.Vellum.DocumentIndexes.Destroy(
		ctx,
		b.modelToDestroyRequest(model),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete document index error",
			fmt.Sprintf("Unable to delete document index, got error: %s", err),
		)
		return
	}
}

func (b *baseResource) modelToCreateRequest(model *ResourceModel) (*vellum.DocumentIndexCreateRequest, error) {
	var status *vellum.EntityStatus
	if !model.Status.IsNull() {
		value, err := vellum.NewEntityStatusFromString(model.Status.ValueString())
		if err != nil {
			return nil, err
		}
		status = value.Ptr()
	}

	var environment *vellum.EnvironmentEnum
	if !model.Environment.IsNull() {
		value, err := vellum.NewEnvironmentEnumFromString(model.Environment.ValueString())
		if err != nil {
			return nil, err
		}
		environment = value.Ptr()
	}

	return &vellum.DocumentIndexCreateRequest{
		Label:       model.Label.ValueString(),
		Name:        model.Name.ValueString(),
		Status:      status,
		Environment: environment,
	}, nil
}

func (b *baseResource) modelToRetrieveRequest(model *ResourceModel) string {
	return model.Id.ValueString()
}

func (b *baseResource) modelToPartialUpdateRequest(model *ResourceModel) (*vellum.PatchedDocumentIndexUpdateRequest, error) {
	var status *vellum.EntityStatus
	if !model.Status.IsNull() {
		value, err := vellum.NewEntityStatusFromString(model.Status.ValueString())
		if err != nil {
			return nil, err
		}
		status = value.Ptr()
	}

	var environment *vellum.EnvironmentEnum
	if !model.Environment.IsNull() {
		value, err := vellum.NewEnvironmentEnumFromString(model.Environment.ValueString())
		if err != nil {
			return nil, err
		}
		environment = value.Ptr()
	}

	return &vellum.PatchedDocumentIndexUpdateRequest{
		Label:       model.Label.ValueStringPointer(),
		Status:      status,
		Environment: environment,
	}, nil
}

func (b *baseResource) modelToDestroyRequest(model *ResourceModel) string {
	return model.Id.ValueString()
}

func (b *baseResource) createResponseToModel(response *vellum.DocumentIndexRead) *ResourceModel {
	return &ResourceModel{
		Id:          types.StringValue(response.Id),
		Name:        types.StringValue(response.Name),
		Created:     types.StringValue(response.Created.Format(time.RFC3339)),
		Environment: types.StringValue(string(*response.Environment)),
		Label:       types.StringValue(response.Label),
		Status:      types.StringValue(string(*response.Status)),
	}
}

func (b *baseResource) retrieveResponseToModel(response *vellum.DocumentIndexRead) *ResourceModel {
	return &ResourceModel{
		Id:          types.StringValue(response.Id),
		Name:        types.StringValue(response.Name),
		Created:     types.StringValue(response.Created.Format(time.RFC3339)),
		Environment: types.StringValue(string(*response.Environment)),
		Label:       types.StringValue(response.Label),
		Status:      types.StringValue(string(*response.Status)),
	}
}

func (b *baseResource) updateResponseToModel(response *vellum.DocumentIndexRead) *ResourceModel {
	return &ResourceModel{
		Id:          types.StringValue(response.Id),
		Name:        types.StringValue(response.Name),
		Created:     types.StringValue(response.Created.Format(time.RFC3339)),
		Environment: types.StringValue(string(*response.Environment)),
		Label:       types.StringValue(response.Label),
		Status:      types.StringValue(string(*response.Status)),
	}
}
