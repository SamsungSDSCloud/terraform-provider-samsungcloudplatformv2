package configinspection

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	scpci "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	configinspection "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/configinspection/1.1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &diagnosisRequestDataSource{}
	_ datasource.DataSourceWithConfigure = &diagnosisRequestDataSource{}
)

// Helper function to simplify the provider implementation.
func NewConfigInspectionDiagnosisRequestDataSource() datasource.DataSource {
	return &diagnosisRequestDataSource{}
}

// Data source implementation.
type diagnosisRequestDataSource struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *scpci.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name
func (d *diagnosisRequestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection_diagnosis_request"
}

// Schema defines the schema for the data source
func (d *diagnosisRequestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Config inspection diagnostic request.",
		Attributes: map[string]schema.Attribute{
			// Input attributes
			common.ToSnakeCase("AccessKey"): schema.StringAttribute{
				Description: "Your API access key.\n" +
					"  - example : 'b19a2ee194744c218xxxxxxxxxxxxxxx'",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
				Description: "Check type of diagnosis.\n" +
					"  - example : 'BP'\n" +
					"  - enum : BP | SSI",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis.\n" +
					"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
				Required: true,
			},
			common.ToSnakeCase("SecretKey"): schema.StringAttribute{
				Description: "Your API secret key.\n" +
					"  - example : 'SAMPLE KEY'",
				Required: true,
			},
			common.ToSnakeCase("TenantId"): schema.StringAttribute{
				Description: "Your tenant ID.\n" +
					"  - example : '1234567890'",
				Required: true,
			},

			// Output attributes
			common.ToSnakeCase("Result"): schema.BoolAttribute{
				Description: "Result of diagnosis request (true, false).\n" +
					"  - example : true",
				Computed: true,
			},
		},
	}
}

// Configure prepares the data source with the provider configuration
func (d *diagnosisRequestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *diagnosisRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scpci.ConfigInspectionDiagnosisRequestDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch data from the API
	res, err := d.client.RequestNewConfigInspectionDiagnosis(ctx, configinspection.DiagnosisRequest{
		AccessKey:          state.AccessKey.ValueString(),
		DiagnosisCheckType: state.DiagnosisCheckType.ValueString(),
		DiagnosisId:        state.DiagnosisId.ValueString(),
		SecretKey:          state.SecretKey.ValueString(),
		TenantId:           state.TenantId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to request config inspection diagnosis", err.Error())
		return
	}

	state.Result = types.BoolValue(res.Result)

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
