package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcPeeringApprovalReource{}
	_ resource.ResourceWithConfigure = &vpcPeeringApprovalReource{}
)

// VpcPeeringApprovalResource is a helper function to simplify the provider implementation.
func NewVpcPeeringApprovalResource() resource.Resource {
	return &vpcPeeringApprovalReource{}
}

// VpcPeeringApprovalResource is the data source implementation.
type vpcPeeringApprovalReource struct {
	_config *scpsdk.Configuration
	client  *vpcv1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPeeringApprovalReource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering_approval"
}

// Schema defines the schema for the data source.
func (d *vpcPeeringApprovalReource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Peering Approval",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("VpcPeeringID"): schema.StringAttribute{
				Description: "VPC Peering ID",
				Required:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Approval Type\n" +
					"  - enum:  CREATE_APPROVE, CREATE_CANCEL, CREATE_REJECT, CREATE_RE_REQUEST, DELETE_APPROVE, DELETE_CANCEL, DELETE_REJECT",
				Required: true,
			},
			common.ToSnakeCase("VpcPeering"): schema.SingleNestedAttribute{
				Description: "VpcPeering",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountType"): schema.StringAttribute{
						Description: "AccountType",
						Computed:    true,
					},
					common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
						Description: "ApproverVpcAccountId",
						Computed:    true,
					},
					common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
						Description: "ApproverVpcId",
						Computed:    true,
					},
					common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
						Description: "ApproverVpcName",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Required:    true,
					},
					common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
						Description: "RequesterVpcAccountId",
						Computed:    true,
					},
					common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
						Description: "RequesterVpcId",
						Computed:    true,
					},
					common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
						Description: "RequesterVpcName",
						Computed:    true,
					},
					common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
						Description: "DeleteRequesterAccountId",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State\n" +
							"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPeeringApprovalReource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.client = inst.Client.VpcV1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPeeringApprovalReource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1.VpcPeeringApprovalResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpcPeering(ctx, state.VpcPeeringID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc peering",
			"Could not read vpc peering ID "+state.VpcPeeringID.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, _ := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)

	state.VpcPeering = vpcObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcPeeringApprovalReource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state vpcv1.VpcPeeringApprovalResource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ApprovalVpcPeering(ctx, state.VpcPeeringID.ValueString(), state.Type.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to get VPC Peering Approval",
			fmt.Sprintf("Error: %s, Detail: %s", err.Error(), detail),
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, _ := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)

	state.VpcPeering = vpcObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)

	mappingApprovalTypeToStatus := map[string][]string{
		"CREATE_APPROVE":    {"ACTIVE"},
		"CREATE_CANCEL":     {"CANCELED"},
		"CREATE_REJECT":     {"CREATE_REJECT"},
		"CREATE_RE_REQUEST": {"CREATING_REQUESTING"},
		"DELETE_APPROVE":    {"DELETED"},
		"DELETE_CANCEL":     {"ACTIVE"},
		"DELETE_REJECT":     {"ACTIVE"},
	}
	expectedStatus := mappingApprovalTypeToStatus[state.Type.ValueString()]

	if state.Type.ValueString() == "DELETE_APPROVE" {
		err = waitForVpcPeeringStatus(ctx, r.client, state.VpcPeeringID.ValueString(), []string{}, expectedStatus)
		if err != nil && !strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError(
				"Error deleting vpc peering",
				"Error waiting for vpc peering to become deleted: "+err.Error(),
			)
			return
		}
		// Set VpcPeering and VpcPeeringID to null
		state.VpcPeering = types.ObjectNull(vpcPeeringModel.AttributeTypes())
		// Set the updated state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	} else {
		err = waitForVpcPeeringStatus(ctx, r.client, state.VpcPeeringID.ValueString(), []string{}, expectedStatus)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error approve vpc peering",
				"Error waiting for vpc peering to be expected value: "+err.Error(),
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
}

func (r *vpcPeeringApprovalReource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1.VpcPeeringApprovalResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing VpcPeering
	err := r.client.DeleteVpcPeering(ctx, state.VpcPeeringID.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc peering",
			"Could not delete vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcPeeringStatus(ctx, r.client, state.VpcPeeringID.ValueString(), []string{}, []string{"DELETED", "DELETING_REQUESTING"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc peering",
			"Error waiting for vpc peering to become deleted: "+err.Error(),
		)
		return
	}
}

func (r *vpcPeeringApprovalReource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state vpcv1.VpcPeeringApprovalResource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ApprovalVpcPeering(ctx, state.VpcPeeringID.ValueString(), state.Type.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to get VPC Peering Approval",
			fmt.Sprintf("Error: %s, Detail: %s", err.Error(), detail),
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, _ := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)

	state.VpcPeering = vpcObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)

	mappingApprovalTypeToStatus := map[string][]string{
		"CREATE_APPROVE":    {"ACTIVE"},
		"CREATE_CANCEL":     {"CANCELED"},
		"CREATE_REJECT":     {"CREATE_REJECT"},
		"CREATE_RE_REQUEST": {"CREATING_REQUESTING"},
		"DELETE_APPROVE":    {"DELETED"},
		"DELETE_CANCEL":     {"ACTIVE"},
		"DELETE_REJECT":     {"ACTIVE"},
	}
	expectedStatus := mappingApprovalTypeToStatus[state.Type.ValueString()]

	if state.Type.ValueString() == "DELETE_APPROVE" {
		err = waitForVpcPeeringStatus(ctx, r.client, state.VpcPeeringID.ValueString(), []string{}, expectedStatus)
		if err != nil && !strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError(
				"Error deleting vpc peering",
				"Error waiting for vpc peering to become deleted: "+err.Error(),
			)
			return
		}
		// Set VpcPeering and VpcPeeringID to null
		state.VpcPeering = types.ObjectNull(vpcPeeringModel.AttributeTypes())
		// Set the updated state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	} else {
		err = waitForVpcPeeringStatus(ctx, r.client, state.VpcPeeringID.ValueString(), []string{}, expectedStatus)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error approve vpc peering",
				"Error waiting for vpc peering to be expected value: "+err.Error(),
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
}
