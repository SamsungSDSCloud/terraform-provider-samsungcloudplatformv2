package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1"
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
	_ resource.Resource                = &vpcPeeringResource{}
	_ resource.ResourceWithConfigure   = &vpcPeeringResource{}
	_ resource.ResourceWithImportState = &vpcPeeringResource{}
	_ resource.ResourceWithModifyPlan  = &vpcPeeringResource{}
)

// NewVpcVpcPeeringResource is a helper function to simplify the provider implementation.
func NewVpcPeeringResource() resource.Resource {
	return &vpcPeeringResource{}
}

// vpcVpcPeeringResource is the data source implementation.
type vpcPeeringResource struct {
	config  *scpsdk.Configuration
	client  *vpcv1.Client
	clients *client.SCPClient
}

func (r *vpcPeeringResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan vpcv1.VpcPeeringResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state vpcv1.VpcPeeringResource
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
		{"approver_vpc_account_id", !plan.ApproverVpcAccountId.Equal(state.ApproverVpcAccountId)},
		{"name", !plan.Name.Equal(state.Name)},
		{"approver_vpc_id", !plan.ApproverVpcId.Equal(state.ApproverVpcId)},
		{"requester_vpc_id", !plan.RequesterVpcId.Equal(state.RequesterVpcId)},
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

	// Reconstruct vpc_peering: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_type, approver_vpc_id, requester_vpc_id, etc. from showing as (known after apply).
	if !state.VpcPeering.IsNull() && !state.VpcPeering.IsUnknown() {
		var statePv vpcv1.VpcPeering
		resp.Diagnostics.Append(state.VpcPeering.As(ctx, &statePv, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when name or description actually changes
		changed := !plan.Description.Equal(state.Description)

		modelVpcPeering := vpcv1.VpcPeering{
			Id:                       statePv.Id,
			AccountType:              statePv.AccountType,
			ApproverVpcAccountId:     statePv.ApproverVpcAccountId,
			ApproverVpcId:            statePv.ApproverVpcId,
			ApproverVpcName:          statePv.ApproverVpcName,
			RequesterVpcAccountId:    statePv.RequesterVpcAccountId,
			RequesterVpcId:           statePv.RequesterVpcId,
			RequesterVpcName:         statePv.RequesterVpcName,
			DeleteRequesterAccountId: statePv.DeleteRequesterAccountId,
			CreatedAt:                statePv.CreatedAt,
			CreatedBy:                statePv.CreatedBy,
			State:                    statePv.State,
			Name:                     statePv.Name,
			Description:              plan.Description,
		}

		if changed {
			modelVpcPeering.ModifiedAt = types.StringUnknown()
			modelVpcPeering.ModifiedBy = types.StringUnknown()
		} else {
			modelVpcPeering.ModifiedAt = statePv.ModifiedAt
			modelVpcPeering.ModifiedBy = statePv.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, modelVpcPeering.AttributeTypes(), modelVpcPeering)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.VpcPeering = mergedObj
		resp.Plan.Set(ctx, plan)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Metadata returns the data source type name.
func (r *vpcPeeringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering"
}

func (r *vpcPeeringResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

// Schema defines the schema for the data source.
func (r *vpcPeeringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC‑to‑VPC peering",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the peering.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
				Description: "The identifier of the account that the approver VPC belongs to.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				Required: true,
			},
			common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
				Description: "The identifier of the approver VPC.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				Required: true,
			},
			common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
				Description: "The identifier of the requester VPC.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the peering.\n" +
					"  - example : peering name\n" +
					"  - Minimum length: 3\n" +
					"  - Maximum length: 20\n" +
					"  - Pattern: ^[a-zA-Z0-9-]*$",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : VPC Peering Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				}},
			common.ToSnakeCase("VpcPeering"): schema.SingleNestedAttribute{
				Description: "VPC‑to‑VPC peering",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountType"): schema.StringAttribute{
						Description: "The type of account.\n" +
							"  - Enum: SAME | DIFFERENT\n" +
							"  - example:SAME",
						Computed: true,
					},
					common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
						Description: "The identifier of the account that the approver VPC belongs to.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
						Description: "The identifier of the approver VPC.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
						Description: "The name of the approver VPC.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
							"  - example : resourceDescription",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the peering.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that modified the resource.\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the peering.\n" +
							"  - example : peering name\n" +
							"  - Minimum length: 3\n" +
							"  - Maximum length: 20\n" +
							"  - Pattern: ^[a-zA-Z0-9-]*$",
						Computed: true,
					},
					common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
						Description: "The identifier of the account that the requester VPC belongs to.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
						Description: "The identifier of the requester VPC.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
						Description: "The name of the requester VPC.\n" +
							"  - example : resourceName",
						Computed: true,
					},
					common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
						Description: "The identifier of account that the deletion requester belongs to.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the peering.\n" +
							"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING\n" +
							"  - example:ACTIVE",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcPeeringResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPeeringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpcv1.VpcPeeringResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpcPeering(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc peering",
			"Could not create vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	plan.Id = types.StringValue(vpcPeering.Id)
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	plan.VpcPeering = vpcObjectValue
	plan.Description = vpcPeeringModel.Description
	plan.ApproverVpcAccountId = vpcPeeringModel.ApproverVpcAccountId
	plan.Name = vpcPeeringModel.Name
	plan.ApproverVpcId = vpcPeeringModel.ApproverVpcId
	plan.RequesterVpcId = vpcPeeringModel.RequesterVpcId

	diags = resp.State.Set(ctx, plan)

	err = waitForVpcPeeringStatus(ctx, r.client, vpcPeering.Id, []string{}, []string{"ACTIVE", "CREATING_REQUESTING"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc peering",
			"Error waiting for vpc peering to become active: "+err.Error(),
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
	//diags = resp.State.Set(ctx, plan)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPeeringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1.VpcPeeringResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc peering",
			"Could not read vpc peering ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	resp.Diagnostics.Append(diags...)

	state.VpcPeering = vpcObjectValue
	state.Description = vpcPeeringModel.Description
	state.ApproverVpcAccountId = vpcPeeringModel.ApproverVpcAccountId
	state.Name = vpcPeeringModel.Name
	state.ApproverVpcId = vpcPeeringModel.ApproverVpcId
	state.RequesterVpcId = vpcPeeringModel.RequesterVpcId

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPeeringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpcv1.VpcPeeringResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpcPeering(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc peering",
			"Could not update vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpcPeering as UpdateVpc items are not populated.
	data, err := r.client.GetVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc peering",
			"Could not read vpc peering ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	resp.Diagnostics.Append(diags...)

	state.VpcPeering = vpcObjectValue
	state.Description = vpcPeeringModel.Description
	state.ApproverVpcAccountId = vpcPeeringModel.ApproverVpcAccountId
	state.Name = vpcPeeringModel.Name
	state.ApproverVpcId = vpcPeeringModel.ApproverVpcId
	state.RequesterVpcId = vpcPeeringModel.RequesterVpcId

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPeeringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1.VpcPeeringResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing VpcPeering
	err := r.client.DeleteVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc peering",
			"Could not delete vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcPeeringStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED", "DELETING_REQUESTING"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc peering",
			"Error waiting for vpc peering to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForVpcPeeringStatus(ctx context.Context, vpcClient *vpcv1.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetVpcPeering(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpcPeering.State), nil
	}, -1, -1, -1, -1)
}
func stringFromNullable(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}
	return types.StringValue(*value)
}
