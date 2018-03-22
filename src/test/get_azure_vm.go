package main

import (
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
        "context"
        "azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
)

var (
        subscriptionID =  "0594cd49-9185-425d-9fe2-8d051e4c6054"
)

func main(){
  ctx := context.Background()
  token, _ := auth.GetServicePrincipalToken()
  vmClient := compute.NewVirtualMachinesClient(subscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := vmClient.Get(ctx, "Dengine", "go", "instanceView")

//  result, _ := response.Result(vmClient)
  result := response.VirtualMachineProperties.NetworkProfile.NetworkInterfaces
  fmt.Printf("%s",result)
  fmt.Printf("\n")
}
