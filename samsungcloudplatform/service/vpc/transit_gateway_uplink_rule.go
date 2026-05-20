package vpc

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcTgwUplinkRuleResource{}
	_ resource.ResourceWithConfigure = &vpcTgwUplinkRuleResource{}
)

// NewVpcTgwUplinkRuleResource is a helper function to simplify the provider implementation.
func NewVpcTgwUplinkRuleResource() resource.Resource {
	return &vpcTgwUplinkRuleResource{}
}

type vpcTgwUplinkRuleResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcTgwUplinkRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_uplink_rule"
}

// Configure adds the provider configured client to the data source.
func (r *vpcTgwUplinkRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

func (r *vpcTgwUplinkRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "vpc",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Description\n  - maxLength: 50\n  - example: Routing Rule Description",
				MarkdownDescription: "Description\n  - maxLength: 50\n  - example: Routing Rule Description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(""),
			},
			"destination_cidr": schema.StringAttribute{
				Required:            true,
				Description:         "Destination CIDR\n  - example: 192.167.5.0/24",
				MarkdownDescription: "Destination CIDR\n  - example: 192.167.5.0/24",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_type": schema.StringAttribute{
				Required:            true,
				Description:         "- enum: [\"TGW\",\"ON_PREMISE\"]",
				MarkdownDescription: "- enum: [\"TGW\",\"ON_PREMISE\"]",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TGW",
						"ON_PREMISE",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"transit_gateway_id": schema.StringAttribute{
				Required:            true,
				Description:         "Transit Gateway ID\n  - example: fe860e0af0c04dcd8182b84f907f31f4",
				MarkdownDescription: "Transit Gateway ID\n  - example: fe860e0af0c04dcd8182b84f907f31f4",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"transit_gateway_rule": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Description\n  - example: Routing Rule Description",
						MarkdownDescription: "Description\n  - example: Routing Rule Description",
					},
					"destination_cidr": schema.StringAttribute{
						Computed:            true,
						Description:         "Destination CIDR\n  - example: 192.167.5.0/24",
						MarkdownDescription: "Destination CIDR\n  - example: 192.167.5.0/24",
					},
					"destination_type": schema.StringAttribute{
						Computed:            true,
						Description:         "- enum: [\"TGW\",\"ON_PREMISE\"]",
						MarkdownDescription: "- enum: [\"TGW\",\"ON_PREMISE\"]",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Routing Rule ID\n  - example: 43772aff4539403d9ba74bf1fdaa00c8",
						MarkdownDescription: "Routing Rule ID\n  - example: 43772aff4539403d9ba74bf1fdaa00c8",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
						MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
					},
				},
				CustomType: vpc.TransitGatewayRuleType{
					ObjectType: types.ObjectType{
						AttrTypes: vpc.TransitGatewayRuleValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcTgwUplinkRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.TransitGatewayUplinkValue
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new uplink rule
	data, err := r.client.CreateTgwUplinkRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Transit Gateway Uplink Rule",
			"Could not create uplink rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	rule := data.TransitGatewayRule

	transitGatewayRule := vpc.TransitGatewayRuleValue{
		Description:     types.StringValue(rule.Description),
		DestinationCidr: types.StringValue(rule.DestinationCidr),
		DestinationType: types.StringValue(string(rule.DestinationType)),
		Id:              types.StringValue(rule.Id),
		State:           types.StringValue(string(rule.State)),
	}

	err = waitForTGWRuleStatus(ctx, r.client, plan.TransitGatewayId.ValueString(), rule.Id, []string{common.CreatingState}, []string{common.ActiveState})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Transit Gateway Uplink Rule",
			"Could not delete GPU Node("+rule.Id+")unexpected error: "+err.Error()+"\n|"+detail,
		)
		return
	}

	transitGatewayRule.State = types.StringValue(common.ActiveState)

	plan.TransitGatewayRule, _ = types.ObjectValueFrom(ctx, transitGatewayRule.AttributeTypes(ctx), transitGatewayRule)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcTgwUplinkRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state vpc.TransitGatewayUplinkValue
	diags := req.State.Get(ctx, &state)

	var transitGatewayId, ruleId types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("transit_gateway_id"), &transitGatewayId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("transit_gateway_rule").AtName("id"), &ruleId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := r.client.GetTGWRuleList(ctx, vpc.TransitGatewayRuleDataSource{Id: ruleId, TransitGatewayId: transitGatewayId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Transit Gateway Uplink Rule",
			"Could not Read Transit Gateway Uplink Rule ID "+ruleId.ValueString()+": "+err.Error()+"\n|Reason: "+detail,
		)
		return
	}

	if data.Count == 0 {
		return
	}

	rule := data.TransitGatewayRules[0]

	transitGatewayRule := vpc.TransitGatewayRuleValue{
		Description:     types.StringValue(rule.Description),
		DestinationCidr: types.StringValue(rule.DestinationCidr),
		DestinationType: types.StringValue(string(rule.DestinationType)),
		Id:              types.StringValue(rule.Id),
		State:           types.StringValue(string(rule.State)),
	}

	state.TransitGatewayRule, _ = types.ObjectValueFrom(ctx, transitGatewayRule.AttributeTypes(ctx), transitGatewayRule)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcTgwUplinkRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update has not been suppported yet.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcTgwUplinkRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var transitGatewayId, ruleId types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("transit_gateway_id"), &transitGatewayId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("transit_gateway_rule").AtName("id"), &ruleId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpc
	err := r.client.DeleteTgwUplinkRule(ctx, transitGatewayId.ValueString(), ruleId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting uplink rule",
			"Could not delete uplink rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForTGWRuleStatus(ctx, r.client, transitGatewayId.ValueString(), ruleId.ValueString(), []string{common.DeletingState}, []string{common.DeletedState})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Transit Gateway Uplink Rule",
			"Could not delete GPU Node("+ruleId.ValueString()+")unexpected error: "+err.Error()+"\n|"+detail,
		)
		return
	}

}

func waitForTGWRuleStatus(ctx context.Context, vpcClient *vpc.Client, tgwId string, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {

		info, _, err := vpcClient.GetTGWRuleList(ctx, vpc.TransitGatewayRuleDataSource{Id: types.StringValue(id), TransitGatewayId: types.StringValue(tgwId)})
		if err != nil {
			return nil, "", err
		}
		if info.Count == 0 {
			return "", common.DeletedState, nil
		}
		return info, string(info.TransitGatewayRules[0].State), nil
	})
}
