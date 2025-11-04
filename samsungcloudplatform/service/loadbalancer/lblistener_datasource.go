package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbListenerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbListenerDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbListenerDataSource() datasource.DataSource {
	return &loadbalancerLbListenerDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLbListenerDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbListenerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_listener" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbListenerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Show Lb Listener.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("LbListener"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Listener.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("InsertClientIp"): schema.BoolAttribute{
						Description: "InsertClientIp",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("Persistence"): schema.StringAttribute{
						Description: "Persistence",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "ServerGroupId",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupName"): schema.StringAttribute{
						Description: "ServerGroupName",
						Optional:    true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "ServicePort",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "ResponseTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "SessionDurationTime",
						Optional:    true,
					},
					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SslCertificate",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "ClientCertId",
								Optional:    true,
							},
							common.ToSnakeCase("ClientCertLevel"): schema.StringAttribute{
								Description: "ClientCertLevel",
								Optional:    true,
							},
							common.ToSnakeCase("ServerCertLevel"): schema.StringAttribute{
								Description: "ServerCertLevel",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("SniCertificate"): schema.ListNestedAttribute{
						Description: "SniCertificate",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("SniCertId"): schema.StringAttribute{
									Description: "SniCertId",
									Optional:    true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "DomainName",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("UrlHandler"): schema.ListNestedAttribute{
						Description: "UrlHandler",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("UrlPattern"): schema.StringAttribute{
									Description: "UrlPattern",
									Optional:    true,
								},
								common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
									Description: "ServerGroupId",
									Optional:    true,
								},
								common.ToSnakeCase("Seq"): schema.Int32Attribute{
									Description: "Seq",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("HttpsRedirection"): schema.SingleNestedAttribute{
						Description: "HttpsRedirection",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "Protocol",
								Optional:    true,
							},
							common.ToSnakeCase("Port"): schema.StringAttribute{
								Description: "Port",
								Optional:    true,
							},
							common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
								Description: "ResponseCode",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "UrlRedirection",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "XForwardedFor",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "XForwardedPort",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "XForwardedProto",
						Optional:    true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "RoutingAction",
						Optional:    true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "ConditionType",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.LoadBalancer
	d.clients = inst.Client
}

func (d *loadbalancerLbListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbListenerDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show LbListener",
			err.Error(),
		)
		return
	}

	lbListenerState := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, _ := types.ObjectValueFrom(ctx, lbListenerState.AttributeTypes(), lbListenerState)
	state.LbListenerDetail = lbListenerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
