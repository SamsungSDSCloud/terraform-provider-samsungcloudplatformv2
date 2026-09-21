package cloudcontrol

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkcloudcontrol "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &landingZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &landingZoneDataSource{}
)

func NewLandingZoneDataSource() datasource.DataSource {
	return &landingZoneDataSource{}
}

type landingZoneDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (d *landingZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_landing_zone"
}

func (d *landingZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudControl
	d.clients = inst.Client
}

func (d *landingZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show CloudControl Landing Zone",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Required: true,
				Description: "Unique identifier of the landing zone. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"landing_zone": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Landing Zone Detail Info. \n" +
					"  - example : '{id: 2c8a138f8d78e1fc29a449dbfa8681, service_name: Cloud Control, ...}' \n",
				Attributes: map[string]schema.Attribute{
					"agree_yn": schema.StringAttribute{
						Computed: true,
						Description: "Terms Agreement YN. \n" +
							"  - example : 'Y' \n" +
							"  - allowed_values : [Y, N] \n",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the landing zone was created. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the landing zone. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the landing zone creator. \n" +
							"  - example : 'John Doe' \n",
					},
					"detective_guardrail_status": schema.StringAttribute{
						Computed: true,
						Description: "Status of detective guardrails. \n" +
							"  - example : 'ENABLED' \n" +
							"  - allowed_values : [ENABLED, DISABLED, ENABLING, DISABLING] \n",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the landing zone. \n" +
							"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
					},
					"identity_center_id": schema.StringAttribute{
						Computed: true,
						Description: "Identity Center Instance ID. \n" +
							"  - example : '0xnw6g1xh2q5' \n",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the landing zone was modified. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who modified the landing zone. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the landing zone modifier. \n" +
							"  - example : 'Alice' \n",
					},
					"organization_id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
					},
					"region": schema.StringAttribute{
						Computed: true,
						Description: "Region of the landing zone. \n" +
							"  - example : 'kr-west1' \n",
					},
					"service_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the CloudControl service. \n" +
							"  - example : 'Cloud Control' \n",
					},
					"srn": schema.StringAttribute{
						Computed: true,
						Description: "Samsung Resource Name (SRN) uniquely identifying this resource. \n" +
							"  - example : 'srn:dev2::a8a2f3c2659646ecaaf28fc8f783921a:::cloudcontrol:landingzone/f7a1ef0b17e34a37811cc2fa7a6bd50b' \n",
					},
					"sso_type": schema.StringAttribute{
						Computed: true,
						Description: "SSO management type of the landing zone - self-managed or managed by Identity Center. \n" +
							"  - allowed_values : [ID_CENTER, SELF] \n" +
							"  - example : 'ID_CENTER' \n",
					},
					"status": schema.StringAttribute{
						Computed: true,
						Description: "Landing Zone Status. \n" +
							"  - example : 'ACTIVE' \n" +
							"  - allowed_values : [ACTIVE, CREATING, CREATE_FAILED, DELETE_FAILED, DELETING, INACTIVE] \n",
					},
					"version_id": schema.StringAttribute{
						Computed: true,
						Description: "Landing Zone Version ID. \n" +
							"  - example : '1.0' \n",
					},
				},
			},
		},
	}
}

func (d *landingZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.LandingZoneDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()
	if landingZoneId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Landing Zone",
			"Landing Zone ID is empty",
		)
		return
	}

	data, err := d.client.GetLandingZone(ctx, landingZoneId)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Landing Zone",
			"Could not read CloudControl Landing Zone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	landingZoneValue, landingZoneDiags := d.buildLandingZoneValue(ctx, &data.LandingZone)
	resp.Diagnostics.Append(landingZoneDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.LandingZone = landingZoneValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *landingZoneDataSource) buildLandingZoneValue(ctx context.Context, lz *sdkcloudcontrol.LandingZoneV1Dot1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var creatorName, identityCenterId, modifierName string
	if lz.CreatorName != nil {
		creatorName = *lz.CreatorName
	}
	if lz.IdentityCenterId.IsSet() && lz.IdentityCenterId.Get() != nil {
		identityCenterId = *lz.IdentityCenterId.Get()
	}
	if lz.ModifierName != nil {
		modifierName = *lz.ModifierName
	}

	landingZoneValue, landingZoneDiags := types.ObjectValue(cloudcontrol.LandingZoneValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"agree_yn":                   types.StringValue(lz.AgreeYn),
		"created_at":                 types.StringValue(lz.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":                 types.StringValue(lz.CreatedBy),
		"creator_name":               types.StringValue(creatorName),
		"detective_guardrail_status": types.StringValue(string(lz.DetectiveGuardrailStatus)),
		"id":                         types.StringValue(lz.Id),
		"identity_center_id":         types.StringValue(identityCenterId),
		"modified_at":                types.StringValue(lz.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":                types.StringValue(lz.ModifiedBy),
		"modifier_name":              types.StringValue(modifierName),
		"organization_id":            types.StringValue(lz.OrganizationId),
		"region":                     types.StringValue(lz.Region),
		"service_name":               types.StringValue(lz.ServiceName),
		"srn":                        types.StringValue(lz.Srn),
		"sso_type":                   types.StringValue(lz.SsoType),
		"status":                     types.StringValue(lz.Status),
		"version_id":                 types.StringValue(lz.VersionId),
	})
	diags.Append(landingZoneDiags...)
	return landingZoneValue, diags
}
