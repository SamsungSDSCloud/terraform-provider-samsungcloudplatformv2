package servicewatch

import (
	"context"
	"math"
	"strings"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *servicewatch.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: servicewatch.NewAPIClient(config),
	}
}

func (client *Client) GetDashboardList(request DashboardDataSources) (*servicewatch.DashboardPageResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.ServicewatchV1DashboardsAPIAPI.ListDashboards(ctx)

	req = req.Size(math.MaxInt32) // 조건에 맞는 모든 리스트 조회를 위해 size 는 Int32 최대값으로 지정한다.
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.NameLike.IsNull() {
		req = req.NameLike(request.NameLike.ValueString())
	}
	if !request.FavoriteEnabled.IsNull() {
		req = req.FavoriteEnabled(request.FavoriteEnabled.ValueBool())
	}
	if !request.Type.IsNull() {
		req = req.Type_(request.Type.ValueString())
	}
	if !request.ServiceCode.IsNull() {
		req = req.ServiceCode(request.ServiceCode.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetDashboard(ctx context.Context, dashboardId string) (*servicewatch.DashboardDetailResponseV1Dot1, error) {
	req := client.sdkClient.ServicewatchV1DashboardsAPIAPI.ShowDashboard(ctx, dashboardId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateDashboard(ctx context.Context, request DashboardResource) (*servicewatch.DashboardDetailResponseV1Dot1, error) {
	req := client.sdkClient.ServicewatchV1DashboardsAPIAPI.CreateDashboard(ctx)

	req = req.DashboardCreationRequest(servicewatch.DashboardCreationRequest{
		Name:    request.Name.ValueString(),
		Widgets: convertWidgets(request.Widgets),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateDashboard(ctx context.Context, dashboardId string, request DashboardResource) (*servicewatch.DashboardDetailResponseV1Dot1, error) {
	req := client.sdkClient.ServicewatchV1DashboardsAPIAPI.SetDashboard(ctx, dashboardId)

	req = req.DashboardUpdateRequest(servicewatch.DashboardUpdateRequest{
		Name:    *servicewatch.NewNullableString(request.Name.ValueStringPointer()),
		Widgets: convertWidgets(request.Widgets),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteDashboard(ctx context.Context, dashboardIds []string) (*servicewatch.DashboardBulkDeleteResponse, error) {
	req := client.sdkClient.ServicewatchV1DashboardsAPIAPI.DeleteBulkDashboards(ctx)

	req = req.DashboardBulkDeleteRequest(servicewatch.DashboardBulkDeleteRequest{
		DashboardIds: dashboardIds,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLogGroupList(ctx context.Context, request LogGroupDataSources) (*servicewatch.LogGroupPageResponse, error) {

	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.ListLogGroups(ctx)
	req = req.Size(math.MaxInt32) // 조건에 맞는 모든 리스트 조회를 위해 size 는 Int32 최대값으로 지정한다.
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if len(request.Ids) > 0 {
		req = req.Ids(toStringSlice(request.Ids))
	}
	if len(request.RetentionPeriods) > 0 {
		req = req.RetentionPeriods(toInt32Slice(request.RetentionPeriods))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLogGroup(ctx context.Context, logGroupId string) (*servicewatch.LogGroupShowResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.ShowLogGroup(ctx, logGroupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLogGroup(ctx context.Context, request LogGroupResource) (*servicewatch.LogGroupShowResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.CreateLogGroup(ctx)

	req = req.LogGroupCreateRequest(servicewatch.LogGroupCreateRequest{
		Name:            request.Name.ValueString(),
		RetentionPeriod: request.RetentionPeriod.ValueInt32(),
		Tags:            convertResourceTag(request.Tags),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLogGroup(ctx context.Context, dashboardId string, request LogGroupResource) (*servicewatch.LogGroupShowResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.SetLogGroup(ctx, dashboardId)

	req = req.LogGroupSetRequest(servicewatch.LogGroupSetRequest{
		RetentionPeriod: request.RetentionPeriod.ValueInt32(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLogGroup(ctx context.Context, logGroupIds []string) (*servicewatch.LogGroupsDeleteResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.DeleteLogGroups(ctx)

	req = req.LogGroupsDeleteRequest(servicewatch.LogGroupsDeleteRequest{
		Ids: logGroupIds,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLogStream(ctx context.Context, logGroupId string, logStreamId string) (*servicewatch.LogStreamShowResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.ShowLogGroupLogStream(ctx, logGroupId, logStreamId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLogStream(ctx context.Context, logGroupId string, logStreamName string) (*servicewatch.LogStreamShowResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.CreateLogGroupLogStream(ctx, logGroupId)
	req = req.LogGroupLogStreamCreateRequest(servicewatch.LogGroupLogStreamCreateRequest{
		Name: logStreamName,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLogStream(ctx context.Context, logGroupId string, logStreamIds []string) (*servicewatch.LogGroupsDeleteResponse, error) {
	req := client.sdkClient.ServicewatchV1LogGroupsAPIsAPI.DeleteLogGroupLogStreams(ctx, logGroupId)
	req = req.LogGroupLogStreamsDeleteRequest(servicewatch.LogGroupLogStreamsDeleteRequest{
		Ids: logStreamIds,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAlert(ctx context.Context, alertId string) (*servicewatch.AlertDetailResponse, error) {
	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.ShowAlert(ctx, alertId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateAlert(ctx context.Context, request AlertResource) (*servicewatch.AlertCreateResponse, error) {
	missingData := servicewatch.MissingDataOptionEnum(request.MissingDataOption.ValueString())

	var recipientIds []string
	request.RecipientIds.ElementsAs(ctx, &recipientIds, false)

	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.CreateAlert(ctx)

	alertCreateRequest := servicewatch.AlertCreateRequest{
		Name:              request.Name.ValueString(),
		Description:       nullableString(request.Description),
		Level:             servicewatch.AlertLevelEnum(request.Level.ValueString()),
		NamespaceId:       request.NamespaceId.ValueString(),
		MetricId:          request.MetricId.ValueString(),
		Dimensions:        convertAlertDimension(ctx, request.Dimensions),
		Period:            request.Period.ValueInt32(),
		Statistic:         servicewatch.StatisticEnum(request.Statistic.ValueString()),
		EvaluationCount:   request.EvaluationCount.ValueInt32Pointer(),
		Threshold:         nullableFloat32(request.Threshold),
		UpperBound:        nullableFloat32(request.UpperBound),
		LowerBound:        nullableFloat32(request.LowerBound),
		Operator:          servicewatch.OperatorEnum(request.Operator.ValueString()),
		ViolationCount:    request.ViolationCount.ValueInt32Pointer(),
		MissingDataOption: *servicewatch.NewNullableMissingDataOptionEnum(&missingData),
		RecipientIds:      recipientIds,
		Tags:              convertAlertTag(request.Tags),
	}
	req = req.AlertCreateRequest(alertCreateRequest)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAlertActivated(ctx context.Context, alertId string, activatedYn string) (*servicewatch.AlertActivatedSetResponse, error) {
	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.SetAlertActivated(ctx, alertId)
	req = req.AlertActivatedSetRequest(servicewatch.AlertActivatedSetRequest{
		ActivatedYn: servicewatch.YNEnum(activatedYn),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAlertDescription(ctx context.Context, alertId string, description string) (*servicewatch.AlertDescriptionSetResponse, error) {
	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.SetAlertDescription(ctx, alertId)
	req = req.AlertDescriptionSetRequest(servicewatch.AlertDescriptionSetRequest{
		Description: &description,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAlert(ctx context.Context, alertId string, request AlertResource) (*servicewatch.AlertSetResponse, error) {
	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.SetAlert(ctx, alertId)

	missingData := servicewatch.MissingDataOptionEnum(request.MissingDataOption.ValueString())
	level := servicewatch.AlertLevelEnum(request.Level.ValueString())

	req = req.AlertSetRequest(servicewatch.AlertSetRequest{
		Level:             &level,
		NamespaceId:       request.NamespaceId.ValueString(),
		MetricId:          request.MetricId.ValueString(),
		Dimensions:        convertAlertDimension(ctx, request.Dimensions),
		Period:            request.Period.ValueInt32(),
		Statistic:         request.Statistic.ValueString(),
		EvaluationCount:   request.EvaluationCount.ValueInt32Pointer(),
		Threshold:         nullableFloat32(request.Threshold),
		UpperBound:        nullableFloat32(request.UpperBound),
		LowerBound:        nullableFloat32(request.LowerBound),
		Operator:          servicewatch.OperatorEnum(request.Operator.ValueString()),
		ViolationCount:    request.ViolationCount.ValueInt32Pointer(),
		MissingDataOption: *servicewatch.NewNullableMissingDataOptionEnum(&missingData),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAlert(ctx context.Context, alertIds []string) (*servicewatch.AlertDeleteResponse, error) {
	req := client.sdkClient.ServicewatchV1AlertsAPIsAPI.DeleteBulkAlerts(ctx)
	req = req.AlertDeleteRequest(servicewatch.AlertDeleteRequest{
		Ids: alertIds,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetMetrics(ctx context.Context, namespaceName string, metricName string, dimensionKeys [][]string) (*servicewatch.MetricsResponseV1Dot1, error) {
	req := client.sdkClient.ServicewatchV1MetricsAPIsAPI.ListMetricInfos(ctx)
	nullableMetricName := servicewatch.NewNullableString(&metricName)

	req = req.MetricSearchRequestV1Dot1(servicewatch.MetricSearchRequestV1Dot1{
		MetricName: *nullableMetricName,
		Namespaces: []servicewatch.SearchOptionNamespaceDtoV1Dot1{
			{
				Name:       namespaceName,
				Dimensions: dimensionKeys,
			},
		},
	})
	resp, _, err := req.Execute()
	return resp, err
}
func (client *Client) GetEventRule(ctx context.Context, eventRuleId string) (*servicewatch.EventRuleShowResponse, error) {
	req := client.sdkClient.ServicewatchV1EventRulesAPIsAPI.ShowEventRule(ctx, eventRuleId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateEventRule(ctx context.Context, request EventRuleResource) (*servicewatch.EventRuleShowResponse, error) {
	req := client.sdkClient.ServicewatchV1EventRulesAPIsAPI.CreateEventRule(ctx)

	req = req.EventRuleCreateRequest(servicewatch.EventRuleCreateRequest{
		Description:    *servicewatch.NewNullableString(request.Description.ValueStringPointer()),
		EventIds:       toStringSlice(request.EventIds),
		Name:           request.Name.ValueString(),
		RecipientIds:   toStringSlice(request.RecipientIds),
		ResourceTypeId: *servicewatch.NewNullableString(request.ResourceTypeId.ValueStringPointer()),
		ServiceId:      request.ServiceId.ValueString(),
		SrnList:        toStringSlice(request.SrnList),
		Tags:           convertResourceTag(request.Tags),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateEventRule(ctx context.Context, eventRuleId string, request EventRuleResource) (*servicewatch.EventRuleShowResponse, error) {
	req := client.sdkClient.ServicewatchV1EventRulesAPIsAPI.SetEventRule(ctx, eventRuleId)

	req = req.EventRuleSetRequest(servicewatch.EventRuleSetRequest{
		ActiveYn:       *servicewatch.NewNullableYNEnum(toYNEnum(request.ActiveYn.ValueString())),
		Description:    *servicewatch.NewNullableString(request.Description.ValueStringPointer()),
		EventIds:       toStringSlice(request.EventIds),
		NoneAttributes: toStringSlice(request.NoneAttributes),
		RecipientIds:   toStringSlice(request.RecipientIds),
		ResourceTypeId: *servicewatch.NewNullableString(request.ResourceTypeId.ValueStringPointer()),
		ServiceId:      request.ServiceId.ValueString(),
		SrnList:        toStringSlice(request.SrnList),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteEventRule(ctx context.Context, eventRuleIds []string) (*servicewatch.EventRulesDeleteResponse, error) {
	req := client.sdkClient.ServicewatchV1EventRulesAPIsAPI.DeleteEventRules(ctx)

	req = req.EventRulesDeleteRequest(servicewatch.EventRulesDeleteRequest{
		Ids: eventRuleIds,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func convertWidgets(widgets types.List) []servicewatch.WidgetDTO {
	if widgets.IsNull() || widgets.IsUnknown() {
		return []servicewatch.WidgetDTO{}
	}

	var items []Widget
	widgets.ElementsAs(context.Background(), &items, false)

	widgetList := make([]servicewatch.WidgetDTO, len(items))
	for i, v := range items {
		properties := convertProperties(v.Properties)
		widgetList[i] = servicewatch.WidgetDTO{
			Type:       v.Type.ValueString(),
			Width:      v.Width.ValueInt32(),
			Height:     v.Height.ValueInt32(),
			Order:      v.Order.ValueInt32(),
			Properties: properties,
		}
	}
	return widgetList
}

func convertProperties(property types.Object) servicewatch.PropertiesDTO {
	if property.IsNull() || property.IsUnknown() {
		return servicewatch.PropertiesDTO{}
	}

	var p Properties
	property.As(context.Background(), &p, basetypes.ObjectAsOptions{})

	return servicewatch.PropertiesDTO{
		Title:         p.Title.ValueString(),
		Stacked:       p.Stacked.ValueBool(),
		View:          p.View.ValueString(),
		Metrics:       convertMetrics(p.Metrics),
		Period:        nullableInt32(p.Period),
		StatisticType: nullableString(p.StatisticType),
	}
}

func convertMetrics(metrics types.List) []servicewatch.MetricDTO {
	if metrics.IsNull() || metrics.IsUnknown() {
		return []servicewatch.MetricDTO{}
	}

	var items []Metric
	metrics.ElementsAs(context.Background(), &items, false)

	metricList := make([]servicewatch.MetricDTO, len(items))
	for i, v := range items {
		metricList[i] = servicewatch.MetricDTO{
			Name:          v.Name.ValueString(),
			NamespaceName: v.NamespaceName.ValueString(),
			DisplayName:   v.DisplayName.ValueString(),
			Color:         v.Color.ValueString(),
			Dimensions:    convertDimensions(v.Dimensions),
			Period:        v.Period.ValueInt32(),
			StatisticType: v.StatisticType.ValueString(),
		}
	}
	return metricList
}

func convertDimensions(dimensions types.List) []servicewatch.DimensionDTO {
	if dimensions.IsNull() || dimensions.IsUnknown() {
		return []servicewatch.DimensionDTO{}
	}
	var items []Dimension
	dimensions.ElementsAs(context.Background(), &items, false)

	dimensionList := make([]servicewatch.DimensionDTO, 0, len(items))
	for _, dimension := range items {
		dimension := servicewatch.DimensionDTO{
			Key:   dimension.Key.ValueString(),
			Value: dimension.Value.ValueString(),
		}
		dimensionList = append(dimensionList, dimension)
	}
	return dimensionList
}

func toStringSlice(items []types.String) []string {
	result := make([]string, len(items))
	for i, v := range items {
		result[i] = string(v.ValueString())
	}
	return result
}

func toInt32Slice(items []types.Int32) []int32 {
	result := make([]int32, len(items))
	for i, v := range items {
		result[i] = v.ValueInt32()
	}
	return result
}

func convertResourceTag(tags types.Map) []servicewatch.ResourceTag {
	var TagsObject []servicewatch.ResourceTag
	for k, v := range tags.Elements() {
		tagObject := servicewatch.ResourceTag{
			Key:   k,
			Value: strings.ReplaceAll(v.String(), "\"", ""),
		}
		TagsObject = append(TagsObject, tagObject)
	}
	return TagsObject
}

func toYNEnum(v string) *servicewatch.YNEnum {
	result, _ := servicewatch.NewYNEnumFromValue(v)
	return result
}
func convertAlertTag(tags types.Map) []servicewatch.TagDTO {
	var TagsObject []servicewatch.TagDTO
	for k, v := range tags.Elements() {
		tagObject := servicewatch.TagDTO{
			Key:   k,
			Value: strings.ReplaceAll(v.String(), "\"", ""),
		}
		TagsObject = append(TagsObject, tagObject)
	}
	return TagsObject
}

func nullableFloat32(v types.Float32) servicewatch.NullableFloat32 {
	if v.IsNull() || v.IsUnknown() {
		return servicewatch.NullableFloat32{}
	}
	val := v.ValueFloat32()
	return *servicewatch.NewNullableFloat32(&val)
}

func nullableInt32(v types.Int32) servicewatch.NullableInt32 {
	if v.IsNull() || v.IsUnknown() {
		return servicewatch.NullableInt32{}
	}
	val := v.ValueInt32()
	return *servicewatch.NewNullableInt32(&val)
}

func nullableString(v types.String) servicewatch.NullableString {
	if v.IsNull() || v.IsUnknown() {
		return servicewatch.NullableString{}
	}
	val := v.ValueString()
	return *servicewatch.NewNullableString(&val)
}

func convertAlertDimension(ctx context.Context, dimensionList types.List) []servicewatch.DimensionDTO {
	var items []Dimension
	dimensionList.ElementsAs(ctx, &items, false)

	dimensions := make([]servicewatch.DimensionDTO, len(items))
	for i, v := range items {
		dimensions[i] = servicewatch.DimensionDTO{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return dimensions
}
