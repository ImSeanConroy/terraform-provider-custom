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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &noteResource{}
	_ resource.ResourceWithConfigure   = &noteResource{}
	_ resource.ResourceWithImportState = &noteResource{}
)

// NewNoteResource is a helper function to simplify the provider implementation.
func NewNoteResource() resource.Resource {
	return &noteResource{}
}

// noteResource is the resource implementation.
type noteResource struct {
	client *client.Client
}

type Note struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text"`
}

// noteResourceModel maps the resource schema data.
type noteResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Text types.String `tfsdk:"text"`
}

// Configure adds the provider configured client to the resource.
func (r *noteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *noteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_note"
}

// Schema defines the schema for the resource.
func (r *noteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"text": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (r *noteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan noteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Text.ValueString()))

	// Generate API request body from plan
	body := Note{
		Text: plan.Text.ValueString(),
	}

	// Create new order
	res, err := r.client.Post("/notes", body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(res.Get("note.id").String())
	plan.Text = types.StringValue(res.Get("note.text").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Create Finished Successfully", plan.Text.ValueString()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *noteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state noteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Text.ValueString()))

	// Get refreshed note value
	res, err := r.client.Get("/notes/" + state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Custom Notes",
			err.Error(),
		)
		return
	}

	// Overwrite notes with refreshed state
	state.ID = types.StringValue(res.Get("id").String())
	state.Text = types.StringValue(res.Get("text").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Read Finished Successfully", state.Text.ValueString()))

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *noteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan noteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Text.ValueString()))

	body := Note{
		Text: plan.Text.ValueString(),
	}

	// Perform the update
	res, err := r.client.Put("/notes/"+plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating order",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(res.Get("note.id").String())
	plan.Text = types.StringValue(res.Get("note.text").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Update Finished Successfully", plan.Text.ValueString()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *noteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state noteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Text.ValueString()))

	// Delete existing note
	_, err := r.client.Delete("/notes/" + state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Custom Notes",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete Finished Successfully", state.Text.ValueString()))
}

func (r *noteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
