package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/iam/1.4"
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
		Description: "Show IAM Users",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)\n" +
					"  - example : 100",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page (between 0 and 10000)\n" +
					"  - example : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort order for results (e.g., 'createdAt,desc').\n" +
					"  - example : 'createdAt,desc'",
				Optional: true,
			},
			common.ToSnakeCase("Email"): schema.StringAttribute{
				Description: "Filter users by email address.\n" +
					"  - example : 'user@example.com'",
				Optional: true,
			},
			common.ToSnakeCase("UserName"): schema.StringAttribute{
				Description: "Filter users by username.\n" +
					"  - example : 'john.doe'",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Filter users by type.\n" +
					"  - example : 'scp'",
				Optional: true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID to filter users.\n" +
					"  - example : '123456789012'",
				Optional: true,
			},
			common.ToSnakeCase("Users"): schema.ListNestedAttribute{
				Description: "List of users matching the filter criteria.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed: true,
							Description: "Account ID of the user.\n" +
								"  - example : '123456789012'",
						},
						"company_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "Company name of the user.\n" +
								"  - example : 'Samsung SDS'",
						},
						"console_url": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "URL to access the console.\n" +
								"  - example : 'https://console.example.com'",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the user was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the user.\n" +
								"  - example : 'admin@example.com'",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Human-readable description of the user.\n" +
								"  - example : 'My user description'",
						},
						"dst_offset": schema.StringAttribute{
							Computed: true,
							Description: "Daylight saving time offset.\n" +
								"  - example : '+09:00'",
						},
						"email": schema.StringAttribute{
							Computed: true,
							Description: "Email address of the user.\n" +
								"  - example : 'user@example.com'",
						},
						"email_authenticated": schema.BoolAttribute{
							Computed: true,
							Description: "Whether the email has been authenticated.\n" +
								"  - example : true",
						},
						"first_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "First name of the user.\n" +
								"  - example : 'John'",
						},
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the user.\n" +
								"  - example : 'usr-1234567890abcdef'",
						},
						"last_login_at": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "Timestamp when the user last logged in.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"last_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "Last name of the user.\n" +
								"  - example : 'Doe'",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the user was last modified.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who last modified the user.\n" +
								"  - example : 'user@example.com'",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "User name.\n" +
								"  - example : 'John Doe'",
						},
						"password": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Description: "User password (masked for security).\n" +
								"  - example : '********'",
							MarkdownDescription: "User password (masked for security).\n" +
								"  - example : '********'",
						},
						"password_reuse_count": schema.Int64Attribute{
							Computed: true,
							Description: "Number of previous passwords that cannot be reused.\n" +
								"  - example : 3",
						},
						"phone_authenticated": schema.BoolAttribute{
							Computed: true,
							Description: "Whether the phone number has been authenticated.\n" +
								"  - example : true",
						},
						"policies": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "Policies",
							MarkdownDescription: "Policies",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Policy ID.\n" +
											"  - example : 'pol-1234567890abcdef'",
									},
									"name": schema.StringAttribute{
										Computed: true,
										Description: "Policy name.\n" +
											"  - example : 'MyPolicy'",
									},
								},
							},
						},
						"timezone": schema.StringAttribute{
							Computed: true,
							Description: "Timezone of the user.\n" +
								"  - example : 'Asia/Seoul'",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type of user.\n" +
								"  - example : 'IAM'",
						},
						"tz_id": schema.StringAttribute{
							Computed: true,
							Description: "Timezone ID.\n" +
								"  - example : 'Asia/Seoul'",
						},
						"user_name": schema.StringAttribute{
							Computed: true,
							Description: "Unique username.\n" +
								"  - example : 'john.doe'",
						},
						"utc_offset": schema.StringAttribute{
							Computed: true,
							Description: "UTC offset of the user.\n" +
								"  - example : '+09:00'",
						},
						"last_password_update_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the password was last updated.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"access_keys": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Access Keys",
							MarkdownDescription: "Access Keys",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"access_key": schema.StringAttribute{
										Computed: true,
										Description: "The access key string value.\n" +
											"  - example : 'ak-example-access-key-id'",
									},
									"created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the access key was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"expiration_timestamp": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the access key expires.\n" +
											"  - example : '2024-01-02T00:00:00Z'",
									},
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Unique identifier of the access key.\n" +
											"  - example : '12345678-1234-1234-1234-1234567890ab'",
									},
									"is_enabled": schema.BoolAttribute{
										Computed: true,
										Description: "Whether the access key is enabled/active.\n" +
											"  - example : true",
									},
								},
							},
						},
						"groups": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Groups",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Group ID.\n" +
											"  - example : 'grp-1234567890abcdef'",
									},
									"name": schema.StringAttribute{
										Computed: true,
										Description: "Group name.\n" +
											"  - example : 'MyGroup'",
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
