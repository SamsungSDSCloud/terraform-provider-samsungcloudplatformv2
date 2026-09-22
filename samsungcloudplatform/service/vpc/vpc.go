package vpc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcVpcResource{}
	_ resource.ResourceWithConfigure   = &vpcVpcResource{}
	_ resource.ResourceWithImportState = &vpcVpcResource{}
	_ resource.ResourceWithModifyPlan  = &vpcVpcResource{}
)

// NewVpcVpcResource is a helper function to simplify the provider implementation.
func NewVpcVpcResource() resource.Resource {
	return &vpcVpcResource{}
}

// vpcVpcResource is the data source implementation.
type vpcVpcResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (r *vpcVpcResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

// Metadata returns the data source type name.
func (r *vpcVpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc"
}

// Schema defines the schema for the data source.
func (r *vpcVpcResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource of vpc",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "The IP address range of the network in CIDR notation.\n" +
					"  - example : 192.167.0.0/18\n" +
					"  - maxMask : /24\n" +
					"  - minMask : /16",
				MarkdownDescription: "The IP address range of the network in CIDR notation.\n" +
					"  - example : 192.167.0.0/18\n" +
					"  - maxMask : /24\n" +
					"  - minMask : /16",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Enter a brief explanation or note about this vpc. This help identify the purpose or usage of the vpc.\n" +
					"  - example : VPC description\n" +
					"  - maxLength : 50",
				MarkdownDescription: "Enter a brief explanation or note about this vpc. This help identify the purpose or usage of the vpc.\n" +
					"  - example : VPC description\n" +
					"  - maxLength : 50",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the vpc.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The unique identifier of the vpc.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the vpc.\n" +
					"  - example : vpcName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "The name of the vpc.\n" +
					"  - example : vpcName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9-]*$"), "Enter 3-20 chars. (English, number, hyphen)"),
				},
				Required: true,
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
				Description:         "Detail information about Vpc.",
				MarkdownDescription: "Detail information about Vpc.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the vpc.\n" +
							"  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						MarkdownDescription: "The identifier of the account that owns the vpc.\n" +
							"  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("CidrCount"): schema.Int32Attribute{
						Description: "The number of CIDR blocks associated with the vpc.\n" +
							"  - example: 20",
						MarkdownDescription: "The number of CIDR blocks associated with the vpc.\n" +
							"  - example: 20",
						Computed: true,
					},
					common.ToSnakeCase("cidrs"): schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"cidr": schema.StringAttribute{
									Computed: true,
									Description: "The IP address range of the network in CIDR notation.\n" +
										"  - example: 192.167.0.0/18",
									MarkdownDescription: "The IP address range of the network in CIDR notation.\n" +
										"  - example: 192.167.0.0/18",
								},
								"created_at": schema.StringAttribute{
									Computed: true,
									Description: "The timestamp when the vpc was created in ISO 8601 format.\n" +
										"  - example: 2024-05-17T00:23:17Z",
									MarkdownDescription: "The timestamp when the vpc was created in ISO 8601 format.\n" +
										"  - example: 2024-05-17T00:23:17Z",
								},
								"created_by": schema.StringAttribute{
									Computed: true,
									Description: "The user id that created the vpc.\n" +
										"  - example: 7df8abb4912e4709b1cb237daccca7a8",
									MarkdownDescription: "The user id that created the vpc.\n" +
										"  - example: 7df8abb4912e4709b1cb237daccca7a8",
								},
								"id": schema.StringAttribute{
									Computed: true,
									Description: "The unique identifier of the vpc.\n" +
										"  - example: 7df8abb4912e4709b1cb237daccca7a8",
									MarkdownDescription: "The unique identifier of the vpc.\n" +
										"  - example: 7df8abb4912e4709b1cb237daccca7a8",
								},
							},
						},
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the vpc was created in ISO 8601 format.\n" +
							"  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "The timestamp when the vpc was created in ISO 8601 format.\n" +
							"  - example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this vpc. This help identify the purpose or usage of the vpc.\n" +
							"  - example: vpcDescription\n" +
							"  - maxLength: 50",
						MarkdownDescription: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the vpc.\n" +
							"  - example: vpcDescription\n" +
							"  - maxLength: 50",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the vpc.\n" +
							"  - example: 7df8abb4912e4709b1cb237daccca7a8",
						MarkdownDescription: "The unique identifier of the vpc.\n" +
							"  - example: 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the vpc was last modified in ISO 8601 format.\n" +
							"  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "The timestamp when the vpc was last modified in ISO 8601 format.\n" +
							"  - example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the vpc.\n" +
							"  - maxLength: 20\n" +
							"  - minLength: 3\n" +
							"  - pattern: `^[a-zA-Z0-9-]*$`\n" +
							"  - example: vpcName",
						MarkdownDescription: "The name of the vpc.\n" +
							"  - maxLength: 20\n" +
							"  - minLength: 3\n" +
							"  - pattern: `^[a-zA-Z0-9-]*$`\n" +
							"  - example: vpcName",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the vpc.\n" +
							"  - enum: [\"ACTIVE\",\"ERROR\"]\n" +
							"  - example : ACTIVE",
						MarkdownDescription: "The current lifecycle state of the vpc.\n" +
							"  - enum: [\"ACTIVE\",\"ERROR\"]\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("ZoneType"): schema.StringAttribute{
						Description: "The zone type of the vpc.\n" +
							"  - example: GLOBAL",
						MarkdownDescription: "The zone type of the vpc.\n" +
							"  - example: GLOBAL",
						Computed: true,
					},
					common.ToSnakeCase("Zones"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "The list of availability zones associated with the vpc.\n" +
							"  - example: [\"kr-1\",\"kr-2\"]",
						MarkdownDescription: "The list of availability zones associated with the vpc.\n" +
							"  - example: [\"kr-1\",\"kr-2\"]",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcVpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1Dot3
	r.clients = inst.Client
}

// Create creates the vpc and sets the initial Terraform state.
func (r *vpcVpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.VpcResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpc(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc",
			"Could not create vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if err = waitForVpcStatus(ctx, r.client, data.Vpc.Id, []string{}, []string{"ACTIVE"}); err != nil {
		resp.Diagnostics.AddError("Error creating vpc", "wait for ACTIVE failed: "+err.Error())
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(data.Vpc.Id)

	vpcModel := vpc.ResponseToVpcValue(data.Vpc)
	vpcObjectValue, d := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Vpc = vpcObjectValue
	plan.Description = vpcModel.Description

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcVpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.Name = types.StringValue(data.Vpc.Name)
	state.Description = types.StringPointerValue(data.Vpc.Description.Get())

	vpcModel := vpc.ResponseToVpcValue(data.Vpc)
	vpcObjectValue, d := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Vpc = vpcObjectValue
	state.Description = vpcModel.Description

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the vpc and sets the updated Terraform state on success.
func (r *vpcVpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.VpcResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpc(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc",
			"Could not update vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcModel := vpc.ResponseToVpcValue(data.Vpc)
	vpcObjectValue, d := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(), vpcModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Vpc = vpcObjectValue
	state.Description = vpcModel.Description

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the vpc and removes the Terraform state on success.
func (r *vpcVpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpc
	err := r.client.DeleteVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc",
			"Could not delete vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func waitForVpcStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetVpc(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Vpc.State), nil
	}, -1, -1, -1, -1)
}

// ModifyPlan validates the plan and merges computed fields to prevent unnecessary (known after apply).
func (r *vpcVpcResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan vpc.VpcResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create
	if req.State.Raw.IsNull() {
		// Enforce cidr on create (API requires it, but schema marks it Optional for import)
		var planCidr types.String

		diags := req.Plan.GetAttribute(ctx, path.Root("cidr"), &planCidr)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if planCidr.IsNull() || planCidr.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("cidr"),
				"Missing Required Attribute on Creation",
				"The 'cidr' attribute is mandatory when creating a new VPC",
			)
		}
		return
	}

	var state vpc.VpcResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fields that cannot be updated via the API — check each for changes
	type fieldCheck struct {
		name    string
		changed bool
	}
	checks := []fieldCheck{
		{"name", !plan.Name.Equal(state.Name)},
	}

	for _, f := range checks {
		if f.changed {
			resp.Diagnostics.AddError(
				"Field changes not supported",
				fmt.Sprintf("Changing `%s` will not update the actual resource. To change %s, recreate the resource.", f.name, f.name),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var stateCidr, configCidr types.String
	req.State.GetAttribute(ctx, path.Root("cidr"), &stateCidr)
	req.Config.GetAttribute(ctx, path.Root("cidr"), &configCidr)

	if stateCidr.IsNull() || stateCidr.IsUnknown() {
		if !configCidr.IsNull() || !configCidr.IsUnknown() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("cidr"), configCidr)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	} else {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("cidr"), stateCidr)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Read fucntion does not return tags, and tags is only use when creating
	// So we omit tags when update
	var stateTags, configTags types.Map
	req.State.GetAttribute(ctx, path.Root("tags"), &stateTags)
	req.Config.GetAttribute(ctx, path.Root("tags"), &configTags)

	if stateTags.IsNull() || stateTags.IsUnknown() {
		if !configTags.IsNull() || !configTags.IsUnknown() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("tags"), configTags)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	} else {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("tags"), stateTags)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Reconstruct vpc: merge stable fields from state, leave rest as unknown.
	// This prevents id, account_id, created_at, created_by from showing as (known after apply).
	if !state.Vpc.IsNull() && !state.Vpc.IsUnknown() {
		var stateAg vpc.VpcDSValue
		resp.Diagnostics.Append(state.Vpc.As(ctx, &stateAg, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only mark modified_at/modified_by as unknown when description actually changes
		descriptionChanged := !plan.Description.Equal(state.Description)

		mergedAg := vpc.VpcDSValue{
			AccountId:   stateAg.AccountId,
			CidrCount:   stateAg.CidrCount,
			Cidrs:       stateAg.Cidrs,
			CreatedAt:   stateAg.CreatedAt,
			CreatedBy:   stateAg.CreatedBy,
			Description: plan.Description,
			Id:          stateAg.Id,
			ModifiedAt:  stateAg.ModifiedAt,
			ModifiedBy:  stateAg.ModifiedBy,
			Name:        stateAg.Name,
			State:       stateAg.State,
			ZoneType:    stateAg.ZoneType,
			Zones:       stateAg.Zones,
		}

		if descriptionChanged {
			mergedAg.ModifiedAt = types.StringUnknown()
			mergedAg.ModifiedBy = types.StringUnknown()
		} else {
			mergedAg.ModifiedAt = stateAg.ModifiedAt
			mergedAg.ModifiedBy = stateAg.ModifiedBy
		}

		mergedObj, diags := types.ObjectValueFrom(ctx, mergedAg.AttributeTypes(), mergedAg)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Vpc = mergedObj
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
