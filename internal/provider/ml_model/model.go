package ml_model

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
)

func NewVellumMLModelCreateRequest(ctx context.Context, mlModelModel *TfMLModelResourceModel) (*vellum.MlModelCreateRequest, diag.Diagnostics) {
	visibility, _ := vellum.NewVisibilityEnumFromString(mlModelModel.Visibility.ValueString())
	hostedBy, _ := vellum.NewHostedByEnumFromString(mlModelModel.HostedBy.ValueString())
	developedBy, _ := vellum.NewMlModelDeveloperFromString(mlModelModel.DevelopedBy.ValueString())
	family, _ := vellum.NewMlModelFamilyFromString(mlModelModel.Family.ValueString())

	features := []vellum.MlModelFeature{}
	for _, feature := range mlModelModel.ExecConfig.Features.Elements() {
		feature, _ := vellum.NewMlModelFeatureFromString(feature.(types.String).ValueString())
		features = append(features, feature)
	}

	metadata := map[string]interface{}{}
	for key, tfvalue := range mlModelModel.ExecConfig.Metadata.Elements() {
		value := tfvalue.(types.String).ValueString()
		var v interface{}
		if err := json.Unmarshal([]byte(value), &v); err != nil {
			metadata[key] = value
		} else {
			metadata[key] = v
		}
	}

	execConfig := vellum.MlModelExecConfigRequest{
		ModelIdentifier: mlModelModel.ExecConfig.ModelIdentifier.ValueString(),
		BaseUrl:         mlModelModel.ExecConfig.BaseUrl.ValueString(),
		Features:        features,
		Metadata:        metadata,
	}

	request := vellum.MlModelCreateRequest{
		Name:        mlModelModel.Name.ValueString(),
		Visibility:  &visibility,
		Family:      family,
		HostedBy:    hostedBy,
		DevelopedBy: developedBy,
		ExecConfig:  &execConfig,
	}

	return &request, nil
}

func NewTfMLModelModel(ctx context.Context, model *TfMLModelResourceModel, mlModel *vellum.MlModelRead) (*TfMLModelResourceModel, diag.Diagnostics) {
	mlModelModel := &TfMLModelResourceModel{
		Id:          types.StringValue(mlModel.Id),
		Name:        types.StringValue(mlModel.Name),
		Visibility:  types.StringValue(string(*mlModel.Visibility)),
		HostedBy:    types.StringValue(string(mlModel.HostedBy)),
		DevelopedBy: types.StringValue(string(mlModel.DevelopedBy.Value)),
		Family:      types.StringValue(string(mlModel.Family.Value)),
		ExecConfig: TfMLModelExecConfig{
			ModelIdentifier: types.StringValue(mlModel.ExecConfig.ModelIdentifier),
			BaseUrl:         types.StringValue(mlModel.ExecConfig.BaseUrl),
			Features: types.ListValueMust(
				types.StringType,
				func() []attr.Value {
					var features []attr.Value
					for _, feature := range mlModel.ExecConfig.Features {
						features = append(features, types.StringValue(string(feature)))
					}
					return features
				}(),
			),
			Metadata: types.MapValueMust(
				types.StringType,
				func() map[string]attr.Value {
					metadata := map[string]attr.Value{}
					for key, value := range mlModel.ExecConfig.Metadata {
						metadata[key] = types.StringValue(value)
					}
					return metadata
				}(),
			),
		},
	}

	return mlModelModel, nil
}

func NewTfMLModelDataSourceModel(ctx context.Context, mlModel *vellum.MlModelRead) (*TfMLModelDataSourceModel, diag.Diagnostics) {
	mlModelModel := &TfMLModelDataSourceModel{
		Id:          types.StringValue(mlModel.Id),
		Name:        types.StringValue(mlModel.Name),
		Visibility:  types.StringValue(string(*mlModel.Visibility)),
		HostedBy:    types.StringValue(string(mlModel.HostedBy)),
		DevelopedBy: types.StringValue(string(mlModel.DevelopedBy.Value)),
		Family:      types.StringValue(string(mlModel.Family.Value)),
	}

	return mlModelModel, nil
}
