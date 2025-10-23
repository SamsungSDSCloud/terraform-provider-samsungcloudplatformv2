package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLbListenerResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLbListenerResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewLoadBalancerListenerResource() resource.Resource {
	return &loadbalancerLbListenerResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type loadbalancerLbListenerResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_listener" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Lb Listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LbListener"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Listener.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("InsertClientIp"): schema.BoolAttribute{
						Description: "InsertClientIp",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("Persistence"): schema.StringAttribute{
						Description: "Persistence",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "ServerGroupId",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupName"): schema.StringAttribute{
						Description: "ServerGroupName",
						Optional:    true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "ServicePort",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "ResponseTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "SessionDurationTime",
						Optional:    true,
					},
					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SslCertificate",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "ClientCertId",
								Optional:    true,
							},
							common.ToSnakeCase("ClientCertLevel"): schema.StringAttribute{
								Description: "ClientCertLevel",
								Optional:    true,
							},
							common.ToSnakeCase("ServerCertLevel"): schema.StringAttribute{
								Description: "ServerCertLevel",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("SniCertificate"): schema.ListNestedAttribute{
						Description: "SniCertificate",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("SniCertId"): schema.StringAttribute{
									Description: "SniCertId",
									Optional:    true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "DomainName",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("UrlHandler"): schema.ListNestedAttribute{
						Description: "UrlHandler",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("UrlPattern"): schema.StringAttribute{
									Description: "UrlPattern",
									Optional:    true,
								},
								common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
									Description: "ServerGroupId",
									Optional:    true,
								},
								common.ToSnakeCase("Seq"): schema.Int32Attribute{
									Description: "Seq",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("HttpsRedirection"): schema.SingleNestedAttribute{
						Description: "HttpsRedirection",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "Protocol",
								Optional:    true,
							},
							common.ToSnakeCase("Port"): schema.StringAttribute{
								Description: "Port",
								Optional:    true,
							},
							common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
								Description: "ResponseCode",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "UrlRedirection",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "XForwardedFor",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "XForwardedPort",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "XForwardedProto",
						Optional:    true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "RoutingAction",
						Optional:    true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "ConditionType",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("LbListenerCreate"): schema.SingleNestedAttribute{
				Description: "Create Lb Listener.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("InsertClientIp"): schema.BoolAttribute{
						Description: "InsertClientIp",
						Optional:    true,
					},
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "LoadbalancerId",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("Persistence"): schema.StringAttribute{
						Description: "Persistence",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "ResponseTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "ServerGroupId",
						Optional:    true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "ServicePort",
						Optional:    true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "SessionDurationTime",
						Optional:    true,
					},

					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SslCertificate",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "ClientCertId",
								Optional:    true,
							},
							common.ToSnakeCase("ClientCertLevel"): schema.StringAttribute{
								Description: "ClientCertLevel",
								Optional:    true,
							},
							common.ToSnakeCase("ServerCertLevel"): schema.StringAttribute{
								Description: "ServerCertLevel",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("SniCertificate"): schema.ListNestedAttribute{
						Description: "SniCertificate",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("SniCertId"): schema.StringAttribute{
									Description: "SniCertId",
									Optional:    true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "DomainName",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("UrlHandler"): schema.ListNestedAttribute{
						Description: "UrlHandler",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("UrlPattern"): schema.StringAttribute{
									Description: "UrlPattern",
									Optional:    true,
								},
								common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
									Description: "ServerGroupId",
									Optional:    true,
								},
								common.ToSnakeCase("Seq"): schema.Int32Attribute{
									Description: "Seq",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("HttpsRedirection"): schema.SingleNestedAttribute{
						Description: "HttpsRedirection",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "Protocol",
								Optional:    true,
							},
							common.ToSnakeCase("Port"): schema.StringAttribute{
								Description: "Port",
								Optional:    true,
							},
							common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
								Description: "ResponseCode",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "UrlRedirection",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "XForwardedFor",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "XForwardedPort",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "XForwardedProto",
						Optional:    true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "RoutingAction",
						Optional:    true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "ConditionType",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan loadbalancer.LbListenerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Listener
	data, err := r.client.CreateLbListener(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not create Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Listener.Id)

	// Map response body to schema and populate Computed attribute values
	lbListenerModel := loadbalancerutil.ConvertResponse(data)
	lbListenerOjbectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	plan.LbListener = lbListenerOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbListenerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Listener
	data, err := r.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Lb Listener",
			"Could not read Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbListenerModel := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	state.LbListener = lbListenerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LbListenerResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateLbListener(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not create Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetFirewallRule as UpdateFirewallRule items are not populated.
	data, err := r.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not update Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	lbListenerModel := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	state.LbListener = lbListenerObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLbListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LbListenerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Server Group
	err := r.client.DeleteLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Listener",
			"Could not delete lb listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
