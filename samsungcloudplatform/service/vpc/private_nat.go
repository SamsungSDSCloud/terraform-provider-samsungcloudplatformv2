package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcPrivateNatResource{}
	_ resource.ResourceWithConfigure = &vpcPrivateNatResource{}
)

// NewVpcPrivateNatResource is a helper function to simplify the provider implementation.
func NewVpcPrivateNatResource() resource.Resource {
	return &vpcPrivateNatResource{}
}

// vpcPrivateNatResource is the data source implementation.
type vpcPrivateNatResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nat"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Private NAT.",
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
				Description: "Private NAT Name \n" +
					"  - example : privateNatName\n" +
					"  - minLength : 3\n" +
					"  - maxLength : 20\n" +
					"  - pattern : ^[a-zA-Z0-9]+$",
				Required: true,
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "Direct Connect ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "CIDR \n" +
					"  - example : 192.168.10.0/24 \n",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Private NAT description\n" +
					"  - maxLength : 50",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNat"): schema.SingleNestedAttribute{
				Description: "Private NAT",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
					common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
						Description: "DirectConnectId",
						Computed:    true,
					},
					common.ToSnakeCase("DirectConnectName"): schema.StringAttribute{
						Description: "DirectConnectName",
						Computed:    true,
					},
					common.ToSnakeCase("Cidr"): schema.StringAttribute{
						Description: "Cidr",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
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
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPrivateNatResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPrivateNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.PrivateNatResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreatePrivateNat(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Private NAT",
			"Could not create Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(privateNat.Id)

	privateNatModel := vpc.PrivateNat{
		Id:                types.StringValue(privateNat.Id),
		Name:              types.StringValue(privateNat.Name),
		VpcId:             types.StringValue(privateNat.VpcId),
		VpcName:           types.StringPointerValue(privateNat.VpcName.Get()),
		DirectConnectId:   types.StringValue(privateNat.DirectConnectId),
		DirectConnectName: types.StringPointerValue(privateNat.DirectConnectName.Get()),
		Cidr:              types.StringValue(privateNat.Cidr),
		State:             types.StringValue(string(privateNat.State)),
		Description:       types.StringPointerValue(privateNat.Description.Get()),
		CreatedAt:         types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(privateNat.CreatedBy),
		ModifiedAt:        types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, diags := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	plan.PrivateNat = privateNatObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForPrivateNatStatus(ctx, r.client, privateNat.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Private NAT",
			"Error waiting for Private NAT to become active: "+err.Error(),
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
func (r *vpcPrivateNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.PrivateNatResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Private NAT
	data, err := r.client.GetPrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Private NAT",
			"Could not read Private NAT ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat

	privateNatModel := vpc.PrivateNat{
		Id:                types.StringValue(privateNat.Id),
		Name:              types.StringValue(privateNat.Name),
		VpcId:             types.StringValue(privateNat.VpcId),
		VpcName:           types.StringPointerValue(privateNat.VpcName.Get()),
		DirectConnectId:   types.StringValue(privateNat.DirectConnectId),
		DirectConnectName: types.StringPointerValue(privateNat.DirectConnectName.Get()),
		Cidr:              types.StringValue(privateNat.Cidr),
		State:             types.StringValue(string(privateNat.State)),
		Description:       types.StringPointerValue(privateNat.Description.Get()),
		CreatedAt:         types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(privateNat.CreatedBy),
		ModifiedAt:        types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, diags := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	state.PrivateNat = privateNatObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPrivateNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.PrivateNatResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Private NAT
	_, err := r.client.UpdatePrivateNat(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Private NAT",
			"Could not update Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetPrivateNat as UpdatePrivateNat items are not populated.
	data, err := r.client.GetPrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Private NAT",
			"Could not read Private NAT ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat

	privateNatModel := vpc.PrivateNat{
		Id:                types.StringValue(privateNat.Id),
		Name:              types.StringValue(privateNat.Name),
		VpcId:             types.StringValue(privateNat.VpcId),
		VpcName:           types.StringPointerValue(privateNat.VpcName.Get()),
		DirectConnectId:   types.StringValue(privateNat.DirectConnectId),
		DirectConnectName: types.StringPointerValue(privateNat.DirectConnectName.Get()),
		Cidr:              types.StringValue(privateNat.Cidr),
		State:             types.StringValue(string(privateNat.State)),
		Description:       types.StringPointerValue(privateNat.Description.Get()),
		CreatedAt:         types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(privateNat.CreatedBy),
		ModifiedAt:        types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, diags := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	state.PrivateNat = privateNatObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPrivateNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.PrivateNatResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Private NAT
	err := r.client.DeletePrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Private NAT",
			"Could not delete Private NAT unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForPrivateNatStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting Private NAT",
			"Error waiting for Private NAT to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForPrivateNatStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetPrivateNat(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.PrivateNat.State), nil
	})
}
