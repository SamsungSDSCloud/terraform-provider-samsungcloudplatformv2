package filestorage

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/filestorage"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &fileStorageAccessRuleResource{}
	_ resource.ResourceWithConfigure   = &fileStorageAccessRuleResource{}
	_ resource.ResourceWithImportState = &fileStorageAccessRuleResource{}
)

func NewFileStorageAccessRuleResource() resource.Resource {
	return &fileStorageAccessRuleResource{}
}

type fileStorageAccessRuleResource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (r *fileStorageAccessRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filestorage_access_rule"
}

func (r *fileStorageAccessRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a File Storage Access Rule. Use this resource to grant access to a File Storage Volume for a specific object (VM, BM, GPU, GPU_NODE, ENDPOINT). Each rule is managed independently, allowing mixed management with external automation (e.g., Kubernetes Auto Scaling).\n" +
				"Since all fields use RequiresReplace(), this resource does not support in-place updates. " +
				"The Read function only checks rule existence (not attribute drift) — if a rule is removed externally, " +
				"Terraform will detect it as gone and recreate it on the next apply.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the access rule.\n" +
					"  - format: {file_storage_id}:{object_type}:{object_id} \n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_storage_id": schema.StringAttribute{
				Description: "The identifier of the File Storage Volume that the access rule belongs to.\n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"object_id": schema.StringAttribute{
				Description: "The identifier of the object to grant access to. \n" +
					"  - example: 43fq3347-02q4-4aa8-ccf9-affe4917bb6f \n",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"object_type": schema.StringAttribute{
				Description: "The type of the object to grant access to.\n" +
					"  - example: VM\n" +
					"  - valid: VM, BM, GPU, GPU_NODE, ENDPOINT \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VM", "BM", "GPU", "GPU_NODE", "ENDPOINT"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *fileStorageAccessRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = inst.Client.FileStorage
	r.clients = inst.Client
}

func (r *fileStorageAccessRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan filestorage.FileStorageAccessRuleResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := filestorage.AccessRuleResource{
		ObjectId:   plan.ObjectId,
		ObjectType: plan.ObjectType,
	}

	err := r.client.UpdateVolumeAccessRule(ctx, plan.FileStorageId.ValueString(), rule, "add")
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating file storage access rule",
			"Could not create file storage access rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(filestorage.BuildAccessRuleResourceID(
		plan.FileStorageId.ValueString(),
		plan.ObjectType.ValueString(),
		plan.ObjectId.ValueString(),
	))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileStorageAccessRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state filestorage.FileStorageAccessRuleResource

    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    data, err := r.client.GetVolumeAccessRules(ctx, state.FileStorageId.ValueString())
    if err != nil {
        // volume 자체가 삭제된 경우: rule도 존재할 수 없으므로 state 제거
        if strings.Contains(err.Error(), "404") {
            resp.State.RemoveResource(ctx)
            return
        }
        detail := client.GetDetailFromError(err)
        resp.Diagnostics.AddError(
            "Error reading file storage access rule",
            "Could not read file storage access rules: "+err.Error()+"\nReason: "+detail,
        )
        return
    }

    found := false
    for _, rule := range data.AccessRules {
        if rule.ObjectId == state.ObjectId.ValueString() && rule.ObjectType == state.ObjectType.ValueString() {
            found = true
            break
        }
    }

    // rule이 외부에서 제거된 경우: 다음 apply에서 재생성되도록 state 제거
    if !found {
        resp.State.RemoveResource(ctx)
        return
    }

    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *fileStorageAccessRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op: all fields use RequiresReplace(), so Terraform will destroy and recreate
	// instead of calling this method. This satisfies the Resource interface.
}

func (r *fileStorageAccessRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state filestorage.FileStorageAccessRuleResource

    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    rule := filestorage.AccessRuleResource{
        ObjectId:   state.ObjectId,
        ObjectType: state.ObjectType,
    }

    err := r.client.UpdateVolumeAccessRule(ctx, state.FileStorageId.ValueString(), rule, "remove")
    if err != nil {
        if strings.Contains(err.Error(), "404") {
            return
        }
        detail := client.GetDetailFromError(err)
        resp.Diagnostics.AddError(
            "Error deleting file storage access rule",
            "Could not delete file storage access rule, unexpected error: "+err.Error()+"\nReason: "+detail,
        )
        return
    }
}

func (r *fileStorageAccessRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fileStorageId, objectType, objectId, err := filestorage.ParseAccessRuleResourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: {file_storage_id}:{object_type}:{object_id}, got: "+req.ID,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("file_storage_id"), fileStorageId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_type"), objectType)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_id"), objectId)...)
}
