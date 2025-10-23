package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"time"
)

var (
	_ resource.Resource              = &virtualServerKeypairResource{}
	_ resource.ResourceWithConfigure = &virtualServerKeypairResource{}
)

func NewVirtualServerKeypairResource() resource.Resource {
	return &virtualServerKeypairResource{}
}

type virtualServerKeypairResource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (r *virtualServerKeypairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_keypair"
}

func (r *virtualServerKeypairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "keypair",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			common.ToSnakeCase("PublicKey"): schema.StringAttribute{
				Description: "Public key",
				Computed:    true,
			},
			common.ToSnakeCase("Fingerprint"): schema.StringAttribute{
				Description: "Fingerprint",
				Computed:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Type",
				Computed:    true,
			},
			common.ToSnakeCase("PrivateKey"): schema.StringAttribute{
				Description: "Private key",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "Created at",
				Computed:    true,
			},
			common.ToSnakeCase("UserId"): schema.StringAttribute{
				Description: "User ID",
				Computed:    true,
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

func (r *virtualServerKeypairResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *virtualServerKeypairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualserver.KeypairResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateKeypair(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating keypair",
			"Could not create keypair, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	getData, err := r.client.GetKeypair(ctx, data.Name)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading keypair",
			"Could not read keypair, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeKeypair, strconv.Itoa(int(getData.Id)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	keypair := getData

	keypairModel := virtualserver.KeypairResource{
		Id:          types.Int32Value(keypair.Id),
		Name:        types.StringValue(keypair.Name),
		PublicKey:   types.StringValue(keypair.PublicKey),
		Fingerprint: types.StringValue(keypair.Fingerprint),
		Type:        types.StringValue(keypair.Type),
		PrivateKey:  types.StringPointerValue(data.PrivateKey.Get()),
		CreatedAt:   types.StringValue(keypair.CreatedAt.Format(time.RFC3339)),
		UserId:      types.StringValue(keypair.UserId),
		Tags:        tagsMap,
	}
	state := keypairModel

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerKeypairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualserver.KeypairResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetKeypair(ctx, state.Name.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading keypair",
			"Could not read keypair name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeKeypair, state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	keypair := data

	keypairModel := virtualserver.KeypairResource{
		Name:        types.StringValue(keypair.Name),
		PublicKey:   types.StringValue(keypair.PublicKey),
		Fingerprint: types.StringValue(keypair.Fingerprint),
		Type:        types.StringValue(keypair.Type),
		Id:          types.Int32Value(keypair.Id),
		UserId:      types.StringValue(keypair.UserId),
		CreatedAt:   types.StringValue(keypair.CreatedAt.Format(time.RFC3339)),
		PrivateKey:  types.StringPointerValue(state.PrivateKey.ValueStringPointer()),
		Tags:        tagsMap,
	}
	newState := keypairModel

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerKeypairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan virtualserver.KeypairResource
	var state virtualserver.KeypairResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	immuntableFields := []string{"Name"}

	if virtualserverutil.IsOverlapFields(immuntableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Keypair",
			"Immutable fields cannot be modified: "+strings.Join(immuntableFields, ", "),
		)
		return
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, ServiceNameVirtualServer, ResourceTypeKeypair, plan.Id.String(), tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tags",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state.Tags = tagsMap

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerKeypairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualserver.KeypairResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteKeypair(ctx, state.Name.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting keypair",
			"Could not delete keypair, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
