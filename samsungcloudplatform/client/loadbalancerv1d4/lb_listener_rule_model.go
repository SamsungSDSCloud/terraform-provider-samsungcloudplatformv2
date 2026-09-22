package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HttpsRedirectionAttrTypes is the attribute type map for HTTPS redirection objects.
var HttpsRedirectionAttrTypes = map[string]attr.Type{
	"protocol":      types.StringType,
	"port":          types.StringType,
	"response_code": types.StringType,
}

// UrlHandlerAttrTypes is the attribute type map for URL handler objects.
var UrlHandlerAttrTypes = map[string]attr.Type{
	"url_pattern":     types.StringType,
	"server_group_id": types.StringType,
	"seq":             types.Int32Type,
}

// LbListenerRuleResource is the top-level Terraform resource model for managing
// listener routing rule configuration.
type LbListenerRuleResource struct {
	ListenerId   types.String  `tfsdk:"listener_id"`
	ListenerRule *ListenerRule `tfsdk:"listener_rule"`
	LbListener   types.Object  `tfsdk:"lb_listener"`
}

// ListenerRule contains the routing rule configuration fields.
// Uses types.Object and types.List to properly handle unknown values during plan.
type ListenerRule struct {
	ConditionType    types.String `tfsdk:"condition_type"`
	HttpsRedirection types.Object `tfsdk:"https_redirection"`
	ServerGroupId    types.String `tfsdk:"server_group_id"`
	UrlHandler       types.List   `tfsdk:"url_handler"`
	UrlRedirection   types.String `tfsdk:"url_redirection"`
}

// HttpsRedirectionInputModel is the Go struct used to extract values from types.Object.
type HttpsRedirectionInputModel struct {
	Protocol     types.String `tfsdk:"protocol"`
	Port         types.String `tfsdk:"port"`
	ResponseCode types.String `tfsdk:"response_code"`
}

// UrlHandlerInputModel is the Go struct used to extract values from URL handler objects.
type UrlHandlerInputModel struct {
	UrlPattern    types.String `tfsdk:"url_pattern"`
	ServerGroupId types.String `tfsdk:"server_group_id"`
	Seq           types.Int32  `tfsdk:"seq"`
}
