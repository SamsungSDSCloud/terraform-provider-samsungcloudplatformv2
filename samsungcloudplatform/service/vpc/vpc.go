package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcVpcResource{}
	_ resource.ResourceWithConfigure = &vpcVpcResource{}
)

// NewVpcVpcResource is a helper function to simplify the provider implementation.
func NewVpcVpcResource() resource.Resource {
	return &vpcVpcResource{}
}

// vpcVpcResource is the data source implementation.
type vpcVpcResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcVpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc"
}

// Schema defines the schema for the data source.
func (r *vpcVpcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "vpc",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "VPC CIDR\n" +
					"  - example : 192.167.0.0/18\n" +
					"  - maxMask : /24\n" +
					"  - minMask : /16",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : VPC description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Required: true,
			},
			common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
				Description: "Vpc",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Cidr"): schema.StringAttribute{
						Description: "Cidr",
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
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcVpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcVpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.VpcResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpc(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc",
			"Could not create vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcElement := data.Vpc
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(vpcElement.Id)

	vpcModel := vpc.Vpc{
		Cidr:        types.StringValue(vpcElement.Cidr),
		CreatedAt:   types.StringValue(vpcElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(vpcElement.CreatedBy),
		Description: types.StringPointerValue(vpcElement.Description.Get()),
		Id:          types.StringValue(vpcElement.Id),
		ModifiedAt:  types.StringValue(vpcElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(vpcElement.ModifiedBy),
		Name:        types.StringValue(vpcElement.Name),
		AccountId:   types.StringValue(vpcElement.AccountId),
		State:       types.StringValue(string(vpcElement.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	plan.Vpc = vpcObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcVpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcElement := data.Vpc

	vpcModel := vpc.Vpc{
		Cidr:        types.StringValue(vpcElement.Cidr),
		CreatedAt:   types.StringValue(vpcElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(vpcElement.CreatedBy),
		Description: types.StringPointerValue(vpcElement.Description.Get()),
		Id:          types.StringValue(vpcElement.Id),
		ModifiedAt:  types.StringValue(vpcElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(vpcElement.ModifiedBy),
		Name:        types.StringValue(vpcElement.Name),
		AccountId:   types.StringValue(vpcElement.AccountId),
		State:       types.StringValue(string(vpcElement.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	state.Vpc = vpcObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcVpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.VpcResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpc(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc",
			"Could not update vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcElement := data.Vpc

	vpcModel := vpc.Vpc{
		Cidr:        types.StringValue(vpcElement.Cidr),
		CreatedAt:   types.StringValue(vpcElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(vpcElement.CreatedBy),
		Description: types.StringPointerValue(vpcElement.Description.Get()),
		Id:          types.StringValue(vpcElement.Id),
		ModifiedAt:  types.StringValue(vpcElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(vpcElement.ModifiedBy),
		Name:        types.StringValue(vpcElement.Name),
		AccountId:   types.StringValue(vpcElement.AccountId),
		State:       types.StringValue(string(vpcElement.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	state.Vpc = vpcObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcVpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpc
	err := r.client.DeleteVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc",
			"Could not delete vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
