package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/loadbalancer/1.3"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loadbalancerLoadbalancerPrivateNatIpResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLoadbalancerPrivateNatIpResource{}
	_ resource.ResourceWithImportState = &loadbalancerLoadbalancerPrivateNatIpResource{}
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
				Description: "Identifier of the resource.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "The LoadBalancer ID associated with the Private NAT IP.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Required: true,
			},
			common.ToSnakeCase("LoadbalancerPrivateNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of private NAT.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Optional: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
						Description: "The external IP address.\n" +
							"  - example : 192.168.0.1\n",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the Private NAT IP.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("InternalIpAddress"): schema.StringAttribute{
						Description: "The internal IP address.\n" +
							"  - example : 10.0.0.1\n",
						Optional: true,
					},
					common.ToSnakeCase("PrivateNatIpId"): schema.StringAttribute{
						Description: "The private NAT IP ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Private NAT IP.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("PrivateStaticNatCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer private static NAT.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
						Description: "The private NAT ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("PrivateNatIpId"): schema.StringAttribute{
						Description: "The private NAT IP ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
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

func (r *loadbalancerLoadbalancerPrivateNatIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("loadbalancer_id"), req, resp)
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

	// Wait for ACTIVE state
	if err := waitForLoadbalancerPrivateNatIpStatus(ctx, r.client, plan.LoadbalancerId.ValueString(), []string{}, []string{"ACTIVE"}); err != nil {
		resp.Diagnostics.AddError(
			"Error creating LoadBalancer Private NAT",
			"Error waiting for Private NAT to become active: "+err.Error(),
		)
		return
	}

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
	var state loadbalancer.LoadbalancerPrivateNatIpResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetLoadbalancerPrivateNatIp(ctx, state.LoadbalancerId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Private NAT", err.Error())
		return
	}

	// Rewrite configurable top-level input fields from API response to detect drift
	privateNat := data.StaticNat.Get()
	state.Id = loadbalancerutil.ToNullableStringValue(privateNat.Id.Get())

	staticNatModel := createLoadbalancerPrivateNatModelFromShow(data)
	staticNatObjectValue, d := types.ObjectValueFrom(ctx, staticNatModel.AttributeTypes(), staticNatModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.LoadbalancerPrivateNatIp = staticNatObjectValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *loadbalancerLoadbalancerPrivateNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"Loadbalancer Private NAT IP does not support in-place updates. To change configuration, recreate the resource.",
	)
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

func createLoadbalancerPrivateNatModelFromShow(data *scploadbalancer.PrivateStaticNatShowResponse) loadbalancer.LoadbalancerPrivateNatIpDetail {
	lbStaticNat := data.StaticNat.Get()
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

func waitForLoadbalancerPrivateNatIpStatus(ctx context.Context, loadbalancerClient *loadbalancer.Client, loadbalancerId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := loadbalancerClient.GetLoadbalancerPrivateNatIp(ctx, loadbalancerId)
		if err != nil {
			return nil, "", err
		}
		state := info.StaticNat.Get().State.Get()
		if state == nil {
			return nil, "", fmt.Errorf("Private NAT state is nil")
		}
		return info, *state, nil
	}, -1, -1, -1, -1)
}
