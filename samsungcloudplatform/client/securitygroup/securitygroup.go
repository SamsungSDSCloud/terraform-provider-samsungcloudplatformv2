package securitygroup

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpsecuritygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/security-group/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	config    *scpsdk.Configuration
	sdkClient *scpsecuritygroup.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: scpsecuritygroup.NewAPIClient(config),
	}
}

//------------------- Security Group -------------------//

func (client *Client) GetSecurityGroupList(page types.Int32, size types.Int32, sort types.String, name types.String, id types.String) (*scpsecuritygroup.SecurityGroupListResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.SecurityGroupV1SecurityGroupApiAPI.ListSecurityGroups(ctx)

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

func (client *Client) CreateSecurityGroup(ctx context.Context, request SecurityGroupResource) (*scpsecuritygroup.SecurityGroupShowResponse, error) {
	req := client.sdkClient.SecurityGroupV1SecurityGroupApiAPI.CreateSecurityGroup(ctx) // 호출을 위한 구조체를 반환 받는다

	name := request.Name.ValueString()
	description := request.Description.ValueString()
	descriptionNS := scpsecuritygroup.NullableString{}
	descriptionNS.Set(&description)
	loggable := request.Loggable.ValueBool()
	loggableNS := scpsecuritygroup.NullableBool{}
	loggableNS.Set(&loggable)
	tags := convertToTags(request.Tags.Elements())

	req = req.SecurityGroupCreateRequest(scpsecuritygroup.SecurityGroupCreateRequest{
		Name:        name,
		Description: descriptionNS,
		Loggable:    loggableNS,
		Tags:        tags,
	})

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetSecurityGroup(ctx context.Context, securityGroupId string) (*scpsecuritygroup.SecurityGroupShowResponse, error) {
	req := client.sdkClient.SecurityGroupV1SecurityGroupApiAPI.ShowSecurityGroup(ctx, securityGroupId) // 호출을 위한 구조체를 반환 받는다.

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateSecurityGroup(ctx context.Context, securityGroupId string, request SecurityGroupResource) error {
	req := client.sdkClient.SecurityGroupV1SecurityGroupApiAPI.SetSecurityGroup(ctx, securityGroupId) // 호출을 위한 구조체를 반환 받는다.
	loggable := request.Loggable.ValueBool()
	loggableNS := scpsecuritygroup.NullableBool{}
	loggableNS.Set(&loggable)

	req = req.SecurityGroupSetRequest(scpsecuritygroup.SecurityGroupSetRequest{
		Description: *scpsecuritygroup.NewNullableString(request.Description.ValueStringPointer()),
		Loggable:    loggableNS,
	})

	_, err := req.Execute()
	return err
}

func (client *Client) DeleteSecurityGroup(ctx context.Context, securityGroupId string) error {
	req := client.sdkClient.SecurityGroupV1SecurityGroupApiAPI.DeleteSecurityGroup(ctx, securityGroupId) // 호출을 위한 구조체를 반환 받는다.

	_, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return err
}

//------------------- Security Group Rule -------------------//

func (client *Client) GetSecurityGroupRuleList(page types.Int32, size types.Int32, sort types.String, id types.String, securityGroupId types.String, remoteIpPrefix types.String, remoteGroupId types.String, description types.String, direction types.String, service types.String) (*scpsecuritygroup.SecurityGroupRuleListResponse, error) {

	ctx := context.Background()
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.ListSecurityGroupRules(ctx) // 호출을 위한 구조체를 반환 받는다.
	if !size.IsNull() {                                                                        // Null 값이 아닌 경우 인자를 추가하도록 한다.
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
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateSecurityGroupRule(ctx context.Context, request SecurityGroupRuleResource) (*scpsecuritygroup.SecurityGroupRuleShowResponse, error) {
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

	req = req.SecurityGroupRuleCreateRequest(scpsecuritygroup.SecurityGroupRuleCreateRequest{
		SecurityGroupId: request.SecurityGroupId.ValueString(),
		Ethertype:       ethertypeNS,
		Protocol:        protocolNS,
		PortRangeMin:    *scpsecuritygroup.NewNullableInt32(request.PortRangeMin.ValueInt32Pointer()),
		PortRangeMax:    *scpsecuritygroup.NewNullableInt32(request.PortRangeMax.ValueInt32Pointer()),
		RemoteIpPrefix:  remoteIpPrefixNS,
		RemoteGroupId:   remoteGroupIdNS,
		Description:     descriptionNS,
		Direction:       request.Direction.ValueString(),
	})

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetSecurityGroupRule(ctx context.Context, securityGroupRuleId string) (*scpsecuritygroup.SecurityGroupRuleShowResponse, error) {
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.ShowSecurityGroupRule(ctx, securityGroupRuleId) // 호출을 위한 구조체를 반환 받는다.

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) DeleteSecurityGroupRule(ctx context.Context, securityGroupRuleId string) error {
	req := client.sdkClient.SecurityGroupV1SecurityGroupRuleApiAPI.DeleteSecurityGroupRule(ctx, securityGroupRuleId) // 호출을 위한 구조체를 반환 받는다.

	_, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return err
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
