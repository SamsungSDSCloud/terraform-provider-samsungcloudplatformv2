package scr_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	containerRegistryResourceType = "samsungcloudplatformv2_scr_container_registry"
	containerRegistryResourceName = "samsungcloudplatformv2_scr_container_registry.registry"
)

func TestAccScrContainerRegistryResource(t *testing.T) {
	registryName := fmt.Sprintf("tf-acc-%s", acctest.RandString(11))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// 1. Create registry with defaults
				Config: testAccScrContainerRegistryConfig(registryName, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(containerRegistryResourceName, "name", registryName),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "public_visible_enabled", "false"),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "public_endpoint_enabled", "false"),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "private_acl_enabled", "true"),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "public_acl_enabled", "false"),
					resource.TestCheckResourceAttrSet(containerRegistryResourceName, "id"),
					resource.TestCheckResourceAttrSet(containerRegistryResourceName, "private_domain"),
					resource.TestCheckResourceAttrSet(containerRegistryResourceName, "bucket_id"),
					resource.TestCheckResourceAttrSet(containerRegistryResourceName, "created_at"),
				),
			},
			{
				// 2. Update: enable public endpoint and public ACL
				Config: testAccScrContainerRegistryConfig(registryName, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(containerRegistryResourceName, "name", registryName),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "public_endpoint_enabled", "true"),
					resource.TestCheckResourceAttr(containerRegistryResourceName, "public_acl_enabled", "true"),
				),
			},
			{
				// 3. Import state
				ResourceName:      containerRegistryResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"public_visible_enabled",
					"public_endpoint_enabled",
					"private_acl_enabled",
					"private_acl_resources",
					"public_acl_enabled",
					"public_acl_resources",
				},
			},
		},
	})
}

func testAccScrContainerRegistryBaseConfig(name string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_scr_container_registry" "registry" {
			name = "%s"
		}`, name)
}

func testAccScrContainerRegistryConfig(name string, publicEndpointEnabled bool, publicAclEnabled bool, publicVisibleEnabled bool) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_scr_container_registry" "registry" {
			name                    = "%s"
			public_visible_enabled  = %t
			public_endpoint_enabled = %t
			private_acl_enabled     = true
			public_acl_enabled      = %t
		}`, name, publicVisibleEnabled, publicEndpointEnabled, publicAclEnabled)
}
