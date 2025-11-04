package loadbalancer

import (
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	loadbalancersdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/loadbalancer/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertResponse(resp *loadbalancersdk.LbListenerShowResponse) loadbalancer.LbListenerDetail {

	var sslCertificate *loadbalancer.SslCertificate

	var sslCertificateFromData = resp.Listener.SslCertificate.Get()
	if sslCertificateFromData != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ClientCertId:    types.StringValue(sslCertificateFromData.GetClientCertId()),
			ServerCertLevel: types.StringValue(sslCertificateFromData.GetServerCertLevel()),
			ClientCertLevel: types.StringValue(sslCertificateFromData.GetClientCertLevel()),
		}
	}

	sniCertificateList := make([]loadbalancer.SniCertificate, len(resp.Listener.SniCertificate))
	for i, sniCertificate := range resp.Listener.SniCertificate {
		sniCertificateList[i] = loadbalancer.SniCertificate{
			SniCertId:  types.StringValue(sniCertificate.GetSniCertId()),
			DomainName: types.StringValue(sniCertificate.GetDomainName()),
		}
	}

	var urlHandlers []loadbalancer.UrlHandler

	for _, urlHandlerInterface := range resp.Listener.UrlHandler {
		urlHandlerMap, _ := urlHandlerInterface.(map[string]interface{})
		urlPattern, _ := urlHandlerMap["url_pattern"].(string)
		serverGroupId, _ := urlHandlerMap["server_group_id"].(string)
		seq, _ := urlHandlerMap["seq"].(float64)
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
		InsertClientIp:      types.BoolValue(resp.Listener.InsertClientIp.IsSet()),
		Name:                types.StringValue(resp.Listener.Name),
		Persistence:         virtualserverutil.ToNullableStringValue(resp.Listener.Persistence.Get()),
		Protocol:            types.StringValue(resp.Listener.Protocol),
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
		RoutingAction:       virtualserverutil.ToNullableStringValue((*string)(resp.Listener.RoutingAction.Get())),
		ConditionType:       virtualserverutil.ToNullableStringValue((*string)(resp.Listener.ConditionType.Get())),
	}
	return rtn
}
