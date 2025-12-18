package certificatemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificateManagerDetailDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateManagerDetailDataSource{}
)

// NewCertificateManagerDetailDataSource is a helper function to simplify the provider implementation.
func NewCertificateManagerDetailDataSource() datasource.DataSource {
	return &certificateManagerDetailDataSource{}
}

type certificateManagerDetailDataSource struct {
	config  *scpsdk.Configuration
	client  *certificatemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *certificateManagerDetailDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_manager"
}

// Schema defines the schema for the data source.
func (d *certificateManagerDetailDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Detail of certificate manager.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Certificate ID",
				Optional:    true,
			},
			common.ToSnakeCase("Certificate"): schema.SingleNestedAttribute{
				Description: "A Detail certificate.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CertKind"): schema.StringAttribute{
						Description: "Certificate type\n" +
							"  - Example: DEV",
						Computed: true,
					},
					common.ToSnakeCase("Cn"): schema.StringAttribute{
						Description: "Certificate Common Name\n" +
							"  - Example: test.go.kr",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID\n" +
							"  - Example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Certificate Name\n" +
							"  - Example: test-certificate",
						Computed: true,
					},
					common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
						Description: "Certificate Expire Date\n" +
							"  - Example: 2026-02-07T18:07:59",
						Computed: true,
					},
					common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
						Description: "Certificate Start Date\n" +
							"  - Example: 2025-02-08T18:07:00",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Certificate State\n" +
							"  - Example: VALID",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *certificateManagerDetailDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.CertificateManager
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *certificateManagerDetailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificatemanager.CertificateManagerDetailDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetCertificateManager(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading internet gateway",
			"Could not read internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	certificate := certificatemanager.Certificate{
		Id:          types.StringValue(data.Certificate.Id),
		Name:        types.StringValue(data.Certificate.Name),
		CertKind:    types.StringValue(*data.Certificate.CertKind),
		Cn:          types.StringValue(data.Certificate.Cn),
		NotBeforeDt: types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		NotAfterDt:  types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		State:       types.StringValue(data.Certificate.State),
	}

	certificateObjectValue, _ := types.ObjectValueFrom(ctx, certificate.AttributeTypes(), certificate)
	state.Certificate = certificateObjectValue
	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
