package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v12"
	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakRealmRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakRealmRoleCreate,
		ReadContext:   resourceKeycloakRealmRoleRead,
		DeleteContext: resourceKeycloakRealmRoleDelete,
		UpdateContext: resourceKeycloakRealmRoleUpdate,
		// This resource can be imported using {{realm}}/{{roleId}}. The role's ID (a GUID) can be found in the URL when viewing the role
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakRealmRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// misc attributes
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapRole(data *schema.ResourceData) (rl gocloak.Role, realm string) {

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

func mapFromRoleToData(data *schema.ResourceData, role gocloak.Role) {
	attributes := map[string]string{}
	for k, v := range *role.Attributes {
		attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
	}

	data.Set("realm_id", data.Get("realm_id").(string))
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("attributes", attributes)
}

func resourceKeycloakRealmRoleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapRole(data)

	id, err := keycloakCLient.CreateRealmRole(ctx, token.AccessToken, realm,
		role)

	if err != nil {
		return diag.Errorf(fmt.Sprintf("could not create realm role %s in realm %s error -> %s", *role.Name, realm, err.Error()))
	}

	data.SetId(id)

	return resourceKeycloakRealmRoleRead(ctx, data, meta)

}

func resourceKeycloakRealmRoleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()

	role, err := keycloakCLient.GetRealmRole(ctx, token.AccessToken, data.Get("realm_id").(string), data.Id())
	if err != nil {
		return diag.Errorf(fmt.Sprint("failed to get realm role %s in realm %s error -> %s", role.Name, data.Get("realm_id").(string), err.Error()))
	}

	//mapFromRoleToData(data, *role)
	return nil
}

func resourceKeycloakRealmRoleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapRole(data)

	err := keycloakCLient.UpdateRealmRole(ctx, token.AccessToken, realm, *role.ID, role)

	if err != nil {
		return diag.Errorf(fmt.Sprintf("could not update realm role %s in realm %s error -> %s", *role.Name, realm, err.Error()))
	}

	return resourceKeycloakRealmRoleRead(ctx, data, meta)
}

func resourceKeycloakRealmRoleDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*embracecloud.EmbraceCloudClient)
	keycloakCLient, token := client.GetKeycloakClient()
	role, realm := mapRole(data)
	err := keycloakCLient.DeleteRealmRole(ctx, token.AccessToken, realm, *role.ID)
	if err != nil {
		return diag.Errorf(fmt.Sprintf("could not delete realm role %s in realm %s error -> %s", *role.Name, realm, err.Error()))
	}
	return nil
}

func resourceKeycloakRealmRoleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
