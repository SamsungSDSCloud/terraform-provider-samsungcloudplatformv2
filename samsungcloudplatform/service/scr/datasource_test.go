package scr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	containerRegistryDataSourceName  = "data.samsungcloudplatformv2_scr_container_registry.test"
	registriesDataSourceName = "data.samsungcloudplatformv2_scr_container_registries.test"
)

func TestAccScrContainerRegistryDataSource(t *testing.T) {
	registryName := fmt.Sprintf("tf-acc-%s", acctest.RandString(11))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccScrContainerRegistryDataSourceConfig(registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(containerRegistryDataSourceName, "name", registryName),
					resource.TestCheckResourceAttrSet(containerRegistryDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(containerRegistryDataSourceName, "private_domain"),
				),
			},
		},
	})
}

func TestAccScrRegistriesDataSource(t *testing.T) {
	registryName := fmt.Sprintf("tf-acc-%s", acctest.RandString(11))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccScrRegistriesDataSourceConfig(registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(registriesDataSourceName, "ids.#"),
				),
			},
		},
	})
}

func testAccScrContainerRegistryDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_scr_container_registry" "registry" {
			name = "%s"
		}

		data "samsungcloudplatformv2_scr_container_registry" "test" {
			id = samsungcloudplatformv2_scr_container_registry.registry.id
		}`, name)
}

func testAccScrRegistriesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_scr_container_registry" "registry" {
			name = "%s"
		}

		data "samsungcloudplatformv2_scr_container_registries" "test" {
			name = "%s"
		}`, name, name)
}

// ---- Image Data Source Tests ----

const (
	imageDataSourceName  = "data.samsungcloudplatformv2_scr_image.test"
	imagesDataSourceName = "data.samsungcloudplatformv2_scr_images.test"
)

func TestAccScrImageDataSource(t *testing.T) {
	imageId := os.Getenv("SCP_SCR_IMAGE_ID")
	if imageId == "" {
		t.Skip("Set SCP_SCR_IMAGE_ID to run this test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccScrImageDataSourceConfig(imageId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(imageDataSourceName, "name"),
					resource.TestCheckResourceAttrSet(imageDataSourceName, "registry_id"),
					resource.TestCheckResourceAttrSet(imageDataSourceName, "repository_id"),
					resource.TestCheckResourceAttrSet(imageDataSourceName, "created_at"),
				),
			},
		},
	})
}

func TestAccScrImagesDataSource(t *testing.T) {
	repositoryId := os.Getenv("SCP_SCR_REPOSITORY_ID")
	if repositoryId == "" {
		t.Skip("Set SCP_SCR_REPOSITORY_ID to run this test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccScrImagesDataSourceConfig(repositoryId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(imagesDataSourceName, "ids.#"),
				),
			},
		},
	})
}

func testAccScrImageDataSourceConfig(imageId string) string {
	return fmt.Sprintf(`
		data "samsungcloudplatformv2_scr_image" "test" {
			id = "%s"
		}`, imageId)
}

func testAccScrImagesDataSourceConfig(repositoryId string) string {
	return fmt.Sprintf(`
		data "samsungcloudplatformv2_scr_images" "test" {
			repository_id = "%s"
		}`, repositoryId)
}
