package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LbListenerResource struct {
	Id               types.String      `tfsdk:"id"`
	LbListener       types.Object      `tfsdk:"lb_listener"`
	LbListenerCreate *LbListenerCreate `tfsdk:"lb_listener_create"`
}

type LbListenerCreate struct {
	Description         types.String      `tfsdk:"description"`
	InsertClientIp      types.Bool        `tfsdk:"insert_client_ip"`
	LoadbalancerId      types.String      `tfsdk:"loadbalancer_id"`
	Name                types.String      `tfsdk:"name"`
	Persistence         types.String      `tfsdk:"persistence"`
	Protocol            types.String      `tfsdk:"protocol"`
	ResponseTimeout     types.Int32       `tfsdk:"response_timeout"`
	ServerGroupId       types.String      `tfsdk:"server_group_id"`
	ServicePort         types.Int32       `tfsdk:"service_port"`
	SessionDurationTime types.Int32       `tfsdk:"session_duration_time"`
	SslCertificate      *SslCertificate   `tfsdk:"ssl_certificate"`
	SniCertificate      []SniCertificate  `tfsdk:"sni_certificate"`
	UrlHandler          []UrlHandler      `tfsdk:"url_handler"`
	UrlRedirection      types.String      `tfsdk:"url_redirection"`
	HttpsRedirection    *HttpsRedirection `tfsdk:"https_redirection"`
	XForwardedFor       types.Bool        `tfsdk:"x_forwarded_for"`
	XForwardedPort      types.Bool        `tfsdk:"x_forwarded_port"`
	XForwardedProto     types.Bool        `tfsdk:"x_forwarded_proto"`
	RoutingAction       types.String      `tfsdk:"routing_action"`
	ConditionType       types.String      `tfsdk:"condition_type"`
	IdleTimeout         types.Int32       `tfsdk:"idle_timeout"`
	SupportHttp2        types.Bool        `tfsdk:"support_http2"`
	HstsConfig          *HstsConfig       `tfsdk:"hsts_config"`
	Tags                types.Map         `tfsdk:"tags"`
}

type LbListenerDataSource struct {
	Size           types.Int32  `tfsdk:"size"`
	Page           types.Int32  `tfsdk:"page"`
	Sort           types.String `tfsdk:"sort"`
	LoadbalancerId types.String `tfsdk:"loadbalancer_id"`
	State          types.String `tfsdk:"state"`
	Name           types.String `tfsdk:"name"`
	ServicePort    types.Int32  `tfsdk:"service_port"`
	LbListeners    []LbListener `tfsdk:"lb_listeners"`
}

type LbListener struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Protocol    types.String `tfsdk:"protocol"`
	State       types.String `tfsdk:"state"`
	ServicePort types.Int32  `tfsdk:"service_port"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
}

type LbListenerDataSourceDetail struct {
	Id               types.String `tfsdk:"id"`
	LbListenerDetail types.Object `tfsdk:"lb_listener"`
}

type LbListenerDetail struct {
	Id                  types.String               `tfsdk:"id"`
	CreatedAt           types.String               `tfsdk:"created_at"`
	CreatedBy           types.String               `tfsdk:"created_by"`
	ModifiedAt          types.String               `tfsdk:"modified_at"`
	ModifiedBy          types.String               `tfsdk:"modified_by"`
	Description         types.String               `tfsdk:"description"`
	InsertClientIp      types.Bool                 `tfsdk:"insert_client_ip"`
	Name                types.String               `tfsdk:"name"`
	Persistence         types.String               `tfsdk:"persistence"`
	Protocol            types.String               `tfsdk:"protocol"`
	ServerGroupId       types.String               `tfsdk:"server_group_id"`
	ServerGroupName     types.String               `tfsdk:"server_group_name"`
	ServicePort         types.Int32                `tfsdk:"service_port"`
	ResponseTimeout     types.Int32                `tfsdk:"response_timeout"`
	SessionDurationTime types.Int32                `tfsdk:"session_duration_time"`
	SslCertificate      *SslCertificate            `tfsdk:"ssl_certificate"`
	SniCertificate      []SniCertificateDataSource `tfsdk:"sni_certificate"`
	State               types.String               `tfsdk:"state"`
	UrlHandler          []UrlHandler               `tfsdk:"url_handler"`
	HttpsRedirection    *HttpsRedirection          `tfsdk:"https_redirection"`
	UrlRedirection      types.String               `tfsdk:"url_redirection"`
	XForwardedFor       types.Bool                 `tfsdk:"x_forwarded_for"`
	XForwardedPort      types.Bool                 `tfsdk:"x_forwarded_port"`
	XForwardedProto     types.Bool                 `tfsdk:"x_forwarded_proto"`
	RoutingAction       types.String               `tfsdk:"routing_action"`
	ConditionType       types.String               `tfsdk:"condition_type"`
	IdleTimeout         types.Int32                `tfsdk:"idle_timeout"`
	HstsConfig          *HstsConfig                `tfsdk:"hsts_config"`
}

func (m LbListenerDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"description":           types.StringType,
		"insert_client_ip":      types.BoolType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
		"name":                  types.StringType,
		"persistence":           types.StringType,
		"protocol":              types.StringType,
		"server_group_id":       types.StringType,
		"server_group_name":     types.StringType,
		"service_port":          types.Int32Type,
		"response_timeout":      types.Int32Type,
		"session_duration_time": types.Int32Type,
		"ssl_certificate": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"client_cert_id":    types.StringType,
				"client_cert_level": types.StringType,
				"server_cert_level": types.StringType,
			},
		},
		"sni_certificate": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"sni_cert_id":  types.StringType,
					"domain_name":  types.StringType,
					"not_after_dt": types.StringType,
				},
			},
		},
		"state": types.StringType,
		"url_handler": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"url_pattern":     types.StringType,
					"server_group_id": types.StringType,
					"seq":             types.Int32Type,
				},
			},
		},
		"https_redirection": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"protocol":      types.StringType,
				"port":          types.StringType,
				"response_code": types.StringType,
			},
		},
		"url_redirection":   types.StringType,
		"x_forwarded_for":   types.BoolType,
		"x_forwarded_port":  types.BoolType,
		"x_forwarded_proto": types.BoolType,
		"routing_action":    types.StringType,
		"condition_type":    types.StringType,
		"idle_timeout":      types.Int32Type,
		"hsts_config": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"max_age":             types.Int32Type,
				"include_sub_domains": types.BoolType,
			},
		},
	}
}

type SslCertificate struct {
	ClientCertId    types.String `tfsdk:"client_cert_id"`
	ClientCertLevel types.String `tfsdk:"client_cert_level"`
	ServerCertLevel types.String `tfsdk:"server_cert_level"`
}

type UrlHandler struct {
	UrlPattern    types.String `tfsdk:"url_pattern"`
	ServerGroupId types.String `tfsdk:"server_group_id"`
	Seq           types.Int32  `tfsdk:"seq"`
}

type HttpsRedirection struct {
	Protocol     types.String `tfsdk:"protocol"`
	Port         types.String `tfsdk:"port"`
	ResponseCode types.String `tfsdk:"response_code"`
}

type SniCertificateDataSource struct {
	SniCertId  types.String `tfsdk:"sni_cert_id"`
	DomainName types.String `tfsdk:"domain_name"`
	NotAfterDt types.String `tfsdk:"not_after_dt"`
}

type SniCertificate struct {
	SniCertId  types.String `tfsdk:"sni_cert_id"`
	DomainName types.String `tfsdk:"domain_name"`
	NotAfterDt types.String `tfsdk:"not_after_dt"`
}

type HstsConfig struct {
	MaxAge            types.Int32 `tfsdk:"max_age"`
	IncludeSubDomains types.Bool  `tfsdk:"include_sub_domains"`
}
