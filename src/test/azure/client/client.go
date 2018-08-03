package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"io/ioutil"
	"log"
	"os/user"
)

type Credentials struct {
	Profile        string `json:"profile,omitempty"`
	ClientID       string `json:"clientID,omitempty"`
	SubscriptionID string `json:"subscriptionID,omitempty"`
	TenantID       string `json:"tenantID,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
}

var result Credentials

func init() {

	user, _ := user.Current()
	file := user.HomeDir + "/.azure/credentials"

	plan, _ := ioutil.ReadFile(file) // filename is the JSON file to read
	var data []Credentials
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Errorf("Cannot unmarshal the json ", err)
	}

	for _, t := range data {
		if t.Profile == "sumanth" {
			result = t
			break
		} else if t.Profile != "sumanth" {
			continue
		}
	}
	if (Credentials{}) == result {
		fmt.Println("I don't know the user")
	}

}

func GetServicePrincipalToken() (adal.OAuthTokenProvider, error, string) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, result.TenantID)
	code, err := adal.NewServicePrincipalToken(
		*oauthConfig,
		result.ClientID,
		result.ClientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
	}

	return code, err, result.SubscriptionID
}
