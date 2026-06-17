package loggingaudit

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loggingaudit" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loggingauditTrailDataSource{}
	_ datasource.DataSourceWithConfigure = &loggingauditTrailDataSource{}
)

func NewLoggingauditTrailDataSource() datasource.DataSource {
	return &loggingauditTrailDataSource{}
}

// resourceManagerResourceGroupDataSource is the data source implementation.
type loggingauditTrailDataSource struct {
	config  *scpsdk.Configuration
	client  *loggingaudit.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loggingauditTrailDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loggingaudit_trail" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loggingauditTrailDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of trail.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description:         "A Number of results displayed per page.",
				MarkdownDescription: "A Number of results displayed per page.\n\nExample: `15`",
				Optional:            true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description:         "A Number of page.",
				MarkdownDescription: "A Number of page.\n\nExample: `1`",
				Optional:            true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description:         "Sorts the query results.",
				MarkdownDescription: "Sorts the query results.\n\nExample: `createdAt:desc`",
				Optional:            true,
			},
			common.ToSnakeCase("TrailName"): schema.StringAttribute{
				Description:         "Name of the trail.",
				MarkdownDescription: "The name of the trail.\n\nExample: `example-trail`",
				Optional:            true,
			},
			common.ToSnakeCase("BucketName"): schema.StringAttribute{
				Description:         "Name of the s3 bucket.",
				MarkdownDescription: "The name of the s3 bucket.\n\nExample: `example-bucket`",
				Optional:            true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "State of the trail.",
				MarkdownDescription: "State of the trail.\n\nExample: `ACTIVE | STOPPED`",
				Optional:            true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description:         "List of resource types.",
				MarkdownDescription: "A list of resource types associated with the group.\n\nExample: `[\"virtual-server\"]`",
				Optional:            true,
			},
			common.ToSnakeCase("IamRoleId"): schema.StringAttribute{
				Description:         "This is the role id (IAM) to use when integrating with ServiceWatch.",
				MarkdownDescription: "This is the role id (IAM) to use when integrating with ServiceWatch.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Optional:            true,
			},
			common.ToSnakeCase("LogGroupName"): schema.StringAttribute{
				Description:         "This is the log group name to use when integrating with ServiceWatch.",
				MarkdownDescription: "This is the log group name to use when integrating with ServiceWatch.\n\nExample: test-log-group",
				Optional:            true,
			},
			common.ToSnakeCase("ServiceWatchYn"): schema.StringAttribute{
				Description:         "Whether to integrate that trail with ServiceWatch.",
				MarkdownDescription: "Whether to integrate that trail with ServiceWatch.\n\nExample: `Y` | `N`",
				Optional:            true,
			},
			common.ToSnakeCase("Trail"): schema.ListNestedAttribute{
				Description: "A list of log.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description:         "Identifier of the account.",
							MarkdownDescription: "The unique identifier of the account.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("AccountName"): schema.StringAttribute{
							Description:         "Name of the account.",
							MarkdownDescription: "The name of the account.\n\nExample: `example-account`",
							Computed:            true,
						},
						common.ToSnakeCase("BucketName"): schema.StringAttribute{
							Description:         "Name of the s3 bucket.",
							MarkdownDescription: "The name of the s3 bucket.\n\nExample: `example-bucket`",
							Computed:            true,
						},
						common.ToSnakeCase("BucketRegion"): schema.StringAttribute{
							Description:         "Region of the s3 bucket.",
							MarkdownDescription: "The region where the S3 bucket is located.\n\nExample: `kr-west1`",
							Computed:            true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description:         "Creation timestamp.",
							MarkdownDescription: "The creation timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description:         "Creator identifier.",
							MarkdownDescription: "The user ID of the creator.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("CreatedUserId"): schema.StringAttribute{
							Description:         "Creator email.",
							MarkdownDescription: "The user's email of the creator.\n\nExample: `test@samsung.com`",
							Computed:            true,
						},
						common.ToSnakeCase("DelYn"): schema.StringAttribute{
							Description:         "Whether to delete Trail.",
							MarkdownDescription: "Whether to delete Trail.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Identifier of the trail.",
							MarkdownDescription: "The unique identifier of the trail.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("LogTypeTotalYn"): schema.StringAttribute{
							Description:         "Whether to collect logs of all resource types.",
							MarkdownDescription: "Whether to collect logs of all resource types.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("LogVerificationYn"): schema.StringAttribute{
							Description:         "Whether to validate collected logs.",
							MarkdownDescription: "Whether to validate collected logs.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description:         "Modification timestamp.",
							MarkdownDescription: "The modification timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description:         "Modifier identifier.",
							MarkdownDescription: "The user ID of the modifier.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("RegionNames"): schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "List of region.",
							MarkdownDescription: "List of region to be collected via trail.\n\nExample: [`kr-west1`, `kr-east1`]",
							Computed:            true,
						},
						common.ToSnakeCase("RegionTotalYn"): schema.StringAttribute{
							Description:         "Whether to collect logs of all region.",
							MarkdownDescription: "Whether to collect logs of all region.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("ResourceTypeTotalYn"): schema.StringAttribute{
							Description:         "Whether to collect logs of all resource types.",
							MarkdownDescription: "Whether to collect logs of all resource types.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description:         "State of the trail.",
							MarkdownDescription: "State of the trail.\n\nExample: `ACTIVE | STOPPED`",
							Computed:            true,
						},
						common.ToSnakeCase("TargetLogTypes"): schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "Type of the log.",
							MarkdownDescription: "The type of the collected log.\n\nExample: `AUDIT` | `EVENT`",
							Computed:            true,
						},
						common.ToSnakeCase("TargetResourceTypes"): schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "List of resource types.",
							MarkdownDescription: "The target to collect logs of resource types.\n\nExample: [`virtual-server`, `vpc`]",
							Computed:            true,
						},
						common.ToSnakeCase("TargetUsers"): schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "List of user ID.",
							MarkdownDescription: "The target to collect logs of user id.\n\nExample: [`e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`, `e4b2c3f8a1f94b6b9f1e8c2dca4f5b67`]",
							Computed:            true,
						},
						common.ToSnakeCase("TrailBatchEndAt"): schema.StringAttribute{
							Description:         "End date and time of the 'trail' collection scheduler.",
							MarkdownDescription: "End date and time of the 'trail' collection scheduler.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailBatchFirstStartAt"): schema.StringAttribute{
							Description:         "Start date at first and time of the 'trail' collection scheduler.",
							MarkdownDescription: "Start date at first and time of the 'trail' collection scheduler.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailBatchLastState"): schema.StringAttribute{
							Description:         "Results of 'Trail' data collection schedule execution.",
							MarkdownDescription: "The results of 'Trail' data collection schedule execution.\n\nExample: `SUCCESS` | `FAILED`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailBatchStartAt"): schema.StringAttribute{
							Description:         "Start date and time of the 'trail' collection scheduler.",
							MarkdownDescription: "Start date and time of the 'trail' collection scheduler.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailBatchSuccessAt"): schema.StringAttribute{
							Description:         "Succeed date and time of the 'trail' collection scheduler.",
							MarkdownDescription: "Succeed date and time of the 'trail' collection scheduler.\n\nExample: `2023-10-27T10:00:00Z`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailDescription"): schema.StringAttribute{
							Description:         "Description of the trail.",
							MarkdownDescription: "The description of the trail.\n\nExample: `a description of the trail`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailName"): schema.StringAttribute{
							Description:         "The name of the trail.",
							MarkdownDescription: "The name of the trail.\n\nExample: `example-trail`",
							Computed:            true,
						},
						common.ToSnakeCase("TrailSaveType"): schema.StringAttribute{
							Description:         "The save type of the trail.",
							MarkdownDescription: "The save type of the trail.\n\nExample: `JSON` | `CSV`",
							Computed:            true,
						},
						common.ToSnakeCase("UserTotalYn"): schema.StringAttribute{
							Description:         "Whether to set the trail collection target to all users.",
							MarkdownDescription: "Whether to set the trail collection target to all users.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("OrganizationTrailYn"): schema.StringAttribute{
							Description:         "Whether Trail collection is performed for sub-accounts of the organization account.",
							MarkdownDescription: "Whether Trail collection is performed for sub-accounts of the organization account.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
						common.ToSnakeCase("LogArchiveAccountId"): schema.StringAttribute{
							Description:         "Organization's administrative account ID",
							MarkdownDescription: "The Organization's administrative account ID.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("IamRoleId"): schema.StringAttribute{
							Description:         "This is the role id (IAM) to use when integrating with ServiceWatch.",
							MarkdownDescription: "This is the role id (IAM) to use when integrating with ServiceWatch.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
							Computed:            true,
						},
						common.ToSnakeCase("LogGroupName"): schema.StringAttribute{
							Description:         "This is the log group name to use when integrating with ServiceWatch.",
							MarkdownDescription: "This is the log group name to use when integrating with ServiceWatch.\n\nExample: test-log-group",
							Computed:            true,
						},
						common.ToSnakeCase("ServiceWatchYn"): schema.StringAttribute{
							Description:         "Whether to integrate that trail with ServiceWatch.",
							MarkdownDescription: "Whether to integrate that trail with ServiceWatch.\n\nExample: `Y` | `N`",
							Computed:            true,
						},
					},
				},
			},
		},
	}

}

// Configure adds the provider configured client to the data source.
func (d *loggingauditTrailDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.LoggingAudit
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *loggingauditTrailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loggingaudit.TrailDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetTrailList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Access Keys",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, trail := range data.Trails {

		trailState := loggingaudit.Trail{
			AccountId:              types.StringPointerValue(trail.AccountId.Get()),
			AccountName:            types.StringPointerValue(trail.AccountName.Get()),
			BucketName:             types.StringValue(trail.BucketName),
			BucketRegion:           types.StringValue(trail.BucketName),
			CreatedAt:              types.StringValue(trail.CreatedAt.Format(time.RFC3339)),
			CreatedBy:              types.StringValue(trail.CreatedBy),
			CreatedUserId:          types.StringPointerValue(trail.CreatedUserId.Get()),
			DelYn:                  types.StringPointerValue(trail.DelYn.Get()),
			Id:                     types.StringValue(trail.Id),
			LogTypeTotalYn:         types.StringPointerValue(trail.LogTypeTotalYn.Get()),
			LogVerificationYn:      types.StringPointerValue(trail.LogVerificationYn.Get()),
			ModifiedAt:             types.StringValue(trail.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:             types.StringValue(trail.ModifiedBy),
			RegionNames:            ConvertInterfaceListToStringList(trail.RegionNames),
			RegionTotalYn:          types.StringPointerValue(trail.RegionTotalYn.Get()),
			ResourceTypeTotalYn:    types.StringPointerValue(trail.ResourceTypeTotalYn.Get()),
			State:                  types.StringPointerValue(trail.State.Get()),
			TargetLogTypes:         ConvertInterfaceListToStringList(trail.TargetLogTypes),
			TargetResourceTypes:    ConvertInterfaceListToStringList(trail.TargetResourceTypes),
			TargetUsers:            ConvertInterfaceListToStringList(trail.TargetUsers),
			TrailBatchEndAt:        TimePointValue(trail.TrailBatchEndAt.Get()),
			TrailBatchFirstStartAt: TimePointValue(trail.TrailBatchFirstStartAt.Get()),
			TrailBatchLastState:    types.StringPointerValue(trail.TrailBatchLastState.Get()),
			TrailBatchStartAt:      TimePointValue(trail.TrailBatchStartAt.Get()),
			TrailBatchSuccessAt:    TimePointValue(trail.TrailBatchSuccessAt.Get()),
			TrailDescription:       types.StringPointerValue(trail.TrailDescription.Get()),
			TrailName:              types.StringValue(trail.TrailName),
			TrailSaveType:          types.StringValue(trail.TrailSaveType),
			UserTotalYn:            types.StringPointerValue(trail.UserTotalYn.Get()),
			OrganizationTrailYn:    types.StringPointerValue(trail.OrganizationTrailYn.Get()),
			LogArchiveAccountId:    types.StringPointerValue(trail.LogArchiveAccountId.Get()),
			IamRoleId:              types.StringPointerValue(trail.IamRoleId.Get()),
			LogGroupName:           types.StringPointerValue(trail.LogGroupName.Get()),
			ServiceWatchYn:         types.StringPointerValue(trail.ServiceWatchYn.Get()),
		}

		state.Trails = append(state.Trails, trailState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
