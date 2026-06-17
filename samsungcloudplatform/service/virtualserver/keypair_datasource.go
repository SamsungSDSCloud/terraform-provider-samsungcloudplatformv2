package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerKeypairDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualServerKeypairDataSource{}
)

func NewVirtualServerKeypairDataSource() datasource.DataSource {
	return &virtualServerKeypairDataSource{}
}

type virtualServerKeypairDataSource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerKeypairDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_keypair"
}

func (d *virtualServerKeypairDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves keypair information.",
		MarkdownDescription: "Retrieves keypair information for SSH access to virtual servers.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Keypair name.\n" +
					"  - example: my-keypair\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				MarkdownDescription: "Keypair name.\n" +
					"  - example: my-keypair\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("Keypair"): schema.SingleNestedAttribute{
				Description:         "Keypair details.",
				MarkdownDescription: "Keypair details including name, public key, fingerprint, and type.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Keypair name.",
						MarkdownDescription: "Keypair name.",
						Computed:            true,
					},
					common.ToSnakeCase("PublicKey"): schema.StringAttribute{
						Description:         "Public key.",
						MarkdownDescription: "Public key in OpenSSH format.",
						Computed:            true,
					},
					common.ToSnakeCase("Fingerprint"): schema.StringAttribute{
						Description:         "Fingerprint.",
						MarkdownDescription: "Fingerprint of the public key.",
						Computed:            true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description:         "Keypair type.",
						MarkdownDescription: "Keypair type.",
						Computed:            true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerKeypairDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VirtualServer
	d.clients = inst.Client
}

func (d *virtualServerKeypairDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.KeypairDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	names, err := GetKeypairs(d.clients, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	if len(names) > 0 || !state.Name.IsNull() {
		name := virtualserverutil.SetResourceIdentifier(state.Name, names, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		keypair, err := d.client.GetKeypair(ctx, name.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Keypair",
				"Could not read Keypair Name "+name.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		keypairModel := virtualserver.Keypair{
			Name:        types.StringValue(keypair.Name),
			PublicKey:   types.StringValue(keypair.PublicKey),
			Fingerprint: types.StringValue(keypair.Fingerprint),
			Type:        types.StringValue(keypair.Type),
		}
		keypairObjectValue, _ := types.ObjectValueFrom(ctx, keypairModel.AttributeTypes(), keypairModel)
		state.Keypair = keypairObjectValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
