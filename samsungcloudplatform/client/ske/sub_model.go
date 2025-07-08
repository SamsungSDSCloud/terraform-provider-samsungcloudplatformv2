package ske

import "github.com/hashicorp/terraform-plugin-framework/types"

type Image struct {
	CustomImageName types.String `tfsdk:"custom_image_name"`
	Os              types.String `tfsdk:"os"`
	OsVersion       types.String `tfsdk:"os_version"`
}

type ServerType struct {
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

type VolumeTypeSummary struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type VolumeType struct {
	Encrypt types.Bool   `tfsdk:"encrypt"`
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
}

type ExternalResource struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type PrivateEndpointAccessControlResource struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type IdMapType struct {
	Id types.String `tfsdk:"id"`
}

type NameMapType struct {
	Name types.String `tfsdk:"name"`
}

type Labels struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type Taints struct {
	Effect types.String `tfsdk:"effect"`
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
}
