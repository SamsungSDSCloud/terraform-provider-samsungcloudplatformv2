package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &skeClusterUserKubeconfigDataSource{}
	_ datasource.DataSourceWithConfigure = &skeClusterUserKubeconfigDataSource{}
)

type skeClusterUserKubeconfigDataSource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func NewSkeClusterUserKubeconfigDataSource() datasource.DataSource {
	return &skeClusterUserKubeconfigDataSource{}
}

//// datasource.DataSource Interface Methods

func (d *skeClusterUserKubeconfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_user_kubeconfig"
}

func (d *skeClusterUserKubeconfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cluster User Kubeconfig",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ClusterId"): schema.StringAttribute{
				Description: "Cluster Id",
				Required:    true,
			},
			common.ToSnakeCase("KubeconfigType"): schema.StringAttribute{
				Description: "Kubeconfig Type",
				Required:    true,
			},
			common.ToSnakeCase("Kubeconfig"): schema.StringAttribute{
				Description: "Kubeconfig",
				Computed:    true,
			},
		},
	}
}

func (d *skeClusterUserKubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.ClusterUserKubeconfigDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userKubeconfig, err := d.clients.Ske.GetUserKubeConfig(ctx, state.ClusterId.ValueString(), state.KubeconfigType.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read User Kubeconfig",
			"Unable to read user kubeconfig: "+detail,
		)
	}

	state.Kubeconfig = types.StringValue(userKubeconfig)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//// datasource.DataSourceWithConfigure Interface Methods

func (d *skeClusterUserKubeconfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
