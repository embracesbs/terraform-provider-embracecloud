package provider

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v12"
	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakClientRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakClientRoleCreate,
		ReadContext:   resourceKeycloakClientRoleRead,
		DeleteContext: resourceKeycloakClientRoleDelete,
		UpdateContext: resourceKeycloakClientRoleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakClientRoleImport,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapClientRole(data *schema.ResourceData) (rl gocloak.Role, realm string) {

	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = strings.Split(value.(string), MULTIVALUE_ATTRIBUTE_SEPARATOR)
		}
	}

	role := gocloak.Role{
		ID:          gocloak.StringP(data.Id()),
		Name:        gocloak.StringP(data.Get("name").(string)),
		Description: gocloak.StringP(data.Get("description").(string)),
		Attributes:  &attributes,
	}
	return role, data.Get("realm_id").(string)
}

func mapFromClientRoleToData(data *schema.ResourceData, role gocloak.Role) {
	attributes := map[string]string{}
	for k, v := range *role.Attributes {
		attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
	}

	data.Set("realm_id", data.Get("realm_id").(string))
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("attributes", attributes)
}

func resourceKeycloakClientRoleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapClientRole(data)
	clientId := data.Get("client_id").(string)

	id, err := keycloakCLient.CreateClientRole(ctx, token.AccessToken, realm, clientId,
		role)

	if err != nil {
		return diag.FromErr((err))
	}

	data.SetId(id)

	return resourceKeycloakRealmRoleRead(ctx, data, meta)

}

func resourceKeycloakClientRoleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	clientId := data.Get("client_id").(string)

	role, err := keycloakCLient.GetClientRole(ctx, token.AccessToken, data.Get("realm_id").(string), clientId, data.Id())
	if err != nil {
		return diag.FromErr((err))
	}

	mapFromRoleToData(data, *role)
	return nil
}

func resourceKeycloakClientRoleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapRole(data)

	err := keycloakCLient.UpdateRole(ctx, token.AccessToken, realm, *role.ID, role)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakRealmRoleRead(ctx, data, meta)
}

func resourceKeycloakClientRoleDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapRole(data)
	clientId := data.Get("client_id").(string)
	keycloakCLient.DeleteClientRole(ctx, token.AccessToken, realm, clientId, *role.ID)
	return nil
}

func resourceKeycloakClientRoleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
