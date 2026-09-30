package smp

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &smpWorkspaceConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &smpWorkspaceConfigurationDataSource{}
)

func NewSmpWorkspaceConfigurationDataSource() datasource.DataSource {
	return &smpWorkspaceConfigurationDataSource{}
}

type smpWorkspaceConfigurationDataSource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpWorkspaceConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_workspace_configuration"
}

func (d *smpWorkspaceConfigurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP WorkspaceConfiguration Data Source",
		Attributes: map[string]schema.Attribute{
			"workspace_configuration": schema.SingleNestedAttribute{
				Description: "Workspace configuration details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":               schema.StringAttribute{Computed: true},
					"state":            schema.StringAttribute{Computed: true},
					"workspace_id":     schema.StringAttribute{Computed: true},
					"retention_period": schema.Int64Attribute{Computed: true},
					"created_at":       schema.StringAttribute{Computed: true},
					"created_by":       schema.StringAttribute{Computed: true},
					"modified_at":      schema.StringAttribute{Computed: true},
					"modified_by":      schema.StringAttribute{Computed: true},
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (d *smpWorkspaceConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(ErrUnexpectedConfigure, fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData))
		return
	}
	d.client = inst.Client.Smp
	d.clients = inst.Client
}

func (d *smpWorkspaceConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.WorkspaceConfigurationDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetWorkspaceConfiguration(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspaceConfiguration, fmt.Sprintf(ErrReadWorkspaceConfigurationFmt, state.WorkspaceId.ValueString(), err.Error(), detail))
		return
	}

	wc := result.GetWorkspaceConfiguration()
	wcObj, diags := types.ObjectValue(
		smp.WorkspaceConfiguration{}.AttributeTypes(),
		map[string]attr.Value{
			"id":               types.StringValue(wc.GetId()),
			"state":            types.StringValue(wc.GetState()),
			"workspace_id":     types.StringValue(wc.GetWorkspaceId()),
			"retention_period": types.Int64Value(int64(wc.GetRetentionPeriod())),
			"created_at":       types.StringValue(wc.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":       types.StringValue(wc.GetCreatedBy()),
			"modified_at":      nullableTimeTypes(wc.GetModifiedAtOk()),
			"modified_by":      nullableStringTypes(wc.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.WorkspaceConfiguration = wcObj
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
