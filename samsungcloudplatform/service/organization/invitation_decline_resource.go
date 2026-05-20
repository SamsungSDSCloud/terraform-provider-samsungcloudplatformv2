package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &invitationDeclineResource{}
	_ resource.ResourceWithConfigure = &invitationDeclineResource{}
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
	resp.TypeName = req.ProviderTypeName + "_invitation_decline"
}

func (r *invitationDeclineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Decline - Use a separate provider alias with target account credentials",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Invitation ID to decline",
				Required:    true,
			},
			"invitation": schema.SingleNestedAttribute{
				Description: "Invitation details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Invitation ID",
						Computed:    true,
					},
					"state": schema.StringAttribute{
						Description: "Invitation State (e.g., DECLINED)",
						Computed:    true,
					},
					"target_account_id": schema.StringAttribute{
						Description: "Target Account ID",
						Computed:    true,
					},
					"organization_id": schema.StringAttribute{
						Description: "Organization ID",
						Computed:    true,
					},
					"master_account_id": schema.StringAttribute{
						Description: "Master Account ID",
						Computed:    true,
					},
					"requested_time": schema.StringAttribute{
						Description: "Requested Time",
						Computed:    true,
					},
					"expired_time": schema.StringAttribute{
						Description: "Expired Time",
						Computed:    true,
					},
					"created_at": schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					"created_by": schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
					"modified_at": schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					"modified_by": schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
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

	state := invitationDeclineStateData{
		Id:         stateData.Id,
		Invitation: types.ObjectValueMust(invitationDeclineInvitationAttributeTypes(), invitationObj),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationDeclineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *invitationDeclineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *invitationDeclineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
