package resourcemanager

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description:         "System Resource Name.",
				MarkdownDescription: "The System Resource Name (SRN) of the resource group.\n\nExample: `srn:s::13d97ad943ca452481d624f78391df13:kr-west1::resourcemanager:resource-group/70636f984e564b3c9e54e74a53f9318d`",
				Required:            true,
			},
			common.ToSnakeCase("EncodedSrn"): schema.StringAttribute{
				Description:         "System Resource Name with base64 encoded.",
				MarkdownDescription: "The System Resource Name (SRN) of the resource group with base64 encoded.\n\nExample: `c3JuOnM6OjEzZDk3YWQ5NDNjYTQ1MjQ4MWQ2MjRmNzgzOTFkZjEzOmtyLXdlc3QxOjpyZXNvdXJjZW1hbmFnZXI6cmVzb3VyY2UtZ3JvdXAvNzA2MzZmOTg0ZTU2NGIzYzllNTRlNzRhNTNmOTMxOGQ=`",
				Computed:            true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description:         "A Number of results displayed per page.",
				MarkdownDescription: "A Number of results displayed per page.\n\nExample: `15`",
				Optional:            true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description:         "A Number of page.",
				MarkdownDescription: "A Number of page.\n\nExample: `1`",
				Optional:            true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description:         "Sorts the query results.",
				MarkdownDescription: "Sorts the query results.\n\nExample: `createdAt:desc`",
				Optional:            true,
			},
			common.ToSnakeCase("Content"): schema.SingleNestedAttribute{
				Description:         "This is set data of SRNs and tags.",
				MarkdownDescription: "This is set data of SRNs and tags.\n\nExample: See nested attributes below.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description:         "System Resource Name.",
						MarkdownDescription: "The System Resource Name (SRN) of the resource group.\n\nExample: `srn:s::13d97ad943ca452481d624f78391df13:kr-west1::resourcemanager:resource-group/70636f984e564b3c9e54e74a53f9318d`",
						Computed:            true,
					},
					common.ToSnakeCase("Tags"): schema.ListNestedAttribute{
						Description:         "A list of tag.",
						MarkdownDescription: "A list of tag.\n\nExample: See nested attributes below.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Key"): schema.StringAttribute{
									Description:         "Tag key.",
									MarkdownDescription: "The key of the tag.\n\nExample: `Environment`",
									Computed:            true,
								},
								common.ToSnakeCase("Value"): schema.StringAttribute{
									Description:         "Tag value.",
									MarkdownDescription: "The value of the tag.\n\nExample: `Production`",
									Computed:            true,
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
