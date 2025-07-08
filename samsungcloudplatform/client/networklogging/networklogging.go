package networklogging

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpnetworklogging "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/network-logging/1.0"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpnetworklogging.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpnetworklogging.NewAPIClient(config),
	}
}

//------------------- Network Logging Storage -------------------//

func (client *Client) CreateNetworkLoggingStorage(ctx context.Context, request NetworkLoggingStorageResource) (*scpnetworklogging.NetworkLoggingStorageShowResponse, error) {
	req := client.sdkClient.NetworkLoggingV1NetworkLoggingStorageApiAPI.CreateNetworkLoggingStorage(ctx)

	req = req.NetworkLoggingStorageCreateRequest(scpnetworklogging.NetworkLoggingStorageCreateRequest{
		BucketName:   request.BucketName.ValueString(),
		ResourceType: scpnetworklogging.NetworkLoggingResourceType(request.ResourceType.ValueString()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteNetworkLoggingStorage(ctx context.Context, networkLoggingStorageId string) error {
	req := client.sdkClient.NetworkLoggingV1NetworkLoggingStorageApiAPI.DeleteNetworkLoggingStorage(ctx, networkLoggingStorageId)

	_, err := req.Execute()
	return err
}

func (client *Client) GetNetworkLoggingStorageList(ctx context.Context, request NetworkLoggingStorageDataSource) (*scpnetworklogging.NetworkLoggingStorageListResponse, error) {
	req := client.sdkClient.NetworkLoggingV1NetworkLoggingStorageApiAPI.ListNetworkLoggingStorages(ctx)

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.ResourceType.IsNull() {
		req = req.ResourceType(scpnetworklogging.NetworkLoggingResourceType(request.ResourceType.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

//------------------- Network Logging Configuration -------------------//

func (client *Client) GetNetworkLoggingConfigurationList(ctx context.Context, request NetworkLoggingConfigurationDataSource) (*scpnetworklogging.NetworkLoggingConfigurationListResponse, error) {
	req := client.sdkClient.NetworkLoggingV1NetworkLoggingConfigurationApiAPI.ListNetworkLoggingConfigurations(ctx)

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.ResourceId.IsNull() {
		req = req.ResourceId(request.ResourceId.ValueString())
	}
	if !request.ResourceType.IsNull() {
		req = req.ResourceType(scpnetworklogging.NetworkLoggingResourceType(request.ResourceType.ValueString()))
	}
	if !request.ResourceName.IsNull() {
		req = req.ResourceName(request.ResourceName.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}
