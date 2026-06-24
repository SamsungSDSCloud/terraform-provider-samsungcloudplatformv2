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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcVpcEndpointResource{}
	_ resource.ResourceWithConfigure = &vpcVpcEndpointResource{}
)

// NewVpcVpcEndpointResource is a helper function to simplify the provider implementation.
func NewVpcVpcEndpointResource() resource.Resource {
	return &vpcVpcEndpointResource{}
}

// vpcVpcEndpointResource is the data source implementation.
type vpcVpcEndpointResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcVpcEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_endpoint"
}

// Schema defines the schema for the data source.
func (r *vpcVpcEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource of vpcendpoint",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "The unique identifier of the endpoint.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the endpoint.\n" +
					"  - example : vpcEndpointName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Required: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the endpoint belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the endpoint belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "The type of the target resource.(File Storage : FS, Object Storage : OBS, Container Registry : SCR, DNS : DNS)\n" +
					"  - example : FS | OBS | SCR | DNS",
				Required: true,
			},
			common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
				Description: "The key identifying the target resource of the endpoint.\n" +
					"  - example(case: SCR/DNS) : 07c5364702384471b650147321b52173 \n" +
					"  - example(case: FS/OBS) : 1.1.1.1",
				Required: true,
			},
			common.ToSnakeCase("ResourceInfo"): schema.StringAttribute{
				Description: "The information about the target resource of the endpoint.\n" +
					"  - example(case: FS) : 192.168.0.1(SSD) \n" +
					"  - example(case: OBS) : https://xxx.samsungsdscloud.com \n" +
					"  - example(case: SCR) : xxx.samsungsdscloud.com(Auth) \n" +
					"  - example(case: DNS) : Private DNS Name",
				Required: true,
			},
			common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
				Description: "The IP address of the endpoint. \n" +
					"  - example : 10.10.10.10",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
					"  - example : VPC Endpoint description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("VpcEndpoint"): schema.SingleNestedAttribute{
				Description: "VpcEndpoint",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the endpoint.\n" +
							"  - example : 12f56e27070248a6a240a497e43fbe18",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the endpoint.\n" +
							"  - example : VpcEndpointName",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the endpoint belongs to.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the endpoint belongs to.\n" +
							"  - example : vpcName",
						Computed: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The identifier of the subnet that the endpoint belongs to.\n" +
							"  - example : 023c57b14f11483689338d085e061492",
						Computed: true,
					},
					common.ToSnakeCase("SubnetName"): schema.StringAttribute{
						Description: "The name of the subnet that the endpoint belongs to.\n" +
							"  - example : subnetName",
						Computed: true,
					},
					common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
						Description: "The IP address of the endpoint.\n" +
							"  - example : 192.167.0.5",
						Computed: true,
					},
					common.ToSnakeCase("ResourceType"): schema.StringAttribute{
						Description: "The type of the target resource.\n" +
							"  - example : FS",
						Computed: true,
					},
					common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
						Description: "The key identifying the target resource of the endpoint.\n" +
							"  - example : 07c5364702384471b650147321b52173",
						Computed: true,
					},
					common.ToSnakeCase("ResourceInfo"): schema.StringAttribute{
						Description: "The information about the target resource of the endpoint.\n" +
							"  - example : x.samsungsdscloud.com(Registry)",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the endpoint.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the endpoint.\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
							"  - example : VpcEndpoint Description",
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcVpcEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcVpcEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.VpcEndpointResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpcEndpoint(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc endpoint",
			"Could not create vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(vpcendpoint.Id)

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		ResourceKey:  types.StringValue(vpcendpoint.AccountId),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	plan.VpcEndpoint = vpcEndpointObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForVpcEndpointStatus(ctx, r.client, vpcendpoint.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc endpoint",
			"Error waiting for vpc endpoint to become active: "+err.Error(),
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
func (r *vpcVpcEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.VpcEndpointResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc endpoint",
			"Could not read vpc endpoint ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		ResourceKey:  types.StringValue(vpcendpoint.AccountId),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	state.VpcEndpoint = vpcEndpointObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcVpcEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.VpcEndpointResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpcEndpoint(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc endpoint",
			"Could not update vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpcEndpoint as UpdateVpc items are not populated.
	data, err := r.client.GetVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc endpoint",
			"Could not read vpc endpoint ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		ResourceKey:  types.StringValue(vpcendpoint.AccountId),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	state.VpcEndpoint = vpcEndpointObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcVpcEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.VpcEndpointResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing VpcEndpoint
	err := r.client.DeleteVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc endpoint",
			"Could not delete vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcEndpointStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc endpoint",
			"Error waiting for vpc endpoint to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForVpcEndpointStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetVpcEndpoint(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpcEndpoint.State), nil
	}, -1, -1, -1, -1)
}
