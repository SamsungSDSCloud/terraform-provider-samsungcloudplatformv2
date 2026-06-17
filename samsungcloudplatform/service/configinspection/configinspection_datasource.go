package configinspection

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	scpci "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &configinspectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &configinspectionsDataSource{}
)

// Helper function to simplify the provider implementation.
func NewConfigInspectionConfigInspectionDataSource() datasource.DataSource {
	return &configinspectionsDataSource{}
}

// Data source implementation.
type configinspectionsDataSource struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *scpci.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name
func (d *configinspectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection_configinspection"
}

// Schema defines the schema for the data source
func (d *configinspectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Detail of config inspection object.",
		Attributes: map[string]schema.Attribute{
			// Input attributes
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis.\n" +
					"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
				Required: true,
			},

			// Output attributes
			common.ToSnakeCase("AuthKeyResponses"): schema.SingleNestedAttribute{
				Description: "Authentication key response.",
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AuthKeyCreatedAt"): schema.StringAttribute{
						Description: "Created date of authkey.\n" +
							"  - example : '2022-01-01T12:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("AuthKeyExpiredAt"): schema.StringAttribute{
						Description: "Expired date of authkey.\n" +
							"  - example : '2022-01-01T12:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("AuthKeyId"): schema.StringAttribute{
						Description: "Id of auth key.\n" +
							"  - example : '9b72a9856e494exxxxxxxxxxxxxxxxxx'",
						Computed: true,
					},
					common.ToSnakeCase("AuthKeyState"): schema.StringAttribute{
						Description: "State of auth key.\n" +
							"  - example : 'ACTIVATED'",
						Computed: true},
					common.ToSnakeCase("UserId"): schema.StringAttribute{
						Description: "Account owner of this diagnosis.\n" +
							"  - example : '4f5d60e9e08b48d0a08xxxxxxxxxxxxx'",
						Computed: true,
					},
				},
				Computed: true,
			},
			common.ToSnakeCase("ScheduleResponse"): schema.SingleNestedAttribute{
				Description: "Diagnosis schedule response.",
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Id of diagnosis.\n" +
							"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisStartTimePattern"): schema.StringAttribute{
						Description: "Start time (5-minute increments, 00 to 23 hours, 00 to 55 minutes).\n" +
							"  - example : '08:00'",
						Computed: true,
					},
					common.ToSnakeCase("FrequencyType"): schema.StringAttribute{
						Description: "Schedule type (monthly, weekly, daily).\n" +
							"  - example : 'MONTH'",
						Computed: true,
					},
					common.ToSnakeCase("FrequencyValue"): schema.StringAttribute{
						Description: "Schedule value (01~31, MONDAY~SUNDAY, everyDay).\n" +
							"  - example : 1",
						Computed: true,
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeBp"): schema.StringAttribute{
						Description: "Checklist Best Practice Use.\n" +
							"  - example : 'y'",
						Computed: true,
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeSsi"): schema.StringAttribute{
						Description: "Checklist SSI usage.\n" +
							"  - example : 'y'",
						Computed: true,
					},
				},
				Computed: true,
			},
			common.ToSnakeCase("SummaryResponses"): schema.SingleNestedAttribute{
				Description: "Summary of the config inspection diagnosis, including identity, configuration, and latest run status.",
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created date.\n" +
							"  - example : '2022-01-01T12:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("CspType"): schema.StringAttribute{
						Description: "Type of cloud service provider.\n" +
							"  - example : 'SCP'\n" +
							"  - enum : SCP | AWS | Azure",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
						Description: "Account Id of diagnosis.\n" +
							"  - example : '0e3dffc50eb247a1adxxxxxxxxxxxxxx'",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
						Description: "Check type of diagnosis.\n" +
							"  - example : 'BP'\n" +
							"  - enum : BP | SSI",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Id of diagnosis.\n" +
							"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
						Description: "Name of diagnosis.\n" +
							"  - example : 'Sample Diagnosis Name'\n" +
							"  - pattern : `^[a-zA-Z0-9-_]+$`",
						Computed: true,
					},
					common.ToSnakeCase("DiagnosisType"): schema.StringAttribute{
						Description: "How the diagnosis is executed.\n" +
							"  - example : 'Console'",
						Computed: true,
					},
					common.ToSnakeCase("PlanType"): schema.StringAttribute{
						Description: "Billing plan for the inspection.\n" +
							"  - example : 'STANDARD'\n" +
							"  - enum : STANDARD | MONTHLY",
						Computed: true,
					},
					common.ToSnakeCase("ErrorState"): schema.StringAttribute{
						Description: "Error type of recent diagnosis.\n" +
							"  - example : 'CONNECTION_FAIL'",
						Computed: true,
					},
					common.ToSnakeCase("RecentDiagnosisAt"): schema.StringAttribute{
						Description: "Recent diagnosis date.\n" +
							"  - example : '2022-01-01T12:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("RecentDiagnosisState"): schema.StringAttribute{
						Description: "Status of the most recent diagnosis run.\n" +
							"  - example : 'Completed'",
						Computed: true,
					},
				},
				Computed: true,
			},
		},
	}
}

// Configure prepares the data source with the provider configuration
func (d *configinspectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = inst.Client.ConfigInspection
	d.clients = inst.Client
}

// Read fetches and reads the resource data
func (d *configinspectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scpci.ConfigInspectionDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch data from the API
	res, err := d.client.GetConfigInspectionObjectDetail(ctx, state.DiagnosisID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch config inspection detail", err.Error())
		return
	}

	// Map API response to the state object
	state.AuthKeyResponses = &scpci.AuthKeyResponse{
		AuthKeyCreatedAt: types.StringPointerValue(res.AuthKeyResponses.AuthKeyCreatedAt),
		AuthKeyExpiredAt: types.StringPointerValue(res.AuthKeyResponses.AuthKeyExpiredAt),
		AuthKeyId:        types.StringPointerValue(res.AuthKeyResponses.AuthKeyId),
		AuthKeyState:     types.StringPointerValue(res.AuthKeyResponses.AuthKeyState),
		UserId:           types.StringPointerValue(res.AuthKeyResponses.UserId),
	}

	state.ScheduleResponse = &scpci.DiagnosisScheduleResponse{
		DiagnosisId:               types.StringPointerValue(res.ScheduleResponse.DiagnosisId),
		DiagnosisStartTimePattern: types.StringPointerValue(res.ScheduleResponse.DiagnosisStartTimePattern),
		FrequencyType:             types.StringPointerValue(res.ScheduleResponse.FrequencyType),
		FrequencyValue:            types.StringPointerValue(res.ScheduleResponse.FrequencyValue),
		UseDiagnosisCheckTypeBp:   types.StringPointerValue(res.ScheduleResponse.UseDiagnosisCheckTypeBp),
		UseDiagnosisCheckTypeSsi:  types.StringPointerValue(res.ScheduleResponse.UseDiagnosisCheckTypeSsi),
	}

	state.SummaryResponses = &scpci.SummaryResponse{
		CreatedAt:            types.StringValue(res.SummaryResponses.CreatedAt.Format(time.RFC1123)),
		CspType:              types.StringValue(res.SummaryResponses.CspType),
		DiagnosisAccountId:   types.StringValue(res.SummaryResponses.DiagnosisAccountId),
		DiagnosisCheckType:   types.StringValue(res.SummaryResponses.DiagnosisCheckType),
		DiagnosisId:          types.StringValue(res.SummaryResponses.DiagnosisId),
		DiagnosisName:        types.StringValue(res.SummaryResponses.DiagnosisName),
		DiagnosisType:        types.StringValue(res.SummaryResponses.DiagnosisType),
		PlanType:             types.StringValue(res.SummaryResponses.PlanType),
		ErrorState:           types.StringPointerValue(res.SummaryResponses.ErrorState.Get()),
		RecentDiagnosisState: types.StringPointerValue(res.SummaryResponses.RecentDiagnosisState.Get()),
	}

	if res.SummaryResponses.RecentDiagnosisAt.IsSet() {
		state.SummaryResponses.RecentDiagnosisAt = types.StringValue(res.SummaryResponses.GetRecentDiagnosisAt().Format(time.RFC1123))
	}

	// Map other fields similarly
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
