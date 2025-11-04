package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/dns/1.2"
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
		Description: "Show PublicDomainName.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("PublicDomainNameDetail"): schema.SingleNestedAttribute{
				Description: "A detail of PublicDomainName.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "AddressType",
						Computed:    true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "AutoExtension",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticAddressEn"): schema.StringAttribute{
						Description: "DomesticAddressEn",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticAddressKo"): schema.StringAttribute{
						Description: "DomesticAddressKo",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "DomesticFirstAddressEn",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "DomesticFirstAddressKo",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "DomesticSecondAddressEn",
						Computed:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "DomesticSecondAddressKo",
						Computed:    true,
					},
					common.ToSnakeCase("ExpiredDate"): schema.StringAttribute{
						Description: "ExpiredDate",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("OverseasAddress"): schema.StringAttribute{
						Description: "OverseasAddress",
						Computed:    true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "OverseasFirstAddress",
						Computed:    true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "OverseasSecondAddress",
						Computed:    true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "OverseasThirdAddress",
						Computed:    true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "PostalCode",
						Computed:    true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "RegisterEmail",
						Computed:    true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "RegisterNameEn",
						Computed:    true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "RegisterNameKo",
						Computed:    true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "RegisterTelno",
						Computed:    true,
					},
					common.ToSnakeCase("StartDate"): schema.StringAttribute{
						Description: "StartDate",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
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
