package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventPolicyNotificationDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventPolicyNotificationDataSources{}
)

func NewCloudMonitoringEventPolicyNotificationDataSources() datasource.DataSource {
	return &cloudMonitoringEventPolicyNotificationDataSources{}
}

type cloudMonitoringEventPolicyNotificationDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventPolicyNotificationDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_policy_notifications" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventPolicyNotificationDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of event.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
				Description: "EventPolicyId",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicyNotifications"): schema.ListNestedAttribute{
				Description: "Event Policy NotificationList.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("RecipientType"): schema.StringAttribute{
							Description: "RecipientType",
							Computed:    true,
						},
						common.ToSnakeCase("RecipientKey"): schema.StringAttribute{
							Description: "RecipientKey",
							Computed:    true,
						},
						common.ToSnakeCase("RecipientName"): schema.StringAttribute{
							Description: "RecipientName",
							Computed:    true,
						},
						common.ToSnakeCase("UserAdditionalInfo"): schema.SingleNestedAttribute{
							Description: "UserAdditionalInfo",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("CompanyName"): schema.StringAttribute{
									Description: "CompanyName",
									Computed:    true,
								},
								common.ToSnakeCase("Email"): schema.StringAttribute{
									Description: "Email",
									Computed:    true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Id",
									Computed:    true,
								},
								common.ToSnakeCase("LastLoginAt"): schema.StringAttribute{
									Description: "LastLoginAt",
									Computed:    true,
								},
								common.ToSnakeCase("MemberCreatedAt"): schema.StringAttribute{
									Description: "MemberCreatedAt",
									Computed:    true,
								},
								common.ToSnakeCase("PasswordReuseCount"): schema.Int64Attribute{
									Description: "PasswordReuseCount",
									Computed:    true,
								},
								common.ToSnakeCase("UserName"): schema.StringAttribute{
									Description: "UserName",
									Computed:    true,
								},
								common.ToSnakeCase("TzId"): schema.StringAttribute{
									Description: "TzId",
									Computed:    true,
								},
								common.ToSnakeCase("Timezone"): schema.StringAttribute{
									Description: "Timezone",
									Computed:    true,
								},
								common.ToSnakeCase("DstOffset"): schema.StringAttribute{
									Description: "DstOffset",
									Computed:    true,
								},
								common.ToSnakeCase("UtcOffset"): schema.StringAttribute{
									Description: "UtcOffset",
									Computed:    true,
								},
								common.ToSnakeCase("AccountId"): schema.StringAttribute{
									Description: "AccountId",
									Computed:    true,
								},
							},
						},
						common.ToSnakeCase("NotificationMethods"): schema.ListNestedAttribute{
							Description: "Event Policy NotificationList.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("EventLevel"): schema.StringAttribute{
										Description: "EventLevel",
										Computed:    true,
									},
									common.ToSnakeCase("SendMethod"): schema.ListAttribute{
										ElementType: types.StringType,
										Description: "SendMethod",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudMonitoringEventPolicyNotificationDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *cloudMonitoringEventPolicyNotificationDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventPolicyNotificationDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEventPolicyNotificationList(state.ResourceType, state.EventPolicyId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Servers : "+detail,
			err.Error(),
		)
		return
	}

	var NotificationMethods []cloudmonitoring.NotificationMethod

	// Map response body to model
	for _, eventPolicyNotificationsElement := range data.Contents {

		//	NotificationMethod
		for _, notificationMethodData := range eventPolicyNotificationsElement.NotificationMethods {

			var sendMethod []types.String
			sendMethod = make([]types.String, len(notificationMethodData.SendMethod))
			for _, sendMethodType := range notificationMethodData.SendMethod {
				sendMethod = append(sendMethod, types.StringValue(sendMethodType))

			}

			notificationMethodDataElement := cloudmonitoring.NotificationMethod{
				EventLevel: types.StringValue(notificationMethodData.EventLevel),
				SendMethod: sendMethod,
			}

			NotificationMethods = append(NotificationMethods, notificationMethodDataElement)
		}

		//	UserAdditionInfo
		userAdditionalInfoElement := eventPolicyNotificationsElement.GetUserAdditionalInfo()
		userAdditionalInfo := cloudmonitoring.CreateBy{
			CompanyName:        types.StringValue(userAdditionalInfoElement.GetCompanyName()),
			Email:              types.StringValue(userAdditionalInfoElement.GetEmail()),
			Id:                 types.StringValue(userAdditionalInfoElement.GetId()),
			LastLoginAt:        types.StringValue(userAdditionalInfoElement.GetLastLoginAt()),
			MemberCreatedAt:    types.StringValue(userAdditionalInfoElement.GetMemberCreatedAt()),
			PasswordReuseCount: types.Int32Value(userAdditionalInfoElement.GetPasswordReuseCount()),
			UserName:           types.StringValue(userAdditionalInfoElement.GetUserName()),
			TzId:               types.StringValue(userAdditionalInfoElement.GetTzId()),
			Timezone:           types.StringValue(userAdditionalInfoElement.GetTimezone()),
			DstOffset:          types.StringValue(userAdditionalInfoElement.GetDstOffset()),
			UtcOffset:          types.StringValue(userAdditionalInfoElement.GetUtcOffset()),
			AccountId:          types.StringValue(userAdditionalInfoElement.GetAccountId()),
		}

		eventPolicyNotificationState := cloudmonitoring.EventPolicyNotification{
			RecipientType:       types.StringValue(eventPolicyNotificationsElement.RecipientType),
			RecipientKey:        types.StringValue(eventPolicyNotificationsElement.RecipientKey),
			RecipientName:       types.StringValue(eventPolicyNotificationsElement.RecipientName),
			NotificationMethods: NotificationMethods,
			UserAdditionalInfo:  userAdditionalInfo,
		}

		state.EventPolicyNotifications = append(state.EventPolicyNotifications, eventPolicyNotificationState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
