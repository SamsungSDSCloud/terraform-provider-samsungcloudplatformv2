package servicewatch

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-servicewatch"

type DashboardDataSources struct {
	Name            types.String `tfsdk:"name"`
	NameLike        types.String `tfsdk:"name_like"`
	FavoriteEnabled types.Bool   `tfsdk:"favorite_enabled"`
	Type            types.String `tfsdk:"type"`
	ServiceCode     types.String `tfsdk:"service_code"`
	Dashboards      types.List   `tfsdk:"dashboards"`
}

type DashboardDataSource struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	FavoriteEnabled types.Bool   `tfsdk:"favorite_enabled"`
	Srn             types.String `tfsdk:"srn"`
	ShareType       types.String `tfsdk:"share_type"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	Widgets         types.List   `tfsdk:"widgets"`
}

type DashboardResource struct {
	LastUpdated types.String `tfsdk:"last_updated"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Srn         types.String `tfsdk:"srn"`
	ShareType   types.String `tfsdk:"share_type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Widgets     types.List   `tfsdk:"widgets"`
}

type Dashboard struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	FavoriteEnabled types.Bool   `tfsdk:"favorite_enabled"`
	CreatedAt       types.String `tfsdk:"created_at"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
}

func (m Dashboard) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"type":             types.StringType,
		"favorite_enabled": types.BoolType,
		"created_at":       types.StringType,
		"modified_at":      types.StringType,
	}
}

type Widget struct {
	Id         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	Width      types.Int32  `tfsdk:"width"`
	Height     types.Int32  `tfsdk:"height"`
	Order      types.Int32  `tfsdk:"order"`
	Properties types.Object `tfsdk:"properties"`
}

func (m Widget) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.StringType,
		"type":   types.StringType,
		"width":  types.Int32Type,
		"height": types.Int32Type,
		"order":  types.Int32Type,
		"properties": types.ObjectType{
			AttrTypes: Properties{}.AttributeTypes(),
		},
	}
}

type Properties struct {
	Title         types.String `tfsdk:"title"`
	Period        types.Int32  `tfsdk:"period"`
	Stacked       types.Bool   `tfsdk:"stacked"`
	StatisticType types.String `tfsdk:"statistic_type"`
	View          types.String `tfsdk:"view"`
	Metrics       types.List   `tfsdk:"metrics"`
}

func (m Properties) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"title":          types.StringType,
		"period":         types.Int32Type,
		"stacked":        types.BoolType,
		"statistic_type": types.StringType,
		"view":           types.StringType,
		"metrics": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Metric{}.AttributeTypes(),
			},
		},
	}
}

type Metric struct {
	Name          types.String `tfsdk:"name"`
	NamespaceName types.String `tfsdk:"namespace_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	Color         types.String `tfsdk:"color"`
	Dimensions    types.List   `tfsdk:"dimensions"`
	Period        types.Int32  `tfsdk:"period"`
	StatisticType types.String `tfsdk:"statistic_type"`
}

func (m Metric) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":           types.StringType,
		"namespace_name": types.StringType,
		"display_name":   types.StringType,
		"color":          types.StringType,
		"dimensions": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Dimension{}.AttributeTypes(),
			},
		},
		"period":         types.Int32Type,
		"statistic_type": types.StringType,
	}
}

type Dimension struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (m Dimension) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}

}

type LogGroupDataSources struct {
	Ids              []types.String `tfsdk:"ids"`
	Name             types.String   `tfsdk:"name"`
	RetentionPeriods []types.Int32  `tfsdk:"retention_periods"`
	LogGroups        types.List     `tfsdk:"log_groups"`
}

type LogGroupDataSource struct {
	LogGroupId types.String `tfsdk:"log_group_id"`
	LogGroup   types.Object `tfsdk:"log_group"`
}

type LogGroupResource struct {
	LastUpdated     types.String `tfsdk:"last_updated"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RetentionPeriod types.Int32  `tfsdk:"retention_period"`
	Tags            types.Map    `tfsdk:"tags"`
	LogGroup        types.Object `tfsdk:"log_group"`
}

type LogGroup struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	AccountId           types.String `tfsdk:"account_id"`
	RetentionPeriod     types.Int32  `tfsdk:"retention_period"`
	RetentionPeriodName types.String `tfsdk:"retention_period_name"`
	Status              types.String `tfsdk:"status"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
}

func (m LogGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"account_id":            types.StringType,
		"retention_period":      types.Int32Type,
		"retention_period_name": types.StringType,
		"status":                types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}

type LogStreamDataSource struct {
	LogGroupId  types.String `tfsdk:"log_group_id"`
	LogStreamId types.String `tfsdk:"log_stream_id"`
	LogStream   types.Object `tfsdk:"log_stream"`
}

type LogStreamResource struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	LogGroupId types.String `tfsdk:"log_group_id"`
	LogStream  types.Object `tfsdk:"log_stream"`
}

type LogStream struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	LogGroupId types.String `tfsdk:"log_group_id"`
	CollectYn  types.String `tfsdk:"collect_yn"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	ModifiedAt types.String `tfsdk:"modified_at"`
	ModifiedBy types.String `tfsdk:"modified_by"`
}

func (m LogStream) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"log_group_id": types.StringType,
		"collect_yn":   types.StringType,
		"created_at":   types.StringType,
		"created_by":   types.StringType,
		"modified_at":  types.StringType,
		"modified_by":  types.StringType,
	}
}

type AlertDataSource struct {
	Id    types.String `tfsdk:"id"`
	Alert types.Object `tfsdk:"alert"`
}

type AlertResource struct {
	LastUpdated       types.String  `tfsdk:"last_updated"`
	Id                types.String  `tfsdk:"id"`
	Name              types.String  `tfsdk:"name"`
	Type              types.String  `tfsdk:"type"`
	Description       types.String  `tfsdk:"description"`
	ActivatedYn       types.String  `tfsdk:"activated_yn"`
	Level             types.String  `tfsdk:"level"`
	NamespaceId       types.String  `tfsdk:"namespace_id"`
	NamespaceName     types.String  `tfsdk:"namespace_name"`
	MetricId          types.String  `tfsdk:"metric_id"`
	MetricName        types.String  `tfsdk:"metric_name"`
	Dimensions        types.List    `tfsdk:"dimensions"`
	Period            types.Int32   `tfsdk:"period"`
	Statistic         types.String  `tfsdk:"statistic"`
	EvaluationCount   types.Int32   `tfsdk:"evaluation_count"`
	Threshold         types.Float32 `tfsdk:"threshold"`
	UpperBound        types.Float32 `tfsdk:"upper_bound"`
	LowerBound        types.Float32 `tfsdk:"lower_bound"`
	Operator          types.String  `tfsdk:"operator"`
	ViolationCount    types.Int32   `tfsdk:"violation_count"`
	MissingDataOption types.String  `tfsdk:"missing_data_option"`
	RecipientIds      types.List    `tfsdk:"recipient_ids"`
	Tags              types.Map     `tfsdk:"tags"`
	CreatedAt         types.String  `tfsdk:"created_at"`
	CreatedBy         types.String  `tfsdk:"created_by"`
	ModifiedAt        types.String  `tfsdk:"modified_at"`
	ModifiedBy        types.String  `tfsdk:"modified_by"`
}

type Alert struct {
	Name                 types.String  `tfsdk:"name"`
	Description          types.String  `tfsdk:"description"`
	Srn                  types.String  `tfsdk:"srn"`
	ActivatedYn          types.String  `tfsdk:"activated_yn"`
	State                types.String  `tfsdk:"state"`
	Level                types.String  `tfsdk:"level"`
	NamespaceId          types.String  `tfsdk:"namespace_id"`
	NamespaceName        types.String  `tfsdk:"namespace_name"`
	MetricId             types.String  `tfsdk:"metric_id"`
	MetricName           types.String  `tfsdk:"metric_name"`
	MetricUnit           types.String  `tfsdk:"metric_unit"`
	Dimensions           types.List    `tfsdk:"dimensions"`
	Period               types.Int32   `tfsdk:"period"`
	Statistic            types.String  `tfsdk:"statistic"`
	EvaluationCount      types.Int32   `tfsdk:"evaluation_count"`
	EvaluationTimeWindow types.Int32   `tfsdk:"evaluation_time_window"`
	Threshold            types.Float32 `tfsdk:"threshold"`
	UpperBound           types.Float32 `tfsdk:"upper_bound"`
	LowerBound           types.Float32 `tfsdk:"lower_bound"`
	Operator             types.String  `tfsdk:"operator"`
	ViolationCount       types.Int32   `tfsdk:"violation_count"`
	MissingDataOption    types.String  `tfsdk:"missing_data_option"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	CreatedBy            types.String  `tfsdk:"created_by"`
	ModifiedAt           types.String  `tfsdk:"modified_at"`
	ModifiedBy           types.String  `tfsdk:"modified_by"`
}

func (m Alert) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":           types.StringType,
		"description":    types.StringType,
		"srn":            types.StringType,
		"activated_yn":   types.StringType,
		"state":          types.StringType,
		"level":          types.StringType,
		"namespace_id":   types.StringType,
		"namespace_name": types.StringType,
		"metric_id":      types.StringType,
		"metric_name":    types.StringType,
		"metric_unit":    types.StringType,
		"dimensions": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Dimension{}.AttributeTypes(),
			},
		},
		"period":                 types.Int32Type,
		"statistic":              types.StringType,
		"evaluation_count":       types.Int32Type,
		"evaluation_time_window": types.Int32Type,
		"threshold":              types.Float32Type,
		"upper_bound":            types.Float32Type,
		"lower_bound":            types.Float32Type,
		"operator":               types.StringType,
		"violation_count":        types.Int32Type,
		"missing_data_option":    types.StringType,
		"created_at":             types.StringType,
		"created_by":             types.StringType,
		"modified_at":            types.StringType,
		"modified_by":            types.StringType,
	}
}

type AlertUpdateModel struct {
	Level             types.String
	NamespaceId       types.String
	MetricId          types.String
	Dimensions        types.List
	Period            types.Int32
	Statistic         types.String
	EvaluationCount   types.Int32
	Threshold         types.Float32
	UpperBound        types.Float32
	LowerBound        types.Float32
	Operator          types.String
	ViolationCount    types.Int32
	MissingDataOption types.String
}

type EventRuleDataSource struct {
	EventRule   types.Object `tfsdk:"event_rule"`
	EventRuleId types.String `tfsdk:"event_rule_id"`
}

type EventRule struct {
	AccountId      types.String `tfsdk:"account_id"`
	ActiveYn       types.String `tfsdk:"active_yn"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Description    types.String `tfsdk:"description"`
	Id             types.String `tfsdk:"id"`
	ModifiedAt     types.String `tfsdk:"modified_at"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	Name           types.String `tfsdk:"name"`
	ResourceTypeId types.String `tfsdk:"resource_type_id"`
	ServiceId      types.String `tfsdk:"service_id"`
}

type EventRuleResource struct {
	Id             types.String   `tfsdk:"id"`
	LastUpdated    types.String   `tfsdk:"last_updated"`
	Description    types.String   `tfsdk:"description"`
	EventIds       []types.String `tfsdk:"event_ids"`
	EventRuleId    types.String   `tfsdk:"event_rule_id"`
	Name           types.String   `tfsdk:"name"`
	RecipientIds   []types.String `tfsdk:"recipient_ids"`
	ResourceTypeId types.String   `tfsdk:"resource_type_id"`
	ServiceId      types.String   `tfsdk:"service_id"`
	SrnList        []types.String `tfsdk:"srn_list"`
	Tags           types.Map      `tfsdk:"tags"`
	EventRule      types.Object   `tfsdk:"event_rule"`
	ActiveYn       types.String   `tfsdk:"active_yn"`
	NoneAttributes []types.String `tfsdk:"none_attributes"`
}

func (v EventRule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":       types.StringType,
		"active_yn":        types.StringType,
		"created_at":       types.StringType,
		"created_by":       types.StringType,
		"description":      types.StringType,
		"id":               types.StringType,
		"modified_at":      types.StringType,
		"modified_by":      types.StringType,
		"name":             types.StringType,
		"resource_type_id": types.StringType,
		"service_id":       types.StringType,
	}
}
