// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotessDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "custom_notes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.custom_notes.test", "notes.0.id", "3f5b4d2e-8c12-4b8b-8cbb-3df2bcad50e5"),
					resource.TestCheckResourceAttr("data.custom_notes.test", "notes.0.text", "Learn Robotics and SLAM"),
					resource.TestCheckResourceAttr("data.custom_notes.test", "notes.0.created_at", "2025-02-03T10:00:00.000Z"),
					resource.TestCheckResourceAttr("data.custom_notes.test", "notes.0.updated_at", "2025-02-03T10:05:00.000Z"),
				),
			},
		},
	})
}
