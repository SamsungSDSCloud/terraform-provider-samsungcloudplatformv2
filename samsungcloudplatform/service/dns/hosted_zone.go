package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/dns"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
	_ resource.Resource                = &dnsHostedZoneResource{}
	_ resource.ResourceWithConfigure   = &dnsHostedZoneResource{}
	_ resource.ResourceWithImportState = &dnsHostedZoneResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsHostedZoneResource() resource.Resource {
	return &dnsHostedZoneResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsHostedZoneResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsHostedZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_hosted_zone" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsHostedZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "A hosted zone that contains DNS records for managing domain name resolution.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the hosted zone to query.\n" +
					"  - example : 3432012nfdksdf03ktrld9234lgfg ",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Zone"): schema.SingleNestedAttribute{
				Description: "Detailed information about the hosted zone.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Optional: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("HostedZoneType"): schema.StringAttribute{
						Description: "The type of the hosted zone (e.g., public or private).\n" +
							"  - example : private ",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Hosted Zone ID\n" +
							"  - example : 3432012nfdksdf03ktrld9234lgfg ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The domain name that a DNS service manages. all DNS records for that domain and its sub‑domains are stored and served within this hosted zone.\n" +
							"  - example : my-zone.com ",
						Optional: true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "Designate Pool ID\n" +
							"  - example : 10fjksdpooliddfsi12389esfdslkdsr32 ",
						Optional: true,
					},
					common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
						Description: "The DNS server ID for registering a Hosted Zone.For a Public‑type Hosted Zone, display it as an empty value.\n" +
							"  - example : 10fjkewefprivatedns3193rud543 ",
						Optional: true,
					},
					common.ToSnakeCase("PrivateDnsName"): schema.StringAttribute{
						Description: "The DNS server name for registering a Hosted Zone.For a Public‑type Hosted Zone, display it as an empty value.\n" +
							"  - example : private-dns01 ",
						Optional: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the hosted zone (e.g., ACTIVE, CREATING, DELETING).\n" +
							"  - example : ACTIVE ",
						Optional: true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "TTL for the zone.\n" +
							"  - example : 3600 ",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("HostedZoneCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new hosted zone.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The domain name that a DNS service manages. all DNS records for that domain and its sub‑domains are stored and served within this hosted zone.\n" +
							"  - example : my-zone.com ",
						Required: true,
					},
					common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
						Description: "The DNS server ID for registering a Hosted Zone. Input this only when the Hosted Zone is of Private type.\n" +
							"  - example : 10fjkewefprivatedns3193rud543 ",
						Optional: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of the hosted zone (e.g., public or private).\n" +
							"  - example : private ",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("public", "private"),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsHostedZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Dns
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *dnsHostedZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.HostedZoneResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateHostedZone(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating HostedZone",
			"Could not create HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	createErr := waitForHostedZoneStatus(ctx, r.client, data.Id, []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating hosted zone",
			"Error creating for hosted zone to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetHostedZone(ctx, data.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Id)

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*dataForShow)

	hostedZoneOjbectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Zone = hostedZoneOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements [resource.ResourceWithImportState].
func (r *dnsHostedZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsHostedZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.HostedZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*data)

	hostedZoneObjectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Zone = hostedZoneObjectValue

	if state.HostedZoneCreate == nil {
		state.HostedZoneCreate = &dns.HostedZoneCreate{}
	}
	state.HostedZoneCreate.Description = hostedZoneModel.Description
	state.HostedZoneCreate.Name = hostedZoneModel.Name
	state.HostedZoneCreate.PrivateDnsId = hostedZoneModel.PrivateDnsId
	state.HostedZoneCreate.Type = hostedZoneModel.HostedZoneType

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsHostedZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var oldState dns.HostedZoneResource
	req.State.Get(ctx, &oldState)
	var state dns.HostedZoneResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if checkModifiedFieldsExcludingDescription(oldState, state) {
		resp.Diagnostics.AddError(
			"Error updating HostedZone",
			"Hosted zones cannot be modified except for the description field.",
		)
		return
	}

	data, err := r.client.UpdateHostedZone(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating HostedZone",
			"Could not update HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForHostedZoneStatus(ctx, r.client, data.Id, []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating hosted zone",
			"Error updating for hosted zone to become active: "+updateErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*dataForShow)

	hostedZoneObjectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Zone = hostedZoneObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsHostedZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.HostedZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting HostedZone",
			"Could not delete HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func waitForHostedZoneStatus(ctx context.Context, hostedZoneClient *dns.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := hostedZoneClient.GetHostedZone(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Status), nil
	}, -1, -1, -1, -1)
}

func checkModifiedFieldsExcludingDescription(oldState dns.HostedZoneResource, newState dns.HostedZoneResource) bool {
	oldHostedZone := oldState.HostedZoneCreate
	newHostedZone := newState.HostedZoneCreate

	if oldHostedZone.Type != newHostedZone.Type {
		return true
	}
	if oldHostedZone.Name != newHostedZone.Name {
		return true
	}
	if oldHostedZone.PrivateDnsId != newHostedZone.PrivateDnsId {
		return true
	}
	return false
}
