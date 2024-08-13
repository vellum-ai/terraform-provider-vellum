package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

func ValueIntPointer(value types.Int64) *int {
	if value.IsNull() {
		return nil
	}
	v := int(value.ValueInt64())
	return &v
}
