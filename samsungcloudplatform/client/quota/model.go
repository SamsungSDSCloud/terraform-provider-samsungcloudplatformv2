package quota

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-quota"

type AccountQuotaDataSource struct {
	AccountQuota types.Object `tfsdk:"account_quota"`
	Id           types.String `tfsdk:"id"`
}

type AccountQuotaType struct {
	basetypes.ObjectType
}

type AccountQuotaDataSourceIds struct {
	Ids    []types.String  `tfsdk:"ids"`
	Filter []filter.Filter `tfsdk:"filter"` // filter field 를 추가한다.
}
type AccountQuota struct {
	AccountId     types.String `tfsdk:"account_id"`
	AccountName   types.String `tfsdk:"account_name"`
	Adjustable    types.Bool   `tfsdk:"adjustable"`
	AppliedValue  types.Int32  `tfsdk:"applied_value"`
	Approval      types.Bool   `tfsdk:"approval"`
	ClassValue    types.String `tfsdk:"class_value"`
	CreatedAt     types.String `tfsdk:"created_at"`
	Description   types.String `tfsdk:"description"`
	FreeRate      types.Int32  `tfsdk:"free_rate"`
	Id            types.String `tfsdk:"id"`
	InitialValue  types.Int32  `tfsdk:"initial_value"`
	MaxPerAccount types.Int64  `tfsdk:"max_per_account"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
	QuotaItem     types.String `tfsdk:"quota_item"`
	Reduction     types.Bool   `tfsdk:"reduction"`
	Request       types.Bool   `tfsdk:"request"`
	RequestClass  types.String `tfsdk:"request_class"`
	ResourceType  types.String `tfsdk:"resource_type"`
	Service       types.String `tfsdk:"service"`
	Srn           types.String `tfsdk:"srn"`
	Unit          types.String `tfsdk:"unit"`
	state         attr.ValueState
}

func (v AccountQuota) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":      types.StringType,
		"account_name":    types.StringType,
		"adjustable":      types.BoolType,
		"applied_value":   types.Int64Type,
		"approval":        types.BoolType,
		"class_value":     types.StringType,
		"created_at":      types.StringType,
		"description":     types.StringType,
		"free_rate":       types.Int64Type,
		"id":              types.StringType,
		"initial_value":   types.Int64Type,
		"max_per_account": types.Int64Type,
		"modified_at":     types.StringType,
		"quota_item":      types.StringType,
		"reduction":       types.BoolType,
		"request":         types.BoolType,
		"request_class":   types.StringType,
		"resource_type":   types.StringType,
		"service":         types.StringType,
		"srn":             types.StringType,
		"unit":            types.StringType,
	}
}
