package gslb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	gslb "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpgslb "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/gslb/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

const reasonPrefix = "\nReason: "

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &gslbGslbResource{}
	_ resource.ResourceWithConfigure = &gslbGslbResource{}
	_ resource.ResourceWithImportState = &gslbGslbResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewGslbGslbResource() resource.Resource {
	return &gslbGslbResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type gslbGslbResource struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *gslbGslbResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *gslbGslbResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Global Server Load Balancer resource for distributing traffic across multiple regions.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the GSLB.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),
			common.ToSnakeCase("Gslb"): schema.SingleNestedAttribute{
				Description: "Details of the Global Server Load Balancer.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "The load balancing algorithm for GSLB traffic distribution (e.g., ROUND_ROBIN, RATIO).\n" +
							"  - example : ROUND_ROBIN",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Example Description for GSLB",
						Computed: true,
					},
					common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
						Description: "The environment usage type for the GSLB (e.g., PUBLIC).\n" +
							"  - example : PUBLIC",
						Computed: true,
					},
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "Health check configuration for monitoring GSLB endpoint availability.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
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
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "The GSLB Health Check Interval.\n" +
									"  - example : 30\n" +
									"  - Range: 5 to 299",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Probe Timeout.\n" +
									"  - example : 10\n" +
									"  - Range: 5 to 300",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "The GSLB Health Check User Name.\n" +
									"  - example : healthcheck_user\n" +
									"  - Max length: 60",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "The GSLB Health Check Password.\n" +
									"  - example : **********",
								Computed: true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "The unique identifier of the health check configuration.\n" +
									"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
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
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The protocol used for health checks (e.g., ICMP, TCP, HTTP, HTTPS).\n" +
									"  - example : TCP",
								Computed: true,
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "The GSLB Health Check Receive String.\n" +
									"  - example : HTTP/1.1 200 OK\n" +
									"  - Max length: 300",
								Computed: true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "The GSLB Health Check Send String. If no input is provided, it operates as a \"GET /\" request.\n" +
									"  - example : GET /",
								Computed: true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "The GSLB Health Check Service Port.\n" +
									"  - example : 80\n" +
									"  - Range: 1 to 65535",
								Computed: true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Timeout. It must be greater than the Interval.\n" +
									"  - example : 40\n" +
									"  - Range: 6 to 300",
								Computed: true,
							},
						},
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the GSLB.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed: true,
					},
					common.ToSnakeCase("LinkedResourceCount"): schema.Int32Attribute{
						Description: "The number of resources linked to this GSLB.\n" +
							"  - example : 2",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the GSLB.\n" +
							"  - example : example.gslb.e.samsungsdscloud.com",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the GSLB (e.g., ACTIVE, CREATING, EDITING, ERROR, DELETING).\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
				},
			},
			common.ToSnakeCase("GslbCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new GSLB.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "The load balancing algorithm for GSLB traffic distribution (e.g., ROUND_ROBIN, RATIO).\n" +
							"  - example : ROUND_ROBIN",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ROUND_ROBIN", "RATIO"),
						},
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Example Description for GSLB",
						Optional: true,
					},
					common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
						Description: "The environment usage type for the GSLB (e.g., PUBLIC).\n" +
							"  - example : PUBLIC",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("PUBLIC"),
						},
					},
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "Health check configuration for monitoring GSLB endpoint availability.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "The GSLB Health Check Interval.\n" +
									"  - example : 30\n" +
									"  - Range: 5 to 299",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Probe Timeout.\n" +
									"  - example : 10\n" +
									"  - Range: 5 to 300",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "The GSLB Health Check User Name.\n" +
									"  - example : healthcheck_user\n" +
									"  - Max Length: 60",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "The GSLB Health Check Password. If the User name is entered, This value is required.\n" +
									"  - example : **********\n" +
									"  - maxLength: 20\n" +
									"  - minLength: 8\n" +
									"  - pattern: ^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@!%*#?&])[A-Za-z\\d$@!%*#?&]$",
								Optional: true,
							},
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The protocol used for health checks (e.g., ICMP, TCP, HTTP, HTTPS).\n" +
									"  - example : TCP",
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("ICMP", "TCP", "HTTP", "HTTPS"),
								},
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "The GSLB Health Check Receive String.\n" +
									"  - example : HTTP/1.1 200 OK\n" +
									"  - Max Length: 300",
								Optional: true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "The GSLB Health Check Send String. If no input is provided, it operates as a \"GET /\" request.\n" +
									"  - example : GET /",
								Optional: true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "The GSLB Health Check Service Port.\n" +
									"  - example : 80\n" +
									"  - Range: 1 to 65535",
								Optional: true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Timeout. It must be greater than the Interval.\n" +
									"  - example : 40\n" +
									"  - Range: 6 to 300",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the GSLB.\n" +
							"  - example : example.gslb.e.samsungsdscloud.com",
						Required: true,
					},
					common.ToSnakeCase("Resources"): schema.ListNestedAttribute{
						Description: "The list of resources for the GSLB.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Description"): schema.StringAttribute{
									Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
										"  - example : Example Description for GSLB Resource",
									Optional: true,
								},
								common.ToSnakeCase("Destination"): schema.StringAttribute{
									Description: "The destination endpoint for the GSLB resource.\n" +
										"  - example : 192.168.1.100",
									Optional: true,
								},
								common.ToSnakeCase("Region"): schema.StringAttribute{
									Description: "The region where the GSLB resource is located.\n" +
										"  - example : kr-west1",
									Optional: true,
								},
								common.ToSnakeCase("Weight"): schema.Int32Attribute{
									Description: "The weight for load balancing distribution (0-100).\n" +
										"  - example : 50",
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *gslbGslbResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Gslb
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *gslbGslbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan gslb.GslbResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateGslb(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Gslb",
			"Could not create Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	err = waitForGslbStatus(ctx, r.client, data.Gslb.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating gslb",
			"Error waiting for gslb to become active: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(data.Gslb.Id)
	data, err = r.client.GetGslb(ctx, data.Gslb.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Could not read Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	if data == nil || data.Gslb.Id == "" {
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Gslb response is nil or empty",
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Gslb = gslbObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *gslbGslbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state gslb.GslbResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetGslb(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Could not read Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Gslb = gslbObjectValue

	// Reconstruct gslb_create from API response so that:
	// 1) Import populates gslb_create (not null)
	// 2) External changes (console/API) are detected as drift
	if state.GslbCreate == nil {
		state.GslbCreate = &gslb.GslbCreate{}
	}

	state.GslbCreate.Name = types.StringValue(data.Gslb.Name)
	state.GslbCreate.Algorithm = types.StringValue(data.Gslb.Algorithm)
	state.GslbCreate.Description = virtualserverutil.ToNullableStringValue(data.Gslb.Description.Get())
	state.GslbCreate.EnvUsage = types.StringValue(data.Gslb.EnvUsage)

	// Rebuild health_check from API response
	healthCheckFromData := data.Gslb.HealthCheck.Get()
	if healthCheckFromData != nil {
		state.GslbCreate.HealthCheck = &gslb.HealthCheckCreate{
			HealthCheckInterval:     types.Int32Value(healthCheckFromData.GetHealthCheckInterval()),
			HealthCheckProbeTimeout: types.Int32Value(healthCheckFromData.GetHealthCheckProbeTimeout()),
			HealthCheckUserId:       types.StringValue(healthCheckFromData.GetHealthCheckUserId()),
			HealthCheckUserPassword: types.StringValue(healthCheckFromData.GetHealthCheckUserPassword()),
			Protocol:                types.StringValue(healthCheckFromData.Protocol),
			ReceiveString:           types.StringValue(healthCheckFromData.GetReceiveString()),
			SendString:              types.StringValue(healthCheckFromData.GetSendString()),
			ServicePort:             types.Int32Value(healthCheckFromData.GetServicePort()),
			Timeout:                 types.Int32Value(healthCheckFromData.GetTimeout()),
		}
	}

	// Rebuild resources from API response.
	// Only overwrite state when the API returns valid data to avoid
	// destroy/recreate on Read round-trip failure (GSLB-FIX-05).
	resourceList, err := r.client.GetGslbResourceList(ctx, gslb.GslbResourceDataSource{GslbId: state.Id})
	if err == nil && resourceList != nil && len(resourceList.GslbResources) > 0 {
		resources := make([]gslb.GslbResourceCreate, 0, len(resourceList.GslbResources))
		for _, res := range resourceList.GslbResources {
			resources = append(resources, gslb.GslbResourceCreate{
				Description: virtualserverutil.ToNullableStringValue(res.Description.Get()),
				Destination: types.StringValue(res.Destination),
				Region:      types.StringValue(res.Region),
				Weight:      common.ToNullableInt32Value(res.Weight.Get()),
			})
		}
		state.GslbCreate.Resources = resources
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *gslbGslbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var oldState gslb.GslbResource
	var state gslb.GslbResource
	req.State.Get(ctx, &oldState)
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if gslbResourceChanged(oldState, state) && gslbHealthCheckChanged(oldState, state) {
		resp.Diagnostics.AddError(
			"Error updating Gslb",
			"Could not change GSLB resources and health checks at the same time",
		)
		return
	}

	// Update existing order
	if gslbChanged(oldState, state) {
		_, err := r.client.UpdateGslb(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	if gslbResourceChanged(oldState, state) {
		_, err := r.client.UpdateGslbResource(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	if gslbHealthCheckChanged(oldState, state) {
		_, err := r.client.UpdateGslbHealthCheck(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	updateErr := waitForGslbStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating Gslb",
			"Error updating for Gslb to become active: "+updateErr.Error(),
		)
		return
	}

	data, err := r.client.GetGslb(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Could not read Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Gslb = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *gslbGslbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state gslb.GslbResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Gslb
	data, err := r.client.DeleteGslb(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Gslb",
			"Could not delete Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Gslb = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertResponseToGslb(data *scpgslb.GslbShowResponse) gslb.GslbDetail {
	var healthCheck *gslb.HealthCheck
	healthCheckFromData := data.Gslb.HealthCheck.Get()
	if healthCheckFromData != nil {
		healthCheck = &gslb.HealthCheck{
			CreatedAt:               types.StringValue(healthCheckFromData.GetCreatedAt().Format(time.RFC3339)),
			CreatedBy:               types.StringValue(healthCheckFromData.CreatedBy),
			HealthCheckInterval:     types.Int32Value(healthCheckFromData.GetHealthCheckInterval()),
			HealthCheckProbeTimeout: types.Int32Value(healthCheckFromData.GetHealthCheckProbeTimeout()),
			HealthCheckUserId:       types.StringValue(healthCheckFromData.GetHealthCheckUserId()),
			HealthCheckUserPassword: types.StringValue(healthCheckFromData.GetHealthCheckUserPassword()),
			Id:                      types.StringValue(healthCheckFromData.Id),
			ModifiedAt:              types.StringValue(healthCheckFromData.GetModifiedAt().Format(time.RFC3339)),
			ModifiedBy:              types.StringValue(healthCheckFromData.ModifiedBy),
			Protocol:                types.StringValue(healthCheckFromData.Protocol),
			ReceiveString:           types.StringValue(healthCheckFromData.GetReceiveString()),
			SendString:              types.StringValue(healthCheckFromData.GetSendString()),
			ServicePort:             types.Int32Value(healthCheckFromData.GetServicePort()),
			Timeout:                 types.Int32Value(healthCheckFromData.GetTimeout()),
		}
	}
	return gslb.GslbDetail{
		Algorithm:           types.StringValue(data.Gslb.Algorithm),
		CreatedAt:           types.StringValue(data.Gslb.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(data.Gslb.CreatedBy),
		Description:         virtualserverutil.ToNullableStringValue(data.Gslb.Description.Get()),
		EnvUsage:            types.StringValue(data.Gslb.EnvUsage),
		HealthCheck:         healthCheck,
		Id:                  types.StringValue(data.Gslb.Id),
		LinkedResourceCount: types.Int32Value(data.Gslb.LinkedResourceCount),
		ModifiedAt:          types.StringValue(data.Gslb.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(data.Gslb.ModifiedBy),
		Name:                types.StringValue(data.Gslb.Name),
		State:               types.StringValue(data.Gslb.State),
	}
}

func gslbChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	if !oldState.GslbCreate.Algorithm.Equal(newState.GslbCreate.Algorithm) {
		return true
	}
	if !oldState.GslbCreate.Description.Equal(newState.GslbCreate.Description) {
		return true
	}
	return false
}

func gslbResourceChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	oldResources := oldState.GslbCreate.Resources
	newResources := newState.GslbCreate.Resources

	// Build key-based maps for order-independent comparison.
	// Key = "Destination/Region" uniquely identifies a GSLB resource endpoint.
	key := func(r gslb.GslbResourceCreate) string {
		return r.Destination.ValueString() + "/" + r.Region.ValueString()
	}

	oldByKey := make(map[string]gslb.GslbResourceCreate, len(oldResources))
	for _, r := range oldResources {
		oldByKey[key(r)] = r
	}

	newByKey := make(map[string]gslb.GslbResourceCreate, len(newResources))
	for _, r := range newResources {
		newByKey[key(r)] = r
	}

	if len(oldByKey) != len(newByKey) {
		return true
	}

	for k, oldRes := range oldByKey {
		newRes, ok := newByKey[k]
		if !ok {
			return true
		}
		if !oldRes.Description.Equal(newRes.Description) {
			return true
		}
		if !oldRes.Weight.Equal(newRes.Weight) {
			return true
		}
	}

	return false
}

func gslbHealthCheckChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	oldHealthCheck := oldState.GslbCreate.HealthCheck
	newHealthCheck := newState.GslbCreate.HealthCheck

	if oldHealthCheck == nil && newHealthCheck == nil {
		return false
	}
	if oldHealthCheck == nil || newHealthCheck == nil {
		return true
	}

	if oldHealthCheck.HealthCheckInterval != newHealthCheck.HealthCheckInterval {
		return true
	}
	if oldHealthCheck.HealthCheckProbeTimeout != newHealthCheck.HealthCheckProbeTimeout {
		return true
	}
	if oldHealthCheck.HealthCheckUserId != newHealthCheck.HealthCheckUserId {
		return true
	}
	if oldHealthCheck.HealthCheckUserPassword != newHealthCheck.HealthCheckUserPassword {
		return true
	}
	if oldHealthCheck.Protocol != newHealthCheck.Protocol {
		return true
	}
	if oldHealthCheck.ReceiveString != newHealthCheck.ReceiveString {
		return true
	}
	if oldHealthCheck.SendString != newHealthCheck.SendString {
		return true
	}
	if oldHealthCheck.ServicePort != newHealthCheck.ServicePort {
		return true
	}
	if oldHealthCheck.Timeout != newHealthCheck.Timeout {
		return true
	}

	return false
}

func waitForGslbStatus(ctx context.Context, gslbClient *gslb.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := gslbClient.GetGslb(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, info.Gslb.State, nil
	}, -1, -1, -1, -1)
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *gslbGslbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
   resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}