package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventNotificationStateDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventNotificationStateDataSources{}
)

func NewCloudMonitoringEventNotificationStateDataSources() datasource.DataSource {
	return &cloudMonitoringEventNotificationStateDataSources{}
}

type cloudMonitoringEventNotificationStateDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventNotificationStateDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_notification_states" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventNotificationStateDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Event Notification States.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Required:    true,
			},
			common.ToSnakeCase("EventId"): schema.StringAttribute{
				Description: "EventId",
				Required:    true,
			},
			common.ToSnakeCase("EventNotificationStates"): schema.ListNestedAttribute{
				Description: "A list of event notification state list.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("UserEmail"): schema.StringAttribute{
							Description: "UserEmail",
							Computed:    true,
						},
						common.ToSnakeCase("UserId"): schema.StringAttribute{
							Description: "UserId",
							Computed:    true,
						},
						common.ToSnakeCase("UserMobileTelNo"): schema.StringAttribute{
							Description: "UserMobileTelNo",
							Computed:    true,
						},
						common.ToSnakeCase("UserNameNotification"): schema.StringAttribute{
							Description: "UserNameNotification",
							Computed:    true,
						},
						common.ToSnakeCase("User"): schema.SingleNestedAttribute{
							Description: "User",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Email"): schema.StringAttribute{
									Description: "Email",
									Optional:    true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Id",
									Optional:    true,
								},
								common.ToSnakeCase("UserName"): schema.StringAttribute{
									Description: "UserName",
									Optional:    true,
								},
							},
						},
						common.ToSnakeCase("NotificationStatus"): schema.ListNestedAttribute{
							Description: "A list of Notification state list.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("ErrorMsg"): schema.StringAttribute{
										Description: "ErrorMsg",
										Computed:    true,
									},
									common.ToSnakeCase("SendDt"): schema.StringAttribute{
										Description: "SendDt",
										Computed:    true,
									},
									common.ToSnakeCase("SendMethod"): schema.StringAttribute{
										Description: "SendMethod",
										Computed:    true,
									},
									common.ToSnakeCase("SendResult"): schema.StringAttribute{
										Description: "SendResult",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

func (d *cloudMonitoringEventNotificationStateDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudMonitoring
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *cloudMonitoringEventNotificationStateDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventNotificationStateDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEventNotificationStateList(state.ResourceType, state.EventId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, eventNotificationStateElement := range data.Contents {

		//TODO	Object User 로 되어 있는 부분에 대한 수정이 필요
		//var userState = eventNotificationStateElement.User
		userState := cloudmonitoring.User{
			Email:    types.StringValue(eventNotificationStateElement.GetUserEmail()),
			Id:       types.StringValue(eventNotificationStateElement.GetUserId()),
			UserName: types.StringValue(eventNotificationStateElement.GetUserNameNotification()),
		}

		var notificationStates []cloudmonitoring.NotificationStatus

		for _, notificationStateElement := range eventNotificationStateElement.NotificationStates {
			notificationState := cloudmonitoring.NotificationStatus{
				ErrorMsg:   types.StringValue(notificationStateElement.GetErrorMsg()),
				SendDt:     types.StringValue(notificationStateElement.GetSendDt().String()),
				SendMethod: types.StringValue(notificationStateElement.GetSendMethod()),
				SendResult: types.StringValue(notificationStateElement.GetSendResult()),
			}

			notificationStates = append(notificationStates, notificationState)

		}

		eventNotificationState := cloudmonitoring.EventNotificationState{
			NotificationStatus:   notificationStates,
			User:                 userState,
			UserEmail:            types.StringValue(eventNotificationStateElement.GetUserEmail()),
			UserId:               types.StringValue(eventNotificationStateElement.UserId),
			UserMobileTelNo:      types.StringValue(eventNotificationStateElement.GetUserMobileTelNo()),
			UserNameNotification: types.StringValue(eventNotificationStateElement.UserNameNotification),
		}

		state.EventNotificationStates = append(state.EventNotificationStates, eventNotificationState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
