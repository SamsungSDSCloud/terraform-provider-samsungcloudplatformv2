package parallelfilestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/parallelfilestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpparallelfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/parallelfilestorage/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

const reasonPrefix = "\nReason: "

var (
	_ resource.Resource              =   &parallelFileStorageVolumeResource{}
	_ resource.ResourceWithConfigure =   &parallelFileStorageVolumeResource{}
	_ resource.ResourceWithImportState = &parallelFileStorageVolumeResource{}
)

func NewParallelFileStorageVolumeResource() resource.Resource {
	return &parallelFileStorageVolumeResource{}
}

type parallelFileStorageVolumeResource struct {
	config  *scpsdk.Configuration
	client  *parallelfilestorage.Client
	clients *client.SCPClient
}

func (r *parallelFileStorageVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parallel_filestorage_volume"
}

func (r *parallelFileStorageVolumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VolumeResourceSchema()
}
func VolumeResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed: true,
				Description: "Account ID \n" +
					"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Created At \n" +
					"  - example : '2024-07-30T04:54:33.219373' \n",
			},
			"capacity_tb": schema.Int32Attribute{
				Description: "Volume capacity(TB). \n" +
					"  - example : 10 \n" +
					"  - maximum : 1000 \n" +
					"  - minimum : 1 \n",
				Required: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 1000),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the resource.",
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
			"name_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "Actual Volume Name on the server (may include a unique suffix appended by the server).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "Zone",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Volume State",
				MarkdownDescription: "Volume State",
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("AccessRules"): schema.SetNestedAttribute{
				Description: "Object of AccessRule",
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

func (r *parallelFileStorageVolumeResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

	r.client = inst.Client.ParallelFileStorage
	r.clients = inst.Client
}

func (r *parallelFileStorageVolumeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Values from plan
	var plan parallelfilestorage.VolumeResource

	diags := request.Config.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Create volume
	data, err := r.client.CreateVolume(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	// Volume State Polling
	getData, err := waitForVolumeStatus(ctx, r.client, data.Id, []string{}, []string{"available", "error"})
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
		err := r.client.UpdateVolumeAccessRule(ctx, plan.Id.ValueString(), plan.AccessRules, "add")
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error Updating AccessRule",
				"Could not update AccessRule, unexpected error: "+err.Error()+reasonPrefix+detail)
			return
		}
	}
	tagsMap, err := tag.GetTags(r.clients, "parallel-filestorage", "volume", volume.Id, false)

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

func (r *parallelFileStorageVolumeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state parallelfilestorage.VolumeResource
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
			"Could not read Volume ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "parallel-filestorage", "volume", state.Id.ValueString(), false)
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

func (r *parallelFileStorageVolumeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state parallelfilestorage.VolumeResource
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
	if plan.CapacityTb != state.CapacityTb {
		err := r.client.UpdateVolumeCapacityTb(ctx, plan.Id.ValueString(), plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError(
				"Error Updating Volume Capacity Tb",
				"Could not update Volume, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	data, detailErr := r.client.GetVolume(ctx, plan.Id.ValueString())
	if detailErr != nil {
		detail := client.GetDetailFromError(detailErr)
		response.Diagnostics.AddError(
			"Error Reading Volume Capacity Tb",
			"Could not read Volume ID "+plan.Id.ValueString()+": "+detailErr.Error()+reasonPrefix+detail,
		)
		return
	}

	//--------------- AccessRule Update ---------------//
	addRule, removeRule := r.ProcessAccessRules(state.AccessRules, plan.AccessRules)
	if addRule != nil {
		err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), addRule, "add")
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error Updating AccessRule",
				"Could not update AccessRule, unexpected error: "+err.Error()+reasonPrefix+detail)
			return
		}
	}

	if removeRule != nil {
		err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), removeRule, "remove")
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error Updating AccessRule",
				"Could not update AccessRule, unexpected error: "+err.Error()+reasonPrefix+detail)
			return
		}
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "parallel-filestorage", "volume", plan.Id.ValueString(), tagElements, false)
	if err != nil {
		response.Diagnostics.AddError(
			"Error Updating Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, plan, tagsMap)
	diags = response.State.Set(ctx, newState)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
}

func (r *parallelFileStorageVolumeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state parallelfilestorage.VolumeResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	// Check Access Rules
	if len(state.AccessRules) != 0 {
		// Update AccessRules (remove)
		err := r.client.UpdateVolumeAccessRule(ctx, state.Id.ValueString(), state.AccessRules, "remove")
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error Updating AccessRule",
				"Could not update AccessRule, unexpected error: "+err.Error()+reasonPrefix+detail)
			return
		}
	}

	// Delete
	err := r.client.DeleteVolume(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Deleting Volume",
			"Could not delete Volume, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}
}

func (r *parallelFileStorageVolumeResource) ProcessAccessRules(stateRules, planRules []parallelfilestorage.AccessRuleResource) ([]parallelfilestorage.AccessRuleResource, []parallelfilestorage.AccessRuleResource) {
	existingRules := make(map[string]parallelfilestorage.AccessRuleResource)
	for _, rule := range stateRules {
		existingRules[rule.ObjectId.ValueString()] = rule
	}

	var toAdd, toRemove []parallelfilestorage.AccessRuleResource

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

func (r *parallelFileStorageVolumeResource) MapGetResponseToState(ctx context.Context, resp *scpparallelfilestorage.VolumeShowResponseV1Dot1, state parallelfilestorage.VolumeResource, tagsMap types.Map) (parallelfilestorage.VolumeResource, error) {

	// AccessRule
	getAccessRule, err := r.client.GetVolumeAccessRules(ctx, resp.Id)

	if err != nil {
		return parallelfilestorage.VolumeResource{}, err
	}

	accessRules := []parallelfilestorage.AccessRuleResource{}
	if len(getAccessRule.AccessRules) == 0 && state.AccessRules != nil {
		accessRules = []parallelfilestorage.AccessRuleResource{}
	} else {
		for _, rules := range getAccessRule.AccessRules {
			rule := parallelfilestorage.AccessRuleResource{
				ObjectId:   types.StringValue(rules.ObjectId),
				ObjectType: types.StringValue(rules.ObjectType),
			}
			accessRules = append(accessRules, rule)
		}
	}
	return parallelfilestorage.VolumeResource{
		AccountId:   types.StringValue(resp.AccountId),
		CreatedAt:   types.StringValue(resp.CreatedAt.Format(time.RFC3339)),
		Id:          types.StringValue(resp.Id),
		Name:        types.StringValue(state.Name.ValueString()),
		NameUuid:    types.StringValue(resp.Name),
		State:       types.StringValue(resp.State),
		Zone:        types.StringValue(resp.Zone),
		CapacityTb:  types.Int32Value(resp.CapacityTb),
		Tags:        tagsMap,
		AccessRules: accessRules,
	}, nil
}

func waitForVolumeStatus(ctx context.Context, parallelFilestorageClient *parallelfilestorage.Client, id string, pendingStates []string, targetStates []string) (*scpparallelfilestorage.VolumeShowResponseV1Dot1, error) {
	var showResponse *scpparallelfilestorage.VolumeShowResponseV1Dot1
	err := client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := parallelFilestorageClient.GetVolume(ctx, id)
		if err != nil {
			return nil, "", err
		}
		showResponse = info
		return info, info.State, nil
	}, -1, -1, -1, -1)
	return showResponse, err
}

func (r *parallelFileStorageVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: <id>,<name>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}
