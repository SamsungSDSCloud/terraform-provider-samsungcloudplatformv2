package loadbalancer

import (
	"context"
	"encoding/json"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/loadbalancer/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *loadbalancer.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: loadbalancer.NewAPIClient(config),
	}
}

func convertToTags(elements map[string]attr.Value) []loadbalancer.Tag {
	var tags []loadbalancer.Tag
	for k, v := range elements {
		tagObject := loadbalancer.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

// ------------ Load Balancer -------------------//
func (client *Client) GetLoadbalancerList(ctx context.Context, request LoadbalancerDataSource) (*loadbalancer.LoadbalancerListResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.ListLoadbalancers(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.ServiceIp.IsNull() {
		req = req.ServiceIp(request.ServiceIp.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLoadbalancer(ctx context.Context, loadbalancerId string) (*loadbalancer.LoadbalancerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.ShowLoadbalancer(ctx, loadbalancerId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateLoadbalancer(ctx context.Context, request LoadbalancerResource) (*loadbalancer.LoadbalancerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.CreateLoadbalancer(ctx)

	loadbalancerCreateRequest := request.LoadbalancerCreate

	loadbalancerElement := loadbalancer.LoadbalancerCreateRequest{
		Loadbalancer: loadbalancer.LoadbalancerCreateRequestDetail{
			Description:            *loadbalancer.NewNullableString(loadbalancerCreateRequest.Description.ValueStringPointer()),
			FirewallEnabled:        *loadbalancer.NewNullableBool(loadbalancerCreateRequest.FirewallEnabled.ValueBoolPointer()),
			FirewallLoggingEnabled: *loadbalancer.NewNullableBool(loadbalancerCreateRequest.FirewallLoggingEnabled.ValueBoolPointer()),
			LayerType:              loadbalancerCreateRequest.LayerType.ValueString(),
			Name:                   loadbalancerCreateRequest.Name.ValueString(),
			PublicipId:             *loadbalancer.NewNullableString(loadbalancerCreateRequest.PublicipId.ValueStringPointer()),
			ServiceIp:              *loadbalancer.NewNullableString(loadbalancerCreateRequest.ServiceIp.ValueStringPointer()),
			SubnetId:               loadbalancerCreateRequest.SubnetId.ValueString(),
			VpcId:                  loadbalancerCreateRequest.VpcId.ValueString(),
		},
	}

	req = req.LoadbalancerCreateRequest(loadbalancerElement)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLoadbalancer(ctx context.Context, loadbalancerId string) error {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.DeleteLoadbalancer(ctx, loadbalancerId)

	_, err := req.Execute()
	return err
}

func (client *Client) UpdateLoadbalancer(ctx context.Context, loadbalancerId string, request LoadbalancerResource) (*loadbalancer.LoadbalancerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.SetLoadbalancer(ctx, loadbalancerId)

	req = req.LoadbalancerUpdateRequest(loadbalancer.LoadbalancerUpdateRequest{
		Loadbalancer: loadbalancer.LoadbalancerUpdateRequestDetail{
			Description: request.LoadbalancerCreate.Description.ValueString(),
		},
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLoadbalancerPublicNatIp(ctx context.Context, request LoadbalancerPublicNatIpResource) (*loadbalancer.StaticNatCreateResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.CreateLoadbalancerPublicNatIp(ctx, request.LoadbalancerId.ValueString())

	loadbalancerNatCreateRequest := request.LoadbalancerNatCreate

	loadbalancerNat := loadbalancer.StaticNatCreateRequest{
		StaticNat: loadbalancer.StaticNatCreateRequestDetail{
			PublicipId: loadbalancerNatCreateRequest.PublicipId.ValueString(),
		},
	}

	req = req.StaticNatCreateRequest(loadbalancerNat)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLoadbalancerPublicNatIp(ctx context.Context, loadbalancerId string) error {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.DeleteLoadbalancerPublicNatIp(ctx, loadbalancerId)

	_, err := req.Execute()
	return err
}

//------------ LB Server Group -------------------//

func (client *Client) GetLbServerGroupList(ctx context.Context, request LbServerGroupDataSource) (*loadbalancer.LbServerGroupListResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.ListLbServerGroups(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	//if !request.Protocol.IsNull() {
	//	req = req.Protocol(loadbalancer.Protocol{
	//		LbServerGroupProtocol:        loadbalancer.LbServerGroupProtocol(request.Protocol.String()).Ptr(),
	//		ArrayOfLbServerGroupProtocol: nil,
	//	})
	//}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbServerGroup(ctx context.Context, lbServerGroupId string) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.ShowLbServerGroup(ctx, lbServerGroupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLbServerGroup(ctx context.Context, request LbServerGroupResource) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.CreateLbServerGroup(ctx)

	lbServerGroup := request.LbServerGroupCreate

	lbServerGroupElement := loadbalancer.LbServerGroupCreate{
		Name:            lbServerGroup.Name.ValueString(),
		VpcId:           lbServerGroup.VpcId.ValueString(),
		SubnetId:        lbServerGroup.SubnetId.ValueString(),
		Protocol:        loadbalancer.LbServerGroupProtocol(lbServerGroup.Protocol.ValueString()),
		LbMethod:        loadbalancer.LbServerGroupLbMethod(lbServerGroup.LbMethod.ValueString()),
		Description:     *loadbalancer.NewNullableString(lbServerGroup.Description.ValueStringPointer()),
		LbHealthCheckId: *loadbalancer.NewNullableString(lbServerGroup.LbHealthCheckId.ValueStringPointer()),
		Tags:            convertToTags(lbServerGroup.Tags.Elements()),
	}

	req = req.LbServerGroupCreateRequest(loadbalancer.LbServerGroupCreateRequest{
		LbServerGroup: lbServerGroupElement,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLbServerGroup(ctx context.Context, lbServerGroupId string, request LbServerGroupResource) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.SetLbServerGroup(ctx, lbServerGroupId)

	lbServerGroup := request.LbServerGroupCreate

	lbServerGroupElement := loadbalancer.LbServerGroupSet{
		LbMethod:        *loadbalancer.NewNullableLbServerGroupLbMethod((*loadbalancer.LbServerGroupLbMethod)(lbServerGroup.LbMethod.ValueStringPointer())),
		Description:     *loadbalancer.NewNullableString(lbServerGroup.Description.ValueStringPointer()),
		LbHealthCheckId: *loadbalancer.NewNullableString(lbServerGroup.LbHealthCheckId.ValueStringPointer()),
	}

	req = req.LbServerGroupSetRequest(loadbalancer.LbServerGroupSetRequest{
		LbServerGroup: lbServerGroupElement,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLbServerGroup(ctx context.Context, lbServerGroupId string) error {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.DeleteLbServerGroup(ctx, lbServerGroupId)

	_, err := req.Execute()
	return err
}

//------------ LB Member -------------------//

func (client *Client) GetLbMemberList(ctx context.Context, request LbMemberDataSource) (*loadbalancer.MemberWithHealthStateListResponse, error) {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.ListLbServerGroupMembers(ctx, request.LbServerGroupId.ValueString())

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.MemberIp.IsNull() {
		req = req.MemberIp(request.MemberIp.ValueString())
	}
	if !request.MemberPort.IsNull() {
		req = req.MemberPort(request.MemberPort.ValueInt32())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbMember(ctx context.Context, lbServerGroupId string, memberId string) (*loadbalancer.MemberShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.ShowLbServerGroupMember(ctx, lbServerGroupId, memberId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLbMember(ctx context.Context, request LbMembersResource) (*loadbalancer.MemberListResponse, error) {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.AddLbServerGroupMembers(ctx, request.LbServerGroupId.ValueString())

	lbMember := request.LbMemberCreate
	var convertedLbMembers []loadbalancer.MemberCreateRequest

	convertedLbMember := loadbalancer.MemberCreateRequest{
		Name:         lbMember.Name.ValueString(),
		MemberIp:     lbMember.MemberIp.ValueString(),
		MemberPort:   lbMember.MemberPort.ValueInt32(),
		ObjectType:   lbMember.ObjectType.ValueString(),
		ObjectId:     lbMember.ObjectId.ValueStringPointer(),
		MemberWeight: lbMember.MemberWeight.ValueInt32Pointer(),
	}
	convertedLbMembers = append(convertedLbMembers, convertedLbMember)
	req = req.MemberListCreateRequest(loadbalancer.MemberListCreateRequest{
		Members: convertedLbMembers,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLbMember(ctx context.Context, lbServerGroupId string, memberId string, request LbMemberResource) (*loadbalancer.MemberShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.SetLbServerGroupMember(ctx, lbServerGroupId, memberId)

	req = req.MemberSetRequest(loadbalancer.MemberSetRequest{
		Member: loadbalancer.MemberSet{
			MemberPort:   *loadbalancer.NewNullableInt32(request.LbMemberSet.MemberPort.ValueInt32Pointer()),
			MemberWeight: *loadbalancer.NewNullableInt32(request.LbMemberSet.MemberWeight.ValueInt32Pointer()),
			//MemberState:  *loadbalancer.NewNullableStatusType((*loadbalancer.StatusType)(request.LbMemberSet.MemberState.ValueStringPointer())),
		},
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLbMember(ctx context.Context, lbServerGroupId string, memberId string) error {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.RemoveLbServerGroupMember(ctx, lbServerGroupId, memberId)

	_, err := req.Execute()
	return err
}

// ------------ LB Listener -------------------//
func (client *Client) GetLbListenerList(ctx context.Context, request LbListenerDataSource) (*loadbalancer.LbListenerListResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.ListLbListeners(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.LoadbalancerId.IsNull() {
		req = req.LoadbalancerId(request.LoadbalancerId.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.ServicePort.IsNull() {
		req = req.ServicePort(request.ServicePort.ValueInt32())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbListener(ctx context.Context, listenerId string) (*loadbalancer.LbListenerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.ShowLbListener(ctx, listenerId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateLbListener(ctx context.Context, request LbListenerResource) (*loadbalancer.LbListenerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.CreateLbListener(ctx)

	lbListener := request.LbListenerCreate

	var sslCertificate *loadbalancer.SslCertificate

	if lbListener.SslCertificate != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ClientCertId:    *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertId.ValueStringPointer()),
			ClientCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertLevel.ValueStringPointer()),
			ServerCertId:    *loadbalancer.NewNullableString(lbListener.SslCertificate.ServerCertId.ValueStringPointer()),
			ServerCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ServerCertLevel.ValueStringPointer()),
		}
	}

	urlHandlerInterfaces := make([]interface{}, len(lbListener.UrlHandler))
	for i, urlHandler := range lbListener.UrlHandler {
		urlHandlerInterfaces[i] = struct {
			UrlPattern    loadbalancer.NullableString `json:"url_pattern"`
			ServerGroupId loadbalancer.NullableString `json:"server_group_id"`
		}{
			*loadbalancer.NewNullableString(urlHandler.UrlPattern.ValueStringPointer()),
			*loadbalancer.NewNullableString(urlHandler.ServerGroupId.ValueStringPointer()),
		}
	}

	urlRedirectionInterfaces := make([]interface{}, len(lbListener.UrlRedirection))
	for i, urlRedirection := range lbListener.UrlRedirection {
		urlRedirectionInterfaces[i] = struct {
			UrlPattern         loadbalancer.NullableString `json:"url_pattern"`
			RedirectUrlPattern loadbalancer.NullableString `json:"redirect_url_pattern"`
		}{
			*loadbalancer.NewNullableString(urlRedirection.UrlPattern.ValueStringPointer()),
			*loadbalancer.NewNullableString(urlRedirection.RedirectUrlPattern.ValueStringPointer()),
		}
	}

	lbListenerElement := loadbalancer.LbListenerCreateRequest{
		Listener: loadbalancer.ListenerForCreate{
			Description:         *loadbalancer.NewNullableString(lbListener.Description.ValueStringPointer()),
			HttpsRedirection:    *loadbalancer.NewNullableBool(lbListener.HttpsRedirection.ValueBoolPointer()),
			InsertClientIp:      *loadbalancer.NewNullableBool(lbListener.InsertClientIp.ValueBoolPointer()),
			LoadbalancerId:      lbListener.LoadbalancerId.ValueString(),
			Name:                lbListener.Name.ValueString(),
			Persistence:         *loadbalancer.NewNullableString(lbListener.Persistence.ValueStringPointer()),
			Protocol:            lbListener.Protocol.ValueString(),
			ResponseTimeout:     *loadbalancer.NewNullableInt32(lbListener.ResponseTimeout.ValueInt32Pointer()),
			ServerGroupId:       *loadbalancer.NewNullableString(lbListener.ServerGroupId.ValueStringPointer()),
			ServicePort:         lbListener.ServicePort.ValueInt32(),
			SessionDurationTime: lbListener.SessionDurationTime.ValueInt32(),
			SslCertificate:      *loadbalancer.NewNullableSslCertificate(sslCertificate),
			UrlHandler:          urlHandlerInterfaces,
			UrlRedirection:      urlRedirectionInterfaces,
			XForwardedFor:       *loadbalancer.NewNullableBool(lbListener.XForwardedFor.ValueBoolPointer()),
			XForwardedPort:      *loadbalancer.NewNullableBool(lbListener.XForwardedPort.ValueBoolPointer()),
			XForwardedProto:     *loadbalancer.NewNullableBool(lbListener.XForwardedProto.ValueBoolPointer()),
		},
	}

	req = req.LbListenerCreateRequest(lbListenerElement)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLbListener(ctx context.Context, lbListenerId string, request LbListenerResource) (*loadbalancer.LbListenerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.SetLbListener(ctx, lbListenerId)

	lbListener := request.LbListenerCreate

	var sslCertificate *loadbalancer.SslCertificate

	if lbListener.SslCertificate != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ClientCertId:    *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertId.ValueStringPointer()),
			ClientCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertLevel.ValueStringPointer()),
			ServerCertId:    *loadbalancer.NewNullableString(lbListener.SslCertificate.ServerCertId.ValueStringPointer()),
			ServerCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ServerCertLevel.ValueStringPointer()),
		}
	}

	urlHandlerInterfaces := make([]interface{}, len(lbListener.UrlHandler))
	for i, urlHandler := range lbListener.UrlHandler {
		urlHandlerInterfaces[i] = struct {
			UrlPattern    loadbalancer.NullableString `json:"url_pattern"`
			ServerGroupId loadbalancer.NullableString `json:"server_group_id"`
		}{
			*loadbalancer.NewNullableString(urlHandler.UrlPattern.ValueStringPointer()),
			*loadbalancer.NewNullableString(urlHandler.ServerGroupId.ValueStringPointer()),
		}
	}

	urlRedirectionInterfaces := make([]interface{}, len(lbListener.UrlRedirection))
	for i, urlRedirection := range lbListener.UrlRedirection {
		urlRedirectionInterfaces[i] = struct {
			UrlPattern         loadbalancer.NullableString `json:"url_pattern"`
			RedirectUrlPattern loadbalancer.NullableString `json:"redirect_url_pattern"`
		}{
			*loadbalancer.NewNullableString(urlRedirection.UrlPattern.ValueStringPointer()),
			*loadbalancer.NewNullableString(urlRedirection.RedirectUrlPattern.ValueStringPointer()),
		}
	}

	lbListenerElement := loadbalancer.LbListenerSetRequest{
		Listener: loadbalancer.ListenerForSet{
			Description:         *loadbalancer.NewNullableString(lbListener.Description.ValueStringPointer()),
			HttpsRedirection:    *loadbalancer.NewNullableBool(lbListener.HttpsRedirection.ValueBoolPointer()),
			InsertClientIp:      *loadbalancer.NewNullableBool(lbListener.InsertClientIp.ValueBoolPointer()),
			Persistence:         *loadbalancer.NewNullableString(lbListener.Persistence.ValueStringPointer()),
			ResponseTimeout:     *loadbalancer.NewNullableInt32(lbListener.ResponseTimeout.ValueInt32Pointer()),
			ServerGroupId:       *loadbalancer.NewNullableString(lbListener.ServerGroupId.ValueStringPointer()),
			SessionDurationTime: *loadbalancer.NewNullableInt32(lbListener.SessionDurationTime.ValueInt32Pointer()),
			SslCertificate:      *loadbalancer.NewNullableSslCertificate(sslCertificate),
			UrlHandler:          urlHandlerInterfaces,
			UrlRedirection:      urlRedirectionInterfaces,
			XForwardedFor:       *loadbalancer.NewNullableBool(lbListener.XForwardedFor.ValueBoolPointer()),
			XForwardedPort:      *loadbalancer.NewNullableBool(lbListener.XForwardedPort.ValueBoolPointer()),
			XForwardedProto:     *loadbalancer.NewNullableBool(lbListener.XForwardedProto.ValueBoolPointer()),
		},
	}

	if lbListenerElement.Listener.Description.Get() == nil {
		lbListenerElement.Listener.Description.Unset()
	}

	if lbListener.Protocol.ValueString() != "HTTP" {
		lbListenerElement.Listener.HttpsRedirection.Unset()
	}

	if lbListener.Protocol.ValueString() != "TCP" {
		lbListenerElement.Listener.InsertClientIp.Unset()
	}

	if lbListener.Protocol.ValueString() == "UDP" {
		lbListenerElement.Listener.Persistence.Unset()
	}

	if lbListenerElement.Listener.ResponseTimeout.Get() == nil {
		lbListenerElement.Listener.ResponseTimeout.Unset()
	}

	if lbListener.Protocol.ValueString() == "HTTP" || lbListener.Protocol.ValueString() == "HTTPS" {
		lbListenerElement.Listener.ServerGroupId.Unset()
	}

	if lbListenerElement.Listener.SessionDurationTime.Get() == nil {
		lbListenerElement.Listener.SessionDurationTime.Unset()
	}

	if lbListenerElement.Listener.SslCertificate.Get() == nil {
		lbListenerElement.Listener.SslCertificate.Unset()
	}

	if lbListener.UrlHandler == nil {
		lbListenerElement.Listener.UrlHandler = nil
	}

	if lbListener.UrlRedirection == nil {
		lbListenerElement.Listener.UrlRedirection = nil
	}

	if lbListener.Protocol.ValueString() == "TCP" || lbListener.Protocol.ValueString() == "UDP" {
		lbListenerElement.Listener.XForwardedFor.Unset()
		lbListenerElement.Listener.XForwardedPort.Unset()
		lbListenerElement.Listener.XForwardedProto.Unset()
	}

	test, err := json.Marshal(lbListenerElement)
	print(test)

	req = req.LbListenerSetRequest(lbListenerElement)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLbListener(ctx context.Context, lbListenerId string) error {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.DeleteLbListener(ctx, lbListenerId)

	_, err := req.Execute()
	return err
}

// ------------ LB Certificate -------------------//
func (client *Client) GetLbCertificateList(ctx context.Context, request LbCertificateDataSource) (*loadbalancer.LbCertificateListResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbCertificatesApiAPI.ListLoadbalancerCertificates(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbCertificate(ctx context.Context, certificateId string) (*loadbalancer.LbCertificateShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LbCertificatesApiAPI.ShowLoadbalancerCertificate(ctx, certificateId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

//------------ LB Server Group -------------------//

func (client *Client) GetLbHealthCheckList(ctx context.Context, request LbHealthCheckDataSource) (*loadbalancer.LbHealthCheckListResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.ListLbHealthChecks(ctx)
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.Protocol.IsNull() {
		req = req.Protocol(loadbalancer.Protocol{
			LbMonitorProtocol:        loadbalancer.LbMonitorProtocol(request.Protocol.String()).Ptr(),
			ArrayOfLbMonitorProtocol: nil,
		})
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbHealthCheck(ctx context.Context, lbHealthCheckId string) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.ShowLbHealthCheck(ctx, lbHealthCheckId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLbHealthCheck(ctx context.Context, request LbHealthCheckResource) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.CreateLbHealthCheck(ctx)

	lbHealthCheck := request.LbHealthCheckCreate

	lbHealthCheckElement := loadbalancer.LbHealthCheckCreate{
		Name:                lbHealthCheck.Name.ValueString(),
		VpcId:               lbHealthCheck.VpcId.ValueString(),
		SubnetId:            lbHealthCheck.SubnetId.ValueString(),
		Protocol:            loadbalancer.LbMonitorProtocol(lbHealthCheck.Protocol.ValueString()),
		HealthCheckPort:     lbHealthCheck.HealthCheckPort.ValueInt32Pointer(),
		HealthCheckInterval: lbHealthCheck.HealthCheckInterval.ValueInt32Pointer(),
		HealthCheckTimeout:  lbHealthCheck.HealthCheckTimeout.ValueInt32Pointer(),
		HealthCheckCount:    lbHealthCheck.HealthCheckCount.ValueInt32Pointer(),
		HttpMethod:          *loadbalancer.NewNullableLbMonitorHttpMethod((*loadbalancer.LbMonitorHttpMethod)(lbHealthCheck.HttpMethod.ValueStringPointer())),
		HealthCheckUrl:      *loadbalancer.NewNullableString(lbHealthCheck.HealthCheckUrl.ValueStringPointer()),
		ResponseCode:        *loadbalancer.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
		RequestData:         *loadbalancer.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
		Description:         *loadbalancer.NewNullableString(lbHealthCheck.Description.ValueStringPointer()),
		Tags:                convertToTags(lbHealthCheck.Tags.Elements()),
	}

	req = req.LbHealthCheckCreateRequest(loadbalancer.LbHealthCheckCreateRequest{
		LbHealthCheck: lbHealthCheckElement,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLbHealthCheck(ctx context.Context, lbHealthCheckId string, request LbHealthCheckResource) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.SetLbHealthCheck(ctx, lbHealthCheckId)

	lbHealthCheck := request.LbHealthCheckCreate

	lbHealthCheckElement := loadbalancer.LbHealthCheckSet{
		Protocol:            *loadbalancer.NewNullableLbMonitorProtocol((*loadbalancer.LbMonitorProtocol)(lbHealthCheck.Protocol.ValueStringPointer())),
		HealthCheckPort:     *loadbalancer.NewNullableInt32(lbHealthCheck.HealthCheckPort.ValueInt32Pointer()),
		HealthCheckInterval: *loadbalancer.NewNullableInt32(lbHealthCheck.HealthCheckInterval.ValueInt32Pointer()),
		HealthCheckTimeout:  *loadbalancer.NewNullableInt32(lbHealthCheck.HealthCheckTimeout.ValueInt32Pointer()),
		HealthCheckCount:    *loadbalancer.NewNullableInt32(lbHealthCheck.HealthCheckCount.ValueInt32Pointer()),
		HttpMethod:          *loadbalancer.NewNullableLbMonitorHttpMethod((*loadbalancer.LbMonitorHttpMethod)(lbHealthCheck.HttpMethod.ValueStringPointer())),
		HealthCheckUrl:      *loadbalancer.NewNullableString(lbHealthCheck.HealthCheckUrl.ValueStringPointer()),
		ResponseCode:        *loadbalancer.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
		RequestData:         *loadbalancer.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
		Description:         *loadbalancer.NewNullableString(lbHealthCheck.Description.ValueStringPointer()),
	}

	req = req.LbHealthCheckSetRequest(loadbalancer.LbHealthCheckSetRequest{
		LbHealthCheck: lbHealthCheckElement,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLbHealthCheck(ctx context.Context, lbHealthCheckId string) error {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.DeleteLbHealthCheck(ctx, lbHealthCheckId)
	_, err := req.Execute()
	return err
}
