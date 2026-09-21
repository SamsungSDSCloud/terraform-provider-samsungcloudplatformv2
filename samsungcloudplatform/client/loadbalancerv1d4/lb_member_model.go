package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LbMemberDataSourceV1d4 represents the input parameters for listing LB members via v1.4 API.
type LbMemberDataSourceV1d4 struct {
	TotalCount      types.Int32    `tfsdk:"total_count"`
	Size            types.Int32    `tfsdk:"size"`
	Page            types.Int32    `tfsdk:"page"`
	Sort            types.String   `tfsdk:"sort"`
	Name            types.String   `tfsdk:"name"`
	MemberIp        types.String   `tfsdk:"member_ip"`
	MemberPort      types.Int32    `tfsdk:"member_port"`
	LbServerGroupId types.String   `tfsdk:"lb_server_group_id"`
	LbMembers       []LbMemberV1d4 `tfsdk:"lb_members"`
}

// LbMemberV1d4 represents a single LB member from the v1.4 API response,
// including the new ObjectAz (Availability Zone) field.
type LbMemberV1d4 struct {
	Id              types.String `tfsdk:"id"`
	LbServerGroupId types.String `tfsdk:"lb_server_group_id"`
	Name            types.String `tfsdk:"name"`
	MemberIp        types.String `tfsdk:"member_ip"`
	MemberPort      types.Int32  `tfsdk:"member_port"`
	MemberState     types.String `tfsdk:"member_state"`
	MemberWeight    types.Int32  `tfsdk:"member_weight"`
	ObjectType      types.String `tfsdk:"object_type"`
	ObjectId        types.String `tfsdk:"object_id"`
	State           types.String `tfsdk:"state"`
	HealthState     types.String `tfsdk:"health_state"`
	ObjectAz        types.String `tfsdk:"object_az"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

// AttributeTypes returns the attribute types for LbMemberV1d4.
func (m LbMemberV1d4) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"lb_server_group_id": types.StringType,
		"name":               types.StringType,
		"member_ip":          types.StringType,
		"member_port":        types.Int32Type,
		"member_state":       types.StringType,
		"member_weight":      types.Int32Type,
		"object_type":        types.StringType,
		"object_id":          types.StringType,
		"state":              types.StringType,
		"health_state":       types.StringType,
		"object_az":          types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}
