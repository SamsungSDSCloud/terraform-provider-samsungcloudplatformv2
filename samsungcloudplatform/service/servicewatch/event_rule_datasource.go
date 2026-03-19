package servicewatch

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchEventRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceWatchEventRuleDataSource{}
)

// serviceWatchEventRuleResource is a helper function to simplify the provider implementation.
func NewServiceWatchEventRuleDataSource() datasource.DataSource {
	return &serviceWatchEventRuleDataSource{}
}

// serviceWatchEventRuleDataSources is the data source implementation.
type serviceWatchEventRuleDataSource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchEventRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_event_rule"
}

func (d *serviceWatchEventRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
	    Description: "Event Rule Data Source",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("EventRule"): schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Event rule",
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID",
					},
					common.ToSnakeCase("ActiveYn"): schema.StringAttribute{
						Computed:            true,
						Description:         "Whether the Event rule is active",
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Computed:            true,
						Description:         "Created date time",
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Computed:            true,
						Description:         "Creator ID",
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Computed:            true,
						Description:         "Event rule description",
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Computed:            true,
						Description:         "Event rule ID",
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Computed:            true,
						Description:         "Modified date time",
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier ID",
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Computed:            true,
						Description:         "Event rule name",
					},
					common.ToSnakeCase("ResourceTypeId"): schema.StringAttribute{
						Computed:            true,
						Description:         "Resource type ID",
					},
					common.ToSnakeCase("ServiceId"): schema.StringAttribute{
						Computed:            true,
						Description:         "Service ID",
					},
				},
			},
			common.ToSnakeCase("EventRuleId"): schema.StringAttribute{
				Required:            true,
				Description:         "Event rule ID",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchEventRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.ServiceWatch
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *serviceWatchEventRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.EventRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEventRule(ctx, state.EventRuleId.ValueString())
	if err != nil && data == nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Event Rule",
			"Could not read Event Rule ID "+state.EventRuleId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	// 존재하지 않는 Event Rule 조회 시 null 로 return
	if data.EventRule.GetId() == "" {
		state.EventRule = types.ObjectNull(state.EventRule.AttributeTypes(ctx))
		resp.State.Set(ctx, &state)
		return
	}

	eventRuleResp := data.EventRule
	eventRule := servicewatch.EventRule{
    	AccountId:           types.StringValue(eventRuleResp.AccountId),
	    ActiveYn:            types.StringValue(string(eventRuleResp.ActiveYn)),
	    CreatedAt:           types.StringValue(eventRuleResp.CreatedAt.Format("2006-01-02 15:04:05")),
		CreatedBy:           types.StringValue(eventRuleResp.CreatedBy),
	    Description:         types.StringPointerValue(eventRuleResp.Description.Get()),
		Id:                  types.StringValue(eventRuleResp.Id),
		ModifiedAt:          types.StringValue(eventRuleResp.ModifiedAt.Format("2006-01-02 15:04:05")),
		ModifiedBy:          types.StringValue(eventRuleResp.ModifiedBy),
		Name:                types.StringValue(eventRuleResp.Name),
		ResourceTypeId:      types.StringPointerValue(eventRuleResp.ResourceTypeId.Get()),
		ServiceId:           types.StringValue(eventRuleResp.ServiceId),
	}
	eventRuleObjectValue, diags := types.ObjectValueFrom(ctx, eventRule.AttributeTypes(), eventRule)
    state.EventRule = eventRuleObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "ObjectValueFrom failed", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}
}
