package ske

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/service/ske/converter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeNodepoolResource{}
	_ resource.ResourceWithConfigure   = &skeNodepoolResource{}
	_ resource.ResourceWithImportState = &skeClusterResource{}
	_ resource.ResourceWithModifyPlan  = &skeNodepoolResource{}
)

func NewSkeNodepoolResource() resource.Resource {
	return &skeNodepoolResource{}
}

type skeNodepoolResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *skeNodepoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepool" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *skeNodepoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "nodepool",
		Attributes: map[string]schema.Attribute{
			"advanced_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"allowed_unsafe_sysctls": schema.StringAttribute{
						Optional:            true,
						Description:         "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
						MarkdownDescription: "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
					},
					"container_log_max_files": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
						MarkdownDescription: "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
						Validators: []validator.Int32{
							int32validator.Between(2, 10),
						},
					},
					"container_log_max_size": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
						MarkdownDescription: "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
						Validators: []validator.Int32{
							int32validator.Between(10, 100),
						},
					},
					"image_gc_high_threshold": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
						MarkdownDescription: "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
						Validators: []validator.Int32{
							int32validator.Between(10, 85),
						},
					},
					"image_gc_low_threshold": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
						MarkdownDescription: "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
						Validators: []validator.Int32{
							int32validator.Between(10, 85),
						},
					},
					"max_pods": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
						MarkdownDescription: "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
						Validators: []validator.Int32{
							int32validator.Between(10, 250),
						},
					},
					"pod_max_pids": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
						MarkdownDescription: "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
						Validators: []validator.Int32{
							int32validator.Between(1024, 4194304),
						},
					},
				},
				Optional:            true,
				Description:         "Node Pool Advanced Settings\n  - example: {max_pods: 110, image_gc_high_threshold: 85, image_gc_low_threshold: 80, container_log_max_size: 10, container_log_max_files: 5, pod_max_pids: 4096, allowed_unsafe_sysctls: 'kernel.msg*'}",
				MarkdownDescription: "Node Pool Advanced Settings\n  - example: {max_pods: 110, image_gc_high_threshold: 85, image_gc_low_threshold: 80, container_log_max_size: 10, container_log_max_files: 5, pod_max_pids: 4096, allowed_unsafe_sysctls: 'kernel.msg*'}",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
			},
			"custom_image_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Custom Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Custom Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
			},
			"desired_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Desired node count (is_auto_scale = false)\n  - example: 2",
				MarkdownDescription: "Desired node count (is_auto_scale = false)\n  - example: 2",
			},
			"image_os": schema.StringAttribute{
				Required:            true,
				Description:         "Image OS\n  - pattern: ubuntu|rhel\n  - example: ubuntu",
				MarkdownDescription: "Image OS\n  - pattern: ubuntu|rhel\n  - example: ubuntu",
			},
			"image_os_version": schema.StringAttribute{
				Required:            true,
				Description:         "Image OS Version\n  - pattern(for ubuntu): 22.04\n  - pattern(for rhel): 8.10|9.4 \n  - example: 22.04",
				MarkdownDescription: "Image OS Version\n  - pattern(for ubuntu): 22.04\n  - pattern(for rhel): 8.10|9.4 \n  - example: 22.04",
			},
			"is_auto_recovery": schema.BoolAttribute{
				Required:            true,
				Description:         "Is Auto Recovery\n  - example: true",
				MarkdownDescription: "Is Auto Recovery\n  - example: true",
			},
			"is_auto_scale": schema.BoolAttribute{
				Required:            true,
				Description:         "Is Auto Scale\n  - example: true",
				MarkdownDescription: "Is Auto Scale\n  - example: true",
			},
			"keypair_name": schema.StringAttribute{
				Required:            true,
				Description:         "Keypair Name for create node\n  - example: test_keypair",
				MarkdownDescription: "Keypair Name for create node\n  - example: test_keypair",
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "Kubernetes Version\n  - example: v1.29.8",
				MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
			},
			"labels": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:            true,
							Description:         "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							MarkdownDescription: "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$"), ""),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							MarkdownDescription: "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							Validators: []validator.String{
								stringvalidator.LengthAtMost(63),
								stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
							},
							Default: stringdefault.StaticString(""),
						},
					},
				},
				Optional: true,
				Description: "Node Pool Labels\n" +
					"  - example: {key='test', value='test'}",
				MarkdownDescription: "Node Pool Labels\n" +
					"  - example: {key='test', value='test'}",
			},
			"linked_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource ID\n  - example: res-12345678",
							MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Name\n  - example: my-resource",
							MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
							MarkdownDescription: "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
						},
					},
				},
				Optional: true,
				Description: "Linked Resources\n" +
					"  - example: {id='res-12345678', name='my-resource', type='fs'}",
				MarkdownDescription: "Linked Resources\n" +
					"  - example: {id='res-12345678', name='my-resource', type='fs'}",
			},
			"max_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Maximum node count (is_auto_scale = true)\n  - example: 5",
				MarkdownDescription: "Maximum node count (is_auto_scale = true)\n  - example: 5",
			},
			"min_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Minimum node count (is_auto_scale = true)\n  - example: 1",
				MarkdownDescription: "Minimum node count (is_auto_scale = true)\n  - example: 1",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Nodepool Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-nodepool",
				MarkdownDescription: "Nodepool Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-nodepool",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9\\-]*[a-z0-9]$"), ""),
				},
			},
			"nodepool": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"advanced_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"allowed_unsafe_sysctls": schema.StringAttribute{
								Computed:            true,
								Description:         "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
								MarkdownDescription: "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
								Default:             stringdefault.StaticString(""),
							},
							"container_log_max_files": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
								MarkdownDescription: "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
							},
							"container_log_max_size": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
								MarkdownDescription: "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
							},
							"image_gc_high_threshold": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
								MarkdownDescription: "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
							},
							"image_gc_low_threshold": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
								MarkdownDescription: "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
							},
							"max_pods": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
								MarkdownDescription: "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
							},
							"pod_max_pids": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
								MarkdownDescription: "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
							},
						},
						Computed:            true,
						Description:         "Node Pool Advanced Settings\n  - example: {max_pods: 110, image_gc_high_threshold: 85, image_gc_low_threshold: 80, container_log_max_size: 10, container_log_max_files: 5, pod_max_pids: 4096, allowed_unsafe_sysctls: 'kernel.msg*'}",
						MarkdownDescription: "Node Pool Advanced Settings\n  - example: {max_pods: 110, image_gc_high_threshold: 85, image_gc_low_threshold: 80, container_log_max_size: 10, container_log_max_files: 5, pod_max_pids: 4096, allowed_unsafe_sysctls: 'kernel.msg*'}",
					},
					"auto_recovery_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Recovery\n  - example: true",
						MarkdownDescription: "Is Auto Recovery\n  - example: true",
					},
					"auto_scale_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Scale\n  - example: true",
						MarkdownDescription: "Is Auto Scale\n  - example: true",
					},
					"cluster": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Cluster",
						MarkdownDescription: "Cluster",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"current_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Current Node Count\n  - example: 1",
						MarkdownDescription: "Current Node Count\n  - example: 1",
					},
					"desired_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Desired Node Count\n  - example: 2",
						MarkdownDescription: "Desired Node Count\n  - example: 2",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
						MarkdownDescription: "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
					},
					"image": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"custom_image_name": schema.StringAttribute{
								Computed:            true,
								Description:         "Custom Image Name\n  - example: custom-image",
								MarkdownDescription: "Custom Image Name\n  - example: custom-image",
							},
							"os": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS\n  - example: ubuntu",
								MarkdownDescription: "Image OS\n  - example: ubuntu",
							},
							"os_version": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS Version\n  - example: 22.04",
								MarkdownDescription: "Image OS Version\n  - example: 22.04",
							},
							"scp_gpu_driver": schema.StringAttribute{
								Computed: true,
							},
						},
						Computed: true,
						Description: "Image\n" +
							"  - example: {custom_image_name='res-12345678', os='my-resource', os_version='fs', scp_gpu_driver='ND_535.183.06'}",
						MarkdownDescription: "Image\n" +
							"  - example: {custom_image_name='res-12345678', os='my-resource', os_version='fs', scp_gpu_driver='ND_535.183.06'}",
					},
					"keypair": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Keypair Name\n  - example: test_keypair",
								MarkdownDescription: "Keypair Name\n  - example: test_keypair",
							},
						},
						Computed: true,
						Description: "Keypair Name\n" +
							"  - example: {name='test_keypair'}",
						MarkdownDescription: "Keypair Name\n" +
							"  - example: {name='test_keypair'}",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed:            true,
						Description:         "Kubernetes Version\n  - example: v1.29.8",
						MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
					},
					"labels": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
								},
							},
						},
						Optional: true,
						Computed: true,
						Description: "Node Pool Labels\n" +
							"  - example: {key='test', value='test'}",
						MarkdownDescription: "Node Pool Labels\n" +
							"  - example: {key='test', value='test'}",
					},
					"linked_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Linked Resource ID\n  - example: res-12345678",
									MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
								},
								"name": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Linked Resource Name\n  - example: my-resource",
									MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
								},
								"type": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Linked Resource Type (fs/obs)\n  - example: fs",
									MarkdownDescription: "Linked Resource Type (fs/obs)\n  - example: fs",
								},
							},
						},
						Optional: true,
						Computed: true,
						Description: "Linked Resources\n" +
							"  - example: {id='res-12345678', name='my-resource', type='fs'}",
						MarkdownDescription: "Linked Resources\n" +
							"  - example: {id='res-12345678', name='my-resource', type='fs'}",
					},
					"max_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Max Node Count\n  - example: 5",
						MarkdownDescription: "Max Node Count\n  - example: 5",
					},
					"min_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Min Node Count\n  - example: 1",
						MarkdownDescription: "Min Node Count\n  - example: 1",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Nodepool Name\n  - example: sample-nodepool",
						MarkdownDescription: "Nodepool Name\n  - example: sample-nodepool",
					},
					"server_group_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
						MarkdownDescription: "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
					},
					"server_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"description": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type Description\n  - example: Standard",
								MarkdownDescription: "Server Type Description\n  - example: Standard",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Server Type",
						MarkdownDescription: "Server Type",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Nodepool Status\n  - pattern: RUNNING|CREATING|SCALINGUP|SCALINGDOWN|DELETING\n  - example: RUNNING",
						MarkdownDescription: "Nodepool Status\n  - pattern: RUNNING|CREATING|SCALINGUP|SCALINGDOWN|DELETING\n  - example: RUNNING",
					},
					"taints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]\n  - example: NoSchedule",
									MarkdownDescription: "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]\n  - example: NoSchedule",
								},
								"key": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
								},
							},
						},
						Optional: true,
						Computed: true,
						Description: "Node Pool Taints\n" +
							"  - example: {effect='NoSchedule', key='example.com/my-app', value='bar'}",
						MarkdownDescription: "Node Pool Taints\n" +
							"  - example: {effect='NoSchedule', key='example.com/my-app', value='bar'}",
					},
					"volume_max_iops": schema.Int32Attribute{
						Computed: true,
					},
					"volume_max_throughput": schema.Int32Attribute{
						Computed: true,
					},
					"volume_size": schema.Int32Attribute{
						Computed:            true,
						Description:         "Volume Size\n  - example: 104",
						MarkdownDescription: "Volume Size\n  - example: 104",
					},
					"volume_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt": schema.BoolAttribute{
								Computed:            true,
								Description:         "Volume Type Encrypt\n  - example: true",
								MarkdownDescription: "Volume Type Encrypt\n  - example: true",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type Name\n  - pattern: SSD|SSD_KMS|HDD|HDD_KMS|SSD_Provisioned\n  - example: SSD",
								MarkdownDescription: "Volume Type Name\n  - pattern: SSD|SSD_KMS|HDD|HDD_KMS|SSD_Provisioned\n  - example: SSD",
							},
						},
						Computed:            true,
						Description:         "Volume Type",
						MarkdownDescription: "Volume Type",
					},
					"preferred_ips": schema.StringAttribute{
						Computed:            true,
						Description:         "Preferred IPs\n  - example: 192.168.0.0-192.168.0.255,192.168.99.0",
						MarkdownDescription: "Preferred IPs\n  - example: 192.168.0.0-192.168.0.255,192.168.99.0",
					},
					"subnet_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
						MarkdownDescription: "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
					},
					"zone": schema.StringAttribute{
						Computed:            true,
						Description:         "Availability zone\n  - example: kr-west1-a",
						MarkdownDescription: "Availability zone\n  - example: kr-west1-a",
					},
				},
				Computed:            true,
				Description:         "Nodepool\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_nodepool#nested-schema-for-nodepool",
				MarkdownDescription: "Nodepool\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_nodepool#nested-schema-for-nodepool",
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scp_gpu_driver": schema.StringAttribute{
				Optional:            true,
				Description:         "Gpu Driver Version For Ubuntu\n  - pattern: 535.183.06\n  - example: 535.183.06",
				MarkdownDescription: "Gpu Driver Version For Ubuntu\n  - pattern: 535.183.06\n  - example: 535.183.06",
			},
			"server_group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Server Group ID For \n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
				MarkdownDescription: "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
			},
			"server_type_id": schema.StringAttribute{
				Required:            true,
				Description:         "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
			},
			"taints": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"effect": schema.StringAttribute{
							Required:            true,
							Description:         "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
							MarkdownDescription: "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"NoSchedule",
									"NoExecute",
									"PreferNoSchedule",
								),
							},
						},
						"key": schema.StringAttribute{
							Required:            true,
							Description:         "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							MarkdownDescription: "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$"), ""),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							MarkdownDescription: "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							Validators: []validator.String{
								stringvalidator.LengthAtMost(63),
								stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
							},
							Default: stringdefault.StaticString(""),
						},
					},
				},
				Optional: true,
				Description: "Node Pool Taints\n" +
					"  - example: {effect='NoSchedule', key='example.com/my-app', value='bar'}",
				MarkdownDescription: "Node Pool Taints\n" +
					"  - example: {effect='NoSchedule', key='example.com/my-app', value='bar'}",
			},
			"volume_max_iops": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(5000, 20000),
				},
				Description:         "Volume Max Iops\n  - example: 5000",
				MarkdownDescription: "Volume Max iops\n  - example: 5000",
			},
			"volume_max_throughput": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(250, 1000),
				},
				Description:         "Volume Type Name\n  - example: 250",
				MarkdownDescription: "Volume Type Name\n  - example: 250",
			},
			"volume_size": schema.Int32Attribute{
				Required:            true,
				Description:         "Volume Size\n  - example: 104",
				MarkdownDescription: "Volume Size\n  - example: 104",
			},
			"volume_type_name": schema.StringAttribute{
				Required:            true,
				Description:         "Volume Type Name\n  - pattern: SSD|SSD_KMS|HDD|HDD_KMS|SSD_Provisioned\n  - example: SSD",
				MarkdownDescription: "Volume Type Name\n  - pattern: SSD|SSD_KMS|HDD|HDD_KMS|SSD_Provisioned\n  - example: SSD",
			},
			"preferred_ips": schema.StringAttribute{
				Optional:            true,
				Description:         "Preferred IPs\n  - example: 192.168.0.0-192.168.0.255,192.168.99.0",
				MarkdownDescription: "Preferred IPs\n  - example: 192.168.0.0-192.168.0.255,192.168.99.0",
			},
			"subnet_id": schema.StringAttribute{
				Required:            true,
				Description:         "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
				MarkdownDescription: "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
			},
			"zone": schema.StringAttribute{
				Required:            true,
				Description:         "Availability zone\n  - example: kr-west1-a",
				MarkdownDescription: "Availability zone\n  - example: kr-west1-a",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *skeNodepoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Ske
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *skeNodepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.NodepoolResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new nodepool
	data, err := r.client.CreateNodepool(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Nodepool",
			"Could not create nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.GetNodepool().Id)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	fillUnconfiguredNodeCounts(&plan, data.GetNodepool().MinNodeCount, data.GetNodepool().DesiredNodeCount)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.GetNodepool().Id, []string{"ScalingUp"}, []string{"Running"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Nodepool",
			"Error waiting for nodepool to become running: "+err.Error(),
		)
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	resp.State = readResp.State
}

func (r *skeNodepoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ske.NodepoolResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Nodepool
	data, httpStatus, err := r.client.GetNodepool(ctx, state.Id.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := converter.NodepoolResponseToNodepoolModel(data)

	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Nodepool = nodepoolObjectValue

	/* import state (determined by state.Name) root level setting start*/
	if state.Name.IsNull() {
		state.Name = nodepoolModel.Name
		state.AdvancedSettings = nodepoolModel.AdvancedSettings
		state.CustomImageId = nodepoolModel.Image.CustomImageName
		state.DesiredNodeCount = nodepoolModel.DesiredNodeCount
		state.ImageOs = nodepoolModel.Image.Os
		state.ImageOsVersion = nodepoolModel.Image.OsVersion
		state.IsAutoScale = nodepoolModel.AutoScaleEnabled
		state.IsAutoRecovery = nodepoolModel.AutoRecoveryEnabled
		state.KeypairName = nodepoolModel.Keypair.Name
		state.KubernetesVersion = nodepoolModel.KubernetesVersion
		state.Labels = nodepoolModel.Labels
		state.LinkedResources = nodepoolModel.LinkedResources
		state.MaxNodeCount = nodepoolModel.MaxNodeCount
		state.MinNodeCount = nodepoolModel.MinNodeCount
		state.ScpGpuDriver = nodepoolModel.Image.ScpGpuDriver
		state.ServerGroupId = nodepoolModel.ServerGroupId
		state.ServerTypeId = nodepoolModel.ServerType.Id
		state.Taints = nodepoolModel.Taints
		state.VolumeSize = nodepoolModel.VolumeSize
		state.VolumeTypeName = nodepoolModel.VolumeType.Name
		state.PreferredIps = nodepoolModel.PreferredIps
		state.SubnetId = nodepoolModel.SubnetId
		state.Zone = nodepoolModel.Zone
	} else {
		refreshNodepoolRootFields(&state, nodepoolModel)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *skeNodepoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state ske.NodepoolResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	if !state.KubernetesVersion.Equal(plan.KubernetesVersion) { // 1. 노드 풀 업그레이드
		/*
			업그레이드 버전 이외 설정 값이 변경 되어서는 안됨 (자동확장, 자동복구, 노드 수)
			자동 확장 사용시 최대/최소 노드 수 변경 불가
			자동 확장 미 사용시 노드 수 변경 불가
		*/
		if state.IsAutoRecovery.Equal(plan.IsAutoRecovery) && state.IsAutoScale.Equal(plan.IsAutoScale) {
			if plan.IsAutoScale.ValueBool() && state.MinNodeCount.Equal(plan.MinNodeCount) && state.MaxNodeCount.Equal(plan.MaxNodeCount) ||
				!plan.IsAutoScale.ValueBool() && state.DesiredNodeCount.Equal(plan.DesiredNodeCount) {
				_, err := r.client.UpgradeNodepool(ctx, plan)
				if err != nil {
					detail := client.GetDetailFromError(err)
					resp.Diagnostics.AddError(
						"Error Upgrade Nodepool",
						"Could not upgrade nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
					)
					return
				}

				err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"Updating"}, []string{"Running"}, true)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Updating Nodepool",
						"Error waiting for nodepool to become running: "+err.Error(),
					)
					return
				}
			} else {
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Version",
					"When nodepool version update, must not modify node count",
				)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Error Updating Nodepool Version",
				"When nodepool version update, must not modify auto recovery and auto scale",
			)
			return
		}
	} else { // label / taint 수정
		if !reflect.DeepEqual(plan.Labels, state.Labels) {
			_, err := r.client.UpdateNodepoolLabels(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Labels",
					"Could not update nodepool labels, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if !reflect.DeepEqual(plan.Taints, state.Taints) {
			_, err := r.client.UpdateNodepoolTaints(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Taints",
					"Could not update nodepool taints, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if !reflect.DeepEqual(plan.LinkedResources, state.LinkedResources) {
			_, err := r.client.UpdateNodepoolLinkedResources(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Taints",
					"Could not update nodepool taints, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if !state.PreferredIps.Equal(plan.PreferredIps) {
			preferredIps := plan.PreferredIps.ValueString()
			_, err := r.client.SetNodepoolPreferredIps(ctx, plan.Id.ValueString(), preferredIps)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool PreferredIps",
					"Could not update nodepool preferred ips, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if (plan.IsAutoScale.ValueBool() && (!state.MinNodeCount.Equal(plan.MinNodeCount) || !state.MaxNodeCount.Equal(plan.MaxNodeCount))) ||
			(!plan.IsAutoScale.ValueBool() && !state.DesiredNodeCount.Equal(plan.DesiredNodeCount)) ||
			!state.IsAutoScale.Equal(plan.IsAutoScale) ||
			!state.IsAutoRecovery.Equal(plan.IsAutoRecovery) {
			_, err := r.client.UpdateNodepool(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool",
					"Could not update nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"ScalingUp", "ScalingDown"}, []string{"Running"}, true)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Nodepool",
					"Error waiting for nodepool to become running: "+err.Error(),
				)
			}
		}
	}

	// Get refreshed value from Nodepool
	data, _, err := r.client.GetNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read Nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := converter.NodepoolResponseToNodepoolModel(data)
	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Nodepool = nodepoolObjectValue
	fillUnconfiguredNodeCounts(&plan, data.GetNodepool().MinNodeCount, data.GetNodepool().DesiredNodeCount)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	resp.State = readResp.State
}

// refreshNodepoolRootFields copies server values into root attributes so that
// changes made outside Terraform show up in the plan. Attributes whose server
// value does not map back to the configuration (custom_image_id is returned as a
// name, advanced_settings and scp_gpu_driver get server defaults) are left as is.
func refreshNodepoolRootFields(state *ske.NodepoolResource, m ske.Nodepool) {
	state.Name = m.Name
	state.KubernetesVersion = m.KubernetesVersion
	state.ImageOs = m.Image.Os
	state.ImageOsVersion = m.Image.OsVersion
	state.IsAutoScale = m.AutoScaleEnabled
	state.IsAutoRecovery = m.AutoRecoveryEnabled
	state.KeypairName = m.Keypair.Name
	state.ServerTypeId = m.ServerType.Id
	state.VolumeSize = m.VolumeSize
	state.VolumeTypeName = m.VolumeType.Name
	state.SubnetId = m.SubnetId
	state.Zone = m.Zone

	// Only the counts used by the scaling mode are set by the user; the others
	// are filled on apply (see fillUnconfiguredNodeCounts).
	if m.AutoScaleEnabled.ValueBool() {
		state.MinNodeCount = m.MinNodeCount
		state.MaxNodeCount = m.MaxNodeCount
	} else {
		state.DesiredNodeCount = m.DesiredNodeCount
	}

	// Keep state order when elements match, so a config order differing from the server's is not a diff.
	if !sameElements(state.Labels, m.Labels) {
		state.Labels = m.Labels
	}
	if !sameElements(state.Taints, m.Taints) {
		state.Taints = m.Taints
	}
	if !sameElements(state.LinkedResources, m.LinkedResources) {
		state.LinkedResources = m.LinkedResources
	}

	// The server returns an empty value for unset strings, which matches a null config.
	if !sameOptionalString(state.PreferredIps, m.PreferredIps) {
		state.PreferredIps = m.PreferredIps
	}
	if !sameOptionalString(state.ServerGroupId, m.ServerGroupId) {
		state.ServerGroupId = m.ServerGroupId
	}
}

// sameElements reports whether a and b contain the same elements, ignoring order.
func sameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[T]int, len(a))
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		if m[v] == 0 {
			return false
		}
		m[v]--
	}
	return true
}

// sameOptionalString treats null and empty string as the same value.
func sameOptionalString(a, b types.String) bool {
	return a.ValueString() == b.ValueString()
}

// fillUnconfiguredNodeCounts fills node counts that do not apply to the current
// scaling mode (desired for auto scale, min/max for fixed size) from the server.
// Only unknown values are filled: overwriting a value set in the configuration
// makes the apply result differ from the plan ("inconsistent result after apply").
func fillUnconfiguredNodeCounts(plan *ske.NodepoolResource, minNodeCount int32, desiredNodeCount int32) {
	if plan.IsAutoScale.ValueBool() {
		if plan.DesiredNodeCount.IsUnknown() {
			plan.DesiredNodeCount = types.Int32Value(minNodeCount)
		}
		return
	}
	if plan.MinNodeCount.IsUnknown() {
		plan.MinNodeCount = types.Int32Value(desiredNodeCount)
	}
	if plan.MaxNodeCount.IsUnknown() {
		plan.MaxNodeCount = types.Int32Value(desiredNodeCount)
	}
}

// nodepoolPlanAttributes holds the attributes ModifyPlan reads. They are read one by
// one as framework types, because the resource model uses Go slices and pointers
// that cannot hold an unknown value (e.g. labels set from another resource).
type nodepoolPlanAttributes struct {
	IsAutoScale       types.Bool
	IsAutoRecovery    types.Bool
	DesiredNodeCount  types.Int32
	MinNodeCount      types.Int32
	MaxNodeCount      types.Int32
	KubernetesVersion types.String
	PreferredIps      types.String
	Labels            types.List
	Taints            types.List
	LinkedResources   types.List
	Nodepool          types.Object
}

func getNodepoolPlanAttributes(ctx context.Context, get func(context.Context, path.Path, interface{}) diag.Diagnostics) (nodepoolPlanAttributes, diag.Diagnostics) {
	var a nodepoolPlanAttributes
	var diags diag.Diagnostics
	diags.Append(get(ctx, path.Root("is_auto_scale"), &a.IsAutoScale)...)
	diags.Append(get(ctx, path.Root("is_auto_recovery"), &a.IsAutoRecovery)...)
	diags.Append(get(ctx, path.Root("desired_node_count"), &a.DesiredNodeCount)...)
	diags.Append(get(ctx, path.Root("min_node_count"), &a.MinNodeCount)...)
	diags.Append(get(ctx, path.Root("max_node_count"), &a.MaxNodeCount)...)
	diags.Append(get(ctx, path.Root("kubernetes_version"), &a.KubernetesVersion)...)
	diags.Append(get(ctx, path.Root("preferred_ips"), &a.PreferredIps)...)
	diags.Append(get(ctx, path.Root("labels"), &a.Labels)...)
	diags.Append(get(ctx, path.Root("taints"), &a.Taints)...)
	diags.Append(get(ctx, path.Root("linked_resources"), &a.LinkedResources)...)
	diags.Append(get(ctx, path.Root("nodepool"), &a.Nodepool)...)
	return a, diags
}

// nodepoolCreateOnlyAttributes have no update API. custom_image_id is left out
// because import stores the image name in it, which never equals the configured ID.
var nodepoolCreateOnlyAttributes = []string{
	"name", "cluster_id", "image_os", "keypair_name", "server_type_id", "server_group_id",
	"volume_type_name", "volume_size", "volume_max_iops", "volume_max_throughput",
	"subnet_id", "zone", "advanced_settings",
}

// nodepoolUpgradeAttributes are sent only with the upgrade API, together with kubernetes_version.
var nodepoolUpgradeAttributes = []string{"image_os_version", "scp_gpu_driver"}

// validateNodepoolUpdate rejects changes Update cannot apply, which would otherwise
// fail after apply with "inconsistent result" because the server keeps the old value.
func validateNodepoolUpdate(ctx context.Context, req resource.ModifyPlanRequest, plan, state, config nodepoolPlanAttributes, resp *resource.ModifyPlanResponse) {
	for _, name := range nodepoolCreateOnlyAttributes {
		changed, diags := configuredAttributeChanged(ctx, req, name)
		resp.Diagnostics.Append(diags...)
		if changed {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Field changes not supported",
				fmt.Sprintf("Changing `%s` will not update the actual nodepool. To change %s, recreate the nodepool.", name, name))
		}
	}

	if plan.KubernetesVersion.Equal(state.KubernetesVersion) {
		for _, name := range nodepoolUpgradeAttributes {
			changed, diags := configuredAttributeChanged(ctx, req, name)
			resp.Diagnostics.Append(diags...)
			if changed {
				resp.Diagnostics.AddAttributeError(path.Root(name), "Field changes not supported",
					fmt.Sprintf("`%s` can only be changed together with `kubernetes_version`.", name))
			}
		}
	}

	// Switching the scaling mode sends the counts of the new mode, so they must be set.
	if !plan.IsAutoScale.IsUnknown() && !plan.IsAutoScale.Equal(state.IsAutoScale) {
		if plan.IsAutoScale.ValueBool() && (config.MinNodeCount.IsNull() || config.MaxNodeCount.IsNull()) {
			resp.Diagnostics.AddAttributeError(path.Root("is_auto_scale"), "Missing node count",
				"`min_node_count` and `max_node_count` must be set when enabling `is_auto_scale`.")
		}
		if !plan.IsAutoScale.ValueBool() && config.DesiredNodeCount.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("is_auto_scale"), "Missing node count",
				"`desired_node_count` must be set when disabling `is_auto_scale`.")
		}
	}
}

// configuredAttributeChanged reports whether a root attribute set in the configuration
// differs from a known prior state value. Null and empty strings are treated as equal.
// A null state (not filled on import) is skipped, as the prior value is not known.
// A value with any unknown part, even inside an object, is skipped until it is known.
func configuredAttributeChanged(ctx context.Context, req resource.ModifyPlanRequest, name string) (bool, diag.Diagnostics) {
	var plan, state, config attr.Value
	var diags diag.Diagnostics
	diags.Append(req.Plan.GetAttribute(ctx, path.Root(name), &plan)...)
	diags.Append(req.State.GetAttribute(ctx, path.Root(name), &state)...)
	diags.Append(req.Config.GetAttribute(ctx, path.Root(name), &config)...)
	if diags.HasError() || config.IsNull() || state.IsNull() || !isFullyKnown(ctx, config) || !isFullyKnown(ctx, plan) {
		return false, diags
	}
	if p, ok := plan.(types.String); ok {
		return !sameOptionalString(p, state.(types.String)), diags
	}
	return !plan.Equal(state), diags
}

// isFullyKnown reports whether v and all values nested in it are known.
func isFullyKnown(ctx context.Context, v attr.Value) bool {
	tfValue, err := v.ToTerraformValue(ctx)
	return err == nil && tfValue.IsFullyKnown()
}

// nodepoolNestedVolatile lists nested nodepool fields an update can change; they stay unknown.
// Every nested field must be in this list or in nodepoolNestedKeep (checked by tests).
var nodepoolNestedVolatile = []string{"status", "current_node_count", "modified_at", "modified_by"}

// nodepoolNestedKeep reports, per nested nodepool field, whether the state value is kept:
// fields fixed at creation always, fields that mirror an input while the input is unchanged.
// An unknown input never equals the state, so fields that depend on it stay unknown.
func nodepoolNestedKeep(plan, state nodepoolPlanAttributes) map[string]bool {
	versionUnchanged := plan.KubernetesVersion.Equal(state.KubernetesVersion)
	scaleUnchanged := plan.IsAutoScale.Equal(state.IsAutoScale) &&
		plan.MinNodeCount.Equal(state.MinNodeCount) &&
		plan.MaxNodeCount.Equal(state.MaxNodeCount) &&
		plan.DesiredNodeCount.Equal(state.DesiredNodeCount)
	return map[string]bool{
		"id":                    true,
		"name":                  true,
		"account_id":            true,
		"cluster":               true,
		"created_at":            true,
		"created_by":            true,
		"keypair":               true,
		"server_type":           true,
		"server_group_id":       true,
		"volume_type":           true,
		"volume_size":           true,
		"volume_max_iops":       true,
		"volume_max_throughput": true,
		"subnet_id":             true,
		"zone":                  true,
		"advanced_settings":     true,
		"kubernetes_version":    versionUnchanged,
		"image":                 versionUnchanged,
		"labels":                plan.Labels.Equal(state.Labels),
		"taints":                plan.Taints.Equal(state.Taints),
		"linked_resources":      plan.LinkedResources.Equal(state.LinkedResources),
		"preferred_ips":         plan.PreferredIps.Equal(state.PreferredIps),
		"auto_recovery_enabled": plan.IsAutoRecovery.Equal(state.IsAutoRecovery),
		"auto_scale_enabled":    plan.IsAutoScale.Equal(state.IsAutoScale),
		// desired_node_count also changes with auto scale
		"min_node_count":     scaleUnchanged,
		"max_node_count":     scaleUnchanged,
		"desired_node_count": scaleUnchanged && !plan.IsAutoScale.ValueBool(),
	}
}

// ModifyPlan keeps values an update cannot change out of "(known after apply)".
func (r *skeNodepoolResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip on create and destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	plan, diags := getNodepoolPlanAttributes(ctx, req.Plan.GetAttribute)
	resp.Diagnostics.Append(diags...)
	state, diags := getNodepoolPlanAttributes(ctx, req.State.GetAttribute)
	resp.Diagnostics.Append(diags...)
	config, diags := getNodepoolPlanAttributes(ctx, req.Config.GetAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	validateNodepoolUpdate(ctx, req, plan, state, config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Node counts not used by the scaling mode are copied from the server on apply
	// (see fillUnconfiguredNodeCounts). Keep them while their source is unchanged.
	if plan.IsAutoScale.Equal(state.IsAutoScale) {
		if plan.IsAutoScale.ValueBool() {
			if plan.MinNodeCount.Equal(state.MinNodeCount) {
				useStateIfUnconfigured(&plan.DesiredNodeCount, state.DesiredNodeCount, config.DesiredNodeCount)
			}
		} else if plan.DesiredNodeCount.Equal(state.DesiredNodeCount) {
			useStateIfUnconfigured(&plan.MinNodeCount, state.MinNodeCount, config.MinNodeCount)
			useStateIfUnconfigured(&plan.MaxNodeCount, state.MaxNodeCount, config.MaxNodeCount)
		}
	}

	if canPinNestedObject(plan.Nodepool, state.Nodepool) {
		keep := nodepoolNestedKeep(plan, state)
		nodepool, diags := pinNestedObject(ctx, state.Nodepool, keep)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Nodepool = nodepool
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("desired_node_count"), plan.DesiredNodeCount)...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("min_node_count"), plan.MinNodeCount)...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("max_node_count"), plan.MaxNodeCount)...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("nodepool"), plan.Nodepool)...)
}

func (r *skeNodepoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.NodepoolResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Resource Group
	data, err := r.client.DeleteNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Nodepool",
			"Could not delete nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.GetResourceId(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting nodepool",
			"Error waiting for nodepool to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForNodepoolStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetNodepool(ctx, id)
		if httpStatus == 200 {
			return info, info.GetNodepool().Status, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETED", nil
		} else if httpStatus == 500 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETING", nil
		} else if err != nil {
			return nil, "", err
		}

		return info, info.GetNodepool().Status, nil
	}, -1, -1, -1, -1)
}

func (r *skeNodepoolResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp) // 유니크한 스키마 입력
}
