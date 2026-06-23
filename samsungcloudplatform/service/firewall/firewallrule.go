package firewall

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/firewall/1.0"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallFirewallRuleResource{}
	_ resource.ResourceWithConfigure   = &firewallFirewallRuleResource{}
	_ resource.ResourceWithImportState = &firewallFirewallRuleResource{}
)

// NewFirewallFirewallRuleResource is a helper function to simplify the provider implementation.
func NewFirewallFirewallRuleResource() resource.Resource {
	return &firewallFirewallRuleResource{}
}

// networkFirewallRuleResource is the data source implementation.
type firewallFirewallRuleResource struct {
	config *scpsdk.Configuration
	client *firewall.Client
}

// Metadata returns the firewallFirewallRuleResource source type name.
func (r *firewallFirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall_rule"
}

// Schema defines the schema for the data source.
func (r *firewallFirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages firewall rules to control network traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 0e2b4ece64944d7d8a72983e945b867b",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "The identifier of the firewall associated with the resource.\n" +
					"  - example: 68db67f78abd405da98a6056a8ee42af",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("FirewallRule"): schema.SingleNestedAttribute{
				Description: "Firewall Rule.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 0e2b4ece64944d7d8a72983e945b867b",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
							"  - example: 0e2b4ece64944d7d8a72983e945b867b",
						Computed: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The identifier of the firewall associated with the resource.\n" +
							"  - example: 68db67f78abd405da98a6056a8ee42af",
						Computed: true,
					},
					common.ToSnakeCase("Sequence"): schema.Int32Attribute{
						Description: "The order in which the rule is evaluated.\n" +
							"  - example: 100",
						Computed: true,
					},
					common.ToSnakeCase("SourceInterface"): schema.StringAttribute{
						Description: "The source interface the rule applies to.\n" +
							"  - example: L2FW-DGW2800up",
						Computed: true,
					},
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "The source IP addresses the rule applies to.\n" +
							"  - example: [10.10.10.0/24, 10.10.11.0/24]",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationInterface"): schema.StringAttribute{
						Description: "The destination interface the rule applies to.\n" +
							"  - example: L2FW-DGW2800dn",
						Computed: true,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "The destination address the rule applies to.\n" +
							"  - example: [192.168.0.0/16, 192.169.0.0/16]",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "The service ports the rule applies to.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "The type of the service.\n" +
										"  - example: TCP",
									Computed: true,
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "The value of the service.\n" +
										"  - example: 80",
									Computed: true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "The action applied to traffic that matches the rule.\n" +
							"  - example: ALLOW",
						Computed: true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "The direction of the traffic the rule applies to.\n" +
							"  - example: INBOUND",
						Computed: true,
					},
					common.ToSnakeCase("VendorRuleId"): schema.StringAttribute{
						Description: "The firewall device's unique identifier for the rule.\n" +
							"  - example: 20",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Firewall rule for web tier",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the resource.\n" +
							"  - example: ENABLE",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
				},
			},
			common.ToSnakeCase("FirewallRuleCreate"): schema.SingleNestedAttribute{
				Description: "Firewall rule create object",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "The source IP addresses the rule applies to.\n" +
							"  - example: [10.10.10.0/24, 10.10.11.0/24]",
						Required:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "The destination address the rule applies to.\n" +
							"  - example: [192.168.0.0/16, 192.169.0.0/16]",
						Required:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "The service ports the rule applies to.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "The type of the service.\n" +
										"  - example: TCP\n" +
										"  - valid: TCP, UDP, ICMP, IP, TCP_ALL, UDP_ALL, ICMP_ALL, ALL",
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("TCP", "UDP", "ICMP", "IP", "TCP_ALL", "UDP_ALL", "ICMP_ALL", "ALL"),
									},
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "The value of the service.\n" +
										"  - example: 80",
									Optional: true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "The action applied to traffic that matches the rule.\n" +
							"  - example: ALLOW\n" +
							"  - valid: ALLOW, DENY",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ALLOW", "DENY"),
						},
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "The direction of the traffic the rule applies to.\n" +
							"  - example: INBOUND\n" +
							"  - valid: INBOUND, OUTBOUND",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("INBOUND", "OUTBOUND"),
						},
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the resource.\n" +
							"  - example: ENABLE\n" +
							"  - valid: ENABLE, DISABLE",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ENABLE", "DISABLE"),
						},
					},
					common.ToSnakeCase("OrderRuleId"): schema.StringAttribute{
						Description: "The target ID used when changing the order of a rule.\n" +
							"  - example: 7087c92d295445cda2785a94aab93c65",
						Optional: true,
					},
					common.ToSnakeCase("OrderDirection"): schema.StringAttribute{
						Description: "The type of ordering change applied to the rule.\n" +
							"  - example: BEFORE\n" +
							"  - valid: BEFORE, AFTER, BOTTOM",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("BEFORE", "AFTER", "BOTTOM"),
						},
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Firewall rule for web tier\n" +
							"  - constraints: maxLength: 100",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(100),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *firewallFirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Firewall
}

// Create creates the resource and sets the initial Terraform state.
func (r *firewallFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan firewall.FirewallRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new firewall rule
	data, err := r.client.CreateFirewallRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating firewall rule",
			"Could not create firewall rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError("Error creating firewall rule", "Received nil response from API")
		return
	}
	// Wait for the firewall rule to reach ACTIVE state before setting Terraform state
	if err := waitForFirewallRuleStatus(ctx, r.client, data.FirewallRule.Id, []string{}, []string{"ACTIVE"}); err != nil {
		resp.Diagnostics.AddError(
			"Error creating firewall rule",
			"Could not wait for firewall rule to become active: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(data.FirewallRule.Id)

	firewallRuleModel := createFirewallRuleModel(data)
	firewallRuleObjectValue, objDiags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.FirewallRule = firewallRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *firewallFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state firewall.FirewallRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from firewall rule
	data, err := r.client.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		if strings.Contains(detail, "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading firewall rule",
			"Could not read firewall rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError("Error processing firewall rule", "Received nil response from API")
		return
	}
	firewallRuleModel := createFirewallRuleModel(data)

	firewallRuleObjectValue, objDiags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.FirewallRule = firewallRuleObjectValue
	// Map additional fields to state (input block & top‑level id)
	state.FirewallId = types.StringValue(data.FirewallRule.FirewallId)
	state.FirewallRuleCreate = firewall.FirewallRuleCreate{
		SourceAddress:      data.FirewallRule.SourceAddress,
		DestinationAddress: data.FirewallRule.DestinationAddress,
		Service:            firewallRuleModel.Service,
		Action:             types.StringValue(string(data.FirewallRule.Action)),
		Direction:          types.StringValue(string(data.FirewallRule.Direction)),
		Status:             types.StringValue(string(data.FirewallRule.Status)),
		Description:        types.StringPointerValue(data.FirewallRule.Description.Get()),
		OrderRuleId:        state.FirewallRuleCreate.OrderRuleId,
		OrderDirection:     state.FirewallRuleCreate.OrderDirection,
	}

	// Set refreshed state
	objDiags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *firewallFirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state firewall.FirewallRuleResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateFirewallRule(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating direct connect",
			"Could not update direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Wait for the firewall rule to reach ACTIVE state after update
	if err := waitForFirewallRuleStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"ACTIVE"}); err != nil {
		resp.Diagnostics.AddError(
			"Error updating firewall rule",
			"Could not wait for firewall rule to become active: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetFirewallRule as UpdateFirewallRule items are not populated.
	data, err := r.client.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		if strings.Contains(detail, "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading firewall rule",
			"Could not read firewall rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError("Error processing firewall rule", "Received nil response from API")
		return
	}
	firewallRuleModel := createFirewallRuleModel(data)

	firewallRuleObjectValue, objDiags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.FirewallRule = firewallRuleObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state firewall.FirewallRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing firewall rule
	err := r.client.DeleteFirewallRule(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting firewall rule",
			"Could not delete firewall rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *firewallFirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForFirewallRuleStatus(ctx context.Context, firewallClient *firewall.Client, id string, pendingStates []string, targetStates []string) error {

	err := client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := firewallClient.GetFirewallRule(id)
		if err != nil {

			return nil, "", err
		}
		currentState := string(info.FirewallRule.State)

		return info, currentState, nil
	}, -1, -1, -1, -1)
	return err
}

func createFirewallRuleModel(data *scpfirewall.FirewallRuleShowResponse) firewall.FirewallRule {
	fwRule := data.FirewallRule
	sourceAddresses := make([]string, 0, len(fwRule.SourceAddress))
	for _, address := range fwRule.SourceAddress {
		sourceAddresses = append(sourceAddresses, address)
	}
	destinationAddresses := make([]string, 0, len(fwRule.DestinationAddress))
	for _, address := range fwRule.DestinationAddress {
		destinationAddresses = append(destinationAddresses, address)
	}
	services := make([]firewall.FirewallPort, 0, len(fwRule.Service))
	for _, service := range fwRule.Service {
		services = append(services, firewall.FirewallPort{
			ServiceType:  types.StringValue(string(service.ServiceType)),
			ServiceValue: types.StringPointerValue(service.ServiceValue),
		})
	}

	return firewall.FirewallRule{
		Id:                   types.StringValue(fwRule.Id),
		Name:                 types.StringPointerValue(fwRule.Name.Get()),
		FirewallId:           types.StringValue(fwRule.FirewallId),
		Sequence:             types.Int32Value(fwRule.Sequence),
		SourceInterface:      types.StringValue(fwRule.SourceInterface),
		SourceAddress:        sourceAddresses,
		DestinationInterface: types.StringValue(fwRule.DestinationInterface),
		DestinationAddress:   destinationAddresses,
		Service:              services,
		Action:               types.StringValue(string(fwRule.Action)),
		Direction:            types.StringValue(string(fwRule.Direction)),
		VendorRuleId:         types.StringValue(fwRule.VendorRuleId),
		Description:          types.StringPointerValue(fwRule.Description.Get()),
		State:                types.StringValue(string(fwRule.State)),
		Status:               types.StringValue(string(fwRule.Status)),
		CreatedAt:            types.StringValue(fwRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(fwRule.CreatedBy),
		ModifiedAt:           types.StringValue(fwRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(fwRule.ModifiedBy),
	}
}
