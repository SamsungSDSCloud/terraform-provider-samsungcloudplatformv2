package resourcemanager

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resourceManagerResourceTagDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceManagerResourceTagDataSource{}
)

// NewResourceManagerResourceTagDataSource is a helper function to simplify the provider implementation.
func NewResourceManagerResourceTagDataSource() datasource.DataSource {
	return &resourceManagerResourceTagDataSource{}
}

// resourceManagerResourceTagDataSource is the data source implementation.
type resourceManagerResourceTagDataSource struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *resourceManagerResourceTagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_resource_tags"
}

// Schema defines the schema for the data source.
func (d *resourceManagerResourceTagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of resource tag.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Srn"): schema.StringAttribute{
				Description: "Srn",
				Required:    true,
			},
			common.ToSnakeCase("EncodedSrn"): schema.StringAttribute{
				Description: "Encoded Srn",
				Computed:    true,
			},
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
			common.ToSnakeCase("Content"): schema.SingleNestedAttribute{
				Description: "Content",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description: "Srn",
						Computed:    true,
					},
					common.ToSnakeCase("Tags"): schema.ListNestedAttribute{
						Description: "A list of tag.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
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
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *resourceManagerResourceTagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *resourceManagerResourceTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resourcemanager.ResourceTagDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.EncodedSrn = types.StringValue(common.EncodeBase64(state.Srn.ValueString()))

	data, err := d.client.GetResourceTagList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Tags",
			err.Error(),
		)
		return
	}

	srnTagList := resourcemanager.SrnTagList{
		Srn: types.StringValue(data.Content.Srn),
	}

	// Map response body to model
	for _, tag := range data.Content.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(tag.Key),
			Value: types.StringPointerValue(tag.Value.Get()),
		}
		srnTagList.Tags = append(srnTagList.Tags, tagState)
	}

	srnTagListObjectValue, diags := types.ObjectValueFrom(ctx, srnTagList.AttributeTypes(), srnTagList)
	state.SrnTagList = srnTagListObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
