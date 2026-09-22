package scr

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrContainerRegistryDataSource{}
	_ datasource.DataSourceWithConfigure = &scrContainerRegistryDataSource{}
)

func NewScrContainerRegistryDataSource() datasource.DataSource {
	return &scrContainerRegistryDataSource{}
}

type scrContainerRegistryDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrContainerRegistryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_container_registry"
}

func (d *scrContainerRegistryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry.",
		MarkdownDescription: "Get details of a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The unique identifier of the registry.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Registry name",
				MarkdownDescription: "The name of the registry.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Registry state",
				MarkdownDescription: "The current state of the registry.",
				Computed:            true,
			},
			"private_domain": schema.StringAttribute{
				Description:         "Private endpoint URL",
				MarkdownDescription: "The private endpoint URL for the registry.",
				Computed:            true,
			},
			"public_domain": schema.StringAttribute{
				Description:         "Public endpoint URL",
				MarkdownDescription: "The public endpoint URL for the registry.",
				Computed:            true,
			},
			"bucket_id": schema.StringAttribute{
				Description:         "Underlying bucket ID",
				MarkdownDescription: "The ID of the underlying object storage bucket.",
				Computed:            true,
			},
			"bucket_name": schema.StringAttribute{
				Description:         "Underlying bucket name",
				MarkdownDescription: "The name of the underlying object storage bucket.",
				Computed:            true,
			},
			"bucket_usage": schema.StringAttribute{
				Description:         "Bucket usage",
				MarkdownDescription: "The current usage of the underlying bucket.",
				Computed:            true,
			},
			"public_visible_enabled": schema.BoolAttribute{
				Description:         "Public visible enabled",
				MarkdownDescription: "Whether the registry is publicly visible.",
				Computed:            true,
			},
			"public_endpoint_enabled": schema.BoolAttribute{
				Description:         "Public endpoint enabled",
				MarkdownDescription: "Whether the public endpoint is enabled.",
				Computed:            true,
			},
			"private_acl_enabled": schema.BoolAttribute{
				Description:         "Private ACL enabled",
				MarkdownDescription: "Whether private ACL is enabled.",
				Computed:            true,
			},
			"private_acl_resources": schema.ListNestedAttribute{
				Description:         "Private ACL resource list",
				MarkdownDescription: "List of resources allowed private access.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
						"resource_name": schema.StringAttribute{
							Description: "Resource name",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type",
							Computed:    true,
						},
						"resource_ips": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "Resource IP addresses",
							Computed:    true,
						},
					},
				},
			},
			"public_acl_enabled": schema.BoolAttribute{
				Description:         "Public ACL enabled",
				MarkdownDescription: "Whether public ACL is enabled.",
				Computed:            true,
			},
			"public_acl_resources": schema.ListNestedAttribute{
				Description:         "Public ACL resource list",
				MarkdownDescription: "List of resources allowed public access.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description: "Resource ID",
							Computed:    true,
						},
						"resource_name": schema.StringAttribute{
							Description: "Resource name",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type",
							Computed:    true,
						},
						"resource_ips": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "Resource IP addresses",
							Computed:    true,
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description:         "Created at",
				MarkdownDescription: "The time the registry was created.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				Description:         "Created by",
				MarkdownDescription: "The user who created the registry.",
				Computed:            true,
			},
			"modified_at": schema.StringAttribute{
				Description:         "Modified at",
				MarkdownDescription: "The time the registry was last modified.",
				Computed:            true,
			},
			"modified_by": schema.StringAttribute{
				Description:         "Modified by",
				MarkdownDescription: "The user who last modified the registry.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrContainerRegistryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Scr
	d.clients = inst.Client
}

func (d *scrContainerRegistryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ContainerRegistryDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := d.client.ShowRegistry(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Registry",
			err.Error(),
		)
		return
	}

	registry := showResp.Registry
	state = ContainerRegistryDataSource{
		Id:                    types.StringValue(registry.Id),
		Name:                  types.StringValue(registry.Name),
		State:                 types.StringValue(registry.State),
		PrivateDomain:         types.StringValue(registry.PrivateDomain),
		PublicDomain:          nullableStringValue(registry.PublicDomain.Get()),
		BucketId:              types.StringValue(registry.BucketId),
		BucketName:            types.StringValue(registry.BucketName),
		BucketUsage:           types.StringValue(registry.BucketUsage),
		PublicVisibleEnabled:  types.BoolValue(registry.PublicVisibleEnabled),
		PublicEndpointEnabled: common.ToNullableBoolValue(registry.PublicEndpointEnabled.Get()),
		PrivateAclEnabled:     types.BoolValue(registry.PrivateAclEnabled),
		PrivateAclResources:   flattenContainerRegistryAclResources(registry.PrivateAclResources),
		PublicAclEnabled:      common.ToNullableBoolValue(registry.PublicAclEnabled.Get()),
		PublicAclResources:    flattenContainerRegistryAclResources(registry.PublicAclResources),
		CreatedAt:             types.StringValue(registry.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringValue(registry.CreatedBy),
		ModifiedAt:            types.StringValue(registry.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:            types.StringValue(registry.ModifiedBy),
		Filter:                state.Filter,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenContainerRegistryAclResources(resources []scr11.Resource) types.List {
	var aclList []AclResource
	for _, r := range resources {
		ipsList, _ := types.ListValueFrom(context.Background(), types.StringType, r.ResourceIps)
		aclList = append(aclList, AclResource{
			ResourceId:   nullableStringValue(r.ResourceId.Get()),
			ResourceName: nullableStringValue(r.ResourceName.Get()),
			ResourceType: nullableStringValue(r.ResourceType.Get()),
			ResourceIps:  ipsList,
		})
	}
	result, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: AclResource{}.AttributeTypes()}, aclList)
	return result
}
