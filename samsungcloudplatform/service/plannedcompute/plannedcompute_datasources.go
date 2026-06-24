package billing

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/billing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description:         "Limit\n  - example: 10",
				MarkdownDescription: "Limit\n  - example: 10",
				Optional:             true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description:         "Page\n  - example: 1",
				MarkdownDescription: "Page\n  - example: 1",
				Optional:             true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			common.ToSnakeCase("StartDate"): schema.StringAttribute{
				Description:         "Start date (YYYY-MM-dd)\n  - example: 2024-01-01",
				MarkdownDescription: "Start date (YYYY-MM-dd)\n  - example: 2024-01-01",
				Optional:             true,
			},
			common.ToSnakeCase("EndDate"): schema.StringAttribute{
				Description:         "End date (YYYY-MM-dd)\n  - example: 2024-12-31",
				MarkdownDescription: "End date (YYYY-MM-dd)\n  - example: 2024-12-31",
				Optional:             true,
			},
			common.ToSnakeCase("ServerType"): schema.StringAttribute{
				Description:         "Server Type\n  - example: s1v1m2",
				MarkdownDescription: "Server Type\n  - example: s1v1m2",
				Optional:             true,
			},
			common.ToSnakeCase("ContractId"): schema.StringAttribute{
				Description:         "Contract Id\n  - example: C1234567",
				MarkdownDescription: "Contract Id\n  - example: C1234567",
				Optional:             true,
			},
			common.ToSnakeCase("ContractType"): schema.ListAttribute{
				Description:         "Contract Type\n  - example: 01",
				MarkdownDescription: "Contract Type\n  - example: 01",
				Optional:             true,
				ElementType:          types.StringType,
			},
			common.ToSnakeCase("NextContractType"): schema.ListAttribute{
				Description:         "Next Contract Type\n  - example: 03",
				MarkdownDescription: "Next Contract Type\n  - example: 03",
				Optional:             true,
				ElementType:          types.StringType,
			},
			common.ToSnakeCase("ServiceId"): schema.ListAttribute{
				Description:         "Service Id\n  - example: VIRTUAL_SERVER",
				MarkdownDescription: "Service Id\n  - example: VIRTUAL_SERVER",
				Optional:             true,
				ElementType:          types.StringType,
			},
			common.ToSnakeCase("OsType"): schema.ListAttribute{
				Description:         "OS Type\n  - example: rhel",
				MarkdownDescription: "OS Type\n  - example: rhel",
				Optional:             true,
				ElementType:          types.StringType,
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description:         "Planned Compute State\n  - example: ACTIVE",
				MarkdownDescription: "Planned Compute State\n  - example: ACTIVE",
				Optional:             true,
				ElementType:          types.StringType,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description:         "Created by\n  - example: 83c3c73d457345e3829ee6d5557c0011",
				MarkdownDescription: "Created by\n  - example: 83c3c73d457345e3829ee6d5557c0011",
				Optional:             true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description:         "Modified by\n  - example: 83c3c73d457345e3829ee6d5557c0011",
				MarkdownDescription: "Modified by\n  - example: 83c3c73d457345e3829ee6d5557c0011",
				Optional:             true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description:         "Sort\n  - example: asc",
				MarkdownDescription: "Sort\n  - example: asc",
				Optional:             true,
				Validators: []validator.String{
					stringvalidator.OneOf("asc", "desc"),
				},
			},

			common.ToSnakeCase("PlannedComputes"): schema.ListNestedAttribute{
                Description: "A list of Planned computes.\n" +
                    "  - example: [{account_id='f5c8e56a4d9b49a8bd89e14758a32d53', contract_id='C1234567, contract_type='01', state='ACTIVE'}]",
                MarkdownDescription: "A list of Planned computes.\n" +
                    "  - example: [{account_id='f5c8e56a4d9b49a8bd89e14758a32d53', contract_id='C1234567', contract_type='01', state='ACTIVE'}]",
                Computed:    true,
                NestedObject: schema.NestedAttributeObject{
                   Attributes: map[string]schema.Attribute{
                      common.ToSnakeCase("AccountId"): schema.StringAttribute{
                         Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
                         MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ContractId"): schema.StringAttribute{
                         Description:         "Contract ID\n  - example: C1234567",
                         MarkdownDescription: "Contract ID\n  - example: C1234567",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ContractType"): schema.StringAttribute{
                         Description:         "Contract Type\n  - example: 01",
                         MarkdownDescription: "Contract Type\n  - example: 01",
                         Computed:            true,
                      },
                      common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                         Description:         "Created at\n  - example: 2024-05-17T00:23:17Z",
                         MarkdownDescription: "Created at\n  - example: 2024-05-17T00:23:17Z",
                         Computed:            true,
                      },
                      common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                         Description:         "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
                         MarkdownDescription: "Created by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
                         Computed:            true,
                      },
                      common.ToSnakeCase("DeleteYn"): schema.StringAttribute{
                         Description:         "Delete Y/N\n  - example: N",
                         MarkdownDescription: "Delete Y/N\n  - example: N",
                         Computed:            true,
                      },
                      common.ToSnakeCase("EndDate"): schema.StringAttribute{
                         Description:         "End date\n  - example: 2025-05-17",
                         MarkdownDescription: "End date\n  - example: 2025-05-17",
                         Computed:            true,
                      },
                      common.ToSnakeCase("FirstContractStartAt"): schema.StringAttribute{
                         Description:         "First contract start at\n  - example: 2023-05-17T00:23:17Z",
                         MarkdownDescription: "First contract start at\n  - example: 2023-05-17T00:23:17Z",
                         Computed:            true,
                      },
                      common.ToSnakeCase("Id"): schema.StringAttribute{
                         Description:         "Planned compute ID\n  - example: Planned-Compute-",
                         MarkdownDescription: "Planned compute ID\n  - example: 83c3c73d457345e3829ee6d5557c0011",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                         Description:         "Modified at\n  - example: 2024-06-24T14:02:10Z",
                         MarkdownDescription: "Modified at\n  - example: 2024-06-24T14:02:10Z",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                         Description:         "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
                         MarkdownDescription: "Modified by\n  - example: ef716e80-1fac-4faa-892d-0132fc7f5583",
                         Computed:            true,
                      },
                      common.ToSnakeCase("NextContractType"): schema.StringAttribute{
                         Description:         "Next contract type\n  - example: 03",
                         MarkdownDescription: "Next contract type\n  - example: 03",
                         Computed:            true,
                      },
                      common.ToSnakeCase("NextEndDate"): schema.StringAttribute{
                         Description:         "Next end date\n  - example: 2026-05-17",
                         MarkdownDescription: "Next end date\n  - example: 2026-05-17",
                         Computed:            true,
                      },
                      common.ToSnakeCase("NextStartDate"): schema.StringAttribute{
                         Description:         "Next start date\n  - example: 2025-05-18",
                         MarkdownDescription: "Next start date\n  - example: 2025-05-18",
                         Computed:            true,
                      },
                      common.ToSnakeCase("OsName"): schema.StringAttribute{
                         Description:         "OS name\n  - example: RHEL",
                         MarkdownDescription: "OS name\n  - example: RHEL",
                         Computed:            true,
                      },
                      common.ToSnakeCase("OsType"): schema.StringAttribute{
                         Description:         "OS type\n  - example: rhel",
                         MarkdownDescription: "OS type\n  - example: rhel",
                         Computed:            true,
                      },
                      common.ToSnakeCase("Region"): schema.StringAttribute{
                         Description:         "Region\n  - example: kr-west1",
                         MarkdownDescription: "Region\n  - example: kr-west1",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ResourceName"): schema.StringAttribute{
                         Description:         "Resource name\n  - example: Planned-compute-01",
                         MarkdownDescription: "Resource name\n  - example: Planned-compute-01",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ResourceType"): schema.StringAttribute{
                         Description:         "Resource type\n  - example: BARE_METAL",
                         MarkdownDescription: "Resource type\n  - example: BARE_METAL",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ServerType"): schema.StringAttribute{
                         Description:         "Server type\n  - example: s1v1m2",
                         MarkdownDescription: "Server type\n  - example: s1v1m2",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ServerTypeDescription"): schema.MapAttribute{
                         Description:         "Server type description\n  - example: {\"cpu\": \"16 Cores\"| \"memory\": \"64 GB\"}",
                         MarkdownDescription: "Server type description\n  - example: {\"cpu\": \"16 Cores\"| \"memory\": \"64 GB\"}",
                         Computed:            true,
                         ElementType:         types.StringType,
                      },
                      common.ToSnakeCase("ServiceId"): schema.StringAttribute{
                         Description:         "Service ID\n  - example: VIRTUAL_SERVER",
                         MarkdownDescription: "Service ID\n  - example: VIRTUAL_SERVER",
                         Computed:            true,
                      },
                      common.ToSnakeCase("ServiceName"): schema.StringAttribute{
                         Description:         "Service Name\n  - example: Virtual Server",
                         MarkdownDescription: "Service Name\n  - example: Virtual Server",
                         Computed:            true,
                      },
                      common.ToSnakeCase("Srn"): schema.StringAttribute{
                         Description:         "srn\n  - example: srn:e::26affb52e16944038a0cd2cc26060e1c:kr1-west1::compute:instance/INSTANCE-UPOg3Z6ZqyiMM0QyC3sI2m",
                         MarkdownDescription: "srn\n  - example: srn:e::26affb52e16944038a0cd2cc26060e1c:kr1-west1::compute:instance/INSTANCE-UPOg3Z6ZqyiMM0QyC3sI2m",
                         Computed:            true,
                      },
                      common.ToSnakeCase("StartDate"): schema.StringAttribute{
                         Description:         "Start date\n  - example: 2024-05-17",
                         MarkdownDescription: "Start date\n  - example: 2024-05-17",
                         Computed:            true,
                      },
                      common.ToSnakeCase("State"): schema.StringAttribute{
                         Description:         "State\n  - example: ACTIVE",
                         MarkdownDescription: "State\n  - example: ACTIVE",
                         Computed:            true,
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
