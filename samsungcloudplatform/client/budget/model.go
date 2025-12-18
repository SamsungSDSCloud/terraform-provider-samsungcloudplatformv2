package budget

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-budget"

type BudgetDataSource struct {
	Id            types.String    `tfsdk:"id"`
	Name          types.String    `tfsdk:"name"`
	Filter        []filter.Filter `tfsdk:"filter"` // filter field 를 추가한다.
	Budget        types.Object    `tfsdk:"budget"`
	Notifications types.Object    `tfsdk:"notifications"`
	Prevention    types.Object    `tfsdk:"prevention"`
}

type BudgetResource struct {
	Id            types.String  `tfsdk:"id"`
	LastUpdated   types.String  `tfsdk:"last_updated"`
	Budget        types.Object  `tfsdk:"budget"`
	Name          types.String  `tfsdk:"name"`
	Amount        types.Int32   `tfsdk:"amount"`
	StartMonth    types.String  `tfsdk:"start_month"`
	Unit          types.String  `tfsdk:"unit"`
	Notifications Notifications `tfsdk:"notifications"`
	Prevention    Prevention    `tfsdk:"prevention"`
}

type BudgetType struct {
	basetypes.ObjectType
}

type BudgetDataSourceIds struct {
	Id     types.String    `tfsdk:"id"`
	Name   types.String    `tfsdk:"name"`
	Filter []filter.Filter `tfsdk:"filter"` // filter field 를 추가한다.
	Ids    []types.String  `tfsdk:"ids"`
}

type Budget struct {
	Amount     types.Int32  `tfsdk:"amount"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	BudgetId   types.String `tfsdk:"budget_id"`
	ModifiedAt types.String `tfsdk:"modified_at"`
	ModifiedBy types.String `tfsdk:"modified_by"`
	Name       types.String `tfsdk:"name"`
	StartMonth types.String `tfsdk:"start_month"`
	BudgetType types.String `tfsdk:"type"`
	Unit       types.String `tfsdk:"unit"`
	state      attr.ValueState
}

func (v Budget) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"amount":      types.Int32Type,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"budget_id":   types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"start_month": types.StringType,
		"type":        types.StringType,
		"unit":        types.StringType,
	}
}

type NotificationsType struct {
	basetypes.ObjectType
}

type Notifications struct {
	IsUseNotification      types.Bool   `tfsdk:"is_use_notification"`
	NotificationSendPeriod types.String `tfsdk:"notification_send_period"`
	Receivers              types.List   `tfsdk:"receivers"`
	Thresholds             types.List   `tfsdk:"thresholds"`
	state                  attr.ValueState
}

func (v Notifications) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"is_use_notification":      types.BoolType,
		"notification_send_period": types.StringType,
		"receivers": types.ListType{
			ElemType: types.StringType,
		},
		"thresholds": types.ListType{
			ElemType: types.Int32Type,
		},
	}
}

type PreventionType struct {
	basetypes.ObjectType
}

type Prevention struct {
	IsUsePrevention types.Bool  `tfsdk:"is_use_prevention"`
	Receivers       types.List  `tfsdk:"receivers"`
	Threshold       types.Int32 `tfsdk:"threshold"`
	state           attr.ValueState
}

func (v Prevention) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"is_use_prevention": types.BoolType,
		"receivers": types.ListType{
			ElemType: types.StringType,
		},
		"threshold": types.Int32Type,
	}
}
