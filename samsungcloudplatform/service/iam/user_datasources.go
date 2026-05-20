package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamUserDataSources{}
	_ datasource.DataSourceWithConfigure = &iamUserDataSources{}
)

func NewIamUserDataSources() datasource.DataSource { return &iamUserDataSources{} }

type iamUserDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamUserDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_users"
}

func (d *iamUserDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamUserDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Users.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page (between 0 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Email"): schema.StringAttribute{
				Description: "Email",
				Optional:    true,
			},
			common.ToSnakeCase("UserName"): schema.StringAttribute{
				Description: "User Name",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "User Type",
				Optional:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID",
				Optional:    true,
			},
			common.ToSnakeCase("Users"): schema.ListNestedAttribute{
				Description: "A list of user.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Account ID",
							MarkdownDescription: "Account ID",
						},
						"company_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Company Name",
							MarkdownDescription: "Company Name",
						},
						"console_url": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Console URL",
							MarkdownDescription: "Console URL",
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
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description",
							MarkdownDescription: "Description",
						},
						"dst_offset": schema.StringAttribute{
							Computed:            true,
							Description:         "Dst Offset",
							MarkdownDescription: "Dst Offset",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							Description:         "Email",
							MarkdownDescription: "Email",
						},
						"email_authenticated": schema.BoolAttribute{
							Computed:            true,
							Description:         "Email Authenticated",
							MarkdownDescription: "Email Authenticated",
						},
						"first_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "First Name",
							MarkdownDescription: "First Name",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "ID",
							MarkdownDescription: "ID",
						},
						"last_login_at": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Last Login At",
							MarkdownDescription: "Last Login At",
						},
						"last_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Last Name",
							MarkdownDescription: "Last Name",
						},
						"last_password_update_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Last Password Update At",
							MarkdownDescription: "Last Password Update At",
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
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name",
							MarkdownDescription: "Name",
						},
						"password": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Password",
							MarkdownDescription: "Password",
						},
						"password_reuse_count": schema.Int64Attribute{
							Computed:            true,
							Description:         "Password Reuse Count",
							MarkdownDescription: "Password Reuse Count",
						},
						"phone_authenticated": schema.BoolAttribute{
							Computed:            true,
							Description:         "Phone Authenticated",
							MarkdownDescription: "Phone Authenticated",
						},
						"policies": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "Policies",
							MarkdownDescription: "Policies",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "ID",
										MarkdownDescription: "ID",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										Description:         "Name",
										MarkdownDescription: "Name",
									},
								},
							},
						},
						"timezone": schema.StringAttribute{
							Computed:            true,
							Description:         "Timezone",
							MarkdownDescription: "Timezone",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type",
							MarkdownDescription: "Type",
						},
						"tz_id": schema.StringAttribute{
							Computed:            true,
							Description:         "TZ ID",
							MarkdownDescription: "TZ ID",
						},
						"user_name": schema.StringAttribute{
							Computed:            true,
							Description:         "User Name",
							MarkdownDescription: "User Name",
						},
						"utc_offset": schema.StringAttribute{
							Computed:            true,
							Description:         "UTC Offset",
							MarkdownDescription: "UTC Offset",
						},
						"access_keys": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Access Keys",
							MarkdownDescription: "Access Keys",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"access_key": schema.StringAttribute{
										Description:         "Access Key",
										MarkdownDescription: "Access Key",
										Computed:            true,
									},
									"created_at": schema.StringAttribute{
										Description:         "Created At",
										MarkdownDescription: "Created At",
										Computed:            true,
									},
									"expiration_timestamp": schema.StringAttribute{
										Description:         "Expiration Timestmap",
										MarkdownDescription: "Expiration Timestmap",
										Computed:            true,
									},
									"id": schema.StringAttribute{
										Description:         "ID",
										MarkdownDescription: "ID",
										Computed:            true,
									},
									"is_enabled": schema.BoolAttribute{
										Description:         "Is Enabled",
										MarkdownDescription: "Is Enabled",
										Computed:            true,
									},
								},
							},
						},
						"groups": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Groups",
							MarkdownDescription: "Groups",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "Group ID",
										MarkdownDescription: "Group ID",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										Description:         "Group Name",
										MarkdownDescription: "Group Name",
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

func (d *iamUserDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.UserDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetUsers(ctx, state.AccountId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Users",
			err.Error(),
		)
		return
	}

	// users
	users, hasError := getUsers(ctx, data.Users)
	if hasError {
		return
	}

	state.Users = users

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getUsers(ctx context.Context, _users []scpsdkiam.IAMUserWithoutPolicyResponseV1Dot4) ([]iam.UserWithoutPolicyDetail, bool) {
	var users []iam.UserWithoutPolicyDetail
	for _, user := range _users {

		// policies(basic)
		var policies []iam.PolicyBasic
		for _, policy := range user.Policies {
			policies = append(policies, iam.PolicyBasic{
				Id:   types.StringValue(policy.Id),
				Name: types.StringValue(policy.GetName()),
			})
		}

		// user nil check
		userCompanyName := user.CompanyName.Get()
		if userCompanyName == nil {
			emptyStr := ""
			userCompanyName = &emptyStr
		}

		userFirstName := user.FirstName.Get()
		if userFirstName == nil {
			emptyStr := ""
			userFirstName = &emptyStr
		}

		userLastName := user.LastName.Get()
		if userLastName == nil {
			emptyStr := ""
			userLastName = &emptyStr
		}

		userLastLoginAt := user.LastLoginAt.Get()
		if userLastLoginAt == nil {
			emptyTime := time.Time{}
			userLastLoginAt = &emptyTime
		}
		userLastPasswordUpdateAt := user.LastPasswordUpdateAt.Get()
		if userLastPasswordUpdateAt == nil {
			emptyTime := time.Time{}
			userLastPasswordUpdateAt = &emptyTime
		}

		// mapped access key info
		var accessKeyInfos []iam.AccessKeyV1Dot4
		for _, accessKeyInfo := range user.AccessKeys {
			accessKeyInfos = append(accessKeyInfos, iam.AccessKeyV1Dot4{
				AccessKey:           types.StringValue(accessKeyInfo.AccessKey),
				CreatedAt:           types.StringValue(accessKeyInfo.CreatedAt.Format(time.RFC3339)),
				ExpirationTimestamp: types.StringValue(accessKeyInfo.ExpirationTimestamp.Format(time.RFC3339)),
				Id:                  types.StringValue(accessKeyInfo.Id),
				IsEnabled:           types.BoolValue(accessKeyInfo.IsEnabled),
			})
		}

		// mapped group info
		var groupInfos []iam.GroupInfo
		for _, groupInfo := range user.Groups {
			groupInfos = append(groupInfos, iam.GroupInfo{
				Id:   types.StringValue(groupInfo.Id),
				Name: types.StringValue(groupInfo.Name),
			})
		}

		userState := iam.UserWithoutPolicyDetail{
			AccountId:            types.StringValue(*user.AccountId.Get()),
			CompanyName:          types.StringValue(*userCompanyName),
			CreatedAt:            types.StringValue(user.CreatedAt.Format(time.RFC3339)),
			CreatedBy:            types.StringValue(user.CreatedBy),
			Description:          types.StringValue(*user.Description.Get()),
			DstOffset:            types.StringValue(*user.DstOffset.Get()),
			Email:                types.StringValue(*user.Email.Get()),
			EmailAuthenticated:   types.BoolValue(user.EmailAuthenticated),
			FirstName:            types.StringValue(*userFirstName),
			Id:                   types.StringValue(user.Id),
			LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
			LastName:             types.StringValue(*userLastName),
			LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
			ModifiedAt:           types.StringValue(user.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:           types.StringValue(user.ModifiedBy),
			Name:                 types.StringValue(*user.Name.Get()),
			PasswordReuseCount:   types.Int32Value(user.PasswordReuseCount),
			PhoneAuthenticated:   types.BoolValue(user.PhoneAuthenticated),
			Policies:             policies,
			Timezone:             types.StringValue(*user.Timezone.Get()),
			Type:                 types.StringValue(user.Type),
			TzId:                 types.StringValue(*user.TzId.Get()),
			UserName:             types.StringValue(*user.UserName.Get()),
			UtcOffset:            types.StringValue(*user.UtcOffset.Get()),
			AccessKeys:           accessKeyInfos,
			Groups:               groupInfos,
		}

		users = append(users, userState)
	}
	return users, false
}
