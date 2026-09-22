package scr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource              = &scrImageResource{}
	_ resource.ResourceWithConfigure = &scrImageResource{}
	_ resource.ResourceWithImportState = &scrImageResource{}
)

func NewScrImageResource() resource.Resource {
	return &scrImageResource{}
}

type scrImageResource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (r *scrImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_image"
}

func (r *scrImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scrImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *scrImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage a Container Registry Image.",
		MarkdownDescription: "Manages properties and policies of a Docker image in a Container Registry. Images are pushed externally; this resource manages description, policies, and lifecycle.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Image ID",
				MarkdownDescription: "The unique identifier of the image.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Image name",
				MarkdownDescription: "The name of the image (digest).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Image description",
				MarkdownDescription: "A description of the image.",
				Optional:            true,
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Image state",
				MarkdownDescription: "The current state of the image.",
				Computed:            true,
			},
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the registry containing the image.",
				Computed:            true,
			},
			"repository_id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The ID of the repository containing the image.",
				Computed:            true,
			},
			"pull_count": schema.Int32Attribute{
				Description:         "Pull count",
				MarkdownDescription: "The number of times the image has been pulled.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				Description:         "Created at",
				MarkdownDescription: "The time the image was created.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				Description:         "Created by",
				MarkdownDescription: "The user who created the image.",
				Computed:            true,
			},
			"modified_at": schema.StringAttribute{
				Description:         "Modified at",
				MarkdownDescription: "The time the image was last modified.",
				Computed:            true,
			},
			"modified_by": schema.StringAttribute{
				Description:         "Modified by",
				MarkdownDescription: "The user who last modified the image.",
				Computed:            true,
			},
			"pull_policy": schema.ListNestedAttribute{
				Description:         "Pull policy configuration",
				MarkdownDescription: "Controls pull behavior based on vulnerability scan results.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"critical_limit": schema.Int32Attribute{
							Description:         "Critical vulnerability limit",
							MarkdownDescription: "Maximum number of critical vulnerabilities allowed for pull.",
							Optional:            true,
						},
						"high_limit": schema.Int32Attribute{
							Description:         "High vulnerability limit",
							MarkdownDescription: "Maximum number of high vulnerabilities allowed for pull.",
							Optional:            true,
						},
						"unmodified_excepted": schema.BoolAttribute{
							Description:         "Unmodified exception",
							MarkdownDescription: "Whether to except unmodified images from pull policy.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"unscanned_image_pull_prevented": schema.BoolAttribute{
							Description:         "Unscanned image pull prevented",
							MarkdownDescription: "Whether to prevent pulling unscanned images.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"vulnerable_image_pull_prevented": schema.BoolAttribute{
							Description:         "Vulnerable image pull prevented",
							MarkdownDescription: "Whether to prevent pulling vulnerable images.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
			"scan_policy": schema.ListNestedAttribute{
				Description:         "Scan policy configuration",
				MarkdownDescription: "Controls automatic vulnerability scanning behavior.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auto_scan_enabled": schema.BoolAttribute{
							Description:         "Auto scan enabled",
							MarkdownDescription: "Whether to automatically scan images on push.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"fixed_version_excepted": schema.BoolAttribute{
							Description:         "Fixed version exception",
							MarkdownDescription: "Whether to except vulnerabilities with fixed versions.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"language_excepted": schema.BoolAttribute{
							Description:         "Language exception",
							MarkdownDescription: "Whether to except language-related vulnerabilities.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"scan_policy_enabled": schema.BoolAttribute{
							Description:         "Scan policy enabled",
							MarkdownDescription: "Whether the scan policy is enabled.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"secret_excepted": schema.BoolAttribute{
							Description:         "Secret exception",
							MarkdownDescription: "Whether to except secret-related vulnerabilities.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"severity_limit": schema.StringAttribute{
							Description:         "Severity limit",
							MarkdownDescription: "Minimum severity level for scan alerts.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"lifecycle_policy": schema.ListNestedAttribute{
				Description:         "Lifecycle policy configuration",
				MarkdownDescription: "Controls automatic cleanup of outdated or untagged images.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"lifecycle_policy_enabled": schema.BoolAttribute{
							Description:         "Lifecycle policy enabled",
							MarkdownDescription: "Whether the lifecycle policy is enabled.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"outdated_rule_duration": schema.Int32Attribute{
							Description:         "Outdated rule duration (days)",
							MarkdownDescription: "Number of days before an outdated image is cleaned up.",
							Optional:            true,
						},
						"outdated_rule_enabled": schema.BoolAttribute{
							Description:         "Outdated rule enabled",
							MarkdownDescription: "Whether the outdated rule is enabled.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"outdated_rule_tag_expression": schema.StringAttribute{
							Description:         "Outdated rule tag expression",
							MarkdownDescription: "Tag expression to match outdated images.",
							Optional:            true,
						},
						"untagged_rule_duration": schema.Int32Attribute{
							Description:         "Untagged rule duration (days)",
							MarkdownDescription: "Number of days before an untagged image is cleaned up.",
							Optional:            true,
						},
						"untagged_rule_enabled": schema.BoolAttribute{
							Description:         "Untagged rule enabled",
							MarkdownDescription: "Whether the untagged rule is enabled.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
			"lock_policy": schema.ListNestedAttribute{
				Description:         "Lock policy configuration",
				MarkdownDescription: "Controls whether the image is locked from deletion.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"locked": schema.BoolAttribute{
							Description:         "Locked",
							MarkdownDescription: "Whether the image is locked.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (r *scrImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageId := plan.Id.ValueString()

	// Image already exists — just apply policies
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		descReq := scr11.ImageSetRequest{Description: plan.Description.ValueString()}
		if err := r.client.UpdateImageDescription(ctx, imageId, descReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating image description",
				"Could not update image description: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	applyImagePolicies(ctx, r.client, imageId, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := r.client.ShowImage(ctx, imageId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading image after creation",
			"Could not read image ID "+imageId+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	plan = flattenImage(showResp.Image)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scrImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := r.client.ShowImage(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Image",
			"Could not read image ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	state = flattenImage(showResp.Image)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *scrImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageId := plan.Id.ValueString()

	// Update description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		descReq := scr11.ImageSetRequest{Description: plan.Description.ValueString()}
		if err := r.client.UpdateImageDescription(ctx, imageId, descReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating image description",
				"Could not update image description: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	applyImagePolicies(ctx, r.client, imageId, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := r.client.ShowImage(ctx, imageId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading image after update",
			"Could not read image ID "+imageId+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	plan = flattenImage(showResp.Image)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scrImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteImage(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Image",
			"Could not delete image, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}
}

func applyImagePolicies(ctx context.Context, scrClient *scr.Client, imageId string, plan ImageResourceModel, diags *diag.Diagnostics) {
	// Update pull policy
	if !plan.PullPolicy.IsNull() && !plan.PullPolicy.IsUnknown() {
		pullPolicy := buildPullPolicy(ctx, plan.PullPolicy)
		pullReq := scr11.PullPolicySetRequest{PullPolicy: pullPolicy}
		if err := scrClient.UpdateImagePullPolicy(ctx, imageId, pullReq); err != nil {
			detail := client.GetDetailFromError(err)
			diags.AddError(
				"Error updating image pull policy",
				"Could not update image pull policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update scan policy
	if !plan.ScanPolicy.IsNull() && !plan.ScanPolicy.IsUnknown() {
		scanPolicy := buildScanPolicy(ctx, plan.ScanPolicy)
		scanReq := scr11.ScanPolicySetRequest{ScanPolicy: scanPolicy}
		if err := scrClient.UpdateImageScanPolicy(ctx, imageId, scanReq); err != nil {
			detail := client.GetDetailFromError(err)
			diags.AddError(
				"Error updating image scan policy",
				"Could not update image scan policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update lifecycle policy
	if !plan.LifecyclePolicy.IsNull() && !plan.LifecyclePolicy.IsUnknown() {
		lifecyclePolicy := buildLifecyclePolicy(ctx, plan.LifecyclePolicy)
		lifecycleReq := scr11.LifecyclePolicySetRequestV11{LifecyclePolicy: lifecyclePolicy}
		if err := scrClient.UpdateImageLifecyclePolicy(ctx, imageId, lifecycleReq); err != nil {
			detail := client.GetDetailFromError(err)
			diags.AddError(
				"Error updating image lifecycle policy",
				"Could not update image lifecycle policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update lock policy
	if !plan.LockPolicy.IsNull() && !plan.LockPolicy.IsUnknown() {
		lockPolicy := buildLockPolicy(ctx, plan.LockPolicy)
		lockReq := scr11.LockPolicySetRequest{LockPolicy: lockPolicy}
		if err := scrClient.UpdateImageLockPolicy(ctx, imageId, lockReq); err != nil {
			detail := client.GetDetailFromError(err)
			diags.AddError(
				"Error updating image lock policy",
				"Could not update image lock policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}
}

func flattenImage(image scr11.ImageV11) ImageResourceModel {
	return ImageResourceModel{
		Id:                types.StringValue(image.Id),
		Name:              types.StringValue(image.Name),
		Description:       nullableStringValue(image.Description.Get()),
		State:             types.StringValue(image.State),
		RegistryId:        types.StringValue(image.RegistryId),
		RepositoryId:      types.StringValue(image.RepositoryId),
		PullCount:         types.Int32Value(image.PullCount),
		PullPolicy:        flattenPullPolicy(image.PullPolicy.Get()),
		ScanPolicy:        flattenScanPolicy(image.ScanPolicy.Get()),
		LifecyclePolicy:   flattenLifecyclePolicy(image.LifecyclePolicy.Get()),
		LockPolicy:        flattenLockPolicy(image.LockPolicy.Get()),
		CreatedAt:         types.StringValue(image.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(image.CreatedBy),
		ModifiedAt:        types.StringValue(image.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(image.ModifiedBy),
	}
}
