// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"os"
	"terraform-provider-foxglove-cloud/internal/foxglove"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &FoxgloveProvider{}
var _ provider.ProviderWithFunctions = &FoxgloveProvider{}

type FoxgloveProvider struct {
	version string
}

type FoxgloveProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *FoxgloveProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "foxglove"
	resp.Version = p.version
}

func (p *FoxgloveProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Foxglove API Key. Can also be set via environment variable FOXGLOVE_API_KEY",
				Optional:            true,
			},
		},
	}
}

func (p *FoxgloveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FoxgloveProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("FOXGLOVE_API_KEY")

	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError("Foxglove api key missing",
			"The provider cannot create the Foxglove client as there is a missing or empty value for the Foxglove API key. "+
				"Set the api_key value in the configuration or use the FOXGLOVE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	foxgloveClient := foxglove.NewClient(apiKey)

	resp.DataSourceData = foxgloveClient
	resp.ResourceData = foxgloveClient
}

func (p *FoxgloveProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewApikeyResource,
	}
}

func (p *FoxgloveProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FoxgloveProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FoxgloveProvider{
			version: version,
		}
	}
}
