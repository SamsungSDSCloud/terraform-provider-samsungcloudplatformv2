package cloudcontrol

import (
	"context"
	"fmt"
	"strings"
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
	_ datasource.DataSource              = &guardrailListDataSource{}
	_ datasource.DataSourceWithConfigure = &guardrailListDataSource{}
)

func NewGuardrailListDataSource() datasource.DataSource {
	return &guardrailListDataSource{}
}

type guardrailListDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (d *guardrailListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_guardrails"
}

func (d *guardrailListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudControl
	d.clients = inst.Client
}

func (d *guardrailListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List CloudControl Guardrails with optional filters. \n" +
			"모든 조회 파라미터는 선택값이다 - 아무것도 지정하지 않으면 전체 목록을 조회한다.",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID to filter by. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Guardrail name to filter by. \n" +
					"  - example : 'foo-guardrail' \n",
			},
			"exclude_unit_id": schema.StringAttribute{
				Optional: true,
				Description: "Linked(Exclusion) Unit ID to filter by. \n" +
					"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
			},
			"guidance": schema.StringAttribute{
				Optional: true,
				Description: "Policy guidance to filter by. \n" +
					"  - allowed values : 'MANDATORY', 'OPTIONAL' \n" +
					"  - example : 'MANDATORY' \n",
			},
			"service_name": schema.StringAttribute{
				Optional: true,
				Description: "Service name to filter by. \n" +
					"  - example : 'iam' \n",
			},
			"status": schema.StringAttribute{
				Optional: true,
				Description: "Guardrail application status to filter by. \n" +
					"  - allowed values : 'ENABLED', 'DISABLED' \n" +
					"  - example : 'ENABLED' \n",
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
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria for the request, formatted as 'field:direction'. \n" +
					"  - example : 'created_at:desc' \n",
			},

			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total number of guardrails matching the filter. \n" +
					"  - example : '20' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria actually applied by the server, as returned in the response. \n" +
					"  - example : ['created_at:asc'] \n",
			},
			"guardrails": schema.ListNestedAttribute{
				Description: "List of guardrails matching the filter. \n" +
					"  - example : '[{id: 2c8a138f8d78e1fc29a449dbfa8681, name: foo-guardrail, status: DISABLED, ...}]' \n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the guardrail. \n" +
								"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Guardrail name. \n" +
								"  - example : 'foo-guardrail' \n",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Guardrail description. \n" +
								"  - example : 'This is an example guardrail.' \n",
						},
						"guidance": schema.StringAttribute{
							Computed: true,
							Description: "Policy guidance. \n" +
								"  - allowed values : 'MANDATORY', 'OPTIONAL' \n" +
								"  - example : 'MANDATORY' \n",
						},
						"service_name": schema.StringAttribute{
							Computed: true,
							Description: "Service name that this guardrail applies to. \n" +
								"  - example : 'iam' \n",
						},
						"status": schema.StringAttribute{
							Computed: true,
							Description: "Guardrail status. \n" +
								"  - allowed values : 'ENABLED', 'DISABLED' \n" +
								"  - example : 'DISABLED' \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Guardrail type. \n" +
								"  - allowed values : 'PREVENTIVE', 'DETECTIVE', 'PROACTIVE' \n" +
								"  - example : 'PREVENTIVE' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Created at timestamp. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "Created by user. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Modified at timestamp. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "Modified by user. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"srn": schema.StringAttribute{
							Computed: true,
							Description: "Samsung Resource Name (SRN). \n" +
								"  - example : 'srn:dev2:::::cloudcontrol:guardrail/2c8a138f8d78e1fc29a449dbfa8681' \n",
						},
						"binding_ous": schema.ListNestedAttribute{
							Computed: true,
							Description: "Organization Units this guardrail is bound to. \n" +
								"  - example : '[{id: ou-c29a138f8f1d78e24dbfa8681fc2fc8, name: Sandbox}]' \n",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Organization Unit ID. \n" +
											"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
									},
									"name": schema.StringAttribute{
										Computed: true,
										Description: "Organization Unit Name. \n" +
											"  - example : 'Sandbox' \n",
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

func (d *guardrailListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.GuardrailListDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()
	name := state.Name.ValueString()
	excludeUnitId := state.ExcludeUnitId.ValueString()
	guidance := state.Guidance.ValueString()
	serviceName := state.ServiceName.ValueString()
	status := state.Status.ValueString()
	sort := state.Sort.ValueString()
	size := state.Size.ValueInt64()
	page := state.Page.ValueInt64()

	data, err := d.client.ListGuardrails(ctx, landingZoneId, name, excludeUnitId, guidance, serviceName, status, sort, int32(size), int32(page))
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Guardrails",
			"Could not read CloudControl Guardrails, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var results []cloudcontrol.GuardrailListItem
	for _, g := range data.Guardrails {
		results = append(results, d.buildGuardrailItem(&g))
	}
	state.Guardrails = results
	state.Size = types.Int64Value(int64(data.Size))
	state.Page = types.Int64Value(int64(data.Page))

	state.TotalCount = types.Int64Value(int64(data.Count))
	if data.Sort != nil {
		sortVals, sortDiags := types.ListValueFrom(ctx, types.StringType, data.Sort)
		resp.Diagnostics.Append(sortDiags...)
		state.SortResult = sortVals
	} else {
		state.SortResult = types.ListNull(types.StringType)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *guardrailListDataSource) buildGuardrailItem(g *sdkcloudcontrol.GuardrailSummary) cloudcontrol.GuardrailListItem {
	item := cloudcontrol.GuardrailListItem{
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

	if g.Description.IsSet() && g.Description.Get() != nil && *g.Description.Get() != "" {
		item.Description = types.StringValue(strings.TrimSpace(*g.Description.Get()))
	} else {
		item.Description = types.StringNull()
	}

	if g.Status != nil && *g.Status != "" {
		item.Status = types.StringValue(*g.Status)
	} else {
		item.Status = types.StringNull()
	}

	if g.Srn.IsSet() && g.Srn.Get() != nil && *g.Srn.Get() != "" {
		item.Srn = types.StringValue(*g.Srn.Get())
	} else {
		item.Srn = types.StringNull()
	}

	if g.BindingOus.IsSet() && g.BindingOus.Get() != nil {
		ous := g.BindingOus.Get()
		if ous.ArrayOfGuardrailBindingOu != nil {
			for _, ou := range *ous.ArrayOfGuardrailBindingOu {
				item.BindingOus = append(item.BindingOus, cloudcontrol.GuardrailBindingOuItem{
					Id:   types.StringValue(ou.Id),
					Name: types.StringValue(ou.Name),
				})
			}
		}
	}

	return item
}
