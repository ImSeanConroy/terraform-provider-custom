// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/imseanconroy/go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tasksDataSource{}
	_ datasource.DataSourceWithConfigure = &tasksDataSource{}
)

// NewTasksDataSource is a helper function to simplify the provider implementation.
func NewTasksDataSource() datasource.DataSource {
	return &tasksDataSource{}
}

// tasksDataSource is the data source implementation.
type tasksDataSource struct {
	client *client.Client
}

// tasksDataSourceModel maps the data source schema data.
type tasksDataSourceModel struct {
	Tasks []taskModel `tfsdk:"tasks"`
}

// coffeesModel maps coffees schema data.
type taskModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Complete    types.Bool   `tfsdk:"complete"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

// Configure adds the provider configured client to the data source.
func (d *tasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *tasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tasks"
}

// Schema defines the schema for the data source.
func (d *tasksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of tasks.",
		Attributes: map[string]schema.Attribute{
			"tasks": schema.ListNestedAttribute{
				Description: "List of tasks.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "String uuid identifier of the task.",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "Task title.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Task description.",
							Optional:    true,
						},
						"complete": schema.BoolAttribute{
							Description: "Task completion status.",
							Optional:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Date and time task was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Date and time task was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tasksDataSourceModel

	res, err := d.client.Get("/tasks")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Custom Tasks",
			err.Error(),
		)
		return
	}

	if tasks := res.Get("tasks").Array(); len(tasks) > 0 {
		for _, task := range tasks {
			taskState := taskModel{
				ID:          types.StringValue(task.Get("id").String()),
				Title:       types.StringValue(task.Get("title").String()),
				Description: types.StringValue(task.Get("description").String()),
				Complete:    types.BoolValue(task.Get("complete").Bool()),
				CreatedAt:   types.StringValue(task.Get("createdAt").String()),
				UpdatedAt:   types.StringValue(task.Get("updatedAt").String()),
			}

			state.Tasks = append(state.Tasks, taskState)
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
