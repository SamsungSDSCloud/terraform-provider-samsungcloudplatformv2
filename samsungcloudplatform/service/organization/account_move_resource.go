package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &accountMoveResource{}
	_ resource.ResourceWithConfigure   = &accountMoveResource{}
	_ resource.ResourceWithImportState = &accountMoveResource{}
)

func NewAccountMoveResource() resource.Resource {
	return &accountMoveResource{}
}

type accountMoveResource struct {
	config  *client.SCPClient
	client  *organization.Client
	clients *client.SCPClient
}

func (r *accountMoveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account_move"
}

func (r *accountMoveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Move Organization Accounts to different OU",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Root ID or Organization Unit ID. \n" +
					"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
				Required: true,
			},
			"target_account_ids": schema.ListAttribute{
				Description: "Account IDs to move. \n" +
					"  - example : ['b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0'] \n",
				Required:    true,
				ElementType: types.StringType,
			},
			"success_ids": schema.ListNestedAttribute{
				Description: "Successfully moved account information. \n" +
					"  - example : '[{success_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, success_name: score-account, ...}]' \n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"success_id": schema.StringAttribute{
							Description: "Organization Account ID. \n" +
								"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
							Computed: true,
						},
						"success_name": schema.StringAttribute{
							Description: "Account Name. \n" +
								"  - example : 'score-account' \n",
							Computed: true,
						},
						"target_id": schema.StringAttribute{
							Description: "Organization Unit ID. \n" +
								"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
							Computed: true,
						},
						"target_name": schema.StringAttribute{
							Description: "Organization Unit Name. \n" +
								"  - example : 'score-organization-unit' \n",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *accountMoveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.Instance, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

func (r *accountMoveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.AccountMoveResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := plan.OrganizationId.ValueString()
	parentUnitId := plan.ParentUnitId.ValueString()

	targetAccountIds := make([]string, 0, 10)
	for _, elem := range plan.TargetAccountIds.Elements() {
		if id, ok := elem.(types.String); ok {
			targetAccountIds = append(targetAccountIds, id.ValueString())
		}
	}

	moveResp, err := r.client.MoveAccounts(ctx, targetAccountIds, parentUnitId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error moving Organization Accounts",
			"Could not move Organization Accounts, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var successIdsList []types.Object
	if moveResp != nil && moveResp.SuccessIds != nil {
		for _, s := range moveResp.SuccessIds {
			successObj, diags := types.ObjectValue(
				organization.MoveAccountSuccessValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"success_id":   types.StringValue(s.SuccessId),
					"success_name": types.StringValue(s.SuccessName),
					"target_id":    types.StringValue(s.TargetId),
					"target_name":  types.StringValue(s.TargetName),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			successIdsList = append(successIdsList, successObj)
		}
	}

	var successIdsValue types.List
	if len(successIdsList) > 0 {
		successIdsValue, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: organization.MoveAccountSuccessValue{}.AttributeTypes(ctx)}, successIdsList)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		successIdsValue = types.ListNull(types.ObjectType{AttrTypes: organization.MoveAccountSuccessValue{}.AttributeTypes(ctx)})
	}

	stateData := organization.AccountMoveResource{
		OrganizationId:   types.StringValue(orgId),
		ParentUnitId:     types.StringValue(parentUnitId),
		TargetAccountIds: plan.TargetAccountIds,
		SuccessIds:       successIdsValue,
	}

	diags = resp.State.Set(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, accountId := range targetAccountIds {
		err = waitForAccountMoveStatus(ctx, r.client, accountId, []string{}, []string{parentUnitId})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error waiting for account move",
				"Error waiting for account "+accountId+" to move: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}
}

func waitForAccountMoveStatus(ctx context.Context, orgClient *organization.Client, accountId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		accountResp, err := orgClient.GetAccount(ctx, accountId)
		if err != nil {
			return nil, "", err
		}
		return accountResp, accountResp.Account.ParentUnitId, nil
	}, -1, -1, -1, -1)
}

func (r *accountMoveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.AccountMoveResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetAccountIds := make([]string, 0, len(state.TargetAccountIds.Elements()))
	for _, elem := range state.TargetAccountIds.Elements() {
		if id, ok := elem.(types.String); ok {
			targetAccountIds = append(targetAccountIds, id.ValueString())
		}
	}

	for _, accountId := range targetAccountIds {
		accountResp, err := r.client.GetAccount(ctx, accountId)
		if err != nil {
			if !strings.Contains(err.Error(), "404") {
				resp.Diagnostics.AddError(
					"Failed to read account",
					fmt.Sprintf("accountId: %s, err: %s", accountId, err.Error()),
				)
				return
			}
			resp.State.RemoveResource(ctx)
			return
		}

		if accountResp.Account.ParentUnitId != state.ParentUnitId.ValueString() {
			resp.Diagnostics.AddWarning(
				"Account drift detected",
				fmt.Sprintf("accountId %s is now in OU %s, expected %s",
					accountId,
					accountResp.Account.ParentUnitId,
					state.ParentUnitId.ValueString(),
				),
			)
			state.ParentUnitId = types.StringValue(accountResp.Account.ParentUnitId)
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *accountMoveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *accountMoveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty.
}

func (r *accountMoveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) < 2 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: organization_id:parent_unit_id:account_id1,account_id2 or use terraform import with config",
		)
		return
	}

	orgId := parts[0]
	parentUnitId := parts[1]

	var accountIds []string
	if len(parts) > 2 {
		accountIds = strings.Split(parts[2], ",")
	}

	var accountIdsList types.List
	var accountIdsDiags diag.Diagnostics
	accountIdsList, accountIdsDiags = types.ListValueFrom(ctx, types.StringType, accountIds)
	resp.Diagnostics.Append(accountIdsDiags...)

	state := organization.AccountMoveResource{
		OrganizationId:   types.StringValue(orgId),
		ParentUnitId:     types.StringValue(parentUnitId),
		TargetAccountIds: accountIdsList,
		SuccessIds:       types.ListNull(types.ObjectType{AttrTypes: organization.MoveAccountSuccessValue{}.AttributeTypes(ctx)}),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

type SCPClient struct {
}
