package loadbalancer

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/virtualserver"
	loadbalancersdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/loadbalancer/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

func ConvertResponse(resp *loadbalancersdk.LbListenerShowResponse) loadbalancer.LbListenerDetail {

	var sslCertificate *loadbalancer.SslCertificate

	var sslCertificateFromData = resp.Listener.SslCertificate.Get()
	if sslCertificateFromData != nil {
		sslCertificate = &loadbalancer.SslCertificate{
			ServerCertId:    types.StringValue(sslCertificateFromData.GetServerCertId()),
			ClientCertId:    types.StringValue(sslCertificateFromData.GetClientCertId()),
			ServerCertLevel: types.StringValue(sslCertificateFromData.GetServerCertLevel()),
			ClientCertLevel: types.StringValue(sslCertificateFromData.GetClientCertLevel()),
		}
	}

	var urlHandlers []loadbalancer.UrlHandler
	var urlRedirections []loadbalancer.UrlRedirection

	for _, urlHandlerInterface := range resp.Listener.UrlHandler {
		urlHandlerMap, _ := urlHandlerInterface.(map[string]interface{})
		urlPattern, _ := urlHandlerMap["url_pattern"].(string)
		serverGroupId, _ := urlHandlerMap["server_group_id"].(string)
		urlHandlers = append(urlHandlers, loadbalancer.UrlHandler{
			UrlPattern:    types.StringValue(urlPattern),
			ServerGroupId: types.StringValue(serverGroupId),
		})
	}

	for _, urlRedirectionInterface := range resp.Listener.UrlRedirection {
		urlRedirectionMap, _ := urlRedirectionInterface.(map[string]interface{})
		urlPattern, _ := urlRedirectionMap["url_pattern"].(string)
		redirectUrlPattern, _ := urlRedirectionMap["redirect_url_pattern"].(string)
		urlRedirections = append(urlRedirections, loadbalancer.UrlRedirection{
			UrlPattern:         types.StringValue(urlPattern),
			RedirectUrlPattern: types.StringValue(redirectUrlPattern),
		})
	}
	rtn := loadbalancer.LbListenerDetail{
		Id:                  types.StringValue(resp.Listener.Id),
		ModifiedBy:          types.StringValue(resp.Listener.ModifiedBy),
		ModifiedAt:          types.StringValue(resp.Listener.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resp.Listener.CreatedBy),
		CreatedAt:           types.StringValue(resp.Listener.CreatedAt.Format(time.RFC3339)),
		Description:         virtualserverutil.ToNullableStringValue(resp.Listener.Description.Get()),
		HttpsRedirection:    types.BoolValue(resp.Listener.HttpsRedirection.IsSet()),
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
		State:               types.StringValue(resp.Listener.State),
		UrlHandler:          urlHandlers,
		UrlRedirection:      urlRedirections,
		XForwardedFor:       ToNullableBoolValue(resp.Listener.XForwardedFor.Get()),
		XForwardedPort:      ToNullableBoolValue(resp.Listener.XForwardedPort.Get()),
		XForwardedProto:     ToNullableBoolValue(resp.Listener.XForwardedProto.Get()),
	}
	return rtn
}
