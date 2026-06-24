package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcTgwRuleResource{}
	_ resource.ResourceWithConfigure = &vpcTgwRuleResource{}
)

// NewVpcTgwRuleResource is a helper function to simplify the provider implementation.
func NewVpcTgwRuleResource() resource.Resource {
	return &vpcTgwRuleResource{}
}

type vpcTgwRuleResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (v vpcTgwRuleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_transit_gateway_rule"
}

// Schema defines the schema for the data source.
func (v *vpcTgwRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit gateway rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the transit gateway rule.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the rule belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
					"  - example : Routing Rule description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "The destination IP address range in CIDR notation.\n" +
					"  - example : 10.10.10.0/24",
				Required: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "The type of the destination.\n" +
					"  - example : VPC | TGW",
				Required: true,
			},
			common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the transit gateway connection belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("RoutingRule"): schema.SingleNestedAttribute{
				Description: "Routing rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the routing rule.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                        Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
                            "  - example : 2024-05-17T00:23:17Z",
                        Computed: true,
                    },
                    common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                        Description: "The user id that created the resource.\n" +
                            "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                        Computed: true,
                    },
                    common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                        Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
                            "  - example : 2024-05-17T00:23:17Z",
                        Computed: true,
                    },
                    common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                        Description: "The user id that modified the resource.\n" +
                            "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                        Computed: true,
                    },
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
							"  - example : resourceDescription",
						Computed: true,
					},
					common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "The destination IP address range in CIDR notation.\n" +
							"  - example : 10.10.10.0/24",
						Computed: true,
					},
					common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
						Description: "The identifier of the destination resource.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
						Description: "The name of the destination resource.\n" +
							"  - example : resourcaName",
						Computed: true,
					},
					common.ToSnakeCase("DestinationType"): schema.StringAttribute{
						Description: "The type of the destination.\n" +
							"  - example : ON-PREM | VPC",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the routing rule.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("SourceResourceId"): schema.StringAttribute{
						Description: "The identifier of the source resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("SourceResourceName"): schema.StringAttribute{
						Description: "The name of the source resource.\n" +
							"  - example : resourcaName",
						Computed: true,
					},
					common.ToSnakeCase("SourceType"): schema.StringAttribute{
						Description: "The type of the source.\n" +
							"  - enum :VPC, TGW\n" +
							"  - example:VPC",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the routing rule.\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the transit gateway connection belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the transit gateway connection belongs to.\n" +
							"  - example : resourceName",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcTgwRuleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

// Create the resource and sets the initial Terraform state.
func (r *vpcTgwRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan vpc.RoutingRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new routing rule
	data, err := r.client.CreateTgwRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Could not create routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRule := data.TransitGatewayRule
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(routingRule.Id)
	diags = resp.State.Set(ctx, plan)

	routingRuleModel := createRoutingRuleModel(&routingRule)

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	plan.RoutingRule = routingRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForRoutingRuleStatus(ctx, r.client, plan.TransitGatewayId.ValueString(), data.TransitGatewayRule.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become active: "+err.Error(),
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
func (r *vpcTgwRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := r.client.GetRoutingRule(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading transit gateway routing rule",
			"Could not read routing rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRuleModel := createRoutingRuleModel(&data.TransitGatewayRules[0])

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	state.RoutingRule = routingRuleObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (v vpcTgwRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r vpcTgwRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing routing rule
	err := r.client.DeleteRoutingRule(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting transit gateway routing rule",
			"Could not delete transit gateway routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForRoutingRuleStatus(ctx, r.client, state.TransitGatewayId.ValueString(), state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become deleted: "+err.Error(),
		)
		return
	}
}

func createRoutingRuleModel(data *scpvpc.TransitGatewayRule) vpc.RoutingRule {
	return vpc.RoutingRule{
		AccountId:               types.StringValue(data.AccountId),
		CreatedAt:               types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(data.CreatedBy),
		Description:             types.StringValue(data.Description),
		DestinationCidr:         types.StringValue(data.DestinationCidr),
		DestinationResourceId:   types.StringPointerValue(data.DestinationResourceId.Get()),
		DestinationResourceName: types.StringPointerValue(data.DestinationResourceName.Get()),
		DestinationType:         types.StringValue(string(data.DestinationType)),
		Id:                      types.StringValue(data.Id),
		ModifiedAt:              types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(data.ModifiedBy),
		SourceResourceId:        types.StringPointerValue(data.SourceResourceId.Get()),
		SourceResourceName:      types.StringPointerValue(data.SourceResourceName.Get()),
		SourceType:              types.StringValue(string(data.SourceType)),
		State:                   types.StringValue(string(data.State)),
		TgwConnectionVpcId:      types.StringPointerValue(data.TgwConnectionVpcId.Get()),
		TgwConnectionVpcName:    types.StringPointerValue(data.TgwConnectionVpcName.Get()),
	}
}

func waitForRoutingRuleStatus(ctx context.Context, routingRuleClient *vpc.Client, transitGatewayId string, routingRuleId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := routingRuleClient.GetRoutingRule(ctx, transitGatewayId, routingRuleId)
		if err != nil {
			return nil, "", err
		}
		if len(info.TransitGatewayRules) == 0 {
			return info, "DELETED", nil
		}
		return info, string(info.TransitGatewayRules[0].State), nil
	}, -1, -1, -1, -1)
}
