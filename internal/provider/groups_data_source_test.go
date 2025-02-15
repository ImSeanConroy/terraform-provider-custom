// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "custom_groups" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.custom_groups.test", "groups.0.id", "42015d3b-3ab6-4355-8361-b65aa5223a6a"),
					resource.TestCheckResourceAttr("data.custom_groups.test", "groups.0.title", "Learn Robotics"),
					resource.TestCheckResourceAttr("data.custom_groups.test", "groups.0.created_at", "2025-02-03T10:00:00.000Z"),
					resource.TestCheckResourceAttr("data.custom_groups.test", "groups.0.updated_at", "2025-02-03T10:05:00.000Z"),
				),
			},
		},
	})
}
