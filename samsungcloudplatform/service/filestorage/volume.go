package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/filestorage/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource              = &fileStorageVolumeResource{}
	_ resource.ResourceWithConfigure = &fileStorageVolumeResource{}
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

func (r *fileStorageVolumeResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "volume",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID \n" +
					"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
				Computed: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "Created At \n" +
					"  - example : '2024-07-30T04:54:33.219373' \n",
				Computed: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Volume Name \n" +
					"  - example : 'my_volume' \n" +
					"  - maxLength: 21  \n" +
					"  - minLength: 3  \n" +
					"  - pattern: '^[a-z]([a-z0-9_]){2,20}$' \n",
				Required: true,
			},
			common.ToSnakeCase("NameUuid"): schema.StringAttribute{
				Description: "Volume Name Uuid \n" +
					"  - example : 'my_volume_2m060u' \n",
				Computed: true,
			},
			common.ToSnakeCase("Protocol"): schema.StringAttribute{
				Description: "Protocol \n" +
					"  - example : 'NFS' \n" +
					"  - pattern: '^(NFS|CIFS)$' \n",
				Required: true,
			},
			common.ToSnakeCase("Purpose"): schema.StringAttribute{
				Description: "Volume Purpose \n" +
					"  - example : 'none' \n",
				Computed: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "Volume State \n" +
					"  - example : 'available' \n",
				Computed: true,
			},
			common.ToSnakeCase("TypeId"): schema.StringAttribute{
				Description: "Volume Type ID \n" +
					"  - example : 'jef22f67-ee83-4gg2-2ab6-3lf774ekfjdu' \n",
				Computed: true,
			},
			common.ToSnakeCase("TypeName"): schema.StringAttribute{
				Description: "Volume Type Name \n" +
					"  - example : 'HDD' \n" +
					"  - pattern: '^(HDD|SSD|HighPerformanceSSD)$' \n",
				Required: true,
			},
			common.ToSnakeCase("FileUnitRecoveryEnabled"): schema.BoolAttribute{
				Description: "File Unit Recovery Enabled \n" +
					"  - example : 'true' \n",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("CifsPassword"): schema.StringAttribute{
				Description: "Cifs Password \n" +
					"  - example : 'cifspwd0!!' \n" +
					"  - maxLength: 20  \n" +
					"  - minLength: 6  \n" +
					"  - pattern: '^(?=.*[a-zA-Z])(?=.*\\d)(?=.*[!#&\\'*+,-.:;<=>?@^_`~/|])[a-zA-Z\\d!#&\\'*+,-.:;<=>?@^_`~/|]{6,20}$' \n",
				Optional:  true,
				WriteOnly: true,
			},
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
							Description: "Object Type" +
								"  - example : 'VM' \n" +
								"  - pattern: '^(VM|BM|GPU|GPU_NODE|ENDPOINT)$' \n",
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

	return filestorage.VolumeResource{
		AccountId:               types.StringValue(resp.AccountId),
		CreatedAt:               types.StringValue(resp.CreatedAt.Format(time.RFC3339)),
		Id:                      types.StringValue(resp.Id),
		Name:                    types.StringValue(state.Name.ValueString()),
		NameUuid:                types.StringValue(resp.Name),
		Protocol:                types.StringValue(resp.Protocol),
		Purpose:                 types.StringValue(resp.Purpose),
		State:                   types.StringValue(resp.State),
		TypeId:                  types.StringValue(resp.TypeId),
		TypeName:                types.StringValue(resp.TypeName),
		FileUnitRecoveryEnabled: types.BoolValue(resp.GetFileUnitRecoveryEnabled()),
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
	})
	return showResponse, err
}
