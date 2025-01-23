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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SpadeUserResource{}
var _ resource.ResourceWithImportState = &SpadeUserResource{}

func NewSpadeUserResource() resource.Resource {
	return &SpadeUserResource{}
}

// SpadeUserResource defines the resource implementation.
type SpadeUserResource struct {
	Client *spade.SpadeClient
}

// SpadeUserResourceModel describes the resource data model.
type SpadeUserResourceModel struct {
	Id        types.Int64  `tfsdk:"id"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Email     types.String `tfsdk:"email"`
	IsActive  types.Bool   `tfsdk:"active"`
	Groups    types.Set    `tfsdk:"groups"`
}

func (r *SpadeUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *SpadeUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents a user within Spade",

		Attributes: map[string]schema.Attribute{
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the account is active",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"groups": schema.SetAttribute{
				MarkdownDescription: "Group identifiers",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(basetypes.NewSetValueMust(types.Int64Type, []attr.Value{})),
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the user",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SpadeUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SpadeUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpadeUserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	groups := data.Groups.Elements()
	groupIDs := make([]int64, len(groups))
	for i, group := range groups {
		id, ok := group.(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert group ID to int, please report issue to provider developers")
			return
		}
		groupIDs[i] = id.ValueInt64()
	}

	spadeResp, err := r.Client.CreateUser(
		data.FirstName.ValueString(),
		data.LastName.ValueString(),
		data.Email.ValueString(),
		data.IsActive.ValueBool(),
		groupIDs,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.FirstName = types.StringValue(spadeResp.FirstName)
	data.LastName = types.StringValue(spadeResp.LastName)
	data.Email = types.StringValue(spadeResp.Email)
	data.IsActive = types.BoolValue(spadeResp.IsActive)
	respGroups, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Groups)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse groups")
		return
	}
	data.Groups = respGroups

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpadeUserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	spadeResp, err := r.Client.ReadUser(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}
	if spadeResp == nil {
		// Resource no longer exists, remove from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.FirstName = types.StringValue(spadeResp.FirstName)
	data.LastName = types.StringValue(spadeResp.LastName)
	data.Email = types.StringValue(spadeResp.Email)
	data.IsActive = types.BoolValue(spadeResp.IsActive)
	respGroups, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Groups)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse groups")
		return
	}
	data.Groups = respGroups

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpadeUserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	groups := data.Groups.Elements()
	groupIDs := make([]int64, len(groups))
	for i, group := range groups {
		id, ok := group.(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Client Error", "Failed to convert group ID to int, please report issue to provider developers")
			return
		}
		groupIDs[i] = id.ValueInt64()
	}

	spadeResp, err := r.Client.UpdateUser(
		data.Id.ValueInt64(),
		data.FirstName.ValueString(),
		data.LastName.ValueString(),
		data.Email.ValueString(),
		data.IsActive.ValueBool(),
		groupIDs,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.Id = types.Int64Value(spadeResp.Id)
	data.FirstName = types.StringValue(spadeResp.FirstName)
	data.LastName = types.StringValue(spadeResp.LastName)
	data.Email = types.StringValue(spadeResp.Email)
	data.IsActive = types.BoolValue(spadeResp.IsActive)
	respGroups, diag := basetypes.NewSetValueFrom(ctx, types.Int64Type, spadeResp.Groups)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Client Error", "Unable to parse groups")
		return
	}
	data.Groups = respGroups

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpadeUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpadeUserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.DeleteUser(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}
}

func (r *SpadeUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			&SpadeUserResourceModel{
				Id: types.Int64Value(id),
				// need to set something here, otherwise terraform can't infer the inner
				// type of the set and panics
				Groups: basetypes.NewSetValueMust(types.Int64Type, []attr.Value{}),
			},
		)...,
	)
}
