package smp

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/smp/1.0"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &smpNotificationGroupResource{}
	_ resource.ResourceWithConfigure   = &smpNotificationGroupResource{}
	_ resource.ResourceWithImportState = &smpNotificationGroupResource{}
)

func NewSmpNotificationGroupResource() resource.Resource {
	return &smpNotificationGroupResource{}
}

type smpNotificationGroupResource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

// --------------------------------------------------------------------------
// METADATA
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_notification_group"
}

// --------------------------------------------------------------------------
// SCHEMA
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP NotificationGroup Resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Notification group name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Notification group description.",
				Optional:    true,
			},
			"recipient_user_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of recipient user IDs.",
				Optional:    true,
			},
			"notification_group": schema.SingleNestedAttribute{
				Description: "Notification group details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":          schema.StringAttribute{Computed: true},
					"name":        schema.StringAttribute{Computed: true},
					"description": schema.StringAttribute{Computed: true},
					"account_id":  schema.StringAttribute{Computed: true},
					"created_at":  schema.StringAttribute{Computed: true},
					"created_by":  schema.StringAttribute{Computed: true},
					"modified_at": schema.StringAttribute{Computed: true},
					"modified_by": schema.StringAttribute{Computed: true},
					"recipients": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id":                schema.StringAttribute{Computed: true},
								"recipient_user_id": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"notification_group_id": schema.StringAttribute{
				Description: "Notification group ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// --------------------------------------------------------------------------
// CONFIGURE
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(ErrUnexpectedConfigure,
			fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData))
		return
	}
	r.client = inst.Client.Smp
	r.clients = inst.Client
}

// --------------------------------------------------------------------------
// CREATE
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smp.NotificationGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateNotificationGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrCreateNotificationGroup,
			fmt.Sprintf(ErrCreateNotificationGroupFmt, err.Error(), detail))
		return
	}

	ng := result.GetNotificationGroup()
	ngId := ng.GetId()
	plan.NotificationGroupId = types.StringValue(ngId)

	// Re-read to get full detail (Create response lacks account_id/recipients)
	showResult, err := r.client.GetNotificationGroup(ctx, ngId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadNotificationGroup,
			fmt.Sprintf(ErrReadNotificationGroupFmt, ngId, err.Error(), detail))
		return
	}

	detail := showResult.GetNotificationGroup()
	ngObj, diags := buildNotificationGroupObject(ctx, detail)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.NotificationGroup = ngObj

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// READ
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smp.NotificationGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetNotificationGroup(ctx, state.NotificationGroupId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadNotificationGroup,
			fmt.Sprintf(ErrReadNotificationGroupFmt,
				state.NotificationGroupId.ValueString(),
				err.Error(),
				detail))
		return
	}

	ng := result.GetNotificationGroup()
	ngObj, diags := buildNotificationGroupObject(ctx, ng)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.NotificationGroupId = types.StringValue(ng.GetId())
	state.NotificationGroup = ngObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --------------------------------------------------------------------------
// UPDATE
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smp.NotificationGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateNotificationGroup(ctx, plan.NotificationGroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrUpdateNotificationGroup,
			fmt.Sprintf(ErrUpdateNotificationGroupFmt, err.Error(), detail))
		return
	}

	// Update recipients if provided
	if !plan.RecipientUserIds.IsNull() && !plan.RecipientUserIds.IsUnknown() {
		err = r.client.UpdateNotificationGroupRecipients(ctx, plan.NotificationGroupId.ValueString(), plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(ErrUpdateNotificationGroup,
				fmt.Sprintf(ErrUpdateNotificationGroupFmt, err.Error(), detail))
			return
		}
	}

	result, err := r.client.GetNotificationGroup(ctx, plan.NotificationGroupId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadNotificationGroup,
			fmt.Sprintf(ErrReadNotificationGroupFmt,
				plan.NotificationGroupId.ValueString(),
				err.Error(),
				detail))
		return
	}

	ng := result.GetNotificationGroup()
	ngObj, diags := buildNotificationGroupObject(ctx, ng)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.NotificationGroupId = types.StringValue(ng.GetId())
	plan.NotificationGroup = ngObj

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// DELETE
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smp.NotificationGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNotificationGroup(ctx, state.NotificationGroupId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrDeleteNotificationGroup,
			fmt.Sprintf(ErrDeleteNotificationGroupFmt, err.Error(), detail))
		return
	}
}

// --------------------------------------------------------------------------
// IMPORT
// --------------------------------------------------------------------------
func (r *smpNotificationGroupResource) ImportState(ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("notification_group_id"), req, resp)
}

// --------------------------------------------------------------------------
// HELPER - build NotificationGroup terraform object from SDK response
// --------------------------------------------------------------------------
func buildNotificationGroupObject(ctx context.Context, ng sdk.NotificationGroupDetailResponse) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build recipients list
	var recipientList []types.Object
	recipients := ng.GetRecipients()
	if recipients != nil && len(recipients) > 0 {
		for _, r := range recipients {
			obj, d := types.ObjectValue(
				map[string]attr.Type{
					"id":                types.StringType,
					"recipient_user_id": types.StringType,
				},
				map[string]attr.Value{
					"id":                types.StringValue(r.GetId()),
					"recipient_user_id": types.StringValue(r.GetRecipientUserId()),
				},
			)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(smp.NotificationGroup{}.AttributeTypes()), diags
			}
			recipientList = append(recipientList, obj)
		}
	} else {
		recipientList = []types.Object{}
	}

	recipientsList, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                types.StringType,
			"recipient_user_id": types.StringType,
		},
	}, recipientList)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(smp.NotificationGroup{}.AttributeTypes()), diags
	}

	ngObj, d := types.ObjectValue(
		smp.NotificationGroup{}.AttributeTypes(),
		map[string]attr.Value{
			"id":          types.StringValue(ng.GetId()),
			"name":        types.StringValue(ng.GetName()),
			"description": nullableStringTypes(ng.GetDescriptionOk()),
			"account_id":  types.StringValue(ng.GetAccountId()),
			"created_at":  types.StringValue(ng.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":  types.StringValue(ng.GetCreatedBy()),
			"modified_at": nullableTimeTypes(ng.GetModifiedAtOk()),
			"modified_by": nullableStringTypes(ng.GetModifiedByOk()),
			"recipients":  recipientsList,
		},
	)
	diags.Append(d...)
	return ngObj, diags
}
