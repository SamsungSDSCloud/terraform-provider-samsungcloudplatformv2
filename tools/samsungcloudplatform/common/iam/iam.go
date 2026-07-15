package iam

import (
	sysuser "os/user"
	"path/filepath"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/config"
)

func SharedClientForRegion(region string) (client.Instance, error) {
	providerConfig := config.ProviderConfig{}
	user, _ := sysuser.Current()

	config.LoadServiceConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.ServiceConfigFile), &providerConfig)
	config.LoadCredentialsConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.CredentialConfigFile), &providerConfig)

	scpClient, _ := client.NewSCPClient(&providerConfig)

	inst := client.Instance{
		Client: scpClient,
	}

	inst.Client.Iam.Config.Region = region

	return inst, nil
}
