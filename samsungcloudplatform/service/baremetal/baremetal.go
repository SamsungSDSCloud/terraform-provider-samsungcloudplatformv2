package baremetal

import (
	"context"
	"errors"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	baremetalcommon "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpbaremetal1d1 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/baremetal/1.1"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource               = &baremetalBaremetalResource{}
	_ resource.ResourceWithConfigure  = &baremetalBaremetalResource{}
	_ resource.ResourceWithModifyPlan = &baremetalBaremetalResource{}
)

// NewBaremetalBaremetalResource is a helper function to simplify the provider implementation.
func NewBaremetalBaremetalResource() resource.Resource {
	return &baremetalBaremetalResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type baremetalBaremetalResource struct {
	config  *scpsdk.Configuration
	client  *baremetal.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *baremetalBaremetalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_baremetal"
}

// Schema defines the schema for the data source.
func (r *baremetalBaremetalResource) Schema(c context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BaremetalResourceSchema(c)
}

// Configure adds the provider configured client to the data source.
func (r *baremetalBaremetalResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Baremetal
	r.clients = inst.Client
}

func BaremetalResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Baremetal",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
				MarkdownDescription: "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Required:            true,
				Description:         "OS image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
				MarkdownDescription: "OS image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
			},
			"image_version": schema.StringAttribute{
				Computed:            true,
				Description:         "Image version\n  - example: RHEL 8.7 for BM",
				MarkdownDescription: "Image version\n  - example: RHEL 8.7 for BM",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"init_script": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "init script\n" +
					"  - example: init script\n" +
					"  - maxLength: 16384",
				MarkdownDescription: "init script\n" +
					"  - example: init script\n" +
					"  - maxLength: 16384",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(16384),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lock_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Use Lock\n  - example: true",
				MarkdownDescription: "Use Lock\n  - example: true",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
				MarkdownDescription: "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
			},
			"network_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_type": schema.StringAttribute{
				Computed:            true,
				Description:         "OS type\n  - example: WINDOWS",
				MarkdownDescription: "OS type\n  - example: WINDOWS",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_user_id": schema.StringAttribute{
				Required: true,
				Description: "OS User Id. When linux image value must be 'root'\n" +
					"  - example: root123\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 5\n" +
					"  - pattern: ^[a-z0-9]{5,20}$",
				MarkdownDescription: "OS User Id. When linux image value must be 'root'\n" +
					"  - example: root123\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 5\n" +
					"  - pattern: ^[a-z0-9]{5,20}$",
			},
			"os_user_password": schema.StringAttribute{
				Required: true,
				Description: "OS user password.\n" +
					"  - example: P@ssword1!2\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 9\n" +
					"  - pattern: ^[A-Za-z0-9@$!%*#&]+$",
				MarkdownDescription: "OS user password.\n" +
					"  - example: P@ssword1!2\n" +
					"  - maxLength: 20\n" +
					"  - minLength: 9\n" +
					"  - pattern: ^[A-Za-z0-9@$!%*#&]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(9, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[A-Za-z0-9@$!%*#&]+$"), "English Capital letters, lowercase letters, numbers, @$!%*#& special characters only."),
					stringvalidator.RegexMatches(regexp.MustCompile("[A-Z]"), "Must contain capital letters."),
					stringvalidator.RegexMatches(regexp.MustCompile("[a-z]"), "Must contain lowercase characters."),
					stringvalidator.RegexMatches(regexp.MustCompile("[0-9]"), "Must contain number."),
					stringvalidator.RegexMatches(regexp.MustCompile("[@$!%*#&]"), "'@$!%*#&' must contain special characters"),
				},
			},
			"placement_group_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "placement group name\n" +
					"  - example: pg-group\n" +
					"  - maxLength: 15\n" +
					"  - minLength: 3\n" +
					"  - pattern: ^[a-z][a-z0-9-]{1,13}[a-z0-9]$\n",
				MarkdownDescription: "placement group name\n" +
					"  - example: pg-group\n" +
					"  - maxLength: 15\n" +
					"  - minLength: 3\n" +
					"  - pattern: ^[a-z][a-z0-9-]{1,13}[a-z0-9]$\n",
				Default: stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Required:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"root_account": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Account\n  - example: root",
				MarkdownDescription: "Root Account\n  - example: root",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_details": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bare_metal_local_subnet_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "local subnet ID\n  - example: 8d0581b1bbde4195a623abbb1b05700d",
							MarkdownDescription: "local subnet ID\n  - example: 8d0581b1bbde4195a623abbb1b05700d",
							Default:             stringdefault.StaticString(""),
						},
						"bare_metal_local_subnet_ip_address": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "local subnet IP\n  - example: 192.168.3.4",
							MarkdownDescription: "local subnet IP\n  - example: 192.168.3.4",
							Validators: []validator.String{
								baremetalcommon.IpStringValidator{},
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"bare_metal_server_name": schema.StringAttribute{
							Required: true,
							Description: "Bare Metal Server name. The rules are different for Windows image types and linux image types.\n" +
								"  - example: bmserver-001\n" +
								"  - minLength: 3\n" +
								"  - maxLength(for window): 15\n" +
								"  - maxLength(for linux): 28\n" +
								"  - pattern(for window): ^[a-z][a-z0-9-]{1,13}[a-z0-9]$\n" +
								"  - pattern(for linux): ^[a-z][a-z0-9-]{1,26}[a-z0-9]$",
							MarkdownDescription: "Bare Metal Server name. The rules are different for Windows image types and linux image types.\n" +
								"  - example: bmserver-001\n" +
								"  - minLength: 3\n" +
								"  - maxLength(for window): 15\n" +
								"  - maxLength(for linux): 28\n" +
								"  - pattern(for window): ^[a-z][a-z0-9-]{1,13}[a-z0-9]$\n" +
								"  - pattern(for linux): ^[a-z][a-z0-9-]{1,26}[a-z0-9]$",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(3),
							},
						},
						"ip_address": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "subnet IP address\n  - example: 192.168.2.4",
							MarkdownDescription: "subnet IP address\n  - example: 192.168.2.4",
							Validators: []validator.String{
								baremetalcommon.IpStringValidator{},
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"nat_enabled": schema.BoolAttribute{
							Required:            true,
							Description:         "Use Public NAT\n  - example: true",
							MarkdownDescription: "Use Public NAT\n  - example: true",
						},
						"public_ip_address_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "public IP address id\n  - example: a765a07e8d9b46f4918fd7d5ed004654",
							MarkdownDescription: "public IP address id\n  - example: a765a07e8d9b46f4918fd7d5ed004654",
							Default:             stringdefault.StaticString(""),
						},
						"server_type_id": schema.StringAttribute{
							Required:            true,
							Description:         "Server Type ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
							MarkdownDescription: "Server Type ID\n  - example: PRODUCT-0iT9dNiLr4lVoYmjlY2Vgg",
						},
						"use_hyper_threading": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Use Hyper Threading\n  - example: true",
							MarkdownDescription: "Use Hyper Threading\n  - example: true",
							Default:             booldefault.StaticBool(false),
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Bare Metal Server ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
							MarkdownDescription: "Bare Metal Server ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"state": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "Bare Metal Server state\n" +
								"  - example: RUNNING\n" +
								"  - pattern: RUNNING | STOPPED",
							MarkdownDescription: "Bare Metal Server state\n" +
								"  - example: RUNNING\n" +
								"  - pattern: RUNNING | STOPPED",
							Default: stringdefault.StaticString(common.RunningState),
							Validators: []validator.String{
								stringvalidator.OneOf(common.RunningState, common.StoppedState),
							},
						},
					},
				},
				Required: true,
				Description: "Detailed settings for each server\n" +
					"  - example: [{bare_metal_server_name='bm-server', server_type_id='83c3c73d457345e3829ee6d5557c0011', nat_enabled='false'}]\n" +
					"  - maxLength: 5\n" +
					"  - minLength: 1\n",
				MarkdownDescription: "Detailed settings for each server\n" +
					"  - example: [{bare_metal_server_name='bm-server', server_type_id='83c3c73d457345e3829ee6d5557c0011', nat_enabled='false'}]\n" +
					"  - maxLength: 5\n" +
					"  - minLength: 1\n",
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 5),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:            true,
				Description:         "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Subnet ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
			},
			"time_zone": schema.StringAttribute{
				Computed:            true,
				Description:         "Time Zone\n  - example: Asia/Seoul",
				MarkdownDescription: "Time Zone\n  - example: Asia/Seoul",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_placement_group": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Use Placement Group\n  - example: false",
				MarkdownDescription: "Use Placement Group\n  - example: false",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
			"tags": tag.ResourceSchema(),
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *baremetalBaremetalResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state baremetal.BaremetalResource
	diags := req.Plan.Get(ctx, &plan)

	// destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.RegionId.IsNull() {
		r.client.Config.Region = plan.RegionId.ValueString()
		r.clients.Iam.Config.Region = plan.RegionId.ValueString()
	}

	// create validate start
	var planServerDetails []baremetal.ServerDetails
	plan.ServerDetails.ElementsAs(ctx, &planServerDetails, false)

	if planServerDetails[0].Id.IsUnknown() {
		// validate placement group
		err := validatePlacementGroup(plan.UsePlacementGroup.ValueBool(), plan.PlacementGroupName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Placement group value error", err.Error())
		}

		imageList, err := r.client.GetImageList(ctx, plan.RegionId.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error reading image",
				"Could not read image data, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
		}

		osType, e := getImageInfo(imageList, plan.ImageId.ValueString())
		if e != nil {
			resp.Diagnostics.AddError("ImageId error", e.Error())
		}

		// validate os user id
		err = validateOsUserId(osType, plan.OsUserId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("osUserId value error", err.Error())
		}

		// validate server info
		errList := validateServerDetailInfo(planServerDetails, osType)
		if errList != nil {
			for _, e := range errList {
				resp.Diagnostics.AddError("server detail info error", e.Error())
			}
		}

		return
	}
	// create validate end

	// update validate start
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateServerDetails []baremetal.ServerDetails
	state.ServerDetails.ElementsAs(ctx, &stateServerDetails, false)

	//validateUpdateValue()
	changeField := getChangeField(ctx, plan, state)

	for _, field := range changeField {
		if !slices.Contains([]string{"State", "ModifiedAt", "ModifiedBy"}, field) {
			resp.Diagnostics.AddError("Cannot update value", "only State can be changed. "+field+" Cannot be changed")
		}
	}

}

// Create create baremetal servers and sets the initial Terraform state.
func (r *baremetalBaremetalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan baremetal.BaremetalResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.RegionId.IsNull() {
		r.client.Config.Region = plan.RegionId.ValueString()
		r.clients.Iam.Config.Region = plan.RegionId.ValueString()
	}

	// Create baremetal servers
	asyncResponse, err := r.client.CreateBaremetal(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Baremetal",
			"Could not create Baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	baremetalIds := asyncResponse.ResourceId
	baremetalIdList := strings.Split(baremetalIds, ", ")

	for _, baremetalId := range baremetalIdList {
		err = waitForBaremetalStatus(ctx, r.client, baremetalId, []string{common.CreatingState}, []string{common.RunningState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error creating Baremetal",
				"Error waiting for baremetal("+baremetalId+")to become RUNNING:"+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	// get common data for server
	baremetalShowResponse, _, err := r.client.GetBaremetal(ctx, baremetalIdList[0])
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Baremetal",
			"Could not read Baremetal ID "+baremetalIdList[0]+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "baremetal", "baremetal", baremetalIdList[0])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Baremetal Tag", err.Error())
		return
	}
	plan.Tags = tagsMap

	setBaremetalDetailToResource(baremetalShowResponse, &plan)

	// get serverDetail data
	var serverDetails []baremetal.ServerDetails
	plan.ServerDetails.ElementsAs(ctx, &serverDetails, false)
	for idx, baremetalId := range baremetalIdList {
		// Get refreshed value from Baremetal
		baremetalShowResponse, _, err := r.client.GetBaremetal(ctx, baremetalId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Baremetal",
				"Could not read Baremetal ID "+baremetalId+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		getServerDetailInfo(baremetalShowResponse, &serverDetails[idx])
	}
	var serverDetail baremetal.ServerDetails
	serverDetailInfoListType, _ := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: serverDetail.AttributeTypes(),
	}, serverDetails)

	plan.ServerDetails = serverDetailInfoListType

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *baremetalBaremetalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state baremetal.BaremetalResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RegionId.IsNull() {
		r.client.Config.Region = state.RegionId.ValueString()
		r.clients.Iam.Config.Region = state.RegionId.ValueString()
	}

	var requestServerDetails []baremetal.ServerDetails
	state.ServerDetails.ElementsAs(ctx, &requestServerDetails, false)

	baremetalId := requestServerDetails[0].Id.ValueString()
	// Get refreshed value from Baremetal
	baremetalShowResponse, _, err := r.client.GetBaremetal(ctx, baremetalId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Baremetal",
			"Could not read Baremetal ID "+baremetalId+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Get Tags
	tagsMap, err := tag.GetTags(r.clients, "baremetal", "baremetal", baremetalId)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Baremetal Tag", err.Error())
		return
	}
	state.Tags = tagsMap

	setBaremetalDetailToResource(baremetalShowResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
// only start, stop can be executed by update
func (r *baremetalBaremetalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan baremetal.BaremetalResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planServerDetails, stateServerDetails []baremetal.ServerDetails
	plan.ServerDetails.ElementsAs(ctx, &planServerDetails, false)
	state.ServerDetails.ElementsAs(ctx, &stateServerDetails, false)

	if len(planServerDetails) != len(stateServerDetails) {
		resp.Diagnostics.AddError("Error updating Baremetal", "Cannot add or delete serverDetail value.")
		return
	}

	var stopBaremetalIds, startBaremetalIds []string

	for idx, planServer := range planServerDetails {
		if !planServer.State.Equal(stateServerDetails[idx].State) {
			// RUNNING -> STOPPED
			if planServer.State.ValueString() == common.StoppedState {
				stopBaremetalIds = append(stopBaremetalIds, planServer.Id.ValueString())
			} else {
				// STOPPED -> RUNNING
				startBaremetalIds = append(startBaremetalIds, planServer.Id.ValueString())
			}
		}
	}

	if len(stopBaremetalIds) != 0 {
		err := r.client.StopBaremetals(ctx, stopBaremetalIds)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Stopping baremetal",
				"Could not stop baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	if len(startBaremetalIds) != 0 {
		err := r.client.StartBaremetals(ctx, startBaremetalIds)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Starting baremetal",
				"Could not start baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	for _, baremetalId := range stopBaremetalIds {
		err := waitForBaremetalStatus(ctx, r.client, baremetalId, []string{common.StoppingState}, []string{common.StoppedState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Stopping baremetal",
				"Could not stop baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	for _, baremetalId := range startBaremetalIds {
		err := waitForBaremetalStatus(ctx, r.client, baremetalId, []string{common.StartingState}, []string{common.RunningState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Stopping baremetal",
				"Could not stop baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	baremetalShowResponse, _, err := r.client.GetBaremetal(ctx, planServerDetails[0].Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Baremetal",
			"Could not read Baremetal ID "+planServerDetails[0].Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	setBaremetalDetailToResource(baremetalShowResponse, &plan)

	for _, serverDetail := range planServerDetails {
		// Get refreshed value from Baremetal
		baremetalShowResponse, _, err := r.client.GetBaremetal(ctx, serverDetail.Id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Baremetal",
				"Could not read Baremetal ID "+serverDetail.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		getServerDetailInfo(baremetalShowResponse, &serverDetail)
	}

	var serverDetail baremetal.ServerDetails
	serverDetailInfoListType, _ := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: serverDetail.AttributeTypes(),
	}, planServerDetails)

	plan.ServerDetails = serverDetailInfoListType

	// update tags
	tagElements := plan.Tags.Elements()
	for _, serverDetail := range planServerDetails {
		tagsMap, err := tag.UpdateTags(r.clients, "baremetal", "baremetal", serverDetail.Id.ValueString(), tagElements)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating tags",
				err.Error(),
			)
			return
		}
		plan.Tags = tagsMap
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete delete baremetal servers and removes the Terraform state on success.
func (r *baremetalBaremetalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state baremetal.BaremetalResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RegionId.IsNull() {
		r.client.Config.Region = state.RegionId.ValueString()
		r.clients.Iam.Config.Region = state.RegionId.ValueString()
	}

	var serverDetails []baremetal.ServerDetails
	state.ServerDetails.ElementsAs(ctx, &serverDetails, false)

	baremetalIds := make([]string, 0)
	for _, serverDetail := range serverDetails {
		baremetalIds = append(baremetalIds, serverDetail.Id.ValueString())
	}

	// Delete baremetal servers
	err := r.client.DeleteBaremetal(ctx, baremetalIds)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting baremetal",
			"Could not delete baremetal, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	for _, baremetalId := range baremetalIds {
		err = waitForBaremetalStatus(ctx, r.client, baremetalId, []string{common.TerminatingState}, []string{common.TerminatedState})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Deleting Baremetal",
				"Could not delete baremetal("+baremetalId+")unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}
}

func waitForBaremetalStatus(ctx context.Context, baremetalClient *baremetal.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpResponse, err := baremetalClient.GetBaremetal(ctx, id)
		if err != nil {
			if httpResponse != nil {
				if httpResponse.StatusCode == 404 {
					return "", common.TerminatedState, nil
				}
			}
			return nil, "", err
		}
		return info, info.State, nil
	})
}

func setBaremetalDetailToResource(baremetalDetail *scpbaremetal1d1.BaremetalShowResponseV1Dot1, baremetalResource *baremetal.BaremetalResource) {
	//basic info
	baremetalResource.AccountId = types.StringValue(baremetalDetail.AccountId)
	baremetalResource.ImageId = types.StringValue(baremetalDetail.ImageId)
	baremetalResource.ImageVersion = types.StringValue(baremetalDetail.ImageVersion)
	baremetalResource.OsType = types.StringValue(baremetalDetail.OsType)
	baremetalResource.RegionId = types.StringValue(baremetalDetail.RegionId)
	baremetalResource.RootAccount = types.StringValue(baremetalDetail.RootAccount)
	baremetalResource.TimeZone = types.StringValue(baremetalDetail.TimeZone)

	// network info
	baremetalResource.NetworkId = types.StringValue(baremetalDetail.NetworkId)
	baremetalResource.VpcId = types.StringValue(baremetalDetail.VpcId)

	// additional info
	baremetalResource.PlacementGroupName = types.StringPointerValue(baremetalDetail.PlacementGroupName.Get())
	baremetalResource.InitScript = types.StringValue(baremetalDetail.InitScript)
	baremetalResource.LockEnabled = types.BoolPointerValue(baremetalDetail.LockEnabled)

	// metadata info
	baremetalResource.CreatedBy = types.StringValue(baremetalDetail.CreatedBy)
	baremetalResource.CreatedAt = types.StringValue(baremetalDetail.CreatedAt.String())
	baremetalResource.ModifiedBy = types.StringValue(baremetalDetail.ModifiedBy)
	baremetalResource.ModifiedAt = types.StringValue(baremetalDetail.ModifiedAt.String())
}

func getServerDetailInfo(baremetalDetail *scpbaremetal1d1.BaremetalShowResponseV1Dot1, serverDetail *baremetal.ServerDetails) {

	var localSubnetIpaddress string
	if len(baremetalDetail.LocalSubnetInfo) == 0 {
		localSubnetIpaddress = ""
	} else {
		localSubnetIpaddress = baremetalDetail.LocalSubnetInfo[0].PolicyLocalSubnetIp
	}

	serverDetail.Id = types.StringValue(baremetalDetail.Id)
	serverDetail.BareMetalLocalSubnetIpAddress = types.StringValue(localSubnetIpaddress)
	serverDetail.IpAddress = types.StringValue(baremetalDetail.PolicyIp)
	serverDetail.State = types.StringValue(baremetalDetail.State)
}

func validatePlacementGroup(usePlacementGroup bool, placementGroupName string) error {
	if usePlacementGroup == true {
		if placementGroupName == "" {
			return errors.New("if usePlacementGroup value is true, the placementGroupName value should not be empty")
		} else if !regexp.MustCompile("^[a-z][a-z0-9-]{1,13}[a-z0-9]$").MatchString(placementGroupName) {
			return errors.New("placementGroupName value must begin with a lowercase letter and must be entered in 3-15 characters using lowercase letters, numbers, and -. (- does not end with.)")
		}
	} else {
		if placementGroupName != "" {
			return errors.New("if usePlacementGroup value is false, the placementGroupName value should be empty")
		}
	}

	return nil
}

func getImageInfo(imageList *scpbaremetal1d1.ImageListResponse, imageId string) (string, error) {
	var osType string

	for _, image := range imageList.Images {
		if image.Id == imageId {
			osType = image.OsDistro
			break
		}
	}

	if osType == "" {
		return "", errors.New("no image exists for imageId")
	}

	return osType, nil
}

func validateOsUserId(osType string, osUserId string) error {

	if osType == "WINDOWS" {
		if osUserId == "administrator" || osUserId == "sdscmpif" {
			return errors.New("for Windows type, osUserId cannot have 'sdscmpif' or 'administrator' values")
		} else if !regexp.MustCompile("^[a-z0-9]{5,20}$").MatchString(osUserId) {
			return errors.New("osUserId value must be entered in 5-20 characters using lowercase letters, numbers")
		}
	} else {
		if osUserId != "root" {
			return errors.New("osUserId of the " + osType + " image must be root")
		}
	}

	return nil
}

func validateServerDetailInfo(serverDetails []baremetal.ServerDetails, osType string) []error {

	var errorList []error
	var ipList, localSubnetIpList, natIpList, serverNameList []string
	for _, serverDetail := range serverDetails {
		serverNameList = append(serverNameList, serverDetail.BareMetalServerName.ValueString())
		ip := serverDetail.IpAddress.ValueString()
		localSubnetIp := serverDetail.BareMetalLocalSubnetIpAddress.ValueString()
		publicIpAddressId := serverDetail.PublicIpAddressId.ValueString()

		if ip != "" {
			ipList = append(ipList, ip)
		}

		// validate localSubnet
		if serverDetail.BareMetalLocalSubnetId.ValueString() == "" {
			if localSubnetIp != "" {
				errorList = append(errorList, errors.New("if local subnet id is empty, the local subnet ip value must also be empty"))
			}
		} else {
			if localSubnetIp != "" {
				localSubnetIpList = append(localSubnetIpList, localSubnetIp)
			}
		}

		// validate publicIpAddress
		if serverDetail.NatEnabled.ValueBool() == true {
			if publicIpAddressId == "" {
				errorList = append(errorList, errors.New("if natEnabled value is true, publicIpAddressId must exist(not empty)"))
			} else {
				natIpList = append(natIpList, publicIpAddressId)
			}
		} else {
			if publicIpAddressId != "" {
				errorList = append(errorList, errors.New("if natEnabled value is false, publicIpAddressId value must be empty"))
			}
		}

		// validate state
		if serverDetail.State.ValueString() != common.RunningState {
			errorList = append(errorList, errors.New("state value must be RUNNING, when server creation request is made"))
		}
	}

	err := validateServerName(osType, serverNameList)
	if err != nil {
		errorList = append(errorList, err)
	}

	// validate duplicate ipList
	if hasDuplicate(ipList) {
		errorList = append(errorList, errors.New("ip cannot be duplicated"))
	}

	// validate duplicate localSubnetIpList
	if hasDuplicate(localSubnetIpList) {
		errorList = append(errorList, errors.New("local subnet ip cannot be duplicated"))
	}

	if hasDuplicate(natIpList) {
		errorList = append(errorList, errors.New("publicIpAddressId cannot be duplicated"))
	}

	return errorList
}

func validateServerName(osType string, serverNameList []string) error {

	if hasDuplicate(serverNameList) {
		return errors.New("server name cannot be duplicated")
	}

	for _, serverName := range serverNameList {
		if osType == "WINDOWS" {
			if !regexp.MustCompile("^[a-z][a-z0-9-]{1,13}[a-z0-9]$").MatchString(serverName) {
				return errors.New("WINDOWS type: serverName value must begin with a lowercase letter and must be entered in 3-15 characters using lowercase letters, numbers, and -. (- does not end with.)")
			}
		} else {
			if !regexp.MustCompile("^[a-z][a-z0-9-]{1,26}[a-z0-9]$").MatchString(serverName) {
				return errors.New(osType + " type: serverName value must begin with a lowercase letter and must be entered in 3-28 characters using lowercase letters, numbers, and -. (- does not end with.)")
			}
		}
	}

	return nil
}

func hasDuplicate(slice []string) bool {
	m := make(map[string]struct{})
	for _, v := range slice {
		if _, exists := m[v]; exists {
			return true
		}
		m[v] = struct{}{}
	}
	return false
}

func getChangeField(ctx context.Context, plan, state baremetal.BaremetalResource) []string {
	var changedFields []string

	planValue := reflect.ValueOf(plan)
	stateValue := reflect.ValueOf(state)

	for i := 0; i < planValue.NumField(); i++ {
		fieldName := planValue.Type().Field(i).Name
		planField := planValue.FieldByName(fieldName)
		stateField := stateValue.FieldByName(fieldName)

		if planField.Type() == reflect.TypeOf(types.List{}) {
			planList := planField.Interface().(types.List)
			stateList := stateField.Interface().(types.List)

			var planServerDetails, stateServerDetails []baremetal.ServerDetails
			planList.ElementsAs(ctx, &planServerDetails, false)
			stateList.ElementsAs(ctx, &stateServerDetails, false)

			if len(planServerDetails) != len(stateServerDetails) {
				changedFields = append(changedFields, "serverDetails")
				return changedFields
			}

			for idx, _ := range planServerDetails {
				planServerDetailsValue := reflect.ValueOf(planServerDetails[idx])
				stateServerDetailsValue := reflect.ValueOf(stateServerDetails[idx])

				for i := 0; i < planServerDetailsValue.NumField(); i++ {
					fieldName := planServerDetailsValue.Type().Field(i).Name
					planField := planServerDetailsValue.FieldByName(fieldName)
					stateField := stateServerDetailsValue.FieldByName(fieldName)

					if !reflect.DeepEqual(planField.Interface(), stateField.Interface()) {
						changedFields = append(changedFields, fieldName)
					}
				}
			}

		} else if planField.Type() == reflect.TypeOf(types.Map{}) {
			// tag can be changed
			continue
		} else {
			if !reflect.DeepEqual(planField.Interface(), stateField.Interface()) {
				changedFields = append(changedFields, fieldName)
			}
		}
	}

	return changedFields
}
