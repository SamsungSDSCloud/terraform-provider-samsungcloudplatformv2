package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsHostedZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsHostedZoneDataSource{}
)

func NewDnsHostedZoneDataSource() datasource.DataSource {
	return &dnsHostedZoneDataSource{}
}

type dnsHostedZoneDataSource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsHostedZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_hosted_zone" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsHostedZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show HostedZone.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("HostedZoneDetail"): schema.SingleNestedAttribute{
				Description: "A detail of HostedZone.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action",
						Optional:    true,
					},
					common.ToSnakeCase("Attributes"): schema.SingleNestedAttribute{
						Description: "Attributes",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ServiceTier"): schema.StringAttribute{
								Description: "ServiceTier",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Email"): schema.StringAttribute{
						Description: "Email",
						Optional:    true,
					},
					common.ToSnakeCase("HostedZoneType"): schema.StringAttribute{
						Description: "HostedZoneType",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
						Description: "Links",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Self"): schema.StringAttribute{
								Description: "Self",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("Masters"): schema.ListAttribute{
						Description: "Masters",
						Optional:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "PoolId",
						Optional:    true,
					},
					common.ToSnakeCase("ProjectId"): schema.StringAttribute{
						Description: "ProjectId",
						Optional:    true,
					},
					common.ToSnakeCase("Serial"): schema.Int32Attribute{
						Description: "Serial",
						Optional:    true,
					},
					common.ToSnakeCase("Shared"): schema.BoolAttribute{
						Description: "Shared",
						Optional:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Optional:    true,
					},
					common.ToSnakeCase("TransferredAt"): schema.StringAttribute{
						Description: "TransferredAt",
						Optional:    true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "Ttl",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
					common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
						Description: "UpdatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("Version"): schema.Int32Attribute{
						Description: "Version",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsHostedZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Dns
	d.clients = inst.Client
}

func (d *dnsHostedZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.HostedZoneDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show HostedZone",
			err.Error(),
		)
		return
	}

	hostedZoneState := convertHostedZone(convertHostedZoneShowResponseToHostedZone(*data))

	hostedZoneObjectValue, _ := types.ObjectValueFrom(ctx, hostedZoneState.AttributeTypes(), hostedZoneState)
	state.HostedZoneDetail = hostedZoneObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
