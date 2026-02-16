package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ssidResource{}
	_ resource.ResourceWithConfigure   = &ssidResource{}
	_ resource.ResourceWithImportState = &ssidResource{}
)

// NewSSIDResource is a helper function to simplify the provider implementation.
func NewSSIDResource() resource.Resource {
	return &ssidResource{}
}

// ssidResource is the resource implementation.
type ssidResource struct {
	client *client.Client
}

// ssidResourceModel maps the resource schema data.
type ssidResourceModel struct {
	ID              types.String `tfsdk:"id"`
	SiteID          types.String `tfsdk:"site_id"`
	Name            types.String `tfsdk:"name"`
	SSID            types.String `tfsdk:"ssid"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	HideSSID        types.Bool   `tfsdk:"hide_ssid"`
	SecurityMode    types.String `tfsdk:"security_mode"`
	Password        types.String `tfsdk:"password"`
	VlanID          types.Int64  `tfsdk:"vlan_id"`
	GuestNetwork    types.Bool   `tfsdk:"guest_network"`
	ClientIsolation types.Bool   `tfsdk:"client_isolation"`
	Band2_4GEnabled types.Bool   `tfsdk:"band_2_4g_enabled"`
	Band5GEnabled   types.Bool   `tfsdk:"band_5g_enabled"`
	Band6GEnabled   types.Bool   `tfsdk:"band_6g_enabled"`
	MaxClients      types.Int64  `tfsdk:"max_clients"`
	RateLimit       types.Bool   `tfsdk:"rate_limit"`
	DownlinkLimit   types.Int64  `tfsdk:"downlink_limit"`
	UplinkLimit     types.Int64  `tfsdk:"uplink_limit"`
	ScheduleEnabled types.Bool   `tfsdk:"schedule_enabled"`
	PortalEnabled   types.Bool   `tfsdk:"portal_enabled"`
	RadiusProfile   types.String `tfsdk:"radius_profile"`
}

// Metadata returns the resource type name.
func (r *ssidResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssid"
}

// Schema defines the schema for the resource.
func (r *ssidResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Omada wireless network (SSID).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SSID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID where the SSID belongs. Defaults to the provider's site_id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name/description of the SSID.",
				Required:    true,
			},
			"ssid": schema.StringAttribute{
				Description: "The SSID (network name) that will be broadcast.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the SSID is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"hide_ssid": schema.BoolAttribute{
				Description: "Whether to hide the SSID from being broadcast.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"security_mode": schema.StringAttribute{
				Description: "The security mode (wpa2-personal, wpa3-personal, wpa2-enterprise, wpa3-enterprise, open).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("wpa2-personal", "wpa3-personal", "wpa2-enterprise", "wpa3-enterprise", "open"),
				},
			},
			"password": schema.StringAttribute{
				Description: "The WiFi password (required for WPA2/WPA3 personal modes).",
				Optional:    true,
				Sensitive:   true,
			},
			"vlan_id": schema.Int64Attribute{
				Description: "The VLAN ID to assign to this SSID.",
				Optional:    true,
			},
			"guest_network": schema.BoolAttribute{
				Description: "Whether this is a guest network.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"client_isolation": schema.BoolAttribute{
				Description: "Whether to isolate clients from each other.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"band_2_4g_enabled": schema.BoolAttribute{
				Description: "Whether to enable this SSID on the 2.4GHz band.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"band_5g_enabled": schema.BoolAttribute{
				Description: "Whether to enable this SSID on the 5GHz band.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"band_6g_enabled": schema.BoolAttribute{
				Description: "Whether to enable this SSID on the 6GHz band.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"max_clients": schema.Int64Attribute{
				Description: "Maximum number of clients allowed to connect.",
				Optional:    true,
			},
			"rate_limit": schema.BoolAttribute{
				Description: "Whether to enable rate limiting.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"downlink_limit": schema.Int64Attribute{
				Description: "Downlink bandwidth limit in Kbps (when rate_limit is enabled).",
				Optional:    true,
			},
			"uplink_limit": schema.Int64Attribute{
				Description: "Uplink bandwidth limit in Kbps (when rate_limit is enabled).",
				Optional:    true,
			},
			"schedule_enabled": schema.BoolAttribute{
				Description: "Whether to enable scheduled SSID activation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"portal_enabled": schema.BoolAttribute{
				Description: "Whether to enable captive portal.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"radius_profile": schema.StringAttribute{
				Description: "RADIUS profile name (for enterprise modes).",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ssidResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *ssidResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ssidResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get site ID from plan or use provider default
	siteID := plan.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
		plan.SiteID = types.StringValue(siteID)
	}

	// Create SSID API request
	ssid := &client.SSID{
		Name:            plan.Name.ValueString(),
		SSID:            plan.SSID.ValueString(),
		Enabled:         plan.Enabled.ValueBool(),
		HideSSID:        plan.HideSSID.ValueBool(),
		SecurityMode:    plan.SecurityMode.ValueString(),
		Password:        plan.Password.ValueString(),
		VlanID:          int(plan.VlanID.ValueInt64()),
		GuestNetwork:    plan.GuestNetwork.ValueBool(),
		ClientIsolation: plan.ClientIsolation.ValueBool(),
		Band2_4GEnabled: plan.Band2_4GEnabled.ValueBool(),
		Band5GEnabled:   plan.Band5GEnabled.ValueBool(),
		Band6GEnabled:   plan.Band6GEnabled.ValueBool(),
		MaxClients:      int(plan.MaxClients.ValueInt64()),
		RateLimit:       plan.RateLimit.ValueBool(),
		DownlinkLimit:   int(plan.DownlinkLimit.ValueInt64()),
		UplinkLimit:     int(plan.UplinkLimit.ValueInt64()),
		ScheduleEnabled: plan.ScheduleEnabled.ValueBool(),
		PortalEnabled:   plan.PortalEnabled.ValueBool(),
		RadiusProfile:   plan.RadiusProfile.ValueString(),
	}

	created, err := r.client.CreateSSID(siteID, ssid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SSID",
			"Could not create SSID, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(created.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ssidResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ssidResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	ssid, err := r.client.GetSSID(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SSID",
			"Could not read SSID ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state (preserve password from state as API doesn't return it)
	state.Name = types.StringValue(ssid.Name)
	state.SSID = types.StringValue(ssid.SSID)
	state.Enabled = types.BoolValue(ssid.Enabled)
	state.HideSSID = types.BoolValue(ssid.HideSSID)
	state.SecurityMode = types.StringValue(ssid.SecurityMode)
	// Password is not returned by API, keep existing value
	if ssid.VlanID > 0 {
		state.VlanID = types.Int64Value(int64(ssid.VlanID))
	}
	state.GuestNetwork = types.BoolValue(ssid.GuestNetwork)
	state.ClientIsolation = types.BoolValue(ssid.ClientIsolation)
	state.Band2_4GEnabled = types.BoolValue(ssid.Band2_4GEnabled)
	state.Band5GEnabled = types.BoolValue(ssid.Band5GEnabled)
	state.Band6GEnabled = types.BoolValue(ssid.Band6GEnabled)
	if ssid.MaxClients > 0 {
		state.MaxClients = types.Int64Value(int64(ssid.MaxClients))
	}
	state.RateLimit = types.BoolValue(ssid.RateLimit)
	if ssid.DownlinkLimit > 0 {
		state.DownlinkLimit = types.Int64Value(int64(ssid.DownlinkLimit))
	}
	if ssid.UplinkLimit > 0 {
		state.UplinkLimit = types.Int64Value(int64(ssid.UplinkLimit))
	}
	state.ScheduleEnabled = types.BoolValue(ssid.ScheduleEnabled)
	state.PortalEnabled = types.BoolValue(ssid.PortalEnabled)
	if ssid.RadiusProfile != "" {
		state.RadiusProfile = types.StringValue(ssid.RadiusProfile)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ssidResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ssidResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	// Update SSID API request
	ssid := &client.SSID{
		ID:              plan.ID.ValueString(),
		Name:            plan.Name.ValueString(),
		SSID:            plan.SSID.ValueString(),
		Enabled:         plan.Enabled.ValueBool(),
		HideSSID:        plan.HideSSID.ValueBool(),
		SecurityMode:    plan.SecurityMode.ValueString(),
		Password:        plan.Password.ValueString(),
		VlanID:          int(plan.VlanID.ValueInt64()),
		GuestNetwork:    plan.GuestNetwork.ValueBool(),
		ClientIsolation: plan.ClientIsolation.ValueBool(),
		Band2_4GEnabled: plan.Band2_4GEnabled.ValueBool(),
		Band5GEnabled:   plan.Band5GEnabled.ValueBool(),
		Band6GEnabled:   plan.Band6GEnabled.ValueBool(),
		MaxClients:      int(plan.MaxClients.ValueInt64()),
		RateLimit:       plan.RateLimit.ValueBool(),
		DownlinkLimit:   int(plan.DownlinkLimit.ValueInt64()),
		UplinkLimit:     int(plan.UplinkLimit.ValueInt64()),
		ScheduleEnabled: plan.ScheduleEnabled.ValueBool(),
		PortalEnabled:   plan.PortalEnabled.ValueBool(),
		RadiusProfile:   plan.RadiusProfile.ValueString(),
	}

	_, err := r.client.UpdateSSID(siteID, plan.ID.ValueString(), ssid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SSID",
			"Could not update SSID, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ssidResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ssidResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	err := r.client.DeleteSSID(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SSID",
			"Could not delete SSID, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform.
func (r *ssidResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
