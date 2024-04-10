package documentindex

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/vellum-ai/terraform-provider-vellum/internal/vellum"
)

// --- Override the following methods to customize the resource ---

func (d *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	d.base.Metadata(ctx, req, resp)
}

func (d *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	d.base.Schema(ctx, req, resp)
}

func (d *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.base.Configure(ctx, req, resp)
}

func (d *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := d.modelToCreateRequest(model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create document index error",
			fmt.Sprintf("Unable to create document index request: %v", err),
		)
		return
	}

	response, err := d.base.Vellum.DocumentIndexes.Create(
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
			d.base.createResponseToModel(response),
		)...,
	)
}

func (d *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	d.base.Read(ctx, req, resp)
}

func (d *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	d.base.Update(ctx, req, resp)
}

func (d *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	d.base.Delete(ctx, req, resp)
}

func (d *Resource) modelToCreateRequest(model *ResourceModel) (*vellum.DocumentIndexCreateRequest, error) {
	request, err := d.base.modelToCreateRequest(model)
	if err != nil {
		return nil, err
	}
	request.IndexingConfig = defaultIndexingConfig
	return request, nil
}

// TODO: Replace this with data.indexing_config, improve indexing_config param in vellum backend.
var defaultIndexingConfig = map[string]interface{}{
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
