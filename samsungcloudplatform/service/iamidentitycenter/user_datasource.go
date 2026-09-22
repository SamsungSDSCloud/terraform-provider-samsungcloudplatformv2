package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamIdentityCenterUserDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterUserDataSource{}
)

func NewIamIdentityCenterUserDataSource() datasource.DataSource {
	return &iamIdentityCenterUserDataSource{}
}

type iamIdentityCenterUserDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_user"
}

func (r *iamIdentityCenterUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves details of an IAM Identity Center User.",
		MarkdownDescription: "Retrieves details of an IAM Identity Center User.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"user_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User Login ID\n  - example: johndoe",
				MarkdownDescription: "User Login ID\n  - example: johndoe",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Real Name\n  - example: John Doe",
				MarkdownDescription: "Real Name\n  - example: John Doe",
			},
			"instance_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				Description:         "User Email\n  - example: john.doe@example.com",
				MarkdownDescription: "User Email\n  - example: john.doe@example.com",
			},
			"phone_number": schema.StringAttribute{
				Computed:            true,
				Description:         "Phone Number\n  - example: 010-1234-5678",
				MarkdownDescription: "Phone Number\n  - example: 010-1234-5678",
			},
			"temporary_password": schema.BoolAttribute{
				Computed:            true,
				Description:         "Temporary Password Status",
				MarkdownDescription: "Temporary Password Status",
			},
			"business_unit": schema.StringAttribute{
				Computed:            true,
				Description:         "Business Unit\n  - example: Business Unit A",
				MarkdownDescription: "Business Unit\n  - example: Business Unit A",
			},
			"department": schema.StringAttribute{
				Computed:            true,
				Description:         "Department\n  - example: Department X",
				MarkdownDescription: "Department\n  - example: Department X",
			},
			"manager": schema.StringAttribute{
				Computed:            true,
				Description:         "Manager\n  - example: Alice Smith",
				MarkdownDescription: "Manager\n  - example: Alice Smith",
			},
			"employee_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Employee ID\n  - example: emp-12345",
				MarkdownDescription: "Employee ID\n  - example: emp-12345",
			},
			"nation_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Nationality ID\n  - example: +82",
				MarkdownDescription: "Nationality ID\n  - example: +82",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "User Description\n  - example: Description of John Doe",
				MarkdownDescription: "User Description\n  - example: Description of John Doe",
			},
			"user_uuid": schema.StringAttribute{
				Computed:            true,
				Description:         "User ID",
				MarkdownDescription: "User ID",
			},
			"user": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "User Details",
				MarkdownDescription: "User Details",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
						MarkdownDescription: "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
					},
					"user_id": schema.StringAttribute{
						Computed:            true,
						Description:         "User Login ID\n  - example: johndoe",
						MarkdownDescription: "User Login ID\n  - example: johndoe",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Real Name\n  - example: John Doe",
						MarkdownDescription: "Real Name\n  - example: John Doe",
					},
					"email": schema.StringAttribute{
						Computed:            true,
						Description:         "User Email\n  - example: john.doe@example.com",
						MarkdownDescription: "User Email\n  - example: john.doe@example.com",
					},
					"phone_number": schema.StringAttribute{
						Computed:            true,
						Description:         "Phone Number\n  - example: 010-1234-5678",
						MarkdownDescription: "Phone Number\n  - example: 010-1234-5678",
					},
					"temporary_password": schema.BoolAttribute{
						Computed:            true,
						Description:         "Temporary Password Status",
						MarkdownDescription: "Temporary Password Status",
					},
					"business_unit": schema.StringAttribute{
						Computed:            true,
						Description:         "Business Unit\n  - example: Business Unit A",
						MarkdownDescription: "Business Unit\n  - example: Business Unit A",
					},
					"department": schema.StringAttribute{
						Computed:            true,
						Description:         "Department\n  - example: Department X",
						MarkdownDescription: "Department\n  - example: Department X",
					},
					"manager": schema.StringAttribute{
						Computed:            true,
						Description:         "Manager\n  - example: Alice Smith",
						MarkdownDescription: "Manager\n  - example: Alice Smith",
					},
					"employee_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Employee ID\n  - example: emp-12345",
						MarkdownDescription: "Employee ID\n  - example: emp-12345",
					},
					"nation_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Nationality ID\n  - example: +82",
						MarkdownDescription: "Nationality ID\n  - example: +82",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "User Description\n  - example: Description of John Doe",
						MarkdownDescription: "User Description\n  - example: Description of John Doe",
					},
					"instance_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Instance ID\n  - example: ssoins-12345",
						MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Name\n  - example: John Doe",
						MarkdownDescription: "Creator Name\n  - example: John Doe",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"modifier_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Name\n  - example: Smith",
						MarkdownDescription: "Modifier Name\n  - example: Smith",
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.clients = inst.Client
}

type userDataSourceModel struct {
	Id                types.String               `tfsdk:"id"`
	UserId            types.String               `tfsdk:"user_id"`
	Name              types.String               `tfsdk:"name"`
	InstanceId        types.String               `tfsdk:"instance_id"`
	Email             types.String               `tfsdk:"email"`
	PhoneNumber       types.String               `tfsdk:"phone_number"`
	TemporaryPassword types.Bool                 `tfsdk:"temporary_password"`
	BusinessUnit      types.String               `tfsdk:"business_unit"`
	Department        types.String               `tfsdk:"department"`
	Manager           types.String               `tfsdk:"manager"`
	EmployeeId        types.String               `tfsdk:"employee_id"`
	NationId          types.String               `tfsdk:"nation_id"`
	Description       types.String               `tfsdk:"description"`
	UserUuid          types.String               `tfsdk:"user_uuid"`
	User              *userDataSourceDetailModel `tfsdk:"user"`
}

type userDataSourceDetailModel struct {
	Id                types.String `tfsdk:"id"`
	UserId            types.String `tfsdk:"user_id"`
	Name              types.String `tfsdk:"name"`
	Email             types.String `tfsdk:"email"`
	PhoneNumber       types.String `tfsdk:"phone_number"`
	TemporaryPassword types.Bool   `tfsdk:"temporary_password"`
	BusinessUnit      types.String `tfsdk:"business_unit"`
	Department        types.String `tfsdk:"department"`
	Manager           types.String `tfsdk:"manager"`
	EmployeeId        types.String `tfsdk:"employee_id"`
	NationId          types.String `tfsdk:"nation_id"`
	Description       types.String `tfsdk:"description"`
	InstanceId        types.String `tfsdk:"instance_id"`
	CreatedAt         types.String `tfsdk:"created_at"`
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatorName       types.String `tfsdk:"creator_name"`
	ModifiedAt        types.String `tfsdk:"modified_at"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
	ModifierName      types.String `tfsdk:"modifier_name"`
}

func (r *iamIdentityCenterUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	userId := state.Id.ValueString()
	result, _, err := r.clients.IamIdentityCenter.GetUser(ctx, instanceId, userId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center User",
			err.Error(),
		)
		return
	}

	if result != nil && result.User.Id != "" {
		state.Id = types.StringValue(result.User.Id)
		state.UserUuid = types.StringValue(result.User.Id)
		state.UserId = types.StringValue(result.User.UserId)
		state.Name = types.StringValue(result.User.Name)
		state.InstanceId = types.StringValue(result.User.InstanceId)
		state.Email = stringValuePtr(result.User.Email.Get())
		state.PhoneNumber = stringValuePtr(result.User.PhoneNumber.Get())
		state.BusinessUnit = stringValuePtr(result.User.BusinessUnit.Get())
		state.Department = stringValuePtr(result.User.Department.Get())
		state.Manager = stringValuePtr(result.User.Manager.Get())
		state.EmployeeId = stringValuePtr(result.User.EmployeeId.Get())
		state.NationId = stringValuePtr(result.User.NationId.Get())
		state.Description = stringValuePtr(result.User.Description.Get())
		state.TemporaryPassword = types.BoolValue(result.User.TemporaryPassword)

		state.User = &userDataSourceDetailModel{}
		state.User.Id = types.StringValue(result.User.Id)
		state.User.UserId = types.StringValue(result.User.UserId)
		state.User.Name = types.StringValue(result.User.Name)
		state.User.Email = stringValuePtr(result.User.Email.Get())
		state.User.PhoneNumber = stringValuePtr(result.User.PhoneNumber.Get())
		state.User.BusinessUnit = stringValuePtr(result.User.BusinessUnit.Get())
		state.User.Department = stringValuePtr(result.User.Department.Get())
		state.User.Manager = stringValuePtr(result.User.Manager.Get())
		state.User.EmployeeId = stringValuePtr(result.User.EmployeeId.Get())
		state.User.NationId = stringValuePtr(result.User.NationId.Get())
		state.User.Description = stringValuePtr(result.User.Description.Get())
		state.User.InstanceId = types.StringValue(result.User.InstanceId)
		state.User.CreatedAt = types.StringValue(result.User.CreatedAt.Format("2006-01-02T15:04:05Z"))
		state.User.CreatedBy = types.StringValue(result.User.CreatedBy)
		state.User.CreatorName = stringValuePtr(result.User.CreatorName)
		state.User.ModifiedAt = types.StringValue(result.User.ModifiedAt.Format("2006-01-02T15:04:05Z"))
		state.User.ModifiedBy = types.StringValue(result.User.ModifiedBy)
		state.User.ModifierName = stringValuePtr(result.User.ModifierName)
		if result.User.TemporaryPassword != false {
			state.TemporaryPassword = types.BoolValue(result.User.TemporaryPassword)
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
