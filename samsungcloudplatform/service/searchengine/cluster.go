package searchengine

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpSearchengine "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/searchengine/1.2"
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
	_ resource.Resource                = &searchengineClusterResource{}
	_ resource.ResourceWithConfigure   = &searchengineClusterResource{}
	_ resource.ResourceWithImportState = &searchengineClusterResource{}
)

// Reusable description fragments to avoid duplicated string literals.
const (
	descExampleMasterData = "  - example: 'MASTER_DATA' \n"
)

func NewSearchengineClusterResource() resource.Resource {
	return &searchengineClusterResource{}
}

type searchengineClusterResource struct {
	config  *scpsdk.Configuration
	client  *searchengine.Client
	clients *client.SCPClient
}

func (r *searchengineClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_searchengine_cluster"
}

func (r *searchengineClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "searchengine",
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
					"  - example: 'b9b9a5547321a65ed980a3dceab2e10b' (Elasticsearch Enterprise 8.15.0)",
				MarkdownDescription: databaseUtils.DescDBaaSEngineVersionID +
					"  - example: 'b9b9a5547321a65ed980a3dceab2e10b' (Elasticsearch Enterprise 8.15.0)",
				Required:  true,
				WriteOnly: true,
			},
			common.ToSnakeCase("IsCombined"): schema.BoolAttribute{
				Description:         "MASTER,DATA combined (IsCombined=true), MASTER,DATA seperated (IsCombined=False)\n  - example: true",
				MarkdownDescription: "MASTER,DATA combined (IsCombined=true), MASTER,DATA seperated (IsCombined=False)\n  - example: true",
				Required:            true,
			},
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
						Description: "BackupOption",
						Required:    true,
						Attributes: map[string]schema.Attribute{
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
					common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
						Description: databaseUtils.DescDatabaseServicePort +
							"  - example: 9201 \n",
						MarkdownDescription: databaseUtils.DescDatabaseServicePort +
							"  - example: 9201 \n",
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
							databaseUtils.DescPatternPassword,
						MarkdownDescription: databaseUtils.DescDatabaseUserPassword +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						Required:  true,
						WriteOnly: true,
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
										Validators: []validator.String{
											stringvalidator.OneOf(databaseUtils.BSRoleTypesOsDataBackup...),
										},
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
											databaseUtils.ComputedNullToUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: databaseUtils.DescRoleType +
											descExampleMasterData +
											"  - pattern: MASTER_DATA / MASTER / DATA / KIBANA / DASHBOARDS \n",
										MarkdownDescription: databaseUtils.DescRoleType +
											descExampleMasterData +
											"  - pattern: MASTER_DATA / MASTER / DATA / KIBANA / DASHBOARDS \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("MASTER_DATA", "MASTER", "DATA", "KIBANA", "DASHBOARDS"),
										},
									},
									common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
										Description:         "User subnet IP address\n  - example: 192.168.4.22",
										MarkdownDescription: "User subnet IP address\n  - example: 192.168.4.22",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
											databaseUtils.ComputedNullToUnknown(),
										},
									},
									common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
										Description:         "Public IP ID (Required when NatEnabled=True & HaEnabled=False)\n  - example: 90a68b14850741598ecacd0eb190873e",
										MarkdownDescription: "Public IP ID (Required when NatEnabled=True & HaEnabled=False)\n  - example: 90a68b14850741598ecacd0eb190873e",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
											databaseUtils.ComputedNullToUnknown(),
										},
									},
								},
							},
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: databaseUtils.DescRoleType +
								descExampleMasterData +
								"  - pattern: MASTER_DATA (IsCombined=True) / MASTER, DATA (IsCombined=False) / KIBANA, DASHBOARDS (required)\n",
							MarkdownDescription: databaseUtils.DescRoleType +
								descExampleMasterData +
								"  - pattern: MASTER_DATA (IsCombined=True) / MASTER, DATA (IsCombined=False) / KIBANA, DASHBOARDS (required)\n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("MASTER_DATA", "MASTER", "DATA", "KIBANA", "DASHBOARDS"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: databaseUtils.DescServerTypeName +
								"  - example: 'se1v2m4' \n",
							MarkdownDescription: databaseUtils.DescServerTypeName +
								"  - example: 'se1v2m4' \n",
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
			common.ToSnakeCase("License"): schema.StringAttribute{
				Description:         "License\n  - example: license",
				MarkdownDescription: "License\n  - example: license",
				Optional:            true,
				WriteOnly:           true,
			},
			common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
				Description: "MaintenanceOption",
				Required:    true,
				Validators: []validator.Object{
					databaseUtils.MaintenanceOptionValidator(),
				},
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
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseNonNullStateForUnknown(),
							databaseUtils.ImmutableBool(),
						},
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
			common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
				Description: databaseUtils.DescNATAvailability +
					databaseUtils.DescExampleFalse,
				MarkdownDescription: databaseUtils.DescNATAvailability +
					databaseUtils.DescExampleFalse,
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

func (r *searchengineClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Searchengine
	r.clients = inst.Client
}

func (r *searchengineClusterResource) nullOutWriteOnlyFields(plan *searchengine.ClusterResource) {
	plan.DbaasEngineVersionId = types.StringNull()
	plan.InstanceNamePrefix = types.StringNull()
	plan.License = types.StringNull()
	if plan.InitConfigOption != nil {
		plan.InitConfigOption.DatabaseUserPassword = types.StringNull()
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *searchengineClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchengine.ClusterResource
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

func (r *searchengineClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpSearchengine.SearchEngineClusterDetailResponseV1Dot1, plan searchengine.ClusterResource, tagsMap types.Map) (searchengine.ClusterResource, error) {

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

	var initConfigOption *searchengine.InitConfigOption
	{
		var dbUserPassword types.String
		if plan.InitConfigOption != nil {
			dbUserPassword = plan.InitConfigOption.DatabaseUserPassword
		} else {
			dbUserPassword = types.StringNull()
		}

		var backupOption = searchengine.BackupOption{}
		if resp.InitConfigOption.BackupOption.Get() != nil {
			backupOption = searchengine.BackupOption{
				RetentionPeriodDay: types.StringValue(resp.InitConfigOption.BackupOption.Get().RetentionPeriodDay),
				StartingTimeHour:   types.StringValue(resp.InitConfigOption.BackupOption.Get().StartingTimeHour),
			}
		}

		initConfigOption = &searchengine.InitConfigOption{
			InitConfigOptionBase: searchengine.InitConfigOptionBase{
				BackupOption:     backupOption,
				DatabasePort:     types.Int32PointerValue(resp.InitConfigOption.DatabasePort.Get()),
				DatabaseUserName: types.StringValue(resp.InitConfigOption.DatabaseUserName),
			},
			DatabaseUserPassword: dbUserPassword,
		}
	}

	instanceGroupsList := databaseUtils.MapInstanceGroupsList(ctx, plan.InstanceGroups, searchengine.MapInstanceGroupResponses(resp.InstanceGroups))

	var maintenanceOption *searchengine.MaintenanceOption
	if resp.MaintenanceOption.IsSet() && resp.MaintenanceOption.Get() != nil {
		maintenanceOption = &searchengine.MaintenanceOption{
			PeriodHour:           types.StringPointerValue(resp.MaintenanceOption.Get().PeriodHour.Get()),
			StartingDayOfWeek:    types.StringPointerValue((*string)(resp.MaintenanceOption.Get().StartingDayOfWeek.Get())),
			StartingTime:         types.StringPointerValue(resp.MaintenanceOption.Get().StartingTime.Get()),
			UseMaintenanceOption: types.BoolPointerValue(resp.MaintenanceOption.Get().UseMaintenanceOption),
		}
	} else {
		// cluster가 failed 상태이면 API가 maintenance_option을 null로 반환하므로,
		// 응답으로 덮어쓰지 않고 직전 plan/state 값을 유지한다.
		maintenanceOption = plan.MaintenanceOption
	}

	return searchengine.ClusterResource{
		Id:                        types.StringValue(resp.Id),
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      plan.DbaasEngineVersionId,
		InitConfigOption:          initConfigOption,
		InstanceGroups:            instanceGroupsList,
		InstanceNamePrefix:        plan.InstanceNamePrefix,
		IsCombined:                types.BoolValue(resp.GetIsCombined()),
		MaintenanceOption:         maintenanceOption,
		License:                   plan.License,
		Name:                      types.StringValue(resp.Name),
		NatEnabled:                types.BoolValue(resp.GetNatEnabled()),
		ServiceState:              types.StringValue(string(resp.ServiceState)),
		SubnetId:                  types.StringValue(resp.SubnetId),
		Tags:                      tagsMap,
		Timezone:                  types.StringValue(resp.Timezone),
		ServiceWatchLogCollection: types.BoolValue(resp.GetServiceWatchLogCollection()),
	}, nil
}

func (r *searchengineClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state searchengine.ClusterResource
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
	tagsMap, err := tag.GetTags(r.clients, "searchengine", "search-engine", state.Id.ValueString(), false)
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

func (r *searchengineClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*searchengine.UpdateHandler{
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

	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
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

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "IsCombined", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "VipPublicIpId", "VirtualIpAddress", "ServiceWatchLogCollection"}

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

func (r *searchengineClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
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

func (r *searchengineClusterResource) handlerUpdateClusterInitConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	clusterId := plan.Id.ValueString()

	var backupState searchengine.BackupOption
	if state.InitConfigOption != nil {
		backupState = state.InitConfigOption.BackupOption
	}
	var backupPlan searchengine.BackupOption
	if plan.InitConfigOption != nil {
		backupPlan = plan.InitConfigOption.BackupOption
	}

	// 1. backup 최초 설정
	if isEmpty(backupState) && !isEmpty(backupPlan) {
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, startingTimeHour, retentionPeriodDay)
		if err != nil {
			return err
		}
	}

	// 2. backup 설정 변경
	if !isEmpty(backupState) && !isEmpty(backupPlan) && !reflect.DeepEqual(backupState, backupPlan) {
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, startingTimeHour, retentionPeriodDay)
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

func isEmpty(sp searchengine.BackupOption) bool {
	return sp.StartingTimeHour.IsNull() && sp.RetentionPeriodDay.IsNull()
}

func (r *searchengineClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
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

func (r *searchengineClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
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
				// 리사이즈/추가를 실행하기 전에 role_type을 먼저 검증한다.
				// 루프 중간에서 실패하면 이미 만들어진 디스크가 state에 남지 않아 drift가 생긴다.
				for _, add := range bsPlan.Adds {
					// AddBlockStorages는 ExtraBlockStorageGroupRoleType(OS 제외)만 허용한다.
					if !databaseUtils.IsExtraBlockStorageRoleType(add.RoleType.ValueString()) {
						return fmt.Errorf("operation not permitted for BLOCK_STORAGE_GROUP type: role_type %q cannot be added to an existing instance group (allowed: %s)", add.RoleType.ValueString(), strings.Join(databaseUtils.BSRoleTypesExtra, ", "))
					}
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

			// Instances Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"Instances"}) {
				// Kibana or DASHBOARDS
				t := currentInstanceGroup.RoleType.ValueString()
				if t == "KIBANA" || t == "DASHBOARDS" {
					resp.Diagnostics.AddError(
						"Invalid Instance Group Type",
						fmt.Sprintf("Instance group of type '%s' does not support  adding instance", t),
					)
					return nil
				}

				var currentIt []databaseUtils.Instance
				currentInstanceGroup.Instances.ElementsAs(ctx, &currentIt, false)
				var desiredIt []databaseUtils.Instance
				desiredInstanceGroup.Instances.ElementsAs(ctx, &desiredIt, false)

				instancePlan := databaseUtils.PlanInstanceUpdate(currentIt, desiredIt)
				if len(instancePlan.Removed) > 0 {
					return fmt.Errorf("operation not permitted for INSTANCE type: removing an existing instance is not supported")
				}
				if len(instancePlan.Adds) > 0 {
					instanceCount := int32(len(instancePlan.Adds))

					var serviceIPAddresses []string

					for _, instance := range instancePlan.Adds {
						if instance.ServiceIpAddress.IsNull() || instance.ServiceIpAddress.IsUnknown() {
							serviceIPAddresses = []string{}
							break
						}

						ip := instance.ServiceIpAddress.ValueString()
						serviceIPAddresses = append(serviceIPAddresses, ip)
					}

					err := r.client.AddInstances(ctx, state.Id.ValueString(), instanceCount, serviceIPAddresses)
					if err != nil {
						return err
					}
				}
			}

			// wait for 구현
			err := waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *searchengineClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan searchengine.ClusterResource
	var state searchengine.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	// Update
	_, err := tag.UpdateTags(r.clients, "searchengine", "search-engine", plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}

	return nil
}

func (r *searchengineClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state searchengine.ClusterResource
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

	// wait for 구현
	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"TERMINATING"}, []string{"TERMINATED"}, false)
	if err != nil {
		if err.Error() != "404 Not Found" {
			resp.Diagnostics.AddError(
				"Error reading server",
				"Could not read server, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *searchengineClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForClusterStatus(ctx context.Context, seClient *searchengine.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := seClient.GetCluster(ctx, id)
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
