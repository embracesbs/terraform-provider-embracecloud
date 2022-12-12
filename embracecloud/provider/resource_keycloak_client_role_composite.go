package provider

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakClientRoleComposite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakClientRoleCompositeCreate,
		ReadContext:   resourceKeycloakClientRoleCompositeRead,
		DeleteContext: resourceKeycloakClientRoleCompositeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakClientRoleCompositeImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"composite_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"composite_role_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceKeycloakClientRoleCompositeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	roleName := data.Get("role_name").(string)
	clientId := data.Get("client_id").(string)
	compositeClientId, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

	role, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, clientId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {

		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, compositeClientId.(string), composteRoleName)
		if err != nil {
			return diag.FromErr(err)
		}

		compRole = append(compRole, *compRoleResponse)

		err = keycloakCLient.AddClientRoleComposite(ctx, token.AccessToken, realm, *role.ID, compRole)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		compRoleResponse, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, composteRoleName)
		if err != nil {
			return diag.FromErr(err)
		}

		compRole = append(compRole, *compRoleResponse)
		err = keycloakCLient.AddRealmRoleComposite(ctx, token.AccessToken, data.Get("realm_id").(string), *role.Name, compRole)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	data.SetId(roleName)

	return nil

}

func resourceKeycloakClientRoleCompositeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceKeycloakClientRoleCompositeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	roleName := data.Get("role_name").(string)
	clientId := data.Get("client_id").(string)
	compositeClientId, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

	role, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, clientId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {
		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, compositeClientId.(string), composteRoleName)
		if err != nil {
			return diag.FromErr(err)
		}

		compRole = append(compRole, *compRoleResponse)

		err = keycloakCLient.DeleteClientRoleComposite(ctx, token.AccessToken, realm, *role.ID, compRole)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		compRoleResponse, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, composteRoleName)
		if err != nil {
			return diag.FromErr(err)
		}

		compRole = append(compRole, *compRoleResponse)
		keycloakCLient.DeleteRealmRoleComposite(ctx, token.AccessToken, realm, *role.Name, compRole)
	}
	return nil
}

func resourceKeycloakClientRoleCompositeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
