package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/organization/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &serviceControlPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceControlPolicyDataSource{}
)

func NewServiceControlPolicyDataSource() datasource.DataSource {
	return &serviceControlPolicyDataSource{}
}

type serviceControlPolicyDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *serviceControlPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_service_control_policy"
}

func (d *serviceControlPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Organization
	d.clients = inst.Client
}

func (d *serviceControlPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization Service Control Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Service Control Policy ID. \n" +
					"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
				Required: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Service Control Policy Name. \n" +
					"  - example : 'MyPolicy' \n",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Policy Description. \n" +
					"  - example : 'This is an example policy.' \n",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "Service Control Policy Type. \n" +
					"  - example : 'MANAGED' \n",
			},
			"document": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Service Control Policy Document. \n" +
					"  - example : '{version: 2024-07-01, statement: [...]}' \n",
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of actions permitted or denied by the policy statement. \n" +
										"  - example : ['s3:PutObject', 's3:GetObject'] \n",
								},
								"condition": schema.MapAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "Policy Condition. \n" +
										"  - example : '{StringEquals: {aws:RequestedRegion: us-east-1}}' \n",
								},
								"effect": schema.StringAttribute{
									Computed: true,
									Description: "Policy Effect. \n" +
										"  - example : 'Allow' \n" +
										"  - allowed_values : ['Allow', 'Deny'] \n",
								},
								"not_action": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "Policy Exclusion Action. \n" +
										"  - example : ['s3:DeleteObject'] \n",
								},
								"principal": schema.StringAttribute{
									Computed: true,
									Description: "Principal entity to which the policy statement applies. \n" +
										"  - example : '*' \n",
								},
								"resource": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of resources to which the policy statement applies. \n" +
										"  - example : ['*'] \n",
								},
								"sid": schema.StringAttribute{
									Computed: true,
									Description: "Syntax ID. \n" +
										"  - example : 'statement1' \n",
								},
							},
						},
						Computed: true,
						Description: "Policy Syntax. \n" +
							"  - example : '[{effect: Allow, action: [*], resource: [*], sid: statement1}]' \n",
					},
					"version": schema.StringAttribute{
						Computed: true,
						Description: "Policy Version. \n" +
							"  - example : '2012-10-17' \n",
					},
				},
			},
			"category": schema.StringAttribute{
				Computed: true,
				Description: "Policy Category. \n" +
					"  - example : 'SCP' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the service control policy was created. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
			},
			"created_by": schema.StringAttribute{
				Computed: true,
				Description: "User who created the service control policy. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
			},
			"creator_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the service control policy creator. \n" +
					"  - example : 'John Doe na' \n",
			},
			"modified_at": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the service control policy was modified. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
				Description: "User who modified the service control policy. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
			},
			"modifier_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the service control policy modifier. \n" +
					"  - example : 'Alice' \n",
			},
			"source": schema.StringAttribute{
				Computed: true,
				Description: "Policy Creation Subject. \n" +
					"  - example : 'ORGANIZATION' \n",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "service control policy state. \n" +
					"  - example : 'ACTIVE' \n",
			},
			"srn": schema.StringAttribute{
				Computed: true,
				Description: "Samsung Resource Name (SRN) uniquely identifying this resource. \n" +
					"  - example : 'srn:dev2::...:servicecontrol:policy/...' \n",
			},
			"service_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the service to which the policy applies. \n" +
					"  - example : 'Organization' \n",
			},
		},
	}
}

func (d *serviceControlPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.ServiceControlPolicyDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyId.ValueString()
	orgId := state.OrganizationId.ValueString()

	if policyId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Service Control Policy",
			"Policy ID is required",
		)
		return
	}

	data, err := d.client.GetServiceControlPolicy(ctx, policyId, orgId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Service Control Policy",
			err.Error(),
		)
		return
	}

	state.PolicyId = types.StringValue(data.Policy.Id)
	state.OrganizationId = types.StringValue(data.Policy.OrganizationId)
	state.Name = types.StringValue(data.Policy.Name)
	if desc := data.Policy.Description.Get(); desc != nil {
		state.Description = types.StringValue(*desc)
	}
	state.Type = types.StringValue(string(data.Policy.Type))
	state.Category = types.StringValue(data.Policy.Category)
	state.CreatedAt = types.StringValue(data.Policy.CreatedAt.Format("2006-01-02T15:04:05.000Z"))
	state.CreatedBy = types.StringValue(data.Policy.CreatedBy)
	state.CreatorName = types.StringValue(data.Policy.GetCreatorName())
	state.ModifiedAt = types.StringValue(data.Policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z"))
	state.ModifiedBy = types.StringValue(data.Policy.ModifiedBy)
	state.ModifierName = types.StringValue(data.Policy.GetModifierName())
	state.Source = types.StringValue(data.Policy.Source)
	state.State = types.StringValue(data.Policy.State)
	state.Srn = types.StringValue(data.Policy.Srn)
	state.ServiceName = types.StringValue(data.Policy.ServiceName)

	documentValue, docDiags := d.buildDocumentValue(ctx, &data.Policy.Document)
	resp.Diagnostics.Append(docDiags...)
	state.Document = documentValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serviceControlPolicyDataSource) buildDocumentValue(ctx context.Context, doc any) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx))
	statementsList := types.ListNull(types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)})
	versionVal := types.StringNull()

	if sdcdoc, ok := doc.(*sdkorganization.ServiceControlPolicyDocument); ok && sdcdoc != nil {
		if sdcdoc.Statement != nil {
			var statementsListDiags diag.Diagnostics
			statementsList, statementsListDiags = d.buildStatementsList(ctx, sdcdoc.Statement)
			diags.Append(statementsListDiags...)
			if diags.HasError() {
				return nullObj, diags
			}
		}
		if sdcdoc.Version != nil {
			versionVal = types.StringValue(*sdcdoc.Version)
		}
	}

	documentValue, valueDiags := types.ObjectValue(
		organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"statement": statementsList,
			"version":   versionVal,
		},
	)
	diags.Append(valueDiags...)
	return documentValue, diags
}

func (d *serviceControlPolicyDataSource) buildStatementValue(ctx context.Context, stmt sdkorganization.ServiceControlPolicyStatement) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	effectStr := ""
	if stmt.Effect != nil {
		effectStr = *stmt.Effect
	}
	sidStr := ""
	if stmt.Sid != nil {
		sidStr = *stmt.Sid
	}
	conditionValue, conditionDiags := d.buildConditionValue(ctx, stmt.Condition)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}
	statementValue, statementDiags := types.ObjectValue(
		organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"action":     d.buildStringList(ctx, stmt.Action),
			"condition":  conditionValue,
			"effect":     types.StringValue(effectStr),
			"not_action": d.buildStringList(ctx, stmt.NotAction),
			"principal":  d.buildPrincipalValue(ctx, stmt.Principal),
			"resource":   d.buildStringList(ctx, stmt.Resource),
			"sid":        types.StringValue(sidStr),
		},
	)
	diags.Append(statementDiags...)
	return statementValue, diags
}

func (d *serviceControlPolicyDataSource) buildStatementsList(ctx context.Context, stmts []sdkorganization.ServiceControlPolicyStatement) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullList := types.ListNull(types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)})
	statements := make([]attr.Value, 0, len(stmts))
	for _, stmt := range stmts {
		statementValue, stmtDiags := d.buildStatementValue(ctx, stmt)
		diags.Append(stmtDiags...)
		if diags.HasError() {
			return nullList, diags
		}
		statements = append(statements, statementValue)
	}
	statementsList, listDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)
	diags.Append(listDiags...)
	return statementsList, diags
}

func (d *serviceControlPolicyDataSource) buildStringList(ctx context.Context, values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	list := make([]attr.Value, 0, len(values))
	for _, v := range values {
		list = append(list, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, list)
}

func (d *serviceControlPolicyDataSource) buildPrincipalValue(ctx context.Context, principal any) types.String {
	if principal == nil {
		return types.StringNull()
	}
	if p, ok := principal.(string); ok {
		return types.StringValue(p)
	}
	if p, ok := principal.(map[string]any); ok {
		if len(p) == 0 {
			return types.StringNull()
		}
		if val, exists := p["scp"]; exists {
			if arr, ok := val.([]string); ok && len(arr) > 0 {
				return types.StringValue(arr[0])
			}
		}
	}
	return types.StringNull()
}

func (d *serviceControlPolicyDataSource) buildConditionValue(ctx context.Context, condition any) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	if condition == nil {
		return types.MapNull(types.StringType), diags
	}
	if c, ok := condition.(map[string]map[string][]string); ok {
		if len(c) == 0 {
			return types.MapNull(types.StringType), diags
		}
		result := make(map[string]attr.Value)
		for k, v := range c {
			result[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		mapVal, mapDiags := types.MapValue(types.StringType, result)
		diags.Append(mapDiags...)
		return mapVal, diags
	}
	return types.MapNull(types.StringType), diags
}
