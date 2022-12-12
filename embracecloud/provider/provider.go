package provider

import (
	"context"

	"github.com/embracesbs/terraform-provider-embracecloud/embracecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MULTIVALUE_ATTRIBUTE_SEPARATOR = "##"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"keycloack_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EMBRACECLOUD_KEYCLOACK_ENABLED", false),
				Description: "Enable keycloak functionality within the provider",
			},
			"keycloak_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EMBRACECLOUD_KEYCLOACK_URL", ""),
				Description: "url of the keycloack intance",
			},
			"keycloak_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EMBRACECLOUD_KEYCLOACK_CLIENT_ID", ""),
				Description: "client id",
			},
			"keycloak_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EMBRACECLOUD_KEYCLOACK_CLIENT_SECRET", ""),
				Description: "client secret",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"embracecloud_realm_role":            resourceKeycloakRealmRole(),
			"embracecloud_realm_role_composite":  resourceKeycloakRealmRoleComposite(),
			"embracecloud_client_role":           resourceKeycloakClientRole(),
			"embracecloud_client_role_composite": resourceKeycloakClientRoleComposite(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	embraceCloudClient := embracecloud.BuildClient()

	if d.Get("keycloack_enabled").(bool) == true {
		embraceCloudClient.InitKeycloak(
			ctx,
			d.Get("keycloak_url").(string),
			d.Get("keycloak_client_id").(string),
			d.Get("keycloak_client_secret").(string))
	}

	return embraceCloudClient, diags

}
