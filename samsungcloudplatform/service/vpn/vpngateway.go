package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpn/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpnVpnGatewayResource{}
	_ resource.ResourceWithConfigure = &vpnVpnGatewayResource{}
	_ resource.ResourceWithImportState = &vpnVpnGatewayResource{}
)

// NewVpnVpnGatewayResource is a helper function to simplify the provider implementation.
func NewVpnVpnGatewayResource() resource.Resource {
	return &vpnVpnGatewayResource{}
}

// vpnVpnGatewayResource is the data source implementation.
type vpnVpnGatewayResource struct {
	config *scpsdk.Configuration
	client *vpn.Client
}

// Metadata returns the data source type name.
func (r *vpnVpnGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_gateway"
}

// Schema defines the schema for the data source.
func (r *vpnVpnGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages VPN gateways to connect external networks securely.",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "A brief explanation or note about this resource.\n" +
					"  - example: VPN test\n" +
					"  - constraints: maxLength: 40",
				Optional: true,
				Computed: true,
				Default: stringdefault.StaticString(""),
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "The IP address assigned to the resource.\n" +
					"  - example: 10.0.0.0/24",
				Required: true,
			},
			common.ToSnakeCase("IpId"): schema.StringAttribute{
				Description: "The identifier of the IP address assigned to the resource.\n" +
					"  - example: bd07e102fe574edf8a1748957c45bdbf",
				Required: true,
			},
			common.ToSnakeCase("IpType"): schema.StringAttribute{
				Description: "The type of the IP address assigned to the resource.\n" +
					"  - example: PUBLIC\n" +
					"  - valid: PUBLIC",
				Required: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the resource.\n" +
					"  - example: vpnGWProd\n" +
					"  - valid: English letters and numbers only\n" +
					"  - constraints: minLength: 1, maxLength: 20",
				Required: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the resource belongs to.\n" +
					"  - example: f32265726b694b32920aa3b111f4c715",
				Required: true,
			},
			common.ToSnakeCase("VpnGateway"): schema.SingleNestedAttribute{
				Description: "The identifier of the VPN gateway that the resource belongs to.\n" +
					"  - example: 01c543eb4b8d42a9a3502345d4025147",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example: 297615908b8e4ec69520a99a6777add3",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: VPN test",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "The IP address assigned to the resource.\n" +
							"  - example: 10.0.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("IpId"): schema.StringAttribute{
						Description: "The identifier of the IP address assigned to the resource.\n" +
							"  - example: bd07e102fe574edf8a1748957c45bdbf",
						Computed: true,
					},
					common.ToSnakeCase("IpType"): schema.StringAttribute{
						Description: "The type of the IP address assigned to the resource.\n" +
							"  - example: PUBLIC",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
							"  - example: vpnGWProd",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the resource belongs to.\n" +
							"  - example: f32265726b694b32920aa3b111f4c715",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the resource belongs to.\n" +
							"  - example: vpcProd01",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpnVpnGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpn
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpnVpnGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpn.VpnGatewayResource
	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpn gateway
	data, err := r.client.CreateVpnGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpn gateway",
			"Could not create vpn gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error creating vpn gateway",
			"Create vpn gateway returned empty response",
		)
		return
	}

	vpnGateway := data.VpnGateway
	plan.Id = types.StringValue(vpnGateway.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForVpnGatewayStatus(ctx, r.client, vpnGateway.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpn gateway",
			"Error waiting for vpn gateway to become active: "+err.Error(),
		)
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *vpnVpnGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpn.VpnGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Get vpn gateway returned empty response",
		)
		return
	}

	vgModel := createVpnGatewayModel(data)
	vg := data.VpnGateway

	vgObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Description = types.StringPointerValue(vg.Description.Get())
	state.IpAddress = types.StringValue(vg.IpAddress)
	state.IpId = types.StringValue(vg.IpId)
	state.IpType = types.StringValue(vg.IpType)
	state.Name = types.StringValue(vg.Name)
	state.VpcId = types.StringValue(vg.VpcId)
	state.VpnGateway = vgObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpnVpnGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpn.VpnGatewayResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpnGateway(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpnGateway as UpdateVpnGateway items are not populated.
	data, err := r.client.GetVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Get vpn gateway returned empty response",
		)
		return
	}

	vgModel := createVpnGatewayModel(data)
	vg := data.VpnGateway

	vgObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Description = types.StringPointerValue(vg.Description.Get())
	state.IpAddress = types.StringValue(vg.IpAddress)
	state.IpId = types.StringValue(vg.IpId)
	state.IpType = types.StringValue(vg.IpType)
	state.Name = types.StringValue(vg.Name)
	state.VpcId = types.StringValue(vg.VpcId)
	state.VpnGateway = vgObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpnVpnGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpn.VpnGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpn gateway
	err := r.client.DeleteVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpn gateway",
			"Could not delete vpn gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnGatewayStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpn gateway",
			"Error waiting for vpn gateway to become deleted: "+err.Error(),
		)
		return
	}
}

func createVpnGatewayModel(data *scpvpn.VpnGatewayShowResponse) vpn.VpnGateway {
	vg := data.VpnGateway
	return vpn.VpnGateway{
		AccountId:   types.StringValue(vg.AccountId),
		CreatedAt:   types.StringValue(vg.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(vg.CreatedBy),
		Description: types.StringPointerValue(vg.Description.Get()),
		Id:          types.StringValue(vg.Id),
		IpAddress:   types.StringValue(vg.IpAddress),
		IpId:        types.StringValue(vg.IpId),
		IpType:      types.StringValue(vg.IpType),
		ModifiedAt:  types.StringValue(vg.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(vg.ModifiedBy),
		Name:        types.StringValue(vg.Name),
		State:       types.StringValue(string(vg.State)),
		VpcId:       types.StringValue(vg.VpcId),
		VpcName:     types.StringValue(vg.VpcName),
	}
}

func waitForVpnGatewayStatus(ctx context.Context, vpnClient *vpn.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpnClient.GetVpnGateway(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpnGateway.State), nil
	}, -1, -1, -1, -1)
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *vpnVpnGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
  resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
