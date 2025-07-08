package loadbalancer

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-loadbalancer"

// ------------ Load Balancer -------------------//
type LoadbalancerDataSource struct {
	Size          types.Int32    `tfsdk:"size"`
	Page          types.Int32    `tfsdk:"page"`
	Sort          types.String   `tfsdk:"sort"`
	Name          types.String   `tfsdk:"name"`
	ServiceIp     types.String   `tfsdk:"service_ip"`
	SubnetId      types.String   `tfsdk:"subnet_id"`
	Loadbalancers []Loadbalancer `tfsdk:"loadbalancers"`
}

// list response
type Loadbalancer struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	PublicNatIp   types.String `tfsdk:"public_nat_ip"`
	ServiceIp     types.String `tfsdk:"service_ip"`
	SourceNatIp   types.String `tfsdk:"source_nat_ip"`
	State         types.String `tfsdk:"state"`
	ListenerCount types.Int32  `tfsdk:"listener_count"`
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
}

type LoadbalancerDataSourceDetail struct {
	Id                 types.String `tfsdk:"id"`
	LoadbalancerDetail types.Object `tfsdk:"loadbalancer"`
}

type LoadbalancerDetail struct {
	AccountId      types.String `tfsdk:"account_id"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Description    types.String `tfsdk:"description"`
	FirewallId     types.String `tfsdk:"firewall_id"`
	FirewallName   types.String `tfsdk:"firewall_name"`
	HealthCheckIp  types.List   `tfsdk:"health_check_ip"`
	Id             types.String `tfsdk:"id"`
	LayerType      types.String `tfsdk:"layer_type"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	Name           types.String `tfsdk:"name"`
	PublicNatIp    types.String `tfsdk:"public_nat_ip"`
	PublicNatState types.String `tfsdk:"public_nat_state"`
	ServiceIp      types.String `tfsdk:"service_ip"`
	SourceNatIp    types.String `tfsdk:"source_nat_ip"`
	State          types.String `tfsdk:"state"`
	SubnetId       types.String `tfsdk:"subnet_id"`
	VpcId          types.String `tfsdk:"vpc_id"`
}

func (m LoadbalancerDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":    types.StringType,
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"description":   types.StringType,
		"firewall_id":   types.StringType,
		"firewall_name": types.StringType,
		"health_check_ip": types.ListType{
			ElemType: types.StringType,
		},
		"id":               types.StringType,
		"layer_type":       types.StringType,
		"modified_at":      types.StringType,
		"modified_by":      types.StringType,
		"name":             types.StringType,
		"public_nat_ip":    types.StringType,
		"public_nat_state": types.StringType,
		"service_ip":       types.StringType,
		"source_nat_ip":    types.StringType,
		"state":            types.StringType,
		"subnet_id":        types.StringType,
		"vpc_id":           types.StringType,
	}
}

type LoadbalancerResource struct {
	Id                 types.String       `tfsdk:"id"`
	Loadbalancer       types.Object       `tfsdk:"loadbalancer"`
	LoadbalancerCreate LoadbalancerCreate `tfsdk:"loadbalancer_create"`
}

type LoadbalancerCreate struct {
	Description            types.String `tfsdk:"description"`
	FirewallEnabled        types.Bool   `tfsdk:"firewall_enabled"`
	FirewallLoggingEnabled types.Bool   `tfsdk:"firewall_logging_enabled"`
	LayerType              types.String `tfsdk:"layer_type"`
	Name                   types.String `tfsdk:"name"`
	PublicipId             types.String `tfsdk:"publicip_id"`
	ServiceIp              types.String `tfsdk:"service_ip"`
	SubnetId               types.String `tfsdk:"subnet_id"`
	VpcId                  types.String `tfsdk:"vpc_id"`
}

type LoadbalancerPublicNatIpResource struct {
	LoadbalancerId          types.String    `tfsdk:"loadbalancer_id"`
	Id                      types.String    `tfsdk:"id"`
	LoadbalancerPublicNatIp types.Object    `tfsdk:"loadbalancer_public_nat_ip"`
	LoadbalancerNatCreate   StaticNatCreate `tfsdk:"static_nat_create"`
}

type StaticNatCreate struct {
	PublicipId types.String `tfsdk:"publicip_id"`
}

type LoadbalancerPublicNatIpDetail struct {
	AccountId         types.String `tfsdk:"account_id"`
	ActionType        types.String `tfsdk:"action_type"`
	CreatedAt         types.String `tfsdk:"created_at"`
	CreatedBy         types.String `tfsdk:"created_by"`
	Description       types.String `tfsdk:"description"`
	ExternalIpAddress types.String `tfsdk:"external_ip_address"`
	Id                types.String `tfsdk:"id"`
	InternalIpAddress types.String `tfsdk:"internal_ip_address"`
	ModifiedAt        types.String `tfsdk:"modified_at"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
	Name              types.String `tfsdk:"name"`
	OwnerId           types.String `tfsdk:"owner_id"`
	OwnerName         types.String `tfsdk:"owner_name"`
	OwnerType         types.String `tfsdk:"owner_type"`
	PublicipId        types.String `tfsdk:"publicip_id"`
	ServiceIpPortId   types.String `tfsdk:"service_ip_port_id"`
	State             types.String `tfsdk:"state"`
	SubnetId          types.String `tfsdk:"subnet_id"`
	Type              types.String `tfsdk:"type"`
	VpcId             types.String `tfsdk:"vpc_id"`
}

func (m LoadbalancerPublicNatIpDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":          types.StringType,
		"action_type":         types.StringType,
		"created_at":          types.StringType,
		"created_by":          types.StringType,
		"description":         types.StringType,
		"external_ip_address": types.StringType,
		"id":                  types.StringType,
		"internal_ip_address": types.StringType,
		"modified_at":         types.StringType,
		"modified_by":         types.StringType,
		"name":                types.StringType,
		"owner_id":            types.StringType,
		"owner_name":          types.StringType,
		"owner_type":          types.StringType,
		"publicip_id":         types.StringType,
		"service_ip_port_id":  types.StringType,
		"state":               types.StringType,
		"subnet_id":           types.StringType,
		"type":                types.StringType,
		"vpc_id":              types.StringType,
	}
}

//------------ LB Server Group -------------------//

// List LB ServerGroup Paramaters
type LbServerGroupDataSource struct {
	Size           types.Int32     `tfsdk:"size"`
	Page           types.Int32     `tfsdk:"page"`
	Sort           types.String    `tfsdk:"sort"`
	Name           types.String    `tfsdk:"name"`
	Protocol       types.List      `tfsdk:"protocol"`
	VpcId          types.String    `tfsdk:"vpc_id"`
	SubnetId       types.String    `tfsdk:"subnet_id"`
	LbServerGroups []LbServerGroup `tfsdk:"lb_server_groups"`
}

type LbServerGroupDataSourceDetail struct {
	Id                  types.String `tfsdk:"id"`
	LbServerGroupDetail types.Object `tfsdk:"lb_server_group"`
}

// Get LB ServerGroup Response
type LbServerGroupDetail struct {
	Name            types.String `tfsdk:"name"`
	Protocol        types.String `tfsdk:"protocol"`
	LoadbalancerId  types.String `tfsdk:"loadbalancer_id"`
	LbName          types.String `tfsdk:"lb_name"`
	LbMethod        types.String `tfsdk:"lb_method"`
	LbHealthCheckId types.String `tfsdk:"lb_health_check_id"`
	State           types.String `tfsdk:"state"`
	VpcId           types.String `tfsdk:"vpc_id"`
	SubnetId        types.String `tfsdk:"subnet_id"`
	AccountId       types.String `tfsdk:"account_id"`
	Description     types.String `tfsdk:"description"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m LbServerGroupDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               types.StringType,
		"protocol":           types.StringType,
		"loadbalancer_id":    types.StringType,
		"lb_name":            types.StringType,
		"lb_method":          types.StringType,
		"lb_health_check_id": types.StringType,
		"state":              types.StringType,
		"vpc_id":             types.StringType,
		"subnet_id":          types.StringType,
		"account_id":         types.StringType,
		"description":        types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}

// List LB ServerGroup Response
type LbServerGroup struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Protocol                 types.String `tfsdk:"protocol"`
	LoadbalancerId           types.String `tfsdk:"loadbalancer_id"`
	LbName                   types.String `tfsdk:"lb_name"`
	State                    types.String `tfsdk:"state"`
	VpcId                    types.String `tfsdk:"vpc_id"`
	LbServerGroupMemberCount types.Int32  `tfsdk:"lb_server_group_member_count"`
	CreatedAt                types.String `tfsdk:"created_at"`
	CreatedBy                types.String `tfsdk:"created_by"`
	ModifiedAt               types.String `tfsdk:"modified_at"`
	ModifiedBy               types.String `tfsdk:"modified_by"`
}

func (m LbServerGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                           types.StringType,
		"name":                         types.StringType,
		"protocol":                     types.StringType,
		"loadbalancer_id":              types.StringType,
		"lb_name":                      types.StringType,
		"state":                        types.StringType,
		"vpc_id":                       types.StringType,
		"lb_server_group_member_count": types.Int32Type,
		"created_at":                   types.StringType,
		"created_by":                   types.StringType,
		"modified_at":                  types.StringType,
		"modified_by":                  types.StringType,
	}
}

type LbServerGroupResource struct {
	Id                  types.String        `tfsdk:"id"`
	LbServerGroup       types.Object        `tfsdk:"lb_server_group"`
	LbServerGroupCreate LbServerGroupCreate `tfsdk:"lb_server_group_create"`
}

// Create LB ServerGroup Request
type LbServerGroupCreate struct {
	Name            types.String `tfsdk:"name"`
	Protocol        types.String `tfsdk:"protocol"`
	LbMethod        types.String `tfsdk:"lb_method"`
	VpcId           types.String `tfsdk:"vpc_id"`
	SubnetId        types.String `tfsdk:"subnet_id"`
	Description     types.String `tfsdk:"description"`
	LbHealthCheckId types.String `tfsdk:"lb_health_check_id"`
	Tags            types.Map    `tfsdk:"tags"`
}

//------------ LB Member -------------------//

type LbMemberDataSource struct {
	Size            types.Int32  `tfsdk:"size"`
	Page            types.Int32  `tfsdk:"page"`
	Sort            types.String `tfsdk:"sort"`
	Name            types.String `tfsdk:"name"`
	MemberIp        types.String `tfsdk:"member_ip"`
	MemberPort      types.Int32  `tfsdk:"member_port"`
	LbServerGroupId types.String `tfsdk:"lb_server_group_id"`
	LbMembers       []LbMember   `tfsdk:"lb_members"`
}

type LbMember struct {
	Id              types.String `tfsdk:"id"`
	LbServerGroupId types.String `tfsdk:"lb_server_group_id"`
	Name            types.String `tfsdk:"name"`
	MemberIp        types.String `tfsdk:"member_ip"`
	MemberPort      types.Int32  `tfsdk:"member_port"`
	MemberState     types.String `tfsdk:"member_state"`
	MemberWeight    types.Int32  `tfsdk:"member_weight"`
	ObjectType      types.String `tfsdk:"object_type"`
	ObjectId        types.String `tfsdk:"object_id"`
	State           types.String `tfsdk:"state"`
	HealthState     types.String `tfsdk:"health_state"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m LbMember) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"lb_server_group_id": types.StringType,
		"name":               types.StringType,
		"member_ip":          types.StringType,
		"member_port":        types.Int32Type,
		"member_state":       types.StringType,
		"member_weight":      types.Int32Type,
		"object_type":        types.StringType,
		"object_id":          types.StringType,
		"state":              types.StringType,
		"health_state":       types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}

type LbMemberDataSourceDetail struct {
	Id              types.String `tfsdk:"id"`
	LbServerGroupId types.String `tfsdk:"lb_server_group_id"`
	LbMemberDetail  types.Object `tfsdk:"lb_member"`
}

type LbMemberDetail struct {
	LbServerGroupId types.String `tfsdk:"lb_server_group_id"`
	Name            types.String `tfsdk:"name"`
	MemberIp        types.String `tfsdk:"member_ip"`
	MemberPort      types.Int32  `tfsdk:"member_port"`
	MemberState     types.String `tfsdk:"member_state"`
	MemberWeight    types.Int32  `tfsdk:"member_weight"`
	ObjectType      types.String `tfsdk:"object_type"`
	ObjectId        types.String `tfsdk:"object_id"`
	State           types.String `tfsdk:"state"`
	SubnetId        types.String `tfsdk:"subnet_id"`
	Uuid            types.String `tfsdk:"uuid"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m LbMemberDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"lb_server_group_id": types.StringType,
		"name":               types.StringType,
		"member_ip":          types.StringType,
		"member_port":        types.Int32Type,
		"member_state":       types.StringType,
		"member_weight":      types.Int32Type,
		"object_type":        types.StringType,
		"object_id":          types.StringType,
		"state":              types.StringType,
		"subnet_id":          types.StringType,
		"uuid":               types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}

type LbMembersResource struct {
	LbServerGroupId types.String   `tfsdk:"lb_server_group_id"`
	Id              types.String   `tfsdk:"id"`
	LbMember        types.Object   `tfsdk:"lb_member"`
	LbMemberCreate  LbMemberCreate `tfsdk:"lb_member_create"`
}

type LbMemberCreate struct {
	Name         types.String `tfsdk:"name"`
	MemberIp     types.String `tfsdk:"member_ip"`
	MemberPort   types.Int32  `tfsdk:"member_port"`
	ObjectType   types.String `tfsdk:"object_type"`
	ObjectId     types.String `tfsdk:"object_id"`
	MemberWeight types.Int32  `tfsdk:"member_weight"`
}

type LbMemberResource struct {
	LbServerGroupId types.String   `tfsdk:"lb_server_group_id"`
	Id              types.String   `tfsdk:"id"`
	LbMember        types.Object   `tfsdk:"lb_member"`
	LbMemberSet     LbMemberCreate `tfsdk:"lb_member_create"`
}

// ------------ LB Listener -------------------//
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
	Id           types.String            `tfsdk:"id"`
	Name         types.String            `tfsdk:"name"`
	Protocol     types.String            `tfsdk:"protocol"`
	State        types.String            `tfsdk:"state"`
	ServicePort  types.Int32             `tfsdk:"service_port"`
	ServerGroups []LbListenerServerGroup `tfsdk:"server_groups"`
	CreatedAt    types.String            `tfsdk:"created_at"`
	CreatedBy    types.String            `tfsdk:"created_by"`
	ModifiedAt   types.String            `tfsdk:"modified_at"`
	ModifiedBy   types.String            `tfsdk:"modified_by"`
}

type LbListenerServerGroup struct {
	ServerGroupId   types.String `tfsdk:"server_group_id"`
	ServerGroupName types.String `tfsdk:"server_group_name"`
}

type LbListenerDataSourceDetail struct {
	Id               types.String `tfsdk:"id"`
	LbListenerDetail types.Object `tfsdk:"lb_listener"`
}

type LbListenerDetail struct {
	Id                  types.String     `tfsdk:"id"`
	CreatedAt           types.String     `tfsdk:"created_at"`
	CreatedBy           types.String     `tfsdk:"created_by"`
	ModifiedAt          types.String     `tfsdk:"modified_at"`
	ModifiedBy          types.String     `tfsdk:"modified_by"`
	Description         types.String     `tfsdk:"description"`
	HttpsRedirection    types.Bool       `tfsdk:"https_redirection"`
	InsertClientIp      types.Bool       `tfsdk:"insert_client_ip"`
	Name                types.String     `tfsdk:"name"`
	Persistence         types.String     `tfsdk:"persistence"`
	Protocol            types.String     `tfsdk:"protocol"`
	ServerGroupId       types.String     `tfsdk:"server_group_id"`
	ServerGroupName     types.String     `tfsdk:"server_group_name"`
	ServicePort         types.Int32      `tfsdk:"service_port"`
	ResponseTimeout     types.Int32      `tfsdk:"response_timeout"`
	SessionDurationTime types.Int32      `tfsdk:"session_duration_time"`
	SslCertificate      *SslCertificate  `tfsdk:"ssl_certificate"`
	State               types.String     `tfsdk:"state"`
	UrlHandler          []UrlHandler     `tfsdk:"url_handler"`
	UrlRedirection      []UrlRedirection `tfsdk:"url_redirection"`
	XForwardedFor       types.Bool       `tfsdk:"x_forwarded_for"`
	XForwardedPort      types.Bool       `tfsdk:"x_forwarded_port"`
	XForwardedProto     types.Bool       `tfsdk:"x_forwarded_proto"`
}

func (m LbListenerDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"description":           types.StringType,
		"https_redirection":     types.BoolType,
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
				"server_cert_id":    types.StringType,
				"server_cert_level": types.StringType,
			},
		},
		"state": types.StringType,
		"url_handler": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"url_pattern":     types.StringType,
					"server_group_id": types.StringType,
				},
			},
		},
		"url_redirection": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"url_pattern":          types.StringType,
					"redirect_url_pattern": types.StringType,
				},
			},
		},
		"x_forwarded_for":   types.BoolType,
		"x_forwarded_port":  types.BoolType,
		"x_forwarded_proto": types.BoolType,
	}
}

type SslCertificate struct {
	ClientCertId    types.String `tfsdk:"client_cert_id"`
	ClientCertLevel types.String `tfsdk:"client_cert_level"`
	ServerCertId    types.String `tfsdk:"server_cert_id"`
	ServerCertLevel types.String `tfsdk:"server_cert_level"`
}

type UrlHandler struct {
	UrlPattern    types.String `tfsdk:"url_pattern"`
	ServerGroupId types.String `tfsdk:"server_group_id"`
}

type UrlRedirection struct {
	UrlPattern         types.String `tfsdk:"url_pattern"`
	RedirectUrlPattern types.String `tfsdk:"redirect_url_pattern"`
}

func convertUrlHandlerToInterface(handlers []UrlHandler) []interface{} {
	result := make([]interface{}, len(handlers))
	for i, handler := range handlers {
		result[i] = map[string]string{
			"url_pattern":     handler.UrlPattern.ValueString(),
			"server_group_id": handler.ServerGroupId.ValueString(),
		}
	}
	return result
}

func convertUrlRedirectionToInterface(redirections []UrlRedirection) []interface{} {
	result := make([]interface{}, len(redirections))
	for i, handler := range redirections {
		result[i] = map[string]string{
			"url_pattern":          handler.UrlPattern.ValueString(),
			"redirect_url_pattern": handler.RedirectUrlPattern.ValueString(),
		}
	}
	return result
}

//------------ LB Health Check -------------------//

type LbHealthCheckDataSource struct {
	Size           types.Int32     `tfsdk:"size"`
	Page           types.Int32     `tfsdk:"page"`
	Sort           types.String    `tfsdk:"sort"`
	Name           types.String    `tfsdk:"name"`
	Protocol       types.List      `tfsdk:"protocol"`
	SubnetId       types.String    `tfsdk:"subnet_id"`
	LbHealthChecks []LbHealthCheck `tfsdk:"lb_health_checks"`
}

type LbHealthCheckDataSourceDetail struct {
	Id                  types.String `tfsdk:"id"`
	LbHealthCheckDetail types.Object `tfsdk:"lb_health_check"`
}

type LbHealthCheckResource struct {
	Id                  types.String        `tfsdk:"id"`
	LbHealthCheck       types.Object        `tfsdk:"lb_health_check"`
	LbHealthCheckCreate LbHealthCheckCreate `tfsdk:"lb_health_check_create"`
}

type LbHealthCheckCreate struct {
	Name                types.String `tfsdk:"name"`
	VpcId               types.String `tfsdk:"vpc_id"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	Protocol            types.String `tfsdk:"protocol"`
	HealthCheckPort     types.Int32  `tfsdk:"health_check_port"`
	HealthCheckInterval types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckTimeout  types.Int32  `tfsdk:"health_check_timeout"`
	HealthCheckCount    types.Int32  `tfsdk:"health_check_count"`
	HttpMethod          types.String `tfsdk:"http_method"`
	HealthCheckUrl      types.String `tfsdk:"health_check_url"`
	ResponseCode        types.String `tfsdk:"response_code"`
	RequestData         types.String `tfsdk:"request_data"`
	Description         types.String `tfsdk:"description"`
	Tags                types.Map    `tfsdk:"tags"`
}

type LbHealthCheck struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	SubnetId           types.String `tfsdk:"subnet_id"`
	Protocol           types.String `tfsdk:"protocol"`
	LbServerGroupCount types.Int32  `tfsdk:"lb_server_group_count"`
	HealthCheckType    types.String `tfsdk:"health_check_type"`
	State              types.String `tfsdk:"state"`
	CreatedAt          types.String `tfsdk:"created_at"`
	CreatedBy          types.String `tfsdk:"created_by"`
	ModifiedAt         types.String `tfsdk:"modified_at"`
	ModifiedBy         types.String `tfsdk:"modified_by"`
}

type LbHealthCheckDetail struct {
	Name                types.String `tfsdk:"name"`
	VpcId               types.String `tfsdk:"vpc_id"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	Protocol            types.String `tfsdk:"protocol"`
	HealthCheckPort     types.Int32  `tfsdk:"health_check_port"`
	HealthCheckInterval types.Int32  `tfsdk:"health_check_interval"`
	HealthCheckTimeout  types.Int32  `tfsdk:"health_check_timeout"`
	HealthCheckCount    types.Int32  `tfsdk:"health_check_count"`
	HttpMethod          types.String `tfsdk:"http_method"`
	HealthCheckUrl      types.String `tfsdk:"health_check_url"`
	ResponseCode        types.String `tfsdk:"response_code"`
	RequestData         types.String `tfsdk:"request_data"`
	HealthCheckType     types.String `tfsdk:"health_check_type"`
	Description         types.String `tfsdk:"description"`
	State               types.String `tfsdk:"state"`
	AccountId           types.String `tfsdk:"account_id"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
}

func (m LbHealthCheckDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"vpc_id":                types.StringType,
		"subnet_id":             types.StringType,
		"protocol":              types.StringType,
		"health_check_port":     types.Int32Type,
		"health_check_interval": types.Int32Type,
		"health_check_timeout":  types.Int32Type,
		"health_check_count":    types.Int32Type,
		"http_method":           types.StringType,
		"health_check_url":      types.StringType,
		"response_code":         types.StringType,
		"request_data":          types.StringType,
		"health_check_type":     types.StringType,
		"description":           types.StringType,
		"state":                 types.StringType,
		"account_id":            types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}

type LbListenerResource struct {
	Id               types.String     `tfsdk:"id"`
	LbListener       types.Object     `tfsdk:"lb_listener"`
	LbListenerCreate LbListenerCreate `tfsdk:"lb_listener_create"`
}

type LbListenerCreate struct {
	Description         types.String     `tfsdk:"description"`
	HttpsRedirection    types.Bool       `tfsdk:"https_redirection"`
	InsertClientIp      types.Bool       `tfsdk:"insert_client_ip"`
	LoadbalancerId      types.String     `tfsdk:"loadbalancer_id"`
	Name                types.String     `tfsdk:"name"`
	Persistence         types.String     `tfsdk:"persistence"`
	Protocol            types.String     `tfsdk:"protocol"`
	ResponseTimeout     types.Int32      `tfsdk:"response_timeout"`
	ServerGroupId       types.String     `tfsdk:"server_group_id"`
	ServicePort         types.Int32      `tfsdk:"service_port"`
	SessionDurationTime types.Int32      `tfsdk:"session_duration_time"`
	SslCertificate      *SslCertificate  `tfsdk:"ssl_certificate"`
	UrlHandler          []UrlHandler     `tfsdk:"url_handler"`
	UrlRedirection      []UrlRedirection `tfsdk:"url_redirection"`
	XForwardedFor       types.Bool       `tfsdk:"x_forwarded_for"`
	XForwardedPort      types.Bool       `tfsdk:"x_forwarded_port"`
	XForwardedProto     types.Bool       `tfsdk:"x_forwarded_proto"`
}

// ------------ LB Certificate -------------------//
type LbCertificateDataSource struct {
	LbCertificates []LbCertificate `tfsdk:"lb_certificates"`
}

type LbCertificate struct {
	CertKind    types.String `tfsdk:"cert_kind"`
	Cn          types.String `tfsdk:"cn"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	NotAfterDt  types.String `tfsdk:"not_after_dt"`
	NotBeforeDt types.String `tfsdk:"not_before_dt"`
	State       types.String `tfsdk:"state"`
}

type LbCertificateDataSourceDetail struct {
	Id                  types.String `tfsdk:"id"`
	LbCertificateDetail types.Object `tfsdk:"lb_certificate"`
}

type LbCertificateDetail struct {
	AccountId    types.String `tfsdk:"account_id"`
	CertBody     types.String `tfsdk:"cert_body"`
	CertChain    types.String `tfsdk:"cert_chain"`
	CertKind     types.String `tfsdk:"cert_kind"`
	Cn           types.String `tfsdk:"cn"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	Id           types.String `tfsdk:"id"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	Name         types.String `tfsdk:"name"`
	NotAfterDt   types.String `tfsdk:"not_after_dt"`
	NotBeforeDt  types.String `tfsdk:"not_before_dt"`
	Organization types.String `tfsdk:"organization"`
	PrivateKey   types.String `tfsdk:"private_key"`
	State        types.String `tfsdk:"state"`
}

func (m LbCertificateDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":    types.StringType,
		"cert_body":     types.StringType,
		"cert_chain":    types.StringType,
		"cert_kind":     types.StringType,
		"cn":            types.StringType,
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"id":            types.StringType,
		"modified_at":   types.StringType,
		"modified_by":   types.StringType,
		"name":          types.StringType,
		"not_after_dt":  types.StringType,
		"not_before_dt": types.StringType,
		"organization":  types.StringType,
		"private_key":   types.StringType,
		"state":         types.StringType,
	}
}
