package vertica

import (
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &verticaEngineVersionDataSources{}
	_ datasource.DataSourceWithConfigure = &verticaEngineVersionDataSources{}
)

func NewVerticaEngineVersionDataSources() datasource.DataSource {
	return &verticaEngineVersionDataSources{}
}

type verticaEngineVersionDataSources struct {
	config  *scpsdk.Configuration
	client  *vertica.Client
	clients *client.SCPClient
}

func (d *verticaEngineVersionDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vertica_engine_version"
}

func (d *verticaEngineVersionDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Engine Versions.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "Engine version ID to filter by\n  - example: d058dc79a86842f9b558933e9d285b9f",
				MarkdownDescription: "Engine version ID to filter by\n  - example: d058dc79a86842f9b558933e9d285b9f",
				Optional:            true,
			},
			common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
				Description:         "Product image type. Omit to return every image type.\n  - example: Vertica Enterprise / Vertica Community",
				MarkdownDescription: "Product image type. Omit to return every image type.\n  - example: Vertica Enterprise / Vertica Community",
				Optional:            true,
			},
			common.ToSnakeCase("EosIncluded"): schema.BoolAttribute{
				Description:         "Whether to include end-of-service versions\n  - example: false",
				MarkdownDescription: "Whether to include end-of-service versions\n  - example: false",
				Optional:            true,
			},
			common.ToSnakeCase("Contents"): schema.ListNestedAttribute{
				Description: "A detail of Engine Version.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("EndOfService"): schema.BoolAttribute{
							Description:         "End of Service\n  - example: false",
							MarkdownDescription: "End of Service\n  - example: false",
							Required:            true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Engine version ID\n  - example: 09c2fe88089040ffa035604e38f7e4e9",
							MarkdownDescription: "Engine version ID\n  - example: 09c2fe88089040ffa035604e38f7e4e9",
							Required:            true,
						},
						common.ToSnakeCase("MajorVersion"): schema.StringAttribute{
							Description:         "Software major version\n  - example: ENTERPRISE 24.2",
							MarkdownDescription: "Software major version\n  - example: ENTERPRISE 24.2",
							Optional:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "Engine version name\n  - example: Vertica ENTERPRISE 24.2.0-2",
							MarkdownDescription: "Engine version name\n  - example: Vertica ENTERPRISE 24.2.0-2",
							Required:            true,
						},
						common.ToSnakeCase("OsType"): schema.StringAttribute{
							Description:         "OS type\n  - example: RHEL",
							MarkdownDescription: "OS type\n  - example: RHEL",
							Required:            true,
						},
						common.ToSnakeCase("OsVersion"): schema.StringAttribute{
							Description:         "OS version\n  - example: 8.5",
							MarkdownDescription: "OS version\n  - example: 8.5",
							Optional:            true,
						},
						common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
							Description:         "Product image type\n  - example: Vertica Enterprise",
							MarkdownDescription: "Product image type\n  - example: Vertica Enterprise",
							Required:            true,
						},
						common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
							Description:         "Software version\n  - example: ENTERPRISE 24.2.0-2",
							MarkdownDescription: "Software version\n  - example: ENTERPRISE 24.2.0-2",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *verticaEngineVersionDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Vertica
	d.clients = inst.Client
}

func (d *verticaEngineVersionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vertica.EngineVersionDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var eosIncluded *bool
	if !state.EosIncluded.IsNull() && !state.EosIncluded.IsUnknown() {
		v := state.EosIncluded.ValueBool()
		eosIncluded = &v
	}

	data, err := d.client.GetEngineVersionList(ctx, state.Id.ValueString(), state.ProductImageType.ValueString(), eosIncluded)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read EngineVersion",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineVersionElement := range data.Contents {
		engineVersionState := vertica.EngineVersion{
			EndOfService:     types.BoolPointerValue(engineVersionElement.EndOfService),
			Id:               types.StringValue(engineVersionElement.Id),
			MajorVersion:     types.StringValue(engineVersionElement.MajorVersion),
			Name:             types.StringValue(engineVersionElement.Name),
			OsType:           types.StringValue(engineVersionElement.OsType),
			OsVersion:        types.StringValue(engineVersionElement.OsVersion),
			ProductImageType: types.StringPointerValue(engineVersionElement.ProductImageType.Get()),
			SoftwareVersion:  types.StringValue(engineVersionElement.SoftwareVersion),
		}
		state.Contents = append(state.Contents, engineVersionState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
