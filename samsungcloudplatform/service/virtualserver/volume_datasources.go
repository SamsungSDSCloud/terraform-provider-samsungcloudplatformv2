package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualServerVolumeDataSources{}
	_ datasource.DataSourceWithConfigure = &virtualServerVolumeDataSources{}
)

// NewComputeVolumeDataSource is a helper function to simplify the provider implementation.
func NewVirtualServerVolumeDataSources() datasource.DataSource {
	return &virtualServerVolumeDataSources{}
}

// virtualServerVolumeDataSource is the data source implementation.
type virtualServerVolumeDataSources struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *virtualServerVolumeDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_volumes"
}

// Schema defines the schema for the data source.
func (d *virtualServerVolumeDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of volumes.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("Bootable"): schema.BoolAttribute{
				Description: "Bootable",
				Optional:    true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Server ID List",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *virtualServerVolumeDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.VirtualServer
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *virtualServerVolumeDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.VolumeDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetVolumes(d.clients, state.Name, state.State, state.Bootable, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Volumes",
			err.Error(),
		)
		return
	}
	state.Ids = ids
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetVolumes(clients *client.SCPClient, Name types.String, State types.String, Bootable types.Bool, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.VirtualServer.GetVolumeListWithParam(Name, State, Bootable)
	if err != nil {
		return nil, err
	}

	contents := data.Volumes
	filteredContents := data.Volumes

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, filters) // 필터링된 컨텐츠의 Index 정보 리턴

		for i, resource := range contents {
			if common.Contains(indices, i) { // Index 정보 기준으로 필터링 진행
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
