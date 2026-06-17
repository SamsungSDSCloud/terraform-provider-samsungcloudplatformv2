package virtualserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		Description:         "Retrieves virtual server information.",
		MarkdownDescription: "Retrieves information about a single virtual server.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "Server ID to query.\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
				MarkdownDescription: "Server ID to query.\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
				Optional:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Server name to query.\n  - example: my-server\n  - minLength: 1\n  - maxLength: 255",
				MarkdownDescription: "Server name to query.\n  - example: my-server\n  - minLength: 1\n  - maxLength: 255",
				Optional:            true,
			},
			common.ToSnakeCase("Ip"): schema.StringAttribute{
				Description:         "IP address to filter servers.\n  - example: 192.168.1.100",
				MarkdownDescription: "IP address to filter servers.\n  - example: 192.168.1.100",
				Optional:            true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "Server state to filter.\n  - Available values: ACTIVE, SHUTOFF, ERROR",
				MarkdownDescription: "Server state to filter.\n  - Available values: ACTIVE, SHUTOFF, ERROR",
				Optional:            true,
			},
			common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
				Description:         "Product category.\n  - Available values: compute, container",
				MarkdownDescription: "Product category.\n  - Available values: compute, container",
				Optional:            true,
			},
			common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
				Description:         "Product offering.\n  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm",
				MarkdownDescription: "Product offering.\n  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm\n  - note: Use gpu_server for GPU instances",
				Optional:            true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description:         "VPC ID.\n  - example: cc976b621087484ea5fd527f4b78708b",
				MarkdownDescription: "VPC ID.\n  - example: cc976b621087484ea5fd527f4b78708b",
				Optional:            true,
			},
			common.ToSnakeCase("ServerTypeId"): schema.StringAttribute{
				Description:         "Server type ID.\n  - example: s1v1m2",
				MarkdownDescription: "Server type ID.\n  - example: s1v1m2",
				Optional:            true,
			},
			common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
				Description:         "Auto Scaling Group ID.\n  - example: 52613bd852b04b39adcb15a8364d856d",
				MarkdownDescription: "Auto Scaling Group ID.\n  - example: 52613bd852b04b39adcb15a8364d856d",
				Optional:            true,
			},
			common.ToSnakeCase("Server"): schema.SingleNestedAttribute{
				Description:         "Retrieved server information.",
				MarkdownDescription: "Retrieved server information including configuration and status.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description:         "Account ID.",
						MarkdownDescription: "Account ID.",
						Computed:            true,
					},
					common.ToSnakeCase("Addresses"): schema.ListNestedAttribute{
						Description:         "Network address list.",
						MarkdownDescription: "Network address list.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("IpAddresses"): schema.ListNestedAttribute{
									Description:         "IP address list.",
									MarkdownDescription: "IP address list.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("IpAddress"): schema.StringAttribute{
												Description:         "IP address.",
												MarkdownDescription: "IP address.",
												Computed:            true,
											},
											common.ToSnakeCase("Version"): schema.Int32Attribute{
												Description:         "IP version.",
												MarkdownDescription: "IP version.",
												Computed:            true,
											},
										},
									},
								},
								common.ToSnakeCase("SubnetName"): schema.StringAttribute{
									Description:         "Subnet name.",
									MarkdownDescription: "Subnet name.",
									Computed:            true,
								},
							},
						},
					},
					common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
						Description:         "Auto Scaling Group ID.",
						MarkdownDescription: "Auto Scaling Group ID.",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description:         "Creation timestamp.",
						MarkdownDescription: "Creation timestamp.",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description:         "Creator ID.",
						MarkdownDescription: "Creator ID.",
						Computed:            true,
					},
					common.ToSnakeCase("DiskConfig"): schema.StringAttribute{
						Description:         "Disk configuration mode.",
						MarkdownDescription: "Disk configuration mode.",
						Computed:            true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "Server ID.",
						MarkdownDescription: "Server ID.",
						Computed:            true,
					},
					common.ToSnakeCase("ImageId"): schema.StringAttribute{
						Description:         "Image ID.",
						MarkdownDescription: "Image ID.",
						Computed:            true,
					},
					common.ToSnakeCase("KeypairName"): schema.StringAttribute{
						Description:         "Keypair name.",
						MarkdownDescription: "Keypair name.",
						Computed:            true,
					},
					common.ToSnakeCase("LaunchConfigurationId"): schema.StringAttribute{
						Description:         "Launch Configuration ID.",
						MarkdownDescription: "Launch Configuration ID.",
						Computed:            true,
					},
					common.ToSnakeCase("Locked"): schema.BoolAttribute{
						Description:         "Lock status.",
						MarkdownDescription: "Lock status.",
						Computed:            true,
					},
					common.ToSnakeCase("Metadata"): schema.MapAttribute{
						Description:         "Metadata.",
						MarkdownDescription: "Metadata.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description:         "Modification timestamp.",
						MarkdownDescription: "Modification timestamp.",
						Computed:            true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Server name.",
						MarkdownDescription: "Server name.",
						Computed:            true,
					},
					common.ToSnakeCase("PlannedComputeOsType"): schema.StringAttribute{
						Description:         "Planned compute OS type.",
						MarkdownDescription: "Planned compute OS type.",
						Computed:            true,
					},
					common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
						Description:         "Product category.",
						MarkdownDescription: "Product category.",
						Computed:            true,
					},
					common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
						Description:         "Product offering.",
						MarkdownDescription: "Product offering.",
						Computed:            true,
					},
					common.ToSnakeCase("SecurityGroups"): schema.ListNestedAttribute{
						Description:         "Security group list.",
						MarkdownDescription: "Security group list.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Name"): schema.StringAttribute{
									Description:         "Security group name.",
									MarkdownDescription: "Security group name.",
									Computed:            true,
								},
							},
						},
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description:         "Server group ID.",
						MarkdownDescription: "Server group ID.",
						Computed:            true,
					},
					common.ToSnakeCase("ServerType"): schema.SingleNestedAttribute{
						Description:         "Server type information.",
						MarkdownDescription: "Server type information.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Disk"): schema.Int32Attribute{
								Description:         "Disk size (GB).",
								MarkdownDescription: "Disk size (GB).",
								Computed:            true,
							},
							common.ToSnakeCase("Ephemeral"): schema.Int32Attribute{
								Description:         "Ephemeral disk size (GB).",
								MarkdownDescription: "Ephemeral disk size (GB).",
								Computed:            true,
							},
							common.ToSnakeCase("ExtraSpecs"): schema.MapAttribute{
								Description:         "Extra specifications.",
								MarkdownDescription: "Extra specifications.",
								Computed:            true,
								ElementType:         types.StringType,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description:         "Server type ID.",
								MarkdownDescription: "Server type ID.",
								Computed:            true,
							},
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description:         "Server type name.",
								MarkdownDescription: "Server type name.",
								Computed:            true,
							},
							common.ToSnakeCase("Ram"): schema.Int32Attribute{
								Description:         "RAM size (MB).",
								MarkdownDescription: "RAM size (MB).",
								Computed:            true,
							},
							common.ToSnakeCase("Swap"): schema.Int32Attribute{
								Description:         "Swap size (MB).",
								MarkdownDescription: "Swap size (MB).",
								Computed:            true,
							},
							common.ToSnakeCase("Vcpus"): schema.Int32Attribute{
								Description:         "Number of vCPUs.",
								MarkdownDescription: "Number of vCPUs.",
								Computed:            true,
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description:         "Server state.",
						MarkdownDescription: "Server state.",
						Computed:            true,
					},
					common.ToSnakeCase("Volumes"): schema.ListNestedAttribute{
						Description:         "Attached volume list.",
						MarkdownDescription: "Attached volume list.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("DeleteOnTermination"): schema.BoolAttribute{
									Description:         "Whether to delete volume when server is terminated.",
									MarkdownDescription: "Whether to delete volume when server is terminated.",
									Computed:            true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description:         "Volume ID.",
									MarkdownDescription: "Volume ID.",
									Computed:            true,
								},
							},
						},
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description:         "VPC ID.",
						MarkdownDescription: "VPC ID.",
						Computed:            true,
					},
					common.ToSnakeCase("PartitionNumber"): schema.Int32Attribute{
						Description:         "Partition number.",
						MarkdownDescription: "Partition number.",
						Computed:            true,
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
			PartitionNumber:       types.Int32PointerValue(server.PartitionNumber.Get()),
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
