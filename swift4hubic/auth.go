package swift4hubic

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
)

func NewOAuthConfig(hubicApplication *HubicApplication) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  hubicApplication.RedirectURL,
		ClientID:     hubicApplication.ClientID,
		ClientSecret: hubicApplication.ClientSecret,
		Scopes:       []string{"usage.r,account.r,credentials.r"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.hubic.com/oauth/auth",
			TokenURL: "https://api.hubic.com/oauth/token",
		},
	}
}

func NewSwiftToken(client *http.Client) (*SwiftV1Token, error) {
	//TODO Store Token
	url := "https://api.hubic.com/1.0/account/credentials"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token SwiftV1Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	Log(LogDebug, "Succesfully fetched Swift Token from %s", url)
	return &token, err
}
