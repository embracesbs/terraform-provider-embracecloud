package embracecloud

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
)

type EmbraceCloudClient struct {
	keycloack         gocloak.GoCloak
	keycloak_token    gocloak.JWT
	keycloack_enabled bool
}

func (cc *EmbraceCloudClient) InitKeycloak(ctx context.Context, url string, clientId string, clientSecret string) {
	cc.keycloack = *gocloak.NewClient(url)
	token, err := cc.keycloack.LoginClient(ctx, clientId, clientSecret, "master")

	if err != nil {

	}
	cc.keycloak_token = *token
	cc.keycloack_enabled = true
}

func (cc *EmbraceCloudClient) GetKeycloakClient() (*gocloak.GoCloak, gocloak.JWT) {
	return &cc.keycloack, cc.keycloak_token
}

func BuildClient() *EmbraceCloudClient {
	return &EmbraceCloudClient{}
}
