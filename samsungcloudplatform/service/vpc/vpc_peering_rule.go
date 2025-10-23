package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		Description: "VPC Peering Rule Resource",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
				Description: "VPC Peering ID",
				Required:    true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "Destination CIDR",
				Required:    true,
			},
			common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
				Description: "Destination VPC Type",
				Required:    true,
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),

			// Output
			common.ToSnakeCase("ID"): schema.StringAttribute{
				Description: "VPC Peering Rule ID",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "Created At",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "Created By",
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
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "VPC Peering Rule ID",
				Computed:    true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description: "Modified At",
				Computed:    true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "Modified By",
				Computed:    true,
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
				Description: "Source VPC Type",
				Computed:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Computed:    true,
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

	plan.Id = types.StringValue(result.VpcPeeringRule.Id)
	plan.CreatedAt = types.StringValue(result.VpcPeeringRule.CreatedAt.Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(result.VpcPeeringRule.CreatedBy)
	plan.DestinationVpcId = types.StringValue(result.VpcPeeringRule.DestinationVpcId)
	plan.DestinationVpcName = types.StringValue(result.VpcPeeringRule.DestinationVpcName)
	plan.ModifiedAt = types.StringValue(result.VpcPeeringRule.ModifiedAt.Format(time.RFC3339))
	plan.ModifiedBy = types.StringValue(result.VpcPeeringRule.ModifiedBy)
	plan.SourceVpcId = types.StringValue(result.VpcPeeringRule.SourceVpcId)
	plan.SourceVpcName = types.StringValue(result.VpcPeeringRule.SourceVpcName)
	plan.SourceVpcType = types.StringValue(string(result.VpcPeeringRule.SourceVpcType))
	plan.State = types.StringValue(string(result.VpcPeeringRule.State))

	diags = resp.State.Set(ctx, plan)

	err = waitForVpcPeeringRuleStatus(ctx, r.client, plan.VpcPeeringId.ValueString(), plan.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc peering rule",
			"Error waiting for vpc peering rule to become active: "+err.Error(),
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
func (r *vpcVpcPeeringRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpc.VpcPeeringRuleResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, status, err := r.client.GetVpcPeeringRule(ctx, state.VpcPeeringId.ValueString(), state.Id.ValueString())
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

	state.Id = types.StringValue(data.Id)
	state.CreatedAt = types.StringValue(data.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(data.CreatedBy)
	state.DestinationVpcId = types.StringValue(data.DestinationVpcId)
	state.DestinationVpcName = types.StringValue(data.DestinationVpcName)
	state.ModifiedAt = types.StringValue(data.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(data.ModifiedBy)
	state.SourceVpcId = types.StringValue(data.SourceVpcId)
	state.SourceVpcName = types.StringValue(data.SourceVpcName)
	state.SourceVpcType = types.StringValue(string(data.SourceVpcType))
	state.State = types.StringValue(string(data.State))

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

	err := r.client.DeleteVpcPeeringRule(ctx, state.VpcPeeringId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to delete VPC peering rule",
			fmt.Sprintf("An error occurred while deleting VPC peering rule: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	err = waitForVpcPeeringRuleStatus(ctx, r.client, state.VpcPeeringId.ValueString(), state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc endpoint",
			"Error waiting for vpc endpoint to become deleted: "+err.Error(),
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
