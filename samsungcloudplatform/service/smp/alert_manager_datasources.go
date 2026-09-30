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
	_ datasource.DataSource              = &smpAlertManagerDataSources{}
	_ datasource.DataSourceWithConfigure = &smpAlertManagerDataSources{}
)

func NewSmpAlertManagerDataSources() datasource.DataSource {
	return &smpAlertManagerDataSources{}
}

type smpAlertManagerDataSources struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpAlertManagerDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_alert_managers"
}

func (d *smpAlertManagerDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP AlertManagers Data Source",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
			"alert_managers": schema.ListNestedAttribute{
				Description: "List of alert managers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
			},
		},
	}
}

func (d *smpAlertManagerDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpAlertManagerDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.AlertManagerDataSources
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
		map[string]attr.Type{
			"id":                    types.StringType,
			"config_data":           types.StringType,
			"state":                 types.StringType,
			"workspace_id":          types.StringType,
			"notification_group_id": types.StringType,
			"created_at":            types.StringType,
			"created_by":            types.StringType,
			"modified_at":           types.StringType,
			"modified_by":           types.StringType,
		},
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

	alertManagers, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":                    types.StringType,
		"config_data":           types.StringType,
		"state":                 types.StringType,
		"workspace_id":          types.StringType,
		"notification_group_id": types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}}, []types.Object{amObj})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AlertManagers = alertManagers
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
