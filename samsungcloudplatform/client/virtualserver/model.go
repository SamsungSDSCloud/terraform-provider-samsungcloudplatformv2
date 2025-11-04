package virtualserver

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-virtualserver"

// -------------------- Volume -------------------- //

type VolumeDataSourceIds struct {
	State    types.String    `tfsdk:"state"`
	Name     types.String    `tfsdk:"name"`
	Bootable types.Bool      `tfsdk:"bootable"`
	Filter   []filter.Filter `tfsdk:"filter"`
	Ids      []types.String  `tfsdk:"ids"`
}
type VolumeDataSource struct {
	Id       types.String    `tfsdk:"id"`
	State    types.String    `tfsdk:"state"`
	Name     types.String    `tfsdk:"name"`
	Bootable types.Bool      `tfsdk:"bootable"`
	Filter   []filter.Filter `tfsdk:"filter"`
	Volume   types.Object    `tfsdk:"volume"`
}

type VolumeResource struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Size        types.Int32    `tfsdk:"size"`
	State       types.String   `tfsdk:"state"`
	UserId      types.String   `tfsdk:"user_id"`
	VolumeType  types.String   `tfsdk:"volume_type"`
	Encrypted   types.Bool     `tfsdk:"encrypted"`
	Bootable    types.Bool     `tfsdk:"bootable"`
	Multiattach types.Bool     `tfsdk:"multiattach"`
	Servers     []VolumeServer `tfsdk:"servers"`
	Tags        types.Map      `tfsdk:"tags"`
}
type Volume struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Size        types.Int32    `tfsdk:"size"`
	State       types.String   `tfsdk:"state"`
	UserId      types.String   `tfsdk:"user_id"`
	VolumeType  types.String   `tfsdk:"volume_type"`
	Encrypted   types.Bool     `tfsdk:"encrypted"`
	Bootable    types.Bool     `tfsdk:"bootable"`
	Multiattach types.Bool     `tfsdk:"multiattach"`
	Servers     []VolumeServer `tfsdk:"servers"`
}

type VolumeServer struct {
	Id types.String `tfsdk:"id"`
}

func (m Volume) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"state":       types.StringType,
		"volume_type": types.StringType,
		"user_id":     types.StringType,
		"encrypted":   types.BoolType,
		"bootable":    types.BoolType,
		"multiattach": types.BoolType,
		"size":        types.Int32Type,
		"servers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id": types.StringType,
				},
			},
		},
	}
}

// -------------------- Keypair -------------------- //

type KeypairDataSourceNames struct {
	Filter []filter.Filter `tfsdk:"filter"`
	Names  []types.String  `tfsdk:"names"`
}

type KeypairDataSource struct {
	Name    types.String    `tfsdk:"name"`
	Filter  []filter.Filter `tfsdk:"filter"`
	Keypair types.Object    `tfsdk:"keypair"`
}

type Keypair struct {
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Type        types.String `tfsdk:"type"`
}

type KeypairResource struct {
	Id          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Type        types.String `tfsdk:"type"`
	PrivateKey  types.String `tfsdk:"private_key"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UserId      types.String `tfsdk:"user_id"`
	Tags        types.Map    `tfsdk:"tags"`
}

func (m Keypair) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"public_key":  types.StringType,
		"fingerprint": types.StringType,
		"type":        types.StringType,
	}
}

// -------------------- Server -------------------- //

type ServerDataSourceIds struct {
	Name               types.String    `tfsdk:"name"`
	Ip                 types.String    `tfsdk:"ip"`
	State              types.String    `tfsdk:"state"`
	ProductCategory    types.String    `tfsdk:"product_category"`
	ProductOffering    types.String    `tfsdk:"product_offering"`
	VpcId              types.String    `tfsdk:"vpc_id"`
	ServerTypeId       types.String    `tfsdk:"server_type_id"`
	AutoScalingGroupId types.String    `tfsdk:"auto_scaling_group_id"`
	Filter             []filter.Filter `tfsdk:"filter"`
	Ids                []types.String  `tfsdk:"ids"`
}

type ServerDataSource struct {
	Id                 types.String    `tfsdk:"id"`
	Name               types.String    `tfsdk:"name"`
	Ip                 types.String    `tfsdk:"ip"`
	State              types.String    `tfsdk:"state"`
	ProductCategory    types.String    `tfsdk:"product_category"`
	ProductOffering    types.String    `tfsdk:"product_offering"`
	VpcId              types.String    `tfsdk:"vpc_id"`
	ServerTypeId       types.String    `tfsdk:"server_type_id"`
	AutoScalingGroupId types.String    `tfsdk:"auto_scaling_group_id"`
	Filter             []filter.Filter `tfsdk:"filter"`
	Server             types.Object    `tfsdk:"server"`
}

type Server struct {
	AccountId             types.String    `tfsdk:"account_id"`
	Addresses             []ServerAddress `tfsdk:"addresses"`
	AutoScalingGroupId    types.String    `tfsdk:"auto_scaling_group_id"`
	CreatedAt             types.String    `tfsdk:"created_at"`
	CreatedBy             types.String    `tfsdk:"created_by"`
	DiskConfig            types.String    `tfsdk:"disk_config"`
	Id                    types.String    `tfsdk:"id"`
	ImageId               types.String    `tfsdk:"image_id"`
	KeypairName           types.String    `tfsdk:"keypair_name"`
	LaunchConfigurationId types.String    `tfsdk:"launch_configuration_id"`
	Locked                types.Bool      `tfsdk:"locked"`
	Metadata              types.Map       `tfsdk:"metadata"`
	ModifiedAt            types.String    `tfsdk:"modified_at"`
	Name                  types.String    `tfsdk:"name"`
	PlannedComputeOsType  types.String    `tfsdk:"planned_compute_os_type"`
	ProductCategory       types.String    `tfsdk:"product_category"`
	ProductOffering       types.String    `tfsdk:"product_offering"`
	SecurityGroups        []SecurityGroup `tfsdk:"security_groups"`
	ServerGroupId         types.String    `tfsdk:"server_group_id"`
	ServerType            ServerType      `tfsdk:"server_type"`
	State                 types.String    `tfsdk:"state"`
	Volumes               []ServerVolume  `tfsdk:"volumes"`
	VpcId                 types.String    `tfsdk:"vpc_id"`
}

type ServerAddress struct {
	IpAddresses []ServerIpAddress `tfsdk:"ip_addresses"`
	SubnetName  types.String      `tfsdk:"subnet_name"`
}

type ServerIpAddress struct {
	IpAddress types.String `tfsdk:"ip_address"`
	Version   types.Int32  `tfsdk:"version"`
}

type SecurityGroup struct {
	Name types.String `tfsdk:"name"`
}

type ServerType struct {
	Disk       types.Int32  `tfsdk:"disk"`
	Ephemeral  types.Int32  `tfsdk:"ephemeral"`
	ExtraSpecs types.Map    `tfsdk:"extra_specs"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Ram        types.Int32  `tfsdk:"ram"`
	Swap       types.Int32  `tfsdk:"swap"`
	Vcpus      types.Int32  `tfsdk:"vcpus"`
}

type ServerVolume struct {
	DeleteOnTermination types.Bool   `tfsdk:"delete_on_termination"`
	Id                  types.String `tfsdk:"id"`
}

type ServerResource struct {
	Id                    types.String         `tfsdk:"id"`
	AccountId             types.String         `tfsdk:"account_id"`
	Networks              types.Map            `tfsdk:"networks"`
	AutoScalingGroupId    types.String         `tfsdk:"auto_scaling_group_id"`
	CreatedAt             types.String         `tfsdk:"created_at"`
	CreatedBy             types.String         `tfsdk:"created_by"`
	DiskConfig            types.String         `tfsdk:"disk_config"`
	ImageId               types.String         `tfsdk:"image_id"`
	KeypairName           types.String         `tfsdk:"keypair_name"`
	LaunchConfigurationId types.String         `tfsdk:"launch_configuration_id"`
	Lock                  types.Bool           `tfsdk:"lock"`
	Metadata              types.Map            `tfsdk:"metadata"`
	ModifiedAt            types.String         `tfsdk:"modified_at"`
	Name                  types.String         `tfsdk:"name"`
	PlannedComputeOsType  types.String         `tfsdk:"planned_compute_os_type"`
	ProductCategory       types.String         `tfsdk:"product_category"`
	ProductOffering       types.String         `tfsdk:"product_offering"`
	SecurityGroups        types.List           `tfsdk:"security_groups"`
	UserData              types.String         `tfsdk:"user_data"`
	ServerGroupId         types.String         `tfsdk:"server_group_id"`
	ServerTypeId          types.String         `tfsdk:"server_type_id"`
	State                 types.String         `tfsdk:"state"`
	BootVolume            ServerResourceVolume `tfsdk:"boot_volume"`
	ExtraVolumes          types.Map            `tfsdk:"extra_volumes"`
	VpcId                 types.String         `tfsdk:"vpc_id"`
	Tags                  types.Map            `tfsdk:"tags"`
}

type ServerResourceNetwork struct {
	SubnetId    types.String `tfsdk:"subnet_id"`
	PortId      types.String `tfsdk:"port_id"`
	FixedIp     types.String `tfsdk:"fixed_ip"`
	PublicIpId  types.String `tfsdk:"public_ip_id"`
	StaticNatId types.String `tfsdk:"static_nat_id"`
}

type Tag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type ServerResourceVolume struct {
	Id                  types.String `tfsdk:"id"`
	DeleteOnTermination types.Bool   `tfsdk:"delete_on_termination"`
	Size                types.Int32  `tfsdk:"size"`
	Type                types.String `tfsdk:"type"`
}

func (m Server) AttributeTypes() map[string]attr.Type {
	ipAddressesObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ip_address": types.StringType,
			"version":    types.Int32Type,
		},
	}

	addressesObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ip_addresses": types.ListType{
				ElemType: ipAddressesObjectType,
			},
			"subnet_name": types.StringType,
		},
	}

	volumesObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"delete_on_termination": types.BoolType,
			"id":                    types.StringType,
		},
	}

	return map[string]attr.Type{
		"account_id": types.StringType,
		"addresses": types.ListType{
			ElemType: addressesObjectType,
		},
		"auto_scaling_group_id":   types.StringType,
		"created_at":              types.StringType,
		"created_by":              types.StringType,
		"disk_config":             types.StringType,
		"id":                      types.StringType,
		"image_id":                types.StringType,
		"keypair_name":            types.StringType,
		"launch_configuration_id": types.StringType,
		"locked":                  types.BoolType,
		"metadata": types.MapType{
			ElemType: types.StringType,
		},
		"modified_at":             types.StringType,
		"name":                    types.StringType,
		"planned_compute_os_type": types.StringType,
		"product_category":        types.StringType,
		"product_offering":        types.StringType,
		"security_groups": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
				},
			},
		},
		"server_group_id": types.StringType,
		"server_type": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"disk":      types.Int32Type,
				"ephemeral": types.Int32Type,
				"extra_specs": types.MapType{
					ElemType: types.StringType,
				},
				"id":    types.StringType,
				"name":  types.StringType,
				"ram":   types.Int32Type,
				"swap":  types.Int32Type,
				"vcpus": types.Int32Type,
			},
		},
		"state": types.StringType,
		"volumes": types.ListType{
			ElemType: volumesObjectType,
		},
		"vpc_id": types.StringType,
	}
}

// -------------------- Server Group -------------------- //

type ServerGroupDataSourceIds struct {
	Filter []filter.Filter `tfsdk:"filter"`
	Ids    []types.String  `tfsdk:"ids"`
}

type ServerGroupDataSource struct {
	Id          types.String    `tfsdk:"id"`
	Filter      []filter.Filter `tfsdk:"filter"`
	ServerGroup types.Object    `tfsdk:"server_group"`
}

type ServerGroup struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Policy    types.String `tfsdk:"policy"`
	AccountId types.String `tfsdk:"account_id"`
	UserId    types.String `tfsdk:"user_id"`
	Members   types.List   `tfsdk:"members"`
}

type ServerGroupResource struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Policy    types.String `tfsdk:"policy"`
	AccountId types.String `tfsdk:"account_id"`
	UserId    types.String `tfsdk:"user_id"`
	Members   types.List   `tfsdk:"members"`
	Tags      types.Map    `tfsdk:"tags"`
}

func (m ServerGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"name":       types.StringType,
		"policy":     types.StringType,
		"account_id": types.StringType,
		"user_id":    types.StringType,
		"members": types.ListType{
			ElemType: types.StringType,
		},
	}
}

// -------------------- Image -------------------- //

type ImageDataSourceIds struct {
	ScpImageType         types.String    `tfsdk:"scp_image_type"`
	ScpOriginalImageType types.String    `tfsdk:"scp_original_image_type"`
	Name                 types.String    `tfsdk:"name"`
	OsDistro             types.String    `tfsdk:"os_distro"`
	Status               types.String    `tfsdk:"status"`
	Visibility           types.String    `tfsdk:"visibility"`
	Filter               []filter.Filter `tfsdk:"filter"`
	Ids                  []types.String  `tfsdk:"ids"`
}

type ImageDataSource struct {
	Id                   types.String    `tfsdk:"id"`
	ScpImageType         types.String    `tfsdk:"scp_image_type"`
	ScpOriginalImageType types.String    `tfsdk:"scp_original_image_type"`
	Name                 types.String    `tfsdk:"name"`
	OsDistro             types.String    `tfsdk:"os_distro"`
	Status               types.String    `tfsdk:"status"`
	Visibility           types.String    `tfsdk:"visibility"`
	Filter               []filter.Filter `tfsdk:"filter"`
	Image                types.Object    `tfsdk:"image"`
}

type Image struct {
	Volumes              types.String `tfsdk:"volumes"`
	Checksum             types.String `tfsdk:"checksum"`
	ContainerFormat      types.String `tfsdk:"container_format"`
	DiskFormat           types.String `tfsdk:"disk_format"`
	File                 types.String `tfsdk:"file"`
	Id                   types.String `tfsdk:"id"`
	MinDisk              types.Int32  `tfsdk:"min_disk"`
	MinRam               types.Int32  `tfsdk:"min_ram"`
	Name                 types.String `tfsdk:"name"`
	OsDistro             types.String `tfsdk:"os_distro"`
	OsHashAlgo           types.String `tfsdk:"os_hash_algo"`
	OsHashValue          types.String `tfsdk:"os_hash_value"`
	OsHidden             types.Bool   `tfsdk:"os_hidden"`
	Owner                types.String `tfsdk:"owner"`
	OwnerAccountName     types.String `tfsdk:"owner_account_name"`
	OwnerUserName        types.String `tfsdk:"owner_user_name"`
	Protected            types.Bool   `tfsdk:"protected"`
	RootDeviceName       types.String `tfsdk:"root_device_name"`
	ScpImageType         types.String `tfsdk:"scp_image_type"`
	ScpK8sVersion        types.String `tfsdk:"scp_k8s_version"`
	ScpOriginalImageType types.String `tfsdk:"scp_original_image_type"`
	ScpOsVersion         types.String `tfsdk:"scp_os_version"`
	Size                 types.Int64  `tfsdk:"size"`
	Status               types.String `tfsdk:"status"`
	VirtualSize          types.Int64  `tfsdk:"virtual_size"`
	Visibility           types.String `tfsdk:"visibility"`
	Url                  types.String `tfsdk:"url"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

type ImageResource struct {
	InstanceId           types.String `tfsdk:"instance_id"`
	Volumes              types.String `tfsdk:"volumes"`
	Checksum             types.String `tfsdk:"checksum"`
	ContainerFormat      types.String `tfsdk:"container_format"`
	DiskFormat           types.String `tfsdk:"disk_format"`
	File                 types.String `tfsdk:"file"`
	Id                   types.String `tfsdk:"id"`
	MinDisk              types.Int32  `tfsdk:"min_disk"`
	MinRam               types.Int32  `tfsdk:"min_ram"`
	Name                 types.String `tfsdk:"name"`
	OsDistro             types.String `tfsdk:"os_distro"`
	OsHashAlgo           types.String `tfsdk:"os_hash_algo"`
	OsHashValue          types.String `tfsdk:"os_hash_value"`
	OsHidden             types.Bool   `tfsdk:"os_hidden"`
	Owner                types.String `tfsdk:"owner"`
	OwnerAccountName     types.String `tfsdk:"owner_account_name"`
	OwnerUserName        types.String `tfsdk:"owner_user_name"`
	Protected            types.Bool   `tfsdk:"protected"`
	RootDeviceName       types.String `tfsdk:"root_device_name"`
	ScpImageType         types.String `tfsdk:"scp_image_type"`
	ScpK8sVersion        types.String `tfsdk:"scp_k8s_version"`
	ScpOriginalImageType types.String `tfsdk:"scp_original_image_type"`
	ScpOsVersion         types.String `tfsdk:"scp_os_version"`
	Size                 types.Int64  `tfsdk:"size"`
	Status               types.String `tfsdk:"status"`
	VirtualSize          types.Int64  `tfsdk:"virtual_size"`
	Visibility           types.String `tfsdk:"visibility"`
	Url                  types.String `tfsdk:"url"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	Tags                 types.Map    `tfsdk:"tags"`
}

func (m Image) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"volumes":                 types.StringType,
		"checksum":                types.StringType,
		"container_format":        types.StringType,
		"disk_format":             types.StringType,
		"file":                    types.StringType,
		"id":                      types.StringType,
		"min_disk":                types.Int32Type,
		"min_ram":                 types.Int32Type,
		"name":                    types.StringType,
		"os_distro":               types.StringType,
		"os_hash_algo":            types.StringType,
		"os_hash_value":           types.StringType,
		"os_hidden":               types.BoolType,
		"owner":                   types.StringType,
		"owner_account_name":      types.StringType,
		"owner_user_name":         types.StringType,
		"protected":               types.BoolType,
		"root_device_name":        types.StringType,
		"scp_image_type":          types.StringType,
		"scp_k8s_version":         types.StringType,
		"scp_original_image_type": types.StringType,
		"scp_os_version":          types.StringType,
		"size":                    types.Int64Type,
		"status":                  types.StringType,
		"virtual_size":            types.Int64Type,
		"visibility":              types.StringType,
		"url":                     types.StringType,
		"created_at":              types.StringType,
		"updated_at":              types.StringType,
	}
}

// -------------------- Handler -------------------- //

type UpdateHandler struct {
	Fields  []string
	Handler func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error
}
