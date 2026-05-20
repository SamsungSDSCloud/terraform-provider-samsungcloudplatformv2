package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchAlertDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceWatchAlertDataSource{}
)

// NewServiceWatchAlertDataSource is a helper function to simplify the provider implementation.
func NewServiceWatchAlertDataSource() datasource.DataSource {
	return &serviceWatchAlertDataSource{}
}

// serviceWatchAlertDataSource is the data source implementation.
type serviceWatchAlertDataSource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchAlertDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_alert"
}

// Schema defines the schema for the data source.
func (d *serviceWatchAlertDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Alert Data Source",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Alert ID",
				Optional:    true,
			},
			common.ToSnakeCase("Alert"): schema.SingleNestedAttribute{
				Description: "Alert",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Alert name",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Alert description",
						Computed:    true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description: "SDS cloud resource name of the Alert",
						Computed:    true,
					},
					common.ToSnakeCase("ActivatedYn"): schema.StringAttribute{
						Description: "Whether the Alert is activated or not",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Alert state - NORMAL, ALERT, INSUFFICIENT_DATA",
						Computed:    true,
					},
					common.ToSnakeCase("Level"): schema.StringAttribute{
						Description: "Alert level - HIGH, MIDDLE, LOW",
						Computed:    true,
					},
					common.ToSnakeCase("NamespaceId"): schema.StringAttribute{
						Description: "Namespace ID",
						Computed:    true,
					},
					common.ToSnakeCase("NamespaceName"): schema.StringAttribute{
						Description: "Namespace name",
						Computed:    true,
					},
					common.ToSnakeCase("MetricId"): schema.StringAttribute{
						Description: "Sharing type",
						Computed:    true,
					},
					common.ToSnakeCase("MetricName"): schema.StringAttribute{
						Description: "Metric name",
						Computed:    true,
					},
					common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
						Description: "Metric unit",
						Computed:    true,
					},
					common.ToSnakeCase("Dimensions"): schema.ListNestedAttribute{
						Description: "List of dimensions",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Key"): schema.StringAttribute{
									Description: "Dimension key",
									Computed:    true,
								},
								common.ToSnakeCase("Value"): schema.StringAttribute{
									Description: "Dimensions value",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("Period"): schema.Int32Attribute{
						Description: "Period (seconds)",
						Computed:    true,
					},
					common.ToSnakeCase("Statistic"): schema.StringAttribute{
						Description: "Statistic - SUM, AVG, MAX, MIN",
						Computed:    true,
					},
					common.ToSnakeCase("EvaluationCount"): schema.Int32Attribute{
						Description: "Evaluation count for the Alert condition",
						Computed:    true,
					},
					common.ToSnakeCase("EvaluationTimeWindow"): schema.Int32Attribute{
						Description: "Evaluation time window (period * evaluation_count)",
						Computed:    true,
					},
					common.ToSnakeCase("Threshold"): schema.Float32Attribute{
						Description: "Threshold for the Alert condition (except for \"RANGE\" operator)",
						Computed:    true,
					},
					common.ToSnakeCase("UpperBound"): schema.Float32Attribute{
						Description: "Upper bound for the Alert \"RANGE\" operator",
						Computed:    true,
					},
					common.ToSnakeCase("LowerBound"): schema.Float32Attribute{
						Description: "Lower bound for the Alert \"RANGE\" operator",
						Computed:    true,
					},
					common.ToSnakeCase("Operator"): schema.StringAttribute{
						Description: "Operator - EQ, NOT_EQ, GT, GTE, LT, LTE, RANGE",
						Computed:    true,
					},
					common.ToSnakeCase("ViolationCount"): schema.Int32Attribute{
						Description: "Violation count for the Alert condition",
						Computed:    true,
					},
					common.ToSnakeCase("MissingDataOption"): schema.StringAttribute{
						Description: "Missing data option - MISSING, BREACHING, NOT_BREACHING, IGNORE",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created date time",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Creator ID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{

						Description: "Modified date time",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modifier ID",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchAlertDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			ErrUnexpectedConfigure,
			fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData),
		)

		return
	}

	d.client = inst.Client.ServiceWatch
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *serviceWatchAlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.AlertDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertResp, err := d.client.GetAlert(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadAlert,
			fmt.Sprintf(ErrReadAlertFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}

	// Dimensions 변환
	var dimensions []servicewatch.Dimension
	for _, item := range alertResp.GetDimensions() {
		dimensions = append(dimensions, servicewatch.Dimension{
			Key:   types.StringValue(item.GetKey()),
			Value: types.StringValue(item.GetValue()),
		})
	}
	dimensionList, diags := types.ListValueFrom(
		ctx, types.ObjectType{
			AttrTypes: servicewatch.Dimension{}.AttributeTypes(),
		},
		dimensions)
	resp.Diagnostics.Append(diags...)
	fmt.Printf("\nDimensions: %+v\n", dimensionList)

	// convert Alert datasource model
	alert := servicewatch.Alert{
		Name:                 types.StringValue(alertResp.GetName()),
		Description:          nullableStringTypes(alertResp.GetDescriptionOk()),
		Srn:                  types.StringValue(alertResp.GetSrn()),
		ActivatedYn:          types.StringValue(string(alertResp.GetActivatedYn())),
		State:                types.StringValue(string(alertResp.GetState())),
		Level:                types.StringValue(string(alertResp.GetLevel())),
		NamespaceId:          types.StringValue(alertResp.GetNamespaceId()),
		NamespaceName:        types.StringValue(alertResp.GetNamespaceName()),
		MetricId:             types.StringValue(alertResp.GetMetricId()),
		MetricName:           types.StringValue(alertResp.GetMetricName()),
		MetricUnit:           types.StringValue(alertResp.GetMetricUnit()),
		Dimensions:           dimensionList,
		Period:               types.Int32Value(alertResp.GetPeriod()),
		Statistic:            types.StringValue(string(alertResp.GetStatistic())),
		EvaluationCount:      types.Int32Value(alertResp.GetEvaluationCount()),
		EvaluationTimeWindow: types.Int32Value(alertResp.GetEvaluationTimeWindow()),
		Threshold:            nullableFloat32Types(alertResp.GetThresholdOk()),
		UpperBound:           nullableFloat32Types(alertResp.GetUpperBoundOk()),
		LowerBound:           nullableFloat32Types(alertResp.GetLowerBoundOk()),
		Operator:             types.StringValue(string(alertResp.GetOperator())),
		ViolationCount:       types.Int32Value(alertResp.GetViolationCount()),
		MissingDataOption:    types.StringValue(string(alertResp.GetMissingDataOption())),
		CreatedAt:            types.StringValue(alertResp.GetCreatedAt().Format(TimeFormatDisplay)),
		ModifiedAt:           types.StringValue(alertResp.GetModifiedAt().Format(TimeFormatDisplay)),
		CreatedBy:            types.StringValue(alertResp.GetCreatedBy()),
		ModifiedBy:           types.StringValue(alertResp.GetModifiedBy()),
	}
	fmt.Printf("\nalert: %+v\n", alert)
	state.Id = types.StringValue(state.Id.ValueString())
	alertObjectValue, diags := types.ObjectValueFrom(ctx, alert.AttributeTypes(), alert)
	fmt.Printf("\nalertObjectValue: %+v\n", alertObjectValue)
	state.Alert = alertObjectValue
	fmt.Printf("\nstate: %+v\n", state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "ObjectValueFrom failed", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}
}

func nullableFloat32Types(val *float32, isSet bool) types.Float32 {
	if isSet && val != nil {
		return types.Float32Value(*val)
	}
	return types.Float32Null()
}

func nullableStringTypes(val *string, isSet bool) types.String {
	if isSet && val != nil {
		return types.StringValue(*val)
	}
	return types.StringNull()
}

func nullableInt32Types(val *int32, isSet bool) types.Int32 {
	if isSet && val != nil {
		return types.Int32Value(*val)
	}
	return types.Int32Null()
}
