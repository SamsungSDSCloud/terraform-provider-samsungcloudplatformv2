package multinodegpucluster

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GpunodeProductDataSources{}
	_ datasource.DataSourceWithConfigure = &GpunodeProductDataSources{}
)

// NewGpunodeProductDataSources is a helper function to simplify the provider implementation.
func NewGpunodeProductDataSources() datasource.DataSource {
	return &GpunodeProductDataSources{}
}

// GpunodeProductDataSources is the data source implementation.
type GpunodeProductDataSources struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *GpunodeProductDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunode_products"
}

// Configure adds the provider configured client to the data source.
func (d *GpunodeProductDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expect *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = inst.Client.Mngc
	d.clients = inst.Client
}

// Schema defines the schema for the data source.
func (d *GpunodeProductDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GpuNodeProductsDataSourcesSchema()
}

func GpuNodeProductsDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "List of GPU Node Product",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional:            true,
				Description:         "Product type\n  - example: SCALE",
				MarkdownDescription: "Product type\n  - example: SCALE",
			},
			"image_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
			},
			"products": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "GPU Node Product List",
				MarkdownDescription: "GPU Node Product List",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Product ID\n  - example: 10a599e031e749b7b260868f441e862b",
							MarkdownDescription: "Product ID\n  - example: 10a599e031e749b7b260868f441e862b",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Product name\n  - example: g3c128b8_metal",
							MarkdownDescription: "Product name\n  - example: g3c128b8_metal",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description\n  - example: CPU 128 | Memory 4096GB | B300(268GB)*8 | Disk 960GB",
							MarkdownDescription: "Description\n  - example: CPU 128 | Memory 4096GB | B300(268GB)*8 | Disk 960GB",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Product type\n  - example: STANDARD",
							MarkdownDescription: "Product type\n  - example: STANDARD",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "State\n  - example: ACTIVE",
							MarkdownDescription: "State\n  - example: ACTIVE",
						},
						"product_attrs": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Product attributes",
							MarkdownDescription: "Product attributes",
							Attributes: map[string]schema.Attribute{
								"compute_class_type_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Compute class type name\n  - example: GPU B300 Metal-3",
									MarkdownDescription: "Compute class type name\n  - example: GPU B300 Metal-3",
								},
								"compute_class_type_value": schema.StringAttribute{
									Computed:            true,
									Description:         "Compute class type unit\n  - example: g3_metal",
									MarkdownDescription: "Compute class type unit\n  - example: g3_metal",
								},
								"cpu_value": schema.StringAttribute{
									Computed:            true,
									Description:         "CPU value\n  - example: 128",
									MarkdownDescription: "CPU value\n  - example: 128",
								},
								"memory_value": schema.StringAttribute{
									Computed:            true,
									Description:         "Memory value\n  - example: 4096",
									MarkdownDescription: "Memory value\n  - example: 4096",
								},
								"disk_value": schema.StringAttribute{
									Computed:            true,
									Description:         "Disk value\n  - example: 960",
									MarkdownDescription: "Disk value\n  - example: 960",
								},
								"disk_unit": schema.StringAttribute{
									Computed:            true,
									Description:         "Disk unit\n  - example: GB",
									MarkdownDescription: "Disk unit\n  - example: GB",
								},
								"gpu_model": schema.StringAttribute{
									Computed:            true,
									Description:         "GPU Model\n  - example: B300",
									MarkdownDescription: "GPU Model\n  - example: B300",
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created at\n  - example: 2024-05-01T00:00:00Z",
							MarkdownDescription: "Created at\n  - example: 2024-05-01T00:00:00Z",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Created by\n  - example: user",
							MarkdownDescription: "Created by\n  - example: user",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified at\n  - example: 2024-05-01T00:00:00Z",
							MarkdownDescription: "Modified at\n  - example: 2024-05-01T00:00:00Z",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified by\n  - example: user",
							MarkdownDescription: "Modified by\n  - example: user",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GpunodeProductDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.GpuNodeProductList

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetProductList(ctx, state.Type.ValueString(), state.ImageId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GPU Node Products",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, product := range data.Products {
		productState := multinodegpuclusterClient.GpuNodeProduct{
			CreatedAt:   types.StringValue(product.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(product.CreatedBy),
			Description: types.StringValue(product.Description),
			Id:          types.StringValue(product.Id),
			ModifiedAt:  types.StringValue(product.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(product.ModifiedBy),
			Name:        types.StringValue(product.Name),
			ProductAttrs: multinodegpuclusterClient.GpuNodeProductAttrs{
				ComputeClassTypeName:  types.StringValue(product.ProductAttrs.ComputeClassTypeName),
				ComputeClassTypeValue: types.StringValue(product.ProductAttrs.ComputeClassTypeValue),
				CpuValue:              types.StringValue(product.ProductAttrs.CpuValue),
				DiskUnit:              types.StringValue(product.ProductAttrs.DiskUnit),
				DiskValue:             types.StringValue(product.ProductAttrs.DiskValue),
				GpuModel:              types.StringValue(product.ProductAttrs.GpuModel),
				MemoryValue:           types.StringValue(product.ProductAttrs.MemoryValue),
			},
			State: types.StringValue(product.State),
			Type:  types.StringValue(product.Type),
		}
		state.Products = append(state.Products, productState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
