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
	_ datasource.DataSource              = &configinspectionsDataSources{}
	_ datasource.DataSourceWithConfigure = &configinspectionsDataSources{}
)

// Helper function to simplify the provider implementation.
func NewConfigInspectionConfigInspectionDataSources() datasource.DataSource {
	return &configinspectionsDataSources{}
}

// Data source implementation.
type configinspectionsDataSources struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *scpci.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name
func (d *configinspectionsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection_configinspections"
}

// Schema defines the schema for the data source
func (d *configinspectionsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of config inspection object.",
		Attributes: map[string]schema.Attribute{
			// Input attributes
			common.ToSnakeCase("WithCount"): schema.StringAttribute{
				Description: "With count\n" +
					"  - Example: true",
				Optional: true,
			},
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit\n" +
					"  - Example: 20",
				Optional: true,
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker\n" +
					"  - Example: 607e0938521643b5b4b266f343fae693",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort\n" +
					"  - Example: created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("IsMine"): schema.BoolAttribute{
				Description: "My Config Inspection\n" +
					"  - Example: false",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis\n" +
					"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis\n" +
					"  - Example: My Diagnosis",
				Optional: true,
			},
			common.ToSnakeCase("CspType"): schema.StringAttribute{
				Description: "Type of cloud service provider\n" +
					"  - Example: SCP",
				Optional: true,
			},
			common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
				Description: "Id of diagnosis\n" +
					"  - Example: 0e3dffc50eb247a1adf4f2e5c82c4f99",
				Optional: true,
			},
			common.ToSnakeCase("RecentDiagnosisState"): schema.ListAttribute{
				Description: "Recent diagnosis state\n" +
					"  - Example: Completed",
				ElementType: types.StringType,
				Optional:    true,
			},
			common.ToSnakeCase("StartDate"): schema.StringAttribute{
				Description: "Start date\n" +
					"  - Example: 2022-01-01 12:00:00",
				Optional: true,
			},
			common.ToSnakeCase("EndDate"): schema.StringAttribute{
				Description: "End date\n" +
					"  - Example: 2022-01-02 12:00:00",
				Optional: true,
			},

			// Output attributes
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count\n" +
					"  - Example: 20",
				Computed: true,
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Links\n" +
					"  - Example: [{\"href\": \"http://scp.samsungsdscloud.com/v1/notices\", \"rel\": \"self\"}]",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Href"): schema.StringAttribute{
							Description: "Href\n" +
								"  - Example : http://scp.samsungsdscloud.com/v1/notices",
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
			common.ToSnakeCase("SummaryResponses"): schema.ListNestedAttribute{
				Description: "Summary responses",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created date\n" +
								"  - Example: 2022-01-01T12:00:00Z",
							Computed: true,
						},
						common.ToSnakeCase("CspType"): schema.StringAttribute{
							Description: "Type of cloud service provider\n" +
								"  - Example: SCP",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
							Description: "Id of diagnosis\n" +
								"  - Example: 0e3dffc50eb247a1adf4f2e5c82c4f99",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
							Description: "Check type of diagnosis\n" +
								"  - Example: BP",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
							Description: "Id of diagnosis\n" +
								"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
							Description: "Name of diagnosis\n" +
								"  - Example: Sample Diagnosis Name",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisType"): schema.StringAttribute{
							Description: "diagnosis Type\n" +
								"  - Example: Console",
							Computed: true,
						},
						common.ToSnakeCase("PlanType"): schema.StringAttribute{
							Description: "plan Type\n" +
								"  - Example: STANDARD",
							Computed: true,
						},
						common.ToSnakeCase("ErrorState"): schema.StringAttribute{
							Description: "Error type of recent diagnosis\n" +
								"  - Example: CONNECTION_FAIL",
							Computed: true,
						},
						common.ToSnakeCase("RecentDiagnosisAt"): schema.StringAttribute{
							Description: "Recent Diagnosis Date\n" +
								"  - Example: 2022-01-01T12:00:00Z",
							Computed: true,
						},
						common.ToSnakeCase("RecentDiagnosisState"): schema.StringAttribute{
							Description: "Recent Diagnosis State\n" +
								"  - Example: Completed",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure prepares the data source with the provider configuration
func (d *configinspectionsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *configinspectionsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scpci.ConfigInspectionDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch data from the API
	res, err := d.client.GetConfigInspectionList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch config inspection list", err.Error())
		return
	}

	// Map API response to the state object
	state.TotalCount = types.Int32PointerValue(res.Count.Get())
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
	if len(res.SummaryResponses) > 0 {
		for _, element := range res.SummaryResponses {
			summaryResponse := scpci.SummaryResponse{
				CreatedAt:            types.StringValue(element.CreatedAt.Format(time.RFC3339)),
				CspType:              types.StringValue(element.CspType),
				DiagnosisAccountId:   types.StringValue(element.DiagnosisAccountId),
				DiagnosisCheckType:   types.StringValue(element.DiagnosisCheckType),
				DiagnosisId:          types.StringValue(element.DiagnosisId),
				DiagnosisName:        types.StringValue(element.DiagnosisName),
				DiagnosisType:        types.StringValue(element.DiagnosisType),
				PlanType:             types.StringValue(element.PlanType),
				ErrorState:           types.StringPointerValue(element.ErrorState.Get()),
				RecentDiagnosisState: types.StringPointerValue(element.RecentDiagnosisState.Get()),
			}
			if element.RecentDiagnosisAt.IsSet() {
				summaryResponse.RecentDiagnosisAt = types.StringValue(element.RecentDiagnosisAt.Get().Format(time.RFC3339))
			}

			state.SummaryResponses = append(state.SummaryResponses, summaryResponse)
		}
	}

	// Map other fields similarly
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
