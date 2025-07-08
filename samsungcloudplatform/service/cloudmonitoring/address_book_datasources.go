package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListAddressBooks
// [GET] /v1/cloudmonitorings/product/v2/users/addrbooks
type cloudMonitoringAddressBookDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringAddressBookDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringAddressBookDataSources{}
)

func NewCloudMonitoringAddressBookSources() datasource.DataSource {
	return &cloudMonitoringAddressBookDataSources{}
}

func (d *cloudMonitoringAddressBookDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_addressbooks"
}

func (d *cloudMonitoringAddressBookDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of AddressBookDataSources.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AddressBooks"): schema.ListNestedAttribute{
				Description: "AddressBooks",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AddrBookName"): schema.StringAttribute{
							Description: "AddrBookName",
							Optional:    true,
						},
						common.ToSnakeCase("AddrbookId"): schema.StringAttribute{
							Description: "AddrbookId",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedByName"): schema.StringAttribute{
							Description: "CreatedByName",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedDt"): schema.StringAttribute{
							Description: "CreatedDt",
							Optional:    true,
						},
						common.ToSnakeCase("MemberCount"): schema.Int32Attribute{
							Description: "MemberCount",
							Optional:    true,
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

func (d *cloudMonitoringAddressBookDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudMonitoringAddressBookDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.AddressBookDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetAddressBookList()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.AddressBook{
			AddrBookName:  types.StringValue(element.GetAddrBookName()),
			AddrbookId:    types.StringValue(element.GetAddrbookId()),
			CreatedBy:     types.StringValue(element.GetCreatedBy()),
			CreatedByName: types.StringValue(element.GetCreatedByName()),
			CreatedDt:     types.StringValue(element.GetCreatedDt().String()),
			MemberCount:   types.Int32Value(element.MemberCount),
		}

		state.AddressBooks = append(state.AddressBooks, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
