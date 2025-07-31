// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	spade "terraform-provider-spade/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &SpadeVariableDataSource{}

func NewSpadeVariableDataSource() datasource.DataSource {
	return &SpadeVariableDataSource{}
}

// SpadeVariableDataSource defines the data source implementation.
type SpadeVariableDataSource struct {
	Client *spade.SpadeClient
}

// SpadeVariableDataSourceModel describes the data source data model.
type SpadeVariableDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
	IsSecret    types.Bool   `tfsdk:"is_secret"`
}

func (d *SpadeVariableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (d *SpadeVariableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Variable data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Identifier of the variable",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the variable",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the variable",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the variable",
				Computed:            true,
			},
			"is_secret": schema.BoolAttribute{
				MarkdownDescription: "Whether the variable is secret",
				Computed:            true,
			},
		},
	}
}

func (d *SpadeVariableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*spade.SpadeClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *spade.SpadeClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.Client = client
}

func (d *SpadeVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SpadeVariableDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := d.Client.SearchVariable(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find variable, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Value = types.StringValue(spadeResp.Value)
	data.IsSecret = types.BoolValue(spadeResp.IsSecret)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
