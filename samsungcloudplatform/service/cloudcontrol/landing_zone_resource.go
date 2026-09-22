package cloudcontrol

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkcloudcontrol "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &landingZoneResource{}
	_ resource.ResourceWithConfigure   = &landingZoneResource{}
	_ resource.ResourceWithImportState = &landingZoneResource{}
)

func NewLandingZoneResource() resource.Resource {
	return &landingZoneResource{}
}

type landingZoneResource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (r *landingZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_landing_zone"
}

func (r *landingZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Landing Zone - creates and manages a secure multi-account cloud environment. \n",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Description: "Unique identifier of the landing zone. \n" +
					"  - example : '79e6a69c28124f6199747f96ce6bd927' \n",
				Computed: true,
			},
			"job_id": schema.StringAttribute{
				Description: "Job ID for the landing zone creation/deletion operation. \n" +
					"  - example : 'f00dbcb9a5104fbfb6da3350ffcd129e' \n",
				Computed: true,
			},
			"additional_ou_name": schema.StringAttribute{
				Description: "Name of the additional organization unit to create. \n" +
					"  - example : 'Sandbox' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"agree_yn": schema.StringAttribute{
				Description: "Agreement to landing zone terms and conditions. \n" +
					"  - example : 'Y' \n" +
					"  - allowed_values : ['Y', 'N'] \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"audit_account_name": schema.StringAttribute{
				Description: "Name for the audit account. \n" +
					"  - example : 'AUDIT' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"audit_login_id": schema.StringAttribute{
				Description: "Login ID for the audit account. \n" +
					"  - example : 'audit@samsung.com' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"basic_ou_name": schema.StringAttribute{
				Description: "Basic Organization Unit Name. \n" +
					"  - example : 'Security' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"detective_guardrail_status": schema.StringAttribute{
				Description: "Status of detective guardrails (detects policy violations). \n" +
					"  - example : 'ENABLED' \n" +
					"  - valid values : ENABLED, DISABLED, ENABLING, DISABLING",
				Required: true,
			},
			"log_archive_account_name": schema.StringAttribute{
				Description: "Name for the log archive account. \n" +
					"  - example : 'LOG_ARCHIVE' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"log_archive_login_id": schema.StringAttribute{
				Description: "Login ID for the log archive account. \n" +
					"  - example : 'log-archive@samsung.com' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sso_type": schema.StringAttribute{
				Description: "Access management type for landing zone accounts. \n" +
					"  - example : 'ID_CENTER' \n" +
					"  - allowed_values : ['ID_CENTER', 'SELF'] \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the landing zone was created. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Description: "User who created the landing zone. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
				Computed: true,
			},
			"creator_name": schema.StringAttribute{
				Description: "Name of the landing zone creator. \n" +
					"  - example : 'John Doe' \n",
				Computed: true,
			},
			"identity_center_id": schema.StringAttribute{
				Description: "Identity Center Instance ID. \n" +
					"  - example : '0xnw6g1xh2q5' \n",
				Computed: true,
			},
			"modified_at": schema.StringAttribute{
				Description: "Timestamp when the landing zone was modified. \n" +
					"  - example : '2025-01-01T00:00:00.000Z' \n",
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Description: "User who modified the landing zone. \n" +
					"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
				Computed: true,
			},
			"modifier_name": schema.StringAttribute{
				Description: "Name of the landing zone modifier. \n" +
					"  - example : 'Alice' \n",
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-2b63982e88b74dbcb71ee972b13e2ce1' \n",
				Computed: true,
			},
			"region": schema.StringAttribute{
				Description: "Region of the landing zone. \n" +
					"  - example : 'kr-west1' \n",
				Computed: true,
			},
			"service_name": schema.StringAttribute{
				Description: "Name of the CloudControl service this landing zone belongs to. \n" +
					"  - example : 'Cloud Control' \n",
				Computed: true,
			},
			"srn": schema.StringAttribute{
				Description: "Samsung Resource Name (SRN) of the landing zone. \n" +
					"  - example : 'srn:dev2::a8a2f3c2659646ecaaf28fc8f783921a:::cloudcontrol:landingzone/f7a1ef0b17e34a37811cc2fa7a6bd50b' \n",
				Computed: true,
			},
			"status": schema.StringAttribute{
				Description: "Landing Zone Status. \n" +
					"  - example : 'ACTIVE' \n" +
					"  - allowed_values : ['ACTIVE', 'INACTIVE', 'CREATING', 'CREATE_FAILED', 'DELETING', 'DELETE_FAILED'] \n",
				Computed: true,
			},
			"version_id": schema.StringAttribute{
				Description: "Landing Zone Version. \n" +
					"  - example : '1.0' \n",
				Computed: true,
			},
		},
	}
}

func (r *landingZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.CloudControl
	r.clients = inst.Client
}

func (r *landingZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudcontrol.LandingZoneResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateLandingZone(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Landing Zone",
			"Could not create Landing Zone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.LandingZoneId = types.StringValue(data.LandingZoneId)
	plan.JobId = types.StringValue(data.JobId)

	err = waitForLandingZoneActiveStatus(ctx, r.client, data.LandingZoneId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for landing zone creation",
			"Error waiting for landing zone "+data.LandingZoneId+" to become active: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func waitForLandingZoneActiveStatus(ctx context.Context, ccClient *cloudcontrol.Client, landingZoneId string) error {
	return client.WaitForStatus(ctx, nil, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		data, err := ccClient.GetLandingZone(ctx, landingZoneId)
		if err != nil {
			return nil, "", err
		}
		status := data.LandingZone.Status
		if status == "CREATE_FAILED" {
			return data, status, fmt.Errorf("landing zone creation failed")
		}
		return data, status, nil
	}, -1, -1, -1, -1)
}

func (r *landingZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudcontrol.LandingZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetLandingZone(ctx, state.LandingZoneId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Landing Zone",
			"Could not read Landing Zone ID "+state.LandingZoneId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lz := data.LandingZone

	state.LandingZoneId = types.StringValue(lz.Id)

	state.AgreeYn = types.StringValue(lz.AgreeYn)
	state.DetectiveGuardrailStatus = types.StringValue(string(lz.DetectiveGuardrailStatus))
	state.SsoType = types.StringValue(lz.SsoType)

	state.CreatedAt = types.StringValue(lz.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(lz.CreatedBy)
	state.CreatorName = types.StringValue(lz.GetCreatorName())

	identityCenterId := lz.IdentityCenterId.Get()
	if identityCenterId == nil || *identityCenterId == "" {
		state.IdentityCenterId = types.StringNull()
	} else {
		state.IdentityCenterId = types.StringValue(*identityCenterId)
	}

	state.ModifiedAt = types.StringValue(lz.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(lz.ModifiedBy)
	state.ModifierName = types.StringValue(lz.GetModifierName())
	state.OrganizationId = types.StringValue(lz.OrganizationId)
	state.Region = types.StringValue(lz.Region)
	state.ServiceName = types.StringValue(lz.ServiceName)
	state.Srn = types.StringValue(lz.Srn)
	state.Status = types.StringValue(lz.Status)
	state.VersionId = types.StringValue(lz.VersionId)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *landingZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state cloudcontrol.LandingZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan cloudcontrol.LandingZoneResource
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()

	updateReq := cloudcontrol.LandingZoneUpdateResource{
		Id:                       types.StringValue(landingZoneId),
		DetectiveGuardrailStatus: plan.DetectiveGuardrailStatus,
	}

	updateResp, err := r.client.UpdateLandingZone(ctx, landingZoneId, updateReq)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Landing Zone",
			"Could not update Landing Zone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	targetStatus := plan.DetectiveGuardrailStatus.ValueString()
	previousStatus := state.DetectiveGuardrailStatus.ValueString()

	var data *sdkcloudcontrol.LandingZoneShowResponseV1Dot1
	maxRetries := 60
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		data, err = r.client.GetLandingZone(ctx, landingZoneId)
		if err != nil {
			break
		}

		currentStatus := string(data.LandingZone.DetectiveGuardrailStatus)

		if currentStatus == targetStatus {
			break
		}
		if currentStatus == "ENABLING" || currentStatus == "DISABLING" {
			time.Sleep(retryInterval)
			continue
		}
		if currentStatus == previousStatus {
			time.Sleep(retryInterval)
			continue
		}
		break
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unable to Read Landing Zone",
			"Could not read Landing Zone ID "+landingZoneId+": "+err.Error(),
		)
		return
	}

	plan.LandingZoneId = types.StringValue(updateResp.LandingZoneId)
	plan.JobId = types.StringValue(updateResp.JobId)
	if data != nil {
		plan.DetectiveGuardrailStatus = types.StringValue(string(data.LandingZone.DetectiveGuardrailStatus))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	resp.State = readResp.State
}

func (r *landingZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudcontrol.LandingZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteLandingZone(ctx, state.LandingZoneId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Landing Zone",
			"Could not delete Landing Zone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *landingZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("landing_zone_id"), req, resp)
}
