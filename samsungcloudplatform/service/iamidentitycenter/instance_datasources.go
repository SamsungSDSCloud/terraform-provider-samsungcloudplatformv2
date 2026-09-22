package iamidentitycenter

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamIdentityCenterInstancesDataSources{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterInstancesDataSources{}
)

func NewIamIdentityCenterInstancesDataSources() datasource.DataSource {
	return &iamIdentityCenterInstancesDataSources{}
}

type iamIdentityCenterInstancesDataSources struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterInstancesDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_instances"
}

func (r *iamIdentityCenterInstancesDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IAM Identity Center Instances.",
		Attributes: map[string]schema.Attribute{
			"size": schema.Int32Attribute{
				Optional: true,
				Description: "Number of results to return per page.\n" +
					"  - example : 100\n" +
					"  - min: 1, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional: true,
				Description: "Page number to retrieve.\n" +
					"  - example : 0\n" +
					"  - min: 0, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort the results by a specific field.\n  - example: name:asc",
				MarkdownDescription: "Sort the results by a specific field.\n  - example: name:asc",
			},
			"instances": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of instances matching the filter criteria.\n" +
					"  - example : [{\"id\": \"instance-id\", \"region\": \"kr-west1\", ...}]",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Instance ID\n  - example: ssoins-12345",
							MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
						},
						"master_account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Master Account ID\n  - example: acc-12345",
							MarkdownDescription: "Master Account ID\n  - example: acc-12345",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Organization ID\n  - example: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5",
							MarkdownDescription: "Organization ID\n  - example: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5",
						},
						"region": schema.StringAttribute{
							Computed:            true,
							Description:         "Region\n  - example: kr-west1",
							MarkdownDescription: "Region\n  - example: kr-west1",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "- enum: [\"ACTIVE\",\"INACTIVE\"]",
							MarkdownDescription: "- enum: [\"ACTIVE\",\"INACTIVE\"]",
						},
						"srn": schema.StringAttribute{
							Computed:            true,
							Description:         "Instance SRN\n  - example: srn:identity-center::prj-01234:kr-west-1::instance/ins-01234",
							MarkdownDescription: "Instance SRN\n  - example: srn:identity-center::prj-01234:kr-west-1::instance/ins-01234",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterInstancesDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.clients = inst.Client
}

type instancesDataSourcesModel struct {
	Size      types.Int32                   `tfsdk:"size"`
	Page      types.Int32                   `tfsdk:"page"`
	Sort      types.String                  `tfsdk:"sort"`
	Instances []instanceDataSourceItemModel `tfsdk:"instances"`
}

type instanceDataSourceItemModel struct {
	Id              types.String `tfsdk:"id"`
	MasterAccountId types.String `tfsdk:"master_account_id"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	Region          types.String `tfsdk:"region"`
	State           types.String `tfsdk:"state"`
	Srn             types.String `tfsdk:"srn"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (r *iamIdentityCenterInstancesDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instancesDataSourcesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	size := state.Size.ValueInt32()
	page := state.Page.ValueInt32()
	sort := state.Sort.ValueString()

	result, err := r.clients.IamIdentityCenter.ListInstances(ctx, size, page, sort)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Instances",
			err.Error(),
		)
		return
	}

	if result != nil && result.Instances != nil {
		state.Instances = make([]instanceDataSourceItemModel, len(result.Instances))
		for i, instance := range result.Instances {
			item := instanceDataSourceItemModel{
				Id:              types.StringValue(instance.Id),
				MasterAccountId: types.StringValue(instance.MasterAccountId),
				OrganizationId:  types.StringValue(instance.OrganizationId),
				Region:          types.StringValue(instance.Region),
				Srn:             types.StringValue(instance.Srn),
				CreatedAt:       types.StringValue(instance.CreatedAt.Format(time.RFC3339)),
				CreatedBy:       types.StringValue(instance.CreatedBy),
				ModifiedAt:      types.StringValue(instance.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:      types.StringValue(instance.ModifiedBy),
			}
			if instance.State != nil {
				item.State = types.StringValue(string(*instance.State))
			}
			state.Instances[i] = item
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
