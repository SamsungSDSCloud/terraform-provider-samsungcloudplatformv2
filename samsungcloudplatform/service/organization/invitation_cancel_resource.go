package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &invitationCancelResource{}
	_ resource.ResourceWithConfigure   = &invitationCancelResource{}
	_ resource.ResourceWithImportState = &invitationCancelResource{}
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
	resp.TypeName = req.ProviderTypeName + "_organization_invitation_cancel"
}

func (r *invitationCancelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Cancel - Cancel pending invitations by IDs",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
			},
			"ids": schema.ListAttribute{
				Description: "Invitation ID. \n" +
					"  - example : ['138c2fc8c29a449dbfa8681f8f1d78e2'] \n",
				ElementType: types.StringType,
				Required:    true,
			},
			"canceled_ids": schema.ListAttribute{
				Description: "Invitation ID. \n" +
					"  - example : ['138c2fc8c29a449dbfa8681f8f1d78e2'] \n",
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
		canceledIds, canceledDiags := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)
		resp.Diagnostics.Append(canceledDiags...)
		state.CanceledIds = canceledIds

		for _, id := range result.SuccessIds {
			err = waitForInvitationCanceled(ctx, r.client, id)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error waiting for invitation cancellation",
					"Error waiting for invitation "+id+" to be canceled: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	} else {
		canceledIds, canceledDiags := types.ListValue(types.StringType, []attr.Value{})
		resp.Diagnostics.Append(canceledDiags...)
		state.CanceledIds = canceledIds
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationCancelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state invitationCancelStateData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	invitationIds := make([]string, 0, 10)
	for _, elem := range state.Ids.Elements() {
		if id, ok := elem.(types.String); ok {
			invitationIds = append(invitationIds, id.ValueString())
		}
	}

	request := organization.OrganizationInvitationsDataSource{
		OrganizationId: types.StringValue(orgId),
	}
	result, err := r.client.GetOrganizationInvitations(ctx, request)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	existingIds := make(map[string]bool)
	for _, inv := range result.GetOrganizationInvitations() {
		existingIds[inv.Id] = true
	}

	for _, id := range invitationIds {
		if !existingIds[id] {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationCancelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *invitationCancelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty.
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

func (r *invitationCancelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) < 1 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: organization_id:invitation_id1,invitation_id2,...",
		)
		return
	}

	orgId := parts[0]
	var ids []string
	if len(parts) > 1 {
		ids = strings.Split(parts[1], ",")
	}
	canceledIds, canceledDiags := types.ListValue(types.StringType, []attr.Value{})
	resp.Diagnostics.Append(canceledDiags...)
	var idsList types.List
	var idsDiags diag.Diagnostics
	idsList, idsDiags = types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(idsDiags...)

	state := invitationCancelStateData{
		OrganizationId: types.StringValue(orgId),
		Ids:            idsList,
		CanceledIds:    canceledIds,
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func waitForInvitationCanceled(ctx context.Context, orgClient *organization.Client, invitationId string) error {
	return client.WaitForStatus(ctx, nil, []string{"INVITING"}, []string{"CANCELED"}, func() (interface{}, string, error) {
		result, err := orgClient.GetAccountInvitations(ctx)
		if err != nil {
			return nil, "", err
		}
		for _, inv := range result.GetAccountInvitations() {
			if inv.Id == invitationId {
				return inv, string(inv.State), nil
			}
		}
		return "CANCELED", "CANCELED", nil
	}, -1, -1, -1, -1)
}
