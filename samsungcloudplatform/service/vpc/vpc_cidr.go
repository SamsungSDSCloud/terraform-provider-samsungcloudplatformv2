package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpcV1Dot2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &VpcCidrResource{}
	_ resource.ResourceWithConfigure = &VpcCidrResource{}
)

// NewVpcCidrResource is a helper function to simplify the provider implementation.
func NewVpcCidrResource() resource.Resource {
	return &VpcCidrResource{}
}

// VpcCidrResource is the resource implementation.
type VpcCidrResource struct {
	_config *scpsdk.Configuration
	client  *vpcV1Dot2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VpcCidrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_cidr"
}

// Schema defines the schema for the resource.
func (r *VpcCidrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC CIDR",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "CIDR \n" +
					"  - example : 192.168.0.0/24",
				Required: true,
			},

			// Output
			common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
				Description: "VPC detail after adding CIDR",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("CidrCount"): schema.Int32Attribute{
						Description: "CIDR Count",
						Computed:    true,
					},
					common.ToSnakeCase("Cidrs"): schema.ListNestedAttribute{
						Description: "CIDRs",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "CIDR ID",
									Computed:    true,
								},
								common.ToSnakeCase("Cidr"): schema.StringAttribute{
									Description: "CIDR",
									Computed:    true,
								},
								common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
									Description: "Created At",
									Computed:    true,
								},
								common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
									Description: "Created By",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VpcCidrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *VpcCidrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcV1Dot2.VpcCidrResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.AddVpcCidr(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to add VPC CIDR",
			fmt.Sprintf("An error occurred while creating VPC CIDR: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	vpcCidr := &vpcV1Dot2.VpcCidrDetail{
		Id:         types.StringValue(data.Vpc.Id),
		Name:       types.StringValue(data.Vpc.Name),
		AccountId:  types.StringValue(data.Vpc.AccountId),
		State:      types.StringValue(string(data.Vpc.State)),
		CidrCount:  types.Int32Value(data.Vpc.CidrCount),
		CreatedAt:  types.StringValue(data.Vpc.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringValue(data.Vpc.CreatedBy),
		ModifiedAt: types.StringValue(data.Vpc.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringValue(data.Vpc.ModifiedBy),
	}

	if data.Vpc.Description.IsSet() {
		if desc := data.Vpc.Description.Get(); desc != nil {
			vpcCidr.Description = types.StringValue(*desc)
		}
	}

	if data.Vpc.Cidrs != nil {
		for _, cidr := range data.Vpc.Cidrs {
			vpcCidr.Cidrs = append(vpcCidr.Cidrs, vpcV1Dot2.VpcCidrInfo{
				Id:        types.StringValue(cidr.Id),
				Cidr:      types.StringValue(cidr.Cidr),
				CreatedAt: types.StringValue(cidr.CreatedAt.Format(time.RFC3339)),
				CreatedBy: types.StringValue(cidr.CreatedBy),
			})
		}
	} else {
		vpcCidr.Cidrs = []vpcV1Dot2.VpcCidrInfo{}
	}

	vpcCidrObjectValue, _ := types.ObjectValueFrom(ctx, vpcCidr.AttributeTypes(), vpcCidr)
	plan.Vpc = vpcCidrObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VpcCidrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO: Implement Read function when needed, remove resource for now

	// resp.Diagnostics.AddError(
	// 	"Read Not Implemented",
	// 	"VPC CIDR Read function is not yet implemented.",
	// )
	resp.State.RemoveResource(ctx)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VpcCidrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: Implement Update function
	resp.Diagnostics.AddError(
		"Update Not Implemented",
		"VPC CIDR Update function is not yet implemented.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VpcCidrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO: Implement Delete function
	resp.Diagnostics.AddError(
		"Delete Not Implemented",
		"VPC CIDR Delete function is not yet implemented.",
	)
}
