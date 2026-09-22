package multinodegpucluster

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &GpunodePublicNatIpResource{}
	_ resource.ResourceWithConfigure   = &GpunodePublicNatIpResource{}
	_ resource.ResourceWithImportState = &GpunodePublicNatIpResource{}
)

func NewGpunodePublicNatIpResource() resource.Resource {
	return &GpunodePublicNatIpResource{}
}

type GpunodePublicNatIpResource struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

func (r *GpunodePublicNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunode_public_nat_ip"
}

func (r *GpunodePublicNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expect *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = inst.Client.Mngc
	r.clients = inst.Client
}

func (r *GpunodePublicNatIpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "GPU Node Public NAT IP Resource. Assign(create)/Release(delete) a public NAT IP to a GPU Node.",
		Attributes: map[string]schema.Attribute{
			"gpu_node_id": schema.StringAttribute{
				Required:            true,
				Description:         "GPU Node ID to assign the public NAT IP to.\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "GPU Node ID to assign the public NAT IP to.\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_ip_address_id": schema.StringAttribute{
				Required:            true,
				Description:         "Public IP address ID to assign as the public NAT IP.\n  - example: b1d1a1c2e3f44a5b8c9d0e1f2a3b4c5d",
				MarkdownDescription: "Public IP address ID to assign as the public NAT IP.\n  - example: b1d1a1c2e3f44a5b8c9d0e1f2a3b4c5d",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_nat": schema.StringAttribute{
				Computed:            true,
				Description:         "Assigned public NAT IP address.\n  - example: 203.0.113.1",
				MarkdownDescription: "Assigned public NAT IP address.\n  - example: 203.0.113.1",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_use_nat": schema.BoolAttribute{
				Computed:            true,
				Description:         "Whether the GPU Node uses public NAT.\n  - example: true",
				MarkdownDescription: "Whether the GPU Node uses public NAT.\n  - example: true",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *GpunodePublicNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan multinodegpuclusterClient.GpuNodePublicNatIpResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	gpunodeId := plan.GpuNodeId.ValueString()

	_, err := r.client.AssignGpuNodePublicNatIp(ctx, gpunodeId, plan.PublicIpAddressId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error assigning GPU Node Public NAT IP",
			"Could not assign public NAT IP, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	err = waitForGpuNodeNatStatus(ctx, r.client, gpunodeId, true)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error assigning GPU Node Public NAT IP",
			"Error waiting for GPU Node("+gpunodeId+") public NAT IP to be assigned: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	if err := setGpuNodeNatInfo(ctx, r, gpunodeId, &plan); err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+gpunodeId+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GpunodePublicNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state multinodegpuclusterClient.GpuNodePublicNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gpunodeId := state.GpuNodeId.ValueString()

	if err := setGpuNodeNatInfo(ctx, r, gpunodeId, &state); err != nil {
		if err == errNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+gpunodeId+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	if !state.PolicyUseNat.ValueBool() {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GpunodePublicNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan multinodegpuclusterClient.GpuNodePublicNatIpResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := setGpuNodeNatInfo(ctx, r, plan.GpuNodeId.ValueString(), &plan); err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading GPU Node",
			"Could not read GPU Node ID "+plan.GpuNodeId.ValueString()+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GpunodePublicNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state multinodegpuclusterClient.GpuNodePublicNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	gpunodeId := state.GpuNodeId.ValueString()

	_, err := r.client.ReleaseGpuNodePublicNatIp(ctx, gpunodeId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error releasing GPU Node Public NAT IP",
			"Could not release public NAT IP, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	err = waitForGpuNodeNatStatus(ctx, r.client, gpunodeId, false)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error releasing GPU Node Public NAT IP",
			"Error waiting for GPU Node("+gpunodeId+") public NAT IP to be released: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}
}

func (r *GpunodePublicNatIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("gpu_node_id"), req, resp)
}

func setGpuNodeNatInfo(ctx context.Context, r *GpunodePublicNatIpResource, gpunodeId string, res *multinodegpuclusterClient.GpuNodePublicNatIpResource) error {
	gpunodeShow, httpResponse, err := r.client.GetGpuNode(ctx, gpunodeId)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			return errNotFound
		}
		return err
	}

	if gpunodeShow.PolicyNat.IsSet() {
		res.PolicyNat = types.StringPointerValue(gpunodeShow.PolicyNat.Get())
	} else {
		res.PolicyNat = types.StringNull()
	}
	res.PolicyUseNat = types.BoolPointerValue(gpunodeShow.PolicyUseNat)

	return nil
}

func waitForGpuNodeNatStatus(ctx context.Context, mngcClient *multinodegpuclusterClient.Client, gpunodeId string, wantNat bool) error {
	const (
		pendingState = "PENDING"
		doneState    = "DONE"
	)

	return client.WaitForStatus(ctx, nil, []string{pendingState}, []string{doneState}, func() (interface{}, string, error) {
		info, httpResponse, err := mngcClient.GetGpuNode(ctx, gpunodeId)
		if err != nil {
			if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
				return nil, "", errNotFound
			}
			return nil, "", err
		}

		if info.State == common.ErrorState {
			return info, info.State, fmt.Errorf("GPU Node %s entered %s state during public NAT IP operation", gpunodeId, common.ErrorState)
		}

		if info.State == common.EditingState {
			return info, pendingState, nil
		}

		useNat := info.PolicyUseNat != nil && *info.PolicyUseNat
		if useNat == wantNat {
			return info, doneState, nil
		}

		action := "assign"
		if !wantNat {
			action = "release"
		}
		return info, info.State, fmt.Errorf(
			"public NAT IP %s for GPU Node %s did not take effect (rolled back by server): state=%s, policy_use_nat=%t",
			action, gpunodeId, info.State, useNat)
	}, -1, -1, -1, -1)
}
