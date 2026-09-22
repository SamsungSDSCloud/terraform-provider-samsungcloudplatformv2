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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &guardrailDataSource{}
	_ datasource.DataSourceWithConfigure = &guardrailDataSource{}
)

func NewGuardrailDataSource() datasource.DataSource {
	return &guardrailDataSource{}
}

type guardrailDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (d *guardrailDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_guardrail"
}

func (d *guardrailDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *guardrailDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Look up a single CloudControl Guardrail by ID.",
		Attributes: map[string]schema.Attribute{
			"guardrail_id": schema.StringAttribute{
				Required: true,
				Description: "Unique identifier of the guardrail to look up. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID to filter by. \n" +
					"  - example : '36f4d54338b74181888a9541d6c1158f' \n",
			},
			"guardrail": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Guardrail details. \n" +
					"  - example : '{id: 2c8a138f8d78e1fc29a449dbfa8681, name: foo-guardrail, status: DISABLED, ...}' \n",
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
						Description: "Guardrail application status. \n" +
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
	}
}

func (d *guardrailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.GuardrailDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	guardrailId := state.GuardrailId.ValueString()
	if guardrailId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Guardrail",
			"Guardrail ID is empty",
		)
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()

	data, err := d.client.GetGuardrail(ctx, guardrailId, landingZoneId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Guardrail",
			"Could not read CloudControl Guardrail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	guardrailValue, guardrailDiags := d.buildGuardrailValue(ctx, data)
	resp.Diagnostics.Append(guardrailDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Guardrail = guardrailValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *guardrailDataSource) buildGuardrailValue(ctx context.Context, g *sdkcloudcontrol.Guardrail) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var description types.String
	if g.Description.IsSet() && g.Description.Get() != nil && *g.Description.Get() != "" {
		description = types.StringValue(strings.TrimSpace(*g.Description.Get()))
	} else {
		description = types.StringNull()
	}

	var status types.String
	if g.Status != nil && *g.Status != "" {
		status = types.StringValue(*g.Status)
	} else {
		status = types.StringNull()
	}

	var srn types.String
	if g.Srn.IsSet() && g.Srn.Get() != nil && *g.Srn.Get() != "" {
		srn = types.StringValue(*g.Srn.Get())
	} else {
		srn = types.StringNull()
	}

	var bindingOuItems []cloudcontrol.GuardrailBindingOuItem
	if g.BindingOus.IsSet() && g.BindingOus.Get() != nil {
		ous := g.BindingOus.Get()
		if ous.ArrayOfGuardrailBindingOu != nil {
			for _, ou := range *ous.ArrayOfGuardrailBindingOu {
				bindingOuItems = append(bindingOuItems, cloudcontrol.GuardrailBindingOuItem{
					Id:   types.StringValue(ou.Id),
					Name: types.StringValue(ou.Name),
				})
			}
		}
	}

	bindingOusValue, bindingOusDiags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingOuItem{}.AttributeTypes(ctx)},
		bindingOuItems,
	)
	diags.Append(bindingOusDiags...)

	guardrailValue, objDiags := types.ObjectValue(
		cloudcontrol.GuardrailValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":           types.StringValue(g.Id),
			"name":         types.StringValue(g.Name),
			"description":  description,
			"guidance":     types.StringValue(g.Guidance),
			"service_name": types.StringValue(g.ServiceName),
			"status":       status,
			"type":         types.StringValue(g.Type),
			"created_at":   types.StringValue(g.CreatedAt.Format(time.RFC3339)),
			"created_by":   types.StringValue(g.CreatedBy),
			"modified_at":  types.StringValue(g.ModifiedAt.Format(time.RFC3339)),
			"modified_by":  types.StringValue(g.ModifiedBy),
			"srn":          srn,
			"binding_ous":  bindingOusValue,
		},
	)
	diags.Append(objDiags...)

	return guardrailValue, diags
}
