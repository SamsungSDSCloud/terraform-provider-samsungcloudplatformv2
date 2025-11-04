package directconnect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpdirectconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/direct-connect/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &directConnectRoutingRuleResource{}
	_ resource.ResourceWithConfigure = &directConnectRoutingRuleResource{}
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
		Description: "routing rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "Direct Connect ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "Destination Type \n" +
					"  - example : ON-PREM | VPC",
				Required: true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "Destination CIDR \n" +
					"  - example : 10.10.10.0/24",
				Required: true,
			},
			common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
				Description: "Destination Resource ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Routing Rule description\n" +
					"  - maxLength : 50\n" +
					"  - minLength : 1",
				Optional: true,
			},
			common.ToSnakeCase("RoutingRule"): schema.SingleNestedAttribute{
				Description: "RoutingRule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("OwnerId"): schema.StringAttribute{
						Description: "OwnerId",
						Computed:    true,
					},
					common.ToSnakeCase("OwnerType"): schema.StringAttribute{
						Description: "OwnerType",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationType"): schema.StringAttribute{
						Description: "DestinationType",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "DestinationCidr",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
						Description: "DestinationResourceId",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
						Description: "DestinationResourceName",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
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
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
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

	routingRuleModel := createRoutingRuleModel(&routingRule)

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	plan.RoutingRule = routingRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

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
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading routing rule",
			"Could not read routing rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRuleModel := createRoutingRuleModel(&data.RoutingRules[0])

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	state.RoutingRule = routingRuleObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *directConnectRoutingRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
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
		DestinationType:         types.StringValue(string(data.OwnerType)),
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
	})
}
