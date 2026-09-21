package securitygroup

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	securitygroupv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &addressGroupCidrDataSources{}
	_ datasource.DataSourceWithConfigure = &addressGroupCidrDataSources{}
)

// NewAddressGroupCidrDataSources is a helper function to simplify the provider implementation.
func NewAddressGroupCidrDataSources() datasource.DataSource {
	return &addressGroupCidrDataSources{}
}

// addressGroupCidrDataSources is the data source implementation.
type addressGroupCidrDataSources struct {
	config  *scpsdk.Configuration
	client  *securitygroupv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *addressGroupCidrDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_address_group_cidrs"
}

// Schema defines the schema for the data source.
func (r *addressGroupCidrDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of CIDRs in an Address Group",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("AddressGroupId"): schema.StringAttribute{
				Description: "The unique identifier of the Address Group.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example: 20\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example: 1\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria.\n" +
					"  - example: created_at:desc\n" +
					"  - valid: field_name:asc or field_name:desc",
				Optional: true,
			},
			common.ToSnakeCase("Address"): schema.StringAttribute{
				Description: "Filter by a specific IP address or CIDR range.\n" +
					"  - example: 192.168.1.0/24",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("Addresses"): schema.SetAttribute{
				Description: "A list of IP addresses or CIDR ranges in the Address Group.",
				ElementType: types.StringType,
				Computed:    true,
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total number of CIDRs matching the criteria.",
				Computed:    true,
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Collection of hypermedia links to related resources or pages.\n" +
					"  - example : [{\"href\": \"http://scp.samsungsdscloud.com/v1/notices\", \"rel\": \"self\"}]",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Href"): schema.StringAttribute{
							Description: "The URL of the linked resource.",
							Computed:    true,
						},
						common.ToSnakeCase("Rel"): schema.StringAttribute{
							Description: "The relationship type of the link.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *addressGroupCidrDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = inst.Client.SecurityGroupV1d1
	r.clients = inst.Client
}

// Read refreshes the data source.
func (r *addressGroupCidrDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config securitygroupv1d1.AddressGroupCidrDataSources
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListAddressGroupCidrs(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading address group CIDRs",
			"Could not read address group CIDRs: "+err.Error(),
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading address group CIDRs",
			"Empty response from API",
		)
		return
	}

	addresses, diags := types.SetValueFrom(ctx, types.StringType, data.Addresses)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Addresses = addresses
	config.TotalCount = types.Int32PointerValue(data.Count.Get())

	if len(data.Links) > 0 {
		for _, element := range data.Links {
			var href, rel string

			if elemMap, ok := element.(map[string]interface{}); ok {
				if h, exists := elemMap["href"]; exists {
					if hStr, ok := h.(string); ok {
						href = hStr
					}
				}
				if r, exists := elemMap["rel"]; exists {
					if rStr, ok := r.(string); ok {
						rel = rStr
					}
				}
			}
			config.Links = append(config.Links, securitygroupv1d1.Link{
				Href: types.StringValue(href),
				Rel:  types.StringValue(rel),
			})
		}
	}

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
