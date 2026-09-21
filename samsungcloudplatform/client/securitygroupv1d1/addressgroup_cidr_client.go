package securitygroupv1d1

import (
	"context"

	scpsecuritygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/security-group/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//------------------- Address Group CIDR -------------------//

func (client *Client) ListAddressGroupCidrs(ctx context.Context, request AddressGroupCidrDataSources) (*scpsecuritygroup.AddressGroupCidrListResponse, error) {
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.ListAddressGroupCidrs(ctx, request.AddressGroupId.ValueString())

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Address.IsNull() {
		req = req.Address(request.Address.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) AddAddressGroupCidrs(ctx context.Context, addressGroupId string, request AddressGroupCidrResource) (*scpsecuritygroup.AddressGroupCidrsAddRequest, error) {
	var addresses []string
	if !request.Addresses.IsNull() && !request.Addresses.IsUnknown() {
		for _, addr := range request.Addresses.Elements() {
			addresses = append(addresses, addr.(types.String).ValueString())
		}
	}

	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.AddAddressGroupCidrs(ctx, addressGroupId)

	resp, _, err := req.AddressGroupCidrsAddRequest(scpsecuritygroup.AddressGroupCidrsAddRequest{
		Addresses: addresses,
	}).Execute()
	return resp, err
}

func (client *Client) RemoveAddressGroupCidrs(ctx context.Context, addressGroupId string, request AddressGroupCidrResource) (*scpsecuritygroup.AddressGroupCidrsAddRequest, error) {
	var addresses []string
	if !request.Addresses.IsNull() && !request.Addresses.IsUnknown() {
		for _, addr := range request.Addresses.Elements() {
			addresses = append(addresses, addr.(types.String).ValueString())
		}
	}

	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.RemoveAddressGroupCidrs(ctx, addressGroupId)

	resp, _, err := req.AddressGroupCidrsAddRequest(scpsecuritygroup.AddressGroupCidrsAddRequest{
		Addresses: addresses,
	}).Execute()
	return resp, err
}
