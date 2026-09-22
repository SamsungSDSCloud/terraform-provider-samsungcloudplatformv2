package cloudcontrol_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const lzIdFilePath = "/tmp/test_lz_ids.txt"

func TestAccLandingZoneResourceTest(t *testing.T) {
	uniqueSuffix := generateRandomString(7)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLandingZoneCreate(uniqueSuffix),
				Check:  saveLandingZoneIdCheck("samsungcloudplatformv2_cloudcontrol_landing_zone.landing_zone"),
			},
			{
				Config: testAccLandingZoneUpdate(uniqueSuffix),
			},
		},
		CheckDestroy: testAccCheckLandingZoneInactive,
	})
}

func testAccLandingZoneCreate(uniqueSuffix string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_cloudcontrol_landing_zone" "landing_zone" {
			additional_ou_name         = "SandboxOu-%s"
			agree_yn                   = "Y"
			audit_account_name         = "Audit-%s"
			audit_login_id             = "audit-%s@samsung.com"
			basic_ou_name              = "SecurityOu-%s"
			detective_guardrail_status = "ENABLED"
			log_archive_account_name   = "Log-%s"
			log_archive_login_id       = "log-%s@samsung.com"
			sso_type                   = "SELF"
		}
	`,
		uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix,
	)
}

func testAccLandingZoneUpdate(uniqueSuffix string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_cloudcontrol_landing_zone" "landing_zone" {
			additional_ou_name         = "SandboxOu-%s"
			agree_yn                   = "Y"
			audit_account_name         = "Audit-%s"
			audit_login_id             = "audit-%s@samsung.com"
			basic_ou_name              = "SecurityOu-%s"
			detective_guardrail_status = "DISABLED"
			log_archive_account_name   = "Log-%s"
			log_archive_login_id       = "log-%s@samsung.com"
			sso_type                   = "SELF"
		}
	`,
		uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix, uniqueSuffix,
	)
}

func saveLandingZoneIdCheck(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		landingZoneId := rs.Primary.Attributes["landing_zone_id"]
		if landingZoneId == "" {
			return fmt.Errorf("landing_zone_id is empty for %s", resourceName)
		}
		return saveLandingZoneId(landingZoneId)
	}
}

func testAccCheckLandingZoneInactive(s *terraform.State) error {
	ids, err := loadLandingZoneIds()
	if err != nil || len(ids) == 0 {
		return nil
	}
	defer clearLandingZoneIds()

	ctx := context.Background()
	scpClient, err := SharedClientForRegion(os.Getenv("SCP_REGION"))
	if err != nil {
		return fmt.Errorf("failed to get client: %v", err)
	}

	maxRetries := 60
	retryInterval := 10 * time.Second

	for _, id := range ids {
		ok := false
		for i := 0; i < maxRetries; i++ {
			data, err := scpClient.Client.CloudControl.GetLandingZone(ctx, id)
			if err != nil {
				ok = true
				break
			}
			status := string(data.LandingZone.Status)
			if status == "INACTIVE" {
				ok = true
				break
			}
			if status == "DELETE_FAILED" {
				return fmt.Errorf("landing zone %s failed to delete: %s", id, status)
			}
			time.Sleep(retryInterval)
		}
		if !ok {
			return fmt.Errorf("timeout waiting for landing zone %s to become INACTIVE", id)
		}
	}
	return nil
}

func saveLandingZoneId(id string) error {
	f, err := os.OpenFile(lzIdFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", id)
	return err
}

func loadLandingZoneIds() ([]string, error) {
	data, err := os.ReadFile(lzIdFilePath)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			ids = append(ids, line)
		}
	}
	return ids, nil
}

func clearLandingZoneIds() error {
	return os.Remove(lzIdFilePath)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_cloudcontrol_landing_zone", &resource.Sweeper{
		Name: "samsungcloudplatformv2_cloudcontrol_landing_zone",
		F:    sweepLandingZone,
	})
}

func sweepLandingZone(region string) error {
	ids, err := loadLandingZoneIds()
	if err != nil {
		return nil
	}

	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %v", err)
	}

	ctx := context.Background()
	var failedIds []string

	for _, id := range ids {
		data, err := scpClient.Client.CloudControl.GetLandingZone(ctx, id)
		if err != nil {
			continue
		}
		status := data.LandingZone.Status
		if status == "INACTIVE" {
			continue
		}
		if status != "ACTIVE" && status != "CREATE_FAILED" && status != "DELETE_FAILED" {
			failedIds = append(failedIds, id)
			continue
		}
		if _, err := scpClient.Client.CloudControl.DeleteLandingZone(ctx, id); err != nil {
			failedIds = append(failedIds, id)
		}
	}

	clearLandingZoneIds()
	for _, id := range failedIds {
		saveLandingZoneId(id)
	}
	return nil
}

func generateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*1103515245 + 12345
		idx := seed % int64(len(charset))
		if idx < 0 {
			idx = -idx
		}
		b[i] = charset[idx]
	}
	return string(b)
}
