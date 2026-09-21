package postgresql

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &postgresqlClusterDataSources{}
	_ datasource.DataSourceWithConfigure = &postgresqlClusterDataSources{}
)

func NewPostgresqlClusterDataSources() datasource.DataSource {
	return &postgresqlClusterDataSources{}
}

type postgresqlClusterDataSources struct {
	config  *scpsdk.Configuration
	client  *postgresql.Client
	clients *client.SCPClient
}

func (d *postgresqlClusterDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_clusters"
}

func (d *postgresqlClusterDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Clusters.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0 ",
				MarkdownDescription: "The page number for pagination.\n" +
					"  - example : 0 ",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20 ",
				MarkdownDescription: "The number of items per page.\n" +
					"  - example : 20 ",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:asc ",
				MarkdownDescription: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:asc ",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Cluster name\n  - example: mytest",
				MarkdownDescription: "Cluster name\n  - example: mytest",
				Optional:            true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description:         "Service state\n  - example: RUNNING",
				MarkdownDescription: "Service state\n  - example: RUNNING",
				Optional:            true,
			},
			common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
				Description:         "Database name\n  - example: mydb",
				MarkdownDescription: "Database name\n  - example: mydb",
				Optional:            true,
			},
			common.ToSnakeCase("Clusters"): schema.ListNestedAttribute{
				Description: "A detail of Cluster.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the endpoint.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							MarkdownDescription: "The identifier of the account that owns the endpoint.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Required: true,
						},
						common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
							Description:         "Database name\n  - example: mydb",
							MarkdownDescription: "Database name\n  - example: mydb",
							Required:            true,
						},
						common.ToSnakeCase("HaEnabled"): schema.BoolAttribute{
							Description:         "HA availability\n  - example: false",
							MarkdownDescription: "HA availability\n  - example: false",
							Optional:            true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
							MarkdownDescription: "Identifier of the resource.\n  - example: 35e21d596d4f41e9b7b66d8f2129213a",
							Required:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "Cluster name\n  - example: mytest",
							MarkdownDescription: "Cluster name\n  - example: mytest",
							Required:            true,
						},
						common.ToSnakeCase("InstanceCount"): schema.Int32Attribute{
							Description:         "Instance Count\n  - example: 1",
							MarkdownDescription: "Instance Count\n  - example: 1",
							Optional:            true,
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description:         "Role type\n  - example: ORIGIN",
							MarkdownDescription: "Role type\n  - example: ORIGIN",
							Required:            true,
						},
						common.ToSnakeCase("ServiceState"): schema.StringAttribute{
							Description:         "Service state\n  - example: RUNNING",
							MarkdownDescription: "Service state\n  - example: RUNNING",
							Required:            true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
							Required:            true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description:         "Created by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
							MarkdownDescription: "Created by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
							Required:            true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
							Required:            true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description:         "Modified by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
							MarkdownDescription: "Modified by\n  - example: 7d21d8f464b54de6a44ebfd2c0a56787",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *postgresqlClusterDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Postgresql
	d.clients = inst.Client
}

func (d *postgresqlClusterDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresql.ClusterDataSource

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
		clusterState := postgresql.Cluster{
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
