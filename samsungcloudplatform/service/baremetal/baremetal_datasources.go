package baremetal

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &baremetalBaremetalDataSources{}
	_ datasource.DataSourceWithConfigure = &baremetalBaremetalDataSources{}
)

// NewBaremetalBaremetalDataSources is a helper function to simplify the provider implementation.
func NewBaremetalBaremetalDataSources() datasource.DataSource {
	return &baremetalBaremetalDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type baremetalBaremetalDataSources struct {
	config  *scpsdk.Configuration
	client  *baremetal.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *baremetalBaremetalDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_baremetals"
}

// Schema defines the schema for the data source.
func (d *baremetalBaremetalDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = BaremetalDataSourcesSchema()
}

// Configure adds the provider configured client to the data source.
func (d *baremetalBaremetalDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Baremetal
	d.clients = inst.Client
}

func BaremetalDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "list of baremetal server",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				Computed:            true,
				Description:         "Bare Metal Server ID List",
				MarkdownDescription: "Bare Metal Server ID List",
				ElementType:         types.StringType,
			},
			"policy_ip": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy IP\n  - example: 192.168.0.1",
				MarkdownDescription: "Policy IP\n  - example: 192.168.0.1",
			},
			"server_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Bare Metal Server name\n  - example: bmserver-001",
				MarkdownDescription: "Bare Metal Server name\n  - example: bmserver-001",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Description:         "Bare Metal Server state\n  - example: RUNNING",
				MarkdownDescription: "Bare Metal Server state\n  - example: RUNNING",
			},
			"vpc_id": schema.StringAttribute{
				Optional:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *baremetalBaremetalDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state baremetal.BaremetalDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetBaremetals(d.clients, state.Ip, state.ServerName, state.State, state.VpcId, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	state.Ids = ids

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetBaremetals(clients *client.SCPClient, ip types.String, serverName types.String, state types.String, vpcId types.String, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.Baremetal.GetBaremetalList(ip, serverName, state, vpcId)
	if err != nil {
		return nil, err
	}

	contents := data.Baremetals
	filteredContents := data.Baremetals

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, filters)

		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
