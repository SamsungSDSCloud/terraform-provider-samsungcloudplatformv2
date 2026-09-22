package scr

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrTagsDataSource{}
	_ datasource.DataSourceWithConfigure = &scrTagsDataSource{}
)

func NewScrTagsDataSource() datasource.DataSource {
	return &scrTagsDataSource{}
}

type scrTagsDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrTagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_tags"
}

func (d *scrTagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of Container Registry Tags.",
		MarkdownDescription: "Get a list of tags for an image in a Container Registry repository.",
		Attributes: map[string]schema.Attribute{
			"image_id": schema.StringAttribute{
				Description:         "Image ID",
				MarkdownDescription: "The ID of the image to list tags from.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
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
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"size": schema.Int32Attribute{
				Description:         "Page size",
				MarkdownDescription: "Number of items per page.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"reference_tags": schema.StringAttribute{
				Description:         "Reference tags filter",
				MarkdownDescription: "Filter by reference tag name.",
				Optional:            true,
			},
			"total_count": schema.Int32Attribute{
				Description:         "Total count",
				MarkdownDescription: "Total number of tags matching the query.",
				Computed:            true,
			},
			"tags": schema.ListNestedAttribute{
				Description:         "Tag list",
				MarkdownDescription: "List of tags for the image.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"hash_digest":          schema.StringAttribute{Computed: true},
						"reference_tags":       schema.ListAttribute{Computed: true, ElementType: types.StringType},
						"lock_policy":          schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"locked": schema.BoolAttribute{Computed: true}}}},
						"last_scanned_at":      schema.StringAttribute{Computed: true},
						"modified_at":          schema.StringAttribute{Computed: true},
						"private_endpoint_url": schema.StringAttribute{Computed: true},
						"public_endpoint_url":  schema.StringAttribute{Computed: true},
						"re_scan_needed":       schema.BoolAttribute{Computed: true},
						"referenced_by":        schema.StringAttribute{Computed: true},
						"scan_state":           schema.StringAttribute{Computed: true},
						"scan_summary":         schema.SingleNestedAttribute{Computed: true, Attributes: scanSummaryDataSourceSchemaAttrs()},
						"size":                 schema.Int32Attribute{Computed: true},
						"state":                schema.StringAttribute{Computed: true},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrTagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrTagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TagsDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sort := ""
	if !state.Sort.IsNull() && !state.Sort.IsUnknown() {
		sort = state.Sort.ValueString()
	}
	var page, size *int32
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		v := state.Page.ValueInt32()
		page = &v
	}
	if !state.Size.IsNull() && !state.Size.IsUnknown() {
		v := state.Size.ValueInt32()
		size = &v
	}
	referenceTags := ""
	if !state.ReferenceTags.IsNull() && !state.ReferenceTags.IsUnknown() {
		referenceTags = state.ReferenceTags.ValueString()
	}

	listResp, err := d.client.ListTagses(ctx, state.ImageId.ValueString(), sort, page, size, referenceTags)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List Tags",
			err.Error(),
		)
		return
	}

	contents := listResp.Tagses
	filteredContents := contents

	if len(state.Filter) > 0 {
		filteredContents = filteredContents[:0]
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError("Filter Error", err.Error())
			return
		}

		for i, tag := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, tag)
			}
		}
		contents = filteredContents
	}

	state.TotalCount = types.Int32Value(listResp.Count)

	tags := make([]TagListModel, 0, len(contents))
	for _, tag := range contents {
		flattened, fDiags := flattenTagList(ctx, tag)
		resp.Diagnostics.Append(fDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tags = append(tags, flattened)
	}
	var tagsDiags diag.Diagnostics
	state.Tags, tagsDiags = types.ListValueFrom(ctx, TagListModelType, tags)
	resp.Diagnostics.Append(tagsDiags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenTagList(ctx context.Context, tag scr11.TagsForList) (TagListModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	refTags, refDiags := types.ListValueFrom(ctx, types.StringType, tag.ReferenceTags)
	diags.Append(refDiags...)
	lockPolicy, lockDiags := flattenTagLockPolicy(ctx, tag.LockPolicy.Get())
	diags.Append(lockDiags...)
	scanSummary, summaryDiags := flattenScanSummary(ctx, tag.ScanSummary)
	diags.Append(summaryDiags...)
	return TagListModel{
		Id:                 types.StringValue(tag.Id),
		HashDigest:         types.StringValue(tag.HashDigest),
		ReferenceTags:      refTags,
		LockPolicy:         lockPolicy,
		LastScannedAt:      types.StringValue(tag.LastScannedAt.Format(time.RFC3339)),
		ModifiedAt:         types.StringValue(tag.ModifiedAt.Format(time.RFC3339)),
		PrivateEndpointUrl: types.StringValue(tag.PrivateEndpointUrl),
		PublicEndpointUrl:  types.StringValue(tag.PublicEndpointUrl),
		ReScanNeeded:       types.BoolValue(tag.ReScanNeeded),
		ReferencedBy:       types.StringValue(tag.ReferencedBy),
		ScanState:          types.StringValue(tag.ScanState),
		ScanSummary:        scanSummary,
		Size:               types.Int32Value(tag.Size),
		State:              types.StringValue(tag.State),
	}, diags
}
