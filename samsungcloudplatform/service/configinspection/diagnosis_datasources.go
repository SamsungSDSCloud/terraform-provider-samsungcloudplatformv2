package configinspection

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	scpci "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &diagnosesDataSource{}
	_ datasource.DataSourceWithConfigure = &diagnosesDataSource{}
)

// Helper function to simplify the provider implementation.
func NewConfigInspectionDiagnosisDataSources() datasource.DataSource {
	return &diagnosesDataSource{}
}

// Data source implementation.
type diagnosesDataSource struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *scpci.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name
func (d *diagnosesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection_diagnoses"
}

// Schema defines the schema for the data source
func (d *diagnosesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Config inspection diagnosis result list",
		Attributes: map[string]schema.Attribute{
			// Input attributes
			common.ToSnakeCase("WithCount"): schema.StringAttribute{
				Description: "With count\n" +
					"  - Example : true",
				Optional: true,
			},
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit\n" +
					"  - Example : 20",
				Optional: true,
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker\n" +
					"  - Example : 607e0938521643b5b4b266f34fae693",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort\n" +
					"  - Example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account Id\n" +
					"  - Example : 0e3dffc50eb247a1adf4f2e5c82c4f99",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisID"): schema.StringAttribute{
				Description: "Id of diagnosis\n" +
					"  - Example : DIA-943731CB8E3045C289BAECAEC3532097",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis\n" +
					"  - Example : Sample Diagnosis Name",
				Optional: true,
			},
			common.ToSnakeCase("CSPType"): schema.StringAttribute{
				Description: "Type of cloud service provider\n" +
					"  - Example : SCP",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisState"): schema.StringAttribute{
				Description: "Status of diagnosis\n" +
					"  - Example : Completed",
				Optional: true,
			},
			common.ToSnakeCase("StartDate"): schema.StringAttribute{
				Description: "Start date\n" +
					"  - Example : 2022-01-01",
				Optional: true,
			},
			common.ToSnakeCase("EndDate"): schema.StringAttribute{
				Description: "End date\n" +
					"  - Example : 2022-12-31",
				Optional: true,
			},
			common.ToSnakeCase("UserId"): schema.StringAttribute{
				Description: "User id\n" +
					"  - Example : 76b563a009584b1380715c00703a02ff",
				Optional: true,
			},

			// Output attributes
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count\n" +
					"  - Example : 20",
				Computed: true,
			},
			common.ToSnakeCase("DiagnosisResultResponses"): schema.ListNestedAttribute{
				Description: "Diagnosis result responses",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CountCheck"): schema.Int32Attribute{
							Description: "Check count\n" +
								"  - Example : 10",
							Computed: true,
						},
						common.ToSnakeCase("CountError"): schema.Int32Attribute{
							Description: "Error count\n" +
								"  - Example : 0",
							Computed: true,
						},
						common.ToSnakeCase("CountFail"): schema.Int32Attribute{
							Description: "Fail count\n" +
								"  - Example : 3",
							Computed: true,
						},
						common.ToSnakeCase("CountNa"): schema.Int32Attribute{
							Description: "Na count\n" +
								"  - Example : 2",
							Computed: true,
						},
						common.ToSnakeCase("CountPass"): schema.Int32Attribute{
							Description: "Pass count\n" +
								"  - Example : 5",
							Computed: true,
						},
						common.ToSnakeCase("CspType"): schema.StringAttribute{
							Description: "Type of cloud service provider\n" +
								"  - Example : SCP",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
							Description: "Account Id of diagnosis\n" +
								"  - Example : 0e3dffc50eb247a1adf4f2e5c82c4f99",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
							Description: "Check type of diagnosis\n" +
								"  - Example : BP",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
							Description: "Id of diagnosis\n" +
								"  - Example : DIA-943731CB8E3045C289BAECAEC3532097",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
							Description: "Name of diagnosis\n" +
								"  - Example : Sample Diagnosis Name",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisRequestSequence"): schema.StringAttribute{
							Description: "Sequence of diagnosis request\n" +
								"  - Example : 1234567890",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisResult"): schema.StringAttribute{
							Description: "Diagnosis Result\n" +
								"  - Example : SUCCESS",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisTotalCount"): schema.Int32Attribute{
							Description: "Diagnosis Total Count\n" +
								"  - Example : 10",
							Computed: true,
						},
						common.ToSnakeCase("ProceedDate"): schema.StringAttribute{
							Description: "Proceed Date\n" +
								"  - Example : 2022-01-01T12:00:00Z",
							Computed: true,
						},
						common.ToSnakeCase("Total"): schema.Int32Attribute{
							Description: "Total count\n" +
								"  - Example : 10",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Links",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Href"): schema.StringAttribute{
							Description: "Href\n" +
								"  - Example : /api/v1/config-inspection/diagnoses?limit=20&marker=607e0938521643b5b4b266f34fae693",
							Computed: true,
						},
						common.ToSnakeCase("Rel"): schema.StringAttribute{
							Description: "Rel\n" +
								"  - Example : self",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure prepares the data source with the provider configuration
func (d *diagnosesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *diagnosesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scpci.ConfigInspectionDiagnosisResultListDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch data from the API
	res, err := d.client.GetConfigInspectionDiagnosisResultList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch config inspection result list", err.Error())
		return
	}

	// Map response data to state
	state.TotalCount = types.Int32PointerValue(res.Count.Get())

	// Handle DiagnosisResultResponses conversion
	if res.DiagnosisResultResponses != nil {
		diagnosisResults := make([]scpci.DiagnosisResultResponse, len(res.DiagnosisResultResponses))
		for i, result := range res.DiagnosisResultResponses {
			diagnosisResults[i] = scpci.DiagnosisResultResponse{
				CountCheck:               types.Int32PointerValue(result.CountCheck),
				CountError:               types.Int32PointerValue(result.CountError),
				CountFail:                types.Int32PointerValue(result.CountFail),
				CountNa:                  types.Int32PointerValue(result.CountNa),
				CountPass:                types.Int32PointerValue(result.CountPass),
				CspType:                  types.StringPointerValue(result.CspType),
				DiagnosisAccountId:       types.StringPointerValue(result.DiagnosisAccountId),
				DiagnosisCheckType:       types.StringPointerValue(result.DiagnosisCheckType),
				DiagnosisId:              types.StringPointerValue(result.DiagnosisId),
				DiagnosisName:            types.StringPointerValue(result.DiagnosisName),
				DiagnosisRequestSequence: types.StringPointerValue(result.DiagnosisRequestSequence),
				DiagnosisResult:          types.StringPointerValue(result.DiagnosisResult),
				DiagnosisTotalCount:      types.Int32PointerValue(result.DiagnosisTotalCount),
				ProceedDate:              types.StringValue(result.ProceedDate.Format(time.RFC3339)),
				Total:                    types.Int32PointerValue(result.Total),
			}
		}
		state.DiagnosisResultResponses = diagnosisResults
	}

	// Handle Links conversion
	if len(res.Links) > 0 {
		for _, element := range res.Links {
			var href, rel string

			if elemMap, ok := element.(map[string]interface{}); ok {
				if h, exists := elemMap["href"]; exists {
					if hStr, ok := h.(string); ok {
						href = hStr
					}
				}
				if r, exists := elemMap["rel"]; exists {
					if rStr, ok := r.(string); ok {
						rel = rStr
					}
				}
				state.Links = append(state.Links, scpci.Link{
					Href: types.StringValue(href),
					Rel:  types.StringValue(rel),
				})
			}
		}
	}

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
