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
	_ datasource.DataSource              = &scrTagPackagesDataSource{}
	_ datasource.DataSourceWithConfigure = &scrTagPackagesDataSource{}
)

func NewScrTagPackagesDataSource() datasource.DataSource {
	return &scrTagPackagesDataSource{}
}

type scrTagPackagesDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrTagPackagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_tag_packages"
}

func (d *scrTagPackagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Tag Packages.",
		MarkdownDescription: "Get a list of packages for a specific tag in a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"tags_id": schema.StringAttribute{
				Description:         "Tags ID",
				MarkdownDescription: "The ID of the tag to get packages for.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"os_language": schema.StringAttribute{
				Description:         "OS or language filter",
				MarkdownDescription: "Filter packages by OS or language, e.g. `Java`.",
				Optional:            true,
			},
			"package_name": schema.StringAttribute{
				Description:         "Package name filter",
				MarkdownDescription: "Filter packages by name.",
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
			"filtered_count": schema.Int32Attribute{
				Description:         "Filtered count",
				MarkdownDescription: "Number of packages matching the filter.",
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
			"packages": schema.ListNestedAttribute{
				Description:         "Package list",
				MarkdownDescription: "List of packages found in the tag.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{Computed: true},
						"os_language": schema.StringAttribute{Computed: true},
						"package_name": schema.StringAttribute{Computed: true},
						"scan_summary": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: scanSummaryDataSourceSchemaAttrs(),
						},
						"type": schema.StringAttribute{Computed: true},
						"version": schema.StringAttribute{Computed: true},
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

func (d *scrTagPackagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrTagPackagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TagPackagesDataSource

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

	opts := scr.TagPackagesQueryOptions{
		OsLanguage:  osLanguage,
		PackageName: packageName,
		Sort:        sort,
		Page:        page,
		Size:        size,
	}

	listResp, err := d.client.ShowTagsPackages(ctx, state.TagsId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tag Packages",
			err.Error(),
		)
		return
	}

	state.FilteredCount = types.Int32Value(listResp.FilteredCount)
	state.LastScannedAt = types.StringValue(listResp.LastScannedAt)
	state.ReleaseVersion = types.StringValue(listResp.ReleaseVersion)

	packages := make([]PackageReportModel, 0, len(listResp.PackageReports))
	for _, pkg := range listResp.PackageReports {
		flattened, fDiags := flattenPackageReport(ctx, pkg)
		resp.Diagnostics.Append(fDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		packages = append(packages, flattened)
	}
	var pkgsDiags diag.Diagnostics
	state.Packages, pkgsDiags = types.ListValueFrom(ctx, PackageReportModelType, packages)
	resp.Diagnostics.Append(pkgsDiags...)

	var summaryDiags diag.Diagnostics
	state.ScanSummary, summaryDiags = flattenScanSummary(ctx, listResp.ScanSummary)
	resp.Diagnostics.Append(summaryDiags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenPackageReport(ctx context.Context, pkg scr11.PackageReport) (PackageReportModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	osLang := types.StringNull()
	if pkg.OsLanguage.IsSet() {
		val := pkg.OsLanguage.Get()
		if val != nil {
			osLang = types.StringValue(*val)
		}
	}
	pkgName := types.StringNull()
	if pkg.PackageName.IsSet() {
		val := pkg.PackageName.Get()
		if val != nil {
			pkgName = types.StringValue(*val)
		}
	}
	scanSummary, summaryDiags := flattenScanSummary(ctx, pkg.ScanSummary)
	diags.Append(summaryDiags...)
	return PackageReportModel{
		Category:    types.StringValue(pkg.Category),
		OsLanguage:  osLang,
		PackageName: pkgName,
		ScanSummary: scanSummary,
		Type:        types.StringValue(pkg.Type),
		Version:     types.StringValue(pkg.Version),
	}, diags
}
