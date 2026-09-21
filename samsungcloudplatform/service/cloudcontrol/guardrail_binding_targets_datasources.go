package cloudcontrol

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &guardrailBindingTargetsDataSources{}
	_ datasource.DataSourceWithConfigure = &guardrailBindingTargetsDataSources{}
)

func NewGuardrailBindingTargetsDataSources() datasource.DataSource {
	return &guardrailBindingTargetsDataSources{}
}

type guardrailBindingTargetsDataSources struct {
	config *scpsdk.Configuration
	client *cloudcontrol.Client
}

func (r *guardrailBindingTargetsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_guardrail_bindings_targets"
}

func (r *guardrailBindingTargetsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of targets (organization units or accounts) bound to a guardrail. \n",
		Attributes: map[string]schema.Attribute{
			"guardrail_id": schema.StringAttribute{
				Required: true,
				Description: "Guardrail ID to query targets for. \n" +
					"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
			},
			"target_type": schema.StringAttribute{
				Required: true,
				Description: "Target type to filter by. \n" +
					"  - example : 'ACCOUNT' \n" +
					"  - allowed_values : ['ACCOUNT', 'OU'] \n",
			},
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID to filter by. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by target name. \n" +
					"  - example : 'ou-test' \n",
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
				Description: "Total number of targets bound to the guardrail. \n" +
					"  - example : '20' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria actually applied by the server. \n" +
					"  - example : ['created_at:asc'] \n",
			},
			"targets": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of targets bound to the guardrail. \n" +
					"  - example : '[{id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, name: score-account, email: score@samsung.com, parent_unit_id: ou-fc8c29a138d78e24bf1fa86812fc8b, parent_unit_name: parent-unit-name}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "ID of the target - Organization Account ID for accounts, " +
								"Organization Unit ID for OUs. \n" +
								"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the target - Organization Account Name for accounts, " +
								"Organization Unit Name for OUs. \n" +
								"  - example : 'score-account' \n",
						},
						"email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the account. Only populated when target_type=ACCOUNT. \n" +
								"  - example : 'score@samsung.com' \n",
						},
						"parent_unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Root or Parent Organization Unit ID. \n" +
								"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
						},
						"parent_unit_name": schema.StringAttribute{
							Computed: true,
							Description: "Root or Parent Organization Unit Name. \n" +
								"  - example : 'parent-unit-name' \n",
						},
					},
				},
			},
		},
	}
}

func (r *guardrailBindingTargetsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *guardrailBindingTargetsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.GuardrailBindingTargetsDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListTargetsForGuardrail(
		ctx,
		state.GuardrailId.ValueString(),
		state.TargetType.ValueString(),
		state.LandingZoneId.ValueString(),
		state.Name.ValueString(),
		state.Sort.ValueString(),
		int32(state.Size.ValueInt64()),
		int32(state.Page.ValueInt64()),
	)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Guardrail Binding Targets",
			"Could not read Guardrail Binding Targets, unexpected error: "+err.Error()+"\nReason: "+detail,
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

	var targets []cloudcontrol.TargetsForGuardrailItem

	// OU 타입 target
	if data.Targets.ArrayOfOrganizationUnitsForGuardrail != nil {
		for _, t := range *data.Targets.ArrayOfOrganizationUnitsForGuardrail {
			parentUnitId := ""
			parentUnitName := ""
			if t.ParentUnitId.IsSet() && t.ParentUnitId.Get() != nil {
				parentUnitId = *t.ParentUnitId.Get()
			}
			if t.ParentUnitName.IsSet() && t.ParentUnitName.Get() != nil {
				parentUnitName = *t.ParentUnitName.Get()
			}
			targets = append(targets, cloudcontrol.TargetsForGuardrailItem{
				Id:             types.StringValue(t.Id),
				Name:           types.StringValue(t.Name),
				Email:          types.StringNull(),
				ParentUnitId:   types.StringValue(parentUnitId),
				ParentUnitName: types.StringValue(parentUnitName),
			})
		}
	}

	// ACCOUNT 타입 target
	if data.Targets.ArrayOfAccountsForGuardrail != nil {
		for _, t := range *data.Targets.ArrayOfAccountsForGuardrail {
			targetName := ""
			if t.Name.IsSet() && t.Name.Get() != nil {
				targetName = *t.Name.Get()
			}
			email := ""
			if t.Email.IsSet() && t.Email.Get() != nil {
				email = *t.Email.Get()
			}
			targets = append(targets, cloudcontrol.TargetsForGuardrailItem{
				Id:             types.StringValue(t.Id),
				Name:           types.StringValue(targetName),
				Email:          types.StringValue(email),
				ParentUnitId:   types.StringValue(t.ParentUnitId),
				ParentUnitName: types.StringValue(t.ParentUnitName),
			})
		}
	}

	state.Targets = targets

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
