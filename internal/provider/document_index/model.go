package document_index

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
)

func NewVellumDocumentIndexCreateRequest(ctx context.Context, documentIndexModel *TfDocumentIndexResourceModel) (*vellum.DocumentIndexCreateRequest, diag.Diagnostics) {
	// TODO: Replace this with data.indexing_config, improve indexing_config param in vellum backend.
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

	request := vellum.DocumentIndexCreateRequest{
		Label:          documentIndexModel.Label.ValueString(),
		Name:           documentIndexModel.Name.ValueString(),
		IndexingConfig: DefaultIndexingConfig,
	}

	return &request, nil
}

func NewTfDocumentIndexModel(ctx context.Context, model *TfDocumentIndexResourceModel, documentIndex *vellum.DocumentIndexRead) (*TfDocumentIndexResourceModel, diag.Diagnostics) {
	documentIndexModel := &TfDocumentIndexResourceModel{
		Id:          types.StringValue(documentIndex.Id),
		Name:        types.StringValue(documentIndex.Name),
		Created:     types.StringValue(documentIndex.Created.String()),
		Environment: types.StringValue(string(*documentIndex.Environment)),
		Label:       types.StringValue(documentIndex.Label),
		Status:      types.StringValue(string(*documentIndex.Status)),
	}

	return documentIndexModel, nil
}

func NewTfDocumentIndexDataSourceModel(ctx context.Context, documentIndex *vellum.DocumentIndexRead) (*TfDocumentIndexDataSourceModel, diag.Diagnostics) {
	documentIndexModel := &TfDocumentIndexDataSourceModel{
		Id:          types.StringValue(documentIndex.Id),
		Name:        types.StringValue(documentIndex.Name),
		Created:     types.StringValue(documentIndex.Created.String()),
		Environment: types.StringValue(string(*documentIndex.Environment)),
		Label:       types.StringValue(documentIndex.Label),
		Status:      types.StringValue(string(*documentIndex.Status)),
	}

	return documentIndexModel, nil
}
