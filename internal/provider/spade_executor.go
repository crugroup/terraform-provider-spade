// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	spade "terraform-provider-spade/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeExecutorResource{}
var _ resource.ResourceWithImportState = &SpadeExecutorResource{}

func NewSpadeExecutorResource() resource.Resource {
	return &SpadeExecutorResource{}
}

// SpadeExecutorResource defines the resource implementation.
type SpadeExecutorResource struct {
	Client *spade.SpadeClient
}

// SpadeExecutorResourceModel describes the resource data model.
type SpadeExecutorResourceModel struct {
	Id                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Callable                types.String `tfsdk:"callable"`
	HistoryProviderCallable types.String `tfsdk:"history_provider_callable"`
}

func (r *SpadeExecutorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_executor"
}

func (r *SpadeExecutorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Spade executor",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the executor",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the executor",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"callable": schema.StringAttribute{
				MarkdownDescription: "Python import path to the Executor class",
				Required:            true,
			},
			"history_provider_callable": schema.StringAttribute{
				MarkdownDescription: "Python import path to the HistoryProvider class",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the executor",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SpadeExecutorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*spade.SpadeClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *spade.SpadeClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.Client = client
}

func (r *SpadeExecutorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeExecutorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.CreateExecutor(
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.Callable.ValueString(),
		data.HistoryProviderCallable.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create executor, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)
	data.HistoryProviderCallable = types.StringValue(spadeResp.HistoryProviderCallable)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeExecutorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeExecutorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadExecutor(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read executor, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)
	data.HistoryProviderCallable = types.StringValue(spadeResp.HistoryProviderCallable)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeExecutorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeExecutorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.UpdateExecutor(
		data.Id.ValueInt64(),
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.Callable.ValueString(),
		data.HistoryProviderCallable.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update executor, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)
	data.HistoryProviderCallable = types.StringValue(spadeResp.HistoryProviderCallable)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeExecutorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeExecutorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteExecutor(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete executor, got error: %s", err))
		return
	}
}

func (r *SpadeExecutorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
