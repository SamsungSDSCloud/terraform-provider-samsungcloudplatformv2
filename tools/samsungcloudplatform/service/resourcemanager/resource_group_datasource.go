package resourcemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resourceManagerResourceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceManagerResourceGroupDataSource{}
)

// NewResourceManagerResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewResourceManagerResourceGroupDataSource() datasource.DataSource {
	return &resourceManagerResourceGroupDataSource{}
}

// resourceManagerResourceGroupDataSource is the data source implementation.
type resourceManagerResourceGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *resourceManagerResourceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_resource_group"
}

// Schema defines the schema for the data source.
func (d *resourceManagerResourceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "resource group.",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			"tags":   tag.DataSourceSchema(),
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "The unique identifier of the resource group.",
				MarkdownDescription: "The unique identifier of the resource group.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Optional:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "The name of the resource group.",
				MarkdownDescription: "The name of the resource group.\n\nExample: `example-rg`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			common.ToSnakeCase("ResourceGroup"): schema.SingleNestedAttribute{
				Description:         "The detailed information of the resource group.",
				MarkdownDescription: "The detailed information of the resource group.\n\nExample: See nested attributes below.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.DataSourceSchema(),
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description:         "Creation timestamp.",
						MarkdownDescription: "The creation timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description:         "Creator identifier.",
						MarkdownDescription: "The user ID of the creator.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description:         "Description of the resource group.",
						MarkdownDescription: "Description of the resource group.\n\nExample: `My resource group`",
						Computed:            true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "Identifier of the resource group.",
						MarkdownDescription: "The unique identifier of the resource group.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description:         "Modification timestamp.",
						MarkdownDescription: "The modification timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description:         "Modifier identifier.",
						MarkdownDescription: "The user ID of the modifier.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Name of the resource group.",
						MarkdownDescription: "The name of the resource group.\n\nExample: `example-rg`",
						Computed:            true,
					},
					common.ToSnakeCase("Region"): schema.StringAttribute{
						Description:         "Region code.",
						MarkdownDescription: "The region code where the resource group is located.\n\nExample: `kr-west1`",
						Computed:            true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description:         "System Resource Name.",
						MarkdownDescription: "The System Resource Name (SRN) of the resource group.\n\nExample: `srn:s::13d97ad943ca452481d624f78391df13:kr-west1::resourcemanager:resource-group/70636f984e564b3c9e54e74a53f9318d`",
						Computed:            true,
					},
					common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
						ElementType:         types.StringType,
						Description:         "List of resource types.",
						MarkdownDescription: "A list of resource types associated with the group.\n\nExample: `[\"virtual-server\"]`",
						Computed:            true,
					},
					common.ToSnakeCase("GroupDefinitionTags"): schema.ListNestedAttribute{
						Description:         "List of group definition tags.",
						MarkdownDescription: "A list of key-value pairs representing group definition tags.\n\nExample: `[{Key: Environment, Value: Production}]`",
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
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *resourceManagerResourceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *resourceManagerResourceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resourcemanager.ResourceGroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetResourceGroups(d.clients, state.Id, state.Name, state.Filter, state.Tags.Elements(), "id")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	if len(ids) > 0 {
		id := ids[0]

		// Get refreshed value from Resource Group
		data, err := d.client.GetResourceGroup(ctx, id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Resource Group",
				"Could not read Resource Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		resourceGroup := data.ResourceGroup

		var resourceTypes []string
		for _, resourceType := range resourceGroup.ResourceTypes {
			resourceTypes = append(resourceTypes, resourceType)
		}
		resourceTypesObject, _ := types.ListValueFrom(ctx, types.StringType, resourceTypes)

		var groupDefinitionTags []resourcemanager.Tag
		for _, t := range resourceGroup.Tags {
			tagState := resourcemanager.Tag{
				Key:   types.StringValue(t.Key),
				Value: types.StringPointerValue(t.Value.Get()),
			}
			groupDefinitionTags = append(groupDefinitionTags, tagState)
		}

		resourceGroupModel := resourcemanager.ResourceGroup{
			CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
			Description:         types.StringPointerValue(resourceGroup.Description.Get()),
			Id:                  types.StringValue(resourceGroup.Id),
			ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
			Name:                types.StringPointerValue(resourceGroup.Name.Get()),
			Region:              types.StringPointerValue(resourceGroup.Region.Get()),
			Srn:                 types.StringValue(resourceGroup.Srn),
			Tags:                state.Tags,
			ResourceTypes:       resourceTypesObject,
			GroupDefinitionTags: groupDefinitionTags,
		}
		resourceGroupObjectValue, _ := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
		state.ResourceGroup = resourceGroupObjectValue
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
