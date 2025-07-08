package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/loadbalancer/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLbMemberResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLbMemberResource{}
)

// NewLoadBalancerLbMemberResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLbMemberResource() resource.Resource {
	return &loadbalancerLbMemberResource{}
}

// loadbalancerLbMemberResource is the data source implementation.
type loadbalancerLbMemberResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_member"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lb Member.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
				Description: "LbServerGroupId",
				Required:    true,
			},
			common.ToSnakeCase("LbMember"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Member.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Computed:    true,
					},
					common.ToSnakeCase("Uuid"): schema.StringAttribute{
						Description: "Uuid",
						Computed:    true,
					},
					common.ToSnakeCase("ObjectId"): schema.StringAttribute{
						Description: "ObjectId",
						Computed:    true,
					},
					common.ToSnakeCase("ObjectType"): schema.StringAttribute{
						Description: "ObjectType",
						Computed:    true,
					},
					common.ToSnakeCase("MemberWeight"): schema.Int32Attribute{
						Description: "MemberWeight",
						Computed:    true,
					},
					common.ToSnakeCase("MemberState"): schema.StringAttribute{
						Description: "MemberState",
						Computed:    true,
					},
					common.ToSnakeCase("MemberPort"): schema.Int32Attribute{
						Description: "MemberPort",
						Computed:    true,
					},
					common.ToSnakeCase("MemberIp"): schema.StringAttribute{
						Description: "MemberIp",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
						Description: "LbServerGroupId",
						Computed:    true,
					},
				},
			},
			common.ToSnakeCase("LbMemberCreate"): schema.SingleNestedAttribute{
				Description: "Create Lb Member.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("ObjectId"): schema.StringAttribute{
						Description: "ObjectId",
						Optional:    true,
					},
					common.ToSnakeCase("ObjectType"): schema.StringAttribute{
						Description: "ObjectType",
						Optional:    true,
					},
					common.ToSnakeCase("MemberPort"): schema.Int32Attribute{
						Description: "MemberPort",
						Optional:    true,
					},
					common.ToSnakeCase("MemberIp"): schema.StringAttribute{
						Description: "MemberIp",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("MemberWeight"): schema.Int32Attribute{
						Description: "MemberWeight",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LbMembersResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Member
	data, err := r.client.CreateLbMember(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Member",
			"Could not create Lb Member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	//for _, member := range data.Members {
	member := data.Members[0]
	plan.Id = types.StringValue(member.Id)

	// Map response body to schema and populate Computed attribute values
	lbMemberModel := createLbMemberModel(member)
	lbMemberOjbectValue, _ := types.ObjectValueFrom(ctx, lbMemberModel.AttributeTypes(), lbMemberModel)
	plan.LbMember = lbMemberOjbectValue

	// Set state to fully populated data
	_ = resp.State.Set(ctx, plan)

	err = waitForMemberStatus(ctx, r.client, member.LbServerGroupId, member.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating lb server group member",
			"Error waiting for lb server group member to become active: "+err.Error(),
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
	//}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbMembersResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Member
	data, err := r.client.GetLbMember(ctx, state.LbServerGroupId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Member",
			"Could not create Lb Member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbMember := data.Member.Get()
	lbMemberModel := loadbalancer.LbMemberDetail{
		LbServerGroupId: types.StringValue(lbMember.LbServerGroupId),
		Name:            types.StringValue(lbMember.Name),
		State:           types.StringValue(string(lbMember.State)),
		MemberIp:        types.StringValue(lbMember.MemberIp),
		MemberPort:      types.Int32Value(lbMember.MemberPort),
		MemberState:     types.StringValue(lbMember.MemberState),
		MemberWeight:    types.Int32Value(lbMember.MemberWeight),
		ObjectType:      types.StringValue(string(lbMember.ObjectType)),
		ObjectId:        virtualserverutil.ToNullableStringValue(lbMember.ObjectId.Get()),
		SubnetId:        types.StringValue(lbMember.SubnetId),
		Uuid:            types.StringValue(lbMember.Uuid),
		ModifiedBy:      types.StringValue(lbMember.ModifiedBy),
		ModifiedAt:      types.StringValue(lbMember.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(lbMember.CreatedBy),
		CreatedAt:       types.StringValue(lbMember.CreatedAt.Format(time.RFC3339)),
	}

	lbMemberOjbectValue, diags := types.ObjectValueFrom(ctx, lbMemberModel.AttributeTypes(), lbMemberModel)
	state.LbMember = lbMemberOjbectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LbMemberResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	data, err := r.client.UpdateLbMember(ctx, state.LbServerGroupId.ValueString(), state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Server Group Member",
			"Could not create Lb Server Group Member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbMemberModel := updateLbMemberModel(data)

	lbMemberObjectValue, diags := types.ObjectValueFrom(ctx, lbMemberModel.AttributeTypes(), lbMemberModel)
	state.LbMember = lbMemberObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, state)

	err = waitForMemberStatus(ctx, r.client, data.Member.Get().LbServerGroupId, data.Member.Get().Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating lb server group member",
			"Error waiting for lb server group member to become active: "+err.Error(),
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
func (r *loadbalancerLbMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LbMembersResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Member
	err := r.client.DeleteLbMember(ctx, state.LbServerGroupId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Member",
			"Could not delete lb member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createLbMemberModel(data scploadbalancer.Member) loadbalancer.LbMemberDetail {
	lbMember := data
	return loadbalancer.LbMemberDetail{
		LbServerGroupId: types.StringValue(lbMember.LbServerGroupId),
		Name:            types.StringValue(lbMember.Name),
		State:           types.StringValue(string(lbMember.State)),
		MemberIp:        types.StringValue(lbMember.MemberIp),
		MemberPort:      types.Int32Value(lbMember.MemberPort),
		MemberState:     types.StringValue(lbMember.MemberState),
		MemberWeight:    types.Int32Value(lbMember.MemberWeight),
		ObjectType:      types.StringValue(string(lbMember.ObjectType)),
		ObjectId:        virtualserverutil.ToNullableStringValue(lbMember.ObjectId.Get()),
		SubnetId:        types.StringValue(lbMember.SubnetId),
		Uuid:            types.StringValue(lbMember.Uuid),
		ModifiedBy:      types.StringValue(lbMember.ModifiedBy),
		ModifiedAt:      types.StringValue(lbMember.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(lbMember.CreatedBy),
		CreatedAt:       types.StringValue(lbMember.CreatedAt.Format(time.RFC3339)),
	}
}
func updateLbMemberModel(data *scploadbalancer.MemberShowResponse) loadbalancer.LbMemberDetail {
	lbMember := data.Member.Get()

	return loadbalancer.LbMemberDetail{
		LbServerGroupId: types.StringValue(lbMember.LbServerGroupId),
		Name:            types.StringValue(lbMember.Name),
		State:           types.StringValue(string(lbMember.State)),
		MemberIp:        types.StringValue(lbMember.MemberIp),
		MemberPort:      types.Int32Value(lbMember.MemberPort),
		MemberState:     types.StringValue(lbMember.MemberState),
		MemberWeight:    types.Int32Value(lbMember.MemberWeight),
		ObjectType:      types.StringValue(string(lbMember.ObjectType)),
		ObjectId:        virtualserverutil.ToNullableStringValue(lbMember.ObjectId.Get()),
		SubnetId:        types.StringValue(lbMember.SubnetId),
		Uuid:            types.StringValue(lbMember.Uuid),
		ModifiedBy:      types.StringValue(lbMember.ModifiedBy),
		ModifiedAt:      types.StringValue(lbMember.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(lbMember.CreatedBy),
		CreatedAt:       types.StringValue(lbMember.CreatedAt.Format(time.RFC3339)),
	}
}

func waitForMemberStatus(ctx context.Context, loadbalancerClient *loadbalancer.Client, lbServerGroupId string, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := loadbalancerClient.GetLbMember(ctx, lbServerGroupId, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.GetMember().State), nil
	})
}
