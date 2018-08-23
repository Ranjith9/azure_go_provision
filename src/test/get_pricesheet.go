package main

import (
        "github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2018-06-30/consumption"
        "context"
        "test/azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
        "encoding/json"
)


func main(){
  ctx := context.Background()
  token, _, subscriptionID := auth.GetServicePrincipalToken()
  vmClient := consumption.NewPriceSheetClient(subscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := vmClient.Get(ctx, "", "", nil)

  fmt.Printf("%+v",response)
//  fmt.Printf("\n")
  json_val, _ := json.MarshalIndent(response,"","  ")
  fmt.Printf("%s\n", string(json_val))

}
