package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerImageDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualServerImageDataSource{}
)

func NewVirtualServerImageDataSource() datasource.DataSource {
	return &virtualServerImageDataSource{}
}

type virtualServerImageDataSource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_image"
}

func (d *virtualServerImageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Image.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
				Description: "SCP Image type",
				Optional:    true,
			},
			common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
				Description: "SCP Original Image type",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("OsDistro"): schema.StringAttribute{
				Description: "OS Distro",
				Optional:    true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Status",
				Optional:    true,
			},
			common.ToSnakeCase("Visibility"): schema.StringAttribute{
				Description: "Visibility",
				Optional:    true,
			},
			common.ToSnakeCase("Image"): schema.SingleNestedAttribute{
				Description: "Image.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Volumes"): schema.StringAttribute{
						Description: "Volumes",
						Computed:    true,
					},
					common.ToSnakeCase("Checksum"): schema.StringAttribute{
						Description: "Checksum",
						Computed:    true,
					},
					common.ToSnakeCase("ContainerFormat"): schema.StringAttribute{
						Description: "Container format",
						Computed:    true,
					},
					common.ToSnakeCase("DiskFormat"): schema.StringAttribute{
						Description: "Disk format",
						Computed:    true,
					},
					common.ToSnakeCase("File"): schema.StringAttribute{
						Description: "File",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID",
						Computed:    true,
					},
					common.ToSnakeCase("MinDisk"): schema.Int32Attribute{
						Description: "Min disk",
						Computed:    true,
					},
					common.ToSnakeCase("MinRam"): schema.Int32Attribute{
						Description: "Min ram",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("OsDistro"): schema.StringAttribute{
						Description: "OS Distro",
						Computed:    true,
					},
					common.ToSnakeCase("OsHashAlgo"): schema.StringAttribute{
						Description: "OS Hash algo",
						Computed:    true,
					},
					common.ToSnakeCase("OsHashValue"): schema.StringAttribute{
						Description: "OS Hash value",
						Computed:    true,
					},
					common.ToSnakeCase("OsHidden"): schema.BoolAttribute{
						Description: "OS hidden",
						Computed:    true,
					},
					common.ToSnakeCase("Owner"): schema.StringAttribute{
						Description: "Owner",
						Computed:    true,
					},
					common.ToSnakeCase("OwnerAccountName"): schema.StringAttribute{
						Description: "Owner account name",
						Computed:    true,
					},
					common.ToSnakeCase("OwnerUserName"): schema.StringAttribute{
						Description: "Owner user name",
						Computed:    true,
					},
					common.ToSnakeCase("Protected"): schema.BoolAttribute{
						Description: "Protected",
						Computed:    true,
					},
					common.ToSnakeCase("RootDeviceName"): schema.StringAttribute{
						Description: "Root device name",
						Computed:    true,
					},
					common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
						Description: "SCP Image type",
						Computed:    true,
					},
					common.ToSnakeCase("ScpK8sVersion"): schema.StringAttribute{
						Description: "SCP K8s version",
						Computed:    true,
					},
					common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
						Description: "SCP original Image type",
						Computed:    true,
					},
					common.ToSnakeCase("ScpOsVersion"): schema.StringAttribute{
						Description: "SCP OS version",
						Computed:    true,
					},
					common.ToSnakeCase("Size"): schema.Int64Attribute{
						Description: "Size",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					common.ToSnakeCase("VirtualSize"): schema.Int64Attribute{
						Description: "Virtual size",
						Computed:    true,
					},
					common.ToSnakeCase("Visibility"): schema.StringAttribute{
						Description: "Visibility",
						Computed:    true,
					},
					common.ToSnakeCase("Url"): schema.StringAttribute{
						Description: "Url",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created at",
						Computed:    true,
					},
					common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
						Description: "Updated at",
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

func (d *virtualServerImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *virtualServerImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ImageDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetImages(d.clients, state.ScpImageType, state.ScpOriginalImageType, state.Name, state.OsDistro, state.Status, state.Visibility, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	if len(ids) > 0 || !state.Id.IsNull() {
		id := virtualserverutil.SetResourceIdentifier(state.Id, ids, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		image, err := d.client.GetImage(ctx, id.ValueString()) // client 를 호출한다.
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Image",
				"Could not read Image ID "+id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		imageModel := virtualserver.Image{
			Volumes:              types.StringPointerValue(image.Volumes.Get()),
			Checksum:             types.StringPointerValue(image.Checksum.Get()),
			ContainerFormat:      types.StringValue(image.ContainerFormat),
			DiskFormat:           types.StringValue(image.DiskFormat),
			File:                 types.StringValue(image.File),
			Id:                   types.StringValue(image.Id),
			MinDisk:              types.Int32Value(image.MinDisk),
			MinRam:               types.Int32Value(image.MinRam),
			Name:                 types.StringValue(image.Name),
			OsDistro:             types.StringPointerValue(image.OsDistro.Get()),
			OsHashAlgo:           types.StringPointerValue(image.OsHashAlgo.Get()),
			OsHashValue:          types.StringPointerValue(image.OsHashValue.Get()),
			OsHidden:             types.BoolValue(image.OsHidden),
			Owner:                types.StringValue(image.Owner),
			OwnerAccountName:     types.StringPointerValue(image.OwnerAccountName.Get()),
			OwnerUserName:        types.StringPointerValue(image.OwnerUserName.Get()),
			Protected:            types.BoolValue(image.Protected),
			RootDeviceName:       types.StringPointerValue(image.RootDeviceName.Get()),
			ScpImageType:         types.StringPointerValue(image.ScpImageType.Get()),
			ScpK8sVersion:        types.StringPointerValue(image.ScpK8sVersion.Get()),
			ScpOriginalImageType: types.StringPointerValue(image.ScpOriginalImageType.Get()),
			ScpOsVersion:         types.StringPointerValue(image.ScpOsVersion.Get()),
			Size:                 types.Int64PointerValue(image.Size.Get()),
			Status:               types.StringValue(image.Status),
			VirtualSize:          types.Int64PointerValue(image.VirtualSize.Get()),
			Visibility:           types.StringValue(image.Visibility),
			Url:                  types.StringPointerValue(image.Url.Get()),
			CreatedAt:            types.StringValue(image.CreatedAt),
			UpdatedAt:            types.StringValue(image.UpdatedAt),
		}
		imageObjectValue, _ := types.ObjectValueFrom(ctx, imageModel.AttributeTypes(), imageModel)
		state.Image = imageObjectValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
