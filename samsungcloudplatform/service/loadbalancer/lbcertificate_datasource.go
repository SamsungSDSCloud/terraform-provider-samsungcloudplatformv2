package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbCertificateDataSource{}
)

// NewResourceManagerResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewLoadbalancerLbCertificateDataSource() datasource.DataSource {
	return &loadbalancerLbCertificateDataSource{}
}

// resourceManagerResourceGroupDataSource is the data source implementation.
type loadbalancerLbCertificateDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbCertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_certificate" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbCertificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific LB Certificate.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the LB Certificate.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Optional: true,
			},
			common.ToSnakeCase("LbCertificate"): schema.SingleNestedAttribute{
				Description: "Details of the LB Certificate.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("CertBody"): schema.StringAttribute{
						Description: "The certificate body in PEM format.\n" +
							"  - example : (sensitive value)\n",
						Optional:  true,
						Sensitive: true,
					},
					common.ToSnakeCase("CertChain"): schema.StringAttribute{
						Description: "The certificate chain in PEM format.\n" +
							"  - example : (sensitive value)\n",
						Optional:  true,
						Sensitive: true,
					},
					common.ToSnakeCase("CertKind"): schema.StringAttribute{
						Description: "The type of certificate.\n" +
							"  - example : SERVER\n" +
							"  - pattern : SERVER | CLIENT\n",
						Optional: true,
					},
					common.ToSnakeCase("Cn"): schema.StringAttribute{
						Description: "The common name (CN) of the certificate.\n" +
							"  - example : example.com\n",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the LB Certificate.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Certificate.\n" +
							"  - example : Certificate01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9._-]+$\n",
						Optional: true,
					},
					common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
						Description: "The expiration date of the certificate.\n" +
							"  - example : 2026-02-12T23:59:59Z\n",
						Optional: true,
					},
					common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
						Description: "The start date of the certificate validity.\n" +
							"  - example : 2025-02-12T00:00:00Z\n",
						Optional: true,
					},
					common.ToSnakeCase("Organization"): schema.StringAttribute{
						Description: "The organization name in the certificate.\n" +
							"  - example : Samsung SDS\n",
						Optional: true,
					},
					common.ToSnakeCase("PrivateKey"): schema.StringAttribute{
						Description: "The private key associated with the certificate.\n" +
							"  - example : (sensitive value)\n",
						Optional:  true,
						Sensitive: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the LB Certificate.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : ACTIVE | ERROR\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.LoadBalancer
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *loadbalancerLbCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadbalancer.LbCertificateDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbCertificate(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show LbCertificate",
			err.Error(),
		)
		return
	}

	lbCertificateState := loadbalancer.LbCertificateDetail{
		AccountId:    types.StringValue(data.Certificate.AccountId),
		CertBody:     types.StringValue(data.Certificate.CertBody),
		CertChain:    virtualserverutil.ToNullableStringValue(data.Certificate.CertChain.Get()),
		CertKind:     types.StringValue(*data.Certificate.CertKind),
		Cn:           types.StringValue(*data.Certificate.Cn),
		CreatedAt:    types.StringValue(data.Certificate.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(data.Certificate.CreatedBy),
		Id:           types.StringValue(data.Certificate.Id),
		ModifiedAt:   types.StringValue(data.Certificate.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(data.Certificate.ModifiedBy),
		Name:         types.StringValue(*data.Certificate.Name),
		NotAfterDt:   types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		NotBeforeDt:  types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		Organization: types.StringValue(data.Certificate.Organization),
		PrivateKey:   types.StringValue(data.Certificate.PrivateKey),
		State:        types.StringValue(*data.Certificate.State),
	}

	lbCertificateObjectValue, _ := types.ObjectValueFrom(ctx, lbCertificateState.AttributeTypes(), lbCertificateState)
	state.LbCertificateDetail = lbCertificateObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
