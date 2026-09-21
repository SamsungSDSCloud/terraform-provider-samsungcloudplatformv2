package securitygroupv1d1

import (
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpsecuritygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/security-group/1.1"
)

const ServiceType = "scp-security-group"

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpsecuritygroup.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpsecuritygroup.NewAPIClient(config),
	}
}
