package main

import (
    "azure/interface"
    "fmt"
)

var (
    resourceGroup = "test"
    location      = "CentralIndia"
)

func main() {
  fmt.Println("Creating VPN")
  vpn := auth.CreateVnet(resourceGroup,"go_vpn","192.168.0.0/16",location)
  fmt.Println("VPN =",*vpn)
  fmt.Println("Creating NSG")
  nsg := auth.CreateNsg(resourceGroup,"go_nsg",location,"192.168.10.0/24","22",101)
  fmt.Println("NSG =",*nsg)
  fmt.Println("Creating ROUTE")
  route := auth.CreateRouteTb(resourceGroup,"go_route",location,"192.168.10.0/24")
  fmt.Println("ROUTE =",*route)
  fmt.Println("Creating SUBNET")
  subnet := auth.CreateSubnet(resourceGroup, "go_sub1","192.168.10.0/24", vpn, nsg, route)
  fmt.Println("Subnet", *subnet)

}
("Subnet =", *subnet.ID)

}
