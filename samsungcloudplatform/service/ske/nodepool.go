package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/ske/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &skeNodepoolResource{}
	_ resource.ResourceWithConfigure = &skeNodepoolResource{}
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
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the nodepool",
				Computed:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			common.ToSnakeCase("ClusterId"): schema.StringAttribute{
				Description: "ClusterId",
				Required:    true,
			},
			common.ToSnakeCase("CustomImageId"): schema.StringAttribute{
				Description: "CustomImageId",
				Optional:    true,
			},
			common.ToSnakeCase("DesiredNodeCount"): schema.Int32Attribute{
				Description: "DesiredNodeCount",
				Computed:    true,
				Optional:    true,
			},
			common.ToSnakeCase("ImageOs"): schema.StringAttribute{
				Description: "ImageOs",
				Required:    true,
			},
			common.ToSnakeCase("ImageOsVersion"): schema.StringAttribute{
				Description: "ImageOsVersion",
				Required:    true,
			},
			common.ToSnakeCase("IsAutoRecovery"): schema.BoolAttribute{
				Description: "IsAutoRecovery",
				Required:    true,
			},
			common.ToSnakeCase("IsAutoScale"): schema.BoolAttribute{
				Description: "IsAutoScale",
				Required:    true,
			},
			common.ToSnakeCase("KeypairName"): schema.StringAttribute{
				Description: "KeypairName",
				Required:    true,
			},
			common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
				Description: "KubernetesVersion",
				Required:    true,
			},
			common.ToSnakeCase("Labels"): schema.ListNestedAttribute{
				Description: "Labels",
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.makeNodepoolLabelsRequestSchema(),
				},
				Optional: true,
			},
			common.ToSnakeCase("Taints"): schema.ListNestedAttribute{
				Description: "Taints",
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.makeNodepoolTaintsRequestSchema(),
				},
				Optional: true,
			},
			common.ToSnakeCase("MaxNodeCount"): schema.Int32Attribute{
				Description: "MaxNodeCount",
				Computed:    true,
				Optional:    true,
			},
			common.ToSnakeCase("MinNodeCount"): schema.Int32Attribute{
				Description: "MinNodeCount",
				Computed:    true,
				Optional:    true,
			},
			common.ToSnakeCase("ServerTypeId"): schema.StringAttribute{
				Description: "ServerTypeId",
				Required:    true,
			},
			common.ToSnakeCase("VolumeTypeName"): schema.StringAttribute{
				Description: "VolumeTypeName",
				Required:    true,
			},
			common.ToSnakeCase("VolumeSize"): schema.Int32Attribute{
				Description: "VolumeSize",
				Required:    true,
			},
			common.ToSnakeCase("NodepoolDetail"): schema.SingleNestedAttribute{
				Description: "NodepoolDetail",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "account id",
						Computed:    true,
					},
					common.ToSnakeCase("AutoRecoveryEnabled"): schema.BoolAttribute{
						Description: "AutoRecoveryEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("AutoScaleEnabled"): schema.BoolAttribute{
						Description: "AutoScaleEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
						Description: "Cluster",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("CurrentNodeCount"): schema.Int32Attribute{
						Description: "CurrentNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("DesiredNodeCount"): schema.Int32Attribute{
						Description: "DesiredNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("Image"): schema.SingleNestedAttribute{
						Description: "Image",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("CustomImageName"): schema.StringAttribute{
								Description: "CustomImageName",
								Computed:    true,
							},
							common.ToSnakeCase("Os"): schema.StringAttribute{
								Description: "Os",
								Computed:    true,
							},
							common.ToSnakeCase("OsVersion"): schema.StringAttribute{
								Description: "OsVersion",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Keypair"): schema.SingleNestedAttribute{
						Description: "Keypair",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description: "Name",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
						Description: "KubernetesVersion",
						Computed:    true,
					},
					common.ToSnakeCase("Labels"): schema.ListNestedAttribute{
						Description: "Labels",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: r.makeNodepoolLabelsSchema(),
						},
					},
					common.ToSnakeCase("Taints"): schema.ListNestedAttribute{
						Description: "Taints",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: r.makeNodepoolTaintsSchema(),
						},
					},
					common.ToSnakeCase("MaxNodeCount"): schema.Int32Attribute{
						Description: "MaxNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("MinNodeCount"): schema.Int32Attribute{
						Description: "MinNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("ServerType"): schema.SingleNestedAttribute{
						Description: "ServerType",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Description"): schema.StringAttribute{
								Description: "Description",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					common.ToSnakeCase("VolumeType"): schema.SingleNestedAttribute{
						Description: "VolumeType",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Encrypt"): schema.BoolAttribute{
								Description: "Encrypt",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description: "Name",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("VolumeSize"): schema.Int32Attribute{
						Description: "VolumeSize",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
				},
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
func (r *skeNodepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	// Retrieve values from plan
	var plan ske.NodepoolResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new nodepool
	data, err := r.client.CreateNodepool(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating nodepool",
			"Could not create nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Nodepool.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	if plan.IsAutoScale.ValueBool() {
		plan.DesiredNodeCount = types.Int32Value(data.Nodepool.MinNodeCount)
	} else {
		plan.MinNodeCount = types.Int32Value(data.Nodepool.DesiredNodeCount)
		plan.MaxNodeCount = types.Int32Value(data.Nodepool.DesiredNodeCount)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.Nodepool.Id, []string{"ScalingUp"}, []string{"Running"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating nodepool",
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
	data, _, err := r.client.GetNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read Nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := createNodepoolDetailModel(data)

	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)
	state.NodepoolDetail = nodepoolObjectValue

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
				_, err := r.client.UpgradeNodepool(ctx, plan.Id.ValueString())
				if err != nil {
					detail := client.GetDetailFromError(err)
					resp.Diagnostics.AddError(
						"Error upgrade nodepool",
						"Could not upgrade nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
					)
					return
				}

				err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"Updating"}, []string{"Running"}, true)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error updating nodepool",
						"Error waiting for nodepool to become running: "+err.Error(),
					)
					return
				}
			} else {
				resp.Diagnostics.AddError(
					"Error updating nodepool version",
					"When nodepool version update, must not modify node count",
				)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Error updating nodepool version",
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
					"Could not update Nodepool Labels, unexpected error: "+err.Error()+"\nReason: "+detail,
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
					"Could not update Nodepool Taints, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if (plan.IsAutoScale.ValueBool() && (!state.MinNodeCount.Equal(plan.MinNodeCount) || !state.MaxNodeCount.Equal(plan.MaxNodeCount))) ||
			(!plan.IsAutoScale.ValueBool() && !state.DesiredNodeCount.Equal(plan.DesiredNodeCount)) ||
			!state.IsAutoRecovery.Equal(plan.IsAutoRecovery) {
			_, err := r.client.UpdateNodepool(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating nodepool",
					"Could not update nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"ScalingUp", "ScalingDown"}, []string{"Running"}, true)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating nodepool",
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

	nodepoolModel := createNodepoolDetailModel(data)
	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.NodepoolDetail = nodepoolObjectValue
	if plan.IsAutoScale.ValueBool() {
		plan.DesiredNodeCount = types.Int32Value(data.Nodepool.MinNodeCount)
	} else {
		plan.MinNodeCount = types.Int32Value(data.Nodepool.DesiredNodeCount)
		plan.MaxNodeCount = types.Int32Value(data.Nodepool.DesiredNodeCount)
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
			"Could not delete Nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.ResourceId, []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting nodepool",
			"Error waiting for nodepool to become deleted: "+err.Error(),
		)
		return
	}
}

func createNodepoolDetailModel(data *scpske.NodepoolShowResponse) ske.NodepoolDetail {
	nodepoolElement := data.Nodepool
	nodepoolModel := ske.NodepoolDetail{
		Id:                  types.StringValue(nodepoolElement.Id),
		Name:                types.StringValue(nodepoolElement.Name),
		AccountId:           types.StringValue(nodepoolElement.AccountId),
		AutoRecoveryEnabled: types.BoolValue(nodepoolElement.AutoRecoveryEnabled),
		AutoScaleEnabled:    types.BoolValue(nodepoolElement.AutoScaleEnabled),
		Cluster: ske.IdMapType{
			Id: types.StringValue(nodepoolElement.Cluster.Id),
		},
		CurrentNodeCount: types.Int32Value(nodepoolElement.CurrentNodeCount),
		DesiredNodeCount: types.Int32Value(nodepoolElement.DesiredNodeCount),
		Image: ske.Image{
			CustomImageName: types.StringPointerValue(nodepoolElement.Image.CustomImageName.Get()),
			Os:              types.StringValue(nodepoolElement.Image.Os),
			OsVersion:       types.StringValue(nodepoolElement.Image.OsVersion),
		},
		Keypair: ske.NameMapType{
			Name: types.StringValue(nodepoolElement.Keypair.Name),
		},
		KubernetesVersion: types.StringValue(nodepoolElement.KubernetesVersion),
		MaxNodeCount:      types.Int32Value(nodepoolElement.MaxNodeCount),
		MinNodeCount:      types.Int32Value(nodepoolElement.MinNodeCount),
		ServerType: ske.ServerType{
			Description: types.StringValue(nodepoolElement.ServerType.Description),
			Id:          types.StringValue(nodepoolElement.ServerType.Id),
		},
		Status: types.StringValue(nodepoolElement.Status),
		VolumeType: ske.VolumeType{
			Encrypt: types.BoolValue(nodepoolElement.VolumeType.Encrypt),
			Id:      types.StringValue(nodepoolElement.VolumeType.Id),
			Name:    types.StringValue(nodepoolElement.VolumeType.Name),
		},
		VolumeSize: types.Int32Value(nodepoolElement.VolumeSize),
		CreatedAt:  types.StringValue(nodepoolElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringValue(nodepoolElement.CreatedBy),
		ModifiedAt: types.StringValue(nodepoolElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringValue(nodepoolElement.ModifiedBy),
	}

	return nodepoolModel
}

func waitForNodepoolStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetNodepool(ctx, id)
		if httpStatus == 200 {
			return info, info.Nodepool.Status, nil
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

		return info, info.Nodepool.Status, nil
	})
}

func (r *skeNodepoolResource) makeNodepoolLabelsRequestSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Required:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Optional:    true,
		},
	}
}

func (r *skeNodepoolResource) makeNodepoolTaintsRequestSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Effect"): schema.StringAttribute{
			Description: "Effect",
			Required:    true,
		},
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Required:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Optional:    true,
		},
	}
}

func (r *skeNodepoolResource) makeNodepoolLabelsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Computed:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Computed:    true,
		},
	}
}

func (r *skeNodepoolResource) makeNodepoolTaintsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Effect"): schema.StringAttribute{
			Description: "Effect",
			Computed:    true,
		},
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Computed:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Computed:    true,
		},
	}
}
