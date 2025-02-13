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
	_ datasource.DataSource              = &notesDataSource{}
	_ datasource.DataSourceWithConfigure = &notesDataSource{}
)

// NewNotesDataSource is a helper function to simplify the provider implementation.
func NewNotesDataSource() datasource.DataSource {
	return &notesDataSource{}
}

// notesDataSource is the data source implementation.
type notesDataSource struct {
	client *client.Client
}

// notesDataSourceModel maps the data source schema data.
type notesDataSourceModel struct {
	Notes []notesModel `tfsdk:"notes"`
}

// coffeesModel maps coffees schema data.
type notesModel struct {
	ID        types.String `tfsdk:"id"`
	Text      types.String `tfsdk:"text"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Configure adds the provider configured client to the data source.
func (d *notesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *notesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notes"
}

// Schema defines the schema for the data source.
func (d *notesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of notes.",
		Attributes: map[string]schema.Attribute{
			"notes": schema.ListNestedAttribute{
				Description: "List of notes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "String uuid identifier of the note.",
							Computed:    true,
						},
						"text": schema.StringAttribute{
							Description: "Note content.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Date and time note was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Date and time note was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *notesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notesDataSourceModel

	res, err := d.client.Get("/notes")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Custom Notes",
			err.Error(),
		)
		return
	}

	if notes := res.Get("notes").Array(); len(notes) > 0 {
		for _, note := range notes {
			noteState := notesModel{
				ID:        types.StringValue(note.Get("id").String()),
				Text:      types.StringValue(note.Get("text").String()),
				CreatedAt: types.StringValue(note.Get("createdAt").String()),
				UpdatedAt: types.StringValue(note.Get("updatedAt").String()),
			}

			state.Notes = append(state.Notes, noteState)
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
