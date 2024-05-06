// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	spade "terraform-provider-spade/internal/client"
)

// Ensure SpadeProvider satisfies various provider interfaces.
var _ provider.Provider = &SpadeProvider{}
var _ provider.ProviderWithFunctions = &SpadeProvider{}

// SpadeProvider defines the provider implementation.
type SpadeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SpadeProviderModel describes the provider data model.
type SpadeProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
}

func (p *SpadeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "spade"
	resp.Version = p.version
}

func (p *SpadeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for managing Spade connectors",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Spade URL",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Login email address",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Login password",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SpadeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SpadeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resource
	client := &spade.SpadeClient{
		ApiUrl:     data.URL.ValueString(),
		HttpClient: http.DefaultClient,
	}
	err := client.Login(data.Email.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Login Error", fmt.Sprintf("Login failed: %s", err))
		return
	}
	tflog.Info(ctx, "Logged in to Spade")
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SpadeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSpadeExecutorResource,
		NewSpadeFileProcessorResource,
		NewSpadeFileFormatResource,
		NewSpadeProcessResource,
		NewSpadeFileResource,
	}
}

func (p *SpadeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *SpadeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SpadeProvider{
			version: version,
		}
	}
}
