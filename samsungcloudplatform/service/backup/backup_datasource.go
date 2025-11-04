package backup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	backuputil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ datasource.DataSource              = &backupBackupDataSource{}
	_ datasource.DataSourceWithConfigure = &backupBackupDataSource{}
)

func NewBackupBackupDataSource() datasource.DataSource {
	return &backupBackupDataSource{}
}

type backupBackupDataSource struct {
	config  *scpsdk.Configuration
	client  *backup.Client
	clients *client.SCPClient
}

func (d *backupBackupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_backup"
}

func (d *backupBackupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Backup.",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Backup name",
				Optional:    true,
			},
			common.ToSnakeCase("ServerName"): schema.StringAttribute{
				Description: "Backup server name",
				Optional:    true,
			},
			common.ToSnakeCase("Backup"): schema.SingleNestedAttribute{
				Description: "A detail of Backup.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Backup server name",
						Computed:    true,
					},
					common.ToSnakeCase("PolicyType"): schema.StringAttribute{
						Description: "Backup policy type",
						Computed:    true,
					},
					common.ToSnakeCase("RetentionPeriod"): schema.StringAttribute{
						Description: "Backup retention period",
						Computed:    true,
					},
					common.ToSnakeCase("RoleType"): schema.StringAttribute{
						Description: "Backup role type",
						Computed:    true,
					},
					common.ToSnakeCase("ServerName"): schema.StringAttribute{
						Description: "Backup server name",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Backup state",
						Computed:    true,
					},
					common.ToSnakeCase("EncryptEnabled"): schema.BoolAttribute{
						Description: "Whether to use Encryption",
						Computed:    true,
					},
					common.ToSnakeCase("PolicyCategory"): schema.StringAttribute{
						Description: "Backup policy category",
						Computed:    true,
					},
					common.ToSnakeCase("ServerCategory"): schema.StringAttribute{
						Description: "Backup server category",
						Computed:    true,
					},
					common.ToSnakeCase("ServerUuid"): schema.StringAttribute{
						Description: "Backup server UUID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},

					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *backupBackupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupBackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backup.BackupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		d.client.Config.Region = state.Region.ValueString()
	}

	ids, err := GetBackups(d.clients, state.Name, state.ServerName, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Backup",
			err.Error(),
		)
	}

	if len(ids) > 0 || !state.Id.IsNull() {
		id := backuputil.SetResourceIdentifier(state.Id, ids, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		data, err := d.client.GetBackup(ctx, id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Backup",
				"Could not read Backup ID "+id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		var backupModel = backup.Backup{
			State:           types.StringValue(data.State),
			ServerName:      types.StringValue(data.ServerName),
			RoleType:        types.StringValue(data.RoleType),
			RetentionPeriod: types.StringValue(data.RetentionPeriod),
			PolicyType:      types.StringValue(data.PolicyType),
			Name:            types.StringValue(data.Name),
			EncryptEnabled:  types.BoolPointerValue(data.EncryptEnabled.Get()),
			PolicyCategory:  types.StringValue(data.PolicyCategory),
			ServerCategory:  types.StringValue(data.ServerCategory),
			ServerUuid:      types.StringValue(data.ServerUuid),
			ModifiedBy:      types.StringValue(data.ModifiedBy),
			ModifiedAt:      types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
			Id:              types.StringValue(data.Id),
			CreatedBy:       types.StringValue(data.CreatedBy),
			CreatedAt:       types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		}
		backupObjectValue, _ := types.ObjectValueFrom(ctx, backupModel.AttributeTypes(), backupModel)
		state.Backup = backupObjectValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
