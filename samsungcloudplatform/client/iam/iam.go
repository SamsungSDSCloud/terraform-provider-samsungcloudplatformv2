package iam

import (
	"context"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/iam/1.0"
	"sort"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *iam.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: iam.NewAPIClient(config),
	}
}

func (client *Client) GetAccessKeyList(ctx context.Context, request AccessKeyDataSource) (*iam.ListAccessKeyResponse, error) {
	req := client.sdkClient.IamV1AccessKeysApiAPI.AccessKeyList(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.AccountId.IsNull() {
		req = req.AccountId(request.AccountId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateAccessKey(ctx context.Context, request AccessKeyResource) (*iam.AccessKeyResponse, error) {
	req := client.sdkClient.IamV1AccessKeysApiAPI.AccessKeyCreate(ctx)

	req = req.AccessKeyCreateRequest(iam.AccessKeyCreateRequest{
		AccessKeyType:     iam.AccessKeyTypeCreateRequestEnum(request.AccessKeyType.ValueString()),
		AccountId:         request.AccountId.ValueStringPointer(),
		Description:       *iam.NewNullableString(request.Description.ValueStringPointer()),
		Duration:          *iam.NewNullableString(request.Duration.ValueStringPointer()),
		ParentAccessKeyId: *iam.NewNullableString(request.ParentAccessKeyId.ValueStringPointer()),
		Passcode:          *iam.NewNullableString(request.Passcode.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccessKey(ctx context.Context, accessKeyId string) (*iam.AccessKeyResponse, error) {
	req := client.sdkClient.IamV1AccessKeysApiAPI.AccessKeyShow(ctx, accessKeyId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAccessKey(ctx context.Context, accessKeyId string, request AccessKeyResource) (*iam.AccessKeyResponse, error) {
	req := client.sdkClient.IamV1AccessKeysApiAPI.AccessKeySet(ctx, accessKeyId)

	req = req.AccessKeyUpdateRequest(iam.AccessKeyUpdateRequest{
		IsEnabled: request.IsEnabled.ValueBoolPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAccessKey(ctx context.Context, accessKeyId string) error {
	req := client.sdkClient.IamV1AccessKeysApiAPI.AccessKeyDelete(ctx, accessKeyId)

	_, err := req.Execute()
	return err
}

func (client *Client) GetEndpointList() (*iam.ListEndpointsResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.IamV1EndpointsApiAPI.ListEndpoints(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

var regions []string

func (client *Client) GetRegionList() []string {
	if len(regions) == 0 {
		ctx := context.Background()

		req := client.sdkClient.IamV1EndpointsApiAPI.ListEndpoints(ctx)

		resp, _, _ := req.Execute()

		regionMap := make(map[string]bool)
		var regions []string

		for _, endpoint := range resp.Endpoints {
			if !regionMap[endpoint.Region] {
				regionMap[endpoint.Region] = true
				regions = append(regions, endpoint.Region)
			}
		}

		sort.Slice(regions, func(i, j int) bool {
			return regions[i] < regions[j]
		})
	}

	return regions
}

func (client *Client) GetAccountId() (string, error) {
	ctx := context.Background()
	data, err := client.GetAccessKeyList(ctx, AccessKeyDataSource{})
	if err != nil {
		return "", err
	}

	if len(data.AccessKeys) == 0 {
		return "", fmt.Errorf("failed to find Account ID")
	}

	accessKey := data.AccessKeys[0]
	return accessKey.AccountId, nil
}
