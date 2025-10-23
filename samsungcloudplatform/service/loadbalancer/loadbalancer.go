package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/loadbalancer/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLoadbalancerResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLoadbalancerResource{}
)

// NewLoadBalancerLoadBalancerResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLoadBalancerResource() resource.Resource {
	return &loadbalancerLoadbalancerResource{}
}

// loadbalancerLoadbalancerResource is the data source implementation.
type loadbalancerLoadbalancerResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLoadbalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLoadbalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Loadbalancer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Loadbalancer"): schema.SingleNestedAttribute{
				Description: "A detail of Loadbalancer.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID",
						Computed:    true,
					},
					common.ToSnakeCase("LayerType"): schema.StringAttribute{
						Description: "Layer type",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "Subnet ID",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VPC ID",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("LoadbalancerCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("FirewallEnabled"): schema.BoolAttribute{
						Description: "FirewallEnabled",
						Optional:    true,
					},
					common.ToSnakeCase("FirewallLoggingEnabled"): schema.BoolAttribute{
						Description: "FirewallLoggingEnabled",
						Optional:    true,
					},
					common.ToSnakeCase("LayerType"): schema.StringAttribute{
						Description: "LayerType",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
						Description: "ServiceIp",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Optional:    true,
					},
					common.ToSnakeCase("SourceNatIp"): schema.StringAttribute{
						Description: "SourceNatIp",
						Optional:    true,
					},
					"health_check_ip_1": schema.StringAttribute{
						Description: "HealthCheckIp1",
						Optional:    true,
					},
					"health_check_ip_2": schema.StringAttribute{
						Description: "HealthCheckIp2",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLoadbalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoadBalancer
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLoadbalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LoadbalancerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Loadbalancer
	data, err := r.client.CreateLoadbalancer(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Loadbalancer",
			"Could not create Loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Loadbalancer.Id)

	// Map response body to schema and populate Computed attribute values
	loadbalancerModel := createLoadbalancerModel(data)
	loadbalancerObjectValue, diags := types.ObjectValueFrom(ctx, loadbalancerModel.AttributeTypes(), loadbalancerModel)
	plan.Loadbalancer = loadbalancerObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLoadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LoadbalancerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Loadbalancer
	data, err := r.client.GetLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Loadbalancer",
			"Could not read Loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	loadbalancerModel := createLoadbalancerModelForRead(data)

	loadbalancerObjectValue, diags := types.ObjectValueFrom(ctx, loadbalancerModel.AttributeTypes(), loadbalancerModel)
	state.Loadbalancer = loadbalancerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLoadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LoadbalancerResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateLoadbalancer(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Loadbalancer",
			"Could not update Loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetLoadbalancer as UpdateLoadbalancer items are not populated.
	data, err := r.client.GetLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Loadbalancer",
			"Could not read Loadbalancer ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	loadbalancerModel := showLoadbalancerModel(data)

	loadbalancerObjectValue, diags := types.ObjectValueFrom(ctx, loadbalancerModel.AttributeTypes(), loadbalancerModel)
	state.Loadbalancer = loadbalancerObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLoadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LoadbalancerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Loadbalancer
	err := r.client.DeleteLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Loadbalancer",
			"Could not delete loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createLoadbalancerModel(data *scploadbalancer.LoadbalancerCreateResponse) loadbalancer.LoadbalancerCreateResponseDetail {
	return loadbalancer.LoadbalancerCreateResponseDetail{
		AccountId:   types.StringValue(data.Loadbalancer.AccountId),
		CreatedAt:   types.StringValue(data.Loadbalancer.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(data.Loadbalancer.CreatedBy),
		Description: virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
		Id:          types.StringValue(data.Loadbalancer.Id),
		LayerType:   types.StringValue(data.Loadbalancer.LayerType),
		ModifiedAt:  types.StringValue(data.Loadbalancer.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(data.Loadbalancer.ModifiedBy),
		Name:        types.StringValue(data.Loadbalancer.Name),
		State:       types.StringValue(data.Loadbalancer.State),
		SubnetId:    types.StringValue(data.Loadbalancer.SubnetId),
		VpcId:       types.StringValue(data.Loadbalancer.VpcId),
	}
}

func createLoadbalancerModelForRead(data *scploadbalancer.LoadbalancerShowResponse) loadbalancer.LoadbalancerCreateResponseDetail {
	return loadbalancer.LoadbalancerCreateResponseDetail{
		AccountId:   types.StringValue(data.Loadbalancer.AccountId),
		CreatedAt:   types.StringValue(data.Loadbalancer.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(data.Loadbalancer.CreatedBy),
		Description: virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
		Id:          types.StringValue(data.Loadbalancer.Id),
		LayerType:   types.StringValue(data.Loadbalancer.LayerType),
		ModifiedAt:  types.StringValue(data.Loadbalancer.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(data.Loadbalancer.ModifiedBy),
		Name:        types.StringValue(data.Loadbalancer.Name),
		State:       types.StringValue(data.Loadbalancer.State),
		SubnetId:    types.StringValue(data.Loadbalancer.SubnetId),
		VpcId:       types.StringValue(data.Loadbalancer.VpcId),
	}
}

func showLoadbalancerModel(data *scploadbalancer.LoadbalancerShowResponse) loadbalancer.LoadbalancerDetail {
	return loadbalancer.LoadbalancerDetail{
		AccountId:        types.StringValue(data.Loadbalancer.AccountId),
		CreatedAt:        types.StringValue(data.Loadbalancer.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.Loadbalancer.CreatedBy),
		Description:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
		FirewallId:       virtualserverutil.ToNullableStringValue(data.Loadbalancer.FirewallId.Get()),
		HealthCheckIp:    ToStringList(data.Loadbalancer.HealthCheckIp),
		Id:               types.StringValue(data.Loadbalancer.Id),
		LayerType:        types.StringValue(data.Loadbalancer.LayerType),
		ModifiedAt:       types.StringValue(data.Loadbalancer.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.Loadbalancer.ModifiedBy),
		Name:             types.StringValue(data.Loadbalancer.Name),
		PublicNatEnabled: common.ToNullableBoolValue(data.Loadbalancer.PublicNatEnabled.Get()),
		ServiceIp:        virtualserverutil.ToNullableStringValue(data.Loadbalancer.ServiceIp.Get()),
		SourceNatIp:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.SourceNatIp.Get()),
		State:            types.StringValue(data.Loadbalancer.State),
		SubnetId:         types.StringValue(data.Loadbalancer.SubnetId),
		VpcId:            types.StringValue(data.Loadbalancer.VpcId),
	}
}
