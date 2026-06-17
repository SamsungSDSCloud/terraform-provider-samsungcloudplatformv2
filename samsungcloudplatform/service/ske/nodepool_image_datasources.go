package ske

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeNodepoolImageDataSources{}
	_ datasource.DataSourceWithConfigure = &skeNodepoolImageDataSources{}
)

func NewSkeNodepoolImageDataSources() datasource.DataSource {
	return &skeNodepoolImageDataSources{}
}

type skeNodepoolImageDataSources struct {
	config *scpsdk.Configuration
	client *ske.Client
}

// Metadata returns the data source type name.
func (d *skeNodepoolImageDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepool_images"
}

// Schema defines the schema for the data source.
func (d *skeNodepoolImageDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of nodepool images available for nodepool creation.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
				Description:         "SCP original image type (k8s, k8s_gpu)\n  - example: k8s",
				MarkdownDescription: "SCP original image type (k8s, k8s_gpu)\n  - example: k8s",
				Required:            true,
			},
			common.ToSnakeCase("Os"): schema.StringAttribute{
				Description:         "Image OS\n  - example: ubuntu",
				MarkdownDescription: "Image OS\n  - example: ubuntu",
				Optional:            true,
			},
			common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
				Description:         "Kubernetes Version\n  - example: v1.29.8",
				MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
				Optional:            true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)\n  - example: 10000",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page\n  - example: 0",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort\n  - example: created_at:desc",
				Optional:    true,
			},
			common.ToSnakeCase("NodepoolImages"): schema.ListNestedAttribute{
				Description: "A list of nodepool images.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
							MarkdownDescription: "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
							Computed:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "Image name\n  - example: ubuntu-22.04-k8s-v1.29.8",
							MarkdownDescription: "Image name\n  - example: ubuntu-22.04-k8s-v1.29.8",
							Computed:            true,
						},
						common.ToSnakeCase("Os"): schema.StringAttribute{
							Description:         "Image OS\n  - example: ubuntu",
							MarkdownDescription: "Image OS\n  - example: ubuntu",
							Computed:            true,
						},
						common.ToSnakeCase("OsVersion"): schema.StringAttribute{
							Description:         "Image OS Version\n  - example: 22.04",
							MarkdownDescription: "Image OS Version\n  - example: 22.04",
							Computed:            true,
						},
						common.ToSnakeCase("ScpGpuDriver"): schema.StringAttribute{
							Description:         "GPU Driver Version\n  - example: ND_570.195.03",
							MarkdownDescription: "GPU Driver Version\n  - example: ND_570.195.03",
							Computed:            true,
						},
						common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
							Description:         "Kubernetes Version\n  - example: v1.29.8",
							MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
							Computed:            true,
							Optional:            true,
						},
						common.ToSnakeCase("EndOfSupport"): schema.BoolAttribute{
							Description:         "Whether this is an EOS (End of Service) image\n  - example: false",
							MarkdownDescription: "Whether this is an EOS (End of Service) image\n  - example: false",
							Computed:            true,
						},
						common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
							Description:         "SCP image type (k8s, custom)\n  - example: k8s",
							MarkdownDescription: "SCP image type (k8s, custom)\n  - example: k8s",
							Computed:            true,
							Optional:            true,
						},
						common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
							Description:         "SCP original image type (k8s, k8s_gpu)\n  - example: k8s",
							MarkdownDescription: "SCP original image type (k8s, k8s_gpu)\n  - example: k8s",
							Computed:            true,
							Optional:            true,
						},
						common.ToSnakeCase("ScpSupportedClassTypes"): schema.ListAttribute{
							ElementType:         types.StringType,
							Description:         "List of supported class types for GPU SKE image\n  - example: [GPU-A100-1]",
							MarkdownDescription: "List of supported class types for GPU SKE image\n  - example: [GPU-A100-1]",
							Computed:            true,
						},
						common.ToSnakeCase("Volume"): schema.SingleNestedAttribute{
							Description:         "Volume information\n  - example: map[size:100]",
							MarkdownDescription: "Volume information\n  - example: map[size:100]",
							Computed:            true,
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Size"): schema.Int64Attribute{
									Description:         "Volume Size\n  - example: 100",
									MarkdownDescription: "Volume Size\n  - example: 100",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *skeNodepoolImageDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.Ske
}

// Read refreshes the Terraform state with the latest data.
func (d *skeNodepoolImageDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.NodepoolImageDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNodepoolImageList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Nodepool Images",
			err.Error(),
		)
		return
	}

	contents := data.NodepoolImages

	// Map response body to model
	for _, nodepoolImage := range contents {
		nodepoolImageState := ske.NodepoolImageSummary{
			Id:                     types.StringValue(nodepoolImage.Id),
			Name:                   types.StringValue(nodepoolImage.Name),
			Os:                     types.StringValue(nodepoolImage.Os),
			OsVersion:              types.StringValue(nodepoolImage.OsVersion),
			KubernetesVersion:      types.StringPointerValue(nodepoolImage.KubernetesVersion.Get()),
			EndOfSupport:           types.BoolPointerValue(nodepoolImage.EndOfSupport),
			ScpImageType:           types.StringPointerValue(nodepoolImage.ScpImageType.Get()),
			ScpOriginalImageType:   types.StringPointerValue(nodepoolImage.ScpOriginalImageType.Get()),
			Volume:                 d.makeVolume(nodepoolImage.Volume),
			ScpGpuDriver:           types.StringPointerValue(nodepoolImage.ScpGpuDriver.Get()),
			ScpSupportedClassTypes: d.makeScpSupportedClassType(nodepoolImage.ScpSupportedClassTypes),
		}
		state.NodepoolImages = append(state.NodepoolImages, nodepoolImageState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *skeNodepoolImageDataSources) makeVolume(volume map[string]interface{}) *ske.NodepoolImageVolume {
	volume_size := volume["size"]
	if volume_size == nil {
		return &ske.NodepoolImageVolume{
			Size: types.Int64PointerValue(nil),
		}
	}
	return &ske.NodepoolImageVolume{
		Size: types.Int64PointerValue(volume_size.(*int64)),
	}
}

func (d *skeNodepoolImageDataSources) makeScpSupportedClassType(scpSupportedClassTypes []interface{}) []types.String {
	var scpSupportedClassTypesModel []types.String
	for _, scpSupportedClassType := range scpSupportedClassTypes {
		scpSupportedClassTypesModel = append(scpSupportedClassTypesModel, types.StringPointerValue(scpSupportedClassType.(*string)))
	}
	return scpSupportedClassTypesModel
}
