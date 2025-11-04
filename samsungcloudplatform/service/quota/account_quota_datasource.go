package quota

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &quotaAccountQuotaDataSource{}
	_ datasource.DataSourceWithConfigure = &quotaAccountQuotaDataSource{}
)

// NewQuotaAccountQuotaDataSource is a helper function to simplify the provider implementation.
func NewQuotaAccountQuotaDataSource() datasource.DataSource {
	return &quotaAccountQuotaDataSource{}
}

// quotaAccountQuotaDataSource is the data source implementation.
type quotaAccountQuotaDataSource struct {
	config  *scpsdk.Configuration
	client  *quota.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *quotaAccountQuotaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quota_account_quota"
}

// Schema defines the schema for the data source.
func (d *quotaAccountQuotaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = AccountQuotaDataSourceSchema()
}

// Configure adds the provider configured client to the data source.
func (d *quotaAccountQuotaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.Quota
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *quotaAccountQuotaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state quota.AccountQuotaDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetAccountQuota(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account Quota",
			err.Error(),
		)
		return
	}

	accountQuotaElement := data.AccountQuota

	accountQuotaModel := quota.AccountQuota{
		AccountId:    types.StringValue(accountQuotaElement.AccountId),
		AccountName:  types.StringValue(accountQuotaElement.AccountName),
		Adjustable:   types.BoolValue(accountQuotaElement.Adjustable),
		AppliedValue: common.ToNullableInt32Value(accountQuotaElement.AppliedValue.Get()),
		Approval:     types.BoolValue(accountQuotaElement.Approval),
		ClassValue:   types.StringValue(accountQuotaElement.ClassValue),
		CreatedAt:    types.StringValue(accountQuotaElement.CreatedAt.Format(time.RFC3339)),
		Description:  types.StringValue(accountQuotaElement.Description),
		FreeRate:     types.Int32Value(accountQuotaElement.FreeRate),
		Id:           types.StringValue(accountQuotaElement.Id),
		InitialValue: types.Int32Value(accountQuotaElement.InitialValue),
		ModifiedAt:   types.StringValue(accountQuotaElement.ModifiedAt.Format(time.RFC3339)),
		QuotaItem:    types.StringValue(accountQuotaElement.QuotaItem),
		Reduction:    types.BoolPointerValue(accountQuotaElement.Reduction),
		Request:      types.BoolPointerValue(accountQuotaElement.Request),
		RequestClass: types.StringValue(accountQuotaElement.RequestClass),
		ResourceType: types.StringValue(accountQuotaElement.ResourceType),
		Service:      types.StringValue(accountQuotaElement.Service),
		Srn:          types.StringValue(accountQuotaElement.Srn),
		Unit:         types.StringValue(accountQuotaElement.Unit),
	}

	AccountQuotaObjectValue, _ := types.ObjectValueFrom(ctx, accountQuotaModel.AttributeTypes(), accountQuotaModel)
	state.AccountQuota = AccountQuotaObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func AccountQuotaDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "The account quota",
		Attributes: map[string]schema.Attribute{
			"account_quota": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique identifier for the account\n  - example: 2b7ed60576ce404bbc734266ff1839a5",
						MarkdownDescription: "Unique identifier for the account\n  - example: 2b7ed60576ce404bbc734266ff1839a5",
					},
					"account_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the account\n  - example: Example Account Inc.",
						MarkdownDescription: "Name of the account\n  - example: Example Account Inc.",
					},
					"adjustable": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag indicating if additional quota is being requested\n  - example: true",
						MarkdownDescription: "Flag indicating if additional quota is being requested\n  - example: true",
					},
					"applied_value": schema.Int64Attribute{
						Computed: true,
					},
					"approval": schema.BoolAttribute{
						Computed:            true,
						Description:         "Approval\n  - example: false",
						MarkdownDescription: "Approval\n  - example: false",
					},
					"class_value": schema.StringAttribute{
						Computed:            true,
						Description:         "Value associated with the request class\n  - example: global",
						MarkdownDescription: "Value associated with the request class\n  - example: global",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Detailed description of the quota item\n  - example: Maximum disk size for virtual servers in the account",
						MarkdownDescription: "Detailed description of the quota item\n  - example: Maximum disk size for virtual servers in the account",
					},
					"free_rate": schema.Int64Attribute{
						Computed:            true,
						Description:         "Free Rate\n  - example: 10",
						MarkdownDescription: "Free Rate\n  - example: 10",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account Quota ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						MarkdownDescription: "Account Quota ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
					},
					"initial_value": schema.Int64Attribute{
						Computed:            true,
						Description:         "Initial quota value allocated\n  - example: 100",
						MarkdownDescription: "Initial quota value allocated\n  - example: 100",
					},
					"max_per_account": schema.Int64Attribute{
						Computed:            true,
						Description:         "Max per Account Value\n  - maximum: 9.99999999e+08\n  - minimum: 1\n  - example: 1000",
						MarkdownDescription: "Max per Account Value\n  - maximum: 9.99999999e+08\n  - minimum: 1\n  - example: 1000",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
					},
					"quota_item": schema.StringAttribute{
						Computed:            true,
						Description:         "Specific quota item within the resource\n  - example: QUOTA.REQUEST.COUNT",
						MarkdownDescription: "Specific quota item within the resource\n  - example: QUOTA.REQUEST.COUNT",
					},
					"reduction": schema.BoolAttribute{
						Computed:            true,
						Description:         "Reduction\n  - example: false",
						MarkdownDescription: "Reduction\n  - example: false",
					},
					"request": schema.BoolAttribute{
						Computed:            true,
						Description:         "Request \n  - example: false",
						MarkdownDescription: "Request \n  - example: false",
					},
					"request_class": schema.StringAttribute{
						Computed:            true,
						Description:         "Request Class\n  - example: Account",
						MarkdownDescription: "Request Class\n  - example: Account",
					},
					"resource_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Type of the resource (e.g., Virtual Server, Storage)\n  - example: Virtual Server Disk",
						MarkdownDescription: "Type of the resource (e.g., Virtual Server, Storage)\n  - example: Virtual Server Disk",
					},
					"service": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the service to which quota applies\n  - example: Virtual Server",
						MarkdownDescription: "Name of the service to which quota applies\n  - example: Virtual Server",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Resource Name for the quota item\n  - example: srn:s::kr-west1:quota:account-quota/123456789",
						MarkdownDescription: "Service Resource Name for the quota item\n  - example: srn:s::kr-west1:quota:account-quota/123456789",
					},
					"unit": schema.StringAttribute{
						Computed:            true,
						Description:         "Unit in which the quota value is measured\n  - example: GB",
						MarkdownDescription: "Unit in which the quota value is measured\n  - example: GB",
					},
				},
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Account Quota ID",
				MarkdownDescription: "Account Quota ID",
			},
		},
	}
}
