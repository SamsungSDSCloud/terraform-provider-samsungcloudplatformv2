package baremetal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"net"
)

type IpStringValidator struct{}

// Description returns a description of the validator.
func (v IpStringValidator) Description(ctx context.Context) string {
	return "Validates that the IP is ip format and not empty value"
}

// MarkdownDescription returns a markdown description of the validator.
func (v IpStringValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that the IP is ip format and not empty value"
}

// ValidateString performs the validation logic.
func (v IpStringValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return // Don't validate unknown or null values
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		resp.Diagnostics.AddError(
			"Invalid IP value",
			"The ip value must not be empty. For automatic allocation, please call without the corresponding field.",
		)
		return
	}

	if net.ParseIP(value) == nil {
		resp.Diagnostics.AddError(
			"Invalid IP value",
			value+" is invalid IP format.",
		)
	}
}
