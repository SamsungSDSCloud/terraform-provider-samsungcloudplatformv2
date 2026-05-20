package organization

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-organization"

type OrganizationDataSource struct {
	Name            types.String `tfsdk:"name"`
	Size            types.Int64  `tfsdk:"size"`
	Page            types.Int64  `tfsdk:"page"`
	Sort            types.String `tfsdk:"sort"`
	MasterAccountId types.String `tfsdk:"master_account_id"`
	OrganizationId  types.String `tfsdk:"organization_id"`
}

// OrganizationListDataSource is for organization list data source (contains both config and result fields)
type OrganizationListDataSource struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Organizations types.List   `tfsdk:"organizations"`
	TotalCount    types.Int64  `tfsdk:"total_count"`
	Page          types.Int64  `tfsdk:"page"`
	Size          types.Int64  `tfsdk:"size"`
	Sort          types.List   `tfsdk:"sort"`
}

// OrganizationSummaryItem represents a single organization in the list
type OrganizationSummaryItem struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	MasterAccountId     types.String `tfsdk:"master_account_id"`
	DelegationAccountId types.String `tfsdk:"delegation_account_id"`
	RootUnitId          types.String `tfsdk:"root_unit_id"`
	UseScpYn            types.Bool   `tfsdk:"use_scp_yn"`
}

func (o OrganizationSummaryItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    basetypes.StringType{},
		"name":                  basetypes.StringType{},
		"created_at":            basetypes.StringType{},
		"created_by":            basetypes.StringType{},
		"modified_at":           basetypes.StringType{},
		"modified_by":           basetypes.StringType{},
		"master_account_id":     basetypes.StringType{},
		"delegation_account_id": basetypes.StringType{},
		"root_unit_id":          basetypes.StringType{},
		"use_scp_yn":            basetypes.BoolType{},
	}
}

type OrganizationResource struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	DelegationAccountId types.String `tfsdk:"delegation_account_id"`
	UseScpYn            types.Bool   `tfsdk:"use_scp_yn"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	CreatorName         types.String `tfsdk:"creator_name"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	ModifierName        types.String `tfsdk:"modifier_name"`
	MasterAccountId     types.String `tfsdk:"master_account_id"`
	MasterAccountEmail  types.String `tfsdk:"master_account_email"`
	RootUnitId          types.String `tfsdk:"root_unit_id"`
	Srn                 types.String `tfsdk:"srn"`
}

func (o OrganizationResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    basetypes.StringType{},
		"name":                  basetypes.StringType{},
		"created_at":            basetypes.StringType{},
		"created_by":            basetypes.StringType{},
		"creator_name":          basetypes.StringType{},
		"modified_at":           basetypes.StringType{},
		"modified_by":           basetypes.StringType{},
		"modifier_name":         basetypes.StringType{},
		"master_account_id":     basetypes.StringType{},
		"master_account_email":  basetypes.StringType{},
		"delegation_account_id": basetypes.StringType{},
		"root_unit_id":          basetypes.StringType{},
		"srn":                   basetypes.StringType{},
		"use_scp_yn":            basetypes.BoolType{},
	}
}

type ServiceControlPolicyListDataSourceRequest struct {
	OrganizationId  types.String
	Name            types.String
	Type            types.String
	ExcludeTargetId types.String
	Page            types.Int64
	Size            types.Int64
	Id              types.String
	Sort            types.String
}

type ControlPoliciesValue struct {
	PolicyId   types.String `tfsdk:"policy_id"`
	PolicyName types.String `tfsdk:"policy_name"`
}

func (o ControlPoliciesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"policy_id":   basetypes.StringType{},
		"policy_name": basetypes.StringType{},
	}
}

type OrganizationUnitValue struct {
	ControlPolicies types.List   `tfsdk:"control_policies"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatorName     types.String `tfsdk:"creator_name"`
	Depth           types.Int64  `tfsdk:"depth"`
	Description     types.String `tfsdk:"description"`
	Id              types.String `tfsdk:"id"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifierName    types.String `tfsdk:"modifier_name"`
	Name            types.String `tfsdk:"name"`
	ParentUnitId    types.String `tfsdk:"parent_unit_id"`
	ServiceName     types.String `tfsdk:"service_name"`
	Srn             types.String `tfsdk:"srn"`
	Type            types.String `tfsdk:"type"`
}

func (o OrganizationUnitValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"control_policies": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: ControlPoliciesValue{}.AttributeTypes(ctx),
			},
		},
		"created_at":     basetypes.StringType{},
		"created_by":     basetypes.StringType{},
		"creator_name":   basetypes.StringType{},
		"depth":          basetypes.Int64Type{},
		"description":    basetypes.StringType{},
		"id":             basetypes.StringType{},
		"modified_at":    basetypes.StringType{},
		"modified_by":    basetypes.StringType{},
		"modifier_name":  basetypes.StringType{},
		"name":           basetypes.StringType{},
		"parent_unit_id": basetypes.StringType{},
		"service_name":   basetypes.StringType{},
		"srn":            basetypes.StringType{},
		"type":           basetypes.StringType{},
	}
}

type OrganizationUnitResource struct {
	Name             types.String   `tfsdk:"name"`
	Description      types.String   `tfsdk:"description"`
	OrganizationId   types.String   `tfsdk:"organization_id"`
	ParentUnitId     types.String   `tfsdk:"parent_unit_id"`
	PolicyIds        []types.String `tfsdk:"policy_ids"`
	OrganizationUnit types.Object   `tfsdk:"organization_unit"`
}

type OrganizationUnitDataSource struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
}

type OrganizationInvitationsDataSource struct {
	OrganizationId types.String
	Size           types.Int64
	Page           types.Int64
	Sort           types.String
	AccountId      types.String
	AccountName    types.String
	AccountEmail   types.String
	State          types.String
	LoginId        types.String
}

// InvitationCreateFailCausedValue for invitation create response
type InvitationCreateFailCausedValue struct {
	ErrorCode    types.String `tfsdk:"error_code"`
	FailedCaused types.String `tfsdk:"failed_caused"`
	FailedIds    types.List   `tfsdk:"failed_ids"`
	Response     types.String `tfsdk:"response"`
}

func (o InvitationCreateFailCausedValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"error_code":    types.StringType,
		"failed_caused": types.StringType,
		"failed_ids": types.ListType{
			ElemType: types.StringType,
		},
		"response": types.StringType,
	}
}

type InvitationResource struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	TargetLoginIds types.List   `tfsdk:"target_login_ids"`
	SuccessIds     types.List   `tfsdk:"success_ids"`
	FailedIds      types.List   `tfsdk:"failed_ids"`
}

type InvitationCancelRequest struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Ids            types.List   `tfsdk:"ids"`
}

func (o InvitationResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": types.StringType,
		"target_login_ids": types.ListType{
			ElemType: types.StringType,
		},
		"success_ids": types.ListType{
			ElemType: types.StringType,
		},
		"failed_ids": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: InvitationCreateFailCausedValue{}.AttributeTypes(ctx),
			},
		},
	}
}

func (o InvitationCancelRequest) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"ids": types.ListType{
			ElemType: types.StringType,
		},
	}
}

type InvitationValue struct {
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ExpiredTime     types.String `tfsdk:"expired_time"`
	Id              types.String `tfsdk:"id"`
	MasterAccountId types.String `tfsdk:"master_account_id"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	RequestedTime   types.String `tfsdk:"requested_time"`
	State           types.String `tfsdk:"state"`
	TargetAccountId types.String `tfsdk:"target_account_id"`
}

func (o InvitationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":        basetypes.StringType{},
		"created_by":        basetypes.StringType{},
		"expired_time":      basetypes.StringType{},
		"id":                basetypes.StringType{},
		"master_account_id": basetypes.StringType{},
		"modified_at":       basetypes.StringType{},
		"modified_by":       basetypes.StringType{},
		"organization_id":   basetypes.StringType{},
		"requested_time":    basetypes.StringType{},
		"state":             basetypes.StringType{},
		"target_account_id": basetypes.StringType{},
	}
}

// InvitationAcceptSuccessId - nested object for success_id
type InvitationAcceptSuccessId struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
}

func (o InvitationAcceptSuccessId) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
}

// InvitationAcceptFailedId - nested object for failed_id
type InvitationAcceptFailedId struct {
	FailedId     types.String `tfsdk:"failed_id"`
	FailedName   types.String `tfsdk:"failed_name"`
	FailedCaused types.String `tfsdk:"failed_caused"`
	ErrorCode    types.String `tfsdk:"error_code"`
	Response     types.Map    `tfsdk:"response"`
}

func (o InvitationAcceptFailedId) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"error_code":    basetypes.StringType{},
		"response":      basetypes.StringType{},
	}
}

type InvitationAcceptResponseValue struct {
	MasterAccountEmail types.String               `tfsdk:"master_account_email"`
	SuccessId          *InvitationAcceptSuccessId `tfsdk:"success_id"`
	FailedId           *InvitationAcceptFailedId  `tfsdk:"failed_id"`
}

func (o InvitationAcceptResponseValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"master_account_email": basetypes.StringType{},
		"success_id":           types.ObjectType{},
		"failed_id":            types.ObjectType{},
	}
}

type AccountInvitationValue struct {
	Id                   types.String `tfsdk:"id"`
	OrganizationId       types.String `tfsdk:"organization_id"`
	OrganizationName     types.String `tfsdk:"organization_name"`
	MasterAccountId      types.String `tfsdk:"master_account_id"`
	MasterAccountName    types.String `tfsdk:"master_account_name"`
	MasterAccountEmail   types.String `tfsdk:"master_account_email"`
	MasterAccountLoginId types.String `tfsdk:"master_account_login_id"`
	TargetAccountId      types.String `tfsdk:"target_account_id"`
	RequestedTime        types.String `tfsdk:"requested_time"`
	ExpiredTime          types.String `tfsdk:"expired_time"`
	State                types.String `tfsdk:"state"`
	CreatedAt            types.String `tfsdk:"created_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
}

func (o AccountInvitationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      basetypes.StringType{},
		"organization_id":         basetypes.StringType{},
		"organization_name":       basetypes.StringType{},
		"master_account_id":       basetypes.StringType{},
		"master_account_name":     basetypes.StringType{},
		"master_account_email":    basetypes.StringType{},
		"master_account_login_id": basetypes.StringType{},
		"target_account_id":       basetypes.StringType{},
		"requested_time":          basetypes.StringType{},
		"expired_time":            basetypes.StringType{},
		"state":                   basetypes.StringType{},
		"created_at":              basetypes.StringType{},
		"created_by":              basetypes.StringType{},
		"modified_at":             basetypes.StringType{},
		"modified_by":             basetypes.StringType{},
	}
}

// InvitationItem is used for the list items in OrganizationInvitations
type InvitationItem struct {
	Id              types.String `tfsdk:"id"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	AccountName     types.String `tfsdk:"account_name"`
	AccountEmail    types.String `tfsdk:"account_email"`
	LoginId         types.String `tfsdk:"login_id"`
	TargetAccountId types.String `tfsdk:"target_account_id"`
	State           types.String `tfsdk:"state"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

// OrganizationInvitationValue is kept for backward compatibility
type OrganizationInvitationValue = InvitationItem

// GetInvitationItemAttributeTypes returns the attribute types for InvitationItem
func GetInvitationItemAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                basetypes.StringType{},
		"organization_id":   basetypes.StringType{},
		"account_name":      basetypes.StringType{},
		"account_email":     basetypes.StringType{},
		"login_id":          basetypes.StringType{},
		"target_account_id": basetypes.StringType{},
		"state":             basetypes.StringType{},
		"created_at":        basetypes.StringType{},
		"created_by":        basetypes.StringType{},
		"modified_at":       basetypes.StringType{},
		"modified_by":       basetypes.StringType{},
	}
}

// AttributeTypes returns the attribute types for OrganizationInvitationValue
func (o OrganizationInvitationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return GetInvitationItemAttributeTypes()
}

// =====================
// Organization Account Types
// =====================

// AccountListValue - output value for organization account list (matches AccountSummary API response)
type AccountListValue struct {
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Email          types.String `tfsdk:"email"`
	Id             types.String `tfsdk:"id"`
	JoinedMethod   types.String `tfsdk:"joined_method"`
	JoinedTime     types.String `tfsdk:"joined_time"`
	LoginId        types.String `tfsdk:"login_id"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	Name           types.String `tfsdk:"name"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ParentUnitId   types.String `tfsdk:"parent_unit_id"`
	ParentUnitName types.String `tfsdk:"parent_unit_name"`
	State          types.String `tfsdk:"state"`
	AccountType    types.String `tfsdk:"type"`
}

func (o AccountListValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":       basetypes.StringType{},
		"created_by":       basetypes.StringType{},
		"email":            basetypes.StringType{},
		"id":               basetypes.StringType{},
		"joined_method":    basetypes.StringType{},
		"joined_time":      basetypes.StringType{},
		"login_id":         basetypes.StringType{},
		"modified_at":      basetypes.StringType{},
		"modified_by":      basetypes.StringType{},
		"name":             basetypes.StringType{},
		"organization_id":  basetypes.StringType{},
		"parent_unit_id":   basetypes.StringType{},
		"parent_unit_name": basetypes.StringType{},
		"state":            basetypes.StringType{},
		"type":             basetypes.StringType{},
	}
}

// AccountValue - computed/output value for organization account (matches OrganizationAccountWithPolicy API response)
type AccountValue struct {
	ControlPolicies types.List   `tfsdk:"control_policies"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatorName     types.String `tfsdk:"creator_name"`
	Email           types.String `tfsdk:"email"`
	Id              types.String `tfsdk:"id"`
	JoinedMethod    types.String `tfsdk:"joined_method"`
	JoinedTime      types.String `tfsdk:"joined_time"`
	LoginId         types.String `tfsdk:"login_id"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifierName    types.String `tfsdk:"modifier_name"`
	Name            types.String `tfsdk:"name"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	ParentUnitId    types.String `tfsdk:"parent_unit_id"`
	ParentUnitName  types.String `tfsdk:"parent_unit_name"`
	Srn             types.String `tfsdk:"srn"`
	State           types.String `tfsdk:"state"`
	AccountType     types.String `tfsdk:"type"`
}

func (o AccountValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"control_policies": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: ControlPoliciesValue{}.AttributeTypes(ctx),
			},
		},
		"created_at":       basetypes.StringType{},
		"created_by":       basetypes.StringType{},
		"creator_name":     basetypes.StringType{},
		"email":            basetypes.StringType{},
		"id":               basetypes.StringType{},
		"joined_method":    basetypes.StringType{},
		"joined_time":      basetypes.StringType{},
		"login_id":         basetypes.StringType{},
		"modified_at":      basetypes.StringType{},
		"modified_by":      basetypes.StringType{},
		"modifier_name":    basetypes.StringType{},
		"name":             basetypes.StringType{},
		"organization_id":  basetypes.StringType{},
		"parent_unit_id":   basetypes.StringType{},
		"parent_unit_name": basetypes.StringType{},
		"srn":              basetypes.StringType{},
		"state":            basetypes.StringType{},
		"type":             basetypes.StringType{},
	}
}

// FailedValue - nested object for failed account creation response
type FailedValue struct {
	ErrorCode    types.String `tfsdk:"error_code"`
	FailedCaused types.String `tfsdk:"failed_caused"`
	Response     types.String `tfsdk:"response"`
}

func (o FailedValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"error_code":    basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"response":      basetypes.StringType{},
	}
}

// SuccessValue - nested object for successful account creation response
type SuccessValue struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
}

func (o SuccessValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
}

// AccountResource - resource model for organization account
type AccountResource struct {
	Id             types.String `tfsdk:"id"`
	Account        types.Object `tfsdk:"account"`
	AccountId      types.String `tfsdk:"account_id"`
	AccountIds     types.List   `tfsdk:"account_ids"`
	Failed         types.Object `tfsdk:"failed"`
	LazyPolicy     types.Bool   `tfsdk:"lazy_policy"`
	LoginId        types.String `tfsdk:"login_id"`
	Name           types.String `tfsdk:"name"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ParentUnitId   types.String `tfsdk:"parent_unit_id"`
	RoleName       types.String `tfsdk:"role_name"`
	Success        types.Object `tfsdk:"success"`
}

// AccountMoveResource - resource model for moving organization accounts to different OU
type AccountMoveResource struct {
	OrganizationId   types.String `tfsdk:"organization_id"`
	ParentUnitId     types.String `tfsdk:"parent_unit_id"`
	TargetAccountIds types.List   `tfsdk:"target_account_ids"`
	SuccessIds       types.List   `tfsdk:"success_ids"`
}

func (o AccountMoveResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":    basetypes.StringType{},
		"parent_unit_id":     basetypes.StringType{},
		"target_account_ids": basetypes.ListType{ElemType: basetypes.StringType{}},
		"success_ids":        basetypes.ListType{ElemType: types.ObjectType{AttrTypes: MoveAccountSuccessValue{}.AttributeTypes(ctx)}},
	}
}

// MoveAccountSuccessValue - nested object for move account success response
type MoveAccountSuccessValue struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
	TargetId    types.String `tfsdk:"target_id"`
	TargetName  types.String `tfsdk:"target_name"`
}

func (o MoveAccountSuccessValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
		"target_id":    basetypes.StringType{},
		"target_name":  basetypes.StringType{},
	}
}

// AccountDataSource - data source model for organization account
type AccountDataSource struct {
	Account   types.Object `tfsdk:"account"`
	AccountId types.String `tfsdk:"account_id"`
}

type AccountListDataSourceRequest struct {
	OrganizationId  types.String
	Size            types.Int64
	Page            types.Int64
	Sort            types.String
	Id              types.String
	Name            types.String
	Email           types.String
	LoginId         types.String
	JoinedStartDate types.String
	JoinedEndDate   types.String
	JoinedMethod    types.String
	ExcludePolicyId types.String
	ParentUnitId    types.String
	ParentUnitName  types.String
	Type            types.String
}

type AccountListDataSource struct {
	Accounts        types.List   `tfsdk:"accounts"`
	TotalCount      types.Int64  `tfsdk:"total_count"`
	Page            types.Int64  `tfsdk:"page"`
	Size            types.Int64  `tfsdk:"size"`
	SortResult      types.List   `tfsdk:"sort_result"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	Sort            types.String `tfsdk:"sort"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Email           types.String `tfsdk:"email"`
	LoginId         types.String `tfsdk:"login_id"`
	JoinedStartDate types.String `tfsdk:"joined_start_date"`
	JoinedEndDate   types.String `tfsdk:"joined_end_date"`
	JoinedMethod    types.String `tfsdk:"joined_method"`
	ExcludePolicyId types.String `tfsdk:"exclude_policy_id"`
	ParentUnitId    types.String `tfsdk:"parent_unit_id"`
	ParentUnitName  types.String `tfsdk:"parent_unit_name"`
	Type            types.String `tfsdk:"type"`
}

// =====================
// Delegation Policy Types
// =====================

// DelegationPolicyStatementValue - nested object for policy statement
type DelegationPolicyStatementValue struct {
	Action    types.List   `tfsdk:"action"`
	Effect    types.String `tfsdk:"effect"`
	Principal types.Object `tfsdk:"principal"`
	Resource  types.List   `tfsdk:"resource"`
	Sid       types.String `tfsdk:"sid"`
}

func (o DelegationPolicyStatementValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"action":    types.ListType{ElemType: types.StringType},
		"effect":    basetypes.StringType{},
		"principal": types.ObjectType{AttrTypes: DelegationPolicyPrincipalValue{}.AttributeTypes(ctx)},
		"resource":  types.ListType{ElemType: types.StringType},
		"sid":       basetypes.StringType{},
	}
}

// DelegationPolicyPrincipalValue - nested object for principal
type DelegationPolicyPrincipalValue struct {
	Scp types.List `tfsdk:"scp"`
}

func (o DelegationPolicyPrincipalValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"scp": types.ListType{ElemType: types.StringType},
	}
}

// DelegationPolicyDocumentValue - nested object for policy document
type DelegationPolicyDocumentValue struct {
	Statement types.List   `tfsdk:"statement"`
	Version   types.String `tfsdk:"version"`
}

func (o DelegationPolicyDocumentValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"statement": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: DelegationPolicyStatementValue{}.AttributeTypes(ctx),
			},
		},
		"version": basetypes.StringType{},
	}
}

func (o DelegationPolicyDocumentValue) toSdkType() (*organization.DelegationPolicyDocument, error) {
	doc := organization.NewDelegationPolicyDocument()

	var statements []organization.DelegationPolicyStatement
	if !o.Statement.IsNull() && !o.Statement.IsUnknown() {
		statementElems := o.Statement.Elements()
		statements = make([]organization.DelegationPolicyStatement, 0, len(statementElems))
		for _, stmtElem := range statementElems {
			stmtObj, ok := stmtElem.(types.Object)
			if !ok {
				continue
			}

			stmtAttrs := stmtObj.Attributes()
			stmt := organization.NewDelegationPolicyStatement()

			if sidAttr, ok := stmtAttrs["sid"]; ok && !sidAttr.IsNull() && !sidAttr.IsUnknown() {
				stmt.SetSid(sidAttr.(types.String).ValueString())
			}

			if effectAttr, ok := stmtAttrs["effect"]; ok && !effectAttr.IsNull() && !effectAttr.IsUnknown() {
				stmt.SetEffect(effectAttr.(types.String).ValueString())
			}

			if actionAttr, ok := stmtAttrs["action"]; ok && !actionAttr.IsNull() && !actionAttr.IsUnknown() {
				actionList, ok := actionAttr.(types.List)
				if ok {
					actionElems := actionList.Elements()
					actions := make([]string, 0, len(actionElems))
					for _, a := range actionElems {
						if s, ok := a.(types.String); ok {
							actions = append(actions, s.ValueString())
						}
					}
					stmt.SetAction(actions)
				}
			}

			if resourceAttr, ok := stmtAttrs["resource"]; ok && !resourceAttr.IsNull() && !resourceAttr.IsUnknown() {
				resourceList, ok := resourceAttr.(types.List)
				if ok {
					resourceElems := resourceList.Elements()
					resources := make([]string, 0, len(resourceElems))
					for _, r := range resourceElems {
						if s, ok := r.(types.String); ok {
							resources = append(resources, s.ValueString())
						}
					}
					stmt.SetResource(resources)
				}
			}

			// Handle principal - supports both string ("*") and object ({scp: [...]})
			// principal_str takes priority if set
			if principalStrAttr, ok := stmtAttrs["principal_str"]; ok && !principalStrAttr.IsNull() && !principalStrAttr.IsUnknown() {
				principalStr := principalStrAttr.(types.String).ValueString()
				principal := organization.Principal{}
				principal.String = &principalStr
				stmt.SetPrincipal(principal)
			} else if principalAttr, ok := stmtAttrs["principal"]; ok && !principalAttr.IsNull() && !principalAttr.IsUnknown() {
				principalObj, ok := principalAttr.(types.Object)
				if ok {
					principalAttrs := principalObj.Attributes()
					if scpAttr, ok := principalAttrs["scp"]; ok && !scpAttr.IsNull() && !scpAttr.IsUnknown() {
						scpList, ok := scpAttr.(types.List)
						if ok {
							scpElems := scpList.Elements()
							scps := make([]string, 0, len(scpElems))
							for _, s := range scpElems {
								if str, ok := s.(types.String); ok {
									scps = append(scps, str.ValueString())
								}
							}
							scpMap := make(map[string][]string)
							scpMap["scp"] = scps
							principal := organization.Principal{}
							principal.MapmapOfStringarrayOfString = &scpMap
							stmt.SetPrincipal(principal)
						}
					}
				}
			}

			statements = append(statements, *stmt)
		}
		doc.SetStatement(statements)
	}

	if !o.Version.IsNull() && !o.Version.IsUnknown() {
		doc.SetVersion(o.Version.ValueString())
	}

	return doc, nil
}

// DelegationPolicyValue - output value for delegation policy (matches SDK DelegationPolicy structure exactly)
type DelegationPolicyValue struct {
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Document       types.Object `tfsdk:"document"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

func (o DelegationPolicyValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":      basetypes.StringType{},
		"created_by":      basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: DelegationPolicyDocumentValue{}.AttributeTypes(ctx)},
		"modified_at":     basetypes.StringType{},
		"modified_by":     basetypes.StringType{},
		"organization_id": basetypes.StringType{},
	}
}

// DelegationPolicyResource - resource model for delegation policy
type DelegationPolicyResource struct {
	OrganizationId types.String                  `tfsdk:"organization_id"`
	Document       DelegationPolicyDocumentValue `tfsdk:"document"`
	Policy         types.Object                  `tfsdk:"policy"`
}

func (o DelegationPolicyResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: DelegationPolicyDocumentValue{}.AttributeTypes(ctx)},
		"policy":          types.ObjectType{AttrTypes: DelegationPolicyValue{}.AttributeTypes(ctx)},
	}
}

// DelegationPolicyDataSource - data source model for delegation policy
type DelegationPolicyDataSource struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Document       types.Object `tfsdk:"document"`
	Policy         types.Object `tfsdk:"policy"`
}

// ToSdkType converts Terraform model to SDK DelegationPolicyCreateRequest
func (o DelegationPolicyResource) ToSdkType() (*organization.DelegationPolicyCreateRequest, error) {
	doc := organization.NewDelegationPolicyDocument()

	var statements []organization.DelegationPolicyStatement
	if !o.Document.Statement.IsNull() && !o.Document.Statement.IsUnknown() {
		statementElems := o.Document.Statement.Elements()
		statements = make([]organization.DelegationPolicyStatement, 0, len(statementElems))
		for _, stmtElem := range statementElems {
			stmtObj, ok := stmtElem.(types.Object)
			if !ok {
				continue
			}

			stmtAttrs := stmtObj.Attributes()
			stmt := organization.NewDelegationPolicyStatement()

			if sidAttr, ok := stmtAttrs["sid"]; ok && !sidAttr.IsNull() && !sidAttr.IsUnknown() {
				stmt.SetSid(sidAttr.(types.String).ValueString())
			}

			if effectAttr, ok := stmtAttrs["effect"]; ok && !effectAttr.IsNull() && !effectAttr.IsUnknown() {
				stmt.SetEffect(effectAttr.(types.String).ValueString())
			}

			if actionAttr, ok := stmtAttrs["action"]; ok && !actionAttr.IsNull() && !actionAttr.IsUnknown() {
				actionList, ok := actionAttr.(types.List)
				if ok {
					actionElems := actionList.Elements()
					actions := make([]string, 0, len(actionElems))
					for _, a := range actionElems {
						if s, ok := a.(types.String); ok {
							actions = append(actions, s.ValueString())
						}
					}
					stmt.SetAction(actions)
				}
			}

			if resourceAttr, ok := stmtAttrs["resource"]; ok && !resourceAttr.IsNull() && !resourceAttr.IsUnknown() {
				resourceList, ok := resourceAttr.(types.List)
				if ok {
					resourceElems := resourceList.Elements()
					resources := make([]string, 0, len(resourceElems))
					for _, r := range resourceElems {
						if s, ok := r.(types.String); ok {
							resources = append(resources, s.ValueString())
						}
					}
					stmt.SetResource(resources)
				}
			}

			// Handle principal - supports both string ("*") and object ({scp: [...]})
			if principalStrAttr, ok := stmtAttrs["principal_str"]; ok && !principalStrAttr.IsNull() && !principalStrAttr.IsUnknown() {
				principalStr := principalStrAttr.(types.String).ValueString()
				principal := organization.Principal{}
				principal.String = &principalStr
				stmt.SetPrincipal(principal)
			} else if principalAttr, ok := stmtAttrs["principal"]; ok && !principalAttr.IsNull() && !principalAttr.IsUnknown() {
				principalObj, ok := principalAttr.(types.Object)
				if ok {
					principalAttrs := principalObj.Attributes()
					if scpAttr, ok := principalAttrs["scp"]; ok && !scpAttr.IsNull() && !scpAttr.IsUnknown() {
						scpList, ok := scpAttr.(types.List)
						if ok {
							scpElems := scpList.Elements()
							scps := make([]string, 0, len(scpElems))
							for _, s := range scpElems {
								if str, ok := s.(types.String); ok {
									scps = append(scps, str.ValueString())
								}
							}
							scpMap := make(map[string][]string)
							scpMap["scp"] = scps
							principal := organization.Principal{}
							principal.MapmapOfStringarrayOfString = &scpMap
							stmt.SetPrincipal(principal)
						}
					}
				}
			}

			statements = append(statements, *stmt)
		}
		doc.SetStatement(statements)
	}

	if !o.Document.Version.IsNull() && !o.Document.Version.IsUnknown() {
		doc.SetVersion(o.Document.Version.ValueString())
	}

	// Create DelegationPolicyCreateRequest
	req := organization.NewDelegationPolicyCreateRequest(*doc)
	if !o.OrganizationId.IsNull() && !o.OrganizationId.IsUnknown() {
		req.SetOrganizationId(o.OrganizationId.ValueString())
	}

	return req, nil
}

// =====================
// Service Control Policy Types
// =====================

type ServiceControlPolicyStatementValue struct {
	Action    types.List   `tfsdk:"action"`
	Condition types.Map    `tfsdk:"condition"`
	Effect    types.String `tfsdk:"effect"`
	NotAction types.List   `tfsdk:"not_action"`
	Principal types.String `tfsdk:"principal"`
	Resource  types.List   `tfsdk:"resource"`
	Sid       types.String `tfsdk:"sid"`
}

func (o ServiceControlPolicyStatementValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"action":     types.ListType{ElemType: types.StringType},
		"condition":  types.MapType{ElemType: types.StringType},
		"effect":     basetypes.StringType{},
		"not_action": types.ListType{ElemType: types.StringType},
		"principal":  basetypes.StringType{},
		"resource":   types.ListType{ElemType: types.StringType},
		"sid":        basetypes.StringType{},
	}
}

type ServiceControlPolicyDocumentValue struct {
	Statement types.List   `tfsdk:"statement"`
	Version   types.String `tfsdk:"version"`
}

func (o ServiceControlPolicyDocumentValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"statement": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: ServiceControlPolicyStatementValue{}.AttributeTypes(ctx),
			},
		},
		"version": basetypes.StringType{},
	}
}

// ServiceControlPolicyValue - output value for service control policy (matches SDK ServiceControlPolicy structure exactly)
type ServiceControlPolicyValue struct {
	Category       types.String `tfsdk:"category"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatorName    types.String `tfsdk:"creator_name"`
	Description    types.String `tfsdk:"description"`
	Document       types.Object `tfsdk:"document"`
	Id             types.String `tfsdk:"id"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	ModifierName   types.String `tfsdk:"modifier_name"`
	Name           types.String `tfsdk:"name"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Source         types.String `tfsdk:"source"`
	State          types.String `tfsdk:"state"`
	Type           types.String `tfsdk:"type"`
}

func (o ServiceControlPolicyValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"category":        basetypes.StringType{},
		"created_at":      basetypes.StringType{},
		"created_by":      basetypes.StringType{},
		"creator_name":    basetypes.StringType{},
		"description":     basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)},
		"id":              basetypes.StringType{},
		"modified_at":     basetypes.StringType{},
		"modified_by":     basetypes.StringType{},
		"modifier_name":   basetypes.StringType{},
		"name":            basetypes.StringType{},
		"organization_id": basetypes.StringType{},
		"source":          basetypes.StringType{},
		"state":           basetypes.StringType{},
		"type":            basetypes.StringType{},
	}
}

type ServiceControlPolicyResource struct {
	PolicyId       types.String `tfsdk:"policy_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Document       types.Object `tfsdk:"document"`
	Type           types.String `tfsdk:"type"`
	Category       types.String `tfsdk:"category"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatorName    types.String `tfsdk:"creator_name"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	ModifierName   types.String `tfsdk:"modifier_name"`
	ServiceName    types.String `tfsdk:"service_name"`
	Source         types.String `tfsdk:"source"`
	State          types.String `tfsdk:"state"`
}

func (o ServiceControlPolicyResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"policy_id":       basetypes.StringType{},
		"organization_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"description":     basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)},
		"type":            basetypes.StringType{},
		"category":        basetypes.StringType{},
		"created_at":      basetypes.StringType{},
		"created_by":      basetypes.StringType{},
		"creator_name":    basetypes.StringType{},
		"modified_at":     basetypes.StringType{},
		"modified_by":     basetypes.StringType{},
		"modifier_name":   basetypes.StringType{},
		"service_name":    basetypes.StringType{},
		"source":          basetypes.StringType{},
		"state":           basetypes.StringType{},
	}
}

// ServiceControlPolicyDataSource - data source model for service control policy
type ServiceControlPolicyDataSource struct {
	PolicyId       types.String `tfsdk:"policy_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Document       types.Object `tfsdk:"document"`
	Category       types.String `tfsdk:"category"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatorName    types.String `tfsdk:"creator_name"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	ModifierName   types.String `tfsdk:"modifier_name"`
	Source         types.String `tfsdk:"source"`
	State          types.String `tfsdk:"state"`
	Srn            types.String `tfsdk:"srn"`
	ServiceName    types.String `tfsdk:"service_name"`
}

func (o ServiceControlPolicyDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"policy_id":       basetypes.StringType{},
		"organization_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"description":     basetypes.StringType{},
		"type":            basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)},
		"category":        basetypes.StringType{},
		"created_at":      basetypes.StringType{},
		"created_by":      basetypes.StringType{},
		"creator_name":    basetypes.StringType{},
		"modified_at":     basetypes.StringType{},
		"modified_by":     basetypes.StringType{},
		"modifier_name":   basetypes.StringType{},
		"source":          basetypes.StringType{},
		"state":           basetypes.StringType{},
		"srn":             basetypes.StringType{},
		"service_name":    basetypes.StringType{},
	}
}

// ServiceControlPolicyListDataSource - data source model for listing service control policies
type ServiceControlPolicyListDataSource struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Policies       types.List   `tfsdk:"policies"`
	TotalCount     types.Int64  `tfsdk:"total_count"`
	Page           types.Int64  `tfsdk:"page"`
	Size           types.Int64  `tfsdk:"size"`
	Sort           types.List   `tfsdk:"sort"`
}

func (o ServiceControlPolicyListDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"type":            basetypes.StringType{},
		"policies":        types.ListType{ElemType: types.ObjectType{AttrTypes: ServiceControlPolicyListValue{}.AttributeTypes(ctx)}},
		"total_count":     basetypes.Int64Type{},
		"page":            basetypes.Int64Type{},
		"size":            basetypes.Int64Type{},
		"sort":            types.ListType{ElemType: types.StringType},
	}
}

type ServiceControlPolicyListValue struct {
	PolicyId       types.String `tfsdk:"policy_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Category       types.String `tfsdk:"category"`
	State          types.String `tfsdk:"state"`
	Source         types.String `tfsdk:"source"`
	Document       types.Object `tfsdk:"document"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatorName    types.String `tfsdk:"creator_name"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	ModifierName   types.String `tfsdk:"modifier_name"`
}

func (o ServiceControlPolicyListValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"policy_id":       basetypes.StringType{},
		"organization_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"description":     basetypes.StringType{},
		"type":            basetypes.StringType{},
		"category":        basetypes.StringType{},
		"state":           basetypes.StringType{},
		"source":          basetypes.StringType{},
		"document":        types.ObjectType{AttrTypes: ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)},
		"created_at":      basetypes.StringType{},
		"created_by":      basetypes.StringType{},
		"creator_name":    basetypes.StringType{},
		"modified_at":     basetypes.StringType{},
		"modified_by":     basetypes.StringType{},
		"modifier_name":   basetypes.StringType{},
	}
}

// ServiceControlPolicyDeleteRequest - request model for deleting service control policies
type ServiceControlPolicyDeleteRequest struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Ids            types.List   `tfsdk:"ids"`
}

func (o ServiceControlPolicyDeleteRequest) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"ids":             types.ListType{ElemType: types.StringType},
	}
}

// ServiceControlPolicyDeleteFailCausedValue - nested object for failed deletion
type ServiceControlPolicyDeleteFailCausedValue struct {
	BindingTargets types.List   `tfsdk:"binding_targets"`
	ErrorCode      types.String `tfsdk:"error_code"`
	FailedCaused   types.String `tfsdk:"failed_caused"`
	Response       types.String `tfsdk:"response"`
}

func (o ServiceControlPolicyDeleteFailCausedValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"binding_targets": types.ListType{ElemType: types.StringType},
		"error_code":      basetypes.StringType{},
		"failed_caused":   basetypes.StringType{},
		"response":        basetypes.StringType{},
	}
}

// ServiceControlPolicyDeleteResponseValue - output value for delete operation
type ServiceControlPolicyDeleteResponseValue struct {
	FailedIds  types.List `tfsdk:"failed_ids"`
	SuccessIds types.List `tfsdk:"success_ids"`
}

func (o ServiceControlPolicyDeleteResponseValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"failed_ids":  types.ListType{ElemType: types.ObjectType{AttrTypes: ServiceControlPolicyDeleteFailCausedValue{}.AttributeTypes(ctx)}},
		"success_ids": types.ListType{ElemType: types.StringType},
	}
}

type ServiceControlPolicyListFilter struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
}

func (o ServiceControlPolicyListFilter) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"type":            basetypes.StringType{},
	}
}

type ServiceControlPolicyListOutput struct {
	OrganizationId  types.String `tfsdk:"organization_id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Id              types.String `tfsdk:"id"`
	ExcludeTargetId types.String `tfsdk:"exclude_target_id"`
	Policies        types.List   `tfsdk:"policies"`
	TotalCount      types.Int64  `tfsdk:"total_count"`
	Page            types.Int64  `tfsdk:"page"`
	Size            types.Int64  `tfsdk:"size"`
	Sort            types.String `tfsdk:"sort"`
	SortResult      types.List   `tfsdk:"sort_result"`
}

func (o ServiceControlPolicyListOutput) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":   basetypes.StringType{},
		"name":              basetypes.StringType{},
		"type":              basetypes.StringType{},
		"id":                basetypes.StringType{},
		"exclude_target_id": basetypes.StringType{},
		"policies":          types.ListType{ElemType: types.ObjectType{AttrTypes: ServiceControlPolicyListValue{}.AttributeTypes(ctx)}},
		"total_count":       basetypes.Int64Type{},
		"page":              basetypes.Int64Type{},
		"size":              basetypes.Int64Type{},
		"sort":              basetypes.StringType{},
		"sort_result":       types.ListType{ElemType: types.StringType},
	}
}
