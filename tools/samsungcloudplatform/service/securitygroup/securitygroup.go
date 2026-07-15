package securitygroup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/securitygroup" // securitygroup client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
)

// NewSecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

// securityGroupResource is the data source implementation.
type securityGroupResource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *securityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group"
}

// Schema defines the schema for the data source.
func (r *securityGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Security Group",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the Security Group.\n" +
					"  - example: sg-web-prod\n" +
					"  - valid: All characters except 'default'\n" +
					"  - constraints: minLength: 1, maxLength: 255, duplicates allowed",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "A brief explanation or note about this resource.\n" +
					"  - example: Security group for web tier\n" +
					"  - constraints: maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("Loggable"): schema.BoolAttribute{
				Description: "Enable flow log for the security group.\n" +
					"  - example: true\n" +
					"  - valid: true / false",
				Optional: true,
			},
			common.ToSnakeCase("SecurityGroup"): schema.SingleNestedAttribute{
		Description: "Manages security groups to protect virtual networks.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example: 297615908b8e4ec69520a99a6777add3",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the Security Group.\n" +
							"  - example: sg-web-prod",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Security group for web tier",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "Enable flow log for the security group.\n" +
							"  - example: true",
						Computed: true,
					},
					common.ToSnakeCase("RuleCount"): schema.Int32Attribute{
						Description: "Number of rules in the Security Group.\n" +
							"  - example: 5",
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
func (r *securityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securitygroup.SecurityGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new SecurityGroup
	data, err := r.client.CreateSecurityGroup(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating security group",
			"Could not create security group, unexpected error: "+err.Error(),
		)
		return
	}

	if data == nil || data.SecurityGroup.Id == "" {
		resp.Diagnostics.AddError(
			"Error creating security group",
			"empty response from API",
		)
		return
	}
	securityGroup := data.SecurityGroup

	plan.Id = types.StringValue(securityGroup.Id)

	sgModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		RuleCount:   types.Int32PointerValue(securityGroup.RuleCount),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, d := types.ObjectValueFrom(ctx, sgModel.AttributeTypes(), sgModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SecurityGroup = sgObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securitygroup.SecurityGroupResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from vpc
	data, err := r.client.GetSecurityGroup(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil || data.SecurityGroup.Id == "" {
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"empty response from API",
		)
		return
	}
	securityGroup := data.SecurityGroup

	sgrModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, d := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Name = types.StringValue(securityGroup.Name)
	state.Description = types.StringPointerValue(securityGroup.Description.Get())
	state.Loggable = types.BoolValue(securityGroup.Loggable)
	state.SecurityGroup = sgObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state securitygroup.SecurityGroupResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.UpdateSecurityGroup(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetSecurityGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	securityGroup := data.SecurityGroup
	if data == nil || securityGroup.Id == "" {
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"empty response from API",
		)
		return
	}

	sgrModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, d := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SecurityGroup = sgObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securitygroup.SecurityGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing network logging storage
	err := r.client.DeleteSecurityGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting security group",
			"Could not delete security group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
