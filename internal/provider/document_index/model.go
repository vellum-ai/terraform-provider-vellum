package document_index

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
)

func NewVellumDocumentIndexCreateRequest(ctx context.Context, documentIndexModel *TfDocumentIndexResourceModel) (*vellum.DocumentIndexCreateRequest, diag.Diagnostics) {
	// TODO: Replace this with data.indexing_config, improve indexing_config param in vellum backend.
	DefaultIndexingConfig := vellum.DocumentIndexIndexingConfigRequest{
		Vectorizer: &vellum.IndexingConfigVectorizerRequest{
			ModelName: "hkunlp/instructor-xl",
			HkunlpInstructorXl: &vellum.HkunlpInstructorXlVectorizerRequest{
				Config: &vellum.InstructorVectorizerConfigRequest{
					InstructionDomain:           "",
					InstructionQueryTextType:    "plain_text",
					InstructionDocumentTextType: "plain_text",
				},
			},
		},
		Chunking: &vellum.DocumentIndexChunkingRequest{
			ChunkerName: "sentence-chunker",
			SentenceChunker: &vellum.SentenceChunkingRequest{
				ChunkerConfig: &vellum.SentenceChunkerConfigRequest{
					CharacterLimit: func() *int {
						v := 1000
						return &v
					}(),
					MinOverlapRatio: func() *float64 {
						v := 0.5
						return &v
					}(),
				},
			},
		},
	}

	request := vellum.DocumentIndexCreateRequest{
		Label:          documentIndexModel.Label.ValueString(),
		Name:           documentIndexModel.Name.ValueString(),
		IndexingConfig: &DefaultIndexingConfig,
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
