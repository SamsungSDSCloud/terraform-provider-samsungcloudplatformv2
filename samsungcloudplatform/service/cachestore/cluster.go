package cachestore

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cachestore"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpCachestore "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cachestore/1.2"
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
	_ resource.Resource                = &cachestoreClusterResource{}
	_ resource.ResourceWithConfigure   = &cachestoreClusterResource{}
	_ resource.ResourceWithImportState = &cachestoreClusterResource{}
)

// Reusable description fragments to avoid duplicated string literals.
const (
	descExampleMaster = "  - example: 'MASTER' \n"
)

func NewCachestoreClusterResource() resource.Resource {
	return &cachestoreClusterResource{}
}

type cachestoreClusterResource struct {
	config  *scpsdk.Configuration
	client  *cachestore.Client
	clients *client.SCPClient
}

func (r *cachestoreClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cachestore_cluster"
}

func (r *cachestoreClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "cachestore",
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
					"  - example: 'aef8e9ace6f54207bdf6266d4028cb74' (Redis OSS Sentinel 7.2.6)",
				MarkdownDescription: databaseUtils.DescDBaaSEngineVersionID +
					"  - example: 'aef8e9ace6f54207bdf6266d4028cb74' (Redis OSS Sentinel 7.2.6)",
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
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
						Description: "Backup option",
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
							databaseUtils.DescExample2866,
						MarkdownDescription: databaseUtils.DescDatabaseServicePort +
							databaseUtils.DescExample2866,
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserPassword"): schema.StringAttribute{
						Description: databaseUtils.DescDatabaseUserPassword +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							"  - pattern: '^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30}' (\"'$제외) \n",
						MarkdownDescription: databaseUtils.DescDatabaseUserPassword +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							"  - pattern: '^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30}' (\"'$제외) \n",
						Required:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("SentinelPort"): schema.Int32Attribute{
						Description: "Sentinel port \n" +
							"  - example: 26378 \n",
						MarkdownDescription: "Sentinel port \n" +
							"  - example: 26378 \n",
						Required: true,
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
											stringvalidator.OneOf(databaseUtils.BSRoleTypesOsData...),
										},
									},
									common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
										Description: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											"  - minLength: 56  \n" +
											databaseUtils.DescMaxLength5120,
										MarkdownDescription: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											"  - minLength: 56  \n" +
											databaseUtils.DescMaxLength5120,
										Required: true,
									},
									common.ToSnakeCase("VolumeType"): schema.StringAttribute{
										Description:         "Volume type \n  - example: 'SSD' \n",
										MarkdownDescription: "Volume type \n  - example: 'SSD' \n",
										Required:            true,
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
										Description:         "Cluster name\n  - example: mytest",
										MarkdownDescription: "Cluster name\n  - example: mytest",
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: databaseUtils.DescRoleType +
											descExampleMaster +
											"  - pattern: MASTER / REPLICA / SENTINEL \n",
										MarkdownDescription: databaseUtils.DescRoleType +
											descExampleMaster +
											"  - pattern: MASTER / REPLICA / SENTINEL \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("MASTER", "REPLICA", "SENTINEL"),
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
										Description:         "Public IP ID (Required when NatEnabled=True)",
										MarkdownDescription: "Public IP ID (Required when NatEnabled=True)",
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
								descExampleMaster +
								"  - pattern: MASTER / MASTER_REPLICA / SENTINEL \n",
							MarkdownDescription: databaseUtils.DescRoleType +
								descExampleMaster +
								"  - pattern: MASTER / MASTER_REPLICA / SENTINEL \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("MASTER", "MASTER_REPLICA", "SENTINEL"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description:         "Server type name \n  - example: 'redis1v1m2' (Redis) / 'css1v1m2' (Valkey) \n",
							MarkdownDescription: "Server type name \n  - example: 'redis1v1m2' (Redis) / 'css1v1m2' (Valkey) \n",
							Required:            true,
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
				Validators: []validator.Object{
					databaseUtils.MaintenanceOptionValidator(),
				},
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
						Description:         "Period in hours \n  - example: 1  \n",
						MarkdownDescription: "Period in hours \n  - example: 1  \n",
						Optional:            true,
					},
					common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
						Description:         "Starting day of week \n  - example: 'MON' \n",
						MarkdownDescription: "Starting day of week \n  - example: 'MON' \n",
						Optional:            true,
					},
					common.ToSnakeCase("StartingTime"): schema.StringAttribute{
						Description:         "Starting time \n  - example: '0000' \n",
						MarkdownDescription: "Starting time \n  - example: '0000' \n",
						Optional:            true,
					},
					common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
						Description:         "Use maintenance option \n  - example: False \n",
						MarkdownDescription: "Use maintenance option \n  - example: False \n",
						Optional:            true,
						Computed:            true,
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
			common.ToSnakeCase("ReplicaCount"): schema.Int32Attribute{
				Description:         "Replica count \n  - example: 0  \n",
				MarkdownDescription: "Replica count \n  - example: 0  \n",
				Required:            true,
				WriteOnly:           true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description:         "Service state \n  - example : 'RUNNING' (Create,Start) / 'STOPPED' (Stop) \n",
				MarkdownDescription: "Service state \n  - example : 'RUNNING' (Create,Start) / 'STOPPED' (Stop) \n",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description:         "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249 \n",
				MarkdownDescription: "Subnet ID\n  - example: 0c6d633730a9470c9cb3c66be1bc9249 \n",
				Required:            true,
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description:         "Timezone \n  - example: 'Asia/Seoul' \n",
				MarkdownDescription: "Timezone \n  - example: 'Asia/Seoul' \n",
				Required:            true,
			},
			common.ToSnakeCase("ServiceWatchLogCollection"): schema.BoolAttribute{
				Description:         "ServiceWatchLogCollection\n  - example: false  \n",
				MarkdownDescription: "ServiceWatchLogCollection\n  - example: false  \n",
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

func (r *cachestoreClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Cachestore
	r.clients = inst.Client
}

func (r *cachestoreClusterResource) nullOutWriteOnlyFields(plan *cachestore.ClusterResource) {
	plan.DbaasEngineVersionId = types.StringNull()
	plan.InstanceNamePrefix = types.StringNull()
	plan.ReplicaCount = types.Int32Null()
	if plan.InitConfigOption != nil {
		plan.InitConfigOption.DatabaseUserPassword = types.StringNull()
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *cachestoreClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cachestore.ClusterResource
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

	// Wait for cluster to become RUNNING
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

	// Call Read to refresh state with actual values
	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *cachestoreClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpCachestore.RedisClusterDetailResponseV1Dot1, plan cachestore.ClusterResource, tagsMap types.Map) (cachestore.ClusterResource, error) {

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

	var initConfigOption *cachestore.InitConfigOption
	{
		var dbUserPassword types.String
		if plan.InitConfigOption != nil {
			dbUserPassword = plan.InitConfigOption.DatabaseUserPassword
		} else {
			dbUserPassword = types.StringNull()
		}

		var backupOption = cachestore.BackupOption{}
		if resp.InitConfigOption.BackupOption.Get() != nil {
			backupOption = cachestore.BackupOption{
				RetentionPeriodDay: types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().RetentionPeriodDay.Get()),
				StartingTimeHour:   types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().StartingTimeHour.Get()),
			}
		}

		initConfigOption = &cachestore.InitConfigOption{
			InitConfigOptionBase: cachestore.InitConfigOptionBase{
				BackupOption: backupOption,
				DatabasePort: types.Int32PointerValue(resp.InitConfigOption.DatabasePort.Get()),
				SentinelPort: types.Int32PointerValue(resp.InitConfigOption.SentinelPort.Get()),
			},
			DatabaseUserPassword: dbUserPassword,
		}
	}

	instanceGroupsList := databaseUtils.MapInstanceGroupsList(ctx, plan.InstanceGroups, cachestore.MapInstanceGroupResponses(resp.InstanceGroups))

	var maintenanceOption *cachestore.MaintenanceOption
	if resp.MaintenanceOption.IsSet() && resp.MaintenanceOption.Get() != nil {
		maintenanceOption = &cachestore.MaintenanceOption{
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

	return cachestore.ClusterResource{
		Id:                        types.StringValue(resp.Id),
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      plan.DbaasEngineVersionId,
		HaEnabled:                 types.BoolPointerValue(resp.HaEnabled),
		InitConfigOption:          initConfigOption,
		InstanceGroups:            instanceGroupsList,
		InstanceNamePrefix:        plan.InstanceNamePrefix,
		MaintenanceOption:         maintenanceOption,
		Name:                      types.StringValue(resp.Name),
		NatEnabled:                types.BoolPointerValue(resp.NatEnabled),
		ReplicaCount:              plan.ReplicaCount,
		ServiceState:              types.StringValue(string(resp.ServiceState)),
		SubnetId:                  types.StringValue(resp.SubnetId),
		Tags:                      tagsMap,
		Timezone:                  types.StringValue(resp.Timezone),
		ServiceWatchLogCollection: types.BoolValue(resp.GetServiceWatchLogCollection()),
	}, nil
}

func (r *cachestoreClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cachestore.ClusterResource
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
	tagsMap, err := tag.GetTags(r.clients, "cachestore", "cache-store", state.Id.ValueString(), false)
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

func (r *cachestoreClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*cachestore.UpdateHandler{
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

	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
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

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "HaEnabled", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "ServiceWatchLogCollection"}

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

func (r *cachestoreClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
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

func (r *cachestoreClusterResource) handlerUpdateClusterInitConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	clusterId := plan.Id.ValueString()

	var backupState cachestore.BackupOption
	if state.InitConfigOption != nil {
		backupState = state.InitConfigOption.BackupOption
	}
	var backupPlan cachestore.BackupOption
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

func isEmpty(sp cachestore.BackupOption) bool {
	return sp.StartingTimeHour.IsNull() && sp.RetentionPeriodDay.IsNull()
}

func (r *cachestoreClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
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

func (r *cachestoreClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
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

	// Match instance groups to prior state by role_type rather than by list position,
	// so reordering (or adding) a group does not misattribute a change to the group that
	// shifted, and appending a group cannot index past stateIGs.
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
			// New instance group with no prior state; nothing to update in place.
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
				if len(bsPlan.Adds) > 0 || len(bsPlan.Removed) > 0 {
					return fmt.Errorf("operation not permitted for BLOCK_STORAGE_GROUP type: the cachestore product does not support the addition of storage, so the addition of block storage in instances is restricted")
				}
				for _, resize := range bsPlan.Resizes {
					if err := r.client.SetBlockStorageSize(ctx, resize.Id, resize.SizeGb); err != nil {
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

func (r *cachestoreClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan cachestore.ClusterResource
	var state cachestore.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	// Update
	_, err := tag.UpdateTags(r.clients, "cachestore", "cache-store", plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}

	return nil
}

func (r *cachestoreClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cachestore.ClusterResource
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

func (r *cachestoreClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForClusterStatus(ctx context.Context, csClient *cachestore.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := csClient.GetCluster(ctx, id)
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
