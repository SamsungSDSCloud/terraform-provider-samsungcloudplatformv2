package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/securitygroup" // securitygroup client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &securityGroupRuleResource{}
	_ resource.ResourceWithConfigure = &securityGroupRuleResource{}
)

// NewSecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupRuleResource() resource.Resource {
	return &securityGroupRuleResource{}
}

// securityGroupResource is the data source implementation.
type securityGroupRuleResource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

func (r *securityGroupRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"This resource does not support in-place updates.",
	)
}

// Metadata returns the data source type name.
func (r *securityGroupRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group_rule"
}

// Schema defines the schema for the data source.
func (r *securityGroupRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages security group rules to filter network traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
				Description: "The identifier of the security group that the resource belongs to.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("ethertype"): schema.StringAttribute{
				Description: "The Ethernet protocol type the rule applies to.\n" +
					"  - example: IPv4\n" +
					"  - valid: IPv4",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4"),
				},
			},
			common.ToSnakeCase("protocol"): schema.StringAttribute{
				Description: "The network protocol the rule applies to.\n" +
					"  - example: TCP\n" +
					"  - valid: None/TCP/UDP/ICMP/1-254(IP Protocol Number)",
				Optional: true,
			},
			common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
				Description: "The minimum port number of the rule's port range.\n" +
					"  - example: 22\n" +
					"  - valid: 1-65535. For ICMP, 0-255. For IP Protocol, None.\n" +
					"  - constraints: None",
				Optional: true,
			},
			common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
				Description: "The maximum port number of the rule's port range.\n" +
					"  - example: 443\n" +
					"  - valid: 1-65535. For ICMP and IP Protocol, None.\n" +
					"  - constraints: None",
				Optional: true,
			},
			common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
				Description: "The remote IP address range the rule applies to in CIDR notation.\n" +
					"  - example: 10.0.0.0/24\n" +
					"  - valid: IPv4 CIDR",
				Optional: true,
			},
			common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
				Description: "The identifier of the remote security group the rule applies to.\n" +
					"  - example: ce5a565f-20fa-48f7-b06d-be0f03d2b50c",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "A brief explanation or note about this resource.\n" +
					"  - example: Security group for web tier\n" +
					"  - constraints: maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("Direction"): schema.StringAttribute{
				Description: "The direction of the traffic the rule applies to.\n" +
					"  - example: ingress\n" +
					"  - valid: ingress, egress",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ingress", "egress"),
				},
			},
			common.ToSnakeCase("SecurityGroupRule"): schema.SingleNestedAttribute{
				Description: "Security Group Rule Object",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
						Description: "The identifier of the security group that the resource belongs to.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ethertype"): schema.StringAttribute{
						Description: "The Ethernet protocol type the rule applies to.\n" +
							"  - example: IPv4",
						Computed: true,
					},
					common.ToSnakeCase("protocol"): schema.StringAttribute{
						Description: "The network protocol the rule applies to.\n" +
							"  - example: TCP",
						Computed: true,
					},
					common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
						Description: "The minimum port number of the rule's port range.\n" +
							"  - example: 5",
						Computed: true,
					},
					common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
						Description: "The maximum port number of the rule's port range.\n" +
							"  - example: 10",
						Computed: true,
					},
					common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
						Description: "The remote IP address range the rule applies to in CIDR notation.\n" +
							"  - example: 10.0.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
						Description: "The identifier of the remote security group the rule applies to.\n" +
							"  - example: ce5a565f-20fa-48f7-b06d-be0f03d2b50c",
						Computed: true,
					},
					common.ToSnakeCase("RemoteGroupName"): schema.StringAttribute{
						Description: "The name of the remote security group the rule applies to.\n" +
							"  - example: sg-db-prod",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Security group for web tier",
						Computed: true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "The direction of the traffic the rule applies to.\n" +
							"  - example: ingress",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *securityGroupRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = inst.Client.SecurityGroup
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securitygroup.SecurityGroupRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new security group rule
	data, err := r.client.CreateSecurityGroupRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating security group rule",
			"Could not create security group rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroupRule := data.SecurityGroupRule
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(data.SecurityGroupRule.Id)

	sgrModel := securitygroup.SecurityGroupRule{
		Id:              types.StringValue(securityGroupRule.Id),
		SecurityGroupId: types.StringValue(securityGroupRule.SecurityGroupId),
		Ethertype:       types.StringPointerValue(securityGroupRule.Ethertype.Get()),
		Protocol:        types.StringPointerValue(securityGroupRule.Protocol.Get()),
		PortRangeMin:    types.Int32PointerValue(securityGroupRule.PortRangeMin.Get()),
		PortRangeMax:    types.Int32PointerValue(securityGroupRule.PortRangeMax.Get()),
		RemoteIpPrefix:  types.StringPointerValue(securityGroupRule.RemoteIpPrefix.Get()),
		RemoteGroupId:   types.StringPointerValue(securityGroupRule.RemoteGroupId.Get()),
		RemoteGroupName: types.StringPointerValue(securityGroupRule.RemoteGroupName.Get()),
		Description:     types.StringPointerValue(securityGroupRule.Description.Get()),
		Direction:       types.StringValue(string(securityGroupRule.Direction)),
		CreatedAt:       types.StringValue(securityGroupRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(securityGroupRule.CreatedBy),
		ModifiedAt:      types.StringValue(securityGroupRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:      types.StringValue(securityGroupRule.ModifiedBy),
	}

	sgrObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	plan.SecurityGroupRule = sgrObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securitygroup.SecurityGroupRuleResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from security group rule
	data, err := r.client.GetSecurityGroupRule(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading security group rule",
			"Could not read security group rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroupRule := data.SecurityGroupRule

	sgrModel := securitygroup.SecurityGroupRule{
		Id:              types.StringValue(securityGroupRule.Id),
		SecurityGroupId: types.StringValue(securityGroupRule.SecurityGroupId),
		Ethertype:       types.StringPointerValue(securityGroupRule.Ethertype.Get()),
		Protocol:        types.StringPointerValue(securityGroupRule.Protocol.Get()),
		PortRangeMin:    types.Int32PointerValue(securityGroupRule.PortRangeMin.Get()),
		PortRangeMax:    types.Int32PointerValue(securityGroupRule.PortRangeMax.Get()),
		RemoteIpPrefix:  types.StringPointerValue(securityGroupRule.RemoteIpPrefix.Get()),
		RemoteGroupId:   types.StringPointerValue(securityGroupRule.RemoteGroupId.Get()),
		RemoteGroupName: types.StringPointerValue(securityGroupRule.RemoteGroupName.Get()),
		Description:     types.StringPointerValue(securityGroupRule.Description.Get()),
		Direction:       types.StringValue(string(securityGroupRule.Direction)),
		CreatedAt:       types.StringValue(securityGroupRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(securityGroupRule.CreatedBy),
		ModifiedAt:      types.StringValue(securityGroupRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:      types.StringValue(securityGroupRule.ModifiedBy),
	}
	sgrObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	state.SecurityGroupRule = sgrObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securitygroup.SecurityGroupRuleResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing security group rule
	err := r.client.DeleteSecurityGroupRule(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting security group rule",
			"Could not delete security group rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
