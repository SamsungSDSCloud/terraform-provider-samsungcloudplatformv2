package loadbalancer

import (
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	loadbalancersdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/loadbalancer/1.3"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertResponse(resp *loadbalancersdk.LbListenerShowResponseV1Dot3) (loadbalancer.LbListenerDetail, int) {

	var sslCertificate *loadbalancer.SslCertificate

	var sslCertificateFromData = resp.Listener.SslCertificate.Get()
	if sslCertificateFromData != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ClientCertId:    types.StringValue(sslCertificateFromData.GetClientCertId()),
			ServerCertLevel: types.StringValue(sslCertificateFromData.GetServerCertLevel()),
			ClientCertLevel: types.StringValue(sslCertificateFromData.GetClientCertLevel()),
		}
	}

	var sniCertificateList []loadbalancer.SniCertificateDataSource
	if resp.Listener.SniCertificate != nil {
		sniCertificateList = make([]loadbalancer.SniCertificateDataSource, len(resp.Listener.SniCertificate))
		for i, sniCertificate := range resp.Listener.SniCertificate {
			sniCertificateList[i] = loadbalancer.SniCertificateDataSource{
				SniCertId:  types.StringValue(sniCertificate.GetSniCertId()),
				DomainName: types.StringValue(sniCertificate.GetDomainName()),
				NotAfterDt: ToNullableTimeString(sniCertificate.NotAfterDt),
			}
		}
	}

	var urlHandlers []loadbalancer.UrlHandler
	skippedUrlHandlers := 0

	for _, urlHandlerInterface := range resp.Listener.UrlHandler {
		urlHandlerMap, ok := urlHandlerInterface.(map[string]interface{})
		if !ok {
			skippedUrlHandlers++
			continue
		}
		urlPattern, okUrl := urlHandlerMap["url_pattern"].(string)
		serverGroupId, okSg := urlHandlerMap["server_group_id"].(string)
		seq, okSeq := urlHandlerMap["seq"].(float64)
		if !okUrl || !okSg || !okSeq {
			skippedUrlHandlers++
			continue
		}
		urlHandlers = append(urlHandlers, loadbalancer.UrlHandler{
			UrlPattern:    types.StringValue(urlPattern),
			ServerGroupId: types.StringValue(serverGroupId),
			Seq:           types.Int32Value(int32(seq)),
		})
	}

	var httpsRedirection *loadbalancer.HttpsRedirection

	var httpsRedirectionFromData = resp.Listener.HttpsRedirection.Get()
	if httpsRedirectionFromData != nil {
		httpsRedirection = &loadbalancer.HttpsRedirection{
			Protocol:     types.StringValue(httpsRedirectionFromData.GetProtocol()),
			Port:         types.StringValue(httpsRedirectionFromData.GetPort()),
			ResponseCode: types.StringValue(httpsRedirectionFromData.GetResponseCode()),
		}
	}

	rtn := loadbalancer.LbListenerDetail{
		Id:                  types.StringValue(resp.Listener.Id),
		ModifiedBy:          types.StringValue(resp.Listener.ModifiedBy),
		ModifiedAt:          types.StringValue(resp.Listener.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resp.Listener.CreatedBy),
		CreatedAt:           types.StringValue(resp.Listener.CreatedAt.Format(time.RFC3339)),
		Description:         virtualserverutil.ToNullableStringValue(resp.Listener.Description.Get()),
		InsertClientIp:      ToNullableBoolValue(resp.Listener.InsertClientIp.Get()),
		Name:                types.StringValue(resp.Listener.Name),
		Persistence:         virtualserverutil.ToNullableStringValue(resp.Listener.Persistence.Get()),
		Protocol:            types.StringValue(string(resp.Listener.Protocol)),
		ServerGroupId:       virtualserverutil.ToNullableStringValue(resp.Listener.ServerGroupId.Get()),
		ServerGroupName:     virtualserverutil.ToNullableStringValue(resp.Listener.ServerGroupName.Get()),
		ServicePort:         types.Int32Value(resp.Listener.ServicePort),
		ResponseTimeout:     ToNullableInt32Value(resp.Listener.ResponseTimeout.Get()),
		SessionDurationTime: ToNullableInt32Value(resp.Listener.SessionDurationTime.Get()),
		SslCertificate:      sslCertificate,
		SniCertificate:      sniCertificateList,
		State:               types.StringValue(resp.Listener.State),
		UrlHandler:          urlHandlers,
		HttpsRedirection:    httpsRedirection,
		UrlRedirection:      virtualserverutil.ToNullableStringValue(resp.Listener.UrlRedirection.Get()),
		XForwardedFor:       ToNullableBoolValue(resp.Listener.XForwardedFor.Get()),
		XForwardedPort:      ToNullableBoolValue(resp.Listener.XForwardedPort.Get()),
		XForwardedProto:     ToNullableBoolValue(resp.Listener.XForwardedProto.Get()),
		RoutingAction:       types.StringValue(string(resp.Listener.RoutingAction)),
		ConditionType:       virtualserverutil.ToNullableStringValue((*string)(resp.Listener.ConditionType.Get())),
		IdleTimeout:         ToNullableInt32Value(resp.Listener.IdleTimeout.Get()),
		HstsMaxAge:          ToNullableInt32Value(resp.Listener.HstsMaxAge.Get()),
	}
	return rtn, skippedUrlHandlers
}
