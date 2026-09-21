package scr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
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
	_ resource.Resource              = &scrRepositoryResource{}
	_ resource.ResourceWithConfigure = &scrRepositoryResource{}
	_ resource.ResourceWithImportState = &scrRepositoryResource{}
)

func NewScrRepositoryResource() resource.Resource {
	return &scrRepositoryResource{}
}

type scrRepositoryResource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (r *scrRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_repository"
}

func (r *scrRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scrRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *scrRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Create and manage a Container Registry Repository.",
		MarkdownDescription: "Manages a Repository resource within a Container Registry for organizing Docker images.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The unique identifier of the repository.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Repository name",
				MarkdownDescription: "The name of the repository. Must start with a lowercase letter or number and contain only lowercase letters, numbers, and hyphens.\n\nExample: `my-repository`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Repository description",
				MarkdownDescription: "A description of the repository.",
				Optional:            true,
				Computed:            true,
			},
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the parent container registry.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "Repository state",
				MarkdownDescription: "The current state of the repository.",
				Computed:            true,
			},
			"private_endpoint_url": schema.StringAttribute{
				Description:         "Private endpoint URL",
				MarkdownDescription: "The private endpoint URL for the repository.",
				Computed:            true,
			},
			"public_endpoint_url": schema.StringAttribute{
				Description:         "Public endpoint URL",
				MarkdownDescription: "The public endpoint URL for the repository.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				Description:         "Created at",
				MarkdownDescription: "The time the repository was created.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				Description:         "Created by",
				MarkdownDescription: "The user who created the repository.",
				Computed:            true,
			},
			"modified_at": schema.StringAttribute{
				Description:         "Modified at",
				MarkdownDescription: "The time the repository was last modified.",
				Computed:            true,
			},
			"modified_by": schema.StringAttribute{
				Description:         "Modified by",
				MarkdownDescription: "The user who last modified the repository.",
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
				MarkdownDescription: "Controls whether the repository is locked from deletion.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"locked": schema.BoolAttribute{
							Description:         "Locked",
							MarkdownDescription: "Whether the repository is locked.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

func (r *scrRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RepositoryResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := scr11.RepositoryCreateRequestV11{
		Name:       plan.Name.ValueString(),
		RegistryId: plan.RegistryId.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description = plan.Description.ValueString()
	}

	// Build pull policy
	if !plan.PullPolicy.IsNull() && !plan.PullPolicy.IsUnknown() {
		pullPolicy := buildPullPolicy(ctx, plan.PullPolicy)
		createReq.PullPolicy = *scr11.NewNullablePullPolicy(&pullPolicy)
	}

	// Build scan policy
	if !plan.ScanPolicy.IsNull() && !plan.ScanPolicy.IsUnknown() {
		scanPolicy := buildScanPolicy(ctx, plan.ScanPolicy)
		createReq.ScanPolicy = *scr11.NewNullableScanPolicy(&scanPolicy)
	}

	// Build lifecycle policy
	if !plan.LifecyclePolicy.IsNull() && !plan.LifecyclePolicy.IsUnknown() {
		lifecyclePolicy := buildLifecyclePolicy(ctx, plan.LifecyclePolicy)
		createReq.LifecyclePolicy = *scr11.NewNullableLifecyclePolicyV11(&lifecyclePolicy)
	}

	// Build lock policy
	if !plan.LockPolicy.IsNull() && !plan.LockPolicy.IsUnknown() {
		lockPolicy := buildLockPolicy(ctx, plan.LockPolicy)
		createReq.LockPolicy = *scr11.NewNullableLockPolicy(&lockPolicy)
	}

	createResp, err := r.client.CreateRepository(ctx, createReq)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating repository",
			"Could not create repository, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	// Update tags if provided
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && len(plan.Tags.Elements()) > 0 {
		_, err := tag.UpdateTags(r.clients, "scr", "repository", createResp.Id, plan.Tags.Elements(), false)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error setting repository tags",
				"Could not set tags on repository "+createResp.Id+": "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Read back the full state
	showResp, err := r.client.ShowRepository(ctx, createResp.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading repository after creation",
			"Could not read repository ID "+createResp.Id+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	originalTags := plan.Tags
	plan = flattenRepository(showResp.Repository)
	// Read tags from resource manager
	tagsMap, err := tag.GetTags(r.clients, "scr", "repository", createResp.Id, false)
	if err != nil {
		plan.Tags = types.MapNull(types.StringType)
	} else {
		plan.Tags = common.NullTagCheck(tagsMap, originalTags)
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scrRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RepositoryResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := r.client.ShowRepository(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Repository",
			"Could not read repository ID "+state.Id.ValueString()+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	originalStateTags := state.Tags
	state = flattenRepository(showResp.Repository)
	// Read tags from resource manager
	tagsMap, err := tag.GetTags(r.clients, "scr", "repository", state.Id.ValueString(), false)
	if err != nil {
		state.Tags = types.MapNull(types.StringType)
	} else {
		state.Tags = common.NullTagCheck(tagsMap, originalStateTags)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *scrRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RepositoryResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state RepositoryResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryId := plan.Id.ValueString()

	// Update description (only if changed)
	if !plan.Description.Equal(state.Description) {
		descReq := scr11.RepositorySetRequest{
			Description: plan.Description.ValueString(),
		}
		if err := r.client.UpdateRepositoryDescription(ctx, repositoryId, descReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating repository description",
				"Could not update repository description: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update pull policy (only if changed)
	if !plan.PullPolicy.Equal(state.PullPolicy) {
		pullPolicy := buildPullPolicy(ctx, plan.PullPolicy)
		pullReq := scr11.PullPolicySetRequest{PullPolicy: pullPolicy}
		if err := r.client.UpdateRepositoryPullPolicy(ctx, repositoryId, pullReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating pull policy",
				"Could not update pull policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update scan policy (only if changed)
	if !plan.ScanPolicy.Equal(state.ScanPolicy) {
		scanPolicy := buildScanPolicy(ctx, plan.ScanPolicy)
		scanReq := scr11.ScanPolicySetRequest{ScanPolicy: scanPolicy}
		if err := r.client.UpdateRepositoryScanPolicy(ctx, repositoryId, scanReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating scan policy",
				"Could not update scan policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update lifecycle policy (only if changed)
	if !plan.LifecyclePolicy.Equal(state.LifecyclePolicy) {
		lifecyclePolicy := buildLifecyclePolicy(ctx, plan.LifecyclePolicy)
		lifecycleReq := scr11.LifecyclePolicySetRequestV11{LifecyclePolicy: lifecyclePolicy}
		if err := r.client.UpdateRepositoryLifecyclePolicy(ctx, repositoryId, lifecycleReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating lifecycle policy",
				"Could not update lifecycle policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Update lock policy (only if changed)
	if !plan.LockPolicy.Equal(state.LockPolicy) {
		lockPolicy := buildLockPolicy(ctx, plan.LockPolicy)
		lockReq := scr11.LockPolicySetRequest{LockPolicy: lockPolicy}
		if err := r.client.UpdateRepositoryLockPolicy(ctx, repositoryId, lockReq); err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating lock policy",
				"Could not update lock policy: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Sync tags if changed
	if !plan.Tags.Equal(state.Tags) {
		_, err := tag.UpdateTags(r.clients, "scr", "repository", repositoryId, plan.Tags.Elements(), false)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating repository tags",
				"Could not update tags: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	// Read back the full state
	showResp, err := r.client.ShowRepository(ctx, repositoryId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading repository after update",
			"Could not read repository ID "+repositoryId+": "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	originalUpdateTags := plan.Tags
	plan = flattenRepository(showResp.Repository)
	// Read tags from resource manager
	tagsMap, err := tag.GetTags(r.clients, "scr", "repository", repositoryId, false)
	if err != nil {
		plan.Tags = types.MapNull(types.StringType)
	} else {
		plan.Tags = common.NullTagCheck(tagsMap, originalUpdateTags)
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scrRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RepositoryResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRepository(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Repository",
			"Could not delete repository, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}
}

func flattenRepository(repo scr11.RepositoryV11) RepositoryResource {
	registryId := types.StringNull()
	if repo.RegistryId != nil {
		registryId = types.StringValue(*repo.RegistryId)
	}

	return RepositoryResource{
		Id:                 types.StringValue(repo.Id),
		Name:               types.StringValue(repo.Name),
		Description:        nullableStringValue(repo.Description.Get()),
		RegistryId:         registryId,
		State:              types.StringValue(repo.State),
		PrivateEndpointUrl: nullableStringValue(repo.PrivateEndpointUrl.Get()),
		PublicEndpointUrl:  nullableStringValue(repo.PublicEndpointUrl.Get()),
		PullPolicy:         flattenPullPolicy(repo.PullPolicy.Get()),
		ScanPolicy:         flattenScanPolicy(repo.ScanPolicy.Get()),
		LifecyclePolicy:    flattenLifecyclePolicy(repo.LifecyclePolicy.Get()),
		LockPolicy:         flattenLockPolicy(repo.LockPolicy.Get()),
		CreatedAt:          types.StringValue(repo.CreatedAt.Format(time.RFC3339)),
		CreatedBy:          types.StringValue(repo.CreatedBy),
		ModifiedAt:         types.StringValue(repo.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:         types.StringValue(repo.ModifiedBy),
	}
}

func buildPullPolicy(ctx context.Context, policyList types.List) scr11.PullPolicy {
	var policyListArr []PullPolicyModel
	policyList.ElementsAs(ctx, &policyListArr, false)
	if len(policyListArr) > 0 {
		p := policyListArr[0]
		pp := scr11.PullPolicy{
			UnmodifiedExcepted:            p.UnmodifiedExcepted.ValueBool(),
			UnscannedImagePullPrevented:   p.UnscannedImagePullPrevented.ValueBool(),
			VulnerableImagePullPrevented:  p.VulnerableImagePullPrevented.ValueBool(),
		}
		if !p.CriticalLimit.IsNull() && !p.CriticalLimit.IsUnknown() {
			v := p.CriticalLimit.ValueInt32()
			pp.CriticalLimit = *scr11.NewNullableInt32(&v)
		}
		if !p.HighLimit.IsNull() && !p.HighLimit.IsUnknown() {
			v := p.HighLimit.ValueInt32()
			pp.HighLimit = *scr11.NewNullableInt32(&v)
		}
		return pp
	}
	return scr11.PullPolicy{
		UnmodifiedExcepted:            false,
		UnscannedImagePullPrevented:   false,
		VulnerableImagePullPrevented:  false,
	}
}

func flattenPullPolicy(policy *scr11.PullPolicy) types.List {
	if policy == nil {
		emptyList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: PullPolicyModel{}.AttributeTypes()}, []PullPolicyModel{})
		return emptyList
	}
	model := PullPolicyModel{
		CriticalLimit:               common.ToNullableInt32Value(policy.CriticalLimit.Get()),
		HighLimit:                   common.ToNullableInt32Value(policy.HighLimit.Get()),
		UnmodifiedExcepted:          types.BoolValue(policy.UnmodifiedExcepted),
		UnscannedImagePullPrevented: types.BoolValue(policy.UnscannedImagePullPrevented),
		VulnerableImagePullPrevented: types.BoolValue(policy.VulnerableImagePullPrevented),
	}
	result, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: PullPolicyModel{}.AttributeTypes()}, []PullPolicyModel{model})
	return result
}

func buildScanPolicy(ctx context.Context, policyList types.List) scr11.ScanPolicy {
	var policyListArr []ScanPolicyModel
	policyList.ElementsAs(ctx, &policyListArr, false)
	if len(policyListArr) > 0 {
		p := policyListArr[0]
		severityLimit := p.SeverityLimit.ValueString()
		if p.SeverityLimit.IsNull() || p.SeverityLimit.IsUnknown() || severityLimit == "" {
			severityLimit = "None"
		}
		return scr11.ScanPolicy{
			AutoScanEnabled:      p.AutoScanEnabled.ValueBool(),
			FixedVersionExcepted: p.FixedVersionExcepted.ValueBool(),
			LanguageExcepted:     p.LanguageExcepted.ValueBool(),
			ScanPolicyEnabled:    p.ScanPolicyEnabled.ValueBool(),
			SecretExcepted:       p.SecretExcepted.ValueBool(),
			SeverityLimit:        severityLimit,
		}
	}
	return scr11.ScanPolicy{
		SeverityLimit: "None",
	}
}

func flattenScanPolicy(policy *scr11.ScanPolicy) types.List {
	if policy == nil {
		emptyList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: ScanPolicyModel{}.AttributeTypes()}, []ScanPolicyModel{})
		return emptyList
	}
	model := ScanPolicyModel{
		AutoScanEnabled:      types.BoolValue(policy.AutoScanEnabled),
		FixedVersionExcepted: types.BoolValue(policy.FixedVersionExcepted),
		LanguageExcepted:     types.BoolValue(policy.LanguageExcepted),
		ScanPolicyEnabled:    types.BoolValue(policy.ScanPolicyEnabled),
		SecretExcepted:       types.BoolValue(policy.SecretExcepted),
		SeverityLimit:        types.StringValue(policy.SeverityLimit),
	}
	result, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: ScanPolicyModel{}.AttributeTypes()}, []ScanPolicyModel{model})
	return result
}

func buildLifecyclePolicy(ctx context.Context, policyList types.List) scr11.LifecyclePolicyV11 {
	var policyListArr []LifecyclePolicy
	policyList.ElementsAs(ctx, &policyListArr, false)
	if len(policyListArr) > 0 {
		p := policyListArr[0]
		return scr11.LifecyclePolicyV11{
			LifecyclePolicyEnabled:    p.LifecyclePolicyEnabled.ValueBool(),
			OutdatedRuleDuration:      p.OutdatedRuleDuration.ValueInt32(),
			OutdatedRuleEnabled:       p.OutdatedRuleEnabled.ValueBool(),
			OutdatedRuleTagExpression: p.OutdatedRuleTagExpression.ValueString(),
			UntaggedRuleDuration:      p.UntaggedRuleDuration.ValueInt32(),
			UntaggedRuleEnabled:       p.UntaggedRuleEnabled.ValueBool(),
		}
	}
	return scr11.LifecyclePolicyV11{}
}

func flattenLifecyclePolicy(policy *scr11.LifecyclePolicyV11) types.List {
	if policy == nil {
		emptyList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: LifecyclePolicy{}.AttributeTypes()}, []LifecyclePolicy{})
		return emptyList
	}
	model := LifecyclePolicy{
		LifecyclePolicyEnabled:    types.BoolValue(policy.LifecyclePolicyEnabled),
		OutdatedRuleDuration:      types.Int32Value(policy.OutdatedRuleDuration),
		OutdatedRuleEnabled:       types.BoolValue(policy.OutdatedRuleEnabled),
		OutdatedRuleTagExpression: types.StringValue(policy.OutdatedRuleTagExpression),
		UntaggedRuleDuration:      types.Int32Value(policy.UntaggedRuleDuration),
		UntaggedRuleEnabled:       types.BoolValue(policy.UntaggedRuleEnabled),
	}
	result, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: LifecyclePolicy{}.AttributeTypes()}, []LifecyclePolicy{model})
	return result
}

func buildLockPolicy(ctx context.Context, policyList types.List) scr11.LockPolicy {
	var policyListArr []LockPolicyModel
	policyList.ElementsAs(ctx, &policyListArr, false)
	if len(policyListArr) > 0 {
		p := policyListArr[0]
		return scr11.LockPolicy{Locked: p.Locked.ValueBool()}
	}
	return scr11.LockPolicy{Locked: false}
}

func flattenLockPolicy(policy *scr11.LockPolicy) types.List {
	if policy == nil {
		emptyList, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: LockPolicyModel{}.AttributeTypes()}, []LockPolicyModel{})
		return emptyList
	}
	model := LockPolicyModel{Locked: types.BoolValue(policy.Locked)}
	result, _ := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: LockPolicyModel{}.AttributeTypes()}, []LockPolicyModel{model})
	return result
}
