package iam

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamGroupPolicyBindingDataSources{}
	_ datasource.DataSourceWithConfigure = &iamGroupPolicyBindingDataSources{}
)

func NewIamGroupPolicyBindingDataSources() datasource.DataSource {
	return &iamGroupPolicyBindingDataSources{}
}

type iamGroupPolicyBindingDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamGroupPolicyBindingDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group_policy_bindings"
}

func (d *iamGroupPolicyBindingDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamGroupPolicyBindingDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Group Policy Bindings",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Group ID",
				MarkdownDescription: "Group ID",
			},
			"size": schema.Int32Attribute{
				Optional:            true,
				Description:         "Size (between 1 and 10000)",
				MarkdownDescription: "Size (between 1 and 10000)",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional:            true,
				Description:         "Page (between 0 and 10000)",
				MarkdownDescription: "Page (between 0 and 10000)",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort",
				MarkdownDescription: "Sort",
			},
			"policy_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy ID",
				MarkdownDescription: "Policy ID",
			},
			"policy_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy Name",
				MarkdownDescription: "Policy Name",
			},
			"policy_type": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy Type",
				MarkdownDescription: "Policy Type",
			},
			"group_policy_bindings": schema.ListNestedAttribute{
				Optional:            true,
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

func (d *iamGroupPolicyBindingDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.GroupPolicyBindingsDataResource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGroupPolicyBindings(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Policies",
			err.Error(),
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	state.GroupPolicyBindings = policies

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
