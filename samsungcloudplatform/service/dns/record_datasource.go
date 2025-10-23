package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/dns/1.1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsRecordDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsRecordDataSource() datasource.DataSource {
	return &dnsRecordDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsRecordDataSource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "show record.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("HostedZoneId"): schema.StringAttribute{
				Description: "HostedZoneId",
				Optional:    true,
			},
			common.ToSnakeCase("RecordDetail"): schema.SingleNestedAttribute{
				Description: "A detail of Record.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
						Description: "Links",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Self"): schema.StringAttribute{
								Description: "Self",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("ProjectId"): schema.StringAttribute{
						Description: "ProjectId",
						Optional:    true,
					},
					common.ToSnakeCase("Records"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "Records",
						Optional:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
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
					common.ToSnakeCase("ZoneId"): schema.StringAttribute{
						Description: "ZoneId",
						Optional:    true,
					},
					common.ToSnakeCase("ZoneName"): schema.StringAttribute{
						Description: "ZoneName",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsRecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.RecordDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Record",
			err.Error(),
		)
		return
	}

	recordState := convertRecordDetail(*data)

	recordObjectValue, _ := types.ObjectValueFrom(ctx, recordState.AttributeTypes(), recordState)
	state.RecordDetail = recordObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func convertRecordDetail(record scpdns.RecordShowResponse) dns.RecordDetail {
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

	return dns.RecordDetail{
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

func convertRecordCreateResponseToRecord(a scpdns.RecordCreateResponse) scpdns.RecordShowResponse {
	b := scpdns.RecordShowResponse{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}

func convertRecordSetResponseToRecord(a scpdns.RecordSetResponse) scpdns.RecordShowResponse {
	b := scpdns.RecordShowResponse{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}
