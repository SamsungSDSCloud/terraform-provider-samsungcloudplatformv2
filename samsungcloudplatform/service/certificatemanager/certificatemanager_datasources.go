package certificatemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificateManagerDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateManagerDataSource{}
)

// NewCertificateManagerDataSource is a helper function to simplify the provider implementation.
func NewCertificateManagerDataSource() datasource.DataSource {
	return &certificateManagerDataSource{}
}

type certificateManagerDataSource struct {
	config  *scpsdk.Configuration
	client  *certificatemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *certificateManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_managers"
}

// Schema defines the schema for the data source.
func (d *certificateManagerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of certificate managers.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "id",
				Optional:    true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("IsMine"): schema.BoolAttribute{
				Description: "IsMine",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("Cn"): schema.StringAttribute{
				Description: "Cn",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "state",
				Optional:    true,
			},
			common.ToSnakeCase("Certificates"): schema.ListNestedAttribute{
				Description: "A list certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CertKind"): schema.StringAttribute{
							Description: "Certificate type\n" +
								"  - Example: PRD",
							Computed: true,
						},
						common.ToSnakeCase("Cn"): schema.StringAttribute{
							Description: "Certificate Common Name\n" +
								"  - Example: test.go.kr",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Certificate ID",
							Computed:    true,
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *certificateManagerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *certificateManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificatemanager.CertificateManagerDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	fmt.Printf("--------------------Start CALL GetCertificateManagerList")

	data, err := d.client.GetCertificateManagerList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading internet gateway",
			"Could not read internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	fmt.Printf("--------------------GetCertificateManagerList %v\n", data)
	fmt.Printf("--------------------GetCertificateManagerList %#v\n", data)

	// Map response body to model
	for _, igw := range data.Certificates {

		vps := certificatemanager.Certificate{
			Id:          types.StringValue(igw.Id),
			Name:        types.StringValue(igw.Name),
			CertKind:    types.StringValue(*igw.CertKind),
			Cn:          types.StringValue(igw.Cn),
			NotBeforeDt: types.StringValue(igw.NotBeforeDt.Format(time.RFC3339)),
			NotAfterDt:  types.StringValue(igw.NotAfterDt.Format(time.RFC3339)),
			State:       types.StringValue(igw.State),
		}

		state.Certificates = append(state.Certificates, vps)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
