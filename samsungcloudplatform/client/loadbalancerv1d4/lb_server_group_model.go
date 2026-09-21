package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LbServerGroupDataSourceDetailV1d4 represents the top-level response for the v1.4 LB Server Group data source.
type LbServerGroupDataSourceDetail struct {
	Id                  types.String `tfsdk:"id"`
	LbServerGroupDetail types.Object `tfsdk:"lb_server_group"`
}

// LbServerGroupDetailV1d4 represents a single LB Server Group from the v1.4 API response.
type LbServerGroupDetail struct {
	Name            types.String `tfsdk:"name"`
	Protocol        types.String `tfsdk:"protocol"`
	LoadbalancerId  types.String `tfsdk:"loadbalancer_id"`
	LbName          types.String `tfsdk:"lb_name"`
	LbMethod        types.String `tfsdk:"lb_method"`
	LbHealthCheckId types.String `tfsdk:"lb_health_check_id"`
	State           types.String `tfsdk:"state"`
	VpcId           types.String `tfsdk:"vpc_id"`
	SubnetId        types.String `tfsdk:"subnet_id"`
	AccountId       types.String `tfsdk:"account_id"`
	Description     types.String `tfsdk:"description"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatedAt       types.String `tfsdk:"created_at"`
}

// AttributeTypes returns the attribute types for LbServerGroupDetail.
func (m LbServerGroupDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               types.StringType,
		"protocol":           types.StringType,
		"loadbalancer_id":    types.StringType,
		"lb_name":            types.StringType,
		"lb_method":          types.StringType,
		"lb_health_check_id": types.StringType,
		"state":              types.StringType,
		"vpc_id":             types.StringType,
		"subnet_id":          types.StringType,
		"account_id":         types.StringType,
		"description":        types.StringType,
		"modified_by":        types.StringType,
		"modified_at":        types.StringType,
		"created_by":         types.StringType,
		"created_at":         types.StringType,
	}
}

// LbServerGroupResourceV1d4 represents the top-level resource for the v1.4 LB Server Group.
type LbServerGroupResourceV1d4 struct {
	Id                  types.String `tfsdk:"id"`
	LbServerGroup       types.Object `tfsdk:"lb_server_group"`
	LbServerGroupCreate types.Object `tfsdk:"lb_server_group_create"`
}

// LbServerGroupResourceDetail represents a single LB Server Group from the v1.4 API create response.
type LbServerGroupResourceDetail struct {
	Name            types.String `tfsdk:"name"`
	Protocol        types.String `tfsdk:"protocol"`
	LoadbalancerId  types.String `tfsdk:"loadbalancer_id"`
	LbName          types.String `tfsdk:"lb_name"`
	LbMethod        types.String `tfsdk:"lb_method"`
	LbHealthCheckId types.String `tfsdk:"lb_health_check_id"`
	State           types.String `tfsdk:"state"`
	VpcId           types.String `tfsdk:"vpc_id"`
	SubnetId        types.String `tfsdk:"subnet_id"`
	AccountId       types.String `tfsdk:"account_id"`
	Description     types.String `tfsdk:"description"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatedAt       types.String `tfsdk:"created_at"`
}

// AttributeTypes returns the attribute types for LbServerGroupResourceDetail.
func (m LbServerGroupResourceDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               types.StringType,
		"protocol":           types.StringType,
		"loadbalancer_id":    types.StringType,
		"lb_name":            types.StringType,
		"lb_method":          types.StringType,
		"lb_health_check_id": types.StringType,
		"state":              types.StringType,
		"vpc_id":             types.StringType,
		"subnet_id":          types.StringType,
		"account_id":         types.StringType,
		"description":        types.StringType,
		"modified_by":        types.StringType,
		"modified_at":        types.StringType,
		"created_by":         types.StringType,
		"created_at":         types.StringType,
	}
}
