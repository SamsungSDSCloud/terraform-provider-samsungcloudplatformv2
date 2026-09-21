package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.3"
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
	_ resource.Resource                = &loadbalancerLoadbalancerPublicNatIpResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLoadbalancerPublicNatIpResource{}
	_ resource.ResourceWithImportState = &loadbalancerLoadbalancerPublicNatIpResource{}
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
				Description: "Identifier of the resource.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "The LoadBalancer ID associated with the Public NAT IP.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Required: true,
			},
			common.ToSnakeCase("LoadbalancerPublicNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of public NAT.",
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
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ActionType"): schema.StringAttribute{
						Description: "The action type.\n" +
							"  - example : NAT_ALL\n",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Public NAT IP for internet access\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
						Description: "The external IP address.\n" +
							"  - example : 203.0.113.1\n",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the Public NAT IP.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("InternalIpAddress"): schema.StringAttribute{
						Description: "The internal IP address.\n" +
							"  - example : 10.0.0.1\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the Public NAT IP.\n" +
							"  - example : PublicNatIp01\n",
						Optional: true,
					},
					common.ToSnakeCase("OwnerId"): schema.StringAttribute{
						Description: "The owner ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("OwnerName"): schema.StringAttribute{
						Description: "The owner name.\n" +
							"  - example : LoadBalancer01\n",
						Optional: true,
					},
					common.ToSnakeCase("OwnerType"): schema.StringAttribute{
						Description: "The owner type.\n" +
							"  - example : ALB\n",
						Optional: true,
					},
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "The public IP ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ServiceIpPortId"): schema.StringAttribute{
						Description: "The service IP port ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Public NAT IP.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Optional: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of static NAT.\n" +
							"  - example : INTERNET\n" +
							"  - pattern : INTERNET | PRIVATE_NAT\n",
						Optional: true,
					},
					common.ToSnakeCase("vpc_id"): schema.StringAttribute{
						Description: "The VPC ID where the LoadBalancer is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("StaticNatCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer static NAT.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "The public IP ID.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
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

func (r *loadbalancerLoadbalancerPublicNatIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("loadbalancer_id"), req, resp)
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

	plan.Id = virtualserverutil.ToNullableStringValue(data.StaticNat.Id.Get())

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
	var state loadbalancer.LoadbalancerPublicNatIpResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call ShowLoadbalancerPublicNatIp to refresh NAT IP details.
	data, err := r.client.GetLoadbalancerPublicNatIp(ctx, state.LoadbalancerId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Public NAT", err.Error())
		return
	}

	// Merge Show response with existing state.
	// Show API returns only 3 fields (ExternalIpAddress, PublicipId, State).
	// Other fields are preserved from existing state to avoid data loss.
	var existingDetail loadbalancer.LoadbalancerPublicNatIpDetail
	if !state.LoadbalancerPublicNatIp.IsNull() {
		diags := state.LoadbalancerPublicNatIp.As(ctx, &existingDetail, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	refreshedDetail := readLoadbalancerNatModel(data, existingDetail)
	staticNatObjectValue, diags := types.ObjectValueFrom(ctx, refreshedDetail.AttributeTypes(), refreshedDetail)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.LoadbalancerPublicNatIp = staticNatObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerLoadbalancerPublicNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"Loadbalancer Public NAT IP does not support in-place updates. To change configuration, recreate the resource.",
	)
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

	// Wait for the Public NAT IP to be fully deleted (404 = success)
	refreshFn := r.getPublicNatIpRefreshFunc(ctx, state.LoadbalancerId.ValueString())
	err = client.WaitForResourceDeleted(ctx, refreshFn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting LB Public NAT",
			"Error waiting for LB Public NAT to become deleted: "+err.Error(),
		)
		return
	}
}

func createLoadbalancerNatModel(data *scploadbalancer.StaticNatCreateResponse) loadbalancer.LoadbalancerPublicNatIpDetail {
	lbStaticNat := data.StaticNat
	return loadbalancer.LoadbalancerPublicNatIpDetail{
		AccountId:         virtualserverutil.ToNullableStringValue(lbStaticNat.AccountId.Get()),
		ActionType:        virtualserverutil.ToNullableStringValue(lbStaticNat.ActionType.Get()),
		CreatedAt:         types.StringValue(lbStaticNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(lbStaticNat.CreatedBy),
		Description:       virtualserverutil.ToNullableStringValue(lbStaticNat.Description.Get()),
		ExternalIpAddress: virtualserverutil.ToNullableStringValue(lbStaticNat.ExternalIpAddress.Get()),
		Id:                virtualserverutil.ToNullableStringValue(lbStaticNat.Id.Get()),
		InternalIpAddress: virtualserverutil.ToNullableStringValue(lbStaticNat.InternalIpAddress.Get()),
		ModifiedAt:        types.StringValue(lbStaticNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(lbStaticNat.ModifiedBy),
		Name:              virtualserverutil.ToNullableStringValue(lbStaticNat.Name.Get()),
		OwnerId:           virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerId.Get()),
		OwnerName:         virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerName.Get()),
		OwnerType:         virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerType.Get()),
		PublicipId:        virtualserverutil.ToNullableStringValue(lbStaticNat.PublicipId.Get()),
		ServiceIpPortId:   virtualserverutil.ToNullableStringValue(lbStaticNat.ServiceIpPortId.Get()),
		State:             virtualserverutil.ToNullableStringValue(lbStaticNat.State.Get()),
		SubnetId:          virtualserverutil.ToNullableStringValue(lbStaticNat.SubnetId.Get()),
		Type:              virtualserverutil.ToNullableStringValue(lbStaticNat.Type.Get()),
		VpcId:             virtualserverutil.ToNullableStringValue(lbStaticNat.VpcId.Get()),
	}
}

// readLoadbalancerNatModel maps the Show API response (StaticNat — 3 fields) to the provider model,
// merging with existing state to preserve fields not returned by the Show API.
func readLoadbalancerNatModel(data *scploadbalancer.LoadbalancerStaticNatResponse, existing loadbalancer.LoadbalancerPublicNatIpDetail) loadbalancer.LoadbalancerPublicNatIpDetail {
	staticNat := data.StaticNat
	return loadbalancer.LoadbalancerPublicNatIpDetail{
		// Fields from Show API (refreshed)
		ExternalIpAddress: types.StringValue(staticNat.ExternalIpAddress),
		PublicipId:        virtualserverutil.ToNullableStringValue(staticNat.PublicipId.Get()),
		State:             types.StringValue(staticNat.State),
		// Fields preserved from existing state (not available from Show API)
		AccountId:         existing.AccountId,
		ActionType:        existing.ActionType,
		CreatedAt:         existing.CreatedAt,
		CreatedBy:         existing.CreatedBy,
		Description:       existing.Description,
		Id:                existing.Id,
		InternalIpAddress: existing.InternalIpAddress,
		ModifiedAt:        existing.ModifiedAt,
		ModifiedBy:        existing.ModifiedBy,
		Name:              existing.Name,
		OwnerId:           existing.OwnerId,
		OwnerName:         existing.OwnerName,
		OwnerType:         existing.OwnerType,
		ServiceIpPortId:   existing.ServiceIpPortId,
		SubnetId:          existing.SubnetId,
		Type:              existing.Type,
		VpcId:             existing.VpcId,
	}
}

func (r *loadbalancerLoadbalancerPublicNatIpResource) getPublicNatIpRefreshFunc(ctx context.Context, loadbalancerId string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		data, err := r.client.GetLoadbalancerPublicNatIp(ctx, loadbalancerId)
		if err != nil {
			return nil, "", err
		}
		return data, data.StaticNat.State, nil
	}
}
