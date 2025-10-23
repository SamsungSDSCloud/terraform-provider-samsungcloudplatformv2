package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/dns"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dnsPublicDomainNameResource{}
	_ resource.ResourceWithConfigure = &dnsPublicDomainNameResource{}
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
		Description: "PublicDomainName.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("PublicDomainName"): schema.SingleNestedAttribute{
				Description: "A detail of PublicDomainName.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "AddressType",
						Optional:    true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "AutoExtension",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticAddressEn"): schema.StringAttribute{
						Description: "DomesticAddressEn",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticAddressKo"): schema.StringAttribute{
						Description: "DomesticAddressKo",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "DomesticFirstAddressEn",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "DomesticFirstAddressKo",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "DomesticSecondAddressEn",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "DomesticSecondAddressKo",
						Optional:    true,
					},
					common.ToSnakeCase("ExpiredDate"): schema.StringAttribute{
						Description: "ExpiredDate",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasAddress"): schema.StringAttribute{
						Description: "OverseasAddress",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "OverseasFirstAddress",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "OverseasSecondAddress",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "OverseasThirdAddress",
						Optional:    true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "PostalCode",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "RegisterEmail",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "RegisterNameEn",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "RegisterNameKo",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "RegisterTelno",
						Optional:    true,
					},
					common.ToSnakeCase("StartDate"): schema.StringAttribute{
						Description: "StartDate",
						Optional:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("PublicDomainNameCreate"): schema.SingleNestedAttribute{
				Description: "Create PublicDomainName.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AddressType"): schema.StringAttribute{
						Description: "AddressType",
						Optional:    true,
					},
					common.ToSnakeCase("AutoExtension"): schema.BoolAttribute{
						Description: "AutoExtension",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressEn"): schema.StringAttribute{
						Description: "DomesticFirstAddressEn",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticFirstAddressKo"): schema.StringAttribute{
						Description: "DomesticFirstAddressKo",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressEn"): schema.StringAttribute{
						Description: "DomesticSecondAddressEn",
						Optional:    true,
					},
					common.ToSnakeCase("DomesticSecondAddressKo"): schema.StringAttribute{
						Description: "DomesticSecondAddressKo",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasFirstAddress"): schema.StringAttribute{
						Description: "OverseasFirstAddress",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasSecondAddress"): schema.StringAttribute{
						Description: "OverseasSecondAddress",
						Optional:    true,
					},
					common.ToSnakeCase("OverseasThirdAddress"): schema.StringAttribute{
						Description: "OverseasThirdAddress",
						Optional:    true,
					},
					common.ToSnakeCase("PostalCode"): schema.StringAttribute{
						Description: "PostalCode",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterEmail"): schema.StringAttribute{
						Description: "RegisterEmail",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterNameEn"): schema.StringAttribute{
						Description: "RegisterNameEn",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterNameKo"): schema.StringAttribute{
						Description: "RegisterNameKo",
						Optional:    true,
					},
					common.ToSnakeCase("RegisterTelno"): schema.StringAttribute{
						Description: "RegisterTelno",
						Optional:    true,
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

	plan.Id = types.StringValue(data.PublicDomainName.Id)

	publicDomainNameModel := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainNameOjbectValue, diags := types.ObjectValueFrom(ctx, publicDomainNameModel.AttributeTypes(), publicDomainNameModel)
	plan.PublicDomainName = publicDomainNameOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Public Domain Name",
			"Could not read Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicDomainNameModel := convertPublicDomainDetail(data.PublicDomainName)

	publicDomainNameObjectValue, diags := types.ObjectValueFrom(ctx, publicDomainNameModel.AttributeTypes(), publicDomainNameModel)
	state.PublicDomainName = publicDomainNameObjectValue

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

	var state dns.PublicDomainNameResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePublicDomainName(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Public Domain Name",
			"Could not update Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	_, err = r.client.UpdatePublicDomainNameInfomation(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Public Domain Name",
			"Could not update Public Domain Name, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
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
