package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/dns/1.3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsPrivateDnsDataSources{}
	_ datasource.DataSourceWithConfigure = &dnsPrivateDnsDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsPrivateDnsDataSources() datasource.DataSource {
	return &dnsPrivateDnsDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsPrivateDnsDataSources struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsPrivateDnsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_private_dnss" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsPrivateDnsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Provides a list of private DNS instances.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20 ",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0 ",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:asc ",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The identifier to filter private DNS by.\n" +
					"  - example : 10fjkewefprivatedns3193rud543 ",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name to filter private DNS by.\n" +
					"  - example : private-dns01 ",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The VPC identifier to filter private DNS by.\n" +
					"  - example : vpc-12345678 ",
				Optional: true,
			},

			common.ToSnakeCase("PrivateDns"): schema.ListNestedAttribute{
				Description: "List of private DNS instances matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsPrivateDnsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsPrivateDnsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.PrivateDnsDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPrivateDnsList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Gslbs",
			err.Error(),
		)
		return
	}

	for _, privateDns := range data.PrivateDns {
		privateDnsState := convertPrivateDns(privateDns)

		state.PrivateDns = append(state.PrivateDns, privateDnsState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func convertPrivateDns(privateDns scpdns.PrivateDns) dns.PrivateDns {

	connectedVpcIds := make([]types.String, len(privateDns.ConnectedVpcIds))

	for idx, connectedVpcId := range privateDns.ConnectedVpcIds {
		connectedVpcIds[idx] = types.StringValue(connectedVpcId)
	}

	return dns.PrivateDns{
		AuthDnsName:      types.StringValue(privateDns.GetAuthDnsName()),
		ConnectedVpcIds:  connectedVpcIds,
		CreatedAt:        types.StringValue(privateDns.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(privateDns.CreatedBy),
		Description:      types.StringValue(privateDns.GetDescription()),
		Id:               types.StringValue(privateDns.Id),
		ModifiedAt:       types.StringValue(privateDns.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(privateDns.ModifiedBy),
		Name:             types.StringValue(privateDns.Name),
		PoolId:           types.StringValue(privateDns.GetPoolId()),
		PoolName:         types.StringValue(privateDns.GetPoolName()),
		RegisteredRegion: types.StringValue(privateDns.GetRegisteredRegion()),
		ResolverIp:       types.StringValue(privateDns.GetResolverIp()),
		ResolverName:     types.StringValue(privateDns.GetResolverName()),
		State:            types.StringValue(string(privateDns.State)),
	}
}
