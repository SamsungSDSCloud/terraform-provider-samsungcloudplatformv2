package loadbalancerv1d4

import (
	"context"

	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
)

func (client *Client) GetLbListener(ctx context.Context, listenerId string) (*loadbalancer.LbListenerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.ShowLbListener(ctx, listenerId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) DeleteLbListener(ctx context.Context, lbListenerId string) error {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.DeleteLbListener(ctx, lbListenerId)

	_, err := req.Execute()
	return err
}

func (client *Client) CreateLbListener(ctx context.Context, request LbListenerResource) (*loadbalancer.LbListenerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.CreateLbListener(ctx)

	lbListener := request.LbListenerCreate

	var sslCertificate *loadbalancer.SslCertificate
	if lbListener.SslCertificate != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ClientCertId:    *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertId.ValueStringPointer()),
			ClientCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ClientCertLevel.ValueStringPointer()),
			ServerCertLevel: *loadbalancer.NewNullableString(lbListener.SslCertificate.ServerCertLevel.ValueStringPointer()),
		}
	}

	var httpsRedirection *loadbalancer.HttpsRedirection

	if lbListener.HttpsRedirection != nil {
		httpsRedirection = &loadbalancer.HttpsRedirection{
			Protocol:     *loadbalancer.NewNullableString(lbListener.HttpsRedirection.Protocol.ValueStringPointer()),
			Port:         *loadbalancer.NewNullableString(lbListener.HttpsRedirection.Port.ValueStringPointer()),
			ResponseCode: *loadbalancer.NewNullableString(lbListener.HttpsRedirection.ResponseCode.ValueStringPointer()),
		}
	}
	var urlHandlers []loadbalancer.UrlHandler

	if lbListener.UrlHandler != nil {
		for _, urlHandler := range lbListener.UrlHandler {

			urlHandlers = append(urlHandlers, loadbalancer.UrlHandler{
				UrlPattern:    urlHandler.UrlPattern.ValueString(),
				ServerGroupId: *loadbalancer.NewNullableString(urlHandler.ServerGroupId.ValueStringPointer()),
				Seq:           urlHandler.Seq.ValueInt32(),
			})
		}
	} else {
		urlHandlers = nil
	}

	tags := convertToTags(lbListener.Tags.Elements())

	lbListenerElement := loadbalancer.LbListenerCreateRequestV1Dot4{
		Listener: loadbalancer.ListenerForCreateV1Dot4{
			Description:         *loadbalancer.NewNullableString(lbListener.Description.ValueStringPointer()),
			InsertClientIp:      *loadbalancer.NewNullableBool(lbListener.InsertClientIp.ValueBoolPointer()),
			LoadbalancerId:      lbListener.LoadbalancerId.ValueString(),
			Name:                lbListener.Name.ValueString(),
			Persistence:         *loadbalancer.NewNullableString(lbListener.Persistence.ValueStringPointer()),
			Protocol:            loadbalancer.LbListenerProtocol(lbListener.Protocol.ValueString()),
			ResponseTimeout:     *loadbalancer.NewNullableInt32(lbListener.ResponseTimeout.ValueInt32Pointer()),
			ServerGroupId:       *loadbalancer.NewNullableString(lbListener.ServerGroupId.ValueStringPointer()),
			ServicePort:         lbListener.ServicePort.ValueInt32(),
			SessionDurationTime: *loadbalancer.NewNullableInt32(lbListener.SessionDurationTime.ValueInt32Pointer()),
			SslCertificate:      *loadbalancer.NewNullableSslCertificate(sslCertificate),
			UrlHandler:          urlHandlers,
			HttpsRedirection:    *loadbalancer.NewNullableHttpsRedirection(httpsRedirection),
			UrlRedirection:      *loadbalancer.NewNullableString(lbListener.UrlRedirection.ValueStringPointer()),
			XForwardedFor:       *loadbalancer.NewNullableBool(lbListener.XForwardedFor.ValueBoolPointer()),
			XForwardedPort:      *loadbalancer.NewNullableBool(lbListener.XForwardedPort.ValueBoolPointer()),
			XForwardedProto:     *loadbalancer.NewNullableBool(lbListener.XForwardedProto.ValueBoolPointer()),
			RoutingAction:       loadbalancer.RoutingAction(lbListener.RoutingAction.ValueString()),
			ConditionType:       *loadbalancer.NewNullableConditionType((*loadbalancer.ConditionType)(lbListener.ConditionType.ValueStringPointer())),
			IdleTimeout:         *loadbalancer.NewNullableInt32(lbListener.IdleTimeout.ValueInt32Pointer()),
			SupportHttp2:        *loadbalancer.NewNullableBool(lbListener.SupportHttp2.ValueBoolPointer()),
			Tags:                tags,
		},
	}

	req = req.LbListenerCreateRequestV1Dot4(lbListenerElement)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbServerGroup(ctx context.Context, lbServerGroupId string) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.ShowLbServerGroup(ctx, lbServerGroupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetLbListenerResource(ctx context.Context, listenerId string, request *loadbalancer.LbListenerResourceSetRequest) (*loadbalancer.LbListenerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.SetLbListenerResource(ctx, listenerId)
	req = req.LbListenerResourceSetRequest(*request)
	resp, _, err := req.Execute()
	return resp, err
}

// SetLbListenerCertificate updates the SSL/SNI certificate for a listener.
// PUT /v1/lb-listeners/{listener_id}/certificate
func (client *Client) SetLbListenerCertificate(ctx context.Context, listenerId string,
	body *loadbalancer.LbListenerCertificateSetRequest) (*loadbalancer.LbListenerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.SetLbListenerCertificate(ctx, listenerId)
	req = req.LbListenerCertificateSetRequest(*body)
	resp, _, err := req.Execute()
	return resp, err
}
