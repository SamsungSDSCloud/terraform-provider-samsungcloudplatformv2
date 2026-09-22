package loadbalancer

import (
	"time"

	loadbalancer "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancerv1d4"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/virtualserver"
	loadbalancersdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertResponse(resp *loadbalancersdk.LbListenerShowResponseV1Dot4) (loadbalancer.LbListenerDetail, int) {

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
				NotAfterDt: ToNullableTimeStringV1d4(sniCertificate.NotAfterDt),
			}
		}
	}

	var urlHandlers []loadbalancer.UrlHandler
	skippedUrlHandlers := 0

	for _, urlHandler := range resp.Listener.UrlHandler {

		urlPattern := urlHandler.GetUrlPattern()
		serverGroupId := urlHandler.GetServerGroupId()
		seq := urlHandler.GetSeq()
		if urlPattern == "" || !urlHandler.ServerGroupId.IsSet() || !urlHandler.Seq.IsSet() {
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

	var hstsConfig *loadbalancer.HstsConfig
	var hstsConfigFromData = resp.Listener.HstsConfig.Get()
	if hstsConfigFromData != nil {
		hstsConfig = &loadbalancer.HstsConfig{
			MaxAge:            types.Int32Value(hstsConfigFromData.GetMaxAge()),
			IncludeSubDomains: types.BoolValue(hstsConfigFromData.GetIncludeSubDomains()),
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
		HstsConfig:          hstsConfig,
	}
	return rtn, skippedUrlHandlers
}

// HasResourceChanged checks if any resource-related fields differ between state and plan.
func HasResourceChanged(state, plan *loadbalancer.LbListenerCreate) bool {
	if state == nil || plan == nil {
		return true
	}

	checks := []struct {
		name    string
		changed bool
	}{
		{"description", !plan.Description.Equal(state.Description)},
		{"insert_client_ip", !plan.InsertClientIp.Equal(state.InsertClientIp)},
		{"persistence", !plan.Persistence.Equal(state.Persistence)},
		{"response_timeout", !plan.ResponseTimeout.Equal(state.ResponseTimeout)},
		{"session_duration_time", !plan.SessionDurationTime.Equal(state.SessionDurationTime)},
		{"idle_timeout", !plan.IdleTimeout.Equal(state.IdleTimeout)},
		{"support_http2", !plan.SupportHttp2.Equal(state.SupportHttp2)},
		{"x_forwarded_for", !plan.XForwardedFor.Equal(state.XForwardedFor)},
		{"x_forwarded_port", !plan.XForwardedPort.Equal(state.XForwardedPort)},
		{"x_forwarded_proto", !plan.XForwardedProto.Equal(state.XForwardedProto)},
	}

	for _, f := range checks {
		if f.changed {
			return true
		}
	}

	return HstsConfigChanged(state.HstsConfig, plan.HstsConfig)
}

// HstsConfigChanged checks if hsts_config fields differ.
func HstsConfigChanged(state, plan *loadbalancer.HstsConfig) bool {
	if state == nil && plan == nil {
		return false
	}
	if (state == nil) != (plan == nil) {
		return true
	}

	return !plan.MaxAge.Equal(state.MaxAge) ||
		!plan.IncludeSubDomains.Equal(state.IncludeSubDomains)
}

// BuildResourceSetRequest converts LbListenerCreate to SDK LbListenerResourceSetRequest.
func BuildResourceSetRequest(create *loadbalancer.LbListenerCreate) *loadbalancersdk.LbListenerResourceSetRequest {
	distinctBase := loadbalancersdk.NewLbListenerDistinctBase()

	if !create.Description.IsNull() && create.Description.ValueString() != "" {
		distinctBase.SetDescription(create.Description.ValueString())
	}
	if !create.InsertClientIp.IsNull() {
		distinctBase.SetInsertClientIp(create.InsertClientIp.ValueBool())
	}
	if !create.Persistence.IsNull() && create.Persistence.ValueString() != "" {
		distinctBase.SetPersistence(create.Persistence.ValueString())
	}
	if !create.ResponseTimeout.IsNull() {
		distinctBase.SetResponseTimeout(create.ResponseTimeout.ValueInt32())
	}
	if !create.SessionDurationTime.IsNull() {
		distinctBase.SetSessionDurationTime(create.SessionDurationTime.ValueInt32())
	}
	if !create.IdleTimeout.IsNull() {
		distinctBase.SetIdleTimeout(create.IdleTimeout.ValueInt32())
	}
	if !create.SupportHttp2.IsNull() {
		distinctBase.SetSupportHttp2(create.SupportHttp2.ValueBool())
	}
	if !create.XForwardedFor.IsNull() {
		distinctBase.SetXForwardedFor(create.XForwardedFor.ValueBool())
	}
	if !create.XForwardedPort.IsNull() {
		distinctBase.SetXForwardedPort(create.XForwardedPort.ValueBool())
	}
	if !create.XForwardedProto.IsNull() {
		distinctBase.SetXForwardedProto(create.XForwardedProto.ValueBool())
	}

	if create.HstsConfig != nil {
		hstsConfig := loadbalancersdk.NewHstsConfig()
		if !create.HstsConfig.MaxAge.IsNull() {
			hstsConfig.SetMaxAge(create.HstsConfig.MaxAge.ValueInt32())
		}
		hstsConfig.SetIncludeSubDomains(create.HstsConfig.IncludeSubDomains.ValueBool())
		distinctBase.SetHstsConfig(*hstsConfig)
	}

	return loadbalancersdk.NewLbListenerResourceSetRequest(*distinctBase)
}

// HasRuleChanged checks if any rule-related fields differ between state and plan.
func HasRuleChanged(state, plan *loadbalancer.LbListenerCreate) bool {
	if state == nil || plan == nil {
		return true
	}

	checks := []struct {
		name    string
		changed bool
	}{
		{"condition_type", !plan.ConditionType.Equal(state.ConditionType)},
		{"server_group_id", !plan.ServerGroupId.Equal(state.ServerGroupId)},
		{"url_redirection", !plan.UrlRedirection.Equal(state.UrlRedirection)},
	}

	for _, f := range checks {
		if f.changed {
			return true
		}
	}

	if HttpsRedirectionChanged(state.HttpsRedirection, plan.HttpsRedirection) {
		return true
	}

	return UrlHandlerChanged(state.UrlHandler, plan.UrlHandler)
}

// HttpsRedirectionChanged checks if https_redirection fields differ.
func HttpsRedirectionChanged(state, plan *loadbalancer.HttpsRedirection) bool {
	if state == nil && plan == nil {
		return false
	}
	if (state == nil) != (plan == nil) {
		return true
	}

	return !plan.Protocol.Equal(state.Protocol) ||
		!plan.Port.Equal(state.Port) ||
		!plan.ResponseCode.Equal(state.ResponseCode)
}

// UrlHandlerChanged checks if url_handler list differs.
func UrlHandlerChanged(state, plan []loadbalancer.UrlHandler) bool {
	if len(state) != len(plan) {
		return true
	}
	for i := range state {
		if !state[i].UrlPattern.Equal(plan[i].UrlPattern) ||
			!state[i].ServerGroupId.Equal(plan[i].ServerGroupId) ||
			!state[i].Seq.Equal(plan[i].Seq) {
			return true
		}
	}
	return false
}

// BuildRuleSetRequest converts LbListenerCreate to SDK LbListenerRuleSetRequest.
func BuildRuleSetRequest(create *loadbalancer.LbListenerCreate) *loadbalancersdk.LbListenerRuleSetRequest {
	ruleBase := loadbalancersdk.NewLbListenerRuleBase()

	if !create.ConditionType.IsNull() && create.ConditionType.ValueString() != "" {
		condType := loadbalancersdk.ConditionType(create.ConditionType.ValueString())
		ruleBase.SetConditionType(condType)
	}
	if !create.ServerGroupId.IsNull() && create.ServerGroupId.ValueString() != "" {
		ruleBase.SetServerGroupId(create.ServerGroupId.ValueString())
	}
	if !create.UrlRedirection.IsNull() && create.UrlRedirection.ValueString() != "" {
		ruleBase.SetUrlRedirection(create.UrlRedirection.ValueString())
	}
	if create.HttpsRedirection != nil {
		httpsRedirect := loadbalancersdk.NewHttpsRedirection()
		if !create.HttpsRedirection.Protocol.IsNull() && create.HttpsRedirection.Protocol.ValueString() != "" {
			httpsRedirect.SetProtocol(create.HttpsRedirection.Protocol.ValueString())
		}
		if !create.HttpsRedirection.Port.IsNull() && create.HttpsRedirection.Port.ValueString() != "" {
			httpsRedirect.SetPort(create.HttpsRedirection.Port.ValueString())
		}
		if !create.HttpsRedirection.ResponseCode.IsNull() && create.HttpsRedirection.ResponseCode.ValueString() != "" {
			httpsRedirect.SetResponseCode(create.HttpsRedirection.ResponseCode.ValueString())
		}
		ruleBase.SetHttpsRedirection(*httpsRedirect)
	}
	if len(create.UrlHandler) > 0 {
		urlHandlers := make([]loadbalancersdk.UrlHandler, len(create.UrlHandler))
		for i, uh := range create.UrlHandler {
			urlHandlers[i] = *loadbalancersdk.NewUrlHandler(
				uh.Seq.ValueInt32(),
				uh.UrlPattern.ValueString(),
			)
			if !uh.ServerGroupId.IsNull() && uh.ServerGroupId.ValueString() != "" {
				urlHandlers[i].SetServerGroupId(uh.ServerGroupId.ValueString())
			}
		}
		ruleBase.SetUrlHandler(urlHandlers)
	}

	return loadbalancersdk.NewLbListenerRuleSetRequest(*ruleBase)
}

// SslCertificateChanged checks if ssl_certificate fields differ.
func SslCertificateChanged(state, plan *loadbalancer.SslCertificate) bool {
	if state == nil && plan == nil {
		return false
	}
	if (state == nil) != (plan == nil) {
		return true
	}

	return !plan.ClientCertId.Equal(state.ClientCertId) ||
		!plan.ClientCertLevel.Equal(state.ClientCertLevel) ||
		!plan.ServerCertLevel.Equal(state.ServerCertLevel)
}

// SniCertificateChanged checks if sni_certificate list differs.
func SniCertificateChanged(state, plan []loadbalancer.SniCertificate) bool {
	if len(state) != len(plan) {
		return true
	}
	for i := range state {
		if !state[i].SniCertId.Equal(plan[i].SniCertId) ||
			!state[i].DomainName.Equal(plan[i].DomainName) {
			return true
		}
	}
	return false
}

// BuildSslCertificateSetRequest converts LbListenerCreate to SDK LbListenerCertificateSetRequest with SSL cert.
func BuildSslCertificateSetRequest(create *loadbalancer.LbListenerCreate) *loadbalancersdk.LbListenerCertificateSetRequest {
	base := loadbalancersdk.LbListenerCertificateBase{}

	if create.SslCertificate != nil {
		sslCert := loadbalancersdk.SslCertificate{}
		if !create.SslCertificate.ClientCertId.IsNull() {
			sslCert.SetClientCertId(create.SslCertificate.ClientCertId.ValueString())
		}
		if !create.SslCertificate.ClientCertLevel.IsNull() {
			sslCert.SetClientCertLevel(create.SslCertificate.ClientCertLevel.ValueString())
		}
		if !create.SslCertificate.ServerCertLevel.IsNull() {
			sslCert.SetServerCertLevel(create.SslCertificate.ServerCertLevel.ValueString())
		}
		base.SslCertificate = *loadbalancersdk.NewNullableSslCertificate(&sslCert)
	}

	return loadbalancersdk.NewLbListenerCertificateSetRequest(base)
}

// BuildSniCertificateSetRequest converts LbListenerCreate to SDK LbListenerCertificateSetRequest with SNI certs.
func BuildSniCertificateSetRequest(create *loadbalancer.LbListenerCreate) *loadbalancersdk.LbListenerCertificateSetRequest {
	base := loadbalancersdk.LbListenerCertificateBase{}

	if len(create.SniCertificate) > 0 {
		sniCerts := make([]loadbalancersdk.SniCertificate, len(create.SniCertificate))
		for i, sni := range create.SniCertificate {
			sniCert := loadbalancersdk.SniCertificate{}
			if !sni.SniCertId.IsNull() {
				sniCert.SetSniCertId(sni.SniCertId.ValueString())
			}
			if !sni.DomainName.IsNull() {
				sniCert.SetDomainName(sni.DomainName.ValueString())
			}
			sniCerts[i] = sniCert
		}
		base.SniCertificate = sniCerts
	}

	return loadbalancersdk.NewLbListenerCertificateSetRequest(base)
}
