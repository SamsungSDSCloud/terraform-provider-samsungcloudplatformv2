package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/virtualserver/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &virtualServerVolumeResource{}
	_ resource.ResourceWithConfigure = &virtualServerVolumeResource{}
)

// NewComputeVolumeResource is a helper function to simplify the provider implementation.
func NewVirtualServerVolumeResource() resource.Resource {
	return &virtualServerVolumeResource{}
}

// virtualServerVolumeResource is the data source implementation.
type virtualServerVolumeResource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *virtualServerVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_volume"
}

// Schema defines the schema for the data source.
func (r *virtualServerVolumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "volume",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Required:    true,
			},
			common.ToSnakeCase("UserId"): schema.StringAttribute{
				Description: "UserId",
				Computed:    true,
			},
			common.ToSnakeCase("VolumeType"): schema.StringAttribute{
				Description: "VolumeType",
				Required:    true,
			},
			common.ToSnakeCase("Encrypted"): schema.BoolAttribute{
				Description: "Encrypted",
				Computed:    true,
			},
			common.ToSnakeCase("Bootable"): schema.BoolAttribute{
				Description: "Bootable",
				Computed:    true,
			},
			common.ToSnakeCase("Multiattach"): schema.BoolAttribute{
				Description: "Multiattach",
				Computed:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Computed:    true,
			},
			common.ToSnakeCase("Servers"): schema.ListNestedAttribute{
				Description: "Servers",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
					},
				},
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *virtualServerVolumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VirtualServer
	r.clients = inst.Client
}

func (r *virtualServerVolumeResource) MapGetResponseToState(resp *scpvirtualserver.VolumeShowResponse, state virtualserver.VolumeResource, tagsMap types.Map) virtualserver.VolumeResource {
	return virtualserver.VolumeResource{
		Id:          types.StringValue(resp.Id),
		Name:        types.StringPointerValue(resp.Name.Get()),
		UserId:      types.StringValue(resp.UserId),
		Bootable:    types.BoolValue(resp.Bootable),
		Multiattach: types.BoolValue(resp.Multiattach),
		Encrypted:   types.BoolValue(resp.Encrypted),
		VolumeType:  types.StringValue(resp.VolumeType),
		Size:        types.Int32Value(resp.Size),
		State:       types.StringValue(resp.State),
		Servers:     state.Servers,
		Tags:        tagsMap,
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *virtualServerVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualserver.VolumeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateVolume(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if plan.Servers != nil {
		for _, addedVm := range plan.Servers {
			_, err := r.client.AttachVolume(ctx, data.Id, addedVm.Id.ValueString())
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error updating volume",
					"Could not update volume, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	}

	getData, err := r.client.GetVolume(ctx, data.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading volume",
			"Could not create volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeVolume, data.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state := r.MapGetResponseToState(getData, plan, tagsMap)
	//state.Tags = plan.Tags

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *virtualServerVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualserver.VolumeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading volume",
			"Could not read volume ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeVolume, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState := r.MapGetResponseToState(data, state, tagsMap)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *virtualServerVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state virtualserver.VolumeResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) {
		// name attribute was changed
		_, err := r.client.UpdateVolume(ctx, state.Id.ValueString(), plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating volume",
				"Could not update volume, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}
	if !plan.Size.Equal(state.Size) {
		// size attribute was changed
		_, err := r.client.ExtendVolume(ctx, state.Id.ValueString(), plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating volume",
				"Could not update volume, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	if !reflect.DeepEqual(plan.Servers, state.Servers) {
		addedVmIds, deletedVmIds := getOldAndNewVmIds(plan, state)

		for _, deletedVmId := range deletedVmIds {
			err := r.client.DetachVolume(ctx, state.Id.ValueString(), deletedVmId)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error updating volume",
					"Could not update volume, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		for _, addedVmId := range addedVmIds {
			_, err := r.client.AttachVolume(ctx, state.Id.ValueString(), addedVmId)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error updating volume",
					"Could not update volume, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	}

	data, err := r.client.GetVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading volume",
			"Could not read volume ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, ServiceNameVirtualServer, ResourceTypeVolume, plan.Id.ValueString(), tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState := r.MapGetResponseToState(data, plan, tagsMap)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *virtualServerVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualserver.VolumeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting volume",
			"Could not delete volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func getOldAndNewVmIds(plan virtualserver.VolumeResource, state virtualserver.VolumeResource) ([]string, []string) {

	addedVmId := diff(plan.Servers, state.Servers)
	deletedVmId := diff(state.Servers, plan.Servers)
	return addedVmId, deletedVmId
}

func diff(a []virtualserver.VolumeServer, b []virtualserver.VolumeServer) []string {
	var result []string
	m := make(map[string]bool)

	for _, v := range b {
		m[v.Id.ValueString()] = true
	}

	for _, v := range a {
		if _, ok := m[v.Id.ValueString()]; !ok {
			result = append(result, v.Id.ValueString())
		}
	}

	return result
}
