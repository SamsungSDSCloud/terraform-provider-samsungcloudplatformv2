package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ datasource.DataSource              = &accountInvitationsDataSource{}
	_ datasource.DataSourceWithConfigure = &accountInvitationsDataSource{}
)

// NewAccountInvitationsDataSource is a helper function to simplify the provider implementation.
func NewAccountInvitationsDataSource() datasource.DataSource {
	return &accountInvitationsDataSource{}
}

type accountInvitationsDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *accountInvitationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account_invitations"
}

// Configure adds the provider configured client to the data source.
func (d *accountInvitationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Organization
	d.clients = inst.Client
}

func (d *accountInvitationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Account Invitations",
		Attributes: map[string]schema.Attribute{
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Received invitation count. \n" +
					"  - example : 5 \n",
			},
			"account_invitations": schema.ListNestedAttribute{
				Computed: true,
				Description: "Received Invitation List. \n" +
					"  - example : '[{id: 0a36e0746dbf4908acf0357829701381, organization_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Invitation ID. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"organization_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the organization. \n" +
								"  - example : 'My Organization' \n",
						},
						"master_account_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the master account that manages the organization. \n" +
								"  - example : '9f8e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
						},
						"master_account_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the master account that manages the organization. \n" +
								"  - example : 'CoreKim' \n",
						},
						"master_account_email": schema.StringAttribute{
							Computed: true,
							Description: "Email address of the master account that manages the organization. \n" +
								"  - example : 'admin@example.com' \n",
						},
						"master_account_login_id": schema.StringAttribute{
							Computed: true,
							Description: "Login ID of the master account that manages the organization. \n" +
								"  - example : 'log-archive@samsung.com' \n",
						},
						"target_account_id": schema.StringAttribute{
							Computed: true,
							Description: "Invitation Receiver Account ID. \n" +
								"  - example : '338c2fc8c29a449dbfa8681f8f1d78e5' \n",
						},
						"requested_time": schema.StringAttribute{
							Computed: true,
							Description: "Invitation Request Datetime. \n" +
								"  - example : '2024-04-17T12:34:56.789Z' \n",
						},
						"expired_time": schema.StringAttribute{
							Computed: true,
							Description: "Invitation Expired Datetime. \n" +
								"  - example : '2024-04-30T12:34:56.789Z' \n",
						},
						"state": schema.StringAttribute{
							Computed: true,
							Description: "Invitation State. \n" +
								"  - example : 'INVITING' \n" +
								"  - allowed_values : ['INVITING', 'REFUSED', 'INVITED', 'CANCELED', 'EXPIRED'] \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the invitation was sent. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who sent the invitation. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the invitation was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the invitation. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
					},
				},
			},
		},
	}
}

func (d *accountInvitationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	result, err := d.client.GetAccountInvitations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Account Invitations",
			err.Error(),
		)
		return
	}

	var invitationObjects []organization.AccountInvitationValue
	for _, inv := range result.GetAccountInvitations() {
		masterAccountName := inv.GetMasterAccountName()
		masterAccountEmail := inv.GetMasterAccountEmail()
		masterAccountLoginId := inv.GetMasterAccountLoginId()

		invitationObjects = append(invitationObjects, organization.AccountInvitationValue{
			Id:                   types.StringValue(inv.Id),
			OrganizationId:       types.StringValue(inv.OrganizationId),
			OrganizationName:     types.StringValue(inv.OrganizationName),
			MasterAccountId:      types.StringValue(inv.MasterAccountId),
			MasterAccountName:    types.StringValue(masterAccountName),
			MasterAccountEmail:   types.StringValue(masterAccountEmail),
			MasterAccountLoginId: types.StringValue(masterAccountLoginId),
			TargetAccountId:      types.StringValue(inv.TargetAccountId),
			RequestedTime:        types.StringValue(inv.RequestedTime.Format(time.RFC3339)),
			ExpiredTime:          types.StringValue(inv.ExpiredTime.Format(time.RFC3339)),
			State:                types.StringValue(string(inv.State)),
			CreatedAt:            types.StringValue(inv.CreatedAt.Format(time.RFC3339)),
			CreatedBy:            types.StringValue(inv.CreatedBy),
			ModifiedAt:           types.StringValue(inv.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:           types.StringValue(inv.ModifiedBy),
		})
	}

	invitationAttrTypes := organization.AccountInvitationValue{}.AttributeTypes(ctx)
	invitationsListVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: invitationAttrTypes}, invitationObjects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := AccountInvitationsDataSource{
		TotalCount:         types.Int64Value(int64(result.GetCount())),
		AccountInvitations: invitationsListVal,
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// AccountInvitationsDataSource is the response data
type AccountInvitationsDataSource struct {
	TotalCount         types.Int64 `tfsdk:"total_count"`
	AccountInvitations types.List  `tfsdk:"account_invitations"`
}

func (o AccountInvitationsDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"total_count":         basetypes.Int64Type{},
		"account_invitations": types.ListType{ElemType: types.ObjectType{AttrTypes: organization.AccountInvitationValue{}.AttributeTypes(ctx)}},
	}
}
