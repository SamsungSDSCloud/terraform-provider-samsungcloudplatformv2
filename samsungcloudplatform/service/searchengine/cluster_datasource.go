package searchengine

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ datasource.DataSource              = &searchengineClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &searchengineClusterDataSource{}
)

func NewSearchengineClusterDataSource() datasource.DataSource {
	return &searchengineClusterDataSource{}
}

type searchengineClusterDataSource struct {
	config  *scpsdk.Configuration
	client  *searchengine.Client
	clients *client.SCPClient
}

func (d *searchengineClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_searchengine_cluster"
}

func (d *searchengineClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Cluster.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
				Description: "A detail of Cluster.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("AllowableIpAddresses"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "AllowableIpAddresses",
						Computed:    true,
					},
					common.ToSnakeCase("DbaasEngine"): schema.StringAttribute{
						Description: "DbaasEngine",
						Computed:    true,
					},
					common.ToSnakeCase("IsCombined"): schema.BoolAttribute{
						Description: "IsCombined",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
						Description: "InitConfigOption.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
								Description: "BackupOption",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("RetentionPeriodDay"): schema.StringAttribute{
										Description: "RetentionPeriodDay",
										Computed:    true,
									},
									common.ToSnakeCase("StartingTimeHour"): schema.StringAttribute{
										Description: "StartingTimeHour",
										Computed:    true,
									},
								},
							},
							common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
								Description: "DatabasePort",
								Computed:    true,
							},
							common.ToSnakeCase("DatabaseUserName"): schema.StringAttribute{
								Description: "DatabaseUserName",
								Computed:    true,
							},
							common.ToSnakeCase("KibanaPort"): schema.Int32Attribute{
								Description: "KibanaPort",
								Computed:    true,
							},
							common.ToSnakeCase("DashboardsPort"): schema.Int32Attribute{
								Description: "DashboardsPort",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("InstanceCount"): schema.Int32Attribute{
						Description: "InstanceCount",
						Computed:    true,
					},
					common.ToSnakeCase("InstanceGroups"): schema.ListNestedAttribute{
						Description: "InstanceGroups",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Id",
									Computed:    true,
								},
								common.ToSnakeCase("RoleType"): schema.StringAttribute{
									Description: "RoleType",
									Computed:    true,
								},
								common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
									Description: "ServerTypeName",
									Computed:    true,
								},
								common.ToSnakeCase("BlockStorageGroups"): schema.ListNestedAttribute{
									Description: "BlockStorageGroups",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Id"): schema.StringAttribute{
												Description: "Id",
												Computed:    true,
											},
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description: "Name",
												Computed:    true,
											},
											common.ToSnakeCase("RoleType"): schema.StringAttribute{
												Description: "RoleType",
												Computed:    true,
											},
											common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
												Description: "SizeGb",
												Computed:    true,
											},
											common.ToSnakeCase("VolumeType"): schema.StringAttribute{
												Description: "VolumeType",
												Computed:    true,
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
												Description: "Name",
												Computed:    true,
											},
											common.ToSnakeCase("RoleType"): schema.StringAttribute{
												Description: "RoleType",
												Computed:    true,
											},
											common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
												Description: "ServiceIpAddress",
												Computed:    true,
											},
											common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
												Description: "PublicIpId",
												Computed:    true,
											},
											//common.ToSnakeCase("PublicIpAddress"): schema.StringAttribute{
											//	Description: "PublicIpAddress",
											//	Computed:    true,
											//},
											//common.ToSnakeCase("ServiceState"): schema.StringAttribute{
											//	Description: "ServiceState",
											//	Computed:    true,
											//},
										},
									},
								},
							},
						},
					},
					common.ToSnakeCase("License"): schema.SingleNestedAttribute{
						Description: "License",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("LicenseExpiryDate"): schema.StringAttribute{
								Description: "LicenseExpiryDate",
								Computed:    true,
							},
							common.ToSnakeCase("LicenseStatus"): schema.StringAttribute{
								Description: "LicenseStatus",
								Computed:    true,
							},
							common.ToSnakeCase("LicenseType"): schema.StringAttribute{
								Description: "LicenseType",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
						Description: "MaintenanceOption",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
								Description: "PeriodHour",
								Computed:    true,
							},
							common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
								Description: "StartingDayOfWeek",
								Computed:    true,
							},
							common.ToSnakeCase("StartingTime"): schema.StringAttribute{
								Description: "StartingTime",
								Computed:    true,
							},
							common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
								Description: "UseMaintenanceOption",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
						Description: "NatEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("ProductType"): schema.StringAttribute{
						Description: "ProductType",
						Computed:    true,
					},
					common.ToSnakeCase("ServiceState"): schema.StringAttribute{
						Description: "ServiceState",
						Computed:    true,
					},
					common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
						Description: "SoftwareVersion",
						Computed:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Computed:    true,
					},
					common.ToSnakeCase("Timezone"): schema.StringAttribute{
						Description: "Timezone",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *searchengineClusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Searchengine
	d.clients = inst.Client
}

func (d *searchengineClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state searchengine.ClusterDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Cluster",
			err.Error(),
		)
		return
	}

	allowableIpAddresses := make([]types.String, len(data.AllowableIpAddresses))
	for i, allowableIpAddress := range data.AllowableIpAddresses {
		allowableIpAddresses[i] = types.StringValue(allowableIpAddress)
	}

	var backupOption = searchengine.BackupOption{}
	if data.InitConfigOption.BackupOption.Get() != nil {
		backupOption = searchengine.BackupOption{
			RetentionPeriodDay: types.StringValue(data.InitConfigOption.BackupOption.Get().RetentionPeriodDay),
			StartingTimeHour:   types.StringValue(data.InitConfigOption.BackupOption.Get().StartingTimeHour),
		}
	}

	var initConfigOption = searchengine.InitConfigOptionResponse{
		BackupOption:     backupOption,
		DatabasePort:     types.Int32PointerValue(data.InitConfigOption.DatabasePort.Get()),
		DatabaseUserName: types.StringValue(data.InitConfigOption.DatabaseUserName),
		KibanaPort:       types.Int32PointerValue(data.InitConfigOption.KibanaPort.Get()),
		DashboardsPort:   types.Int32PointerValue(data.InitConfigOption.DashboardsPort.Get()),
	}

	var InstanceGroups []searchengine.InstanceGroup
	for _, instanceGroup := range data.InstanceGroups {
		var BlockStorage []searchengine.BlockStorageGroup
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			BlockStorage = append(BlockStorage, searchengine.BlockStorageGroup{
				Id:         types.StringValue(blockStorage.Id),
				Name:       types.StringValue(blockStorage.Name),
				RoleType:   types.StringValue(string(blockStorage.RoleType)),
				SizeGb:     types.Int32Value(blockStorage.SizeGb),
				VolumeType: types.StringValue(string(blockStorage.VolumeType)),
			})
		}

		var Instance []searchengine.Instance
		for _, instance := range instanceGroup.Instances {
			Instance = append(Instance, searchengine.Instance{
				Name:             types.StringValue(instance.Name),
				RoleType:         types.StringValue(string(instance.RoleType)),
				ServiceIpAddress: types.StringPointerValue(instance.ServiceIpAddress.Get()),
				PublicIpId:       types.StringPointerValue(instance.PublicIpId.Get()),
				//PublicIpAddress:  types.StringPointerValue(instance.PublicIpAddress.Get()),
				//ServiceState:     types.StringValue(string(instance.ServiceState)),
			})
		}

		InstanceGroups = append(InstanceGroups, searchengine.InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: BlockStorage,
			Instances:          Instance,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		})
	}

	var MaintenanceOption = searchengine.MaintenanceOption{
		PeriodHour:           types.StringPointerValue(data.MaintenanceOption.Get().PeriodHour.Get()),
		StartingDayOfWeek:    types.StringPointerValue((*string)(data.MaintenanceOption.Get().StartingDayOfWeek.Get())),
		StartingTime:         types.StringPointerValue(data.MaintenanceOption.Get().StartingTime.Get()),
		UseMaintenanceOption: types.BoolPointerValue(data.MaintenanceOption.Get().UseMaintenanceOption),
	}

	var License = searchengine.License{
		LicenseExpiryDate: types.StringValue(data.License.Get().LicenseExpiryDate),
		LicenseStatus:     types.StringValue(data.License.Get().LicenseStatus),
		LicenseType:       types.StringValue(data.License.Get().LicenseType),
	}

	var searchengineState = searchengine.ClusterDetail{
		AccountId:            types.StringValue(data.AccountId),
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngine:          types.StringValue(string(data.DbaasEngine)),
		IsCombined:           types.BoolPointerValue(data.IsCombined.Get()),
		Id:                   types.StringValue(data.Id),
		InitConfigOption:     initConfigOption,
		InstanceCount:        types.Int32PointerValue(data.InstanceCount),
		InstanceGroups:       InstanceGroups,
		License:              License,
		MaintenanceOption:    MaintenanceOption,
		Name:                 types.StringValue(data.Name),
		NatEnabled:           types.BoolPointerValue(data.NatEnabled.Get()),
		ProductType:          types.StringValue(string(data.ProductType)),
		ServiceState:         types.StringValue(string(data.ServiceState)),
		SoftwareVersion:      types.StringValue(data.SoftwareVersion),
		SubnetId:             types.StringValue(data.SubnetId),
		Timezone:             types.StringValue(data.Timezone),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
	}
	searchengineObjectValue, diags := types.ObjectValueFrom(ctx, searchengineState.AttributeTypes(), searchengineState)
	state.ClusterDetail = searchengineObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
