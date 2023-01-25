package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

var permissionsAttribute = &schema.Resource{
	Description: `The permissions of the access. Note: Not all permissions are available in all plans, and some permissions require others. More info in [1password docs](https://developer.1password.com/docs/cli/vault-permissions/).`,
	Schema: map[string]*schema.Schema{
		"allow_viewing":           {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"allow_editing":           {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"allow_managing":          {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"view_items":              {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"create_items":            {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"edit_items":              {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"archive_items":           {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"delete_items":            {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"view_and_copy_passwords": {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"view_item_history":       {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"import_items":            {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"export_items":            {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"copy_and_share_items":    {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"print_items":             {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
		"manage_vault":            {Type: schema.TypeBool, Computed: true, Optional: true, Default: false},
	},
}

func mapTfToModelAccessPermissions(ap AccessPermissions) model.AccessPermissions {
	return model.AccessPermissions{
		AllowViewing:         ap.AllowViewing,
		AllowEditing:         ap.AllowEditing,
		AllowManaging:        ap.AllowManaging,
		ViewItems:            ap.ViewItems,
		CreateItems:          ap.CreateItems,
		EditItems:            ap.EditItems,
		ArchiveItems:         ap.ArchiveItems,
		DeleteItems:          ap.DeleteItems,
		ViewAndCopyPasswords: ap.ViewAndCopyPasswords,
		ViewItemHistory:      ap.ViewItemHistory,
		ImportItems:          ap.ImportItems,
		ExportItems:          ap.ExportItems,
		CopyAndShareItems:    ap.CopyAndShareItems,
		PrintItems:           ap.PrintItems,
		ManageVault:          ap.ManageVault,
	}
}

func mapModelToTfAccessPermissions(m model.AccessPermissions) *AccessPermissions {
	return &AccessPermissions{
		AllowViewing:         types.Bool{Value: m.AllowViewing},
		AllowEditing:         types.Bool{Value: m.AllowEditing},
		AllowManaging:        types.Bool{Value: m.AllowManaging},
		ViewItems:            types.Bool{Value: m.ViewItems},
		CreateItems:          types.Bool{Value: m.CreateItems},
		EditItems:            types.Bool{Value: m.EditItems},
		ArchiveItems:         types.Bool{Value: m.ArchiveItems},
		DeleteItems:          types.Bool{Value: m.DeleteItems},
		ViewAndCopyPasswords: types.Bool{Value: m.ViewAndCopyPasswords},
		ViewItemHistory:      types.Bool{Value: m.ViewItemHistory},
		ImportItems:          types.Bool{Value: m.ImportItems},
		ExportItems:          types.Bool{Value: m.ExportItems},
		CopyAndShareItems:    types.Bool{Value: m.CopyAndShareItems},
		PrintItems:           types.Bool{Value: m.PrintItems},
		ManageVault:          types.Bool{Value: m.ManageVault},
	}
}
