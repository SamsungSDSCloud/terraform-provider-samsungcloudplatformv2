package resourcemanager

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resourceManagerResourceGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &resourceManagerResourceGroupDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewResourceManagerResourceGroupDataSources() datasource.DataSource {
	return &resourceManagerResourceGroupDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type resourceManagerResourceGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *resourceManagerResourceGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_resource_groups"
}

// Schema defines the schema for the data source.
func (d *resourceManagerResourceGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of resource group.",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			"tags":   tag.DataSourceSchema(),
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id (between 1 and 64 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name (between 1 and 256 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "ID List",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *resourceManagerResourceGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.ResourceManager
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *resourceManagerResourceGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resourcemanager.ResourceGroupDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		d.client.Config.Region = state.Region.ValueString()
	}

	ids, err := GetResourceGroups(d.clients, state.Id, state.Name, state.Filter, state.Tags.Elements(), "id")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	state.Ids = ids

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetResourceGroups(clients *client.SCPClient, id types.String, name types.String, filters []filter.Filter, tags map[string]attr.Value, idField string) ([]types.String, error) {
	data, err := clients.ResourceManager.GetResourceGroupList(id, name)
	if err != nil {
		return nil, err
	}

	contents := data.ResourceGroups
	filteredContents := data.ResourceGroups

	if len(tags) > 0 {
		filteredContents = filteredContents[:0]
		indices := tag.GetTagIndices(clients, contents, tags, idField)

		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, filters)

		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
