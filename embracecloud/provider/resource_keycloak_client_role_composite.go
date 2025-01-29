package provider

import (
	"context"
	"fmt"

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
			"parent_role_name": {
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
	roleName := data.Get("parent_role_name").(string)
	clientId := data.Get("client_id").(string)
	compositeClientId, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

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

	role, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *clients[0].ID, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {
		var compClientId = compositeClientId.(string)
		var params = gocloak.GetClientsParams{
			ClientID: &compClientId,
		}

		compClients, err := keycloakCLient.GetClients(ctx, token.AccessToken, realm, params)

		if err != nil {
			return diag.Errorf(fmt.Sprintf("cannot find client %s in realm %s", clientId, realm))
		}
		if len(clients) < 1 {
			return diag.Errorf(fmt.Sprintf("Client %s not found in realm %s", clientId, realm))
		}

		if len(clients) > 1 {
			return diag.Errorf("multiple clients found")
		}

		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *compClients[0].ID, composteRoleName)
		if err != nil {
			return diag.Errorf(fmt.Sprintf("Could not find client role in client %s with name %s in realm %s error -> %s", *compClients[0].ID, composteRoleName, realm, err.Error()))
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

	data.SetId(roleName)

	return nil

}

func resourceKeycloakClientRoleCompositeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakClient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	roleName := data.Get("parent_role_name").(string)
	clientId := data.Get("client_id").(string)
	compositeClientId, isClient := data.GetOkExists("composite_client_id")
	compositeRoleName := data.Get("composite_role_name").(string)

	// Fetch Client
	var params = gocloak.GetClientsParams{
		ClientID: &clientId,
	}

	clients, err := keycloakClient.GetClients(ctx, token.AccessToken, realm, params)
	if err != nil || len(clients) < 1 {
		return diag.Errorf(fmt.Sprintf("Client %s not found in realm %s. Marking resource for recreation.", clientId, realm))
	}

	if len(clients) > 1 {
		return diag.Errorf("Multiple clients found for ID %s in realm %s", clientId, realm)
	}

	role, err := keycloakClient.GetClientRole(ctx, token.AccessToken, realm, *clients[0].ID, roleName)
	if err != nil {
		data.SetId("")
	}
	return nil
}

func resourceKeycloakClientRoleCompositeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	roleName := data.Get("parent_role_name").(string)
	clientId := data.Get("client_id").(string)
	compositeClientId, isClient := data.GetOkExists("composite_client_id")
	composteRoleName := data.Get("composite_role_name").(string)

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

	role, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *clients[0].ID, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	var compRole []gocloak.Role

	if isClient == true {
		var compClientId = compositeClientId.(string)
		var params = gocloak.GetClientsParams{
			ClientID: &compClientId,
		}

		compClients, err := keycloakCLient.GetClients(ctx, token.AccessToken, realm, params)

		if err != nil {

			return diag.Errorf(fmt.Sprintf("Could not find client %s in realm %s error -> %s", clientId, realm, err.Error()))
		}
		if len(clients) < 1 {
			return diag.Errorf(fmt.Sprintf("Client %s not found in realm %s", clientId, realm))
		}

		if len(clients) > 1 {
			return diag.Errorf("multiple clients found")
		}
		compRoleResponse, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, realm, *compClients[0].ID, composteRoleName)
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

func resourceKeycloakClientRoleCompositeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
