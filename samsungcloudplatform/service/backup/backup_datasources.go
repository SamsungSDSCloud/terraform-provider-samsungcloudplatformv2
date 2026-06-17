package backup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &backupBackupDataSources{}
	_ datasource.DataSourceWithConfigure = &backupBackupDataSources{}
)

func NewBackupBackupDataSources() datasource.DataSource {
	return &backupBackupDataSources{}
}

type backupBackupDataSources struct {
	config  *scpsdk.Configuration
	client  *backup.Client
	clients *client.SCPClient
}

func (d *backupBackupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_backups"
}

func (d *backupBackupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Backups.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Backup name \n" +
					"  - example: 'terraformtestbackup01'",
				Optional: true,
			},
			common.ToSnakeCase("ServerName"): schema.StringAttribute{
				Description: "Backup server name \n" +
					"  - example: 'terraformbackupserver01'",
				Optional: true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Backup ID List \n " +
					"  - example: ['401e7e1489b94a849714195ee50a1918','55f39954c2684d68ae9691ab0a38d86f']",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *backupBackupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Backup
	d.clients = inst.Client
}

func (d *backupBackupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backup.BackupDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetBackups(d.clients, state.Name, state.ServerName, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Backups",
			err.Error(),
		)
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetBackups(clients *client.SCPClient, name types.String, serverName types.String, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.Backup.GetBackupList(name, serverName)
	if err != nil {
		return nil, err
	}

	contents := data.Contents
	filteredContents := data.Contents

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices, err := filter.GetFilterIndices(contents, filters)
		if err != nil {
			return nil, err
		}

		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
