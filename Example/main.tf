terraform {
  required_providers {
    embracecloud = {
      source = "terraform.example.com/local/embracecloud"
      version = "0.1"
    }
  }


}
provider "embracecloud" {
    keycloack_enabled = true
    keycloak_client_id =""
    keycloak_client_secret = ""
    keycloak_url =""
  
}

resource "embracecloud_realm_role" "comp1" {

    realm_id = "embracecloud"
    name = "testembracecloudprovider223"
    description = "test223"
    attributes = {
        test = "test",
        test2 = "test2"
    }
  
}


resource "embracecloud_realm_role" "comp2" {

    realm_id = "embracecloud"
    name = "testembracecloudprovider223composite"
    description = "test223composite"
    attributes = {
        test = "test",
        test2 = "test2"
    }
  
}


resource "embracecloud_realm_role_composite" "test" {
  realm_id = "embracecloud"
  parent_role_name = embracecloud_realm_role23.comp1.name
  composite_role_name = embracecloud_realm_role23.comp2.name
}


resource "embracecloud_realm_role_composite" "test_client_role" {
  realm_id = "embracecloud"
  parent_role_name = embracecloud_realm_role23.comp1.name
  composite_client_id = "5a113750-8100-4e2d-ac63-20f5d1315028"
  composite_role_name = "Huurder"
}


