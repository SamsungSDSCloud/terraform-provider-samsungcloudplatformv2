package certificatemanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &certificateManagerSelfSignResource{}
	_ resource.ResourceWithConfigure   = &certificateManagerSelfSignResource{}
	_ resource.ResourceWithImportState = &certificateManagerSelfSignResource{}
)

func NewCertificateManagerSelfSignResource() resource.Resource {
	return &certificateManagerSelfSignResource{}
}

type certificateManagerSelfSignResource struct {
	config *scpsdk.Configuration
	client *certificatemanager.Client
}

func (r *certificateManagerSelfSignResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// This is a no-op implementation
	response.Diagnostics.AddError(
		"Update not supported",
		"This resource does not support in-place updates.",
	)
}

// Metadata returns the data source type name.
func (r *certificateManagerSelfSignResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_manager_self_sign"
}

// Schema defines the schema for the data source.
func (r *certificateManagerSelfSignResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "certificate manager",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Cn"): schema.StringAttribute{
				Description: "Certificate Common Name.\n" +
					"  - example : 'test.go.kr'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Certificate Name.\n" +
					"  - example : 'test-certificate'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
				Description: "Certificate Expire Date.\n" +
					"  - example : '20251212'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
				Description: "Certificate Start Date.\n" +
					"  - example : '20250101'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Organization"): schema.StringAttribute{
				Description: "Certificate Organization Name.\n" +
					"  - example : 'samsungSDS'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Recipients"): schema.ListAttribute{
				Description: "Expired certificates Recipients",
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Region"): schema.StringAttribute{
				Description: "Name of region.\n" +
					"  - example : 'west1'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description: "Timezone indentifier.\n" +
					"  - example : 'Asia/Seoul'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Certificate"): schema.SingleNestedAttribute{
				Description: "Certificate detail",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CertKind"): schema.StringAttribute{
						Description: "Certificate type.\n" +
							"  - example : 'DEV'",
						Computed: true,
					},
					common.ToSnakeCase("Cn"): schema.StringAttribute{
						Description: "Certificate Common Name.\n" +
							"  - example : 'test.go.kr'",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Certificate ID.\n" +
							"  - example : '0fdd87aab8cb46f59b7c1f81ed03fb3e'",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Certificate Name.\n" +
							"  - example : 'test-certificate'",
						Computed: true,
					},
					common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
						Description: "Certificate Expire Date.\n" +
							"  - example : '2026-02-07T18:07:59'",
						Computed: true,
					},
					common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
						Description: "Certificate Start Date.\n" +
							"  - example : '2025-02-08T18:07:00'",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Certificate State.\n" +
							"  - example : 'VALID'",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *certificateManagerSelfSignResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.CertificateManager
}

func (r *certificateManagerSelfSignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificatemanager.CertificateManagerSelfSignResource

	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new certificate manager
	data, err := r.client.SelfSignCreateCertificateManager(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating certificate manager self sign",
			"Could not create certificate manager, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error creating certificate manager self sign",
			"An error occurred while creating certificate manager self sign. Empty response",
		)
		return
	}

	plan.Id = types.StringValue(data.Certificate.Id)
	vgModel := certificatemanager.Certificate{
		Id:          types.StringValue(data.Certificate.Id),
		Name:        types.StringValue(data.Certificate.Name),
		CertKind:    types.StringPointerValue(data.Certificate.CertKind),
		Cn:          types.StringValue(data.Certificate.Cn),
		NotBeforeDt: types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		NotAfterDt:  types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		State:       types.StringValue(data.Certificate.State),
	}

	certificateObjectValue, dia := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Certificate = certificateObjectValue

	diags = resp.State.Set(ctx, plan)

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *certificateManagerSelfSignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state certificatemanager.CertificateManagerSelfSignResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetCertificateManager(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading certificate manager",
			"Could not read certificate manager ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading certificate manager",
			"An error occurred while reading certificate manager. Empty response",
		)
		return
	}

	vgModel := certificatemanager.Certificate{
		Id:          types.StringValue(data.Certificate.Id),
		Name:        types.StringValue(data.Certificate.Name),
		CertKind:    types.StringPointerValue(data.Certificate.CertKind),
		Cn:          types.StringValue(data.Certificate.Cn),
		NotBeforeDt: types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		NotAfterDt:  types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		State:       types.StringValue(data.Certificate.State),
	}

	vgObjectValue, dia := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Certificate = vgObjectValue

	// Update top-level input fields from API response for drift detection.
	// region/timezone are not returned by the API
	state.Cn = types.StringValue(data.Certificate.Cn)
	state.Name = types.StringValue(data.Certificate.Name)
	state.NotBeforeDt = types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339))
	state.NotAfterDt = types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339))
	state.Organization = types.StringValue(data.Certificate.Organization)

	// Region and Timezone are not included in the API response — retain the input values.

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificateManagerSelfSignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state certificatemanager.CertificateManagerSelfSignResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificateManager(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting certificate manager",
			"Could not delete certificate manager, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *certificateManagerSelfSignResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
