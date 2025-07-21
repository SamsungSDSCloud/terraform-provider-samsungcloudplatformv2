package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcPrivateNatIpResource{}
	_ resource.ResourceWithConfigure = &vpcPrivateNatIpResource{}
)

// NewVpcPrivateNatIpResource is a helper function to simplify the provider implementation.
func NewVpcPrivateNatIpResource() resource.Resource {
	return &vpcPrivateNatIpResource{}
}

// vpcPrivateNatIpResource is the data source implementation.
type vpcPrivateNatIpResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nat_ip"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Private NAT IP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
				Description: "Private NAT ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "IP address \n" +
					"  - example : 192.168.10.0",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Private NAT IP description\n" +
					"  - maxLength : 50",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNatIp"): schema.SingleNestedAttribute{
				Description: "Private NAT IP",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "IpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
						Description: "PrivateNatId",
						Computed:    true,
					},
					common.ToSnakeCase("PrivateNatName"): schema.StringAttribute{
						Description: "PrivateNatName",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
						Description: "AttachedResourceName",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
						Description: "AttachedResourceType",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
						Description: "AttachedResourceId",
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
func (d *vpcPrivateNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcPrivateNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.PrivateNatIpResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreatePrivateNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Private NAT IP",
			"Could not create Private NAT IP, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNatIp := data.PrivateNatIp
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(privateNatIp.Id)

	ptr := privateNatIp.AttachedResourceType.Get()
	var attachedResourceType *string
	if ptr != nil {
		s := string(*ptr)
		attachedResourceType = &s
	} else {
		attachedResourceType = nil
	}
	privateNatIpModel := vpc.PrivateNatIp{
		Id:                   types.StringValue(privateNatIp.Id),
		IpAddress:            types.StringValue(privateNatIp.IpAddress),
		PrivateNatId:         types.StringValue(privateNatIp.PrivateNatId),
		AttachedResourceName: types.StringPointerValue(privateNatIp.AttachedResourceName.Get()),
		AttachedResourceType: types.StringPointerValue(attachedResourceType),
		AttachedResourceId:   types.StringPointerValue(privateNatIp.AttachedResourceId.Get()),
		State:                types.StringValue(string(privateNatIp.State)),
		CreatedAt:            types.StringValue(privateNatIp.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(privateNatIp.CreatedBy),
		ModifiedAt:           types.StringValue(privateNatIp.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(privateNatIp.ModifiedBy),
	}
	privateNatIpObjectValue, diags := types.ObjectValueFrom(ctx, privateNatIpModel.AttributeTypes(), privateNatIpModel)
	plan.PrivateNatIp = privateNatIpObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcPrivateNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *vpcPrivateNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPrivateNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.PrivateNatIpResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Private NAT
	err := r.client.DeletePrivateNatIp(ctx, state.PrivateNatId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Private NAT IP",
			"Could not delete Private NAT IP unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
