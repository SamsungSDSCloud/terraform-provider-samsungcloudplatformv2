package cloudmonitoring

import "github.com/hashicorp/terraform-plugin-framework/datasource"

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListMetrics
// [GET] /v1/cloudmonitorings/product/v2/metrics
// 파일명은 model.go 파일의 MetricDataSourceIds를 참고해서 metric_datasources.go 로 생성함

type cloudMonitoringMetricDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// DataSource, DataSourceWithConfigure 인터페이스를 상속받아서 cloudMonitoringMetric에 대한 DataSources를 정의
var (
	_ datasource.DataSource              = &cloudMonitoringMetricDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringMetricDataSources{}
)

// service.go 파일에 해당 데이터소스를 서비스로 등록하는데 사용
func NewCloudMonitoringMetricDataSources() datasource.DataSource {
	return &cloudMonitoringMetricDataSources{}
}

// Metadata returns the data source type name.
func (d *cloudMonitoringMetricDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// main.tf의 "samsungcloudplatformv2_cloudmonitoring_metrics" 에서 사용
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_metrics"
}

func (d *cloudMonitoringMetricDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	// model.go에서 작성한 필드에 대해 상세한 속성 작성
	resp.Schema = schema.Schema{
		Description: "The Schema of cloudMonitoringMetricDataSources.",
		Attributes: map[string]schema.Attribute{
			// model.go의 tfsdk 포맷과 일치하도록 필드명 입력
			common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
				// Required: 필수로 값이 있어야하는 것
				// Optional: 값이 없을 수도 있는 것
				// Computed: 하위 값으로 구성되는 필드
				Description: "Attribute",
				Optional:    true,
			},
			common.ToSnakeCase("ObjectType"): schema.StringAttribute{
				Description: "ObjectType",
				Optional:    true,
			},
			// 하위 스키마가 존재하는 스키마는 아래와 같이 ListNestedAttribute로 정의
			common.ToSnakeCase("Metrics"): schema.ListNestedAttribute{
				Description: "Metrics.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("DisableObject"): schema.StringAttribute{
							Description: "DisableObject",
							Optional:    true,
						},
						common.ToSnakeCase("DisplayUnit"): schema.StringAttribute{
							Description: "DisplayUnit",
							Optional:    true,
						},
						common.ToSnakeCase("FixedUnit"): schema.StringAttribute{
							Description: "FixedUnit",
							Optional:    true,
						},
						common.ToSnakeCase("IsLogMetric"): schema.StringAttribute{
							Description: "IsLogMetric",
							Required:    true,
						},
						common.ToSnakeCase("IsObjectExist"): schema.StringAttribute{
							Description: "IsObjectExist",
							Required:    true,
						},
						common.ToSnakeCase("MetricDescription"): schema.StringAttribute{
							Description: "MetricDescription",
							Optional:    true,
						},
						common.ToSnakeCase("MetricDescriptionEn"): schema.StringAttribute{
							Description: "MetricDescriptionEn",
							Optional:    true,
						},
						common.ToSnakeCase("MetricKey"): schema.StringAttribute{
							Description: "MetricKey",
							Required:    true,
						},
						common.ToSnakeCase("MetricName"): schema.StringAttribute{
							Description: "MetricName",
							Required:    true,
						},
						common.ToSnakeCase("MetricOrder"): schema.Int32Attribute{
							Description: "MetricOrder",
							Optional:    true,
						},
						common.ToSnakeCase("MetricSetKey"): schema.StringAttribute{
							Description: "MetricSetKey",
							Required:    true,
						},
						common.ToSnakeCase("MetricSetName"): schema.StringAttribute{
							Description: "MetricSetName",
							Required:    true,
						},
						common.ToSnakeCase("MetricType"): schema.StringAttribute{
							Description: "MetricType",
							Required:    true,
						},
						common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
							Description: "MetricUnit",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectKeyName"): schema.StringAttribute{
							Description: "ObjectKeyName",
							Required:    true,
						},
						common.ToSnakeCase("ObjectType"): schema.StringAttribute{
							Description: "ObjectType",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectTypeNameEng"): schema.StringAttribute{
							Description: "ObjectTypeNameEng",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectTypeNameLoc"): schema.StringAttribute{
							Description: "ObjectTypeNameLoc",
							Optional:    true,
						},
						common.ToSnakeCase("PerfTitle"): schema.StringAttribute{
							Description: "PerfTitle",
							Required:    true,
						},
						common.ToSnakeCase("ProductTargetType"): schema.StringAttribute{
							Description: "ProductTargetType",
							Optional:    true,
						},
						common.ToSnakeCase("ProductTargetTypeEn"): schema.StringAttribute{
							Description: "ProductTargetTypeEn",
							Optional:    true,
						},
						common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
							Description: "ProductTypeCode",
							Required:    true,
						},
						common.ToSnakeCase("ProductTypeName"): schema.StringAttribute{
							Description: "ProductTypeName",
							Required:    true,
						},
					},
				},
			},
		},
		// model.go에서 filter를 정의했다면 이것도 입력한다.
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

func (d *cloudMonitoringMetricDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudMonitoring
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *cloudMonitoringMetricDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// 모델에서 정의한 인풋필드
	var state cloudmonitoring.MetricDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 클라이언트를 통해 실행한 결과 저장
	res, err := d.client.GetMetricList(state.ProductTypeCode, state.ObjectType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.Metric{
			// 필수값은 바로 변수명을 가져올 수 있고 옵션값은 get을 통해 가져올 수 있는 것으로 보임(empty에 대한 로직때문으로 추측)
			DisableObject:       types.StringValue(element.GetDisableObject()),
			DisplayUnit:         types.StringValue(element.GetDisplayUnit()),
			FixedUnit:           types.StringValue(element.GetFixedUnit()),
			IsLogMetric:         types.StringValue(element.IsLogMetric),
			IsObjectExist:       types.StringValue(element.IsObjectExist),
			MetricDescription:   types.StringValue(element.GetMetricDescription()),
			MetricDescriptionEn: types.StringValue(element.GetMetricDescriptionEn()),
			MetricKey:           types.StringValue(element.MetricKey),
			MetricName:          types.StringValue(element.MetricName),
			MetricOrder:         types.Int32Value(element.GetMetricOrder()),
			MetricSetKey:        types.StringValue(element.MetricSetKey),
			MetricSetName:       types.StringValue(element.MetricSetName),
			MetricType:          types.StringValue(element.MetricType),
			MetricUnit:          types.StringValue(element.GetFixedUnit()),
			ObjectKeyName:       types.StringValue(element.ObjectKeyName),
			ObjectType:          types.StringValue(element.GetObjectType()),
			ObjectTypeNameEng:   types.StringValue(element.GetObjectTypeNameEng()),
			ObjectTypeNameLoc:   types.StringValue(element.GetObjectTypeNameLoc()),
			PerfTitle:           types.StringValue(element.PerfTitle),
			ProductTargetType:   types.StringValue(element.GetMetricDescriptionEn()),
			ProductTargetTypeEn: types.StringValue(element.GetMetricDescriptionEn()),
			ProductTypeCode:     types.StringValue(element.ProductTypeCode),
			ProductTypeName:     types.StringValue(element.ProductTypeName),
		}

		state.Metrics = append(state.Metrics, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
