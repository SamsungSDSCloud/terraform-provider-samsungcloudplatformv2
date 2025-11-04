package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcVpcPeeringRuleResource{}
	_ resource.ResourceWithConfigure = &vpcVpcPeeringRuleResource{}
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

// Schema defines the schema for the resource.
func (r *vpcVpcPeeringRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Peering Rule",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
				Description: "VPC Peering ID",
				Required:    true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "Destination CIDR\n" +
					"  - Example : 192.168.1.0/24 \n",
				Required: true,
			},
			common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
				Description: "Destination VPC Type \n" +
					"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
					"  - Reference : VpcPeeringRuleDestinationVpcType",
				Required: true,
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),

			// Output
			common.ToSnakeCase("VpcPeeringRule"): schema.SingleNestedAttribute{
				Description: "VPC Peering Rule details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "VPC Peering Rule ID",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At\n" +
							"  - Example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By\n" +
							"  - Example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "Destination CIDR",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
						Description: "Destination VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationVpcName"): schema.StringAttribute{
						Description: "Destination VPC Name",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
						Description: "Destination VPC Type\n" +
							"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
							"  - Reference : VpcPeeringRuleDestinationVpcType",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At\n" +
							"  - Example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By\n" +
							"  - Example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
						Description: "Source VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("SourceVpcName"): schema.StringAttribute{
						Description: "Source VPC Name",
						Computed:    true,
					},
					common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
						Description: "Source VPC Type\n" +
							"  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
							"  - Reference : VpcPeeringRuleDestinationVpcType",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State\n" +
							"  - Example : CREATING | ACTIVE | DELETING | DELETED | ERROR\n" +
							"  - Reference : RoutingRuleState",
						Computed: true,
					},
					common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
						Description: "VPC Peering ID",
						Computed:    true,
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
	vpcPeeringRuleValue, _ := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	plan.VpcPeeringRule = vpcPeeringRuleValue

	err = waitForVpcPeeringRuleStatus(ctx, r.client, plan.VpcPeeringId.ValueString(), result.VpcPeeringRule.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc peering rule",
			"Error waiting for vpc peering rule to become active: "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcVpcPeeringRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	data, status, err := r.client.GetVpcPeeringRule(ctx, state.VpcPeeringId.ValueString(), rule.Id.ValueString())
	if err != nil {
		// Check if the error indicates the resource was not found
		if status == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to read VPC peering rule",
			fmt.Sprintf("An error occurred while reading VPC peering rule: %s. Details: %s", err.Error(), detail),
		)
		return
	}

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
	vpcPeeringRuleValue, _ := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
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
	})
}
