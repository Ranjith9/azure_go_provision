package main

import (
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
        "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
        "context"
        "test/azure/client"
        "fmt"
        "github.com/Azure/go-autorest/autorest"
        "strings"
//        "encoding/json"
)


func main(){
  ctx := context.Background()
  token, _, subscriptionID := auth.GetServicePrincipalToken()
  vmClient := compute.NewVirtualMachinesClient(subscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

//   Getting VM components

  vm, _  := vmClient.Get(ctx, "M1038273", "test-vm", "")

  fmt.Println(*vm.ID)
  net := vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces

//   Getting NIC associated with the VM
  var nic string
  for _, nictemp := range *net {
//      fmt.Println(*nictemp.ID)
      nic = *nictemp.ID
  }

  nic_slice := strings.Split(nic, "/")
//  fmt.Println(nic_slice[len(nic_slice) - 1])
  nic_card := nic_slice[len(nic_slice) - 1]

//  Getting Ip associated with the VM
  nicClient := network.NewInterfacesClient(subscriptionID)
  nicClient.Authorizer = autorest.NewBearerAuthorizer(token)

  net_nic, _ := nicClient.Get(ctx , "M1038273", nic_card, "")

  fmt.Printf(*net_nic.Name)

  var ip string
  ip_val := net_nic.InterfacePropertiesFormat.IPConfigurations
  for _, iptemp := range *ip_val {
//      fmt.Println(*iptemp.ID)
      ip = *iptemp.ID
  }

  nsg := net_nic.InterfacePropertiesFormat.NetworkSecurityGroup

//  fmt.Println("\n",ip,*nsg.ID)
  ipconfig_slice := strings.Split(ip, "/")
  ip_name := ipconfig_slice[len(ipconfig_slice) - 1]

  nsg_slice := strings.Split(*nsg.ID, "/")
  nsg_name := nsg_slice[len(nsg_slice) - 1]

  fmt.Println("\n",ip_name,nsg_name)

  ipconfigClient := network.NewInterfaceIPConfigurationsClient(subscriptionID)
  ipconfigClient.Authorizer = autorest.NewBearerAuthorizer(token)

  pubip,_ := ipconfigClient.Get(ctx, "M1038273", nic_card, ip_name)


  true_ip := pubip.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress

  fmt.Println(*true_ip.ID)
}
