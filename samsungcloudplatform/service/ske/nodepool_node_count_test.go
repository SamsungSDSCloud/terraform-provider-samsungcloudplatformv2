package ske

import (
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFillUnconfiguredNodeCounts(t *testing.T) {
	tests := []struct {
		name                          string
		plan                          ske.NodepoolResource
		wantDesired, wantMin, wantMax types.Int32
	}{
		{
			name: "fixed size, min/max not configured: filled from server desired",
			plan: ske.NodepoolResource{
				IsAutoScale:      types.BoolValue(false),
				DesiredNodeCount: types.Int32Value(2),
				MinNodeCount:     types.Int32Unknown(),
				MaxNodeCount:     types.Int32Unknown(),
			},
			wantDesired: types.Int32Value(2), wantMin: types.Int32Value(2), wantMax: types.Int32Value(2),
		},
		{
			name: "fixed size, min/max configured: kept as planned",
			plan: ske.NodepoolResource{
				IsAutoScale:      types.BoolValue(false),
				DesiredNodeCount: types.Int32Value(2),
				MinNodeCount:     types.Int32Value(1),
				MaxNodeCount:     types.Int32Value(3),
			},
			wantDesired: types.Int32Value(2), wantMin: types.Int32Value(1), wantMax: types.Int32Value(3),
		},
		{
			name: "auto scale, desired not configured: filled from server min",
			plan: ske.NodepoolResource{
				IsAutoScale:      types.BoolValue(true),
				DesiredNodeCount: types.Int32Unknown(),
				MinNodeCount:     types.Int32Value(1),
				MaxNodeCount:     types.Int32Value(3),
			},
			wantDesired: types.Int32Value(1), wantMin: types.Int32Value(1), wantMax: types.Int32Value(3),
		},
		{
			name: "auto scale, desired configured: kept as planned",
			plan: ske.NodepoolResource{
				IsAutoScale:      types.BoolValue(true),
				DesiredNodeCount: types.Int32Value(2),
				MinNodeCount:     types.Int32Value(1),
				MaxNodeCount:     types.Int32Value(3),
			},
			wantDesired: types.Int32Value(2), wantMin: types.Int32Value(1), wantMax: types.Int32Value(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := tt.plan
			// server values: min=1, desired=2
			fillUnconfiguredNodeCounts(&plan, 1, 2)
			if !plan.DesiredNodeCount.Equal(tt.wantDesired) {
				t.Errorf("desired_node_count = %v, want %v", plan.DesiredNodeCount, tt.wantDesired)
			}
			if !plan.MinNodeCount.Equal(tt.wantMin) {
				t.Errorf("min_node_count = %v, want %v", plan.MinNodeCount, tt.wantMin)
			}
			if !plan.MaxNodeCount.Equal(tt.wantMax) {
				t.Errorf("max_node_count = %v, want %v", plan.MaxNodeCount, tt.wantMax)
			}
		})
	}
}
