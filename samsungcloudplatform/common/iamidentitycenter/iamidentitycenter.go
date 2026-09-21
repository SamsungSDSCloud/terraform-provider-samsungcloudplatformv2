package iamidentitycenter

import (
	"fmt"
	sysuser "os/user"
	"path/filepath"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/config"
)

func SharedClientForRegion(region string) (client.Instance, error) {
	user, err := sysuser.Current()
	if err != nil {
		return client.Instance{}, fmt.Errorf("failed to get current OS user: %w", err)
	}

	providerConfig := config.ProviderConfig{}
	config.LoadServiceConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.ServiceConfigFile), &providerConfig)
	config.LoadCredentialsConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.CredentialConfigFile), &providerConfig)

	scpClient, err := client.NewSCPClient(&providerConfig)
	if err != nil {
		return client.Instance{}, fmt.Errorf("failed to create SCP client: %w", err)
	}

	inst := client.Instance{
		Client: scpClient,
	}

	inst.Client.Iam.Config.Region = region

	return inst, nil
}
