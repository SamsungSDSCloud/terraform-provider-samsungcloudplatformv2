package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancerv1d4"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scploadbalancerv1d4 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loadbalancerLbHealthCheckResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLbHealthCheckResource{}
	_ resource.ResourceWithImportState = &loadbalancerLbHealthCheckResource{}
)

// NewLoadBalancerLbHealthCheckResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLbHealthCheckResource() resource.Resource {
	return &loadbalancerLbHealthCheckResource{}
}

// loadbalancerLbHealthCheckResource is the data source implementation.
type loadbalancerLbHealthCheckResource struct {
	config     *scpsdk.Configuration
	clientv1d4 *loadbalancerv1d4.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbHealthCheckResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_health_check"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbHealthCheckResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "LB Health Check resource for monitoring server health.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LbHealthCheck"): schema.SingleNestedAttribute{
				Description: "Details of the LB Health Check.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : this is an lb server group\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 8acceeb6920c4fc494490d864f67f0b5\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 60fba45cb6c811efba41ba92e4fe7200\n",
						Optional: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol used for the health check.\n" +
							"  - example : TCP\n" +
							"  - pattern : TCP | HTTP | HTTPS\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Health Check.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Health Check.\n" +
							"  - example : ServerGroup01\n" +
							"  - minLength : 3\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9][-a-zA-Z0-9_]*[a-zA-Z0-9]$\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "The port number used for health checks.\n" +
							"  - example : 80\n" +
							"  - minimum : 1\n" +
							"  - maximum : 65534\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "The interval between health checks in seconds.\n" +
							"  - example : 5\n" +
							"  - minimum : 1\n" +
							"  - maximum : 180\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "The timeout for health check responses in seconds. Must be less than or equal to the interval.\n" +
							"  - example : 5\n" +
							"  - minimum : 1\n" +
							"  - maximum : 180\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "The number of consecutive health check failures before marking as unhealthy.\n" +
							"  - example : 3\n" +
							"  - minimum : 1\n" +
							"  - maximum : 10\n",
						Optional: true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "The HTTP method used for health checks.\n" +
							"  - example : GET\n" +
							"  - pattern : GET | POST\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "The URL path for HTTP health checks.\n" +
							"  - example : /test\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 50\n" +
							"  - pattern : ^/[A-Za-z0-9/._?&=-]*$\n",
						Optional: true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "The expected HTTP response code for health checks.\n" +
							"  - example : 200\n" +
							"  - minimum : 200\n" +
							"  - maximum : 599\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckType"): schema.StringAttribute{
						Description: "The type of health check.\n" +
							"  - example : DEFAULT\n" +
							"  - pattern : DEFAULT | CUSTOM\n",
						Optional: true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "The request data sent during health checks.\n" +
							"  - example : username=admin&password=1234\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("LbHealthCheckCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new LB Health Check.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : this is an lb server group\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 8acceeb6920c4fc494490d864f67f0b5\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 60fba45cb6c811efba41ba92e4fe7200\n",
						Optional: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol used for the health check.\n" +
							"  - example : TCP\n" +
							"  - pattern : TCP | HTTP | HTTPS\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Health Check.\n" +
							"  - example : ServerGroup01\n" +
							"  - minLength : 3\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9][-a-zA-Z0-9_]*[a-zA-Z0-9]$\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "The port number used for health checks.\n" +
							"  - example : 80\n" +
							"  - minimum : 1\n" +
							"  - maximum : 65534\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "The interval between health checks in seconds.\n" +
							"  - example : 5\n" +
							"  - minimum : 1\n" +
							"  - maximum : 180\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "The timeout for health check responses in seconds. Must be less than or equal to the interval.\n" +
							"  - example : 5\n" +
							"  - minimum : 1\n" +
							"  - maximum : 180\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "The number of consecutive health check failures before marking as unhealthy.\n" +
							"  - example : 3\n" +
							"  - minimum : 1\n" +
							"  - maximum : 10\n",
						Optional: true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "The HTTP method used for health checks.\n" +
							"  - example : GET\n" +
							"  - pattern : GET | POST\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "The URL path for HTTP health checks.\n" +
							"  - example : /test\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 50\n" +
							"  - pattern : ^/[A-Za-z0-9/._?&=-]*$\n",
						Optional: true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "The expected HTTP response code for health checks.\n" +
							"  - example : 200\n" +
							"  - minimum : 200\n" +
							"  - maximum : 599\n",
						Optional: true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "The request data sent during health checks.\n" +
							"  - example : username=admin&password=1234\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbHealthCheckResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clientv1d4 = inst.Client.LoadBalancerV1d4
	r.clients = inst.Client
}

func (r *loadbalancerLbHealthCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbHealthCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancerv1d4.LbHealthCheckResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Health Check using v1.4 API (async - returns 202 Accepted)
	lbHealthCheck := plan.LbHealthCheckCreate

	body := &loadbalancerv1d4.LbHealthCheckCreate{
		Name:                lbHealthCheck.Name,
		VpcId:               lbHealthCheck.VpcId,
		SubnetId:            lbHealthCheck.SubnetId,
		Protocol:            lbHealthCheck.Protocol,
		HealthCheckPort:     lbHealthCheck.HealthCheckPort,
		HealthCheckInterval: lbHealthCheck.HealthCheckInterval,
		HealthCheckTimeout:  lbHealthCheck.HealthCheckTimeout,
		HealthCheckCount:    lbHealthCheck.HealthCheckCount,
		HttpMethod:          lbHealthCheck.HttpMethod,
		HealthCheckUrl:      lbHealthCheck.HealthCheckUrl,
		ResponseCode:        lbHealthCheck.ResponseCode,
		RequestData:         lbHealthCheck.RequestData,
		Description:         lbHealthCheck.Description,
		Tags:                lbHealthCheck.Tags,
	}

	data, err := r.clientv1d4.CreateLbHealthCheck(ctx, body)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Could not create Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.LbHealthCheck.Id)

	// Map response body to schema and populate Computed attribute values
	lbHealthCheckModel := createLbHealthCheckModelV1d4(data)
	lbHealthCheckObjectValue, diags := types.ObjectValueFrom(ctx, lbHealthCheckModel.AttributeTypes(), lbHealthCheckModel)
	plan.LbHealthCheck = lbHealthCheckObjectValue

	// Set state before waiter — if waiter fails, ID remains in state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for ACTIVE state (async operation)
	refreshFn := r.getLbHealthCheckRefreshFunc(ctx, plan.Id.ValueString())
	err = client.WaitForResourceCreated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Error waiting for Lb Health Check to become active: "+err.Error(),
		)
		return
	}

	// Final Read to refresh state
	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbHealthCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancerv1d4.LbHealthCheckResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Health Check
	data, err := r.clientv1d4.GetLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Lb Health Check",
			"Could not read Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbHealthCheckModel := createLbHealthCheckModelV1d4(data)

	lbHealthCheckObjectValue, diags := types.ObjectValueFrom(ctx, lbHealthCheckModel.AttributeTypes(), lbHealthCheckModel)
	state.LbHealthCheck = lbHealthCheckObjectValue

	// Reconcile lb_health_check_create input block with API response to detect drift
	// Only populate if nil (e.g., after import) — preserve user config values otherwise
	if state.LbHealthCheckCreate == nil {
		state.LbHealthCheckCreate = &loadbalancerv1d4.LbHealthCheckCreate{
			Name:                types.StringValue(data.LbHealthCheck.Name),
			Description:         loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.Description.Get()),
			Protocol:            loadbalancerutil.ToNullableStringValue((*string)(data.LbHealthCheck.Protocol)),
			VpcId:               loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.VpcId.Get()),
			SubnetId:            loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.SubnetId.Get()),
			HealthCheckPort:     ToNullableInt32Value(data.LbHealthCheck.HealthCheckPort.Get()),
			HealthCheckInterval: ToNullableInt32Value(data.LbHealthCheck.HealthCheckInterval),
			HealthCheckTimeout:  ToNullableInt32Value(data.LbHealthCheck.HealthCheckTimeout),
			HealthCheckCount:    ToNullableInt32Value(data.LbHealthCheck.HealthCheckCount),
			HttpMethod:          loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.HttpMethod.Get()),
			HealthCheckUrl:      loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.HealthCheckUrl.Get()),
			ResponseCode:        loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.ResponseCode.Get()),
			RequestData:         loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.RequestData.Get()),
			Tags:                types.MapNull(types.StringType),
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbHealthCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancerv1d4.LbHealthCheckResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lbHealthCheck := state.LbHealthCheckCreate

	body := scploadbalancerv1d4.LbHealthCheckSetRequest{
		LbHealthCheck: scploadbalancerv1d4.LbHealthCheckSet{
			Protocol:            *scploadbalancerv1d4.NewNullableLbMonitorProtocol((*scploadbalancerv1d4.LbMonitorProtocol)(lbHealthCheck.Protocol.ValueStringPointer())),
			HealthCheckPort:     *scploadbalancerv1d4.NewNullableInt32(lbHealthCheck.HealthCheckPort.ValueInt32Pointer()),
			HealthCheckInterval: *scploadbalancerv1d4.NewNullableInt32(lbHealthCheck.HealthCheckInterval.ValueInt32Pointer()),
			HealthCheckTimeout:  *scploadbalancerv1d4.NewNullableInt32(lbHealthCheck.HealthCheckTimeout.ValueInt32Pointer()),
			HealthCheckCount:    *scploadbalancerv1d4.NewNullableInt32(lbHealthCheck.HealthCheckCount.ValueInt32Pointer()),
			HttpMethod:          *scploadbalancerv1d4.NewNullableLbMonitorHttpMethod((*scploadbalancerv1d4.LbMonitorHttpMethod)(lbHealthCheck.HttpMethod.ValueStringPointer())),
			HealthCheckUrl:      *scploadbalancerv1d4.NewNullableString(lbHealthCheck.HealthCheckUrl.ValueStringPointer()),
			ResponseCode:        *scploadbalancerv1d4.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
			RequestData:         *scploadbalancerv1d4.NewNullableString(lbHealthCheck.RequestData.ValueStringPointer()),
			Description:         *scploadbalancerv1d4.NewNullableString(lbHealthCheck.Description.ValueStringPointer()),
		},
	}

	// Update using v1.4 API (async - returns 202 Accepted)
	_, err := r.clientv1d4.UpdateLbHealthCheck(ctx, state.Id.ValueString(), &body)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Lb Health Check",
			"Could not update Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Set plan state before waiter
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for ACTIVE state (async operation)
	refreshFn := r.getLbHealthCheckRefreshFunc(ctx, state.Id.ValueString())
	err = client.WaitForResourceUpdated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Lb Health Check",
			"Error waiting for Lb Health Check to become active: "+err.Error(),
		)
		return
	}

	// Final Read to refresh state
	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLbHealthCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancerv1d4.LbHealthCheckResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Health Check using v1.4 API (async operation)
	err := r.clientv1d4.DeleteLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Health Check",
			"Could not delete lb health check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Wait for deletion to complete (async operation)
	refreshFn := r.getLbHealthCheckRefreshFunc(ctx, state.Id.ValueString())
	err = client.WaitForResourceDeleted(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Lb Health Check",
			"Error waiting for Lb Health Check to be deleted: "+err.Error(),
		)
		return
	}
}

func ToNullableInt32Value(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}

func (r *loadbalancerLbHealthCheckResource) getLbHealthCheckRefreshFunc(ctx context.Context, id string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		info, err := r.clientv1d4.GetLbHealthCheck(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, info.LbHealthCheck.State, nil
	}
}

func createLbHealthCheckModelV1d4(data *scploadbalancerv1d4.LbHealthCheckShowResponse) loadbalancerv1d4.LbHealthCheckDetail {
	lbHealthCheck := data.LbHealthCheck

	return loadbalancerv1d4.LbHealthCheckDetail{
		Name:                types.StringValue(lbHealthCheck.Name),
		VpcId:               loadbalancerutil.ToNullableStringValue(lbHealthCheck.VpcId.Get()),
		SubnetId:            loadbalancerutil.ToNullableStringValue(lbHealthCheck.SubnetId.Get()),
		Protocol:            loadbalancerutil.ToNullableStringValue((*string)(lbHealthCheck.Protocol)),
		HealthCheckPort:     ToNullableInt32Value(lbHealthCheck.HealthCheckPort.Get()),
		HealthCheckInterval: ToNullableInt32Value(lbHealthCheck.HealthCheckInterval),
		HealthCheckTimeout:  ToNullableInt32Value(lbHealthCheck.HealthCheckTimeout),
		HealthCheckCount:    ToNullableInt32Value(lbHealthCheck.HealthCheckCount),
		HealthCheckUrl:      loadbalancerutil.ToNullableStringValue(lbHealthCheck.HealthCheckUrl.Get()),
		HttpMethod:          loadbalancerutil.ToNullableStringValue(lbHealthCheck.HttpMethod.Get()),
		ResponseCode:        loadbalancerutil.ToNullableStringValue(lbHealthCheck.ResponseCode.Get()),
		RequestData:         loadbalancerutil.ToNullableStringValue(lbHealthCheck.RequestData.Get()),
		HealthCheckType:     types.StringValue(string(lbHealthCheck.HealthCheckType)),
		State:               types.StringValue(lbHealthCheck.State),
		AccountId:           loadbalancerutil.ToNullableStringValue(lbHealthCheck.AccountId.Get()),
		Description:         loadbalancerutil.ToNullableStringValue(lbHealthCheck.Description.Get()),
		ModifiedBy:          types.StringValue(lbHealthCheck.ModifiedBy),
		ModifiedAt:          types.StringValue(lbHealthCheck.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(lbHealthCheck.CreatedBy),
		CreatedAt:           types.StringValue(lbHealthCheck.CreatedAt.Format(time.RFC3339)),
	}
}
