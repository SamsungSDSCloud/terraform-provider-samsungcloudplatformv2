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
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/dns/1.3"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnsPrivateDnsResource{}
	_ resource.ResourceWithConfigure   = &dnsPrivateDnsResource{}
	_ resource.ResourceWithImportState = &dnsPrivateDnsResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsPrivateDnsResource() resource.Resource {
	return &dnsPrivateDnsResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsPrivateDnsResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsPrivateDnsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_private_dns" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsPrivateDnsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "A private DNS instance for managing internal DNS resolution.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the private DNS.\n" +
					"  - example : 10fjkewefprivatedns3193rud543 ",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("PrivateDns"): schema.SingleNestedAttribute{
				Description: "Detailed information about the private DNS.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AuthDnsName"): schema.StringAttribute{
						Description: "The authoritative DNS name of the private DNS.\n" +
							"  - example : auth.dns.example.com ",
						Optional: true,
					},
					common.ToSnakeCase("ConnectedVpcIds"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of VPC identifiers connected to this private DNS.Only VPCs that are connected to the DNS can query the domain information registered in it.\n" +
							"  - example : [\"10fjkewefvpc3193rud543\"] ",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
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
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the private DNS.\n" +
							"  - example : 10fjkewefprivatedns3193rud543 ",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the private DNS.\n" +
							"  - example : private-dns01 ",
						Optional: true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "The resource pool identifier associated with the private DNS.\n" +
							"  - example : 10fjksdpooliddfsi12389esfdslkdsr32 ",
						Optional: true,
					},
					common.ToSnakeCase("PoolName"): schema.StringAttribute{
						Description: "The name of the resource pool.\n" +
							"  - example : pool-01 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisteredRegion"): schema.StringAttribute{
						Description: "The region where the private DNS is registered.\n" +
							"  - example : KR-WEST1 ",
						Optional: true,
					},
					common.ToSnakeCase("ResolverIp"): schema.StringAttribute{
						Description: "The IP address of the DNS resolver.\n" +
							"  - example : 198.19.0.101 ",
						Optional: true,
					},
					common.ToSnakeCase("ResolverName"): schema.StringAttribute{
						Description: "The name of the DNS resolver.\n" +
							"  - example : resolver-01 ",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the private DNS.\n" +
							"  - example : ACTIVE ",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("PrivateDnsCreate"): schema.SingleNestedAttribute{
				Description: "Configuration for creating a new private DNS.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("ConnectedVpcIds"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of VPC identifiers connected to this private DNS.Only VPCs that are connected to the DNS can query the domain information registered in it.\n" +
							"  - example : [\"10fjkewefvpc3193rud543\"] ",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name for the private DNS to be created.\n" +
							"  - example : private-dns01 ",
						Required: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsPrivateDnsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dnsPrivateDnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.PrivateDnsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state dns.PrivateDnsDataSource
	listData, listErr := r.client.GetPrivateDnsList(ctx, state)

	if listErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Private Dns List",
			listErr.Error(),
		)
		return
	}

	var data *scpdns.PrivateDnsShowResponse
	var err error

	isActivated := false

	for _, privateDns := range listData.PrivateDns {
		if privateDns.State == "INACTIVE" && privateDns.Name == plan.PrivateDnsCreate.Name.ValueString() {
			// ACTIVATE
			data, err = r.client.ActivatePrivateDns(ctx, plan)
			isActivated = true
			break
		}
	}

	if !isActivated {
		// CREATE
		data, err = r.client.CreatePrivateDns(ctx, plan)
	}

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating(activating) Private Dns",
			"Could not create(activate) Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	id := data.PrivateDns.Id
	if id == "" {
		resp.Diagnostics.AddError("Error creating Private Dns", "API returned record without id")
		return
	}
	createErr := waitForPrivateDnsStatus(ctx, r.client, id, []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating(activating) private dns",
			"Error creating(activating) for private dns to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetPrivateDns(ctx, id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(dataForShow.PrivateDns.Id)

	privateDnsModel := convertPrivateDns(dataForShow.PrivateDns)

	privateDnsOjbectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.PrivateDns = privateDnsOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements [resource.ResourceWithImportState].
func (r *dnsPrivateDnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsPrivateDnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.PrivateDnsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetPrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateDnsModel := convertPrivateDns(data.PrivateDns)

	privateDnsObjectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.PrivateDns = privateDnsObjectValue

	if state.PrivateDnsCreate == nil {
		state.PrivateDnsCreate = &dns.PrivateDnsCreate{}
	}
	state.PrivateDnsCreate.Description = privateDnsModel.Description
	state.PrivateDnsCreate.Name = privateDnsModel.Name
	state.PrivateDnsCreate.ConnectedVpcIds = privateDnsModel.ConnectedVpcIds

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsPrivateDnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan

	var state dns.PrivateDnsResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdatePrivateDns(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Private Dns",
			"Could not update Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForPrivateDnsStatus(ctx, r.client, data.PrivateDns.Id, []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating private dns",
			"Error updating for private dns to become active: "+updateErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetPrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateDnsModel := convertPrivateDns(dataForShow.PrivateDns)

	privateDnsObjectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.PrivateDns = privateDnsObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsPrivateDnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.PrivateDnsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Gslb
	err := r.client.DeletePrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Private Dns",
			"Could not delete Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func waitForPrivateDnsStatus(ctx context.Context, privateDnsClient *dns.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := privateDnsClient.GetPrivateDns(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.PrivateDns.State), nil
	}, -1, -1, -1, -1)
}
