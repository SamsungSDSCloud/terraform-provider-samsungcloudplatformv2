package virtualserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/vpc"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/virtualserver/1.4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource               = &virtualServerServerResource{}
	_ resource.ResourceWithConfigure  = &virtualServerServerResource{}
	_ resource.ResourceWithModifyPlan = &virtualServerServerResource{}
)

var osDiskDeviceNames = []string{"/dev/vda", "/dev/sda"}

// import 처럼 prior state 에 맵 키가 없을 때 API 결과에 붙일 키의 접두어.
// interface_1, interface_2 / volume_1, volume_2 ... 형태가 된다.
const (
	networkKeyPrefix = "interface"
	volumeKeyPrefix  = "volume"
)

// nextGeneratedKey 는 prefix_1, prefix_2 ... 중 아직 쓰이지 않은 첫 키를 돌려준다.
// 호출한 쪽에서 곧바로 existing 에 값을 넣어야 다음 호출이 다른 키를 받는다.
func nextGeneratedKey[T any](prefix string, existing map[string]T) string {
	for i := 1; ; i++ {
		key := fmt.Sprintf("%s_%d", prefix, i)
		if _, taken := existing[key]; !taken {
			return key
		}
	}
}

func isOsDisk(device string) bool {
	for _, name := range osDiskDeviceNames {
		if device == name {
			return true
		}
	}
	return false
}

func bootVolumeToObject(bootVolume virtualserver.ServerResourceVolume) types.Object {
	obj, _ := types.ObjectValueFrom(context.Background(),
		map[string]attr.Type{
			"id":                    types.StringType,
			"delete_on_termination": types.BoolType,
			"size":                  types.Int32Type,
			"type":                  types.StringType,
			"max_iops":              types.Int32Type,
			"max_throughput":        types.Int32Type,
		},
		bootVolume)
	return obj
}

func objectToBootVolume(obj types.Object) (virtualserver.ServerResourceVolume, error) {
	var bootVolume virtualserver.ServerResourceVolume
	if obj.IsNull() || obj.IsUnknown() {
		return bootVolume, nil
	}
	attrs := obj.Attributes()
	if id, ok := attrs["id"]; ok {
		bootVolume.Id = id.(types.String)
	}
	if del, ok := attrs["delete_on_termination"]; ok {
		bootVolume.DeleteOnTermination = del.(types.Bool)
	}
	if size, ok := attrs["size"]; ok {
		bootVolume.Size = size.(types.Int32)
	}
	if volType, ok := attrs["type"]; ok {
		bootVolume.Type = volType.(types.String)
	}
	if maxIops, ok := attrs["max_iops"]; ok {
		bootVolume.MaxIops = maxIops.(types.Int32)
	}
	if maxThroughput, ok := attrs["max_throughput"]; ok {
		bootVolume.MaxThroughput = maxThroughput.(types.Int32)
	}
	return bootVolume, nil
}

func NewVirtualServerServerResource() resource.Resource {
	return &virtualServerServerResource{}
}

type virtualServerServerResource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (r *virtualServerServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_server"
}

func (r *virtualServerServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Virtual Server resource.\n\n" +
			"**GPU Server Creation Guide:**\n" +
			"- GPU Server can only use GPU images (images with `scp_image_type` of `gpu_standard`, 'gpu_custom').\n" +
			"- GPU Server types, refer to SCP documentation for available options in your environment.",
		MarkdownDescription: "Virtual Server resource.\n\n" +
			"**GPU Server Creation Guide:**\n" +
			"- GPU Server can only use GPU images (images with `scp_image_type` of `gpu_standard`, 'gpu_custom').\n" +
			"- GPU Server types, refer to SCP documentation for available options in your environment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Resource ID.",
				MarkdownDescription: "Resource ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description:         "Account ID.",
				MarkdownDescription: "Account ID.",
				Computed:            true,
			},
			common.ToSnakeCase("Networks"): schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("PortId"): schema.StringAttribute{
							Description:         "Port ID.",
							MarkdownDescription: "Port ID.\n  - example: 12345678-1234-1234-1234-123456789012",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description:         "Subnet ID.",
							MarkdownDescription: "Subnet ID.\n  - example: ab313c43291e4b678f4bacffe10768ae",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
						},
						common.ToSnakeCase("FixedIp"): schema.StringAttribute{
							Description:         "Fixed IP address.",
							MarkdownDescription: "Fixed IP address.\n  - example: 192.168.1.100",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
						},
						common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
							Description:         "Public IP ID.",
							MarkdownDescription: "Public IP ID.\n  - example: a765a07e8d9b46f4918fd7d5ed004654",
							Optional:            true,
						},
						common.ToSnakeCase("StaticNatId"): schema.StringAttribute{
							Description:         "Static NAT ID.",
							MarkdownDescription: "Static NAT ID.",
							Computed:            true,
						},
						common.ToSnakeCase("IsDefault"): schema.BoolAttribute{
							Description:         "Whether this is the default port.",
							MarkdownDescription: "Whether this is the default port.",
							Computed:            true,
						},
					},
				},
				Required: true,
				Description: "Network settings. Defines network interfaces to attach to the server.\n" +
					"  - example: {\"network-1\": {\"subnet_id\": \"ab313c43291e4b678f4bacffe10768ae\"}}",
				MarkdownDescription: "Network settings. Defines network interfaces to attach to the server.\n" +
					"  - example: {\"network-1\": {\"subnet_id\": \"ab313c43291e4b678f4bacffe10768ae\"}}",
			},
			common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
				Description:         "Auto Scaling Group ID. Only has value for servers created by Auto Scaling Group.",
				MarkdownDescription: "Auto Scaling Group ID. Only has value for servers created by Auto Scaling Group.",
				Computed:            true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description:         "Creation timestamp.",
				MarkdownDescription: "Creation timestamp.",
				Computed:            true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description:         "Creator ID.",
				MarkdownDescription: "Creator ID.",
				Computed:            true,
			},
			common.ToSnakeCase("DiskConfig"): schema.StringAttribute{
				Description:         "Disk configuration mode.",
				MarkdownDescription: "Disk configuration mode.",
				Computed:            true,
			},
			common.ToSnakeCase("ImageId"): schema.StringAttribute{
				Description: "Image ID. Specifies the OS image to use.\n" +
					"  - example: 70a599e0-31e7-49b7-b260-868f441e862b\n" +
					"  - note: For GPU Server, only GPU standard images (scp_image_type=gpu_standard, gpu_custom) can be used.",
				MarkdownDescription: "Image ID. Specifies the OS image to use.\n" +
					"  - example: 70a599e0-31e7-49b7-b260-868f441e862b\n" +
					"  - note: For GPU Server, only GPU standard images (scp_image_type=gpu_standard, gpu_custom) can be used.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("KeypairName"): schema.StringAttribute{
				Description:         "Keypair name. Specifies the keypair for SSH access.",
				MarkdownDescription: "Keypair name. Specifies the keypair for SSH access.\n  - example: my-keypair",
				Required:            true,
			},
			common.ToSnakeCase("LaunchConfigurationId"): schema.StringAttribute{
				Description:         "Launch Configuration ID. Only has value for servers created by Auto Scaling Group.",
				MarkdownDescription: "Launch Configuration ID. Only has value for servers created by Auto Scaling Group.",
				Computed:            true,
			},
			common.ToSnakeCase("Lock"): schema.BoolAttribute{
				Description:         "Lock status. When locked, most user operations are blocked.",
				MarkdownDescription: "Lock status. When locked, most user operations are blocked.\n  - example: false",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			common.ToSnakeCase("Metadata"): schema.MapAttribute{
				Description:         "Metadata. Specifies key-value pairs to store on the server.",
				MarkdownDescription: "Metadata. Specifies key-value pairs to store on the server.\n  - example: {\"key\": \"value\"}",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description:         "Modification timestamp.",
				MarkdownDescription: "Modification timestamp.",
				Computed:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Server name.\n" +
					"  - example: my-server\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 63\n" +
					"  - pattern: ^[a-zA-Z0-9-_ ]*$",
				MarkdownDescription: "Server name.\n" +
					"  - example: my-server\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 63\n" +
					"  - pattern: ^[a-zA-Z0-9-_ ]*$",
				Required: true,
			},
			common.ToSnakeCase("PlannedComputeOsType"): schema.StringAttribute{
				Description:         "Planned compute OS type.",
				MarkdownDescription: "Planned compute OS type.",
				Computed:            true,
			},
			common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
				Description: "Product category.\n" +
					"  - example: compute\n" +
					"  - Available values: compute, container",
				MarkdownDescription: "Product category.\n" +
					"  - example: compute\n" +
					"  - Available values: compute, container",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
				Description: "Product offering. Determines the server type.\n" +
					"  - example: virtual_server\n" +
					"  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm\n",
				MarkdownDescription: "Product offering. Determines the server type.\n" +
					"  - example: virtual_server\n" +
					"  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm\n",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("SecurityGroups"): schema.ListAttribute{
				Description:         "Security group ID list.",
				MarkdownDescription: "Security group ID list.\n  - example: [\"c09c3f05-03d9-443f-b27a-40e0f973c75f\"]",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			common.ToSnakeCase("UserData"): schema.StringAttribute{
				Description: "User data script. Base64-encoded script to run on server startup.\n" +
					"  - example: IyEvYmluL2Jhc2gKL2Jpbi9zdQplY2hvICJJIGFtIGluIHlvdSEiCg==\n" +
					"  - maxLength: 65535 bytes",
				MarkdownDescription: "User data script. Base64-encoded script to run on server startup.\n" +
					"  - example: IyEvYmluL2Jhc2gKL2Jpbi9zdQplY2hvICJJIGFtIGluIHlvdSEiCg==\n" +
					"  - maxLength: 65535 bytes",
				Optional: true,
			},
			common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
				Description:         "Server group ID. Places the server in a specific server group.",
				MarkdownDescription: "Server group ID. Places the server in a specific server group.\n  - example: 616fb98f-46ca-475e-917e-2563e5a8cd19",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ServerTypeId"): schema.StringAttribute{
				Description: "Server type ID. Determines CPU, memory, and other server specifications.\n" +
					"  - example: s1v1m2\n",
				MarkdownDescription: "Server type ID. Determines CPU, memory, and other server specifications.\n" +
					"  - example: s1v1m2\n",
				Required: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "Server state.\n" +
					"  - example: ACTIVE\n" +
					"  - Available values: ACTIVE, SHUTOFF",
				MarkdownDescription: "Server state.\n" +
					"  - example: ACTIVE\n" +
					"  - Available values: ACTIVE, SHUTOFF",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			common.ToSnakeCase("BootVolume"): schema.SingleNestedAttribute{
				Description:         "Boot volume settings. Defines the root disk where OS is installed.",
				MarkdownDescription: "Boot volume settings. Defines the root disk where OS is installed.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "Volume ID.",
						MarkdownDescription: "Volume ID.",
						Computed:            true,
					},
					common.ToSnakeCase("DeleteOnTermination"): schema.BoolAttribute{
						Description:         "Whether to delete volume when server is terminated.",
						MarkdownDescription: "Whether to delete volume when server is terminated.\n  - example: true",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					common.ToSnakeCase("Size"): schema.Int32Attribute{
						Description:         "Volume size (GiB). Must be a multiple of 8.",
						MarkdownDescription: "Volume size (GiB). Must be a multiple of 8.\n  - example: 104\n  - minimum: 8",
						Required:            true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description:         "Volume type. Defaults to SSD if not specified.",
						MarkdownDescription: "Volume type. Defaults to SSD if not specified.\n  - example: SSD",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					common.ToSnakeCase("MaxIops"): schema.Int32Attribute{
						Description:         "Maximum IOPS. Number of read/write operations per second.",
						MarkdownDescription: "Maximum IOPS. Number of read/write operations per second.\n  - example: 10000",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
					},
					common.ToSnakeCase("MaxThroughput"): schema.Int32Attribute{
						Description:         "Maximum throughput (MB/s). Amount of data transferred per second.",
						MarkdownDescription: "Maximum throughput (MB/s). Amount of data transferred per second.\n  - example: 500",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
					},
				},
			},
			common.ToSnakeCase("ExtraVolumes"): schema.MapNestedAttribute{
				Description:         "Extra volume settings. Defines additional volumes to attach besides boot volume.",
				MarkdownDescription: "Extra volume settings. Defines additional volumes to attach besides boot volume.",
				Computed:            true,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "Volume ID.",
							MarkdownDescription: "Volume ID.",
							Computed:            true,
						},
						common.ToSnakeCase("DeleteOnTermination"): schema.BoolAttribute{
							Description:         "Whether to delete volume when server is terminated.",
							MarkdownDescription: "Whether to delete volume when server is terminated.\n  - example: true",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseNonNullStateForUnknown()},
						},
						common.ToSnakeCase("Size"): schema.Int32Attribute{
							Description:         "Volume size (GiB). Must be a multiple of 8.",
							MarkdownDescription: "Volume size (GiB). Must be a multiple of 8.\n  - example: 104\n  - minimum: 8",
							Required:            true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description:         "Volume type. Defaults to SSD_Provisioned if not specified.",
							MarkdownDescription: "Volume type. Defaults to SSD_Provisioned if not specified.\n  - example: SSD",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseNonNullStateForUnknown(),
							},
						},
						common.ToSnakeCase("MaxIops"): schema.Int32Attribute{
							Description:         "Maximum IOPS. Number of read/write operations per second.",
							MarkdownDescription: "Maximum IOPS. Number of read/write operations per second.\n  - example: 10000",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseNonNullStateForUnknown()},
						},
						common.ToSnakeCase("MaxThroughput"): schema.Int32Attribute{
							Description:         "Maximum throughput (MB/s). Amount of data transferred per second.",
							MarkdownDescription: "Maximum throughput (MB/s). Amount of data transferred per second.\n  - example: 500",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseNonNullStateForUnknown()},
						},
					},
				},
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description:         "VPC ID.",
				MarkdownDescription: "VPC ID.",
				Computed:            true,
			},
			common.ToSnakeCase("PartitionNumber"): schema.Int32Attribute{
				Description:         "Partition number. Only used when server group type is partition.",
				MarkdownDescription: "Partition number. Only used when server group type is partition.\n  - example: 1",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Zone"): schema.StringAttribute{
				Required:            true,
				Description:         "Zone ID\n  - example: kr-west1-a",
				MarkdownDescription: "Zone ID\n  - example: kr-west1-a",
			},
		},
	}

}

func (r *virtualServerServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VirtualServer
	r.clients = inst.Client
}

func (r *virtualServerServerResource) AsyncPollingTags(ctx context.Context, resourceId string, serviceName string,
	resourceType string, maxAttempts int, internal time.Duration) (types.Map, error) {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, resourceId, false)

		if err != nil {
			return types.Map{}, fmt.Errorf("attempt %d/%d failed: %w",
				attempt, maxAttempts, err)
		}

		if len(tagsMap.Elements()) > 0 {
			return tagsMap, nil
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return types.Map{}, fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return types.Map{}, fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func (r *virtualServerServerResource) AsyncPollingNetworkInterface(ctx context.Context, serverId string, portId string,
	requestType string, maxAttempts int, internal time.Duration) error {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := r.client.GetServerInterface(ctx, serverId, portId)

		if err != nil {
			var statusCode float64
			statusCode = virtualserverutil.GetStatusCodeFromError(err)

			if requestType == "InterfaceDelete" {
				if statusCode == 404 {
					return nil
				}
			}
			return fmt.Errorf("attempt %d/%d failed: %w", attempt, maxAttempts, err)
		}

		if len(resp.Interfaces) == 0 {
			if attempt < maxAttempts {
				select {
				case <-ticker.C:
					continue
				case <-ctx.Done():
					return fmt.Errorf("polling canceled: %w", ctx.Err())
				}
			}
			continue
		}

		if requestType == "NatCreate" {
			staticNat, natExist := resp.Interfaces[0].GetStaticNatOk()
			if natExist {
				if staticNat.State == "ACTIVE" {
					return nil
				}
			}
		}

		if requestType == "NatDelete" {
			_, natExist := resp.Interfaces[0].GetStaticNatOk()
			if !natExist {
				return nil
			}
		}

		if requestType == "InterfaceCreate" {
			portState := resp.Interfaces[0].GetPortState()
			if portState == "ACTIVE" {
				return nil
			}
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func (r *virtualServerServerResource) AsyncPollingQosUpdate(ctx context.Context, volumeId string, desiredIops, desiredThroughput int32) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 60; i++ {
		vol, err := r.client.GetVolume(ctx, volumeId)
		if err != nil {
			return fmt.Errorf("failed to get volume during polling: %w", err)
		}

		iopsMatch := *vol.MaxIops.Get() == desiredIops
		throughputMatch := *vol.MaxThroughput.Get() == desiredThroughput

		if iopsMatch && throughputMatch {
			return nil
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout waiting for volume update (ID: %s)", volumeId)
}

func (r *virtualServerServerResource) AsyncPollingVolumeIops(ctx context.Context, serverId string,
	bootExtraVolume types.Object, stateExtraVolumes types.Map) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 1000; i++ {
		_, _, allMatched, err := r.ResolveServerVolumes(ctx, serverId, bootExtraVolume, stateExtraVolumes)
		if err != nil {
			return fmt.Errorf("failed to get volume during polling: %w", err)
		}

		if allMatched {
			return nil
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout waiting for volume update")
}

func (r *virtualServerServerResource) matchVolumeAttributes(
	unmappedVolume virtualserver.ServerResourceVolume,
	planVolume virtualserver.ServerResourceVolume,
	defaultVolumeTypeName string,
) bool {
	if !(planVolume.Size.IsNull() || planVolume.Size.IsUnknown()) &&
		unmappedVolume.Size.ValueInt32() != planVolume.Size.ValueInt32() {
		return false
	}

	if !(planVolume.Type.IsNull() || planVolume.Type.IsUnknown()) &&
		unmappedVolume.Type.ValueString() != planVolume.Type.ValueString() {
		return false
	}
	if (planVolume.Type.IsNull() || planVolume.Type.IsUnknown()) &&
		unmappedVolume.Type.ValueString() != defaultVolumeTypeName {
		return false
	}

	if !(planVolume.DeleteOnTermination.IsNull() || planVolume.DeleteOnTermination.IsUnknown()) &&
		unmappedVolume.DeleteOnTermination.ValueBool() != planVolume.DeleteOnTermination.ValueBool() {
		return false
	}
	if (planVolume.DeleteOnTermination.IsNull() || planVolume.DeleteOnTermination.IsUnknown()) &&
		unmappedVolume.DeleteOnTermination.ValueBool() != false {
		return false
	}

	if !(planVolume.MaxIops.IsNull() || planVolume.MaxIops.IsUnknown()) &&
		unmappedVolume.MaxIops.ValueInt32() != planVolume.MaxIops.ValueInt32() {
		return false
	}

	if !(planVolume.MaxThroughput.IsNull() || planVolume.MaxThroughput.IsUnknown()) &&
		unmappedVolume.MaxThroughput.ValueInt32() != planVolume.MaxThroughput.ValueInt32() {
		return false
	}

	return true
}

func (r *virtualServerServerResource) findMatchingVolumeKey(
	unmappedVolume virtualserver.ServerResourceVolume,
	extraVolumesMap map[string]virtualserver.ServerResourceVolume,
	mappedVolumeKeys map[string]bool,
	defaultVolumeTypeName string,
) string {
	for key, planVolume := range extraVolumesMap {
		if mappedVolumeKeys[key] {
			continue
		}

		if r.matchVolumeAttributes(unmappedVolume, planVolume, defaultVolumeTypeName) {
			return key
		}
	}
	return ""
}

// MapUnmappedExtraVolumes 는 state 키에 짝지어지지 않은 볼륨들을 속성 일치로 배치하고,
// 끝까지 짝을 못 찾은 볼륨 목록을 돌려준다.
// (이전에는 allMatched bool 만 돌려줘서 남은 볼륨이 그대로 유실됐다.)
func (r *virtualServerServerResource) MapUnmappedExtraVolumes(
	unmappedExtraVolumes []virtualserver.ServerResourceVolume,
	extraVolumesMap map[string]virtualserver.ServerResourceVolume,
	mappedVolumeKeys map[string]bool,
	defaultVolumeTypeName string,
) []virtualserver.ServerResourceVolume {
	var remaining []virtualserver.ServerResourceVolume

	for _, unmappedVolume := range unmappedExtraVolumes {
		bestKey := r.findMatchingVolumeKey(unmappedVolume, extraVolumesMap, mappedVolumeKeys, defaultVolumeTypeName)

		if bestKey != "" {
			extraVolumesMap[bestKey] = unmappedVolume
			mappedVolumeKeys[bestKey] = true
		} else {
			remaining = append(remaining, unmappedVolume)
		}
	}

	return remaining
}

func (r *virtualServerServerResource) isBootVolumeMatched(
	bootVolume virtualserver.ServerResourceVolume,
	stateBootVolume virtualserver.ServerResourceVolume,
) bool {
	if bootVolume.Id.IsNull() || bootVolume.Id.IsUnknown() {
		return false
	}

	if !stateBootVolume.MaxIops.IsNull() && !stateBootVolume.MaxIops.IsUnknown() &&
		bootVolume.MaxIops.ValueInt32() != stateBootVolume.MaxIops.ValueInt32() {
		return false
	}
	if !stateBootVolume.MaxThroughput.IsNull() && !stateBootVolume.MaxThroughput.IsUnknown() &&
		bootVolume.MaxThroughput.ValueInt32() != stateBootVolume.MaxThroughput.ValueInt32() {
		return false
	}

	return true
}

func (r *virtualServerServerResource) buildVolumeSets(getServerVolumes *scpvirtualserver.ServerVolumesResponse) (
	map[string]bool, map[string]bool, map[string]bool) {
	volumeBootVolumeSet := make(map[string]bool)
	volumeIdsSet := make(map[string]bool)
	volumeDeleteOnTerminationSet := make(map[string]bool)

	for _, volume := range getServerVolumes.Volumes {
		volumeIdsSet[volume.VolumeId] = true
		if isOsDisk(volume.Device) {
			volumeBootVolumeSet[volume.VolumeId] = true
		} else {
			volumeBootVolumeSet[volume.VolumeId] = false
		}
		volumeDeleteOnTerminationSet[volume.Id] = volume.DeleteOnTermination
	}

	return volumeBootVolumeSet, volumeIdsSet, volumeDeleteOnTerminationSet
}

func (r *virtualServerServerResource) categorizeVolumes(
	getVolumes *scpvirtualserver.VolumeListResponseV1Dot4,
	volumeIdsSet map[string]bool,
	volumeDeleteOnTerminationSet map[string]bool,
	volumeBootVolumeSet map[string]bool,
) (virtualserver.ServerResourceVolume, []virtualserver.ServerResourceVolume) {
	var bootVolume virtualserver.ServerResourceVolume
	var extraVolumes []virtualserver.ServerResourceVolume

	for _, volume := range getVolumes.Volumes {
		if volumeIdsSet[volume.Id] {
			volumeObject := virtualserver.ServerResourceVolume{
				Id:                  types.StringValue(volume.Id),
				DeleteOnTermination: types.BoolValue(volumeDeleteOnTerminationSet[volume.Id]),
				Size:                types.Int32Value(volume.Size),
				Type:                types.StringValue(volume.VolumeType),
				MaxIops:             types.Int32PointerValue(volume.MaxIops.Get()),
				MaxThroughput:       types.Int32PointerValue(volume.MaxThroughput.Get()),
			}
			if volumeBootVolumeSet[volume.Id] {
				bootVolume = volumeObject
			} else {
				extraVolumes = append(extraVolumes, volumeObject)
			}
		}
	}

	return bootVolume, extraVolumes
}

func (r *virtualServerServerResource) processExtraVolumes(
	ctx context.Context,
	extraVolumes []virtualserver.ServerResourceVolume,
	stateExtraVolumes types.Map,
) (map[string]virtualserver.ServerResourceVolume, []virtualserver.ServerResourceVolume, map[string]bool) {
	var extraVolumesMap map[string]virtualserver.ServerResourceVolume
	stateExtraVolumes.ElementsAs(ctx, &extraVolumesMap, false)
	if extraVolumesMap == nil {
		// import 처럼 prior state 가 null 이면 ElementsAs 는 nil map 을 남긴다.
		// 이 상태로 두면 뒤에서 키를 넣을 수 없어 extra_volumes 가 전부 유실된다.
		extraVolumesMap = make(map[string]virtualserver.ServerResourceVolume)
	}

	extraVolumeIdKeyMap := make(map[string]string)
	for key, volume := range extraVolumesMap {
		if !(volume.Id.IsUnknown() || volume.Id.IsNull()) {
			extraVolumeIdKeyMap[volume.Id.ValueString()] = key
		}
	}

	var unmappedExtraVolumes []virtualserver.ServerResourceVolume
	mappedVolumeKeys := make(map[string]bool)
	for _, volume := range extraVolumes {
		if key, exist := extraVolumeIdKeyMap[volume.Id.ValueString()]; exist {
			extraVolumesMap[key] = volume
			mappedVolumeKeys[key] = true
		} else {
			unmappedExtraVolumes = append(unmappedExtraVolumes, volume)
		}
	}

	for _, key := range extraVolumeIdKeyMap {
		mappedVolumeKeys[key] = true
	}

	return extraVolumesMap, unmappedExtraVolumes, mappedVolumeKeys
}

func (r *virtualServerServerResource) ResolveServerVolumes(
	ctx context.Context,
	serverId string,
	stateBootVolume types.Object,
	stateExtraVolumes types.Map) (virtualserver.ServerResourceVolume, types.Map, bool, error) {

	bootVolumeFromState, err := objectToBootVolume(stateBootVolume)
	if err != nil {
		return virtualserver.ServerResourceVolume{}, types.Map{}, false, err
	}

	getServerVolumes, err := r.client.GetServerVolumeList(ctx, serverId)
	if err != nil {
		return virtualserver.ServerResourceVolume{}, types.Map{}, false, err
	}

	volumeBootVolumeSet, volumeIdsSet, volumeDeleteOnTerminationSet := r.buildVolumeSets(getServerVolumes)

	getVolumes, err := r.client.GetVolumeList()
	if err != nil {
		return virtualserver.ServerResourceVolume{}, types.Map{}, false, err
	}

	bootVolume, extraVolumes := r.categorizeVolumes(getVolumes, volumeIdsSet, volumeDeleteOnTerminationSet, volumeBootVolumeSet)

	bootVolumeMatched := r.isBootVolumeMatched(bootVolume, bootVolumeFromState)

	extraVolumesMap, unmappedExtraVolumes, mappedVolumeKeys := r.processExtraVolumes(ctx, extraVolumes, stateExtraVolumes)

	defaultVolumeType, err := r.client.GetDefaultVolumeType(ctx)
	if err != nil {
		return virtualserver.ServerResourceVolume{}, types.Map{}, false, err
	}
	defaultVolumeTypeName := ""
	if name := defaultVolumeType.Name.Get(); name != nil {
		defaultVolumeTypeName = *name
	}

	unmappedExtraVolumes = r.MapUnmappedExtraVolumes(unmappedExtraVolumes, extraVolumesMap, mappedVolumeKeys, defaultVolumeTypeName)

	// 끝까지 state 키에 짝지어지지 않은 볼륨은 새 키(volume_1, volume_2 ...)로 담는다.
	// import 이거나 콘솔에서 볼륨을 붙인 경우가 여기에 해당한다.
	for _, volume := range unmappedExtraVolumes {
		key := nextGeneratedKey(volumeKeyPrefix, extraVolumesMap)
		extraVolumesMap[key] = volume
		mappedVolumeKeys[key] = true
	}

	extraVolumeMatched := true
	for _, volume := range extraVolumesMap {
		if volume.Id.IsNull() || volume.Id.IsUnknown() {
			extraVolumeMatched = false
			break
		}
	}

	volumeElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                    types.StringType,
			"delete_on_termination": types.BoolType,
			"size":                  types.Int32Type,
			"type":                  types.StringType,
			"max_iops":              types.Int32Type,
			"max_throughput":        types.Int32Type,
		},
	}

	extraVolumeObject, _ := types.MapValueFrom(ctx, volumeElemType, extraVolumesMap)
	allMatched := bootVolumeMatched && extraVolumeMatched

	return bootVolume, extraVolumeObject, allMatched, nil
}

// mapNetworksBySubnetAndFixedIp 는 subnet_id + fixed_ip 가 모두 일치하는 state 키에 배치하고,
// 짝을 못 찾은 인터페이스 목록을 돌려준다.
func (r *virtualServerServerResource) mapNetworksBySubnetAndFixedIp(
	unmappedNetworks []virtualserver.ServerResourceNetwork,
	networkMap map[string]virtualserver.ServerResourceNetwork,
	mappedNetworkKeys map[string]bool,
) []virtualserver.ServerResourceNetwork {
	var remaining []virtualserver.ServerResourceNetwork

	for _, unmappedNetwork := range unmappedNetworks {
		matched := false

		for key, planNetwork := range networkMap {
			if mappedNetworkKeys[key] {
				continue
			}

			if !(planNetwork.SubnetId.IsNull() || planNetwork.SubnetId.IsUnknown()) &&
				!(planNetwork.FixedIp.IsNull() || planNetwork.FixedIp.IsUnknown()) {
				if unmappedNetwork.SubnetId.ValueString() == planNetwork.SubnetId.ValueString() &&
					unmappedNetwork.FixedIp.ValueString() == planNetwork.FixedIp.ValueString() {
					networkMap[key] = unmappedNetwork
					mappedNetworkKeys[key] = true
					matched = true
					break
				}
			}
		}

		if !matched {
			remaining = append(remaining, unmappedNetwork)
		}
	}

	return remaining
}

// mapNetworksBySubnetOnly 는 fixed_ip 를 안 적은 state 키에 subnet_id 만으로 배치하고,
// 짝을 못 찾은 인터페이스 목록을 돌려준다.
func (r *virtualServerServerResource) mapNetworksBySubnetOnly(
	unmappedNetworks []virtualserver.ServerResourceNetwork,
	networkMap map[string]virtualserver.ServerResourceNetwork,
	mappedNetworkKeys map[string]bool,
) []virtualserver.ServerResourceNetwork {
	var remaining []virtualserver.ServerResourceNetwork

	for _, unmappedNetwork := range unmappedNetworks {
		matched := false

		for key, planNetwork := range networkMap {
			if mappedNetworkKeys[key] {
				continue
			}

			if !(planNetwork.SubnetId.IsNull() || planNetwork.SubnetId.IsUnknown()) &&
				(planNetwork.FixedIp.IsNull() || planNetwork.FixedIp.IsUnknown()) {
				if unmappedNetwork.SubnetId.ValueString() == planNetwork.SubnetId.ValueString() {
					networkMap[key] = unmappedNetwork
					mappedNetworkKeys[key] = true
					matched = true
					break
				}
			}
		}

		if !matched {
			remaining = append(remaining, unmappedNetwork)
		}
	}

	return remaining
}

func (r *virtualServerServerResource) processNetworks(
	ctx context.Context,
	resp *scpvirtualserver.ServerShowResponseV1Dot4,
	state virtualserver.ServerResource,
) (types.Map, error) {
	var networkMap map[string]virtualserver.ServerResourceNetwork
	state.Networks.ElementsAs(ctx, &networkMap, false)
	if networkMap == nil {
		// import 처럼 prior state 가 null 이면 ElementsAs 는 nil map 을 남긴다.
		// 이 상태로 두면 뒤에서 키를 넣을 수 없어 networks 가 전부 유실된다.
		networkMap = make(map[string]virtualserver.ServerResourceNetwork)
	}

	var networkKeyPortMap = make(map[string]string)
	for key, network := range networkMap {
		if !network.PortId.IsUnknown() && !network.PortId.IsNull() {
			networkKeyPortMap[network.PortId.ValueString()] = key
		}
	}

	getInterfaces, err := r.client.GetServerInterfaceList(ctx, resp.Id)
	if err != nil {
		return types.Map{}, err
	}

	var unmappedNetworks []virtualserver.ServerResourceNetwork
	for _, itf := range getInterfaces.Interfaces {
		var publicIpId, staticNatId *string
		staticNat, natExist := itf.GetStaticNatOk()
		if natExist {
			publicIpId = staticNat.PublicipId.Get()
			staticNatId = &staticNat.Id
		}

		var fixedIp string
		if len(itf.FixedIps) > 0 {
			fixedIp = itf.FixedIps[0].IpAddress
		}

		network := virtualserver.ServerResourceNetwork{
			SubnetId:    types.StringValue(itf.SubnetId),
			PortId:      types.StringValue(itf.PortId),
			FixedIp:     types.StringValue(fixedIp),
			PublicIpId:  types.StringPointerValue(publicIpId),
			StaticNatId: types.StringPointerValue(staticNatId),
			IsDefault:   types.BoolValue(itf.IsDefault),
		}

		if _, exist := networkKeyPortMap[itf.PortId]; exist {
			networkMap[networkKeyPortMap[itf.PortId]] = network
		} else {
			unmappedNetworks = append(unmappedNetworks, network)
		}
	}

	mappedNetworkKeys := make(map[string]bool)
	for _, key := range networkKeyPortMap {
		mappedNetworkKeys[key] = true
	}

	unmappedNetworks = r.mapNetworksBySubnetAndFixedIp(unmappedNetworks, networkMap, mappedNetworkKeys)
	unmappedNetworks = r.mapNetworksBySubnetOnly(unmappedNetworks, networkMap, mappedNetworkKeys)

	// 끝까지 state 키에 짝지어지지 않은 인터페이스는 새 키(interface_1, interface_2 ...)로 담는다.
	// import 이거나 콘솔에서 NIC 을 붙인 경우가 여기에 해당한다.
	for _, network := range unmappedNetworks {
		key := nextGeneratedKey(networkKeyPrefix, networkMap)
		networkMap[key] = network
		mappedNetworkKeys[key] = true
	}

	networkElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"subnet_id":     types.StringType,
			"port_id":       types.StringType,
			"fixed_ip":      types.StringType,
			"public_ip_id":  types.StringType,
			"static_nat_id": types.StringType,
			"is_default":    types.BoolType,
		},
	}

	networks, _ := types.MapValueFrom(ctx, networkElemType, networkMap)
	return networks, nil
}

func (r *virtualServerServerResource) processMetadata(resp *scpvirtualserver.ServerShowResponseV1Dot4) types.Map {
	metadataMap := make(map[string]attr.Value)
	for k, v := range resp.Metadata {
		// v.(string) 단정은 API 가 string 이 아닌 값을 주면 panic 한다.
		// 스키마상 Map[String] 이므로 문자열로 표현해서 담는다.
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprintf("%v", v)
		}
		metadataMap[k] = types.StringValue(s)
	}
	metadata, _ := types.MapValue(types.StringType, metadataMap)
	return metadata
}

func (r *virtualServerServerResource) processSecurityGroups(
	ctx context.Context,
	resp *scpvirtualserver.ServerShowResponseV1Dot4,
	state virtualserver.ServerResource,
) ([]attr.Value, error) {
	getSecurityGroups, err := r.client.GetServerSecurityGroupList(ctx, resp.Id)
	if err != nil {
		return nil, err
	}

	attached := make(map[string]bool, len(getSecurityGroups.SecurityGroups))
	for _, securityGroup := range getSecurityGroups.SecurityGroups {
		attached[securityGroup.Id] = true
	}

	// 이전 구현은 슬라이스 길이를 API 개수로 잡고 prior state 원소 루프로만 채웠다.
	// 그래서 state 가 API 보다 적으면(= import 는 항상 0개) 원소가 nil 로 남아
	// ListValueMust 에서 nil pointer dereference panic 이 났다.
	// 여기서는 state 순서를 유지하면서(순서만 바뀌는 phantom diff 방지)
	// 붙어 있는 것만 남기고, state 에 없던 것은 뒤에 붙인다. nil 원소가 생기지 않는다.
	securityGroups := make([]attr.Value, 0, len(getSecurityGroups.SecurityGroups))
	included := make(map[string]bool, len(getSecurityGroups.SecurityGroups))

	for _, stateSecurityGroup := range state.SecurityGroups.Elements() {
		id, ok := stateSecurityGroup.(types.String)
		if !ok || id.IsNull() || id.IsUnknown() {
			continue
		}
		if attached[id.ValueString()] && !included[id.ValueString()] {
			securityGroups = append(securityGroups, id)
			included[id.ValueString()] = true
		}
	}

	for _, securityGroup := range getSecurityGroups.SecurityGroups {
		if !included[securityGroup.Id] {
			securityGroups = append(securityGroups, types.StringValue(securityGroup.Id))
			included[securityGroup.Id] = true
		}
	}

	return securityGroups, nil
}

func (r *virtualServerServerResource) MapGetResponseToState(ctx context.Context,
	resp *scpvirtualserver.ServerShowResponseV1Dot4, state virtualserver.ServerResource, tagsMap types.Map) (virtualserver.ServerResource, error) {
	networks, err := r.processNetworks(ctx, resp, state)
	if err != nil {
		return virtualserver.ServerResource{}, err
	}

	metadata := r.processMetadata(resp)

	plannedComputeOsTypeJson, _ := resp.PlannedComputeOsType.MarshalJSON()
	plannedComputeOsType := strings.Trim(string(plannedComputeOsTypeJson), "\"")

	productCategoryJson, _ := resp.ProductCategory.MarshalJSON()
	productCategory := strings.Trim(string(productCategoryJson), "\"")

	productOfferingJson, _ := resp.ProductOffering.MarshalJSON()
	productOffering := strings.Trim(string(productOfferingJson), "\"")

	securityGroups, err := r.processSecurityGroups(ctx, resp, state)
	if err != nil {
		return virtualserver.ServerResource{}, err
	}

	bootVolume, extraVolumeObject, _, err := r.ResolveServerVolumes(ctx, resp.Id, state.BootVolume, state.ExtraVolumes)
	if err != nil {
		// 이전에는 err 를 받아놓고 검사하지 않아, 볼륨 조회가 실패해도
		// 에러 없이 빈 boot_volume / extra_volumes 가 state 에 들어갔다.
		return virtualserver.ServerResource{}, err
	}

	return virtualserver.ServerResource{
		Id:                    types.StringValue(resp.Id),
		AccountId:             types.StringValue(resp.AccountId),
		Networks:              networks,
		AutoScalingGroupId:    virtualserverutil.ToNullableStringValue(resp.AutoScalingGroupId.Get()),
		CreatedAt:             types.StringValue(resp.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringValue(resp.CreatedBy),
		DiskConfig:            types.StringValue(resp.DiskConfig),
		ImageId:               virtualserverutil.ToNullableStringValue(resp.ImageId.Get()),
		KeypairName:           virtualserverutil.ToNullableStringValue(resp.KeypairName.Get()),
		LaunchConfigurationId: virtualserverutil.ToNullableStringValue(resp.LaunchConfigurationId.Get()),
		Lock:                  types.BoolValue(resp.Locked),
		Metadata:              metadata,
		ModifiedAt:            types.StringValue(resp.ModifiedAt.Format(time.RFC3339)),
		Name:                  types.StringValue(resp.Name),
		PlannedComputeOsType:  types.StringValue(plannedComputeOsType),
		ProductCategory:       types.StringValue(productCategory),
		ProductOffering:       types.StringValue(productOffering),
		SecurityGroups:        types.ListValueMust(types.StringType, securityGroups),
		UserData:              state.UserData,
		ServerGroupId:         virtualserverutil.ToNullableStringValue(resp.ServerGroupId.Get()),
		ServerTypeId:          virtualserverutil.ToNullableStringValue(resp.ServerType.Id.Get()),
		State:                 types.StringValue(resp.State),
		BootVolume:            bootVolumeToObject(bootVolume),
		ExtraVolumes:          extraVolumeObject,
		VpcId:                 virtualserverutil.ToNullableStringValue(resp.VpcId.Get()),
		PartitionNumber:       virtualserverutil.ToNullableInt32Value(resp.PartitionNumber.Get()),
		Tags:                  tagsMap,
		Zone:                  types.StringValue(resp.Zone),
	}, nil
}

func (r *virtualServerServerResource) normalizePlan(plan virtualserver.ServerResource, state virtualserver.ServerResource) virtualserver.ServerResource {
	// nomarlizePlan nomalizes the plan by replacing unknown values with their previous state values

	// Metadata
	if plan.Metadata.IsUnknown() {
		plan.Metadata = state.Metadata
	}
	// ProductCategory
	if plan.ProductCategory.IsUnknown() {
		plan.ProductCategory = state.ProductCategory
	}
	// ProductOffering
	if plan.ProductOffering.IsUnknown() {
		plan.ProductOffering = state.ProductOffering
	}
	// BootVolume
	planBootVolume, _ := objectToBootVolume(plan.BootVolume)
	stateBootVolume, _ := objectToBootVolume(state.BootVolume)
	if planBootVolume.Id.IsUnknown() {
		planBootVolume.Id = stateBootVolume.Id
	}
	if planBootVolume.Type.IsUnknown() {
		planBootVolume.Type = stateBootVolume.Type
	}
	if planBootVolume.DeleteOnTermination.IsUnknown() {
		planBootVolume.DeleteOnTermination = stateBootVolume.DeleteOnTermination
	}
	plan.BootVolume = bootVolumeToObject(planBootVolume)

	return plan
}

func (r *virtualServerServerResource) handlerUpdateServerName(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)

	_, err := r.client.UpdateServer(ctx, plan)
	if err != nil {
		return err
	}

	return nil
}

func (r *virtualServerServerResource) handlerUpdateServerLock(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)

	var err error
	if plan.Lock.ValueBool() {
		err = r.client.UpdateServerLock(ctx, plan.Id.ValueString())
	} else {
		err = r.client.UpdateServerUnlock(ctx, plan.Id.ValueString())
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *virtualServerServerResource) handlerUpdateServerSecurityGroup(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	var securityGroups []string
	diags := plan.SecurityGroups.ElementsAs(ctx, &securityGroups, false)
	resp.Diagnostics.Append(diags...)

	var networkMap map[string]virtualserver.ServerResourceNetwork
	state.Networks.ElementsAs(ctx, &networkMap, false)

	for _, network := range networkMap {
		_, err := r.clients.Vpc.UpdatePort(ctx, network.PortId.ValueString(), vpc.PortResource{SecurityGroups: securityGroups})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *virtualServerServerResource) handlerUpdateServerType(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	err := r.client.UpdateServerType(ctx, plan)
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpvirtualserver.ServerShowResponseV1Dot4, error) {
		return r.client.GetServer(ctx, id)
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 1000, 3*time.Second,
		"ServerType.Name", plan.ServerTypeId.ValueString(), "", getFunc)
	if err != nil {
		return err
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 1000, 3*time.Second,
		"State", state.State.ValueString(), "ERROR", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *virtualServerServerResource) normalizeNetworkFields(
	planNetwork virtualserver.ServerResourceNetwork,
	stateNetwork virtualserver.ServerResourceNetwork,
) virtualserver.ServerResourceNetwork {
	if planNetwork.PublicIpId.IsUnknown() {
		planNetwork.PublicIpId = stateNetwork.PublicIpId
	}
	if planNetwork.SubnetId.IsUnknown() {
		planNetwork.SubnetId = stateNetwork.SubnetId
	}
	if planNetwork.FixedIp.IsUnknown() {
		planNetwork.FixedIp = stateNetwork.FixedIp
	}
	if planNetwork.PortId.IsUnknown() {
		planNetwork.PortId = stateNetwork.PortId
	}
	if planNetwork.IsDefault.IsUnknown() {
		planNetwork.IsDefault = stateNetwork.IsDefault
	}
	return planNetwork
}

func (r *virtualServerServerResource) updateNetworkPublicIp(
	ctx context.Context,
	serverId string,
	planNetwork virtualserver.ServerResourceNetwork,
	stateNetwork virtualserver.ServerResourceNetwork,
) error {
	if planNetwork.PublicIpId.IsNull() {
		err := r.client.DeleteServerInterfaceNat(ctx, serverId, planNetwork.PortId.ValueString(), stateNetwork.StaticNatId.ValueString())
		if err != nil {
			return err
		}
		return r.AsyncPollingNetworkInterface(ctx, serverId, planNetwork.PortId.ValueString(), "NatDelete", 1000, 3*time.Second)
	} else {
		_, err := r.client.CreateServerInterfaceNat(ctx, serverId, planNetwork.PortId.ValueString(), planNetwork.PublicIpId.ValueString())
		if err != nil {
			return err
		}
		return r.AsyncPollingNetworkInterface(ctx, serverId, planNetwork.PortId.ValueString(), "NatCreate", 1000, 3*time.Second)
	}
}

func (r *virtualServerServerResource) updateExistingNetwork(
	ctx context.Context,
	serverId string,
	planNetwork virtualserver.ServerResourceNetwork,
	stateNetwork virtualserver.ServerResourceNetwork,
) error {
	normalizedPlanNetwork := r.normalizeNetworkFields(planNetwork, stateNetwork)

	networksFields := []string{"PortId", "SubnetId", "FixedIp", "PublicIpId", "IsDefault"}
	immutableNetworksFiled := []string{"PortId", "SubnetId", "FixedIp", "IsDefault"}

	changesFields, err := virtualserverutil.GetChangedFields(normalizedPlanNetwork, stateNetwork, networksFields)
	if err != nil {
		return err
	}

	if len(changesFields) == 0 {
		return nil
	}

	if virtualserverutil.IsOverlapFields(changesFields, immutableNetworksFiled) {
		detailed := fmt.Sprintf("Plan: %v,\nState: %v", normalizedPlanNetwork, stateNetwork)
		return fmt.Errorf("Immutable Networks Field changes: %s\n%s", strings.Join(immutableNetworksFiled, ", "), detailed)
	}

	if virtualserverutil.IsOverlapFields(changesFields, []string{"PublicIpId"}) {
		return r.updateNetworkPublicIp(ctx, serverId, normalizedPlanNetwork, stateNetwork)
	}

	return nil
}

func (r *virtualServerServerResource) createNewNetwork(
	ctx context.Context,
	serverId string,
	planNetwork virtualserver.ServerResourceNetwork,
) error {
	resp, err := r.client.CreateServerInterface(ctx, serverId, planNetwork)
	if err != nil {
		return err
	}
	return r.AsyncPollingNetworkInterface(ctx, serverId, resp.PortId, "InterfaceCreate", 1000, 3*time.Second)
}

func (r *virtualServerServerResource) deleteNetwork(
	ctx context.Context,
	serverId string,
	stateNetwork virtualserver.ServerResourceNetwork,
) error {
	err := r.client.DeleteServerInterface(ctx, serverId, stateNetwork.PortId.ValueString())
	if err != nil {
		return err
	}
	return r.AsyncPollingNetworkInterface(ctx, serverId, stateNetwork.PortId.ValueString(), "InterfaceDelete", 1000, 3*time.Second)
}

func (r *virtualServerServerResource) handlerUpdateServerNetwork(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	var planNetworks, stateNetworks map[string]virtualserver.ServerResourceNetwork
	plan.Networks.ElementsAs(ctx, &planNetworks, false)
	state.Networks.ElementsAs(ctx, &stateNetworks, false)

	serverId := plan.Id.ValueString()

	for key, planNetwork := range planNetworks {
		if stateNetwork, exists := stateNetworks[key]; exists {
			if err := r.updateExistingNetwork(ctx, serverId, planNetwork, stateNetwork); err != nil {
				return err
			}
		}
	}

	for key, planNetwork := range planNetworks {
		if _, exists := stateNetworks[key]; !exists {
			if err := r.createNewNetwork(ctx, serverId, planNetwork); err != nil {
				return err
			}
		}
	}

	for key, stateNetwork := range stateNetworks {
		if _, exists := planNetworks[key]; !exists {
			if err := r.deleteNetwork(ctx, serverId, stateNetwork); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *virtualServerServerResource) updateVolumeSize(
	ctx context.Context,
	serverId string,
	planVolume virtualserver.ServerResourceVolume,
	stateVolume virtualserver.ServerResourceVolume,
) error {
	if planVolume.Size.ValueInt32() <= stateVolume.Size.ValueInt32() {
		return fmt.Errorf("Volume size must be expanded. \nPlan : %v, State: %v",
			planVolume.Size.ValueInt32(), stateVolume.Size.ValueInt32())
	}
	if planVolume.Size.ValueInt32()%8 != 0 {
		return fmt.Errorf("volume size must be a multiple of 8")
	}

	_, err := r.client.ExtendVolume(ctx, stateVolume.Id.ValueString(), virtualserver.VolumeResource{
		Size: planVolume.Size,
	})
	return err
}

func (r *virtualServerServerResource) updateVolumeQos(
	ctx context.Context,
	planVolume virtualserver.ServerResourceVolume,
	stateVolume virtualserver.ServerResourceVolume,
) error {
	volumeResource := virtualserver.VolumeResource{
		MaxIops:       planVolume.MaxIops,
		MaxThroughput: planVolume.MaxThroughput,
	}
	err := r.client.UpdateVolumeQos(ctx, stateVolume.Id.ValueString(), volumeResource)
	if err != nil {
		return err
	}

	if !((planVolume.MaxThroughput.IsNull() || planVolume.MaxThroughput.IsUnknown()) && (planVolume.MaxIops.IsNull() || planVolume.MaxIops.IsUnknown())) {
		return r.AsyncPollingQosUpdate(ctx, stateVolume.Id.ValueString(), planVolume.MaxIops.ValueInt32(), planVolume.MaxThroughput.ValueInt32())
	}

	return nil
}

func (r *virtualServerServerResource) updateSingleVolume(
	ctx context.Context,
	serverId string,
	planVolume virtualserver.ServerResourceVolume,
	stateVolume virtualserver.ServerResourceVolume,
) error {
	if !planVolume.Type.Equal(stateVolume.Type) {
		return fmt.Errorf("immutable Volume Field changes: Type")
	}

	if !planVolume.Size.Equal(stateVolume.Size) {
		if err := r.updateVolumeSize(ctx, serverId, planVolume, stateVolume); err != nil {
			return err
		}
	}

	if !planVolume.DeleteOnTermination.Equal(stateVolume.DeleteOnTermination) {
		if err := r.client.UpdateServerVolume(ctx, serverId, stateVolume.Id.ValueString(), planVolume.DeleteOnTermination.ValueBool()); err != nil {
			return err
		}
	}

	if !planVolume.MaxThroughput.Equal(stateVolume.MaxThroughput) || !planVolume.MaxIops.Equal(stateVolume.MaxIops) {
		if err := r.updateVolumeQos(ctx, planVolume, stateVolume); err != nil {
			return err
		}
	}

	return nil
}

func (r *virtualServerServerResource) createExtraVolume(
	ctx context.Context,
	serverId string,
	plan virtualserver.ServerResource,
	planExtraVolume virtualserver.ServerResourceVolume,
	getVolumeFunc func(string) (*scpvirtualserver.VolumeShowResponseV1Dot4, error),
) error {
	volumeResource := virtualserver.VolumeResource{
		Name:          types.StringValue(serverId + "-blank-vol"),
		Size:          planExtraVolume.Size,
		VolumeType:    planExtraVolume.Type,
		Tags:          plan.Tags,
		MaxIops:       planExtraVolume.MaxIops,
		MaxThroughput: planExtraVolume.MaxThroughput,
		Zone:          plan.Zone,
	}
	resp, err := r.client.CreateVolume(ctx, volumeResource)
	if err != nil {
		return err
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, resp.Id, 100, 3*time.Second,
		"State", "available", "error", getVolumeFunc)
	if err != nil {
		return err
	}

	_, err = r.client.AttachVolume(ctx, resp.Id, serverId)
	if err != nil {
		return err
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, resp.Id, 100, 3*time.Second,
		"State", "in-use", "error", getVolumeFunc)
	if err != nil {
		return err
	}

	if planExtraVolume.DeleteOnTermination.ValueBool() {
		return r.client.UpdateServerVolume(ctx, serverId, resp.Id, planExtraVolume.DeleteOnTermination.ValueBool())
	}

	return nil
}

func (r *virtualServerServerResource) deleteExtraVolume(
	ctx context.Context,
	stateVolume virtualserver.ServerResourceVolume,
	serverId string,
	getVolumeFunc func(string) (*scpvirtualserver.VolumeShowResponseV1Dot4, error),
) error {
	err := r.client.DetachVolume(ctx, stateVolume.Id.ValueString(), serverId)
	if err != nil {
		return err
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, stateVolume.Id.ValueString(), 100, 3*time.Second,
		"State", "available", "error", getVolumeFunc)
	if err != nil {
		return err
	}

	return r.client.DeleteVolume(ctx, stateVolume.Id.ValueString())
}

func (r *virtualServerServerResource) handlerUpdateServerVolume(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	serverId := plan.Id.ValueString()

	if !plan.BootVolume.Equal(state.BootVolume) {
		planBootVolume, _ := objectToBootVolume(plan.BootVolume)
		stateBootVolume, _ := objectToBootVolume(state.BootVolume)
		if err := r.updateSingleVolume(ctx, serverId, planBootVolume, stateBootVolume); err != nil {
			return err
		}
	}

	if !plan.ExtraVolumes.Equal(state.ExtraVolumes) {
		var planExtraVolumes, stateExtraVolumes map[string]virtualserver.ServerResourceVolume
		plan.ExtraVolumes.ElementsAs(ctx, &planExtraVolumes, false)
		state.ExtraVolumes.ElementsAs(ctx, &stateExtraVolumes, false)

		for key, volume := range planExtraVolumes {
			if stateExtraVolume, exists := stateExtraVolumes[key]; exists {
				if err := r.updateSingleVolume(ctx, serverId, volume, stateExtraVolume); err != nil {
					return err
				}
			}
		}

		getVolumeFunc := func(id string) (*scpvirtualserver.VolumeShowResponseV1Dot4, error) {
			return r.client.GetVolume(ctx, id)
		}

		for key, planExtraVolume := range planExtraVolumes {
			if _, exists := stateExtraVolumes[key]; !exists {
				if err := r.createExtraVolume(ctx, serverId, plan, planExtraVolume, getVolumeFunc); err != nil {
					return err
				}
			}
		}

		for key, stateVolume := range stateExtraVolumes {
			if _, exists := planExtraVolumes[key]; !exists {
				if err := r.deleteExtraVolume(ctx, stateVolume, serverId, getVolumeFunc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *virtualServerServerResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	serviceName, resourceType := r.resolveServerServiceInfoFromModel(state)

	// Server
	_, err := tag.UpdateTags(r.clients, serviceName, resourceType, plan.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}
	// Port
	var networkMap map[string]virtualserver.ServerResourceNetwork
	state.Networks.ElementsAs(ctx, &networkMap, false)
	for _, network := range networkMap {
		_, err := tag.UpdateTags(r.clients, ServiceNameVpc, ResourceTypePort, network.PortId.ValueString(), plan.Tags.Elements(), false)
		if err != nil {
			return err
		}
	}
	// Volume
	stateBootVolume, _ := objectToBootVolume(state.BootVolume)
	_, err = tag.UpdateTags(r.clients, ServiceNameVirtualServer, ResourceTypeVolume, stateBootVolume.Id.ValueString(), plan.Tags.Elements(), false)
	if err != nil {
		return err
	}
	var extraVolumeMap map[string]virtualserver.ServerResourceVolume
	state.ExtraVolumes.ElementsAs(ctx, &extraVolumeMap, false)
	for _, volume := range extraVolumeMap {
		_, err := tag.UpdateTags(r.clients, ServiceNameVirtualServer, ResourceTypeVolume, volume.Id.ValueString(), plan.Tags.Elements(), false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *virtualServerServerResource) getStateTransitions() map[string]map[string]func(ctx context.Context, serverId string) error {
	transitions := make(map[string]map[string]func(ctx context.Context, serverId string) error)

	addState := func(from string, to string, callFunc func(ctx context.Context, serverId string) error) {
		if transitions[from] == nil {
			transitions[from] = make(map[string]func(ctx context.Context, serverId string) error)
		}
		transitions[from][to] = callFunc
	}

	// State Transition Map
	addState("SHUTOFF", "ACTIVE", r.client.StartServer)
	addState("ACTIVE", "SHUTOFF", r.client.StopServer)

	return transitions
}

func (r *virtualServerServerResource) handlerUpdateServerState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	currentState := state.State.ValueString()
	desiredState := plan.State.ValueString()

	if currentState == desiredState {
		return nil
	}

	err := r.getStateTransitions()[currentState][desiredState](ctx, plan.Id.ValueString())
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpvirtualserver.ServerShowResponseV1Dot4, error) {
		return r.client.GetServer(ctx, id)
	}

	_, err = virtualserverutil.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 1000, 3*time.Second,
		"State", desiredState, "ERROR", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *virtualServerServerResource) AsyncPollingServerDeleted(ctx context.Context, serverId string, maxAttempts int, internal time.Duration) error {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := r.client.GetServer(ctx, serverId)

		if err != nil {
			var statusCode float64
			statusCode = virtualserverutil.GetStatusCodeFromError(err)
			if statusCode == 404 {
				return nil
			}
			return fmt.Errorf("attempt %d/%d failed: %w", attempt, maxAttempts, err)
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func (r *virtualServerServerResource) resolveServerServiceInfoFromResponse(response *scpvirtualserver.ServerShowResponseV1Dot4) (serviceName, resourceType string) {
	if response.ProductOffering.Get().Ptr() != nil &&
		(*response.ProductOffering.Get().Ptr() == ProductOfferingGpuServer || *response.ProductOffering.Get().Ptr() == ProductOfferingK8sGpuServer) {
		return ServiceNameGpuServer, ResourceTypeGpuServer
	}
	return ServiceNameVirtualServer, ResourceTypeVirtualServer
}

func (r *virtualServerServerResource) resolveServerServiceInfoFromModel(model virtualserver.ServerResource) (serviceName, resourceType string) {
	if !model.ProductOffering.IsNull() &&
		(model.ProductOffering.ValueString() == ProductOfferingGpuServer || model.ProductOffering.ValueString() == ProductOfferingK8sGpuServer) {
		return ServiceNameGpuServer, ResourceTypeGpuServer
	}
	return ServiceNameVirtualServer, ResourceTypeVirtualServer
}

func (r *virtualServerServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualserver.ServerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.State.IsNull() {
		if plan.State.ValueString() != "ACTIVE" {
			resp.Diagnostics.AddError(
				"Error Creating Server",
				"Invalid server state. Server state must be 'ACTIVE' during creation.\nState: "+plan.State.ValueString())
			return
		}
	}

	data, err := r.client.CreateServer(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if len(data.Servers) == 0 {
		resp.Diagnostics.AddError(
			"Error creating server",
			"CreateServer returned empty Servers list",
		)
		return
	}

	getFunc := func(id string) (*scpvirtualserver.ServerShowResponseV1Dot4, error) {
		return r.client.GetServer(ctx, id)
	}

	getData, err := virtualserverutil.AsyncRequestPollingWithState(ctx, data.Servers[0].Id, 600, 3*time.Second,
		"State", "ACTIVE", "ERROR", getFunc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server",
			"Could not read server, unexpected error: "+err.Error(),
		)
		return
	}

	serviceName, resourceType := r.resolveServerServiceInfoFromResponse(getData)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, data.Servers[0].Id, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	err = r.AsyncPollingVolumeIops(ctx, data.Servers[0].Id, plan.BootVolume, plan.ExtraVolumes)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Volume",
			err.Error(),
		)
		return
	}

	state, err := r.MapGetResponseToState(ctx, getData, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualserver.ServerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetServer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading server",
			"Could not read server name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	serviceName, resourceType := r.resolveServerServiceInfoFromResponse(data)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, state.Id.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, state, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*virtualserver.UpdateHandler{
		{
			Fields:  []string{"Name"},
			Handler: r.handlerUpdateServerName,
		},
		{
			Fields:  []string{"SecurityGroups"},
			Handler: r.handlerUpdateServerSecurityGroup,
		},
		{
			Fields:  []string{"ServerTypeId"},
			Handler: r.handlerUpdateServerType,
		},
		{
			Fields:  []string{"Networks"},
			Handler: r.handlerUpdateServerNetwork,
		},
		{
			Fields:  []string{"BootVolume", "ExtraVolumes"},
			Handler: r.handlerUpdateServerVolume,
		},
		{
			Fields:  []string{"Tags"},
			Handler: r.handlerUpdateTag,
		},
		{
			Fields:  []string{"State"},
			Handler: r.handlerUpdateServerState,
		},
	}

	var plan virtualserver.ServerResource
	var state virtualserver.ServerResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = r.normalizePlan(plan, state)

	var settableFileds []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFileds = append(settableFileds, virtualserverutil.SnakeToPascal(attrName))
		}
	}

	changeFields, err := virtualserverutil.GetChangedFields(plan, state, settableFileds)
	if err != nil {
		return
	}

	if virtualserverutil.IsOverlapFields([]string{"Lock"}, changeFields) {
		if plan.Lock.ValueBool() {
			handlers = append(handlers, &virtualserver.UpdateHandler{
				Fields:  []string{"Lock"},
				Handler: r.handlerUpdateServerLock,
			})
		} else {
			handlers = append([]*virtualserver.UpdateHandler{
				{
					Fields:  []string{"Lock"},
					Handler: r.handlerUpdateServerLock,
				},
			}, handlers...)
		}
	}

	for _, h := range handlers {
		if virtualserverutil.IsOverlapFields(h.Fields, changeFields) {
			if err := h.Handler(ctx, req, resp); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Server",
					"Could not update server, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	data, err := r.client.GetServer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading server",
			"Could not read server name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	serviceName, resourceType := r.resolveServerServiceInfoFromResponse(data)
	tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, state.Id.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualserver.ServerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServer(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting server",
			"Could not delete server, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = r.AsyncPollingServerDeleted(ctx, state.Id.ValueString(), 1000, 3*time.Second)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Server",
			"Error wating for server to become deleted\n"+err.Error(),
		)
		return
	}
}

func (r *virtualServerServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *virtualServerServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan virtualserver.ServerResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.State.Raw.IsNull() {
		return
	}

	var state virtualserver.ServerResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	immutableFields := []string{"ImageId", "KeypairName", "Metadata", "ProductCategory", "ProductOffering", "UserData", "ServerGroupId", "Zone"}

	changeFields, err := virtualserverutil.GetChangedFields(plan, state, immutableFields)
	if err != nil {
		return
	}

	if virtualserverutil.IsOverlapFields(immutableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Server",
			"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
		)
		return
	}
}
