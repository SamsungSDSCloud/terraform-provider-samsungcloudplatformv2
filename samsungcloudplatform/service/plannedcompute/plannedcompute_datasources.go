package billing

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/billing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &billingPlannedComputeDataSource{}
	_ datasource.DataSourceWithConfigure = &billingPlannedComputeDataSource{}
)

func NewBillingPlannedComputeDataSource() datasource.DataSource {
	return &billingPlannedComputeDataSource{}
}

type billingPlannedComputeDataSource struct {
	config  *scpsdk.Configuration
	client  *billing.Client
	clients *client.SCPClient
}

func (d *billingPlannedComputeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_planned_computes"
}

func (d *billingPlannedComputeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of planned computes.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			common.ToSnakeCase("StartDate"): schema.StringAttribute{
				Description: "Start date (YYYY-MM-dd)",
				Optional:    true,
			},
			common.ToSnakeCase("EndDate"): schema.StringAttribute{
				Description: "End date (YYYY-MM-dd)",
				Optional:    true,
			},
			common.ToSnakeCase("ServerType"): schema.StringAttribute{
				Description: "Server Type",
				Optional:    true,
			},
			common.ToSnakeCase("ContractId"): schema.StringAttribute{
				Description: "Contract Id",
				Optional:    true,
			},
			common.ToSnakeCase("ContractType"): schema.ListAttribute{
				Description: "Contract Type",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("NextContractType"): schema.ListAttribute{
				Description: "Next Contract Type",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("ServiceId"): schema.ListAttribute{
				Description: "Service Id",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("OsType"): schema.ListAttribute{
				Description: "OS Type",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description: "Planned Compute State",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "Created by",
				Optional:    true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "Modified by",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("asc", "desc"),
				},
			},

			common.ToSnakeCase("PlannedComputes"): schema.ListNestedAttribute{
				Description: "A list of Planned computes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("ContractId"): schema.StringAttribute{
							Description: "Contract ID",
							Computed:    true,
						},
						common.ToSnakeCase("ContractType"): schema.StringAttribute{
							Description: "Contract Type",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created at",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "Created by",
							Computed:    true,
						},
						common.ToSnakeCase("DeleteYn"): schema.StringAttribute{
							Description: "Delete Y/N",
							Computed:    true,
						},
						common.ToSnakeCase("EndDate"): schema.StringAttribute{
							Description: "End date",
							Computed:    true,
						},
						common.ToSnakeCase("FirstContractStartAt"): schema.StringAttribute{
							Description: "First contract start at",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Planned compute ID",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "Modified at",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "Modified by",
							Computed:    true,
						},
						common.ToSnakeCase("NextContractType"): schema.StringAttribute{
							Description: "Next contract type",
							Computed:    true,
						},
						common.ToSnakeCase("NextEndDate"): schema.StringAttribute{
							Description: "Next end date",
							Computed:    true,
						},
						common.ToSnakeCase("NextStartDate"): schema.StringAttribute{
							Description: "Next end date",
							Computed:    true,
						},
						common.ToSnakeCase("OsName"): schema.StringAttribute{
							Description: "OS name",
							Computed:    true,
						},
						common.ToSnakeCase("OsType"): schema.StringAttribute{
							Description: "OS type",
							Computed:    true,
						},
						common.ToSnakeCase("Region"): schema.StringAttribute{
							Description: "Region",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceName"): schema.StringAttribute{
							Description: "Resource name",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceType"): schema.StringAttribute{
							Description: "Resource type",
							Computed:    true,
						},
						common.ToSnakeCase("ServerType"): schema.StringAttribute{
							Description: "Server type",
							Computed:    true,
						},
						common.ToSnakeCase("ServerTypeDescription"): schema.MapAttribute{
							Description: "Server type description",
							Computed:    true,
							ElementType: types.StringType,
						},
						common.ToSnakeCase("ServiceId"): schema.StringAttribute{
							Description: "Service ID",
							Computed:    true,
						},
						common.ToSnakeCase("ServiceName"): schema.StringAttribute{
							Description: "Service Name",
							Computed:    true,
						},
						common.ToSnakeCase("Srn"): schema.StringAttribute{
							Description: "srn",
							Computed:    true,
						},
						common.ToSnakeCase("StartDate"): schema.StringAttribute{
							Description: "Start date",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *billingPlannedComputeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Billing
	d.clients = inst.Client
}

func (d *billingPlannedComputeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state billing.PlannedComputeDataSourceIds
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPlannedComputeList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PlannedComputes",
			err.Error(),
		)
		return
	}

	for _, plannedCompute := range data.PlannedComputes {
		serverTypeDesc, _ := convertMapStringInterfaceToTypesMap(plannedCompute.GetServerTypeDescription())

		plannedComputeState := billing.PlannedCompute{
			AccountId:             types.StringValue(plannedCompute.GetAccountId()),
			ContractId:            types.StringValue(plannedCompute.GetContractId()),
			ContractType:          types.StringValue(plannedCompute.GetContractType()),
			CreatedAt:             types.StringValue(plannedCompute.GetCreatedAt().Format("2006-01-02 15:04:05")),
			CreatedBy:             types.StringValue(plannedCompute.GetCreatedBy()),
			DeleteYn:              types.StringValue(plannedCompute.GetDeleteYn()),
			EndDate:               types.StringValue(plannedCompute.GetEndDate()),
			FirstContractStartAt:  types.StringValue(plannedCompute.GetFirstContractStartAt()),
			Id:                    types.StringValue(plannedCompute.GetId()),
			ModifiedAt:            types.StringValue(plannedCompute.GetModifiedAt().Format("2006-01-02 15:04:05")),
			ModifiedBy:            types.StringValue(plannedCompute.GetModifiedBy()),
			NextContractType:      types.StringValue(plannedCompute.GetNextContractType()),
			NextEndDate:           types.StringValue(plannedCompute.GetNextEndDate()),
			NextStartDate:         types.StringValue(plannedCompute.GetNextStartDate()),
			OsName:                types.StringValue(plannedCompute.GetOsName()),
			OsType:                types.StringValue(plannedCompute.GetOsType()),
			Region:                types.StringValue(plannedCompute.GetRegion()),
			ResourceName:          types.StringValue(plannedCompute.GetResourceName()),
			ResourceType:          types.StringValue(plannedCompute.GetResourceType()),
			ServerType:            types.StringValue(plannedCompute.GetServerType()),
			ServerTypeDescription: serverTypeDesc,
			ServiceId:             types.StringValue(plannedCompute.GetServiceId()),
			ServiceName:           types.StringValue(plannedCompute.GetServiceName()),
			Srn:                   types.StringValue(plannedCompute.GetSrn()),
			StartDate:             types.StringValue(plannedCompute.GetStartDate()),
			State:                 types.StringValue(plannedCompute.GetState()),
		}
		state.PlannedComputes = append(state.PlannedComputes, plannedComputeState)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
