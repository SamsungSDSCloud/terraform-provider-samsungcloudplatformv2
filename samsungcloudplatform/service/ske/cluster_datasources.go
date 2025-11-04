package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeClusterDataSources{}
	_ datasource.DataSourceWithConfigure = &skeClusterDataSources{}
)

// skeClusterDataSources is the data source implementation.
type skeClusterDataSources struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func NewSkeClusterDataSources() datasource.DataSource {
	return &skeClusterDataSources{}
}

//// datasource.DataSource Interface Methods

func (d *skeClusterDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_clusters"
}

// Schema 에는 tf 파일에서 허용되는 설정, 리소스의 상태와 속성을 정의한다.
func (d *skeClusterDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of cluster.",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			"tags":   tag.DataSourceSchema(),
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
				Validators:  []validator.String{
					// regex: r"^[a-z0-9\-]*$",
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "SubnetId (validation)",
				Optional:    true,
				Validators:  []validator.String{},
			},
			common.ToSnakeCase("Status"): schema.ListAttribute{
				Description: "Status List",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("KubernetesVersion"): schema.ListAttribute{
				Description: "KubernetesVersion List",
				Optional:    true,
				ElementType: types.StringType,
			},
			// response
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "ID List",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *skeClusterDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.ClusterDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		d.client.Config.Region = state.Region.ValueString()
	}

	ids, err := getClusters(ctx, d.clients, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Cluster",
			err.Error(),
		)
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//// datasource.DataSourceWithConfigure Interface Methods

func (d *skeClusterDataSources) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Ske
	d.clients = inst.Client
}

//// private

const idField = "id"

func getClusters(ctx context.Context, clients *client.SCPClient, state ske.ClusterDataSourceIds) ([]types.String, error) {
	data, err := clients.Ske.GetClusterList(ctx, state)
	if err != nil {
		return nil, err
	}
	contents := data.Clusters
	filteredContents := data.Clusters

	tags := state.Tags.Elements()
	if len(tags) > 0 {
		filteredContents = filteredContents[:0]
		indices := tag.GetTagIndices(clients, contents, tags, idField)
		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	filters := state.Filter
	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, state.Filter)
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
		ids = append(ids, types.StringPointerValue(content.Id))
	}
	return ids, nil
}
