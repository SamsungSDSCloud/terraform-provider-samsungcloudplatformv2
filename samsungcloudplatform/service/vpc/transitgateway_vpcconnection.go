package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcTgwVpcConnectionResource{}
	_ resource.ResourceWithConfigure   = &vpcTgwVpcConnectionResource{}
	_ resource.ResourceWithImportState = &vpcTgwVpcConnectionResource{}
)

// NewVpcTgwVpcConnectionResource is a helper function to simplify the provider implementation.
func NewVpcTgwVpcConnectionResource() resource.Resource {
	return &vpcTgwVpcConnectionResource{}
}

type vpcTgwVpcConnectionResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (r *vpcTgwVpcConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

func (v vpcTgwVpcConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_vpc_connection"
}

func (r *vpcTgwVpcConnectionResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	parts := strings.Split(request.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		response.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: TransitGaewayId/VpcConnectionId, got: %q", request.ID),
		)
		return
	}

	response.State.SetAttribute(ctx, path.Root("transit_gateway_id"), types.StringValue(parts[0]))
	response.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))
}

func (v *vpcTgwVpcConnectionResource) Schema(ctx context.Context, request resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The connection between a VPC and a Transit Gateway",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the connection.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the connection belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the connection belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("TransitGatewayVpcConnection"): schema.SingleNestedAttribute{
				Description: "The connection between a VPC and a Transit Gateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The unique identifier of the connection.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the connection.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the connection.\n " +
							"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
						Description: "The identifier of the transit gateway that the connection belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the connection belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the connection belongs to.\n" +
							"  - example : vpcName",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r vpcTgwVpcConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.TgwVpcConnectionResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new direct connect
	data, _, err := r.client.CreateTransitGatewayVpcConnection(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating tgw vpc connection",
			"Could not create tgw vpc connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.TransitGatewayVpcConnection.Id)
	diags = resp.State.Set(ctx, plan)
	tgwvpcconModel := createTgwVpcConnectionModel(&data.TransitGatewayVpcConnection)
	vpcconObjectValue, diags := types.ObjectValueFrom(ctx, tgwvpcconModel.AttributeTypes(), tgwvpcconModel)
	plan.TgwVpcConnection = vpcconObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	err = waitForVpcConnectStatus(ctx, r.client, plan.TransitGatewayId.ValueString(), data.TransitGatewayVpcConnection.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tgw vpc connection",
			"Error waiting for tgw vpc connection to become active: "+err.Error(),
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

func (r vpcTgwVpcConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state vpc.TgwVpcConnectionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := r.client.GetTransitGatewayVpcConnection(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading tgw vpc connection",
			"Could not read tgw vpc connection ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// No data return from List API <=> Detail data not found
	if len(data.TransitGatewayVpcConnections) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	conn := data.TransitGatewayVpcConnections[0]

	// Rewrite configurable top-level input fields from API response to detect drift
	state.TransitGatewayId = types.StringValue(conn.TransitGatewayId)
	state.VpcId = types.StringValue(conn.VpcId)

	vpcconnectionModel := createTgwVpcConnectionModel(&conn)

	vpcconnectionObjectValue, diags := types.ObjectValueFrom(ctx, vpcconnectionModel.AttributeTypes(), vpcconnectionModel)
	state.TgwVpcConnection = vpcconnectionObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (v vpcTgwVpcConnectionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// Non-compliant: Empty function
}

func (r *vpcTgwVpcConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.TgwVpcConnectionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing tgw vpc connection
	err := r.client.DeleteTransitGatewayVpcConnection(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting tgw vpc connection",
			"Could not delete tgw vpc connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcConnectStatus(ctx, r.client, state.TransitGatewayId.ValueString(), state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting tgw vpc connection",
			"Error waiting for tgw vpc connection to become deleted: "+err.Error(),
		)
		return
	}
}

func createTgwVpcConnectionModel(data *scpvpc.TransitGatewayVpcConnection) vpc.TgwVpcConnection {
	return vpc.TgwVpcConnection{
		AccountId:        types.StringValue(data.AccountId),
		CreatedAt:        types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.CreatedBy),
		Id:               types.StringValue(data.Id),
		ModifiedAt:       types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.ModifiedBy),
		State:            types.StringValue(string(data.State)),
		TransitGatewayId: types.StringValue(data.TransitGatewayId),
		VpcId:            types.StringValue(data.VpcId),
		VpcName:          types.StringValue(data.VpcName),
	}
}

func waitForVpcConnectStatus(ctx context.Context, vpcConnectClient *vpc.Client, transitGatewayId string, vpcConnectionId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcConnectClient.GetTransitGatewayVpcConnection(ctx, transitGatewayId, vpcConnectionId)
		if err != nil {
			return nil, "", err
		}
		if len(info.TransitGatewayVpcConnections) == 0 {
			return info, "DELETED", nil
		}
		return info, string(info.TransitGatewayVpcConnections[0].State), nil
	}, -1, -1, -1, -1)
}
