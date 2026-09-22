package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcNatGatewayResource{}
	_ resource.ResourceWithConfigure   = &vpcNatGatewayResource{}
	_ resource.ResourceWithImportState = &vpcNatGatewayResource{}
	_ resource.ResourceWithModifyPlan  = &vpcNatGatewayResource{}
)

// NewVpcNatGatewayResource is a helper function to simplify the provider implementation.
func NewVpcNatGatewayResource() resource.Resource {
	return &vpcNatGatewayResource{}
}

// vpcNatGatewayResource is the data source implementation.
type vpcNatGatewayResource struct {
	config     *scpsdk.Configuration
	client     *vpc.Client
	clientV1d3 *vpcv1d3.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcNatGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_nat_gateway"
}

// Schema defines the schema for the data source.
func (r *vpcNatGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Nat Gateway resource",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the nat gateway.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the nat gateway belongs to.\n" +
					"  - example : 607e0938521643b5b4b266f343fae693",
				Required: true,
			},
			common.ToSnakeCase("PublicipIds"): schema.ListAttribute{
				Description: "A list of public IP address identifiers.\n" +
					"  - example : [\"023c57b14f11483689338d085e061492\", \"a7c2cac139b24183bf92a0af9633a316\"]",
				Required:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("MultiZoneEnabled"): schema.BoolAttribute{
				Description: "Indicates whether Multi-AZ is enabled for the NAT gateway.\n" +
					"  - example : true",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : NAT Gateway Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("NatGateway"): schema.SingleNestedAttribute{
				Description: "NatGateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the NAT gateway.\n" +
							"  - example : 12f56e27070248a6a240a497e43fbe18",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the NAT gateway.\n" +
							"  - example : NatGatewayName",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the NAT gateway belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the NAT gateway belongs to.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The identifier of the subnet that the NAT gateway belongs to.\n" +
							"  - example : 023c57b14f11483689338d085e061492",
						Computed: true,
					},
					common.ToSnakeCase("SubnetName"): schema.StringAttribute{
						Description: "The name of the subnet that the NAT gateway belongs to.\n" +
							"  - example : subnetName",
						Computed: true,
					},
					common.ToSnakeCase("SubnetCidr"): schema.StringAttribute{
						Description: "The IP address range of the subnet in CIDR notation.\n" +
							"  - example : 192.167.1.0/24",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the NAT gateway.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the NAT gateway.\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("MultiZoneEnabled"): schema.BoolAttribute{
						Description: "Indicates whether Multi-AZ is enabled for the NAT gateway.\n" +
							"  - example : true",
						Computed: true,
					},
					common.ToSnakeCase("NatGatewayIps"): schema.ListNestedAttribute{
						Description: "A list of NAT gateway IP addresses.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("IpAddress"): schema.StringAttribute{
									Description: "The IP address of the NAT gateway.\n" +
										"  - example : 42.15.165.56",
									Computed: true,
								},
								common.ToSnakeCase("PublicipId"): schema.StringAttribute{
									Description: "The identifier of the public IP address.\n" +
										"  - example : 390133c4259d43aebb24c2e0d66524f6",
									Computed: true,
								},
							},
						},
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : NAT Gateway Description",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcNatGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpc
	r.clientV1d3 = inst.Client.VpcV1Dot3
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcNatGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpcv1d3.NatGatewayResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new nat gateway using v1.3 API
	data, err := r.clientV1d3.CreateNatGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating nat gateway",
			"Could not create nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	natgateway := data.NatGateway
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(natgateway.Id)

	natGatewayModel := vpcv1d3.ResponseToNatGatewayValue(natgateway)
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.NatGateway = natGatewayObjectValue

	// Description might be default to "" in API response
	plan.Description = natGatewayModel.Description

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForNatGatewayStatus(ctx, r.client, natgateway.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating nat gateway",
			"Error waiting for nat gateway to become active: "+err.Error(),
		)
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcNatGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1d3.NatGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.clientV1d3.GetNatGateway(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading nat gateway",
			"Could not read nat gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error reading data",
			"An error occurred while reading data. Empty response",
		)
		return
	}

	natgateway := data.NatGateway

	natGatewayModel := vpcv1d3.ResponseToNatGatewayValue(natgateway)
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.NatGateway = natGatewayObjectValue

	// Refresh input attributes from API response
	state.SubnetId = types.StringValue(natgateway.SubnetId)
	state.Description = types.StringPointerValue(natgateway.Description.Get())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcNatGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan vpcv1d3.NatGatewayResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing nat gateway
	_, err := r.clientV1d3.UpdateNatGateway(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating nat gateway",
			"Could not update nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated value from GetNatGateway as UpdateNatGateway response is limited.
	data, err := r.clientV1d3.GetNatGateway(ctx, plan.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading nat gateway",
			"Could not read nat gateway ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	natgateway := data.NatGateway

	natGatewayModel := vpcv1d3.ResponseToNatGatewayValue(natgateway)
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.NatGateway = natGatewayObjectValue

	// Description might be default to "" in API response
	plan.Description = natGatewayModel.Description

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcNatGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1d3.NatGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing NatGateway
	err := r.client.DeleteNatGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting nat gateway",
			"Could not delete nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForNatGatewayStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting nat gateway",
			"Error waiting for nat gateway to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForNatGatewayStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetNatGateway(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, info.NatGateway.State, nil
	}, -1, -1, -1, -1)
}

func (r *vpcNatGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 1 || parts[0] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: nat_gateway_id, got: %q", req.ID),
		)
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[0]))
}

func (r *vpcNatGatewayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan vpcv1d3.NatGatewayResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state vpcv1d3.NatGatewayResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fields that cannot be updated via the API — check each for changes
	type fieldCheck struct {
		name    string
		changed bool
	}
	checks := []fieldCheck{
		{"subnet_id", !plan.SubnetId.Equal(state.SubnetId) && !state.SubnetId.IsNull()},
		{"publicip_ids", !plan.PublicipIds.Equal(state.PublicipIds) && !state.PublicipIds.IsNull()},
		{"tags", !plan.Tags.Equal(state.Tags)},
	}

	for _, f := range checks {
		if f.changed {
			resp.Diagnostics.AddError(
				"Field changes not supported",
				fmt.Sprintf("Changing `%s` will not update the actual resource. To change %s, recreate the resource.", f.name, f.name),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Reconstruct nat_gateway: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_id, created_at, created_by from showing as (known after apply).
	if !state.NatGateway.IsNull() && !state.NatGateway.IsUnknown() {
		var stateNg vpcv1d3.NatGatewayValue
		resp.Diagnostics.Append(state.NatGateway.As(ctx, &stateNg, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when description actually changes
		descriptionChanged := !plan.Description.Equal(state.Description)

		mergedNg := vpcv1d3.NatGatewayValue{
			Id:               stateNg.Id,
			Name:             stateNg.Name,
			NatGatewayIps:    stateNg.NatGatewayIps,
			VpcId:            stateNg.VpcId,
			VpcName:          stateNg.VpcName,
			SubnetId:         stateNg.SubnetId,
			SubnetName:       stateNg.SubnetName,
			SubnetCidr:       stateNg.SubnetCidr,
			AccountId:        stateNg.AccountId,
			State:            stateNg.State,
			MultiZoneEnabled: stateNg.MultiZoneEnabled,
			CreatedAt:        stateNg.CreatedAt,
			CreatedBy:        stateNg.CreatedBy,
			Description:      plan.Description,
		}

		if descriptionChanged {
			mergedNg.ModifiedAt = types.StringUnknown()
			mergedNg.ModifiedBy = types.StringUnknown()
		} else {
			mergedNg.ModifiedAt = stateNg.ModifiedAt
			mergedNg.ModifiedBy = stateNg.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, mergedNg.AttributeTypes(), mergedNg)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.NatGateway = mergedObj
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
