// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/imseanconroy/go-client"
	"github.com/tidwall/gjson"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &taskResource{}
	_ resource.ResourceWithConfigure   = &taskResource{}
	_ resource.ResourceWithImportState = &taskResource{}
)

// NewTaskResource is a helper function to simplify the provider implementation.
func NewTaskResource() resource.Resource {
	return &taskResource{}
}

// taskResource is the resource implementation.
type taskResource struct {
	client *client.Client
}

type Task struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Complete    bool   `json:"complete,omitempty"`
}

// taskResourceModel maps the resource schema data.
type taskResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Complete    types.Bool   `tfsdk:"complete"`
}

// Configure adds the provider configured client to the resource.
func (r *taskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *taskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

// Schema defines the schema for the resource.
func (r *taskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages tasks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "String uuid identifier of the task.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Task title.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Task description.",
				Optional:    true,
				Computed:    true,
			},
			"complete": schema.BoolAttribute{
				Description: "Task completion status.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Create a new resource.
func (r *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan taskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Title.ValueString()))

	// Generate API request body from plan
	body := Task{
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Complete:    plan.Complete.ValueBool(),
	}

	// Create new task
	res, err := r.client.Post("/tasks", body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating task",
			"Could not create task, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(res.Get("task.id").String())
	plan.Title = types.StringValue(res.Get("task.title").String())
	plan.Description = types.StringValue(res.Get("task.description").String())
	if value := res.Get("task.complete"); !value.Exists() || value.Type == gjson.Null {
		plan.Complete = types.BoolValue(false)
	} else {
		plan.Complete = types.BoolValue(value.Bool())
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Create Finished Successfully", plan.Title.ValueString()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state taskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Title.ValueString()))

	// Get refreshed task value
	res, err := r.client.Get("/tasks/" + state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading task",
			"Could not read task, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite task with refreshed state
	state.ID = types.StringValue(res.Get("id").String())
	state.Title = types.StringValue(res.Get("title").String())
	state.Description = types.StringValue(res.Get("description").String())
	state.Complete = types.BoolValue(res.Get("complete").Bool())

	tflog.Debug(ctx, fmt.Sprintf("%s: Read Finished Successfully", state.Title.ValueString()))

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan taskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Title.ValueString()))

	// Generate API request body from plan
	body := Task{
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Complete:    plan.Complete.ValueBool(),
	}

	// Perform the update
	res, err := r.client.Put("/tasks/"+plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating task",
			"Could not update task, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(res.Get("task.id").String())
	plan.Title = types.StringValue(res.Get("task.title").String())
	plan.Description = types.StringValue(res.Get("task.description").String())
	if value := res.Get("task.complete"); !value.Exists() || value.Type == gjson.Null {
		plan.Complete = types.BoolValue(false)
	} else {
		plan.Complete = types.BoolValue(value.Bool())
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update Finished Successfully", plan.Title.ValueString()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state taskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Title.ValueString()))

	// Delete existing task
	_, err := r.client.Delete("/tasks/" + state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting task",
			"Could not delete task, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete Finished Successfully", state.Title.ValueString()))
}

func (r *taskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
