package samsungcloudplatform

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/config"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/service"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sysuser "os/user"
	"path/filepath"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &samsungcloudplatformv2Provider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &samsungcloudplatformv2Provider{
			version: version,
		}
	}
}

// samsungcloudplatformv2Provider is the provider implementation.
type samsungcloudplatformv2Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *samsungcloudplatformv2Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "samsungcloudplatformv2"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *samsungcloudplatformv2Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with samsungcloudplatformv2.",
		Attributes: map[string]schema.Attribute{
			"auth_url": schema.StringAttribute{
				Description: "Authentication URL for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_AUTH_URL environment variable.",
				Optional:    true,
			},
			"endpoint_override": schema.StringAttribute{
				Description: "Endpoint override configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ENDPOINT_OVERRIDE environment variable.",
				Optional:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ACCOUNT_ID environment variable.",
				Optional:    true,
			},
			"default_region": schema.StringAttribute{
				Description: "Default region configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_DEFAULT_REGION environment variable.",
				Optional:    true,
			},
			"access_key": schema.StringAttribute{
				Description: "Access key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ACCESS_KEY environment variable.",
				Optional:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "Secret key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_SECRET_KEY environment variable.",
				Optional:    true,
			},
			"auth_token": schema.StringAttribute{
				Description: "Auth token for calling samsungcloudplatformv2 API",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a samsungcloudplatformv2 API client for data sources and resources.
func (p *samsungcloudplatformv2Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring samsungcloudplatformv2 client")

	// Retrieve provider data from configuration
	providerConfig := config.ProviderConfig{}
	diags := req.Config.Get(ctx, &providerConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := sysuser.Current()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve current user",
			"Error details: "+err.Error(),
		)
		return
	}

	config.LoadServiceConfig(resp, filepath.Join(user.HomeDir, ".scpconf", config.ServiceConfigFile), &providerConfig)
	if resp.Diagnostics.HasError() {
		return
	}

	config.LoadCredentialsConfig(resp, filepath.Join(user.HomeDir, ".scpconf", config.CredentialConfigFile), &providerConfig)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ConfigureServiceAndCredentials(resp, &providerConfig)
	if resp.Diagnostics.HasError() {
		return
	}

	scpClient, err := client.NewSCPClient(&providerConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Samsungcloudplatform EngineClient",
			"Error details: "+err.Error(),
		)
		return
	}

	inst := client.Instance{
		Client: scpClient,
	}

	// Make the samsungcloudplatformv2 client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = inst
	resp.ResourceData = inst

	tflog.Info(ctx, "Configured samsungcloudplatformv2 client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *samsungcloudplatformv2Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return service.DataSourceConstructors
}

// Resources defines the resources implemented in the provider.
func (p *samsungcloudplatformv2Provider) Resources(_ context.Context) []func() resource.Resource {
	return service.ResourceConstructors
}
