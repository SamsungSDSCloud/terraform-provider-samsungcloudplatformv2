package directconnect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpdirectconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/direct-connect/1.0"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &directConnectRoutingRuleResource{}
	_ resource.ResourceWithConfigure   = &directConnectRoutingRuleResource{}
	_ resource.ResourceWithImportState = &directConnectRoutingRuleResource{}
)

// NewDirectConnectRoutingRuleResource is a helper function to simplify the provider implementation.
func NewDirectConnectRoutingRuleResource() resource.Resource {
	return &directConnectRoutingRuleResource{}
}

// networkRoutingRuleResource is the data source implementation.
type directConnectRoutingRuleResource struct {
	config  *scpsdk.Configuration
	client  *directconnect.Client
	clients *client.SCPClient
}

// Metadata returns the directConnectRoutingRuleResource source type name.
func (r *directConnectRoutingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directconnect_routing_rule"
}

// Schema defines the schema for the data source.
func (r *directConnectRoutingRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Direct Connect Routing Rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the routing rule.\n" +
					"  - example : fe860e0af0c04dcd8182b84f907f31f4",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "The identifier of the direct Connect.\n " +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "The type of the routing destination.In the VPC, the Direct Connect direction is ON_PREMISE, in the opposite direction—from Direct Connect toward the VPC—the direction is VPC.\n" +
					"  -  example : ON-PREMISE | VPC",
				Required: true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "The destination IP address range in CIDR notation.\n" +
					"  - example : 10.10.10.0/24",
				Required: true,
			},
			common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
				Description: "The identifier of the destination resource.When the Destination Type is VPC, provide the VpcId.\n " +
					"  -  example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this routing rule. This help identify the purpose or usage of the resource.\n" +
					"  - example : Routing Rule description\n" +
					"  - maxLength : 50\n" +
					"  - minLength : 1",
				Optional: true,
			},
			common.ToSnakeCase("RoutingRule"): schema.SingleNestedAttribute{
				Description: "Direct Connect Routing Rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the routing rule.\n" +
							"  - example : fe860e0af0c04dcd8182b84f907f31f4",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the direct connect.\n " +
							"  -  example: 27bb070b564349f8a31cc60734cc36a5",
						Computed: true,
					},
					common.ToSnakeCase("OwnerId"): schema.StringAttribute{
						Description: "The identifier of the routing rule owner.\n " +
							"  -  example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed: true,
					},
					common.ToSnakeCase("OwnerType"): schema.StringAttribute{
						Description: "The type of the routing rule owner.\n" +
							"  -  example: DIRECT_CONNECT",
						Computed: true,
					},
					common.ToSnakeCase("DestinationType"): schema.StringAttribute{
						Description: "The type of the routing destination.In the VPC, the Direct Connect direction is ON_PREMISE, in the opposite direction—from Direct Connect toward the VPC—the direction is VPC.\n" +
							"  -  example : ON-PREMISE | VPC",
						Computed: true,
					},
					common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "The destination IP address range in CIDR notation.\n" +
							"  - example : 10.10.10.0/24",
						Computed: true,
					},
					common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
						Description: "The identifier of the destination resource.When the Destination Type is VPC, provide the VpcId.\n " +
							"  -  example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
						Description: "The name of the destination resource.When the Destination Type is VPC, provide the Vpc Name.\n " +
							"  -  example : Resource Name",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this routing rule. This help identify the purpose or usage of the resource.\n" +
							"  - example : Routing Rule description\n" +
							"  - maxLength : 50\n" +
							"  - minLength : 1",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the routing rule.\n" +
							"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *directConnectRoutingRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.DirectConnect
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *directConnectRoutingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan directconnect.RoutingRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new routing rule
	data, err := r.client.CreateRoutingRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating routing rule",
			"Could not create routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRule := data.RoutingRule
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(routingRule.Id)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	routingRuleModel := createRoutingRuleModel(&routingRule)

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.RoutingRule = routingRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForRoutingRuleStatus(ctx, r.client, plan.DirectConnectId.ValueString(), data.RoutingRule.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing rule",
			"Error waiting for routing rule to become active: "+err.Error(),
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

func (r *directConnectRoutingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: DirectConnectId/RoutingRuleId, got: %q", req.ID),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("direct_connect_id"), types.StringValue(parts[0]))
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))
}

func (r *directConnectRoutingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state directconnect.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := r.client.GetRoutingRule(ctx, state.DirectConnectId.ValueString(), state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading routing rule",
			"Could not read routing rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// No data return from List API <=> Detail data not found
	if len(data.RoutingRules) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	routingRuleModel := createRoutingRuleModel(&data.RoutingRules[0])

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.RoutingRule = routingRuleObjectValue

	// Update top-level input fields from API response to detect drift
	rr := data.RoutingRules[0]
	state.DestinationType = types.StringValue(string(rr.DestinationType))
	state.DestinationCidr = types.StringValue(rr.DestinationCidr)
	state.DestinationResourceId = types.StringPointerValue(rr.DestinationResourceId.Get())
	if rr.Description != "" {
		state.Description = types.StringValue(rr.Description)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update is a no-op: the API has no update endpoint.
// RequiresReplace on all input fields ensures Terraform destroys and recreates instead.
func (r *directConnectRoutingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directconnect.RoutingRuleResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.AddWarning(
		"Routing rule update not supported",
		"Routing rule attributes are immutable. Terraform will replace this resource on the next apply.",
	)
	resp.State.Set(ctx, plan)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *directConnectRoutingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state directconnect.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing routing rule
	err := r.client.DeleteRoutingRule(ctx, state.DirectConnectId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting routing rule",
			"Could not delete routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForRoutingRuleStatus(ctx, r.client, state.DirectConnectId.ValueString(), state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting routing rule",
			"Error waiting for routing rule to become deleted: "+err.Error(),
		)
		return
	}
}

func createRoutingRuleModel(data *scpdirectconnect.RoutingRule) directconnect.RoutingRule {
	return directconnect.RoutingRule{
		Id:                      types.StringValue(data.Id),
		AccountId:               types.StringValue(data.AccountId),
		OwnerId:                 types.StringValue(data.OwnerId),
		OwnerType:               types.StringValue(string(data.OwnerType)),
		DestinationType:         types.StringValue(string(data.DestinationType)),
		DestinationCidr:         types.StringValue(data.DestinationCidr),
		DestinationResourceId:   types.StringPointerValue(data.DestinationResourceId.Get()),
		DestinationResourceName: types.StringPointerValue(data.DestinationResourceName.Get()),
		Description:             types.StringValue(data.Description),
		CreatedAt:               types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(data.CreatedBy),
		ModifiedAt:              types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(data.ModifiedBy),
		State:                   types.StringValue(string(data.State)),
	}
}

func waitForRoutingRuleStatus(ctx context.Context, routingRuleClient *directconnect.Client, directConnectId string, routingRuleId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := routingRuleClient.GetRoutingRule(ctx, directConnectId, routingRuleId)
		if err != nil {
			return nil, "", err
		}
		if len(info.RoutingRules) == 0 {
			return info, "DELETED", nil
		}
		return info, string(info.RoutingRules[0].State), nil
	}, -1, -1, -1, -1)
}
