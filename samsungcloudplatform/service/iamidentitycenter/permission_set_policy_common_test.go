package iamidentitycenter

import (
	"strings"
	"testing"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
)

func TestExistingInlineContentArrayBranch(t *testing.T) {
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
	arr := []interface{}{map[string]interface{}{"content": doc}}
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr}),
	}

	got, ok := existingInlineContent(p)
	if !ok {
		t.Fatalf("expected inline content to be extracted from array branch, got ok=false")
	}
	if got == "" {
		t.Fatalf("expected non-empty extracted content")
	}
	if !strings.Contains(got, "s3:GetObject") {
		t.Fatalf("extracted content missing s3:GetObject: %q", got)
	}
	// extracted JSON must round-trip through normalization
	norm, err := normalizePolicyDocument(got)
	if err != nil {
		t.Fatalf("normalizePolicyDocument(%q) error: %v", got, err)
	}
	if norm == "" {
		t.Fatalf("normalized content is empty")
	}
}

func TestExistingInlineContentArrayStringElement(t *testing.T) {
	arr := []interface{}{map[string]interface{}{"content": "{\"Version\":\"2012-10-17\"}"}}
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr}),
	}

	got, ok := existingInlineContent(p)
	if !ok {
		t.Fatalf("expected ok=true for string element")
	}
	if got != `{"Version":"2012-10-17"}` {
		t.Fatalf("unexpected extracted content: %q", got)
	}
}

func TestExistingInlineContentMapBranch(t *testing.T) {
	m := map[string]interface{}{"content": "{\"Version\":\"2012-10-17\"}"}
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{MapmapOfStringAny: &m}),
	}

	got, ok := existingInlineContent(p)
	if !ok {
		t.Fatalf("expected ok=true for map branch")
	}
	if got != `{"Version":"2012-10-17"}` {
		t.Fatalf("unexpected extracted content: %q", got)
	}
}

func TestExistingInlineContentEmptyContent(t *testing.T) {
	arr := []interface{}{map[string]interface{}{"other": "value"}}
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr}),
	}
	if _, ok := existingInlineContent(p); ok {
		t.Fatalf("expected ok=false when no content key present")
	}
}

func TestExistingInlineContentNil(t *testing.T) {
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{}),
	}
	if _, ok := existingInlineContent(p); ok {
		t.Fatalf("expected ok=false when Contents has no data")
	}
}

func TestExistingInlineContentUnwrappedArray(t *testing.T) {
	// Console stores the inline policy document directly as the array element.
	doc := map[string]interface{}{
		"Version": "2024-07-01",
		"Statement": []interface{}{
			map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []interface{}{"objectstorage:CreateBucket"},
				"Resource": []interface{}{"*"},
				"Sid":      "VisualEditor0",
			},
		},
	}
	arr := []interface{}{doc}
	p := sdk.PolicyV1Dot2{
		Contents: *sdk.NewNullableContents(&sdk.Contents{ArrayOfAny: &arr}),
	}

	got, ok := existingInlineContent(p)
	if !ok {
		t.Fatalf("expected inline content extracted from unwrapped array element, got ok=false")
	}
	if !strings.Contains(got, "VisualEditor0") || !strings.Contains(got, "2024-07-01") {
		t.Fatalf("unexpected extracted content: %q", got)
	}
	norm, err := normalizePolicyDocument(got)
	if err != nil {
		t.Fatalf("normalizePolicyDocument(%q) error: %v", got, err)
	}
	if norm == "" {
		t.Fatalf("normalized content is empty")
	}
}

func TestBuildInlinePoliciesMapUnwrapped(t *testing.T) {
	doc := `{"Version":"2024-07-01","Statement":[{"Effect":"Allow","Action":["objectstorage:CreateBucket"],"Resource":["*"],"Sid":"VisualEditor0"}]}`
	got := buildInlinePoliciesMap([]string{doc})
	if len(got) != 1 {
		t.Fatalf("expected 1 element, got %d", len(got))
	}
	el := got[0]
	if _, hasContent := el["content"]; hasContent {
		t.Fatalf("expected unwrapped element, but got wrapped with content key: %#v", el)
	}
	v, ok := el["Version"]
	if !ok || v != "2024-07-01" {
		t.Fatalf("expected top-level Version in unwrapped element, got %#v", el)
	}
}
