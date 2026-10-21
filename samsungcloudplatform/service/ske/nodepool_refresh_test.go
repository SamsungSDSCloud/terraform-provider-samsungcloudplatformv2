package ske

import (
	"reflect"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func label(k, v string) ske.Label {
	return ske.Label{Key: types.StringValue(k), Value: types.StringValue(v)}
}

func TestRefreshNodepoolRootFields(t *testing.T) {
	server := ske.Nodepool{
		Name:                types.StringValue("scale"),
		AutoScaleEnabled:    types.BoolValue(false),
		AutoRecoveryEnabled: types.BoolValue(true),
		DesiredNodeCount:    types.Int32Value(3),
		MinNodeCount:        types.Int32Value(3),
		MaxNodeCount:        types.Int32Value(3),
		Labels:              []ske.Label{label("a", "1"), label("b", "2")},
		PreferredIps:        types.StringValue(""),
		ServerGroupId:       types.StringNull(),
	}

	t.Run("external changes are reflected", func(t *testing.T) {
		state := ske.NodepoolResource{
			IsAutoScale:      types.BoolValue(false),
			IsAutoRecovery:   types.BoolValue(false),
			DesiredNodeCount: types.Int32Value(1),
			MinNodeCount:     types.Int32Value(1),
			Labels:           []ske.Label{label("a", "1")},
		}
		refreshNodepoolRootFields(&state, server)

		if !state.IsAutoRecovery.Equal(types.BoolValue(true)) {
			t.Errorf("is_auto_recovery = %v, want true", state.IsAutoRecovery)
		}
		if !state.DesiredNodeCount.Equal(types.Int32Value(3)) {
			t.Errorf("desired_node_count = %v, want 3", state.DesiredNodeCount)
		}
		if !reflect.DeepEqual(state.Labels, server.Labels) {
			t.Errorf("labels = %v, want %v", state.Labels, server.Labels)
		}
	})

	t.Run("values matching the configuration are kept", func(t *testing.T) {
		state := ske.NodepoolResource{
			IsAutoScale:      types.BoolValue(false),
			DesiredNodeCount: types.Int32Value(3),
			MinNodeCount:     types.Int32Value(1), // configured, not used for fixed size
			Labels:           []ske.Label{label("b", "2"), label("a", "1")},
			PreferredIps:     types.StringNull(),
			ServerGroupId:    types.StringNull(),
		}
		refreshNodepoolRootFields(&state, server)

		if !state.MinNodeCount.Equal(types.Int32Value(1)) {
			t.Errorf("min_node_count = %v, want configured 1", state.MinNodeCount)
		}
		if !reflect.DeepEqual(state.Labels, []ske.Label{label("b", "2"), label("a", "1")}) {
			t.Errorf("labels = %v, want config order kept", state.Labels)
		}
		if !state.PreferredIps.IsNull() {
			t.Errorf("preferred_ips = %v, want null kept", state.PreferredIps)
		}
	})
}
