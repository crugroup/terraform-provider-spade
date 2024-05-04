// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeFileResource{}
var _ resource.ResourceWithImportState = &SpadeFileResource{}

func NewSpadeFileResource() resource.Resource {
	return &SpadeFileResource{}
}

// SpadeFileResource defines the resource implementation.
type SpadeFileResource struct {
	Client *SpadeClient
}

// SpadeFileResourceModel describes the resource data model.
type SpadeFileResourceModel struct {
	Id           types.Int64          `tfsdk:"id"`
	Code         types.String         `tfsdk:"code"`
	Description  types.String         `tfsdk:"description"`
	Tags         types.List           `tfsdk:"tags"`
	Format       types.Int64          `tfsdk:"format"`
	Processor    types.Int64          `tfsdk:"processor"`
	SystemParams jsontypes.Normalized `tfsdk:"system_params"`
	UserParams   jsontypes.Normalized `tfsdk:"user_params"`
}

func (r *SpadeFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *SpadeFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Spade file",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				MarkdownDescription: "Name of the processor",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the processor",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags for the processor",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(basetypes.NewListValueMust(types.StringType, []attr.Value{})),
			},
			"format": schema.Int64Attribute{
				MarkdownDescription: "Identifier for file format",
				Required:            true,
			},
			"processor": schema.Int64Attribute{
				MarkdownDescription: "Identifier for processor",
				Required:            true,
			},
			"system_params": schema.StringAttribute{
				MarkdownDescription: "JSON of system parameters",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				Default:             stringdefault.StaticString("{}"),
			},
			"user_params": schema.StringAttribute{
				MarkdownDescription: "JSON of user parameters (JsonSchema form)",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				Default:             stringdefault.StaticString("{}"),
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

func (r *SpadeFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SpadeFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeFileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var systemParamsJson map[string]interface{}
	err := json.Unmarshal([]byte(data.SystemParams.ValueString()), &systemParamsJson)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse system_params, got error: %s", err))
		return
	}
	var userParamsJson map[string]interface{}
	err = json.Unmarshal([]byte(data.UserParams.ValueString()), &userParamsJson)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse user_params, got error: %s", err))
		return
	}
	tags := data.Tags.Elements()
	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		str := tag.(types.String).ValueString()
		tagStrings[i] = str
	}

	spadeResp, err := r.Client.CreateFile(
		data.Code.ValueString(),
		data.Description.ValueString(),
		tagStrings,
		data.Format.ValueInt64(),
		data.Processor.ValueInt64(),
		systemParamsJson,
		userParamsJson,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create file, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewListValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "fuck1")
		return
	}
	data.Tags = respTags
	data.Format = types.Int64Value(spadeResp.Format)
	data.Processor = types.Int64Value(spadeResp.Processor)
	respSystemParams, err := json.Marshal(spadeResp.SystemParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal system_params, got error: %s", err))
		return
	}
	respUserParams, err := json.Marshal(spadeResp.UserParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal user_params, got error: %s", err))
		return
	}
	data.SystemParams = jsontypes.NewNormalizedValue(string(respSystemParams))
	data.UserParams = jsontypes.NewNormalizedValue(string(respUserParams))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeFileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadFile(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewListValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Tags = respTags
	data.Format = types.Int64Value(spadeResp.Format)
	data.Processor = types.Int64Value(spadeResp.Processor)
	respSystemParams, err := json.Marshal(spadeResp.SystemParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal system_params, got error: %s", err))
		return
	}
	respUserParams, err := json.Marshal(spadeResp.UserParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal user_params, got error: %s", err))
		return
	}
	data.SystemParams = jsontypes.NewNormalizedValue(string(respSystemParams))
	data.UserParams = jsontypes.NewNormalizedValue(string(respUserParams))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeFileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var systemParamsJson map[string]interface{}
	err := json.Unmarshal([]byte(data.SystemParams.ValueString()), &systemParamsJson)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse system_params, got error: %s", err))
		return
	}
	var userParamsJson map[string]interface{}
	err = json.Unmarshal([]byte(data.UserParams.ValueString()), &userParamsJson)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse user_params, got error: %s", err))
		return
	}
	tags := data.Tags.Elements()
	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		str := tag.(types.String).ValueString()
		tagStrings[i] = str
	}

	spadeResp, err := r.Client.UpdateFile(
		data.Id.ValueInt64(),
		data.Code.ValueString(),
		data.Description.ValueString(),
		tagStrings,
		data.Format.ValueInt64(),
		data.Processor.ValueInt64(),
		systemParamsJson,
		userParamsJson,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update file, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewListValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "fuck1")
		return
	}
	data.Tags = respTags
	data.Format = types.Int64Value(spadeResp.Format)
	data.Processor = types.Int64Value(spadeResp.Processor)
	respSystemParams, err := json.Marshal(spadeResp.SystemParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal system_params, got error: %s", err))
		return
	}
	respUserParams, err := json.Marshal(spadeResp.UserParams)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to marshal user_params, got error: %s", err))
		return
	}
	data.SystemParams = jsontypes.NewNormalizedValue(string(respSystemParams))
	data.UserParams = jsontypes.NewNormalizedValue(string(respUserParams))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeFileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteFile(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete file, got error: %s", err))
		return
	}
}

func (r *SpadeFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
