// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"
	spade "terraform-provider-spade/internal/client"

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
var _ resource.Resource = &SpadeVariableSetResource{}
var _ resource.ResourceWithImportState = &SpadeVariableSetResource{}

func NewSpadeVariableSetResource() resource.Resource {
	return &SpadeVariableSetResource{}
}

// SpadeVariableSetResource defines the resource implementation.
type SpadeVariableSetResource struct {
	Client *spade.SpadeClient
}

// SpadeVariableSetResourceModel describes the resource data model.
type SpadeVariableSetResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Variables   types.Set    `tfsdk:"variables"`
}

func (r *SpadeVariableSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_set"
}

func (r *SpadeVariableSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents a variable set within Spade",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the variable set",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the variable set",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the variable set",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"variables": schema.SetAttribute{
				MarkdownDescription: "Variable identifiers",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(basetypes.NewSetValueMust(types.Int64Type, []attr.Value{})),
			},
		},
	}
}

func (r *SpadeVariableSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SpadeVariableSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeVariableSetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	variables := data.Variables.Elements()
	variableIDs := make([]int64, len(variables))
	for i, variable := range variables {
		id, ok := variable.(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert variable ID to int, please report issue to provider developers")
			return
		}
		variableIDs[i] = id.ValueInt64()
	}

	spadeResp, err := r.Client.CreateVariableSet(
		data.Name.ValueString(),
		data.Description.ValueString(),
		variableIDs,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variable set, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	respVariables, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Variables)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse variables")
		return
	}
	data.Variables = respVariables

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeVariableSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeVariableSetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadVariableSet(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read variable set, got error: %s", err))
		return
	}
	if spadeResp == nil {
		// Resource no longer exists, remove from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	respVariables, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Variables)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse variables")
		return
	}
	data.Variables = respVariables

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeVariableSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeVariableSetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	variables := data.Variables.Elements()
	variableIDs := make([]int64, len(variables))
	for i, variable := range variables {
		id, ok := variable.(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert variable ID to int, please report issue to provider developers")
			return
		}
		variableIDs[i] = id.ValueInt64()
	}

	spadeResp, err := r.Client.UpdateVariableSet(
		data.Id.ValueInt64(),
		data.Name.ValueString(),
		data.Description.ValueString(),
		variableIDs,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update variable set, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.Name = types.StringValue(spadeResp.Name)
	data.Description = types.StringValue(spadeResp.Description)
	respVariables, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Variables)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse variables")
		return
	}
	data.Variables = respVariables

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeVariableSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeVariableSetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteVariableSet(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete variable set, got error: %s", err))
		return
	}
}

func (r *SpadeVariableSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			&SpadeVariableSetResourceModel{
				Id: types.Int64Value(id),
				// need to set something here, otherwise terraform can't infer the inner
				// type of the set and panics
				Variables: basetypes.NewSetValueMust(types.Int64Type, []attr.Value{}),
			},
		)...,
	)
}
