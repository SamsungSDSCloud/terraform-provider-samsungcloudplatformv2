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
	_ datasource.DataSource              = &diagnosisDataSource{}
	_ datasource.DataSourceWithConfigure = &diagnosisDataSource{}
)

// Helper function to simplify the provider implementation.
func NewConfigInspectionDiagnosisDataSource() datasource.DataSource {
	return &diagnosisDataSource{}
}

// Data source implementation.
type diagnosisDataSource struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *scpci.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name
func (d *diagnosisDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection_diagnosis"
}

// Schema defines the schema for the data source
func (d *diagnosisDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Config inspection diagnosis result detail.",
		Attributes: map[string]schema.Attribute{
			// Input attributes
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis.\n" +
					"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisRequestSequence"): schema.StringAttribute{
				Description: "Sequence of diagnosis request.\n" +
					"  - example : 'SCPCIS-E75FD21CA524441C9C1B1B381D5974F7'",
				Required: true,
			},
			common.ToSnakeCase("WithCount"): schema.StringAttribute{
				Description: "Whether to include the total item count in the response.\n" +
					"  - example : true",
				Optional: true,
			},
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Maximum number of items to return per page.\n" +
					"  - example : 20",
				Optional: true,
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Pagination token from a previous response to fetch the next page.\n" +
					"  - example : '607e0938521643b5b4b266f34fae693'",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : 'created_at:desc'",
				Optional: true,
			},

			// Output attributes
			common.ToSnakeCase("ChecklistName"): schema.StringAttribute{
				Description: "Name of the checklist used for the diagnosis.",
				Computed:    true,
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total number of items available across all pages.\n" +
					"  - example : 20",
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
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis.\n" +
					"  - example : 'Sample Diagnosis Name'\n" +
					"  - pattern : `^[a-zA-Z0-9-_]+$`",
				Computed: true,
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Collection of hypermedia links to related resources or pages.\n" +
					"  - example : [{\"href\": \"http://scp.samsungsdscloud.com/v1/notices\", \"rel\": \"self\"}]",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Href"): schema.StringAttribute{
							Description: "URL of the linked resource.\n" +
								"  - example : 'http://scp.samsungsdscloud.com/v1/notices'",
							Computed: true,
						},
						common.ToSnakeCase("Rel"): schema.StringAttribute{
							Description: "Relationship type of the link.\n" +
								"  - example : 'self'",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("ProceedDate"): schema.StringAttribute{
				Description: "Date the diagnosis was performed.\n" +
					"  - example : '2022-01-01 12:00:00'",
				Computed: true,
			},
			common.ToSnakeCase("ResultDetailList"): schema.ListNestedAttribute{
				Description: "Result detail list",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ActionGuide"): schema.StringAttribute{
							Description: "Measure guide description.\n" +
								"  - example : 'Delete security group rules that allow access from ANY (0.0.0.0/0) or overly permissive source IPs to principle-based access ports, as shown below.'",
							Computed: true,
						},
						common.ToSnakeCase("Changed"): schema.BoolAttribute{
							Description: "Is changed?\n" +
								"  - example : true",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
							Description: "Check type of diagnosis.\n" +
								"  - example : 'BP'\n" +
								"  - enum : BP | SSI",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCriteria"): schema.StringAttribute{
							Description: "Decision standard description.\n" +
								"  - example : '[Security Group Rule]\nThere should be no rules that allow any IP access to principle-based access ports.'",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisItem"): schema.StringAttribute{
							Description: "Sub category description.\n" +
								"  - example : '2.NW_003. Principle-based access ports for each protocol should specify the required IPs and allow access only to them.'",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisLayer"): schema.StringAttribute{
							Description: "Inspector item category description.\n" +
								"  - example : '2.NETWORK'",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisMethod"): schema.StringAttribute{
							Description: "Inspector method description.\n" +
								"  - example : 'Check that security group inbound rules allow access only to users or systems that require principle-based access.'",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisResult"): schema.StringAttribute{
							Description: "Overall diagnosis execution status.\n" +
								"  - example : '03'",
							Computed: true,
						},
						common.ToSnakeCase("ResultContents"): schema.StringAttribute{
							Description: "Detailed finding of the check result.\n" +
								"  - example : 'Sample Result Contents'",
							Computed: true,
						},
						common.ToSnakeCase("SubCategory"): schema.StringAttribute{
							Description: "Sub category.\n" +
								"  - example : 'NURIBP_SCP_02.NW_004'",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("Total"): schema.Int32Attribute{
				Description: "Total number of diagnosis result items.",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the data source with the provider configuration
func (d *diagnosisDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *diagnosisDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scpci.ConfigInspectionDiagnosisResultDetailDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch data from the API
	res, err := d.client.GetConfigInspectionDiagnosisResultDetail(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch config inspection result detail", err.Error())
		return
	}

	// Map response data to state
	state.ChecklistName = types.StringValue(res.ChecklistName)
	state.DiagnosisAccountId = types.StringValue(res.DiagnosisAccountId)
	state.DiagnosisCheckType = types.StringValue(res.DiagnosisCheckType)
	state.DiagnosisName = types.StringValue(res.DiagnosisName)
	state.ProceedDate = types.StringValue(res.ProceedDate.Format(time.RFC3339))
	state.Total = types.Int32Value(res.Total)

	// Handle ResultDetailList conversion
	if res.ResultDetailList != nil {
		resultDetails := make([]scpci.DiagnosisResultDetail, len(res.ResultDetailList))
		for i, detail := range res.ResultDetailList {
			resultDetails[i] = scpci.DiagnosisResultDetail{
				ActionGuide:        detail.ActionGuide,
				DiagnosisCheckType: detail.DiagnosisCheckType,
				DiagnosisCriteria:  detail.DiagnosisCriteria,
				DiagnosisItem:      detail.DiagnosisItem,
				DiagnosisLayer:     detail.DiagnosisLayer,
				DiagnosisMethod:    detail.DiagnosisMethod,
				DiagnosisResult:    detail.DiagnosisResult,
				ResultContents:     detail.ResultContents,
			}
			if detail.Changed.IsSet() {
				resultDetails[i].Changed = types.BoolPointerValue(detail.Changed.Get())
			}
			if detail.SubCategory.IsSet() {
				resultDetails[i].SubCategory = types.StringPointerValue(detail.SubCategory.Get())
			}
		}
		state.ResultDetailList = resultDetails
	}

	if res.Count.IsSet() {
		state.TotalCount = types.Int32PointerValue(res.Count.Get())
	}

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
