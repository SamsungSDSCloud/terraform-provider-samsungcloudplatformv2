package iamidentitycenter

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamIdentityCenterInstanceDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterInstanceDataSource{}
)

// NewIamIdentityCenterInstanceDataSource is a helper function to simplify the provider implementation.
func NewIamIdentityCenterInstanceDataSource() datasource.DataSource {
	return &iamIdentityCenterInstanceDataSource{}
}

type iamIdentityCenterInstanceDataSource struct {
	clients *client.SCPClient
}

type instanceDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	IdentityStoreType       types.String `tfsdk:"identity_store_type"`
	IdentityStoreConfig     types.Object `tfsdk:"identity_store_config"`
	SelfManagedPassword     types.Bool   `tfsdk:"self_managed_password"`
	Region                  types.String `tfsdk:"region"`
	AccountId               types.String `tfsdk:"account_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	CreatorName             types.String `tfsdk:"creator_name"`
	DelegatedAccountId      types.String `tfsdk:"delegated_account_id"`
	DelegatedAt             types.String `tfsdk:"delegated_at"`
	DelegatedBy             types.String `tfsdk:"delegated_by"`
	IdentityStoreId         types.String `tfsdk:"identity_store_id"`
	LastSyncAt              types.String `tfsdk:"last_sync_at"`
	LastSyncStatus          types.String `tfsdk:"last_sync_status"`
	MasterAccountId         types.String `tfsdk:"master_account_id"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifierName            types.String `tfsdk:"modifier_name"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	Port                    types.Int64  `tfsdk:"port"`
	PublicIp                types.String `tfsdk:"public_ip"`
	ResourceName            types.String `tfsdk:"resource_name"`
	ResourceType            types.String `tfsdk:"resource_type"`
	ResourceTypeDisplayName types.String `tfsdk:"resource_type_display_name"`
	Service                 types.String `tfsdk:"service"`
	ServiceName             types.String `tfsdk:"service_name"`
	Srn                     types.String `tfsdk:"srn"`
	SsoStartUrl             types.String `tfsdk:"sso_start_url"`
	State                   types.String `tfsdk:"state"`
}

// Metadata returns the data source type name.
func (d *iamIdentityCenterInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_instance"
}

// Configure adds the provider configured client to the data source.
func (d *iamIdentityCenterInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.clients = inst.Client
}

func (d *iamIdentityCenterInstanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM Identity Center Instance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Instance ID\n  - example: ssoins-12345",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Instance Name\n  - example: My Instance",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Instance Description\n  - example: My Instance Description",
			},
			"identity_store_type": schema.StringAttribute{
				Computed:    true,
				Description: "- enum: [\"IDENTITY_CENTER_DIRECTORY\",\"ACTIVE_DIRECTORY\",\"OPENLDAP\",\"DIRECTORY_389DS\"]",
			},
			"identity_store_config": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Identity Store Config\n" +
					"  - example: map[bind_credential:password bind_dn:CN=Administrator,CN=Users,DC=test,DC=add connection_url:ldaps://test.subdomain.company.com:123]",
				Attributes: map[string]schema.Attribute{
					"bind_credential": schema.StringAttribute{
						Computed:    true,
						Description: "Bind Password\n  - maxLength: 128",
					},
					"bind_dn": schema.StringAttribute{
						Computed:    true,
						Description: "Bind DN\n  - maxLength: 255",
					},
					"connection_url": schema.StringAttribute{
						Computed:    true,
						Description: "LDAP Server URL\n  - maxLength: 255",
					},
					"rdn_ldap_attribute": schema.StringAttribute{
						Computed:    true,
						Description: "RDN LDAP Attribute\n  - maxLength: 64",
					},
					"user_object_classes": schema.StringAttribute{
						Computed:    true,
						Description: "User Object Classes\n  - maxLength: 100",
					},
					"username_ldap_attribute": schema.StringAttribute{
						Computed:    true,
						Description: "Username LDAP Attribute\n  - maxLength: 64",
					},
					"users_dn": schema.StringAttribute{
						Computed:    true,
						Description: "Users DN\n  - maxLength: 255",
					},
				},
			},
			"self_managed_password": schema.BoolAttribute{
				Computed:    true,
				Description: "Self Managed Password",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Region",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Account ID\n  - example: 1234ab567cd64769e8f9g490hi304891",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Created At\n  - example: 2024-05-17T00:23:17Z",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
			},
			"creator_name": schema.StringAttribute{
				Computed:    true,
				Description: "Creator Name\n  - example: John Doe",
			},
			"delegated_account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Delegated Account ID\n  - example: acc-2353",
			},
			"delegated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Delegation Date\n  - example: 2023-10-15 14:30:00",
			},
			"delegated_by": schema.StringAttribute{
				Computed:    true,
				Description: "Delegated By\n  - example: sef2w3c8c29a449dbfa8681f8f1d78e2",
			},
			"identity_store_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identity Store ID\n  - example: d-12345",
			},
			"last_sync_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last Sync Date\n  - example: 2025-08-06 10:30:00",
			},
			"last_sync_status": schema.StringAttribute{
				Computed:    true,
				Description: "- enum: [\"CREATED\",\"RUNNING\",\"FINISHED\",\"FAILED\",\"NONE\",\"ABORTED\"]",
			},
			"master_account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Master Account ID\n  - example: acc-12345",
			},
			"modified_at": schema.StringAttribute{
				Computed:    true,
				Description: "Modified At\n  - example: 2024-05-17T00:23:17Z",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
			},
			"modifier_name": schema.StringAttribute{
				Computed:    true,
				Description: "Modifier Name\n  - example: Smith",
			},
			"organization_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization ID\n  - example: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5",
			},
			"port": schema.Int64Attribute{
				Computed:    true,
				Description: "Port\n  - example: 443",
			},
			"public_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public IP",
			},
			"resource_name": schema.StringAttribute{
				Computed:    true,
				Description: "Resource Name\n  - example: resource-name-1",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "Resource Type [instance|permission-set]\n  - example: virtual-server",
			},
			"resource_type_display_name": schema.StringAttribute{
				Computed:    true,
				Description: "Resource Type Display Name\n  - example: Virtual Server",
			},
			"service": schema.StringAttribute{
				Computed:    true,
				Description: "Service\n  - example: virtualserver",
			},
			"service_name": schema.StringAttribute{
				Computed:    true,
				Description: "Service Name\n  - example: Virtual Server",
			},
			"srn": schema.StringAttribute{
				Computed:    true,
				Description: "Instance SRN\n  - example: srn:identity-center::prj-01234:kr-west-1::instance/ins-01234",
			},
			"sso_start_url": schema.StringAttribute{
				Computed:    true,
				Description: "SSO Start URL\n  - example: https://login.e.samsungsdscloud.com/d-12345",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "- enum: [\"ACTIVE\",\"INACTIVE\"]",
			},
		},
	}
}

func (d *iamIdentityCenterInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data instanceDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := data.ID.ValueString()
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"ID Required",
			"An 'id' must be provided.",
		)
		return
	}

	result, httpResp, err := d.clients.IamIdentityCenter.GetInstance(ctx, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.Diagnostics.AddWarning(
				"IAM Identity Center Instance not found",
				fmt.Sprintf("Instance with ID '%s' was not found. This may indicate the resource has been deleted.", instanceId),
			)
			resp.State.Set(ctx, &data)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Instance",
			err.Error(),
		)
		return
	}

	instanceDetail := result.Instance

	state := instanceDataSourceModel{
		ID:                      data.ID,
		Name:                    types.StringValue(instanceDetail.Name),
		Description:             toStringValue(instanceDetail.Description),
		IdentityStoreType:       toStringValueFromEnum(instanceDetail.IdentityStoreType),
		SelfManagedPassword:     toBoolValue(instanceDetail.SelfManagedPassword),
		Region:                  types.StringValue(instanceDetail.Region),
		AccountId:               toStringValue(instanceDetail.AccountId),
		CreatedAt:               types.StringValue(instanceDetail.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(instanceDetail.CreatedBy),
		CreatorName:             toStringValue(instanceDetail.CreatorName),
		DelegatedAccountId:      toStringValue(instanceDetail.DelegatedAccountId),
		DelegatedAt:             toTimeStringValue(instanceDetail.DelegatedAt),
		DelegatedBy:             toStringValue(instanceDetail.DelegatedBy),
		IdentityStoreId:         types.StringValue(instanceDetail.IdentityStoreId),
		LastSyncAt:              toStringValue(instanceDetail.LastSyncAt),
		LastSyncStatus:          toStringValueFromEnum(instanceDetail.LastSyncStatus),
		MasterAccountId:         types.StringValue(instanceDetail.MasterAccountId),
		ModifiedAt:              types.StringValue(instanceDetail.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(instanceDetail.ModifiedBy),
		ModifierName:            toStringValue(instanceDetail.ModifierName),
		OrganizationId:          types.StringValue(instanceDetail.OrganizationId),
		Port:                    toInt64Value(instanceDetail.Port),
		PublicIp:                toStringValue(instanceDetail.PublicIp),
		ResourceName:            toStringValue(instanceDetail.ResourceName),
		ResourceType:            toStringValue(instanceDetail.ResourceType),
		ResourceTypeDisplayName: toStringValue(instanceDetail.ResourceTypeDisplayName),
		Service:                 toStringValue(instanceDetail.Service),
		ServiceName:             toStringValue(instanceDetail.ServiceName),
		Srn:                     types.StringValue(instanceDetail.Srn),
		SsoStartUrl:             types.StringValue(instanceDetail.SsoStartUrl),
		State:                   toStringValueFromEnum(instanceDetail.State),
	}

	// Convert identity store config
	identityStoreConfigType := planmodel{
		BindCredential:        types.StringNull(),
		BindDn:                types.StringNull(),
		ConnectionUrl:         types.StringNull(),
		RdnLdapAttribute:      types.StringNull(),
		UserObjectClasses:     types.StringNull(),
		UsernameLdapAttribute: types.StringNull(),
		UsersDn:               types.StringNull(),
	}

	if instanceDetail.IdentityStoreConfig.IsSet() {
		config := instanceDetail.IdentityStoreConfig.Get()
		if config != nil {
			identityStoreConfigType.BindCredential = toStringValueFromPtr(config.BindCredential)
			identityStoreConfigType.BindDn = toStringValueFromPtr(config.BindDn)
			identityStoreConfigType.ConnectionUrl = toStringValueFromPtr(config.ConnectionUrl)
			identityStoreConfigType.RdnLdapAttribute = toStringValueFromPtr(config.RdnLdapAttribute)
			identityStoreConfigType.UserObjectClasses = toStringValueFromPtr(config.UserObjectClasses)
			identityStoreConfigType.UsernameLdapAttribute = toStringValueFromPtr(config.UsernameLdapAttribute)
			identityStoreConfigType.UsersDn = toStringValueFromPtr(config.UsersDn)
		}
	}

	identityStoreConfigObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"bind_credential":         types.StringType,
		"bind_dn":                 types.StringType,
		"connection_url":          types.StringType,
		"rdn_ldap_attribute":      types.StringType,
		"user_object_classes":     types.StringType,
		"username_ldap_attribute": types.StringType,
		"users_dn":                types.StringType,
	}, identityStoreConfigType)
	resp.Diagnostics.Append(diags...)
	state.IdentityStoreConfig = identityStoreConfigObj

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
