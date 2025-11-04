package baremetalblockstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/baremetalblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.ResourceWithConfigure  = &baremetalBlockStorageVolume{}
	_ resource.ResourceWithModifyPlan = &baremetalBlockStorageVolume{}
)

func NewBaremetalBlockStorageVolumeResource() resource.Resource {
	return &baremetalBlockStorageVolume{}
}

type baremetalBlockStorageVolume struct {
	config  *scpsdk.Configuration
	client  *baremetalblockstorage.Client
	clients *client.SCPClient
}

func (r *baremetalBlockStorageVolume) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_baremetal_blockstorage_volume"
}

func (r *baremetalBlockStorageVolume) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Block Storage(BM)",
		Attributes: map[string]schema.Attribute{
			"region": region.ResourceSchema(),
			"tags":   tag.ResourceSchema(),
			common.ToSnakeCase("id"): schema.StringAttribute{
				Description: "Volume id. \n" +
					"  - example: 8bf55e738d4e44b5a21dbe133a42ecbe \n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("name"): schema.StringAttribute{
				Description: "Volume name. \n" +
					"  - example : my-bs-01 \n" +
					"  - maxLength : 28 \n" +
					"  - minLength : 3 \n" +
					"  - pattern : '^[a-zA-Z][a-zA-Z0-9-]*$'\n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 28),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*$"),
						"It starts with English and must end in numbers, '-', or English."),
				},
			},
			common.ToSnakeCase("diskType"): schema.StringAttribute{
				Description: "Disk type. \n" +
					"  - example : SSD \n" +
					"  - pattern : SSD|HDD \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SSD", "HDD"),
				},
			},
			common.ToSnakeCase("sizeGb"): schema.Int32Attribute{
				Description: "Volume capacity(GB). \n" +
					"  - example : 10 \n" +
					"  - maximum : 16384 \n" +
					"  - minimum : 1 \n",
				Required: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 16384),
				},
			},
			common.ToSnakeCase("attachments"): schema.ListNestedAttribute{
				Description: "List of server id and type. \n" +
					"  - example : [{object_type='BM', object_id='83c3c73d457345e3829ee6d5557c0011'}] \n" +
					"  - maxLength : 5 \n" +
					"  - minLength : 1 \n",

				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"object_type": schema.StringAttribute{
							Description: "Object type. \n" +
								"  - example : 'BM' \n" +
								"  - pattern : 'BM|MNGC'",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("BM", "MNGC"),
							},
						},
						"object_id": schema.StringAttribute{
							Description: "Object id. \n" +
								"  - example : 83c3c73d457345e3829ee6d5557c0016",
							Required: true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(0, 5),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *baremetalBlockStorageVolume) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var state, plan baremetalblockstorage.VolumeResource
	diags := request.Plan.Get(ctx, &plan)
	// for destroy
	if reflect.ValueOf(plan).IsZero() {
		response.Diagnostics.Append(request.State.Get(ctx, &state)...)
		if response.Diagnostics.HasError() {
			return
		}
		if len(state.Attachments) > 0 {
			response.Diagnostics.AddError("Could not delete Block Storage(BM)",
				"Could not delete Block Storage(BM).\nAttachments must be all detach.")
		}
		return
	}

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// for creating
	if plan.Id.IsUnknown() {
		return
	}

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !state.Name.Equal(plan.Name) {
		response.Diagnostics.AddError("Could not change name",
			"Could not change name.\nIf you want to change, create a new resource or change it through the console.")
	}

	if !state.SizeGb.Equal(plan.SizeGb) {
		response.Diagnostics.AddError("Could not change size_gb",
			"Could not change size_gb.\nIf you want to change, create a new resource.")
	}
	if !state.DiskType.Equal(plan.DiskType) {
		response.Diagnostics.AddError("Could not change disk_type",
			"Could not change disk_type.\nIf you want to change, create a new resource.")
	}
	if !state.Region.Equal(plan.Region) {
		response.Diagnostics.AddError("Could not change region",
			"Could not change region.\nIf you want to change, create a new resource.")
	}
}

func (r *baremetalBlockStorageVolume) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan baremetalblockstorage.VolumeResource
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Region.IsNull() {
		r.client.Config.Region = plan.Region.ValueString()
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 30*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	timeoutErrCause, _ := json.Marshal(map[string]interface{}{
		"errors": map[string]interface{}{
			"detail": "Create Timeout.",
		},
	})
	ctx, cancel := context.WithTimeoutCause(ctx, createTimeout, scpsdk.GenericOpenAPIError{
		ResponseBody: timeoutErrCause,
		ErrorMessage: "Create Timeout.",
	})
	defer cancel()

	asyncResponse, err := r.client.CreateBlockStorage(ctx, plan)
	if err != nil {
		var detail string
		if reflect.PointerTo(reflect.TypeOf(scpsdk.GenericOpenAPIError{})) == reflect.TypeOf(err) {
			detail = client.GetDetailFromError(err)
		} else {
			detail = err.Error()
		}
		response.Diagnostics.AddError("Error creating Block Storage(BM)",
			"Could not create Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
		return
	}

	blockStorageId := *asyncResponse.ResourceId

	plan.Id = types.StringValue(blockStorageId)

	err = r.waitForBlockStorageState(ctx, blockStorageId, []string{common.CreatingState}, []string{common.InUseState}, createTimeout)

	if err != nil {
		var detail string
		if reflect.PointerTo(reflect.TypeOf(scpsdk.GenericOpenAPIError{})) == reflect.TypeOf(err) {
			detail = client.GetDetailFromError(err)
		} else {
			detail = err.Error()
		}
		response.Diagnostics.AddError("Error creating Block Storage(BM)",
			"Could not create Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
		return
	}
	response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *baremetalBlockStorageVolume) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state baremetalblockstorage.VolumeResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	volumeResponse, _, err := r.client.GetBlockStorage(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError("Error Reading Block Storage(BM)",
			"Could not read Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
		return
	}
	blockStorage := volumeResponse.Result

	refreshBlockStorageName := types.StringPointerValue(blockStorage.Name)
	if !state.Name.Equal(refreshBlockStorageName) {
		//state.Name = refreshBlockStorageName
		response.Diagnostics.AddWarning("Warning Reading Block Storage(BM)",
			"The name of the Block Storage has been modified.\n"+
				"Cannot modify the name in Terraform.\n"+
				"You must manually modify the \".tfstate\", \".tf\" file or roll back the name in the console.\n"+
				"Current: "+refreshBlockStorageName.String()+"\n"+
				"Previous: "+state.Name.String())
	}

	attachments := make([]baremetalblockstorage.Attachment, 0)
	for _, attachment := range blockStorage.Attachments {
		elems := baremetalblockstorage.Attachment{
			ObjectId:   types.StringPointerValue(attachment.Id),
			ObjectType: types.StringValue(string(*attachment.Type)),
		}

		attachments = append(attachments, elems)
	}
	state.Attachments = attachments
	tagsMap, err := tag.GetTags(r.clients, "baremetal-blockstorage", "volume", *blockStorage.Id)
	if err != nil {
		response.Diagnostics.AddError("Error Reading Block Storage(BM)", err.Error())
	}
	state.Tags = tagsMap

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *baremetalBlockStorageVolume) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state, plan baremetalblockstorage.VolumeResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	beforeAttachmentMap := make(map[string]string)
	for _, attachment := range state.Attachments {
		beforeAttachmentMap[attachment.ObjectId.ValueString()] = attachment.ObjectType.ValueString()
	}

	attachObjectList := make([]baremetalblockstorage.Attachment, 0)
	for _, attachment := range plan.Attachments {
		if _, ok := beforeAttachmentMap[attachment.ObjectId.ValueString()]; !ok {
			attachObjectList = append(attachObjectList, attachment)
		} else {
			delete(beforeAttachmentMap, attachment.ObjectId.ValueString())
		}
	}

	if len(attachObjectList) > 0 {
		_, _, err := r.client.AttachBlockStorages(ctx, plan.Id.ValueString(), attachObjectList)
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error updating Block Storage(BM)",
				"Could not update Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
			return
		}
	}

	detachedObjectIdList := make([]string, 0)
	for k := range beforeAttachmentMap {
		detachedObjectIdList = append(detachedObjectIdList, k)
	}
	if len(detachedObjectIdList) > 0 {
		_, _, err := r.client.DetachBlockStorages(ctx, plan.Id.ValueString(), detachedObjectIdList)
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError("Error updating Block Storage(BM)",
				"Could not update Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
			return
		}
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "baremetal-blockstorage", "volume", plan.Id.ValueString(), tagElements)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating tags",
			err.Error(),
		)
		return
	} else {
		plan.Tags = tagsMap
	}

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *baremetalBlockStorageVolume) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state baremetalblockstorage.VolumeResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	blockStorageId := state.Id.ValueString()
	_, _, err := r.client.DeleteBlockStorage(ctx, blockStorageId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError("Error Deleting Block Storage(BM)",
			"Could not delete Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 30*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	timeoutErrCause, _ := json.Marshal(map[string]interface{}{
		"errors": map[string]interface{}{
			"detail": "Create Timeout.",
		},
	})

	context.WithTimeoutCause(ctx, deleteTimeout, scpsdk.GenericOpenAPIError{
		ResponseBody: timeoutErrCause,
		ErrorMessage: "Delete Timeout",
	})

	err = r.waitForBlockStorageState(ctx, blockStorageId, []string{common.DeletingState, common.AvailableState}, []string{common.DeletedState}, deleteTimeout)
	if err != nil {
		var detail string
		if reflect.PointerTo(reflect.TypeOf(scpsdk.GenericOpenAPIError{})) == reflect.TypeOf(err) {
			detail = client.GetDetailFromError(err)
		} else {
			detail = err.Error()
		}
		response.Diagnostics.AddError("Error Deleting Block Storage(BM)",
			"Could not delete Block Storage(BM), unexpected error:"+err.Error()+"\nReason: "+detail)
		return
	}
}

func (r *baremetalBlockStorageVolume) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = inst.Client.BaremetalBlockStorage
	r.clients = inst.Client
}

func (r *baremetalBlockStorageVolume) waitForBlockStorageState(ctx context.Context, blockStorageId string,
	pendingStates []string, targetStates []string, timeout time.Duration) error {
	return r.client.WaitForStatus(ctx, pendingStates, targetStates, timeout, func() (interface{}, string, error) {
		info, c, err := r.client.GetBlockStorage(ctx, blockStorageId)
		if err != nil {
			if c == 404 {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, strings.ToUpper(string(*info.Result.State)), nil
	})
}
