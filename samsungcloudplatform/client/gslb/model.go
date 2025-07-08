package gslb

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-gslb" // 해당 서비스의 서비스 타입(keystone 에 등록된 service type)을 추가한다.

type GslbDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Size  types.Int32  `tfsdk:"size"`
	Page  types.Int32  `tfsdk:"page"`
	Sort  types.String `tfsdk:"sort"`
	State types.String `tfsdk:"state"`
	Name  types.String `tfsdk:"name"`
	Gslbs []Gslb       `tfsdk:"gslbs"`
}

type Gslb struct {
	Algorithm           types.String `tfsdk:"algorithm"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	Description         types.String `tfsdk:"description"`
	EnvUsage            types.String `tfsdk:"env_usage"`
	Id                  types.String `tfsdk:"id"`
	LinkedResourceCount types.Int32  `tfsdk:"linked_resource_count"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	Name                types.String `tfsdk:"name"`
	State               types.String `tfsdk:"state"`
}

type GslbDataSourceDetail struct {
	Id         types.String `tfsdk:"id"`
	GslbDetail types.Object `tfsdk:"gslb_detail"`
}

type GslbDetail struct {
	Algorithm           types.String `tfsdk:"algorithm"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	Description         types.String `tfsdk:"description"`
	EnvUsage            types.String `tfsdk:"env_usage"`
	HealthCheck         *HealthCheck `tfsdk:"health_check"`
	Id                  types.String `tfsdk:"id"`
	LinkedResourceCount types.Int32  `tfsdk:"linked_resource_count"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	Name                types.String `tfsdk:"name"`
	State               types.String `tfsdk:"state"`
}

type HealthCheck struct {
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	HealthCheckInterval     types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckProbeTimeout types.Int32  `tfsdk:"health_check_probe_timeout"`
	HealthCheckUserId       types.String `tfsdk:"health_check_user_id"`
	HealthCheckUserPassword types.String `tfsdk:"health_check_user_password"`
	Id                      types.String `tfsdk:"id"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	Protocol                types.String `tfsdk:"protocol"`
	ReceiveString           types.String `tfsdk:"receive_string"`
	SendString              types.String `tfsdk:"send_string"`
	ServicePort             types.Int32  `tfsdk:"service_port"`
	Timeout                 types.Int32  `tfsdk:"timeout"`
}

func (m GslbDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"algorithm":   types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"env_usage":   types.StringType,
		"health_check": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"created_at":                 types.StringType,
				"created_by":                 types.StringType,
				"health_check_interval":      types.Int32Type,
				"health_check_probe_timeout": types.Int32Type,
				"health_check_user_id":       types.StringType,
				"health_check_user_password": types.StringType,
				"id":                         types.StringType,
				"modified_at":                types.StringType,
				"modified_by":                types.StringType,
				"protocol":                   types.StringType,
				"receive_string":             types.StringType,
				"send_string":                types.StringType,
				"service_port":               types.Int32Type,
				"timeout":                    types.Int32Type,
			},
		},
		"id":                    types.StringType,
		"linked_resource_count": types.Int32Type,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
		"name":                  types.StringType,
		"state":                 types.StringType,
	}
}

type GslbResourceDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Size          types.Int32          `tfsdk:"size"`
	Page          types.Int32          `tfsdk:"page"`
	Sort          types.String         `tfsdk:"sort"`
	GslbId        types.String         `tfsdk:"gslb_id"`
	GslbResources []GslbResourceDetail `tfsdk:"gslb_resources"`
}

type GslbResourceDetail struct {
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Destination types.String `tfsdk:"destination"`
	Disabled    types.Bool   `tfsdk:"disabled"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Region      types.String `tfsdk:"region"`
	Weight      types.Int32  `tfsdk:"weight"`
}

type GslbResource struct {
	Id         types.String `tfsdk:"id"`
	Gslb       types.Object `tfsdk:"gslb"`
	GslbCreate GslbCreate   `tfsdk:"gslb_create"`
	Tags       types.Map    `tfsdk:"tags"`
}

type GslbCreate struct {
	Algorithm   types.String         `tfsdk:"algorithm"`
	Description types.String         `tfsdk:"description"`
	EnvUsage    types.String         `tfsdk:"env_usage"`
	HealthCheck *HealthCheckCreate   `tfsdk:"health_check"`
	Name        types.String         `tfsdk:"name"`
	Resources   []GslbResourceCreate `tfsdk:"resources"`
}

type HealthCheckCreate struct {
	HealthCheckInterval     types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckProbeTimeout types.Int32  `tfsdk:"health_check_probe_timeout"`
	HealthCheckUserId       types.String `tfsdk:"health_check_user_id"`
	HealthCheckUserPassword types.String `tfsdk:"health_check_user_password"`
	Protocol                types.String `tfsdk:"protocol"`
	ReceiveString           types.String `tfsdk:"receive_string"`
	SendString              types.String `tfsdk:"send_string"`
	ServicePort             types.Int32  `tfsdk:"service_port"`
	Timeout                 types.Int32  `tfsdk:"timeout"`
}

type GslbResourceCreate struct {
	Description types.String `tfsdk:"description"`
	Destination types.String `tfsdk:"destination"`
	Disabled    types.Bool   `tfsdk:"disabled"`
	Region      types.String `tfsdk:"region"`
	Weight      types.Int32  `tfsdk:"weight"`
}
