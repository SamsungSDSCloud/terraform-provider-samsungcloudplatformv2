package parallelfilestorage

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-parallel-filestorage" // 해당 서비스의 서비스 타입(keystone 에 등록된 service type)을 추가한다.

type VolumeDataSourceIds struct {
	Offset   types.Int32    `tfsdk:"offset"`
	Limit    types.Int32    `tfsdk:"limit"`
	Sort     types.String   `tfsdk:"sort"`
	Name     types.String   `tfsdk:"name"`
	Ids      []types.String `tfsdk:"ids"`
}

type VolumeDataSource struct {
	AccountId               types.String		 `tfsdk:"account_id"`
	CreatedAt               types.String		 `tfsdk:"created_at"`
	Id                      types.String		 `tfsdk:"id"`
	Name                    types.String		 `tfsdk:"name"`
	MountPath               types.String		 `tfsdk:"mount_path"`
	State                   types.String		 `tfsdk:"state"`
	Zone                	types.String		 `tfsdk:"zone"`
	CapacityTb              types.Int32  		 `tfsdk:"capacity_tb"`
	AccessRules             []AccessRuleResource `tfsdk:"access_rules"`
}

type VolumeResource struct {
	AccountId               types.String            `tfsdk:"account_id"`
	CapacityTb              types.Int32             `tfsdk:"capacity_tb"`
	CreatedAt               types.String            `tfsdk:"created_at"`
	Id                      types.String            `tfsdk:"id"`
	Name                    types.String            `tfsdk:"name"`
	State                   types.String            `tfsdk:"state"`
	Zone					types.String			`tfsdk:"zone"`
	Tags                    types.Map               `tfsdk:"tags"`
	AccessRules             []AccessRuleResource	`tfsdk:"access_rules"`
}

type AccessRuleResource struct {
	ObjectId   types.String `tfsdk:"object_id"`
	ObjectType types.String `tfsdk:"object_type"`
}
