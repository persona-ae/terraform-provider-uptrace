package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
	"github.com/persona-ae/terraform-provider-uptrace/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &monitorResource{}
	_ resource.ResourceWithConfigure   = &monitorResource{}
	_ resource.ResourceWithImportState = &monitorResource{}
)

func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

// monitorResource is the resource implementation.
type monitorResource struct {
	// this client is set by the provider
	client *uptrace.UptraceClient
}

// Metadata returns the resource type name.
func (r *monitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "monitorResource.Metadata", map[string]any{"req": req, "resp": resp})

	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Schema defines the schema for the resource.
func (r *monitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "monitorResource.Schema", map[string]any{"req": req, "resp": resp})

	resp.Schema = schema.Schema{
		Description: "Manages an monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitor.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The monitor type.",
			},
			"query": schema.StringAttribute{
				Required:    true,
				Description: "The monitor's query.",
			},
			"metrics": schema.ListAttribute{
				Required:    true,
				Description: "List of metrics to monitor.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":  types.StringType,
						"alias": types.StringType,
					},
				},
			},
			"column": schema.StringAttribute{
				Required:    true,
				Description: "TODO",
			},
			"min_allowed_value": schema.Float32Attribute{
				Required:    true,
				Description: "TODO",
			},
			"max_allowed_value": schema.Float32Attribute{
				Required:    true,
				Description: "TODO",
			},
			"notify_everyone_by_email": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to notify everyone by email.",
			},
			"team_ids": schema.SetAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
				Description: "List of team ids to be notified by email. Overrides notifyEveryoneByEmail.",
			},
			"channel_ids": schema.SetAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
				Description: "List of channel ids to send notifications.",
			},
			"project_id": schema.Int32Attribute{
				Computed:    true,
				Description: "Project ID the monitor belongs to.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the monitor.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *monitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "monitorResource.Configure", map[string]any{"req": req, "resp": resp})
	if req.ProviderData == nil {
		return
	}

	// extract the client from the provider data
	client, ok := req.ProviderData.(*uptrace.UptraceClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *uptrace.UptraceClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create a new resource.
func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "monitorResource.Create", map[string]any{"req": req, "resp": resp})
	var plan models.TFMonitorData

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "planned query", map[string]any{"query": plan.Query.ValueString()})

	// Generate API request body from plan
	var monitor uptrace.Monitor
	diags = utils.TFMonitorToUptraceMonitor(ctx, plan, &monitor)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating monitor", map[string]any{"monitor": monitor, "query": monitor.Params.Query})

	// Create new monitor
	response := uptrace.MonitorIdResponse{}
	err := r.client.CreateMonitor(ctx, monitor, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create monitor",
			fmt.Sprintf("Failed to create monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "CreateMonitor OK: %s", map[string]any{"response": response})

	plan.ID = types.StringValue(response.Monitor.Id)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read resource information.
func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "monitorResource.Read", map[string]any{"req": req, "resp": resp})

	// Get current state
	// Read data from Terraform state
	var state models.TFMonitorData
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get fresh state from uptrace
	// Generate API request body from plan
	var response uptrace.GetMonitorByIdResponse
	err := r.client.GetMonitorById(ctx, state.ID.ValueString(), &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get monitor",
			fmt.Sprintf("Failed to get monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "GetMonitorById OK", map[string]any{"response": response})

	// Set refreshed state
	diags := utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource information.
func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "monitorResource.Update", map[string]any{"req": req, "resp": resp})

	var plan models.TFMonitorData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	var monitor uptrace.Monitor
	diags := utils.TFMonitorToUptraceMonitor(ctx, plan, &monitor)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response uptrace.MonitorIdResponse
	err := r.client.UpdateMonitor(ctx, id, monitor, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update monitor",
			fmt.Sprintf("Failed to update monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "UpdateMonitor OK", map[string]any{"response": response})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete resource information.
func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "monitorResource.Delete", map[string]any{"req": req, "resp": resp})

	var state models.TFMonitorData
	// Read Terraform plan data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	err := r.client.DeleteMonitor(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete monitor",
			fmt.Sprintf("Failed to delete monitor: %s", err),
		)
		return
	}

	// log the response
	var response any
	tflog.Info(ctx, "DeleteMonitor OK", map[string]any{"response": response})
}

func (r *monitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "monitorResource.ImportState", map[string]any{"req": req, "resp": resp})

	// Get fresh state from Uptrace
	id := req.ID
	var response uptrace.GetMonitorByIdResponse
	err := r.client.GetMonitorById(ctx, id, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get monitor",
			fmt.Sprintf("Failed to get monitor: %s", err),
		)
		return
	}

	// Save data into Terraform state
	var state models.TFMonitorData
	diags := utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
