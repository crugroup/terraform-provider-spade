// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeFileFormatResource{}
var _ resource.ResourceWithImportState = &SpadeFileFormatResource{}

func NewSpadeFileFormatResource() resource.Resource {
	return &SpadeFileFormatResource{}
}

// SpadeFileFormatResource defines the resource implementation.
type SpadeFileFormatResource struct {
	Client *SpadeClient
}

// SpadeFileFormatResourceModel describes the resource data model.
type SpadeFileFormatResourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Format types.String `tfsdk:"format"`
}

func (r *SpadeFileFormatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_format"
}

func (r *SpadeFileFormatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Spade file format",

		Attributes: map[string]schema.Attribute{
			"format": schema.StringAttribute{
				MarkdownDescription: "File format name",
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

func (r *SpadeFileFormatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*SpadeClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *SpadeClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.Client = client
}

func (r *SpadeFileFormatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeFileFormatResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.CreateFileFormat(data.Format.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create file format, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Format = types.StringValue(spadeResp.Format)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileFormatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeFileFormatResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadFileFormat(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file format, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Format = types.StringValue(spadeResp.Format)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileFormatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeFileFormatResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.UpdateFileFormat(
		data.Id.ValueInt64(),
		data.Format.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update file format, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Format = types.StringValue(spadeResp.Format)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileFormatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeFileFormatResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteFileFormat(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete file format, got error: %s", err))
		return
	}
}

func (r *SpadeFileFormatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
