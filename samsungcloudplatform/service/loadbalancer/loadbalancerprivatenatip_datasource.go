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
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLoadbalancerPrivateNatIpDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLoadbalancerPrivateNatIpDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerPrivateNatIpDataSource() datasource.DataSource {
	return &loadbalancerLoadbalancerPrivateNatIpDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLoadbalancerPrivateNatIpDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLoadbalancerPrivateNatIpDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer_private_nat_ip" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLoadbalancerPrivateNatIpDataSource) Schema(_ context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Show Loadbalancer Private NAT.",
		Attributes: map[string]schema.Attribute{

			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "loadbalancer id",
				Optional:    true,
			},
			//  - Computed:true 로 읽기 전용 값임을 명시
			//  - Optional:true 로 “값이 없을 경우 null” 을 허용
			common.ToSnakeCase("LoadbalancerPrivateNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of private NAT.",
				Optional:    true, // null 허용
				Computed:    true, // 읽기 전용이지만 값이 없을 때는 null 로 반환
				Attributes: map[string]schema.Attribute{

					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
						Description: "ExternalIpAddress",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("InternalIpAddress"): schema.StringAttribute{
						Description: "InternalIpAddress",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("PrivateNatIpId"): schema.StringAttribute{
						Description: "PrivateNatIpId",
						Optional:    true,
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLoadbalancerPrivateNatIpDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *loadbalancerLoadbalancerPrivateNatIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LoadbalancerPrivateNatDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLoadbalancerPrivateNatIp(ctx, state.LoadbalancerId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Loadbalancer private NAT IP",
			err.Error(),
		)
		return
	}

	// lbStaticNat이 nil일 경우
	if data.StaticNat.Get() == nil {
		// LoadbalancerPrivateNatIp을 null로 설정
		state.LoadbalancerPrivateNatIp = types.ObjectNull(
			loadbalancer.LoadbalancerPrivateNatIpDetail{}.AttributeTypes(),
		)
	} else {
		lbStaticNat := data.StaticNat
		var loadbalancerPrivateNatState = loadbalancer.LoadbalancerPrivateNatIpDetail{
			CreatedAt:         types.StringValue(lbStaticNat.Get().CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(lbStaticNat.Get().CreatedBy),
			ExternalIpAddress: loadbalancerutil.ToNullableStringValue(lbStaticNat.Get().ExternalIpAddress.Get()),
			Id:                loadbalancerutil.ToNullableStringValue(lbStaticNat.Get().Id.Get()),
			InternalIpAddress: loadbalancerutil.ToNullableStringValue(lbStaticNat.Get().InternalIpAddress.Get()),
			ModifiedAt:        types.StringValue(lbStaticNat.Get().ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(lbStaticNat.Get().ModifiedBy),
			PrivateNatIpId:    loadbalancerutil.ToNullableStringValue(lbStaticNat.Get().PrivateNatIpId.Get()),
			State:             loadbalancerutil.ToNullableStringValue(lbStaticNat.Get().State.Get()),
		}

		loadbalancerObjectValue, _ := types.ObjectValueFrom(ctx, loadbalancerPrivateNatState.AttributeTypes(), loadbalancerPrivateNatState)
		state.LoadbalancerPrivateNatIp = loadbalancerObjectValue
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
