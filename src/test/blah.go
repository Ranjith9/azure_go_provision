package main

import (
        "github.com/Azure/go-autorest/autorest/adal"
        "github.com/Azure/go-autorest/autorest/azure"
		"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
        "github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
        "github.com/Azure/go-autorest/autorest"
        "context"
        "log"
        "fmt"
)

var (
        clientID =        "1b390a83-5255-47e3-bcd3-a5e41061e661"
//        subscriptionID =  "0594cd49-9185-425d-9fe2-8d051e4c6054"
        tenantID =        "85c997b9-f494-46b3-a11d-772983cf6f11"
        clientSecret =    "ceZBgXQoryOMGvK6txScc/TruRGaHucs9uayj8d/OtI="
        subscriptionID string
)




func getServicePrincipalToken() (adal.OAuthTokenProvider, error) {
        oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
        code, err := adal.NewServicePrincipalToken(
                *oauthConfig,
                clientID,
                clientSecret,
                azure.PublicCloud.ResourceManagerEndpoint)
        if err != nil {
                log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
        }

       return code, err
}


func main(){
  ctx := context.Background()
  token, _ := getServicePrincipalToken()
  subsciptionClient := subscription.NewSubscriptionsClient()
  subsciptionClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := subsciptionClient.List(ctx)

//  fmt.Printf("%+v",response.Values())

    result := response.Values()

    for _, p := range result {
         fmt.Println(*p.SubscriptionID)
    }


//////////////////////////////////////////////////////////////////////////
  fmt.Println("\n Enter any value from above")
  fmt.Scanln(&subscriptionID)
  rgClient := resources.NewGroupsClient(subscriptionID)
  rgClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response1, _ := rgClient.List(ctx, "", nil)

  result1 := response1.Values()

      for _, p := range result1 {
         fmt.Println(*p.Name)
    }
}
