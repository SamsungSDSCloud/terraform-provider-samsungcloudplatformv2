package networklogging

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-network-logging"

//------------------- Network Logging Storage -------------------//

type NetworkLoggingStorageDataSource struct {
	Limit                  types.Int32                     `tfsdk:"limit"`
	Marker                 types.String                    `tfsdk:"marker"`
	Sort                   types.String                    `tfsdk:"sort"`
	ResourceType           types.String                    `tfsdk:"resource_type"`
	NetworkLoggingStorages []NetworkLoggingStorageResource `tfsdk:"network_logging_storages"`
}

type NetworkLoggingStorageResource struct {
	Id           types.String `tfsdk:"id"`
	AccountId    types.String `tfsdk:"account_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	BucketName   types.String `tfsdk:"bucket_name"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
}

func (m NetworkLoggingStorageResource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"account_id":    types.StringType,
		"resource_type": types.StringType,
		"bucket_name":   types.StringType,
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"modified_at":   types.StringType,
		"modified_by":   types.StringType,
	}
}

//------------------- Network Logging Configuration -------------------//

type NetworkLoggingConfigurationDataSource struct {
	Limit                        types.Int32                   `tfsdk:"limit"`
	Marker                       types.String                  `tfsdk:"marker"`
	Sort                         types.String                  `tfsdk:"sort"`
	ResourceId                   types.String                  `tfsdk:"resource_id"`
	ResourceType                 types.String                  `tfsdk:"resource_type"`
	ResourceName                 types.String                  `tfsdk:"resource_name"`
	NetworkLoggingConfigurations []NetworkLoggingConfiguration `tfsdk:"network_logging_configurations"`
}

type NetworkLoggingConfiguration struct {
	Id                 types.String `tfsdk:"id"`
	AccountId          types.String `tfsdk:"account_id"`
	ResourceId         types.String `tfsdk:"resource_id"`
	ResourceType       types.String `tfsdk:"resource_type"`
	ResourceName       types.String `tfsdk:"resource_name"`
	BucketName         types.String `tfsdk:"bucket_name"`
	SecurityGroupLogId types.String `tfsdk:"security_group_log_id"`
	UpInterface        types.String `tfsdk:"up_interface"`
	DownInterface      types.String `tfsdk:"down_interface"`
	CreatedAt          types.String `tfsdk:"created_at"`
	CreatedBy          types.String `tfsdk:"created_by"`
	ModifiedAt         types.String `tfsdk:"modified_at"`
	ModifiedBy         types.String `tfsdk:"modified_by"`
}

func (m NetworkLoggingConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"account_id":            types.StringType,
		"resource_id":           types.StringType,
		"resource_type":         types.StringType,
		"resource_name":         types.StringType,
		"bucket_name":           types.StringType,
		"security_group_log_id": types.StringType,
		"up_interface":          types.StringType,
		"down_interface":        types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}
