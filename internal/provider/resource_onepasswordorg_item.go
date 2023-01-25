package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type resourceItemType struct{}

func (r resourceItemType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides a Item resource.

When a 1password user resources is created, it will be invited  by email.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"vault_id": {
				Type:        types.StringType,
				Required:    true,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description: "The name of the user.",
			},
			"title": {
				Type:        types.StringType,
				Required:    true,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description: "The name of the user.",
			},
			"fields": {
				Computed: true,
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"label": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"type": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"purpose": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"value": {
						Type:      types.StringType,
						Optional:  true,
						Computed:  true,
						Sensitive: true,
					},
				}),
			},
		},
	}, nil
}

func (r resourceItemType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceItem{
		p: *prv,
	}, nil
}

type resourceItem struct {
	p provider
}

func (r resourceItem) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfItem Item
	diags := req.Plan.Get(ctx, &tfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create user.
	u := mapTfToModelItem(tfItem)
	newItem, err := r.p.repo.CreateItem(ctx, u)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", "Could not create user, unexpected error: "+err.Error())
		return
	}

	// Map user to tf model.
	newTfItem := mapModelToTfItem(*newItem)

	diags = resp.State.Set(ctx, newTfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceItem) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfItem Item
	diags := req.State.Get(ctx, &tfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get user.
	id := tfItem.ID.Value
	user, err := r.p.repo.GetItemByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map user to tf model.
	readTfItem := mapModelToTfItem(*user)

	diags = resp.State.Set(ctx, readTfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceItem) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Get plan values.
	var plan Item
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state.
	var state Item
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan user as the new data and set ID from state.
	u := mapTfToModelItem(plan)
	u.ID = state.ID.Value

	newItem, err := r.p.repo.EnsureItem(ctx, u)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", "Could not update user, unexpected error: "+err.Error())
		return
	}

	// Map user to tf model.
	readTfItem := mapModelToTfItem(*newItem)

	diags = resp.State.Set(ctx, readTfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceItem) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfItem Item
	diags := req.State.Get(ctx, &tfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete user.
	id := tfItem.ID.Value
	err := r.p.repo.DeleteItem(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("Could not delete user %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceItem) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTfToModelItem(u Item) model.Item {
	return model.Item{
		ID:      u.ID.Value,
		VaultID: u.VaultID.Value,
		Fields:  mapTfToModelItemField(u.Section),
		Title:   u.Title.Value,
	}
}

func mapTfToModelItemField(u []Field) []model.Field {
	return []model.Field{}
	//	return model.Field{
	//		Label: u.Label.Value,
	//	}
}

func mapModelToTfItem(u model.Item) Item {
	return Item{
		ID:      types.String{Value: u.ID},
		VaultID: types.String{Value: u.VaultID},
		Title:   types.String{Value: u.Title},
		Section: mapModelToTfItemField(u.Fields),
	}
}

func mapModelToTfItemField(u []model.Field) []Field {
	fields := []Field{}
	for _, f := range u {
		fields = append(fields, Field{
			ID:      types.String{Value: f.ID},
			Label:   types.String{Value: f.Label},
			Type:    types.String{Value: f.Type},
			Value:   types.String{Value: f.Value},
			Purpose: types.String{Value: f.Purpose},
		})
	}
	return fields
}
