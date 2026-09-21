package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancerv1d4"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scploadbalancerv1d4 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loadbalancerLoadbalancerResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLoadbalancerResource{}
	_ resource.ResourceWithImportState = &loadbalancerLoadbalancerResource{}
)

// NewLoadBalancerLoadBalancerResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLoadBalancerResource() resource.Resource {
	return &loadbalancerLoadbalancerResource{}
}

// loadbalancerLoadbalancerResource is the data source implementation.
type loadbalancerLoadbalancerResource struct {
	config     *scpsdk.Configuration
	clientv1d4 *loadbalancerv1d4.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLoadbalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLoadbalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "LoadBalancer resource for distributing traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Loadbalancer"): schema.SingleNestedAttribute{
				Description: "Details of the LoadBalancer.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LoadBalancer for web traffic\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The firewall ID associated with the LoadBalancer.\n" +
							"  - example : 6af449ba19134eb1929241e0ea0d718d\n",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("LayerType"): schema.StringAttribute{
						Description: "The layer type of the Load Balancer.\n" +
							"  - example : L7\n" +
							"  - pattern : L4 | L7\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LoadBalancer.\n" +
							"  - example : LoadBalancer01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Load Balancer.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the LoadBalancer is deployed.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the LoadBalancer is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("PublicNatEnabled"): schema.BoolAttribute{
						Description: "Whether public NAT is enabled.\n" +
							"  - example : true\n",
						Computed: true,
					},
					common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
						Description: "The service IP address of the LoadBalancer.\n" +
							"  - example : 20.20.1.87\n",
						Computed: true,
					},
					common.ToSnakeCase("SourceNatIp"): schema.StringAttribute{
						Description: "The source NAT IP address of the LoadBalancer.\n" +
							"  - example : 20.20.0.127\n",
						Computed: true,
					},
					common.ToSnakeCase("HealthCheckIps"): schema.ListAttribute{
						Description: "The list of health check IP addresses.\n" +
							"  - example : [\"10.0.0.1\", \"10.0.0.2\"]\n",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Zones"): schema.ListAttribute{
						Description: "The list of availability zones where the subnet is located.\n" +
							"  - example : [\"zone-1\", \"zone-2\"]\n",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			common.ToSnakeCase("LoadbalancerCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new LoadBalancer.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LoadBalancer for web traffic\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("FirewallEnabled"): schema.BoolAttribute{
						Description: "Whether firewall is enabled.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("FirewallLoggingEnabled"): schema.BoolAttribute{
						Description: "Whether firewall logging is enabled.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("LayerType"): schema.StringAttribute{
						Description: "The layer type of the Load Balancer.\n" +
							"  - example : L7\n" +
							"  - pattern : L4 | L7\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LoadBalancer.\n" +
							"  - example : LoadBalancer01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n",
						Optional: true,
					},
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "The Public IP ID address.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
						Description: "The service IP address.\n" +
							"  - example : 192.168.0.1\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("SourceNatIp"): schema.StringAttribute{
						Description: "The source NAT IP address.\n" +
							"  - example : 192.168.0.1\n",
						Optional: true,
					},
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("Zones"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of availability zones where the subnet is located.\n" +
							"  - example : [\"zone-1\", \"zone-2\"]",
						MarkdownDescription: "The list of availability zones where the subnet is located.\n" +
							"  - example : [\"zone-1\", \"zone-2\"]",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckIps"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of health check IP addresses.\n" +
							"  - example : [\"192.168.0.1\", \"192.168.0.1\"]",
						MarkdownDescription: "The list of health check IP addresses.\n" +
							"  - example : [\"192.168.0.1\", \"192.168.0.1\"]",
						Optional: true,
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

	r.clientv1d4 = inst.Client.LoadBalancerV1d4
}

func (r *loadbalancerLoadbalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLoadbalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancerv1d4.LoadbalancerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Loadbalancer
	data, err := r.clientv1d4.CreateLoadbalancer(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Loadbalancer",
			"Could not create Loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Loadbalancer.Id)

	// Set ID into state before waiting so that if the waiter fails,
	// Terraform retains the resource ID and does not recreate it on the next run.
	resp.State.Set(ctx, plan)

	refreshFn := r.getLoadbalancerRefreshFunc(ctx, data.Loadbalancer.Id)
	err = client.WaitForResourceCreated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Loadbalancer",
			"Error waiting for Loadbalancer to become active: "+err.Error(),
		)
		return
	}

	// Refresh resource state from API
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}
	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLoadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancerv1d4.LoadbalancerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Loadbalancer
	data, err := r.clientv1d4.GetLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Reconcile loadbalancer_create input block with API response to detect drift
	// Only populate if nil (e.g., after import) — preserve user config values otherwise
	if state.LoadbalancerCreate == nil {
		emptyTags, _ := types.MapValue(types.StringType, map[string]attr.Value{})
		state.LoadbalancerCreate = &loadbalancerv1d4.LoadbalancerCreate{
			Name:           types.StringValue(data.Loadbalancer.Name),
			Description:    virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
			LayerType:      types.StringValue(data.Loadbalancer.LayerType),
			Zones:          convertList(data.Loadbalancer.Zones),
			HealthCheckIps: convertList(data.Loadbalancer.HealthCheckIps),
			Tags:           emptyTags,
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
func (r *loadbalancerLoadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancerv1d4.LoadbalancerResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.clientv1d4.UpdateLoadbalancer(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Loadbalancer",
			"Could not update Loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Restore loadbalancer_create from plan so Read does not overwrite it with API values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshFn := r.getLoadbalancerRefreshFunc(ctx, state.Id.ValueString())
	err = client.WaitForResourceUpdated(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Loadbalancer",
			"Error waiting for Loadbalancer to become active: "+err.Error(),
		)
		return
	}

	// Refresh resource state from API
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLoadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancerv1d4.LoadbalancerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Loadbalancer
	err := r.clientv1d4.DeleteLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Loadbalancer",
			"Could not delete loadbalancer, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	refreshFn := r.getLoadbalancerRefreshFunc(ctx, state.Id.ValueString())
	err = client.WaitForResourceDeleted(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Loadbalancer",
			"Error waiting for Loadbalancer to be deleted: "+err.Error(),
		)
		return
	}

	// Final Read after delete — confirms resource is gone (404) and rebuilds state
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	resp.State = readResp.State
}

func convertList(zones []string) types.List {
	if zones == nil {
		zones = []string{}
	}
	var elements []attr.Value
	for _, z := range zones {
		elements = append(elements, types.StringValue(z))
	}
	return types.ListValueMust(types.StringType, elements)
}

func createLoadbalancerModelForRead(data *scploadbalancerv1d4.LoadbalancerShowResponseV1Dot4) loadbalancerv1d4.LoadbalancerCreateResponseDetail {
	return loadbalancerv1d4.LoadbalancerCreateResponseDetail{
		AccountId:        types.StringValue(data.Loadbalancer.AccountId),
		CreatedAt:        types.StringValue(data.Loadbalancer.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.Loadbalancer.CreatedBy),
		Description:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
		Id:               types.StringValue(data.Loadbalancer.Id),
		LayerType:        types.StringValue(data.Loadbalancer.LayerType),
		ModifiedAt:       types.StringValue(data.Loadbalancer.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.Loadbalancer.ModifiedBy),
		Name:             types.StringValue(data.Loadbalancer.Name),
		State:            types.StringValue(data.Loadbalancer.State),
		SubnetId:         types.StringValue(data.Loadbalancer.SubnetId),
		VpcId:            types.StringValue(data.Loadbalancer.VpcId),
		FirewallId:       virtualserverutil.ToNullableStringValue(data.Loadbalancer.FirewallId.Get()),
		PublicNatEnabled: common.ToNullableBoolValue(data.Loadbalancer.PublicNatEnabled.Get()),
		ServiceIp:        virtualserverutil.ToNullableStringValue(data.Loadbalancer.ServiceIp.Get()),
		SourceNatIp:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.SourceNatIp.Get()),
		Zones:            convertList(data.Loadbalancer.Zones),
		HealthCheckIps:   convertList(data.Loadbalancer.HealthCheckIps),
	}
}

func (r *loadbalancerLoadbalancerResource) getLoadbalancerRefreshFunc(ctx context.Context, id string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		data, err := r.clientv1d4.GetLoadbalancer(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return data, string(data.Loadbalancer.State), nil
	}
}


