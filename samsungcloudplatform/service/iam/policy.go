package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamPolicyResource{}
	_ resource.ResourceWithConfigure   = &iamPolicyResource{}
	_ resource.ResourceWithImportState = &iamPolicyResource{}
)

// NewIamPolicyResource is a helper function to simplify the provider implementation.
func NewIamPolicyResource() resource.Resource {
	return &iamPolicyResource{}
}

// iamPolicyResource is the data source implementation.
type iamPolicyResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy"
}

func (r *iamPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the policy.\n" +
					"  - example : 'pol-1234567890abcdef'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_name": schema.StringAttribute{
				Optional: true,
				Description: "Name of the policy.\n" +
					"  - example : 'MyPolicy'\n" +
					"  - maxLength: 128\n" +
					"  - minLength: 3",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "Human-readable description of the policy.\n" +
					"  - example : 'My policy description'",
			},
			"tags": tag.ResourceSchema(),
			"policy_version": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Policy version to create or update.\n" +
					"  - example : '{policy_document: {statement: [{action: [iam:CreateRole], effect: Allow, resource: [*], ...}]}}'",
				Attributes: map[string]schema.Attribute{
					"policy_document": schema.SingleNestedAttribute{
						Optional: true,
						Description: "The policy document containing the permission definitions.\n" +
							"  - example : '{statement: [{action: [iam:CreateRole], effect: Allow, resource: [*], ...}]}'",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								Optional: true,
								Description: "List of policy statements defining the permissions.\n" +
									"  - example : '[{action: [iam:CreateRole], effect: Allow, resource: [*], ...}]'",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.ListAttribute{
											Optional: true,
											Description: "List of actions allowed by this statement (e.g., iam:CreateRole, iam:ListUsers).\n" +
												"  - example : ['iam:CreateRole']",
											ElementType: types.StringType,
										},
										"not_action": schema.ListAttribute{
											Optional: true,
											Description: "List of actions that are explicitly excluded from this statement.\n" +
												"  - example : ['iam:DeleteRole']",
											ElementType: types.StringType,
										},
										"effect": schema.StringAttribute{
											Optional: true,
											Description: "Effect of the statement - either Allow or Deny.\n" +
												"  - example : 'Allow'",
										},
										"resource": schema.ListAttribute{
											Optional: true,
											Description: "List of resources (ARNs or wildcards) that the statement applies to.\n" +
												"  - example : ['*']",
											ElementType: types.StringType,
										},
										"sid": schema.StringAttribute{
											Optional: true,
											Description: "Statement ID (SID) - unique identifier for this policy statement.\n" +
												"  - example : 'Stmt1'",
										},
										"condition": schema.MapAttribute{
											ElementType: types.MapType{
												ElemType: types.ListType{
													ElemType: types.StringType,
												},
											},
											Optional: true,
											Description: "Conditions that must be met for the policy statement to take effect.\n" +
												"  - example : {\"StringEquals\": {\"scp:PrincipalTag/department\": [\"finance\"]}}",
										},
									},
								},
							},
							"version": schema.StringAttribute{
								Optional: true,
								Description: "Policy document version.\n" +
									"  - example : '2024-07-01'",
							},
						},
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Description: "Detailed information about the policy.\n" +
					"  - example : '{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, creator_email: user@example.com, ...}'",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed: true,
						Description: "Account ID that owns the policy.\n" +
							"  - example : '123456789012'",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the policy was created.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the policy.\n" +
							"  - example : 'user@example.com'",
					},
					"creator_email": schema.StringAttribute{
						Computed: true,
						Description: "Email of the policy creator.\n" +
							"  - example : 'user@example.com'",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the policy creator.\n" +
							"  - example : 'John Doe'",
					},
					"default_version_id": schema.StringAttribute{
						Computed: true,
						Description: "Default version ID of the policy.\n" +
							"  - example : 'pol-1234567890abcdef'",
					},
					"description": schema.StringAttribute{
						Computed: true,
						Description: "Human-readable description of the policy.\n" +
							"  - example : 'My policy description'",
					},
					"domain_name": schema.StringAttribute{
						Computed: true,
						Description: "Domain name associated with the policy.\n" +
							"  - example : 'scp'",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the policy.\n" +
							"  - example : 'pol-1234567890abcdef'",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the policy was last modified.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who last modified the policy.\n" +
							"  - example : 'user@example.com'",
					},
					"modifier_email": schema.StringAttribute{
						Computed: true,
						Description: "Email of the user who last modified the policy.\n" +
							"  - example : 'user@example.com'",
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the user who last modified the policy.\n" +
							"  - example : 'John Doe'",
					},
					"policy_category": schema.StringAttribute{
						Computed: true,
						Description: "Category of the policy.\n" +
							"  - example : 'IDENTITY_BASED' | 'RESOURCE_BASED'",
					},
					"policy_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the policy.\n" +
							"  - example : 'PolicyName'",
					},
					"policy_type": schema.StringAttribute{
						Computed: true,
						Description: "Type of policy.\n" +
							"  - example : 'USER_DEFINED' | 'SYSTEM_MANAGED'",
					},
					"policy_versions": schema.ListNestedAttribute{
						Optional: true,
						Description: "List of versions associated with the policy.\n" +
							"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: pol-1234567890abcdef, modified_at: 2024-05-17T00:23:17Z, ...}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the policy version was created.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
									MarkdownDescription: "Timestamp when the policy version was created.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
								},
								"created_by": schema.StringAttribute{
									Computed: true,
									Description: "User who created the policy version.\n" +
										"  - example : 'user@example.com'",
									MarkdownDescription: "User who created the policy version.\n" +
										"  - example : 'user@example.com'",
								},
								"id": schema.StringAttribute{
									Computed: true,
									Description: "Unique identifier of the policy version.\n" +
										"  - example : 'pol-1234567890abcdef'",
									MarkdownDescription: "Unique identifier of the policy version.\n" +
										"  - example : 'pol-1234567890abcdef'",
								},
								"modified_at": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the policy version was last modified.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
									MarkdownDescription: "Timestamp when the policy version was last modified.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
								},
								"modified_by": schema.StringAttribute{
									Computed: true,
									Description: "User who last modified the policy version.\n" +
										"  - example : 'user@example.com'",
									MarkdownDescription: "User who last modified the policy version.\n" +
										"  - example : 'user@example.com'",
								},
								"policy_document": schema.SingleNestedAttribute{
									Computed: true,
									Description: "The policy document containing the permission definitions.\n" +
										"  - example : '{statement: [{action: [iam:CreateRole], effect: Allow, resource: [*], ...}], version: 2024-07-01}'",
									Attributes: map[string]schema.Attribute{
										"statement": schema.ListNestedAttribute{
											Computed: true,
											Description: "List of policy statements defining the permissions.\n" +
												"  - example : '[{action: [iam:CreateRole], effect: Allow, resource: [*], sid: Sid1, ...}]'",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.ListAttribute{
														Optional: true,
														Description: "Actions permitted by the policy statement (e.g., iam:CreateRole, iam:ListUsers).\n" +
															"  - example : ['iam:CreateRole']",
														ElementType: types.StringType,
													},
													"not_action": schema.ListAttribute{
														Optional: true,
														Description: "Actions explicitly excluded from the policy statement.\n" +
															"  - example : ['iam:DeleteRole']",
														ElementType: types.StringType,
													},
													"effect": schema.StringAttribute{
														Computed: true,
														Description: "Effect of the policy statement - either Allow or Deny.\n" +
															"  - example : 'Allow'",
													},
													"resource": schema.ListAttribute{
														Optional: true,
														Description: "Resources that the policy statement applies to (ARNs or wildcards).\n" +
															"  - example : ['*']",
														ElementType: types.StringType,
													},
													"sid": schema.StringAttribute{
														Computed: true,
														Description: "Statement ID - unique identifier for this policy statement.\n" +
															"  - example : 'Sid1'",
													},
													"condition": schema.MapAttribute{
														ElementType: types.MapType{
															ElemType: types.ListType{
																ElemType: types.StringType,
															},
														},
														Optional: true,
														Description: "Conditions that must be met for the policy statement to take effect.\n" +
															"  - example : {'StringEquals': {'scp:PrincipalTag/department': ['finance']}}",
													},
													"principal": schema.SingleNestedAttribute{
														Optional: true,
														Description: "Principal - The entity (user, service, or account) that the policy statement applies to.\n" +
															"  - example : '{principal_string: 123456789012, principal_map: {AWS: [arn:aws:iam::123456789012:root]}}'",
														Attributes: map[string]schema.Attribute{
															"principal_string": schema.StringAttribute{
																Optional:            true,
																Description:         "Principal as a string value (e.g., AWS account ID or IAM user ARN).\n  - example : '123456789012'",
																MarkdownDescription: "Principal as a string value (e.g., AWS account ID or IAM user ARN).\n  - example : '123456789012'",
															},
															"principal_map": schema.MapAttribute{
																Optional: true,
																ElementType: types.ListType{
																	ElemType: types.StringType,
																},
																Description:         "Principal as a map - supports multiple principal types (e.g., AWS, Federated, etc.).\n  - example : {'AWS': ['arn:aws:iam::123456789012:root']}",
																MarkdownDescription: "Principal as a map - supports multiple principal types (e.g., AWS, Federated, etc.).\n  - example : {'AWS': ['arn:aws:iam::123456789012:root']}",
															},
														},
													},
												},
											},
										},
										"version": schema.StringAttribute{
											Computed: true,
											Description: "Policy Version\n" +
												"  - example : '2024-07-01'",
											MarkdownDescription: "Policy Version\n" +
												"  - example : '2024-07-01'",
										},
									},
								},

								"policy_id": schema.StringAttribute{
									Computed: true,
									Description: "Unique identifier of the policy.\n" +
										"  - example : 'pol-1234567890abcdef'",
									MarkdownDescription: "Unique identifier of the policy.\n" +
										"  - example : 'pol-1234567890abcdef'",
								},
								"policy_version_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the policy version.\n" +
										"  - example : 'POLICY_VERSION_1'",
									MarkdownDescription: "Name of the policy version.\n" +
										"  - example : 'POLICY_VERSION_1'",
								},
							},
						},
					},
					"resource_type": schema.StringAttribute{
						Computed: true,
						Description: "Type of resource the policy applies to.\n" +
							"  - example : 'policy'",
					},
					"service_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the service.\n" +
							"  - example : 'Identity Access Management'",
					},
					"service_type": schema.StringAttribute{
						Computed: true,
						Description: "Type of service.\n" +
							"  - example : 'iam'",
					},
					"srn": schema.StringAttribute{
						Computed: true,
						Description: "Service Resource Name (SRN) - Unique identifier for the policy in the SCP system.\n" +
							"  - example : 'srn:e:::::iam:policy/policy-12345678'",
						MarkdownDescription: "Service Resource Name (SRN) - Unique identifier for the policy in the SCP system.\n" +
							"  - example : 'srn:e:::::iam:policy/policy-12345678'",
					},
					"state": schema.StringAttribute{
						Computed: true,
						Description: "State of the policy.\n" +
							"  - example : 'ACTIVE'",
						MarkdownDescription: "State of the policy.\n" +
							"  - example : 'ACTIVE'",
					},
				},
			},
		},
	}
}

func (r *iamPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Iam
	r.clients = inst.Client
}

func (r *iamPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *iamPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.PolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreatePolicy(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating policy",
			"Could not create policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Polling is not needed because the CreatePolicy API returns a complete policy object
	// with all details (including State, PolicyVersions, etc.) synchronously.
	// The policy is fully created when the API returns, no pending state exists.

	var policyVersions []iam.PolicyVersion
	//policy versions
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

	plan.Id = types.StringPointerValue(data.Id)

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

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	plan.Policy = policyObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.PolicyResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Show Policy",
			err.Error(),
		)
		return
	}

	var policyVersions []iam.PolicyVersion
	//policy versions
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

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	state.Policy = policyObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state iam.PolicyResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePolicy(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Policy",
			"Could not update Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Policy",
			"Could not read policy ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var policyVersions []iam.PolicyVersion
	//policy versions
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

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	state.Policy = policyObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.PolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicy(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam policy",
			"Could not delete Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertCondition(ctx context.Context, rawCondition map[string]map[string][]interface{}) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionMapType := types.MapType{
		ElemType: types.ListType{
			ElemType: types.StringType,
		},
	}

	if rawCondition == nil {
		emptyOuterMap := map[string]attr.Value{}

		emptyConditionMap, emptyConditionDiags := types.MapValueFrom(ctx, conditionMapType, emptyOuterMap)
		diags.Append(emptyConditionDiags...)
		if emptyConditionDiags.HasError() {
			emptyConditionMap = types.MapUnknown(conditionMapType)
		}

		return emptyConditionMap, diags
	}

	outerMap := map[string]attr.Value{}
	for condType, innerMap := range rawCondition {
		if innerMap == nil {
			outerMap[condType] = types.MapNull(types.MapType{
				ElemType: types.ListType{ElemType: types.StringType},
			})
			continue
		}

		inner := map[string]attr.Value{}
		for key, values := range innerMap {
			stringValues := make([]attr.Value, len(values))
			for i, v := range values {
				if s, ok := v.(string); ok {
					stringValues[i] = types.StringValue(s)
				} else {
					stringValues[i] = types.StringNull()
					diags.AddAttributeWarning(
						path.Root("condition"),
						"Invalid Condition Value",
						"Value is not a string. Using null instead.",
					)
				}
			}

			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, stringValues)
			diags.Append(listDiags...)
			if listDiags.HasError() {
				listValue = types.ListNull(types.StringType)
			}

			inner[key] = listValue
		}

		innerMapType := types.ListType{
			ElemType: types.StringType,
		}

		mapValue, mapDiags := types.MapValueFrom(ctx, innerMapType, inner)
		if mapDiags.HasError() {
			mapValue = types.MapUnknown(innerMapType)
		}
		outerMap[condType] = mapValue
	}

	conditionMap, condDiags := types.MapValueFrom(ctx, conditionMapType, outerMap)
	diags.Append(condDiags...)
	if condDiags.HasError() {
		conditionMap = types.MapUnknown(conditionMapType)
	}

	return conditionMap, diags
}
