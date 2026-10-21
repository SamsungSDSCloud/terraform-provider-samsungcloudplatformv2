package ske

import (
	"context"
	"sort"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// These tests keep the field lists used by ModifyPlan in sync with the schema.
// When one fails after a schema change, classify the new field in the code, not here
// (except nodepoolUpdatableAttributes below, which mirrors what Update sends).

func nestedAttributeNames(t *testing.T, s schema.Schema, name string) []string {
	t.Helper()
	nested, ok := s.Attributes[name].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("%s is not a single nested attribute", name)
	}
	var names []string
	for n := range nested.Attributes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// assertClassified checks that every nested field is either kept or volatile, exactly once,
// and that both lists only name existing fields.
func assertClassified(t *testing.T, fields []string, keep map[string]bool, volatile []string) {
	t.Helper()
	count := map[string]int{}
	for n := range keep {
		count[n]++
	}
	for _, n := range volatile {
		count[n]++
	}
	exists := map[string]bool{}
	for _, f := range fields {
		exists[f] = true
		switch count[f] {
		case 0:
			t.Errorf("nested field %q is not classified: add it to the keep map or the volatile list", f)
		case 2:
			t.Errorf("nested field %q is in both the keep map and the volatile list", f)
		}
	}
	for n := range count {
		if !exists[n] {
			t.Errorf("%q is listed but is not a nested field in the schema", n)
		}
	}
}

func TestClusterNestedFieldsClassified(t *testing.T) {
	fields := nestedAttributeNames(t, ClusterResourceSchema(), "cluster")
	assertClassified(t, fields, clusterNestedKeep(ske.ClusterResource{}, ske.ClusterResource{}), clusterNestedVolatile)
}

func TestNodepoolNestedFieldsClassified(t *testing.T) {
	fields := nestedAttributeNames(t, nodepoolSchema(), "nodepool")
	assertClassified(t, fields, nodepoolNestedKeep(nodepoolPlanAttributes{}, nodepoolPlanAttributes{}), nodepoolNestedVolatile)
}

// nodepoolUpdatableAttributes are the configurable attributes Update sends to the API.
// If the API gains support for another attribute, move it here from nodepoolCreateOnlyAttributes
// and handle it in Update; otherwise plan keeps rejecting the change.
var nodepoolUpdatableAttributes = []string{
	"kubernetes_version", "labels", "taints", "linked_resources", "preferred_ips",
	"is_auto_recovery", "is_auto_scale", "desired_node_count", "min_node_count", "max_node_count",
}

// nodepoolUncheckedAttributes are configurable but intentionally not checked for changes.
var nodepoolUncheckedAttributes = []string{
	"custom_image_id", // import stores the image name, which never equals the configured ID
}

func TestNodepoolRootAttributesClassified(t *testing.T) {
	s := nodepoolSchema()
	count := map[string]int{}
	for _, list := range [][]string{nodepoolUpdatableAttributes, nodepoolCreateOnlyAttributes, nodepoolUpgradeAttributes, nodepoolUncheckedAttributes} {
		for _, n := range list {
			count[n]++
		}
	}
	for name, a := range s.Attributes {
		if !a.IsRequired() && !a.IsOptional() {
			continue // computed only, not set by the user
		}
		switch count[name] {
		case 0:
			t.Errorf("attribute %q is not classified: add it to nodepoolCreateOnlyAttributes, or handle it in Update and add it to nodepoolUpdatableAttributes", name)
		case 1:
		default:
			t.Errorf("attribute %q is classified more than once", name)
		}
	}
	for n := range count {
		a, ok := s.Attributes[n]
		if !ok {
			t.Errorf("%q is listed but is not an attribute in the schema", n)
		} else if !a.IsRequired() && !a.IsOptional() {
			t.Errorf("%q is listed but cannot be set in the configuration", n)
		}
	}
}

// TestNodepoolPlanAttributesPaths checks that every path read by ModifyPlan exists in the schema.
func TestNodepoolPlanAttributesPaths(t *testing.T) {
	state := nodepoolState(t)
	plan, config := nodepoolUpdate(state, func(*ske.NodepoolResource) {})
	req := modifyPlanRequest(t, nodepoolSchema(), state, plan, config)
	if _, diags := getNodepoolPlanAttributes(context.Background(), req.Plan.GetAttribute); diags.HasError() {
		t.Fatalf("getNodepoolPlanAttributes: %v", diags)
	}
}
