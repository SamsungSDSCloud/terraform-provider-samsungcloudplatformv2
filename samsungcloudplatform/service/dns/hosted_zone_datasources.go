package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/dns/1.2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsHostedZoneDataSources{}
	_ datasource.DataSourceWithConfigure = &dnsHostedZoneDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewDnsHostedZoneDataSources() datasource.DataSource {
	return &dnsHostedZoneDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type dnsHostedZoneDataSources struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *dnsHostedZoneDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_hosted_zones" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *dnsHostedZoneDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of hosted zone.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("ExactName"): schema.StringAttribute{
				Description: "ExactName",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Type",
				Optional:    true,
			},
			common.ToSnakeCase("Email"): schema.StringAttribute{
				Description: "Email",
				Optional:    true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Status",
				Optional:    true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description",
				Optional:    true,
			},
			common.ToSnakeCase("Ttl"): schema.Int32Attribute{
				Description: "Ttl",
				Optional:    true,
			},
			common.ToSnakeCase("HostedZones"): schema.ListNestedAttribute{
				Description: "A list of HostedZone.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Action"): schema.StringAttribute{
							Description: "Action",
							Optional:    true,
						},
						common.ToSnakeCase("Attributes"): schema.SingleNestedAttribute{
							Description: "Attributes",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceTier"): schema.StringAttribute{
									Description: "ServiceTier",
									Optional:    true,
								},
							},
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "CreatedAt",
							Optional:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Optional:    true,
						},
						common.ToSnakeCase("Email"): schema.StringAttribute{
							Description: "Email",
							Optional:    true,
						},
						common.ToSnakeCase("HostedZoneType"): schema.StringAttribute{
							Description: "HostedZoneType",
							Optional:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
						common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
							Description: "Links",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Self"): schema.StringAttribute{
									Description: "Self",
									Optional:    true,
								},
							},
						},
						common.ToSnakeCase("Masters"): schema.ListAttribute{
							Description: "Masters",
							Optional:    true,
							ElementType: types.StringType,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Optional:    true,
						},
						common.ToSnakeCase("PoolId"): schema.StringAttribute{
							Description: "PoolId",
							Optional:    true,
						},
						common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
							Description: "PrivateDnsId",
							Optional:    true,
						},
						common.ToSnakeCase("PrivateDnsName"): schema.StringAttribute{
							Description: "PrivateDnsName",
							Optional:    true,
						},
						common.ToSnakeCase("ProjectId"): schema.StringAttribute{
							Description: "ProjectId",
							Optional:    true,
						},
						common.ToSnakeCase("Serial"): schema.Int32Attribute{
							Description: "Serial",
							Optional:    true,
						},
						common.ToSnakeCase("Shared"): schema.BoolAttribute{
							Description: "Shared",
							Optional:    true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "Status",
							Optional:    true,
						},
						common.ToSnakeCase("TransferredAt"): schema.StringAttribute{
							Description: "TransferredAt",
							Optional:    true,
						},
						common.ToSnakeCase("Ttl"): schema.Int32Attribute{
							Description: "Ttl",
							Optional:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Type",
							Optional:    true,
						},
						common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
							Description: "UpdatedAt",
							Optional:    true,
						},
						common.ToSnakeCase("Version"): schema.Int32Attribute{
							Description: "Version",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *dnsHostedZoneDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Dns
	d.clients = inst.Client
}

func (d *dnsHostedZoneDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state dns.HostedZoneDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetHostedZoneList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HostedZones",
			err.Error(),
		)
		return
	}

	for _, hostedZone := range data.HostedZones {
		hostedZoneState := convertHostedZoneShowResponseV1Dot2ToHostedZone(convertHostedZoneV1Dot2ToHostedZoneShowResponseV1Dot2(hostedZone))

		state.HostedZones = append(state.HostedZones, hostedZoneState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func convertHostedZoneV1Dot2ToHostedZoneShowResponseV1Dot2(hostedZoneV1Dot2 scpdns.HostedZoneV1Dot2) scpdns.HostedZoneShowResponseV1Dot2 {
	hostedZoneShowResponseV1Dot2 := scpdns.HostedZoneShowResponseV1Dot2{}
	data, _ := json.Marshal(hostedZoneV1Dot2)
	json.Unmarshal(data, &hostedZoneShowResponseV1Dot2)
	return hostedZoneShowResponseV1Dot2
}

func convertHostedZoneShowResponseV1Dot2ToHostedZone(hostedZoneShowResponseV1Dot2 scpdns.HostedZoneShowResponseV1Dot2) dns.HostedZone {

	var attributes *dns.Attributes
	if hostedZoneShowResponseV1Dot2.Attributes != nil {
		serviceTier, ok := hostedZoneShowResponseV1Dot2.Attributes["service_tier"].(string)
		if ok {
			attributes = &dns.Attributes{
				ServiceTier: types.StringValue(serviceTier),
			}
		}
	}

	var links *dns.Links
	if hostedZoneShowResponseV1Dot2.Links != nil {
		self, ok := hostedZoneShowResponseV1Dot2.Links["self"].(string)
		if ok {
			links = &dns.Links{
				Self: types.StringValue(self),
			}
		}
	}

	// Masters 슬라이스 변환
	masters := make([]types.String, len(hostedZoneShowResponseV1Dot2.Masters))
	for i, s := range hostedZoneShowResponseV1Dot2.Masters {
		masters[i] = types.StringValue(s)
	}

	return dns.HostedZone{
		Action:         types.StringValue(hostedZoneShowResponseV1Dot2.Action),
		Attributes:     attributes,
		CreatedAt:      virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.CreatedAt.Get()),
		Description:    virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.Description.Get()),
		Email:          types.StringValue(hostedZoneShowResponseV1Dot2.Email),
		HostedZoneType: virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.HostedZoneType.Get()),
		Id:             types.StringValue(hostedZoneShowResponseV1Dot2.Id),
		Links:          links,
		Masters:        masters,
		Name:           types.StringValue(hostedZoneShowResponseV1Dot2.Name),
		PoolId:         types.StringValue(hostedZoneShowResponseV1Dot2.PoolId),
		PrivateDnsId:   virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.PrivateDnsId.Get()),
		PrivateDnsName: virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.PrivateDnsName.Get()),
		ProjectId:      types.StringValue(hostedZoneShowResponseV1Dot2.ProjectId),
		Serial:         types.Int32Value(hostedZoneShowResponseV1Dot2.Serial),
		Shared:         common.ToNullableBoolValue(hostedZoneShowResponseV1Dot2.Shared.Get()),
		Status:         types.StringValue(hostedZoneShowResponseV1Dot2.Status),
		TransferredAt:  virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.TransferredAt.Get()),
		Ttl:            common.ToNullableInt32Value(hostedZoneShowResponseV1Dot2.Ttl.Get()),
		Type:           virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.Type.Get()),
		UpdatedAt:      virtualserverutil.ToNullableStringValue(hostedZoneShowResponseV1Dot2.UpdatedAt.Get()),
		Version:        common.ToNullableInt32Value(hostedZoneShowResponseV1Dot2.Version.Get()),
	}
}

func convertHostedZoneDeleteResponseToHostedZone(hostedZoneDeleteResponse scpdns.HostedZoneDeleteResponse) dns.HostedZone {
	var attributes *dns.Attributes
	if hostedZoneDeleteResponse.Attributes != nil {
		serviceTier, ok := hostedZoneDeleteResponse.Attributes["service_tier"].(string)
		if ok {
			attributes = &dns.Attributes{
				ServiceTier: types.StringValue(serviceTier),
			}
		}
	}

	var links *dns.Links
	if hostedZoneDeleteResponse.Links != nil {
		self, ok := hostedZoneDeleteResponse.Links["self"].(string)
		if ok {
			links = &dns.Links{
				Self: types.StringValue(self),
			}
		}
	}

	// Masters 슬라이스 변환
	masters := make([]types.String, len(hostedZoneDeleteResponse.Masters))
	for i, s := range hostedZoneDeleteResponse.Masters {
		masters[i] = types.StringValue(s)
	}

	return dns.HostedZone{
		Action:         types.StringValue(hostedZoneDeleteResponse.Action),
		Attributes:     attributes,
		CreatedAt:      virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.CreatedAt.Get()),
		Description:    virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.Description.Get()),
		Email:          types.StringValue(hostedZoneDeleteResponse.Email),
		HostedZoneType: virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.HostedZoneType.Get()),
		Id:             types.StringValue(hostedZoneDeleteResponse.Id),
		Links:          links,
		Masters:        masters,
		Name:           types.StringValue(hostedZoneDeleteResponse.Name),
		PoolId:         types.StringValue(hostedZoneDeleteResponse.PoolId),
		ProjectId:      types.StringValue(hostedZoneDeleteResponse.ProjectId),
		Serial:         types.Int32Value(hostedZoneDeleteResponse.Serial),
		Shared:         common.ToNullableBoolValue(hostedZoneDeleteResponse.Shared.Get()),
		Status:         types.StringValue(hostedZoneDeleteResponse.Status),
		TransferredAt:  virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.TransferredAt.Get()),
		Ttl:            common.ToNullableInt32Value(hostedZoneDeleteResponse.Ttl.Get()),
		Type:           virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.Type.Get()),
		UpdatedAt:      virtualserverutil.ToNullableStringValue(hostedZoneDeleteResponse.UpdatedAt.Get()),
		Version:        common.ToNullableInt32Value(hostedZoneDeleteResponse.Version.Get()),
	}
}
