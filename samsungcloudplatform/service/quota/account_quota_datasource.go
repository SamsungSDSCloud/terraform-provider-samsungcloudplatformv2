package quota

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
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
						Description:         "Unique identifier for the account",
						MarkdownDescription: "Unique identifier for the account",
					},
					"account_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the account",
						MarkdownDescription: "Name of the account",
					},
					"adjustable": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag indicating if additional quota is being requested",
						MarkdownDescription: "Flag indicating if additional quota is being requested",
					},
					"applied_value": schema.Int64Attribute{
						Computed: true,
					},
					"approval": schema.BoolAttribute{
						Computed:            true,
						Description:         "Approval",
						MarkdownDescription: "Approval",
					},
					"class_value": schema.StringAttribute{
						Computed:            true,
						Description:         "Value associated with the request class",
						MarkdownDescription: "Value associated with the request class",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At",
						MarkdownDescription: "Created At",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Detailed description of the quota item",
						MarkdownDescription: "Detailed description of the quota item",
					},
					"free_rate": schema.Int64Attribute{
						Computed:            true,
						Description:         "Free Rate",
						MarkdownDescription: "Free Rate",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account Quota ID",
						MarkdownDescription: "Account Quota ID",
					},
					"initial_value": schema.Int64Attribute{
						Computed:            true,
						Description:         "Initial quota value allocated",
						MarkdownDescription: "Initial quota value allocated",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At",
						MarkdownDescription: "Modified At",
					},
					"quota_item": schema.StringAttribute{
						Computed:            true,
						Description:         "Specific quota item within the resource",
						MarkdownDescription: "Specific quota item within the resource",
					},
					"reduction": schema.BoolAttribute{
						Computed:            true,
						Description:         "Reduction",
						MarkdownDescription: "Reduction",
					},
					"request": schema.BoolAttribute{
						Computed:            true,
						Description:         "Reqeust ",
						MarkdownDescription: "Reqeust ",
					},
					"request_class": schema.StringAttribute{
						Computed:            true,
						Description:         "Classification of the quota request (e.g., Account, Region)",
						MarkdownDescription: "Classification of the quota request (e.g., Account, Region)",
					},
					"resource_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Type of the resource (e.g., Virtual Server, Storage)",
						MarkdownDescription: "Type of the resource (e.g., Virtual Server, Storage)",
					},
					"service": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the service to which quota applies",
						MarkdownDescription: "Name of the service to which quota applies",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Resource Name for the quota item",
						MarkdownDescription: "Service Resource Name for the quota item",
					},
					"unit": schema.StringAttribute{
						Computed:            true,
						Description:         "Unit in which the quota value is measured (e.g., EA, GB)",
						MarkdownDescription: "Unit in which the quota value is measured (e.g., EA, GB)",
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
