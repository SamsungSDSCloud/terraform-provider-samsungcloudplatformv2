package loggingaudit

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loggingaudit" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("TrailName"): schema.StringAttribute{
				Description: "TrailName",
				Optional:    true,
			},
			common.ToSnakeCase("BucketName"): schema.StringAttribute{
				Description: "BucketName",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "ResourceName",
				Optional:    true,
			},
			common.ToSnakeCase("Trail"): schema.ListNestedAttribute{
				Description: "A list of log.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("AccountName"): schema.StringAttribute{
							Description: "AccountName",
							Computed:    true,
						},
						common.ToSnakeCase("BucketName"): schema.StringAttribute{
							Description: "BucketName",
							Computed:    true,
						},
						common.ToSnakeCase("BucketRegion"): schema.StringAttribute{
							Description: "BucketRegion",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "CreatedAt",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedUserId"): schema.StringAttribute{
							Description: "CreatedUserId",
							Computed:    true,
						},
						common.ToSnakeCase("DelYn"): schema.StringAttribute{
							Description: "DelYn",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("LogTypeTotalYn"): schema.StringAttribute{
							Description: "LogTypeTotalYn",
							Computed:    true,
						},
						common.ToSnakeCase("LogVerificationYn"): schema.StringAttribute{
							Description: "LogVerificationYn",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("RegionNames"): schema.ListAttribute{
							ElementType: types.StringType,
							Description: "RegionNames",
							Computed:    true,
						},
						common.ToSnakeCase("RegionTotalYn"): schema.StringAttribute{
							Description: "RegionTotalYn",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceTypeTotalYn"): schema.StringAttribute{
							Description: "ResourceTypeTotalYn",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
						},
						common.ToSnakeCase("TargetLogTypes"): schema.ListAttribute{
							ElementType: types.StringType,
							Description: "TargetLogTypes",
							Computed:    true,
						},
						common.ToSnakeCase("TargetResourceTypes"): schema.ListAttribute{
							ElementType: types.StringType,
							Description: "TargetResourceTypes",
							Computed:    true,
						},
						common.ToSnakeCase("TargetUsers"): schema.ListAttribute{
							ElementType: types.StringType,
							Description: "TargetUsers",
							Computed:    true,
						},
						common.ToSnakeCase("TrailBatchEndAt"): schema.StringAttribute{
							Description: "TrailBatchEndAt",
							Computed:    true,
						},
						common.ToSnakeCase("TrailBatchFirstStartAt"): schema.StringAttribute{
							Description: "TrailBatchFirstStartAt",
							Computed:    true,
						},
						common.ToSnakeCase("TrailBatchLastState"): schema.StringAttribute{
							Description: "TrailBatchLastState",
							Computed:    true,
						},
						common.ToSnakeCase("TrailBatchStartAt"): schema.StringAttribute{
							Description: "TrailBatchStartAt",
							Computed:    true,
						},
						common.ToSnakeCase("TrailBatchSuccessAt"): schema.StringAttribute{
							Description: "TrailBatchSuccessAt",
							Computed:    true,
						},
						common.ToSnakeCase("TrailDescription"): schema.StringAttribute{
							Description: "TrailDescription",
							Computed:    true,
						},
						common.ToSnakeCase("TrailName"): schema.StringAttribute{
							Description: "TrailName",
							Computed:    true,
						},
						common.ToSnakeCase("TrailSaveType"): schema.StringAttribute{
							Description: "TrailSaveType",
							Computed:    true,
						},
						common.ToSnakeCase("UserTotalYn"): schema.StringAttribute{
							Description: "UserTotalYn",
							Computed:    true,
						},
						common.ToSnakeCase("OrganizationTrailYn"): schema.StringAttribute{
							Description: "OrganizationTrailYn",
							Computed:    true,
						},
						common.ToSnakeCase("LogArchiveAccountId"): schema.StringAttribute{
							Description: "LogArchiveAccountId",
							Computed:    true,
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
