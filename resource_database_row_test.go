package main

import (
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/datasource"
	resource2 "github.com/BESTSELLER/terraform-provider-servicenow-data/internal/resource"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRowBasic(t *testing.T) {
	tableId := "x_beas_team_engi_0_lasse"
	team := "engineering-services"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRowBasic(tableId, team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRowExists("servicenow-data_table_row.eng-services-test"),
				),
			},
		},
	})
}

func testAccCheckRowDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resource2.TableRowResourceName {
			continue
		}

		tableID, rowID, err := datasource.ExtractIDs(rs.Primary.ID)

		if err != nil {
			return err
		}

		err = c.DeleteTableRow(tableID, rowID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckRowBasic(tableId, team string) string {
	return fmt.Sprintf(`
resource "servicenow-data_table_row" "eng-services-test" {
  table_id = "%s"
  row_data = {
    "team": "%s",
  }
}
	`, tableId, team)
}

func testAccCheckRowExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		return nil
	}
}
