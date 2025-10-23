package cloudmonitoring

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpcloudmonitoring "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/cloudmonitoring"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpcloudmonitoring.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpcloudmonitoring.NewAPIClient(config),
	}
}

// ------------ EVENT -------------------//

func (client *Client) GetEventList(xResourceType types.String, ProductResourceId types.String, EventState types.String, QueryStartDt types.String, QueryEndDt types.String) (*scpcloudmonitoring.PageResponseOpenApiEventResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventV2API.GetProductEventList(ctx)
	//req = req.Size(math.MaxInt32)
	req = req.Size(10)

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}
	if !ProductResourceId.IsNull() {
		req = req.ProductResourceId(ProductResourceId.ValueString())
	}
	if !QueryStartDt.IsNull() {
		req = req.QueryStartDt(QueryStartDt.ValueString())
	}
	if !QueryEndDt.IsNull() {
		req = req.QueryEndDt(QueryEndDt.ValueString())
	}
	if !EventState.IsNull() {
		req = req.EventState(EventState.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccountEventList(EventState types.String, QueryStartDt types.String, QueryEndDt types.String) (*scpcloudmonitoring.PageResponseOpenApiEventResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventV2API.GetAccountEventList(ctx)
	//req = req.Size(math.MaxInt32)
	req = req.Size(10)

	if !QueryStartDt.IsNull() {
		req = req.QueryStartDt(QueryStartDt.ValueString())
	}
	if !QueryEndDt.IsNull() {
		req = req.QueryEndDt(QueryEndDt.ValueString())
	}
	if !EventState.IsNull() {
		req = req.EventState(EventState.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetEvent(eventId types.String, xResourceType types.String) (*scpcloudmonitoring.EventDetailResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventV2API.GetEventDetail(ctx, eventId.ValueString())

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetEventNotificationStateList(xResourceType types.String, EventId types.String) (*scpcloudmonitoring.PageResponseEventNotificationResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventV2API.GetEventNotificationStates(ctx, EventId.ValueString())
	//req = req.Size(math.MaxInt32)
	req = req.Size(10)

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetEventPolicyHistoryList(xResourceType types.String, EventPolicyId types.Int64, QueryStartDt types.String, QueryEndDt types.String) (*scpcloudmonitoring.PageResponseEventPolicyHistoryResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.GetEventPolicyHistories(ctx, EventPolicyId.ValueInt64())
	//req = req.Size(math.MaxInt32)
	req = req.Size(10)

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}
	if !QueryStartDt.IsNull() {
		req = req.QueryStartDt(QueryStartDt.ValueString())
	}
	if !QueryEndDt.IsNull() {
		req = req.QueryEndDt(QueryEndDt.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetEventPolicyNotificationList(xResourceType types.String, EventPolicyId types.Int64) (*scpcloudmonitoring.PageResponseNotificationResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.GetEventPolicyNotification(ctx, EventPolicyId.ValueInt64())
	//req = req.Size(math.MaxInt32)
	req = req.Size(10)

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err
}

// ------------ EVENT POLICY-------------------//

func (client *Client) GetProductEventPolicyList(request EventPolicyDataSourceIds) (*scpcloudmonitoring.PageResponseEventPolicyResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.GetProductEventPolicyList(ctx)

	if !request.ResourceType.IsNull() {
		req = req.XResourceType(request.ResourceType.ValueString())
	}

	if !request.ProductResourceId.IsNull() {
		req = req.ProductResourceId(request.ProductResourceId.ValueString())
	}

	if !request.Size.IsNull() {
		req = req.Size(10)
	}

	//criteria := &scpcloudmonitoring.EventPolicySearchCriteria{
	//	MetricKey:         request.MetricKey.ValueStringPointer(),
	//	ProductResourceId: request.ProductResourceId.ValueString(),
	//	Page:              request.Page.ValueInt32Pointer(),
	//	Size:              request.Size.ValueInt32Pointer(),
	//}
	//
	//req = req.Criteria(*criteria)

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetEventPolicyDetail(xResourceType types.String, EventPolicyId types.Int64) (*scpcloudmonitoring.EventPolicyDetailResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.GetEventPolicyDetail(ctx, EventPolicyId.ValueInt64())

	if !xResourceType.IsNull() {
		req = req.XResourceType(xResourceType.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetEventPolicy(request EventPolicyResource) (*scpcloudmonitoring.EventPolicyDetailResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.GetEventPolicyDetail(ctx, request.EventPolicyId.ValueInt64())

	if !request.ResourceType.IsNull() {
		req = req.XResourceType(request.ResourceType.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) CreateEventPolicy(request EventPolicyResource) (*scpcloudmonitoring.EventPolicyDetailResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.PutEventPolicy(ctx)

	if !request.ResourceType.IsNull() {
		req = req.XResourceType(request.ResourceType.ValueString())
	}

	eventThreshold := scpcloudmonitoring.EventThreshold{
		SingleThreshold: scpcloudmonitoring.NewSingleThreshold(
			request.EventPolicy.EventThreshold.SingleThreshold.ComparisonOperator.ValueString(),
			request.EventPolicy.EventThreshold.SingleThreshold.Value.ValueFloat64()),
		ThresholdType:  request.EventPolicy.EventThreshold.ThresholdType.ValueString(),
		MetricFunction: request.EventPolicy.EventThreshold.MetricFunction.ValueStringPointer(),
	}

	eventPolicyInfo := &scpcloudmonitoring.EventPolicyInfo{
		DisableYn:          request.DisableYn.ValueString(),
		EventLevel:         request.EventLevel.ValueString(),
		EventMessagePrefix: request.EventMessagePrefix.ValueStringPointer(),
		EventOccurTimeZone: request.EventOccurTimeZone.ValueStringPointer(),
		//EventPolicyStatistics *EventPolicyStatistics `json:"eventPolicyStatistics,omitempty"`
		EventThreshold: eventThreshold,
		FtCount:        request.FtCount.ValueInt64(),
		IsLogMetric:    request.IsLogMetric.ValueString(),
		MetricKey:      request.MetricKey.ValueString(),
		//MetricName:           request.EventPolicy.MetricName.ValueStringPointer(),
		//ObjectDisplayName:    request.EventPolicy.ObjectDisplayName.ValueStringPointer(),
		//ObjectName:           request.EventPolicy.ObjectName.ValueStringPointer(),
		//ObjectType:           request.EventPolicy.ObjectType.ValueStringPointer(),
		//PodObjectDisplayName: request.EventPolicy.PodObjectDisplayName.ValueStringPointer(),
		//PodObjectName:        request.EventPolicy.PodObjectName.ValueStringPointer(),
	}

	req = req.EventPolicyCreateRequest(scpcloudmonitoring.EventPolicyCreateRequest{
		EventPolicyRequest: *eventPolicyInfo,
		ProductResourceId:  request.ProductResourceId.ValueString(),
		//NotificationRecipients: notificationRecipients,
	})

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) UpdateEventPolicy(request EventPolicyResource, eventPolicyId types.Int64) (*scpcloudmonitoring.EventPolicyDetailResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.ModifyEventPolicy(ctx, eventPolicyId.ValueInt64())

	if !request.ResourceType.IsNull() {
		req = req.XResourceType(request.ResourceType.ValueString())
	}

	eventThreshold := scpcloudmonitoring.EventThreshold{
		SingleThreshold: scpcloudmonitoring.NewSingleThreshold(
			request.EventPolicy.EventThreshold.SingleThreshold.ComparisonOperator.ValueString(),
			request.EventPolicy.EventThreshold.SingleThreshold.Value.ValueFloat64()),
		ThresholdType:  request.EventPolicy.EventThreshold.ThresholdType.ValueString(),
		MetricFunction: request.EventPolicy.EventThreshold.MetricFunction.ValueStringPointer(),
	}

	eventPolicyInfoEditable := &scpcloudmonitoring.EventPolicyInfoEditable{
		DisableYn:          request.DisableYn.ValueString(),
		EventLevel:         request.EventLevel.ValueString(),
		EventMessagePrefix: request.EventMessagePrefix.ValueStringPointer(),
		EventOccurTimeZone: request.EventOccurTimeZone.ValueStringPointer(),
		//EventPolicyStatistics *EventPolicyStatistics `json:"eventPolicyStatistics,omitempty"`
		EventThreshold: eventThreshold,
		FtCount:        request.FtCount.ValueInt64(),
		//IsLogMetric:    request.IsLogMetric.ValueString(),
		//MetricKey:      request.MetricKey.ValueString(),
		//MetricName:           request.EventPolicy.MetricName.ValueStringPointer(),
		//ObjectDisplayName:    request.EventPolicy.ObjectDisplayName.ValueStringPointer(),
		//ObjectName:           request.EventPolicy.ObjectName.ValueStringPointer(),
		//ObjectType:           request.EventPolicy.ObjectType.ValueStringPointer(),
		//PodObjectDisplayName: request.EventPolicy.PodObjectDisplayName.ValueStringPointer(),
		//PodObjectName:        request.EventPolicy.PodObjectName.ValueStringPointer(),
	}

	req = req.EventPolicyUpdateRequest(scpcloudmonitoring.EventPolicyUpdateRequest{
		EventPolicyRequest: *eventPolicyInfoEditable,
		//NotificationRecipients: notificationRecipients,
	})

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) DeleteEventPolicy(eventPolicyId types.Int64, resourceType types.String) error {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIEventPolicyV2API.DeleteEventPolicy(ctx, eventPolicyId.ValueInt64())
	if !resourceType.IsNull() {
		req = req.XResourceType(resourceType.ValueString())
	}
	_, err := req.Execute()
	return err

}

// ------------ PRODUCT -------------------//

// ListMetrics
// [GET] /v1/cloudmonitorings/product/v2/metrics
// input: swagger의 page, size, sort, token, requestId 등의 공통필드를 제외한 나머지 비즈니스 필드 입력
// output: swagger response schema명 입력(로컬의 SDK를 사용함)
func (client *Client) GetMetricList(productTypeCode types.String, objectType types.String) (*scpcloudmonitoring.PageResponseMetricInfoDto, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIMetricV2API.GetMetricList(ctx)

	// input paremeter -> req 매핑
	if !productTypeCode.IsNull() {
		req = req.ProductTypeCode(productTypeCode.ValueString())
	}
	if !objectType.IsNull() {
		req = req.ObjectType(objectType.ValueString())
	}

	// req 실행 후 response 리턴
	resp, _, err := req.Execute()
	return resp, err
}

// ListService
// [GET] /v1/cloudmonitorings/product/v1/product-types
func (client *Client) GetProductTypeList(productCategoryCode types.String) (*scpcloudmonitoring.PageResponseProductTypeInfoDto, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIProductTypeV1API.GetProductTypeList(ctx)

	if !productCategoryCode.IsNull() {
		req = req.ProductCategoryCode(productCategoryCode.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ListAccountResources
// [GET] /v1/cloudmonitorings/product/v2/accounts/products
func (client *Client) GetAccountProductList(resourceType types.String) (*scpcloudmonitoring.PageResponseAccountProductDto, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIAccountV2API.GetAccountProductList(ctx)

	if !resourceType.IsNull() {
		req = req.XResourceType(resourceType.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ListAccountMember
// [GET] /v1/cloudmonitorings/product/v1/accounts/members
func (client *Client) GetAccountMembers() (*scpcloudmonitoring.PageResponseJsonArrayProjectMemberResponseForOpenAPI, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIAccountV1API.GetAccountMembers(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

// ListAddressBooks
// [GET] /v1/cloudmonitorings/product/v2/users/addrbooks
func (client *Client) GetAddressBookList() (*scpcloudmonitoring.PageResponseAlarmAddrBookDto, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIUserV2API.GetAdressBookList(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

// ListAddressBookMembers
// [GET] /v1/cloudmonitorings/product/v2/addrbooks/{{addrbookId}}/members
func (client *Client) GetAddressBookMemberList(addrbookId int32) (*scpcloudmonitoring.PageResponseAlarmAddrBookMemberDto, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIAddressBookV2API.GetAdressBookMemberList(ctx, addrbookId)

	resp, _, err := req.Execute()
	return resp, err
}

// ListMetricPerfData
// [POST] /v1/cloudmonitorings/product/v2/metric-data
func (client *Client) GetMetricPerfDataList(request MetricPerfDataDataSourceIds) (*scpcloudmonitoring.ListResponseMetricStatisticsDataDtoOpenAPIV2, error) {
	ctx := context.Background()
	req := client.sdkClient.OpenAPIMetricDataV2API.GetMetricPerfDataList(ctx)

	if !request.XResourceType.IsNull() {
		req = req.XResourceType(request.XResourceType.ValueString())
	}

	//metricDataConditionOpenAPIV2 := make([]scpcloudmonitoring.MetricDataConditionOpenAPIV2, len(request.MetricDataConditions))
	//
	//for i, mdc := range request.MetricDataConditions {
	//	productResourceInfos := make([]scpcloudmonitoring.ProductResourceInfo, len(mdc.ProductResourceInfos))
	//
	//	for j, pri := range mdc.ProductResourceInfos {
	//		productResourceInfos[j] = scpcloudmonitoring.ProductResourceInfo{
	//			ObjectList:        pri.ObjectList,
	//			ProductResourceId: pri.ProductResourceId.ValueString(),
	//		}
	//	}
	//
	//	metricDataConditionOpenAPIV2[i] = scpcloudmonitoring.MetricDataConditionOpenAPIV2{
	//		MetricKey:            mdc.MetricKey.ValueString(),
	//		ObjectType:           mdc.ObjectType.ValueStringPointer(),
	//		ProductResourceInfos: productResourceInfos,
	//		StatisticsPeriod:     mdc.StatisticsPeriod.ValueInt32Pointer(),
	//		StatisticsTypeList:   mdc.StatisticsTypeList,
	//	}
	//}

	// 계층이 아닌 구조로 변경
	productResourceInfos := make([]scpcloudmonitoring.ProductResourceInfo, 1)
	productResourceInfos[0] = scpcloudmonitoring.ProductResourceInfo{
		ObjectList:        request.ObjectList,
		ProductResourceId: request.ProductResourceId.ValueString(),
	}

	metricDataConditionOpenAPIV2 := make([]scpcloudmonitoring.MetricDataConditionOpenAPIV2, 1)
	metricDataConditionOpenAPIV2[0] = scpcloudmonitoring.MetricDataConditionOpenAPIV2{
		MetricKey:            request.MetricKey.ValueString(),
		ObjectType:           request.ObjectType.ValueStringPointer(),
		ProductResourceInfos: productResourceInfos,
		StatisticsPeriod:     request.StatisticsPeriod.ValueInt32Pointer(),
		StatisticsTypeList:   request.StatisticsTypeList,
	}
	req = req.MetricDataSearchCriteriaOpenAPIV2(scpcloudmonitoring.MetricDataSearchCriteriaOpenAPIV2{
		IgnoreInvalid:        request.IgnoreInvalid.ValueStringPointer(),
		MetricDataConditions: metricDataConditionOpenAPIV2,
		QueryStartDt:         request.QueryStartDt.ValueString(),
		QueryEndDt:           request.QueryEndDt.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}
