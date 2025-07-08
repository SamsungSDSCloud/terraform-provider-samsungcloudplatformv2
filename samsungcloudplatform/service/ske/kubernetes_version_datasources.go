package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/region"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/ske/1.0"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &skeKubernetesVersionDataSources{}
	_ datasource.DataSourceWithConfigure = &skeKubernetesVersionDataSources{}
)

type skeKubernetesVersionDataSources struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func NewSkeKubernetesVersionDataSources() datasource.DataSource {
	return &skeKubernetesVersionDataSources{}
}

//// datasource.DataSource Interface Methods

func (d *skeKubernetesVersionDataSources) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_kubernetes_versions"
}

func (d *skeKubernetesVersionDataSources) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kubernetes Version List",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			common.ToSnakeCase("KubernetesVersions"): schema.ListNestedAttribute{
				Description: "Kubernetes Version List",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
							Description: "Kubernetes Version",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *skeKubernetesVersionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.KubernetesVersionDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !state.Region.IsNull() {
		d.client.Config.Region = state.Region.ValueString()
	}

	kubernetesVersionListResponse, err := d.clients.Ske.GetKubernetesVersionList(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Kubernetes Versions",
			err.Error(),
		)
		return
	}
	kubernetesVersions := kubernetesVersionListResponse.GetKubernetesVersions()
	var kubernetesVersionsModels []ske.KubernetesVersionSummary
	for _, kubernetesVersion := range kubernetesVersions {
		kubernetesVersionsModels = append(kubernetesVersionsModels, d.makeKubernetesVersionModel(&kubernetesVersion))
	}
	state.KubernetesVersions = kubernetesVersionsModels
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *skeKubernetesVersionDataSources) makeKubernetesVersionModel(kubernetesVersion *scpske.KubernetesVersionSummary) ske.KubernetesVersionSummary {
	return ske.KubernetesVersionSummary{
		Description:       types.StringValue(kubernetesVersion.GetDescription()),
		KubernetesVersion: types.StringValue(kubernetesVersion.GetKubernetesVersion()),
	}
}

//// datasource.DataSourceWithConfigure Interface Methods

func (d *skeKubernetesVersionDataSources) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
