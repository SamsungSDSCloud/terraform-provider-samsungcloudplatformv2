package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcTgwResource{}
	_ resource.ResourceWithConfigure   = &vpcTgwResource{}
	_ resource.ResourceWithImportState = &vpcTgwResource{}
	_ resource.ResourceWithModifyPlan  = &vpcTgwResource{}
)

// NewVpcTgwResource is a helper function to simplify the provider implementation.
func NewVpcTgwResource() resource.Resource {
	return &vpcTgwResource{}
}

type vpcTgwResource struct {
	config  *scpsdk.Configuration
	client  *vpcv1d3.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcTgwResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway"
}

func (r *vpcTgwResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

// Schema defines the schema for the data source.
func (r *vpcTgwResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit Gateway",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the transit gateway.\n" +
					"  - example : fe860e0af0c04dcd8182b84f907f31f4",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the transit gateway.\n" +
					"  - example : TransitGatewayName",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : Transit Gateway Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Tgw"): schema.SingleNestedAttribute{
				Description: "Tgw",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the transit gateway.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the transit gateway.\n" +
							"  - example : fe860e0af0c04dcd8182b84f907f31f4",
						Computed: true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "The bandwidth capacity of the connection.\n" +
							"  - example : 1",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the transit gateway was created in ISO 8601 format.\n " +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the transit gateway. \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this transit gateway. This help identify the purpose or usage of the resource.\n" +
							"  - example : Tgw description\n" +
							"  - maxLength : 50\n" +
							"  - minLength : 1",
						Computed: true,
					},
					common.ToSnakeCase("firewall_connection_state"): schema.StringAttribute{
						Description: "The current lifecycle state of the firewall connection. \n" +
							"  - example : INACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
						Description: "List of the FirewallIds\n" +
							"  - example : bbb93aca123f4bb2b2c0f206f4a86b2b",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the transit gateway was last modified in ISO 8601 format. \n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that modified the transit gateway. \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the transit gateway.\n" +
							"  - example : Tgw name",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the transit gateway.\n" +
							"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING \n" +
							"  - example : CREATING",
						Computed: true,
					},
					common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
						Description: "Whether the uplink is enabled.\n" +
							"  - example : false",
						Computed: true,
					},
					common.ToSnakeCase("UplinkActiveZone"): schema.StringAttribute{
						Description: "The active availability zone for the uplink connection.\n" +
							"  - example : apigw-kr-1",
						Computed: true,
					},
					common.ToSnakeCase("UplinkStandbyZone"): schema.StringAttribute{
						Description: "The standby availability zone for the uplink connection.\n" +
							"  - example : apigw-kr-2",
						Computed: true,
					},
					common.ToSnakeCase("UplinkZoneState"): schema.StringAttribute{
						Description: "The current state of the uplink zone.\n" +
							"  - enum: ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcTgwResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1Dot3
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcTgwResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan vpcv1d3.TgwResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subnet
	data, _, err := r.client.CreateTransitGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating CreateTransitGateway",
			"Could not create CreateTransitGateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tgw := data.TransitGateway
	plan.Id = types.StringValue(tgw.Id)
	desc := tgw.Description.Get()
	if desc != nil {
		plan.Description = types.StringValue(*desc)
	} else {
		plan.Description = types.StringValue("")
	}
	diags = resp.State.Set(ctx, plan)

	err = waitForTgwtStatus(ctx, r.client, tgw.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tgw",
			"Error waiting for CreateTransitGateway to become active: "+err.Error(),
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
func (r *vpcTgwResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	var state vpcv1d3.TgwResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetTransitGatewayInfo(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading tgw",
			"Could not read transitGateway  "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	transitGateway := data.TransitGateway

	tgwModel := vpcv1d3.MapToTgw(transitGateway)
	tgwObjectValue, diags := types.ObjectValueFrom(ctx, tgwModel.AttributeTypes(), tgwModel)
	state.Name = tgwModel.Name
	state.Tgw = tgwObjectValue
	state.Description = tgwModel.Description

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcTgwResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpcv1d3.TgwResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing tgw
	err := r.client.UpdateTransitGateway(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating tgw",
			"Could not update tgw, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetTransitGatewayInfo(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GetTransitGatewayInfo",
			"Could not read GetTransitGatewayInfo ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tgwElem := data.TransitGateway

	tgwModel := vpcv1d3.MapToTgw(tgwElem)
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, tgwModel.AttributeTypes(), tgwModel)
	state.Tgw = vpcObjectValue
	state.Description = tgwModel.Description

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcTgwResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state

	var state vpcv1d3.TgwResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subnet
	err := r.client.DeleteTransitGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting TgwResource",
			"Could not delete TgwResource, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForTgwtStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting TgwResource",
			"Error waiting for TgwResource to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForTgwtStatus(ctx context.Context, vpcClient *vpcv1d3.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetTransitGatewayInfo(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.TransitGateway.State), nil
	}, -1, -1, -1, -1)
}

func (r *vpcTgwResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan vpcv1d3.TgwResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state vpcv1d3.TgwResource
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
		{"name", !plan.Name.Equal(state.Name)},
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

	// Reconstruct tgw: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_id, created_at, created_by from showing as (known after apply).
	if !state.Tgw.IsNull() && !state.Tgw.IsUnknown() {
		var stateTgw vpcv1d3.Tgw
		resp.Diagnostics.Append(state.Tgw.As(ctx, &stateTgw, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when description actually changes
		descriptionChanged := !plan.Description.Equal(state.Description)

		mergedTgw := vpcv1d3.Tgw{
			Id:                      stateTgw.Id,
			AccountId:               stateTgw.AccountId,
			Bandwidth:               stateTgw.Bandwidth,
			CreatedAt:               stateTgw.CreatedAt,
			CreatedBy:               stateTgw.CreatedBy,
			FirewallConnectionState: stateTgw.FirewallConnectionState,
			FirewallIds:             stateTgw.FirewallIds,
			Name:                    stateTgw.Name,
			State:                   stateTgw.State,
			UplinkEnabled:           stateTgw.UplinkEnabled,
			UplinkActiveZone:        stateTgw.UplinkActiveZone,
			UplinkStandbyZone:       stateTgw.UplinkStandbyZone,
			UplinkZoneState:         stateTgw.UplinkZoneState,
			Description:             plan.Description,
		}

		if descriptionChanged {
			mergedTgw.ModifiedAt = types.StringUnknown()
			mergedTgw.ModifiedBy = types.StringUnknown()
		} else {
			mergedTgw.ModifiedAt = stateTgw.ModifiedAt
			mergedTgw.ModifiedBy = stateTgw.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, mergedTgw.AttributeTypes(), mergedTgw)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Tgw = mergedObj
		resp.Plan.Set(ctx, plan)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
