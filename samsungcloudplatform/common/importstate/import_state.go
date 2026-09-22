package importstate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const separator = ":"

// ImportState handles import for composite keys (2+ fields).
func ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
	fieldPaths ...path.Path,
) {
	if len(fieldPaths) == 0 {
		resp.Diagnostics.AddError(
			"Invalid ImportState call",
			"at least one fieldPath must be provided",
		)
		return
	}

	parts, err := parseImportID(req.ID, len(fieldPaths))
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}

	for i, p := range fieldPaths {
		resp.Diagnostics.Append(
			resp.State.SetAttribute(ctx, p, parts[i])...,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// parseImportID splits the import ID into field parts.
func parseImportID(importID string, fieldCount int) (parts []string, err error) {
	if strings.TrimSpace(importID) == "" {
		return nil, fmt.Errorf("import ID must not be empty")
	}

	parts, err = splitAndValidate(importID, fieldCount)
	return parts, err
}

// splitAndValidate splits importID by separator into exactly count parts and validates them.
func splitAndValidate(importID string, count int) ([]string, error) {
	parts := strings.SplitN(importID, separator, count)
	if len(parts) != count {
		return nil, fmt.Errorf(
			"invalid import ID %q: expected %d segment(s) separated by %q, got %d",
			importID, count, separator, len(parts),
		)
	}
	if err := validateSegments(importID, parts); err != nil {
		return nil, err
	}
	return parts, nil
}

// validateSegments checks that no segment is empty or whitespace-only.
func validateSegments(importID string, segments []string) error {
	for i, s := range segments {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf(
				"invalid import ID %q: segment %d is empty",
				importID, i+1,
			)
		}
	}
	return nil
}
