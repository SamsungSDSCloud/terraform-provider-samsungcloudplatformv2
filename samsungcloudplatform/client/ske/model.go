package ske

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-ske"

//------------ Cluster -------------------//

type ClusterDataSourceIds struct {
	Size              types.Int32     `tfsdk:"size"`
	Page              *int32          `tfsdk:"page"`
	Sort              types.String    `tfsdk:"sort"`
	Name              types.String    `tfsdk:"name"`
	Status            []types.String  `tfsdk:"status"`
	KubernetesVersion []types.String  `tfsdk:"kubernetes_version"`
	SubnetId          types.String    `tfsdk:"subnet_id"`
	Region            types.String    `tfsdk:"region"` // region field 를 추가한다.
	Filter            []filter.Filter `tfsdk:"filter"` // filter field 를 추가한다.
	Tags              types.Map       `tfsdk:"tags"`   // tags  field 를 추가한다.
	Ids               []types.String  `tfsdk:"ids"`
}

type ClusterDataSource struct {
	Id      types.String `tfsdk:"id"`
	Cluster types.Object `tfsdk:"cluster"`
}

type Cluster struct {
	Id                                    types.String                           `tfsdk:"id"`
	Name                                  types.String                           `tfsdk:"name"`
	AccountId                             types.String                           `tfsdk:"account_id"`
	CloudLoggingEnabled                   types.Bool                             `tfsdk:"cloud_logging_enabled"`
	KubernetesVersion                     types.String                           `tfsdk:"kubernetes_version"`
	ClusterNamespace                      types.String                           `tfsdk:"cluster_namespace"`
	MaxNodeCount                          types.Int32                            `tfsdk:"max_node_count"`
	NodeCount                             types.Int32                            `tfsdk:"node_count"`
	PrivateEndpointUrl                    types.String                           `tfsdk:"private_endpoint_url"`
	PrivateKubeconfigDownloadYn           types.String                           `tfsdk:"private_kubeconfig_download_yn"`
	PrivateEndpointAccessControlResources []PrivateEndpointAccessControlResource `tfsdk:"private_endpoint_access_control_resources"`
	PublicEndpointUrl                     types.String                           `tfsdk:"public_endpoint_url"`
	PublicKubeconfigDownloadYn            types.String                           `tfsdk:"public_kubeconfig_download_yn"`
	PublicEndpointAccessControlIp         types.String                           `tfsdk:"public_endpoint_access_control_ip"`
	Vpc                                   ExternalResource                       `tfsdk:"vpc"`
	Subnet                                ExternalResource                       `tfsdk:"subnet"`
	Volume                                ExternalResource                       `tfsdk:"volume"`
	SecurityGroupList                     []ExternalResource                     `tfsdk:"security_group_list"`
	ManagedSecurityGroup                  ExternalResource                       `tfsdk:"managed_security_group"`
	CreatedAt                             types.String                           `tfsdk:"created_at"`
	CreatedBy                             types.String                           `tfsdk:"created_by"`
	ModifiedAt                            types.String                           `tfsdk:"modified_at"`
	ModifiedBy                            types.String                           `tfsdk:"modified_by"`
	Status                                types.String                           `tfsdk:"status"`
	ServiceWatchLoggingEnabled            types.Bool                             `tfsdk:"service_watch_logging_enabled"` //v1.1
}

func (m Cluster) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                             types.StringType,
		"name":                           types.StringType,
		"account_id":                     types.StringType,
		"cloud_logging_enabled":          types.BoolType,
		"kubernetes_version":             types.StringType,
		"cluster_namespace":              types.StringType,
		"max_node_count":                 types.Int32Type,
		"node_count":                     types.Int32Type,
		"private_endpoint_url":           types.StringType,
		"private_kubeconfig_download_yn": types.StringType,
		"private_endpoint_access_control_resources": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
		}},
		"public_endpoint_url":               types.StringType,
		"public_kubeconfig_download_yn":     types.StringType,
		"public_endpoint_access_control_ip": types.StringType,
		"vpc": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
			},
		},
		"subnet": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
			},
		},
		"volume": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
			},
		},
		"security_group_list": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
			},
		}},
		"managed_security_group": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
			},
		},
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"status":      types.StringType,
		// v1.1
		"service_watch_logging_enabled": types.BoolType,
	}
}

type ClusterResource struct {
	Id                                    types.String                           `tfsdk:"id"`
	LastUpdated                           types.String                           `tfsdk:"last_updated"`
	Name                                  types.String                           `tfsdk:"name"`
	CloudLoggingEnabled                   types.Bool                             `tfsdk:"cloud_logging_enabled"`
	KubernetesVersion                     types.String                           `tfsdk:"kubernetes_version"`
	PrivateEndpointAccessControlResources []PrivateEndpointAccessControlResource `tfsdk:"private_endpoint_access_control_resources"`
	PublicEndpointAccessControlIp         types.String                           `tfsdk:"public_endpoint_access_control_ip"`
	SecurityGroupIdList                   []types.String                         `tfsdk:"security_group_id_list"`
	SubnetId                              types.String                           `tfsdk:"subnet_id"`
	VolumeId                              types.String                           `tfsdk:"volume_id"`
	VpcId                                 types.String                           `tfsdk:"vpc_id"`
	ServiceWatchLoggingEnabled            types.Bool                             `tfsdk:"service_watch_logging_enabled"` //v1.1
	Region                                types.String                           `tfsdk:"region"`                        // region field 를 추가한다.
	Tags                                  types.Map                              `tfsdk:"tags"`                          // tags field 필드를 추가한다.
	Cluster                               types.Object                           `tfsdk:"cluster"`
}

type ClusterKubeconfigDataSource struct {
	ClusterId      types.String `tfsdk:"cluster_id"`
	KubeconfigType types.String `tfsdk:"kubeconfig_type"`
	Kubeconfig     types.String `tfsdk:"kubeconfig"`
}

type KubernetesVersionDataSources struct {
	Region             types.String               `tfsdk:"region"` // region field 를 추가한다.
	KubernetesVersions []KubernetesVersionSummary `tfsdk:"kubernetes_versions"`
}

type KubernetesVersion struct {
	Description       types.String `tfsdk:"description"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
}

type KubernetesVersionSummary struct {
	Description       types.String `tfsdk:"description"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
}

//------------ Nodepool -------------------//

type NodepoolDataSource struct {
	Id       types.String `tfsdk:"id"`
	Nodepool types.Object `tfsdk:"nodepool"`
}

type NodepoolDataSources struct {
	ClusterId types.String `tfsdk:"cluster_id"`
	Nodepools []Nodepool   `tfsdk:"nodepools"`
}

type NodepoolnodeDataSources struct {
	NodepoolId types.String     `tfsdk:"nodepool_id"`
	Nodes      []NodeInNodepool `tfsdk:"nodes"`
}

type Nodepool struct {
	Id                  types.String      `tfsdk:"id"`
	Name                types.String      `tfsdk:"name"`
	AccountId           types.String      `tfsdk:"account_id"`
	AutoRecoveryEnabled types.Bool        `tfsdk:"auto_recovery_enabled"`
	AutoScaleEnabled    types.Bool        `tfsdk:"auto_scale_enabled"`
	CurrentNodeCount    types.Int32       `tfsdk:"current_node_count"`
	DesiredNodeCount    types.Int32       `tfsdk:"desired_node_count"`
	Image               Image             `tfsdk:"image"`
	KubernetesVersion   types.String      `tfsdk:"kubernetes_version"`
	ServerType          ServerType        `tfsdk:"server_type"`
	Status              types.String      `tfsdk:"status"`
	VolumeType          VolumeTypeSummary `tfsdk:"volume_type"`
}

type NodeInNodepool struct {
	Name              types.String `tfsdk:"name"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Status            types.String `tfsdk:"status"`
}

type NodepoolResource struct {
	Id                types.String      `tfsdk:"id"`
	LastUpdated       types.String      `tfsdk:"last_updated"`
	Name              types.String      `tfsdk:"name"`
	ClusterId         types.String      `tfsdk:"cluster_id"`
	CustomImageId     types.String      `tfsdk:"custom_image_id"`
	DesiredNodeCount  types.Int32       `tfsdk:"desired_node_count"`
	ImageOs           types.String      `tfsdk:"image_os"`
	ImageOsVersion    types.String      `tfsdk:"image_os_version"`
	Labels            []Labels          `tfsdk:"labels"`
	Taints            []Taints          `tfsdk:"taints"`
	IsAutoRecovery    types.Bool        `tfsdk:"is_auto_recovery"`
	IsAutoScale       types.Bool        `tfsdk:"is_auto_scale"`
	KeypairName       types.String      `tfsdk:"keypair_name"`
	KubernetesVersion types.String      `tfsdk:"kubernetes_version"`
	MaxNodeCount      types.Int32       `tfsdk:"max_node_count"`
	MinNodeCount      types.Int32       `tfsdk:"min_node_count"`
	ServerTypeId      types.String      `tfsdk:"server_type_id"`
	VolumeTypeName    types.String      `tfsdk:"volume_type_name"`
	VolumeSize        types.Int32       `tfsdk:"volume_size"`
	ServerGroupId     types.String      `tfsdk:"server_group_id"`   // v1.1
	AdvancedSettings  *AdvancedSettings `tfsdk:"advanced_settings"` // v1.1
	NodepoolDetail    types.Object      `tfsdk:"nodepool_detail"`
}

type NodepoolDetail struct {
	Id                  types.String      `tfsdk:"id"`
	Name                types.String      `tfsdk:"name"`
	AccountId           types.String      `tfsdk:"account_id"`
	AutoRecoveryEnabled types.Bool        `tfsdk:"auto_recovery_enabled"`
	AutoScaleEnabled    types.Bool        `tfsdk:"auto_scale_enabled"`
	Cluster             IdMapType         `tfsdk:"cluster"`
	CurrentNodeCount    types.Int32       `tfsdk:"current_node_count"`
	DesiredNodeCount    types.Int32       `tfsdk:"desired_node_count"`
	Image               Image             `tfsdk:"image"`
	Keypair             NameMapType       `tfsdk:"keypair"`
	KubernetesVersion   types.String      `tfsdk:"kubernetes_version"`
	Labels              []Labels          `tfsdk:"labels"`
	Taints              []Taints          `tfsdk:"taints"`
	MaxNodeCount        types.Int32       `tfsdk:"max_node_count"`
	MinNodeCount        types.Int32       `tfsdk:"min_node_count"`
	ServerType          ServerType        `tfsdk:"server_type"`
	Status              types.String      `tfsdk:"status"`
	VolumeType          VolumeType        `tfsdk:"volume_type"`
	VolumeSize          types.Int32       `tfsdk:"volume_size"`
	CreatedAt           types.String      `tfsdk:"created_at"`
	CreatedBy           types.String      `tfsdk:"created_by"`
	ModifiedAt          types.String      `tfsdk:"modified_at"`
	ModifiedBy          types.String      `tfsdk:"modified_by"`
	ServerGroupId       types.String      `tfsdk:"server_group_id"`   // v1.1
	AdvancedSettings    *AdvancedSettings `tfsdk:"advanced_settings"` // v1.1
}

func (m NodepoolDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"account_id":            types.StringType,
		"auto_recovery_enabled": types.BoolType,
		"auto_scale_enabled":    types.BoolType,
		"cluster": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		},
		"current_node_count": types.Int32Type,
		"desired_node_count": types.Int32Type,
		"image": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"custom_image_name": types.StringType,
				"os":                types.StringType,
				"os_version":        types.StringType,
			},
		},
		"keypair": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
			},
		},
		"kubernetes_version": types.StringType,
		"labels": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}},
		"taints": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"effect": types.StringType,
				"key":    types.StringType,
				"value":  types.StringType,
			},
		}},
		"max_node_count": types.Int32Type,
		"min_node_count": types.Int32Type,
		"server_type": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"description": types.StringType,
				"id":          types.StringType,
			},
		},
		"status": types.StringType,
		"volume_type": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"encrypt": types.BoolType,
				"id":      types.StringType,
				"name":    types.StringType,
			},
		},
		"volume_size": types.Int32Type,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		//v1.1
		"server_group_id": types.StringType,
		"advanced_settings": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"allowed_unsafe_sysctls":  types.StringType,
				"container_log_max_files": types.Int32Type,
				"container_log_max_size":  types.Int32Type,
				"image_gc_high_threshold": types.Int32Type,
				"image_gc_low_threshold":  types.Int32Type,
				"max_pods":                types.Int32Type,
				"pod_max_pids":            types.Int32Type,
			},
		},
	}
}
