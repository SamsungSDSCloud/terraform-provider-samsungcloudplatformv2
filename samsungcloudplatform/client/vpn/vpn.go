package vpn

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/vpn/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	config    *scpsdk.Configuration
	sdkClient *scpvpn.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		config:    config,
		sdkClient: scpvpn.NewAPIClient(config),
	}
}

//------------ VPN Gateway -------------------//

func (client *Client) CreateVpnGateway(ctx context.Context, request VpnGatewayResource) (*scpvpn.VpnGatewayShowResponse, error) {
	req := client.sdkClient.VpnV1VpnGatewaysApiAPI.CreateVpnGateway(ctx)
	tags := convertToTags(request.Tags.Elements())

	description := request.Description.ValueString()
	descriptionES := ""
	descriptionNS := scpvpn.NullableString{}

	if description == "" {
		descriptionNS.Set(&descriptionES)
	} else {
		descriptionNS.Set(&description)
	}

	req = req.VpnGatewayCreateRequest(scpvpn.VpnGatewayCreateRequest{
		Description: descriptionNS,
		IpAddress:   request.IpAddress.ValueString(),
		IpId:        request.IpId.ValueString(),
		IpType:      request.IpType.ValueString(),
		Name:        request.Name.ValueString(),
		Tags:        tags,
		VpcId:       request.VpcId.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpnGateway(ctx context.Context, vpnGatewayId string) (*scpvpn.VpnGatewayShowResponse, error) {
	req := client.sdkClient.VpnV1VpnGatewaysApiAPI.ShowVpnGateway(ctx, vpnGatewayId) // 호출을 위한 구조체를 반환 받는다.

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateVpnGateway(ctx context.Context, vpnGatewayId string, request VpnGatewayResource) (*scpvpn.VpnGatewayShowResponse, error) {
	req := client.sdkClient.VpnV1VpnGatewaysApiAPI.SetVpnGateway(ctx, vpnGatewayId) // 호출을 위한 구조체를 반환 받는다.
	description := request.Description.ValueString()
	descriptionNS := scpvpn.NullableString{}
	descriptionNS.Set(&description)

	req = req.VpnGatewaySetRequest(scpvpn.VpnGatewaySetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) DeleteVpnGateway(ctx context.Context, vpnGatewayId string) error {
	req := client.sdkClient.VpnV1VpnGatewaysApiAPI.DeleteVpnGateway(ctx, vpnGatewayId) // 호출을 위한 구조체를 반환 받는다.

	_, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return err
}

func (client *Client) GetVpnGatewayList(ctx context.Context, size types.Int32, page types.Int32, sort types.String, name types.String,
	ipAddress types.String, vpcName types.String, vpcId types.String) (*scpvpn.VpnGatewayListResponse, error) {
	req := client.sdkClient.VpnV1VpnGatewaysApiAPI.ListVpnGateways(ctx)
	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !name.IsNull() {
		req = req.Name(name.ValueString())
	}
	if !ipAddress.IsNull() {
		req = req.IpAddress(ipAddress.ValueString())
	}
	if !vpcName.IsNull() {
		req = req.VpcName(vpcName.ValueString())
	}
	if !vpcId.IsNull() {
		req = req.VpcId(vpcId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ------------ VPN Tunnel -------------------//
func (client *Client) CreateVpnTunnel(ctx context.Context, request VpnTunnelResource) (*scpvpn.VpnTunnelShowResponse, error) {
	req := client.sdkClient.VpnV1VpnTunnelsApiAPI.CreateVpnTunnel(ctx)
	tags := convertToTags(request.Tags.Elements())

	description := request.Description.ValueString()
	descriptionES := ""
	descriptionNS := scpvpn.NullableString{}

	if description == "" {
		descriptionNS.Set(&descriptionES)
	} else {
		descriptionNS.Set(&description)
	}

	createRequest := scpvpn.VpnTunnelCreateRequest{
		Description:  descriptionNS,
		Name:         request.Name.ValueString(),
		VpnGatewayId: request.VpnGatewayId.ValueString(),
		Phase1: scpvpn.VpnPhase1CreateRequest{
			Phase1DiffieHellmanGroups: convertDiffieHellmanGroups(request.Phase1.DiffieHellmanGroups),
			Phase1Encryptions:         convertEncryptions(request.Phase1.Encryptions),
			DpdRetryInterval:          request.Phase1.DpdRetryInterval.ValueInt32(),
			IkeVersion:                request.Phase1.IkeVersion.ValueInt32(),
			Phase1LifeTime:            request.Phase1.LifeTime.ValueInt32(),
			PeerGatewayIp:             request.Phase1.PeerGatewayIp.ValueString(),
			PreSharedKey:              request.Phase1.PreSharedKey.ValueString(),
		},
		Phase2: scpvpn.VpnPhase2CreateRequest{
			Phase2DiffieHellmanGroups: convertDiffieHellmanGroups(request.Phase2.DiffieHellmanGroups),
			Phase2Encryptions:         convertEncryptions(request.Phase2.Encryptions),
			Phase2LifeTime:            request.Phase2.LifeTime.ValueInt32(),
			PerfectForwardSecrecy:     scpvpn.VpnPerfectForwardSecrecyType(request.Phase2.PerfectForwardSecrecy.ValueString()),
			RemoteSubnet:              request.Phase2.RemoteSubnet.ValueString(),
		},
		Tags: tags,
	}

	req = req.VpnTunnelCreateRequest(createRequest)

	resp, _, err := req.Execute()
	return resp, err
}

func convertToTags(elements map[string]attr.Value) []scpvpn.Tag {
	var tags []scpvpn.Tag
	for k, v := range elements {
		tagObject := scpvpn.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func convertDiffieHellmanGroups(groups []types.Int32) []int32 {
	result := make([]int32, len(groups))
	for i, group := range groups {
		result[i] = group.ValueInt32()
	}
	return result
}

func convertEncryptions(encryptions []types.String) []string {
	result := make([]string, len(encryptions))
	for i, encryption := range encryptions {
		result[i] = encryption.ValueString()
	}
	return result
}

func (client *Client) GetVpnTunnel(ctx context.Context, vpnTunnelId string) (*scpvpn.VpnTunnelShowWithStatusResponse, error) {
	req := client.sdkClient.VpnV1VpnTunnelsApiAPI.ShowVpnTunnel(ctx, vpnTunnelId) // 호출을 위한 구조체를 반환 받는다.

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateVpnTunnel(ctx context.Context, vpnTunnelId string, request VpnTunnelResource) (*scpvpn.VpnTunnelShowResponse, error) {
	req := client.sdkClient.VpnV1VpnTunnelsApiAPI.SetVpnTunnel(ctx, vpnTunnelId) // 호출을 위한 구조체를 반환 받는다.
	description := request.Description.ValueString()
	descriptionNS := scpvpn.NullableString{}
	descriptionNS.Set(&description)

	perfectForwardSecrecy := scpvpn.VpnPerfectForwardSecrecyType(request.Phase2.PerfectForwardSecrecy.ValueString())

	req = req.VpnTunnelSetRequest(scpvpn.VpnTunnelSetRequest{
		Description: *scpvpn.NewNullableString(request.Description.ValueStringPointer()),
		Phase1: *scpvpn.NewNullableVpnPhase1SetRequest(&scpvpn.VpnPhase1SetRequest{
			Phase1DiffieHellmanGroups: convertDiffieHellmanGroups(request.Phase1.DiffieHellmanGroups),
			Phase1Encryptions:         convertEncryptions(request.Phase1.Encryptions),
			DpdRetryInterval:          *scpvpn.NewNullableInt32(request.Phase1.DpdRetryInterval.ValueInt32Pointer()),
			IkeVersion:                *scpvpn.NewNullableInt32(request.Phase1.IkeVersion.ValueInt32Pointer()),
			Phase1LifeTime:            *scpvpn.NewNullableInt32(request.Phase1.LifeTime.ValueInt32Pointer()),
			PeerGatewayIp:             *scpvpn.NewNullableString(request.Phase1.PeerGatewayIp.ValueStringPointer()),
			PreSharedKey:              *scpvpn.NewNullableString(request.Phase1.PreSharedKey.ValueStringPointer()),
		}),
		Phase2: *scpvpn.NewNullableVpnPhase2SetRequest(&scpvpn.VpnPhase2SetRequest{
			Phase2DiffieHellmanGroups: convertDiffieHellmanGroups(request.Phase2.DiffieHellmanGroups),
			Phase2Encryptions:         convertEncryptions(request.Phase2.Encryptions),
			Phase2LifeTime:            *scpvpn.NewNullableInt32(request.Phase2.LifeTime.ValueInt32Pointer()),
			PerfectForwardSecrecy:     *scpvpn.NewNullableVpnPerfectForwardSecrecyType(&perfectForwardSecrecy),
			RemoteSubnet:              *scpvpn.NewNullableString(request.Phase2.RemoteSubnet.ValueStringPointer()),
		}),
	})

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) DeleteVpnTunnel(ctx context.Context, vpnTunnelId string) error {
	req := client.sdkClient.VpnV1VpnTunnelsApiAPI.DeleteVpnTunnel(ctx, vpnTunnelId) // 호출을 위한 구조체를 반환 받는다.

	_, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return err
}

func (client *Client) GetVpnTunnelList(ctx context.Context, page types.Int32, size types.Int32, sort types.String, name types.String,
	vpnGatewayId types.String, vpnGatewayName types.String, peerGatewayIp types.String, remoteSubnet types.String) (*scpvpn.VpnTunnelListResponse, error) {

	req := client.sdkClient.VpnV1VpnTunnelsApiAPI.ListVpnTunnels(ctx)
	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !name.IsNull() {
		req = req.Name(name.ValueString())
	}
	if !vpnGatewayId.IsNull() {
		req = req.VpnGatewayId(vpnGatewayId.ValueString())
	}
	if !vpnGatewayName.IsNull() {
		req = req.VpnGatewayName(vpnGatewayName.ValueString())
	}
	if !peerGatewayIp.IsNull() {
		req = req.PeerGatewayIp(peerGatewayIp.ValueString())
	}
	if !remoteSubnet.IsNull() {
		req = req.RemoteSubnet(remoteSubnet.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}
