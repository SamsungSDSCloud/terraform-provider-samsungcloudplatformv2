package parallelfilestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/parallelfilestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &parallelFileStorageVolumeDataSource{}
	_ datasource.DataSourceWithConfigure = &parallelFileStorageVolumeDataSource{}
)

func NewParallelFileStorageVolumeDataSource() datasource.DataSource {
	return &parallelFileStorageVolumeDataSource{}
}

type parallelFileStorageVolumeDataSource struct {
	config  *scpsdk.Configuration
	client  *parallelfilestorage.Client
	clients *client.SCPClient
}

func (d *parallelFileStorageVolumeDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parallel_filestorage_volume"
}

func (d *parallelFileStorageVolumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VolumeDataSourceSchema()
}
func VolumeDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed: true,
				Description: "Account ID \n" +
					"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Created At \n" +
					"  - example : '2024-07-30T04:54:33.219373' \n",
			},
			"id": schema.StringAttribute{
				Required: true,
				Description: "ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Volume Name \n" +
					"  - example : 'my_volume' \n",
			},
			"mount_path": schema.StringAttribute{
				Computed: true,
				Description: "Volume Mount Path \n" +
					"  - example : 'xxx.xx.xxx.xxx'",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Volume State \n" +
					"  - example : 'available' \n",
			},
			"zone": schema.StringAttribute{
				Computed: true,
				Description: "Zone \n" +
					"  - example : 'kr-west1-a' \n",
			},
			"capacity_tb": schema.Int32Attribute{
				Computed: true,
				Description: "Volume Capacity \n" +
					"  - example : '100' \n",
			},
			common.ToSnakeCase("AccessRules"): schema.SetNestedAttribute{
				Description: "Object of AccessRule",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ObjectId"): schema.StringAttribute{
							Description: "Object Id \n" +
								"  - example : '43fq3347-02q4-4aa8-ccf9-affe4917bb6f' \n",
							Optional: true,
						},
						common.ToSnakeCase("ObjectType"): schema.StringAttribute{
							Description: "Object Type" +
								"  - example : 'VM' \n" +
								"  - pattern: `^(VM|BM|GPU|GPU_NODE|ENDPOINT)$` \n",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (d *parallelFileStorageVolumeDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Plase report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	d.client = inst.Client.ParallelFileStorage
	d.clients = inst.Client
}

func (d *parallelFileStorageVolumeDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state parallelfilestorage.VolumeDataSource
	diags := request.Config.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if state.Id.IsNull() || state.Id.IsUnknown() || state.Id.ValueString() == "" {
		response.Diagnostics.AddError("Missing Volume ID", "The volume id must be provided.")
		return
	}
	
	volume, err := d.client.GetVolume(ctx, state.Id.ValueString())

	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError("Error Reading Volume",
			"Could not read Volume Id "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail)    
		return
	}
	if volume.AccountId != "" {
		state.AccountId = types.StringValue(volume.AccountId)
	}
	if !volume.CreatedAt.IsZero() {
		state.CreatedAt = types.StringValue(volume.CreatedAt.Format(time.RFC3339))
	}
	if volume.Id != "" {
		state.Id = types.StringValue(volume.Id)
	}
	if volume.Name != "" {
		state.Name = types.StringValue(volume.Name)
	}
	if volume.Zone != "" {
		state.Zone = types.StringValue(volume.Zone)
	}

	if volume.State != "" {
		state.State = types.StringValue(volume.State)
	}
	if volume.MountPath.Get() != nil {
		state.MountPath = types.StringValue(*volume.MountPath.Get())
	}
	state.CapacityTb = types.Int32Value(volume.CapacityTb)

	// AccessRule
	getAccessRule, err := d.client.GetVolumeAccessRules(ctx, state.Id.ValueString())

	var accessRules []parallelfilestorage.AccessRuleResource
	if err != nil || getAccessRule == nil || len(getAccessRule.AccessRules) == 0 {
		accessRules = []parallelfilestorage.AccessRuleResource{}
	} else {
		for _, rules := range getAccessRule.AccessRules {
			rule := parallelfilestorage.AccessRuleResource{
				ObjectId:   types.StringValue(rules.ObjectId),
				ObjectType: types.StringValue(rules.ObjectType),
			}
			accessRules = append(accessRules, rule)
		}
	}
	state.AccessRules = accessRules

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
