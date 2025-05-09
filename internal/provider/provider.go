// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/persona-ae/terraform-provider-uptrace/internal/resources"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

// Ensure UptraceProvider satisfies various provider interfaces.
var _ provider.Provider = &UptraceProvider{}
var _ provider.ProviderWithFunctions = &UptraceProvider{}
var _ provider.ProviderWithEphemeralResources = &UptraceProvider{}

// UptraceProvider defines the provider implementation.
type UptraceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *UptraceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uptrace"
	resp.Version = p.version
}

func (p *UptraceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for authentication.",
				Required:            true,
				Sensitive:           true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Uptrace project ID.",
				Required:            true,
			},
		},
	}
}

func (p *UptraceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config struct {
		APIKey    types.String `tfsdk:"api_key"`
		ProjectID types.String `tfsdk:"project_id"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := uptrace.NewUptraceClient(
		config.ProjectID.ValueString(),
		config.APIKey.ValueString(),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UptraceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewMonitorResource,
	}
}

func (p *UptraceProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *UptraceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *UptraceProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UptraceProvider{
			version: version,
		}
	}
}
