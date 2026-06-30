package vpc

import (
	"context"
	"fmt"

	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcPortResource{}
	_ resource.ResourceWithConfigure   = &vpcPortResource{}
	_ resource.ResourceWithImportState = &vpcPortResource{}
)

// NewVpcPortResource is a helper function to simplify the provider implementation.
func NewVpcPortResource() resource.Resource {
	return &vpcPortResource{}
}

// vpcPortResource is the data source implementation.
type vpcPortResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcPortResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_port"
}

// Schema defines the schema for the data source.
func (r *vpcPortResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Port resource for network",
		MarkdownDescription: "Port",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the port.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The unique identifier of the port.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "The identifier of the account that owns the port.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				MarkdownDescription: "The identifier of the account that owns the port.\n" +
					"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
				Computed: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "The identifier of the resource that this port is attached to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The identifier of the resource that this port is attached to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description: "The type of the resource that this port is attached to.\n" +
					"  - example : VM",
				MarkdownDescription: "The type of the resource that this port is attached to.\n" +
					"  - example : VM",
				Computed: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
					"  - example : 2024-05-17T00:23:17Z",
				MarkdownDescription: "The timestamp when the resource was created, in ISO 8601 format.\n" +
					"  - example : 2024-05-17T00:23:17Z",
				Computed: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : Port Description\n" +
					"  - maxLength : 50",
				MarkdownDescription: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : Port Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Default: stringdefault.StaticString(""),
			},
			common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
				Description: "The fixed IP address assigned to the port.\n" +
					"  - example : 172.24.4.2",
				MarkdownDescription: "The fixed IP address assigned to the port.\n" +
					"  - example : 172.24.4.2",
				Optional: true,
			},
			common.ToSnakeCase("MacAddress"): schema.StringAttribute{
				Description: "The MAC address of the port.\n" +
					"  - example : fa:16:3e:5c:9b:7a",
				MarkdownDescription: "The MAC address of the port.\n" +
					"  - example : fa:16:3e:5c:9b:7a",
				Computed: true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
					"  - example : 2024-05-17T00:23:17Z",
				MarkdownDescription: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
					"  - example : 2024-05-17T00:23:17Z",
				Computed: true,
			},
			common.ToSnakeCase("SecurityGroups"): schema.ListNestedAttribute{
				Description:         "The list of security groups associated with the port.",
				MarkdownDescription: "The list of security groups associated with the port.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the port.\n" +
					"  - example : portName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "The name of the port.\n" +
					"  - example : portName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Required: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the port.\n" +
					"  - example : ACTIVE",
				MarkdownDescription: "The current lifecycle state of the port.\n" +
					"  - example : ACTIVE",
				Computed: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the port belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				MarkdownDescription: "The identifier of the subnet that the port belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("SubnetName"): schema.StringAttribute{
				Description: "The name of the subnet that the port belongs to.\n" +
					"  - example : subnetName",
				MarkdownDescription: "The name of the subnet that the port belongs to.\n" +
					"  - example : subnetName",
				Computed: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the port belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The identifier of the VPC that the port belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the port belongs to.\n" +
					"  - example : vpcName",
				MarkdownDescription: "The name of the VPC that the port belongs to.\n" +
					"  - example : vpcName",
				Computed: true,
			},
			common.ToSnakeCase("VirtualIpAddresses"): schema.ListAttribute{
				Description: "Virtual IP Addresses\n" +
					"  - example : [\"192.168.1.100\"]",
				MarkdownDescription: "Virtual IP Addresses\n" +
					"  - example : [\"192.168.1.100\"]",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcPortResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.PortResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new port
	data, err := r.client.CreatePort(ctx, plan)
	if err != nil || !data.Port.IsSet() {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating port",
			"Could not create port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpc.MapPort(data.Port.Get(), &plan)
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.PortResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetPort(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading port",
			"Could not read port ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil || !data.Port.IsSet() {
		resp.Diagnostics.AddError(
			"Error reading data",
			"An error occurred while reading data. Empty response",
		)
		return
	}

	// Set refreshed state
	vpc.MapPort(data.Port.Get(), &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.PortResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	data, err := r.client.UpdatePort(ctx, state.Id.ValueString(), state)
	if err != nil || !data.Port.IsSet() {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating port",
			"Could not update port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpc.MapPort(data.Port.Get(), &state)
	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.PortResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing port
	err := r.client.DeletePort(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting port",
			"Could not delete port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *vpcPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))
}
