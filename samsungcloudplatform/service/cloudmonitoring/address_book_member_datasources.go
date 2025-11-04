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

// ListAddressBookMembers
// [GET] /v1/cloudmonitorings/product/v2/addrbooks/{{addrbookId}}/members
type cloudMonitoringAddressBookMemberDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringAddressBookMemberDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringAddressBookMemberDataSources{}
)

func NewCloudMonitoringAddressMemberBookSources() datasource.DataSource {
	return &cloudMonitoringAddressBookMemberDataSources{}
}

func (d *cloudMonitoringAddressBookMemberDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_addressbookmembers"
}

func (d *cloudMonitoringAddressBookMemberDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of AddressBookMemberDataSources.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AddrbookId"): schema.Int32Attribute{
				Description: "AddrbookId",
				Required:    true,
			},
			common.ToSnakeCase("AddressBookMembers"): schema.ListNestedAttribute{
				Description: "AddressBookMembers",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("UserEmail"): schema.StringAttribute{
							Description: "UserEmail",
							Required:    true,
						},
						common.ToSnakeCase("UserId"): schema.StringAttribute{
							Description: "UserId",
							Required:    true,
						},
						common.ToSnakeCase("UserLanguage"): schema.StringAttribute{
							Description: "UserLanguage",
							Required:    true,
						},
						common.ToSnakeCase("UserMobileTelNo"): schema.StringAttribute{
							Description: "UserMobileTelNo",
							Required:    true,
						},
						common.ToSnakeCase("UserName"): schema.StringAttribute{
							Description: "UserName",
							Required:    true,
						},
						common.ToSnakeCase("UserTimezone"): schema.StringAttribute{
							Description: "UserTimezone",
							Required:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *cloudMonitoringAddressBookMemberDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudMonitoringAddressBookMemberDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.AddressBookMemberDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetAddressBookMemberList(state.AddrbookId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.AddressBookMember{
			UserEmail:       types.StringValue(element.GetUserEmail()),
			UserId:          types.StringValue(element.GetUserId()),
			UserLanguage:    types.StringValue(element.GetUserLanguage()),
			UserMobileTelNo: types.StringValue(element.GetUserMobileTelNo()),
			UserName:        types.StringValue(element.GetUserName()),
			UserTimezone:    types.StringValue(element.GetUserTimezone()),
		}

		state.AddressBookMembers = append(state.AddressBookMembers, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
