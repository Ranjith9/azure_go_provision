package main

import (
        "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
        "context"
        "test/azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
        "encoding/json"
)


func main(){
  ctx := context.Background()
  token, _, subscriptionID := auth.GetServicePrincipalToken()
  rgClient := resources.NewGroupsClient(subscriptionID)
  rgClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := rgClient.List(ctx, "", nil)

//  value := response.
//  fmt.Printf("%+v",response)
//  fmt.Println(response.Response)

  json_val, _ := json.Marshal(response.Values())
  fmt.Printf("%s\n", string(json_val))

}
