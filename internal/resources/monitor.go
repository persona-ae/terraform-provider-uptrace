package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
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
				Computed:    true,
				Description: "The name of the monitor.",
			},
			"project_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Project ID the monitor belongs to.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the monitor.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The monitor type.",
			},
			"query": schema.StringAttribute{
				Computed:    true,
				Description: "The monitor's query.",
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

	// Generate API request body from plan

	// Create new monitor
	// poll the describe monitor endpoint until the monitor is ready
	// Poll every n seconds
	// log the response
	// Save data into Terraform state
}

// Read resource information.
func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	// Read data from Terraform state
	// Get fresh state from Pinecone
	// Generate API request body from plan
	// log the response
	// Set refreshed state
}

// Update resource information.
func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	// log the response
	// Save updated data into Terraform state
}

// Delete resource information.
func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read Terraform plan data into the model
	// log the response
}

func (r *monitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// fetch a fresh index from pinecone
	// req.ID appears to be our only way to get the index name
	// note that what gets fetched from pinecone, based on purely
	// the index name, may differ from the rest of whatever is
	// specified in the resource stanza in HCL
	// Get fresh state from Pinecone
	// Save data into Terraform state
	// Retrieve import ID and save to id attribute
}
