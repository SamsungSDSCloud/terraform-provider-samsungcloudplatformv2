package filestorage

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/filestorage"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileStorageAccessRulesDataSources{}
	_ datasource.DataSourceWithConfigure = &fileStorageAccessRulesDataSources{}
)

func NewFileStorageAccessRulesDataSources() datasource.DataSource {
	return &fileStorageAccessRulesDataSources{}
}

type fileStorageAccessRulesDataSources struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (d *fileStorageAccessRulesDataSources) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	d.client = inst.Client.FileStorage
	d.clients = inst.Client
}

func (d *fileStorageAccessRulesDataSources) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_access_rules"
}

func (d *fileStorageAccessRulesDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Lists all access rules for a File Storage Volume. This data source is read-only and includes rules managed by external automation (e.g., Kubernetes Auto Scaling). Do NOT use this data source output with `for_each` to create `filestorage_access_rule` resources — doing so will re-introduce the forced reconciliation problem.",
		Attributes: map[string]schema.Attribute{
			"file_storage_id": schema.StringAttribute{
				Description: "The identifier of the File Storage Volume to query access rules for.\n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required: true,
			},
			"access_rules": schema.ListNestedAttribute{
				Description: "The complete list of access rules at the time of the query, including rules added by external automation.\n" +
					"  - This is read-only and does not trigger any drift or deletion.\n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"object_id": schema.StringAttribute{
							Description: "The identifier of the object granted access.\n" +
								"  - example : '43fq3347-02q4-4aa8-ccf9-affe4917bb6f' \n",
							Computed: true,
						},
						"object_type": schema.StringAttribute{
							Description: "The type of the object granted access.\n" +
								"  - example : 'VM' \n" +
								"  - valid : VM, BM, GPU, GPU_NODE, ENDPOINT \n",
							Computed: true,
						},
					},
				},
			},
		},
	}
}


func (d *fileStorageAccessRulesDataSources) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state filestorage.FileStorageAccessRulesDataSource

	diags := request.Config.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetVolumeAccessRules(ctx, state.FileStorageId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading File Storage Access Rules",
			"Could not read file storage access rules: "+err.Error(),
		)
		return
	}
	
	accessRules := []filestorage.AccessRuleResource{}
	for _, rule := range data.AccessRules {
		accessRules = append(accessRules, filestorage.AccessRuleResource{
			ObjectId:   types.StringValue(rule.ObjectId),
			ObjectType: types.StringValue(rule.ObjectType),
		})
	}

	state.AccessRules = accessRules

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}
