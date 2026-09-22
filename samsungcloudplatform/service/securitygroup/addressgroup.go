package securitygroup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	securitygroupv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	securitygroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/security-group/1.1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = &addressGroupResource{}
	_ resource.ResourceWithConfigure   = &addressGroupResource{}
	_ resource.ResourceWithImportState = &addressGroupResource{}
	_ resource.ResourceWithModifyPlan  = &addressGroupResource{}
)

// NewAddressGroupResource is a helper function to simplify the provider implementation.
func NewAddressGroupResource() resource.Resource {
	return &addressGroupResource{}
}

// addressGroupResource is the resource implementation.
type addressGroupResource struct {
	config  *scpsdk.Configuration
	client  *securitygroupv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *addressGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_address_group"
}

// Schema defines the schema for the resource.
func (r *addressGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Address Group",
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
				Description: "The name of the Address Group.\n" +
					"  - example: ag-web-prod\n" +
					"  - valid: All characters except 'default'\n" +
					"  - constraints: minLength: 1, maxLength: 255, duplicates allowed",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "A brief explanation or note about this resource.\n" +
					"  - example: Address group for web tier\n" +
					"  - constraints: maxLength: 255",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Addresses"): schema.SetAttribute{
				Description: "A list of IP addresses or CIDR ranges in the Address Group.\n" +
					"  - example: [\"192.168.1.0/24\", \"10.0.0.1\"]",
				ElementType: types.StringType,
				Optional:    true,
			},
			common.ToSnakeCase("AddressGroup"): schema.SingleNestedAttribute{
				Description: "Manages address groups for security group rules.",
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
						Description: "The name of the Address Group.\n" +
							"  - example: ag-web-prod",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Address group for web tier",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Addresses"): schema.SetAttribute{
						Description: "A list of IP addresses or CIDR ranges in the Address Group.",
						ElementType: types.StringType,
						Computed:    true,
					},
					common.ToSnakeCase("AddressCount"): schema.Int32Attribute{
						Description: "Number of addresses in the Address Group.\n" +
							"  - example: 5",
						Computed: true,
					},
					common.ToSnakeCase("AddressLimit"): schema.Int32Attribute{
						Description: "Maximum number of addresses allowed in the Address Group.\n" +
							"  - example: 100",
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

// Configure adds the provider configured client to the resource.
func (r *addressGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = inst.Client.SecurityGroupV1d1
	r.clients = inst.Client
}

// buildAddressGroupModel creates the AddressGroup model, address set, and object value from API data.
func buildAddressGroupModel(ctx context.Context, ag securitygroup.AddressGroup, diags *diag.Diagnostics) (securitygroupv1d1.AddressGroup, types.Set, types.Object) {
	agModel := securitygroupv1d1.AddressGroup{
		AccountId:    types.StringValue(ag.AccountId),
		AddressCount: types.Int32Value(ag.GetAddressCount()),
		AddressLimit: types.Int32Value(ag.GetAddressLimit()),
		CreatedAt:    types.StringValue(ag.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(ag.CreatedBy),
		Description:  types.StringPointerValue(ag.Description.Get()),
		Id:           types.StringValue(ag.Id),
		ModifiedAt:   types.StringValue(ag.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(ag.ModifiedBy),
		Name:         types.StringValue(ag.GetName()),
		State:        types.StringValue(ag.GetState()),
	}
	addressLst, dia := types.SetValueFrom(ctx, types.StringType, ag.Addresses)
	diags.Append(dia...)
	agModel.Addresses = addressLst
	agObjectValue, d := types.ObjectValueFrom(ctx, agModel.AttributeTypes(), agModel)
	diags.Append(d...)
	return agModel, addressLst, agObjectValue
}

// extractAddressList converts a types.Set of strings to a []string.
func extractAddressList(addrs types.Set) []string {
	var out []string
	if !addrs.IsNull() && !addrs.IsUnknown() {
		for _, addr := range addrs.Elements() {
			out = append(out, addr.(types.String).ValueString())
		}
	}
	return out
}

// computeAddressDiff returns addresses to add and remove.
func computeAddressDiff(planAddrs, stateAddrs []string) (toAdd []string, toRemove []string) {
	stateMap := make(map[string]struct{}, len(stateAddrs))
	for _, a := range stateAddrs {
		stateMap[a] = struct{}{}
	}
	planMap := make(map[string]struct{}, len(planAddrs))
	for _, a := range planAddrs {
		planMap[a] = struct{}{}
	}
	for _, a := range planAddrs {
		if _, exists := stateMap[a]; !exists {
			toAdd = append(toAdd, a)
		}
	}
	for _, a := range stateAddrs {
		if _, exists := planMap[a]; !exists {
			toRemove = append(toRemove, a)
		}
	}
	return
}

// syncAddressGroupCidrs removes then adds addresses. Returns error on failure.
func (r *addressGroupResource) syncAddressGroupCidrs(ctx context.Context, addressId string, toAdd, toRemove []string, diags *diag.Diagnostics) error {
	if len(toRemove) > 0 {
		removeReq := securitygroupv1d1.AddressGroupCidrResource{}
		addrSet, d := types.SetValueFrom(ctx, types.StringType, toRemove)
		diags.Append(d...)
		if diags.HasError() {
			return fmt.Errorf("diagnostic error preparing remove request")
		}
		removeReq.Addresses = addrSet
		_, err := r.client.RemoveAddressGroupCidrs(ctx, addressId, removeReq)
		if err != nil {
			return fmt.Errorf("removing addresses: %w", err)
		}
	}
	if len(toAdd) > 0 {
		addReq := securitygroupv1d1.AddressGroupCidrResource{}
		addrSet, d := types.SetValueFrom(ctx, types.StringType, toAdd)
		diags.Append(d...)
		if diags.HasError() {
			return fmt.Errorf("diagnostic error preparing add request")
		}
		addReq.Addresses = addrSet
		_, err := r.client.AddAddressGroupCidrs(ctx, addressId, addReq)
		if err != nil {
			return fmt.Errorf("adding addresses: %w", err)
		}
	}
	return nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *addressGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securitygroupv1d1.AddressGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateAddressGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating address group",
			"Could not create address group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil || data.AddressGroup.Id == "" {
		resp.Diagnostics.AddError(
			"Error creating address group",
			"empty response from API",
		)
		return
	}
	addressGroup := data.AddressGroup
	plan.Id = types.StringValue(addressGroup.Id)

	agModel, addresses, agObjectValue := buildAddressGroupModel(ctx, addressGroup, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.Addresses.IsNull() {
		// Avoid inconsistence data when creation
		plan.Addresses = addresses
	}
	plan.AddressGroup = agObjectValue
	plan.Description = agModel.Description

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *addressGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securitygroupv1d1.AddressGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetAddressGroup(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Could not read address group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Empty response from API",
		)
		return
	}
	addressGroup := data.AddressGroup

	_, addressLst, agObjectValue := buildAddressGroupModel(ctx, addressGroup, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Name = types.StringValue(addressGroup.GetName())
	state.Description = types.StringPointerValue(addressGroup.Description.Get())
	state.Addresses = addressLst
	state.AddressGroup = agObjectValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *addressGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan securitygroupv1d1.AddressGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state securitygroupv1d1.AddressGroupResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compare plan addresses vs state addresses to determine adds and removes
	addressId := plan.Id.ValueString()

	planAddrSet := extractAddressList(plan.Addresses)
	stateAddrSet := extractAddressList(state.Addresses)
	toAdd, toRemove := computeAddressDiff(planAddrSet, stateAddrSet)

	if err := r.syncAddressGroupCidrs(ctx, addressId, toAdd, toRemove, &resp.Diagnostics); err != nil {
		detail := client.GetDetailFromError(err)
		action := err.Error()
		resp.Diagnostics.AddError(
			"Error syncing addresses for address group",
			"Could not sync addresses for address group ID "+addressId+": "+action+"\nReason: "+detail,
		)
		return
	}

	// Update description
	err := r.client.UpdateAddressGroup(ctx, addressId, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating address group",
			"Could not update address group ID "+addressId+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetAddressGroup(ctx, addressId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Could not read address group ID "+addressId+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Empty response from API",
		)
		return
	}

	agModel, addressLst, agObjectValue := buildAddressGroupModel(ctx, data.AddressGroup, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AddressGroup = agObjectValue
	state.Addresses = addressLst
	state.Description = agModel.Description

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *addressGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securitygroupv1d1.AddressGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAddressGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting address group",
			"Could not delete address group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *addressGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *addressGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan securitygroupv1d1.AddressGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state securitygroupv1d1.AddressGroupResource
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
	if resp.Diagnostics.HasError() {
		return
	}

	// Reconstruct address_group: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_id, created_at, created_by from showing as (known after apply).
	if !state.AddressGroup.IsNull() && !state.AddressGroup.IsUnknown() {
		var stateAg securitygroupv1d1.AddressGroup
		resp.Diagnostics.Append(state.AddressGroup.As(ctx, &stateAg, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when description or addresses actually change
		descriptionChanged := !plan.Description.Equal(state.Description)
		addressesChanged := !plan.Addresses.Equal(state.Addresses)

		mergedAg := securitygroupv1d1.AddressGroup{
			Id:           stateAg.Id,
			AccountId:    stateAg.AccountId,
			CreatedAt:    stateAg.CreatedAt,
			CreatedBy:    stateAg.CreatedBy,
			AddressCount: stateAg.AddressCount,
			AddressLimit: stateAg.AddressLimit,
			Addresses:    plan.Addresses,
			Name:         stateAg.Name,
			State:        stateAg.State,
			Description:  plan.Description,
		}

		if descriptionChanged || addressesChanged {
			mergedAg.ModifiedAt = types.StringUnknown()
			mergedAg.ModifiedBy = types.StringUnknown()
			if addressesChanged {
				mergedAg.AddressCount = types.Int32Unknown()
			}
		} else {
			mergedAg.ModifiedAt = stateAg.ModifiedAt
			mergedAg.ModifiedBy = stateAg.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, mergedAg.AttributeTypes(), mergedAg)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.AddressGroup = mergedObj
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
