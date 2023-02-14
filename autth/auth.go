package auth

import (
	"encoding/json"
	"fmt"
	"os"
)

type MyJsonName struct {
	Webb Web `json:"web"`
}

type Web struct {
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	AuthURI                 string   `json:"auth_uri"`
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret"`
	JavascriptOrigins       []string `json:"javascript_origins"`
	ProjectID               string   `json:"project_id"`
	RedirectUris            []string `json:"redirect_uris"`
	TokenURI                string   `json:"token_uri"`
}

func MashallXIV() Web {
	file, err := os.Open("./middleware/google.json")

	if err != nil {
		fmt.Println("Error opening JSON file:", err)
	}
	defer file.Close()

	data := MyJsonName{}

	web := Web{}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
	}

	web.AuthProviderX509CertURL = data.Webb.AuthProviderX509CertURL
	web.AuthURI = data.Webb.AuthURI
	web.ClientID = data.Webb.ClientID
	web.ClientSecret = data.Webb.ClientSecret
	web.RedirectUris = data.Webb.RedirectUris

	return web
}
