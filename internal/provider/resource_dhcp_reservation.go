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
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpReservationResource{}
	_ resource.ResourceWithConfigure   = &dhcpReservationResource{}
	_ resource.ResourceWithImportState = &dhcpReservationResource{}
)

// NewDHCPReservationResource is a helper function to simplify the provider implementation.
func NewDHCPReservationResource() resource.Resource {
	return &dhcpReservationResource{}
}

// dhcpReservationResource is the resource implementation.
type dhcpReservationResource struct {
	client *client.Client
}

// dhcpReservationResourceModel maps the resource schema data.
type dhcpReservationResourceModel struct {
	ID         types.String `tfsdk:"id"`
	SiteID     types.String `tfsdk:"site_id"`
	Name       types.String `tfsdk:"name"`
	MACAddress types.String `tfsdk:"mac_address"`
	IPAddress  types.String `tfsdk:"ip_address"`
	NetworkID  types.String `tfsdk:"network_id"`
	Comment    types.String `tfsdk:"comment"`
}

// Metadata returns the resource type name.
func (r *dhcpReservationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_reservation"
}

// Schema defines the schema for the resource.
func (r *dhcpReservationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DHCP reservation (static IP assignment) in Omada Controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the DHCP reservation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID where the reservation belongs. Defaults to the provider's site_id.",
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
			"mac_address": schema.StringAttribute{
				Description: "The MAC address of the device (format: AA:BB:CC:DD:EE:FF).",
				Required:    true,
			},
			"ip_address": schema.StringAttribute{
				Description: "The IP address to reserve for this device.",
				Required:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "The network ID where this reservation applies.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Optional comment about this reservation.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *dhcpReservationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dhcpReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dhcpReservationResourceModel
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

	// Create DHCP reservation API request
	reservation := &client.DHCPReservation{
		Name:       plan.Name.ValueString(),
		MACAddress: plan.MACAddress.ValueString(),
		IPAddress:  plan.IPAddress.ValueString(),
		NetworkID:  plan.NetworkID.ValueString(),
		Comment:    plan.Comment.ValueString(),
	}

	created, err := r.client.CreateDHCPReservation(siteID, reservation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating DHCP Reservation",
			"Could not create DHCP reservation, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(created.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dhcpReservationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	reservation, err := r.client.GetDHCPReservation(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading DHCP Reservation",
			"Could not read DHCP reservation ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	state.Name = types.StringValue(reservation.Name)
	state.MACAddress = types.StringValue(reservation.MACAddress)
	state.IPAddress = types.StringValue(reservation.IPAddress)
	state.NetworkID = types.StringValue(reservation.NetworkID)
	if reservation.Comment != "" {
		state.Comment = types.StringValue(reservation.Comment)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dhcpReservationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	// Update DHCP reservation API request
	reservation := &client.DHCPReservation{
		ID:         plan.ID.ValueString(),
		Name:       plan.Name.ValueString(),
		MACAddress: plan.MACAddress.ValueString(),
		IPAddress:  plan.IPAddress.ValueString(),
		NetworkID:  plan.NetworkID.ValueString(),
		Comment:    plan.Comment.ValueString(),
	}

	_, err := r.client.UpdateDHCPReservation(siteID, plan.ID.ValueString(), reservation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating DHCP Reservation",
			"Could not update DHCP reservation, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dhcpReservationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	err := r.client.DeleteDHCPReservation(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting DHCP Reservation",
			"Could not delete DHCP reservation, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform.
func (r *dhcpReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
