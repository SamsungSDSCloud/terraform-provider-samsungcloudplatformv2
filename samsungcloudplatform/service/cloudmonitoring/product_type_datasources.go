package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListService
// [GET] /v1/cloudmonitorings/product/v1/product-types
type cloudMonitoringProductTypeDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringProductTypeDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringProductTypeDataSources{}
)

func NewCloudMonitoringProductTypeSources() datasource.DataSource {
	return &cloudMonitoringProductTypeDataSources{}
}

func (d *cloudMonitoringProductTypeDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_producttypes"
}

func (d *cloudMonitoringProductTypeDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of cloudMonitoringProductTypeDataSources.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ProductCategoryCode"): schema.StringAttribute{
				Description: "ProductCategoryCode",
				Optional:    true,
			},
			common.ToSnakeCase("ProductTypes"): schema.ListNestedAttribute{
				Description: "ProductTypes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ParentProductTypeName"): schema.StringAttribute{
							Description: "ParentProductTypeName",
							Required:    true,
						},
						common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
							Description: "ProductTypeCode",
							Required:    true,
						},
						common.ToSnakeCase("ProductTypeName"): schema.StringAttribute{
							Description: "ProductTypeName",
							Required:    true,
						},
						common.ToSnakeCase("StateMetricKey"): schema.StringAttribute{
							Description: "StateMetricKey",
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

func (d *cloudMonitoringProductTypeDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudMonitoringProductTypeDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.ProductTypeDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 클라이언트를 통해 실행한 결과 저장
	res, err := d.client.GetProductTypeList(state.ProductCategoryCode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.ProductType{
			ParentProductTypeName: types.StringValue(element.ParentProductTypeName),
			ProductTypeCode:       types.StringValue(element.ProductTypeCode),
			ProductTypeName:       types.StringValue(element.ProductTypeName),
			StateMetricKey:        types.StringValue(element.StateMetricKey),
		}

		state.ProductTypes = append(state.ProductTypes, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
