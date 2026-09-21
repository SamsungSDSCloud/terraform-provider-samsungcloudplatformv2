package eventstreams

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpEventstreams "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/eventstreams/1.2"
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
	_ resource.Resource                = &eventstreamsClusterResource{}
	_ resource.ResourceWithConfigure   = &eventstreamsClusterResource{}
	_ resource.ResourceWithImportState = &eventstreamsClusterResource{}
)

// Reusable description fragments to avoid duplicated string literals.
const (
	descExampleZookeeperBroker = "  - example: 'ZOOKEEPER_BROKER' \n"
	descPatternLowerAsc        = "  - pattern: ^[a-z]+$ \n"
)

func NewEventstreamsClusterResource() resource.Resource {
	return &eventstreamsClusterResource{}
}

type eventstreamsClusterResource struct {
	config  *scpsdk.Configuration
	client  *eventstreams.Client
	clients *client.SCPClient
}

func (r *eventstreamsClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventstreams_cluster"
}

func (r *eventstreamsClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "eventstreams",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
				MarkdownDescription: "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AkhqEnabled"): schema.BoolAttribute{
				Description:         "AHKQ Enabled\n- example: false",
				MarkdownDescription: "AHKQ Enabled\n- example: false",
				Required:            true,
				WriteOnly:           true,
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
					"  - example: '189299a34f464cac94a24f2d8d57afec' (Kafka 3.8.0)",
				MarkdownDescription: databaseUtils.DescDBaaSEngineVersionID +
					"  - example: '189299a34f464cac94a24f2d8d57afec' (Kafka 3.8.0)",
				Required:  true,
				WriteOnly: true,
			},
			common.ToSnakeCase("IsCombined"): schema.BoolAttribute{
				Description:         "ZOOKEEPER,BROKER combined (IsCombined=true), ZOOKEEPER,BROKER seperated (IsCombined=False) \n- example: false",
				MarkdownDescription: "ZOOKEEPER,BROKER combined (IsCombined=true), ZOOKEEPER,BROKER seperated (IsCombined=False) \n- example: false",
				Required:            true,
			},
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AkhqId"): schema.StringAttribute{
						Description: "Akhq ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						MarkdownDescription: "Akhq ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						Optional:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("AkhqPassword"): schema.StringAttribute{
						Description: "Akhq password \n" +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						MarkdownDescription: "Akhq password password \n" +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						Optional:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("BrokerPort"): schema.Int32Attribute{
						Description: "Broker port \n" +
							"  - example: 9091 \n",
						MarkdownDescription: "Broker port \n" +
							"  - example: 9091 \n",
						Required: true,
					},
					common.ToSnakeCase("BrokerSaslId"): schema.StringAttribute{
						Description: "Broker Sasl ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						MarkdownDescription: "Broker Sasl ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						Required:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("BrokerSaslPassword"): schema.StringAttribute{
						Description: "Broker Sasl password \n" +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						MarkdownDescription: "Broker Sasl password \n" +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						Required:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("ZookeeperPort"): schema.Int32Attribute{
						Description: "Zookeeper port \n" +
							"  - example: 2180 \n",
						MarkdownDescription: "Zookeeper port \n" +
							"  - example: 2180 \n",
						Required: true,
					},
					common.ToSnakeCase("ZookeeperSaslId"): schema.StringAttribute{
						Description: "Zookeeper Sasl ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						MarkdownDescription: "Zookeeper Sasl ID \n" +
							databaseUtils.DescMinLength2 +
							databaseUtils.DescMaxLength20 +
							descPatternLowerAsc,
						Required:  true,
						WriteOnly: true,
					},
					common.ToSnakeCase("ZookeeperSaslPassword"): schema.StringAttribute{
						Description: "Zookeeper Sasl password \n" +
							databaseUtils.DescMinLength8 +
							databaseUtils.DescMaxLength30 +
							databaseUtils.DescPatternPassword,
						MarkdownDescription: "Zookeeper Sasl password \n" +
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
										Description:         "Block storage role type\n  - example: OS",
										MarkdownDescription: "Block storage role type\n  - example: OS",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(databaseUtils.BSRoleTypesOsData...),
										},
									},
									common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
										Description: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											databaseUtils.DescMinLength16 +
											databaseUtils.DescMaxLength5120 +
											"  - example: 104",
										MarkdownDescription: databaseUtils.DescSizeInGB +
											databaseUtils.DescExample104 +
											databaseUtils.DescMinLength16 +
											databaseUtils.DescMaxLength5120 +
											"  - example: 104",
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
											descExampleZookeeperBroker +
											"  - pattern: ZOOKEEPER_BROKER / ZOOKEEPER / BROKER / AKHQ \n",
										MarkdownDescription: databaseUtils.DescRoleType +
											descExampleZookeeperBroker +
											"  - pattern: ZOOKEEPER_BROKER / ZOOKEEPER / BROKER / AKHQ \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("ZOOKEEPER_BROKER", "ZOOKEEPER", "BROKER", "AKHQ"),
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
										Description:         "Public IP ID (Required when NatEnabled=True)\n  - example: 90a68b14850741598ecacd0eb190873e",
										MarkdownDescription: "Public IP ID (Required when NatEnabled=True)\n  - example: 90a68b14850741598ecacd0eb190873e",
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
								descExampleZookeeperBroker +
								"  - pattern: ZOOKEEPER_BROKER (IsCombined=True) / ZOOKEEPER, BROKER (IsCombined=False) / AKHQ (optional) \n",
							MarkdownDescription: databaseUtils.DescRoleType +
								descExampleZookeeperBroker +
								"  - pattern: ZOOKEEPER_BROKER (IsCombined=True) / ZOOKEEPER, BROKER (IsCombined=False) / AKHQ (optional) \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ZOOKEEPER_BROKER", "ZOOKEEPER", "BROKER", "AKHQ"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: databaseUtils.DescServerTypeName +
								"  - example: 'es1v2m4' \n",
							MarkdownDescription: databaseUtils.DescServerTypeName +
								"  - example: 'es1v2m4' \n",
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
				Description:         "Timezone\n  - example: Asia/Seoul",
				MarkdownDescription: "Timezone\n  - example: Asia/Seoul",
				Required:            true,
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

func (r *eventstreamsClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Eventstreams
	r.clients = inst.Client
}

func (r *eventstreamsClusterResource) nullOutWriteOnlyFields(plan *eventstreams.ClusterResource) {
	plan.DbaasEngineVersionId = types.StringNull()
	plan.InstanceNamePrefix = types.StringNull()
	plan.AkhqEnabled = types.BoolNull()
	if plan.InitConfigOption != nil {
		plan.InitConfigOption.AkhqId = types.StringNull()
		plan.InitConfigOption.AkhqPassword = types.StringNull()
		plan.InitConfigOption.BrokerSaslId = types.StringNull()
		plan.InitConfigOption.BrokerSaslPassword = types.StringNull()
		plan.InitConfigOption.ZookeeperSaslId = types.StringNull()
		plan.InitConfigOption.ZookeeperSaslPassword = types.StringNull()
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *eventstreamsClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan eventstreams.ClusterResource
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

func (r *eventstreamsClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, plan eventstreams.ClusterResource, tagsMap types.Map) (eventstreams.ClusterResource, error) {

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

	var initConfigOption *eventstreams.InitConfigOption
	{
		var akhqId types.String
		var akhqPassword types.String
		var brokerSaslId types.String
		var brokerSaslPassword types.String
		var zookeeperSaslId types.String
		var zookeeperSaslPassword types.String
		if plan.InitConfigOption != nil {
			akhqId = plan.InitConfigOption.AkhqId
			akhqPassword = plan.InitConfigOption.AkhqPassword
			brokerSaslId = plan.InitConfigOption.BrokerSaslId
			brokerSaslPassword = plan.InitConfigOption.BrokerSaslPassword
			zookeeperSaslId = plan.InitConfigOption.ZookeeperSaslId
			zookeeperSaslPassword = plan.InitConfigOption.ZookeeperSaslPassword
		} else {
			akhqId = types.StringNull()
			akhqPassword = types.StringNull()
			brokerSaslId = types.StringNull()
			brokerSaslPassword = types.StringNull()
			zookeeperSaslId = types.StringNull()
			zookeeperSaslPassword = types.StringNull()
		}

		initConfigOption = &eventstreams.InitConfigOption{
			InitConfigOptionBase: eventstreams.InitConfigOptionBase{
				BrokerPort:    types.Int32PointerValue(resp.InitConfigOption.BrokerPort),
				ZookeeperPort: types.Int32PointerValue(resp.InitConfigOption.ZookeeperPort),
			},
			AkhqId:                akhqId,
			AkhqPassword:          akhqPassword,
			BrokerSaslId:          brokerSaslId,
			BrokerSaslPassword:    brokerSaslPassword,
			ZookeeperSaslId:       zookeeperSaslId,
			ZookeeperSaslPassword: zookeeperSaslPassword,
		}
	}

	instanceGroupsList := databaseUtils.MapInstanceGroupsList(ctx, plan.InstanceGroups, eventstreams.MapInstanceGroupResponses(resp.InstanceGroups))

	var maintenanceOption *eventstreams.MaintenanceOption
	if resp.MaintenanceOption.IsSet() && resp.MaintenanceOption.Get() != nil {
		maintenanceOption = &eventstreams.MaintenanceOption{
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

	return eventstreams.ClusterResource{
		Id:                        types.StringValue(resp.Id),
		AkhqEnabled:               plan.AkhqEnabled,
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      plan.DbaasEngineVersionId,
		InitConfigOption:          initConfigOption,
		InstanceGroups:            instanceGroupsList,
		InstanceNamePrefix:        plan.InstanceNamePrefix,
		IsCombined:                types.BoolPointerValue(resp.IsCombined.Get()),
		MaintenanceOption:         maintenanceOption,
		Name:                      types.StringValue(resp.Name),
		NatEnabled:                types.BoolPointerValue(resp.NatEnabled.Get()),
		ServiceState:              types.StringValue(string(resp.ServiceState)),
		SubnetId:                  types.StringValue(resp.SubnetId),
		Tags:                      tagsMap,
		Timezone:                  types.StringValue(resp.Timezone),
		ServiceWatchLogCollection: types.BoolPointerValue(resp.ServiceWatchLogCollection),
	}, nil
}

func (r *eventstreamsClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventstreams.ClusterResource
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
	tagsMap, err := tag.GetTags(r.clients, "eventstreams", "event-streams", state.Id.ValueString(), false)
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

func (r *eventstreamsClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*eventstreams.UpdateHandler{
		{
			Fields:  []string{"ServiceState"},
			Handler: r.handlerUpdateClusterState,
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

	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
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

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "IsCombined", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "VipPublicIpId", "VirtualIpAddress", "ServiceWatchLogCollection", "InitConfigOption"}

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

func (r *eventstreamsClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
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

func (r *eventstreamsClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
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

func (r *eventstreamsClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
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
				if len(bsPlan.Adds) > 0 || len(bsPlan.Removed) > 0 {
					return fmt.Errorf("operation not permitted for BLOCK_STORAGE_GROUP type: the eventstreams product does not support the addition of storage, so the addition of block storage in instances is restricted")
				}
				for _, resize := range bsPlan.Resizes {
					if err := r.client.SetBlockStorageSize(ctx, resize.Id, resize.SizeGb); err != nil {
						return err
					}
				}
			}

			// Instances Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"Instances"}) {
				// Kibana or DASHBOARDS
				t := currentInstanceGroup.RoleType.ValueString()
				if t == "KIBANA" || t == "DASHBOARDS" {
					return fmt.Errorf("instance group of type '%s' does not support adding instance", t)
				}

				var currentInst []databaseUtils.Instance
				currentInstanceGroup.Instances.ElementsAs(ctx, &currentInst, false)
				var desiredInst []databaseUtils.Instance
				desiredInstanceGroup.Instances.ElementsAs(ctx, &desiredInst, false)

				instancePlan := databaseUtils.PlanInstanceUpdate(currentInst, desiredInst)
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

			err := waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"EDITING"}, []string{"RUNNING"}, true)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *eventstreamsClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	diags.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("failed to read plan or state")
	}

	// Update
	_, err := tag.UpdateTags(r.clients, "eventstreams", "event-streams", plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventstreamsClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventstreams.ClusterResource
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
		resp.Diagnostics.AddError(
			"Error waiting for cluster deletion",
			"Could not wait for cluster deletion, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *eventstreamsClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForClusterStatus(ctx context.Context, esClient *eventstreams.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := esClient.GetCluster(ctx, id)
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
