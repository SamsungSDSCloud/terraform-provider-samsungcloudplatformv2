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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnsPublicDomainNameResource{}
	_ resource.ResourceWithConfigure   = &dnsPublicDomainNameResource{}
	_ resource.ResourceWithImportState = &dnsPublicDomainNameResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsPublicDomainNameResource() resource.Resource {
	return &dnsPublicDomainNameResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsPublicDomainNameResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsPublicDomainNameResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_public_domain_name" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsPublicDomainNameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "A public domain name registration resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the public domain name.\n" +
					"  - example : 125jkdkt5fpublicdomain3193rud546 ",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("PublicDomainName"): schema.SingleNestedAttribute{
				Description: "Detailed information about the public domain name.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "The type of address for the domain registration.\n" +
							"  - example : DOMESTIC ",
						Optional: true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "Indicates whether automatic extension is enabled for the domain.\n" +
							"  - example : true ",
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
					common.ToSnakeCase("DomesticAddressEn"): schema.StringAttribute{
						Description: "Domestic address in English\n" +
							"  - example : Samsung-ro 123, Suwon-si, Gyeonggi-do, Korea ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticAddressKo"): schema.StringAttribute{
						Description: "Domestic address in Korean\n" +
							"  - example : 경기도 수원시 삼성로 123 ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "Domestic first address in English\n" +
							"  - example : Samsung-ro 123 ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "Domestic first address in Korean\n" +
							"  - example : 삼성로 123 ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "Domestic second address in English\n" +
							"  - example : Suwon-si, Gyeonggi-do ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "Domestic second address in Korean\n" +
							"  - example : 경기도 수원시 ",
						Optional: true,
					},
					common.ToSnakeCase("ExpiredDate"): schema.StringAttribute{
						Description: "The expiration date of the domain registration.\n" +
							"  - example : 2025-12-31 ",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the public domain name.\n" +
							"  - example : 10fjkeweffpublicdomain3193rud543 ",
						Computed: true,
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
						Description: "The name of the public domain name.\n" +
							"  - example : example.com ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasAddress"): schema.StringAttribute{
						Description: "The overseas address for the domain registration.\n" +
							"  - example : 123 Main Street, City, Country ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "Overseas first address for the domain registration.\n" +
							"  - example : 123 Main Street ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "Overseas second address for the domain registration.\n" +
							"  - example : Suite 100 ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "Overseas third address for the domain registration.\n" +
							"  - example : City, State 12345 ",
						Optional: true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "The postal code for the domain registration.\n" +
							"  - example : 12345 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "The email address of the domain registrant.\n" +
							"  - example : user@example.com ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "The name of the domain registrant in English\n" +
							"  - example : John Doe ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "The name of the domain registrant in Korean\n" +
							"  - example : 홍길동 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "The telephone number of the domain registrant.\n" +
							"  - example : 82-10-1234-5678 ",
						Optional: true,
					},
					common.ToSnakeCase("StartDate"): schema.StringAttribute{
						Description: "The start date of the domain registration.\n" +
							"  - example : 2024-01-01 ",
						Optional: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the public domain name.\n" +
							"  - example : REGISTERED ",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("PublicDomainNameCreate"): schema.SingleNestedAttribute{
				Description: "Configuration for creating a new public domain name registration.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "The type of address for the domain registration.\n" +
							"  - example : DOMESTIC ",
						Optional: true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "Indicates whether automatic extension is enabled for the domain.\n" +
							"  - example : true ",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : This is description ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "Domestic first address in English\n" +
							"  - example : Samsung-ro 123 ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "Domestic first address in Korean\n" +
							"  - example : 삼성로 123 ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "Domestic second address in English\n" +
							"  - example : Suwon-si, Gyeonggi-do ",
						Optional: true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "Domestic second address in Korean\n" +
							"  - example : 경기도 수원시 ",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name for the public domain name to be created.\n" +
							"  - example : example.com ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "Overseas first address for the domain registration.\n" +
							"  - example : 123 Main Street ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "Overseas second address for the domain registration.\n" +
							"  - example : Suite 100 ",
						Optional: true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "Overseas third address for the domain registration.\n" +
							"  - example : City, State 12345 ",
						Optional: true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "The postal code for the domain registration.\n" +
							"  - example : 12345 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "The email address of the domain registrant.\n" +
							"  - example : user@example.com ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "The name of the domain registrant in English\n" +
							"  - example : John Doe ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "The name of the domain registrant in Korean\n" +
							"  - example : 홍길동 ",
						Optional: true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "The telephone number of the domain registrant.\n" +
							"  - example : 82-10-1234-5678 ",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsPublicDomainNameResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dnsPublicDomainNameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.PublicDomainNameResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreatePublicDomainName(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Public Domain Name",
			"Could not create Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	createErr := waitForPublicDomainNameStatus(ctx, r.client, data.PublicDomainName.Id, []string{}, []string{"REGISTERED"})

	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating(activating) Public Domain Name",
			"Error creating(activating) for Public Domain Name to become active: "+createErr.Error())
		return
	}

	plan.Id = types.StringValue(data.PublicDomainName.Id)

	publicDomainNameModel := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainNameOjbectValue, diags := types.ObjectValueFrom(ctx, publicDomainNameModel.AttributeTypes(), publicDomainNameModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.PublicDomainName = publicDomainNameOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements [resource.ResourceWithImportState].
func (r *dnsPublicDomainNameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsPublicDomainNameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.PublicDomainNameResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetPublicDomainName(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Public Domain Name",
			"Could not read Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicDomainNameModel := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainNameObjectValue, diags := types.ObjectValueFrom(ctx, publicDomainNameModel.AttributeTypes(), publicDomainNameModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.PublicDomainName = publicDomainNameObjectValue

	if state.PublicDomainNameCreate == nil {
		state.PublicDomainNameCreate = &dns.PublicDomainNameCreate{}
	}
	state.PublicDomainNameCreate.AddressType = publicDomainNameModel.AddressType
	state.PublicDomainNameCreate.AutoExtension = publicDomainNameModel.AutoExtension
	state.PublicDomainNameCreate.Description = publicDomainNameModel.Description
	state.PublicDomainNameCreate.DomesticFirstAddressEn = publicDomainNameModel.DomesticFirstAddressEn
	state.PublicDomainNameCreate.DomesticFirstAddressKo = publicDomainNameModel.DomesticFirstAddressKo
	state.PublicDomainNameCreate.DomesticSecondAddressEn = publicDomainNameModel.DomesticSecondAddressEn
	state.PublicDomainNameCreate.DomesticSecondAddressKo = publicDomainNameModel.DomesticSecondAddressKo
	state.PublicDomainNameCreate.Name = publicDomainNameModel.Name
	state.PublicDomainNameCreate.OverseasFirstAddress = publicDomainNameModel.OverseasFirstAddress
	state.PublicDomainNameCreate.OverseasSecondAddress = publicDomainNameModel.OverseasSecondAddress
	state.PublicDomainNameCreate.OverseasThirdAddress = publicDomainNameModel.OverseasThirdAddress
	state.PublicDomainNameCreate.PostalCode = publicDomainNameModel.PostalCode
	state.PublicDomainNameCreate.RegisterEmail = publicDomainNameModel.RegisterEmail
	state.PublicDomainNameCreate.RegisterNameEn = publicDomainNameModel.RegisterNameEn
	state.PublicDomainNameCreate.RegisterNameKo = publicDomainNameModel.RegisterNameKo
	state.PublicDomainNameCreate.RegisterTelno = publicDomainNameModel.RegisterTelno

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsPublicDomainNameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var oldState dns.PublicDomainNameResource
	req.State.Get(ctx, &oldState)
	var state dns.PublicDomainNameResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if checkPublicDomainNameImmutableFields(oldState, state) {
		resp.Diagnostics.AddError(
			"Error updating Public Domain Name",
			"Public domain name fields (name) cannot be modified.",
		)
		return
	}

	oldCreate := oldState.PublicDomainNameCreate
	newCreate := state.PublicDomainNameCreate

	updateDomainChanged := newCreate.AutoExtension != oldCreate.AutoExtension ||
		newCreate.Description != oldCreate.Description

	if updateDomainChanged {
		_, err := r.client.UpdatePublicDomainName(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Public Domain Name",
				"Could not update Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	whoisChanged := newCreate.AddressType != oldCreate.AddressType ||
		newCreate.DomesticFirstAddressEn != oldCreate.DomesticFirstAddressEn ||
		newCreate.DomesticFirstAddressKo != oldCreate.DomesticFirstAddressKo ||
		newCreate.DomesticSecondAddressEn != oldCreate.DomesticSecondAddressEn ||
		newCreate.DomesticSecondAddressKo != oldCreate.DomesticSecondAddressKo ||
		newCreate.OverseasFirstAddress != oldCreate.OverseasFirstAddress ||
		newCreate.OverseasSecondAddress != oldCreate.OverseasSecondAddress ||
		newCreate.OverseasThirdAddress != oldCreate.OverseasThirdAddress ||
		newCreate.PostalCode != oldCreate.PostalCode ||
		newCreate.RegisterEmail != oldCreate.RegisterEmail ||
		newCreate.RegisterTelno != oldCreate.RegisterTelno

	if whoisChanged {
		_, err := r.client.UpdatePublicDomainNameInfomation(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Public Domain Name",
				"Could not update Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	data, err := r.client.GetPublicDomainName(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Public Domain Name",
			"Could not read Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicDomainNameModel := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainNameObjectValue, diags := types.ObjectValueFrom(ctx, publicDomainNameModel.AttributeTypes(), publicDomainNameModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.PublicDomainName = publicDomainNameObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsPublicDomainNameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	resp.Diagnostics.AddError(
		"Error Deleting Public Domain Name",
		"Public Domain Name does not support delete method.",
	)
	return
}

func waitForPublicDomainNameStatus(ctx context.Context, privateDnsClient *dns.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := privateDnsClient.GetPublicDomainName(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.PublicDomainName.Status), nil
	}, -1, -1, -1, -1)
}

func checkPublicDomainNameImmutableFields(oldState dns.PublicDomainNameResource, newState dns.PublicDomainNameResource) bool {
	oldCreate := oldState.PublicDomainNameCreate
	newCreate := newState.PublicDomainNameCreate

	if oldCreate.Name != newCreate.Name {
		return true
	}
	return false
}
