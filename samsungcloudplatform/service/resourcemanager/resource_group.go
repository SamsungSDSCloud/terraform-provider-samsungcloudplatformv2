package resourcemanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceManagerResourceGroupResource{}
	_ resource.ResourceWithConfigure   = &resourceManagerResourceGroupResource{}
	_ resource.ResourceWithImportState = &resourceManagerResourceGroupResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewResourceManagerResourceGroupResource() resource.Resource {
	return &resourceManagerResourceGroupResource{}
}

func (r *resourceManagerResourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID 기반 단일 리소스 import 시 표준 패턴
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// resourceManagerResourceGroupResource is the data source implementation.
type resourceManagerResourceGroupResource struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *resourceManagerResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_resource_group"
}

// Schema defines the schema for the data source.
func (r *resourceManagerResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource Group",
		Attributes: map[string]schema.Attribute{
			"region": region.ResourceSchema(),
			"tags":   tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the resource group.",
				MarkdownDescription: "The unique identifier of the resource group.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description:         "Timestamp of the last Terraform update of the Resource Group",
				MarkdownDescription: "The timestamp of the last terraform update of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
				Computed:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Name of the resource group.",
				MarkdownDescription: "The name of the resource group.\n\nExample: `example-rg`",
				Optional:            true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description:         "Description of the resource group.",
				MarkdownDescription: "The description of the resource group.\n\nExample: `This is a discription of resource group`",
				Optional:            true,
			},
			common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "List of resource types.",
				MarkdownDescription: "A list of resource types associated with the group.\n\nExample: `[\"virtual-server\"]`",
				Optional:            true,
			},
			common.ToSnakeCase("GroupDefinitionTags"): schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Group Definition Tags",
				MarkdownDescription: "Creates a resource group for resources with the specified tag.\n\nExample: `[{Key: Environment, Value: Production}]`",
			},
			common.ToSnakeCase("ResourceGroup"): schema.SingleNestedAttribute{
				Description:         "List of group definition tags.",
				MarkdownDescription: "A list of key-value pairs representing group definition tags.\n\nExample: `[{Key: Environment, Value: Production}]`",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description:         "Creation timestamp.",
						MarkdownDescription: "The creation timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description:         "Creator identifier.",
						MarkdownDescription: "The user ID of the creator.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description:         "Description of the resource group.",
						MarkdownDescription: "The description of the resource group.\n\nExample: `This is a discription of resource group`",
						Computed:            true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "The unique identifier of the resource group.",
						MarkdownDescription: "The unique identifier of the resource group.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description:         "Modification timestamp.",
						MarkdownDescription: "The modification timestamp of the resource group.\n\nExample: `2023-10-27T10:00:00Z`",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description:         "Modifier identifier.",
						MarkdownDescription: "The user ID of the modifier.\n\nExample: `e4b2c3f8a1d94b6b9f7e8c2d3a4f5b67`",
						Computed:            true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Name of the resource group.",
						MarkdownDescription: "The name of the resource group.\n\nExample: `example-rg`",
						Computed:            true,
					},
					common.ToSnakeCase("Region"): schema.StringAttribute{
						Description:         "Region code.",
						MarkdownDescription: "The region code where the resource group is located.\n\nExample: `kr-west1`",
						Computed:            true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description:         "System Resource Name.",
						MarkdownDescription: "The System Resource Name (SRN) of the resource group.\n\nExample: `srn:s::13d97ad943ca452481d624f78391df13:kr-west1::resourcemanager:resource-group/70636f984e564b3c9e54e74a53f9318d`",
						Computed:            true,
					},
					common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
						ElementType:         types.StringType,
						Description:         "List of resource types.",
						MarkdownDescription: "A list of resource types associated with the group.\n\nExample: `[\"virtual-server\"]`",
						Computed:            true,
					},
					common.ToSnakeCase("GroupDefinitionTags"): schema.ListNestedAttribute{
						Description:         "List of group definition tags.",
						MarkdownDescription: "A list of key-value pairs representing group definition tags.\n\nExample: `[{Key: Environment, Value: Production}]`",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Key"): schema.StringAttribute{
									Description:         "Tag key.",
									MarkdownDescription: "The key of the tag.\n\nExample: `Environment`",
									Computed:            true,
								},
								common.ToSnakeCase("Value"): schema.StringAttribute{
									Description:         "Tag value.",
									MarkdownDescription: "The value of the tag.\n\nExample: `Production`",
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
func (r *resourceManagerResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.ResourceManager
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceManagerResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resourcemanager.ResourceGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Resource Group
	data, err := r.client.CreateResourceGroup(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Resource Group",
			"Could not create Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(resourceGroup.Id)

	// ID가 유효한 경우에만 폴링
	if !plan.Id.IsNull() && plan.Id.ValueString() != "" {
		// 생성 후 리소스가 준비될 때까지 대기
		err = r.waitForResourceGroupReady(ctx, plan.Id.ValueString(), 60*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for resource group",
				"Resource Group was created but failed to become ready: "+err.Error(),
			)
			return
		}

		// 최신 상태 다시 조회
		// 1. 변수명을 showData로 변경하여 타입 충돌 방지
		showData, err := r.client.GetResourceGroup(ctx, plan.Id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error reading resource group after creation",
				// 2. plan.Id를 문자열로 변환 (ValueString 사용)
				"Could not read Resource Group ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		// 3. 변수명 동기화
		resourceGroup = showData.ResourceGroup
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id, tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tags",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}
	resourceTypesObject, diags := types.ListValueFrom(ctx, types.StringType, resourceTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	plan.ResourceGroup = resourceGroupObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceManagerResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state resourcemanager.ResourceGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Resource Group
	data, err := r.client.GetResourceGroup(ctx, state.Id.ValueString())

	if err != nil {
		// 404 Not Found 감지 시 상태 제거 (SDK 에러 처리 방식에 따라 조건 수정 필요)
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			"Could not read Resource Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup

	// Get Tags
	tagsMap, err := tag.GetTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}

	resourceTypesObject, diags := types.ListValueFrom(ctx, types.StringType, resourceTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	state.ResourceGroup = resourceGroupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceManagerResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state resourcemanager.ResourceGroupResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Resource Group
	_, err := r.client.UpdateResourceGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Resource Group",
			"Could not update Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetResourceGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading resourceGroup",
			"Could not read resourceGroup ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup

	tagElements := state.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id, tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tags",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}
	resourceTypesObject, diags := types.ListValueFrom(ctx, types.StringType, resourceTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	state.ResourceGroup = resourceGroupObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceManagerResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state resourcemanager.ResourceGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag.UpdateTags(r.clients, "resourcemanager", "resource-group", state.Id.ValueString(), make(map[string]attr.Value))

	// Delete existing Resource Group
	err := r.client.DeleteResourceGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Resource Group",
			"Could not delete Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *resourceManagerResourceGroupResource) waitForResourceGroupReady(ctx context.Context, id string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timeout waiting for resource group %s to be ready", id)
		case <-ticker.C:
			data, err := r.client.GetResourceGroup(ctx, id)
			if err != nil {
				// 404 Not Found - 생성 중일 수 있음, 계속 대기
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
					continue
				}
				// 일시적 오류는 무시하고 계속 폴링 (네트워크 문제 등)
				if client.IsTransientError(err) {
					continue
				}
				// 그 외의 영구적인 오류는 반환
				return err
			}

			// 상태 확인 - 실제 API 응답에 맞게 수정 필요
			resourceGroup := data.ResourceGroup
			if resourceGroup.Id != "" {
				return nil
			}
		}
	}
}
