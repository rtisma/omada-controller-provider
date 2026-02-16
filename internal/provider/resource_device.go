package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

// NewDeviceResource is a helper function to simplify the provider implementation.
func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client *client.Client
}

// deviceResourceModel maps the resource schema data.
type deviceResourceModel struct {
	MAC        types.String `tfsdk:"mac"`
	SiteID     types.String `tfsdk:"site_id"`
	Name       types.String `tfsdk:"name"`
	LEDEnabled types.Bool   `tfsdk:"led_enabled"`
	Location   types.String `tfsdk:"location"`
	Type       types.String `tfsdk:"type"`
	Model      types.String `tfsdk:"model"`
	Status     types.String `tfsdk:"status"`
}

// Metadata returns the resource type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Omada device configuration (AP, switch, gateway).",
		Attributes: map[string]schema.Attribute{
			"mac": schema.StringAttribute{
				Description: "The MAC address of the device (format: AA:BB:CC:DD:EE:FF). This is used as the device identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID where the device belongs. Defaults to the provider's site_id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name/description of the device.",
				Required:    true,
			},
			"led_enabled": schema.BoolAttribute{
				Description: "Whether the LED is enabled on the device.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"location": schema.StringAttribute{
				Description: "The location of the device.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of device (ap, switch, gateway). This is read-only.",
				Computed:    true,
			},
			"model": schema.StringAttribute{
				Description: "The model of the device. This is read-only.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the device. This is read-only.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deviceResourceModel
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

	// First, check if device exists
	device, err := r.client.GetDevice(siteID, plan.MAC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device",
			"Could not find device with MAC "+plan.MAC.ValueString()+": "+err.Error()+
				"\n\nMake sure the device is connected to the controller before managing it with Terraform.",
		)
		return
	}

	// If device needs adoption, adopt it
	if device.Adoption {
		if err := r.client.AdoptDevice(siteID, plan.MAC.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error Adopting Device",
				"Could not adopt device, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Update device configuration
	config := &client.DeviceConfig{
		Name:       plan.Name.ValueString(),
		LEDEnabled: plan.LEDEnabled.ValueBool(),
		Location:   plan.Location.ValueString(),
	}

	updated, err := r.client.UpdateDevice(siteID, plan.MAC.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Configuring Device",
			"Could not configure device, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Type = types.StringValue(updated.Type)
	plan.Model = types.StringValue(updated.Model)
	plan.Status = types.StringValue(updated.Status)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	device, err := r.client.GetDevice(siteID, state.MAC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device",
			"Could not read device MAC "+state.MAC.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	state.Name = types.StringValue(device.Name)
	state.LEDEnabled = types.BoolValue(device.LEDEnabled)
	state.Location = types.StringValue(device.Location)
	state.Type = types.StringValue(device.Type)
	state.Model = types.StringValue(device.Model)
	state.Status = types.StringValue(device.Status)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	// Update device configuration
	config := &client.DeviceConfig{
		Name:       plan.Name.ValueString(),
		LEDEnabled: plan.LEDEnabled.ValueBool(),
		Location:   plan.Location.ValueString(),
	}

	updated, err := r.client.UpdateDevice(siteID, plan.MAC.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Device",
			"Could not update device, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Type = types.StringValue(updated.Type)
	plan.Model = types.StringValue(updated.Model)
	plan.Status = types.StringValue(updated.Status)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	// Forget the device (remove from controller)
	err := r.client.ForgetDevice(siteID, state.MAC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Device",
			"Could not forget device, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform.
func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by MAC address
	resource.ImportStatePassthroughID(ctx, path.Root("mac"), req, resp)
}
