package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/dns/1.3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsHostedZoneDataSources{}
	_ datasource.DataSourceWithConfigure = &dnsHostedZoneDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsHostedZoneDataSources() datasource.DataSource {
	return &dnsHostedZoneDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsHostedZoneDataSources struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsHostedZoneDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_hosted_zones" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsHostedZoneDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of hosted zones.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0 ",
				Optional: true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20 ",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:asc ",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name to filter hosted zones by.\n" +
					"  - example : my-zone.com ",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type to filter hosted zones by (PUBLIC or PRIVATE).\n" +
					"  - example : private ",
				Optional: true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "The status to filter hosted zones by.\n" +
					"  - example : ACTIVE ",
				Optional: true,
			},
			common.ToSnakeCase("HostedZones"): schema.ListNestedAttribute{
				Description: "List of hosted zones matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsHostedZoneDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsHostedZoneDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.HostedZoneDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetHostedZoneList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HostedZones",
			err.Error(),
		)
		return
	}

	for _, hostedZone := range data.HostedZones {
		hostedZoneState := convertHostedZoneShowResponseV1Dot3ToHostedZone(convertHostedZoneV1Dot3ToHostedZoneShowResponseV1Dot3(hostedZone))

		state.HostedZones = append(state.HostedZones, hostedZoneState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func convertHostedZoneV1Dot3ToHostedZoneShowResponseV1Dot3(hostedZoneV1Dot3 scpdns.HostedZoneV1Dot3) scpdns.HostedZoneShowResponseV1Dot3 {
	hostedZoneShowResponseV1Dot3 := scpdns.HostedZoneShowResponseV1Dot3{}
	data, _ := json.Marshal(hostedZoneV1Dot3)
	json.Unmarshal(data, &hostedZoneShowResponseV1Dot3)
	return hostedZoneShowResponseV1Dot3
}

func convertHostedZoneShowResponseV1Dot3ToHostedZone(hostedZoneShowResponseV1Dot3 scpdns.HostedZoneShowResponseV1Dot3) dns.HostedZone {

	var hostedZoneType *string
	if hostedZoneShowResponseV1Dot3.HostedZoneType != nil {
		s := string(*hostedZoneShowResponseV1Dot3.HostedZoneType)
		hostedZoneType = &s
	} else {
		hostedZoneType = nil
	}

	return dns.HostedZone{
		CreatedAt:      types.StringValue(hostedZoneShowResponseV1Dot3.CreatedAt.Format(time.RFC3339)),
		CreatedBy:      types.StringValue(hostedZoneShowResponseV1Dot3.CreatedBy),
		Description:    virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot3.Description.Get()),
		HostedZoneType: virtualserverutil.ToNullableStringValue(hostedZoneType),
		Id:             types.StringValue(hostedZoneShowResponseV1Dot3.Id),
		ModifiedAt:     types.StringValue(hostedZoneShowResponseV1Dot3.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:     types.StringValue(hostedZoneShowResponseV1Dot3.ModifiedBy),
		Name:           types.StringValue(hostedZoneShowResponseV1Dot3.Name),
		PoolId:         types.StringValue(hostedZoneShowResponseV1Dot3.PoolId),
		PrivateDnsId:   virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot3.PrivateDnsId.Get()),
		PrivateDnsName: virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot3.PrivateDnsName.Get()),
		Status:         types.StringValue(string(hostedZoneShowResponseV1Dot3.Status)),
		Ttl:            common.ToNullableInt32Value(hostedZoneShowResponseV1Dot3.Ttl.Get()),
	}
}
