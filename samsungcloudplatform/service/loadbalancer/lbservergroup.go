package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/loadbalancer/1.3"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLbServerGroupResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLbServerGroupResource{}
)

// NewLoadBalancerLbServerGroupResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLbServerGroupResource() resource.Resource {
	return &loadbalancerLbServerGroupResource{}
}

// loadbalancerLbServerGroupResource is the data source implementation.
type loadbalancerLbServerGroupResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbServerGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_server_group"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbServerGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "LB Server Group resource for managing server pools.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LbServerGroup"): schema.SingleNestedAttribute{
				Description: "Details of the LB Server Group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LB Server Group for web servers\n" +
							"  - maxLength : 255\n",
						Computed: true,
					},
					common.ToSnakeCase("LbMethod"): schema.StringAttribute{
						Description: "The load balancing method.\n" +
							"  - example : ROUND_ROBIN\n" +
							"  - pattern : ROUND_ROBIN | LEAST_CONNECTION | IP_HASH | WEIGHTED_ROUND_ROBIN | WEIGHTED_LEAST_CONNECTION\n",
						Computed: true,
					},
					common.ToSnakeCase("LbName"): schema.StringAttribute{
						Description: "The name of the LoadBalancer.\n" +
							"  - example : LoadBalancer01\n",
						Computed: true,
					},
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "The LoadBalancer ID associated with the server group.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the LB Server Group.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR | EDITING\n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Server Group.\n" +
							"  - example : ServerGroup01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n",
						Computed: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol for the server group.\n" +
							"  - example : TCP\n" +
							"  - pattern : TCP | UDP\n",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("LbHealthCheckId"): schema.StringAttribute{
						Description: "The LB Health Check ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("LbServerGroupCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new LB Server Group.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Server Group.\n" +
							"  - example : ServerGroup01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n",
						Optional: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol for the server group.\n" +
							"  - example : TCP\n" +
							"  - pattern : TCP | UDP\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LB Server Group for web servers\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("LbMethod"): schema.StringAttribute{
						Description: "The load balancing method.\n" +
							"  - example : ROUND_ROBIN\n" +
							"  - pattern : ROUND_ROBIN | LEAST_CONNECTION | IP_HASH | WEIGHTED_ROUND_ROBIN | WEIGHTED_LEAST_CONNECTION\n",
						Optional: true,
					},
					common.ToSnakeCase("LbHealthCheckId"): schema.StringAttribute{
						Description: "The LB Health Check ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbServerGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = inst.Client.LoadBalancer
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbServerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LbServerGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Server Group
	data, err := r.client.CreateLbServerGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Server Group",
			"Could not create Lb Server Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.LbServerGroup.Id)

	// Map response body to schema and populate Computed attribute values
	lbServerGroupModel := createLbServerGroupModel(data)
	lbServerGroupOjbectValue, diags := types.ObjectValueFrom(ctx, lbServerGroupModel.AttributeTypes(), lbServerGroupModel)
	plan.LbServerGroup = lbServerGroupOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbServerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbServerGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Server Group
	data, err := r.client.GetLbServerGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Server Group",
			"Could not create Lb Server Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbServerGroupModel := createLbServerGroupModel(data)

	lbServerGroupObjectValue, diags := types.ObjectValueFrom(ctx, lbServerGroupModel.AttributeTypes(), lbServerGroupModel)
	state.LbServerGroup = lbServerGroupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbServerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LbServerGroupResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	data, err := r.client.UpdateLbServerGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Server Group",
			"Could not create Lb Server Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbServerGroupModel := createLbServerGroupModel(data)

	lbServerGroupObjectValue, diags := types.ObjectValueFrom(ctx, lbServerGroupModel.AttributeTypes(), lbServerGroupModel)
	state.LbServerGroup = lbServerGroupObjectValue

	diags = resp.State.Set(ctx, state)

	err = waitForLbServerGroupStatus(ctx, r.client, data.LbServerGroup.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating lb server group",
			"Error waiting for lb server group to become active: "+err.Error(),
		)
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLbServerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LbServerGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Server Group
	err := r.client.DeleteLbServerGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Server Group",
			"Could not delete lb server group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForLbServerGroupStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error Deleting LB Server Group",
			"Error waiting for direct connect to become deleted: "+err.Error(),
		)
		return
	}

}

func createLbServerGroupModel(data *scploadbalancer.LbServerGroupShowResponse) loadbalancer.LbServerGroupDetail {
	lbServerGroup := data.LbServerGroup

	return loadbalancer.LbServerGroupDetail{
		Name:            types.StringValue(lbServerGroup.Name),
		Protocol:        types.StringValue(string(lbServerGroup.Protocol)),
		LoadbalancerId:  types.StringPointerValue(lbServerGroup.LoadbalancerId.Get()),
		LbName:          virtualserverutil.ToNullableStringValue(lbServerGroup.LbName.Get()),
		LbMethod:        types.StringValue(string(lbServerGroup.LbMethod)),
		LbHealthCheckId: virtualserverutil.ToNullableStringValue(lbServerGroup.LbHealthCheckId.Get()),
		State:           types.StringValue(lbServerGroup.State),
		VpcId:           types.StringValue(lbServerGroup.VpcId),
		SubnetId:        types.StringValue(lbServerGroup.SubnetId),
		AccountId:       types.StringValue(lbServerGroup.AccountId),
		Description:     virtualserverutil.ToNullableStringValue(lbServerGroup.Description.Get()),
		ModifiedBy:      types.StringValue(lbServerGroup.ModifiedBy),
		ModifiedAt:      types.StringValue(lbServerGroup.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(lbServerGroup.CreatedBy),
		CreatedAt:       types.StringValue(lbServerGroup.CreatedAt.Format(time.RFC3339)),
	}
}

func waitForLbServerGroupStatus(ctx context.Context, loadbalancerClient *loadbalancer.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := loadbalancerClient.GetLbServerGroup(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.LbServerGroup.State), nil
	}, -1, -1, -1, -1)
}
