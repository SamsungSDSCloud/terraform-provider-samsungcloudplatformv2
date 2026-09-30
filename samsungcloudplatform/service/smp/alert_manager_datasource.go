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
	_ datasource.DataSource              = &smpAlertManagerDataSource{}
	_ datasource.DataSourceWithConfigure = &smpAlertManagerDataSource{}
)

func NewSmpAlertManagerDataSource() datasource.DataSource {
	return &smpAlertManagerDataSource{}
}

type smpAlertManagerDataSource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpAlertManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_alert_manager"
}

func (d *smpAlertManagerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP AlertManager Data Source",
		Attributes: map[string]schema.Attribute{
			"alert_manager": schema.SingleNestedAttribute{
				Description: "Alert manager details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":                    schema.StringAttribute{Computed: true},
					"config_data":           schema.StringAttribute{Computed: true},
					"state":                 schema.StringAttribute{Computed: true},
					"workspace_id":          schema.StringAttribute{Computed: true},
					"notification_group_id": schema.StringAttribute{Computed: true},
					"created_at":            schema.StringAttribute{Computed: true},
					"created_by":            schema.StringAttribute{Computed: true},
					"modified_at":           schema.StringAttribute{Computed: true},
					"modified_by":           schema.StringAttribute{Computed: true},
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (d *smpAlertManagerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpAlertManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.AlertManagerDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetAlertManager(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadAlertManager, fmt.Sprintf(ErrReadAlertManagerFmt, state.WorkspaceId.ValueString(), err.Error(), detail))
		return
	}

	am := result.GetAlertManager()
	amObj, diags := types.ObjectValue(
		smp.AlertManager{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                    types.StringValue(am.GetId()),
			"config_data":           types.StringValue(am.GetConfigData()),
			"state":                 types.StringValue(am.GetState()),
			"workspace_id":          types.StringValue(am.GetWorkspaceId()),
			"notification_group_id": types.StringValue(am.GetNotificationGroupId()),
			"created_at":            types.StringValue(am.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":            types.StringValue(am.GetCreatedBy()),
			"modified_at":           nullableTimeTypes(am.GetModifiedAtOk()),
			"modified_by":           nullableStringTypes(am.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AlertManager = amObj
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
