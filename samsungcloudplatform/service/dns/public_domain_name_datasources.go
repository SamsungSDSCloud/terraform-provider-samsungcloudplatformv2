package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsPublicDomainNameDataSources{}
	_ datasource.DataSourceWithConfigure = &dnsPublicDomainNameDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsPublicDomainNameDataSources() datasource.DataSource {
	return &dnsPublicDomainNameDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsPublicDomainNameDataSources struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsPublicDomainNameDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_public_domain_names" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsPublicDomainNameDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of public domain name.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "CreatedBy",
				Optional:    true,
			},
			common.ToSnakeCase("PublicDomainNames"): schema.ListNestedAttribute{
				Description: "A list of PublicDomainName.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "created at",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "created by",
							Computed:    true,
						},
						common.ToSnakeCase("ExpiredDate"): schema.StringAttribute{
							Description: "ExpiredDate",
							Optional:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
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
							Optional:    true,
						},
						common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
							Description: "Register email",
							Optional:    true,
						},
						common.ToSnakeCase("StartDate"): schema.StringAttribute{
							Description: "StartDate",
							Optional:    true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "Status",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsPublicDomainNameDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsPublicDomainNameDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.PublicDomainNameDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPublicDomainNameList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PublicDomainNames",
			err.Error(),
		)
		return
	}

	for _, publicDomainName := range data.PublicDomainNames {
		publicDomainNameState := dns.PublicDomainName{
			CreatedAt:     types.StringValue(publicDomainName.CreatedAt),
			CreatedBy:     types.StringValue(publicDomainName.CreatedBy),
			ExpiredDate:   types.StringValue(publicDomainName.ExpiredDate),
			Id:            types.StringValue(publicDomainName.Id),
			ModifiedAt:    types.StringValue(publicDomainName.ModifiedAt),
			ModifiedBy:    types.StringValue(publicDomainName.ModifiedBy),
			Name:          types.StringValue(publicDomainName.Name),
			RegisterEmail: types.StringValue(publicDomainName.RegisterEmail),
			StartDate:     types.StringValue(publicDomainName.StartDate),
			Status:        types.StringValue(publicDomainName.Status),
		}

		state.PublicDomainNames = append(state.PublicDomainNames, publicDomainNameState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}
