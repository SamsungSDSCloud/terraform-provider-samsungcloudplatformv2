package iamidentitycenter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Permission set policy categories as returned by the SCP IAM Identity Center
// policies API. These are also used as the GetPolicies filter values.
const (
	permissionSetPolicyCategoryCustom  = "CUSTOM_POLICY"
	permissionSetPolicyCategoryManaged = "MANAGED_POLICY"
	permissionSetPolicyCategoryInline  = "INLINE_POLICY"
)

// composePermissionSetPolicyID builds the Terraform resource id for a single
// permission set policy. The identifier chunk differs per policy type:
// custom -> policy name, managed -> IAM managed policy id, inline -> content hash.
func composePermissionSetPolicyID(permissionSetID, instanceID, identifier string) string {
	return strings.Join([]string{permissionSetID, instanceID, identifier}, ",")
}

// parsePermissionSetPolicyID splits a composed id back into its parts. The
// identifier chunk (3rd) may itself contain commas, so it is left untouched.
func parsePermissionSetPolicyID(id string) (permissionSetID, instanceID, identifier string, err error) {
	parts := strings.SplitN(id, ",", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid permission set policy id %q: expected <permission-set-id>,<instance-id>,<identifier>", id)
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid permission set policy id %q: empty segment", id)
	}
	return parts[0], parts[1], parts[2], nil
}

// fetchPermissionSetPolicies retrieves the policies of a permission set,
// optionally filtered by category. An empty category means "all policies".
func fetchPermissionSetPolicies(ctx context.Context, clients *client.SCPClient, permissionSetID, instanceID, category string) ([]sdk.PolicyV1Dot2, diag.Diagnostics) {
	var diags diag.Diagnostics
	result, err := clients.IamIdentityCenter.GetPolicies(ctx, permissionSetID, instanceID, category, "", 0, 0)
	if err != nil {
		diags.AddError(
			"Failed to read IAM Identity Center Permission Set Policies",
			err.Error(),
		)
		return nil, diags
	}
	if result == nil {
		return nil, diags
	}
	return result.Policies, diags
}

// findCustomPolicy returns the first CUSTOM_POLICY whose name matches, or nil.
func findCustomPolicy(policies []sdk.PolicyV1Dot2, name string) *sdk.PolicyV1Dot2 {
	for i := range policies {
		if string(policies[i].Category) != permissionSetPolicyCategoryCustom {
			continue
		}
		if policies[i].GetName() == name {
			return &policies[i]
		}
	}
	return nil
}

// managedPolicyIAMID extracts the IAM managed policy id that the backend stores
// in the contents of a managed policy. For MANAGED_POLICY the contents map
// carries the IAM id under the "id" key (distinct from the IDC-internal
// policy.Id used for deletes).
func managedPolicyIAMID(p sdk.PolicyV1Dot2) (string, bool) {
	if !p.Contents.IsSet() || p.Contents.Get() == nil || p.Contents.Get().MapmapOfStringAny == nil {
		return "", false
	}
	id, ok := (*p.Contents.Get().MapmapOfStringAny)["id"].(string)
	return id, ok
}

// managedPolicyDisplayName extracts the display name stored in the contents of
// a managed policy (the "name" key next to the IAM id).
func managedPolicyDisplayName(p sdk.PolicyV1Dot2) (string, bool) {
	if !p.Contents.IsSet() || p.Contents.Get() == nil || p.Contents.Get().MapmapOfStringAny == nil {
		return "", false
	}
	name, ok := (*p.Contents.Get().MapmapOfStringAny)["name"].(string)
	return name, ok
}

// findManagedPolicy returns the first MANAGED_POLICY whose contents id (the IAM
// managed policy id) matches iamPolicyID, or nil.
func findManagedPolicy(policies []sdk.PolicyV1Dot2, iamPolicyID string) *sdk.PolicyV1Dot2 {
	for i := range policies {
		if string(policies[i].Category) != permissionSetPolicyCategoryManaged {
			continue
		}
		if id, ok := managedPolicyIAMID(policies[i]); ok && id == iamPolicyID {
			return &policies[i]
		}
	}
	return nil
}

// existingInlineContent returns the raw policy document stored in the contents
// of an inline policy. The backend stores the JSON document either as an array
// of {"content": <doc>} elements (legacy wrapped shape), under the "content"
// key of a map (MapmapOfStringAny), or - matching the console's native format -
// as the document itself in an array element (unwrapped shape). All shapes are
// handled, with the unwrapped array element only treated as a document when it
// resembles a policy (has a Version or Statement key).
func existingInlineContent(p sdk.PolicyV1Dot2) (string, bool) {
	if !p.Contents.IsSet() || p.Contents.Get() == nil {
		return "", false
	}
	contents := p.Contents.Get()

	if contents.ArrayOfAny != nil {
		for _, el := range *contents.ArrayOfAny {
			m, ok := el.(map[string]interface{})
			if !ok {
				continue
			}
			if raw, ok := m["content"]; ok {
				if s, ok := raw.(string); ok {
					return s, true
				}
				if b, err := json.Marshal(raw); err == nil {
					return string(b), true
				}
				continue
			}
			if _, hasVersion := m["Version"]; !hasVersion {
				if _, hasStatement := m["Statement"]; !hasStatement {
					continue
				}
			}
			if b, err := json.Marshal(m); err == nil {
				return string(b), true
			}
		}
		return "", false
	}

	if contents.MapmapOfStringAny == nil {
		return "", false
	}
	c, ok := (*contents.MapmapOfStringAny)["content"].(string)
	return c, ok
}

// findInlinePolicy returns the first INLINE_POLICY whose normalized document
// equals the given normalized document, or nil.
func findInlinePolicy(policies []sdk.PolicyV1Dot2, normalizedDoc string) *sdk.PolicyV1Dot2 {
	for i := range policies {
		if string(policies[i].Category) != permissionSetPolicyCategoryInline {
			continue
		}
		raw, ok := existingInlineContent(policies[i])
		if !ok {
			continue
		}
		norm, err := normalizePolicyDocument(raw)
		if err != nil {
			continue
		}
		if norm == normalizedDoc {
			return &policies[i]
		}
	}
	return nil
}

// normalizePolicyDocument canonically re-marshals a JSON document so that key
// ordering and formatting differences do not affect equality. It validates that
// the document is valid JSON.
func normalizePolicyDocument(doc string) (string, error) {
	var parsed interface{}
	if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
		return "", fmt.Errorf("invalid policy document JSON: %w", err)
	}
	canonical, err := json.Marshal(parsed)
	if err != nil {
		return "", err
	}
	return string(canonical), nil
}

// inlinePolicyIdentifier returns a deterministic hash of a normalized policy
// document, used as the identifier chunk of an inline policy resource id.
func inlinePolicyIdentifier(normalizedDoc string) string {
	sum := sha256.Sum256([]byte(normalizedDoc))
	return hex.EncodeToString(sum[:])
}

// buildInlinePoliciesMap converts inline policy JSON strings into the per-policy
// content shape expected by the set policies request. Documents are stored
// unwrapped (the document itself as each array element) to match the console's
// native inline policy format so provider-created policies stay editable in the
// console basic-mode editor. Non-object JSON falls back to the wrapped shape.
func buildInlinePoliciesMap(inlinePolicies []string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, ip := range inlinePolicies {
		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(ip), &doc); err == nil {
			result = append(result, doc)
		} else {
			result = append(result, map[string]interface{}{"content": ip})
		}
	}
	return result
}
