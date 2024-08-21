package ml_model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"math/big"

	vellum "terraform-provider-vellum/internal/sdk"
)

func ValueIntPointer(value types.Int64) *int {
	if value.IsNull() {
		return nil
	}
	v := int(value.ValueInt64())
	return &v
}

func ValueFloat64Pointer(value types.Number) *float64 {
	if value.IsNull() {
		return nil
	}
	v, _ := value.ValueBigFloat().Float64()
	return &v
}

func NewOpenApiNumberPropertyRequest(ctx context.Context, openApiNumberProperty *TfOpenApiNumberProperty) (*vellum.OpenApiNumberPropertyRequest, diag.Diagnostics) {
	if openApiNumberProperty == nil {
		return nil, diag.Diagnostics{}
	}

	minimum := ValueFloat64Pointer(openApiNumberProperty.Minimum)
	maximum := ValueFloat64Pointer(openApiNumberProperty.Maximum)
	format := openApiNumberProperty.Format.ValueStringPointer()
	exclusiveMinimum := openApiNumberProperty.ExclusiveMinimum.ValueBoolPointer()
	exclusiveMaximum := openApiNumberProperty.ExclusiveMaximum.ValueBoolPointer()
	title := openApiNumberProperty.Title.ValueStringPointer()
	description := openApiNumberProperty.Description.ValueStringPointer()

	request := vellum.OpenApiNumberPropertyRequest{
		Minimum:          minimum,
		Maximum:          maximum,
		Format:           format,
		ExclusiveMinimum: exclusiveMinimum,
		ExclusiveMaximum: exclusiveMaximum,
		Title:            title,
		Description:      description,
	}

	return &request, nil
}

func NewOpenApiIntegerPropertyRequest(ctx context.Context, openApiIntegerProperty *TfOpenApiIntegerProperty) (*vellum.OpenApiIntegerPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiIntegerProperty == nil {
		return nil, diagnostics
	}

	minimum := ValueIntPointer(openApiIntegerProperty.Minimum)
	maximum := ValueIntPointer(openApiIntegerProperty.Maximum)
	exclusiveMinimum := openApiIntegerProperty.ExclusiveMinimum.ValueBoolPointer()
	exclusiveMaximum := openApiIntegerProperty.ExclusiveMaximum.ValueBoolPointer()
	title := openApiIntegerProperty.Title.ValueStringPointer()
	description := openApiIntegerProperty.Description.ValueStringPointer()

	request := vellum.OpenApiIntegerPropertyRequest{
		Minimum:          minimum,
		Maximum:          maximum,
		ExclusiveMinimum: exclusiveMinimum,
		ExclusiveMaximum: exclusiveMaximum,
		Title:            title,
		Description:      description,
	}

	return &request, diagnostics
}

func NewOpenApiStringPropertyRequest(ctx context.Context, openApiStringProperty *TfOpenApiStringProperty) (*vellum.OpenApiStringPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiStringProperty == nil {
		return nil, diagnostics
	}

	minLength := ValueIntPointer(openApiStringProperty.MinLength)
	maxLength := ValueIntPointer(openApiStringProperty.MaxLength)
	pattern := openApiStringProperty.Pattern.ValueStringPointer()
	format := openApiStringProperty.Format.ValueStringPointer()
	title := openApiStringProperty.Title.ValueStringPointer()
	description := openApiStringProperty.Description.ValueStringPointer()

	request := vellum.OpenApiStringPropertyRequest{
		MinLength:   minLength,
		MaxLength:   maxLength,
		Pattern:     pattern,
		Format:      format,
		Title:       title,
		Description: description,
	}

	return &request, diagnostics
}

func NewOpenApiBooleanPropertyRequest(ctx context.Context, openApiBooleanProperty *TfOpenApiBooleanProperty) (*vellum.OpenApiBooleanPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiBooleanProperty == nil {
		return nil, diagnostics
	}

	title := openApiBooleanProperty.Title.ValueStringPointer()
	description := openApiBooleanProperty.Description.ValueStringPointer()

	request := vellum.OpenApiBooleanPropertyRequest{
		Title:       title,
		Description: description,
	}

	return &request, diagnostics
}

func NewOpenApiArrayPropertyRequest(ctx context.Context, openApiArrayProperty *TfOpenApiArrayProperty) (*vellum.OpenApiArrayPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiArrayProperty == nil {
		return nil, diagnostics
	}

	minItems := ValueIntPointer(openApiArrayProperty.MinItems)
	maxItems := ValueIntPointer(openApiArrayProperty.MaxItems)
	uniqueItems := openApiArrayProperty.UniqueItems.ValueBoolPointer()
	items, _ := NewOpenApiPropertyRequest(ctx, openApiArrayProperty.Items)
	prefixItems := func() []*vellum.OpenApiPropertyRequest {
		if openApiArrayProperty.PrefixItems == nil {
			return nil
		}

		var prefixItems []*vellum.OpenApiPropertyRequest
		for _, rawPrefixItem := range openApiArrayProperty.PrefixItems {
			prefixItem, _ := NewOpenApiPropertyRequest(ctx, rawPrefixItem)
			prefixItems = append(prefixItems, prefixItem)
		}
		return prefixItems
	}()
	contains, _ := NewOpenApiPropertyRequest(ctx, openApiArrayProperty.Contains)
	minContains := ValueIntPointer(openApiArrayProperty.MinContains)
	maxContains := ValueIntPointer(openApiArrayProperty.MaxContains)
	title := openApiArrayProperty.Title.ValueStringPointer()
	description := openApiArrayProperty.Description.ValueStringPointer()

	request := vellum.OpenApiArrayPropertyRequest{
		MinItems:    minItems,
		MaxItems:    maxItems,
		UniqueItems: uniqueItems,
		Items:       items,
		PrefixItems: prefixItems,
		Contains:    contains,
		MinContains: minContains,
		MaxContains: maxContains,
		Title:       title,
		Description: description,
	}

	return &request, diagnostics
}

func NewOpenApiObjectPropertyRequest(ctx context.Context, openApiObjectProperty *TfOpenApiObjectProperty) (*vellum.OpenApiObjectPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiObjectProperty == nil {
		return nil, diagnostics
	}

	properties := func() map[string]*vellum.OpenApiPropertyRequest {
		if openApiObjectProperty.Properties == nil {
			return nil
		}

		customParameters := map[string]*vellum.OpenApiPropertyRequest{}
		for key, rawValue := range openApiObjectProperty.Properties {
			value, _ := NewOpenApiPropertyRequest(ctx, rawValue)
			customParameters[key] = value
		}
		return customParameters
	}()
	required := func() []string {
		if openApiObjectProperty.Required.IsNull() {
			return nil
		}

		var features []string
		for _, feature := range openApiObjectProperty.Required.Elements() {
			features = append(features, feature.(types.String).ValueString())
		}
		return features
	}()
	minProperties := ValueIntPointer(openApiObjectProperty.MinProperties)
	maxProperties := ValueIntPointer(openApiObjectProperty.MaxProperties)
	propertyNames, _ := NewOpenApiPropertyRequest(ctx, openApiObjectProperty.PropertyNames)
	patternProperties := func() map[string]*vellum.OpenApiPropertyRequest {
		if openApiObjectProperty.PatternProperties == nil {
			return nil
		}

		customParameters := map[string]*vellum.OpenApiPropertyRequest{}
		for key, rawValue := range openApiObjectProperty.PatternProperties {
			value, _ := NewOpenApiPropertyRequest(ctx, rawValue)
			customParameters[key] = value
		}
		return customParameters
	}()
	additionalProperties, _ := NewOpenApiPropertyRequest(ctx, openApiObjectProperty.AdditionalProperties)
	title := openApiObjectProperty.Title.ValueStringPointer()
	description := openApiObjectProperty.Description.ValueStringPointer()

	request := vellum.OpenApiObjectPropertyRequest{
		Properties:           properties,
		Required:             required,
		MinProperties:        minProperties,
		MaxProperties:        maxProperties,
		PropertyNames:        propertyNames,
		AdditionalProperties: additionalProperties,
		PatternProperties:    patternProperties,
		Title:                title,
		Description:          description,
	}

	return &request, diagnostics
}

func NewOpenApiPropertyRequest(ctx context.Context, openApiProperty *TfOpenApiProperty) (*vellum.OpenApiPropertyRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if openApiProperty == nil {
		return nil, diagnostics
	}

	if openApiProperty.Number != nil {
		f, _ := NewOpenApiNumberPropertyRequest(ctx, openApiProperty.Number)
		return &vellum.OpenApiPropertyRequest{OpenApiNumberPropertyRequest: f}, nil
	} else if openApiProperty.Integer != nil {
		f, _ := NewOpenApiIntegerPropertyRequest(ctx, openApiProperty.Integer)
		return &vellum.OpenApiPropertyRequest{OpenApiIntegerPropertyRequest: f}, nil
	} else if openApiProperty.String != nil {
		f, _ := NewOpenApiStringPropertyRequest(ctx, openApiProperty.String)
		return &vellum.OpenApiPropertyRequest{OpenApiStringPropertyRequest: f}, nil
	} else if openApiProperty.Boolean != nil {
		f, _ := NewOpenApiBooleanPropertyRequest(ctx, openApiProperty.Boolean)
		return &vellum.OpenApiPropertyRequest{OpenApiBooleanPropertyRequest: f}, nil
	}
	// TODO: Add object and array
	// 	Note: Terraform plugin doesn't support recursive schemas

	diagnostics.AddError(
		"Unsupported custom parameter type",
		fmt.Sprintf("Encountered property %s", openApiProperty),
	)
	return nil, diagnostics
}

func NewVellumMLModelCreateRequest(ctx context.Context, mlModelModel *TfMLModelResourceModel) (*vellum.MlModelCreateRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

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

	parameterConfig := vellum.MlModelParameterConfigRequest{
		Temperature: func() *vellum.OpenApiNumberPropertyRequest {
			f, _ := NewOpenApiNumberPropertyRequest(ctx, mlModelModel.ParameterConfig.Temperature)
			return f
		}(),
		MaxTokens: func() *vellum.OpenApiIntegerPropertyRequest {
			f, _ := NewOpenApiIntegerPropertyRequest(ctx, mlModelModel.ParameterConfig.MaxTokens)
			return f
		}(),
		Stop: func() *vellum.OpenApiArrayPropertyRequest {
			f, _ := NewOpenApiArrayPropertyRequest(ctx, mlModelModel.ParameterConfig.Stop)
			return f
		}(),
		TopK: func() *vellum.OpenApiIntegerPropertyRequest {
			f, _ := NewOpenApiIntegerPropertyRequest(ctx, mlModelModel.ParameterConfig.TopK)
			return f
		}(),
		TopP: func() *vellum.OpenApiNumberPropertyRequest {
			f, _ := NewOpenApiNumberPropertyRequest(ctx, mlModelModel.ParameterConfig.TopP)
			return f
		}(),
		FrequencyPenalty: func() *vellum.OpenApiNumberPropertyRequest {
			f, _ := NewOpenApiNumberPropertyRequest(ctx, mlModelModel.ParameterConfig.FrequencyPenalty)
			return f
		}(),
		PresencePenalty: func() *vellum.OpenApiNumberPropertyRequest {
			f, _ := NewOpenApiNumberPropertyRequest(ctx, mlModelModel.ParameterConfig.PresencePenalty)
			return f
		}(),
		LogitBias: func() *vellum.OpenApiObjectPropertyRequest {
			f, _ := NewOpenApiObjectPropertyRequest(ctx, mlModelModel.ParameterConfig.LogitBias)
			return f
		}(),
		CustomParameters: func() map[string]*vellum.OpenApiPropertyRequest {
			if mlModelModel.ParameterConfig.CustomParameters == nil {
				return nil
			}

			customParameters := map[string]*vellum.OpenApiPropertyRequest{}
			for key, rawValue := range mlModelModel.ParameterConfig.CustomParameters {
				value, _ := NewOpenApiPropertyRequest(ctx, rawValue)
				customParameters[key] = value
			}
			return customParameters
		}(),
	}

	displayTags := []vellum.MlModelDisplayTag{}
	for _, rawDisplayTag := range mlModelModel.DisplayConfig.Tags.Elements() {
		displayTag, _ := vellum.NewMlModelDisplayTagFromString(rawDisplayTag.(types.String).ValueString())
		displayTags = append(displayTags, displayTag)
	}

	displayConfig := vellum.MlModelDisplayConfigRequest{
		Label:       mlModelModel.DisplayConfig.Label.ValueString(),
		Description: mlModelModel.DisplayConfig.Description.ValueString(),
		Tags:        displayTags,
		DefaultDisplayPriority: func() *float64 {
			f, _ := mlModelModel.DisplayConfig.DefaultDisplayPriority.ValueBigFloat().Float64()
			return &f
		}(),
	}

	request := vellum.MlModelCreateRequest{
		Name:            mlModelModel.Name.ValueString(),
		Visibility:      &visibility,
		Family:          family,
		HostedBy:        hostedBy,
		DevelopedBy:     developedBy,
		ExecConfig:      &execConfig,
		ParameterConfig: &parameterConfig,
		DisplayConfig:   &displayConfig,
	}

	return &request, diagnostics
}

func NewTfOpenApiNumberPropertyModel(ctx context.Context, property *vellum.OpenApiNumberProperty) (*TfOpenApiNumberProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	minimum := func() basetypes.NumberValue {
		if property.Minimum == nil {
			return types.NumberNull()
		}
		return types.NumberValue(big.NewFloat((*property.Minimum)))
	}()
	maximum := func() basetypes.NumberValue {
		if property.Maximum == nil {
			return types.NumberNull()
		}
		return types.NumberValue(big.NewFloat((*property.Maximum)))
	}()

	propertyModel := &TfOpenApiNumberProperty{
		Minimum:          minimum,
		Maximum:          maximum,
		Format:           types.StringPointerValue(property.Format),
		ExclusiveMinimum: types.BoolPointerValue(property.ExclusiveMinimum),
		ExclusiveMaximum: types.BoolPointerValue(property.ExclusiveMaximum),
		Title:            types.StringPointerValue(property.Title),
		Description:      types.StringPointerValue(property.Description),
	}

	return propertyModel, diagnostics
}

func NewTfOpenApiIntegerPropertyModel(ctx context.Context, property *vellum.OpenApiIntegerProperty) (*TfOpenApiIntegerProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	minimum := func() basetypes.Int64Value {
		if property.Minimum == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.Minimum))
	}()
	maximum := func() basetypes.Int64Value {
		if property.Maximum == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.Maximum))
	}()

	propertyModel := &TfOpenApiIntegerProperty{
		Minimum:          minimum,
		Maximum:          maximum,
		ExclusiveMinimum: types.BoolPointerValue(property.ExclusiveMinimum),
		ExclusiveMaximum: types.BoolPointerValue(property.ExclusiveMaximum),
		Title:            types.StringPointerValue(property.Title),
		Description:      types.StringPointerValue(property.Description),
	}

	return propertyModel, diagnostics
}

func NewTfOpenApiStringPropertyModel(ctx context.Context, property *vellum.OpenApiStringProperty) (*TfOpenApiStringProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		// Shouldn't be null
		diagnostics.AddError(
			"OpenApiStringProperty is null",
			"OpenApiStringProperty is null",
		)
		return nil, diagnostics
	}

	minLength := func() basetypes.Int64Value {
		if property.MinLength == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MinLength))
	}()
	maxLength := func() basetypes.Int64Value {
		if property.MaxLength == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MaxLength))
	}()

	propertyModel := &TfOpenApiStringProperty{
		MinLength:   minLength,
		MaxLength:   maxLength,
		Pattern:     types.StringPointerValue(property.Pattern),
		Format:      types.StringPointerValue(property.Format),
		Title:       types.StringPointerValue(property.Title),
		Description: types.StringPointerValue(property.Description),
	}

	return propertyModel, diagnostics
}

func NewTfOpenApiBooleanPropertyModel(ctx context.Context, property *vellum.OpenApiBooleanProperty) (*TfOpenApiBooleanProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	propertyModel := &TfOpenApiBooleanProperty{
		Title:       types.StringPointerValue(property.Title),
		Description: types.StringPointerValue(property.Description),
	}

	return propertyModel, diagnostics
}

func NewTfOpenApiPropertyModel(ctx context.Context, property *vellum.OpenApiProperty) (*TfOpenApiProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	if property.Type == "string" {
		val, valDiagnostics := NewTfOpenApiStringPropertyModel(ctx, property.String)
		diagnostics = append(diagnostics, valDiagnostics...)
		return &TfOpenApiProperty{
			String: val,
		}, diagnostics
	} else if property.Type == "integer" {
		val, valDiagnostics := NewTfOpenApiIntegerPropertyModel(ctx, property.Integer)
		diagnostics = append(diagnostics, valDiagnostics...)
		return &TfOpenApiProperty{
			Integer: val,
		}, diagnostics
	} else if property.Type == "number" {
		val, valDiagnostics := NewTfOpenApiNumberPropertyModel(ctx, property.Number)
		diagnostics = append(diagnostics, valDiagnostics...)
		return &TfOpenApiProperty{
			Number: val,
		}, diagnostics
	} else if property.Type == "boolean" {
		val, valDiagnostics := NewTfOpenApiBooleanPropertyModel(ctx, property.Boolean)
		diagnostics = append(diagnostics, valDiagnostics...)
		return &TfOpenApiProperty{
			Boolean: val,
		}, diagnostics
	}

	diagnostics.AddError(
		"Unsupported custom parameter type",
		fmt.Sprintf("Encountered property %s", property),
	)
	return nil, diagnostics
}

func NewTfOpenApiArrayPropertyModel(ctx context.Context, property *vellum.OpenApiArrayProperty) (*TfOpenApiArrayProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	minItems := func() basetypes.Int64Value {
		if property.MinItems == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MinItems))
	}()
	maxItems := func() basetypes.Int64Value {
		if property.MaxItems == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MaxItems))
	}()

	items, _ := NewTfOpenApiPropertyModel(ctx, property.Items)
	contains, _ := NewTfOpenApiPropertyModel(ctx, property.Contains)

	minContains := func() basetypes.Int64Value {
		if property.MinContains == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MinContains))
	}()
	maxContains := func() basetypes.Int64Value {
		if property.MaxContains == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MaxContains))
	}()

	propertyModel := &TfOpenApiArrayProperty{
		MinItems:    minItems,
		MaxItems:    maxItems,
		UniqueItems: types.BoolPointerValue(property.UniqueItems),
		Items:       items,
		PrefixItems: func() []*TfOpenApiProperty {
			var prefixItems []*TfOpenApiProperty
			for _, rawPrefixItem := range property.PrefixItems {
				prefixItem, _ := NewTfOpenApiPropertyModel(ctx, rawPrefixItem)
				prefixItems = append(prefixItems, prefixItem)
			}
			return prefixItems
		}(),
		Contains:    contains,
		MinContains: minContains,
		MaxContains: maxContains,
		Title:       types.StringPointerValue(property.Title),
		Description: types.StringPointerValue(property.Description),
	}

	return propertyModel, diagnostics
}

func NewTfOpenApiObjectPropertyModel(ctx context.Context, property *vellum.OpenApiObjectProperty) (*TfOpenApiObjectProperty, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	properties := func() map[string]*TfOpenApiProperty {
		if property.Properties == nil {
			return nil
		}

		customParameters := map[string]*TfOpenApiProperty{}
		for key, rawValue := range property.Properties {
			value, valueDiagnostics := NewTfOpenApiPropertyModel(ctx, rawValue)
			diagnostics = append(diagnostics, valueDiagnostics...)
			customParameters[key] = value
		}
		return customParameters
	}()

	required := func() types.List {
		if property.Required == nil {
			return types.ListNull(types.StringType)
		}

		return types.ListValueMust(
			types.StringType,
			func() []attr.Value {
				var features []attr.Value
				for _, feature := range property.Required {
					features = append(features, types.StringValue(feature))
				}
				return features
			}(),
		)
	}()

	minProperties := func() basetypes.Int64Value {
		if property.MinProperties == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MinProperties))
	}()
	maxProperties := func() basetypes.Int64Value {
		if property.MaxProperties == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*property.MaxProperties))
	}()

	propertyNames := func() *TfOpenApiProperty {
		if property.PropertyNames == nil {
			return nil
		}

		value, _ := NewTfOpenApiPropertyModel(ctx, property.PropertyNames)
		return value
	}()

	additionalProperties := func() *TfOpenApiProperty {
		if property.AdditionalProperties == nil {
			return nil
		}

		value, _ := NewTfOpenApiPropertyModel(ctx, property.AdditionalProperties)
		return value
	}()

	patternProperties := func() map[string]*TfOpenApiProperty {
		if property.PatternProperties == nil {
			return nil
		}

		customParameters := map[string]*TfOpenApiProperty{}
		for key, rawValue := range property.PatternProperties {
			value, _ := NewTfOpenApiPropertyModel(ctx, rawValue)
			customParameters[key] = value
		}
		return customParameters
	}()

	propertyModel := &TfOpenApiObjectProperty{
		Properties:           properties,
		Required:             required,
		MinProperties:        minProperties,
		MaxProperties:        maxProperties,
		PropertyNames:        propertyNames,
		AdditionalProperties: additionalProperties,
		PatternProperties:    patternProperties,
		Title:                types.StringPointerValue(property.Title),
		Description:          types.StringPointerValue(property.Description),
	}

	return propertyModel, nil
}

func NewTfMlModelTokenizerConfig(ctx context.Context, property *vellum.MlModelTokenizerConfig) (*TfMlModelTokenizerConfig, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	if property == nil {
		return nil, diagnostics
	}

	if property.Type == "HUGGING_FACE" {
		return &TfMlModelTokenizerConfig{
			Type: types.StringValue("HUGGING_FACE"),
			HuggingFace: &TfHuggingFaceTokenizerConfig{
				Name: types.StringValue(property.HuggingFace.Name),
				Path: types.StringPointerValue(property.HuggingFace.Path),
			},
		}, diagnostics
	} else if property.Type == "TIKTOKEN" {
		return &TfMlModelTokenizerConfig{
			Type: types.StringValue("TIKTOKEN"),
			Tiktoken: &TiktokenTokenizerConfig{
				Name: types.StringValue(property.Tiktoken.Name),
			},
		}, diagnostics
	}

	return nil, diagnostics
}

func NewTfMLModelModel(ctx context.Context, model *TfMLModelResourceModel, mlModel *vellum.MlModelRead) (*TfMLModelResourceModel, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	temperature, _ := NewTfOpenApiNumberPropertyModel(ctx, mlModel.ParameterConfig.Temperature)
	maxTokens, _ := NewTfOpenApiIntegerPropertyModel(ctx, mlModel.ParameterConfig.MaxTokens)
	stop, _ := NewTfOpenApiArrayPropertyModel(ctx, mlModel.ParameterConfig.Stop)
	topP, _ := NewTfOpenApiNumberPropertyModel(ctx, mlModel.ParameterConfig.TopP)
	topK, _ := NewTfOpenApiIntegerPropertyModel(ctx, mlModel.ParameterConfig.TopK)
	frequencyPenalty, _ := NewTfOpenApiNumberPropertyModel(ctx, mlModel.ParameterConfig.FrequencyPenalty)
	presencePenalty, _ := NewTfOpenApiNumberPropertyModel(ctx, mlModel.ParameterConfig.PresencePenalty)
	logitBias, _ := NewTfOpenApiObjectPropertyModel(ctx, mlModel.ParameterConfig.LogitBias)

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
			ForceSystemCredentials: types.BoolPointerValue(mlModel.ExecConfig.ForceSystemCredentials),
			TokenizerConfig: func() *TfMlModelTokenizerConfig {
				value, _ := NewTfMlModelTokenizerConfig(ctx, mlModel.ExecConfig.TokenizerConfig)
				return value
			}(),
		},
		ParameterConfig: &TfMLModelParameterConfig{
			Temperature:      temperature,
			MaxTokens:        maxTokens,
			Stop:             stop,
			TopP:             topP,
			TopK:             topK,
			FrequencyPenalty: frequencyPenalty,
			PresencePenalty:  presencePenalty,
			LogitBias:        logitBias,
			CustomParameters: func() map[string]*TfOpenApiProperty {
				if mlModel.ParameterConfig.CustomParameters == nil {
					return nil
				}

				customParameters := map[string]*TfOpenApiProperty{}
				for key, rawValue := range mlModel.ParameterConfig.CustomParameters {
					value, valueDiagnostics := NewTfOpenApiPropertyModel(ctx, rawValue)
					diagnostics = append(diagnostics, valueDiagnostics...)

					customParameters[key] = value
				}
				return customParameters
			}(),
		},
		DisplayConfig: &TfMLModelDisplayConfig{
			Label:       types.StringValue(mlModel.DisplayConfig.Label),
			Description: types.StringValue(mlModel.DisplayConfig.Description),
			Tags: types.ListValueMust(
				types.StringType,
				func() []attr.Value {
					var tags []attr.Value
					for _, tag := range mlModel.DisplayConfig.Tags {
						tags = append(tags, types.StringValue(string(tag.Value)))
					}
					return tags
				}(),
			),
			DefaultDisplayPriority: types.NumberValue(big.NewFloat(*mlModel.DisplayConfig.DefaultDisplayPriority)),
		},
	}

	return mlModelModel, diagnostics
}

func NewTfMLModelDataSourceModel(ctx context.Context, mlModel *vellum.MlModelRead) (*TfMLModelDataSourceModel, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	mlModelModel := &TfMLModelDataSourceModel{
		Id:          types.StringValue(mlModel.Id),
		Name:        types.StringValue(mlModel.Name),
		Visibility:  types.StringValue(string(*mlModel.Visibility)),
		HostedBy:    types.StringValue(string(mlModel.HostedBy)),
		DevelopedBy: types.StringValue(string(mlModel.DevelopedBy.Value)),
		Family:      types.StringValue(string(mlModel.Family.Value)),
	}

	return mlModelModel, diagnostics
}
