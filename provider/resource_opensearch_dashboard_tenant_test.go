package provider

import (
	"context"
	"fmt"
	"testing"

	elastic7 "github.com/olivere/elastic/v7"
	elastic6 "gopkg.in/olivere/elastic.v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpensearchOpenDistroDashboardTenant(t *testing.T) {
	provider := Provider()
	diags := provider.Configure(context.Background(), &terraform.ResourceConfig{})
	if diags.HasError() {
		t.Skipf("err: %#v", diags)
	}
	meta := provider.Meta()
	esClient, err := getClient(meta.(*ProviderConf))
	if err != nil {
		t.Skipf("err: %s", err)
	}
	var allowed bool
	switch esClient.(type) {
	case *elastic6.Client:
		allowed = false
	default:
		allowed = true
	}

	randomName := "test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if !allowed {
				t.Skip("Allowed only for ES >= 7")
			}
		},
		Providers:    testAccOpendistroProviders,
		CheckDestroy: testAccCheckOpensearchDashboardTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenDistroDashboardTenantResource(randomName),
				Check: resource.ComposeTestCheckFunc(
					testCheckOpensearchDashboardTenantExists("opensearch_dashboard_tenant.test"),
					resource.TestCheckResourceAttr(
						"opensearch_dashboard_tenant.test",
						"id",
						randomName,
					),
					resource.TestCheckResourceAttr(
						"opensearch_dashboard_tenant.test",
						"description",
						"test",
					),
				),
			},
			{
				Config: testAccOpenDistroDashboardTenantResourceUpdated(randomName),
				Check: resource.ComposeTestCheckFunc(
					testCheckOpensearchDashboardTenantExists("opensearch_dashboard_tenant.test"),
					resource.TestCheckResourceAttr(
						"opensearch_dashboard_tenant.test",
						"description",
						"test2",
					),
				),
			},
		},
	})
}

func testAccCheckOpensearchDashboardTenantDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opensearch_dashboard_tenant" {
			continue
		}

		meta := testAccOpendistroProvider.Meta()

		var err error
		esClient, err := getClient(meta.(*ProviderConf))
		if err != nil {
			return err
		}
		switch esClient.(type) {
		case *elastic7.Client:
			_, err = resourceOpensearchGetOpenDistroDashboardTenant(rs.Primary.ID, meta.(*ProviderConf))
		default:
		}

		if err != nil {
			return nil // should be not found error
		}

		return fmt.Errorf("DashboardTenant %q still exists", rs.Primary.ID)
	}

	return nil
}
func testCheckOpensearchDashboardTenantExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opensearch_dashboard_tenant" {
				continue
			}

			meta := testAccOpendistroProvider.Meta()

			var err error
			esClient, err := getClient(meta.(*ProviderConf))
			if err != nil {
				return err
			}
			switch esClient.(type) {
			case *elastic7.Client:
				_, err = resourceOpensearchGetOpenDistroDashboardTenant(rs.Primary.ID, meta.(*ProviderConf))
			default:
			}

			if err != nil {
				return err
			}

			return nil
		}

		return nil
	}
}

func testAccOpenDistroDashboardTenantResource(resourceName string) string {
	return fmt.Sprintf(`
resource "opensearch_dashboard_tenant" "test" {
  tenant_name = "%s"
  description = "test"
}
	`, resourceName)
}

func testAccOpenDistroDashboardTenantResourceUpdated(resourceName string) string {
	return fmt.Sprintf(`
resource "opensearch_dashboard_tenant" "test" {
  tenant_name = "%s"
  description = "test2"
}
	`, resourceName)
}
