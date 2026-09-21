package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancerv1d4"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.3"
	scploadbalancerv1d4 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loadbalancerLbServerGroupResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLbServerGroupResource{}
	_ resource.ResourceWithImportState = &loadbalancerLbServerGroupResource{}
)

// NewLoadBalancerLbServerGroupResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLbServerGroupResource() resource.Resource {
	return &loadbalancerLbServerGroupResource{}
}

// loadbalancerLbServerGroupResource is the data source implementation.
type loadbalancerLbServerGroupResource struct {
	config     *scpsdk.Configuration
	client     *loadbalancer.Client
	clientv1d4 *loadbalancerv1d4.Client
	clients    *client.SCPClient
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
							"  - pattern : ROUND_ROBIN | LEAST_CONNECTION | SOURCE_IP_PORT_HASH | SOURCE_IP_HASH | WEIGHTED_ROUND_ROBIN | WEIGHTED_LEAST_CONNECTION\n",
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
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR | EDITING | TERMINATING\n",
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
							"  - pattern : ROUND_ROBIN | LEAST_CONNECTION | SOURCE_IP_PORT_HASH | SOURCE_IP_HASH | WEIGHTED_ROUND_ROBIN | WEIGHTED_LEAST_CONNECTION\n",
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
	r.clientv1d4 = inst.Client.LoadBalancerV1d4
	r.clients = inst.Client
}

func (r *loadbalancerLbServerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

	lbServerGroup := plan.LbServerGroupCreate

	lbServerGroupElement := scploadbalancerv1d4.LbServerGroupCreate{
		Name:            lbServerGroup.Name.ValueString(),
		VpcId:           lbServerGroup.VpcId.ValueString(),
		SubnetId:        lbServerGroup.SubnetId.ValueString(),
		Protocol:        scploadbalancerv1d4.LbServerGroupProtocol(lbServerGroup.Protocol.ValueString()),
		LbMethod:        scploadbalancerv1d4.LbServerGroupLbMethod(lbServerGroup.LbMethod.ValueString()),
		Description:     *scploadbalancerv1d4.NewNullableString(lbServerGroup.Description.ValueStringPointer()),
		LbHealthCheckId: *scploadbalancerv1d4.NewNullableString(lbServerGroup.LbHealthCheckId.ValueStringPointer()),
		Tags:            convertToTags(lbServerGroup.Tags.Elements()),
	}

	data, err := r.clientv1d4.CreateLbServerGroupV1d4(ctx, scploadbalancerv1d4.LbServerGroupCreateRequest{
		LbServerGroup: lbServerGroupElement,
	})
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
	lbServerGroupModel := loadbalancerv1d4.LbServerGroupResourceDetail{
		Name:            types.StringValue(data.LbServerGroup.Name),
		Protocol:        types.StringValue(string(data.LbServerGroup.Protocol)),
		LoadbalancerId:  types.StringPointerValue(data.LbServerGroup.LoadbalancerId.Get()),
		LbName:          virtualserverutil.ToNullableStringValue(data.LbServerGroup.LbName.Get()),
		LbMethod:        types.StringValue(string(data.LbServerGroup.LbMethod)),
		LbHealthCheckId: virtualserverutil.ToNullableStringValue(data.LbServerGroup.LbHealthCheckId.Get()),
		State:           types.StringValue(data.LbServerGroup.State),
		VpcId:           types.StringValue(data.LbServerGroup.VpcId),
		SubnetId:        types.StringValue(data.LbServerGroup.SubnetId),
		AccountId:       types.StringValue(data.LbServerGroup.AccountId),
		Description:     virtualserverutil.ToNullableStringValue(data.LbServerGroup.Description.Get()),
		ModifiedBy:      types.StringValue(data.LbServerGroup.ModifiedBy),
		ModifiedAt:      types.StringValue(data.LbServerGroup.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(data.LbServerGroup.CreatedBy),
		CreatedAt:       types.StringValue(data.LbServerGroup.CreatedAt.Format(time.RFC3339)),
	}

	lbServerGroupOjbectValue, diags := types.ObjectValueFrom(ctx, lbServerGroupModel.AttributeTypes(), lbServerGroupModel)
	plan.LbServerGroup = lbServerGroupOjbectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshFn := r.getLbServerGroupRefreshFunc(ctx, data.LbServerGroup.Id)
	err = client.WaitForResourceCreated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating lbServerGroup ",
			"Error waiting for lbServerGroup to become active: "+err.Error(),
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
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Reconcile lb_server_group_create input block with API response to detect drift
	// Only populate if nil (e.g., after import) — preserve user config values otherwise
	if state.LbServerGroupCreate == nil {
		state.LbServerGroupCreate = &loadbalancer.LbServerGroupCreate{
			Name:            types.StringValue(data.LbServerGroup.Name),
			Protocol:        types.StringValue(string(data.LbServerGroup.Protocol)),
			VpcId:           types.StringValue(data.LbServerGroup.VpcId),
			SubnetId:        types.StringValue(data.LbServerGroup.SubnetId),
			Description:     virtualserverutil.ToNullableStringValue(data.LbServerGroup.Description.Get()),
			LbMethod:        types.StringValue(string(data.LbServerGroup.LbMethod)),
			LbHealthCheckId: virtualserverutil.ToNullableStringValue(data.LbServerGroup.LbHealthCheckId.Get()),
			Tags:            types.MapNull(types.StringType),
		}
	}

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

	refreshFn := r.getLbServerGroupRefreshFunc(ctx, data.LbServerGroup.Id)
	err = client.WaitForResourceUpdated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating LbServerGroup",
			"Error waiting for LbServerGroup to become ACTIVE: "+err.Error(),
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

	refreshFn := r.getLbServerGroupRefreshFunc(ctx, state.Id.ValueString())
	err = client.WaitForResourceDeleted(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting LbServerGroup",
			"Error waiting for LbServerGroup to become deleted: "+err.Error(),
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

func convertToTags(elements map[string]attr.Value) []scploadbalancerv1d4.Tag {
	var tags []scploadbalancerv1d4.Tag
	for k, v := range elements {
		tagObject := scploadbalancerv1d4.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func (r *loadbalancerLbServerGroupResource) getLbServerGroupRefreshFunc(ctx context.Context, id string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		data, err := r.client.GetLbServerGroup(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return data, data.LbServerGroup.State, nil
	}
}
