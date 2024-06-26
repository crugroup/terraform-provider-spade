// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	spade "terraform-provider-spade/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeProcessResource{}
var _ resource.ResourceWithImportState = &SpadeProcessResource{}

func NewSpadeProcessResource() resource.Resource {
	return &SpadeProcessResource{}
}

// SpadeProcessResource defines the resource implementation.
type SpadeProcessResource struct {
	Client *spade.SpadeClient
}

// SpadeProcessResourceModel describes the resource data model.
type SpadeProcessResourceModel struct {
	Id           types.Int64          `tfsdk:"id"`
	Code         types.String         `tfsdk:"code"`
	Description  types.String         `tfsdk:"description"`
	Tags         types.Set            `tfsdk:"tags"`
	Executor     types.Int64          `tfsdk:"executor"`
	SystemParams jsontypes.Normalized `tfsdk:"system_params"`
	UserParams   jsontypes.Normalized `tfsdk:"user_params"`
}

func (r *SpadeProcessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_process"
}

func (r *SpadeProcessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents a user-facing process within Spade",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				MarkdownDescription: "Name of the process",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the process",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags for the process",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(basetypes.NewSetValueMust(types.StringType, []attr.Value{})),
			},
			"executor": schema.Int64Attribute{
				MarkdownDescription: "Identifier to the underlying executor",
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
				MarkdownDescription: "Identifier of the process",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SpadeProcessResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SpadeProcessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeProcessResourceModel

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
		str, ok := tag.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert tag to string, please report issue to provider developers")
			return
		}
		tagStrings[i] = str.ValueString()
	}

	spadeResp, err := r.Client.CreateProcess(
		data.Code.ValueString(),
		data.Description.ValueString(),
		tagStrings,
		data.Executor.ValueInt64(),
		systemParamsJson,
		userParamsJson,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create process, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewSetValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse tags")
		return
	}
	data.Tags = respTags
	data.Executor = types.Int64Value(spadeResp.Executor)
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

func (r *SpadeProcessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeProcessResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadProcess(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read process, got error: %s", err))
		return
	}
	if spadeResp == nil {
		// Resource no longer exists, remove from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewSetValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Tags = respTags
	data.Executor = types.Int64Value(spadeResp.Executor)
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

func (r *SpadeProcessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeProcessResourceModel

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
		str, ok := tag.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert tag to string, please report issue to provider developers")
			return
		}
		tagStrings[i] = str.ValueString()
	}

	spadeResp, err := r.Client.UpdateProcess(
		data.Id.ValueInt64(),
		data.Code.ValueString(),
		data.Description.ValueString(),
		tagStrings,
		data.Executor.ValueInt64(),
		systemParamsJson,
		userParamsJson,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update process, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Code = types.StringValue(spadeResp.Code)
	data.Description = types.StringValue(spadeResp.Description)
	respTags, diag := basetypes.NewSetValueFrom(ctx, types.StringType, spadeResp.Tags)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse tags")
		return
	}
	data.Tags = respTags
	data.Executor = types.Int64Value(spadeResp.Executor)
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

func (r *SpadeProcessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeProcessResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteProcess(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete process, got error: %s", err))
		return
	}
}

func (r *SpadeProcessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			&SpadeProcessResourceModel{
				Id: types.Int64Value(id),
				// need to set something here, otherwise terraform can't infer the inner
				// type of the set and panics
				Tags: basetypes.NewSetValueMust(types.StringType, []attr.Value{}),
			},
		)...,
	)
}
