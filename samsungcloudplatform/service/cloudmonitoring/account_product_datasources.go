package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListAccountResources
// [GET] /v1/cloudmonitorings/product/v2/accounts/products
type cloudMonitoringAccountProductDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringAccountProductDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringAccountProductDataSources{}
)

func NewCloudMonitoringAccountProductSources() datasource.DataSource {
	return &cloudMonitoringAccountProductDataSources{}
}

func (d *cloudMonitoringAccountProductDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_accountproducts"
}

func (d *cloudMonitoringAccountProductDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of cloudMonitoringAccountProductDataSources.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Required:    true,
			},
			common.ToSnakeCase("AccountProducts"): schema.ListNestedAttribute{
				Description: "AccountProducts",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "accountId",
							Required:    true,
						},
						common.ToSnakeCase("EndDt"): schema.StringAttribute{
							Description: "endDt",
							Optional:    true,
						},
						common.ToSnakeCase("LastEventLevel"): schema.StringAttribute{
							Description: "lastEventLevel",
							Optional:    true,
						},
						common.ToSnakeCase("PoolName"): schema.StringAttribute{
							Description: "poolName",
							Optional:    true,
						},
						common.ToSnakeCase("ProductIpAddress"): schema.StringAttribute{
							Description: "productIpAddress",
							Optional:    true,
						},
						common.ToSnakeCase("ProductName"): schema.StringAttribute{
							Description: "productName",
							Required:    true,
						},
						common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
							Description: "productResourceId",
							Required:    true,
						},
						common.ToSnakeCase("ProductState"): schema.StringAttribute{
							Description: "productState",
							Required:    true,
						},
						common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
							Description: "productTypeCode",
							Required:    true,
						},
						common.ToSnakeCase("ProductTypeName"): schema.StringAttribute{
							Description: "productTypeName",
							Required:    true,
						},
						common.ToSnakeCase("StartDt"): schema.StringAttribute{
							Description: "startDt",
							Required:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *cloudMonitoringAccountProductDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudMonitoring
	d.clients = inst.Client
}

func (d *cloudMonitoringAccountProductDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.AccountProductDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetAccountProductList(state.ResourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.AccountProduct{
			AccountId:         types.StringValue(element.AccountId),
			EndDt:             types.StringValue(element.GetEndDt().String()),
			LastEventLevel:    types.StringValue(element.GetLastEventLevel()),
			PoolName:          types.StringValue(element.GetPoolName()),
			ProductIpAddress:  types.StringValue(element.GetProductIpAddress()),
			ProductName:       types.StringValue(element.ProductName),
			ProductResourceId: types.StringValue(element.ProductResourceId),
			ProductState:      types.StringValue(element.ProductState),
			ProductTypeCode:   types.StringValue(element.ProductTypeCode),
			ProductTypeName:   types.StringValue(element.ProductTypeName),
			StartDt:           types.StringValue(element.GetStartDt().String()),
		}

		state.AccountProducts = append(state.AccountProducts, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
