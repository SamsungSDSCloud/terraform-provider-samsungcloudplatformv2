package securitygroupv1d1

import (
	"context"

	scpsecuritygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/security-group/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (client *Client) CreateSecurityGroupRule(ctx context.Context, request SecurityGroupRuleResource) (*scpsecuritygroup.SecurityGroupRuleShowResponseV1Dot1, error) {
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.CreateSecurityGroupRule(ctx) // 호출을 위한 구조체를 반환 받는다
	description := request.Description.ValueString()
	descriptionNS := scpsecuritygroup.NullableString{}
	descriptionNS.Set(&description)

	protocol := request.Protocol.ValueString()
	protocolNS := scpsecuritygroup.NullableString{}
	if protocol != "" {
		protocolNS.Set(&protocol)
	}

	remoteIpPrefix := request.RemoteIpPrefix.ValueString()
	remoteIpPrefixNS := scpsecuritygroup.NullableString{}
	if remoteIpPrefix != "" {
		remoteIpPrefixNS.Set(&remoteIpPrefix)
	}

	remoteGroupId := request.RemoteGroupId.ValueString()
	remoteGroupIdNS := scpsecuritygroup.NullableString{}
	if remoteGroupId != "" {
		remoteGroupIdNS.Set(&remoteGroupId)
	}

	ethertype := request.Ethertype.ValueString()
	ethertypeNS := scpsecuritygroup.NullableString{}
	if ethertype != "" {
		ethertypeNS.Set(&ethertype)
	}

	req = req.SecurityGroupRuleCreateRequestV1Dot1(scpsecuritygroup.SecurityGroupRuleCreateRequestV1Dot1{
		SecurityGroupId:      request.SecurityGroupId.ValueString(),
		Ethertype:            ethertypeNS,
		Protocol:             protocolNS,
		PortRangeMin:         *scpsecuritygroup.NewNullableInt32(request.PortRangeMin.ValueInt32Pointer()),
		PortRangeMax:         *scpsecuritygroup.NewNullableInt32(request.PortRangeMax.ValueInt32Pointer()),
		RemoteAddressGroupId: *scpsecuritygroup.NewNullableString(request.RemoteAddressGroupId.ValueStringPointer()),
		RemoteIpPrefix:       remoteIpPrefixNS,
		RemoteGroupId:        remoteGroupIdNS,
		Description:          descriptionNS,
		Direction:            request.Direction.ValueString(),
	})

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetSecurityGroupRule(ctx context.Context, securityGroupRuleId string) (*scpsecuritygroup.SecurityGroupRuleShowResponseV1Dot1, error) {
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.ShowSecurityGroupRule(ctx, securityGroupRuleId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetSecurityGroupRuleList(page types.Int32, size types.Int32, sort types.String, id types.String, securityGroupId types.String, remoteIpPrefix types.String, remoteGroupId types.String, description types.String, direction types.String, service types.String) (*scpsecuritygroup.SecurityGroupRuleListResponseV1Dot1, error) {

	ctx := context.Background()
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.ListSecurityGroupRules(ctx)
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
	if !securityGroupId.IsNull() {
		req = req.SecurityGroupId(securityGroupId.ValueString())
	}
	if !remoteIpPrefix.IsNull() {
		req = req.RemoteIpPrefix(remoteIpPrefix.ValueString())
	}
	if !remoteGroupId.IsNull() {
		req = req.RemoteGroupId(remoteGroupId.ValueString())
	}
	if !description.IsNull() {
		req = req.Description(description.ValueString())
	}
	if !direction.IsNull() {
		req = req.Direction(direction.ValueString())
	}
	if !service.IsNull() {
		req = req.Service(service.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}
