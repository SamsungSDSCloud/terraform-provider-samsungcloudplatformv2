package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcTgwResource{}
	_ resource.ResourceWithConfigure = &vpcTgwResource{}
)

// NewVpcTgwResource is a helper function to simplify the provider implementation.
func NewVpcTgwResource() resource.Resource {
	return &vpcTgwResource{}
}

type vpcTgwResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcTgwResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transitgateway"
}

// Schema defines the schema for the data source.
func (r *vpcTgwResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "transitgateway",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "TGW name",
				Optional:    true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description",
				Required:    true,
			},
			common.ToSnakeCase("Tgw"): schema.SingleNestedAttribute{
				Description: "Tgw",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "Bandwidth",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description\n" +
							"  - example : Tgw description\n" +
							"  - maxLength : 50\n" +
							"  - minLength : 1",
						Computed: true,
					},
					common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
						Description: "FirewallIds",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name\n" +
							"  - example : Tgw name\n" +
							"  - pattern : ^[a-zA-Z0-9]*$\n" +
							"  - maxLength : 20\n" +
							"  - minLength : 3",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State" +
							" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
						Computed: true,
					},
					common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
						Description: "UplinkEnabled" +
							"  - example : false\n",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcTgwResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *vpcTgwResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	fmt.Println("createData", req)

	// Retrieve values from plan
	var plan vpc.TgwResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subnet
	data, _, err := r.client.CreateTransitGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating CreateTransitGateway",
			"Could not create CreateTransitGateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tgw := data.TransitGateway
	plan.Id = types.StringValue(tgw.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForTgwtStatus(ctx, r.client, tgw.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tgw",
			"Error waiting for CreateTransitGateway to become active: "+err.Error(),
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
func (r *vpcTgwResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	fmt.Println("getRequestDetail", req)

	var state vpc.TgwResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetTransitGatewayInfo(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading tgw",
			"Could not read transitGateway  "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	transitGateway := data.TransitGateway

	tgwModel := vpc.Tgw{
		Id:            types.StringValue(transitGateway.Id),
		Description:   types.StringPointerValue(transitGateway.Description.Get()),
		Name:          types.StringValue(transitGateway.Name),
		AccountId:     types.StringValue(transitGateway.AccountId),
		Bandwidth:     types.Int32PointerValue(transitGateway.Bandwidth.Get()),
		CreatedAt:     types.StringValue(transitGateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(transitGateway.CreatedBy),
		FirewallIds:   types.StringPointerValue(transitGateway.FirewallIds.Get()),
		ModifiedAt:    types.StringValue(transitGateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(transitGateway.ModifiedBy),
		State:         types.StringValue(string(transitGateway.State)),
		UplinkEnabled: types.BoolPointerValue(transitGateway.UplinkEnabled),
	}
	tgwObjectValue, diags := types.ObjectValueFrom(ctx, tgwModel.AttributeTypes(), tgwModel)
	state.Tgw = tgwObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcTgwResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	fmt.Println("getRequestUpdate", req)
	var state vpc.TgwResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing tgw
	err := r.client.UpdateTransitGateway(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating tgw",
			"Could not update tgw, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetTransitGatewayInfo(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GetTransitGatewayInfo",
			"Could not read GetTransitGatewayInfo ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tgwElem := data.TransitGateway

	tgwModel := vpc.Tgw{
		Id:            types.StringValue(tgwElem.Id),
		Description:   types.StringPointerValue(tgwElem.Description.Get()),
		Name:          types.StringValue(tgwElem.Name),
		AccountId:     types.StringValue(tgwElem.AccountId),
		Bandwidth:     types.Int32PointerValue(tgwElem.Bandwidth.Get()),
		CreatedAt:     types.StringValue(tgwElem.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(tgwElem.CreatedBy),
		FirewallIds:   types.StringPointerValue(tgwElem.FirewallIds.Get()),
		ModifiedAt:    types.StringValue(tgwElem.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(tgwElem.ModifiedBy),
		State:         types.StringValue(string(tgwElem.State)),
		UplinkEnabled: types.BoolPointerValue(tgwElem.UplinkEnabled),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, tgwModel.AttributeTypes(), tgwModel)
	state.Tgw = vpcObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcTgwResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state

	fmt.Println("getRequestDelete", req)

	var state vpc.TgwResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subnet
	err := r.client.DeleteTransitGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting TgwResource",
			"Could not delete TgwResource, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForTgwtStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting TgwResource",
			"Error waiting for TgwResource to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForTgwtStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetTransitGatewayInfo(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.TransitGateway.State), nil
	})
}
