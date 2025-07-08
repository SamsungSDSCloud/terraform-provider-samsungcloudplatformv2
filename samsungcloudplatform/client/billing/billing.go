package billing

import (
	"context"
	"fmt"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpbilling "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/billing"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpbilling.APIClient // scpnetwork 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: scpbilling.NewAPIClient(config),
	}
}

func (client *Client) GetPlannedComputeList(ctx context.Context, request PlannedComputeDataSourceIds) (*scpbilling.PlannedComputeListResponse, error) {
	req := client.sdkClient.BillingplanV1PlannedComputeApiAPI.ListPlannedComputes(ctx)

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.StartDate.IsNull() {
		req = req.StartDate(request.StartDate.String())
	}
	if !request.EndDate.IsNull() {
		req = req.EndDate(request.EndDate.String())
	}
	if !request.ServerType.IsNull() {
		req = req.ServerType(request.ServerType.String())
	}
	if !request.ContractId.IsNull() {
		req = req.ContractId(request.ContractId.String())
	}
	if len(request.ContractType) > 0 {
		var contractTypes []*string
		for _, contractType := range request.ContractType {
			tempString := contractType.String()
			contractTypes = append(contractTypes, &tempString)
		}
		req = req.ContractType(contractTypes)
	}
	if len(request.NextContractType) > 0 {
		var nextContractTypes []*string
		for _, nextContractType := range request.NextContractType {
			tempString := nextContractType.String()
			nextContractTypes = append(nextContractTypes, &tempString)
		}
		req = req.NextContractType(nextContractTypes)
	}
	if len(request.ServiceId) > 0 {
		var serviceIds []*string
		for _, serviceId := range request.ServiceId {
			tempString := serviceId.String()
			serviceIds = append(serviceIds, &tempString)
		}
		req = req.ServiceId(serviceIds)
	}
	if len(request.OsType) > 0 {
		var osTypes []*string
		for _, osType := range request.OsType {
			tempString := osType.String()
			osTypes = append(osTypes, &tempString)
		}
		req = req.OsType(osTypes)
	}
	if len(request.State) > 0 {
		var states []*string
		for _, state := range request.State {
			tempString := state.String()
			states = append(states, &tempString)
		}
		req = req.State(states)
	}
	if !request.CreatedBy.IsNull() {
		req = req.CreatedBy(request.CreatedBy.String())
	}
	if !request.ModifiedBy.IsNull() {
		req = req.ModifiedBy(request.ModifiedBy.String())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.String())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePlannedCompute(ctx context.Context, request PlannedComputeResource) (*scpbilling.PlannedComputeResponse, error) {
	req := client.sdkClient.BillingplanV1PlannedComputeApiAPI.CreatePlannedComputes(ctx)
	contractType, err := scpbilling.NewPlannedComputeContractEnumFromValue(request.ContractType.ValueString())
	if err != nil {
		return nil, err
	}
	osType, err := scpbilling.NewPlannedComputeOSTypeEnumFromValue(request.OsType.ValueString())
	if err != nil {
		return nil, err
	}

	var TagsObject []scpbilling.TagDTO

	for k, v := range request.Tags.Elements() {
		tagObject := scpbilling.TagDTO{
			Key:   types.StringValue(k).ValueStringPointer(),
			Value: v.(types.String).ValueStringPointer(),
		}

		TagsObject = append(TagsObject, tagObject)
	}
	req = req.PlannedComputeCreateRequest(scpbilling.PlannedComputeCreateRequest{
		AccountId:    request.AccountId.ValueString(),
		ContractType: contractType,
		OsType:       osType,
		ServerType:   request.ServerType.ValueString(),
		ServiceId:    request.ServiceId.ValueString(),
		ServiceName:  request.ServiceName.ValueStringPointer(),
		Tag:          TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPlannedCompute(ctx context.Context, plannedComputeId string) (*scpbilling.PlannedComputeResponse, error) {
	req := client.sdkClient.BillingplanV1PlannedComputeApiAPI.ShowPlannedCompute(ctx, plannedComputeId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePlannedCompute(ctx context.Context, plannedComputeId string, request PlannedComputeResource) (*scpbilling.PlannedComputeResponse, error) {
	req := client.sdkClient.BillingplanV1PlannedComputeApiAPI.UpdatePlannedCompute(ctx, plannedComputeId)
	contractType, err := scpbilling.NewPlannedComputeContractEnumFromValue(request.ContractType.ValueString())
	if err != nil {
		return nil, err
	}
	action, err := scpbilling.NewPlannedComputeChangeActionEnumFromValue(request.Action.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid action: %v", err)
	}
	nullableContractType := scpbilling.NewNullablePlannedComputeContractEnum(contractType)
	req = req.PlannedComputeChangeRequest(scpbilling.PlannedComputeChangeRequest{
		Action:       action,
		ContractType: *nullableContractType,
		ServerType:   *scpbilling.NewNullableString(request.ServerType.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}
