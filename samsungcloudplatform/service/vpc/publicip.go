package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/vpc/1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcPublicipResource{}
	_ resource.ResourceWithConfigure = &vpcPublicipResource{}
)

// NewVpcPublicipResource is a helper function to simplify the provider implementation.
func NewVpcPublicipResource() resource.Resource {
	return &vpcPublicipResource{}
}

// vpcPublicipResource is the data source implementation.
type vpcPublicipResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcPublicipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_publicip"
}

// Schema defines the schema for the data source.
func (r *vpcPublicipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "publicip",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Type \n" +
					"  - example : IGW | GGW | SIGW",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Public IP description\n" +
					"  - maxLength : 50\n" +
					"  - minLength : 1",
				Optional: true,
			},
			common.ToSnakeCase("Publicip"): schema.SingleNestedAttribute{
				Description: "Publicip",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "IpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
						Description: "AttachedResourceType",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
						Description: "AttachedResourceName",
						Computed:    true,
					},
					common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
						Description: "AttachedResourceId",
						Computed:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
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
func (r *vpcPublicipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPublicipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.PublicipResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new publicip
	data, err := r.client.CreatePublicip(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating publicip",
			"Could not create publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Publicip.Id)
	publicipModel := createPublicipModel(data)
	publicipObjectValue, diags := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	plan.Publicip = publicipObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPublicipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.PublicipResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from publicip
	data, err := r.client.GetPublicip(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading publicip",
			"Could not read publicip ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicipModel := createPublicipModel(data)
	publicipObjectValue, diags := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	state.Publicip = publicipObjectValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPublicipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var state vpc.PublicipResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdatePublicip(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating publicip",
			"Could not update publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetPublicip as UpdatePublicip items are not populated.
	data, err := r.client.GetPublicip(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading publicip",
			"Could not read publicip ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	publicipModel := createPublicipModel(data)
	publicipObjectValue, diags := types.ObjectValueFrom(ctx, publicipModel.AttributeTypes(), publicipModel)
	state.Publicip = publicipObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPublicipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.PublicipResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing publicip
	err := r.client.DeletePublicip(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting publicip",
			"Could not delete publicip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createPublicipModel(data *scpvpc.PublicipShowResponse) vpc.Publicip {
	publicip := data.Publicip
	publicipModel := vpc.Publicip{
		IpAddress:            types.StringValue(data.GetPublicip().IpAddress),
		AccountId:            types.StringValue(data.GetPublicip().AccountId),
		AttachedResourceName: types.StringPointerValue(data.GetPublicip().AttachedResourceName.Get()),
		AttachedResourceId:   types.StringPointerValue(data.GetPublicip().AttachedResourceId.Get()),
		Type:                 types.StringValue(string(data.GetPublicip().Type)),
		State:                types.StringValue(string(data.GetPublicip().State)),
		Description:          types.StringPointerValue(data.GetPublicip().Description.Get()),
		CreatedAt:            types.StringValue(data.GetPublicip().CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.GetPublicip().CreatedBy),
		ModifiedAt:           types.StringValue(data.GetPublicip().ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.GetPublicip().ModifiedBy),
	}
	attachedResourceType := publicip.AttachedResourceType.Get()
	if attachedResourceType != nil {
		attachedResourceTypeStr := string(*attachedResourceType)
		publicipModel.AttachedResourceType = types.StringPointerValue(&attachedResourceTypeStr)
	} else {
		publicipModel.AttachedResourceType = types.StringPointerValue(nil)
	}
	return publicipModel
}
