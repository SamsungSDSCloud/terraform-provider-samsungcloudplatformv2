package sqlserver

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpSqlserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/sqlserver/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &sqlserverClusterResource{}
	_ resource.ResourceWithConfigure   = &sqlserverClusterResource{}
	_ resource.ResourceWithImportState = &sqlserverClusterResource{}
)

func NewSqlserverClusterResource() resource.Resource {
	return &sqlserverClusterResource{}
}

type sqlserverClusterResource struct {
	config  *scpsdk.Configuration
	client  *sqlserver.Client
	clients *client.SCPClient
}

func (r *sqlserverClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sqlserver_cluster"
}

func (r *sqlserverClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "sqlserver",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
				MarkdownDescription: "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AllowableIpAddresses"): schema.SetAttribute{
				Description: databaseUtils.DescAllowedIPAddressesList +
					"  - example: ['192.168.10.1/32']",
				MarkdownDescription: databaseUtils.DescAllowedIPAddressesList +
					"  - example: ['192.168.10.1/32']",
				Required:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("DbaasEngineVersionId"): schema.StringAttribute{
				Description: databaseUtils.DescDBaaSEngineVersionID +
					"  - example: 'b15f062f45f54c3fbf28bfa14d87d1b8' (Microsoft SQL Server 2022 Enterprise ENG-KB5048033-x64)",
				MarkdownDescription: databaseUtils.DescDBaaSEngineVersionID +
					"  - example: 'b15f062f45f54c3fbf28bfa14d87d1b8' (Microsoft SQL Server 2022 Enterprise ENG-KB5048033-x64)",
				Required:  true,
				WriteOnly: true,
			},
			common.ToSnakeCase("HaEnabled"): schema.BoolAttribute{
				Description: databaseUtils.DescHAAvailability +
					databaseUtils.DescExampleFalse,
				MarkdownDescription: databaseUtils.DescHAAvailability +
					databaseUtils.DescExampleFalse,
				Required: true,
			},
			common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
				Description: databaseUtils.DescNATAvailability +
					databaseUtils.DescExampleFalse,
				MarkdownDescription: databaseUtils.DescNATAvailability +
					databaseUtils.DescExampleFalse,
				Required: true,
			},
			common.ToSnakeCase("OriginClusterId"): schema.StringAttribute{
				Description:         "Origin Cluster Id.\n  - example: 25e21d596d4f41e9b7b66d8f2129213b",
				MarkdownDescription: "Origin Cluster Id.\n  - example: 25e21d596d4f41e9b7b66d8f2129213b",
				Computed:            true,
			},
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AuditEnabled"): schema.BoolAttribute{
						Description:         "Audit Log Setting\n  - example: true",
						MarkdownDescription: "Audit Log Setting\n  - example: true",
						Required:            true,
					},
					common.ToSnakeCase("AdConfig"): schema.SingleNestedAttribute{
						Description: "AdConfig",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("AdDnsServers"): schema.SetAttribute{
								ElementType: types.StringType,
								Description: "AD DNS Servers\n  - example: 192.168.10.10",
								Optional:    true,
							},
							common.ToSnakeCase("AdDomainName"): schema.StringAttribute{
								Description: "AD Domain Name\n  - example: test",
								Optional:    true,
							},
							common.ToSnakeCase("AdNetbiosName"): schema.StringAttribute{
								Description: "AD NetBIOS Name\n  - example: test",
								Optional:    true,
							},
							common.ToSnakeCase("AdUserId"): schema.StringAttribute{
								Description: "AD User ID\n  - example: test",
								Optional:    true,
							},
							common.ToSnakeCase("AdUserPassword"): schema.StringAttribute{
								Description: "AD User Password \n" +
									"  - example: test123!",
								MarkdownDescription: "AD User Password \n" +
									"  - example: test123!",
								Optional:  true,
								WriteOnly: true,
							},
							common.ToSnakeCase("FailoverClusterName"): schema.StringAttribute{
								Description:         "Failover Cluster Name\n  - example: testcluster",
								MarkdownDescription: "Failover Cluster Name\n  - example: testcluster",
								Optional:            true,
							},
						},
					},
					common.ToSnakeCase("AdEnabled"): schema.BoolAttribute{
						Description:         "AdEnabled\n  - example: false",
						MarkdownDescription: "AdEnabled\n  - example: false",
						Required:            true,
					},
					common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
						Description: "Backup option",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ArchiveFrequencyMinute"): schema.StringAttribute{
								Description: databaseUtils.DescBackupStartingTimeMinute +
									databaseUtils.DescExample60 +
									databaseUtils.DescPatternMinuteOptions,
								MarkdownDescription: databaseUtils.DescBackupStartingTimeMinute +
									databaseUtils.DescExample60 +
									databaseUtils.DescPatternMinuteOptions,
								Optional: true,
							},
							common.ToSnakeCase("FullBackupDayOfWeek"): schema.StringAttribute{
								Description: "Days of the week when full backups are performed \n" +
									"  - allowed values: 'MON','TUE','WED','THU','FRI','SAT','SUN' \n" +
									"  - default : 'SUN'  \n",
								MarkdownDescription: "Days of the week when full backups are performed \n" +
									"  - allowed values: 'MON','TUE','WED','THU','FRI','SAT','SUN' \n" +
									"  - default : 'SUN'  \n",
								Optional: true,
							},
							common.ToSnakeCase("RetentionPeriodDay"): schema.StringAttribute{
								Description: databaseUtils.DescBackupRetentionPeriodDay +
									databaseUtils.DescExample7 +
									databaseUtils.DescMin7 +
									databaseUtils.DescMax35,
								MarkdownDescription: databaseUtils.DescBackupRetentionPeriodDay +
									databaseUtils.DescExample7 +
									databaseUtils.DescMin7 +
									databaseUtils.DescMax35,
								Optional: true,
							},
							common.ToSnakeCase("StartingTimeHour"): schema.StringAttribute{
								Description: databaseUtils.DescBackupStartingTimeHour +
									databaseUtils.DescExample12 +
									databaseUtils.DescMin00 +
									databaseUtils.DescMax23,
								MarkdownDescription: databaseUtils.DescBackupStartingTimeHour +
									databaseUtils.DescExample12 +
									databaseUtils.DescMin00 +
									databaseUtils.DescMax23,
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("DatabaseCollation"): schema.StringAttribute{
						Description: "Database collation \n" +
							"  - allowed values: 'SQL_Latin1_General_CP1_CI_AS','Korean_Wansung_CS_AS','Chinese_PRC_CI_AS'",
						MarkdownDescription: "Database collation \n" +
							"  - allowed values: 'SQL_Latin1_General_CP1_CI_AS','Korean_Wansung_CS_AS','Chinese_PRC_CI_AS'",
						Required: true,
					},
					common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
						Description: databaseUtils.DescDatabaseServicePort +
							databaseUtils.DescExample2866,
						MarkdownDescription: databaseUtils.DescDatabaseServicePort +
							databaseUtils.DescExample2866,
						Required: true,
					},
					common.ToSnakeCase("DatabaseServiceName"): schema.StringAttribute{
						Description: "Database service name \n" +
							"  - example: 'Test' \n" +
							"  - minLength: 1  \n" +
							"  - maxLength: 15  \n" +
							"  - pattern: ^[A-Z][a-zA-Z]*$ \n",
						MarkdownDescription: "Database service name \n" +
							"  - example: 'Test' \n" +
							"  - minLength: 1  \n" +
							"  - maxLength: 15  \n" +
							"  - pattern: ^[A-Z][a-zA-Z]*$ \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserName"): schema.StringAttribute{
						Description: databaseUtils.DescDatabaseUserName +
							databaseUtils.DescExampleTest2 +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							databaseUtils.DescPatternLowerAlpha,
						MarkdownDescription: databaseUtils.DescDatabaseUserName +
							databaseUtils.DescExampleTest2 +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							databaseUtils.DescPatternLowerAlpha,
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserPassword"): schema.StringAttribute{
						Description: databaseUtils.DescDatabaseUserPassword +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPasswordEng,
						MarkdownDescription: databaseUtils.DescDatabaseUserPassword +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPasswordEng,
						Required:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("Databases"): schema.ListNestedAttribute{
						Description: "Databases",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
									Description: databaseUtils.DescDatabaseName +
										databaseUtils.DescExampleTest2 +
										databaseUtils.DescMinLength3 +
										databaseUtils.DescMaxLength20 +
										databaseUtils.DescPatternAlphaAlnum,
									MarkdownDescription: databaseUtils.DescDatabaseName +
										databaseUtils.DescExampleTest2 +
										databaseUtils.DescMinLength3 +
										databaseUtils.DescMaxLength20 +
										databaseUtils.DescPatternAlphaAlnum,
									Required: true,
								},
								common.ToSnakeCase("DriveLetter"): schema.StringAttribute{
									Description:         "Drive Letter\n  - example: C",
									MarkdownDescription: "Drive Letter\n  - example: C",
									Required:            true,
								},
							},
						},
					},
					common.ToSnakeCase("License"): schema.StringAttribute{
						Description:         "License\n  - example: license",
						MarkdownDescription: "License\n  - example: license",
						Required:            true,
						WriteOnly:           true,
					},
					common.ToSnakeCase("OriginRegion"): schema.StringAttribute{
						Description:         "Origin Region\n -example: kr-west1",
						MarkdownDescription: "Origin Region\n -example: kr-west1",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			common.ToSnakeCase("InstanceGroups"): schema.ListNestedAttribute{
				Description: "Instance groups",
				Required:    true,
				PlanModifiers: []planmodifier.List{
					databaseUtils.InstanceGroupsPlanModifier(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("BlockStorageGroups"): schema.ListNestedAttribute{
							Description: "BlockStorage groups",
							Required:    true,
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
										Description: databaseUtils.DescRoleType +
											databaseUtils.DescExampleOS,
										MarkdownDescription: databaseUtils.DescRoleType +
											databaseUtils.DescExampleOS,
										Required: true,
									},
									common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
										Description: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											databaseUtils.DescMinLength16 +
											databaseUtils.DescMaxLength5120,
										MarkdownDescription: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											databaseUtils.DescMinLength16 +
											databaseUtils.DescMaxLength5120,
										Required: true,
									},
									common.ToSnakeCase("VolumeType"): schema.StringAttribute{
										Description: databaseUtils.DescVolumeType +
											databaseUtils.DescExampleSSD,
										MarkdownDescription: databaseUtils.DescVolumeType +
											databaseUtils.DescExampleSSD,
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("SSD", "SSD_KMS", "HDD", "HDD_KMS"),
										},
									},
								},
							},
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Instance group ID.\n  - example: ee48b333d5a84097adc079dec17ab872",
							MarkdownDescription: "Instance group ID.\n  - example: ee48b333d5a84097adc079dec17ab872",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						common.ToSnakeCase("Instances"): schema.ListNestedAttribute{
							Description: "Instances",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("Name"): schema.StringAttribute{
										Description:         "Instance name\n  - example: test001",
										MarkdownDescription: "Instance name\n  - example: test001",
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: databaseUtils.DescRoleType +
											"  - example: 'PRIMARY' \n" +
											"  - pattern: ACTIVE / PRIMARY / SECONDARY \n",
										MarkdownDescription: databaseUtils.DescRoleType +
											"  - example: 'PRIMARY' \n" +
											"  - pattern: ACTIVE / PRIMARY / SECONDARY \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("ACTIVE", "PRIMARY", "SECONDARY"),
										},
									},
									common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
										Description:         "User subnet IP address\n  - example: 192.168.4.22",
										MarkdownDescription: "User subnet IP address\n  - example: 192.168.4.22",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
										Description:         "Public IP ID (Required when NatEnabled=True & HaEnabled=False)\n  - example: 90a68b14850741598ecacd0eb190873e",
										MarkdownDescription: "Public IP ID (Required when NatEnabled=True & HaEnabled=False)\n  - example: 90a68b14850741598ecacd0eb190873e",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: databaseUtils.DescRoleType +
								databaseUtils.DescExampleACTIVE +
								"  - pattern: ACTIVE (HaEnabled=False) / PRIMARY_SECONDARY (HaEnabled=True) \n",
							MarkdownDescription: databaseUtils.DescRoleType +
								databaseUtils.DescExampleACTIVE +
								"  - pattern: ACTIVE (HaEnabled=False) / PRIMARY_SECONDARY (HaEnabled=True) \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ACTIVE", "PRIMARY_SECONDARY"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: databaseUtils.DescServerTypeName +
								databaseUtils.DescExampleDb1v2m4,
							MarkdownDescription: databaseUtils.DescServerTypeName +
								databaseUtils.DescExampleDb1v2m4,
							Required: true,
						},
					},
				},
			},
			common.ToSnakeCase("InstanceNamePrefix"): schema.StringAttribute{
				Description: databaseUtils.DescInstanceNamePrefix +
					databaseUtils.DescExampleTest +
					databaseUtils.DescMinLength3 +
					databaseUtils.DescMaxLength13 +
					databaseUtils.DescPatternLowerAlnumDash,
				MarkdownDescription: databaseUtils.DescInstanceNamePrefix +
					databaseUtils.DescExampleTest +
					databaseUtils.DescMinLength3 +
					databaseUtils.DescMaxLength13 +
					databaseUtils.DescPatternLowerAlnumDash,
				Required:  true,
				WriteOnly: true,
			},
			common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
				Description: "Maintenance option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
						Description: databaseUtils.DescPeriodInHours +
							databaseUtils.DescExample1,
						MarkdownDescription: databaseUtils.DescPeriodInHours +
							databaseUtils.DescExample1,
						Optional: true,
					},
					common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
						Description: databaseUtils.DescStartingDayOfWeek +
							databaseUtils.DescExampleMON,
						MarkdownDescription: databaseUtils.DescStartingDayOfWeek +
							databaseUtils.DescExampleMON,
						Optional: true,
					},
					common.ToSnakeCase("StartingTime"): schema.StringAttribute{
						Description: databaseUtils.DescStartingTime +
							databaseUtils.DescExample0000,
						MarkdownDescription: databaseUtils.DescStartingTime +
							databaseUtils.DescExample0000,
						Optional: true,
					},
					common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
						Description: databaseUtils.DescUseMaintenanceOption +
							databaseUtils.DescExampleFalse,
						MarkdownDescription: databaseUtils.DescUseMaintenanceOption +
							databaseUtils.DescExampleFalse,
						Optional: true,
						Computed: true,
					},
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: databaseUtils.DescClusterName +
					databaseUtils.DescExampleTest +
					databaseUtils.DescMinLength3 +
					databaseUtils.DescMaxLength20 +
					databaseUtils.DescPatternAlpha,
				MarkdownDescription: databaseUtils.DescClusterName +
					databaseUtils.DescExampleTest +
					databaseUtils.DescMinLength3 +
					databaseUtils.DescMaxLength20 +
					databaseUtils.DescPatternAlpha,
				Required: true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description: databaseUtils.DescServiceState +
					databaseUtils.DescExampleRunningStopped,
				MarkdownDescription: databaseUtils.DescServiceState +
					databaseUtils.DescExampleRunningStopped,
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description:         "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249",
				MarkdownDescription: "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249",
				Required:            true,
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description: databaseUtils.DescTimezone +
					databaseUtils.DescExampleAsiaSeoul,
				MarkdownDescription: databaseUtils.DescTimezone +
					databaseUtils.DescExampleAsiaSeoul,
				Required: true,
			},
			common.ToSnakeCase("VipPublicIpId"): schema.StringAttribute{
				Description:         "VIP Public IP ID (Required when NatEnabled=True & HaEnabled=True)\n  - example: 88a68b14850741599ecacd0eb190999a",
				MarkdownDescription: "VIP Public IP ID (Required when NatEnabled=True & HaEnabled=True)\n  - example: 88a68b14850741599ecacd0eb190999a",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					databaseUtils.ImmutableString(),
				},
			},
			common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
				Description:         "Virtual IP address\n  - example: 192.168.4.30",
				MarkdownDescription: "Virtual IP address\n  - example: 192.168.4.30",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					databaseUtils.ImmutableString(),
				},
			},
			common.ToSnakeCase("ServiceWatchLogCollection"): schema.BoolAttribute{
				Description:         "ServiceWatchLogCollection\n - example: false",
				MarkdownDescription: "ServiceWatchLogCollection\n - example: false",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					databaseUtils.ImmutableBool(),
				},
			},
		},
	}
}

func (r *sqlserverClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Sqlserver
	r.clients = inst.Client
}

func (r *sqlserverClusterResource) nullOutWriteOnlyFields(plan *sqlserver.ClusterResource) {
	plan.DbaasEngineVersionId = types.StringNull()
	plan.InstanceNamePrefix = types.StringNull()
	if plan.InitConfigOption != nil {
		plan.InitConfigOption.DatabaseUserPassword = types.StringNull()
		plan.InitConfigOption.License = types.StringNull()
		plan.InitConfigOption.AdConfig.AdUserPassword = types.StringNull()
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *sqlserverClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sqlserver.ClusterResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		r.nullOutWriteOnlyFields(&plan)

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// cluster id 반환
	clusterId := data.Resource.Id

	// Save state immediately after creation to prevent orphan resources
	plan.Id = types.StringValue(clusterId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"CREATING"}, []string{"RUNNING"}, true)
	if err != nil {
		r.nullOutWriteOnlyFields(&plan)
		resp.State.Set(ctx, plan)

		resp.Diagnostics.AddError(
			"Error waiting for Cluster",
			"Cluster was created but failed to become RUNNING: "+err.Error(),
		)
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *sqlserverClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpSqlserver.SqlserverClusterDetailResponseV1Dot1, plan sqlserver.ClusterResource, tagsMap types.Map) (sqlserver.ClusterResource, error) {

	var allowableIpAddresses types.Set
	if len(resp.AllowableIpAddresses) == 0 {
		allowableIpAddresses, _ = types.SetValue(types.StringType, []attr.Value{})
	} else {
		ipAddresses := make([]attr.Value, len(resp.AllowableIpAddresses))
		for i, ipAddress := range resp.AllowableIpAddresses {
			ipAddresses[i] = types.StringValue(ipAddress)
		}
		allowableIpAddresses, _ = types.SetValue(types.StringType, ipAddresses)
	}

	var initConfigOption *sqlserver.InitConfigOption
	{
		var dbUserPassword types.String
		var adUserPassword types.String
		var license types.String
		if plan.InitConfigOption != nil {
			dbUserPassword = plan.InitConfigOption.DatabaseUserPassword
			adUserPassword = plan.InitConfigOption.AdConfig.AdUserPassword
			license = plan.InitConfigOption.License
		} else {
			dbUserPassword = types.StringNull()
			adUserPassword = types.StringNull()
			license = types.StringNull()
		}

		var backupOption = sqlserver.BackupOption{}
		if resp.InitConfigOption.BackupOption.Get() != nil {
			backupOption = sqlserver.BackupOption{
				ArchiveFrequencyMinute: types.StringValue(resp.InitConfigOption.BackupOption.Get().ArchiveFrequencyMinute),
				FullBackupDayOfWeek:    types.StringPointerValue((*string)(resp.InitConfigOption.BackupOption.Get().FullBackupDayOfWeek)),
				RetentionPeriodDay:     types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().RetentionPeriodDay.Get()),
				StartingTimeHour:       types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().StartingTimeHour.Get()),
			}
		}

		var adDnsServers types.Set
		var adConfig = sqlserver.AdConfig{}
		// AdConfig가 없을 때도 ad_dns_servers set이 element type을 갖도록 기본값을 설정한다.
		adConfig.AdDnsServers = types.SetNull(types.StringType)
		if resp.InitConfigOption.AdConfig.Get() != nil {
			adCfg := resp.InitConfigOption.AdConfig.Get()
			adDnsServerList := make([]attr.Value, len(adCfg.AdDnsServers))
			for i, adDnsServer := range adCfg.AdDnsServers {
				adDnsServerList[i] = types.StringValue(adDnsServer)
			}
			adDnsServers, _ = types.SetValue(types.StringType, adDnsServerList)

			adConfig = sqlserver.AdConfig{
				AdConfigBase: sqlserver.AdConfigBase{
					AdDnsServers:        adDnsServers,
					AdDomainName:        types.StringPointerValue(adCfg.AdDomainName.Get()),
					AdNetbiosName:       types.StringPointerValue(adCfg.AdNetbiosName.Get()),
					AdUserId:            types.StringPointerValue(adCfg.AdUserId.Get()),
					FailoverClusterName: types.StringPointerValue(adCfg.FailoverClusterName.Get()),
				},
				AdUserPassword: adUserPassword,
			}
		}

		dataBases := make([]sqlserver.Databases, len(resp.InitConfigOption.Databases))
		for i, database := range resp.InitConfigOption.Databases {
			dataBases[i] = sqlserver.Databases{DatabaseName: types.StringValue(database.GetDatabaseName()), DriveLetter: types.StringValue(database.GetDriveLetter())}
		}

		initConfigOption = &sqlserver.InitConfigOption{
			InitConfigOptionBase: sqlserver.InitConfigOptionBase{
				AdEnabled:           types.BoolPointerValue(resp.InitConfigOption.AdEnabled),
				AuditEnabled:        types.BoolPointerValue(resp.InitConfigOption.AuditEnabled),
				BackupOption:        backupOption,
				DatabaseCollation:   types.StringPointerValue((*string)(resp.InitConfigOption.DatabaseCollation)),
				DatabasePort:        types.Int32PointerValue(resp.InitConfigOption.DatabasePort),
				DatabaseServiceName: types.StringValue(resp.InitConfigOption.DatabaseServiceName),
				DatabaseUserName:    types.StringValue(resp.InitConfigOption.DatabaseUserName),
				Databases:           dataBases,
				OriginRegion:        types.StringPointerValue(resp.InitConfigOption.OriginRegion.Get()),
			},
			AdConfig:             adConfig,
			DatabaseUserPassword: dbUserPassword,
			License:              license,
		}

		instanceGroupsList := databaseUtils.MapInstanceGroupsList(ctx, plan.InstanceGroups, sqlserver.MapInstanceGroupResponses(resp.InstanceGroups))

		var maintenanceOption *sqlserver.MaintenanceOption
		if resp.MaintenanceOption.IsSet() && resp.MaintenanceOption.Get() != nil {
			maintenanceOption = &sqlserver.MaintenanceOption{
				PeriodHour:           types.StringPointerValue(resp.MaintenanceOption.Get().PeriodHour.Get()),
				StartingDayOfWeek:    types.StringPointerValue((*string)(resp.MaintenanceOption.Get().StartingDayOfWeek.Get())),
				StartingTime:         types.StringPointerValue(resp.MaintenanceOption.Get().StartingTime.Get()),
				UseMaintenanceOption: types.BoolPointerValue(resp.MaintenanceOption.Get().UseMaintenanceOption),
			}
		} else {
			maintenanceOption = plan.MaintenanceOption
		}

		var originClusterId string
		if resp.OriginClusterId.Get() != nil {
			originClusterId = *resp.OriginClusterId.Get()
		}
		return sqlserver.ClusterResource{
			Id:                   types.StringValue(resp.Id),
			AllowableIpAddresses: allowableIpAddresses,
			DbaasEngineVersionId: plan.DbaasEngineVersionId,
			NatEnabled:           types.BoolPointerValue(resp.NatEnabled),
			HaEnabled:            types.BoolPointerValue(resp.HaEnabled),
			InitConfigOption:     initConfigOption,
			InstanceGroups:       instanceGroupsList,
			InstanceNamePrefix:   plan.InstanceNamePrefix,
			MaintenanceOption:    maintenanceOption,
			OriginClusterId:      types.StringValue(originClusterId),
			Name:                 types.StringValue(resp.Name),
			ServiceState:         types.StringValue(string(resp.ServiceState)),
			SubnetId:             types.StringValue(resp.SubnetId),
			Tags:                 tagsMap,
			Timezone:             types.StringValue(resp.Timezone),
			VipPublicIpId:        types.StringValue(resp.GetVipPublicIpId()),
			//VipPublicIpAddress:   types.StringValue(resp.GetVipPublicIpAddress()),
			VirtualIpAddress:          types.StringValue(resp.GetVirtualIpAddress()),
			ServiceWatchLogCollection: types.BoolValue(resp.GetServiceWatchLogCollection()),
		}, nil
	}
}

func (r *sqlserverClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sqlserver.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read Cluster name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "sqlserver", "sqlserver", state.Id.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, state, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sqlserverClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*sqlserver.UpdateHandler{
		{
			Fields:  []string{"ServiceState"},
			Handler: r.handlerUpdateClusterState,
		},
		{
			Fields:  []string{"InitConfigOption"},
			Handler: r.handlerUpdateClusterInitConfig,
		},
		{
			Fields:  []string{"AllowableIpAddresses"},
			Handler: r.handlerUpdateClusterAllowableIpAddresses,
		},
		{
			Fields:  []string{"InstanceGroups"},
			Handler: r.handlerUpdateInstanceGroups,
		},
		{
			Fields:  []string{"Tags"},
			Handler: r.handlerUpdateTag,
		},
	}

	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settableFields []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFields = append(settableFields, databaseUtils.SnakeToPascal(attrName))
		}
	}

	changeFields, err := databaseUtils.GetChangedFields(plan, state, settableFields)
	if err != nil {
		return
	}

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "HaEnabled", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "VipPublicIpId", "VirtualIpAddress", "ServiceWatchLogCollection"}

	// InitConfigOption is immutable except for BackupOption: guard it only when a
	// field other than BackupOption changed.
	initConfigOnlyBackup := plan.InitConfigOption != nil && state.InitConfigOption != nil &&
		databaseUtils.OnlyBackupOptionChanged(*plan.InitConfigOption, *state.InitConfigOption)
	if !initConfigOnlyBackup {
		immutableFields = append(immutableFields, "InitConfigOption")
	}

	// Reject changes to immutable fields, reporting only the fields actually changed.
	if violated := databaseUtils.OverlapFields(immutableFields, changeFields); len(violated) > 0 {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Immutable fields cannot be modified: "+strings.Join(violated, ", "),
		)
		return
	}

	// Dispatch each handler whose fields changed.
	for _, h := range handlers {
		if !databaseUtils.IsOverlapFields(h.Fields, changeFields) {
			continue
		}
		if err := h.Handler(ctx, req, resp); err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Cluster",
				"Could not update cluster, unexpected error: "+err.Error(),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *sqlserverClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	currentState := state.ServiceState.ValueString()
	desiredState := plan.ServiceState.ValueString()

	if currentState == desiredState {
		return nil
	}

	// 현재 상태가 전이 중(STOPPING/STARTING)이면 별도 명령 없이 종료 상태가 될 때까지 대기한다.
	currentState, settleErr := databaseUtils.WaitForSettledState(ctx, currentState, plan.Id.ValueString(),
		func(ctx context.Context, clusterId string, pendingStates, targetStates []string) error {
			return waitForClusterStatus(ctx, r.client, clusterId, pendingStates, targetStates, true)
		})
	if settleErr != nil {
		return settleErr
	}

	// 전이 대기 후 이미 목표 상태에 도달했으면 종료한다.
	if currentState == desiredState {
		return nil
	}

	// state에 따라 start, stop 구분
	transition, ok := databaseUtils.GetStateTransitions(r.client)[currentState][desiredState]
	if !ok || transition == nil {
		return fmt.Errorf("unsupported service_state transition: %q -> %q (allowed transitions: STOPPED->RUNNING, RUNNING->STOPPED)", currentState, desiredState)
	}

	err := transition(ctx, plan.Id.ValueString())
	if err != nil {
		return err
	}

	pendingStates := databaseUtils.GetPendingStates(currentState)
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), pendingStates, []string{desiredState}, true)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlserverClusterResource) handlerUpdateClusterInitConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	clusterId := plan.Id.ValueString()

	var backupState sqlserver.BackupOption
	if state.InitConfigOption != nil {
		backupState = state.InitConfigOption.BackupOption
	}
	var backupPlan sqlserver.BackupOption
	if plan.InitConfigOption != nil {
		backupPlan = plan.InitConfigOption.BackupOption
	}

	// 1. backup 최초 설정
	if isEmpty(backupState) && !isEmpty(backupPlan) {
		archiveFrequencyMinute := backupPlan.ArchiveFrequencyMinute.ValueString()
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()
		fullBackupDayOfWeek := backupPlan.FullBackupDayOfWeek.ValueString()

		err := r.client.SetBackup(ctx, clusterId, archiveFrequencyMinute, startingTimeHour, retentionPeriodDay, fullBackupDayOfWeek)
		if err != nil {
			return err
		}
	}

	// 2. backup 설정 변경
	if !isEmpty(backupState) && !isEmpty(backupPlan) && !reflect.DeepEqual(backupState, backupPlan) {
		archiveFrequencyMinute := backupPlan.ArchiveFrequencyMinute.ValueString()
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()
		fullBackupDayOfWeek := backupPlan.FullBackupDayOfWeek.ValueString()

		err := r.client.SetBackup(ctx, clusterId, archiveFrequencyMinute, startingTimeHour, retentionPeriodDay, fullBackupDayOfWeek)
		if err != nil {
			return err
		}
	}

	// 3. bacup 설정 삭제
	if !isEmpty(backupState) && isEmpty(backupPlan) {
		err := r.client.UnSetBackup(ctx, clusterId)
		if err != nil {
			return err
		}
	}

	err := waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true)
	if err != nil {
		return err
	}

	return nil
}

func isEmpty(sp sqlserver.BackupOption) bool {
	return sp.ArchiveFrequencyMinute.IsNull() && sp.StartingTimeHour.IsNull() && sp.RetentionPeriodDay.IsNull() && sp.FullBackupDayOfWeek.IsNull()
}

func (r *sqlserverClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	clusterId := plan.Id.ValueString()

	addedIPs, removedIps := databaseUtils.CompareIPAddresses(state.AllowableIpAddresses, plan.AllowableIpAddresses)

	err := r.client.SetSecurityGroupRules(ctx, clusterId, addedIPs, removedIps)
	if err != nil {
		return err
	}

	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlserverClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	var planIGs []databaseUtils.InstanceGroup
	plan.InstanceGroups.ElementsAs(ctx, &planIGs, false)
	var stateIGs []databaseUtils.InstanceGroup
	state.InstanceGroups.ElementsAs(ctx, &stateIGs, false)

	stateIGByRole := make(map[string]databaseUtils.InstanceGroup, len(stateIGs))
	for _, sg := range stateIGs {
		if _, exists := stateIGByRole[sg.RoleType.ValueString()]; !exists {
			stateIGByRole[sg.RoleType.ValueString()] = sg
		}
	}

	for i := 0; i < len(planIGs); i++ {
		desiredInstanceGroup := planIGs[i]
		currentInstanceGroup, matched := stateIGByRole[desiredInstanceGroup.RoleType.ValueString()]
		if !matched {
			continue
		}

		instanceGroupFields := []string{"Instances", "RoleType", "ServerTypeName"}

		changedFields, err := databaseUtils.GetChangedFields(desiredInstanceGroup, currentInstanceGroup, instanceGroupFields)
		if err != nil {
			return err
		}

		var currentBS []databaseUtils.BlockStorageGroup
		currentInstanceGroup.BlockStorageGroups.ElementsAs(ctx, &currentBS, false)
		var desiredBS []databaseUtils.BlockStorageGroup
		desiredInstanceGroup.BlockStorageGroups.ElementsAs(ctx, &desiredBS, false)

		if !reflect.DeepEqual(currentBS, desiredBS) {
			changedFields = append(changedFields, "BlockStorageGroups")
		}

		immutableFields := []string{"RoleType"}

		if databaseUtils.IsOverlapFields(immutableFields, changedFields) {
			return fmt.Errorf("immutable fields cannot be modified: %s", strings.Join(immutableFields, ", "))
		}

		if databaseUtils.IsOverlapFields([]string{"Instances"}, changedFields) {
			return fmt.Errorf("operation not permitted for INSTANCE type: modifying instances is not supported")
		}

		if len(changedFields) > 0 {
			// ServerTypeName Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"ServerTypeName"}) {
				err := r.client.SetServerType(ctx, currentInstanceGroup.Id.ValueString(), desiredInstanceGroup.ServerTypeName.ValueString())
				if err != nil {
					return err
				}
			}

			// BlockStorageGroups Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"BlockStorageGroups"}) {

				bsPlan, err := databaseUtils.PlanBlockStorageUpdate(currentBS, desiredBS)
				if err != nil {
					return err
				}
				if len(bsPlan.Removed) > 0 {
					return fmt.Errorf("operation not permitted for BLOCK_STORAGE_GROUP type: removing an existing block storage is not supported")
				}

				// Resize existing Block Storages
				for _, resize := range bsPlan.Resizes {
					if err := r.client.SetBlockStorageSize(ctx, resize.Id, resize.SizeGb); err != nil {
						return err
					}
				}
				if len(bsPlan.Resizes) > 0 {
					if err := waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true); err != nil {
						return err
					}
				}

				// Add new Block Storages
				for _, add := range bsPlan.Adds {
					if err := r.client.AddBlockStorages(ctx, currentInstanceGroup.Id.ValueString(), add.RoleType.ValueString(), add.SizeGb.ValueInt32(), add.VolumeType.ValueString()); err != nil {
						return err
					}
				}
				if len(bsPlan.Adds) > 0 {
					if err := waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true); err != nil {
						return err
					}
				}

			}
			// wait for 구현
			err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *sqlserverClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan sqlserver.ClusterResource
	var state sqlserver.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	// Update
	_, err := tag.UpdateTags(r.clients, "sqlserver", "sqlserver", plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlserverClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sqlserver.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// cluster id 반환
	clusterId := state.Id.ValueString()

	// Delete cluster
	err := r.client.DeleteCluster(ctx, clusterId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"TERMINATING"}, []string{"TERMINATED"}, false)
	if err != nil {
		if err.Error() != "cluster with id="+clusterId+" not found" {
			resp.Diagnostics.AddError(
				"Error deleting cluster",
				"Could not delete cluster, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *sqlserverClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForClusterStatus(ctx context.Context, sqlClient *sqlserver.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := sqlClient.GetCluster(ctx, id)
		if httpStatus == 200 {
			currentState := string(info.ServiceState)
			for _, s := range pendingStates {
				if s == currentState {
					return info, currentState, nil
				}
			}
			for _, s := range targetStates {
				if s == currentState {
					return info, currentState, nil
				}
			}
			return nil, "", fmt.Errorf("cluster with id=%s transitioned to unexpected state: %s", id, currentState)
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}
			return info, "TERMINATED", nil
		} else if err != nil {
			return nil, "", err
		}
		return info, string(info.ServiceState), nil
	}, -1, -1, -1, -1)
}
