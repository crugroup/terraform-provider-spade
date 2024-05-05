// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"
	spade "terraform-provider-spade/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeFileProcessorResource{}
var _ resource.ResourceWithImportState = &SpadeFileProcessorResource{}

func NewSpadeFileProcessorResource() resource.Resource {
	return &SpadeFileProcessorResource{}
}

// SpadeFileProcessorResource defines the resource implementation.
type SpadeFileProcessorResource struct {
	Client *spade.SpadeClient
}

// SpadeFileProcessorResourceModel describes the resource data model.
type SpadeFileProcessorResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Callable    types.String `tfsdk:"callable"`
}

func (r *SpadeFileProcessorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_processor"
}

func (r *SpadeFileProcessorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Spade file processor",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the file processor",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the file processor",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"callable": schema.StringAttribute{
				MarkdownDescription: "Python import path to the FileProcessor class",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SpadeFileProcessorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SpadeFileProcessorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeFileProcessorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.CreateFileProcessor(
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.Callable.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create file processor, got error: %s", err))
		return
	}
	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileProcessorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeFileProcessorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadFileProcessor(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file processor, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileProcessorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeFileProcessorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.UpdateFileProcessor(
		data.Id.ValueInt64(),
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.Callable.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update file processor, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	data.Callable = types.StringValue(spadeResp.Callable)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileProcessorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeFileProcessorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteFileProcessor(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete file processor, got error: %s", err))
		return
	}
}

func (r *SpadeFileProcessorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected a numeric resource ID, got: %s", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			&SpadeFileProcessorResourceModel{
				Id: types.Int64Value(id),
			},
		)...,
	)
}
