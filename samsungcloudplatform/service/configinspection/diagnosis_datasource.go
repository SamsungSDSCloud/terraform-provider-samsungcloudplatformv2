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
				Description: "Id of diagnosis\n" +
					"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisRequestSequence"): schema.StringAttribute{
				Description: "Sequence of diagnosis request\n" +
					"  - Example: SCPCIS-E75FD21CA524441C9C1B1B381D5974F7",
				Required: true,
			},
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

			// Output attributes
			common.ToSnakeCase("ChecklistName"): schema.StringAttribute{
				Description: "Checklist Name\n" +
					"  - Example: Sample Checklist",
				Computed: true,
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Count\n" +
					"  - Example: 20",
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
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis\n" +
					"  - Example: Sample Diagnosis Name",
				Computed: true,
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Links",
				Computed:    true,
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
			common.ToSnakeCase("ProceedDate"): schema.StringAttribute{
				Description: "Proceed Date",
				Computed:    true,
			},
			common.ToSnakeCase("ResultDetailList"): schema.ListNestedAttribute{
				Description: "Result detail list",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ActionGuide"): schema.StringAttribute{
							Description: "Measure guide description\n" +
								"  - Example: 원칙 접근 Port에 Source IP가 ANY(0.0.0.0/0)로 접근되어 있거나 과도하게 허용되는 Security Group 추가는 아래와 같이 삭제합니다.",
							Computed: true,
						},
						common.ToSnakeCase("Changed"): schema.BoolAttribute{
							Description: "Is changed?\n" +
								"  - Example: true",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
							Description: "Check type of diagnosis\n" +
								"  - Example: BP",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisCriteria"): schema.StringAttribute{
							Description: "Decision standard description\n" +
								"  - Example: 【 Security Group 추가 】\n" +
								"① 원칙접근 Port에 Any IP 접근을 허용하는 추가가 존재하지 않아야 합니다.",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisItem"): schema.StringAttribute{
							Description: "Sub category description\n" +
								"  - Example: 2.NW_003. 프로토콜 별 원칙접근 Port는 접근이 필요한 IP를 지정하여 접근을 허용해야 합니다.",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisLayer"): schema.StringAttribute{
							Description: "Inspector item category description\n" +
								"  - Example: 2.NETWORK",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisMethod"): schema.StringAttribute{
							Description: "Inspector method description\n" +
								"  - Example: Security Group의 Inbound 추가에 원칙 접근이 필요한 사용자 또는 시스템만 접근을 허용하는 추가를 확인합니다.",
							Computed: true,
						},
						common.ToSnakeCase("DiagnosisResult"): schema.StringAttribute{
							Description: "Verify state\n" +
								"  - Example: 03",
							Computed: true,
						},
						common.ToSnakeCase("ResultContents"): schema.StringAttribute{
							Description: "Result Contents\n" +
								"  - Example: 상세 내용",
							Computed: true,
						},
						common.ToSnakeCase("SubCategory"): schema.StringAttribute{
							Description: "Sub category\n" +
								"  - Example: NURIBP_SCP_02.NW_004",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("Total"): schema.Int32Attribute{
				Description: "Total\n" +
					"  - Example: 10",
				Computed: true,
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
