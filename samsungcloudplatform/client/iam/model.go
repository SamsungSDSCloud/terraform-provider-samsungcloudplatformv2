package iam

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-iam"

type AccessKeyDataSource struct {
	Limit      types.Int32  `tfsdk:"limit"`
	Marker     types.String `tfsdk:"marker"`
	Sort       types.String `tfsdk:"sort"`
	AccountId  types.String `tfsdk:"account_id"`
	AccessKeys []AccessKey  `tfsdk:"access_keys"`
}

type AccessKeyResource struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	AccessKeyType     types.String `tfsdk:"access_key_type"`
	AccountId         types.String `tfsdk:"account_id"`
	Description       types.String `tfsdk:"description"`
	Duration          types.String `tfsdk:"duration"`
	ParentAccessKeyId types.String `tfsdk:"parent_access_key_id"`
	Passcode          types.String `tfsdk:"passcode"`
	AccessKey         types.Object `tfsdk:"access_key"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
}

type AccessKey struct {
	AccessKey           types.String `tfsdk:"access_key"`
	AccessKeyType       types.String `tfsdk:"access_key_type"`
	AccountId           types.String `tfsdk:"account_id"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	Description         types.String `tfsdk:"description"`
	ExpirationTimestamp types.String `tfsdk:"expiration_timestamp"`
	Id                  types.String `tfsdk:"id"`
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	ParentAccessKeyId   types.String `tfsdk:"parent_access_key_id"`
	SecretKey           types.String `tfsdk:"secret_key"`
}

func (m AccessKey) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_key":           types.StringType,
		"access_key_type":      types.StringType,
		"account_id":           types.StringType,
		"created_at":           types.StringType,
		"created_by":           types.StringType,
		"description":          types.StringType,
		"expiration_timestamp": types.StringType,
		"id":                   types.StringType,
		"is_enabled":           types.BoolType,
		"modified_at":          types.StringType,
		"modified_by":          types.StringType,
		"parent_access_key_id": types.StringType,
		"secret_key":           types.StringType,
	}
}
