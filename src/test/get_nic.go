package main

import (
        "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
        "context"
        "test/azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
        "encoding/json"
)


func main(){
  ctx := context.Background()
  token, _, subscriptionID := auth.GetServicePrincipalToken()
  vmClient := network.NewInterfacesClient(subscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := vmClient.Get(ctx, "M1038273", "db39", "")

//  fmt.Printf("%+v",response)
//  fmt.Printf("\n")
  json_val, _ := json.Marshal(response)
  fmt.Printf("%s\n", string(json_val))

}
