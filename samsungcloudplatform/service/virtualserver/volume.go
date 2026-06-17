package virtualserver

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/virtualserver/1.3"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const reasonPrefix = "\nReason: "

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
		Description:         "Creates a volume.",
		MarkdownDescription: "Creates a block storage volume for virtual servers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Resource ID.",
				MarkdownDescription: "Resource ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Volume name.\n" +
					"  - example: my-volume\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				MarkdownDescription: "Volume name.\n" +
					"  - example: my-volume\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description:         "Volume size (GiB). Must be a multiple of 8.",
				MarkdownDescription: "Volume size (GiB). Must be a multiple of 8.\n  - example: 104\n  - minimum: 8",
				Required:            true,
			},
			common.ToSnakeCase("UserId"): schema.StringAttribute{
				Description:         "User ID.",
				MarkdownDescription: "User ID.",
				Computed:            true,
			},
			common.ToSnakeCase("VolumeType"): schema.StringAttribute{
				Description: "Volume type.\n" +
					"  - example: ssd\n" +
					"  - Available values: ssd_provisioned, ssd, hdd",
				MarkdownDescription: "Volume type.\n" +
					"  - example: ssd\n" +
					"  - Available values: ssd_provisioned, ssd, hdd",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Encrypted"): schema.BoolAttribute{
				Description:         "Encryption flag.",
				MarkdownDescription: "Whether the volume is encrypted.",
				Computed:            true,
			},
			common.ToSnakeCase("Bootable"): schema.BoolAttribute{
				Description:         "Bootable flag.",
				MarkdownDescription: "Whether the volume is bootable.",
				Computed:            true,
			},
			common.ToSnakeCase("Multiattach"): schema.BoolAttribute{
				Description:         "Multi-attach flag.",
				MarkdownDescription: "Whether the volume can be attached to multiple servers.",
				Computed:            true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "Volume state.",
				MarkdownDescription: "Volume state.",
				Computed:            true,
			},
			common.ToSnakeCase("Servers"): schema.ListNestedAttribute{
				Description:         "List of attached servers.",
				MarkdownDescription: "List of attached servers.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Server ID.",
							MarkdownDescription: "Server ID.\n  - example: 97e6b22c-9a41-4378-9ca5-22df7457a32f",
							Optional:            true,
						},
					},
				},
			},
			common.ToSnakeCase("MaxIops"): schema.Int32Attribute{
				Description: "Maximum IOPS per second.\n" +
					"  - example: 10000\n" +
					"  - note: Number of read/write operations a volume can process per second",
				MarkdownDescription: "Maximum IOPS per second.\n" +
					"  - example: 10000\n" +
					"  - note: Number of read/write operations a volume can process per second",
				Optional: true,
			},
			common.ToSnakeCase("MaxThroughput"): schema.Int32Attribute{
				Description: "Maximum throughput per second (MB/s).\n" +
					"  - example: 500\n" +
					"  - note: Actual amount of data transferred to/from storage device per second",
				MarkdownDescription: "Maximum throughput per second (MB/s).\n" +
					"  - example: 500\n" +
					"  - note: Actual amount of data transferred to/from storage device per second",
				Optional: true,
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

func (r *virtualServerVolumeResource) AsyncPollingQosUpdate(ctx context.Context, volumeId string, desiredIops, desiredThroughput int32) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 60; i++ {
		vol, err := r.client.GetVolume(ctx, volumeId)
		if err != nil {
			return fmt.Errorf("failed to get volume during polling: %w", err)
		}

		iopsMatch := *vol.MaxIops.Get() == desiredIops
		throughputMatch := *vol.MaxThroughput.Get() == desiredThroughput

		if iopsMatch && throughputMatch {
			return nil
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout waiting for volume update (ID: %s)", volumeId)
}

func (r *virtualServerVolumeResource) MapGetResponseToState(resp *scpvirtualserver.VolumeShowResponseV1Dot2, state virtualserver.VolumeResource, tagsMap types.Map) virtualserver.VolumeResource {
	return virtualserver.VolumeResource{
		Id:            types.StringValue(resp.Id),
		Name:          types.StringPointerValue(resp.Name.Get()),
		UserId:        types.StringValue(resp.UserId),
		Bootable:      types.BoolValue(resp.Bootable),
		Multiattach:   types.BoolValue(resp.Multiattach),
		Encrypted:     types.BoolValue(resp.Encrypted),
		VolumeType:    types.StringValue(resp.VolumeType),
		Size:          types.Int32Value(resp.Size),
		State:         types.StringValue(resp.State),
		Servers:       state.Servers,
		MaxThroughput: types.Int32PointerValue(resp.MaxThroughput.Get()),
		MaxIops:       types.Int32PointerValue(resp.MaxIops.Get()),
		Tags:          tagsMap,
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
			"Could not create volume, unexpected error: "+err.Error()+reasonPrefix+detail,
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
					"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
				)
				return
			}
		}
	}

	getVolumeFunc := func(id string) (*scpvirtualserver.VolumeShowResponseV1Dot2, error) {
		return r.client.GetVolume(ctx, id)
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, data.Id, 10, 3*time.Second,
		"State", "available", "error", getVolumeFunc)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	getData, err := r.client.GetVolume(ctx, data.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading volume",
			"Could not create volume, unexpected error: "+err.Error()+reasonPrefix+detail,
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
			"Could not read volume ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
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

func (r *virtualServerVolumeResource) updateVolumeName(
	ctx context.Context,
	plan virtualserver.VolumeResource,
	state virtualserver.VolumeResource,
	resp *resource.UpdateResponse,
) bool {
	if plan.Name.Equal(state.Name) {
		return true
	}

	_, err := r.client.UpdateVolume(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating volume",
			"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return false
	}
	return true
}

func (r *virtualServerVolumeResource) updateVolumeSize(
	ctx context.Context,
	plan virtualserver.VolumeResource,
	state virtualserver.VolumeResource,
	resp *resource.UpdateResponse,
) bool {
	if plan.Size.Equal(state.Size) {
		return true
	}

	_, err := r.client.ExtendVolume(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating volume",
			"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return false
	}
	return true
}

func (r *virtualServerVolumeResource) updateVolumeServers(
	ctx context.Context,
	plan virtualserver.VolumeResource,
	state virtualserver.VolumeResource,
	resp *resource.UpdateResponse,
) bool {
	if reflect.DeepEqual(plan.Servers, state.Servers) {
		return true
	}

	addedVmIds, deletedVmIds := getOldAndNewVmIds(plan, state)

	for _, deletedVmId := range deletedVmIds {
		err := r.client.DetachVolume(ctx, state.Id.ValueString(), deletedVmId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating volume",
				"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return false
		}
	}

	for _, addedVmId := range addedVmIds {
		_, err := r.client.AttachVolume(ctx, state.Id.ValueString(), addedVmId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating volume",
				"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return false
		}
	}

	return true
}

func (r *virtualServerVolumeResource) updateVolumeQos(
	ctx context.Context,
	plan virtualserver.VolumeResource,
	state virtualserver.VolumeResource,
	resp *resource.UpdateResponse,
) bool {
	if plan.MaxThroughput.Equal(state.MaxThroughput) && plan.MaxIops.Equal(state.MaxIops) {
		return true
	}

	err := r.client.UpdateVolumeQos(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating volume",
			"Could not update volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return false
	}

	if !((plan.MaxThroughput.IsNull() || plan.MaxThroughput.IsUnknown()) && (plan.MaxIops.IsNull() || plan.MaxIops.IsUnknown())) {
		err = r.AsyncPollingQosUpdate(ctx, state.Id.ValueString(), plan.MaxIops.ValueInt32(), plan.MaxThroughput.ValueInt32())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error waiting for volume update",
				"Timed out waiting for values to apply: "+err.Error()+reasonPrefix+detail,
			)
			return false
		}
	}

	return true
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *virtualServerVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state virtualserver.VolumeResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !r.updateVolumeName(ctx, plan, state, resp) {
		return
	}

	if !r.updateVolumeSize(ctx, plan, state, resp) {
		return
	}

	if !r.updateVolumeServers(ctx, plan, state, resp) {
		return
	}

	if !r.updateVolumeQos(ctx, plan, state, resp) {
		return
	}

	data, err := r.client.GetVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading volume",
			"Could not read volume ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
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

func (r *virtualServerVolumeResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
