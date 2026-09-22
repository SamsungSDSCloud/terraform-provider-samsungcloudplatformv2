package scr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &scrContainerRegistryResource{}
	_ resource.ResourceWithConfigure = &scrContainerRegistryResource{}
	_ resource.ResourceWithImportState = &scrContainerRegistryResource{}
)

func NewScrContainerRegistryResource() resource.Resource {
	return &scrContainerRegistryResource{}
}

type scrContainerRegistryResource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (r *scrContainerRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_container_registry"
}

func (r *scrContainerRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = inst.Client.Scr
	r.clients = inst.Client
}

func (r *scrContainerRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *scrContainerRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Create and manage a Container Registry.",
		MarkdownDescription: "Manages a Container Registry resource for storing Docker images.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The unique identifier of the registry.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Registry name",
				MarkdownDescription: "The name of the container registry. Must start with a lowercase letter and contain only lowercase letters and numbers.\n\nExample: `my-registry`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
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
				MarkdownDescription: "Whether the registry is publicly visible.\n\nAllowed values: `true` | `false`\n\nExample: `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"public_endpoint_enabled": schema.BoolAttribute{
				Description:         "Public endpoint enabled",
				MarkdownDescription: "Whether the public endpoint is enabled.\n\nAllowed values: `true` | `false`\n\nExample: `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"private_acl_enabled": schema.BoolAttribute{
				Description:         "Private ACL enabled",
				MarkdownDescription: "Whether private ACL is enabled.\n\nAllowed values: `true` | `false`\n\nExample: `true`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"private_acl_resources": schema.ListNestedAttribute{
				Description:         "Private ACL resource list",
				MarkdownDescription: "List of resources allowed private access to the registry.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description:         "Resource ID",
							MarkdownDescription: "The ID of the resource.",
							Optional:            true,
						},
						"resource_name": schema.StringAttribute{
							Description:         "Resource name",
							MarkdownDescription: "The name of the resource.",
							Optional:            true,
						},
						"resource_type": schema.StringAttribute{
							Description:         "Resource type",
							MarkdownDescription: "The type of the resource.\n\nExample: `virtualServer`",
							Optional:            true,
						},
						"resource_ips": schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "Resource IP addresses",
							MarkdownDescription: "List of IP addresses for the resource.",
							Optional:            true,
						},
					},
				},
			},
			"public_acl_enabled": schema.BoolAttribute{
				Description:         "Public ACL enabled",
				MarkdownDescription: "Whether public ACL is enabled.\n\nAllowed values: `true` | `false`\n\nExample: `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"public_acl_resources": schema.ListNestedAttribute{
				Description:         "Public ACL resource list",
				MarkdownDescription: "List of resources allowed public access to the registry.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description:         "Resource ID",
							MarkdownDescription: "The ID of the resource.",
							Optional:            true,
						},
						"resource_name": schema.StringAttribute{
							Description:         "Resource name",
							MarkdownDescription: "The name of the resource.",
							Optional:            true,
						},
						"resource_type": schema.StringAttribute{
							Description:         "Resource type",
							MarkdownDescription: "The type of the resource.\n\nExample: `virtualServer`",
							Optional:            true,
						},
						"resource_ips": schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "Resource IP addresses",
							MarkdownDescription: "List of IP addresses for the resource.",
							Optional:            true,
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
	}
}

func (r *scrContainerRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerRegistryResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateAclResources := buildAclResources(ctx, plan.PrivateAclResources)

	publicAclResources := buildAclResources(ctx, plan.PublicAclResources)
	if resp.Diagnostics.HasError() {
		return
	}

	publicAclEnabled := plan.PublicAclEnabled.ValueBool()
	publicEndpointEnabled := plan.PublicEndpointEnabled.ValueBool()

	createReq := scr11.ContainerRegistryCreateRequest{
		Name:                   plan.Name.ValueString(),
		PrivateAclEnabled:      plan.PrivateAclEnabled.ValueBool(),
		PrivateAclResources:    privateAclResources,
		PublicAclEnabled:       *scr11.NewNullableBool(&publicAclEnabled),
		PublicAclResources:     publicAclResources,
		PublicEndpointEnabled:  *scr11.NewNullableBool(&publicEndpointEnabled),
		PublicVisibleEnabled:   plan.PublicVisibleEnabled.ValueBool(),
	}

	createResp, err := r.client.CreateRegistry(ctx, createReq)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating registry",
			"Could not create registry, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	// Poll until registry is ready
	if createResp.Id != "" {
		err = r.waitForRegistryReady(ctx, createResp.Id, 120*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for registry",
				"Registry was created but failed to become ready: "+err.Error(),
			)
			return
		}
	}

	// Read back the full state
	showResp, err := r.client.ShowRegistry(ctx, createResp.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading registry after creation",
			"Could not read registry ID "+createResp.Id+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	plan = flattenContainerRegistry(showResp.Registry)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scrContainerRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerRegistryResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := r.client.ShowRegistry(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Registry",
			"Could not read registry ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	state = flattenContainerRegistry(showResp.Registry)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *scrContainerRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerRegistryResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState ContainerRegistryResource
	diags = req.State.Get(ctx, &priorState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registryId := plan.Id.ValueString()

	// Update private ACL
	privateAclResources := buildAclResources(ctx, plan.PrivateAclResources)
	if resp.Diagnostics.HasError() {
		return
	}
	privateAclReq := scr11.PrivateAclSetRequest{
		PrivateAclEnabled:   plan.PrivateAclEnabled.ValueBool(),
		PrivateAclResources: privateAclResources,
	}
	if err := r.client.UpdatePrivateAcl(ctx, registryId, privateAclReq); err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating private ACL",
			"Could not update private ACL: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	// Update public endpoint enabled before public ACL
	publicEndpointReq := scr11.PublicEndpointEnabledSetRequest{
		PublicEndpointEnabled: plan.PublicEndpointEnabled.ValueBool(),
	}
	if err := r.client.UpdatePublicEndpointEnabled(ctx, registryId, publicEndpointReq); err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating public endpoint",
			"Could not update public endpoint: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	// Update public ACL when there is a change
	publicAclChanged := plan.PublicAclEnabled.ValueBool() != priorState.PublicAclEnabled.ValueBool() ||
		!plan.PublicAclResources.Equal(priorState.PublicAclResources)
	if publicAclChanged {
		publicAclResources := buildAclResources(ctx, plan.PublicAclResources)
		if resp.Diagnostics.HasError() {
			return
		}
		publicAclReq := scr11.PublicAclSetRequest{
			PublicAclEnabled:   plan.PublicAclEnabled.ValueBool(),
			PublicAclResources: publicAclResources,
		}
		if err := r.client.UpdatePublicAcl(ctx, registryId, publicAclReq); err != nil {
			// API returns 409 when public endpoint is disabled — ignore it
			if !strings.Contains(err.Error(), "409") {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error updating public ACL",
					"Could not update public ACL: "+err.Error()+reasonPrefix+detail,
				)
				return
			}
		}
	}

	state := plan

	// Read back the full state
	showResp, err := r.client.ShowRegistry(ctx, registryId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading registry after update",
			"Could not read registry ID "+registryId+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	state = flattenContainerRegistry(showResp.Registry)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *scrContainerRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ContainerRegistryResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRegistry(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Registry",
			"Could not delete registry, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}
}

func (r *scrContainerRegistryResource) waitForRegistryReady(ctx context.Context, id string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timeout waiting for registry %s to be ready", id)
		case <-ticker.C:
			showResp, err := r.client.ShowRegistry(ctx, id)
			if err != nil {
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
					continue
				}
				if client.IsTransientError(err) {
					continue
				}
				return err
			}
			registry := showResp.Registry
			if registry.Id != "" {
				return nil
			}
		}
	}
}

func nullableStringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func flattenContainerRegistry(registry scr11.ContainerRegistry) ContainerRegistryResource {
	privateAclList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: AclResource{}.AttributeTypes()}, convertAclResources(registry.PrivateAclResources))
	publicAclList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: AclResource{}.AttributeTypes()}, convertAclResources(registry.PublicAclResources))

	return ContainerRegistryResource{
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
		PrivateAclResources:   privateAclList,
		PublicAclEnabled:      common.ToNullableBoolValue(registry.PublicAclEnabled.Get()),
		PublicAclResources:    publicAclList,
		CreatedAt:             types.StringValue(registry.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringValue(registry.CreatedBy),
		ModifiedAt:            types.StringValue(registry.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:            types.StringValue(registry.ModifiedBy),
	}
}

func buildAclResources(ctx context.Context, aclList types.List) []scr11.Resource {
	if aclList.IsNull() || aclList.IsUnknown() {
		return []scr11.Resource{}
	}

	var aclResourceList []AclResource
	diags := aclList.ElementsAs(ctx, &aclResourceList, false)
	if diags.HasError() {
		return []scr11.Resource{}
	}

	resources := []scr11.Resource{}
	for _, ar := range aclResourceList {
		var ips []string
		diags := ar.ResourceIps.ElementsAs(ctx, &ips, false)
		if diags.HasError() {
			continue
		}

		res := scr11.Resource{
			ResourceIps: []string{},
		}
		if ips != nil {
			res.ResourceIps = ips
		}
		if !ar.ResourceId.IsNull() {
			v := ar.ResourceId.ValueString()
			res.ResourceId = *scr11.NewNullableString(&v)
		}
		if !ar.ResourceName.IsNull() {
			v := ar.ResourceName.ValueString()
			res.ResourceName = *scr11.NewNullableString(&v)
		}
		if !ar.ResourceType.IsNull() {
			v := ar.ResourceType.ValueString()
			res.ResourceType = *scr11.NewNullableString(&v)
		}
		resources = append(resources, res)
	}
	return resources
}

func convertAclResources(resources []scr11.Resource) []AclResource {
	result := []AclResource{}
	for _, r := range resources {
		ar := AclResource{
			ResourceId:   nullableStringValue(r.ResourceId.Get()),
			ResourceName: nullableStringValue(r.ResourceName.Get()),
			ResourceType: nullableStringValue(r.ResourceType.Get()),
		}
		ips := r.ResourceIps
		if ips == nil {
			ips = []string{}
		}
		ipsList, _ := types.ListValueFrom(context.Background(), types.StringType, ips)
		ar.ResourceIps = ipsList
		result = append(result, ar)
	}
	return result
}
