package cloudcontrol

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &accountFactoryResource{}
	_ resource.ResourceWithConfigure   = &accountFactoryResource{}
	_ resource.ResourceWithImportState = &accountFactoryResource{}
)

func NewAccountFactoryResource() resource.Resource {
	return &accountFactoryResource{}
}

type accountFactoryResource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (r *accountFactoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_account_factory"
}

func (r *accountFactoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages CloudControl Account Factory - creates and manages an account in the landing zone.",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Description: "Landing Zone ID where the account will be created. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
				Required: true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the account to create. \n" +
					"  - example : 'foo-account' \n",
				Required: true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Parent Organization Unit ID where the account will be placed. \n" +
					"  - example : 'ou-b30e9fcc39f84a20bf9e7458e5ec3801' \n",
				Required: true,
			},
			"email": schema.StringAttribute{
				Description: "Email address for the account. \n" +
					"  - example : 'myaccount@example.com' \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sso_user_email": schema.StringAttribute{
				Description: "SSO user email address. \n" +
					"  - example : 'user_email@example.com' \n",
				Optional: true,
			},
			"sso_user_name": schema.StringAttribute{
				Description: "Username of the Identity Center user. \n" +
					"  - example : 'John Doe' \n",
				Optional: true,
			},
			"sso_user_real_name": schema.StringAttribute{
				Description: "Real name of the SSO User. \n" +
					"  - example : 'John Doe' \n",
				Optional: true,
			},
			"account_id": schema.StringAttribute{
				Computed: true,
				Description: "Account ID in the organization. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
			},
			"job_id": schema.StringAttribute{
				Computed: true,
				Description: "Job ID for the account creation operation. \n" +
					"  - example : '0a36e0746dbf4908acf0357829701381' \n",
			},
		},
	}
}

func (r *accountFactoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.CloudControl
	r.clients = inst.Client
}

func (r *accountFactoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudcontrol.AccountFactoryResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateAccountFactoryAccount(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Account Factory Account",
			"Could not create Account Factory Account, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.AccountId = types.StringValue(data.AccountId)
	plan.JobId = types.StringValue(data.JobId)

	err = waitForAccountFactoryActive(ctx, r.clients.Organization, data.AccountId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for account creation",
			"Error waiting for account "+data.AccountId+" to become active: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *accountFactoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudcontrol.AccountFactoryResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := state.AccountId.ValueString()
	if accountId == "" {
		accountId = state.AccountId.ValueString()
	}

	if accountId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Account Factory Account",
			"Account ID is empty",
		)
		return
	}

	data, err := r.clients.Organization.GetAccount(ctx, accountId, "", nil)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.AccountId = types.StringValue(data.Account.Id)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accountFactoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported Operation",
		"Update is not supported for Account Factory resource",
	)
}

func (r *accountFactoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Unsupported Operation",
		"Delete is not supported for Account Factory resource. The account will remain in the cloud.",
	)
}

func (r *accountFactoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	accountId := req.ID
	if accountId == "" {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"The import ID cannot be empty",
		)
		return
	}

	state := cloudcontrol.AccountFactoryResource{
		AccountId:     types.StringValue(accountId),
		LandingZoneId: types.StringValue(""),
		Name:          types.StringValue(""),
		ParentUnitId:  types.StringValue(""),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func waitForAccountFactoryActive(ctx context.Context, orgClient *organization.Client, accountId string) error {
	return client.WaitForStatus(ctx, nil, []string{"PENDING"}, []string{"ACTIVE", "RUNNING"}, func() (interface{}, string, error) {
		accountResp, err := orgClient.GetAccount(ctx, accountId, "", nil)
		if err != nil {
			return nil, "", err
		}
		state := string(accountResp.Account.State)
		if state == "CREATE_FAILED" {
			return accountResp, state, fmt.Errorf("account creation failed")
		}
		return accountResp, state, nil
	}, -1, -1, -1, -1)
}
