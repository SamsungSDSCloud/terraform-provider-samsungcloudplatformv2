package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LbHealthCheckResource struct {
	Id                  types.String         `tfsdk:"id"`
	LbHealthCheck       types.Object         `tfsdk:"lb_health_check"`
	LbHealthCheckCreate *LbHealthCheckCreate `tfsdk:"lb_health_check_create"`
}

type LbHealthCheckCreate struct {
	Name                types.String `tfsdk:"name"`
	VpcId               types.String `tfsdk:"vpc_id"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	Protocol            types.String `tfsdk:"protocol"`
	HealthCheckPort     types.Int32  `tfsdk:"health_check_port"`
	HealthCheckInterval types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckTimeout  types.Int32  `tfsdk:"health_check_timeout"`
	HealthCheckCount    types.Int32  `tfsdk:"health_check_count"`
	HttpMethod          types.String `tfsdk:"http_method"`
	HealthCheckUrl      types.String `tfsdk:"health_check_url"`
	ResponseCode        types.String `tfsdk:"response_code"`
	RequestData         types.String `tfsdk:"request_data"`
	Description         types.String `tfsdk:"description"`
	Tags                types.Map    `tfsdk:"tags"`
}

type LbHealthCheckDetail struct {
	Name                types.String `tfsdk:"name"`
	VpcId               types.String `tfsdk:"vpc_id"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	Protocol            types.String `tfsdk:"protocol"`
	HealthCheckPort     types.Int32  `tfsdk:"health_check_port"`
	HealthCheckInterval types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckTimeout  types.Int32  `tfsdk:"health_check_timeout"`
	HealthCheckCount    types.Int32  `tfsdk:"health_check_count"`
	HttpMethod          types.String `tfsdk:"http_method"`
	HealthCheckUrl      types.String `tfsdk:"health_check_url"`
	ResponseCode        types.String `tfsdk:"response_code"`
	RequestData         types.String `tfsdk:"request_data"`
	HealthCheckType     types.String `tfsdk:"health_check_type"`
	Description         types.String `tfsdk:"description"`
	State               types.String `tfsdk:"state"`
	AccountId           types.String `tfsdk:"account_id"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
}

func (m LbHealthCheckDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"vpc_id":                types.StringType,
		"subnet_id":             types.StringType,
		"protocol":              types.StringType,
		"health_check_port":     types.Int32Type,
		"health_check_interval": types.Int32Type,
		"health_check_timeout":  types.Int32Type,
		"health_check_count":    types.Int32Type,
		"http_method":           types.StringType,
		"health_check_url":      types.StringType,
		"response_code":         types.StringType,
		"request_data":          types.StringType,
		"health_check_type":     types.StringType,
		"description":           types.StringType,
		"state":                 types.StringType,
		"account_id":            types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}
