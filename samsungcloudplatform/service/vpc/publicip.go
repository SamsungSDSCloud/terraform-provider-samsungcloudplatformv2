package vpc

import (
	"context"
	"fmt"
	"strings"

	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpvpcv1d3 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &vpcPublicipResource{}
	_ resource.ResourceWithConfigure   = &vpcPublicipResource{}
	_ resource.ResourceWithImportState = &vpcPublicipResource{}
	_ resource.ResourceWithModifyPlan  = &vpcPublicipResource{}
)

// NewVpcPublicipResource is a helper function to simplify the provider implementation.
func NewVpcPublicipResource() resource.Resource {
	return &vpcPublicipResource{}
}

// vpcPublicipResource is the data source implementation.
type vpcPublicipResource struct {
	config     *scpsdk.Configuration
	clientv1d3 *vpcv1d3.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcPublicipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_publicip"
}

// Schema defines the schema for the data source.
func (r *vpcPublicipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "publicip",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the public ip.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type of the public ip.\n" +
					"  - example : IGW | GGW | SIGW",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : Public IP Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Zone"): schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Zone\n  - example: kr-west1-a",
				MarkdownDescription: "Zone\n  - example: kr-west1-a",
			},
			common.ToSnakeCase("Publicip"): schema.SingleNestedAttribute{
				Description: "Publicip",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the public ip.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "The IP address assigned to the resource.\n" +
							"  - example : 203.0.113.10",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the public ip.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
						Description: "The type of the resource that this public ip is attached to.\n" +
							"  - example : VM",
						Computed: true,
					},
					common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
						Description: "The name of the resource that this public ip is attached to.\n" +
							"  - example : my-server",
						Computed: true,
					},
					common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
						Description: "The identifier of the resource that this public ip is attached to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of the public ip.\n" +
							"  - example : IGW",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the public ip.\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Public IP Description",
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
					common.ToSnakeCase("Zone"): schema.StringAttribute{
						Computed:    true,
						Description: "Zone\n  - example: kr-west1-a",
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcPublicipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clientv1d3 = inst.Client.VpcV1Dot3
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPublicipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpcv1d3.PublicipResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new publicip
	data, err := r.clientv1d3.CreatePublicip(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating publicip",
			"Could not create publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Publicip.Id)
	publicipModel := createPublicipModelV1Dot3(data)
	publicipObjectValue, diags := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Publicip = publicipObjectValue
	// Description might be default to "" in API response
	plan.Description = publicipModel.Description
	plan.Type = publicipModel.Type
	plan.Zone = publicipModel.Zone

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPublicipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1d3.PublicipResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from publicip
	data, err := r.clientv1d3.GetPublicip(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading publicip",
			"Could not read publicip ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
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

	publicipModel := createPublicipModelV1Dot3(data)
	publicipObjectValue, diag := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Publicip = publicipObjectValue
	// Refresh input attributes
	state.Description = publicipModel.Description
	state.Type = publicipModel.Type
	state.Zone = publicipModel.Zone

	tagsMap := getPublicipTags(r.clients, state.Id.ValueString())
	state.Tags = tagsMap

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPublicipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var plan vpcv1d3.PublicipResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.clientv1d3.UpdatePublicip(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating publicip",
			"Could not update publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated value from GetPublicip as UpdatePublicip response is limited.
	data, err := r.clientv1d3.GetPublicip(ctx, plan.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading publicip",
			"Could not read publicip ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicipModel := createPublicipModelV1Dot3(data)
	publicipObjectValue, diag := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Publicip = publicipObjectValue
	plan.Description = publicipModel.Description
	plan.Type = publicipModel.Type
	plan.Zone = publicipModel.Zone

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPublicipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1d3.PublicipResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing publicip
	err := r.clientv1d3.DeletePublicip(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting publicip",
			"Could not delete publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createPublicipModelV1Dot3(data *scpvpcv1d3.PublicipShowResponseV1Dot3) vpcv1d3.Publicip {
	publicip := data.Publicip
	publicipModel := vpcv1d3.Publicip{
		Id:                   types.StringValue(data.GetPublicip().Id),
		IpAddress:            types.StringValue(data.GetPublicip().IpAddress),
		AccountId:            types.StringValue(data.GetPublicip().AccountId),
		AttachedResourceName: types.StringPointerValue(data.GetPublicip().AttachedResourceName.Get()),
		AttachedResourceId:   types.StringPointerValue(data.GetPublicip().AttachedResourceId.Get()),
		Type:                 types.StringValue(string(data.GetPublicip().Type)),
		State:                types.StringValue(string(data.GetPublicip().State)),
		Description:          types.StringPointerValue(data.GetPublicip().Description.Get()),
		CreatedAt:            types.StringValue(data.GetPublicip().CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.GetPublicip().CreatedBy),
		ModifiedAt:           types.StringValue(data.GetPublicip().ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.GetPublicip().ModifiedBy),
		Zone:                 types.StringValue(data.GetPublicip().Zone),
	}
	attachedResourceType := publicip.AttachedResourceType.Get()
	if attachedResourceType != nil {
		attachedResourceTypeStr := string(*attachedResourceType)
		publicipModel.AttachedResourceType = types.StringPointerValue(&attachedResourceTypeStr)
	} else {
		publicipModel.AttachedResourceType = types.StringPointerValue(nil)
	}
	return publicipModel
}

func getPublicipTags(clients *client.SCPClient, resourceIdentifier string) types.Map {
	serviceName := "vpc"
	resourceType := "publicip"

	srn, err := tag.GetSRN(clients, serviceName, resourceType, resourceIdentifier, false)
	if err != nil {
		return types.MapNull(types.StringType)
	}

	resp, err := clients.ResourceManager.GetResourceTags(srn)
	if err != nil {
		return types.MapNull(types.StringType)
	}

	tags := make(map[string]attr.Value)
	for _, t := range resp.Content.Tags {
		tags[t.Key] = types.StringValue(*t.Value.Get())
	}

	if len(tags) == 0 {
		return types.MapNull(types.StringType)
	}

	tagsMap, _ := types.MapValue(types.StringType, tags)
	return tagsMap
}

func (r *vpcPublicipResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan vpcv1d3.PublicipResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state vpcv1d3.PublicipResource
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
		{"type", !plan.Type.Equal(state.Type)},
		{"zone", !plan.Zone.Equal(state.Zone)},
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

	// Reconstruct publicip: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_id, created_at, created_by from showing as (known after apply).
	if !state.Publicip.IsNull() && !state.Publicip.IsUnknown() {
		var statePi vpcv1d3.Publicip
		resp.Diagnostics.Append(state.Publicip.As(ctx, &statePi, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when description actually changes
		descriptionChanged := !plan.Description.Equal(state.Description)

		mergedPi := vpcv1d3.Publicip{
			Id:                   statePi.Id,
			IpAddress:            statePi.IpAddress,
			AccountId:            statePi.AccountId,
			AttachedResourceType: statePi.AttachedResourceType,
			AttachedResourceName: statePi.AttachedResourceName,
			AttachedResourceId:   statePi.AttachedResourceId,
			Type:                 statePi.Type,
			State:                statePi.State,
			CreatedAt:            statePi.CreatedAt,
			CreatedBy:            statePi.CreatedBy,
			Zone:                 statePi.Zone,
			Description:          plan.Description,
		}

		if descriptionChanged {
			mergedPi.ModifiedAt = types.StringUnknown()
			mergedPi.ModifiedBy = types.StringUnknown()
		} else {
			mergedPi.ModifiedAt = statePi.ModifiedAt
			mergedPi.ModifiedBy = statePi.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, mergedPi.AttributeTypes(), mergedPi)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Publicip = mergedObj
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// ImportState imports an existing resource into Terraform state.
func (r *vpcPublicipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))
}
