package sqlserver

import (
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &sqlserverEngineVersionDataSources{}
	_ datasource.DataSourceWithConfigure = &sqlserverEngineVersionDataSources{}
)

func NewSqlserverlEngineVersionDataSources() datasource.DataSource {
	return &sqlserverEngineVersionDataSources{}
}

type sqlserverEngineVersionDataSources struct {
	config  *scpsdk.Configuration
	client  *sqlserver.Client
	clients *client.SCPClient
}

func (d *sqlserverEngineVersionDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sqlserver_engine_version"
}

func (d *sqlserverEngineVersionDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Engine Versions.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
				Description:         "Product image type\n  - example: Microsoft SQL Server Standard",
				MarkdownDescription: "Product image type\n  - example: Microsoft SQL Server Standard",
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
							Description:         "Engine version ID\n  - example: 3692630f08884127a7168821b5654871",
							MarkdownDescription: "Engine version ID\n  - example: 3692630f08884127a7168821b5654871",
							Required:            true,
						},
						common.ToSnakeCase("MajorVersion"): schema.StringAttribute{
							Description:         "Software major version\n  - example: 2022 Standard ENG",
							MarkdownDescription: "Software major version\n  - example: 2022 Standard ENG",
							Optional:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "Engine version name\n  - example: Microsoft SQL Server 2022 Standard ENG-KB5065865-x64",
							MarkdownDescription: "Engine version name\n  - example: Microsoft SQL Server 2022 Standard ENG-KB5065865-x64",
							Required:            true,
						},
						common.ToSnakeCase("OsType"): schema.StringAttribute{
							Description:         "OS type\n  - example: WINDOWS",
							MarkdownDescription: "OS type\n  - example: WINDOWS",
							Required:            true,
						},
						common.ToSnakeCase("OsVersion"): schema.StringAttribute{
							Description:         "OS version\n  - example: 2019 Std.",
							MarkdownDescription: "OS version\n  - example: 2019 Std.",
							Optional:            true,
						},
						common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
							Description:         "Product image type\n  - example: Microsoft SQL Server Standard",
							MarkdownDescription: "Product image type\n  - example: Microsoft SQL Server Standard",
							Required:            true,
						},
						common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
							Description:         "Software version\n  - example: 2022 Standard ENG-KB5065865",
							MarkdownDescription: "Software version\n  - example: 2022 Standard ENG-KB5065865",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *sqlserverEngineVersionDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Sqlserver
	d.clients = inst.Client
}

func (d *sqlserverEngineVersionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sqlserver.EngineVersionDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 사용자가 product_image_type을 지정하지 않은 경우 기본값을 사용한다.
	if state.ProductImageType.IsNull() || state.ProductImageType.ValueString() == "" {
		state.ProductImageType = types.StringValue("Microsoft SQL Server Standard")
	}

	data, err := d.client.GetEngineVersionList(ctx, state.ProductImageType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read EngineVersion",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineVersionElement := range data.Contents {
		engineVersionState := sqlserver.EngineVersion{
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
