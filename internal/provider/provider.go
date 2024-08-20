package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/imseanconroy/go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &customProvider{}
)

// customProviderModel maps provider schema data to a Go type.
type customProviderModel struct {
	Url   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &customProvider{
			version: version,
		}
	}
}

// customProvider is the provider implementation.
type customProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *customProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "custom"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *customProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *customProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Custom client")

	// Retrieve provider data from configuration
	var config customProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Custom API Url",
			"The provider cannot create the Custom API client as there is an unknown configuration value for the Custom API url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CUSTOM_URL environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Custom API Token",
			"The provider cannot create the Custom API client as there is an unknown configuration value for the Custom API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CUSTOM_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("CUSTOM_URL")
	token := os.Getenv("CUSTOM_TOKEN")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Custom API Url",
			"The provider cannot create the Custom API client as there is a missing or empty value for the Custom API url. "+
				"Set the url value in the configuration or use the CUSTOM_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Custom API Token",
			"The provider cannot create the Custom API client as there is a missing or empty value for the Custom API token. "+
				"Set the token value in the configuration or use the CUSTOM_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "custom_url", url)
	ctx = tflog.SetField(ctx, "custom_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "custom_token")

	tflog.Debug(ctx, "Creating Custom client")

	// Create a new Custom client using the configuration values
	client, err := client.NewClient(url, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Custom API Client",
			"An unexpected error occurred when creating the Custom API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Custom Client Error: "+err.Error(),
		)
		return
	}

	// Make the Custom client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Custom client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *customProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *customProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
