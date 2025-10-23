package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &iamGroupDataSources{}
)

// NewIamGroupDataSources is a helper function to simplify the provider implementation.
func NewIamGroupDataSources() datasource.DataSource {
	return &iamGroupDataSources{}
}

type iamGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_groups"
}

// Configure adds the provider configured client to the data source.
func (d *iamGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Groups",
		Attributes: map[string]schema.Attribute{
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
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Group Name",
				MarkdownDescription: "Group Name",
			},
			"groups": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "A list of group",
				MarkdownDescription: "A list of group",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "생성 일시",
							MarkdownDescription: "생성 일시",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "생성자",
							MarkdownDescription: "생성자",
						},
						"creator_email": schema.StringAttribute{
							Computed:            true,
							Description:         "생성자 Email",
							MarkdownDescription: "생성자 Email",
						},
						"creator_name": schema.StringAttribute{
							Computed:            true,
							Description:         "생성자 성, 이름",
							MarkdownDescription: "생성자 성, 이름",
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"domain_name": schema.StringAttribute{
							Computed:            true,
							Description:         "도메인 이름",
							MarkdownDescription: "도메인 이름",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "ID",
							MarkdownDescription: "ID",
						},
						"members": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Members",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "생성 일시",
										MarkdownDescription: "생성 일시",
									},
									"created_by": schema.StringAttribute{
										Computed:            true,
										Description:         "생성자",
										MarkdownDescription: "생성자",
									},
									"creator_created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "생성 일시",
										MarkdownDescription: "생성 일시",
									},
									"creator_email": schema.StringAttribute{
										Computed:            true,
										Description:         "생성자 Email",
										MarkdownDescription: "생성자 Email",
									},
									"creator_last_login_at": schema.StringAttribute{
										Optional:            true,
										Description:         "생성자 마지막 로그인 일시",
										MarkdownDescription: "생성자 마지막 로그인 일시",
									},
									"creator_name": schema.StringAttribute{
										Computed:            true,
										Description:         "생성자 성, 이름",
										MarkdownDescription: "생성자 성, 이름",
									},
									"group_names": schema.ListAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										Description:         "Group names",
										MarkdownDescription: "Group names",
									},
									"user_created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "생성 일시",
										MarkdownDescription: "생성 일시",
									},
									"user_email": schema.StringAttribute{
										Computed:            true,
										Description:         "User Email",
										MarkdownDescription: "User Email",
									},
									"user_id": schema.StringAttribute{
										Computed:            true,
										Description:         "User ID",
										MarkdownDescription: "User ID",
									},
									"user_last_login_at": schema.StringAttribute{
										Optional:            true,
										Description:         "User 마지막 로그인 일시",
										MarkdownDescription: "User 마지막 로그인 일시",
									},
									"user_name": schema.StringAttribute{
										Computed:            true,
										Description:         "User 성, 이름",
										MarkdownDescription: "User 성, 이름",
									},
								},
							},
						},
						"policies": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_id": schema.StringAttribute{
										Computed:            true,
										Optional:            true,
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

						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "수정 일시",
							MarkdownDescription: "수정 일시",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "수정자",
							MarkdownDescription: "수정자",
						},
						"modifier_email": schema.StringAttribute{
							Computed:            true,
							Description:         "수정자 Email",
							MarkdownDescription: "수정자 Email",
						},
						"modifier_name": schema.StringAttribute{
							Computed:            true,
							Description:         "수정자 성, 이름",
							MarkdownDescription: "수정자 성, 이름",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Group 이름",
							MarkdownDescription: "Group 이름",
						},
						"resource_type": schema.StringAttribute{
							Computed: true,
						},
						"service_name": schema.StringAttribute{
							Computed: true,
						},
						"service_type": schema.StringAttribute{
							Computed: true,
						},
						"srn": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Group Type",
							MarkdownDescription: "Group Type",
						},
					},
				},
			},
		},
	}
}

func (d *iamGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.GroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGroups(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Groups",
			err.Error(),
		)
		return
	}

	for _, group := range data.Groups {

		// group members
		members := getGroupMembers(group.Members)

		// policies
		policies, hasError := getPolicies(ctx, group.Policies)
		if hasError {
			return
		}

		groupState := iam.Group{
			CreatedAt:     types.StringValue(group.CreatedAt.Format(time.RFC3339)),
			CreatedBy:     types.StringValue(group.CreatedBy),
			CreatorEmail:  types.StringPointerValue(group.CreatorEmail),
			CreatorName:   types.StringPointerValue(group.CreatorName),
			Description:   types.StringPointerValue(group.Description.Get()),
			DomainName:    types.StringValue(group.DomainName),
			Id:            types.StringValue(group.Id),
			Members:       members,
			Policies:      policies,
			ModifiedAt:    types.StringValue(group.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:    types.StringValue(group.ModifiedBy),
			ModifierEmail: types.StringPointerValue(group.ModifierEmail),
			ModifierName:  types.StringPointerValue(group.ModifierName),
			Name:          types.StringValue(group.Name),
			ResourceType:  types.StringPointerValue(group.ResourceType.Get()),
			ServiceName:   types.StringPointerValue(group.ServiceName.Get()),
			ServiceType:   types.StringPointerValue(group.ServiceType.Get()),
			Srn:           types.StringPointerValue(group.Srn.Get()),
			GroupType:     types.StringValue(group.Type),
		}

		state.Groups = append(state.Groups, groupState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
