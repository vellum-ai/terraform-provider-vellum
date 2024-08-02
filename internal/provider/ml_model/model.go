package ml_model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vellum "terraform-provider-vellum/internal/sdk"
)

func NewVellumMLModelCreateRequest(ctx context.Context, mlModelModel *TfMLModelResourceModel) (*vellum.MlModelCreateRequest, diag.Diagnostics) {
	visibility, _ := vellum.NewVisibilityEnumFromString(mlModelModel.Visibility.ValueString())
	hostedBy, _ := vellum.NewHostedByEnumFromString(mlModelModel.HostedBy.ValueString())
	developedBy, _ := vellum.NewMlModelDeveloperFromString(mlModelModel.DevelopedBy.ValueString())
	family, _ := vellum.NewMlModelFamilyFromString(mlModelModel.Family.ValueString())

	// TODO: Pass in actual values rather than dummy keys
	// Create an empty slice for features
	features := []vellum.MlModelFeature{"CHAT_MESSAGE_SYSTEM"}
	execConfig := vellum.MlModelExecConfigRequest{
		ModelIdentifier: "test",
		BaseUrl:         "http://localhost:8080",
		Metadata:        map[string]interface{}{"key": "value"},
		Features:        features,
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
