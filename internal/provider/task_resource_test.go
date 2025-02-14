// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTaskResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "custom_task" "test" {
					title = "New task"
					description = "Create a new task"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_task.test", "title", "New task"),
					resource.TestCheckResourceAttr("custom_task.test", "description", "Create a new task"),
					resource.TestCheckResourceAttrSet("custom_task.test", "complete"),
					resource.TestCheckResourceAttrSet("custom_task.test", "id"),
				),
			},
			{
				ResourceName:      "custom_task.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
				resource "custom_task" "test" {
					title = "Updated task"
					description = "Update a new task"
					complete = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_task.test", "title", "Updated task"),
					resource.TestCheckResourceAttr("custom_task.test", "description", "Update a new task"),
					resource.TestCheckResourceAttr("custom_task.test", "complete", "true"),
				),
			},
		},
	})
}
