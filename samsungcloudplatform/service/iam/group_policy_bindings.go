package iam

import (
	"context"
	"fmt"
	"slices"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &iamGroupPolicyBindingsResource{}
	_ resource.ResourceWithConfigure = &iamGroupPolicyBindingsResource{}
)

func NewIamGroupPolicyBindingsResource() resource.Resource {
	return &iamGroupPolicyBindingsResource{}
}

type iamGroupPolicyBindingsResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (r *iamGroupPolicyBindingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group_policy_bindings"
}

func (r *iamGroupPolicyBindingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamGroupPolicyBindingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Group Policy",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Group ID",
				MarkdownDescription: "Group ID",
			},
			"policy_ids": schema.ListAttribute{
				Optional:            true,
				Description:         "Policy IDs",
				MarkdownDescription: "Policy IDs",
				ElementType:         types.StringType,
			},
			"group_policy_bindings": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "Group Policy Bindings",
				MarkdownDescription: "Group Policy Bindings",
				NestedObject: schema.NestedAttributeObject{
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
		},
	}

}

func (r *iamGroupPolicyBindingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.GroupPolicyBindingsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.client.AddGroupPolicyBindings(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group policies",
			"Could not add group policies, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	groupPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.GroupPolicyBindings = groupPolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupPolicyBindingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.GroupPolicyBindingsResource
	var state iam.GroupPolicyBindingsResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// detach
	err := r.client.RemoveGroupPolicyBindings(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group policy",
			"Could not delete Group Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// attach
	data, err := r.client.AddGroupPolicyBindings(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group policies",
			"Could not add group policies, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	groupPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.GroupPolicyBindings = groupPolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupPolicyBindingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.GroupPolicyBindingsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveGroupPolicyBindings(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group policy",
			"Could not delete Group Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *iamGroupPolicyBindingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.GroupPolicyBindingsResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetGroupPolicyBindings(ctx, state.GroupId.ValueString(), iam.GroupPolicyBindingsDataResource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Policies",
			err.Error(),
		)
		return
	}

	getAllGroupPolicies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	var groupPolicy []iam.Policy
	for _, policy := range getAllGroupPolicies {
		if slices.Contains(state.PolicyIds, policy.Id) {
			groupPolicy = append(groupPolicy, policy)
		}
	}

	groupPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, groupPolicy)
	state.GroupPolicyBindings = groupPolicyObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
