package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
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
			"policy_id": schema.StringAttribute{
				Description:         "Policy ID",
				MarkdownDescription: "Policy ID",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Service Control Policy Name",
				MarkdownDescription: "Service Control Policy Name",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "Description",
				MarkdownDescription: "Description",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy Type (SYSTEM_MANAGED or USER_DEFINED)",
				MarkdownDescription: "Policy Type (SYSTEM_MANAGED or USER_DEFINED)",
			},
			"document": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Service Control Policy Document",
				MarkdownDescription: "Service Control Policy Document",
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									Description:         "Action",
									MarkdownDescription: "Action",
								},
								"condition": schema.MapAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									Description:         "Condition",
									MarkdownDescription: "Condition",
								},
								"effect": schema.StringAttribute{
									Computed:            true,
									Description:         "Effect (Allow/Deny)",
									MarkdownDescription: "Effect (Allow/Deny)",
								},
								"not_action": schema.ListAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									Description:         "Not Action",
									MarkdownDescription: "Not Action",
								},
								"principal": schema.StringAttribute{
									Computed:            true,
									Description:         "Principal (e.g., \"*\")",
									MarkdownDescription: "Principal (e.g., \"*\")",
								},
								"resource": schema.ListAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									Description:         "Resource",
									MarkdownDescription: "Resource",
								},
								"sid": schema.StringAttribute{
									Computed:            true,
									Description:         "Statement ID",
									MarkdownDescription: "Statement ID",
								},
							},
						},
						Computed:            true,
						Description:         "Policy Statement",
						MarkdownDescription: "Policy Statement",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Version",
						MarkdownDescription: "Policy Version",
					},
				},
			},
			"category": schema.StringAttribute{
				Computed:            true,
				Description:         "Category",
				MarkdownDescription: "Category",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At",
				MarkdownDescription: "Created At",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created By",
				MarkdownDescription: "Created By",
			},
			"creator_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Creator Name",
				MarkdownDescription: "Creator Name",
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At",
				MarkdownDescription: "Modified At",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified By",
				MarkdownDescription: "Modified By",
			},
			"modifier_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Modifier Name",
				MarkdownDescription: "Modifier Name",
			},
			"source": schema.StringAttribute{
				Computed:            true,
				Description:         "Source",
				MarkdownDescription: "Source",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "State",
				MarkdownDescription: "State",
			},
			"srn": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy SRN",
				MarkdownDescription: "Policy SRN",
			},
			"service_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Service Name",
				MarkdownDescription: "Service Name",
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

	statementsList := types.ListNull(types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)})
	versionVal := types.StringNull()

	if doc != nil {
		if sdcdoc, ok := doc.(*sdkorganization.ServiceControlPolicyDocument); ok && sdcdoc != nil {
			if sdcdoc.Statement != nil {
				statements := make([]attr.Value, 0, len(sdcdoc.Statement))
				for _, stmt := range sdcdoc.Statement {
					var effectStr string
					if stmt.Effect != nil {
						effectStr = *stmt.Effect
					}
					var sidStr string
					if stmt.Sid != nil {
						sidStr = *stmt.Sid
					}

					principalValue := d.buildPrincipalValue(ctx, stmt.Principal)
					conditionValue := d.buildConditionValue(ctx, stmt.Condition)

					statementValue, _ := types.ObjectValue(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
						"action":     d.buildStringList(ctx, stmt.Action),
						"condition":  conditionValue,
						"effect":     types.StringValue(effectStr),
						"not_action": d.buildStringList(ctx, stmt.NotAction),
						"principal":  principalValue,
						"resource":   d.buildStringList(ctx, stmt.Resource),
						"sid":        types.StringValue(sidStr),
					})
					statements = append(statements, statementValue)
				}
				statementsList = types.ListValueMust(
					types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
					statements,
				)
			}
			if sdcdoc.Version != nil {
				versionVal = types.StringValue(*sdcdoc.Version)
			}
		}
	}

	documentValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   versionVal,
	})
	diags.Append(valueDiags...)

	return documentValue, diags
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

func (d *serviceControlPolicyDataSource) buildConditionValue(ctx context.Context, condition any) types.Map {
	if condition == nil {
		return types.MapNull(types.StringType)
	}
	if c, ok := condition.(map[string]map[string][]string); ok {
		if len(c) == 0 {
			return types.MapNull(types.StringType)
		}
		result := make(map[string]attr.Value)
		for k, v := range c {
			result[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		return types.MapValueMust(types.StringType, result)
	}
	return types.MapNull(types.StringType)
}
