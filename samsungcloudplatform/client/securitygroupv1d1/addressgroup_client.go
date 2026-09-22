package securitygroupv1d1

import (
	"context"

	scpsecuritygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/security-group/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//------------------- Address Group -------------------//

func (client *Client) CreateAddressGroup(ctx context.Context, request AddressGroupResource) (*scpsecuritygroup.AddressGroupShowResponse, error) {
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.CreateAddressGroup(ctx)

	name := request.Name.ValueString()
	description := request.Description.ValueString()
	descriptionNS := scpsecuritygroup.NullableString{}
	descriptionNS.Set(&description)

	var addresses []string
	if !request.Addresses.IsNull() && !request.Addresses.IsUnknown() {
		for _, addr := range request.Addresses.Elements() {
			addresses = append(addresses, addr.(types.String).ValueString())
		}
	}

	tags := convertToTags(request.Tags.Elements())

	req = req.AddressGroupCreateRequest(scpsecuritygroup.AddressGroupCreateRequest{
		Name:        name,
		Description: descriptionNS,
		Addresses:   addresses,
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAddressGroup(ctx context.Context, addressGroupId string) (*scpsecuritygroup.AddressGroupShowResponse, error) {
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.ShowAddressGroup(ctx, addressGroupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAddressGroup(ctx context.Context, addressGroupId string, request AddressGroupResource) error {
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.SetAddressGroup(ctx, addressGroupId)

	req = req.AddressGroupSetRequest(scpsecuritygroup.AddressGroupSetRequest{
		Description: *scpsecuritygroup.NewNullableString(request.Description.ValueStringPointer()),
	})

	_, err := req.Execute()
	return err
}

func (client *Client) DeleteAddressGroup(ctx context.Context, addressGroupId string) error {
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.DeleteAddressGroup(ctx, addressGroupId)

	_, err := req.Execute()
	return err
}

func (client *Client) ListAddressGroups(page types.Int32, size types.Int32, sort types.String, id types.String, name types.String) (*scpsecuritygroup.AddressGroupListResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.SecurityGroupV1AddressGroupApiAPI.ListAddressGroups(ctx)

	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !id.IsNull() {
		req = req.Id(id.ValueString())
	}
	if !name.IsNull() {
		req = req.Name(name.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func convertToTags(elements map[string]attr.Value) []scpsecuritygroup.Tag {
	var tags []scpsecuritygroup.Tag
	for k, v := range elements {
		tagObject := scpsecuritygroup.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}
