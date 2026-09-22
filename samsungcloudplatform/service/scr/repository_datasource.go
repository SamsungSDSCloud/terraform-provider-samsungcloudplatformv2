package scr

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrRepositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &scrRepositoryDataSource{}
)

func NewScrRepositoryDataSource() datasource.DataSource {
	return &scrRepositoryDataSource{}
}

type scrRepositoryDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrRepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_repository"
}

func (d *scrRepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Repository.",
		MarkdownDescription: "Get details of a Container Registry Repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The unique identifier of the repository.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Repository name",
				MarkdownDescription: "The name of the repository.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Repository description",
				MarkdownDescription: "A description of the repository.",
				Computed:            true,
			},
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the parent container registry.",
				Computed:            true,
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
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"critical_limit": schema.Int32Attribute{Computed: true},
						"high_limit": schema.Int32Attribute{Computed: true},
						"unmodified_excepted": schema.BoolAttribute{Computed: true},
						"unscanned_image_pull_prevented": schema.BoolAttribute{Computed: true},
						"vulnerable_image_pull_prevented": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"scan_policy": schema.ListNestedAttribute{
				Description:         "Scan policy configuration",
				MarkdownDescription: "Controls automatic vulnerability scanning behavior.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auto_scan_enabled": schema.BoolAttribute{Computed: true},
						"fixed_version_excepted": schema.BoolAttribute{Computed: true},
						"language_excepted": schema.BoolAttribute{Computed: true},
						"scan_policy_enabled": schema.BoolAttribute{Computed: true},
						"secret_excepted": schema.BoolAttribute{Computed: true},
						"severity_limit": schema.StringAttribute{Computed: true},
					},
				},
			},
			"lifecycle_policy": schema.ListNestedAttribute{
				Description:         "Lifecycle policy configuration",
				MarkdownDescription: "Controls automatic cleanup of outdated or untagged images.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"lifecycle_policy_enabled": schema.BoolAttribute{Computed: true},
						"outdated_rule_duration": schema.Int32Attribute{Computed: true},
						"outdated_rule_enabled": schema.BoolAttribute{Computed: true},
						"outdated_rule_tag_expression": schema.StringAttribute{Computed: true},
						"untagged_rule_duration": schema.Int32Attribute{Computed: true},
						"untagged_rule_enabled": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"lock_policy": schema.ListNestedAttribute{
				Description:         "Lock policy configuration",
				MarkdownDescription: "Controls whether the repository is locked from deletion.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"locked": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"tags": tag.DataSourceSchema(),
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Scr
	d.clients = inst.Client
}

func (d *scrRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RepositoryDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := d.client.ShowRepository(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Repository",
			err.Error(),
		)
		return
	}

	repo := showResp.Repository
	registryId := types.StringNull()
	if repo.RegistryId != nil {
		registryId = types.StringValue(*repo.RegistryId)
	}

	originalDsTags := state.Tags
	state = RepositoryDataSource{
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
		Filter:             state.Filter,
	}

	// Read tags from resource manager
	tagsMap, err := tag.GetTags(d.clients, "scr", "repository", repo.Id, false)
	if err != nil {
		state.Tags = types.MapNull(types.StringType)
	} else {
		state.Tags = common.NullTagCheck(tagsMap, originalDsTags)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ---- Repositories Data Source (list) ----

var (
	_ datasource.DataSource              = &scrRepositoryDataSources{}
	_ datasource.DataSourceWithConfigure = &scrRepositoryDataSources{}
)

func NewScrRepositoryDataSources() datasource.DataSource {
	return &scrRepositoryDataSources{}
}

type scrRepositoryDataSources struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrRepositoryDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_repositories"
}

func (d *scrRepositoryDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of Container Registry Repositories.",
		MarkdownDescription: "Get a list of Repository IDs in a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the parent container registry.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Repository name filter",
				MarkdownDescription: "Filter repositories by name.",
				Optional:            true,
			},
			"sort": schema.StringAttribute{
				Description:         "Sort",
				MarkdownDescription: "Sort order, e.g. `name:asc`.",
				Optional:            true,
			},
			"page": schema.Int32Attribute{
				Description:         "Page number",
				MarkdownDescription: "Page number (0-based).",
				Optional:            true,
			},
			"size": schema.Int32Attribute{
				Description:         "Page size",
				MarkdownDescription: "Number of items per page.",
				Optional:            true,
			},
			"ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Repository ID list",
				MarkdownDescription: "List of repository IDs.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrRepositoryDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Scr
	d.clients = inst.Client
}

func (d *scrRepositoryDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RepositoryDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name = state.Name.ValueString()
	}
	sort := ""
	if !state.Sort.IsNull() && !state.Sort.IsUnknown() {
		sort = state.Sort.ValueString()
	}
	var page, size int32
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		page = state.Page.ValueInt32()
	}
	if !state.Size.IsNull() && !state.Size.IsUnknown() {
		size = state.Size.ValueInt32()
	}

	listResp, err := d.client.ListRepositories(ctx, state.RegistryId.ValueString(), name, sort, page, size)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Repositories",
			err.Error(),
		)
		return
	}

	contents := listResp.Repositories
	filteredContents := contents

	if len(state.Filter) > 0 {
		filteredContents = filteredContents[:0]
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError("Filter Error", err.Error())
			return
		}

		for i, repo := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, repo)
			}
		}
		contents = filteredContents
	}

	var ids []types.String
	for _, repo := range contents {
		ids = append(ids, types.StringValue(repo.Id))
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
