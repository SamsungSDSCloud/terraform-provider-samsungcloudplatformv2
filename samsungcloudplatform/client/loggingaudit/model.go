package loggingaudit

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-loggingaudit"

//Trail

type TrailDataSource struct {
	Size         types.Int32  `tfsdk:"size"`
	Page         types.Int32  `tfsdk:"page"`
	Sort         types.String `tfsdk:"sort"`
	TrailName    types.String `tfsdk:"trail_name"`
	BucketName   types.String `tfsdk:"bucket_name"`
	State        types.String `tfsdk:"state"`
	ResourceType types.String `tfsdk:"resource_type"`
	Trails       []Trail      `tfsdk:"trail"`
}

type TagCreateRequest struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type TrailResource struct {
	Id                  types.String       `tfsdk:"id"`
	LastUpdated         types.String       `tfsdk:"last_updated"`
	AccountId           types.String       `tfsdk:"account_id"`
	BucketName          types.String       `tfsdk:"bucket_name"`
	BucketRegion        types.String       `tfsdk:"bucket_region"`
	LogTypeTotalYn      types.String       `tfsdk:"log_type_total_yn"`
	LogVerificationYn   types.String       `tfsdk:"log_verification_yn"`
	RegionNames         []types.String     `tfsdk:"region_names"`
	RegionTotalYn       types.String       `tfsdk:"region_total_yn"`
	ResourceTypeTotalYn types.String       `tfsdk:"resource_type_total_yn"`
	TagCreateRequests   []TagCreateRequest `tfsdk:"tag_create_requests"`
	TargetLogTypes      []types.String     `tfsdk:"target_log_types"`
	TargetResourceTypes []types.String     `tfsdk:"target_resource_types"`
	TargetUsers         []types.String     `tfsdk:"target_users"`
	TrailDescription    types.String       `tfsdk:"trail_description"`
	TrailName           types.String       `tfsdk:"trail_name"`
	TrailSaveType       types.String       `tfsdk:"trail_save_type"`
	UserTotalYn         types.String       `tfsdk:"user_total_yn"`
	Trail               types.Object       `tfsdk:"trail"`
}

type TrailResourceId struct {
	LogTypeTotalYn      types.String `tfsdk:"log_type_total_yn"`
	LogVerificationYn   types.String `tfsdk:"log_verification_yn"`
	RegionNames         types.List   `tfsdk:"region_names"`
	RegionTotalYn       types.String `tfsdk:"region_total_yn"`
	ResourceTypeTotalYn types.String `tfsdk:"resource_type_total_yn"`
	TargetLogTypes      types.List   `tfsdk:"target_log_types"`
	TargetResourceTypes types.List   `tfsdk:"target_resource_types"`
	TargetUsers         types.List   `tfsdk:"target_users"`
	TrailDescription    types.String `tfsdk:"trail_description"`
	TrailSaveType       types.String `tfsdk:"trail_save_type"`
	UserTotalYn         types.String `tfsdk:"user_total_yn"`
	Trail               types.Object `tfsdk:"trail"`
}

func (m Trail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":                 types.StringType,
		"account_name":               types.StringType,
		"bucket_name":                types.StringType,
		"bucket_region":              types.StringType,
		"created_at":                 types.StringType,
		"created_by":                 types.StringType,
		"created_user_id":            types.StringType,
		"del_yn":                     types.StringType,
		"id":                         types.StringType,
		"log_type_total_yn":          types.StringType,
		"log_verification_yn":        types.StringType,
		"modified_at":                types.StringType,
		"modified_by":                types.StringType,
		"region_names":               types.ListType{ElemType: types.StringType},
		"region_total_yn":            types.StringType,
		"resource_type_total_yn":     types.StringType,
		"state":                      types.StringType,
		"target_log_types":           types.ListType{ElemType: types.StringType},
		"target_resource_types":      types.ListType{ElemType: types.StringType},
		"target_users":               types.ListType{ElemType: types.StringType},
		"trail_batch_end_at":         types.StringType,
		"trail_batch_first_start_at": types.StringType,
		"trail_batch_last_state":     types.StringType,
		"trail_batch_start_at":       types.StringType,
		"trail_batch_success_at":     types.StringType,
		"trail_description":          types.StringType,
		"trail_name":                 types.StringType,
		"trail_save_type":            types.StringType,
		"user_total_yn":              types.StringType,
	}
}

type Trail struct {
	AccountId              types.String   `tfsdk:"account_id"`
	AccountName            types.String   `tfsdk:"account_name"`
	BucketName             types.String   `tfsdk:"bucket_name"`
	BucketRegion           types.String   `tfsdk:"bucket_region"`
	CreatedAt              types.String   `tfsdk:"created_at"`
	CreatedBy              types.String   `tfsdk:"created_by"`
	CreatedUserId          types.String   `tfsdk:"created_user_id"`
	DelYn                  types.String   `tfsdk:"del_yn"`
	Id                     types.String   `tfsdk:"id"`
	LogTypeTotalYn         types.String   `tfsdk:"log_type_total_yn"`
	LogVerificationYn      types.String   `tfsdk:"log_verification_yn"`
	ModifiedAt             types.String   `tfsdk:"modified_at"`
	ModifiedBy             types.String   `tfsdk:"modified_by"`
	RegionNames            []types.String `tfsdk:"region_names"`
	RegionTotalYn          types.String   `tfsdk:"region_total_yn"`
	ResourceTypeTotalYn    types.String   `tfsdk:"resource_type_total_yn"`
	State                  types.String   `tfsdk:"state"`
	TargetLogTypes         []types.String `tfsdk:"target_log_types"`
	TargetResourceTypes    []types.String `tfsdk:"target_resource_types"`
	TargetUsers            []types.String `tfsdk:"target_users"`
	TrailBatchEndAt        types.String   `tfsdk:"trail_batch_end_at"`
	TrailBatchFirstStartAt types.String   `tfsdk:"trail_batch_first_start_at"`
	TrailBatchLastState    types.String   `tfsdk:"trail_batch_last_state"`
	TrailBatchStartAt      types.String   `tfsdk:"trail_batch_start_at"`
	TrailBatchSuccessAt    types.String   `tfsdk:"trail_batch_success_at"`
	TrailDescription       types.String   `tfsdk:"trail_description"`
	TrailName              types.String   `tfsdk:"trail_name"`
	TrailSaveType          types.String   `tfsdk:"trail_save_type"`
	UserTotalYn            types.String   `tfsdk:"user_total_yn"`
}

func ConvertStringListToInterfaceList(strs []types.String) []interface{} {
	result := make([]interface{}, 0, len(strs))
	for _, s := range strs {
		if s.IsNull() || s.IsUnknown() {
			result = append(result, nil)
		} else {
			result = append(result, s.ValueString())
		}
	}
	return result
}
