package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &delegationAccountResource{}
	_ resource.ResourceWithConfigure   = &delegationAccountResource{}
	_ resource.ResourceWithImportState = &delegationAccountResource{}
)

func NewDelegationAccountResource() resource.Resource {
	return &delegationAccountResource{}
}

type delegationAccountResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *delegationAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_delegation_account"
}

func (r *delegationAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages delegation accounts for the organization",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "Delegation Account ID. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
				Required: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
			},
			"service_type": schema.StringAttribute{
				Description: "Service type of the delegated account. \n" +
					"  - allowed_values : ['identity-center', 'resource-optimizer', 'infrastructure-builder'] \n" +
					"  - example : 'identity-center' \n",
				Required: true,
			},
			"delegation_account": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Delegation Account Info. \n" +
					"  - example : '{id: 0a36e0746dbf4908acf0357829701381, account_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, ...}' \n",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the delegation record. \n" +
							"  - example : '0a36e0746dbf4908acf0357829701381' \n",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Created timestamp. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "Creator ID. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"account_id": schema.StringAttribute{
						Description: "Delegation Account ID. \n" +
							"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						Computed: true,
					},
					"organization_id": schema.StringAttribute{
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						Computed: true,
					},
					"service_type": schema.StringAttribute{
						Description: "Service type of the delegated account. \n" +
							"  - allowed_values : ['identity-center', 'resource-optimizer', 'infrastructure-builder'] \n" +
							"  - example : 'identity-center' \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *delegationAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *delegationAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.DelegationAccountResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateDelegationAccount(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Delegation Account",
			"Could not create Delegation Account, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	account := data.DelegationAccount
	delegationValue, objDiags := types.ObjectValue(
		organization.DelegationAccountValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":              types.StringValue(account.Id),
			"account_id":      types.StringValue(account.AccountId),
			"organization_id": types.StringValue(account.OrganizationId),
			"service_type":    types.StringValue(account.ServiceType),
			"created_at":      types.StringValue(account.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
			"created_by":      types.StringValue(account.CreatedBy),
		},
	)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.DelegationAccount = delegationValue

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *delegationAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.DelegationAccountResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	accountId := state.AccountId.ValueString()
	serviceType := state.ServiceType.ValueString()

	data, err := r.client.ListDelegationAccounts(ctx, orgId, serviceType, 0, 0, "", accountId)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	if data != nil {
		for _, item := range data.DelegationAccounts {
			if item.AccountId == accountId {
				delegationValue, objDiags := types.ObjectValue(
					organization.DelegationAccountValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":              types.StringValue(item.Id),
						"account_id":      types.StringValue(item.AccountId),
						"organization_id": types.StringValue(item.OrganizationId),
						"service_type":    types.StringValue(item.ServiceType),
						"created_at":      types.StringValue(item.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
						"created_by":      types.StringValue(item.CreatedBy),
					},
				)
				resp.Diagnostics.Append(objDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				state.DelegationAccount = delegationValue
				found = true
				break
			}
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *delegationAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"Delegation Account resource does not support update operations",
	)
}

func (r *delegationAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.DelegationAccountResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteDelegationAccount(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Delegation Account",
			"Could not delete Delegation Account, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *delegationAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID

	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"Expected format: account_id:organization_id:service_type",
		)
		return
	}

	state := organization.DelegationAccountResource{
		AccountId:         types.StringValue(parts[0]),
		OrganizationId:    types.StringValue(parts[1]),
		ServiceType:       types.StringValue(parts[2]),
		DelegationAccount: types.ObjectNull(organization.DelegationAccountValue{}.AttributeTypes(ctx)),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
