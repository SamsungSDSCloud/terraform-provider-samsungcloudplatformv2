package directconnect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpdirectconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/direct-connect/1.0"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &directConnectDirectConnectResource{}
	_ resource.ResourceWithConfigure   = &directConnectDirectConnectResource{}
	_ resource.ResourceWithImportState = &directConnectDirectConnectResource{}
)

// NewDirectConnectDirectConnectResource is a helper function to simplify the provider implementation.
func NewDirectConnectDirectConnectResource() resource.Resource {
	return &directConnectDirectConnectResource{}
}

// directConnectDirectConnectResource is the data source implementation.
type directConnectDirectConnectResource struct {
	config  *scpsdk.Configuration
	client  *directconnect.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *directConnectDirectConnectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directconnect_direct_connect"
}

// Schema defines the schema for the data source.
func (r *directConnectDirectConnectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Direct Connect for connect between a customer's premises and a vpc.",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the direct connect.\n" +
					"  - example : fe860e0af0c04dcd8182b84f907f31f4",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
				Description: "The bandwidth capacity(1Gpbs, 10Gpbs, 20Gpbs or 40Gpbs) of the connection.\n" +
					"  - example : 1 | 10 | 20 | 40",
				Required: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this direct connect. This help identify the purpose or usage of the resource. \n" +
					"  - example : Direct Connect description\n" +
					"  - maxLength : 50\n" +
					"  - minLength : 1",
				Optional: true,
			},
			common.ToSnakeCase("FirewallEnabled"): schema.BoolAttribute{
				Description: "Whether the firewall is enabled for the direct connect.(firewall Enable : true, firewall Diable : false)\n" +
					"  - example : true | false",
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("FirewallLoggable"): schema.BoolAttribute{
				Description: "Whether firewall logging is enabled for the direct connect.(firewall logging Enable : true, firewall logging Diable : false) \n" +
					"  - example : true | false",
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the direct connect.\n" +
					"- example : directConnectName",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the direct connect belongs to.\n" +
					"- example : 023c57b14f11483689338d085e061492",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DirectConnect"): schema.SingleNestedAttribute{
				Description: "DirectConnect",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Identifier of the direct connect.\n" +
							"  - example : fe860e0af0c04dcd8182b84f907f31f4",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the direct connect.\n" +
							"  - example : directConnectName",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the direct connect.\n " +
							"  - example: 27bb070b564349f8a31cc60734cc36a5",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this direct connect. This help identify the purpose or usage of the resource. \n" +
							"  - example : Direct Connect description\n" +
							"  - maxLength : 50\n" +
							"  - minLength : 1",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the direct connect belongs to.\n" +
							"- example : 023c57b14f11483689338d085e061492",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the direct connect belongs to.\n" +
							"  - example : vpc-prod-01",
						Computed: true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "The bandwidth capacity(1Gpbs, 10Gpbs, 20Gpbs or 40Gpbs) of the connection.\n" +
							"  - example : 1 | 10 | 20 | 40",
						Computed: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The identifier of the firewall associated with the direct connect.\n" +
							"  - example: 68db67f78abd405da98a6056a8ee42af",
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
						Description: "The current lifecycle state of the direct connect. \n" +
							"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *directConnectDirectConnectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.DirectConnect
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *directConnectDirectConnectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan directconnect.DirectConnectResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new direct connect
	data, err := r.client.CreateDirectConnect(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating direct connect",
			"Could not create direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.DirectConnect.Id)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.DirectConnect = dconObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForDirectConnectStatus(ctx, r.client, data.DirectConnect.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating direct connect",
			"Error waiting for direct connect to become active: "+err.Error(),
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

func (r *directConnectDirectConnectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *directConnectDirectConnectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state directconnect.DirectConnectResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from direct connect
	data, err := r.client.GetDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading direct connect",
			"Could not read direct connect ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.DirectConnect = dconObjectValue

	// Update top-level input fields from API response to detect drift
	state.Name = types.StringValue(data.DirectConnect.Name)
	state.Description = types.StringPointerValue(data.DirectConnect.Description.Get())
	state.Bandwidth = types.Int32Value(data.DirectConnect.Bandwidth)
	state.VpcId = types.StringValue(data.DirectConnect.VpcId)
	if data.DirectConnect.FirewallId.IsSet() && data.DirectConnect.FirewallId.Get() != nil && *data.DirectConnect.FirewallId.Get() != "" {
		state.FirewallEnabled = types.BoolValue(true)
	} else {
		state.FirewallEnabled = types.BoolValue(false)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *directConnectDirectConnectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state directconnect.DirectConnectResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateDirectConnect(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating direct connect",
			"Could not update direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetDirectConnect as UpdateDirectConnect items are not populated.
	data, err := r.client.GetDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading direct connect",
			"Could not read direct connect ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.DirectConnect = dconObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *directConnectDirectConnectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state directconnect.DirectConnectResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing direct connect
	err := r.client.DeleteDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting direct connect",
			"Could not delete direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForDirectConnectStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting direct connect",
			"Error waiting for direct connect to become deleted: "+err.Error(),
		)
		return
	}
}

func createDirectConnectModel(data *scpdirectconnect.DirectConnectShowResponse) directconnect.DirectConnect {
	dcon := data.DirectConnect
	return directconnect.DirectConnect{
		Id:          types.StringValue(dcon.Id),
		Name:        types.StringValue(dcon.Name),
		AccountId:   types.StringValue(dcon.AccountId),
		Description: types.StringPointerValue(dcon.Description.Get()),
		VpcId:       types.StringValue(dcon.VpcId),
		VpcName:     types.StringValue(dcon.VpcName),
		Bandwidth:   types.Int32Value(dcon.Bandwidth),
		FirewallId:  types.StringPointerValue(dcon.FirewallId.Get()),
		CreatedAt:   types.StringValue(dcon.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(dcon.CreatedBy),
		ModifiedAt:  types.StringValue(dcon.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(dcon.ModifiedBy),
		State:       types.StringValue(string(dcon.State)),
	}
}

func waitForDirectConnectStatus(ctx context.Context, directConnectClient *directconnect.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := directConnectClient.GetDirectConnect(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.DirectConnect.State), nil
	}, -1, -1, -1, -1)
}
