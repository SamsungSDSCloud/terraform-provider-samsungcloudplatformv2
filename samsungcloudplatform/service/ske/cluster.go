package ske

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/ske/1.6"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeClusterResource{}
	_ resource.ResourceWithConfigure   = &skeClusterResource{}
	_ resource.ResourceWithImportState = &skeClusterResource{}
)

// linkedResourceObjectType is the tfsdk object type for the linked_resources nested block.
var linkedResourceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	},
}

// Example ID used in schema descriptions for nested resource objects.
const exampleNestedResourceID = "  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f'}"

// listStringsEqualSet checks if two types.List of strings contain the same elements, ignoring order.
func listStringsEqualSet(a, b types.List) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	if len(a.Elements()) != len(b.Elements()) {
		return false
	}
	var aSlice, bSlice []string
	_ = a.ElementsAs(context.Background(), &aSlice, false)
	_ = b.ElementsAs(context.Background(), &bSlice, false)
	return stringSliceEqualSet(aSlice, bSlice)
}

// stringSliceEqualSet checks if two string slices contain the same elements, ignoring order.
func stringSliceEqualSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		m[v]--
		if m[v] == 0 {
			delete(m, v)
		}
	}
	return len(m) == 0
}

// peacrListEqualSet checks if two types.List of PrivateEndpointAccessControlResource contain the same elements, ignoring order.
func peacrListEqualSet(ctx context.Context, a, b types.List) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	if len(a.Elements()) != len(b.Elements()) {
		return false
	}
	var aSlice, bSlice []ske.PrivateEndpointAccessControlResource
	_ = a.ElementsAs(ctx, &aSlice, false)
	_ = b.ElementsAs(ctx, &bSlice, false)
	return peacrSliceEqualSet(aSlice, bSlice)
}

// peacrSliceEqualSet checks if two PrivateEndpointAccessControlResource slices contain the same elements, ignoring order.
func peacrSliceEqualSet(a, b []ske.PrivateEndpointAccessControlResource) bool {
	if len(a) != len(b) {
		return false
	}
	key := func(r ske.PrivateEndpointAccessControlResource) string {
		return r.Id.ValueString() + "\x1f" + r.Name.ValueString() + "\x1f" + r.Type.ValueString()
	}
	m := make(map[string]int, len(a))
	for _, v := range a {
		m[key(v)]++
	}
	for _, v := range b {
		k := key(v)
		m[k]--
		if m[k] == 0 {
			delete(m, k)
		}
	}
	return len(m) == 0
}

// linkedResourceListEqualSet checks if two types.List of LinkedResource contain the same elements, ignoring order.
func linkedResourceListEqualSet(ctx context.Context, a, b types.List) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	if len(a.Elements()) != len(b.Elements()) {
		return false
	}
	var aSlice, bSlice []ske.LinkedResource
	_ = a.ElementsAs(ctx, &aSlice, false)
	_ = b.ElementsAs(ctx, &bSlice, false)
	return linkedResourceSliceEqualSet(aSlice, bSlice)
}

// linkedResourceSliceEqualSet checks if two LinkedResource slices contain the same elements, ignoring order.
func linkedResourceSliceEqualSet(a, b []ske.LinkedResource) bool {
	if len(a) != len(b) {
		return false
	}
	key := func(r ske.LinkedResource) string {
		return r.Id.ValueString() + "\x1f" + r.Name.ValueString() + "\x1f" + r.Type.ValueString()
	}
	m := make(map[string]int, len(a))
	for _, v := range a {
		m[key(v)]++
	}
	for _, v := range b {
		k := key(v)
		m[k]--
		if m[k] == 0 {
			delete(m, k)
		}
	}
	return len(m) == 0
}

// NewSkeClusterResource is a helper function to simplify the provider implementation.
func NewSkeClusterResource() resource.Resource {
	return &skeClusterResource{}
}

// skeClusterResource is the data source implementation.
type skeClusterResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *skeClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster"
}

// Schema defines the schema for the data source.
func (r *skeClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterResourceSchema()
}

func ClusterResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "cluster",
		Attributes: map[string]schema.Attribute{
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"additional_subnet_id_list": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						Description:         "Additional Subnet ID List",
						MarkdownDescription: "List of additional subnet IDs associated with the cluster.",
					},
					"cluster_namespace": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Namespace\n  - example: sample-cluster-12345",
						MarkdownDescription: "Cluster Namespace\n  - example: sample-cluster-12345",
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
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						MarkdownDescription: "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed: true,
						Description: "Cluster Version\n" +
							"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
							"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
							"  - example: v1.34.3\n" +
							"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
						MarkdownDescription: "Cluster Version\n" +
							"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
							"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
							"  - example: v1.34.3\n" +
							"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
					},
					"linked_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource ID\n  - example: res-12345678",
									MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Name\n  - example: my-resource",
									MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
									MarkdownDescription: "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
								},
							},
						},
						Computed:            true,
						Description:         "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
						MarkdownDescription: "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
					},
					"managed_security_group": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Managed Security Group\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Managed Security Group\n" +
							exampleNestedResourceID,
					},
					"max_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Max Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Max Node Count\n  - example: 5",
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
						Description:         "Cluster Name\n  - example: sample-cluster",
						MarkdownDescription: "Cluster Name\n  - example: sample-cluster",
					},
					"node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Node Count\n  - example: 5",
					},
					"private_endpoint_access_control_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Name\n  - example: sample-name",
									MarkdownDescription: "Private Endpoint Access Control Resource Name\n  - example: sample-name",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Type\n  - pattern: vm |bm|gpuvm|mngc|devops\n  - example: vm",
									MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - pattern: vm |bm|gpuvm|mngc|devops\n  - example: vm",
								},
							},
						},
						Computed: true,
						Description: "Private Endpoint Access Control Resources\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
						MarkdownDescription: "Private Endpoint Access Control Resources\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
					},
					"private_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
					},
					"private_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
						MarkdownDescription: "Private Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
					},
					"public_endpoint_access_control_ip": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
						MarkdownDescription: "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
					},
					"public_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
					},
					"public_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
						MarkdownDescription: "Public Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
					},
					"security_group_list": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
							},
						},
						Computed: true,
						Description: "Connected Security Group List\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Connected Security Group List\n" +
							exampleNestedResourceID,
					},
					"service_watch_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Service Watch Enabled\n  - example: true",
						MarkdownDescription: "Service Watch Enabled\n  - example: true",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
						MarkdownDescription: "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
					},
					"subnet": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Subnet of Cluster\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Subnet of Cluster\n" +
							exampleNestedResourceID,
					},
					"volume": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "NFS Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "NFS Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Connected File Storage\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Connected File Storage\n" +
							exampleNestedResourceID,
					},
					"vpc": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource name\n  - example: sample-name",
								MarkdownDescription: "External Resource name\n  - example: sample-name",
							},
						},
						Computed: true,
						Description: "VPC of Cluster\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name'}",
						MarkdownDescription: "VPC of Cluster\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name'}",
					},
				},
				Computed:            true,
				Description:         "Cluster\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_cluster#nested-schema-for-cluster",
				MarkdownDescription: "Cluster\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_cluster#nested-schema-for-cluster",
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Required: true,
				Description: "Cluster Version\n" +
					"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
					"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
					"  - example: v1.34.3\n" +
					"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
				MarkdownDescription: "Cluster Version\n" +
					"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
					"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
					"  - example: v1.34.3\n" +
					"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$"), ""),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Name\n  - maxLength: 30\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-cluster",
				MarkdownDescription: "Cluster Name\n  - maxLength: 30\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-cluster",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 30),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9\\-]*[a-z0-9]$"), ""),
				},
			},
			"private_endpoint_access_control_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							MarkdownDescription: "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource Name\n  - example: sample-name",
							MarkdownDescription: "Private Endpoint Access Control Resource Name\n  - example: sample-name",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource Type\n  - pattern: vm|bm|gpuvm|mngc|devops\n  - example: vm",
							MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - pattern: vm|bm|gpuvm|mngc|devops\n  - example: vm",
						},
					},
				},
				Optional: true,
				Computed: true,
				Description: "Private Endpoint Access Control Resources\n" +
					"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
				MarkdownDescription: "Private Endpoint Access Control Resources\n" +
					"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
			},
			"public_endpoint_access_control_ip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
				MarkdownDescription: "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
			},
			"security_group_id_list": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Security Group ID List\n  - example: [bdfda539-bd2e-4a5c-9021-ec6d52d1ca79]",
				MarkdownDescription: "Security Group ID List\n  - example: [bdfda539-bd2e-4a5c-9021-ec6d52d1ca79]",
			},
			"service_watch_logging_enabled": schema.BoolAttribute{
				Required:            true,
				Description:         "Service Watch Enabled\n  - example: true",
				MarkdownDescription: "Service Watch Enabled\n  - example: true",
			},
			"default_subnet_id": schema.StringAttribute{
				Required:            true,
				Description:         "Default Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
				MarkdownDescription: "Default Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
			},
			"nfs_volume_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NFS Volume (File Storage) ID\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
				MarkdownDescription: "NFS Volume (File Storage) ID\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
			},
			"vpc_id": schema.StringAttribute{
				Required:            true,
				Description:         "VPC ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "VPC ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
			},
			"additional_subnet_id_list": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Additional Subnet ID List (max 2)\n  - example: [88b9be0f9a7f4c6e8a5f0e6b9e6b6f1a]",
				MarkdownDescription: "Additional Subnet ID List (max 2)\n  - example: [88b9be0f9a7f4c6e8a5f0e6b9e6b6f1a]",
				Validators: []validator.List{
					listvalidator.SizeAtMost(2),
				},
			},
			"deletion_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cluster Deletion Protection Enabled\n  - example: true",
				MarkdownDescription: "Cluster Deletion Protection Enabled\n  - example: true",
			},
			"linked_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource ID\n  - example: 8a6c40d5-270b-4ef9-aa6b-ac19e89a71e9",
							MarkdownDescription: "Linked Resource ID\n  - example: 8a6c40d5-270b-4ef9-aa6b-ac19e89a71e9",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Name\n  - example: terraformskenfs_xjviey",
							MarkdownDescription: "Linked Resource Name\n  - example: terraformskenfs_xjviey",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Type\n  - pattern: fs|obs\n  - example: fs",
							MarkdownDescription: "Linked Resource Type\n  - pattern: fs|obs\n  - example: fs",
							Validators: []validator.String{
								stringvalidator.OneOf("fs", "obs"),
							},
						},
					},
				},
				Optional: true,
				Computed: true,
				Description: "Linked Resources\n" +
					"  - example: {id='8a6c40d5-270b-4ef9-aa6b-ac19e89a71e9', name='terraformskenfs_xjviey', type='fs'}",
				MarkdownDescription: "Linked Resources\n" +
					"  - example: {id='8a6c40d5-270b-4ef9-aa6b-ac19e89a71e9', name='terraformskenfs_xjviey', type='fs'}",
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *skeClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *skeClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{"CREATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		// Still attempt Read to populate known values (e.g. status=ERROR) before returning,
		// so Terraform state doesn't retain unknown Computed values.
		readReq := resource.ReadRequest{State: resp.State}
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, readReq, readResp)
		resp.Diagnostics.Append(readResp.Diagnostics...)
		resp.State = readResp.State
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *skeClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from cluster
	data, httpStatus, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	cluster := data.Cluster

	var securityGroups []ske.ExternalResourceId
	for _, sgId := range cluster.SecurityGroupIdList {
		securityGroups = append(securityGroups, ske.ExternalResourceId{
			Id: types.StringValue(sgId),
		})
	}

	privateEndpointAccessControlResources := make([]ske.PrivateEndpointAccessControlResource, 0, len(cluster.PrivateEndpointAccessControlResources))
	for _, privateEndpointAccessControlResource := range cluster.PrivateEndpointAccessControlResources {
		privateEndpointAccessControlResources = append(privateEndpointAccessControlResources, r.makePrivateEndpointAccessControlResourceModel((*scpske.PrivateEndpointAccessControlResource)(&privateEndpointAccessControlResource)))
	}

	clusterModel := ske.Cluster{
		Id:                                    types.StringValue(cluster.Id),
		Name:                                  types.StringValue(cluster.Name),
		AccountId:                             types.StringValue(cluster.AccountId),
		AdditionalSubnetIdList:                func() types.List { v, _ := types.ListValueFrom(ctx, types.StringType, func() []string { ids := make([]string, len(cluster.AdditionalSubnetIdList)); copy(ids, cluster.AdditionalSubnetIdList); return ids }()); return v }(),
		KubernetesVersion:                     types.StringValue(cluster.KubernetesVersion),
		ClusterNamespace:                      types.StringValue(cluster.ClusterNamespace),
		MaxNodeCount:                          types.Int32PointerValue(cluster.MaxNodeCount.Get()),
		NodeCount:                             types.Int32PointerValue(cluster.NodeCount.Get()),
		PrivateEndpointUrl:                    types.StringValue(cluster.PrivateEndpointUrl),
		PrivateKubeconfigDownloadYn:           types.StringValue(cluster.PrivateKubeconfigDownloadYn),
		PrivateEndpointAccessControlResources: privateEndpointAccessControlResources,
		PublicEndpointUrl:                     types.StringValue(cluster.GetPublicEndpointUrl()),
		PublicKubeconfigDownloadYn:            types.StringValue(cluster.PublicKubeconfigDownloadYn),
		PublicEndpointAccessControlIp:         r.nullableStringToStringValue(cluster.PublicEndpointAccessControlIp),
		Vpc:                                   r.nullableStringToExternalResourceId(cluster.VpcId),
		Subnet:                                r.nullableStringToExternalResourceId(cluster.DefaultSubnetId),
		Volume:                                r.nullableStringToExternalResourceId(cluster.NfsVolumeId),
		SecurityGroupList:                     securityGroups,
		ManagedSecurityGroup:                  r.nullableStringToExternalResourceId(cluster.ManagedSecurityGroupId),
		CreatedAt:                             types.StringValue(cluster.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                             types.StringValue(cluster.CreatedBy),
		ModifiedAt:                            types.StringValue(cluster.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                            types.StringValue(cluster.ModifiedBy),
		Status:                                types.StringValue(cluster.Status),
		ServiceWatchLoggingEnabled:            types.BoolValue(cluster.GetServiceWatchLoggingEnabled()),
		LinkedResources:                       convertLinkedResourcesFromSDK(cluster.LinkedResources),
	}
	clusterObjectValue, diags := types.ObjectValueFrom(ctx, clusterModel.AttributeTypes(), clusterModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Cluster = clusterObjectValue

	/* root level setting from API response (covers both import and refresh) */
	state.Name = types.StringValue(cluster.Name)
	state.KubernetesVersion = types.StringValue(cluster.KubernetesVersion)

	// Preserve config order: if element set matches current state, keep state value.
	peacrListValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		},
	}, privateEndpointAccessControlResources)
	resp.Diagnostics.Append(diags...)
	if !state.PrivateEndpointAccessControlResources.IsNull() && !state.PrivateEndpointAccessControlResources.IsUnknown() &&
		peacrListEqualSet(ctx, state.PrivateEndpointAccessControlResources, peacrListValue) {
		// Keep existing state to preserve config order.
	} else {
		state.PrivateEndpointAccessControlResources = peacrListValue
	}
	state.PublicEndpointAccessControlIp = r.nullableStringToStringValue(cluster.PublicEndpointAccessControlIp)

	securityGroupsStringList := make([]types.String, 0, len(cluster.SecurityGroupIdList))
	for _, sgId := range cluster.SecurityGroupIdList {
		securityGroupsStringList = append(securityGroupsStringList, types.StringValue(sgId))
	}
	securityGroupIdListValue, diags := types.ListValueFrom(ctx, types.StringType, securityGroupsStringList)
	resp.Diagnostics.Append(diags...)
	if !state.SecurityGroupIdList.IsNull() && !state.SecurityGroupIdList.IsUnknown() &&
		listStringsEqualSet(state.SecurityGroupIdList, securityGroupIdListValue) {
		// Keep existing state to preserve config order.
	} else {
		state.SecurityGroupIdList = securityGroupIdListValue
	}
	state.ServiceWatchLoggingEnabled = types.BoolValue(cluster.ServiceWatchLoggingEnabled)
	state.DefaultSubnetId = types.StringValue(cluster.GetDefaultSubnetId())
	state.NfsVolumeId = r.nullableStringToStringValue(cluster.NfsVolumeId)
	state.VpcId = types.StringValue(cluster.GetVpcId())

	// v1.6 fields
	state.DeletionProtectionEnabled = types.BoolValue(cluster.DeletionProtectionEnabled)
	additionalSubnetIdList := make([]types.String, 0, len(cluster.AdditionalSubnetIdList))
	for _, id := range cluster.AdditionalSubnetIdList {
		additionalSubnetIdList = append(additionalSubnetIdList, types.StringValue(id))
	}
	additionalSubnetIdListValue, _ := types.ListValueFrom(ctx, types.StringType, additionalSubnetIdList)
	if !state.AdditionalSubnetIdList.IsNull() && !state.AdditionalSubnetIdList.IsUnknown() &&
		listStringsEqualSet(state.AdditionalSubnetIdList, additionalSubnetIdListValue) {
		// Keep existing state to preserve config order.
	} else {
		state.AdditionalSubnetIdList = additionalSubnetIdListValue
	}
	linkedResources := make([]ske.LinkedResource, 0, len(cluster.LinkedResources))
	for _, lr := range cluster.LinkedResources {
		linkedResources = append(linkedResources, ske.LinkedResource{
			Id:   types.StringValue(lr.Id),
			Name: types.StringValue(lr.Name),
			Type: types.StringValue(lr.Type),
		})
	}
	linkedResourcesList, _ := types.ListValueFrom(ctx, linkedResourceObjectType, linkedResources)
	if !state.LinkedResources.IsNull() && !state.LinkedResources.IsUnknown() &&
		linkedResourceListEqualSet(ctx, state.LinkedResources, linkedResourcesList) {
		// Keep existing state to preserve config order.
	} else {
		state.LinkedResources = linkedResourcesList
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *skeClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state ske.ClusterResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	err := r.syncSecurityGroupList(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncKubernetesVersion(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPrivateEndpointAccessControlResources(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPublicEndpointAccessControlIp(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncServiceWatchLoggingEnabled(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncTags(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncClusterSubnets(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncNfsVolume(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncDeletionProtection(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncLinkedResources(ctx, state, &plan, resp)
	if err != nil {
		return
	}

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
	resp.State = readResp.State
}

func (r *skeClusterResource) syncSecurityGroupList(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.SecurityGroupIdList.IsUnknown() || listStringsEqualSet(plan.SecurityGroupIdList, state.SecurityGroupIdList) {
		return nil
	}
	data, err := r.client.UpdateClusterSecurityGroups(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating SecurityGroupList",
			"Could not update security group list, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.Cluster.Id)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncKubernetesVersion(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.KubernetesVersion.Equal(state.KubernetesVersion) {
		return nil
	}
	data, err := r.client.UpgradeCluster(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating KubernetesVersion",
			"Could not update kubernetes version, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPrivateEndpointAccessControlResources(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.PrivateEndpointAccessControlResources.IsUnknown() || peacrListEqualSet(ctx, plan.PrivateEndpointAccessControlResources, state.PrivateEndpointAccessControlResources) {
		return nil
	}
	data, err := r.client.UpdatePrivateEndpointAccessControlResources(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PrivateEndpointAccessControlResources",
			"Could not update cluster private endpoint access control resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPublicEndpointAccessControlIp(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.PublicEndpointAccessControlIp.IsUnknown() || plan.PublicEndpointAccessControlIp.Equal(state.PublicEndpointAccessControlIp) {
		return nil
	}
	data, err := r.client.UpdatePublicEndpointAccessControlIps(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PublicEndpointAccessControlIp",
			"Could not update public endpoint access control ip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncServiceWatchLoggingEnabled(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.ServiceWatchLoggingEnabled.Equal(state.ServiceWatchLoggingEnabled) {
		return nil
	}
	data, err := r.client.UpdateServiceWatchLoggingEnabled(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating ServiceWatchLoggingEnabled",
			"Could not update service watch logging enabled, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncTags(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.Tags.Equal(state.Tags) {
		return nil
	}
	_, err := tag.UpdateTags(r.clients, "ske", "cluster", plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Tags",
			"Could not update tags, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(plan.Id.ValueString())
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncClusterSubnets(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.AdditionalSubnetIdList.IsUnknown() || listStringsEqualSet(plan.AdditionalSubnetIdList, state.AdditionalSubnetIdList) {
		return nil
	}
	data, err := r.client.UpdateClusterSubnets(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Cluster Subnets",
			"Could not update cluster subnets, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	return nil
}

func (r *skeClusterResource) syncNfsVolume(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	// Skip when unset in config (unknown) to avoid sending an update that the API rejects.
	if plan.NfsVolumeId.IsUnknown() || plan.NfsVolumeId.Equal(state.NfsVolumeId) {
		return nil
	}
	data, err := r.client.UpdateClusterNfsVolume(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating NFS Volume",
			"Could not update NFS volume, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	return nil
}

func (r *skeClusterResource) syncDeletionProtection(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	// Skip when unset in config to avoid sending false (which would silently disable deletion protection).
	if plan.DeletionProtectionEnabled.IsNull() || plan.DeletionProtectionEnabled.Equal(state.DeletionProtectionEnabled) {
		return nil
	}
	data, _, err := r.client.UpdateClusterDeletionProtection(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Deletion Protection",
			"Could not update deletion protection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	return nil
}

func (r *skeClusterResource) syncLinkedResources(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.LinkedResources.IsUnknown() || linkedResourceListEqualSet(ctx, plan.LinkedResources, state.LinkedResources) {
		return nil
	}
	data, err := r.client.UpdateClusterLinkedResources(ctx, plan.Id.ValueString(), ske.LinkedResourcesFromList(ctx, plan.LinkedResources))
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Linked Resources",
			"Could not update linked resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	return nil
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *skeClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing cluster
	data, err := r.client.DeleteCluster(ctx, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Error waiting for cluster to become deleted: "+err.Error(),
		)
		return
	}
}

func (r *skeClusterResource) nullableStringToExternalResourceId(ns scpske.NullableString) ske.ExternalResourceId {
	if ns.Get() == nil {
		return ske.ExternalResourceId{
			Id: types.StringNull(),
		}
	}
	return ske.ExternalResourceId{
		Id: types.StringValue(*ns.Get()),
	}
}

// nullableStringToStringValue converts a NullableString to types.String (null-safe, treats empty as null).
func (r *skeClusterResource) nullableStringToStringValue(ns scpske.NullableString) types.String {
	if !ns.IsSet() || ns.Get() == nil || *ns.Get() == "" {
		return types.StringNull()
	}
	return types.StringValue(*ns.Get())
}

func (r *skeClusterResource) makePrivateEndpointAccessControlResourceModel(privateEndpointAccessControlResource *scpske.PrivateEndpointAccessControlResource) ske.PrivateEndpointAccessControlResource {
	return ske.PrivateEndpointAccessControlResource{
		Id:   types.StringValue(privateEndpointAccessControlResource.GetId()),
		Name: types.StringValue(privateEndpointAccessControlResource.GetName()),
		Type: types.StringValue(privateEndpointAccessControlResource.GetType()),
	}
}

//

func waitForClusterStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetCluster(ctx, id)
		if httpStatus == 200 {
			return info, info.Cluster.Status, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETED", nil
		} else if err != nil {
			return nil, "", err
		}

		return info, info.Cluster.Status, nil
	}, -1, -1, -1, -1)
}

func (r *skeClusterResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp) // 유니크한 스키마 입력
}
