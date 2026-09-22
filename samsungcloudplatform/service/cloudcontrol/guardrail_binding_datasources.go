package cloudcontrol

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkcloudcontrol "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &guardrailBindingDataSources{}
	_ datasource.DataSourceWithConfigure = &guardrailBindingDataSources{}
)

func NewGuardrailBindingDataSources() datasource.DataSource {
	return &guardrailBindingDataSources{}
}

type guardrailBindingDataSources struct {
	config *scpsdk.Configuration
	client *cloudcontrol.Client
}

func (r *guardrailBindingDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_guardrail_bindings_guardrails"
}

func (r *guardrailBindingDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of guardrails bound to a target (organization unit). \n",
		Attributes: map[string]schema.Attribute{
			"target_id": schema.StringAttribute{
				Required: true,
				Description: "Target organization unit ID to query guardrails for. \n" +
					"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
			},
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID that contains the organization unit. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by guardrail name. \n" +
					"  - example : 'guardrail-example' \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria for the request, formatted as 'field:direction'. \n" +
					"  - example : 'created_at:desc' \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size to request. Server default applies if omitted. \n" +
					"  - example : '20' \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number to request. Server default applies if omitted. \n" +
					"  - example : '0' \n",
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total number of guardrails bound to the target. \n" +
					"  - example : '20' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria actually applied by the server. \n" +
					"  - example : ['created_at:asc'] \n",
			},
			"guardrails": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of guardrails bound to the target. \n" +
					"  - example : '[{id: 07fd90ddcc9c4f2f84de093c65535392, name: guardrail-example, guidance: MANDATORY, service_name: iam, type: PREVENTIVE, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the guardrail. \n" +
								"  - example : '07fd90ddcc9c4f2f84de093c65535392' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the guardrail. \n" +
								"  - example : 'guardrail-example' \n",
						},
						"guidance": schema.StringAttribute{
							Computed: true,
							Description: "Policy guidance of the guardrail. \n" +
								"  - allowed values : 'MANDATORY', 'OPTIONAL' \n" +
								"  - example : 'MANDATORY' \n",
						},
						"service_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the service the guardrail belongs to. \n" +
								"  - example : 'iam' \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type of the guardrail. \n" +
								"  - allowed values : 'PREVENTIVE', 'DETECTIVE', 'PROACTIVE' \n" +
								"  - example : 'PREVENTIVE' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the guardrail was bound to the target. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who bound the guardrail to the target. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the guardrail binding was last modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who last modified the guardrail binding. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"link_types": schema.SingleNestedAttribute{
							Computed: true,
							Description: "How this guardrail is linked to the target - directly bound, " +
								"or inherited from a parent OU. \n" +
								"  - example : '{directed: [{target_id: 1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: example-target}], inherited: [{target_id: ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: root}]}' \n",
							Attributes: map[string]schema.Attribute{
								"directed": schema.ListNestedAttribute{
									Computed: true,
									Description: "Targets this guardrail is directly bound to. \n" +
										"  - example : '[{target_id: 1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: example-target}]' \n",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_id": schema.StringAttribute{
												Computed: true,
												Description: "ID of the target the guardrail is bound to. \n" +
													"  - example : '1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
											},
											"target_name": schema.StringAttribute{
												Computed: true,
												Description: "Name of the target the guardrail is bound to. \n" +
													"  - example : 'example-target' \n",
											},
										},
									},
								},
								"inherited": schema.ListNestedAttribute{
									Computed: true,
									Description: "Targets this guardrail is inherited from. \n" +
										"  - example : '[{target_id: ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: root}]' \n",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_id": schema.StringAttribute{
												Computed: true,
												Description: "ID of the target the guardrail is inherited from. \n" +
													"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
											},
											"target_name": schema.StringAttribute{
												Computed: true,
												Description: "Name of the target the guardrail is inherited from. \n" +
													"  - example : 'root' \n",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *guardrailBindingDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = inst.Client.CloudControl
}

func (r *guardrailBindingDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.GuardrailBindingDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListGuardrailsForTarget(
		ctx,
		state.TargetId.ValueString(),
		state.LandingZoneId.ValueString(),
		state.Name.ValueString(),
		state.Sort.ValueString(),
		int32(state.Size.ValueInt64()),
		int32(state.Page.ValueInt64()),
	)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Guardrail Bindings",
			"Could not read Guardrail Bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.TotalCount = types.Int64Value(int64(data.Count))
	state.Size = types.Int64Value(int64(data.Size))
	state.Page = types.Int64Value(int64(data.Page))

	if data.Sort != nil {
		sortVals, sortDiags := types.ListValueFrom(ctx, types.StringType, data.Sort)
		resp.Diagnostics.Append(sortDiags...)
		state.SortResult = sortVals
	} else {
		state.SortResult = types.ListNull(types.StringType)
	}

	var guardrails []cloudcontrol.GuardrailsForTargetItem
	for _, g := range data.Guardrails {
		guardrails = append(guardrails, buildGuardrailForTargetItem(&g))
	}
	state.Guardrails = guardrails

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildGuardrailForTargetItem(g *sdkcloudcontrol.GuardrailsForTargetSummary) cloudcontrol.GuardrailsForTargetItem {
	item := cloudcontrol.GuardrailsForTargetItem{
		Id:          types.StringValue(g.Id),
		Name:        types.StringValue(g.Name),
		Guidance:    types.StringValue(g.Guidance),
		ServiceName: types.StringValue(g.ServiceName),
		Type:        types.StringValue(g.Type),
		CreatedAt:   types.StringValue(g.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(g.CreatedBy),
		ModifiedAt:  types.StringValue(g.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(g.ModifiedBy),
	}

	buildLinkTargets := func(infos []sdkcloudcontrol.GuardrailLinkInfo) []cloudcontrol.GuardrailLinkTargetItem {
		var result []cloudcontrol.GuardrailLinkTargetItem
		for _, info := range infos {
			targetName := ""
			if info.TargetName.IsSet() && info.TargetName.Get() != nil {
				targetName = *info.TargetName.Get()
			}
			result = append(result, cloudcontrol.GuardrailLinkTargetItem{
				TargetId:   types.StringValue(info.TargetId),
				TargetName: types.StringValue(targetName),
			})
		}
		return result
	}

	item.LinkTypes = &cloudcontrol.GuardrailLinkTypesItem{
		Directed:  buildLinkTargets(g.LinkTypes["directed"]),
		Inherited: buildLinkTargets(g.LinkTypes["inherited"]),
	}

	return item
}
