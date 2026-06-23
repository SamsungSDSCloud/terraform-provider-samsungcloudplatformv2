package database

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InstanceGroup represents a group of instances within a cluster.
type InstanceGroup struct {
	BlockStorageGroups types.List   `tfsdk:"block_storage_groups"`
	Id                 types.String `tfsdk:"id"`
	Instances          types.List   `tfsdk:"instances"`
	RoleType           types.String `tfsdk:"role_type"`
	ServerTypeName     types.String `tfsdk:"server_type_name"`
}

func (m InstanceGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"role_type":          types.StringType,
		"server_type_name":   types.StringType,
		"block_storage_groups": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":          types.StringType,
					"name":        types.StringType,
					"role_type":   types.StringType,
					"size_gb":     types.Int32Type,
					"volume_type": types.StringType,
				},
			},
		},
		"instances": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":               types.StringType,
					"role_type":          types.StringType,
					"service_ip_address": types.StringType,
					"public_ip_id":       types.StringType,
				},
			},
		},
	}
}

// BlockStorageGroup represents a block storage group within an instance group.
type BlockStorageGroup struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	RoleType   types.String `tfsdk:"role_type"`
	SizeGb     types.Int32  `tfsdk:"size_gb"`
	VolumeType types.String `tfsdk:"volume_type"`
}

func (m BlockStorageGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"role_type":   types.StringType,
		"size_gb":     types.Int32Type,
		"volume_type": types.StringType,
	}
}

// Instance represents a single instance within an instance group.
type Instance struct {
	Name             types.String `tfsdk:"name"`
	RoleType         types.String `tfsdk:"role_type"`
	ServiceIpAddress types.String `tfsdk:"service_ip_address"`
	PublicIpId       types.String `tfsdk:"public_ip_id"`
}

func (m Instance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               types.StringType,
		"role_type":          types.StringType,
		"service_ip_address": types.StringType,
		"public_ip_id":       types.StringType,
	}
}
