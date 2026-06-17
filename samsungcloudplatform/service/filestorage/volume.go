package filestorage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/filestorage/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &fileStorageVolumeResource{}
	_ resource.ResourceWithConfigure   = &fileStorageVolumeResource{}
	_ resource.ResourceWithImportState = &fileStorageVolumeResource{}
)

func NewFileStorageVolumeResource() resource.Resource {
	return &fileStorageVolumeResource{}
}

type fileStorageVolumeResource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (r *fileStorageVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filestorage_volume"
}

func (r *fileStorageVolumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VolumeResourceSchema()
}
func VolumeResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a File Storage Volume on Samsung Cloud Platform.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed: true,
				Description: "Account ID \n" +
					"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
			},
			"cifs_password": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
				Sensitive: true,
				Description: "Cifs Password \n" +
					"  - example : '<YOUR_CIFS_PASSWORD>' \n" +
					"  - maxLength: 20  \n" +
					"  - minLength: 6  \n" +
					"  - pattern: `^(?=.*[a-zA-Z])(?=.*\\d)(?=.*[!#&\\'*+,-.:;<=>?@^_`~/|])[a-zA-Z\\d!#&\\'*+,-.:;<=>?@^_`~/|]{6,20}$` \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Created At \n" +
					"  - example : '2024-07-30T04:54:33.219373' \n",
			},
			"encryption_enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Volume Encryption Enabled \n" +
					"  - example : true",
			},
			"endpoint_path": schema.StringAttribute{
				Computed: true,
				Description: "The network endpoint path used to mount and access the file storage volume. \n" +
					"  - example : 'xxx.xx.xxx.xxx' \n",
			},
			"file_unit_recovery_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Description: "File Unit Recovery Enabled \n" +
					"  - example : true \n",
			},
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Identifier of the resource. \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				// planmodifier 별도 추가
				PlanModifiers: []planmodifier.String{ //  PlanModifiers 추가
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: "Volume Name \n" +
					"  - example : 'my_volume' \n" +
					"  - maxLength: 21  \n" +
					"  - minLength: 3  \n" +
					"  - pattern: `^[a-z]([a-z0-9_]){2,20}$` \n",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z]([a-z0-9_]){2,20}$"), "Enter 3~21 char.(lower case, numbers, _) starting with lower case."),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("NameUuid"): schema.StringAttribute{
				Description: "The unique system-assigned name (UUID format) for the volume. \n" +
					"  - example : 'my_volume_2m060u' \n",
				Computed: true,
			},
			"path": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Description: "Volume Mount Path \n" +
					"  - example : 'xxx.xx.xxx.xxx'",
			},
			"protocol": schema.StringAttribute{
				Required: true,
				Description: "Protocol \n" +
					"  - example : 'NFS' \n" +
					"  - pattern: `^(NFS|CIFS)$` \n",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^(NFS|CIFS)$"), "Protocol must be one of (NFS, CIFS)."),
				},
			},
			"purpose": schema.StringAttribute{
				Computed: true,
				Description: "The designated purpose or workload type of the volume (e.g., general, backup). \n" +
					"  - example : 'none' \n",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "The current lifecycle state of the volume. Valid values: creating, available, error, deleting.",
				MarkdownDescription: "The current lifecycle state of the volume. Valid values: creating, available, error, deleting.",
			},
			"type_id": schema.StringAttribute{
				Computed: true,
				Description: "The unique identifier of the storage tier (volume type) assigned to this volume. \n" +
					"  - example : 'jef22f67-ee83-4gg2-2ab6-3lf774ekfjdu' \n",
			},
			"type_name": schema.StringAttribute{
				Required: true,
				Description: "Volume Type Name \n" +
					"  - example : 'HDD' \n" +
					"  - pattern: `^(HDD|SSD|HighPerformanceSSD|SSD_SAP_S|SSD_SAP_E)$` \n",
			},
			// 별도로 Optional: true 추가
			"usage": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The current usage of the volume in GiB.",
			},
			// custom으로 추가
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("AccessRules"): schema.SetNestedAttribute{
				Description: "List of AccessRule",
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ObjectId"): schema.StringAttribute{
							Description: "Object Id \n" +
								"  - example : '43fq3347-02q4-4aa8-ccf9-affe4917bb6f' \n",
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						common.ToSnakeCase("ObjectType"): schema.StringAttribute{
							Description: "Object Type \n" +
								"  - example : 'VM' \n" +
								"  - pattern: `^(VM|BM|GPU|GPU_NODE|ENDPOINT)$` \n",
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *fileStorageVolumeResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = inst.Client.FileStorage
	r.clients = inst.Client
}

func (r *fileStorageVolumeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Values from plan
	var plan filestorage.VolumeResource
	diags := request.Config.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Protocol, CifsPassword validation
	if plan.Protocol.ValueString() == "NFS" && plan.CifsPassword.ValueStringPointer() != nil {
		response.Diagnostics.AddError("Error creating volume", "Could not create volume, NFS Protocol doesn't need cifs_password")
		return
	}

	// Create volume
	data, err := r.client.CreateVolume(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Volume State Polling
	getData, err := waitForVolumeStatus(ctx, r.client, data.VolumeId, []string{}, []string{"available", "error"})
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating volume",
			"Error waiting for volume to become available: "+err.Error(),
		)
		return
	}

	if getData.State == "error" {
		response.Diagnostics.AddError(
			"Error creating volume",
			"Error States for volume to become error",
		)
		return
	}

	volume := getData
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(volume.Id)

	// Update Access Rule
	if len(plan.AccessRules) != 0 {
		for _, rule := range plan.AccessRules {
			err := r.client.UpdateVolumeAccessRule(ctx, plan.Id.ValueString(), rule, "add")
			if err != nil {
				detail := client.GetDetailFromError(err)
				response.Diagnostics.AddError("Error Updating AccessRule",
					"Could not update AccessRule, unexpected error: "+err.Error()+"\nReason: "+detail)
				return
			}
		}
	}

	tagsMap, err := tag.GetTags(r.clients, "filestorage", "volume", volume.Id)
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state, err := r.MapGetResponseToState(ctx, volume, plan, tagsMap)
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}

func (r *fileStorageVolumeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state filestorage.VolumeResource
	diags := request.State.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetVolume(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			response.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Volume",
			"Could not read Volume ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "filestorage", "volume", state.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, state, tagsMap)

	// Set refreshed state
	diags = response.State.Set(ctx, &newState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *fileStorageVolumeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state filestorage.VolumeResource

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	//--------------- Volume Update ---------------//
	err := r.client.UpdateVolume(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Updating Volume",
			"Could not update Volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	data, detailErr := r.client.GetVolume(ctx, plan.Id.ValueString())
	if detailErr != nil {
		detail := client.GetDetailFromError(detailErr)
		response.Diagnostics.AddError(
			"Error Reading Volume",
			"Could not read Volume ID "+plan.Id.ValueString()+": "+detailErr.Error()+"\nReason: "+detail,
		)
		return
	}

	//--------------- AccessRule Update ---------------//
	addRule, removeRule := r.ProcessAccessRules(state.AccessRules, plan.AccessRules)
	if addRule != nil {
		for _, rule := range addRule {
			err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), rule, "add")
			if err != nil {
				detail := client.GetDetailFromError(err)
				response.Diagnostics.AddError("Error Updating AccessRule",
					"Could not update AccessRule, unexpected error: "+err.Error()+"\nReason: "+detail)
				return
			}
		}
	}

	if removeRule != nil {
		for _, rule := range removeRule {
			err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), rule, "remove")
			if err != nil {
				detail := client.GetDetailFromError(err)
				response.Diagnostics.AddError("Error Updating AccessRule",
					"Could not update AccessRule, unexpected error: "+err.Error()+"\nReason: "+detail)
				return
			}
		}
	}

	tagsMap, err := tag.GetTags(r.clients, "filestorage", "volume", state.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	newState, err := r.MapGetResponseToState(ctx, data, state, tagsMap)
	diags = response.State.Set(ctx, newState)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
}

func (r *fileStorageVolumeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state filestorage.VolumeResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Check Access Rules
	if len(state.AccessRules) != 0 {
		// Update AccessRules (remove)
		for _, rule := range state.AccessRules {
			err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), rule, "remove")
			if err != nil {
				detail := client.GetDetailFromError(err)
				response.Diagnostics.AddError("Error Updating AccessRule",
					"Could not update AccessRule, unexpected error: "+err.Error()+"\nReason: "+detail)
				return
			}
		}
	}

	// Delete
	err := r.client.DeleteVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Deleting Volume",
			"Could not delete Volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *fileStorageVolumeResource) ProcessAccessRules(stateRules, planRules []filestorage.AccessRuleResource) ([]filestorage.AccessRuleResource, []filestorage.AccessRuleResource) {
	existingRules := make(map[string]filestorage.AccessRuleResource)
	for _, rule := range stateRules {
		existingRules[rule.ObjectId.ValueString()] = rule
	}

	var toAdd, toRemove []filestorage.AccessRuleResource

	// Add or Existing Access Rule
	for _, planRule := range planRules {
		objectId := planRule.ObjectId.ValueString()
		if _, exists := existingRules[objectId]; !exists {
			toAdd = append(toAdd, planRule)
		}
		delete(existingRules, objectId)
	}

	// Check Remove Access Rule
	for _, rule := range existingRules {
		toRemove = append(toRemove, rule)
	}

	return toAdd, toRemove
}

func (r *fileStorageVolumeResource) MapGetResponseToState(ctx context.Context, resp *scpfilestorage.VolumeShowResponse, state filestorage.VolumeResource, tagsMap types.Map) (filestorage.VolumeResource, error) {

	// AccessRule
	getAccessRule, err := r.client.GetVolumeAccessRules(ctx, resp.Id)

	if err != nil {
		return filestorage.VolumeResource{}, err
	}

	var accessRules []filestorage.AccessRuleResource
	if len(getAccessRule.AccessRules) == 0 && state.AccessRules != nil {
		accessRules = []filestorage.AccessRuleResource{}
	} else {
		for _, rules := range getAccessRule.AccessRules {
			rule := filestorage.AccessRuleResource{
				ObjectId:   types.StringValue(rules.ObjectId),
				ObjectType: types.StringValue(rules.ObjectType),
			}
			accessRules = append(accessRules, rule)
		}
	}

	endpointPath := types.StringNull()
	if v := resp.EndpointPath.Get(); v != nil {
		endpointPath = types.StringValue(*v)
	}

	pathValue := types.StringNull()
	if v := resp.Path.Get(); v != nil {
		pathValue = types.StringValue(*v)
	}

	usage := types.Int64Null()
	if v := resp.Usage.Get(); v != nil {
		usage = types.Int64Value(*v)
	}

	return filestorage.VolumeResource{
		AccountId:               types.StringValue(resp.AccountId),
		CreatedAt:               types.StringValue(resp.CreatedAt.Format(time.RFC3339)),
		EncryptionEnabled:       types.BoolValue(resp.EncryptionEnabled),
		EndpointPath:            endpointPath,
		FileUnitRecoveryEnabled: types.BoolValue(resp.GetFileUnitRecoveryEnabled()),
		Id:                      types.StringValue(resp.Id),
		Name:                    types.StringValue(state.Name.ValueString()),
		NameUuid:                types.StringValue(resp.Name),
		Path:                    pathValue,
		Protocol:                types.StringValue(resp.Protocol),
		Purpose:                 types.StringValue(resp.Purpose),
		State:                   types.StringValue(resp.State),
		TypeId:                  types.StringValue(resp.TypeId),
		TypeName:                types.StringValue(resp.TypeName),
		Usage:                   usage,
		Tags:                    tagsMap,
		AccessRules:             accessRules,
	}, nil
}

func waitForVolumeStatus(ctx context.Context, fileStorageClient *filestorage.Client, id string, pendingStates []string, targetStates []string) (*scpfilestorage.VolumeShowResponse, error) {
	var showResponse *scpfilestorage.VolumeShowResponse
	err := client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := fileStorageClient.GetVolume(ctx, id)
		if err != nil {
			return nil, "", err
		}
		showResponse = info
		return info, info.State, nil
	}, -1, -1, -1, -1)
	return showResponse, err
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *fileStorageVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
