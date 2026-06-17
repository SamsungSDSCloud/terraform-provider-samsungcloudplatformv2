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
	_ datasource.DataSource              = &dnsPrivateDnsDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsPrivateDnsDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsPrivateDnsDataSource() datasource.DataSource {
	return &dnsPrivateDnsDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsPrivateDnsDataSource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsPrivateDnsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_private_dns" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsPrivateDnsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Provides details about a specific private DNS.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the private DNS to query.\n" +
					"  - example : 10fjkewefprivatedns3193rud543 ",
				Optional: true,
			},
			common.ToSnakeCase("PrivateDnsDetail"): schema.SingleNestedAttribute{
				Description: "Detailed information about the private DNS.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AuthDnsName"): schema.StringAttribute{
						Description: "The authoritative DNS name of the private DNS.\n" +
							"  - example : auth.dns.example.com ",
						Optional: true,
					},
					common.ToSnakeCase("ConnectedVpcIds"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of VPC identifiers connected to this private DNS.Only VPCs that are connected to the DNS can query the domain information registered in it.\n" +
							"  - example : ['vpc-12345678', 'vpc-87654321'] ",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
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
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the private DNS.\n" +
							"  - example : 10fjkewefprivatedns3193rud543 ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the private DNS.\n" +
							"  - example : private-dns01 ",
						Optional: true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "The resource pool identifier associated with the private DNS.\n" +
							"  - example : 10fjksdpooliddfsi12389esfdslkdsr32 ",
						Optional: true,
					},
					common.ToSnakeCase("PoolName"): schema.StringAttribute{
						Description: "The name of the resource pool.\n" +
							"  - example : pool-01 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisteredRegion"): schema.StringAttribute{
						Description: "The region where the private DNS is registered.\n" +
							"  - example : KR-WEST1 ",
						Optional: true,
					},
					common.ToSnakeCase("ResolverIp"): schema.StringAttribute{
						Description: "The IP address of the DNS resolver.\n" +
							"  - example : 198.19.0.101 ",
						Optional: true,
					},
					common.ToSnakeCase("ResolverName"): schema.StringAttribute{
						Description: "The name of the DNS resolver.\n" +
							"  - example : resolver-01 ",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the private DNS.\n" +
							"  - example : ACTIVE ",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsPrivateDnsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsPrivateDnsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.PrivateDnsDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show PrivateDns",
			err.Error(),
		)
		return
	}

	dnsState := convertPrivateDns(data.PrivateDns)

	dnsObjectValue, _ := types.ObjectValueFrom(ctx, dnsState.AttributeTypes(), dnsState)
	state.PrivateDns = dnsObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
