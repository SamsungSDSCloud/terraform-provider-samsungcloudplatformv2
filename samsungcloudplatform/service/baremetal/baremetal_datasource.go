package baremetal

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/baremetal"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &baremetalBaremetalDataSource{}
	_ datasource.DataSourceWithConfigure = &baremetalBaremetalDataSource{}
)

// NewBaremetalBaremetalDataSource is a helper function to simplify the provider implementation.
func NewBaremetalBaremetalDataSource() datasource.DataSource {
	return &baremetalBaremetalDataSource{}
}

// baremetalBaremetalDataSource is the data source implementation.
type baremetalBaremetalDataSource struct {
	config  *scpsdk.Configuration
	client  *baremetal.Client
	clients *client.SCPClient
}

// Schema defines the schema for the data source.
func (d *baremetalBaremetalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = BaremetalDataSourceSchema()
}

// Metadata returns the data source type name.
func (d *baremetalBaremetalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_baremetal"
}

// Configure adds the provider configured client to the data source.
func (d *baremetalBaremetalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = inst.Client.Baremetal
	d.clients = inst.Client
}

func BaremetalDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Show Baremetal.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
				MarkdownDescription: "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
			},
			"hyper_threading_use": schema.BoolAttribute{
				Computed:            true,
				Description:         "Use Hyper Threading\n  - example: true",
				MarkdownDescription: "Use Hyper Threading\n  - example: true",
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Bare Metal Server ID",
				MarkdownDescription: "Bare Metal Server ID",
			},
			"image_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
				MarkdownDescription: "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
			},
			"image_version": schema.StringAttribute{
				Computed:            true,
				Description:         "Image version\n  - example: RHEL 8.7 for BM",
				MarkdownDescription: "Image version\n  - example: RHEL 8.7 for BM",
			},
			"init_script": schema.StringAttribute{
				Computed:            true,
				Description:         "Init script\n  - example: init script",
				MarkdownDescription: "Init script\n  - example: init script",
			},
			"local_subnet_info": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Interface Name\n  - example: ens8f1,bond_serv.2",
							MarkdownDescription: "Interface Name\n  - example: ens8f1,bond_serv.2",
						},
						"local_subnet_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Local Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
							MarkdownDescription: "Local Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
						},
						"policy_local_subnet_ip": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy Local Subnet IP\n  - example: 192.168.0.1",
							MarkdownDescription: "Policy Local Subnet IP\n  - example: 192.168.0.1",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "LocalSubnet IP state\n  - example: CREATING",
							MarkdownDescription: "LocalSubnet IP state\n  - example: CREATING",
						},
						"vlan_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Vlan Id\n  - example: 500",
							MarkdownDescription: "Vlan Id\n  - example: 500",
						},
						"vni_role_name": schema.StringAttribute{
							Computed:            true,
							Description:         "VNI Role Name\n  - example: mgmtJ4dzQBo",
							MarkdownDescription: "VNI Role Name\n  - example: mgmtJ4dzQBo",
						},
					},
				},
				Computed: true,
			},
			"lock_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Use Lock\n  - example: true",
				MarkdownDescription: "Use Lock\n  - example: true",
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
				MarkdownDescription: "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
			},
			"network_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
			},
			"os_type": schema.StringAttribute{
				Computed:            true,
				Description:         "OS type\n  - example: WINDOWS",
				MarkdownDescription: "OS type\n  - example: WINDOWS",
			},
			"placement_group_name": schema.StringAttribute{
				Computed:            true,
				Description:         "placement group name\n  - example: pg-group",
				MarkdownDescription: "placement group name\n  - example: pg-group",
			},
			"policy_ip": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy IP\n  - example: 192.168.0.1",
				MarkdownDescription: "Policy IP\n  - example: 192.168.0.1",
			},
			"private_nat_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"nat_id": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT ID\n  - example: cadcc59a44264eca9674be339cf3e857\n",
						MarkdownDescription: "NAT ID\n  - example: cadcc59a44264eca9674be339cf3e857\n",
					},
					"nat_ip": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT IP\n  - example: 192.170.2.10\n",
						MarkdownDescription: "NAT IP\n  - example: 192.170.2.10\n",
					},
					"nat_ip_id": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT IP ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac\n",
						MarkdownDescription: "NAT IP ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac\n",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "Static NAT state\n  - example: ACTIVE\n",
						MarkdownDescription: "Static NAT state\n  - example: ACTIVE\n",
					},
					"static_nat_id": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT ID\n  - example: 997e99959c7b415b84bbd250c9fe716c\n",
						MarkdownDescription: "NAT ID\n  - example: 997e99959c7b415b84bbd250c9fe716c\n",
					},
				},
				Description:         "Private Nat Info\n",
				MarkdownDescription: "Private Nat Info\n",
				Computed:            true,
			},
			"product_type_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Bare Metal Server scale ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
				MarkdownDescription: "Bare Metal Server scale ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
			},
			"public_nat_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"nat_id": schema.StringAttribute{
						Computed: true,
					},
					"nat_ip": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT IP\n  - example: 192.170.2.10\n",
						MarkdownDescription: "NAT IP\n  - example: 192.170.2.10\n",
					},
					"nat_ip_id": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT IP ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac\n",
						MarkdownDescription: "NAT IP ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac\n",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "Static NAT state\n  - example: ACTIVE\n",
						MarkdownDescription: "Static NAT state\n  - example: ACTIVE\n",
					},
					"static_nat_id": schema.StringAttribute{
						Computed:            true,
						Description:         "NAT ID\n  - example: 997e99959c7b415b84bbd250c9fe716c\n",
						MarkdownDescription: "NAT ID\n  - example: 997e99959c7b415b84bbd250c9fe716c\n",
					},
				},
				Description:         "Public Nat Info\n",
				MarkdownDescription: "Public Nat Info\n",
				Computed:            true,
			},
			"region_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"root_account": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Account\n  - example: rootaccount",
				MarkdownDescription: "Root Account\n  - example: rootaccount",
			},
			"server_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Bare Metal Server name\n  - example: bmserver-001",
				MarkdownDescription: "Bare Metal Server name\n  - example: bmserver-001",
			},
			"server_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Bare Metal Server scale type\n  - example: s1v8m32_metal",
				MarkdownDescription: "Bare Metal Server scale type\n  - example: s1v8m32_metal",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Bare Metal Server state\n  - example: RUNNING",
				MarkdownDescription: "Bare Metal Server state\n  - example: RUNNING",
			},
			"time_zone": schema.StringAttribute{
				Computed:            true,
				Description:         "Time Zone\n  - example: Asia/Seoul",
				MarkdownDescription: "Time Zone\n  - example: Asia/Seoul",
			},
			"use_local_subnet": schema.BoolAttribute{
				Computed:            true,
				Description:         "Use Local Subnet\n  - example: true",
				MarkdownDescription: "Use Local Subnet\n  - example: true",
			},
			"vpc_id": schema.StringAttribute{
				Computed:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
		},
	}
}

func (d *baremetalBaremetalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state baremetal.BaremetalDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetBaremetal(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Baremetal",
			err.Error(),
		)
		return
	}

	//basic info
	state.AccountId = types.StringValue(data.AccountId)
	state.ImageId = types.StringValue(data.ImageId)
	state.ImageVersion = types.StringValue(data.ImageVersion)
	state.OsType = types.StringValue(data.OsType)
	state.ProductTypeId = types.StringValue(data.ProductTypeId)
	state.RegionId = types.StringValue(data.RegionId)
	state.RootAccount = types.StringValue(data.RootAccount)
	state.ServerName = types.StringValue(data.ServerName)
	state.ServerType = types.StringValue(data.ServerType)
	state.State = types.StringValue(data.State)
	state.TimeZone = types.StringValue(data.TimeZone)

	// network info
	state.NetworkId = types.StringValue(data.NetworkId)
	state.VpcId = types.StringValue(data.VpcId)
	state.PolicyIp = types.StringValue(data.PolicyIp)
	state.UseLocalSubnet = types.BoolPointerValue(data.UseLocalSubnet)

	if data.PublicNatInfo.Get() != nil {
		publicNatInfo := baremetal.PublicNatInfoValue{
			NatId:       types.StringPointerValue(data.PublicNatInfo.Get().NatId.Get()),
			NatIp:       types.StringValue(data.PublicNatInfo.Get().NatIp),
			NatIpId:     types.StringValue(data.PublicNatInfo.Get().NatIpId),
			State:       types.StringValue(data.PublicNatInfo.Get().State),
			StaticNatId: types.StringPointerValue(data.PublicNatInfo.Get().StaticNatId.Get()),
		}
		publicNatInfoObject, _ := types.ObjectValueFrom(ctx, publicNatInfo.AttributeTypes(), publicNatInfo)
		state.PublicNatInfo = publicNatInfoObject
	}

	if data.PrivateNatInfo.Get() != nil {
		privateNatInfo := baremetal.PrivateNatInfoValue{
			NatId:       types.StringPointerValue(data.PrivateNatInfo.Get().NatId.Get()),
			NatIp:       types.StringValue(data.PrivateNatInfo.Get().NatIp),
			NatIpId:     types.StringValue(data.PrivateNatInfo.Get().NatIpId),
			State:       types.StringValue(data.PrivateNatInfo.Get().State),
			StaticNatId: types.StringPointerValue(data.PrivateNatInfo.Get().StaticNatId.Get()),
		}
		privateNatInfoObject, _ := types.ObjectValueFrom(ctx, privateNatInfo.AttributeTypes(), privateNatInfo)
		state.PrivateNatInfo = privateNatInfoObject
	}

	var localSubnetInfo baremetal.LocalSubnetInfo
	var localSubnetInfoList []baremetal.LocalSubnetInfo
	for _, localSubnetInfoResponse := range data.LocalSubnetInfo {
		localSubnetInfo = baremetal.LocalSubnetInfo{
			InterfaceName:       types.StringValue(localSubnetInfoResponse.InterfaceName),
			LocalSubnetId:       types.StringValue(localSubnetInfoResponse.LocalSubnetId),
			PolicyLocalSubnetIp: types.StringValue(localSubnetInfoResponse.PolicyLocalSubnetIp),
			State:               types.StringValue(localSubnetInfoResponse.State),
			VlanId:              types.StringValue(localSubnetInfoResponse.VlanId),
			VniRoleName:         types.StringValue(localSubnetInfoResponse.VniRoleName),
		}
		localSubnetInfoList = append(localSubnetInfoList, localSubnetInfo)
	}

	localSubnetInfoListType, _ := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: localSubnetInfo.AttributeTypes(),
	}, localSubnetInfoList)

	state.LocalSubnetInfo = localSubnetInfoListType

	// additional info
	state.PlacementGroupName = types.StringPointerValue(data.PlacementGroupName.Get())
	state.HyperThreadingUse = types.BoolValue(data.HyperThreadingUse)
	state.InitScript = types.StringValue(data.InitScript)
	state.LockEnabled = types.BoolPointerValue(data.LockEnabled)

	// metadata info
	state.CreatedBy = types.StringValue(data.CreatedBy)
	state.CreatedAt = types.StringValue(data.CreatedAt.String())
	state.ModifiedBy = types.StringValue(data.ModifiedBy)
	state.ModifiedAt = types.StringValue(data.ModifiedAt.String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
