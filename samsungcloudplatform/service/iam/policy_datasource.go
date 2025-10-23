package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &iamPolicyDataSource{}
)

// NewIamPolicyDataSource is a helper function to simplify the provider implementation.
func NewIamPolicyDataSource() datasource.DataSource {
	return &iamPolicyDataSource{}
}

type iamPolicyDataSource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy"
}

// Configure adds the provider configured client to the data source.
func (d *iamPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.Iam
	d.clients = inst.Client
}

func (d *iamPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy ID",
				MarkdownDescription: "Policy ID",
			},
			"policy": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "A detail of Policy",
				MarkdownDescription: "A detail of Policy",
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID",
						MarkdownDescription: "Account ID",
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
					"creator_email": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Email",
						MarkdownDescription: "Creator Email",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Name",
						MarkdownDescription: "Creator Name",
					},
					"default_version_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Default Version ID",
						MarkdownDescription: "Default Version ID",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Description",
						MarkdownDescription: "Description",
					},
					"domain_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Domain Name",
						MarkdownDescription: "Domain Name",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
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
					"modifier_email": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Email",
						MarkdownDescription: "Modifier Email",
					},
					"modifier_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Name",
						MarkdownDescription: "Modifier Name",
					},
					"policy_category": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Category",
						MarkdownDescription: "Policy Category",
					},
					"policy_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Name",
						MarkdownDescription: "Policy Name",
					},
					"policy_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Type",
						MarkdownDescription: "Policy Type",
					},
					"policy_versions": schema.ListNestedAttribute{
						Optional:            true,
						Description:         "Policy Versions",
						MarkdownDescription: "Policy Versions",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
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
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "ID",
									MarkdownDescription: "ID",
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
								"policy_document": schema.SingleNestedAttribute{
									Computed:            true,
									Description:         "Policy Document",
									MarkdownDescription: "Policy Document",
									Attributes: map[string]schema.Attribute{
										"statement": schema.ListNestedAttribute{
											Computed:            true,
											Description:         "Statement",
											MarkdownDescription: "Statement",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.ListAttribute{
														Optional:            true,
														Description:         "Action",
														MarkdownDescription: "Action",
														ElementType:         types.StringType,
													},
													"not_action": schema.ListAttribute{
														Optional:            true,
														Description:         "Not Action",
														MarkdownDescription: "Not Action",
														ElementType:         types.StringType,
													},
													"effect": schema.StringAttribute{
														Computed:            true,
														Description:         "Effect",
														MarkdownDescription: "Effect",
													},
													"resource": schema.ListAttribute{
														Optional:            true,
														Description:         "Resource",
														MarkdownDescription: "Resource",
														ElementType:         types.StringType,
													},
													"sid": schema.StringAttribute{
														Computed:            true,
														Description:         "SID",
														MarkdownDescription: "SID",
													},
													"condition": schema.MapAttribute{
														ElementType: types.MapType{
															ElemType: types.ListType{
																ElemType: types.StringType,
															},
														},
														Optional: true,
													},
													"principal": schema.SingleNestedAttribute{
														Optional:            true,
														Description:         "Principal",
														MarkdownDescription: "Principal",
														Attributes: map[string]schema.Attribute{
															"principal_string": schema.StringAttribute{
																Optional: true,
															},
															"principal_map": schema.MapAttribute{
																Optional: true,
																ElementType: types.ListType{
																	ElemType: types.StringType,
																},
															},
														},
													},
												},
											},
										},
										"version": schema.StringAttribute{
											Computed:            true,
											Description:         "Policy Version",
											MarkdownDescription: "Policy Version",
										},
									},
								},

								"policy_id": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy ID",
									MarkdownDescription: "Policy ID",
								},
								"policy_version_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy Version Name",
									MarkdownDescription: "Policy Version Name",
								},
							},
						},
					},
					"resource_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Resource Type",
						MarkdownDescription: "Resource Type",
					},
					"service_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Name",
						MarkdownDescription: "Service Name",
					},
					"service_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Type",
						MarkdownDescription: "Service Type",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "SRN",
						MarkdownDescription: "SRN",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "State",
						MarkdownDescription: "State",
					},
				},
			},
		},
	}
}

func (d *iamPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.PolicyDatasourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Policy",
			err.Error(),
		)
		return
	}
	//policy versions
	var policyVersions []iam.PolicyVersion

	for _, policyVersion := range data.PolicyVersions {

		var statements []iam.Statement
		for _, _statement := range policyVersion.PolicyDocument.Statement {

			// resource
			resources := make([]types.String, 0, len(_statement.Resource))
			for _, _resource := range _statement.Resource {
				resources = append(resources, types.StringValue(_resource))
			}

			// action
			actions := make([]types.String, 0, len(_statement.Action))
			for _, _action := range _statement.Action {
				actions = append(actions, types.StringValue(_action))
			}

			// not action
			notActions := make([]types.String, 0, len(_statement.NotAction))
			for _, _notAction := range _statement.NotAction {
				notActions = append(notActions, types.StringValue(_notAction))
			}

			// principal
			principal, _ := convertPrincipal(ctx, _statement.Principal)

			// condition
			condition, _ := convertCondition(ctx, _statement.Condition)

			statement := iam.Statement{
				Sid:       types.StringPointerValue(_statement.Sid),
				Effect:    types.StringValue(_statement.Effect),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: principal,
				Condition: condition,
			}

			statements = append(statements, statement)
		}

		policyDocument := iam.PolicyDocument{
			Version:   types.StringValue(policyVersion.PolicyDocument.Version),
			Statement: statements,
		}

		policyVersionState := iam.PolicyVersion{
			CreatedAt:         types.StringValue(policyVersion.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(policyVersion.CreatedBy),
			Id:                types.StringValue(*policyVersion.Id),
			ModifiedAt:        types.StringValue(policyVersion.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(policyVersion.ModifiedBy),
			PolicyDocument:    policyDocument,
			PolicyId:          types.StringValue(*policyVersion.PolicyId),
			PolicyVersionName: types.StringValue(*policyVersion.PolicyVersionName),
		}
		policyVersions = append(policyVersions, policyVersionState)

	}

	policyState := iam.Policy{
		AccountId:        types.StringPointerValue(data.AccountId.Get()),
		CreatedAt:        types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.CreatedBy),
		CreatorEmail:     types.StringPointerValue(data.CreatorEmail.Get()),
		CreatorName:      types.StringPointerValue(data.CreatorName.Get()),
		DefaultVersionId: types.StringValue(*data.DefaultVersionId),
		Description:      types.StringPointerValue(data.Description.Get()),
		DomainName:       types.StringValue(data.DomainName),
		Id:               types.StringValue(*data.Id),
		ModifiedAt:       types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.ModifiedBy),
		ModifierEmail:    types.StringPointerValue(data.ModifierEmail.Get()),
		ModifierName:     types.StringPointerValue(data.ModifierName.Get()),
		PolicyCategory:   types.StringValue(string(*data.PolicyCategory)),
		PolicyName:       types.StringValue(*data.PolicyName),
		PolicyType:       types.StringValue(string(*data.PolicyType)),
		PolicyVersions:   policyVersions,
		ResourceType:     types.StringPointerValue(data.ResourceType.Get()),
		ServiceName:      types.StringPointerValue(data.ServiceName.Get()),
		ServiceType:      types.StringPointerValue(data.ServiceType.Get()),
		Srn:              types.StringValue(data.Srn),
		State:            types.StringValue(string(*data.State)),
	}

	policyObjectValue, _ := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	state.Policy = policyObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
