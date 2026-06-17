package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchLogGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceWatchLogGroupDataSource{}
)

// serviceWatchLogGroupResource is a helper function to simplify the provider implementation.
func NewServiceWatchLogGroupDataSource() datasource.DataSource {
	return &serviceWatchLogGroupDataSource{}
}

// serviceWatchLogGroupDataSources is the data source implementation.
type serviceWatchLogGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchLogGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_log_group"
}

// Schema defines the schema for the data source.
func (d *serviceWatchLogGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Group Data Source",
		Attributes: map[string]schema.Attribute{
		common.ToSnakeCase("LogGroupId"): schema.StringAttribute{
			Description: "The unique identifier of the log group.\n" +
				" - example : bce52822147744b4afe0187164caa2e8\n",
			Required: true,
		},
			common.ToSnakeCase("LogGroup"): schema.SingleNestedAttribute{
				Description: "List of log group.\n" +
					" - example : {\"id\": \"bce52822147744b4afe0187164caa2e8\", \"name\": \"testlg01\"}\n",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
				common.ToSnakeCase("Id"): schema.StringAttribute{
					Description: "The unique identifier of the log group.\n" +
						" - example : bce52822147744b4afe0187164caa2e8\n",
					Computed: true,
				},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Log group name.\n" +
							" - example : testlg01\n",
						Computed: true,
					},
				common.ToSnakeCase("AccountId"): schema.StringAttribute{
					Description: "The unique identifier of the account.\n" +
						" - example : 1bcf39b344ac41cbaf0466ff0d2bebad\n",
					Computed: true,
				},
					common.ToSnakeCase("RetentionPeriod"): schema.Int32Attribute{
						Description: "Log group retention period.\n" +
							" - example : 365\n",
						Computed: true,
					},
					common.ToSnakeCase("RetentionPeriodName"): schema.StringAttribute{
						Description: "Log group retention period name.\n" +
							" - example : 1 year\n",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Log group status.\n" +
							"Allowed values: ACTIVE, DELETING, DELETED.\n" +
							" - example : ACTIVE\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							" - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							" - example : 2024-05-17T00:23:17Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchLogGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceWatchLogGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.LogGroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLogGroup(ctx, state.LogGroupId.ValueString())
	if err != nil {
		var detailMessage string
		if data != nil { // 존재하지 않는 log group 조회 시, Error message 수정
			detailMessage = "Could not read Log Group ID " + state.LogGroupId.ValueString() +
				": 404 Not Found\nReason: No Log Group found with ID " + state.LogGroupId.ValueString()
		} else { // 일반적인 Error 처리
			detail := client.GetDetailFromError(err)
			detailMessage = "Could not read Dashboard ID " + state.LogGroupId.ValueString() + ": " + err.Error() + "\nReason: " + detail
		}
		resp.Diagnostics.AddError("Error Reading Log Group", detailMessage)
		return
	}

	logGroupResp := data.LogGroup
	logGroup := servicewatch.LogGroup{
		Id:                  types.StringValue(logGroupResp.Id),
		Name:                types.StringValue(logGroupResp.Name),
		AccountId:           types.StringValue(logGroupResp.AccountId),
		RetentionPeriod:     types.Int32Value(logGroupResp.RetentionPeriod),
		RetentionPeriodName: types.StringValue(logGroupResp.RetentionPeriodName),
		Status:              types.StringValue(string(logGroupResp.Status)),
		CreatedAt:           types.StringValue(logGroupResp.CreatedAt.Format("2006-01-02 15:04:05")),
		CreatedBy:           types.StringValue(logGroupResp.CreatedBy),
		ModifiedAt:          types.StringValue(logGroupResp.ModifiedAt.Format("2006-01-02 15:04:05")),
		ModifiedBy:          types.StringValue(logGroupResp.ModifiedBy),
	}
	logGroupObjectValue, diags := types.ObjectValueFrom(ctx, logGroup.AttributeTypes(), logGroup)
	state.LogGroup = logGroupObjectValue

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
