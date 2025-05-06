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
	data_sources "github.com/persona-ae/terraform-provider-uptrace/internal/data-sources"
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

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *UptraceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *UptraceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "API token for authentication.",
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
		Token     types.String `tfsdk:"token"`
		ProjectID types.String `tfsdk:"project_id"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := uptrace.NewUptraceClient(
		config.ProjectID.ValueString(),
		config.Token.ValueString(),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UptraceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *UptraceProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *UptraceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data_sources.NewMonitorDataSource,
	}
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
