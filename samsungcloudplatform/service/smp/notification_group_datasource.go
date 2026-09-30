package smp

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &smpNotificationGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &smpNotificationGroupDataSource{}
)

func NewSmpNotificationGroupDataSource() datasource.DataSource {
	return &smpNotificationGroupDataSource{}
}

type smpNotificationGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpNotificationGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_notification_group"
}

func (d *smpNotificationGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP NotificationGroup Data Source",
		Attributes: map[string]schema.Attribute{
			"notification_group": schema.SingleNestedAttribute{
				Description: "Notification group details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":          schema.StringAttribute{Computed: true},
					"name":        schema.StringAttribute{Computed: true},
					"description": schema.StringAttribute{Computed: true},
					"account_id":  schema.StringAttribute{Computed: true},
					"created_at":  schema.StringAttribute{Computed: true},
					"created_by":  schema.StringAttribute{Computed: true},
					"modified_at": schema.StringAttribute{Computed: true},
					"modified_by": schema.StringAttribute{Computed: true},
					"recipients": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id":                schema.StringAttribute{Computed: true},
								"recipient_user_id": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"notification_group_id": schema.StringAttribute{
				Description: "Notification group ID.",
				Required:    true,
			},
		},
	}
}

func (d *smpNotificationGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpNotificationGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.NotificationGroupDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetNotificationGroup(ctx, state.NotificationGroupId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadNotificationGroup, fmt.Sprintf(ErrReadNotificationGroupFmt, state.NotificationGroupId.ValueString(), err.Error(), detail))
		return
	}

	ng := result.GetNotificationGroup()
	ngObj, diags := buildNotificationGroupObject(ctx, ng)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.NotificationGroup = ngObj

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
