package main

import (
    "azure/interface"
    "fmt"
)

 var (
    resourceGroup = "Dengine"
    location      = "CentralIndia"
)

func main() {
  fmt.Println("Creating LB\n")

  fmt.Println("Creating IP for LB")
  ip := auth.CreatePublicIp(resourceGroup,"test-lb-ip",location)
  fmt.Println("Created IP for LB")
  lb := auth.CreateLB(resourceGroup,"test-lb",location,ip)
  fmt.Println(*lb)
}
