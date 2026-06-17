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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsPublicDomainNameDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsPublicDomainNameDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsPublicDomainNameDataSource() datasource.DataSource {
	return &dnsPublicDomainNameDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsPublicDomainNameDataSource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsPublicDomainNameDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_public_domain_name" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsPublicDomainNameDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides details about a specific public domain name.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the public domain name.\n" +
					"  - example : 125jkdkt5fpublicdomain3193rud546 ",
				Optional: true,
			},
			common.ToSnakeCase("PublicDomainNameDetail"): schema.SingleNestedAttribute{
				Description: "Detailed information about the public domain name.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "The type of address for the domain registration.\n" +
							"  - example : DOMESTIC ",
						Computed: true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "Indicates whether automatic extension is enabled for the domain.\n" +
							"  - example : true ",
						Computed: true,
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
						Computed: true,
					},
					common.ToSnakeCase("DomesticAddressEn"): schema.StringAttribute{
						Description: "The domestic address in English for the domain registration.\n" +
							"  - example : Samsung-ro 123, Suwon-si, Gyeonggi-do, Korea ",
						Computed: true,
					},
					common.ToSnakeCase("DomesticAddressKo"): schema.StringAttribute{
						Description: "The domestic address in Korean for the domain registration.\n" +
							"  - example : 경기도 수원시 삼성로 123 ",
						Computed: true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "The first line of domestic address in English.\n" +
							"  - example : Samsung-ro 123 ",
						Computed: true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "The first line of domestic address in Korean.\n" +
							"  - example : 삼성로 123 ",
						Computed: true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "The second line of domestic address in English.\n" +
							"  - example : Suwon-si, Gyeonggi-do ",
						Computed: true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "The second line of domestic address in Korean.\n" +
							"  - example : 경기도 수원시 ",
						Computed: true,
					},
					common.ToSnakeCase("ExpiredDate"): schema.StringAttribute{
						Description: "The expiration date of the domain registration.\n" +
							"  - example : 2025-12-31 ",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the public domain name.\n" +
							"  - example : pdn-abc123def456 ",
						Computed: true,
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
						Description: "The name of the public domain name.\n" +
							"  - example : example.com ",
						Computed: true,
					},
					common.ToSnakeCase("OverseasAddress"): schema.StringAttribute{
						Description: "The overseas address for the domain registration.\n" +
							"  - example : 123 Main Street, City, Country ",
						Computed: true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "The first line of overseas address.\n" +
							"  - example : 123 Main Street ",
						Computed: true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "The second line of overseas address.\n" +
							"  - example : Suite 100 ",
						Computed: true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "The third line of overseas address.\n" +
							"  - example : City, State 12345 ",
						Computed: true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "The postal code for the domain registration.\n" +
							"  - example : 12345 ",
						Computed: true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "The email address of the domain registrant.\n" +
							"  - example : user@example.com ",
						Computed: true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "The name of the domain registrant in English.\n" +
							"  - example : John Doe ",
						Computed: true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "The name of the domain registrant in Korean.\n" +
							"  - example : 홍길동 ",
						Computed: true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "The telephone number of the domain registrant.\n" +
							"  - example : 82-10-1234-5678 ",
						Computed: true,
					},
					common.ToSnakeCase("StartDate"): schema.StringAttribute{
						Description: "The start date of the domain registration.\n" +
							"  - example : 2024-01-01 ",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the public domain name.\n" +
							"  - example : REGISTERED ",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsPublicDomainNameDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsPublicDomainNameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.PublicDomainNameDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPublicDomainName(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show PublicDomainName",
			err.Error(),
		)
		return
	}

	publicDomainState := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainObjectValue, _ := types.ObjectValueFrom(ctx, publicDomainState.AttributeTypes(), publicDomainState)
	state.PublicDomainNameDetail = publicDomainObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertPublicDomainDetail(publicDomainDetail scpdns.PublicDomainDetail) dns.PublicDomainNameDetail {
	return dns.PublicDomainNameDetail{
		AddressType:             types.StringValue(publicDomainDetail.AddressType),
		AutoExtension:           types.BoolValue(publicDomainDetail.AutoExtension),
		CreatedAt:               types.StringValue(publicDomainDetail.CreatedAt),
		CreatedBy:               types.StringValue(publicDomainDetail.CreatedBy),
		Description:             types.StringValue(publicDomainDetail.GetDescription()),
		DomesticAddressEn:       types.StringValue(publicDomainDetail.DomesticAddressEn),
		DomesticAddressKo:       types.StringValue(publicDomainDetail.DomesticAddressKo),
		DomesticFirstAddressEn:  types.StringValue(publicDomainDetail.DomesticFirstAddressEn),
		DomesticFirstAddressKo:  types.StringValue(publicDomainDetail.DomesticFirstAddressKo),
		DomesticSecondAddressEn: types.StringValue(publicDomainDetail.DomesticSecondAddressEn),
		DomesticSecondAddressKo: types.StringValue(publicDomainDetail.DomesticSecondAddressKo),
		ExpiredDate:             types.StringValue(publicDomainDetail.ExpiredDate),
		Id:                      types.StringValue(publicDomainDetail.Id),
		ModifiedAt:              types.StringValue(publicDomainDetail.ModifiedAt),
		ModifiedBy:              types.StringValue(publicDomainDetail.ModifiedBy),
		Name:                    types.StringValue(publicDomainDetail.Name),
		OverseasAddress:         types.StringValue(*publicDomainDetail.OverseasAddress),
		OverseasFirstAddress:    types.StringValue(publicDomainDetail.OverseasFirstAddress),
		OverseasSecondAddress:   types.StringValue(publicDomainDetail.OverseasSecondAddress),
		OverseasThirdAddress:    types.StringValue(publicDomainDetail.OverseasThirdAddress),
		PostalCode:              types.StringValue(publicDomainDetail.PostalCode),
		RegisterEmail:           types.StringValue(publicDomainDetail.RegisterEmail),
		RegisterNameEn:          types.StringValue(publicDomainDetail.RegisterNameEn),
		RegisterNameKo:          types.StringValue(publicDomainDetail.RegisterNameKo),
		RegisterTelno:           types.StringValue(publicDomainDetail.RegisterTelno),
		StartDate:               types.StringValue(publicDomainDetail.StartDate),
		Status:                  types.StringValue(publicDomainDetail.Status),
	}
}
