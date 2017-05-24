package swift4hubic

import (
	"fmt"
	"io/ioutil"

	"encoding/json"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
	"time"
)

type SwiftV1Token struct {
	AuthToken  string    `json:"token"`
	StorageURL string    `json:"endpoint"`
	Expires    time.Time `json:"expires"`
}

type HubicApplication struct {
	Account      string       `yaml:"account"`
	Password     string       `yaml:"password"`
	ClientID     string       `yaml:"client_id"`
	ClientSecret string       `yaml:"client_secret"`
	RedirectURL  string       `yaml:"redirect_url"`
	OAuthToken   string       `yaml:"oauth_token"`
	SwiftToken   SwiftV1Token `yaml:"swift_token"`
}

type Configuration struct {
	Port              string              `yaml:"port"`
	HubicApplications []*HubicApplication `yaml:"hubic"`
}

func GetConfiguration(fileName string) (*Configuration, error) {
	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}

	if !cfg.validate() {
		return nil, fmt.Errorf("Configuration is invalid")
	}
	return &cfg, nil
}

func (cfg Configuration) validate() (success bool) {
	success = true
	invalid := func(msg string) {
		Log(LogError, msg)
		success = false
	}

	if cfg.Port == "" {
		invalid("missing value for port")
	}

	for index, hubicApplication := range cfg.HubicApplications {
		if hubicApplication.ClientID == "" {
			invalid(fmt.Sprintf("hubic.[%d].client_id", index))
		}
		if hubicApplication.ClientSecret == "" {
			invalid(fmt.Sprintf("hubic.[%d].client_secret", index))
		}
		if hubicApplication.RedirectURL == "" {
			invalid(fmt.Sprintf("hubic.[%d].redirect_url", index))
		}
	}

	return success
}

//TODO Not Thread Safe
func (hubicApp HubicApplication) getToken() (*oauth2.Token, error) {
	tokenString := hubicApp.OAuthToken
	if tokenString == "" {
		return nil, fmt.Errorf("No OAuth Token stored. Please register...")
	}
	token := new(oauth2.Token)
	err := json.Unmarshal([]byte(tokenString), token)
	return token, err
}

//TODO Implement Save Token
