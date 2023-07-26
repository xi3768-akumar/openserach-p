package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	elastic7 "github.com/olivere/elastic/v7"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpensearchComponentTemplate(t *testing.T) {
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
	case *elastic7.Client:
		allowed = true
	default:
		allowed = false
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if !allowed {
				t.Skip("/_component_template endpoint only supported on ES >= 7.8")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpensearchComponentTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpensearchComponentTemplate,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpensearchComponentTemplateExists("opensearch_component_template.test"),
				),
			},
		},
	})
}

func TestAccOpensearchComponentTemplate_importBasic(t *testing.T) {
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
	case *elastic7.Client:
		allowed = true
	default:
		allowed = false
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if !allowed {
				t.Skip("/_component_template endpoint only supported on ES >= 7.8")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpensearchComponentTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpensearchComponentTemplate,
			},
			{
				ResourceName:      "opensearch_component_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckOpensearchComponentTemplateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No component template ID is set")
		}

		meta := testAccProvider.Meta()

		esClient, err := getClient(meta.(*ProviderConf))
		if err != nil {
			return err
		}

		switch client := esClient.(type) {
		case *elastic7.Client:
			_, err = client.IndexGetComponentTemplate(rs.Primary.ID).Do(context.TODO())
		default:
			err = errors.New("/_component_template endpoint only supported on ES >= 7.8")
		}

		if err != nil {
			return err
		}

		return nil
	}
}

func testCheckOpensearchComponentTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opensearch_component_template" {
			continue
		}

		meta := testAccProvider.Meta()

		esClient, err := getClient(meta.(*ProviderConf))
		if err != nil {
			return err
		}

		switch client := esClient.(type) {
		case *elastic7.Client:
			_, err = client.IndexGetComponentTemplate(rs.Primary.ID).Do(context.TODO())
		default:
			err = errors.New("/_component_template endpoint only supported on ES >= 7.8")
		}

		if err != nil {
			return nil // should be not found error
		}

		return fmt.Errorf("Component template %q still exists", rs.Primary.ID)
	}

	return nil
}

var testAccOpensearchComponentTemplate = `
resource "opensearch_component_template" "test" {
  name = "terraform-test"
  body = <<EOF
{
  "template": {
    "settings": {
      "index": {
        "number_of_shards": 1
      }
    },
    "mappings": {
      "properties": {
        "host_name": {
          "type": "keyword"
        },
        "created_at": {
          "type": "date",
          "format": "EEE MMM dd HH:mm:ss Z yyyy"
        }
      }
    },
    "aliases": {
      "mydata": { }
    }
  }
}
EOF
}
`
