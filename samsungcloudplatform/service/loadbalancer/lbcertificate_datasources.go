package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbCertificateDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbCertificateDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbCertificateDataSources() datasource.DataSource {
	return &loadbalancerLbCertificateDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLbCertificateDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbCertificateDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_certificates" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbCertificateDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of lb certificate.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("LbCertificates"): schema.ListNestedAttribute{
				Description: "A list of Lb Certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CertKind"): schema.StringAttribute{
							Description: "CertKind",
							Optional:    true,
						},
						common.ToSnakeCase("Cn"): schema.StringAttribute{
							Description: "Cn",
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
						common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
							Description: "NotAfterDt",
							Optional:    true,
						},
						common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
							Description: "NotBeforeDt",
							Optional:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbCertificateDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadbalancerLbCertificateDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadbalancer.LbCertificateDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbCertificateList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read LbCertificates",
			err.Error(),
		)
		return
	}

	for _, lbcertificate := range data.Certificates {

		lbcertificateState := loadbalancer.LbCertificate{
			CertKind:    types.StringValue(*lbcertificate.CertKind),
			Cn:          types.StringValue(*lbcertificate.Cn),
			CreatedAt:   types.StringValue(lbcertificate.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(lbcertificate.CreatedBy),
			Id:          types.StringValue(lbcertificate.Id),
			ModifiedAt:  types.StringValue(lbcertificate.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(lbcertificate.ModifiedBy),
			Name:        types.StringValue(*lbcertificate.Name),
			NotAfterDt:  types.StringValue(lbcertificate.NotAfterDt.Format(time.RFC3339)),
			NotBeforeDt: types.StringValue(lbcertificate.NotBeforeDt.Format(time.RFC3339)),
			State:       types.StringValue(*lbcertificate.State),
		}

		state.LbCertificates = append(state.LbCertificates, lbcertificateState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
