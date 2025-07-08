package quota

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &quotaAccountQuotaDataSources{}
	_ datasource.DataSourceWithConfigure = &quotaAccountQuotaDataSources{}
)

// NewQuotaAccountQuotaDataSources is a helper function to simplify the provider implementation.
func NewQuotaAccountQuotaDataSources() datasource.DataSource {
	return &quotaAccountQuotaDataSources{}
}

// quotaAccountQuotaDataSources is the data source implementation.
type quotaAccountQuotaDataSources struct {
	config  *scpsdk.Configuration
	client  *quota.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *quotaAccountQuotaDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quota_account_quotas"
}

// Schema defines the schema for the data source.
func (d *quotaAccountQuotaDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of account quotas",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "ID List",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *quotaAccountQuotaDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *quotaAccountQuotaDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state quota.AccountQuotaDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetAccountQuotas(d.clients, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account quota",
			err.Error(),
		)
	}

	state.Ids = ids // 상태 데이터에 ID 리스트를 추가한다.

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetAccountQuotas(clients *client.SCPClient, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.Quota.GetAccountQuotaList()
	if err != nil {
		return nil, err
	}
	contents := data.AccountQuotas
	filteredContents := data.AccountQuotas // 필터링된 컨텐츠

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, filters) // 필터링된 컨텐츠의 Index 정보 리턴

		for i, resource := range contents {
			if common.Contains(indices, i) { // Index 정보 기준으로 필터링 진행
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
