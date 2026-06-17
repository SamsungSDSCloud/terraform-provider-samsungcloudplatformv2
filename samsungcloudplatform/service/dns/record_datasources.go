package dns

import (
	"context"
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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsRecordDataSources{}
	_ datasource.DataSourceWithConfigure = &dnsRecordDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsRecordDataSources() datasource.DataSource {
	return &dnsRecordDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsRecordDataSources struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsRecordDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsRecordDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of DNS records.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "The maximum number of items to return.\n" +
					"  - example : 10 ",
				Optional: true,
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "The last record ID of the previous page.\n" +
					"  - example : 6ed7bc1-4b05-3cc7-7105-c1b71f7f30a7 ",
				Optional: true,
			},
			common.ToSnakeCase("SortDir"): schema.StringAttribute{
				Description: "The sort direction for the list (ASC or DESC).\n" +
					"  - example : ASC ",
				Optional: true,
			},
			common.ToSnakeCase("SortKey"): schema.StringAttribute{
				Description: "The field to sort by.\n" +
					"  - example : name ",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name to filter DNS records by.\n" +
					"  - example : test.app ",
				Optional: true,
			},
			common.ToSnakeCase("ExactName"): schema.StringAttribute{
				Description: "Filter by exact name match.\n" +
					"  - example : exact-record-name ",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type of the DNS record (e.g., A, AAAA, CNAME, MX, TXT, SPF).\n" +
					"  - example : A ",
				Optional: true,
			},
			common.ToSnakeCase("Data"): schema.StringAttribute{
				Description: "The data to filter DNS records by.\n" +
					"  - example : 192.168.1.1 ",
				Optional: true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "The status to filter DNS records by.\n" +
					"  - example : ACTIVE ",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : This is description ",
				Optional: true,
			},
			common.ToSnakeCase("Ttl"): schema.Int32Attribute{
				Description: "The Time-To-Live (TTL) value in seconds for the DNS record.\n" +
					"  - example : 3600 ",
				Optional: true,
			},
			common.ToSnakeCase("HostedZoneId"): schema.StringAttribute{
				Description: "ID for the zone that contains this record\n" +
					"  - example : 3432012nfdksdf03ktrld9234lgfg ",
				Optional: true,
			},
			common.ToSnakeCase("Records"): schema.ListNestedAttribute{
				Description: "List of DNS records matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Action"): schema.StringAttribute{
							Description: "The action performed on the DNS record.\n" +
								"  - example : NONE ",
							Optional: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z ",
							Optional: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
								"  - example : This is description ",
							Optional: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the DNS record.\n" +
                                "  - example : 6ed7bc1-4b05-3cc7-7105-c1b71f7f30a7 ",
							Optional: true,
						},
						common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
							Description: "The links related to the DNS record.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Self"): schema.StringAttribute{
									Description: "The self-referential link of the DNS record.\n" +
										"  - example : https://api.samsungsdscloud.com/dns/v1/records/3432012nfdksdf03ktrld9234lgfg ",
									Optional: true,
								},
							},
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the DNS record.\n" +
								"  - example : test.app ",
							Optional: true,
						},
						common.ToSnakeCase("ProjectId"): schema.StringAttribute{
							Description: "The project identifier associated with the DNS record.\n" +
								"  - example : 003dffc50eb123a1cbf4f2e5c71d4f15 ",
							Optional: true,
						},
						common.ToSnakeCase("Records"): schema.ListAttribute{
							ElementType: types.StringType,
							Description: "The list of record values for this DNS record.\n" +
								"  - example : [\"12.34.45.67\"]",
							Optional: true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "The current status of the DNS record.\n" +
								"  - example : ACTIVE ",
							Optional: true,
						},
						common.ToSnakeCase("Ttl"): schema.Int32Attribute{
							Description: "The Time-To-Live (TTL) value in seconds for the DNS record.\n" +
								"  - example : 3600 ",
							Optional: true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "The type of the DNS record (e.g., A, AAAA, CNAME, MX, TXT, SPF).\n" +
								"  - example : A ",
							Optional: true,
						},
						common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last updated, in ISO 8601 format.\n" +
								"  - example : 2026-02-09T08:00:40Z ",
							Optional: true,
						},
						common.ToSnakeCase("Version"): schema.Int32Attribute{
							Description: "The version of the DNS record.\n" +
								"  - example : 1 ",
							Optional: true,
						},
						common.ToSnakeCase("ZoneId"): schema.StringAttribute{
							Description: "The zone identifier of the DNS record.\n" +
								"  - example : 3432012nfdksdf03ktrld9234lgfg ",
							Optional: true,
						},
						common.ToSnakeCase("ZoneName"): schema.StringAttribute{
							Description: "The name of the zone containing the DNS record.\n" +
								"  - example : my-zone.com ",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsRecordDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsRecordDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.RecordDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRecordList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Records",
			err.Error(),
		)
		return
	}

	for _, record := range data.Records {
		recordState := convertRecord(record)

		state.Records = append(state.Records, recordState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func convertRecord(record scpdns.Record) dns.Record {
	var links *dns.Links
	if record.Links != nil {
		self, ok := record.Links["self"].(string)
		if ok {
			links = &dns.Links{
				Self: types.StringValue(self),
			}
		}
	}

	records := make([]types.String, len(record.Records))

	for idx, recordFromRes := range record.Records {
		records[idx] = types.StringValue(recordFromRes)
	}

	return dns.Record{
		Action:      virtualserverutil.ToNullableStringValue(record.Action.Get()),
		CreatedAt:   virtualserverutil.ToNullableStringValue(record.CreatedAt.Get()),
		Description: virtualserverutil.ToNullableStringValue(record.Description.Get()),
		Id:          virtualserverutil.ToNullableStringValue(record.Id.Get()),
		Links:       links,
		Name:        virtualserverutil.ToNullableStringValue(record.Name.Get()),
		ProjectId:   virtualserverutil.ToNullableStringValue(record.ProjectId.Get()),
		Records:     records,
		Status:      virtualserverutil.ToNullableStringValue(record.Status.Get()),
		Ttl:         common.ToNullableInt32Value(record.Ttl.Get()),
		Type:        virtualserverutil.ToNullableStringValue(record.Type.Get()),
		UpdatedAt:   virtualserverutil.ToNullableStringValue(record.UpdatedAt.Get()),
		Version:     common.ToNullableInt32Value(record.Version.Get()),
		ZoneId:      virtualserverutil.ToNullableStringValue(record.ZoneId.Get()),
		ZoneName:    virtualserverutil.ToNullableStringValue(record.ZoneName.Get()),
	}
}
