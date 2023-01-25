package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type dataSourceItemType struct{}

func (d dataSourceItemType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about a 1password Item.
`,

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"uuid": {
				Description: "The UUID of the item to retrieve. This field will be populated with the UUID of the item if the item it looked up by its title.",
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Optional:    true,
				Type:        types.StringType,
			},
			"title": {
				Description: "The title of the item to retrieve. This field will be populated with the title of the item if the item it looked up by its UUID.",
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Optional:    true,
				Type:        types.StringType,
			},
			"vault": {
				Required: true,
				Type:     types.StringType,
			},
			"category": {
				Description: "The category of the item. One of [\"login\" \"password\" \"database\"]",
				Computed:    true,
				Type:        types.StringType,
			},
			//    database (String, Read-only) (Only applies to the database category) The name of the database.
			"hostname": {
				Description: "(Only applies to the database category) The address where the database can be found",
				Computed:    true,
				Type:        types.StringType,
			},
			//    password (String, Read-only) Password for this item.
			"port": {
				Description: "(Only applies to the database category) The port the database is listening on.",
				Computed:    true,
				Type:        types.StringType,
			},
			"type": {
				Description: "(Only applies to the database category) The type of database. One of [\"db2\" \"filemaker\" \"msaccess\" \"mssql\" \"mysql\" \"oracle\" \"postgresql\" \"sqlite\" \"other\"]",
				Computed:    true,
				Type:        types.StringType,
			},
			"url": {
				Description: "The primary URL for the item.",
				Computed:    true,
				Type:        types.StringType,
			},
			"username": {
				Description: "Username for this item.",
				Computed:    true,
				Type:        types.StringType,
			},
			"password": {
				Description: "Password for this item.",
				Computed:    true,
				Type:        types.StringType,
				Sensitive:   true,
			},
			"section": {
				Computed: true,
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"field": {
						Computed: true,
						Optional: true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								Type:     types.StringType,
								Optional: true,
								Computed: true,
							},
						}),
					},
				}),
			},
			//			"fields": {
			//				Computed: true,
			//				Optional: true,
			//				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
			//					"id": {
			//						Type:     types.StringType,
			//						Optional: true,
			//						Computed: true,
			//					},
			//					"label": {
			//						Type:     types.StringType,
			//						Optional: true,
			//						Computed: true,
			//					},
			//					"type": {
			//						Type:     types.StringType,
			//						Optional: true,
			//						Computed: true,
			//					},
			//					"purpose": {
			//						Type:     types.StringType,
			//						Optional: true,
			//						Computed: true,
			//					},
			//					"value": {
			//						Type:     types.StringType,
			//						Optional: true,
			//						Computed: true,
			//					},
			//				}),
			//			},
		},
	}, nil
}

func (d dataSourceItemType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceItem{
		p: *prv,
	}, nil
}

type dataSourceItem struct {
	p provider
}

func (d dataSourceItem) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfItem Item
	diags := req.Config.Get(ctx, &tfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Item.
	item, err := d.p.repo.GetItemByTitle(ctx, tfItem.VaultID.Value, tfItem.Title.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error getting item", "Could not get item, unexpected error: "+err.Error())
		return
	}

	newTfItem := mapModelToTfItem(*item)

	diags = resp.State.Set(ctx, newTfItem)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
