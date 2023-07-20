package provider

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v12"
	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakServiceAccountDetails() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakServiceAccountDetailsCreate,
		ReadContext:   resourceKeycloakServiceAccountDetailsRead,
		DeleteContext: resourceKeycloakServiceAccountDetailsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakServiceAccountDetailsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"last_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
				
			},
		},
	}
}

func resourceKeycloakServiceAccountDetailsCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	firstName := data.Get("first_name").(string)
	lastName := data.Get("last_name").(string)

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

	serviceAccountUser, err := keycloakCLient.GetClientServiceAccount(ctx, token.AccessToken, realm, *clients[0].ID)
	if err != nil {

	}

	serviceAccountUser.FirstName = &firstName
	serviceAccountUser.LastName = &lastName

	err = keycloakCLient.UpdateUser(ctx, token.AccessToken, realm, *serviceAccountUser)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("username", *serviceAccountUser.Username)
	data.SetId(*serviceAccountUser.ID)

	return nil
}

func resourceKeycloakServiceAccountDetailsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	userId := data.Get("id").(string)
	realm := data.Get("realm_id").(string)

	user, err := keycloakCLient.GetUserByID(ctx, token.AccessToken, realm, userId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("first_name", user.FirstName)
	data.Set("last_name", user.LastName)

	return nil
}

func resourceKeycloakServiceAccountDetailsDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	realm := data.Get("realm_id").(string)
	userId := data.Get("id").(string)

	data.Set("first_name", "")
	data.Set("last_name", "")

	firstName := data.Get("first_name").(string)
	lastName := data.Get("last_name").(string)

	user, err := keycloakCLient.GetUserByID(ctx, token.AccessToken, realm, userId)
	if err != nil {
		return diag.FromErr(err)
	}

	user.FirstName = &firstName
	user.LastName = &lastName

	keycloakCLient.UpdateUser(ctx, token.AccessToken, realm, *user)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKeycloakServiceAccountDetailsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
