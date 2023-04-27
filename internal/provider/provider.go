package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TrelloProvider defines the provider implementation.
type TrelloProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TrelloProviderModel describes the provider data model.
type TrelloProviderModel struct {}

func (p *TrelloProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.Version = p.version
}

func (p *TrelloProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *TrelloProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TrelloProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TrelloProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTerraformResource,
	}
}

func (p *TrelloProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TrelloProvider{
			version: version,
		}
	}
}
