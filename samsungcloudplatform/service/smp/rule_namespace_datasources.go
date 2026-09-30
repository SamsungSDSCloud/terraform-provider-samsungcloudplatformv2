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
	_ datasource.DataSource              = &smpRuleNamespaceDataSources{}
	_ datasource.DataSourceWithConfigure = &smpRuleNamespaceDataSources{}
)

func NewSmpRuleNamespaceDataSources() datasource.DataSource {
	return &smpRuleNamespaceDataSources{}
}

type smpRuleNamespaceDataSources struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpRuleNamespaceDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_rule_namespaces"
}

func (d *smpRuleNamespaceDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP RuleNamespaces Data Source",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter by name.",
				Optional:    true,
			},
			"name_like": schema.StringAttribute{
				Description: "Wildcard search for rule namespace names.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "Filter by rule namespace ID.",
				Optional:    true,
			},
			"id_like": schema.StringAttribute{
				Description: "Wildcard search for rule namespace IDs.",
				Optional:    true,
			},
			"page": schema.Int64Attribute{
				Description: "Page number for pagination (default: 1).",
				Optional:    true,
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Number of results per page.",
				Optional:    true,
				Computed:    true,
			},
			"sort_by": schema.StringAttribute{
				Description: "Field to sort by (e.g. name, created_at).",
				Optional:    true,
			},
			"sort_order": schema.StringAttribute{
				Description: "Sort order: asc or desc.",
				Optional:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID (required to list rule namespaces).",
				Required:    true,
			},
			"total_count": schema.Int64Attribute{
				Description: "Total number of results.",
				Computed:    true,
			},
			"sort": schema.ListNestedAttribute{
				Description: "Sort order returned by the API.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{Computed: true},
					},
				},
			},
			"rule_namespaces": schema.ListNestedAttribute{
				Description: "List of rule namespaces.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true},
						"name":         schema.StringAttribute{Computed: true},
						"state":        schema.StringAttribute{Computed: true},
						"workspace_id": schema.StringAttribute{Computed: true},
						"created_at":   schema.StringAttribute{Computed: true},
						"modified_at":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *smpRuleNamespaceDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smpRuleNamespaceDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.RuleNamespaceDataSources
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRuleNamespaceList(ctx, state.WorkspaceId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read RuleNamespaces", err.Error())
		return
	}

	var rnList []types.Object
	var sortList []types.Object
	if data != nil {
		state.TotalCount = types.Int64Value(int64(data.GetCount()))
		state.Page = types.Int64Value(int64(data.GetPage()))
		state.Size = types.Int64Value(int64(data.GetSize()))

		if sort, ok := data.GetSortOk(); ok && sort != nil {
			for _, s := range sort {
				obj, d := types.ObjectValue(
					map[string]attr.Type{"value": types.StringType},
					map[string]attr.Value{"value": types.StringValue(s)},
				)
				resp.Diagnostics.Append(d...)
				if resp.Diagnostics.HasError() {
					return
				}
				sortList = append(sortList, obj)
			}
		}
		if sortList == nil {
			sortList = []types.Object{}
		}
		sortVal, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.StringType}}, sortList)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Sort = sortVal

		for _, rn := range data.GetRuleNamespaces() {
			obj, diags := types.ObjectValue(
				map[string]attr.Type{
					"id":           types.StringType,
					"name":         types.StringType,
					"state":        types.StringType,
					"workspace_id": types.StringType,
					"created_at":   types.StringType,
					"modified_at":  types.StringType,
				},
				map[string]attr.Value{
					"id":           types.StringValue(rn.GetId()),
					"name":         types.StringValue(rn.GetName()),
					"state":        types.StringValue(rn.GetState()),
					"workspace_id": types.StringValue(rn.GetWorkspaceId()),
					"created_at":   types.StringValue(rn.GetCreatedAt().Format(TimeFormatDisplay)),
					"modified_at":  nullableTimeTypes(rn.GetModifiedAtOk()),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			rnList = append(rnList, obj)
		}
	}
	if rnList == nil {
		rnList = []types.Object{}
	}
	state.RuleNamespaces, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"state":        types.StringType,
		"workspace_id": types.StringType,
		"created_at":   types.StringType,
		"modified_at":  types.StringType,
	}}, rnList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
