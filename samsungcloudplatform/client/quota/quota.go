package quota

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/quota/1.1"
	"math"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *quota.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: quota.NewAPIClient(config),
	}
}

func (client *Client) GetAccountQuotaList() (*quota.AccountQuotaListResponseV1dot1, error) {
	ctx := context.Background()

	req := client.sdkClient.QuotaV1AccountQuotasAPIsAPI.ListAccountQuota(ctx)
	req = req.Size(math.MaxInt32)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetAccountQuota(accountQuotaId string) (*quota.AccountQuotaShowResponseV1dot1, error) {
	ctx := context.Background()

	req := client.sdkClient.QuotaV1AccountQuotasAPIsAPI.ShowAccountQuota(ctx, accountQuotaId)
	resp, _, err := req.Execute()
	return resp, err
}
