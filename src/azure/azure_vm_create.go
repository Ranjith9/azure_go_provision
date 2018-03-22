package main

import (
    "azure/interface"
    "fmt"
)

func main() {
  fmt.Println("Creating VM")
  auth.CreateVm()
}
lIndia"
    password      = "ubuntu@12345"
)

func main() {
  fmt.Println("Creating VM\n")
  fmt.Println("Creating IP for VM")
  ip := auth.CreatePublicIp(resourceGroup,"zander-ip",location)
  fmt.Println("Creating NIC for VM")
  nic := auth.CreateNic(resourceGroup,"zander-nic",location, "/subscriptions/0594cd49-9185-425d-9fe2-8d051e4c6054/resourceGroups/Dengine/providers/Microsoft.Network/virtualNetworks/go_vpn/subnets/go_sub1", "/subscriptions/0594cd49-9185-425d-9fe2-8d051e4c6054/resourceGroups/Dengine/providers/Microsoft.Network/networkSecurityGroups/go_nsg", "test-lb", ip)
  vm := auth.CreateVm(resourceGroup,"zander-vm",location,password,nic)
  fmt.Println(*vm)
}
