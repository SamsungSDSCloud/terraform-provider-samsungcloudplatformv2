package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// ProviderConfig maps provider schema data to a Go type.
type ProviderConfig struct {
	AuthUrl                  types.String `tfsdk:"auth_url"`
	EndpointOverride         types.String `tfsdk:"endpoint_override"`
	AccountId                types.String `tfsdk:"account_id"`
	DefaultRegion            types.String `tfsdk:"default_region"`
	AccessKey                types.String `tfsdk:"access_key"`
	SecretKey                types.String `tfsdk:"secret_key"`
	AuthToken                types.String `tfsdk:"auth_token"`
	MaxRemainDays            types.Int64  `tfsdk:"max_remain_days"`
	MicroversionCheckTimeout types.Int64  `tfsdk:"microversion_check_timeout"`
}

const ServiceConfigFile = "config.json"
const CredentialConfigFile = "credentials.json"
const DefaultMaxRemainDays = 90
const DefaultMicroversionCheckTimeout = 15

func LoadServiceConfig(resp *provider.ConfigureResponse, path string, providerConfig *ProviderConfig) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil || len(data) == 0 {
		data = []byte("{}")
	}

	tempConfig := struct {
		AuthUrl                  string `json:"auth-url"`
		EndpointOverride         string `json:"endpoint-override"`
		AccountId                string `json:"account-id"`
		DefaultRegion            string `json:"default-region"`
		MaxRemainDays            int64  `json:"max-remain-days"`
		MicroversionCheckTimeout int64  `json:"microversion-check-timeout"`
	}{}

	if err := json.Unmarshal(data, &tempConfig); err != nil {
		resp.Diagnostics.AddError(
			"Unable to load service configuration file",
			"Error details: "+err.Error(),
		)
		return
	}

	if tempConfig.EndpointOverride != "" {
		providerConfig.EndpointOverride = types.StringValue(tempConfig.EndpointOverride)
	}
	if tempConfig.AuthUrl != "" {
		providerConfig.AuthUrl = types.StringValue(tempConfig.AuthUrl)
	}
	if tempConfig.AccountId != "" {
		providerConfig.AccountId = types.StringValue(tempConfig.AccountId)
	}
	if tempConfig.DefaultRegion != "" {
		providerConfig.DefaultRegion = types.StringValue(tempConfig.DefaultRegion)
	}
	if tempConfig.MaxRemainDays != 0 {
		providerConfig.MaxRemainDays = types.Int64Value(tempConfig.MaxRemainDays)
	}
	if tempConfig.MicroversionCheckTimeout != 0 {
		providerConfig.MicroversionCheckTimeout = types.Int64Value(tempConfig.MicroversionCheckTimeout)
	}
}

func LoadCredentialsConfig(resp *provider.ConfigureResponse, path string, providerConfig *ProviderConfig) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil || len(data) == 0 {
		data = []byte("{}")
	}

	tempConfig := struct {
		AccessKey string `json:"access-key"`
		SecretKey string `json:"secret-key"`
		AuthToken string `json:"auth-token"`
	}{}

	if err := json.Unmarshal(data, &tempConfig); err != nil {
		resp.Diagnostics.AddError(
			"Unable to load credential configuration file",
			"Error details: "+err.Error(),
		)
		return
	}

	if tempConfig.AccessKey != "" {
		providerConfig.AccessKey = types.StringValue(tempConfig.AccessKey)
	}
	if tempConfig.SecretKey != "" {
		providerConfig.SecretKey = types.StringValue(tempConfig.SecretKey)
	}
	if tempConfig.AuthToken != "" {
		providerConfig.AuthToken = types.StringValue(tempConfig.AuthToken)
	}
}

func ConfigureServiceAndCredentials(resp *provider.ConfigureResponse, providerConfig *ProviderConfig) {
	if providerConfig.AuthUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_url"),
			"Unknown samsungcloudplatformv2 Auth URL",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 AUTH URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_AUTH_URL environment variable.",
		)
	}
	if providerConfig.EndpointOverride.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint override"),
			"Unknown samsungcloudplatformv2 Endpoint Override",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 endpoint override. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_ENDPOINT_OVERRIDE environment variable.",
		)
	}
	if providerConfig.AccountId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("account_id"),
			"Unknown samsungcloudplatformv2 Account Id",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Account Id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_ACCOUNT_ID environment variable.",
		)
	}
	if providerConfig.DefaultRegion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("default_region"),
			"Unknown samsungcloudplatformv2 Default Region",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Default Region. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_DEFAULT_REGION environment variable.",
		)
	}
	if providerConfig.AccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Unknown samsungcloudplatformv2 Access Key",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Access Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_ACCESS_KEY environment variable.",
		)
	}
	if providerConfig.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Unknown samsungcloudplatformv2 Secret Key",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Secret Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_SECRET_KEY environment variable.",
		)
	}
	if providerConfig.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("max_remain_days"),
			"Unknown samsungcloudplatformv2 Max Remain Days",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Max Remain Days. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_MAX_REMAIN_DAYS environment variable.",
		)
	}
	if providerConfig.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("microversion_check_timeout"),
			"Unknown samsungcloudplatformv2 Microversion Check Timeout",
			"The provider cannot create the samsungcloudplatformv2 API client as there is an unknown configuration value for the samsungcloudplatformv2 Microversion Check Timeout. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SCP_TF_MICROVERSION_CHECK_TIMEOUT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	authUrl := os.Getenv("SCP_TF_AUTH_URL")
	endpointOverride := os.Getenv("SCP_TF_ENDPOINT_OVERRIDE")
	AccountId := os.Getenv("SCP_TF_ACCOUNT_ID")
	DefaultRegion := os.Getenv("SCP_TF_DEFAULT_REGION")
	accessKey := os.Getenv("SCP_TF_ACCESS_KEY")
	secretKey := os.Getenv("SCP_TF_SECRET_KEY")
	authToken := os.Getenv("SCP_TF_AUTH_TOKEN")
	maxRemainDays, _ := strconv.ParseInt(os.Getenv("SCP_TF_MAX_REMAIN_DAYS"), 10, 32)
	microversionCheckTimeout, _ := strconv.ParseInt(os.Getenv("SCP_TF_MICROVERSION_CHECK_TIMEOUT"), 10, 32)
	if !providerConfig.AuthUrl.IsNull() {
		authUrl = providerConfig.AuthUrl.ValueString()
	}
	if !providerConfig.EndpointOverride.IsNull() {
		endpointOverride = providerConfig.EndpointOverride.ValueString()
	}
	if !providerConfig.AccountId.IsNull() {
		AccountId = providerConfig.AccountId.ValueString()
	}
	if !providerConfig.DefaultRegion.IsNull() {
		DefaultRegion = providerConfig.DefaultRegion.ValueString()
	}
	if !providerConfig.AccessKey.IsNull() {
		accessKey = providerConfig.AccessKey.ValueString()
	}
	if !providerConfig.SecretKey.IsNull() {
		secretKey = providerConfig.SecretKey.ValueString()
	}
	if !providerConfig.AuthToken.IsNull() {
		authToken = providerConfig.AuthToken.ValueString()
	}
	if !providerConfig.MaxRemainDays.IsNull() {
		maxRemainDays = providerConfig.MaxRemainDays.ValueInt64()
	} else {
		maxRemainDays = DefaultMaxRemainDays
	}
	if !providerConfig.MicroversionCheckTimeout.IsNull() {
		microversionCheckTimeout = providerConfig.MicroversionCheckTimeout.ValueInt64()
	} else {
		microversionCheckTimeout = DefaultMicroversionCheckTimeout
	}

	if authUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_url"),
			"Missing samsungcloudplatformv2 Auth URL",
			"The provider cannot create the samsungcloudplatformv2 API client as there is a missing or empty value for the samsungcloudplatformv2 Auth URL. "+
				"Set the auth_url value in the configuration or use the SCP_TF_AUTH_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if accessKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing samsungcloudplatformv2 Access Key",
			"The provider cannot create the samsungcloudplatformv2 API client as there is a missing or empty value for the samsungcloudplatformv2 Access Key. "+
				"Set the access_key value in the configuration or use the SCP_TF_ACCESS_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if secretKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing samsungcloudplatformv2 Secret Key",
			"The provider cannot create the samsungcloudplatformv2 API client as there is a missing or empty value for the samsungcloudplatformv2 Secret Key. "+
				"Set the secret_key value in the configuration or use the SCP_TF_SECRET_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	providerConfig.AuthUrl = types.StringValue(authUrl)
	providerConfig.EndpointOverride = types.StringValue(endpointOverride)
	providerConfig.AccountId = types.StringValue(AccountId)
	providerConfig.DefaultRegion = types.StringValue(DefaultRegion)
	providerConfig.AccessKey = types.StringValue(accessKey)
	providerConfig.SecretKey = types.StringValue(secretKey)
	providerConfig.AuthToken = types.StringValue(authToken)
	providerConfig.MaxRemainDays = types.Int64Value(maxRemainDays)
	providerConfig.MicroversionCheckTimeout = types.Int64Value(microversionCheckTimeout)
}

func getAuthToken(authUrl string, accessKey string, secretKey string) (string, error) {
	jsonData := fmt.Sprintf(`{
		"auth": {
			"identity": {
				"methods": ["application_credential"],
				"application_credential": {
					"id": "%s",
					"secret": "%s"
				}
			}
		}
	}`, accessKey, secretKey)

	req, err := http.NewRequest("POST", authUrl+"/v3/auth/tokens?nocatalog=null", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	tlsConfig, _ := common.CreateTlsConfig()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	xSubjectToken := resp.Header.Get("X-Subject-Token")

	return xSubjectToken, nil
}
