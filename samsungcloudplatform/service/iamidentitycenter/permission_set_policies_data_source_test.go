package iamidentitycenter

import (
	"strings"
	"testing"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
)

func TestInlinePolicyDocumentUnset(t *testing.T) {
	var contents sdk.NullableContents
	if got := inlinePolicyDocument(contents); got != "" {
		t.Fatalf("expected empty string for unset contents, got %q", got)
	}
}

func TestInlinePolicyDocumentNilContents(t *testing.T) {
	contents := *sdk.NewNullableContents(&sdk.Contents{})
	if got := inlinePolicyDocument(contents); got != "" {
		t.Fatalf("expected empty string for empty contents, got %q", got)
	}
}

func TestInlinePolicyDocumentArraySingleElement(t *testing.T) {
	doc := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []interface{}{
			map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []interface{}{"s3:GetObject"},
				"Resource": []interface{}{"*"},
			},
		},
	}
	arr := []interface{}{doc}
	contents := *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr})

	got := inlinePolicyDocument(contents)
	if got == "" {
		t.Fatalf("expected non-empty inline policy document")
	}
	if !strings.Contains(got, "s3:GetObject") {
		t.Fatalf("inline policy missing s3:GetObject: %q", got)
	}
}

func TestInlinePolicyDocumentUnwrappedArrayMultipleElements(t *testing.T) {
	// With multiple unwrapped elements the resource parser returns the first
	// document (an inline policy is a singleton), not a JSON array.
	arr := []interface{}{
		map[string]interface{}{"Version": "2012-10-17"},
		map[string]interface{}{"Version": "2024-07-01"},
	}
	contents := *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr})

	got := inlinePolicyDocument(contents)
	if got != `{"Version":"2012-10-17"}` {
		t.Fatalf("expected first document, got %q", got)
	}
}

func TestInlinePolicyDocumentWrappedArray(t *testing.T) {
	// Legacy wrapped shape: array of {"content": <raw doc string>}.
	raw := `{"Version":"2012-10-17"}`
	arr := []interface{}{map[string]interface{}{"content": raw}}
	contents := *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr})

	got := inlinePolicyDocument(contents)
	if got != raw {
		t.Fatalf("expected inner document from wrapped array, got %q", got)
	}
}

func TestInlinePolicyDocumentWrappedMap(t *testing.T) {
	// Wrapped map shape: contents map carries the document under "content".
	m := map[string]interface{}{"content": `{"Version":"2012-10-17"}`}
	contents := *sdk.NewNullableContents(&sdk.Contents{MapmapOfStringAny: &m})

	got := inlinePolicyDocument(contents)
	if got != `{"Version":"2012-10-17"}` {
		t.Fatalf("expected inner document from wrapped map, got %q", got)
	}
}

func TestReferenceOf(t *testing.T) {
	name := "MyManagedPolicy"
	policy := sdk.PolicyV1Dot2{
		Id:   "policy-id-1",
		Name: *sdk.NewNullableString(&name),
	}

	got := referenceOf(policy)
	if got.Id.ValueString() != "policy-id-1" {
		t.Fatalf("unexpected id: %q", got.Id.ValueString())
	}
	if got.Name.ValueString() != "MyManagedPolicy" {
		t.Fatalf("unexpected name: %q", got.Name.ValueString())
	}
}

func TestReferenceOfEmptyName(t *testing.T) {
	policy := sdk.PolicyV1Dot2{Id: "policy-id-2"}

	got := referenceOf(policy)
	if got.Id.ValueString() != "policy-id-2" {
		t.Fatalf("unexpected id: %q", got.Id.ValueString())
	}
	if !got.Name.IsNull() {
		t.Fatalf("expected null name when policy name is unset, got %#v", got.Name)
	}
}

func TestNullableStringSet(t *testing.T) {
	raw := "hello"
	value := sdk.NewNullableString(&raw)
	got := nullableString(*value)
	if got.ValueString() != "hello" {
		t.Fatalf("unexpected string value: %q", got.ValueString())
	}
}

func TestNullableStringUnset(t *testing.T) {
	got := nullableString(sdk.NullableString{})
	if !got.IsNull() {
		t.Fatalf("expected null for unset nullable string, got %#v", got)
	}
}

func TestCollectPolicyReferenceRoutesByCategory(t *testing.T) {
	customName := "custom-policy"
	managedName := "managed-policy"

	policies := []sdk.PolicyV1Dot2{
		{Category: "CUSTOM_POLICY", Id: "custom-id", Name: *sdk.NewNullableString(&customName)},
		{Category: "MANAGED_POLICY", Id: "managed-id", Name: *sdk.NewNullableString(&managedName)},
	}

	var custom []permissionSetPolicyReference
	var managed []permissionSetPolicyReference
	var inline string

	for _, p := range policies {
		collectPolicyReference(p, &custom, &managed, &inline)
	}

	if len(custom) != 1 || custom[0].Id.ValueString() != "custom-id" || custom[0].Name.ValueString() != "custom-policy" {
		t.Fatalf("unexpected custom policies: %#v", custom)
	}
	if len(managed) != 1 || managed[0].Id.ValueString() != "managed-id" || managed[0].Name.ValueString() != "managed-policy" {
		t.Fatalf("unexpected managed policies: %#v", managed)
	}
	if custom[0].Id.ValueString() == managed[0].Id.ValueString() {
		t.Fatalf("custom and managed references must not share identities")
	}
}

func TestCollectPolicyReferenceSkipsEmptyId(t *testing.T) {
	var custom []permissionSetPolicyReference
	var managed []permissionSetPolicyReference
	var inline string

	collectPolicyReference(sdk.PolicyV1Dot2{Category: "CUSTOM_POLICY"}, &custom, &managed, &inline)
	collectPolicyReference(sdk.PolicyV1Dot2{Category: "MANAGED_POLICY"}, &managed, &managed, &inline)

	if len(custom) != 0 || len(managed) != 0 {
		t.Fatalf("expected empty references for policies without ids, got custom=%#v managed=%#v", custom, managed)
	}
}
