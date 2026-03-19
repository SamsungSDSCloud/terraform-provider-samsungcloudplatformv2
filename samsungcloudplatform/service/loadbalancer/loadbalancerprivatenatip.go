package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/loadbalancer/1.3"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLoadbalancerPrivateNatIpResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLoadbalancerPrivateNatIpResource{}
)

// loadbalancerLoadbalancerPrivateNatIpResource is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerPrivateNatIpResource() resource.Resource {
	return &loadbalancerLoadbalancerPrivateNatIpResource{}
}

// loadbalancerLoadbalancerPrivateNatIpResource is the data source implementation.
type loadbalancerLoadbalancerPrivateNatIpResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLoadbalancerPrivateNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer_private_nat_ip"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLoadbalancerPrivateNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Loadbalancer Private NAT.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "LoadbalancerId",
				Required:    true,
			},
			common.ToSnakeCase("LoadbalancerPrivateNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of private NAT.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Optional:    true,
					},
					common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
						Description: "ExternalIpAddress",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("InternalIpAddress"): schema.StringAttribute{
						Description: "InternalIpAddress",
						Optional:    true,
					},
					common.ToSnakeCase("PrivateNatIpId"): schema.StringAttribute{
						Description: "PrivateNatIpId",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("PrivateStaticNatCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer private static NAT.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
						Description: "PrivateNatId",
						Optional:    true,
					},
					common.ToSnakeCase("PrivateNatIpId"): schema.StringAttribute{
						Description: "PrivateNatIpId",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLoadbalancerPrivateNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *loadbalancerLoadbalancerPrivateNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LoadbalancerPrivateNatIpResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Static NAT
	data, err := r.client.CreateLoadbalancerPrivateNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating LoadBalancer Private NAT",
			"Could not create LoadBalancer Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = loadbalancerutil.ToNullableStringValue(data.StaticNat.Id.Get())

	// Map response body to schema and populate Computed attribute values
	staticNatModel := createLoadbalancerPrivateNatModel(data)
	staticNatObjectValue, diags := types.ObjectValueFrom(ctx, staticNatModel.AttributeTypes(), staticNatModel)
	plan.LoadbalancerPrivateNatIp = staticNatObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerLoadbalancerPrivateNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read refreshes the Terraform state with the latest data.
}

func (r *loadbalancerLoadbalancerPrivateNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update updates the resource and sets the updated Terraform state on success.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLoadbalancerPrivateNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LoadbalancerPrivateNatIpResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Static NAT
	err := r.client.DeleteLoadbalancerPrivateNatIp(ctx, state.LoadbalancerId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB PrivateNAT",
			"Could not delete LB Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createLoadbalancerPrivateNatModel(data *scploadbalancer.PrivateStaticNatCreateResponse) loadbalancer.LoadbalancerPrivateNatIpDetail {
	lbStaticNat := data.StaticNat
	return loadbalancer.LoadbalancerPrivateNatIpDetail{
		CreatedAt:         types.StringValue(lbStaticNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(lbStaticNat.CreatedBy),
		ExternalIpAddress: loadbalancerutil.ToNullableStringValue(lbStaticNat.ExternalIpAddress.Get()),
		Id:                loadbalancerutil.ToNullableStringValue(lbStaticNat.Id.Get()),
		InternalIpAddress: loadbalancerutil.ToNullableStringValue(lbStaticNat.InternalIpAddress.Get()),
		ModifiedAt:        types.StringValue(lbStaticNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(lbStaticNat.ModifiedBy),
		PrivateNatIpId:    loadbalancerutil.ToNullableStringValue(lbStaticNat.PrivateNatIpId.Get()),
		State:             loadbalancerutil.ToNullableStringValue(lbStaticNat.State.Get()),
	}
}
