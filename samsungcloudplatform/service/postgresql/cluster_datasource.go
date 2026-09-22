package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/database"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &postgresqlClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlClusterDataSource{}
)

func NewPostgresqlClusterDataSource() datasource.DataSource {
	return &postgresqlClusterDataSource{}
}

type postgresqlClusterDataSource struct {
	config  *scpsdk.Configuration
	client  *postgresql.Client
	clients *client.SCPClient
}

func (d *postgresqlClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_cluster"
}

func (d *postgresqlClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Cluster.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "Cluster ID\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
				MarkdownDescription: "Cluster ID\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
				Optional:            true,
			},
			common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
				Description: "A detail of Cluster.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the endpoint.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						MarkdownDescription: "The identifier of the account that owns the endpoint.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("AllowableIpAddresses"): schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         "Allowed IP addresses list  \n  - example: ['192.168.10.1/32']",
						MarkdownDescription: "Allowed IP addresses list  \n  - example: ['192.168.10.1/32']",
						Computed:            true,
					},
					common.ToSnakeCase("DbaasEngine"): schema.StringAttribute{
						Description:         "DBaaS engine\n  - example: PostgreSQL",
						MarkdownDescription: "DBaaS engine\n  - example: PostgreSQL",
						Computed:            true,
					},
					common.ToSnakeCase("DbaasEngineVersionName"): schema.StringAttribute{
						Description:         "DBaaS engine version name\n  - example: PostgreSQL 16.4",
						MarkdownDescription: "DBaaS engine version name\n  - example: PostgreSQL 16.4",
						Computed:            true,
					},
					common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
						Description:         "NAT availability\n  - example: false",
						MarkdownDescription: "NAT availability\n  - example: false",
						Computed:            true,
					},
					common.ToSnakeCase("HaEnabled"): schema.BoolAttribute{
						Description:         "HA availability\n  - example: false",
						MarkdownDescription: "HA availability\n  - example: false",
						Computed:            true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "Identifier of the resource.\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
						MarkdownDescription: "Identifier of the resource.\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
						Computed:            true,
					},
					common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
						Description: "InitConfigOption.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("AuditEnabled"): schema.BoolAttribute{
								Description:         "Audit Log Setting\n  - example: true",
								MarkdownDescription: "Audit Log Setting\n  - example: true",
								Computed:            true,
							},
							common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
								Description: "BackupOption",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("ArchiveFrequencyMinute"): schema.StringAttribute{
										Description:         "Backup starting time (minute) \n  - example: 60",
										MarkdownDescription: "Backup starting time (minute) \n  - example: 60",
										Computed:            true,
									},
									common.ToSnakeCase("RetentionPeriodDay"): schema.StringAttribute{
										Description:         "Backup retention period (day) \n  - example: 7",
										MarkdownDescription: "Backup retention period (day) \n  - example: 7",
										Computed:            true,
									},
									common.ToSnakeCase("StartingTimeHour"): schema.StringAttribute{
										Description:         "Backup starting time (hour) \n  - example: 12",
										MarkdownDescription: "Backup starting time (hour) \n  - example: 12",
										Computed:            true,
									},
								},
							},
							common.ToSnakeCase("DatabaseEncoding"): schema.StringAttribute{
								Description:         "Database encoding\n  - example: UTF-8",
								MarkdownDescription: "Database encoding\n  - example: UTF-8",
								Computed:            true,
							},
							common.ToSnakeCase("DatabaseLocale"): schema.StringAttribute{
								Description:         "Database locale\n  - example: C",
								MarkdownDescription: "Database locale\n  - example: C",
								Computed:            true,
							},
							common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
								Description:         "Database name\n  - example: mydb",
								MarkdownDescription: "Database name\n  - example: mydb",
								Computed:            true,
							},
							common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
								Description:         "Database service port\n  - example: 2866",
								MarkdownDescription: "Database service port\n  - example: 2866",
								Computed:            true,
							},
							common.ToSnakeCase("DatabaseUserName"): schema.StringAttribute{
								Description:         "Database user name\n  - example: mydb",
								MarkdownDescription: "Database user name\n  - example: mydb",
								Computed:            true,
							},
							common.ToSnakeCase("OriginRegion"): schema.StringAttribute{
								Description:         "Origin region of the source cluster (set for restored/replica clusters)\n  - example: kr-west1",
								MarkdownDescription: "Origin region of the source cluster (set for restored/replica clusters)\n  - example: kr-west1",
								Computed:            true,
							},
						},
					},
					common.ToSnakeCase("InstanceCount"): schema.Int32Attribute{
						Description:         "Instance Count\n  - example: 1",
						MarkdownDescription: "Instance Count\n  - example: 1",
						Computed:            true,
					},
					common.ToSnakeCase("InstanceGroups"): schema.ListNestedAttribute{
						Description: "InstanceGroups",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description:         "Instance group ID\n  - example: ee48b333d5a84097adc079dec17ab872",
									MarkdownDescription: "Instance group ID\n  - example: ee48b333d5a84097adc079dec17ab872",
									Computed:            true,
								},
								common.ToSnakeCase("RoleType"): schema.StringAttribute{
									Description:         "Role type\n  - example: ACTIVE",
									MarkdownDescription: "Role type\n  - example: ACTIVE",
									Computed:            true,
								},
								common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
									Description:         "Server type name\n  - example: db1v2m4",
									MarkdownDescription: "Server type name\n  - example: db1v2m4",
									Computed:            true,
								},
								common.ToSnakeCase("BlockStorageGroups"): schema.ListNestedAttribute{
									Description: "BlockStorageGroups",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Id"): schema.StringAttribute{
												Description:         "Block storage group ID\n  - example: 1cf2c013bace4960878dfff31f6feec5",
												MarkdownDescription: "Block storage group ID\n  - example: 1cf2c013bace4960878dfff31f6feec5",
												Computed:            true,
											},
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description:         "Block storage group name\n  - example: cluster-Disk-00",
												MarkdownDescription: "Block storage group name\n  - example: cluster-Disk-00",
												Computed:            true,
											},
											common.ToSnakeCase("RoleType"): schema.StringAttribute{
												Description:         "Block storage role type\n  - example: OS",
												MarkdownDescription: "Block storage role type\n  - example: OS",
												Computed:            true,
											},
											common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
												Description:         "Size (GB)\n  - example: 104",
												MarkdownDescription: "Size (GB)\n  - example: 104",
												Computed:            true,
											},
											common.ToSnakeCase("VolumeType"): schema.StringAttribute{
												Description:         "Volume type\n  - example: SSD",
												MarkdownDescription: "Volume type\n  - example: SSD",
												Computed:            true,
											},
										},
									},
								},
								common.ToSnakeCase("Instances"): schema.ListNestedAttribute{
									Description: "Instances",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description:         "Instance name\n  - example: test001",
												MarkdownDescription: "Instance name\n  - example: test001",
												Computed:            true,
											},
											common.ToSnakeCase("RoleType"): schema.StringAttribute{
												Description:         "Role type\n  - example: ACTIVE",
												MarkdownDescription: "Role type\n  - example: ACTIVE",
												Computed:            true,
											},
											common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
												Description:         "User subnet IP address\n  - example: 192.168.4.22",
												MarkdownDescription: "User subnet IP address\n  - example: 192.168.4.22",
												Computed:            true,
											},
											common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
												Description:         "Public IP ID\n  - example: 90a68b14850741598ecacd0eb190873e",
												MarkdownDescription: "Public IP ID\n  - example: 90a68b14850741598ecacd0eb190873e",
												Computed:            true,
											},
										},
									},
								},
							},
						},
					},
					common.ToSnakeCase("IsKernelPatchable"): schema.BoolAttribute{
						Description:         "Whether the kernel can be patched\n  - example: true",
						MarkdownDescription: "Whether the kernel can be patched\n  - example: true",
						Computed:            true,
					},
					common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
						Description: "MaintenanceOption",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
								Description:         "Period in hours\n  - example: 1",
								MarkdownDescription: "Period in hours\n  - example: 1",
								Computed:            true,
							},
							common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
								Description:         "Starting day of week\n  - example: MON",
								MarkdownDescription: "Starting day of week\n  - example: MON",
								Computed:            true,
							},
							common.ToSnakeCase("StartingTime"): schema.StringAttribute{
								Description:         "Starting time\n  - example: 0000",
								MarkdownDescription: "Starting time\n  - example: 0000",
								Computed:            true,
							},
							common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
								Description:         "Use maintenance option\n  - example: true",
								MarkdownDescription: "Use maintenance option\n  - example: true",
								Computed:            true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Cluster name\n  - example: mytest",
						MarkdownDescription: "Cluster name\n  - example: mytest",
						Computed:            true,
					},
					common.ToSnakeCase("OriginClusterId"): schema.StringAttribute{
						Description:         "Origin cluster ID\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
						MarkdownDescription: "Origin cluster ID\n  - example: 8443a2afe9ce47c8958aa9dd010934d9",
						Computed:            true,
					},
					common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
						Description:         "Product image type\n  - example: PostgreSQL Community",
						MarkdownDescription: "Product image type\n  - example: PostgreSQL Community",
						Computed:            true,
					},
					common.ToSnakeCase("ProductType"): schema.StringAttribute{
						Description:         "Product type\n  - example: PostgreSQL Community",
						MarkdownDescription: "Product type\n  - example: PostgreSQL Community",
						Computed:            true,
					},
					common.ToSnakeCase("Replicas"): schema.SetAttribute{
						Description:         "Replicas list IDs\n  - example: ['6bc0cecc7a21488f974602574da86071']",
						MarkdownDescription: "Replicas list IDs\n  - example: ['6bc0cecc7a21488f974602574da86071']",
						Computed:            true,
						ElementType:         types.StringType,
					},
					common.ToSnakeCase("RoleType"): schema.StringAttribute{
						Description:         "Role type\n  - example: ORIGIN",
						MarkdownDescription: "Role type\n  - example: ORIGIN",
						Computed:            true,
					},
					common.ToSnakeCase("ServiceState"): schema.StringAttribute{
						Description:         "Service state\n  - example: RUNNING",
						MarkdownDescription: "Service state\n  - example: RUNNING",
						Computed:            true,
					},
					common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
						Description:         "Software version\n  - example: COMMUNITY 17.7",
						MarkdownDescription: "Software version\n  - example: COMMUNITY 17.7",
						Computed:            true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description:         "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249",
						MarkdownDescription: "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249",
						Computed:            true,
					},
					common.ToSnakeCase("Timezone"): schema.StringAttribute{
						Description:         "Timezone\n  - example: Asia/Seoul",
						MarkdownDescription: "Timezone\n  - example: Asia/Seoul",
						Computed:            true,
					},
					common.ToSnakeCase("VipPublicIpAddress"): schema.StringAttribute{
						Description:         "(VIP) Public IP address\n  - example: 10.10.10.10",
						MarkdownDescription: "(VIP) Public IP address\n  - example: 10.10.10.10",
						Computed:            true,
					},
					common.ToSnakeCase("VipPublicIpId"): schema.StringAttribute{
						Description:         "(VIP) Public IP ID\n  - example: 88a68b14850741599ecacd0eb190999a",
						MarkdownDescription: "(VIP) Public IP ID\n  - example: 88a68b14850741599ecacd0eb190999a",
						Computed:            true,
					},
					common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
						Description:         "Virtual IP address\n  - example: 192.168.4.30",
						MarkdownDescription: "Virtual IP address\n  - example: 192.168.4.30",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description:         "Created by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
						MarkdownDescription: "Created by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description:         "Modified by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
						MarkdownDescription: "Modified by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
						Computed:            true,
					},
					common.ToSnakeCase("ServiceWatchLogCollection"): schema.BoolAttribute{
						Description:         "ServiceWatchLogCollection\n - example: false",
						MarkdownDescription: "ServiceWatchLogCollection\n - example: false",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *postgresqlClusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Postgresql
	d.clients = inst.Client
}

func (d *postgresqlClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresql.ClusterDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Cluster",
			err.Error(),
		)
		return
	}

	var allowableIpAddresses types.Set
	if len(data.AllowableIpAddresses) == 0 {
		allowableIpAddresses, _ = types.SetValue(types.StringType, []attr.Value{})
	} else {
		ipAddresses := make([]attr.Value, len(data.AllowableIpAddresses))
		for i, ipAddress := range data.AllowableIpAddresses {
			ipAddresses[i] = types.StringValue(ipAddress)
		}
		allowableIpAddresses, _ = types.SetValue(types.StringType, ipAddresses)
	}

	var BackupOption = postgresql.BackupOption{}
	if data.InitConfigOption.BackupOption.Get() != nil {
		BackupOption = postgresql.BackupOption{
			ArchiveFrequencyMinute: types.StringPointerValue(data.InitConfigOption.BackupOption.Get().ArchiveFrequencyMinute.Get()),
			RetentionPeriodDay:     types.StringPointerValue(data.InitConfigOption.BackupOption.Get().RetentionPeriodDay.Get()),
			StartingTimeHour:       types.StringPointerValue(data.InitConfigOption.BackupOption.Get().StartingTimeHour.Get()),
		}
	}

	var initConfigOption = &postgresql.InitConfigOptionBase{
		AuditEnabled:     types.BoolPointerValue(data.InitConfigOption.AuditEnabled),
		BackupOption:     BackupOption,
		DatabaseEncoding: types.StringPointerValue(data.InitConfigOption.DatabaseEncoding.Get()),
		DatabaseLocale:   types.StringPointerValue(data.InitConfigOption.DatabaseLocale.Get()),
		DatabaseName:     types.StringValue(data.InitConfigOption.DatabaseName),
		DatabasePort:     types.Int32PointerValue(data.InitConfigOption.DatabasePort.Get()),
		DatabaseUserName: types.StringValue(data.InitConfigOption.DatabaseUserName),
		OriginRegion:     types.StringPointerValue(data.InitConfigOption.OriginRegion.Get()),
	}

	var InstanceGroups []database.InstanceGroup
	for _, instanceGroup := range data.InstanceGroups {
		var BlockStorage []database.BlockStorageGroup
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			BlockStorage = append(BlockStorage, database.BlockStorageGroup{
				Id:         types.StringValue(blockStorage.Id),
				Name:       types.StringValue(blockStorage.Name),
				RoleType:   types.StringValue(string(blockStorage.RoleType)),
				SizeGb:     types.Int32Value(blockStorage.SizeGb),
				VolumeType: types.StringValue(string(blockStorage.VolumeType)),
			})
		}

		var Instance []database.Instance
		for _, instance := range instanceGroup.Instances {
			Instance = append(Instance, database.Instance{
				Name:             types.StringValue(instance.Name),
				RoleType:         types.StringValue(string(instance.RoleType)),
				ServiceIpAddress: types.StringPointerValue(instance.ServiceIpAddress.Get()),
				PublicIpId:       types.StringPointerValue(instance.PublicIpId.Get()),
			})
		}

		blockStorageGroupList, blockStorageDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: database.BlockStorageGroup{}.AttributeTypes()}, BlockStorage)
		resp.Diagnostics.Append(blockStorageDiags...)
		instanceList, instanceDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: database.Instance{}.AttributeTypes()}, Instance)
		resp.Diagnostics.Append(instanceDiags...)

		InstanceGroups = append(InstanceGroups, database.InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: blockStorageGroupList,
			Instances:          instanceList,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		})
	}

	var MaintenanceOption *postgresql.MaintenanceOption
	if data.MaintenanceOption.IsSet() && data.MaintenanceOption.Get() != nil {
		MaintenanceOption = &postgresql.MaintenanceOption{
			PeriodHour:           types.StringPointerValue(data.MaintenanceOption.Get().PeriodHour.Get()),
			StartingDayOfWeek:    types.StringPointerValue((*string)(data.MaintenanceOption.Get().StartingDayOfWeek.Get())),
			StartingTime:         types.StringPointerValue(data.MaintenanceOption.Get().StartingTime.Get()),
			UseMaintenanceOption: types.BoolPointerValue(data.MaintenanceOption.Get().UseMaintenanceOption),
		}
	}

	var replicas types.Set
	if len(data.Replicas) == 0 {
		replicas, _ = types.SetValue(types.StringType, []attr.Value{})
	} else {
		replicaIds := make([]attr.Value, len(data.Replicas))
		for i, replicaId := range data.Replicas {
			replicaIds[i] = types.StringValue(replicaId)
		}
		replicas, _ = types.SetValue(types.StringType, replicaIds)
	}

	var postgresqlState = postgresql.ClusterDetail{
		AccountId:                 types.StringValue(data.AccountId),
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngine:               types.StringValue(data.DbaasEngine),
		NatEnabled:                types.BoolPointerValue(data.NatEnabled),
		HaEnabled:                 types.BoolPointerValue(data.HaEnabled),
		Id:                        types.StringValue(data.Id),
		InitConfigOption:          initConfigOption,
		InstanceCount:             types.Int32PointerValue(data.InstanceCount),
		InstanceGroups:            InstanceGroups,
		IsKernelPatchable:         types.BoolValue(data.IsKernelPatchable),
		MaintenanceOption:         MaintenanceOption,
		Name:                      types.StringValue(data.Name),
		OriginClusterId:           types.StringPointerValue(data.OriginClusterId.Get()),
		ProductType:               types.StringValue(string(data.ProductType)),
		Replicas:                  replicas,
		RoleType:                  types.StringPointerValue((*string)(data.RoleType.Get())),
		ServiceState:              types.StringValue(string(data.ServiceState)),
		SoftwareVersion:           types.StringValue(data.SoftwareVersion),
		SubnetId:                  types.StringValue(data.SubnetId),
		Timezone:                  types.StringValue(data.Timezone),
		VipPublicIpId:             types.StringPointerValue(data.VipPublicIpId.Get()),
		VirtualIpAddress:          types.StringPointerValue(data.VirtualIpAddress.Get()),
		CreatedAt:                 types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                 types.StringValue(data.CreatedBy),
		ModifiedAt:                types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                types.StringValue(data.ModifiedBy),
		ServiceWatchLogCollection: types.BoolValue(data.GetServiceWatchLogCollection()),
	}

	state.ClusterDetail = &postgresqlState

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
