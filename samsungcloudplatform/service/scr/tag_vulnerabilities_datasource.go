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
	_ datasource.DataSource              = &scrTagVulnerabilitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &scrTagVulnerabilitiesDataSource{}
)

func NewScrTagVulnerabilitiesDataSource() datasource.DataSource {
	return &scrTagVulnerabilitiesDataSource{}
}

type scrTagVulnerabilitiesDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrTagVulnerabilitiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_tag_vulnerabilities"
}

func (d *scrTagVulnerabilitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Tag Vulnerabilities.",
		MarkdownDescription: "Get a list of vulnerabilities (CVEs) for a specific tag in a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"tags_id": schema.StringAttribute{
				Description:         "Tags ID",
				MarkdownDescription: "The ID of the tag to get vulnerabilities for.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"os_language": schema.StringAttribute{
				Description:         "OS or language filter",
				MarkdownDescription: "Filter vulnerabilities by OS or language, e.g. `Java`.",
				Optional:            true,
			},
			"package_name": schema.StringAttribute{
				Description:         "Package name filter",
				MarkdownDescription: "Filter vulnerabilities by package name.",
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
			"update_version_available": schema.BoolAttribute{
				Description:         "Update version available filter",
				MarkdownDescription: "Filter vulnerabilities where an updated version is available.",
				Optional:            true,
			},
			"severity": schema.StringAttribute{
				Description:         "Severity filter",
				MarkdownDescription: "Filter by severity, e.g. `High`, `Critical`.",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				Description:         "Category filter",
				MarkdownDescription: "Filter by category, e.g. `Language`, `OS`.",
				Optional:            true,
			},
			"filtered_count": schema.Int32Attribute{
				Description:         "Filtered count",
				MarkdownDescription: "Number of vulnerabilities matching the filter.",
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
			"vulnerabilities": schema.ListNestedAttribute{
				Description:         "Vulnerability list",
				MarkdownDescription: "List of CVEs found in the tag.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category":        schema.StringAttribute{Computed: true},
						"current_version": schema.StringAttribute{Computed: true},
						"cve_code":        schema.StringAttribute{Computed: true},
						"description":     schema.StringAttribute{Computed: true},
						"links":           schema.StringAttribute{Computed: true},
						"os_language":     schema.StringAttribute{Computed: true},
						"package_name":    schema.StringAttribute{Computed: true},
						"severity":        schema.StringAttribute{Computed: true},
						"update_version":  schema.StringAttribute{Computed: true},
						"vectors": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"attack_complexity":   schema.StringAttribute{Computed: true},
								"attack_vector":       schema.StringAttribute{Computed: true},
								"availability":        schema.StringAttribute{Computed: true},
								"base_severity":       schema.StringAttribute{Computed: true},
								"confidentiality":     schema.StringAttribute{Computed: true},
								"cvss":                schema.Float64Attribute{Computed: true},
								"integrity":           schema.StringAttribute{Computed: true},
								"privileges_required": schema.StringAttribute{Computed: true},
								"scope":               schema.StringAttribute{Computed: true},
								"user_interaction":    schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"scan_summary": schema.SingleNestedAttribute{
				Description:         "Scan summary",
				MarkdownDescription: "Aggregate vulnerability scan summary.",
				Computed:            true,
				Attributes:          scanSummaryDataSourceSchemaAttrs(),
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrTagVulnerabilitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrTagVulnerabilitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TagVulnerabilitiesDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	osLanguage := ""
	if !state.OsLanguage.IsNull() && !state.OsLanguage.IsUnknown() {
		osLanguage = state.OsLanguage.ValueString()
	}
	packageName := ""
	if !state.PackageName.IsNull() && !state.PackageName.IsUnknown() {
		packageName = state.PackageName.ValueString()
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
	var updateVersionAvailable *bool
	if !state.UpdateVersionAvailable.IsNull() && !state.UpdateVersionAvailable.IsUnknown() {
		v := state.UpdateVersionAvailable.ValueBool()
		updateVersionAvailable = &v
	}
	severity := ""
	if !state.Severity.IsNull() && !state.Severity.IsUnknown() {
		severity = state.Severity.ValueString()
	}
	category := ""
	if !state.Category.IsNull() && !state.Category.IsUnknown() {
		category = state.Category.ValueString()
	}

	opts := scr.TagVulnerabilityQueryOptions{
		OsLanguage:             osLanguage,
		PackageName:            packageName,
		Sort:                   sort,
		Page:                   page,
		Size:                   size,
		UpdateVersionAvailable: updateVersionAvailable,
		Severity:               severity,
		Category:               category,
	}

	listResp, err := d.client.ShowTagsVulnerabilities(ctx, state.TagsId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tag Vulnerabilities",
			err.Error(),
		)
		return
	}

	state.FilteredCount = types.Int32Value(listResp.FilteredCount)
	report := listResp.VulnerabilityReport
	state.LastScannedAt = types.StringValue(report.LastScannedAt)
	state.ReleaseVersion = types.StringValue(report.ReleaseVersion)

	cves := make([]CveModel, 0, len(report.Cves))
	for _, cve := range report.Cves {
		flattened, fDiags := flattenCve(ctx, cve)
		resp.Diagnostics.Append(fDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		cves = append(cves, flattened)
	}
	var vulnDiags diag.Diagnostics
	state.Vulnerabilities, vulnDiags = types.ListValueFrom(ctx, CveModelType, cves)
	resp.Diagnostics.Append(vulnDiags...)

	var summaryDiags diag.Diagnostics
	state.ScanSummary, summaryDiags = flattenScanSummary(ctx, report.ScanSummary)
	resp.Diagnostics.Append(summaryDiags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenCve(ctx context.Context, cve scr11.Cve) (CveModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	osLang := types.StringNull()
	if cve.OsLanguage.IsSet() {
		val := cve.OsLanguage.Get()
		if val != nil {
			osLang = types.StringValue(*val)
		}
	}
	pkgName := types.StringNull()
	if cve.PackageName.IsSet() {
		val := cve.PackageName.Get()
		if val != nil {
			pkgName = types.StringValue(*val)
		}
	}
	vectors, vecDiags := flattenVectors(ctx, cve.Vectors)
	diags.Append(vecDiags...)
	return CveModel{
		Category:       types.StringValue(cve.Category),
		CurrentVersion: types.StringValue(cve.CurrentVersion),
		CveCode:        types.StringValue(cve.CveCode),
		Description:    types.StringValue(cve.Description),
		Links:          types.StringValue(cve.Links),
		OsLanguage:     osLang,
		PackageName:    pkgName,
		Severity:       types.StringValue(cve.Severity),
		UpdateVersion:  types.StringValue(cve.UpdateVersion),
		Vectors:        vectors,
	}, diags
}

func flattenVectors(ctx context.Context, vectors scr11.Vectors) (types.Object, diag.Diagnostics) {
	model := VectorsModel{
		AttackComplexity:   types.StringValue(vectors.AttackComplexity),
		AttackVector:       types.StringValue(vectors.AttackVector),
		Availability:       types.StringValue(vectors.Availability),
		BaseSeverity:       types.StringValue(vectors.BaseSeverity),
		Confidentiality:    types.StringValue(vectors.Confidentiality),
		Cvss:               types.Float64Value(float64(vectors.Cvss)),
		Integrity:          types.StringValue(vectors.Integrity),
		PrivilegesRequired: types.StringValue(vectors.PrivilegesRequired),
		Scope:              types.StringValue(vectors.Scope),
		UserInteraction:    types.StringValue(vectors.UserInteraction),
	}
	return types.ObjectValueFrom(ctx, VectorsModelType.AttrTypes, model)
}



