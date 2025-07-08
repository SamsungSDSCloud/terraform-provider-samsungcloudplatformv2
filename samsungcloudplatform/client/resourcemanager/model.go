package resourcemanager

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-resourcemanager"

// Tag

type TagDataSource struct {
	Key                types.String `tfsdk:"key"`
	Value              types.String `tfsdk:"value"`
	ResourceIdentifier types.String `tfsdk:"resource_identifier"`
	Contents           []TagContent `tfsdk:"content"`
}

type TagContent struct {
	Srn   types.String `tfsdk:"srn"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (m TagContent) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"srn":   types.StringType,
		"key":   types.StringType,
		"value": types.StringType,
	}
}

// Resource Tag

type ResourceTagDataSource struct {
	Srn        types.String `tfsdk:"srn"`
	EncodedSrn types.String `tfsdk:"encoded_srn"`
	Size       types.Int32  `tfsdk:"size"`
	Page       types.Int32  `tfsdk:"page"`
	Sort       types.String `tfsdk:"sort"`
	SrnTagList types.Object `tfsdk:"content"`
}

type SrnTagList struct {
	Srn  types.String `tfsdk:"srn"`
	Tags []Tag        `tfsdk:"tags"`
}

func (m SrnTagList) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"srn": types.StringType,
		"tags": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}},
	}
}

type Tag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// Resource Group

type ResourceGroupDataSource struct {
	Id            types.String    `tfsdk:"id"`
	Name          types.String    `tfsdk:"name"`
	Region        types.String    `tfsdk:"region"`
	Filter        []filter.Filter `tfsdk:"filter"`
	Tags          types.Map       `tfsdk:"tags"`
	ResourceGroup types.Object    `tfsdk:"resource_group"`
}

type ResourceGroupDataSourceIds struct {
	Id     types.String    `tfsdk:"id"`
	Name   types.String    `tfsdk:"name"`
	Region types.String    `tfsdk:"region"`
	Filter []filter.Filter `tfsdk:"filter"`
	Tags   types.Map       `tfsdk:"tags"`
	Ids    []types.String  `tfsdk:"ids"`
}

type ResourceGroupResource struct {
	Id                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	ResourceTypes       types.List   `tfsdk:"resource_types"`
	GroupDefinitionTags types.Map    `tfsdk:"group_definition_tags"`
	Region              types.String `tfsdk:"region"`
	Tags                types.Map    `tfsdk:"tags"`
	ResourceGroup       types.Object `tfsdk:"resource_group"`
}

type ResourceGroup struct {
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	Description         types.String `tfsdk:"description"`
	Id                  types.String `tfsdk:"id"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	Name                types.String `tfsdk:"name"`
	Region              types.String `tfsdk:"region"`
	Srn                 types.String `tfsdk:"srn"`
	ResourceTypes       types.List   `tfsdk:"resource_types"`
	Tags                types.Map    `tfsdk:"tags"`
	GroupDefinitionTags []Tag        `tfsdk:"group_definition_tags"`
}

func (m ResourceGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":     types.StringType,
		"created_by":     types.StringType,
		"description":    types.StringType,
		"id":             types.StringType,
		"modified_at":    types.StringType,
		"modified_by":    types.StringType,
		"name":           types.StringType,
		"region":         types.StringType,
		"srn":            types.StringType,
		"resource_types": types.ListType{ElemType: types.StringType},
		"tags":           types.MapType{ElemType: types.StringType},
		"group_definition_tags": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}},
	}
}
