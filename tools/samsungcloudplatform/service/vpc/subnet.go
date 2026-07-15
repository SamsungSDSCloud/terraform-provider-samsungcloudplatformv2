package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"regexp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcSubnetResource{}
	_ resource.ResourceWithConfigure = &vpcSubnetResource{}
)

// NewVpcSubnetResource is a helper function to simplify the provider implementation.
func NewVpcSubnetResource() resource.Resource {
	return &vpcSubnetResource{}
}

// vpcSubnetResource is the data source implementation.
type vpcSubnetResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcSubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet"
}

// Schema defines the schema for the data source.
func (r *vpcSubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Subnet resource.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AllocationPools"): schema.ListNestedAttribute{
				Description: "The ranges of IP addresses available for allocation within the subnet.\n" +
					"  - example : [{ \"start\": \"192.168.0.3\", \"end\": \"192.168.0.254\" }]",
				MarkdownDescription: "The ranges of IP addresses available for allocation within the subnet.\n" +
					"  - example : [{ \"start\": \"192.168.0.3\", \"end\": \"192.168.0.254\" }]",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("end"): schema.StringAttribute{
							Description: "The end IP address of the allocation range.\n" +
								"  - example : 192.168.0.1",
							Required: true,
						},
						common.ToSnakeCase("start"): schema.StringAttribute{
							Description: "The start IP address of the allocation range.\n" +
								"  - example : 192.168.0.1",
							Required: true,
						},
					},
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "The identifier of the account that owns the subnet.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				MarkdownDescription: "The identifier of the account that owns the subnet.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				Computed: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "The IP address range of the subnet in CIDR notation.\n" +
					"  - example : 192.168.0.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
				MarkdownDescription: "The IP address range of the subnet in CIDR notation.\n" +
					"  - example : 192.168.0.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
				Required: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "The timestamp when the subnet was created in ISO 8601 format.\n" +
					"  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "The timestamp when the subnet was created in ISO 8601 format.\n" +
					"  - example: 2024-05-17T00:23:17Z",
				Computed: true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "The user id that created the resource.\n" +
					"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "The user id that created the resource.\n" +
					"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
				Computed: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this subnet. This help identify the purpose or usage of the subnet.\n" +
					"  - maxLength: 50\n" +
					"  - example: Subnet Description",
				MarkdownDescription: "Enter a brief explanation or note about this subnet. This help identify the purpose or usage of the subnet.\n" +
					"  - maxLength: 50\n" +
					"  - example: Subnet Description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("dhcp_ip_address"): schema.StringAttribute{
				Computed: true,
				Optional: true,
				Description: "The IP address automatically assigned by DHCP.\n" +
					"  - example: 192.168.0.2",
				MarkdownDescription: "The IP address automatically assigned by DHCP.\n" +
					"  - example: 192.168.0.2",
			},
			common.ToSnakeCase("DnsNameservers"): schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The list of DNS name server addresses for the subnet.\n" +
					"  - example: [\"1.1.1.1\", \"2.2.2.2\"]",
				MarkdownDescription: "The list of DNS name server addresses for the subnet.\n" +
					"  - example: [\"1.1.1.1\", \"2.2.2.2\"]",
			},
			common.ToSnakeCase("GatewayIpAddress"): schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The gateway IP address of the subnet.\n" +
					"  - example: 192.168.0.1",
				MarkdownDescription: "The gateway IP address of the subnet.\n" +
					"  - example: 192.168.0.1",
			},
			common.ToSnakeCase("HostRoutes"): schema.ListNestedAttribute{
				Description: "The static host routes configured for the subnet.\n" +
					"  - example : [{ \"destination\": \"192.168.24.0/24\", \"nexthop\": \"192.168.0.5\" }]",
				MarkdownDescription: "The static host routes configured for the subnet.\n" +
					"  - example : [{ \"destination\": \"192.168.24.0/24\", \"nexthop\": \"192.168.0.5\" }]",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Destination"): schema.StringAttribute{
							Description: "the target IP address range (CIDR) for which the route should be applied\n" +
								"  - example : 192.168.0.1",
							Required: true,
						},
						common.ToSnakeCase("Nexthop"): schema.StringAttribute{
							Description: "The IP address of the next router/VM that the traffic should be sent to in order to reach the destination\n" +
								"  - example : 192.168.0.1",
							Required: true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the subnet.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The unique identifier of the subnet.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description: "The timestamp when the subnet was last modified in ISO 8601 format.\n" +
					"  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "The timestamp when the subnet was last modified in ISO 8601 format.\n" +
					"  - example: 2024-05-17T00:23:17Z",
				Computed: true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "The user id that last modified the resource.\n" +
					"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "The user id that last modified the resource.\n" +
					"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
				Computed: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the subnet.\n" +
					"  - example : subnetName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "The name of the subnet.\n" +
					"  - example : subnetName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9-]*$"), "Enter 3 -20 chars. (English, number, hyphen)"),
				},
				Required: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the subnet.\n" +
					"  - example : ACTIVE",
				MarkdownDescription: "The current lifecycle state of the subnet.\n" +
					"  - example : ACTIVE",
				Computed: true,
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type of the subnet.\n" +
					"  - example : GENERAL | LOCAL | VPC_ENDPOINT",
				MarkdownDescription: "The type of the subnet.\n" +
					"  - example : GENERAL | LOCAL | VPC_ENDPOINT",
				Required: true,
			},
			common.ToSnakeCase("VpcID"): schema.StringAttribute{
				Description: "The identifier of the VPC that the subnet belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The identifier of the VPC that the subnet belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the subnet belongs to.\n" +
					"  - example : VpcName",
				MarkdownDescription: "The name of the VPC that the subnet belongs to.\n" +
					"  - example : VpcName",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcSubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

// Create creates the subnet and sets the initial Terraform state.
func (r *vpcSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.SubnetResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subnet
	data, err := r.client.CreateSubnet(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating subnet",
			"Could not create subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	subnet := data.Subnet
	plan.Id = types.StringValue(subnet.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForSubnetStatus(ctx, r.client, subnet.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subnet",
			"Error waiting for subnet to become active: "+err.Error(),
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
func (r *vpcSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.SubnetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetSubnet(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading subnet",
			"Could not read subnet ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Set refreshed state
	subnet := data.Subnet
	state.Id = types.StringValue(subnet.Id)
	state.AccountId = types.StringValue(subnet.AccountId)
	state.GatewayIpAddress = types.StringPointerValue(subnet.GatewayIpAddress.Get())
	state.VpcName = types.StringValue(subnet.VpcName)
	state.Description = types.StringPointerValue(subnet.Description.Get())
	state.State = types.StringValue(string(subnet.State))
	state.CreatedAt = types.StringValue(subnet.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(subnet.CreatedBy)
	state.ModifiedAt = types.StringValue(subnet.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(subnet.ModifiedBy)
	state.DhcpIpAddress = types.StringPointerValue(subnet.DhcpIpAddress.Get())

	dnsNameservers, _ := types.SetValueFrom(ctx, types.StringType, subnet.GetDnsNameservers())

	state.DnsNameservers = dnsNameservers

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the subnet and sets the updated Terraform state on success.
func (r *vpcSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.SubnetResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateSubnet(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating subnet",
			"Could not update subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
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

// Delete deletes the subnet and removes the Terraform state on success.
func (r *vpcSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.SubnetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subnet
	err := r.client.DeleteSubnet(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting subnet",
			"Could not delete subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForSubnetStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting subnet",
			"Error waiting for subnet to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForSubnetStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetSubnet(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Subnet.State), nil
	}, -1, -1, -1, -1)
}
