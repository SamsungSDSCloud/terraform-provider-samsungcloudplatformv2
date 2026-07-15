package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &invitationDeclineResource{}
	_ resource.ResourceWithConfigure   = &invitationDeclineResource{}
	_ resource.ResourceWithImportState = &invitationDeclineResource{}
)

func NewInvitationDeclineResource() resource.Resource {
	return &invitationDeclineResource{}
}

type invitationDeclineResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *invitationDeclineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_invitation_decline"
}

func (r *invitationDeclineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Decline - Use a separate provider alias with target account credentials",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Invitation ID to decline. \n" +
					"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
				Required: true,
			},
			"invitation": schema.SingleNestedAttribute{
				Description: "Invitation details. \n" +
					"  - example : '{id: 138c2fc8c29a449dbfa8681f8f1d78e2, organization_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, ...}' \n",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Invitation ID to decline. \n" +
							"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
						Computed: true,
					},
					"state": schema.StringAttribute{
						Description: "Invitation State. \n" +
							"  - example : 'REFUSED' \n" +
							"  - allowed_values : ['INVITING', 'REFUSED', 'INVITED', 'CANCELED', 'EXPIRED'] \n",
						Computed: true,
					},
					"target_account_id": schema.StringAttribute{
						Description: "Invitation Receiver Account ID. \n" +
							"  - example : '338c2fc8c29a449dbfa8681f8f1d78e5' \n",
						Computed: true,
					},
					"organization_id": schema.StringAttribute{
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						Computed: true,
					},
					"master_account_id": schema.StringAttribute{
						Description: "Unique identifier of the master account that manages the organization. \n" +
							"  - example : '9f8e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
						Computed: true,
					},
					"requested_time": schema.StringAttribute{
						Description: "Invitation Request Datetime. \n" +
							"  - example : '2024-04-17T12:34:56.789Z' \n",
						Computed: true,
					},
					"expired_time": schema.StringAttribute{
						Description: "Invitation Expired Datetime. \n" +
							"  - example : '2024-04-30T12:34:56.789Z' \n",
						Computed: true,
					},
					"created_at": schema.StringAttribute{
						Description: "Timestamp when the invitation was sent. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Description: "User who sent the invitation. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Description: "Timestamp when the invitation was declined. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Description: "User who declined the invitation. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *invitationDeclineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *invitationDeclineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData invitationDeclineStateData
	diags := req.Plan.Get(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.DeclineInvitation(ctx, stateData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error declining invitation",
			err.Error(),
		)
		return
	}

	invitation := result.GetInvitation()
	invitationObj := map[string]attr.Value{
		"id":                types.StringValue(invitation.Id),
		"state":             types.StringValue(string(invitation.State)),
		"target_account_id": types.StringValue(invitation.TargetAccountId),
		"organization_id":   types.StringValue(invitation.OrganizationId),
		"master_account_id": types.StringValue(invitation.MasterAccountId),
		"requested_time":    types.StringValue(invitation.RequestedTime.Format(time.RFC3339)),
		"expired_time":      types.StringValue(invitation.ExpiredTime.Format(time.RFC3339)),
		"created_at":        types.StringValue(invitation.CreatedAt.Format(time.RFC3339)),
		"created_by":        types.StringValue(invitation.CreatedBy),
		"modified_at":       types.StringValue(invitation.ModifiedAt.Format(time.RFC3339)),
		"modified_by":       types.StringValue(invitation.ModifiedBy),
	}

	invitationAttr, invitationDiags := types.ObjectValue(invitationDeclineInvitationAttributeTypes(), invitationObj)
	resp.Diagnostics.Append(invitationDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := invitationDeclineStateData{
		Id:         stateData.Id,
		Invitation: invitationAttr,
	}

	invitationId := stateData.Id.ValueString()
	err = waitForInvitationDeclined(ctx, r.client, invitationId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for invitation decline",
			"Error waiting for invitation to be declined: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationDeclineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state invitationDeclineStateData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitationId := state.Id.ValueString()
	result, err := r.client.GetAccountInvitations(ctx)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, inv := range result.GetAccountInvitations() {
		if inv.Id == invitationId {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationDeclineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *invitationDeclineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty.
}

type invitationDeclineStateData struct {
	Id         types.String `tfsdk:"id"`
	Invitation types.Object `tfsdk:"invitation"`
}

func (o invitationDeclineStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"invitation": basetypes.ObjectType{},
	}
}

func invitationDeclineInvitationAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                basetypes.StringType{},
		"state":             basetypes.StringType{},
		"target_account_id": basetypes.StringType{},
		"organization_id":   basetypes.StringType{},
		"master_account_id": basetypes.StringType{},
		"requested_time":    basetypes.StringType{},
		"expired_time":      basetypes.StringType{},
		"created_at":        basetypes.StringType{},
		"created_by":        basetypes.StringType{},
		"modified_at":       basetypes.StringType{},
		"modified_by":       basetypes.StringType{},
	}
}

func (r *invitationDeclineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	invitationId := req.ID
	if invitationId == "" {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: invitation_id",
		)
		return
	}

	invitationObj := map[string]attr.Value{
		"id":                types.StringValue(""),
		"state":             types.StringValue(""),
		"target_account_id": types.StringValue(""),
		"organization_id":   types.StringValue(""),
		"master_account_id": types.StringValue(""),
		"requested_time":    types.StringValue(""),
		"expired_time":      types.StringValue(""),
		"created_at":        types.StringValue(""),
		"created_by":        types.StringValue(""),
		"modified_at":       types.StringValue(""),
		"modified_by":       types.StringValue(""),
	}

	invitationAttr, invitationDiags := types.ObjectValue(invitationDeclineInvitationAttributeTypes(), invitationObj)
	resp.Diagnostics.Append(invitationDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := invitationDeclineStateData{
		Id:         types.StringValue(invitationId),
		Invitation: invitationAttr,
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func waitForInvitationDeclined(ctx context.Context, orgClient *organization.Client, invitationId string) error {
	return client.WaitForStatus(ctx, nil, []string{"INVITING"}, []string{"REFUSED"}, func() (interface{}, string, error) {
		result, err := orgClient.GetAccountInvitations(ctx)
		if err != nil {
			return nil, "", err
		}
		for _, inv := range result.GetAccountInvitations() {
			if inv.Id == invitationId {
				return inv, string(inv.State), nil
			}
		}
		return "REFUSED", "REFUSED", nil
	}, -1, -1, -1, -1)
}
