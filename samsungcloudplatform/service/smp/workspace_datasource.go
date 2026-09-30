package smp

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

var (
	_ datasource.DataSource              = &smpWorkspaceDataSource{}
	_ datasource.DataSourceWithConfigure = &smpWorkspaceDataSource{}
)

func NewSmpWorkspaceDataSource() datasource.DataSource {
	return &smpWorkspaceDataSource{}
}

type smpWorkspaceDataSource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpWorkspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_workspace"
}

func (d *smpWorkspaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP Workspace Data Source",
		Attributes: map[string]schema.Attribute{
			"workspace": schema.SingleNestedAttribute{
				Description: "Workspace details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
					"state": schema.StringAttribute{Computed: true},
					"account_id": schema.StringAttribute{Computed: true},
					"created_at": schema.StringAttribute{Computed: true},
					"created_by": schema.StringAttribute{Computed: true},
					"modified_at": schema.StringAttribute{Computed: true},
					"modified_by": schema.StringAttribute{Computed: true},
					"private_acl_enabled": schema.BoolAttribute{Computed: true},
					"private_acl_resources": schema.ListNestedAttribute{
						Computed:     true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"resource_id":   schema.StringAttribute{Computed: true},
								"resource_name": schema.StringAttribute{Computed: true},
								"resource_type": schema.StringAttribute{Computed: true},
							},
						},
					},
					"private_prometheus_endpoint": schema.StringAttribute{Computed: true},
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (d *smpWorkspaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpWorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.WorkspaceDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetWorkspace(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspace, fmt.Sprintf(ErrReadWorkspaceFmt, state.WorkspaceId.ValueString(), err.Error(), detail))
		return
	}

	workspace := result.GetWorkspace()
	workspaceObj, diags := types.ObjectValue(
		smp.Workspace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                          types.StringValue(workspace.GetId()),
			"name":                        types.StringValue(workspace.GetName()),
			"state":                       types.StringValue(workspace.GetState()),
			"account_id":                  types.StringValue(workspace.GetAccountId()),
			"created_at":                  types.StringValue(workspace.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":                  types.StringValue(workspace.GetCreatedBy()),
			"modified_at":                 nullableTimeTypes(workspace.GetModifiedAtOk()),
			"modified_by":                 nullableStringTypes(workspace.GetModifiedByOk()),
			"private_acl_enabled":         nullableBoolTypes(workspace.GetPrivateAclEnabledOk()),
			"private_acl_resources":       privateAclResourcesToList(ctx, workspace.GetPrivateAclResources(), &resp.Diagnostics),
			"private_prometheus_endpoint": nullableStringTypes(workspace.GetPrivatePrometheusEndpointOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Workspace = workspaceObj

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
