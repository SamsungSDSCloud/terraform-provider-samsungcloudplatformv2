package resourcemanager

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resourceManagerTagDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceManagerTagDataSource{}
)

// NewResourceManagerTagDataSource is a helper function to simplify the provider implementation.
func NewResourceManagerTagDataSource() datasource.DataSource {
	return &resourceManagerTagDataSource{}
}

// resourceManagerTagDataSource is the data source implementation.
type resourceManagerTagDataSource struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *resourceManagerTagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_tags"
}

// Schema defines the schema for the data source.
func (d *resourceManagerTagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of resource tag.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Key"): schema.StringAttribute{
				Description: "Tag Key",
				Optional:    true,
			},
			common.ToSnakeCase("Value"): schema.StringAttribute{
				Description: "Tag Value",
				Optional:    true,
			},
			common.ToSnakeCase("ResourceIdentifier"): schema.StringAttribute{
				Description: "Resource Identifier",
				Optional:    true,
			},
			common.ToSnakeCase("Content"): schema.ListNestedAttribute{
				Description: "A list of tag.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Srn"): schema.StringAttribute{
							Description: "Srn",
							Computed:    true,
						},
						common.ToSnakeCase("Key"): schema.StringAttribute{
							Description: "Key",
							Computed:    true,
						},
						common.ToSnakeCase("Value"): schema.StringAttribute{
							Description: "Value",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *resourceManagerTagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *resourceManagerTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resourcemanager.TagDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetTagList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tags",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, tag := range data.Content {
		tagState := resourcemanager.TagContent{
			Srn:   types.StringValue(tag.Srn),
			Key:   types.StringValue(tag.Key),
			Value: types.StringValue(tag.Value),
		}
		state.Contents = append(state.Contents, tagState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
