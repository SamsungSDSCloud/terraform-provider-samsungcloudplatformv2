package smp

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &smpNotificationGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &smpNotificationGroupDataSources{}
)

func NewSmpNotificationGroupDataSources() datasource.DataSource {
	return &smpNotificationGroupDataSources{}
}

type smpNotificationGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (d *smpNotificationGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_notification_groups"
}

func (d *smpNotificationGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP NotificationGroups Data Source",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter by name.",
				Optional:    true,
			},
			"name_like": schema.StringAttribute{
				Description: "Wildcard search for notification group names.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "Filter by notification group ID.",
				Optional:    true,
			},
			"id_like": schema.StringAttribute{
				Description: "Wildcard search for notification group IDs.",
				Optional:    true,
			},
			"page": schema.Int64Attribute{
				Description: "Page number for pagination (default: 1).",
				Optional:    true,
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Number of results per page.",
				Optional:    true,
				Computed:    true,
			},
			"sort_by": schema.StringAttribute{
				Description: "Field to sort by (e.g. name, created_at).",
				Optional:    true,
			},
			"sort_order": schema.StringAttribute{
				Description: "Sort order: asc or desc.",
				Optional:    true,
			},
			"total_count": schema.Int64Attribute{
				Description: "Total number of results.",
				Computed:    true,
			},
			"sort": schema.ListNestedAttribute{
				Description: "Sort order returned by the API.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{Computed: true},
					},
				},
			},
			"notification_groups": schema.ListNestedAttribute{
				Description: "List of notification groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"created_at":  schema.StringAttribute{Computed: true},
						"modified_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *smpNotificationGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(ErrUnexpectedConfigure, fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData))
		return
	}
	d.client = inst.Client.Smp
	d.clients = inst.Client
}

func (d *smpNotificationGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smp.NotificationGroupDataSources
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNotificationGroupList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read NotificationGroups", err.Error())
		return
	}

	var ngList []types.Object
	var sortList []types.Object
	if data != nil {
		state.TotalCount = types.Int64Value(int64(data.GetCount()))
		state.Page = types.Int64Value(int64(data.GetPage()))
		state.Size = types.Int64Value(int64(data.GetSize()))

		if sort, ok := data.GetSortOk(); ok && sort != nil {
			for _, s := range sort {
				obj, d := types.ObjectValue(
					map[string]attr.Type{"value": types.StringType},
					map[string]attr.Value{"value": types.StringValue(s)},
				)
				resp.Diagnostics.Append(d...)
				if resp.Diagnostics.HasError() {
					return
				}
				sortList = append(sortList, obj)
			}
		}
		if sortList == nil {
			sortList = []types.Object{}
		}
		sortVal, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.StringType}}, sortList)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Sort = sortVal

		for _, ng := range data.GetNotificationGroups() {
			obj, d := types.ObjectValue(
				map[string]attr.Type{
					"id":          types.StringType,
					"name":        types.StringType,
					"description": types.StringType,
					"created_at":  types.StringType,
					"modified_at": types.StringType,
				},
				map[string]attr.Value{
					"id":          types.StringValue(ng.GetId()),
					"name":        types.StringValue(ng.GetName()),
					"description": nullableStringTypes(ng.GetDescriptionOk()),
					"created_at":  types.StringValue(ng.GetCreatedAt().Format(TimeFormatDisplay)),
					"modified_at": nullableTimeTypes(ng.GetModifiedAtOk()),
				},
			)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			ngList = append(ngList, obj)
		}
	}

	if ngList == nil {
		ngList = []types.Object{}
	}
	state.NotificationGroups, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"created_at":  types.StringType,
		"modified_at": types.StringType,
	}}, ngList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
