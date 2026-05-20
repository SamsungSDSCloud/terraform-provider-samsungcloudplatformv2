package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &VPCSubnetVipResource{}
	_ resource.ResourceWithConfigure = &VPCSubnetVipResource{}
)

// NewVPCSubnetVipResource is a helper function to simplify the provider implementation.
func NewVPCSubnetVipResource() resource.Resource {
	return &VPCSubnetVipResource{}
}

// VPCSubnetVipResource is the resource implementation.
type VPCSubnetVipResource struct {
	_config *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VPCSubnetVipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vip"
}

// Schema defines the schema for the resource.
func (r *VPCSubnetVipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Subnet Vip",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet Id \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : Subnet VIP Description",
				Optional: true,
			},
			common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
				Description: "Virtual IP Address \n" +
					"  - example : 192.168.20.6",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("SubnetVip"): schema.SingleNestedAttribute{
				Description: "Subnet vip created",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Subnet Vip Id",
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
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State \n" +
							"  - enum : CREATING, ACTIVE, DELETING, DELETED, ERROR",
						Computed: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "Subnet ID",
						Computed:    true,
					},
					common.ToSnakeCase("VipPortId"): schema.StringAttribute{
						Description: "Vip Port Id",
						Computed:    true,
					},
					common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
						Description: "Virtual IP Address \n" +
							"  - example : 192.168.20.6",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("ConnectedPorts"): schema.ListNestedAttribute{
						Description: "Connected Ports",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Connected Port Id",
									Computed:    true,
								},
								common.ToSnakeCase("PortId"): schema.StringAttribute{
									Description: "Port ID",
									Computed:    true,
								},
								common.ToSnakeCase("PortName"): schema.StringAttribute{
									Description: "Port Name",
									Computed:    true,
								},
								common.ToSnakeCase("PortIpAddress"): schema.StringAttribute{
									Description: "Port IP Address",
									Computed:    true,
								},
								common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
									Description: "Connected resource ID",
									Computed:    true,
								},
								common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
									Description: "Connected resource Type",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("StaticNat"): schema.SingleNestedAttribute{
						Description: "Static NAT Info",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
								Description: "Static Nat External Ip Address",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Static Nat Id",
								Computed:    true,
							},
							common.ToSnakeCase("PublicipId"): schema.StringAttribute{
								Description: "Publicip ID",
								Computed:    true,
							},
							common.ToSnakeCase("State"): schema.StringAttribute{
								Description: "Static Nat State",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VPCSubnetVipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VPCSubnetVipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.SubnetVipResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateSubnetVIP(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create VPC Subnet VIP",
			fmt.Sprintf("An error occurred while creating VPC Subnet VIP: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	subnetVip := &vpcv1d2.VpcSubnetVipDetail{
		Id:               types.StringValue(data.SubnetVip.Id),
		CreatedAt:        types.StringValue(data.SubnetVip.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.SubnetVip.CreatedBy),
		ModifiedAt:       types.StringValue(data.SubnetVip.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.SubnetVip.ModifiedBy),
		State:            types.StringValue((string)(*data.SubnetVip.State.Ptr())),
		SubnetId:         types.StringValue(data.SubnetVip.SubnetId),
		VipPortId:        types.StringValue(data.SubnetVip.VipPortId),
		VirtualIpAddress: types.StringValue(data.SubnetVip.VirtualIpAddress),
	}

	if data.SubnetVip.Description.IsSet() {
		if desc := data.SubnetVip.Description.Get(); desc != nil {
			subnetVip.Description = types.StringValue(*desc)
		}
	}

	if data.SubnetVip.ConnectedPorts != nil {
		for _, port := range data.SubnetVip.ConnectedPorts {
			subnetVip.ConnectedPorts = append(subnetVip.ConnectedPorts, vpcv1d2.ConnectedPortInfo{
				Id:                   types.StringValue(port.Id),
				PortId:               types.StringValue(port.PortId),
				PortName:             types.StringValue(port.PortName),
				PortIpAddress:        types.StringValue(port.PortIpAddress),
				AttachedResourceId:   types.StringValue(port.AttachedResourceId),
				AttachedResourceType: types.StringValue(port.AttachedResourceType),
			})
		}
	} else {
		subnetVip.ConnectedPorts = []vpcv1d2.ConnectedPortInfo{}
	}

	if data.SubnetVip.StaticNat.IsSet() {
		resultStaticNat := data.SubnetVip.StaticNat.Get()
		if resultStaticNat != nil {
			subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{
				ExternalIpAddress: types.StringValue(resultStaticNat.ExternalIpAddress),
				Id:                types.StringValue(resultStaticNat.Id),
				PublicipId:        types.StringValue(resultStaticNat.PublicipId),
				State:             types.StringValue(resultStaticNat.State),
			}
		}
	} else {
		subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{}
	}

	subnetVipObjectValue, _ := types.ObjectValueFrom(ctx, subnetVip.AttributeTypes(), subnetVip)
	plan.SubnetVip = subnetVipObjectValue
	plan.Description = subnetVip.Description

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VPCSubnetVipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.SubnetVipResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var subnetVipDetail vpcv1d2.VpcSubnetVipDetail
	errR := state.SubnetVip.As(ctx, &subnetVipDetail, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse VPC Subnet VIP detail",
			fmt.Sprintf("An error occurred while parsing VPC Subnet VIP detail: %s", errR),
		)
		return
	}

	data, err := r.client.ShowSubnetVip(ctx, state.SubnetId.ValueString(), subnetVipDetail.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Subnet VIP",
			"Could not read Subnet VIP ID "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	// Map API response to object
	subnetVip := &vpcv1d2.VpcSubnetVipDetail{
		Id:               types.StringValue(data.SubnetVip.Id),
		CreatedAt:        types.StringValue(data.SubnetVip.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.SubnetVip.CreatedBy),
		ModifiedAt:       types.StringValue(data.SubnetVip.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.SubnetVip.ModifiedBy),
		State:            types.StringValue((string)(*data.SubnetVip.State.Ptr())),
		SubnetId:         types.StringValue(data.SubnetVip.SubnetId),
		VipPortId:        types.StringValue(data.SubnetVip.VipPortId),
		VirtualIpAddress: types.StringValue(data.SubnetVip.VirtualIpAddress),
	}

	if data.SubnetVip.Description.IsSet() {
		if desc := data.SubnetVip.Description.Get(); desc != nil {
			subnetVip.Description = types.StringValue(*desc)
		}
	}

	if data.SubnetVip.ConnectedPorts != nil {
		for _, port := range data.SubnetVip.ConnectedPorts {
			subnetVip.ConnectedPorts = append(subnetVip.ConnectedPorts, vpcv1d2.ConnectedPortInfo{
				Id:                   types.StringValue(port.Id),
				PortId:               types.StringValue(port.PortId),
				PortName:             types.StringValue(port.PortName),
				PortIpAddress:        types.StringValue(port.PortIpAddress),
				AttachedResourceId:   types.StringValue(port.AttachedResourceId),
				AttachedResourceType: types.StringValue(port.AttachedResourceType),
			})
		}
	} else {
		subnetVip.ConnectedPorts = []vpcv1d2.ConnectedPortInfo{}
	}

	if data.SubnetVip.StaticNat.IsSet() {
		resultStaticNat := data.SubnetVip.StaticNat.Get()
		if resultStaticNat != nil {
			subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{
				ExternalIpAddress: types.StringValue(resultStaticNat.ExternalIpAddress),
				Id:                types.StringValue(resultStaticNat.Id),
				PublicipId:        types.StringValue(resultStaticNat.PublicipId),
				State:             types.StringValue(resultStaticNat.State),
			}
		}
	} else {
		subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{}
	}

	subnetVipObjectValue, _ := types.ObjectValueFrom(ctx, subnetVip.AttributeTypes(), subnetVip)
	state.SubnetVip = subnetVipObjectValue
	state.Description = subnetVip.Description

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VPCSubnetVipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d2.SubnetVipResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var subnetVipDetail vpcv1d2.VpcSubnetVipDetail
	errR := state.SubnetVip.As(ctx, &subnetVipDetail, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse VPC Subnet VIP detail",
			fmt.Sprintf("An error occurred while parsing VPC Subnet VIP detail: %s", errR),
		)
		return
	}

	err := r.client.DeleteSubnetVIP(ctx, state.SubnetId.ValueString(), subnetVipDetail.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting subnet",
			"Could not delete subnet Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VPCSubnetVipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state vpcv1d2.SubnetVipResource
	var plan vpcv1d2.SubnetVipResource

	// Retrieve values from plan
	req.Plan.Get(ctx, &plan)
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var subnetVipDetailState vpcv1d2.VpcSubnetVipDetail
	errR := state.SubnetVip.As(ctx, &subnetVipDetailState, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse VPC Subnet VIP detail",
			fmt.Sprintf("An error occurred while parsing VPC Subnet VIP detail: %s", errR),
		)
		return
	}

	// Update existing order
	data, err := r.client.UpdateSubnetVIP(ctx, state.SubnetId.ValueString(), subnetVipDetailState.Id.ValueString(), plan.Description.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating subnet",
			"Could not update subnet Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map API response to object
	subnetVip := &vpcv1d2.VpcSubnetVipDetail{
		Id:               types.StringValue(data.SubnetVip.Id),
		CreatedAt:        types.StringValue(data.SubnetVip.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.SubnetVip.CreatedBy),
		ModifiedAt:       types.StringValue(data.SubnetVip.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.SubnetVip.ModifiedBy),
		State:            types.StringValue((string)(*data.SubnetVip.State.Ptr())),
		SubnetId:         types.StringValue(data.SubnetVip.SubnetId),
		VipPortId:        types.StringValue(data.SubnetVip.VipPortId),
		VirtualIpAddress: types.StringValue(data.SubnetVip.VirtualIpAddress),
	}

	fmt.Printf("data.SubnetVip.Description %v", string(*data.SubnetVip.Description.Get()))
	if data.SubnetVip.Description.IsSet() {
		if desc := data.SubnetVip.Description.Get(); desc != nil {
			subnetVip.Description = types.StringValue(*desc)
		}
	}

	if data.SubnetVip.ConnectedPorts != nil {
		for _, port := range data.SubnetVip.ConnectedPorts {
			subnetVip.ConnectedPorts = append(subnetVip.ConnectedPorts, vpcv1d2.ConnectedPortInfo{
				Id:                   types.StringValue(port.Id),
				PortId:               types.StringValue(port.PortId),
				PortName:             types.StringValue(port.PortName),
				PortIpAddress:        types.StringValue(port.PortIpAddress),
				AttachedResourceId:   types.StringValue(port.AttachedResourceId),
				AttachedResourceType: types.StringValue(port.AttachedResourceType),
			})
		}
	} else {
		subnetVip.ConnectedPorts = []vpcv1d2.ConnectedPortInfo{}
	}

	if data.SubnetVip.StaticNat.IsSet() {
		resultStaticNat := data.SubnetVip.StaticNat.Get()
		if resultStaticNat != nil {
			subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{
				ExternalIpAddress: types.StringValue(resultStaticNat.ExternalIpAddress),
				Id:                types.StringValue(resultStaticNat.Id),
				PublicipId:        types.StringValue(resultStaticNat.PublicipId),
				State:             types.StringValue(resultStaticNat.State),
			}
		}
	} else {
		subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{}
	}

	subnetVipObjectValue, _ := types.ObjectValueFrom(ctx, subnetVip.AttributeTypes(), subnetVip)
	state.SubnetVip = subnetVipObjectValue
	state.Description = plan.Description

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
