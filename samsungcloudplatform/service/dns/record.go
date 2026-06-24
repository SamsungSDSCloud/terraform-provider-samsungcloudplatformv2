package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/dns"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
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
	_ resource.Resource                = &dnsRecordResource{}
	_ resource.ResourceWithConfigure   = &dnsRecordResource{}
	_ resource.ResourceWithImportState = &dnsRecordResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsRecordResource() resource.Resource {
	return &dnsRecordResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsRecordResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "A DNS record resource for managing domain name resolution.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the DNS record.\n" +
					"  - example : 6ed7bc1-4b05-3cc7-7105-c1b71f7f30a7 ",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("HostedZoneId"): schema.StringAttribute{
				Description: "The identifier of the hosted zone that contains this DNS record.\n" +
					"  - example : 3432012nfdksdf03ktrld9234lgfg ",
				Optional: true,
			},
			common.ToSnakeCase("Record"): schema.SingleNestedAttribute{
				Description: "Detailed information about the DNS record.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "The action performed on the DNS record.\n" +
							"  - example : NONE",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the DNS record.\n" +
							"  - example : 6ed7bc1-4b05-3cc7-7105-c1b71f7f30a7 ",
						Optional: true,
					},
					common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
						Description: "The links related to the DNS record.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Self"): schema.StringAttribute{
								Description: "The self-referential link of the DNS record.\n" +
									"  - example : https://api.samsungsdscloud.com/dns/v1/records/3432012nfdksdf03ktrld9234lgfg ",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the DNS record.\n" +
							"  - example : test.app ",
						Optional: true,
					},
					common.ToSnakeCase("ProjectId"): schema.StringAttribute{
						Description: "The project identifier associated with the DNS record.\n" +
							"  - example : 003dffc50eb123a1cbf4f2e5c71d4f15 ",
						Optional: true,
					},
					common.ToSnakeCase("Records"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "A list of data for this record\n" +
							"  - example : [\"12.34.45.67\"] ",
						Optional: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the DNS record.\n" +
							"  - example : ACTIVE ",
						Optional: true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "The Time-To-Live (TTL) value in seconds for the DNS record.\n" +
							"  - example : 3600 ",
						Optional: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of the DNS record (e.g., A, AAAA, CNAME, MX, TXT, SPF).\n" +
							"  - example : A ",
						Optional: true,
					},
					common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last updated, in ISO 8601 format.\n" +
							"  - example : 2026-02-09T08:00:40Z ",
						Optional: true,
					},
					common.ToSnakeCase("Version"): schema.Int32Attribute{
						Description: "The version of the DNS record.\n" +
							"  - example : 1 ",
						Optional: true,
					},
					common.ToSnakeCase("ZoneId"): schema.StringAttribute{
						Description: "ID for the zone that contains this record\n" +
							"  - example : 3432012nfdksdf03ktrld9234lgfg ",
						Optional: true,
					},
					common.ToSnakeCase("ZoneName"): schema.StringAttribute{
						Description: "The name of the zone that contains this record\n" +
							"  - example : my-zone.com ",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("RecordCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new DNS record.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name for the DNS record to be created.\n" +
							"  - example : test.app ",
						Optional: true,
					},
					common.ToSnakeCase("Records"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "A list of data for this record\n" +
							"  - example : [\"12.34.45.67\"]",
						Optional: true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "The Time-To-Live (TTL) value in seconds for the DNS record.\n" +
							"  - example : 3600 ",
						Optional: true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "The type of the DNS record to create (e.g., A, AAAA, CNAME, MX, TXT).\n" +
							"  - example : A ",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("A", "AAAA", "CNAME", "MX", "TXT", "SPF"),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.RecordResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateRecord(ctx, plan.HostedZoneId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Record",
			"Could not create Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	idPtr := data.Id.Get()
	if idPtr == nil {
		resp.Diagnostics.AddError("Error creating Record", "API returned record without id")
		return
	}

	createErr := waitForRecordStatus(ctx, r.client, plan.HostedZoneId.ValueString(), *idPtr, []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating record",
			"Error creating for record to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetRecord(ctx, plan.HostedZoneId.ValueString(), *idPtr)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(*idPtr)

	recordModel := convertRecordDetail(*dataForShow)

	recordOjbectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Record = recordOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: hostedZoneId/recordId, got: %q", req.ID),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("hosted_zone_id"), types.StringValue(parts[0]))
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.RecordResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	recordModel := convertRecordDetail(*data)

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Record = recordObjectValue
	state.HostedZoneId = recordModel.ZoneId

	if state.RecordCreate == nil {
		state.RecordCreate = &dns.RecordCreate{}
	}
	state.RecordCreate.Description = recordModel.Description
	state.RecordCreate.Name = recordModel.Name
	state.RecordCreate.Records = recordModel.Records
	state.RecordCreate.Ttl = recordModel.Ttl
	state.RecordCreate.Type = recordModel.Type

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan

	var state dns.RecordResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Record",
			"Could not update Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForRecordStatus(ctx, r.client, state.HostedZoneId.ValueString(), state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating record",
			"Error updating for record to become active: "+updateErr.Error(),
		)
		return
	}

	idPtr := data.Id.Get()
	if idPtr == nil {
		resp.Diagnostics.AddError("Error updating Record", "API returned record without id")
		return
	}

	dataForShow, err := r.client.GetRecord(ctx, state.HostedZoneId.ValueString(), *idPtr)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	recordModel := convertRecordDetail(*dataForShow)

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Record = recordObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.RecordResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.DeleteRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Record",
			"Could not delete Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	recordModel := convertRecordDetail(convertRecordCreateResponseToRecord(*data))

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Record = recordObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func waitForRecordStatus(ctx context.Context, recordClient *dns.Client, hostedZoneId string, recordId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := recordClient.GetRecord(ctx, hostedZoneId, recordId)
		if err != nil {
			return nil, "", err
		}
		return info, *info.Status.Get(), nil
	}, -1, -1, -1, -1)
}
