package epas

import (
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &epasEngineVersionDataSources{}
	_ datasource.DataSourceWithConfigure = &epasEngineVersionDataSources{}
)

func NewEpasEngineVersionDataSources() datasource.DataSource {
	return &epasEngineVersionDataSources{}
}

type epasEngineVersionDataSources struct {
	config  *scpsdk.Configuration
	client  *epas.Client
	clients *client.SCPClient
}

func (d *epasEngineVersionDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_epas_engine_version"
}

func (d *epasEngineVersionDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Engine Versions.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Contents"): schema.ListNestedAttribute{
				Description: "A detail of Engine Version.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("EndOfService"): schema.BoolAttribute{
							Description: "EndOfService",
							Required:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Required:    true,
						},
						common.ToSnakeCase("MajorVersion"): schema.StringAttribute{
							Description: "MajorVersion",
							Optional:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Required:    true,
						},
						common.ToSnakeCase("OsType"): schema.StringAttribute{
							Description: "OsType",
							Required:    true,
						},
						common.ToSnakeCase("OsVersion"): schema.StringAttribute{
							Description: "OsVersion",
							Optional:    true,
						},
						common.ToSnakeCase("ProductImageType"): schema.StringAttribute{
							Description: "ProductImageType",
							Required:    true,
						},
						common.ToSnakeCase("SoftwareVersion"): schema.StringAttribute{
							Description: "SoftwareVersion",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (d *epasEngineVersionDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Epas
	d.clients = inst.Client
}

func (d *epasEngineVersionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state epas.EngineVersionDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEngineVersionList(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read EngineVersion",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineVersionElement := range data.Contents {
		engineVersionState := epas.EngineVersion{
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
