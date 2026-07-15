package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Configure adds the provider configured client to the data source.
func (d *organizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Organization Name. \n" +
					"  - example : 'My Organization' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the organization was created. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
			},
			"created_by": schema.StringAttribute{
				Computed: true,
				Description: "User who created the organization. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
			},
			"master_account_id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the master account that manages the organization. \n" +
					"  - example : '9f8e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
			},
			"master_account_email": schema.StringAttribute{
				Computed: true,
				Description: "Email address of the master account that manages the organization. \n" +
					"  - example : 'admin@example.com' \n",
			},
			"delegation_account_id": schema.StringAttribute{
				Computed: true,
				Description: "Delegation Account. \n" +
					"  - example : '124e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
			},
			"root_unit_id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the root organizational unit. \n" +
					"  - example : 'r-abcd7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
			},
			"srn": schema.StringAttribute{
				Computed: true,
				Description: "Samsung Resource Name (SRN) uniquely identifying this resource. \n" +
					"  - example : 'srn:dev2::1b8c29a138d78e24bf1fa86812fcaa:kr-west1::organizations/organization/o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"use_scp_yn": schema.BoolAttribute{
				Computed: true,
				Description: "Control Policy Usage YN. \n" +
					"  - example : true \n",
			},
			"creator_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the organization creator. \n" +
					"  - example : 'John Doe na' \n",
			},
			"modified_at": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the organization was modified. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
				Description: "User who modified the organization. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
			},
			"modifier_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the organization modifier. \n" +
					"  - example : 'Alice' \n",
			},
		},
	}
}

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organization.OrganizationResource

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := data.Id.ValueString()

	result, err := d.client.GetOrganization(ctx, organizationId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization",
			err.Error(),
		)
		return
	}

	org := result.Organization

	data.Id = types.StringValue(org.Id)
	data.Name = types.StringValue(org.Name)
	data.CreatedAt = types.StringValue(org.CreatedAt.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(org.CreatedBy)
	data.MasterAccountId = types.StringValue(org.MasterAccountId)
	data.MasterAccountEmail = types.StringValue(org.MasterAccountEmail)
	data.DelegationAccountId = types.StringPointerValue(org.DelegationAccountId.Get())
	data.RootUnitId = types.StringValue(org.RootUnitId)
	data.Srn = types.StringValue(org.Srn)
	data.UseScpYn = types.BoolValue(org.UseScpYn)
	data.CreatorName = types.StringValue(org.GetCreatorName())
	data.ModifiedAt = types.StringValue(org.ModifiedAt.Format(time.RFC3339))
	data.ModifiedBy = types.StringValue(org.ModifiedBy)
	data.ModifierName = types.StringValue(org.GetModifierName())

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
