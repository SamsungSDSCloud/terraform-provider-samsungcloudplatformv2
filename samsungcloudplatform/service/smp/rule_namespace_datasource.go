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
	_ datasource.DataSource              = &smpRuleNamespaceDataSource{}
	_ datasource.DataSourceWithConfigure = &smpRuleNamespaceDataSource{}
)

func NewSmpRuleNamespaceDataSource() datasource.DataSource {
	return &smpRuleNamespaceDataSource{}
}

type smpRuleNamespaceDataSource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpRuleNamespaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_rule_namespace"
}

func (d *smpRuleNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP RuleNamespace Data Source",
		Attributes: map[string]schema.Attribute{
			"rule_namespace": schema.SingleNestedAttribute{
				Description: "Rule namespace details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Computed: true},
					"name":         schema.StringAttribute{Computed: true},
					"config_data":  schema.StringAttribute{Computed: true},
					"state":        schema.StringAttribute{Computed: true},
					"workspace_id": schema.StringAttribute{Computed: true},
					"created_at":   schema.StringAttribute{Computed: true},
					"created_by":   schema.StringAttribute{Computed: true},
					"modified_at":  schema.StringAttribute{Computed: true},
					"modified_by":  schema.StringAttribute{Computed: true},
				},
			},
			"rule_namespace_id": schema.StringAttribute{
				Description: "Rule namespace ID.",
				Required:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (d *smpRuleNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpRuleNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.RuleNamespaceDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetRuleNamespace(ctx, state.WorkspaceId.ValueString(), state.RuleNamespaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadRuleNamespace, fmt.Sprintf(ErrReadRuleNamespaceFmt, state.RuleNamespaceId.ValueString(), err.Error(), detail))
		return
	}

	rn := result.GetRuleNamespace()
	rnObj, diags := types.ObjectValue(
		smp.RuleNamespace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":           types.StringValue(rn.GetId()),
			"name":         types.StringValue(rn.GetName()),
			"config_data":  types.StringValue(rn.GetConfigData()),
			"state":        types.StringValue(rn.GetState()),
			"workspace_id": types.StringValue(rn.GetWorkspaceId()),
			"created_at":   types.StringValue(rn.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":   types.StringValue(rn.GetCreatedBy()),
			"modified_at":  nullableTimeTypes(rn.GetModifiedAtOk()),
			"modified_by":  nullableStringTypes(rn.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.RuleNamespace = rnObj
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
