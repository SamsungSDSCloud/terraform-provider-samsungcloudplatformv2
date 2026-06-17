package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	SkipVersionCheck         types.Bool   `tfsdk:"skip_version_check"`
}

const (
	ServiceConfigFile               = "config.json"
	CredentialConfigFile            = "credentials.json"
	DefaultMaxRemainDays            = 90
	DefaultMicroversionCheckTimeout = 15
)

func ConfigureServiceAndCredentials(resp *provider.ConfigureResponse, providerConfig *ProviderConfig, serviceConfigPath string, credentialsConfigPath string) {
	serviceConfig := loadJsonFile[serviceConfig](serviceConfigPath)
	credConfig := loadJsonFile[credentialsConfig](credentialsConfigPath)

	authUrl := getStringValue(providerConfig.AuthUrl, "SCP_TF_AUTH_URL", serviceConfig.AuthUrl)
	endpointOverride := getStringValue(providerConfig.EndpointOverride, "SCP_TF_ENDPOINT_OVERRIDE", serviceConfig.EndpointOverride)
	accountId := getStringValue(providerConfig.AccountId, "SCP_TF_ACCOUNT_ID", serviceConfig.AccountId)
	defaultRegion := getStringValue(providerConfig.DefaultRegion, "SCP_TF_DEFAULT_REGION", serviceConfig.DefaultRegion)
	accessKey := getStringValue(providerConfig.AccessKey, "SCP_TF_ACCESS_KEY", credConfig.AccessKey)
	secretKey := getStringValue(providerConfig.SecretKey, "SCP_TF_SECRET_KEY", credConfig.SecretKey)
	authToken := getStringValue(providerConfig.AuthToken, "SCP_TF_AUTH_TOKEN", credConfig.AuthToken)
	maxRemainDays := getIntValue(providerConfig.MaxRemainDays, "SCP_TF_MAX_REMAIN_DAYS", serviceConfig.MaxRemainDays, DefaultMaxRemainDays)
	microversionCheckTimeout := getIntValue(providerConfig.MicroversionCheckTimeout, "SCP_TF_MICROVERSION_CHECK_TIMEOUT", serviceConfig.MicroversionCheckTimeout, DefaultMicroversionCheckTimeout)
	skipVersionCheck := getBoolValue(providerConfig.SkipVersionCheck, "SCP_TF_SKIP_VERSION_CHECK", true)

	if authUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_url"),
			"Missing samsungcloudplatformv2 Auth URL",
			"Set the auth_url value in the configuration, use SCP_TF_AUTH_URL environment variable, or set in config.json",
		)
	}
	if accessKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing samsungcloudplatformv2 Access Key",
			"Set the access_key value in the configuration, use SCP_TF_ACCESS_KEY environment variable, or set in credentials.json",
		)
	}
	if secretKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing samsungcloudplatformv2 Secret Key",
			"Set the secret_key value in the configuration, use SCP_TF_SECRET_KEY environment variable, or set in credentials.json",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	providerConfig.AuthUrl = types.StringValue(authUrl)
	providerConfig.EndpointOverride = types.StringValue(endpointOverride)
	providerConfig.AccountId = types.StringValue(accountId)
	providerConfig.DefaultRegion = types.StringValue(defaultRegion)
	providerConfig.AccessKey = types.StringValue(accessKey)
	providerConfig.SecretKey = types.StringValue(secretKey)
	providerConfig.AuthToken = types.StringValue(authToken)
	providerConfig.MaxRemainDays = types.Int64Value(maxRemainDays)
	providerConfig.MicroversionCheckTimeout = types.Int64Value(microversionCheckTimeout)
	providerConfig.SkipVersionCheck = types.BoolValue(skipVersionCheck)
}

type serviceConfig struct {
	AuthUrl                  string `json:"auth-url"`
	EndpointOverride         string `json:"endpoint-override"`
	AccountId                string `json:"account-id"`
	DefaultRegion            string `json:"default-region"`
	MaxRemainDays            int64  `json:"max-remain-days"`
	MicroversionCheckTimeout int64  `json:"microversion-check-timeout"`
}

type credentialsConfig struct {
	AccessKey string `json:"access-key"`
	SecretKey string `json:"secret-key"`
	AuthToken string `json:"auth-token"`
}

type jsonConfig[T any] struct {
	data T
}

func loadJsonFile[T any](path string) T {
	var zero T
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil || len(data) == 0 {
		return zero
	}
	var config T
	if err := json.Unmarshal(data, &config); err != nil {
		return zero
	}
	return config
}

func getStringValue(tfValue types.String, envKey string, fileValue string) string {
	if !tfValue.IsNull() && tfValue.ValueString() != "" {
		return tfValue.ValueString()
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return fileValue
}

func getIntValue(tfValue types.Int64, envKey string, fileValue int64, defaultValue int64) int64 {
	if !tfValue.IsNull() {
		return tfValue.ValueInt64()
	}
	if v := os.Getenv(envKey); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 32); err == nil {
			return parsed
		}
	}
	if fileValue != 0 {
		return fileValue
	}
	return defaultValue
}

func getBoolValue(tfValue types.Bool, envKey string, defaultValue bool) bool {
	if !tfValue.IsNull() {
		return tfValue.ValueBool()
	}
	if v := os.Getenv(envKey); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func LoadServiceConfig(resp *provider.ConfigureResponse, path string, providerConfig *ProviderConfig) {
	config := loadJsonFile[serviceConfig](path)
	if config.AuthUrl != "" {
		providerConfig.AuthUrl = types.StringValue(config.AuthUrl)
	}
	if config.EndpointOverride != "" {
		providerConfig.EndpointOverride = types.StringValue(config.EndpointOverride)
	}
	if config.AccountId != "" {
		providerConfig.AccountId = types.StringValue(config.AccountId)
	}
	if config.DefaultRegion != "" {
		providerConfig.DefaultRegion = types.StringValue(config.DefaultRegion)
	}
	if config.MaxRemainDays != 0 {
		providerConfig.MaxRemainDays = types.Int64Value(config.MaxRemainDays)
	}
	if config.MicroversionCheckTimeout != 0 {
		providerConfig.MicroversionCheckTimeout = types.Int64Value(config.MicroversionCheckTimeout)
	}
}

func LoadCredentialsConfig(resp *provider.ConfigureResponse, path string, providerConfig *ProviderConfig) {
	config := loadJsonFile[credentialsConfig](path)
	if config.AccessKey != "" {
		providerConfig.AccessKey = types.StringValue(config.AccessKey)
	}
	if config.SecretKey != "" {
		providerConfig.SecretKey = types.StringValue(config.SecretKey)
	}
	if config.AuthToken != "" {
		providerConfig.AuthToken = types.StringValue(config.AuthToken)
	}
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
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	tlsConfig, err := common.CreateTlsConfig()
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	xSubjectToken := resp.Header.Get("X-Subject-Token")

	return xSubjectToken, nil
}
