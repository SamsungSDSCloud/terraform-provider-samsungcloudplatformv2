package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
	resp.TypeName = req.ProviderTypeName + "_account_invitations"
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
				Computed:            true,
				Description:         "Total Count",
				MarkdownDescription: "Total Count",
			},
			"account_invitations": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "A list of Account Invitations",
				MarkdownDescription: "A list of Account Invitations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Invitation ID",
							MarkdownDescription: "Invitation ID",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Organization ID",
							MarkdownDescription: "Organization ID",
						},
						"organization_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Organization Name",
							MarkdownDescription: "Organization Name",
						},
						"master_account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Master Account ID",
							MarkdownDescription: "Master Account ID",
						},
						"master_account_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Master Account Name",
							MarkdownDescription: "Master Account Name",
						},
						"master_account_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Master Account Email",
							MarkdownDescription: "Master Account Email",
						},
						"master_account_login_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Master Account Login ID",
							MarkdownDescription: "Master Account Login ID",
						},
						"target_account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Target Account ID",
							MarkdownDescription: "Target Account ID",
						},
						"requested_time": schema.StringAttribute{
							Computed:            true,
							Description:         "Requested Time",
							MarkdownDescription: "Requested Time",
						},
						"expired_time": schema.StringAttribute{
							Computed:            true,
							Description:         "Expired Time",
							MarkdownDescription: "Expired Time",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "State",
							MarkdownDescription: "State",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created At",
							MarkdownDescription: "Created At",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Created By",
							MarkdownDescription: "Created By",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified At",
							MarkdownDescription: "Modified At",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified By",
							MarkdownDescription: "Modified By",
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
