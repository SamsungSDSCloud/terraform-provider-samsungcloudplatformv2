package multinodegpucluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/multinodegpucluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	multinodegpuclustersdk1d2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/multinodegpucluster/1.2"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &GpunodeResource{}
	_ resource.ResourceWithConfigure   = &GpunodeResource{}
	_ resource.ResourceWithModifyPlan  = &GpunodeResource{}
	_ resource.ResourceWithImportState = &GpunodeResource{}
)

// NewGpunodeResource is a helper function to simplify the provider implementation.
func NewGpunodeResource() resource.Resource {
	return &GpunodeResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type GpunodeResource struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

var ERROR_EXPLAIN = "\nReason: "
var errNotFound = fmt.Errorf("resource not found")

// Metadata returns the resource type name.
func (multinodegpuclusterRS *GpunodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunode"
}

// Configure adds the provider configured client to the resource.
func (multinodegpuclusterRS *GpunodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expect *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	multinodegpuclusterRS.client = inst.Client.Mngc
	multinodegpuclusterRS.clients = inst.Client
}

// Schema defines the schema for the resource.
func (multinodegpuclusterRS *GpunodeResource) Schema(c context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GpuNodeResourceSchema(c)
}

func GpuNodeResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "GPU Node Resource",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_fabric_details": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"cluster_fabric_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Cluster Fabric ID\n - example: 20c507a036c447cdb3b19468d8ea62ac",
						MarkdownDescription: "Cluster Fabric ID\n - example: 20c507a036c447cdb3b19468d8ea62ac",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"cluster_fabric_name": schema.StringAttribute{
						Required: true,
						Description: "Cluster Fabric Name. Enter 3 to 15 char (English, number and -) starting with English." +
							" Do not use '-' for the last char." +
							"  - minLength: 3\n" +
							"  - maxLength: 15\n" +
							"  - pattern: ^([a-zA-Z]{1})([a-zA-Z0-9-]{1,13})([a-zA-Z0-9]{1})$\n" +
							"  - example: cluster001",
						MarkdownDescription: "Cluster Fabric Name. Enter 3 to 15 char (English, number and -) starting with English." +
							" Do not use '-' for the last char." +
							"  - minLength: 3\n" +
							"  - maxLength: 15\n" +
							"  - pattern: ^([a-zA-Z]{1})([a-zA-Z0-9-]{1,13})([a-zA-Z0-9]{1})$\n" +
							"  - example: cluster001",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^([a-zA-Z]{1})([a-zA-Z0-9-]{1,13})([a-zA-Z0-9]{1})$"), "Cluster Fabric Name. Enter 3 to 15 char (English, number and -) starting with English. Do not use '-' for the last char."),
						},
					},
					"node_pool_id": schema.StringAttribute{
						Required:            true,
						Description:         "Node Pool ID\n  - example: POOL001-krw1a",
						MarkdownDescription: "Node Pool ID\n  - example: POOL001-krw1a",
					},
				},
				CustomType: multinodegpuclusterClient.ClusterFabricDetailsType{
					ObjectType: types.ObjectType{
						AttrTypes: multinodegpuclusterClient.ClusterFabricDetailsValue{}.AttributeTypes(ctx),
					},
				},
				Required: true,
				Description: "Cluster Fabric Details  \n" +
					"  - example: {cluster_fabric_id='20c507a036c447cdb3b19468d8ea62ac', cluster_fabric_name='cluster001', node_pool_id='POOL001-krw1a'}",
				MarkdownDescription: "Cluster Fabric Details  \n" +
					"  - example: {cluster_fabric_id='YOUR RESOURCE'S CLUSTER_FABRIC_ID', cluster_fabric_name='cluster001', node_pool_id='YOUR RESOURCE'S NODE_POOL_ID'}",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gpu_node_name_prefix": schema.StringAttribute{
				Required:  true,
				WriteOnly: true,
				Description: "GPU Node name start with a lowercase, and enter 3 to 24 using lowercase, number and -." +
					" (It does not end with -.)\n" +
					"  - example: gpunode-1\n" +
					"  - minLength: 3\n" +
					"  - maxLength(for linux): 24\n" +
					"  - pattern(for linux): ^[a-z][a-z0-9-]{1,22}[a-z0-9]$",
				MarkdownDescription: "GPU Node name start with a lowercase, and enter 3 to 24 using lowercase, number and -." +
					" (It does not end with -.)\n" +
					"  - example: gpunode-1\n" +
					"  - minLength: 3\n" +
					"  - maxLength(for linux): 24\n" +
					"  - pattern(for linux): ^[a-z][a-z0-9-]{1,22}[a-z0-9]$",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9-]{1,22}[a-z0-9]$"), "Server Name Prefix must begin with a lowercase letter and must be entered in 3-24 characters using lowercase letters, numbers, and -. (It does not end with -.)"),
				},
			},
			"image_id": schema.StringAttribute{
				Required:            true,
				Description:         "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
				MarkdownDescription: "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
			},
			"image_version": schema.StringAttribute{
				Computed:            true,
				Description:         "Image Version\n  - example: RHEL 8.7 for BM",
				MarkdownDescription: "Image Version\n  - example: RHEL 8.7 for BM",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"init_script": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Init Script\n  - maxLength: 16384\n  - example: #!/bin/bash\\necho 'Hello World!'",
				MarkdownDescription: "Init Script\n  - maxLength: 16384\n  - example: #!/bin/bash\\necho 'Hello World!'",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(16384),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lock_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Use Lock\n  - example: true",
				MarkdownDescription: "Use Lock\n  - example: true",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
			},
			"os_type": schema.StringAttribute{
				Computed:            true,
				Description:         "OS type\n  - example: WINDOWS",
				MarkdownDescription: "OS type\n  - example: WINDOWS",
			},
			"os_user_id": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "OS User Id. When linux image value must be 'root'\n  - example: root",
				MarkdownDescription: "OS User Id. When linux image value must be 'root'\n  - example: root",
				Default:             stringdefault.StaticString("root"),
			},
			"os_user_password": schema.StringAttribute{
				Required:  true,
				WriteOnly: true,
				Description: "OS user password.\n" +
					"  - example: P@ssword1!2\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 9\n" +
					"  - pattern: ^[A-Za-z0-9@$!%*#&]+$",
				MarkdownDescription: "OS user password.\n" +
					"  - example: P@ssword1!2\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 9\n" +
					"  - pattern: ^[A-Za-z0-9@$!%*#&]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(9, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[A-Za-z0-9@$!%*#&]+$"), "English Capital letters, lowercase letters, numbers, @$!%*#& special characters only."),
					stringvalidator.RegexMatches(regexp.MustCompile("[A-Z]"), "Must contain capital letters."),
					stringvalidator.RegexMatches(regexp.MustCompile("[a-z]"), "Must contain lowercase characters."),
					stringvalidator.RegexMatches(regexp.MustCompile("[0-9]"), "Must contain number."),
					stringvalidator.RegexMatches(regexp.MustCompile("[@$!%*#&]"), "'@$!%*#&' must contain special characters"),
				},
			},
			"product_type_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
				MarkdownDescription: "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
			},
			"region_id": schema.StringAttribute{
				Required:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"root_account": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Account\n  - example: root",
				MarkdownDescription: "Root Account\n  - example: root",
			},
			"server_details": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gpu_node_name": schema.StringAttribute{
							Computed: true,
							Description: "GPU Node name in format prefix-###.\n" +
								"  - example: gpunode-001\n",
							MarkdownDescription: "GPU Node name in format prefix-###.\n" +
								"  - example: gpunode-001\n",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "GPU Node ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
							MarkdownDescription: "GPU Node ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"ip_address": schema.StringAttribute{
							Computed:            true,
							Description:         "subnet IP address\n  - example: 192.168.2.4",
							MarkdownDescription: "subnet IP address\n  - example: 192.168.2.4",
						},
						"nat_enabled": schema.BoolAttribute{
							Computed:            true,
							Description:         "Use Public NAT\n  - example: true",
							MarkdownDescription: "Use Public NAT\n  - example: true",
						},
						"pfs_ip": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "PFS IP list\n  - example: [10.252.128.2 10.252.128.3]",
							MarkdownDescription: "PFS IP list\n  - example: [10.252.128.2 10.252.128.3]",
						},
						"policy_ip": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy IP\n  - example: 192.168.0.1",
							MarkdownDescription: "Policy IP\n  - example: 192.168.0.1",
						},
						"policy_nat": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy NAT\n  - example: 192.168.0.1",
							MarkdownDescription: "Policy NAT\n  - example: 192.168.0.1",
						},
						"policy_use_nat": schema.BoolAttribute{
							Computed:            true,
							Description:         "Policy use NAT\n  - example: true",
							MarkdownDescription: "Policy use NAT\n  - example: true",
						},
						"server_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Server type\n  - example: s1v8m32_metal",
							MarkdownDescription: "Server type\n  - example: s1v8m32_metal",
						},
						"state": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "Server state\n" +
								"  - example: RUNNING\n" +
								"  - pattern: RUNNING | STOPPED",
							MarkdownDescription: "Server state\n" +
								"  - example: RUNNING\n" +
								"  - pattern: RUNNING | STOPPED",
							Default: stringdefault.StaticString(common.RunningState),
							Validators: []validator.String{
								stringvalidator.OneOf(common.RunningState, common.StoppedState),
							},
						},
					},
				},
				Required: true,
				Description: "Detailed settings for each server, 2 or more server on creation\n" +
					"  - example: [{state: RUNNING}, {state: RUNNING}]\n" +
					"  - maxLength: 5\n" +
					"  - minLength: 2\n",
				MarkdownDescription: "Detailed settings for each server, 2 or more server on creation\n" +
					"  - example: [{state: RUNNING}, {state: RUNNING}]\n" +
					"  - maxLength: 5\n" +
					"  - minLength: 2\n",
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 5),
				},
			},
			"server_type_id": schema.StringAttribute{
				Required:            true,
				WriteOnly:           true,
				Description:         "Server Type ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
				MarkdownDescription: "Server Type ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
			},
			"subnet_id": schema.StringAttribute{
				Required:            true,
				Description:         "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
			},
			"tags": tag.ResourceSchema(),
			"time_zone": schema.StringAttribute{
				Computed:            true,
				Description:         "Time Zone\n  - example: Asia/Seoul",
				MarkdownDescription: "Time Zone\n  - example: Asia/Seoul",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (multinodegpuclusterRS *GpunodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// OPEN: read plan and user config
	var plan, draft multinodegpuclusterClient.GpuNodeResource

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)

	diags = req.Config.Get(ctx, &draft)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// MAIN LOGIC
	asyncResponse, err := multinodegpuclusterRS.client.CreateGpuNode(ctx, plan, draft)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating GPU Nodes",
			"Could not create GPU Nodes, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	space := regexp.MustCompile(`\s+`)
	ids := space.ReplaceAllString(asyncResponse.ResourceId, "")
	idList := strings.Split(ids, ",")
	if len(idList) == 0 {
		resp.Diagnostics.AddError(
			"Error creating GPU Node",
			"Error in creating : GPU Node was not created.",
		)
		return
	}

	for _, id := range idList {
		err = waitForGpuNodeStatus(ctx, multinodegpuclusterRS.client, id, []string{common.CreatingState}, []string{common.RunningState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error creating GPU Node",
				"Error waiting for GPU Node("+id+")to become RUNNING:"+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}
	gpunodeId := idList[0]

	err = setGpuNodeCommonInfo(ctx, multinodegpuclusterRS, gpunodeId, &plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+gpunodeId+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	// Get Server Details
	serverDetails, diags, err := findNodesInFabric(ctx, multinodegpuclusterRS, idList)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching GPU Nodes in Fabric "+plan.ClusterFabricDetails.ClusterFabricId.ValueString(), err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ServerDetails = serverDetails

	// CLOSE: Set state to plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (multinodegpuclusterRS *GpunodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// OPEN: read state
	var state multinodegpuclusterClient.GpuNodeResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// MAIN LOGIC
	var serversState []multinodegpuclusterClient.ServerDetailsValue
	state.ServerDetails.ElementsAs(ctx, &serversState, false)

	if len(serversState) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	gpunodeId := serversState[0].Id.ValueString()

	err := setGpuNodeCommonInfo(ctx, multinodegpuclusterRS, gpunodeId, &state)
	if err != nil {
		if err == errNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+gpunodeId+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	// Get Tags
	tagsMap, err := tag.GetTags(multinodegpuclusterRS.clients, "multinodegpucluster", "gpu-node", gpunodeId)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading GPU Node Tag", err.Error())
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)
	state.Tags = tagsMap

	// Get Server Details
	var gpunodeIds []string
	for _, server := range serversState {
		gpunodeIds = append(gpunodeIds, server.Id.ValueString())
	}

	serverDetails, diags, err := findNodesInFabric(ctx, multinodegpuclusterRS, gpunodeIds)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		if err == errNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Searching GPU Nodes in Fabric "+state.ClusterFabricDetails.ClusterFabricId.ValueString(), err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state.ServerDetails = serverDetails

	// CLOSE: Update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (multinodegpuclusterRS *GpunodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// OPEN: read state and plan
	var state, plan multinodegpuclusterClient.GpuNodeResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planServerDetails, stateServerDetails []multinodegpuclusterClient.ServerDetailsValue
	plan.ServerDetails.ElementsAs(ctx, &planServerDetails, false)
	state.ServerDetails.ElementsAs(ctx, &stateServerDetails, false)

	if len(planServerDetails) != len(stateServerDetails) {
		resp.Diagnostics.AddError("Error updating ", "Cannot add or remove fabric member.")
		return
	}

	planServerDetailsMap := make(map[string]multinodegpuclusterClient.ServerDetailsValue)
	for _, sd := range planServerDetails {
		planServerDetailsMap[sd.Id.ValueString()] = sd
	}

	var stopIds, startIds []string

	// MAIN LOGIC

	for _, stateServer := range stateServerDetails {
		planServer, ok := planServerDetailsMap[stateServer.Id.ValueString()]
		if !ok {
			continue
		}
		if !planServer.State.Equal(stateServer.State) {
			// RUNNING -> STOPPED
			if planServer.State.ValueString() == common.StoppedState {
				stopIds = append(stopIds, planServer.Id.ValueString())
			} else {
				// STOPPED -> RUNNING
				startIds = append(startIds, planServer.Id.ValueString())
			}
		}
	}

	if len(stopIds) != 0 {
		err := multinodegpuclusterRS.client.StopGpuNodes(ctx, stopIds)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Stopping GPU Node",
				"Could not stop GPU Node, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}

	if len(startIds) != 0 {
		err := multinodegpuclusterRS.client.StartGpuNodes(ctx, startIds)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Starting GPU Node",
				"Could not start GPU Node, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}

	for _, Id := range stopIds {
		err := waitForGpuNodeStatus(ctx, multinodegpuclusterRS.client, Id, []string{common.StoppingState}, []string{common.StoppedState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Stopping GPU Node",
				"Could not stop GPU Node, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}

	for _, Id := range startIds {
		err := waitForGpuNodeStatus(ctx, multinodegpuclusterRS.client, Id, []string{common.StartingState}, []string{common.RunningState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Starting GPU Node",
				"Could not start GPU Node, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}

	// Refesh data
	if len(planServerDetails) == 0 {
		resp.Diagnostics.AddError(
			"Error Refresh GPU Node",
			"Could not refresh Plan Server Details",
		)
		return
	}
	gpunodeId := planServerDetails[0].Id.ValueString()

	err := setGpuNodeCommonInfo(ctx, multinodegpuclusterRS, gpunodeId, &plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+gpunodeId+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	// Get Server Details
	var gpunodeIds []string
	for _, server := range planServerDetails {
		gpunodeIds = append(gpunodeIds, server.Id.ValueString())
	}
	serverDetails, diags, err := findNodesInFabric(ctx, multinodegpuclusterRS, gpunodeIds)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching GPU Nodes in Fabric "+plan.ClusterFabricDetails.ClusterFabricId.ValueString(), err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ServerDetails = serverDetails

	// update tags
	tagElements := plan.Tags.Elements()
	for _, serverDetail := range planServerDetails {
		tagsMap, err := tag.UpdateTags(multinodegpuclusterRS.clients, "multinodegpucluster", "gpu-node", serverDetail.Id.ValueString(), tagElements)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating tags",
				err.Error(),
			)
			return
		}
		tagsMap = common.NullTagCheck(tagsMap, plan.Tags)
		plan.Tags = tagsMap
	}

	// CLOSE: Set state to plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (multinodegpuclusterRS *GpunodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// OPEN: read state
	var state multinodegpuclusterClient.GpuNodeResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 40*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	// MAIN LOGIC
	var serverDetails []multinodegpuclusterClient.ServerDetailsValue
	state.ServerDetails.ElementsAs(ctx, &serverDetails, false)

	ids := make([]string, 0)
	for _, serverDetail := range serverDetails {
		ids = append(ids, serverDetail.Id.ValueString())
	}

	// Delete GPU Nodes
	err := multinodegpuclusterRS.client.DeleteGpuNodes(ctx, ids)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting GPU Node",
			"Could not delete GPU Node, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	for _, nodeId := range ids {
		err = waitForGpuNodeStatus(ctx, multinodegpuclusterRS.client, nodeId, []string{common.TerminatingState}, []string{common.TerminatedState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Deleting GPU Node",
				"Could not delete GPU Node("+nodeId+")unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
			return
		}
	}

}

func (multinodegpuclusterRS *GpunodeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// OPEN: read state (if verify update) and plan
	var plan multinodegpuclusterClient.GpuNodeResource

	// destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// MAIN LOGIC
	var planServerDetails []multinodegpuclusterClient.ServerDetailsValue
	plan.ServerDetails.ElementsAs(ctx, &planServerDetails, false)

	if len(planServerDetails) == 0 {
		resp.Diagnostics.AddError("Error", "server_details cannot be empty")
		return
	}

	if planServerDetails[0].Id.IsUnknown() {
		// create validate start

		if len(planServerDetails) < 2 {
			resp.Diagnostics.AddError(
				"Error creating server",
				"The minimum number of generated servers in the GPU Node is 2",
			)
		}

		imageList, err := multinodegpuclusterRS.client.GetImageList(ctx, plan.RegionId.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error reading image",
				"Could not read image data, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
			)
		}

		_, e := getImageInfo(imageList, plan.ImageId.ValueString())
		if e != nil {
			resp.Diagnostics.AddError("ImageId error", e.Error())
		}

		// validate server info
		errList := validateServerDetailInfo(planServerDetails)

		for _, e := range errList {
			resp.Diagnostics.AddError("server detail info error", e.Error())
		}

		// create validate end
		return
	}

	// update validation - required fields cannot be modified
	var state multinodegpuclusterClient.GpuNodeResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.ClusterFabricDetails.ClusterFabricId.Equal(plan.ClusterFabricDetails.ClusterFabricId) ||
		!state.ClusterFabricDetails.ClusterFabricName.Equal(plan.ClusterFabricDetails.ClusterFabricName) ||
		!state.ClusterFabricDetails.NodePoolId.Equal(plan.ClusterFabricDetails.NodePoolId) {
		resp.Diagnostics.AddError("Could not change cluster_fabric_details",
			"Could not change cluster_fabric_details.\nIf you want to change, create a new resource.")
	}

	if !state.ImageId.Equal(plan.ImageId) {
		resp.Diagnostics.AddError("Could not change image_id",
			"Could not change image_id.\nIf you want to change, create a new resource.")
	}

	if !state.RegionId.Equal(plan.RegionId) {
		resp.Diagnostics.AddError("Could not change region_id",
			"Could not change region_id.\nIf you want to change, create a new resource.")
	}

	var stateServerDetails []multinodegpuclusterClient.ServerDetailsValue
	state.ServerDetails.ElementsAs(ctx, &stateServerDetails, false)

	if len(planServerDetails) != len(stateServerDetails) {
		resp.Diagnostics.AddError("Could not change server_details count",
			"Could not add or remove servers.\nIf you want to change, create a new resource.")
	}

	if !state.SubnetId.Equal(plan.SubnetId) {
		resp.Diagnostics.AddError("Could not change subnet_id",
			"Could not change subnet_id.\nIf you want to change, create a new resource.")
	}

	if !state.VpcId.Equal(plan.VpcId) {
		resp.Diagnostics.AddError("Could not change vpc_id",
			"Could not change vpc_id.\nIf you want to change, create a new resource.")
	}

	if !state.LockEnabled.Equal(plan.LockEnabled) {
		resp.Diagnostics.AddError("Could not change lock_enabled",
			"Could not change lock_enabled.\nIf you want to change, create a new resource")
	}

	if !state.InitScript.Equal(plan.InitScript) {
		resp.Diagnostics.AddError("Could not change init_script",
			"Could not change init_script.\nIf you want to change, create a new resource.")
	}

	if !state.OsUserId.Equal(plan.OsUserId) {
		resp.Diagnostics.AddError("Could not change os_user_id",
			"Could not change os_user_id.\nIf you want to change, create a new resource.")
	}

}

func (multinodegpuclusterRS *GpunodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The raw ID supplied by the user on the CLI
	gpunodeId := req.ID
	if gpunodeId == "" {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"The import ID cannot be empty",
		)
		return
	}

	var state multinodegpuclusterClient.GpuNodeResource

	// populate default as new creation, be refreshed later in read phase
	var diags diag.Diagnostics
	state.ServerDetails, diags = types.ListValueFrom(
		ctx,
		types.ObjectType{
			AttrTypes: multinodegpuclusterClient.ServerDetailsValue{}.AttributeTypes(),
		},
		[]multinodegpuclusterClient.ServerDetailsValue{{
			Id:    types.StringValue(gpunodeId),
			PfsIp: types.ListNull(types.StringType),
		}},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Tags, diags = types.MapValue(types.StringType, make(map[string]attr.Value, 0))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"delete": types.StringType,
		}),
	}

	// CLOSE: Update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func setGpuNodeDetailToResource(nodeDetail *multinodegpuclustersdk1d2.GpuNodeShowResponse, resource *multinodegpuclusterClient.GpuNodeResource) {

	// basic info
	resource.AccountId = types.StringValue(nodeDetail.AccountId)
	resource.ImageId = types.StringValue(nodeDetail.ImageId)
	resource.ImageVersion = types.StringValue(nodeDetail.ImageVersion)
	resource.OsType = types.StringValue(nodeDetail.OsType)
	resource.RegionId = types.StringValue(nodeDetail.RegionId)
	resource.RootAccount = types.StringValue(nodeDetail.RootAccount)
	resource.OsUserId = types.StringValue(nodeDetail.RootAccount)
	resource.TimeZone = types.StringValue(nodeDetail.TimeZone)
	resource.ProductTypeId = types.StringValue(nodeDetail.ProductTypeId)

	// cluster info
	clusterFabricDetails := multinodegpuclusterClient.ClusterFabricDetailsValue{}
	clusterFabricDetails.ClusterFabricName = types.StringValue(nodeDetail.ClusterFabricName)
	if nodeDetail.ClusterFabricId.IsSet() {
		clusterFabricDetails.ClusterFabricId = types.StringPointerValue(nodeDetail.ClusterFabricId.Get())
	}
	clusterFabricDetails.NodePoolId = types.StringValue(nodeDetail.NodePoolId)
	resource.ClusterFabricDetails = clusterFabricDetails

	// network info
	resource.SubnetId = types.StringValue(nodeDetail.NetworkId)
	resource.VpcId = types.StringValue(nodeDetail.VpcId)

	// additional info
	resource.InitScript = types.StringValue(nodeDetail.InitScript)
	resource.LockEnabled = types.BoolPointerValue(nodeDetail.LockEnabled)

	// metadata info
	resource.CreatedBy = types.StringValue(nodeDetail.CreatedBy)
	resource.CreatedAt = types.StringValue(nodeDetail.CreatedAt.String())
	resource.ModifiedBy = types.StringValue(nodeDetail.ModifiedBy)
	resource.ModifiedAt = types.StringValue(nodeDetail.ModifiedAt.String())

}

func waitForGpuNodeStatus(ctx context.Context, mngcClient *multinodegpuclusterClient.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpResponse, err := mngcClient.GetGpuNode(ctx, id)
		if err != nil {
			if httpResponse != nil {
				if httpResponse.StatusCode == 404 {
					return "", common.TerminatedState, nil
				}
			}
			return nil, "", err
		}
		return info, info.State, nil
	}, -1, -1, -1, -1)
}

func setGpuNodeCommonInfo(ctx context.Context, multinodegpuclusterRS *GpunodeResource, gpunodeId string, resource *multinodegpuclusterClient.GpuNodeResource) error {
	gpunodeShow, httpResponse, err := multinodegpuclusterRS.client.GetGpuNode(ctx, gpunodeId)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			return errNotFound
		}
		return err
	}

	setGpuNodeDetailToResource(gpunodeShow, resource)

	return nil

}

func findNodesInFabric(ctx context.Context, multinodegpuclusterRS *GpunodeResource, gpunodeIds []string) (types.List, diag.Diagnostics, error) {
	serverDetails := make([]multinodegpuclusterClient.ServerDetailsValue, len(gpunodeIds))
	var allDiags diag.Diagnostics
	for pos, gpunodeId := range gpunodeIds {
		gpunodeShow, _, err := multinodegpuclusterRS.client.GetGpuNode(ctx, gpunodeId)
		if err != nil {
			return types.ListNull(types.ObjectType{
				AttrTypes: multinodegpuclusterClient.ServerDetailsValue{}.AttributeTypes(),
			}), nil, err
		}
		serverDetail, d := setServerDetailInfo(ctx, *gpunodeShow)
		allDiags.Append(d...)
		serverDetails[pos] = serverDetail
	}

	if len(serverDetails) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: multinodegpuclusterClient.ServerDetailsValue{}.AttributeTypes(),
		}), allDiags, nil
	}

	serverDetailsListType, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: serverDetails[0].AttributeTypes(),
	}, serverDetails)
	allDiags.Append(d...)
	return serverDetailsListType, allDiags, nil
}

func setServerDetailInfo(ctx context.Context, nodeDetail multinodegpuclustersdk1d2.GpuNodeShowResponse) (multinodegpuclusterClient.ServerDetailsValue, diag.Diagnostics) {
	var serverDetail multinodegpuclusterClient.ServerDetailsValue
	var diags diag.Diagnostics

	serverDetail.GpuNodeName = types.StringValue(nodeDetail.GpuNodeName)
	serverDetail.Id = types.StringValue(nodeDetail.Id)
	serverDetail.IpAddress = types.StringValue(nodeDetail.PolicyIp)
	serverDetail.NatEnabled = types.BoolPointerValue(nodeDetail.PolicyUseNat)
	serverDetail.ServerType = types.StringValue(nodeDetail.ServerType)
	serverDetail.State = types.StringValue(nodeDetail.State)

	serverDetail.PfsIp, diags = types.ListValueFrom(ctx, types.StringType, nodeDetail.PfsIp)
	serverDetail.PolicyIp = types.StringValue(nodeDetail.PolicyIp)
	if nodeDetail.PolicyNat.IsSet() {
		serverDetail.PolicyNat = types.StringPointerValue(nodeDetail.PolicyNat.Get())
	}
	serverDetail.PolicyUseNat = types.BoolPointerValue(nodeDetail.PolicyUseNat)

	return serverDetail, diags
}

func getImageInfo(imageList *multinodegpuclustersdk1d2.GpuNodeImageListResponse, imageId string) (string, error) {
	var osType string

	for _, image := range imageList.Images {
		if image.Id == imageId {
			osType = image.OsDistro
			break
		}
	}

	if osType == "" {
		return "", errors.New("no image exists for imageId")
	}

	return osType, nil
}

func validateServerDetailInfo(serverDetails []multinodegpuclusterClient.ServerDetailsValue) []error {

	var errorList []error
	for _, serverDetail := range serverDetails {

		// validate state
		if serverDetail.State.ValueString() != common.RunningState {
			errorList = append(errorList, errors.New("state value must be RUNNING, when server creation request is made"))
		}
	}

	return errorList

}
