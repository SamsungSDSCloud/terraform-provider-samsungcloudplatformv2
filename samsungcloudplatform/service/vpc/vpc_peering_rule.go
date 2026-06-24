package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcVpcPeeringRuleResource{}
	_ resource.ResourceWithConfigure   = &vpcVpcPeeringRuleResource{}
	_ resource.ResourceWithImportState = &vpcVpcPeeringRuleResource{}
)

// NewVpcVpcPeeringRuleResource is a helper function to simplify the provider implementation.
func NewVpcVpcPeeringRuleResource() resource.Resource {
	return &vpcVpcPeeringRuleResource{}
}

// vpcVpcPeeringRuleResource is the resource implementation.
type vpcVpcPeeringRuleResource struct {
	_config *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *vpcVpcPeeringRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering_rule"
}

func (r *vpcVpcPeeringRuleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("vpc_peering_id"), request, response)

	parts := strings.Split(request.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		response.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: vpcPeeringId/ruleId, got: %q", request.ID),
		)
		return
	}

	response.State.SetAttribute(ctx, path.Root("vpc_peering_id"), types.StringValue(parts[0]))
	response.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))
}

// Schema defines the schema for the resource.
func (r *vpcVpcPeeringRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Peering Rule",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
				Description: "The identifier of the VPC peering.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "The destination IP address range in CIDR notation.\n  - Example : 192.168.1.0/24",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
				Description: "The type of the destination VPC.\n" +
					"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
					"  - Reference : VpcPeeringRuleDestinationVpcType",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the VPC peering rule.\n" +
					"  - example : 6fdb3ffbe5ae4878a1ef6f55e208aa05",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// Output
			common.ToSnakeCase("VpcPeeringRule"): schema.SingleNestedAttribute{
				Description: "VPC Peering Rule details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the VPC peering rule.\n" +
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
					}, common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "The destination IP address range in CIDR notation.\n" +
							"  - example : 10.10.10.0/24",
						Computed: true,
					},
					common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
						Description: "The identifier of the destination VPC.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("DestinationVpcName"): schema.StringAttribute{
						Description: "The name of the destination VPC.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
						Description: "The type of the destination VPC.\n" +
							"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
							"  - Reference : VpcPeeringRuleDestinationVpcType",
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
					}, common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
						Description: "The identifier of the source VPC.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("SourceVpcName"): schema.StringAttribute{
						Description: "The name of the source VPC.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
						Description: "The type of the source VPC.\n" +
							"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
							"  - Reference : VpcPeeringRuleDestinationVpcType",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the VPC peering rule.\n" +
							"  - Example : CREATING | ACTIVE | DELETING | DELETED | ERROR\n" +
							"  - Reference : RoutingRuleState",
						Computed: true,
					},
					common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
						Description: "The identifier of the VPC peering.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *vpcVpcPeeringRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcVpcPeeringRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpc.VpcPeeringRuleResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateVpcPeeringRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create VPC peering rule",
			fmt.Sprintf("An error occurred while creating VPC peering rule: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Set the nested structure in the plan
	vpcPeeringModel := vpc.VpcPeeringRule{
		Id:                 types.StringValue(result.VpcPeeringRule.Id),
		CreatedAt:          types.StringValue(result.VpcPeeringRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:          types.StringValue(result.VpcPeeringRule.CreatedBy),
		DestinationVpcId:   types.StringValue(result.VpcPeeringRule.DestinationVpcId),
		DestinationVpcName: types.StringValue(result.VpcPeeringRule.DestinationVpcName),
		DestinationCidr:    types.StringValue(result.VpcPeeringRule.DestinationCidr),
		ModifiedAt:         types.StringValue(result.VpcPeeringRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:         types.StringValue(result.VpcPeeringRule.ModifiedBy),
		SourceVpcId:        types.StringValue(result.VpcPeeringRule.SourceVpcId),
		SourceVpcName:      types.StringValue(result.VpcPeeringRule.SourceVpcName),
		SourceVpcType:      types.StringValue(string(result.VpcPeeringRule.SourceVpcType)),
		State:              types.StringValue(string(result.VpcPeeringRule.State)),
		DestinationVpcType: types.StringValue(string(result.VpcPeeringRule.DestinationVpcType)),
		VpcPeeringId:       types.StringValue(plan.VpcPeeringId.ValueString()),
	}
	vpcPeeringRuleValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	resp.Diagnostics.Append(diags...)

	plan.VpcPeeringRule = vpcPeeringRuleValue

	err = waitForVpcPeeringRuleStatus(ctx, r.client, plan.VpcPeeringId.ValueString(), result.VpcPeeringRule.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc peering rule",
			"Error waiting for vpc peering rule to become active: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)

	// Refresh resource state
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
func (r *vpcVpcPeeringRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpc.VpcPeeringRuleResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If VpcPeeringRule is null (e.g. after import), use the rule ID from state
	ruleId := types.StringNull()
	if !state.VpcPeeringRule.IsNull() && !state.VpcPeeringRule.IsUnknown() {
		var rule vpc.VpcPeeringRule
		errR := state.VpcPeeringRule.As(ctx, &rule, basetypes.ObjectAsOptions{})
		if errR != nil {
			resp.Diagnostics.AddError(
				"Failed to parse VPC peering rule",
				fmt.Sprintf("An error occurred while parsing VPC peering rule: %s", errR),
			)
			return
		}
		ruleId = rule.Id
	} else {
		ruleId = state.Id
	}

	data, status, err := r.client.GetVpcPeeringRule(ctx, state.VpcPeeringId.ValueString(), ruleId.ValueString())
	if err != nil {
		// Check if the error indicates the resource was not found
		if status == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to read VPC peering rule with"+" VpcPeeringId  "+state.VpcPeeringId.ValueString()+" ruleId "+ruleId.ValueString(),
			fmt.Sprintf("An error occurred while reading VPC peering rule: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Rewrite configurable top-level input fields from API response to detect drift
	state.DestinationCidr = types.StringValue(data.DestinationCidr)
	state.DestinationVpcType = types.StringValue(string(data.DestinationVpcType))

	// Set the nested structure in the plan
	vpcPeeringModel := vpc.VpcPeeringRule{
		Id:                 types.StringValue(data.Id),
		CreatedAt:          types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:          types.StringValue(data.CreatedBy),
		DestinationVpcId:   types.StringValue(data.DestinationVpcId),
		DestinationVpcName: types.StringValue(data.DestinationVpcName),
		DestinationCidr:    types.StringValue(data.DestinationCidr),
		ModifiedAt:         types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:         types.StringValue(data.ModifiedBy),
		SourceVpcId:        types.StringValue(data.SourceVpcId),
		SourceVpcName:      types.StringValue(data.SourceVpcName),
		SourceVpcType:      types.StringValue(string(data.SourceVpcType)),
		State:              types.StringValue(string(data.State)),
		DestinationVpcType: types.StringValue(string(data.DestinationVpcType)),
		VpcPeeringId:       types.StringValue(data.VpcPeeringId),
	}
	vpcPeeringRuleValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	resp.Diagnostics.Append(diags...)

	state.VpcPeeringRule = vpcPeeringRuleValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcVpcPeeringRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpc.VpcPeeringRuleResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rule vpc.VpcPeeringRule
	errR := state.VpcPeeringRule.As(ctx, &rule, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse VPC peering rule",
			fmt.Sprintf("An error occurred while parsing VPC peering rule: %s", errR),
		)
		return
	}

	err := r.client.DeleteVpcPeeringRule(ctx, state.VpcPeeringId.ValueString(), rule.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to delete VPC peering rule",
			fmt.Sprintf("An error occurred while deleting VPC peering rule: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	err = waitForVpcPeeringRuleStatus(ctx, r.client, state.VpcPeeringId.ValueString(), rule.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting VPC peering rule",
			"Error waiting for VPC peering rule to become deleted: "+err.Error(),
		)
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcVpcPeeringRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// VPC peering rule does not support update operations
	// This is a no-op implementation
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"VPC peering rule resources do not support update operations. The resource will not be updated.",
	)
}

func waitForVpcPeeringRuleStatus(ctx context.Context, vpcClient *vpc.Client, vpcPeeringId string, ruleId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		data, status, err := vpcClient.GetVpcPeeringRule(ctx, vpcPeeringId, ruleId)
		if err != nil {
			if status == 404 {
				return vpc.VpcPeeringRule{}, "DELETED", nil
			}
			return nil, "", err
		}
		return data, string(data.State), nil
	}, -1, -1, -1, -1)
}
