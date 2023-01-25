package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func dataSourceVault() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides information about a 1password vault.
`,
		ReadContext: dataSourceVaultRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the vault.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"uuid": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataSourceVaultRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	name := data.Get("name").(string)
	vault, err := p.repo.GetVaultByName(ctx, name)
	if err != nil {
		return diag.Errorf("Error getting user: Could not get user, unexpected error: " + err.Error())
	}

	data.SetId(vaultTerraformID(vault))
	data.Set("uuid", vault.ID)
	data.Set("name", vault.Name)
	data.Set("description", vault.Description)
	return diags
}

func vaultTerraformID(vault *model.Vault) string {
	return fmt.Sprintf("vaults/%s", vault.ID)
}
