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
	_ datasource.DataSource              = &smpWorkspaceConfigurationDataSources{}
	_ datasource.DataSourceWithConfigure = &smpWorkspaceConfigurationDataSources{}
)

func NewSmpWorkspaceConfigurationDataSources() datasource.DataSource {
	return &smpWorkspaceConfigurationDataSources{}
}

type smpWorkspaceConfigurationDataSources struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpWorkspaceConfigurationDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_workspace_configurations"
}

func (d *smpWorkspaceConfigurationDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP WorkspaceConfigurations Data Source",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
			"workspace_configurations": schema.ListNestedAttribute{
				Description: "List of workspace configurations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
			},
		},
	}
}

func (d *smpWorkspaceConfigurationDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpWorkspaceConfigurationDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.WorkspaceConfigurationDataSources
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
		map[string]attr.Type{
			"id":               types.StringType,
			"state":            types.StringType,
			"workspace_id":     types.StringType,
			"retention_period": types.Int64Type,
			"created_at":       types.StringType,
			"created_by":       types.StringType,
			"modified_at":      types.StringType,
			"modified_by":      types.StringType,
		},
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

	configurations, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":               types.StringType,
		"state":            types.StringType,
		"workspace_id":     types.StringType,
		"retention_period": types.Int64Type,
		"created_at":       types.StringType,
		"created_by":       types.StringType,
		"modified_at":      types.StringType,
		"modified_by":      types.StringType,
	}}, []types.Object{wcObj})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.WorkspaceConfigurations = configurations
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
