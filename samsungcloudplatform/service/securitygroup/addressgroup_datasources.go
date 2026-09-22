package securitygroup

import (
	"context"
	"fmt"
	"time"

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
	_ datasource.DataSource              = &addressGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &addressGroupDataSources{}
)

// NewAddressGroupDataSources is a helper function to simplify the provider implementation.
func NewAddressGroupDataSources() datasource.DataSource {
	return &addressGroupDataSources{}
}

// addressGroupDataSources is the data source implementation.
type addressGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *securitygroupv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *addressGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_address_groups"
}

// Schema defines the schema for the data source.
func (r *addressGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of address groups",
		Attributes: map[string]schema.Attribute{
			// Input
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
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the Address Group.\n" +
					"  - example: ag-web-prod",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total number of Address Groups matching the criteria.",
				Computed:    true,
			},
			common.ToSnakeCase("Links"): schema.ListNestedAttribute{
				Description: "Collection of hypermedia links to related resources or pages.\n" +
					"  - example : [{\"href\": \"http://scp.samsungsdscloud.com/v1/notices\", \"rel\": \"self\"}]",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Href"): schema.StringAttribute{
							Description: "URL of the linked resource.\n" +
								"  - example : 'http://scp.samsungsdscloud.com/v1/notices'",
							Computed: true,
						},
						common.ToSnakeCase("Rel"): schema.StringAttribute{
							Description: "Relationship type of the link.\n" +
								"  - example : 'self'",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("AddressGroups"): schema.ListNestedAttribute{
				Description: "List of Address Groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The account ID associated with the resource.",
							Computed:    true,
						},
						common.ToSnakeCase("AddressCount"): schema.Int32Attribute{
							Description: "Number of addresses in the Address Group.",
							Computed:    true,
						},
						common.ToSnakeCase("AddressLimit"): schema.Int32Attribute{
							Description: "Maximum number of addresses allowed in the Address Group.",
							Computed:    true,
						},
						common.ToSnakeCase("Addresses"): schema.SetAttribute{
							Description: "A list of IP addresses or CIDR ranges in the Address Group.",
							ElementType: types.StringType,
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created in ISO 8601 format.",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user ID that created the resource.",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "A brief explanation or note about this resource.",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the resource.",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified in ISO 8601 format.",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user ID that modified the resource.",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the Address Group.",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current state of the resource.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *addressGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (r *addressGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config securitygroupv1d1.AddressGroupDataSources
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListAddressGroups(config.Page, config.Size, config.Sort, config.Id, config.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading address groups",
			"Could not read address groups: "+err.Error(),
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading address groups",
			"empty response from API",
		)
		return
	}

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
				config.Links = append(config.Links, securitygroupv1d1.Link{
					Href: types.StringValue(href),
					Rel:  types.StringValue(rel),
				})
			}
		}
	}

	var addressGroups []securitygroupv1d1.AddressGroup
	for _, ag := range data.AddressGroups {
		agObj := securitygroupv1d1.AddressGroup{
			AccountId:    types.StringValue(ag.AccountId),
			AddressCount: types.Int32Value(ag.GetAddressCount()),
			AddressLimit: types.Int32Value(ag.GetAddressLimit()),
			CreatedAt:    types.StringValue(ag.CreatedAt.Format(time.RFC3339)),
			CreatedBy:    types.StringValue(ag.CreatedBy),
			Description:  types.StringPointerValue(ag.Description.Get()),
			Id:           types.StringValue(ag.Id),
			ModifiedAt:   types.StringValue(ag.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:   types.StringValue(ag.ModifiedBy),
			Name:         types.StringValue(ag.GetName()),
			State:        types.StringValue(ag.GetState()),
		}

		addressLst, dia := types.SetValueFrom(ctx, types.StringType, ag.Addresses)
		resp.Diagnostics.Append(dia...)
		if resp.Diagnostics.HasError() {
			return
		}

		agObj.Addresses = addressLst
		addressGroups = append(addressGroups, agObj)
	}

	config.AddressGroups = addressGroups

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
