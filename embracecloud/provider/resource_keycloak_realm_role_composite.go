package provider

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v12"
	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakRealmRoleComposite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakRealmRoleCompositeCreate,
		ReadContext:   resourceKeycloakRealmRoleCompositeRead,
		DeleteContext: resourceKeycloakRealmRoleCompositeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakRealmRoleCompositeImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_role_name": {
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

func resourceKeycloakRealmRoleCompositeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	role_name := data.Get("parent_role_name").(string)
	composite_client_id, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

	role, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, role_name)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {
		var clientId = composite_client_id.(string)
		var params = gocloak.GetClientsParams{
			ClientID: &clientId,
		}

		clients, err := keycloakCLient.GetClients(ctx, token.AccessToken, realm, params)

		if err != nil {
			return diag.Errorf(fmt.Sprintf("cannot find client %s in realm %s", clientId, realm))
		}
		if len(clients) < 1 {
			return diag.Errorf(fmt.Sprintf("Client %s not found in realm %s", clientId, realm))
		}

		if len(clients) > 1 {
			return diag.Errorf("multiple clients found")
		}

		res, err := keycloakCLient.GetClient(ctx, token.AccessToken, realm, *clients[0].ID)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("cannot find client %s in realm %s", clientId, realm))
		}
		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *res.ID, composteRoleName)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not find client role in client %s with name %s in realm %s error -> %s", *res.ID, composteRoleName, realm, err.Error()))
		}

		compRole = append(compRole, *compRoleResponse)

		err = keycloakCLient.AddClientRoleComposite(ctx, token.AccessToken, realm, *role.ID, compRole)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Cannot add composite client role %s from client %s in realm %s error -> %s", *compRole[0].Name, clientId, realm, err.Error()))
		}

	} else {
		compRoleResponse, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, composteRoleName)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not find realm role %s in realm %s error -> %s", composteRoleName, realm, err.Error()))
		}

		compRole = append(compRole, *compRoleResponse)
		err = keycloakCLient.AddRealmRoleComposite(ctx, token.AccessToken, data.Get("realm_id").(string), *role.Name, compRole)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Cannot add composite %s to realmrole %s in realm %s error -> %s", *compRole[0].Name, *role.Name, realm, err.Error()))
		}
	}

	data.SetId(role_name)

	return nil

}

func resourceKeycloakRealmRoleCompositeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceKeycloakRealmRoleCompositeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	role_name := data.Get("parent_role_name").(string)
	composite_client_id, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

	role, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, role_name)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {
		var clientId = composite_client_id.(string)
		var params = gocloak.GetClientsParams{
			ClientID: &clientId,
		}

		clients, err := keycloakCLient.GetClients(ctx, token.AccessToken, realm, params)

		if err != nil {
			return diag.Errorf(fmt.Sprintf("cannot find client %s in realm %s", clientId, realm))
		}
		if len(clients) < 1 {
			return diag.Errorf(fmt.Sprintf("Client %s not found in realm %s", clientId, realm))
		}

		if len(clients) > 1 {
			return diag.Errorf("multiple clients found")
		}

		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *clients[0].ID, composteRoleName)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not find client role in client %s with name %s in realm %s error -> %s", clientId, composteRoleName, realm, err.Error()))
		}

		compRole = append(compRole, *compRoleResponse)

		err = keycloakCLient.DeleteClientRoleComposite(ctx, token.AccessToken, realm, *role.ID, compRole)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Cannot delete composite client role %s from client %s in realm %s error -> %s", *compRole[0].Name, clientId, realm, err.Error()))
		}

	} else {
		compRoleResponse, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, realm, composteRoleName)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not find realm role %s in realm %s error -> %s", composteRoleName, realm, err.Error()))
		}

		compRole = append(compRole, *compRoleResponse)
		err = keycloakCLient.DeleteRealmRoleComposite(ctx, token.AccessToken, realm, *role.Name, compRole)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not delete composite role %s from realmrole %s in realm %s error -> %s", *compRole[0].Name, *role.Name, realm, err.Error()))

		}
	}
	return nil
}

func resourceKeycloakRealmRoleCompositeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
