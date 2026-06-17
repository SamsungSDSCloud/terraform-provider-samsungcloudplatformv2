package loggingaudit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loggingauditTrailResource{}
	_ resource.ResourceWithConfigure   = &loggingauditTrailResource{}
	_ resource.ResourceWithImportState = &loggingauditTrailResource{}
)

// NewLoggingauditLogResource is a helper function to simplify the provider implementation.
func NewLoggingauditTrailResource() resource.Resource {
	return &loggingauditTrailResource{}
}

// loggingauditLogResource is the data source implementation.
type loggingauditTrailResource struct {
	config  *scpsdk.Configuration
	client  *loggingaudit.Client
	clients *client.SCPClient
}

func (r *loggingauditTrailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID 기반 단일 리소스 import 시 표준 패턴
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata returns the data source type name.
func (r *loggingauditTrailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loggingaudit_trail"
}

func (r *loggingauditTrailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Trail",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the trail.",
				MarkdownDescription: "The unique identifier of the trail.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description:         "Timestamp of the last Terraform update of the access key.",
				MarkdownDescription: "Timestamp of the last Terraform update of the access key.\n\nExample: `2023-10-27T10:00:00Z`",
				Computed:            true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description:         "Identifier of the account.",
				MarkdownDescription: "The unique identifier of the account.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Required:            true,
			},
			common.ToSnakeCase("BucketName"): schema.StringAttribute{
				Description:         "Name of the s3 bucket.",
				MarkdownDescription: "The name of the s3 bucket.\n\nExample: `example-bucket`",
				Optional:            true,
			},
			common.ToSnakeCase("BucketRegion"): schema.StringAttribute{
				Description:         "Region of the s3 bucket.",
				MarkdownDescription: "The region where the S3 bucket is located.\n\nExample: `kr-west1`",
				Optional:            true,
			},
			common.ToSnakeCase("LogTypeTotalYn"): schema.StringAttribute{
				Description:         "Whether to collect logs of all resource types.",
				MarkdownDescription: "Whether to collect logs of all resource types.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("LogVerificationYn"): schema.StringAttribute{
				Description:         "Whether to validate collected logs.",
				MarkdownDescription: "Whether to validate collected logs.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("RegionNames"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "List of region.",
				MarkdownDescription: "List of region to be collected via trail.\n\nExample: [`kr-west1`, `kr-east1`]",
				Optional:            true,
			},
			common.ToSnakeCase("RegionTotalYn"): schema.StringAttribute{
				Description:         "Whether to collect logs of all region.",
				MarkdownDescription: "Whether to collect logs of all region.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("ResourceTypeTotalYn"): schema.StringAttribute{
				Description:         "Whether to collect logs of all resource types.",
				MarkdownDescription: "Whether to collect logs of all resource types.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("TagCreateRequests"): schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":   types.StringType,
						"value": types.StringType,
					},
				},
				Description:         "Sets of tag.",
				MarkdownDescription: "Sets of tag.\n\nExample: [{\"key1\", \"value1\"}, {\"key2\", \"value2\"}]",
				Optional:            true,
			},
			common.ToSnakeCase("TargetLogTypes"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Type of the log.",
				MarkdownDescription: "The type of the collected log.\n\nExample: `AUDIT` | `EVENT`",
				Optional:            true,
			},
			common.ToSnakeCase("TargetResourceTypes"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "List of resource types.",
				MarkdownDescription: "The target to collect logs of resource types.\n\nExample: [`virtual-server`, `vpc`]",
				Optional:            true,
			},
			common.ToSnakeCase("TargetUsers"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "List of user ID.",
				MarkdownDescription: "The target to collect logs of user id.\n\nExample: [`e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`, `e4b2c3f8a1f94b6b9f1e8c2dca4f5b67`]",
				Optional:            true,
			},
			common.ToSnakeCase("TrailDescription"): schema.StringAttribute{
				Description:         "Description of the trail.",
				MarkdownDescription: "The description of the trail.\n\nExample: `a description of the trail`",
				Required:            true,
			},
			common.ToSnakeCase("TrailName"): schema.StringAttribute{
				Description:         "The name of the trail.",
				MarkdownDescription: "The name of the trail.\n\nExample: `example-trail`",
				Optional:            true,
			},
			common.ToSnakeCase("TrailSaveType"): schema.StringAttribute{
				Description:         "The save type of the trail.",
				MarkdownDescription: "The save type of the trail.\n\nExample: `JSON` | `CSV`",
				Optional:            true,
			},
			common.ToSnakeCase("UserTotalYn"): schema.StringAttribute{
				Description:         "Whether to set the trail collection target to all users.",
				MarkdownDescription: "Whether to set the trail collection target to all users.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("OrganizationTrailYn"): schema.StringAttribute{
				Description:         "Whether Trail collection is performed for sub-accounts of the organization account.",
				MarkdownDescription: "Whether Trail collection is performed for sub-accounts of the organization account.\n\nExample: `Y` | `N`",
				Required:            true,
			},
			common.ToSnakeCase("LogArchiveAccountId"): schema.StringAttribute{
				Description:         "Organization's administrative account ID",
				MarkdownDescription: "The Organization's administrative account ID.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Required:            true,
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
			common.ToSnakeCase("Trail"): schema.SingleNestedAttribute{
				Description: "A list of log.",
				Computed:    true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (r *loggingauditTrailResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoggingAudit
	r.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (r *loggingauditTrailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loggingaudit.TrailResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetTrail(ctx, state.Id.ValueString())

	if err != nil {
		// 404 Not Found 감지 시 상태 제거 (SDK 에러 처리 방식에 따라 조건 수정 필요)
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Trail",
			"Could not read Trail ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail

	trailModel := loggingaudit.Trail{
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

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)
	state.Trail = trailObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func TimePointValue(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func (r *loggingauditTrailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loggingaudit.TrailResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateTrail(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Trail",
			"Could not create Trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail
	plan.Id = types.StringValue(trail.Id)

	// ID가 유효한 경우에만 폴링
	if !plan.Id.IsNull() && plan.Id.ValueString() != "" {
		// 생성 후 리소스가 준비될 때까지 대기
		err = r.waitForTrailReady(ctx, plan.Id.ValueString(), 60*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for trail",
				"Trail was created but failed to become ready: "+err.Error(),
			)
			return
		}

		// 최신 상태 다시 조회
		// 1. 변수명을 showData로 변경하여 타입 충돌 방지
		showData, err := r.client.GetTrail(ctx, plan.Id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error reading trail after creation",
				// 2. plan.Id를 문자열로 변환 (ValueString 사용)
				"Could not read Trail ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		// 3. 변수명 동기화
		trail = showData.Trail
	}

	trailModel := loggingaudit.Trail{
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

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)

	plan.Trail = trailObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loggingauditTrailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state loggingaudit.TrailResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SetTrail(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Trail",
			"Could not update Trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetTrail(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading trail",
			"Could not read trail ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail

	trailModel := loggingaudit.Trail{
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

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)
	state.Trail = trailObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loggingauditTrailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loggingaudit.TrailResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing iam
	err := r.client.DeleteTrailKey(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting trail",
			"Could not delete trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func ConvertInterfaceListToStringList(rawList []interface{}) []types.String {
	result := make([]types.String, 0, len(rawList))
	for _, v := range rawList {
		if v == nil {
			result = append(result, types.StringNull())
		} else {
			strVal, _ := v.(string)
			result = append(result, types.StringValue(strVal))
		}
	}
	return result
}

func (r *loggingauditTrailResource) waitForTrailReady(ctx context.Context, id string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timeout waiting for trail %s to be ready", id)
		case <-ticker.C:
			data, err := r.client.GetTrail(ctx, id)
			if err != nil {
				// 404 Not Found - 생성 중일 수 있음, 계속 대기
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
					continue
				}
				// 일시적 오류는 무시하고 계속 폴링 (네트워크 문제 등)
				if client.IsTransientError(err) {
					continue
				}
				// 그 외의 영구적인 오류는 반환
				return err
			}

			// 상태 확인 - 실제 API 응답에 맞게 수정 필요
			trail := data.Trail
			if trail.Id != "" {
				return nil
			}
		}
	}
}
