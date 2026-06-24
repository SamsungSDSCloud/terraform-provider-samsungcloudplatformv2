package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcInternetGatewayResource{}
	_ resource.ResourceWithConfigure   = &vpcInternetGatewayResource{}
	_ resource.ResourceWithImportState = &vpcInternetGatewayResource{}
)

// NewVpcInternetGatewayResource is a helper function to simplify the provider implementation.
func NewVpcInternetGatewayResource() resource.Resource {
	return &vpcInternetGatewayResource{}
}

// vpcInternetGatewayResource is the data source implementation.
type vpcInternetGatewayResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcInternetGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_internet_gateway"
}

// Schema defines the schema for the data source.
func (r *vpcInternetGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Internet Gateway resource for Internet traffic from vpc.",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the internet gateway.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type of the internet gateway.GGW is only supported on SCP for Samsung. SIGW is only supported on SCP for Enterprise.\n" +
					"  - example : IGW | GGW | SIGW",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
					"  - example : Internet Gateway Description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			}, common.ToSnakeCase("Loggable"): schema.BoolAttribute{
				Description: "Whether logging is enabled for the NAT.(NAT logging Enable : true, NAT logging Diable : false) \n" +
					"  - example : true | false",
				Optional: true,
			},
			common.ToSnakeCase("FirewallEnabled"): schema.BoolAttribute{
				Description: "Whether the firewall is enabled for the internet gateway.(Enable : true, Disable : false)\n" +
					"  - example : true | false",
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("FirewallLoggable"): schema.BoolAttribute{
				Description: "Whether firewall logging is enabled for the internet gateway.(Enable : true, Disable : false)\n" +
					"  - example : true | false",
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the internet gateway belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("InternetGateway"): schema.SingleNestedAttribute{
				Description: "InternetGateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the internet gateway.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the internet gateway.\n" +
							"  - example : my-internet-gateway",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the internet gateway.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of the internet gateway.GGW is only supported on SCP for Samsung. SIGW is only supported on SCP for Enterprise.\n" +
							"  - example : IGW | GGW | SIGW",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Internet Gateway Description",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the internet gateway belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the internet gateway belongs to.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "Whether logging is enabled for the NAT.(NAT logging Enable : true, NAT logging Diable : false)\n" +
							"  - example : true",
						Computed: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The identifier of the firewall associated with the internet gateway.\n" +
							"  - example : 68db67f78abd405da98a6056a8ee42af",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the internet gateway.\n" +
							"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcInternetGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcInternetGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.InternetGatewayResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new internet gateway
	data, err := r.client.CreateInternetGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating internet gateway",
			"Could not create internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(data.InternetGateway.Id)

	igwModel := createInternetGatewayModel(data)

	igwObjectValue, diags := types.ObjectValueFrom(ctx, igwModel.AttributeTypes(), igwModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.InternetGateway = igwObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForInternetGatewayStatus(ctx, r.client, data.InternetGateway.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating internet gateway",
			"Error waiting for internet gateway to become active: "+err.Error(),
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
func (r *vpcInternetGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.InternetGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from internet gateway
	data, err := r.client.GetInternetGateway(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading internet gateway",
			"Could not read internet gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error reading data",
			"An error occurred while reading data. Empty response",
		)
		return
	}

	// Refresh input attributes from API response
	state.Type = types.StringValue(string(data.InternetGateway.Type))
	state.Description = types.StringPointerValue(data.InternetGateway.Description.Get())
	state.Loggable = types.BoolValue(data.InternetGateway.GetLoggable())
	state.VpcId = types.StringValue(data.InternetGateway.VpcId)

	igwModel := createInternetGatewayModel(data)

	igwObjectValue, diags := types.ObjectValueFrom(ctx, igwModel.AttributeTypes(), igwModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.InternetGateway = igwObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcInternetGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.InternetGatewayResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateInternetGateway(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating internet gateway",
			"Could not update internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetInternetGateway as UpdateInternetGateway items are not populated.
	data, err := r.client.GetInternetGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading internet gateway",
			"Could not read internet gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	igwModel := createInternetGatewayModel(data)

	igwObjectValue, diags := types.ObjectValueFrom(ctx, igwModel.AttributeTypes(), igwModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.InternetGateway = igwObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcInternetGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.InternetGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing internet gateway
	err := r.client.DeleteInternetGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting internet gateway",
			"Could not delete internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForInternetGatewayStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting internet gateway",
			"Error waiting for internet gateway to become deleted: "+err.Error(),
		)
		return
	}
}

func createInternetGatewayModel(data *scpvpc.InternetGatewayShowResponse) vpc.InternetGateway {
	igw := data.InternetGateway
	return vpc.InternetGateway{
		Id:          types.StringValue(igw.Id),
		Name:        types.StringValue(igw.Name),
		AccountId:   types.StringValue(igw.AccountId),
		Description: types.StringPointerValue(igw.Description.Get()),
		VpcId:       types.StringValue(igw.VpcId),
		VpcName:     types.StringValue(igw.VpcName),
		Type:        types.StringValue(string(igw.Type)),
		Loggable:    types.BoolValue(igw.GetLoggable()),
		FirewallId:  types.StringPointerValue(igw.FirewallId.Get()),
		CreatedAt:   types.StringValue(igw.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(igw.CreatedBy),
		ModifiedAt:  types.StringValue(igw.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(igw.ModifiedBy),
		State:       types.StringValue(string(igw.State)),
	}
}

func waitForInternetGatewayStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetInternetGateway(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.InternetGateway.State), nil
	}, -1, -1, -1, -1)
}

func (r *vpcInternetGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 1 || parts[0] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: internet_gateway_id, got: %q", req.ID),
		)
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[0]))
}
