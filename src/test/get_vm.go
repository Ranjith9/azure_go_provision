package main

import (
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
        "context"
        "test/azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
        "encoding/json"
)


func main(){
  ctx := context.Background()
  token, _, subscriptionID := auth.GetServicePrincipalToken()
  vmClient := compute.NewVirtualMachinesClient(subscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := vmClient.InstanceView(ctx, "M1038273", "ubuntu")

  fmt.Printf("%+v",response)
//  fmt.Printf("\n")
  json_val, _ := json.MarshalIndent(response,"","  ")
  fmt.Printf("%s\n", string(json_val))

}
