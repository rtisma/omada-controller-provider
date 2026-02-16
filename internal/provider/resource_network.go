package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-omada/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

// NewNetworkResource is a helper function to simplify the provider implementation.
func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

// networkResource is the resource implementation.
type networkResource struct {
	client *client.Client
}

// networkResourceModel maps the resource schema data.
type networkResourceModel struct {
	ID           types.String `tfsdk:"id"`
	SiteID       types.String `tfsdk:"site_id"`
	Name         types.String `tfsdk:"name"`
	VlanID       types.Int64  `tfsdk:"vlan_id"`
	Gateway      types.String `tfsdk:"gateway"`
	Netmask      types.String `tfsdk:"netmask"`
	DHCPEnabled  types.Bool   `tfsdk:"dhcp_enabled"`
	DHCPStart    types.String `tfsdk:"dhcp_start"`
	DHCPEnd      types.String `tfsdk:"dhcp_end"`
	LeaseTime    types.Int64  `tfsdk:"lease_time"`
	DNSPrimary   types.String `tfsdk:"dns_primary"`
	DNSSecondary types.String `tfsdk:"dns_secondary"`
	DomainName   types.String `tfsdk:"domain_name"`
	Purpose      types.String `tfsdk:"purpose"`
}

// Metadata returns the resource type name.
func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema defines the schema for the resource.
func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Omada network/VLAN.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the network.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID where the network belongs. Defaults to the provider's site_id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the network.",
				Required:    true,
			},
			"vlan_id": schema.Int64Attribute{
				Description: "The VLAN ID for the network (1-4094).",
				Required:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "The gateway IP address for the network.",
				Required:    true,
			},
			"netmask": schema.StringAttribute{
				Description: "The subnet mask for the network (e.g., 255.255.255.0).",
				Required:    true,
			},
			"dhcp_enabled": schema.BoolAttribute{
				Description: "Whether DHCP is enabled for this network.",
				Optional:    true,
			},
			"dhcp_start": schema.StringAttribute{
				Description: "The starting IP address of the DHCP range.",
				Optional:    true,
			},
			"dhcp_end": schema.StringAttribute{
				Description: "The ending IP address of the DHCP range.",
				Optional:    true,
			},
			"lease_time": schema.Int64Attribute{
				Description: "The DHCP lease time in minutes. Default is 1440 (24 hours).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1440),
			},
			"dns_primary": schema.StringAttribute{
				Description: "The primary DNS server.",
				Optional:    true,
			},
			"dns_secondary": schema.StringAttribute{
				Description: "The secondary DNS server.",
				Optional:    true,
			},
			"domain_name": schema.StringAttribute{
				Description: "The domain name for the network.",
				Optional:    true,
			},
			"purpose": schema.StringAttribute{
				Description: "The purpose of the network (e.g., 'general', 'guest', 'iot').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *networkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkResourceModel
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

	// Create network API request
	network := &client.Network{
		Name:         plan.Name.ValueString(),
		VlanID:       int(plan.VlanID.ValueInt64()),
		Gateway:      plan.Gateway.ValueString(),
		Netmask:      plan.Netmask.ValueString(),
		DHCPEnabled:  plan.DHCPEnabled.ValueBool(),
		DHCPStart:    plan.DHCPStart.ValueString(),
		DHCPEnd:      plan.DHCPEnd.ValueString(),
		LeaseTime:    int(plan.LeaseTime.ValueInt64()),
		DNSPrimary:   plan.DNSPrimary.ValueString(),
		DNSSecondary: plan.DNSSecondary.ValueString(),
		DomainName:   plan.DomainName.ValueString(),
		Purpose:      plan.Purpose.ValueString(),
	}

	created, err := r.client.CreateNetwork(siteID, network)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Network",
			"Could not create network, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(created.ID)
	if created.Purpose != "" {
		plan.Purpose = types.StringValue(created.Purpose)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	network, err := r.client.GetNetwork(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Network",
			"Could not read network ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	state.Name = types.StringValue(network.Name)
	state.VlanID = types.Int64Value(int64(network.VlanID))
	state.Gateway = types.StringValue(network.Gateway)
	state.Netmask = types.StringValue(network.Netmask)
	state.DHCPEnabled = types.BoolValue(network.DHCPEnabled)
	state.DHCPStart = types.StringValue(network.DHCPStart)
	state.DHCPEnd = types.StringValue(network.DHCPEnd)
	state.LeaseTime = types.Int64Value(int64(network.LeaseTime))
	state.DNSPrimary = types.StringValue(network.DNSPrimary)
	state.DNSSecondary = types.StringValue(network.DNSSecondary)
	state.DomainName = types.StringValue(network.DomainName)
	if network.Purpose != "" {
		state.Purpose = types.StringValue(network.Purpose)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	// Update network API request
	network := &client.Network{
		ID:           plan.ID.ValueString(),
		Name:         plan.Name.ValueString(),
		VlanID:       int(plan.VlanID.ValueInt64()),
		Gateway:      plan.Gateway.ValueString(),
		Netmask:      plan.Netmask.ValueString(),
		DHCPEnabled:  plan.DHCPEnabled.ValueBool(),
		DHCPStart:    plan.DHCPStart.ValueString(),
		DHCPEnd:      plan.DHCPEnd.ValueString(),
		LeaseTime:    int(plan.LeaseTime.ValueInt64()),
		DNSPrimary:   plan.DNSPrimary.ValueString(),
		DNSSecondary: plan.DNSSecondary.ValueString(),
		DomainName:   plan.DomainName.ValueString(),
		Purpose:      plan.Purpose.ValueString(),
	}

	_, err := r.client.UpdateNetwork(siteID, plan.ID.ValueString(), network)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Network",
			"Could not update network, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.SiteID.ValueString()
	if siteID == "" {
		siteID = r.client.GetSiteID()
	}

	err := r.client.DeleteNetwork(siteID, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Network",
			"Could not delete network, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform.
func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID: format is "site_id/network_id" or just "network_id" to use provider's site_id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
