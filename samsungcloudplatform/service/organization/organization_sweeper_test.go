package organization_test

import (
	sysuser "os/user"
	"path/filepath"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func SharedClientForRegion(region string) (client.Instance, error) {
	providerConfig := config.ProviderConfig{}
	user, _ := sysuser.Current()

	config.LoadServiceConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.ServiceConfigFile), &providerConfig)
	config.LoadCredentialsConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.CredentialConfigFile), &providerConfig)

	scpClient, _ := client.NewSCPClient(&providerConfig)

	inst := client.Instance{
		Client: scpClient,
	}

	inst.Client.Organization.Config.Region = region

	return inst, nil
}
