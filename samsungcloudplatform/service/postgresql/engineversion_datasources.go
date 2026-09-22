package postgresql

import (
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &postgresqlEngineVersionDataSources{}
	_ datasource.DataSourceWithConfigure = &postgresqlEngineVersionDataSources{}
)

func NewPostgresqlEngineVersionDataSources() datasource.DataSource {
	return &postgresqlEngineVersionDataSources{}
}

type postgresqlEngineVersionDataSources struct {
	config  *scpsdk.Configuration
	client  *postgresql.Client
	clients *client.SCPClient
}

func (d *postgresqlEngineVersionDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_engine_version"
}

func (d *postgresqlEngineVersionDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Engine Versions.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "Engine version ID to filter by\n  - example: d058dc79a86842f9b558933e9d285b9f",
				MarkdownDescription: "Engine version ID to filter by\n  - example: d058dc79a86842f9b558933e9d285b9f",
				Optional:            true,
			},
			common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
				Description:         "Product image type\n  - example: PostgreSQL Community",
				MarkdownDescription: "Product image type\n  - example: PostgreSQL Community",
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
							Description:         "Engine version ID\n  - example: d058dc79a86842f9b558933e9d285b9f",
							MarkdownDescription: "Engine version ID\n  - example: d058dc79a86842f9b558933e9d285b9f",
							Required:            true,
						},
						common.ToSnakeCase("MajorVersion"): schema.StringAttribute{
							Description:         "Software major version\n  - example: COMMUNITY 17",
							MarkdownDescription: "Software major version\n  - example: COMMUNITY 17",
							Optional:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "Engine version name\n  - example: PostgreSQL Community 17.7",
							MarkdownDescription: "Engine version name\n  - example: PostgreSQL Community 17.7",
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
							Description:         "Product image type\n  - example: PostgreSQL Community",
							MarkdownDescription: "Product image type\n  - example: PostgreSQL Community",
							Required:            true,
						},
						common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
							Description:         "Software version\n  - example: COMMUNITY 17.7",
							MarkdownDescription: "Software version\n  - example: COMMUNITY 17.7",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *postgresqlEngineVersionDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Postgresql
	d.clients = inst.Client
}

func (d *postgresqlEngineVersionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresql.EngineVersionDataSource

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
		engineVersionState := postgresql.EngineVersion{
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
