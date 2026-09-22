package firewall

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/firewall/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	config    *scpsdk.Configuration
	sdkClient *scpfirewall.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: scpfirewall.NewAPIClient(config),
	}
}

//------------------- Firewall -------------------//

func (client *Client) ShowFirewall(ctx context.Context, firewallId string) (*scpfirewall.FirewallShowResponseV1Dot1, error) {
	req := client.sdkClient.FirewallV1FirewallApiAPI.ShowFirewall(ctx, firewallId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetFirewall(ctx context.Context, firewallResource Resource) (*scpfirewall.FirewallShowResponseV1Dot1, error) {
	req := client.sdkClient.FirewallV1FirewallApiAPI.SetFirewall(ctx, firewallResource.Id.ValueString())

	fwFlavorName := scpfirewall.FirewallFlavorType(firewallResource.FlavorName.ValueString())
	fwSetReq := scpfirewall.FirewallSetRequest{
		FlavorName: &fwFlavorName,
		Loggable:   firewallResource.Loggable.ValueBoolPointer(),
	}

	req = req.FirewallSetRequest(fwSetReq)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetFirewallList(page types.Int32, size types.Int32, sort types.String, name types.String,
	vpcName types.String, productType types.List, state types.List) (*scpfirewall.FirewallListResponseV1Dot1, error) {

	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallApiAPI.ListFirewalls(ctx)

	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !name.IsNull() {
		req = req.Name(name.ValueString())
	}
	if !vpcName.IsNull() {
		req = req.VpcName(vpcName.ValueString())
	}
	if !productType.IsNull() {
		var productTypes []string
		productType.ElementsAs(ctx, &productTypes, false)
		reqProductTypeList := make([]scpfirewall.FirewallProductType, 0, len(productTypes))
		for _, productType := range productTypes {
			reqProductTypeList = append(reqProductTypeList, scpfirewall.FirewallProductType(productType))
		}
		req = req.ProductType(reqProductTypeList)
	}
	if !state.IsNull() {
		var states []string
		state.ElementsAs(ctx, &states, false)
		reqStateList := make([]scpfirewall.FirewallState, 0, len(states))
		for _, state := range states {
			reqStateList = append(reqStateList, scpfirewall.FirewallState(state))
		}
		req = req.State(reqStateList)
	}

	resp, _, err := req.Execute()
	return resp, err
}
