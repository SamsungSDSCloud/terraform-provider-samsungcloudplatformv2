package servicewatch

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchLogGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &serviceWatchLogGroupDataSources{}
)

// NewServiceWatchLogGroupDataSources is a helper function to simplify the provider implementation.
func NewServiceWatchLogGroupDataSources() datasource.DataSource {
	return &serviceWatchLogGroupDataSources{}
}

// serviceWatchLogGroupDataSources is the data source implementation.
type serviceWatchLogGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchLogGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_log_groups"
}

// Schema defines the schema for the data sources.
func (d *serviceWatchLogGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Group Data Sources",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Ids"): schema.StringAttribute{
				Description: "List of Log group IDs.\n" +
					" - example : bce52822147744b4afe0187164caa2e8\n",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of log group.\n" +
					" - example : testlg01\n" +
					" - minLength: 3\n" +
					" - maxLength: 512\n",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 512),
				},
			},
			common.ToSnakeCase("RetentionPeriods"): schema.ListAttribute{
				ElementType: types.Int32Type,
				Description: "List of Log group retention periods.\n" +
					" - example : [365]\n",
				Optional:    true,
			},
			common.ToSnakeCase("LogGroups"): schema.ListNestedAttribute{
				Description: "List of Log groups.\n" +
					" - example : [{\"id\": \"bce52822147744b4afe0187164caa2e8\", \"name\": \"testlg01\"}]\n",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the log group.\n" +
								" - example : bce52822147744b4afe0187164caa2e8\n",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of log group.\n" +
								" - example : testlg01\n",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The unique identifier of the account.\n" +
								" - example : 1bcf39b344ac41cbaf0466ff0d2bebad\n",
							Computed:    true,
						},
						common.ToSnakeCase("RetentionPeriod"): schema.Int32Attribute{
							Description: "Log group retention period.\n" +
								" - example : 365\n",
							Computed:    true,
						},
						common.ToSnakeCase("RetentionPeriodName"): schema.StringAttribute{
							Description: "Log group retention period name.\n" +
								" - example : 1 year\n",
							Computed:    true,
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
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
								" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
								" - example : 2024-05-17T00:23:17Z\n",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that last modified the resource.\n" +
								" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data sources.
func (d *serviceWatchLogGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceWatchLogGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.LogGroupDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLogGroupList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Log Group List",
			err.Error(),
		)
		return
	}

	if data == nil {
		state.LogGroups = types.ListValueMust(types.ObjectType{AttrTypes: servicewatch.LogGroup{}.AttributeTypes()}, []attr.Value{})
	} else {
		var logGroupList []servicewatch.LogGroup
		for _, logGroup := range data.LogGroups {
			logGroups := servicewatch.LogGroup{
				Id:                  types.StringValue(logGroup.GetId()),
				Name:                types.StringValue(logGroup.GetName()),
				AccountId:           types.StringValue(logGroup.GetAccountId()),
				RetentionPeriod:     types.Int32Value(logGroup.GetRetentionPeriod()),
				RetentionPeriodName: types.StringValue(logGroup.GetRetentionPeriodName()),
				Status:              types.StringValue(string(logGroup.GetStatus())),
				CreatedAt:           types.StringValue(logGroup.GetCreatedAt().Format("2006-01-02 15:04:05")),
				CreatedBy:           types.StringValue(logGroup.GetCreatedBy()),
				ModifiedAt:          types.StringValue(logGroup.GetModifiedAt().Format("2006-01-02 15:04:05")),
				ModifiedBy:          types.StringValue(logGroup.GetModifiedBy()),
			}
			logGroupList = append(logGroupList, logGroups)
			state.LogGroups, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.LogGroup{}.AttributeTypes()}, logGroupList)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
