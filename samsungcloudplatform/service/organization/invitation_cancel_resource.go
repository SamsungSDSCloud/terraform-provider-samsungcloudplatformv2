package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &invitationCancelResource{}
	_ resource.ResourceWithConfigure = &invitationCancelResource{}
)

func NewInvitationCancelResource() resource.Resource {
	return &invitationCancelResource{}
}

type invitationCancelResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *invitationCancelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invitation_cancel"
}

func (r *invitationCancelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Cancel - Cancel pending invitations by IDs",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Organization ID",
				Required:    true,
			},
			"ids": schema.ListAttribute{
				Description: "Invitation IDs to cancel",
				ElementType: types.StringType,
				Required:    true,
			},
			"canceled_ids": schema.ListAttribute{
				Description: "Successfully canceled invitation IDs",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (r *invitationCancelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

func (r *invitationCancelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData invitationCancelStateData
	diags := req.Plan.Get(ctx, &planData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitationReq := organization.InvitationCancelRequest{
		OrganizationId: planData.OrganizationId,
		Ids:            planData.Ids,
	}
	tflog.Debug(ctx, "CancelInvitations request", map[string]interface{}{
		"org_id": planData.OrganizationId.ValueString(),
		"ids":    planData.Ids,
	})
	result, err := r.client.CancelInvitations(ctx, invitationReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error canceling invitation",
			err.Error(),
		)
		return
	}

	state := invitationCancelStateData{
		OrganizationId: planData.OrganizationId,
		Ids:            planData.Ids,
	}
	if len(result.SuccessIds) > 0 {
		canceledIds, _ := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)
		state.CanceledIds = canceledIds
	} else {
		state.CanceledIds = types.ListValueMust(types.StringType, []attr.Value{})
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationCancelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *invitationCancelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *invitationCancelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

type invitationCancelStateData struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Ids            types.List   `tfsdk:"ids"`
	CanceledIds    types.List   `tfsdk:"canceled_ids"`
}

func (o invitationCancelStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"ids": types.ListType{
			ElemType: types.StringType,
		},
		"canceled_ids": types.ListType{
			ElemType: types.StringType,
		},
	}
}
