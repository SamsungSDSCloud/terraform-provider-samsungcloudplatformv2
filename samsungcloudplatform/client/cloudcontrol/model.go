package cloudcontrol

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-cloudcontrol"

// LandingZoneResource - resource model for cloudcontrol landing zone
type LandingZoneResource struct {
	AdditionalOuName         types.String `tfsdk:"additional_ou_name"`
	AgreeYn                  types.String `tfsdk:"agree_yn"`
	AuditAccountName         types.String `tfsdk:"audit_account_name"`
	AuditLoginId             types.String `tfsdk:"audit_login_id"`
	BasicOuName              types.String `tfsdk:"basic_ou_name"`
	DetectiveGuardrailStatus types.String `tfsdk:"detective_guardrail_status"`
	LogArchiveAccountName    types.String `tfsdk:"log_archive_account_name"`
	LogArchiveLoginId        types.String `tfsdk:"log_archive_login_id"`
	SsoType                  types.String `tfsdk:"sso_type"`

	// Computed fields from response
	JobId            types.String `tfsdk:"job_id"`
	LandingZoneId    types.String `tfsdk:"landing_zone_id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	CreatorName      types.String `tfsdk:"creator_name"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	ModifierName     types.String `tfsdk:"modifier_name"`
	OrganizationId   types.String `tfsdk:"organization_id"`
	Region           types.String `tfsdk:"region"`
	ServiceName      types.String `tfsdk:"service_name"`
	Srn              types.String `tfsdk:"srn"`
	Status           types.String `tfsdk:"status"`
	VersionId        types.String `tfsdk:"version_id"`
	IdentityCenterId types.String `tfsdk:"identity_center_id"`
}

func (o LandingZoneResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"additional_ou_name":         basetypes.StringType{},
		"agree_yn":                   basetypes.StringType{},
		"audit_account_name":         basetypes.StringType{},
		"audit_login_id":             basetypes.StringType{},
		"basic_ou_name":              basetypes.StringType{},
		"detective_guardrail_status": basetypes.StringType{},
		"log_archive_account_name":   basetypes.StringType{},
		"log_archive_login_id":       basetypes.StringType{},
		"sso_type":                   basetypes.StringType{},
		"job_id":                     basetypes.StringType{},
		"landing_zone_id":            basetypes.StringType{},
		"created_at":                 basetypes.StringType{},
		"created_by":                 basetypes.StringType{},
		"creator_name":               basetypes.StringType{},
		"modified_at":                basetypes.StringType{},
		"modified_by":                basetypes.StringType{},
		"modifier_name":              basetypes.StringType{},
		"organization_id":            basetypes.StringType{},
		"region":                     basetypes.StringType{},
		"service_name":               basetypes.StringType{},
		"srn":                        basetypes.StringType{},
		"status":                     basetypes.StringType{},
		"version_id":                 basetypes.StringType{},
		"identity_center_id":         basetypes.StringType{},
	}
}

// ToSdkType converts Terraform model to SDK LandingZoneCreateRequestV1Dot1
func (o LandingZoneResource) ToSdkType() (*cloudcontrol.LandingZoneCreateRequestV1Dot1, error) {
	req := cloudcontrol.NewLandingZoneCreateRequestV1Dot1(
		o.AdditionalOuName.ValueString(),
		cloudcontrol.CommonEnum(o.AgreeYn.ValueString()),
		o.AuditAccountName.ValueString(),
		o.AuditLoginId.ValueString(),
		o.BasicOuName.ValueString(),
		cloudcontrol.DetectiveGuardrailStatusEnum(o.DetectiveGuardrailStatus.ValueString()),
		o.LogArchiveAccountName.ValueString(),
		o.LogArchiveLoginId.ValueString(),
		cloudcontrol.SsoTypeEnum(o.SsoType.ValueString()),
	)
	return req, nil
}

// LandingZoneUpdateResource - resource model for updating cloudcontrol landing zone
type LandingZoneUpdateResource struct {
	Id                       types.String `tfsdk:"id"`
	DetectiveGuardrailStatus types.String `tfsdk:"detective_guardrail_status"`
}

func (o LandingZoneUpdateResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                         basetypes.StringType{},
		"detective_guardrail_status": basetypes.StringType{},
	}
}

// ToSdkType converts Terraform model to SDK LandingZoneUpdateRequest
func (o LandingZoneUpdateResource) ToSdkType() (*cloudcontrol.LandingZoneUpdateRequest, error) {
	req := cloudcontrol.NewLandingZoneUpdateRequest(
		o.DetectiveGuardrailStatus.ValueString(),
	)
	return req, nil
}

// LandingZoneDataSource - data source model for cloudcontrol landing zone
type LandingZoneDataSource struct {
	LandingZoneId types.String `tfsdk:"landing_zone_id"`
	LandingZone   types.Object `tfsdk:"landing_zone"`
}

func (o LandingZoneDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"landing_zone_id": basetypes.StringType{},
		"landing_zone":    types.ObjectType{AttrTypes: LandingZoneValue{}.AttributeTypes(ctx)},
	}
}

// LandingZoneValue - computed/output value for cloudcontrol landing zone
type LandingZoneValue struct {
	AgreeYn                  types.String `tfsdk:"agree_yn"`
	CreatedAt                types.String `tfsdk:"created_at"`
	CreatedBy                types.String `tfsdk:"created_by"`
	CreatorName              types.String `tfsdk:"creator_name"`
	DetectiveGuardrailStatus types.String `tfsdk:"detective_guardrail_status"`
	Id                       types.String `tfsdk:"id"`
	IdentityCenterId         types.String `tfsdk:"identity_center_id"`
	ModifiedAt               types.String `tfsdk:"modified_at"`
	ModifiedBy               types.String `tfsdk:"modified_by"`
	ModifierName             types.String `tfsdk:"modifier_name"`
	OrganizationId           types.String `tfsdk:"organization_id"`
	Region                   types.String `tfsdk:"region"`
	ServiceName              types.String `tfsdk:"service_name"`
	Srn                      types.String `tfsdk:"srn"`
	SsoType                  types.String `tfsdk:"sso_type"`
	Status                   types.String `tfsdk:"status"`
	VersionId                types.String `tfsdk:"version_id"`
}

func (o LandingZoneValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"agree_yn":                   basetypes.StringType{},
		"created_at":                 basetypes.StringType{},
		"created_by":                 basetypes.StringType{},
		"creator_name":               basetypes.StringType{},
		"detective_guardrail_status": basetypes.StringType{},
		"id":                         basetypes.StringType{},
		"identity_center_id":         basetypes.StringType{},
		"modified_at":                basetypes.StringType{},
		"modified_by":                basetypes.StringType{},
		"modifier_name":              basetypes.StringType{},
		"organization_id":            basetypes.StringType{},
		"region":                     basetypes.StringType{},
		"service_name":               basetypes.StringType{},
		"srn":                        basetypes.StringType{},
		"sso_type":                   basetypes.StringType{},
		"status":                     basetypes.StringType{},
		"version_id":                 basetypes.StringType{},
	}
}

// AccountFactoryResource - resource model for cloudcontrol account factory
type AccountFactoryResource struct {
	LandingZoneId   types.String `tfsdk:"landing_zone_id"`
	Name            types.String `tfsdk:"name"`
	ParentUnitId    types.String `tfsdk:"parent_unit_id"`
	Email           types.String `tfsdk:"email"`
	SsoUserEmail    types.String `tfsdk:"sso_user_email"`
	SsoUserName     types.String `tfsdk:"sso_user_name"`
	SsoUserRealName types.String `tfsdk:"sso_user_real_name"`

	// Computed fields from response
	AccountId types.String `tfsdk:"account_id"`
	JobId     types.String `tfsdk:"job_id"`
}

func (o AccountFactoryResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"landing_zone_id":    basetypes.StringType{},
		"name":               basetypes.StringType{},
		"parent_unit_id":     basetypes.StringType{},
		"email":              basetypes.StringType{},
		"sso_user_email":     basetypes.StringType{},
		"sso_user_name":      basetypes.StringType{},
		"sso_user_real_name": basetypes.StringType{},
		"account_id":         basetypes.StringType{},
		"job_id":             basetypes.StringType{},
	}
}

// ToSdkType converts Terraform model to SDK AccountFactoryCreateRequest
func (o AccountFactoryResource) ToSdkType() (*cloudcontrol.AccountFactoryCreateRequest, error) {
	req := cloudcontrol.NewAccountFactoryCreateRequest(
		o.Email.ValueString(),
		o.LandingZoneId.ValueString(),
		o.Name.ValueString(),
		o.ParentUnitId.ValueString(),
	)

	// Set optional fields if provided
	if !o.SsoUserEmail.IsNull() && o.SsoUserEmail.ValueString() != "" {
		req.SetSsoUserEmail(o.SsoUserEmail.ValueString())
	}
	if !o.SsoUserName.IsNull() && o.SsoUserName.ValueString() != "" {
		req.SetSsoUserName(o.SsoUserName.ValueString())
	}
	if !o.SsoUserRealName.IsNull() && o.SsoUserRealName.ValueString() != "" {
		req.SetSsoUserRealName(o.SsoUserRealName.ValueString())
	}

	return req, nil
}

// BaselineAssignmentResource - resource model for cloudcontrol baseline assignment
type BaselineAssignmentResource struct {
	AssignmentId    types.String `tfsdk:"assignment_id"`
	LandingZoneId   types.String `tfsdk:"landing_zone_id"`
	ResourceType    types.String `tfsdk:"resource_type"`
	AgreeYn         types.String `tfsdk:"agree_yn"`
	ParentUnitId    types.String `tfsdk:"parent_unit_id"`
	SsoUserName     types.String `tfsdk:"sso_user_name"`
	SsoUserRealName types.String `tfsdk:"sso_user_real_name"`

	// Computed fields from response
	JobId                    types.String `tfsdk:"job_id"`
	Status                   types.String `tfsdk:"status"`
	AccountCount             types.Int64  `tfsdk:"account_count"`
	AccountAssignedCount     types.Int64  `tfsdk:"account_assigned_count"`
	OuCount                  types.Int64  `tfsdk:"ou_count"`
	OuAssignedCount          types.Int64  `tfsdk:"ou_assigned_count"`
	DetectiveGuardrailStatus types.String `tfsdk:"detective_guardrail_status"`
	DetectiveGuardrailType   types.String `tfsdk:"detective_guardrail_type"`
	ReregisterTrigger        types.String `tfsdk:"reregister_trigger"`
}

func (o BaselineAssignmentResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"assignment_id":              basetypes.StringType{},
		"landing_zone_id":            basetypes.StringType{},
		"resource_type":              basetypes.StringType{},
		"agree_yn":                   basetypes.StringType{},
		"reregister_trigger":         basetypes.StringType{},
		"parent_unit_id":             basetypes.StringType{},
		"sso_user_name":              basetypes.StringType{},
		"sso_user_real_name":         basetypes.StringType{},
		"job_id":                     basetypes.StringType{},
		"status":                     basetypes.StringType{},
		"account_count":              basetypes.Int64Type{},
		"account_assigned_count":     basetypes.Int64Type{},
		"ou_count":                   basetypes.Int64Type{},
		"ou_assigned_count":          basetypes.Int64Type{},
		"detective_guardrail_status": basetypes.StringType{},
		"detective_guardrail_type":   basetypes.StringType{},
	}
}

// ToSdkType converts Terraform model to SDK BaselineAssignmentAddRequest
func (o BaselineAssignmentResource) ToSdkType() (*cloudcontrol.BaselineAssignmentAddRequest, error) {
	req := cloudcontrol.NewBaselineAssignmentAddRequest(
		o.LandingZoneId.ValueString(),
		cloudcontrol.BaselineAssignmentTypeEnum(o.ResourceType.ValueString()),
	)

	// Set optional fields if provided
	if !o.AgreeYn.IsNull() && o.AgreeYn.ValueString() != "" {
		req.SetAgreeYn(cloudcontrol.CommonEnum(o.AgreeYn.ValueString()))
	}
	if !o.ParentUnitId.IsNull() && o.ParentUnitId.ValueString() != "" {
		req.SetParentUnitId(o.ParentUnitId.ValueString())
	}
	if !o.SsoUserName.IsNull() && o.SsoUserName.ValueString() != "" {
		req.SetSsoUserName(o.SsoUserName.ValueString())
	}
	if !o.SsoUserRealName.IsNull() && o.SsoUserRealName.ValueString() != "" {
		req.SetSsoUserRealName(o.SsoUserRealName.ValueString())
	}

	return req, nil
}

// BaselineAssignmentUpdateResource - resource model for updating cloudcontrol baseline assignment
type BaselineAssignmentUpdateResource struct {
	LandingZoneId types.String `tfsdk:"landing_zone_id"`
	AgreeYn       types.String `tfsdk:"agree_yn"`
}

func (o BaselineAssignmentUpdateResource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"landing_zone_id": basetypes.StringType{},
		"agree_yn":        basetypes.StringType{},
	}
}

// ToSdkType converts Terraform model to SDK BaselineAssignmentUpdateRequest
func (o BaselineAssignmentUpdateResource) ToSdkType() (*cloudcontrol.BaselineAssignmentUpdateRequest, error) {
	req := cloudcontrol.NewBaselineAssignmentUpdateRequest(
		o.LandingZoneId.ValueString(),
	)

	// Set optional fields if provided
	if !o.AgreeYn.IsNull() && o.AgreeYn.ValueString() != "" {
		req.SetAgreeYn(cloudcontrol.CommonEnum(o.AgreeYn.ValueString()))
	}

	return req, nil
}

// BaselineAssignmentDataSource - data source model for listing cloudcontrol baseline assignments
// (모든 필터는 선택값 - API 스펙 그대로)
type BaselineAssignmentDataSource struct {
	LandingZoneId       types.String                 `tfsdk:"landing_zone_id"`
	ResourceType        types.String                 `tfsdk:"resource_type"`
	AssignmentId        types.String                 `tfsdk:"assignment_id"`
	Status              types.String                 `tfsdk:"status"`
	BaselineAssignments []BaselineAssignmentListItem `tfsdk:"baseline_assignments"`
}

// BaselineAssignmentListItem is a single item in the baseline_assignments list
type BaselineAssignmentListItem struct {
	Id                                types.String `tfsdk:"id"`
	Type                              types.String `tfsdk:"type"`
	AccountCount                      types.Int64  `tfsdk:"account_count"`
	AccountAssignedCount              types.Int64  `tfsdk:"account_assigned_count"`
	OuCount                           types.Int64  `tfsdk:"ou_count"`
	OuAssignedCount                   types.Int64  `tfsdk:"ou_assigned_count"`
	PreventiveGuardrailActivatedCount types.Int64  `tfsdk:"preventive_guardrail_activated_count"`
	DetectiveGuardrailStatus          types.String `tfsdk:"detective_guardrail_status"`
	DetectiveGuardrailType            types.String `tfsdk:"detective_guardrail_type"`
	SsoUserName                       types.String `tfsdk:"sso_user_name"`
	Status                            types.String `tfsdk:"status"`
}

func (o BaselineAssignmentListItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                   basetypes.StringType{},
		"type":                                 basetypes.StringType{},
		"account_count":                        basetypes.Int64Type{},
		"account_assigned_count":               basetypes.Int64Type{},
		"ou_count":                             basetypes.Int64Type{},
		"ou_assigned_count":                    basetypes.Int64Type{},
		"preventive_guardrail_activated_count": basetypes.Int64Type{},
		"detective_guardrail_status":           basetypes.StringType{},
		"detective_guardrail_type":             basetypes.StringType{},
		"sso_user_name":                        basetypes.StringType{},
		"status":                               basetypes.StringType{},
	}
}

type GuardrailListDataSource struct {
	LandingZoneId types.String        `tfsdk:"landing_zone_id"`
	Name          types.String        `tfsdk:"name"`
	ExcludeUnitId types.String        `tfsdk:"exclude_unit_id"`
	Guidance      types.String        `tfsdk:"guidance"`
	ServiceName   types.String        `tfsdk:"service_name"`
	Status        types.String        `tfsdk:"status"`
	Size          types.Int64         `tfsdk:"size"`
	Page          types.Int64         `tfsdk:"page"`
	Sort          types.String        `tfsdk:"sort"`
	TotalCount    types.Int64         `tfsdk:"total_count"`
	SortResult    types.List          `tfsdk:"sort_result"`
	Guardrails    []GuardrailListItem `tfsdk:"guardrails"`
}

// GuardrailListItem is a single item in the guardrails list
type GuardrailListItem struct {
	Id          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Description types.String             `tfsdk:"description"`
	Guidance    types.String             `tfsdk:"guidance"`
	ServiceName types.String             `tfsdk:"service_name"`
	Status      types.String             `tfsdk:"status"`
	Type        types.String             `tfsdk:"type"`
	CreatedAt   types.String             `tfsdk:"created_at"`
	CreatedBy   types.String             `tfsdk:"created_by"`
	ModifiedAt  types.String             `tfsdk:"modified_at"`
	ModifiedBy  types.String             `tfsdk:"modified_by"`
	Srn         types.String             `tfsdk:"srn"`
	BindingOus  []GuardrailBindingOuItem `tfsdk:"binding_ous"`
}

func (o GuardrailListItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":           basetypes.StringType{},
		"name":         basetypes.StringType{},
		"description":  basetypes.StringType{},
		"guidance":     basetypes.StringType{},
		"service_name": basetypes.StringType{},
		"status":       basetypes.StringType{},
		"type":         basetypes.StringType{},
		"created_at":   basetypes.StringType{},
		"created_by":   basetypes.StringType{},
		"modified_at":  basetypes.StringType{},
		"modified_by":  basetypes.StringType{},
		"srn":          basetypes.StringType{},
		"binding_ous": types.ListType{
			ElemType: types.ObjectType{AttrTypes: GuardrailBindingOuItem{}.AttributeTypes(ctx)},
		},
	}
}

// GuardrailBindingOuItem represents a single organization unit that a guardrail is bound to
type GuardrailBindingOuItem struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (o GuardrailBindingOuItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":   basetypes.StringType{},
		"name": basetypes.StringType{},
	}
}

// GuardrailDataSource - data source model for single guardrail lookup (상세조회)
type GuardrailDataSource struct {
	GuardrailId   types.String `tfsdk:"guardrail_id"`
	LandingZoneId types.String `tfsdk:"landing_zone_id"`
	Guardrail     types.Object `tfsdk:"guardrail"`
}

func (o GuardrailDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"guardrail_id":    basetypes.StringType{},
		"landing_zone_id": basetypes.StringType{},
		"guardrail":       types.ObjectType{AttrTypes: GuardrailValue{}.AttributeTypes(ctx)},
	}
}

// GuardrailValue - single guardrail object returned by the 상세조회 API
type GuardrailValue struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Guidance    types.String `tfsdk:"guidance"`
	ServiceName types.String `tfsdk:"service_name"`
	Status      types.String `tfsdk:"status"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Srn         types.String `tfsdk:"srn"`
	BindingOus  types.List   `tfsdk:"binding_ous"`
}

func (o GuardrailValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":           basetypes.StringType{},
		"name":         basetypes.StringType{},
		"description":  basetypes.StringType{},
		"guidance":     basetypes.StringType{},
		"service_name": basetypes.StringType{},
		"status":       basetypes.StringType{},
		"type":         basetypes.StringType{},
		"created_at":   basetypes.StringType{},
		"created_by":   basetypes.StringType{},
		"modified_at":  basetypes.StringType{},
		"modified_by":  basetypes.StringType{},
		"srn":          basetypes.StringType{},
		"binding_ous": types.ListType{
			ElemType: types.ObjectType{AttrTypes: GuardrailBindingOuItem{}.AttributeTypes(ctx)},
		},
	}
}

type GuardrailBindingResource struct {
	LandingZoneId types.String `tfsdk:"landing_zone_id"`
	GuardrailIds  types.List   `tfsdk:"guardrail_ids"`
	UnitIds       types.List   `tfsdk:"unit_ids"`
	SuccessIds    types.List   `tfsdk:"success_ids"`
	FailedIds     types.List   `tfsdk:"failed_ids"`
}

// GuardrailBindingSuccessItem - a single successfully bound guardrail/unit pair
type GuardrailBindingSuccessItem struct {
	GuardrailId types.String `tfsdk:"guardrail_id"`
	UnitId      types.String `tfsdk:"unit_id"`
}

func (o GuardrailBindingSuccessItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"guardrail_id": basetypes.StringType{},
		"unit_id":      basetypes.StringType{},
	}
}

// GuardrailBindingFailedItem - a single failed guardrail/unit binding attempt
type GuardrailBindingFailedItem struct {
	GuardrailId  types.String `tfsdk:"guardrail_id"`
	UnitId       types.String `tfsdk:"unit_id"`
	ErrorCode    types.String `tfsdk:"error_code"`
	FailedCaused types.String `tfsdk:"failed_caused"`
	Response     types.Map    `tfsdk:"response"`
}

func (o GuardrailBindingFailedItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"guardrail_id":  basetypes.StringType{},
		"unit_id":       basetypes.StringType{},
		"error_code":    basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"response":      types.MapType{ElemType: basetypes.StringType{}},
	}
}

// ============================================================
// GET /cloudcontrol/v1/guardrail-bindings/guardrails
// (특정 target에 걸린 가드레일 목록 조회)
// ============================================================

// GuardrailBindingDataSource - data source model for listing guardrails bound to a target
type GuardrailBindingDataSource struct {
	TargetId      types.String              `tfsdk:"target_id"`
	LandingZoneId types.String              `tfsdk:"landing_zone_id"`
	Name          types.String              `tfsdk:"name"`
	Sort          types.String              `tfsdk:"sort"`
	Size          types.Int64               `tfsdk:"size"`
	Page          types.Int64               `tfsdk:"page"`
	TotalCount    types.Int64               `tfsdk:"total_count"`
	SortResult    types.List                `tfsdk:"sort_result"`
	Guardrails    []GuardrailsForTargetItem `tfsdk:"guardrails"`
}

func (o GuardrailBindingDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"target_id":       basetypes.StringType{},
		"landing_zone_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"sort":            basetypes.StringType{},
		"size":            basetypes.Int64Type{},
		"page":            basetypes.Int64Type{},
		"count":           basetypes.Int64Type{},
		"sort_result": types.ListType{
			ElemType: basetypes.StringType{},
		},
		"guardrails": types.ListType{
			ElemType: types.ObjectType{AttrTypes: GuardrailsForTargetItem{}.AttributeTypes(ctx)},
		},
	}
}

// GuardrailsForTargetItem - single guardrail entry in the guardrails-for-target list
type GuardrailsForTargetItem struct {
	Id          types.String            `tfsdk:"id"`
	Name        types.String            `tfsdk:"name"`
	Guidance    types.String            `tfsdk:"guidance"`
	ServiceName types.String            `tfsdk:"service_name"`
	Type        types.String            `tfsdk:"type"`
	CreatedAt   types.String            `tfsdk:"created_at"`
	CreatedBy   types.String            `tfsdk:"created_by"`
	ModifiedAt  types.String            `tfsdk:"modified_at"`
	ModifiedBy  types.String            `tfsdk:"modified_by"`
	LinkTypes   *GuardrailLinkTypesItem `tfsdk:"link_types"`
}

func (o GuardrailsForTargetItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":           basetypes.StringType{},
		"name":         basetypes.StringType{},
		"guidance":     basetypes.StringType{},
		"service_name": basetypes.StringType{},
		"type":         basetypes.StringType{},
		"created_at":   basetypes.StringType{},
		"created_by":   basetypes.StringType{},
		"modified_at":  basetypes.StringType{},
		"modified_by":  basetypes.StringType{},
		"link_types": types.ObjectType{
			AttrTypes: GuardrailLinkTypesItem{}.AttributeTypes(ctx),
		},
	}
}

// GuardrailLinkTypesItem - the link_types object (directed / inherited)
type GuardrailLinkTypesItem struct {
	Directed  []GuardrailLinkTargetItem `tfsdk:"directed"`
	Inherited []GuardrailLinkTargetItem `tfsdk:"inherited"`
}

func (o GuardrailLinkTypesItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"directed": types.ListType{
			ElemType: types.ObjectType{AttrTypes: GuardrailLinkTargetItem{}.AttributeTypes(ctx)},
		},
		"inherited": types.ListType{
			ElemType: types.ObjectType{AttrTypes: GuardrailLinkTargetItem{}.AttributeTypes(ctx)},
		},
	}
}

// GuardrailLinkTargetItem - a single {target_id, target_name} entry inside directed/inherited
type GuardrailLinkTargetItem struct {
	TargetId   types.String `tfsdk:"target_id"`
	TargetName types.String `tfsdk:"target_name"`
}

func (o GuardrailLinkTargetItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"target_id":   basetypes.StringType{},
		"target_name": basetypes.StringType{},
	}
}

// ============================================================
// GET /cloudcontrol/v1/guardrail-bindings/targets
// (특정 가드레일이 걸린 target(ACCOUNT/OU) 목록 조회)
// ============================================================

// GuardrailBindingTargetsDataSource - data source model for listing targets bound to a guardrail
type GuardrailBindingTargetsDataSource struct {
	GuardrailId   types.String              `tfsdk:"guardrail_id"`
	TargetType    types.String              `tfsdk:"target_type"`
	LandingZoneId types.String              `tfsdk:"landing_zone_id"`
	Name          types.String              `tfsdk:"name"`
	Sort          types.String              `tfsdk:"sort"`
	Size          types.Int64               `tfsdk:"size"`
	Page          types.Int64               `tfsdk:"page"`
	TotalCount    types.Int64               `tfsdk:"total_count"`
	SortResult    types.List                `tfsdk:"sort_result"`
	Targets       []TargetsForGuardrailItem `tfsdk:"targets"`
}

func (o GuardrailBindingTargetsDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"guardrail_id":    basetypes.StringType{},
		"target_type":     basetypes.StringType{},
		"landing_zone_id": basetypes.StringType{},
		"name":            basetypes.StringType{},
		"sort":            basetypes.StringType{},
		"size":            basetypes.Int64Type{},
		"page":            basetypes.Int64Type{},
		"count":           basetypes.Int64Type{},
		"sort_result": types.ListType{
			ElemType: basetypes.StringType{},
		},
		"targets": types.ListType{
			ElemType: types.ObjectType{AttrTypes: TargetsForGuardrailItem{}.AttributeTypes(ctx)},
		},
	}
}

type TargetsForGuardrailItem struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	ParentUnitId   types.String `tfsdk:"parent_unit_id"`
	ParentUnitName types.String `tfsdk:"parent_unit_name"`
}

func (o TargetsForGuardrailItem) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":               basetypes.StringType{},
		"name":             basetypes.StringType{},
		"email":            basetypes.StringType{},
		"parent_unit_id":   basetypes.StringType{},
		"parent_unit_name": basetypes.StringType{},
	}
}
