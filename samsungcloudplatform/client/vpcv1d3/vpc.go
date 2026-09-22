package vpcv1d3

import (
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

const ServiceType = "scp-vpc"

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *vpc.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: vpc.NewAPIClient(config),
	}
}
