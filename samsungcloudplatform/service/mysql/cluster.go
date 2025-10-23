package mysql

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/mysql"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpMysql "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/mysql/1.0"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
	"time"
)

var (
	_ resource.Resource              = &mysqlClusterResource{}
	_ resource.ResourceWithConfigure = &mysqlClusterResource{}
)

func NewMysqlClusterResource() resource.Resource {
	return &mysqlClusterResource{}
}

type mysqlClusterResource struct {
	config  *scpsdk.Configuration
	client  *mysql.Client
	clients *client.SCPClient
}

func (r *mysqlClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_cluster"
}

func (r *mysqlClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "mysql",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AllowableIpAddresses"): schema.ListAttribute{
				Description: "Allowed IP addresses list  \n" +
					"  - example: ['192.168.10.1/32']",
				Required:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("DbaasEngineVersionId"): schema.StringAttribute{
				Description: "DBaaS engine version ID \n" +
					"  - example: '03cf36ee3162415e854c411e3690c85e' (MySQL Community 8.0.41)",
				Required: true,
			},
			common.ToSnakeCase("HaEnabled"): schema.BoolAttribute{
				Description: "HA availability \n" +
					"  - example: False \n",
				Required: true,
			},
			common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
				Description: "NAT availability \n" +
					"  - example: False \n",
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
							common.ToSnakeCase("ArchiveFrequencyMinute"): schema.StringAttribute{
								Description: "Backup starting time (minute) \n" +
									"  - example: 60 \n" +
									"  - pattern: 5 / 10 / 30 / 60 \n",
								Optional: true,
							},
							common.ToSnakeCase("RetentionPeriodDay"): schema.StringAttribute{
								Description: "Backup retention period (day) \n" +
									"  - example: 7 \n" +
									"  - min: 7 \n" +
									"  - max: 35 \n",
								Optional: true,
							},
							common.ToSnakeCase("StartingTimeHour"): schema.StringAttribute{
								Description: "Backup starting time (hour) \n" +
									"  - example: 12 \n" +
									"  - min: 00 \n" +
									"  - max: 23 \n",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("DatabaseCaseSensitive"): schema.BoolAttribute{
						Description: "Database case sensitivity \n" +
							"  - example: False \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseCharacterSet"): schema.StringAttribute{
						Description: "Database encoding \n" +
							"  - example: 'utf8mb3' \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
						Description: "Database name \n" +
							"  - example: 'test' \n" +
							"  - minLength: 3  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-zA-Z][a-zA-Z0-9]*$ \n",
						Required: true,
					},
					common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
						Description: "Database service port \n" +
							"  - example: 2866 \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserName"): schema.StringAttribute{
						Description: "Database user name \n" +
							"  - example: 'test' \n" +
							"  - minLength: 2  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-z]*$ \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserPassword"): schema.StringAttribute{
						Description: "Database user password \n" +
							"  - minLength: 8  \n" +
							"  - maxLength: 30  \n" +
							"  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30} (\"'제외) \n",
						Required: true,
					},
				},
			},
			common.ToSnakeCase("InstanceGroups"): schema.ListNestedAttribute{
				Description: "Instance groups",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("BlockStorageGroups"): schema.ListNestedAttribute{
							Description: "BlockStorage groups",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("Id"): schema.StringAttribute{
										Description: "Id",
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("Name"): schema.StringAttribute{
										Description: "Name",
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: "Role type \n" +
											"  - example: 'OS' \n",
										Required: true,
									},
									common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
										Description: "Size in GB \n" +
											"  - example: 104 \n" +
											"  - minLength: 16  \n" +
											"  - maxLength: 5120  \n",
										Required: true,
									},
									common.ToSnakeCase("VolumeType"): schema.StringAttribute{
										Description: "Volume type \n" +
											"  - example: 'SSD' \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("SSD", "SSD_KMS", "HDD", "HDD_KMS"),
										},
									},
								},
							},
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
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
										Description: "Name",
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: "Role type \n" +
											"  - example: 'ACTIVE' \n" +
											"  - pattern: ACTIVE / STANDBY \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("ACTIVE", "STANDBY"),
										},
									},
									common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
										Description: "User subnet IP address",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
										Description: "Public IP ID (Required when NatEnabled=True & HaEnabled=False)",
										Optional:    true,
									},
									//common.ToSnakeCase("PublicIpAddress"): schema.StringAttribute{
									//	Description: "Public IP address",
									//	Computed:    true,
									//	PlanModifiers: []planmodifier.String{
									//		stringplanmodifier.UseStateForUnknown(),
									//	},
									//},
								},
							},
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: "Role type \n" +
								"  - example: 'ACTIVE' \n" +
								"  - pattern: ACTIVE (HaEnabled=False) / ACTIVE_STANDBY (HaEnabled=True) \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ACTIVE", "ACTIVE_STANDBY"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: "Server type name \n" +
								"  - example: 'db1v2m4' \n",
							Required: true,
						},
					},
				},
			},
			common.ToSnakeCase("InstanceNamePrefix"): schema.StringAttribute{
				Description: "Instance name prefix \n" +
					"  - example: 'test'  \n" +
					"  - minLength: 3  \n" +
					"  - maxLength: 13  \n" +
					"  - pattern: ^[a-z][a-zA-Z0-9\\-]*$ \n",
				Required: true,
			},
			common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
				Description: "Maintenance option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
						Description: "Period in hours \n" +
							"  - example: 1  \n",
						Optional: true,
					},
					common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
						Description: "Starting day of week \n" +
							"  - example: 'MON' \n",
						Optional: true,
					},
					common.ToSnakeCase("StartingTime"): schema.StringAttribute{
						Description: "Starting time \n" +
							"  - example: '0000' \n",
						Optional: true,
					},
					common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
						Description: "Use maintenance option \n" +
							"  - example: False \n",
						Optional: true,
						Computed: true,
					},
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Cluster name \n" +
					"  - example: 'test'  \n" +
					"  - minLength: 3  \n" +
					"  - maxLength: 20  \n" +
					"  - pattern: ^[a-zA-Z]*$ \n",
				Required: true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description: "Service state \n" +
					"  - example : 'RUNNING' (Create,Start) / 'STOPPED' (Stop) \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID",
				Required:    true,
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description: "Timezone \n" +
					"  - example: 'Asia/Seoul' \n",
				Required: true,
			},
			common.ToSnakeCase("VipPublicIpId"): schema.StringAttribute{
				Description: "VIP Public IP ID (Required when NatEnabled=True & HaEnabled=True)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			//common.ToSnakeCase("VipPublicIpAddress"): schema.StringAttribute{
			//	Description: "Virtual IP address",
			//	Computed:    true,
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.UseStateForUnknown(),
			//	},
			//},
			common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
				Description: "Virtual IP address",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *mysqlClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Mysql
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *mysqlClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan mysql.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// cluster id 반환
	clusterId := data.Resource.Id

	// cluster 조회 func
	getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	//wait for 구현
	getData, err := databaseUtils.AsyncRequestPollingWithState(ctx, clusterId, 500, 10*time.Second,
		"ServiceState", "RUNNING", "FAILED", getFunc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Cluster",
			"Could not read Cluster, unexpected error: "+err.Error(),
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "mysql", "mysql", clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	if len(plan.Tags.Elements()) > 0 {
		getTags, err := r.AsyncPollingTags(ctx, clusterId, "mysql", "mysql",
			100, 3*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Tag",
				err.Error(),
			)
			return
		}
		tagsMap = getTags
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	//Metadata 처리
	state, err := r.MapGetResponseToState(ctx, getData, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mysqlClusterResource) AsyncPollingTags(ctx context.Context, clusterId string, serviceName string,
	resourceType string, maxAttempts int, internal time.Duration) (types.Map, error) {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, clusterId)

		if err != nil {
			return types.Map{}, fmt.Errorf("attempt %d/%d failed: %w",
				attempt, maxAttempts, err)
		}

		if len(tagsMap.Elements()) > 0 {
			return tagsMap, nil
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return types.Map{}, fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return types.Map{}, fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func (r *mysqlClusterResource) MapGetResponseToState(ctx context.Context, resp *scpMysql.MysqlClusterDetailResponse, plan mysql.ClusterResource, tagsMap types.Map) (mysql.ClusterResource, error) {

	var allowableIpAddresses types.List
	if len(resp.AllowableIpAddresses) == 0 {
		allowableIpAddresses, _ = types.ListValue(types.StringType, []attr.Value{})
	} else {
		ipAddresses := make([]attr.Value, len(resp.AllowableIpAddresses))
		for i, ipAddress := range resp.AllowableIpAddresses {
			ipAddresses[i] = types.StringValue(ipAddress)
		}
		allowableIpAddresses, _ = types.ListValue(types.StringType, ipAddresses)
	}

	var backupOption = mysql.BackupOption{}
	if resp.InitConfigOption.BackupOption.Get() != nil {
		backupOption = mysql.BackupOption{
			ArchiveFrequencyMinute: types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().ArchiveFrequencyMinute.Get()),
			RetentionPeriodDay:     types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().RetentionPeriodDay.Get()),
			StartingTimeHour:       types.StringPointerValue(resp.InitConfigOption.BackupOption.Get().StartingTimeHour.Get()),
		}
	}

	var initConfigOption = mysql.InitConfigOption{
		BackupOption:          backupOption,
		DatabaseCaseSensitive: plan.InitConfigOption.DatabaseCaseSensitive,
		DatabaseCharacterSet:  types.StringPointerValue(resp.InitConfigOption.DatabaseCharacterSet.Get()),
		DatabaseName:          types.StringValue(resp.InitConfigOption.DatabaseName),
		DatabasePort:          types.Int32PointerValue(resp.InitConfigOption.DatabasePort.Get()),
		DatabaseUserName:      types.StringValue(resp.InitConfigOption.DatabaseUserName),
		DatabaseUserPassword:  plan.InitConfigOption.DatabaseUserPassword,
	}

	var InstanceGroups []mysql.InstanceGroup
	for _, instanceGroup := range resp.InstanceGroups {
		var BlockStorage []mysql.BlockStorageGroup
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			BlockStorage = append(BlockStorage, mysql.BlockStorageGroup{
				Id:         types.StringValue(blockStorage.Id),
				Name:       types.StringValue(blockStorage.Name),
				RoleType:   types.StringValue(string(blockStorage.RoleType)),
				SizeGb:     types.Int32Value(blockStorage.SizeGb),
				VolumeType: types.StringValue(string(blockStorage.VolumeType)),
			})
		}

		var Instance []mysql.Instance
		for _, instance := range instanceGroup.Instances {
			Instance = append(Instance, mysql.Instance{
				Name:             types.StringValue(instance.Name),
				RoleType:         types.StringValue(string(instance.RoleType)),
				ServiceIpAddress: types.StringPointerValue(instance.ServiceIpAddress.Get()),
				PublicIpId:       types.StringPointerValue(instance.PublicIpId.Get()),
				//PublicIpAddress:  types.StringPointerValue(instance.PublicIpAddress.Get()),
				//ServiceState:     types.StringValue(string(instance.ServiceState)),
			})
		}

		InstanceGroups = append(InstanceGroups, mysql.InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: BlockStorage,
			Instances:          Instance,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		})
	}

	var maintenanceOption = mysql.MaintenanceOption{}
	if resp.MaintenanceOption.Get() != nil {
		maintenanceOption = mysql.MaintenanceOption{
			PeriodHour:           types.StringPointerValue(resp.MaintenanceOption.Get().PeriodHour.Get()),
			StartingDayOfWeek:    types.StringPointerValue((*string)(resp.MaintenanceOption.Get().StartingDayOfWeek.Get())),
			StartingTime:         types.StringPointerValue(resp.MaintenanceOption.Get().StartingTime.Get()),
			UseMaintenanceOption: types.BoolPointerValue(resp.MaintenanceOption.Get().UseMaintenanceOption),
		}
	}

	return mysql.ClusterResource{
		Id:                   types.StringValue(resp.Id),
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: plan.DbaasEngineVersionId,
		NatEnabled:           types.BoolPointerValue(resp.NatEnabled),
		HaEnabled:            types.BoolPointerValue(resp.HaEnabled),
		InitConfigOption:     initConfigOption,
		InstanceGroups:       InstanceGroups,
		InstanceNamePrefix:   plan.InstanceNamePrefix,
		MaintenanceOption:    maintenanceOption,
		Name:                 types.StringValue(resp.Name),
		ServiceState:         types.StringValue(string(resp.ServiceState)),
		SubnetId:             types.StringValue(resp.SubnetId),
		Tags:                 tagsMap,
		Timezone:             types.StringValue(resp.Timezone),
		VipPublicIpId:        types.StringValue(resp.GetVipPublicIpId()),
		//VipPublicIpAddress:   types.StringValue(resp.GetVipPublicIpAddress()),
		VirtualIpAddress: types.StringValue(resp.GetVirtualIpAddress()),
	}, nil
}

func (r *mysqlClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysql.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read Cluster name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "mysql", "mysql", state.Id.ValueString())
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

func (r *mysqlClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*mysql.UpdateHandler{
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

	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
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

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "HaEnabled", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "VipPublicIpId", "VirtualIpAddress"}

	if databaseUtils.IsOverlapFields(immutableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
		)
		return
	}

	// 변경 확인
	for _, h := range handlers {
		if databaseUtils.IsOverlapFields(h.Fields, changeFields) {
			if err := h.Handler(ctx, req, resp); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Cluster",
					"Could not update cluster, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	data, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading cluster",
			"Could not read cluster name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "mysql", "mysql", state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState, _ := r.MapGetResponseToState(ctx, data, plan, tagsMap)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *mysqlClusterResource) getStateTransitions() map[string]map[string]func(ctx context.Context, clusterId string) error {
	transitions := make(map[string]map[string]func(ctx context.Context, clusterId string) error)

	addState := func(from string, to string, callFunc func(ctx context.Context, clusterId string) error) {
		// from map 이 구성 되지 않았을때 초기화
		if transitions[from] == nil {
			transitions[from] = make(map[string]func(ctx context.Context, clusterId string) error)
		}
		transitions[from][to] = callFunc
	}

	// State Transition Map
	addState("STOPPED", "RUNNING", r.client.StartCluster)
	addState("RUNNING", "STOPPED", r.client.StopCluster)

	return transitions
}

func (r *mysqlClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	currentState := state.ServiceState.ValueString()
	desiredState := plan.ServiceState.ValueString()

	if currentState == desiredState {
		return nil
	}

	// state에 따라 start, stop 구분
	err := r.getStateTransitions()[currentState][desiredState](ctx, plan.Id.ValueString())
	if err != nil {
		return err
	}

	// wait for 구현
	getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", desiredState, "ERROR", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlClusterResource) handlerUpdateClusterInitConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	clusterId := plan.Id.ValueString()

	backupState := state.InitConfigOption.BackupOption
	backupPlan := plan.InitConfigOption.BackupOption

	// 1. backup 최초 설정
	if isEmpty(backupState) && !isEmpty(backupPlan) {
		archiveFrequencyMinute := backupPlan.ArchiveFrequencyMinute.ValueString()
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, archiveFrequencyMinute, startingTimeHour, retentionPeriodDay)
		if err != nil {
			return err
		}
	}

	// 2. backup 설정 변경
	if !isEmpty(backupState) && !isEmpty(backupPlan) && !reflect.DeepEqual(backupState, backupPlan) {
		archiveFrequencyMinute := backupPlan.ArchiveFrequencyMinute.ValueString()
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, archiveFrequencyMinute, startingTimeHour, retentionPeriodDay)
		if err != nil {
			return err
		}
	}

	// 3. backup 설정 삭제
	if !isEmpty(backupState) && isEmpty(backupPlan) {
		err := r.client.UnSetBackup(ctx, clusterId)
		if err != nil {
			return err
		}
	}

	// wait for 구현
	getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err := databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", "RUNNING", "ERROR", getFunc)
	if err != nil {
		return err
	}
	return nil
}

func isEmpty(sp mysql.BackupOption) bool {
	return sp.ArchiveFrequencyMinute.IsNull() && sp.StartingTimeHour.IsNull() && sp.RetentionPeriodDay.IsNull()
}

func (r *mysqlClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	clusterId := plan.Id.ValueString()

	ipState, _ := databaseUtils.ConvertListtoStringSlice(state.AllowableIpAddresses)
	ipPlan, _ := databaseUtils.ConvertListtoStringSlice(plan.AllowableIpAddresses)

	addedIPs, removedIps := databaseUtils.CompareIPAddresses(ipState, ipPlan)

	err := r.client.SetSecurityGroupRules(ctx, clusterId, addedIPs, removedIps)
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", "RUNNING", "FAILED", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func compareIPAddresses(state []types.String, plan []types.String) ([]string, []string) {
	toStringSlice := func(input []types.String) []string {
		var result []string
		for _, item := range input {
			if !item.IsNull() {
				result = append(result, item.ValueString())
			}
		}
		return result
	}

	toSet := func(items []string) map[string]bool {
		set := make(map[string]bool)
		for _, item := range items {
			set[item] = true
		}
		return set
	}

	stateIPSet := toSet(toStringSlice(state))
	planIPSet := toSet(toStringSlice(plan))

	diff := func(sourceSet map[string]bool, targetSet map[string]bool) []string {
		var result []string
		for key := range sourceSet {
			if !targetSet[key] {
				result = append(result, key)
			}
		}
		return result
	}

	addedIPs := diff(planIPSet, stateIPSet)
	removedIPs := diff(stateIPSet, planIPSet)

	return addedIPs, removedIPs
}

func (r *mysqlClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	for i := 0; i < len(plan.InstanceGroups); i++ {
		currentInstanceGroup := state.InstanceGroups[i]
		desiredInstanceGroup := plan.InstanceGroups[i]

		instanceGroupFields := []string{"BlockStorageGroups", "Id", "Instances", "RoleType", "ServerTypeName"}

		changedFields, err := databaseUtils.GetChangedFields(desiredInstanceGroup, currentInstanceGroup, instanceGroupFields)
		if err != nil {
			return err
		}

		immutableFields := []string{"Id", "RoleType"}

		if databaseUtils.IsOverlapFields(immutableFields, changedFields) {
			resp.Diagnostics.AddError(
				"Error Updating Cluster",
				"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
			)
			return nil
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
				if len(currentInstanceGroup.BlockStorageGroups) == len(desiredInstanceGroup.BlockStorageGroups) {
					// Resize Block Storage
					for i := 0; i < len(currentInstanceGroup.BlockStorageGroups); i++ {
						currentBlockStorage := currentInstanceGroup.BlockStorageGroups[i]
						desiredBlockStorage := desiredInstanceGroup.BlockStorageGroups[i]

						bsFields := []string{"Id", "Name", "RoleType", "SizeGb", "VolumeType"}
						changedBsFields, err := databaseUtils.GetChangedFields(currentBlockStorage, desiredBlockStorage, bsFields)
						if err != nil {
							return err
						}

						immutableBsFields := []string{"Id", "Name", "RoleType", "VolumeType"}

						if databaseUtils.IsOverlapFields(immutableBsFields, changedBsFields) {
							resp.Diagnostics.AddError(
								"Error Updating Cluster",
								"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
							)
							return nil
						}

						if databaseUtils.IsOverlapFields(changedBsFields, []string{"SizeGb"}) {
							//client
							err := r.client.SetBlockStorageSize(ctx, currentBlockStorage.Id.ValueString(), desiredBlockStorage.SizeGb.ValueInt32())
							if err != nil {
								return err
							}
						}
					}
				} else {
					// Add Block Storage
					addBlockStorage := desiredInstanceGroup.BlockStorageGroups[len(desiredInstanceGroup.BlockStorageGroups)-1]
					err := r.client.AddBlockStorages(ctx, currentInstanceGroup.Id.ValueString(), addBlockStorage.RoleType.ValueString(), addBlockStorage.SizeGb.ValueInt32(), addBlockStorage.VolumeType.ValueString())
					if err != nil {
						return err
					}
				}

			}

			// wait for 구현
			getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
				return r.client.GetCluster(ctx, id)
			}

			_, err := databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
				"ServiceState", "RUNNING", "ERROR", getFunc)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *mysqlClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan mysql.ClusterResource
	var state mysql.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	// Update
	_, err := tag.UpdateTags(r.clients, "mysql", "mysql", plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysql.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// cluster id 반환
	clusterId := state.Id.ValueString()

	// Delete existing cluster
	err := r.client.DeleteCluster(ctx, clusterId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// cluster 조회 func
	getFunc := func(id string) (*scpMysql.MysqlClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}
	// wait for 구현
	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, clusterId, 200, 20*time.Second,
		"ServiceState", "TERMINATED", "FAILED", getFunc)
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
