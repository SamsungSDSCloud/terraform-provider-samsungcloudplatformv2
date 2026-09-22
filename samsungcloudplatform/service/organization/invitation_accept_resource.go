package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &invitationAcceptResource{}
	_ resource.ResourceWithConfigure   = &invitationAcceptResource{}
	_ resource.ResourceWithImportState = &invitationAcceptResource{}
)

func NewInvitationAcceptResource() resource.Resource {
	return &invitationAcceptResource{}
}

type invitationAcceptResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *invitationAcceptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_invitation_accept"
}

func (r *invitationAcceptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Accept - Use a separate provider alias with target account credentials",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Invitation ID. \n" +
					"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
				Required: true,
			},
			"master_account_email": schema.StringAttribute{
				Description: "Email address of the master account that manages the organization. \n" +
					"  - example : 'master@samsung.com' \n",
				Computed: true,
			},
			"success_id": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Success Info. \n" +
					"  - example : '{success_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, success_name: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5}' \n",
				Attributes: map[string]schema.Attribute{
					"success_id": schema.StringAttribute{
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						Computed: true,
					},
					"success_name": schema.StringAttribute{
						Description: "Organization Name. \n" +
							"  - example : 'My Organization' \n",
						Computed: true,
					},
				},
			},
			"failed_id": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Failed Info. \n" +
					"  - example : '{error_code: Invitation.AlreadySentError, failed_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, ...}' \n",
				Attributes: map[string]schema.Attribute{
					"failed_id": schema.StringAttribute{
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						Computed: true,
					},
					"failed_name": schema.StringAttribute{
						Description: "Organization Name. \n" +
							"  - example : 'My Organization' \n",
						Computed: true,
					},
					"failed_caused": schema.StringAttribute{
						Description: "Failure Reason. \n" +
							"  - example : 'Account b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0 is not found' \n",
						Computed: true,
					},
					"error_code": schema.StringAttribute{
						Description: "Error code returned when the operation fails. \n" +
							"  - example : 'Invitation.AlreadySentError' \n",
						Computed: true,
					},
					"response": schema.MapAttribute{
						ElementType: types.StringType,
						Description: "Quota Response. \n" +
							"  - example : {} \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *invitationAcceptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *invitationAcceptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData invitationAcceptStateData
	diags := req.Plan.Get(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.AcceptInvitation(ctx, stateData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error accepting invitation",
			err.Error(),
		)
		return
	}

	if result.HasFailedId() {
		failedObj := result.GetFailedId()
		if failedObj.GetErrorCode() != "" {
			resp.Diagnostics.AddError(
				"Error accepting invitation",
				fmt.Sprintf("[%s] %s", failedObj.GetErrorCode(), failedObj.GetFailedCaused()),
			)
			return
		}
	}

	masterAccountEmail := result.GetMasterAccountEmail()

	var successIdVal types.Object
	if result.HasSuccessId() {
		successObj := result.GetSuccessId()
		var successDiags diag.Diagnostics
		successIdVal, successDiags = types.ObjectValue(
			map[string]attr.Type{
				"success_id":   basetypes.StringType{},
				"success_name": basetypes.StringType{},
			},
			map[string]attr.Value{
				"success_id":   types.StringValue(successObj.GetSuccessId()),
				"success_name": types.StringValue(successObj.GetSuccessName()),
			},
		)
		resp.Diagnostics.Append(successDiags...)
	} else {
		successIdVal = types.ObjectNull(map[string]attr.Type{
			"success_id":   basetypes.StringType{},
			"success_name": basetypes.StringType{},
		})
	}

	var failedIdVal types.Object
	if result.HasFailedId() {
		failedObj := result.GetFailedId()
		responseMap, responseDiags := types.MapValue(types.StringType, map[string]attr.Value{})
		resp.Diagnostics.Append(responseDiags...)
		var failedDiags diag.Diagnostics
		failedIdVal, failedDiags = types.ObjectValue(
			map[string]attr.Type{
				"failed_id":     basetypes.StringType{},
				"failed_name":   basetypes.StringType{},
				"failed_caused": basetypes.StringType{},
				"error_code":    basetypes.StringType{},
				"response":      types.MapType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"failed_id":     types.StringValue(failedObj.GetFailedId()),
				"failed_name":   types.StringValue(failedObj.GetFailedName()),
				"failed_caused": types.StringValue(failedObj.GetFailedCaused()),
				"error_code":    types.StringValue(failedObj.GetErrorCode()),
				"response":      responseMap,
			},
		)
		resp.Diagnostics.Append(failedDiags...)
	} else {
		failedIdVal = types.ObjectNull(map[string]attr.Type{
			"failed_id":     basetypes.StringType{},
			"failed_name":   basetypes.StringType{},
			"failed_caused": basetypes.StringType{},
			"error_code":    basetypes.StringType{},
			"response":      types.MapType{ElemType: types.StringType},
		})
	}

	state := invitationAcceptStateData{
		Id:                 stateData.Id,
		MasterAccountEmail: types.StringValue(masterAccountEmail),
		SuccessId:          successIdVal,
		FailedId:           failedIdVal,
	}

	invitationId := stateData.Id.ValueString()
	err = waitForInvitationAccepted(ctx, r.client, invitationId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for invitation acceptance",
			"Error waiting for invitation to be accepted: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationAcceptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state invitationAcceptStateData
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

func (r *invitationAcceptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *invitationAcceptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty.
}

type InvitationAcceptSuccessStateData struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
}

func (o InvitationAcceptSuccessStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
}

type InvitationAcceptFailedStateData struct {
	FailedId     types.String      `tfsdk:"failed_id"`
	FailedName   types.String      `tfsdk:"failed_name"`
	FailedCaused types.String      `tfsdk:"failed_caused"`
	ErrorCode    types.String      `tfsdk:"error_code"`
	Response     map[string]string `tfsdk:"response"`
}

func (o InvitationAcceptFailedStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"error_code":    basetypes.StringType{},
		"response":      types.MapType{ElemType: types.StringType},
	}
}

type invitationAcceptStateData struct {
	Id                 types.String `tfsdk:"id"`
	MasterAccountEmail types.String `tfsdk:"master_account_email"`
	SuccessId          types.Object `tfsdk:"success_id"`
	FailedId           types.Object `tfsdk:"failed_id"`
}

func (o invitationAcceptStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   basetypes.StringType{},
		"master_account_email": basetypes.StringType{},
		"success_id":           types.ObjectType{},
		"failed_id":            types.ObjectType{},
	}
}

func (r *invitationAcceptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForInvitationAccepted(ctx context.Context, orgClient *organization.Client, invitationId string) error {
	return client.WaitForStatus(ctx, nil, []string{"INVITING"}, []string{"INVITED"}, func() (interface{}, string, error) {
		result, err := orgClient.GetAccountInvitations(ctx)
		if err != nil {
			return nil, "", err
		}
		for _, inv := range result.GetAccountInvitations() {
			if inv.Id == invitationId {
				return inv, string(inv.State), nil
			}
		}
		return "INVITED", "INVITED", nil
	}, -1, -1, -1, -1)
}
