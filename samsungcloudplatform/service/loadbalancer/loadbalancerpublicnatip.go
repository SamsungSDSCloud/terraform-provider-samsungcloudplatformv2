package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/loadbalancer/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLoadbalancerPublicNatIpResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLoadbalancerPublicNatIpResource{}
)

// NewLoadbalancerLoadbalancerPublicNatIpResource is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerPublicNatIpResource() resource.Resource {
	return &loadbalancerLoadbalancerPublicNatIpResource{}
}

// loadbalancerLoadbalancerPublicNatIpResource is the data source implementation.
type loadbalancerLoadbalancerPublicNatIpResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer_public_nat_ip"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Loadbalancer Public NAT.",
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
			common.ToSnakeCase("LoadbalancerPublicNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of public NAT.",
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
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Optional:    true,
					},
					common.ToSnakeCase("ActionType"): schema.StringAttribute{
						Description: "ActionType",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerId"): schema.StringAttribute{
						Description: "OwnerId",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerName"): schema.StringAttribute{
						Description: "OwnerName",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerType"): schema.StringAttribute{
						Description: "OwnerType",
						Optional:    true,
					},
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "PublicipId",
						Optional:    true,
					},
					common.ToSnakeCase("ServiceIpPortId"): schema.StringAttribute{
						Description: "ServiceIpPortId",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
					common.ToSnakeCase("vpc_id"): schema.StringAttribute{
						Description: "vpc_id",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("StaticNatCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer static NAT.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "PublicipId",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *loadbalancerLoadbalancerPublicNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LoadbalancerPublicNatIpResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Static NAT
	data, err := r.client.CreateLoadbalancerPublicNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Public NAT",
			"Could not create Public NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.StaticNat.Id)

	// Map response body to schema and populate Computed attribute values
	staticNatModel := createLoadbalancerNatModel(data)
	staticNatObjectValue, diags := types.ObjectValueFrom(ctx, staticNatModel.AttributeTypes(), staticNatModel)
	plan.LoadbalancerPublicNatIp = staticNatObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LoadbalancerPublicNatIpResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Static NAT
	err := r.client.DeleteLoadbalancerPublicNatIp(ctx, state.LoadbalancerId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Public NAT",
			"Could not delete LB Public NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createLoadbalancerNatModel(data *scploadbalancer.StaticNatCreateResponse) loadbalancer.LoadbalancerPublicNatIpDetail {
	lbStaticNat := data.StaticNat
	return loadbalancer.LoadbalancerPublicNatIpDetail{
		AccountId:         types.StringValue(lbStaticNat.AccountId),
		ActionType:        types.StringValue(lbStaticNat.ActionType),
		CreatedAt:         types.StringValue(lbStaticNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(lbStaticNat.CreatedBy),
		Description:       types.StringValue(lbStaticNat.Description),
		ExternalIpAddress: types.StringValue(lbStaticNat.ExternalIpAddress),
		Id:                types.StringValue(lbStaticNat.Id),
		InternalIpAddress: types.StringValue(lbStaticNat.InternalIpAddress),
		ModifiedAt:        types.StringValue(lbStaticNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(lbStaticNat.ModifiedBy),
		Name:              types.StringValue(lbStaticNat.Name),
		OwnerId:           types.StringValue(lbStaticNat.OwnerId),
		OwnerName:         types.StringValue(lbStaticNat.OwnerName),
		OwnerType:         types.StringValue(lbStaticNat.OwnerType),
		PublicipId:        types.StringValue(lbStaticNat.PublicipId),
		ServiceIpPortId:   types.StringValue(lbStaticNat.ServiceIpPortId),
		State:             types.StringValue(lbStaticNat.State),
		SubnetId:          types.StringValue(lbStaticNat.AccountId),
		Type:              types.StringValue(lbStaticNat.Type),
		VpcId:             types.StringValue(lbStaticNat.VpcId),
	}
}
