package scr

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
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
	_ datasource.DataSource              = &scrTagSecretsDataSource{}
	_ datasource.DataSourceWithConfigure = &scrTagSecretsDataSource{}
)

func NewScrTagSecretsDataSource() datasource.DataSource {
	return &scrTagSecretsDataSource{}
}

type scrTagSecretsDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrTagSecretsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_tag_secrets"
}

func (d *scrTagSecretsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Tag Secrets.",
		MarkdownDescription: "Get a list of secrets for a specific tag in a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"tags_id": schema.StringAttribute{
				Description:         "Tags ID",
				MarkdownDescription: "The ID of the tag to get secrets for.",
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
			"file_name": schema.StringAttribute{
				Description:         "File name filter",
				MarkdownDescription: "Filter secrets by file name.",
				Optional:            true,
			},
			"filtered_count": schema.Int32Attribute{
				Description:         "Filtered count",
				MarkdownDescription: "Number of secrets matching the filter.",
				Computed:            true,
			},
			"last_scanned_at": schema.StringAttribute{
				Description:         "Last scanned at",
				MarkdownDescription: "The time the tag was last scanned.",
				Computed:            true,
			},
			"release_version": schema.StringAttribute{
				Description:         "Release version",
				MarkdownDescription: "The OS release version of the tag.",
				Computed:            true,
			},
			"secrets": schema.ListNestedAttribute{
				Description:         "Secret list",
				MarkdownDescription: "List of secrets found in the tag.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category":            schema.StringAttribute{Computed: true},
						"file_name":           schema.StringAttribute{Computed: true},
						"match":               schema.StringAttribute{Computed: true, Sensitive: true},
						"rule_id":             schema.StringAttribute{Computed: true},
						"severity":            schema.StringAttribute{Computed: true},
						"start_line":          schema.Int32Attribute{Computed: true},
						"target":              schema.StringAttribute{Computed: true, Sensitive: true},
						"title":               schema.StringAttribute{Computed: true},
						"vulnerability_class": schema.StringAttribute{Computed: true},
					},
				},
			},
			"secret_summary": schema.SingleNestedAttribute{
				Description:         "Secret summary",
				MarkdownDescription: "Aggregate secret severity summary.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"critical":     schema.Int32Attribute{Computed: true},
					"high":         schema.Int32Attribute{Computed: true},
					"low":          schema.Int32Attribute{Computed: true},
					"medium":       schema.Int32Attribute{Computed: true},
					"total_secret": schema.Int32Attribute{Computed: true},
					"unknown":      schema.Int32Attribute{Computed: true},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrTagSecretsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrTagSecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TagSecretsDataSource

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
	fileName := ""
	if !state.FileName.IsNull() && !state.FileName.IsUnknown() {
		fileName = state.FileName.ValueString()
	}

	opts := scr.TagSecretsQueryOptions{
		Sort:     sort,
		Page:     page,
		Size:     size,
		FileName: fileName,
	}

	listResp, err := d.client.ShowTagsSecrets(ctx, state.TagsId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tag Secrets",
			err.Error(),
		)
		return
	}

	state.FilteredCount = types.Int32Value(listResp.FilteredCount)
	state.LastScannedAt = types.StringValue(listResp.LastScannedAt)
	state.ReleaseVersion = types.StringValue(listResp.ReleaseVersion)

	secrets := make([]SecretReportModel, 0, len(listResp.SecretReports))
	for _, secret := range listResp.SecretReports {
		flattened, fDiags := flattenSecretReport(ctx, secret)
		resp.Diagnostics.Append(fDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		secrets = append(secrets, flattened)
	}
	var secDiags diag.Diagnostics
	state.Secrets, secDiags = types.ListValueFrom(ctx, SecretReportModelType, secrets)
	resp.Diagnostics.Append(secDiags...)

	var summaryDiags diag.Diagnostics
	state.SecretSummary, summaryDiags = flattenSecretSummary(ctx, listResp.SecretSummary)
	resp.Diagnostics.Append(summaryDiags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenSecretReport(ctx context.Context, secret scr11.SecretReport) (SecretReportModel, diag.Diagnostics) {
	return SecretReportModel{
		Category:           types.StringValue(secret.Category),
		FileName:           types.StringValue(secret.FileName),
		Match:              types.StringValue(secret.Match),
		RuleId:             types.StringValue(secret.RuleId),
		Severity:           types.StringValue(secret.Severity),
		StartLine:          types.Int32Value(secret.StartLine),
		Target:             types.StringValue(secret.Target),
		Title:              types.StringValue(secret.Title),
		VulnerabilityClass: types.StringValue(secret.VulnerabilityClass),
	}, nil
}

func flattenSecretSummary(ctx context.Context, summary scr11.SecretSummary) (types.Object, diag.Diagnostics) {
	model := SecretSummaryModel{
		Critical:    types.Int32Value(summary.Critical),
		High:        types.Int32Value(summary.High),
		Low:         types.Int32Value(summary.Low),
		Medium:      types.Int32Value(summary.Medium),
		TotalSecret: types.Int32Value(summary.TotalSecret),
		Unknown:     types.Int32Value(summary.Unknown),
	}
	return types.ObjectValueFrom(ctx, SecretSummaryModelType.AttrTypes, model)
}
