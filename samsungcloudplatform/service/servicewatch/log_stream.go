package servicewatch

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &serviceWatchLogStreamResource{}
	_ resource.ResourceWithConfigure = &serviceWatchLogStreamResource{}
)

// NewServiceWatchLogStreamResource is a helper function to simplify the provider implementation.
func NewServiceWatchLogStreamResource() resource.Resource {
	return &serviceWatchLogStreamResource{}
}

// serviceWatchLogStreamResource is the data source implementation.
type serviceWatchLogStreamResource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *serviceWatchLogStreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_log_stream"
}

// Schema defines the schema for the resource.
func (r *serviceWatchLogStreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Stream Resource",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the log stream.\n" +
					" - example : b0fc99a46a2c44b8b22739b130853dd8\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Log stream Name.\n" +
					" - example : testls01\n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("LogGroupId"): schema.StringAttribute{
				Description: "The unique identifier of the log group.\n" +
					" - example : bce52822147744b4afe0187164caa2e8\n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("LogStream"): schema.SingleNestedAttribute{
				Description: "List of log stream.\n" +
					" - example : {\"id\": \"b0fc99a46a2c44b8b22739b130853dd8\", \"name\": \"testls01\"}\n",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Log stream ID.\n" +
							" - example : b0fc99a46a2c44b8b22739b130853dd8\n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Log stream name.\n" +
							" - example : testls01\n",
						Computed: true,
					},
					common.ToSnakeCase("LogGroupId"): schema.StringAttribute{
						Description: "The unique identifier of the log group.\n" +
							" - example : bce52822147744b4afe0187164caa2e8\n",
						Computed: true,
					},
					common.ToSnakeCase("CollectYn"): schema.StringAttribute{
						Description: "Whether to collect logs or not.\n" +
							" - example : Y\n",
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
func (r *serviceWatchLogStreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.ServiceWatch
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *serviceWatchLogStreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicewatch.LogStreamResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new LogStream
	data, err := r.client.CreateLogStream(ctx, plan.LogGroupId.ValueString(), plan.Name.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating LogStream",
			"Could not create LogStream, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	logStream := convertLogStream(&data.LogStream)
	logStreamObjectValue, diags := types.ObjectValueFrom(ctx, logStream.AttributeTypes(), logStream)

	plan.Id = types.StringValue(logStream.Id.ValueString())
	plan.LogStream = logStreamObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceWatchLogStreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicewatch.LogStreamResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Log Stream
	data, err := r.client.GetLogStream(ctx, state.LogGroupId.ValueString(), state.Id.ValueString())
	if err != nil {
		var detailMessage string
		if data != nil { // 존재하지 않는 log stream 조회 시 Error message 수정
			detailMessage = "Could not read Log Group ID " + state.LogGroupId.ValueString() + ", Log Stream ID " + state.Id.ValueString() +
				": 404 Not Found\nReason: No Log Stream found with ID " + state.Id.ValueString()
		} else { // 일반적인 Error 처리
			detail := client.GetDetailFromError(err)
			detailMessage = "Could not read Log Group ID " + state.LogGroupId.ValueString() +
				", Log Stream Id " + state.Id.ValueString() + ": " + err.Error() + "\nReason: " + detail
		}

		resp.Diagnostics.AddError("Error Reading Log Stream", detailMessage)
		return
	}

	// Map response body to schema and populate Computed attribute values
	logStream := convertLogStream(&data.LogStream)
	logStreamObjectValue, diags := types.ObjectValueFrom(ctx, logStream.AttributeTypes(), logStream)
	state.LogStream = logStreamObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceWatchLogStreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Config inspection does not support update operations
	// This is a no-op implementation
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"Log stream resources do not support \"UPDATE\" operations. The resource will not be updated.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceWatchLogStreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicewatch.LogStreamResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	logGroupId := state.LogGroupId.ValueString()
	logStreamId := state.Id.ValueString()

	// Delete existing Resource Group
	_, err := r.client.DeleteLogStream(ctx, logGroupId, []string{logStreamId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Log Stream",
			"Could not delete Log Stream, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertLogStream(logStreamResp *servicewatch2.LogStreamDTO) servicewatch.LogStream {
	return servicewatch.LogStream{
		Id:         types.StringValue(logStreamResp.Id),
		Name:       types.StringValue(logStreamResp.Name),
		LogGroupId: types.StringValue(logStreamResp.LogGroupId),
		CollectYn:  types.StringValue(string(logStreamResp.CollectYn)),
		CreatedAt:  types.StringValue(logStreamResp.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringValue(logStreamResp.CreatedBy),
		ModifiedAt: types.StringValue(logStreamResp.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringValue(logStreamResp.ModifiedBy),
	}
}
