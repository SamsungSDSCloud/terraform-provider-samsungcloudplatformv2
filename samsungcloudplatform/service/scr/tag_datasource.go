package scr

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrTagDataSource{}
	_ datasource.DataSourceWithConfigure = &scrTagDataSource{}
)

func NewScrTagDataSource() datasource.DataSource {
	return &scrTagDataSource{}
}

type scrTagDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrTagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_tag"
}

func (d *scrTagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Tag.",
		MarkdownDescription: "Get details of a specific tag in a Container Registry repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Tags ID",
				MarkdownDescription: "The unique identifier of the tag.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"image_id": schema.StringAttribute{
				Description:         "Image ID",
				MarkdownDescription: "The ID of the image this tag belongs to.",
				Computed:            true,
			},
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the registry containing this tag.",
				Computed:            true,
			},
			"repository_id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The ID of the repository containing this tag.",
				Computed:            true,
			},
			"hash_digest": schema.StringAttribute{
				Description:         "Hash digest",
				MarkdownDescription: "The hash digest of the tag.",
				Computed:            true,
			},
			"manifest": schema.StringAttribute{
				Description:         "Manifest",
				MarkdownDescription: "The manifest content of the tag.",
				Computed:            true,
				Sensitive:           true,
			},
			"manifest_media_type": schema.StringAttribute{
				Description:         "Manifest media type",
				MarkdownDescription: "The media type of the manifest.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Tag state",
				MarkdownDescription: "The current state of the tag.",
				Computed:            true,
			},
			"reference_tags": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Reference tags",
				MarkdownDescription: "List of tags referencing this tag.",
				Computed:            true,
			},
			"lock_policy": schema.ListNestedAttribute{
				Description:         "Lock policy",
				MarkdownDescription: "The lock policy applied to this tag.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"locked": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description:         "Created at",
				MarkdownDescription: "The time the tag was created.",
				Computed:            true,
			},
			"modified_at": schema.StringAttribute{
				Description:         "Modified at",
				MarkdownDescription: "The time the tag was last modified.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrTagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TagDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := d.client.ShowTags(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tag",
			err.Error(),
		)
		return
	}

	tag := showResp.Tags
	refTags, refDiags := types.ListValueFrom(ctx, types.StringType, tag.ReferenceTags)
	resp.Diagnostics.Append(refDiags...)
	lockPolicy, lockDiags := flattenTagLockPolicy(ctx, tag.LockPolicy.Get())
	resp.Diagnostics.Append(lockDiags...)

	state = TagDataSource{
		Id:                types.StringValue(tag.Id),
		ImageId:           types.StringValue(tag.ImageId),
		RegistryId:        types.StringValue(tag.RegistryId),
		RepositoryId:      types.StringValue(tag.RepositoryId),
		HashDigest:        types.StringValue(tag.HashDigest),
		Manifest:          types.StringValue(tag.Manifest),
		ManifestMediaType: types.StringValue(tag.ManifestMediaType),
		State:             types.StringValue(tag.State),
		ReferenceTags:     refTags,
		LockPolicy:        lockPolicy,
		CreatedAt:         types.StringValue(tag.CreatedAt.Format(time.RFC3339)),
		ModifiedAt:        types.StringValue(tag.ModifiedAt.Format(time.RFC3339)),
		Filter:            state.Filter,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenTagLockPolicy(ctx context.Context, policy *scr11.LockPolicy) (types.List, diag.Diagnostics) {
	if policy == nil {
		return types.ListValueFrom(ctx, LockPolicyModelType, []LockPolicyModel{})
	}
	model := LockPolicyModel{Locked: types.BoolValue(policy.Locked)}
	return types.ListValueFrom(ctx, LockPolicyModelType, []LockPolicyModel{model})
}
