package ske

import (
	"context"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// modifyPlanRequest builds a ModifyPlan request for an update from state to plan.
func modifyPlanRequest(t *testing.T, s schema.Schema, state, plan, config any) resource.ModifyPlanRequest {
	t.Helper()
	ctx := context.Background()
	empty := tftypes.NewValue(s.Type().TerraformType(ctx), nil)

	req := resource.ModifyPlanRequest{
		State:  tfsdk.State{Schema: s, Raw: empty},
		Plan:   tfsdk.Plan{Schema: s, Raw: empty},
		Config: tfsdk.Config{Schema: s, Raw: empty},
	}
	if d := req.State.Set(ctx, state); d.HasError() {
		t.Fatalf("set state: %v", d)
	}
	if d := req.Plan.Set(ctx, plan); d.HasError() {
		t.Fatalf("set plan: %v", d)
	}
	cfgPlan := tfsdk.Plan{Schema: s, Raw: empty}
	if d := cfgPlan.Set(ctx, config); d.HasError() {
		t.Fatalf("set config: %v", d)
	}
	req.Config.Raw = cfgPlan.Raw
	return req
}

// modifyPlan runs ModifyPlan of r for an update from state to plan and returns the resulting plan.
func modifyPlan(t *testing.T, r resource.ResourceWithModifyPlan, s schema.Schema, state, plan, config any) tfsdk.Plan {
	t.Helper()
	req := modifyPlanRequest(t, s, state, plan, config)
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ModifyPlan: %v", resp.Diagnostics)
	}
	return resp.Plan
}

func nestedAttrs(t *testing.T, obj types.Object) map[string]attr.Value {
	t.Helper()
	if obj.IsUnknown() || obj.IsNull() {
		t.Fatalf("nested object is not known: %v", obj)
	}
	return obj.Attributes()
}

func assertKnown(t *testing.T, attrs map[string]attr.Value, names ...string) {
	t.Helper()
	for _, n := range names {
		if attrs[n].IsUnknown() {
			t.Errorf("%s is unknown, want state value", n)
		}
	}
}

func assertUnknown(t *testing.T, attrs map[string]attr.Value, names ...string) {
	t.Helper()
	for _, n := range names {
		if !attrs[n].IsUnknown() {
			t.Errorf("%s = %v, want unknown", n, attrs[n])
		}
	}
}

func nodepoolSchema() schema.Schema {
	resp := &resource.SchemaResponse{}
	(&skeNodepoolResource{}).Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

func nodepoolState(t *testing.T) ske.NodepoolResource {
	t.Helper()
	nested := ske.Nodepool{
		Id:               types.StringValue("np-1"),
		Name:             types.StringValue("scale"),
		AutoScaleEnabled: types.BoolValue(false),
		CurrentNodeCount: types.Int32Value(1),
		DesiredNodeCount: types.Int32Value(1),
		MinNodeCount:     types.Int32Value(1),
		MaxNodeCount:     types.Int32Value(1),
		Status:           types.StringValue("Running"),
		ModifiedAt:       types.StringValue("2026-10-15T00:00:00Z"),
		Labels:           []ske.Label{{Key: types.StringValue("a"), Value: types.StringValue("1")}},
	}
	obj, d := types.ObjectValueFrom(context.Background(), nested.AttributeTypes(), nested)
	if d.HasError() {
		t.Fatalf("nested object: %v", d)
	}
	return ske.NodepoolResource{
		Id:               types.StringValue("np-1"),
		Name:             types.StringValue("scale"),
		IsAutoScale:      types.BoolValue(false),
		IsAutoRecovery:   types.BoolValue(true),
		DesiredNodeCount: types.Int32Value(1),
		MinNodeCount:     types.Int32Value(1),
		MaxNodeCount:     types.Int32Value(1),
		Labels:           []ske.Label{{Key: types.StringValue("a"), Value: types.StringValue("1")}},
		Nodepool:         obj,
	}
}

// update plan as Terraform builds it: computed values not set in the configuration are unknown.
func nodepoolUpdate(state ske.NodepoolResource, edit func(*ske.NodepoolResource)) (plan, config ske.NodepoolResource) {
	config = state
	config.Id = types.StringNull()
	config.MinNodeCount = types.Int32Null()
	config.MaxNodeCount = types.Int32Null()
	config.Nodepool = types.ObjectNull(ske.Nodepool{}.AttributeTypes())
	edit(&config)

	plan = config
	plan.Id = state.Id
	plan.MinNodeCount = types.Int32Unknown()
	plan.MaxNodeCount = types.Int32Unknown()
	plan.Nodepool = types.ObjectUnknown(ske.Nodepool{}.AttributeTypes())
	return plan, config
}

func TestNodepoolModifyPlan_LabelChange(t *testing.T) {
	state := nodepoolState(t)
	plan, config := nodepoolUpdate(state, func(p *ske.NodepoolResource) {
		p.Labels = []ske.Label{{Key: types.StringValue("a"), Value: types.StringValue("2")}}
	})

	got := modifyPlan(t, &skeNodepoolResource{}, nodepoolSchema(), state, plan, config)
	var result ske.NodepoolResource
	got.Get(context.Background(), &result)

	if !result.MinNodeCount.Equal(types.Int32Value(1)) || !result.MaxNodeCount.Equal(types.Int32Value(1)) {
		t.Errorf("min/max = %v/%v, want state 1/1", result.MinNodeCount, result.MaxNodeCount)
	}
	attrs := nestedAttrs(t, result.Nodepool)
	assertKnown(t, attrs, "id", "name", "auto_scale_enabled", "desired_node_count", "min_node_count", "max_node_count")
	assertUnknown(t, attrs, "labels", "status", "current_node_count", "modified_at")
}

func TestNodepoolModifyPlan_Scale(t *testing.T) {
	state := nodepoolState(t)
	plan, config := nodepoolUpdate(state, func(p *ske.NodepoolResource) {
		p.DesiredNodeCount = types.Int32Value(2)
	})

	got := modifyPlan(t, &skeNodepoolResource{}, nodepoolSchema(), state, plan, config)
	var result ske.NodepoolResource
	got.Get(context.Background(), &result)

	// min/max follow the new desired count on apply
	if !result.MinNodeCount.IsUnknown() || !result.MaxNodeCount.IsUnknown() {
		t.Errorf("min/max = %v/%v, want unknown", result.MinNodeCount, result.MaxNodeCount)
	}
	attrs := nestedAttrs(t, result.Nodepool)
	assertKnown(t, attrs, "id", "name", "labels")
	assertUnknown(t, attrs, "desired_node_count", "min_node_count", "max_node_count", "current_node_count", "status")
}

// setUnknown marks a root attribute unknown in both plan and config, as when it
// references an attribute of a resource that is not created yet.
func setUnknown(t *testing.T, req *resource.ModifyPlanRequest, s schema.Schema, name string) {
	t.Helper()
	ctx := context.Background()
	typ := s.Attributes[name].GetType()
	unknown, err := typ.ValueFromTerraform(ctx, tftypes.NewValue(typ.TerraformType(ctx), tftypes.UnknownValue))
	if err != nil {
		t.Fatal(err)
	}
	if d := req.Plan.SetAttribute(ctx, path.Root(name), unknown); d.HasError() {
		t.Fatalf("set plan %s: %v", name, d)
	}
	cfg := tfsdk.Plan{Schema: s, Raw: req.Config.Raw}
	if d := cfg.SetAttribute(ctx, path.Root(name), unknown); d.HasError() {
		t.Fatalf("set config %s: %v", name, d)
	}
	req.Config.Raw = cfg.Raw
}

func TestNodepoolModifyPlan_UnknownInput(t *testing.T) {
	s := nodepoolSchema()
	for _, name := range []string{"labels", "taints", "linked_resources", "advanced_settings"} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			state := nodepoolState(t)
			plan, config := nodepoolUpdate(state, func(p *ske.NodepoolResource) {
				p.IsAutoRecovery = types.BoolValue(false)
			})
			req := modifyPlanRequest(t, s, state, plan, config)
			setUnknown(t, &req, s, name)

			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&skeNodepoolResource{}).ModifyPlan(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("ModifyPlan: %v", resp.Diagnostics)
			}

			var got types.Object
			resp.Plan.GetAttribute(ctx, path.Root("nodepool"), &got)
			attrs := nestedAttrs(t, got)
			assertKnown(t, attrs, "id", "name")
			if name != "advanced_settings" {
				assertUnknown(t, attrs, name)
			}
		})
	}
}

func TestClusterModifyPlan_UnknownInput(t *testing.T) {
	s := ClusterResourceSchema()
	for _, name := range []string{"linked_resources", "private_endpoint_access_control_resources", "security_group_id_list", "additional_subnet_id_list", "tags"} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			state, plan, config := clusterDeletionProtectionUpdate(t)
			req := modifyPlanRequest(t, s, state, plan, config)
			setUnknown(t, &req, s, name)

			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&skeClusterResource{}).ModifyPlan(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("ModifyPlan: %v", resp.Diagnostics)
			}

			var got types.Object
			resp.Plan.GetAttribute(ctx, path.Root("cluster"), &got)
			assertKnown(t, nestedAttrs(t, got), "id", "name")
		})
	}
}

func TestClusterModifyPlan_DeletionProtection(t *testing.T) {
	ctx := context.Background()
	state, plan, config := clusterDeletionProtectionUpdate(t)
	subnets := state.AdditionalSubnetIdList

	got := modifyPlan(t, &skeClusterResource{}, ClusterResourceSchema(), state, plan, config)
	var result ske.ClusterResource
	got.Get(ctx, &result)

	if !result.AdditionalSubnetIdList.Equal(subnets) {
		t.Errorf("additional_subnet_id_list = %v, want state value", result.AdditionalSubnetIdList)
	}
	if result.LinkedResources.IsUnknown() || result.PrivateEndpointAccessControlResources.IsUnknown() {
		t.Errorf("linked_resources / private_endpoint_access_control_resources are unknown, want state value")
	}
	if !result.PublicEndpointAccessControlIp.IsNull() || !result.NfsVolumeId.IsNull() {
		t.Errorf("public_endpoint_access_control_ip / nfs_volume_id = %v / %v, want state null", result.PublicEndpointAccessControlIp, result.NfsVolumeId)
	}
	if !result.DeletionProtectionEnabled.Equal(types.BoolValue(true)) {
		t.Errorf("deletion_protection_enabled = %v, want true", result.DeletionProtectionEnabled)
	}
	attrs := nestedAttrs(t, result.Cluster)
	assertKnown(t, attrs, "id", "name", "vpc", "kubernetes_version", "additional_subnet_id_list", "linked_resources")
	assertUnknown(t, attrs, "status", "modified_at", "node_count")
}

// clusterDeletionProtectionUpdate returns state, plan and config for a deletion protection toggle.
func clusterDeletionProtectionUpdate(t *testing.T) (state, plan, config ske.ClusterResource) {
	t.Helper()
	ctx := context.Background()
	subnets, _ := types.ListValueFrom(ctx, types.StringType, []string{"sub-2"})
	nested := ske.Cluster{
		Id:                     types.StringValue("cl-1"),
		Name:                   types.StringValue("tf-scale-w"),
		KubernetesVersion:      types.StringValue("v1.36.3"),
		AdditionalSubnetIdList: subnets,
		Status:                 types.StringValue("RUNNING"),
		ModifiedAt:             types.StringValue("2026-10-15T00:00:00Z"),
		Vpc:                    ske.ExternalResourceId{Id: types.StringValue("vpc-1")},
	}
	obj, d := types.ObjectValueFrom(ctx, nested.AttributeTypes(), nested)
	if d.HasError() {
		t.Fatalf("nested object: %v", d)
	}
	emptyLinked, _ := types.ListValueFrom(ctx, linkedResourceObjectType, []ske.LinkedResource{})
	emptyPeacr, _ := types.ListValueFrom(ctx, linkedResourceObjectType, []ske.PrivateEndpointAccessControlResource{})
	emptySg, _ := types.ListValueFrom(ctx, types.StringType, []string{})
	state = ske.ClusterResource{
		Id:                                    types.StringValue("cl-1"),
		Name:                                  types.StringValue("tf-scale-w"),
		KubernetesVersion:                     types.StringValue("v1.36.3"),
		ServiceWatchLoggingEnabled:            types.BoolValue(false),
		DefaultSubnetId:                       types.StringValue("sub-1"),
		VpcId:                                 types.StringValue("vpc-1"),
		PrivateEndpointAccessControlResources: emptyPeacr,
		PublicEndpointAccessControlIp:         types.StringNull(),
		SecurityGroupIdList:                   emptySg,
		NfsVolumeId:                           types.StringNull(),
		AdditionalSubnetIdList:                subnets,
		DeletionProtectionEnabled:             types.BoolValue(false),
		LinkedResources:                       emptyLinked,
		Tags:                                  types.MapNull(types.StringType),
		Cluster:                               obj,
	}

	config = state
	config.Id = types.StringNull()
	config.DeletionProtectionEnabled = types.BoolValue(true)
	config.PrivateEndpointAccessControlResources = types.ListNull(linkedResourceObjectType)
	config.SecurityGroupIdList = types.ListNull(types.StringType)
	config.AdditionalSubnetIdList = types.ListNull(types.StringType)
	config.LinkedResources = types.ListNull(linkedResourceObjectType)
	config.Cluster = types.ObjectNull(ske.Cluster{}.AttributeTypes())

	plan = config
	plan.Id = state.Id
	plan.PrivateEndpointAccessControlResources = types.ListUnknown(linkedResourceObjectType)
	plan.PublicEndpointAccessControlIp = types.StringUnknown()
	plan.SecurityGroupIdList = types.ListUnknown(types.StringType)
	plan.NfsVolumeId = types.StringUnknown()
	plan.AdditionalSubnetIdList = types.ListUnknown(types.StringType)
	plan.LinkedResources = types.ListUnknown(linkedResourceObjectType)
	plan.Cluster = types.ObjectUnknown(ske.Cluster{}.AttributeTypes())

	return state, plan, config
}

func TestNodepoolModifyPlan_UnsupportedChange(t *testing.T) {
	tests := []struct {
		name    string
		edit    func(*ske.NodepoolResource)
		wantErr bool
	}{
		{"name change", func(p *ske.NodepoolResource) { p.Name = types.StringValue("renamed") }, true},
		{"os version without kubernetes version", func(p *ske.NodepoolResource) { p.ImageOsVersion = types.StringValue("24.04") }, true},
		{"os version with kubernetes version", func(p *ske.NodepoolResource) {
			p.ImageOsVersion = types.StringValue("24.04")
			p.KubernetesVersion = types.StringValue("v1.37.0")
		}, false},
		{"enable auto scale without min/max", func(p *ske.NodepoolResource) { p.IsAutoScale = types.BoolValue(true) }, true},
		{"enable auto scale with min/max", func(p *ske.NodepoolResource) {
			p.IsAutoScale = types.BoolValue(true)
			p.MinNodeCount = types.Int32Value(1)
			p.MaxNodeCount = types.Int32Value(3)
		}, false},
		// not filled on import: prior value unknown, so not rejected
		{"attribute missing in state", func(p *ske.NodepoolResource) { p.VolumeMaxIops = types.Int32Value(5000) }, false},
		{"label change", func(p *ske.NodepoolResource) { p.Labels = []ske.Label{label("a", "2")} }, false},
	}
	s := nodepoolSchema()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := nodepoolState(t)
			state.KubernetesVersion = types.StringValue("v1.36.3")
			state.ImageOsVersion = types.StringValue("22.04")
			plan, config := nodepoolUpdate(state, tt.edit)
			if !config.MinNodeCount.IsNull() {
				plan.MinNodeCount, plan.MaxNodeCount = config.MinNodeCount, config.MaxNodeCount
			}
			req := modifyPlanRequest(t, s, state, plan, config)
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&skeNodepoolResource{}).ModifyPlan(context.Background(), req, resp)
			if got := resp.Diagnostics.HasError(); got != tt.wantErr {
				t.Errorf("error = %v, want %v: %v", got, tt.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestNodepoolModifyPlan_PartiallyUnknownAdvancedSettings(t *testing.T) {
	known := ske.AdvancedSettings{
		MaxPods:              types.Int32Value(110),
		ImageGcHighThreshold: types.Int32Value(85),
		ImageGcLowThreshold:  types.Int32Value(80),
		ContainerLogMaxSize:  types.Int32Value(10),
		ContainerLogMaxFiles: types.Int32Value(5),
		PodMaxPids:           types.Int32Value(4096),
		AllowedUnsafeSysctls: types.StringNull(),
	}
	tests := []struct {
		name    string
		maxPods types.Int32
		wantErr bool
	}{
		{"unchanged", types.Int32Value(110), false},
		// set from a resource not created yet: change is not known until apply
		{"one attribute unknown", types.Int32Unknown(), false},
		{"changed", types.Int32Value(200), true},
	}
	s := nodepoolSchema()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			state := nodepoolState(t)
			state.AdvancedSettings = &known
			planned := known
			planned.MaxPods = tt.maxPods
			plan, config := nodepoolUpdate(state, func(p *ske.NodepoolResource) {
				p.Labels = []ske.Label{label("a", "2")}
				p.AdvancedSettings = &planned
			})
			req := modifyPlanRequest(t, s, state, plan, config)
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&skeNodepoolResource{}).ModifyPlan(ctx, req, resp)
			if got := resp.Diagnostics.HasError(); got != tt.wantErr {
				t.Errorf("error = %v, want %v: %v", got, tt.wantErr, resp.Diagnostics)
			}
		})
	}
}

// C-12 in the PP full test: linked_resources set on the server but left out of
// the configuration made every plan show an update.
func TestClusterModifyPlan_UnsetLinkedResourcesNoChange(t *testing.T) {
	ctx := context.Background()
	state, _, _ := clusterDeletionProtectionUpdate(t)
	linked, _ := types.ListValueFrom(ctx, linkedResourceObjectType, []ske.LinkedResource{{
		Id: types.StringValue("fs-1"), Name: types.StringValue("vol"), Type: types.StringValue("fs"),
	}})
	state.LinkedResources = linked

	config := state
	config.Id = types.StringNull()
	config.LinkedResources = types.ListNull(linkedResourceObjectType)
	config.Cluster = types.ObjectNull(ske.Cluster{}.AttributeTypes())

	// as proposed by Terraform: the unset nested list is unknown, so the computed object is too
	plan := state
	plan.LinkedResources = types.ListUnknown(linkedResourceObjectType)
	plan.Cluster = types.ObjectUnknown(ske.Cluster{}.AttributeTypes())

	got := modifyPlan(t, &skeClusterResource{}, ClusterResourceSchema(), state, plan, config)
	var result ske.ClusterResource
	got.Get(ctx, &result)

	if !result.LinkedResources.Equal(linked) {
		t.Errorf("linked_resources = %v, want state value", result.LinkedResources)
	}
	gotTf, _ := result.Cluster.ToTerraformValue(ctx)
	wantTf, _ := state.Cluster.ToTerraformValue(ctx)
	if !gotTf.Equal(wantTf) {
		t.Errorf("cluster = %v, want state value %v (no change)", gotTf, wantTf)
	}
}

func TestClusterModifyPlan_ChangeKeepsVolatileUnknown(t *testing.T) {
	// a real change must still leave volatile nested fields unknown
	state, plan, config := clusterDeletionProtectionUpdate(t)
	got := modifyPlan(t, &skeClusterResource{}, ClusterResourceSchema(), state, plan, config)
	var result ske.ClusterResource
	got.Get(context.Background(), &result)
	assertUnknown(t, nestedAttrs(t, result.Cluster), "status", "modified_at")
}
