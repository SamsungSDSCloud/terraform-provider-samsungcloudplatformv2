package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/dns/1.1"
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
		hostedZoneState := convertHostedZone(hostedZone)

		state.HostedZones = append(state.HostedZones, hostedZoneState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func convertHostedZone(hostedZone scpdns.HostedZone) dns.HostedZone {

	var attributes *dns.Attributes
	if hostedZone.Attributes != nil {
		serviceTier, ok := hostedZone.Attributes["service_tier"].(string)
		if ok {
			attributes = &dns.Attributes{
				ServiceTier: types.StringValue(serviceTier),
			}
		}
	}

	var links *dns.Links
	if hostedZone.Links != nil {
		self, ok := hostedZone.Links["self"].(string)
		if ok {
			links = &dns.Links{
				Self: types.StringValue(self),
			}
		}
	}

	// Masters 슬라이스 변환
	masters := make([]types.String, len(hostedZone.Masters))
	for i, s := range hostedZone.Masters {
		masters[i] = types.StringValue(s)
	}

	return dns.HostedZone{
		Action:         types.StringValue(hostedZone.Action),
		Attributes:     attributes,
		CreatedAt:      virtualserverutil.ToNullableStringValue(hostedZone.CreatedAt.Get()),
		Description:    virtualserverutil.ToNullableStringValue(hostedZone.Description.Get()),
		Email:          types.StringValue(hostedZone.Email),
		HostedZoneType: virtualserverutil.ToNullableStringValue(hostedZone.HostedZoneType.Get()),
		Id:             types.StringValue(hostedZone.Id),
		Links:          links,
		Masters:        masters,
		Name:           types.StringValue(hostedZone.Name),
		PoolId:         types.StringValue(hostedZone.PoolId),
		ProjectId:      types.StringValue(hostedZone.ProjectId),
		Serial:         types.Int32Value(hostedZone.Serial),
		Shared:         common.ToNullableBoolValue(hostedZone.Shared.Get()),
		Status:         types.StringValue(hostedZone.Status),
		TransferredAt:  virtualserverutil.ToNullableStringValue(hostedZone.TransferredAt.Get()),
		Ttl:            common.ToNullableInt32Value(hostedZone.Ttl.Get()),
		Type:           virtualserverutil.ToNullableStringValue(hostedZone.Type.Get()),
		UpdatedAt:      virtualserverutil.ToNullableStringValue(hostedZone.UpdatedAt.Get()),
		Version:        common.ToNullableInt32Value(hostedZone.Version.Get()),
	}
}

func convertHostedZoneShowResponseToHostedZone(a scpdns.HostedZoneShowResponse) scpdns.HostedZone {
	b := scpdns.HostedZone{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}

func convertHostedZoneCreateResponseToHostedZone(a scpdns.HostedZoneCreateResponse) scpdns.HostedZone {
	b := scpdns.HostedZone{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}

func convertHostedZoneSetResponseToHostedZone(a scpdns.HostedZoneSetResponse) scpdns.HostedZone {
	b := scpdns.HostedZone{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}

func convertHostedZoneDeleteResponseToHostedZone(a scpdns.HostedZoneDeleteResponse) scpdns.HostedZone {
	b := scpdns.HostedZone{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &b)
	return b
}
