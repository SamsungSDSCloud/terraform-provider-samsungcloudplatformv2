package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
		Description: "Provides details about a specific hosted zone.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the hosted zone to query.\n" +
					"  - example : 3432012nfdksdf03ktrld9234lgfg ",
				Optional: true,
			},
			common.ToSnakeCase("HostedZoneDetail"): schema.SingleNestedAttribute{
				Description: "Detailed information about the hosted zone.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Optional: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("HostedZoneType"): schema.StringAttribute{
						Description: "The type of the hosted zone (e.g., public or private).\n" +
							"  - example : private ",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the hosted zone.\n" +
							"  - example : 3432012nfdksdf03ktrld9234lgfg ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The domain name that a DNS service manages. all DNS records for that domain and its sub‑domains are stored and served within this hosted zone.\n" +
							"  - example : my-zone.com ",
						Optional: true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "The resource pool identifier associated with the hosted zone.\n" +
							"  - example : 10fjksdpooliddfsi12389esfdslkdsr32 ",
						Optional: true,
					},
					common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
						Description: "The DNS server ID for registering a Hosted Zone.For a Public‑type Hosted Zone, display it as an empty value.\n" +
							"  - example : 10fjkewefprivatedns3193rud543 ",
						Optional: true,
					},
					common.ToSnakeCase("PrivateDnsName"): schema.StringAttribute{
						Description: "The DNS server name for registering a Hosted Zone.For a Public‑type Hosted Zone, display it as an empty value.\n" +
							"  - example : private-dns01 ",
						Optional: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the hosted zone (e.g., ACTIVE, CREATING, DELETING).\n" +
							"  - example : ACTIVE ",
						Optional: true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "The Time-To-Live (TTL) value in seconds for DNS records in this zone.\n" +
							"  - example : 3600 ",
						Optional: true,
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

	hostedZoneState := convertHostedZoneShowResponseV1Dot3ToHostedZone(*data)

	hostedZoneObjectValue, _ := types.ObjectValueFrom(ctx, hostedZoneState.AttributeTypes(), hostedZoneState)
	state.HostedZoneDetail = hostedZoneObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
