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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name (between 1 and 256 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			common.ToSnakeCase("ResourceGroup"): schema.SingleNestedAttribute{
				Description: "Resource Group",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.DataSourceSchema(),
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Region"): schema.StringAttribute{
						Description: "Region",
						Computed:    true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description: "Srn",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "ResourceTypes",
						Computed:    true,
					},
					common.ToSnakeCase("GroupDefinitionTags"): schema.ListNestedAttribute{
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
