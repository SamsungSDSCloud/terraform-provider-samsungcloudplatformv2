package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ datasource.DataSource              = &virtualServerServerDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualServerServerDataSource{}
)

func NewVirtualServerServerDataSource() datasource.DataSource {
	return &virtualServerServerDataSource{}
}

type virtualServerServerDataSource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_server"
}

func (d *virtualServerServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of servers.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("Ip"): schema.StringAttribute{
				Description: "Ip",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
				Description: "Product category",
				Optional:    true,
			},
			common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
				Description: "Product offering",
				Optional:    true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID",
				Optional:    true,
			},
			common.ToSnakeCase("ServerTypeId"): schema.StringAttribute{
				Description: "Server type ID",
				Optional:    true,
			},
			common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
				Description: "Auto scaling group ID",
				Optional:    true,
			},
			common.ToSnakeCase("Server"): schema.SingleNestedAttribute{
				Description: "Server.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("Addresses"): schema.ListNestedAttribute{
						Description: "Addresses",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("IpAddresses"): schema.ListNestedAttribute{
									Description: "IP addresses",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("IpAddress"): schema.StringAttribute{
												Description: "IP address",
												Computed:    true,
											},
											common.ToSnakeCase("Version"): schema.Int32Attribute{
												Description: "Version",
												Computed:    true,
											},
										},
									},
								},
								common.ToSnakeCase("SubnetName"): schema.StringAttribute{
									Description: "Subnet name",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
						Description: "Auto scaling group ID",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created by",
						Computed:    true,
					},
					common.ToSnakeCase("DiskConfig"): schema.StringAttribute{
						Description: "Disk config",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID",
						Computed:    true,
					},
					common.ToSnakeCase("ImageId"): schema.StringAttribute{
						Description: "Image ID",
						Computed:    true,
					},
					common.ToSnakeCase("KeypairName"): schema.StringAttribute{
						Description: "Keypair name",
						Computed:    true,
					},
					common.ToSnakeCase("LaunchConfigurationId"): schema.StringAttribute{
						Description: "Launch Configuration ID",
						Computed:    true,
					},
					common.ToSnakeCase("Locked"): schema.BoolAttribute{
						Description: "Locked",
						Computed:    true,
					},
					common.ToSnakeCase("Metadata"): schema.MapAttribute{
						Description: "Metadata",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified at",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("PlannedComputeOsType"): schema.StringAttribute{
						Description: "Planned compute os type",
						Computed:    true,
					},
					common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
						Description: "Product category",
						Computed:    true,
					},
					common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
						Description: "Product offering",
						Computed:    true,
					},
					common.ToSnakeCase("SecurityGroups"): schema.ListNestedAttribute{
						Description: "Security groups",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Name"): schema.StringAttribute{
									Description: "Name",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "Server group ID",
						Computed:    true,
					},
					common.ToSnakeCase("ServerType"): schema.SingleNestedAttribute{
						Description: "Server type",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Disk"): schema.Int32Attribute{
								Description: "Disk",
								Computed:    true,
							},
							common.ToSnakeCase("Ephemeral"): schema.Int32Attribute{
								Description: "Ephemeral",
								Computed:    true,
							},
							common.ToSnakeCase("ExtraSpecs"): schema.MapAttribute{
								Description: "Extra specs",
								Computed:    true,
								ElementType: types.StringType,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "ID",
								Computed:    true,
							},
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description: "Name",
								Computed:    true,
							},
							common.ToSnakeCase("Ram"): schema.Int32Attribute{
								Description: "Ram",
								Computed:    true,
							},
							common.ToSnakeCase("Swap"): schema.Int32Attribute{
								Description: "Swap",
								Computed:    true,
							},
							common.ToSnakeCase("Vcpus"): schema.Int32Attribute{
								Description: "Vcpus",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Volumes"): schema.ListNestedAttribute{
						Description: "Volumes",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("DeleteOnTermination"): schema.BoolAttribute{
									Description: "Delete on termination",
									Computed:    true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "ID",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "Vpc ID",
						Computed:    true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VirtualServer
	d.clients = inst.Client
}

func (d *virtualServerServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ServerDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetServers(d.clients, state.Name, state.Ip, state.State, state.ProductCategory, state.ProductOffering,
		state.VpcId, state.ServerTypeId, state.AutoScalingGroupId, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Server",
			err.Error(),
		)
		return
	}

	if len(ids) > 0 || !state.Id.IsNull() {
		id := virtualserverutil.SetResourceIdentifier(state.Id, ids, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		server, err := d.client.GetServer(ctx, id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Server",
				"Could not read Server ID "+id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		// Addresses
		var addresses []virtualserver.ServerAddress
		for _, address := range server.Addresses {
			var ipAddresses []virtualserver.ServerIpAddress
			for _, ipAddress := range address.IpAddresses {
				ipAddressState := virtualserver.ServerIpAddress{
					IpAddress: types.StringValue(ipAddress.IpAddress),
					Version:   types.Int32Value(ipAddress.Version),
				}
				ipAddresses = append(ipAddresses, ipAddressState)
			}
			addressState := virtualserver.ServerAddress{
				IpAddresses: ipAddresses,
				SubnetName:  types.StringValue(address.SubnetName),
			}
			addresses = append(addresses, addressState)
		}

		// Metadata
		metadataMap := make(map[string]attr.Value)

		for k, v := range server.Metadata {
			metadataMap[k] = types.StringValue(v.(string))
		}
		metadata, _ := types.MapValue(types.StringType, metadataMap)

		// PlannedComputeOsType
		plannedComputeOsTypeJson, _ := server.PlannedComputeOsType.MarshalJSON()
		plannedComputeOsType := strings.Trim(string(plannedComputeOsTypeJson), "\"")

		// ProductCategory
		productCategoryJson, _ := server.ProductCategory.MarshalJSON()
		productCategory := strings.Trim(string(productCategoryJson), "\"")

		// ProductOffering
		productOfferingJson, _ := server.ProductOffering.MarshalJSON()
		productOffering := strings.Trim(string(productOfferingJson), "\"")

		//SecurityGroups
		var securityGroups []virtualserver.SecurityGroup
		for _, securityGroup := range server.SecurityGroups {
			securityGroupState := virtualserver.SecurityGroup{
				Name: types.StringValue(securityGroup.Name),
			}
			securityGroups = append(securityGroups, securityGroupState)
		}

		//ServerType
		extraSpecsMap := make(map[string]attr.Value)

		for k, v := range server.ServerType.ExtraSpecs {
			extraSpecsMap[k] = types.StringValue(v.(string))
		}
		extraSpecs, _ := types.MapValue(types.StringType, extraSpecsMap)

		serverTypeState := virtualserver.ServerType{
			Disk:       types.Int32Value(server.ServerType.Disk),
			Ephemeral:  types.Int32Value(server.ServerType.Ephemeral),
			ExtraSpecs: extraSpecs,
			Id:         types.StringPointerValue(server.ServerType.Id.Get()),
			Name:       types.StringValue(server.ServerType.Name),
			Ram:        types.Int32Value(server.ServerType.Disk),
			Swap:       types.Int32Value(server.ServerType.Disk),
			Vcpus:      types.Int32Value(server.ServerType.Disk),
		}

		//Volumes
		var volumes []virtualserver.ServerVolume
		for _, volume := range server.Volumes {
			volumeState := virtualserver.ServerVolume{
				DeleteOnTermination: types.BoolValue(volume.DeleteOnTermination),
				Id:                  types.StringValue(volume.Id),
			}
			volumes = append(volumes, volumeState)
		}

		serverModel := virtualserver.Server{
			AccountId:             types.StringValue(server.AccountId),
			Addresses:             addresses,
			AutoScalingGroupId:    types.StringPointerValue(server.AutoScalingGroupId.Get()),
			CreatedAt:             types.StringValue(server.CreatedAt.Format(time.RFC3339)),
			CreatedBy:             types.StringValue(server.CreatedBy),
			DiskConfig:            types.StringValue(server.DiskConfig),
			Id:                    types.StringValue(server.Id),
			ImageId:               types.StringPointerValue(server.ImageId.Get()),
			KeypairName:           types.StringPointerValue(server.KeypairName.Get()),
			LaunchConfigurationId: types.StringPointerValue(server.LaunchConfigurationId.Get()),
			Locked:                types.BoolValue(server.Locked),
			Metadata:              metadata,
			ModifiedAt:            types.StringValue(server.ModifiedAt.Format(time.RFC3339)),
			Name:                  types.StringValue(server.Name),
			PlannedComputeOsType:  types.StringValue(plannedComputeOsType),
			ProductCategory:       types.StringValue(productCategory),
			ProductOffering:       types.StringValue(productOffering),
			SecurityGroups:        securityGroups,
			ServerGroupId:         types.StringPointerValue(server.ServerGroupId.Get()),
			ServerType:            serverTypeState,
			State:                 types.StringValue(server.State),
			Volumes:               volumes,
			VpcId:                 types.StringPointerValue(server.VpcId.Get()),
		}
		serverObjectValue, _ := types.ObjectValueFrom(ctx, serverModel.AttributeTypes(), serverModel)
		state.Server = serverObjectValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
