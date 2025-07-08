package billing

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-billingplan"


type PlannedCompute struct {
	AccountId             types.String          `tfsdk:"account_id"`
	ContractId            types.String          `tfsdk:"contract_id"`
	ContractType          types.String          `tfsdk:"contract_type"`
	CreatedAt             types.String          `tfsdk:"created_at"`
	CreatedBy             types.String          `tfsdk:"created_by"`
	DeleteYn              types.String          `tfsdk:"delete_yn"`
	EndDate               types.String          `tfsdk:"end_date"`
	FirstContractStartAt  types.String          `tfsdk:"first_contract_start_at"`
	Id                    types.String          `tfsdk:"id"`
	ModifiedAt            types.String          `tfsdk:"modified_at"`
	ModifiedBy            types.String          `tfsdk:"modified_by"`
	NextContractType      types.String          `tfsdk:"next_contract_type"`
	NextEndDate           types.String          `tfsdk:"next_end_date"`
	NextStartDate         types.String          `tfsdk:"next_start_date"`
	OsName                types.String          `tfsdk:"os_name"`
	OsType                types.String          `tfsdk:"os_type"`
	Region                types.String          `tfsdk:"region"`
	ResourceName          types.String          `tfsdk:"resource_name"`
	ResourceType          types.String          `tfsdk:"resource_type"`
	ServerType            types.String          `tfsdk:"server_type"`
	ServerTypeDescription types.Map				`tfsdk:"server_type_description"`
	ServiceId             types.String          `tfsdk:"service_id"`
	ServiceName           types.String          `tfsdk:"service_name"`
	Srn                   types.String          `tfsdk:"srn"`
	StartDate             types.String          `tfsdk:"start_date"`
	State                 types.String          `tfsdk:"state"`
}

type PlannedComputeDataSourceIds struct {
	Limit            types.Int32      `tfsdk:"limit"`
	Page             types.Int32      `tfsdk:"page"`
	StartDate        types.String     `tfsdk:"start_date"`
	EndDate          types.String     `tfsdk:"end_date"`
	ServerType       types.String     `tfsdk:"server_type"`
	ContractId       types.String     `tfsdk:"contract_id"`
	ContractType     []types.String   `tfsdk:"contract_type"`
	NextContractType []types.String   `tfsdk:"next_contract_type"`
	ServiceId        []types.String   `tfsdk:"service_id"`
	OsType           []types.String   `tfsdk:"os_type"`
	State            []types.String   `tfsdk:"state"`
	CreatedBy        types.String     `tfsdk:"created_by"`
	ModifiedBy       types.String     `tfsdk:"modified_by"`
	Sort             types.String     `tfsdk:"sort"`
	PlannedComputes  []PlannedCompute `tfsdk:"planned_computes"`
}

type PlannedComputeResource struct {
	Id             types.String `tfsdk:"id"`
	LastUpdated    types.String `tfsdk:"last_updated"`
	AccountId      types.String `tfsdk:"account_id"`
	ContractType   types.String `tfsdk:"contract_type"`
	OsType         types.String `tfsdk:"os_type"`
	ServerType     types.String `tfsdk:"server_type"`
	ServiceId      types.String `tfsdk:"service_id"`
	ServiceName    types.String `tfsdk:"service_name"`
	Action		   types.String `tfsdk:"action"`
	Tags           types.Map    `tfsdk:"tags"`
	Region         types.String `tfsdk:"region"`
	PlannedCompute types.Object `tfsdk:"planned_compute"`
}

func (m PlannedCompute) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":              types.StringType,
		"contract_id":             types.StringType,
		"contract_type":           types.StringType,
		"created_at":              types.StringType,
		"created_by":              types.StringType,
		"delete_yn":               types.StringType,
		"end_date":                types.StringType,
		"first_contract_start_at": types.StringType,
		"id":                      types.StringType,
		"modified_at":             types.StringType,
		"modified_by":             types.StringType,
		"next_contract_type":      types.StringType,
		"next_end_date":           types.StringType,
		"next_start_date":         types.StringType,
		"os_name":                 types.StringType,
		"os_type":                 types.StringType,
		"region":                  types.StringType,
		"resource_name":           types.StringType,
		"resource_type":           types.StringType,
		"server_type":             types.StringType,
		"server_type_description": types.MapType{ElemType: types.StringType},
		"service_id":              types.StringType,
		"service_name":            types.StringType,
		"srn":                     types.StringType,
		"start_date":              types.StringType,
		"state":                   types.StringType,
	}
}
