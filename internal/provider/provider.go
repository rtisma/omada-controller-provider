package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &omadaProvider{}
)

// omadaProviderModel maps provider schema data to a Go type.
type omadaProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	SiteID   types.String `tfsdk:"site_id"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

// omadaProvider is the provider implementation.
type omadaProvider struct {
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &omadaProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *omadaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omada"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *omadaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with TP-Link Omada Controller for managing network infrastructure.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The URL of the Omada Controller. Example: https://192.168.1.1:8043",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Omada Controller authentication.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Omada Controller authentication.",
				Required:    true,
				Sensitive:   true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID (name) to manage. Default is 'Default'.",
				Optional:    true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to skip TLS certificate verification. Default is false.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a Omada API client for data sources and resources.
func (p *omadaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config omadaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Omada Controller Host",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada Controller host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Omada Controller Username",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada Controller username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Omada Controller Password",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada Controller password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := config.Host.ValueString()
	username := config.Username.ValueString()
	password := config.Password.ValueString()
	siteID := config.SiteID.ValueString()
	insecure := config.Insecure.ValueBool()

	// Default site ID to "Default"
	if siteID == "" {
		siteID = "Default"
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Omada Controller Host",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada Controller host. "+
				"Set the host value in the configuration.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Omada Controller Username",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada Controller username. "+
				"Set the username value in the configuration.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Omada Controller Password",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada Controller password. "+
				"Set the password value in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Omada client using the configuration values
	apiClient, err := client.NewClient(host, username, password, siteID, insecure)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Omada API Client",
			"An unexpected error occurred when creating the Omada API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Omada Client Error: "+err.Error(),
		)
		return
	}

	// Make the Omada client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// DataSources defines the data sources implemented in the provider.
func (p *omadaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSiteDataSource,
		NewDevicesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *omadaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
		NewSSIDResource,
		NewDHCPReservationResource,
		NewDeviceResource,
	}
}
