package mysql

import (
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
	"time"
)

var (
	_ datasource.DataSource              = &mysqlClusterDataSources{}
	_ datasource.DataSourceWithConfigure = &mysqlClusterDataSources{}
)

func NewMysqlClusterDataSources() datasource.DataSource {
	return &mysqlClusterDataSources{}
}

type mysqlClusterDataSources struct {
	config  *scpsdk.Configuration
	client  *mysql.Client
	clients *client.SCPClient
}

func (d *mysqlClusterDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_clusters"
}

func (d *mysqlClusterDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Clusters.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description: "ServiceState",
				Optional:    true,
			},
			common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
				Description: "DatabaseName",
				Optional:    true,
			},
			common.ToSnakeCase("Clusters"): schema.ListNestedAttribute{
				Description: "A detail of Cluster.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Required:    true,
						},
						common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
							Description: "DatabaseName",
							Required:    true,
						},
						common.ToSnakeCase("HaEnabled"): schema.BoolAttribute{
							Description: "HaEnabled",
							Optional:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Required:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Required:    true,
						},
						common.ToSnakeCase("InstanceCount"): schema.Int32Attribute{
							Description: "InstanceCount",
							Optional:    true,
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: "RoleType",
							Required:    true,
						},
						common.ToSnakeCase("ServiceState"): schema.StringAttribute{
							Description: "ServiceState",
							Required:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "CreatedAt",
							Required:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Required:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Required:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (d *mysqlClusterDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Mysql
	d.clients = inst.Client
}

func (d *mysqlClusterDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mysql.ClusterDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetClusterList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Clusters",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, clusterElement := range data.Contents {
		var clusterState = mysql.Cluster{
			AccountId:     types.StringValue(clusterElement.AccountId),
			DatabaseName:  types.StringPointerValue(clusterElement.DatabaseName.Get()),
			HaEnabled:     types.BoolPointerValue(clusterElement.HaEnabled),
			Id:            types.StringValue(clusterElement.Id),
			Name:          types.StringValue(clusterElement.Name),
			InstanceCount: types.Int32PointerValue(clusterElement.InstanceCount),
			ServiceState:  types.StringValue(string(clusterElement.ServiceState)),
			RoleType:      types.StringPointerValue((*string)(clusterElement.RoleType.Get())),
			CreatedAt:     types.StringValue(clusterElement.CreatedAt.Format(time.RFC3339)),
			CreatedBy:     types.StringValue(clusterElement.CreatedBy),
			ModifiedAt:    types.StringValue(clusterElement.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:    types.StringValue(clusterElement.ModifiedBy),
		}
		state.Clusters = append(state.Clusters, clusterState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
