package virtualserver

import (
	"context"
	"fmt"
	"strings"
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

var (
	_ resource.Resource              = &virtualServerImageResource{}
	_ resource.ResourceWithConfigure = &virtualServerImageResource{}
)

func NewVirtualServerImageResource() resource.Resource {
	return &virtualServerImageResource{}
}

type virtualServerImageResource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (r *virtualServerImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_image"
}

func (r *virtualServerImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Image resource.\n\n" +
			"**Image Creation Methods:**\n" +
			"- URL-based creation: Specify the URL of a qcow2 image file uploaded to Object Storage\n" +
			"- Server snapshot: Create an image from an existing server snapshot (specify instance_id)\n\n",
		MarkdownDescription: "Image resource for virtual servers.\n\n" +
			"**Image Creation Methods:**\n" +
			"- URL-based creation: Specify the URL of a qcow2 image file uploaded to Object Storage\n" +
			"- Server snapshot: Create an image from an existing server snapshot (specify instance_id)\n\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Resource ID.",
				MarkdownDescription: "Resource ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("InstanceId"): schema.StringAttribute{
				Description: "Server ID. Specify when creating an image from an existing server.\n" +
					"  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f\n" +
					"  - note: Do not specify when creating an image via URL.",
				MarkdownDescription: "Server ID. Specify when creating an image from an existing server.\n" +
					"  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f\n" +
					"  - note: Do not specify when creating an image via URL.",
				Optional: true,
			},
			common.ToSnakeCase("Volumes"): schema.StringAttribute{
				Description:         "Volume information.",
				MarkdownDescription: "Volume information.",
				Computed:            true,
			},
			common.ToSnakeCase("Checksum"): schema.StringAttribute{
				Description:         "MD5 checksum of image data. Used for image integrity verification.",
				MarkdownDescription: "MD5 checksum of image data. Used for image integrity verification.",
				Computed:            true,
			},
			common.ToSnakeCase("ContainerFormat"): schema.StringAttribute{
				Description:         "Container format.",
				MarkdownDescription: "Container format.\n  - example: bare",
				Optional:            true,
				Computed:            true,
			},
			common.ToSnakeCase("DiskFormat"): schema.StringAttribute{
				Description:         "Disk format.",
				MarkdownDescription: "Disk format.\n  - example: qcow2",
				Optional:            true,
				Computed:            true,
			},
			common.ToSnakeCase("File"): schema.StringAttribute{
				Description:         "Image file URL.",
				MarkdownDescription: "Image file URL.",
				Computed:            true,
			},
			common.ToSnakeCase("MinDisk"): schema.Int32Attribute{
				Description:         "Minimum disk size (GB).",
				MarkdownDescription: "Minimum disk size (GB).\n  - example: 100",
				Optional:            true,
				Computed:            true,
			},
			common.ToSnakeCase("MinRam"): schema.Int32Attribute{
				Description:         "Minimum RAM size (MB).",
				MarkdownDescription: "Minimum RAM size (MB).\n  - example: 2048",
				Optional:            true,
				Computed:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Image name.\n" +
					"  - example: ubuntu-22.04\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				MarkdownDescription: "Image name.\n" +
					"  - example: ubuntu-22.04\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				Required: true,
			},
			common.ToSnakeCase("OsDistro"): schema.StringAttribute{
				Description: "OS distribution.\n" +
					"  - example: ubuntu\n" +
					"  - Available values: alma, centos, rhel, rocky, ubuntu, windows, oracle",
				MarkdownDescription: "OS distribution.\n" +
					"  - example: ubuntu\n" +
					"  - Available values: alma, centos, rhel, rocky, ubuntu, windows, oracle",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("OsHashAlgo"): schema.StringAttribute{
				Description:         "Hash algorithm for image integrity verification.",
				MarkdownDescription: "Hash algorithm for image integrity verification.",
				Computed:            true,
			},
			common.ToSnakeCase("OsHashValue"): schema.StringAttribute{
				Description:         "Hash value of image binary.",
				MarkdownDescription: "Hash value of image binary.",
				Computed:            true,
			},
			common.ToSnakeCase("OsHidden"): schema.BoolAttribute{
				Description:         "Image hidden status.",
				MarkdownDescription: "Image hidden status.",
				Computed:            true,
			},
			common.ToSnakeCase("Owner"): schema.StringAttribute{
				Description:         "Owner account ID.",
				MarkdownDescription: "Owner account ID.",
				Computed:            true,
			},
			common.ToSnakeCase("OwnerAccountName"): schema.StringAttribute{
				Description:         "Owner account name.",
				MarkdownDescription: "Owner account name.",
				Computed:            true,
			},
			common.ToSnakeCase("OwnerUserName"): schema.StringAttribute{
				Description:         "Owner user name.",
				MarkdownDescription: "Owner user name.",
				Computed:            true,
			},
			common.ToSnakeCase("Protected"): schema.BoolAttribute{
				Description:         "Deletion protection. When set to true, prevents image deletion.",
				MarkdownDescription: "Deletion protection. When set to true, prevents image deletion.\n  - example: false",
				Optional:            true,
				Computed:            true,
			},
			common.ToSnakeCase("RootDeviceName"): schema.StringAttribute{
				Description:         "Root device name.",
				MarkdownDescription: "Root device name.",
				Computed:            true,
			},
			common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
				Description:         "SCP image type.",
				MarkdownDescription: "SCP image type.",
				Computed:            true,
			},
			common.ToSnakeCase("ScpK8sVersion"): schema.StringAttribute{
				Description:         "K8S version. Only has value for K8S images.",
				MarkdownDescription: "K8S version. Only has value for K8S images.",
				Computed:            true,
			},
			common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
				Description:         "Original image type.",
				MarkdownDescription: "Original image type.",
				Computed:            true,
			},
			common.ToSnakeCase("ScpOsVersion"): schema.StringAttribute{
				Description:         "OS version.",
				MarkdownDescription: "OS version.",
				Computed:            true,
			},
			common.ToSnakeCase("Size"): schema.Int64Attribute{
				Description:         "Image size (bytes).",
				MarkdownDescription: "Image size (bytes).",
				Computed:            true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description:         "Image status.",
				MarkdownDescription: "Image status.",
				Computed:            true,
			},
			common.ToSnakeCase("VirtualSize"): schema.Int64Attribute{
				Description:         "Virtual disk size (bytes).",
				MarkdownDescription: "Virtual disk size (bytes).",
				Computed:            true,
			},
			common.ToSnakeCase("Visibility"): schema.StringAttribute{
				Description: "Image visibility.\n" +
					"  - example: private\n" +
					"  - Available values: shared, private",
				MarkdownDescription: "Image visibility.\n" +
					"  - example: private\n" +
					"  - Available values: shared, private",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Url"): schema.StringAttribute{
				Description: "Object Storage URL. Only qcow2 format is allowed.\n" +
					"  - example: https://object-store.kr-west1.s.samsungsdscloud.com/bucket/image.qcow2\n" +
					"  - note: Specify only when creating an image via URL.",
				MarkdownDescription: "Object Storage URL. Only qcow2 format is allowed.\n" +
					"  - example: https://object-store.kr-west1.s.samsungsdscloud.com/bucket/image.qcow2\n" +
					"  - note: Specify only when creating an image via URL.",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description:         "Creation timestamp.",
				MarkdownDescription: "Creation timestamp.",
				Computed:            true,
			},
			common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
				Description:         "Update timestamp.",
				MarkdownDescription: "Update timestamp.",
				Computed:            true,
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

func (r *virtualServerImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *virtualServerImageResource) MapGetResponseToState(resp *scpvirtualserver.ImageShowResponseV1Dot2, tagsMap types.Map) virtualserver.ImageResource {
	return virtualserver.ImageResource{
		Volumes:              types.StringPointerValue(resp.Volumes.Get()),
		Checksum:             types.StringPointerValue(resp.Checksum.Get()),
		ContainerFormat:      types.StringValue(resp.ContainerFormat),
		DiskFormat:           types.StringValue(resp.DiskFormat),
		File:                 types.StringValue(resp.File),
		Id:                   types.StringValue(resp.Id),
		MinDisk:              types.Int32Value(resp.MinDisk),
		MinRam:               types.Int32Value(resp.MinRam),
		Name:                 types.StringValue(resp.Name),
		OsDistro:             types.StringPointerValue(resp.OsDistro.Get()),
		OsHashAlgo:           types.StringPointerValue(resp.OsHashAlgo.Get()),
		OsHashValue:          types.StringPointerValue(resp.OsHashValue.Get()),
		OsHidden:             types.BoolValue(resp.OsHidden),
		Owner:                types.StringValue(resp.Owner),
		OwnerAccountName:     types.StringPointerValue(resp.OwnerAccountName.Get()),
		OwnerUserName:        types.StringPointerValue(resp.OwnerUserName.Get()),
		Protected:            types.BoolValue(resp.Protected),
		RootDeviceName:       types.StringPointerValue(resp.RootDeviceName.Get()),
		ScpImageType:         types.StringPointerValue(resp.ScpImageType.Get()),
		ScpK8sVersion:        types.StringPointerValue(resp.ScpK8sVersion.Get()),
		ScpOriginalImageType: types.StringPointerValue(resp.ScpOriginalImageType.Get()),
		ScpOsVersion:         types.StringPointerValue(resp.ScpOsVersion.Get()),
		Size:                 types.Int64PointerValue(resp.Size.Get()),
		Status:               types.StringValue(resp.Status),
		VirtualSize:          types.Int64PointerValue(resp.VirtualSize.Get()),
		Visibility:           types.StringValue(resp.Visibility),
		Url:                  types.StringPointerValue(resp.Url.Get()),
		CreatedAt:            types.StringValue(resp.CreatedAt),
		UpdatedAt:            types.StringValue(resp.UpdatedAt),
		Tags:                 tagsMap,
	}
}

func (r *virtualServerImageResource) handlerUpdateImage(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ImageResource
	req.Plan.Get(ctx, &plan)

	_, err := r.client.UpdateImage(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		return err
	}
	return nil
}

func (r *virtualServerImageResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ImageResource
	var state virtualserver.ImageResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	serviceName, resourceType := r.resolveImageServiceInfoFromModel(state)
	_, err := tag.UpdateTags(r.clients, serviceName, resourceType, plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		return err
	}
	return nil
}

func (r *virtualServerImageResource) resolveImageServiceInfoFromResponse(response *scpvirtualserver.ImageShowResponseV1Dot2) (serviceName, resourceType string) {
	if response.ScpImageType.Get() != nil && *response.ScpImageType.Get() == ScpImageTypeGpuCustom {
		return ServiceNameGpuServer, ResourceTypeImage
	}
	return ServiceNameVirtualServer, ResourceTypeImage
}

func (r *virtualServerImageResource) resolveImageServiceInfoFromModel(model virtualserver.ImageResource) (serviceName, resourceType string) {
	if !model.ScpImageType.IsNull() && model.ScpImageType.ValueString() == ScpImageTypeGpuCustom {
		return ServiceNameGpuServer, ResourceTypeImage
	}
	return ServiceNameVirtualServer, ResourceTypeImage
}

func (r *virtualServerImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualserver.ImageResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var imageId string

	// Server ID를 통한 이미지 생성
	if !plan.InstanceId.IsNull() {
		data, err := r.client.CreateImageFromServer(ctx, plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error creating image",
				"Could not create image, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		imageId = data.ImageId
	}

	// URL을 통한 이미지 생성
	if plan.InstanceId.IsNull() {
		data, err := r.client.CreateImage(ctx, plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error creating image",
				"Could not create image, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		imageId = data.Id
	}

	getFunc := func(id string) (*scpvirtualserver.ImageShowResponseV1Dot2, error) {
		return r.client.GetImage(ctx, id)
	}

	getData, err := virtualserverutil.AsyncRequestPollingWithState(ctx, imageId, 10, 3*time.Second,
		"Status", "active", "queued", getFunc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			"Could not create image, unexpected error: "+err.Error(),
		)
		return
	}

	serviceName, resourceType := r.resolveImageServiceInfoFromResponse(getData)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, imageId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state := r.MapGetResponseToState(getData, tagsMap)
	if !plan.InstanceId.IsNull() {
		state.InstanceId = plan.InstanceId
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualserver.ImageResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetImage(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading image",
			"Could not read image id "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	serviceName, resourceType := r.resolveImageServiceInfoFromResponse(data)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState := r.MapGetResponseToState(data, tagsMap)
	if !state.InstanceId.IsNull() {
		newState.InstanceId = state.InstanceId
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*virtualserver.UpdateHandler{
		{
			Fields:  []string{"MinDisk", "MinRam", "Protected", "Visibility"},
			Handler: r.handlerUpdateImage,
		},
		{
			Fields:  []string{"Tags"},
			Handler: r.handlerUpdateTag,
		},
	}

	var plan virtualserver.ImageResource
	var state virtualserver.ImageResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settableFileds []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFileds = append(settableFileds, virtualserverutil.SnakeToPascal(attrName))
		}
	}

	changeFields, err := virtualserverutil.GetChangedFields(plan, state, settableFileds)
	if err != nil {
		return
	}

	immuntableFields := []string{"Name", "OsDistro", "DiskFormat", "ContainerFormat", "Url", "InstanceId"}

	if virtualserverutil.IsOverlapFields(immuntableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Image",
			"Immutable fields cannot be modified: "+strings.Join(immuntableFields, ", "),
		)
		return
	}

	for _, h := range handlers {
		if virtualserverutil.IsOverlapFields(h.Fields, changeFields) {
			if err := h.Handler(ctx, req, resp); err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Image",
					"Could not update image, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	}

	getData, err := r.client.GetImage(ctx, plan.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Image",
			"Could not read image, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	serviceName, resourceType := r.resolveImageServiceInfoFromResponse(getData)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState := r.MapGetResponseToState(getData, tagsMap)
	if !state.InstanceId.IsNull() {
		newState.InstanceId = state.InstanceId
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *virtualServerImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualserver.ImageResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteImage(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting image",
			"Could not delete image, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *virtualServerImageResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
