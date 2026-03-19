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
	_ datasource.DataSource              = &serviceWatchLogStreamDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceWatchLogStreamDataSource{}
)

// serviceWatchLogStreamDataSource is a helper function to simplify the provider implementation.
func NewServiceWatchLogStreamDataSource() datasource.DataSource {
	return &serviceWatchLogStreamDataSource{}
}

// serviceWatchLogStreamDataSource is the data source implementation.
type serviceWatchLogStreamDataSource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchLogStreamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_log_stream"
}

// Schema defines the schema for the data source.
func (d *serviceWatchLogStreamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Stream Data Source",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("LogGroupId"): schema.StringAttribute{
				Description: "Log group ID",
				Required:    true,
			},
			common.ToSnakeCase("LogStreamId"): schema.StringAttribute{
				Description: "Log stream ID",
				Required:    true,
			},
			common.ToSnakeCase("LogStream"): schema.SingleNestedAttribute{
				Description: "Log stream",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Log stream ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Log stream name",
						Computed:    true,
					},
					common.ToSnakeCase("LogGroupId"): schema.StringAttribute{
						Description: "Log group ID",
						Computed:    true,
					},
					common.ToSnakeCase("CollectYn"): schema.StringAttribute{
						Description: "Whether to collect logs or not",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created date time",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Creator ID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified date time",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modifier ID",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchLogStreamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceWatchLogStreamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.LogStreamDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLogStream(ctx, state.LogGroupId.ValueString(), state.LogStreamId.ValueString())
	if err != nil {
		var detailMessage string
		if data != nil { // 존재하지 않는 log stream 조회 시 Error message 수정
			detailMessage = "Could not read Log Group ID " + state.LogGroupId.ValueString() + ", Log Stream ID " + state.LogStreamId.ValueString() +
				": 404 Not Found\nReason: No Log Stream found with ID " + state.LogStreamId.ValueString()
		} else { // 일반적인 Error 처리
			detail := client.GetDetailFromError(err)
			detailMessage = "Could not read Log Group ID " + state.LogGroupId.ValueString() +
				", Log Stream Id " + state.LogStreamId.ValueString() + ": " + err.Error() + "\nReason: " + detail
		}

		resp.Diagnostics.AddError("Error Reading Log Stream", detailMessage)
		return
	}

	logStreamResp := data.LogStream
	logStream := servicewatch.LogStream{
		Id:         types.StringValue(logStreamResp.Id),
		Name:       types.StringValue(logStreamResp.Name),
		LogGroupId: types.StringValue(logStreamResp.LogGroupId),
		CollectYn:  types.StringValue(string(logStreamResp.CollectYn)),
		CreatedAt:  types.StringValue(logStreamResp.CreatedAt.Format("2006-01-02 15:04:05")),
		CreatedBy:  types.StringValue(logStreamResp.CreatedBy),
		ModifiedAt: types.StringValue(logStreamResp.ModifiedAt.Format("2006-01-02 15:04:05")),
		ModifiedBy: types.StringValue(logStreamResp.ModifiedBy),
	}
	logStreamObjectValue, diags := types.ObjectValueFrom(ctx, logStream.AttributeTypes(), logStream)
	state.LogStream = logStreamObjectValue

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
