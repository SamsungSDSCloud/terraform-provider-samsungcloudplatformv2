package ske

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-ske"

//------------ Cluster -------------------//

type ClusterDataSourceIds struct {
	Size              *int32          `tfsdk:"size"`
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
	Id                        types.String `tfsdk:"id"`
	Cluster                   types.Object `tfsdk:"cluster"`
	DeletionProtectionEnabled types.Bool   `tfsdk:"deletion_protection_enabled"`
}

type Cluster struct {
	Id                                    types.String                           `tfsdk:"id"`
	Name                                  types.String                           `tfsdk:"name"`
	AccountId                             types.String                           `tfsdk:"account_id"`
	AdditionalSubnetIdList                types.List                             `tfsdk:"additional_subnet_id_list"`
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
	Vpc                                   ExternalResourceId                       `tfsdk:"vpc"`
	Subnet                                ExternalResourceId                       `tfsdk:"subnet"`
	Volume                                ExternalResourceId                       `tfsdk:"volume"`
	SecurityGroupList                     []ExternalResourceId                     `tfsdk:"security_group_list"`
	ManagedSecurityGroup                  ExternalResourceId                       `tfsdk:"managed_security_group"`
	CreatedAt                             types.String                           `tfsdk:"created_at"`
	CreatedBy                             types.String                           `tfsdk:"created_by"`
	ModifiedAt                            types.String                           `tfsdk:"modified_at"`
	ModifiedBy                            types.String                           `tfsdk:"modified_by"`
	Status                                types.String                           `tfsdk:"status"`
	ServiceWatchLoggingEnabled            types.Bool                             `tfsdk:"service_watch_logging_enabled"`
	LinkedResources                       []LinkedResource                       `tfsdk:"linked_resources"`
}

func (m Cluster) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                             types.StringType,
		"name":                           types.StringType,
		"account_id":                     types.StringType,
		"additional_subnet_id_list":      types.ListType{ElemType: types.StringType},
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
				"id": types.StringType,
			},
		},
		"subnet": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		},
		"volume": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		},
		"security_group_list": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		}},
		"managed_security_group": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		},
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"status":      types.StringType,
		// v1.1
		"service_watch_logging_enabled": types.BoolType,
		// v1.6
		"linked_resources": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
		}},
	}
}

type ClusterResource struct {
	Id                                    types.String `tfsdk:"id"`
	Name                                  types.String `tfsdk:"name"`
	KubernetesVersion                     types.String `tfsdk:"kubernetes_version"`
	PrivateEndpointAccessControlResources types.List   `tfsdk:"private_endpoint_access_control_resources"`
	PublicEndpointAccessControlIp         types.String `tfsdk:"public_endpoint_access_control_ip"`
	SecurityGroupIdList                   types.List   `tfsdk:"security_group_id_list"`
	DefaultSubnetId                       types.String `tfsdk:"default_subnet_id"`
	NfsVolumeId                           types.String `tfsdk:"nfs_volume_id"`
	VpcId                                 types.String `tfsdk:"vpc_id"`
	ServiceWatchLoggingEnabled            types.Bool   `tfsdk:"service_watch_logging_enabled"` //v1.1
	Tags                                  types.Map    `tfsdk:"tags"`                          // tags field 필드를 추가한다.
	Cluster                               types.Object `tfsdk:"cluster"`
	AdditionalSubnetIdList                types.List   `tfsdk:"additional_subnet_id_list"`   // v1.6
	DeletionProtectionEnabled             types.Bool   `tfsdk:"deletion_protection_enabled"` // v1.6
	LinkedResources                       types.List   `tfsdk:"linked_resources"`            // v1.6
	//LastUpdated                           types.String                           `tfsdk:"last_updated"`
	//Region                                types.String                           `tfsdk:"region"` // region field 를 추가한다.
}

type ClusterKubeconfigDataSource struct {
	ClusterId      types.String `tfsdk:"cluster_id"`
	KubeconfigType types.String `tfsdk:"kubeconfig_type"`
	Kubeconfig     types.String `tfsdk:"kubeconfig"`
}

//------------ Cluster Subnets -------------------//

type ClusterSubnetsResource struct {
	Id                     types.String `tfsdk:"id"`
	ClusterId              types.String `tfsdk:"cluster_id"`
	AdditionalSubnetIdList types.Set    `tfsdk:"additional_subnet_id_list"`
}

type ClusterUserKubeconfigDataSource struct {
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
	ClusterId types.String      `tfsdk:"cluster_id"`
	Nodepools []NodepoolSummary `tfsdk:"nodepools"`
}

type NodepoolnodeDataSources struct {
	NodepoolId types.String     `tfsdk:"nodepool_id"`
	Nodes      []NodeInNodepool `tfsdk:"nodes"`
}

type NodepoolSummary struct {
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
	SubnetId            types.String      `tfsdk:"subnet_id"`
	VolumeType          VolumeTypeSummary `tfsdk:"volume_type"`
}

type NodeInNodepool struct {
	Name              types.String `tfsdk:"name"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Status            types.String `tfsdk:"status"`
}

type NodepoolResource struct {
	Id types.String `tfsdk:"id"`
	//LastUpdated       types.String      `tfsdk:"last_updated"`
	Name                types.String      `tfsdk:"name"`
	ClusterId           types.String      `tfsdk:"cluster_id"`
	CustomImageId       types.String      `tfsdk:"custom_image_id"`
	DesiredNodeCount    types.Int32       `tfsdk:"desired_node_count"`
	ImageOs             types.String      `tfsdk:"image_os"`
	ImageOsVersion      types.String      `tfsdk:"image_os_version"`
	Labels              []Label           `tfsdk:"labels"`
	Taints              []Taint           `tfsdk:"taints"`
	IsAutoRecovery      types.Bool        `tfsdk:"is_auto_recovery"`
	IsAutoScale         types.Bool        `tfsdk:"is_auto_scale"`
	KeypairName         types.String      `tfsdk:"keypair_name"`
	KubernetesVersion   types.String      `tfsdk:"kubernetes_version"`
	MaxNodeCount        types.Int32       `tfsdk:"max_node_count"`
	MinNodeCount        types.Int32       `tfsdk:"min_node_count"`
	ServerTypeId        types.String      `tfsdk:"server_type_id"`
	VolumeTypeName      types.String      `tfsdk:"volume_type_name"`
	VolumeSize          types.Int32       `tfsdk:"volume_size"`
	ServerGroupId       types.String      `tfsdk:"server_group_id"`       // v1.1
	AdvancedSettings    *AdvancedSettings `tfsdk:"advanced_settings"`     // v1.1
	LinkedResources     []LinkedResource  `tfsdk:"linked_resources"`      // v1.3
	VolumeMaxIops       types.Int32       `tfsdk:"volume_max_iops"`       // v1.4
	VolumeMaxThroughput types.Int32       `tfsdk:"volume_max_throughput"` // v1.4
	ScpGpuDriver        types.String      `tfsdk:"scp_gpu_driver"`        // v1.4
	PreferredIps        types.String      `tfsdk:"preferred_ips"`         // v1.5
	SubnetId            types.String      `tfsdk:"subnet_id"`             // v1.5
	Zone                types.String      `tfsdk:"zone"`                  // v1.5
	Nodepool            types.Object      `tfsdk:"nodepool"`
}

type Nodepool struct {
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
	Labels              []Label           `tfsdk:"labels"`
	Taints              []Taint           `tfsdk:"taints"`
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
	ServerGroupId       types.String      `tfsdk:"server_group_id"`       // v1.1
	AdvancedSettings    *AdvancedSettings `tfsdk:"advanced_settings"`     // v1.1
	LinkedResources     []LinkedResource  `tfsdk:"linked_resources"`      // v1.3
	VolumeMaxIops       *int32            `tfsdk:"volume_max_iops"`       // v1.4
	VolumeMaxThroughput *int32            `tfsdk:"volume_max_throughput"` // v1.4
	PreferredIps        types.String      `tfsdk:"preferred_ips"`         // v1.5
	SubnetId            types.String      `tfsdk:"subnet_id"`             // v1.5
	Zone                types.String      `tfsdk:"zone"`                  // v1.5
}

func (m Nodepool) AttributeTypes() map[string]attr.Type {
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
				"scp_gpu_driver":    types.StringType, // v1.4
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
		// v1.1
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
		// v1.3
		"linked_resources": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
		}},
		// v1.4
		"volume_max_iops":       types.Int32Type,
		"volume_max_throughput": types.Int32Type,
		"preferred_ips":         types.StringType,
		"subnet_id":             types.StringType,
		"zone":                  types.StringType,
	}
}

//------------ NodepoolImage -------------------//

type NodepoolImageDataSources struct {
	Size                 types.Int32            `tfsdk:"size"`
	Page                 types.Int32            `tfsdk:"page"`
	Sort                 types.String           `tfsdk:"sort"`
	KubernetesVersion    types.String           `tfsdk:"kubernetes_version"`
	Os                   types.String           `tfsdk:"os"`
	ScpOriginalImageType types.String           `tfsdk:"scp_original_image_type"`
	NodepoolImages       []NodepoolImageSummary `tfsdk:"nodepool_images"`
}

type NodepoolImageSummary struct {
	Id                     types.String         `tfsdk:"id"`
	Name                   types.String         `tfsdk:"name"`
	Os                     types.String         `tfsdk:"os"`
	OsVersion              types.String         `tfsdk:"os_version"`
	KubernetesVersion      types.String         `tfsdk:"kubernetes_version"`
	EndOfSupport           types.Bool           `tfsdk:"end_of_support"`
	ScpImageType           types.String         `tfsdk:"scp_image_type"`
	ScpOriginalImageType   types.String         `tfsdk:"scp_original_image_type"`
	Volume                 *NodepoolImageVolume `tfsdk:"volume"`
	ScpGpuDriver           types.String         `tfsdk:"scp_gpu_driver"`
	ScpSupportedClassTypes []types.String       `tfsdk:"scp_supported_class_types"`
	Visibility             types.String         `tfsdk:"visibility"`
	Zone                   types.String         `tfsdk:"zone"`
}

type NodepoolImageVolume struct {
	VolumeSize types.Int64 `tfsdk:"volume_size"`
}
