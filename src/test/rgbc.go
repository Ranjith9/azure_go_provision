package main

import (
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
//	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
//        "github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
        "github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
        "context"
        "log"
//        "github.com/marstr/randname"
//        "encoding/json"
//        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
//        "reflect"
)

var (
        clientID =        "1b390a83-5255-47e3-bcd3-a5e41061e661"
        subscriptionID =  "0594cd49-9185-425d-9fe2-8d051e4c6054"
        tenantID =        "85c997b9-f494-46b3-a11d-772983cf6f11"
        clientSecret =    "ceZBgXQoryOMGvK6txScc/TruRGaHucs9uayj8d/OtI="
        userID =          "M1038273@mindtree.com"
        password =        "Lovingme9"
)


func getServicePrincipalToken() (adal.OAuthTokenProvider, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
/*	code, err := adal.NewServicePrincipalToken(
		*oauthConfig,
		clientID,
		clientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)*/

         code, err := adal.NewServicePrincipalTokenFromUsernamePassword(
               *oauthConfig,
                clientID,
                userID,
                password,
                azure.PublicCloud.ResourceManagerEndpoint)

        if err != nil {
                log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
        }

       return code, err
}


func main(){
  ctx := context.Background()
  token, _ := getServicePrincipalToken()
  rbacClient := graphrbac.NewUsersClient(tenantID)
  rbacClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := rbacClient.Get(ctx,"054f4d4b-fb33-411f-ad15-cf2b81dc17b4")

/*  response, _ := rbacClient.Create(ctx,
                   graphrbac.ApplicationCreateParameters{
		   AvailableToOtherTenants: to.BoolPtr(false),
		   DisplayName:             to.StringPtr("trying"),
		   Homepage:                to.StringPtr("http://sample"),
		   IdentifierUris:          &[]string{randname.GenerateWithPrefix("http://sample", 10)},
	           },
                 )*/

  fmt.Printf("%+v",response)

//  json_val, _ := json.Marshal(response.Values())
//  fmt.Printf("%s\n", string(json_val))

}
