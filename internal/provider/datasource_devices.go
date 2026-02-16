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
	_ datasource.DataSource              = &devicesDataSource{}
	_ datasource.DataSourceWithConfigure = &devicesDataSource{}
)

// NewDevicesDataSource is a helper function to simplify the provider implementation.
func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

// devicesDataSource is the data source implementation.
type devicesDataSource struct {
	client *client.Client
}

// devicesDataSourceModel maps the data source schema data.
type devicesDataSourceModel struct {
	SiteID  types.String  `tfsdk:"site_id"`
	Devices []deviceModel `tfsdk:"devices"`
}

// deviceModel maps device data
type deviceModel struct {
	MAC             types.String `tfsdk:"mac"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Model           types.String `tfsdk:"model"`
	Status          types.String `tfsdk:"status"`
	LEDEnabled      types.Bool   `tfsdk:"led_enabled"`
	Location        types.String `tfsdk:"location"`
	Site            types.String `tfsdk:"site"`
	IP              types.String `tfsdk:"ip"`
	Uptime          types.Int64  `tfsdk:"uptime"`
	FirmwareVersion types.String `tfsdk:"firmware_version"`
	NeedAdopt       types.Bool   `tfsdk:"need_adopt"`
}

// Metadata returns the data source type name.
func (d *devicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

// Schema defines the schema for the data source.
func (d *devicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about Omada devices (APs, switches, gateways).",
		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				Description: "The site ID to query devices from. Defaults to the provider's site_id.",
				Optional:    true,
			},
			"devices": schema.ListNestedAttribute{
				Description: "List of devices in the site.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mac": schema.StringAttribute{
							Description: "The MAC address of the device.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the device.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of device (ap, switch, gateway).",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "The model of the device.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the device (connected, disconnected, etc).",
							Computed:    true,
						},
						"led_enabled": schema.BoolAttribute{
							Description: "Whether the LED is enabled on the device.",
							Computed:    true,
						},
						"location": schema.StringAttribute{
							Description: "The location of the device.",
							Computed:    true,
						},
						"site": schema.StringAttribute{
							Description: "The site where the device is located.",
							Computed:    true,
						},
						"ip": schema.StringAttribute{
							Description: "The IP address of the device.",
							Computed:    true,
						},
						"uptime": schema.Int64Attribute{
							Description: "The uptime of the device in seconds.",
							Computed:    true,
						},
						"firmware_version": schema.StringAttribute{
							Description: "The firmware version of the device.",
							Computed:    true,
						},
						"need_adopt": schema.BoolAttribute{
							Description: "Whether the device needs to be adopted.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state devicesDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get site ID from configuration or use the provider's default site
	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = d.client.GetSiteID()
	}

	// Fetch devices
	devices, err := d.client.GetDevices(siteID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Omada Devices",
			"An unexpected error occurred when reading Omada devices. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Map devices to state
	state.SiteID = types.StringValue(siteID)
	state.Devices = make([]deviceModel, len(devices))

	for i, device := range devices {
		state.Devices[i] = deviceModel{
			MAC:             types.StringValue(device.MAC),
			Name:            types.StringValue(device.Name),
			Type:            types.StringValue(device.Type),
			Model:           types.StringValue(device.Model),
			Status:          types.StringValue(device.Status),
			LEDEnabled:      types.BoolValue(device.LEDEnabled),
			Location:        types.StringValue(device.Location),
			Site:            types.StringValue(device.Site),
			IP:              types.StringValue(device.IP),
			Uptime:          types.Int64Value(device.Uptime),
			FirmwareVersion: types.StringValue(device.FirmwareVersion),
			NeedAdopt:       types.BoolValue(device.Adoption),
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
