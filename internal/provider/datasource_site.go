package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &siteDataSource{}
	_ datasource.DataSourceWithConfigure = &siteDataSource{}
)

// NewSiteDataSource is a helper function to simplify the provider implementation.
func NewSiteDataSource() datasource.DataSource {
	return &siteDataSource{}
}

// siteDataSource is the data source implementation.
type siteDataSource struct {
	client *client.Client
}

// siteDataSourceModel maps the data source schema data.
type siteDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Location types.String `tfsdk:"location"`
	TimeZone types.String `tfsdk:"timezone"`
	Scenario types.String `tfsdk:"scenario"`
}

// Metadata returns the data source type name.
func (d *siteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema defines the schema for the data source.
func (d *siteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about an Omada site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the site.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the site.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the site.",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "The location of the site.",
				Computed:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone of the site.",
				Computed:    true,
			},
			"scenario": schema.StringAttribute{
				Description: "The scenario of the site.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *siteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *siteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state siteDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get site name from configuration or use the provider's default site
	siteName := state.Name.ValueString()
	if siteName == "" {
		siteName = d.client.GetSiteID()
	}

	// Fetch site information
	site, err := d.client.GetSite(siteName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Omada Site",
			"An unexpected error occurred when reading the Omada site. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(site.ID)
	state.Name = types.StringValue(site.Name)
	state.Type = types.StringValue(site.Type)
	state.Location = types.StringValue(site.Location)
	state.TimeZone = types.StringValue(site.TimeZone)
	state.Scenario = types.StringValue(site.Scenario)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
