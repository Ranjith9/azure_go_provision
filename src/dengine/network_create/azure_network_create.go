package NetworkCreate

import (
  "fmt"
 // "reflect"
//  "encoding/json"
  "strings"
  "dengine/azureinterface"
  "strconv"
//  "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
)

var (

  Network_Response NetworkResponse

)

type NetworkCreateInput struct {
  Name    string
  VpnCidr string
  SubCidr []string
  Type    string
  Ports   []string
  Cloud   string
  Location  string
  ResourceGroup string
}

type VpcReponse struct {
  Name string
  Id   string
  Type string
  Nsg  string
}

type SubnetResponse struct {
  Name string
  Id   string
//    Subnet []network.Subnet
}

type NetworkResponse struct {
  VpnReponse    VpcReponse
  SubnetResponse []SubnetResponse
}

// being create_network my job is to create network and give back the response who called me
func CreateNetwork(net NetworkCreateInput) NetworkResponse {

  switch strings.ToLower(net.Cloud) {
  case "aws" :
  case "azure" :
      fmt.Println("creating vpn")
      vpn := DengineAzureInterface.CreateVnet(net.ResourceGroup, net.Name, net.VpnCidr, net.Location)
      nsg := DengineAzureInterface.CreateNsg(net.ResourceGroup, *vpn.Name+"_nsg", net.Location, net.VpnCidr)

      vpn_response := VpcReponse{*vpn.Name, *vpn.ID , net.Type, *nsg.ID}

      fmt.Println("creating subnets")
      response_subnet := []SubnetResponse{}
      for i, sub := range net.SubCidr {
          route := DengineAzureInterface.CreateRouteTb(net.ResourceGroup, net.Name+ "_sub_" + strconv.Itoa(i)+"_route", net.Location, sub)
          subname, subid := DengineAzureInterface.CreateSubnet(net.ResourceGroup, net.Name+ "_sub_" + strconv.Itoa(i), sub, vpn.Name, nsg.ID, route.ID)
//      subnet = SubnetResponse{subname,subid}
      response_subnet = append(response_subnet, SubnetResponse{*subname,*subid})
      }
//      network_response := NetworkResponse{vpn_response,response_subnet}
      network_repsonseptr := &Network_Response
     *network_repsonseptr = NetworkResponse{vpn_response,response_subnet}
  case "gcp"   :
  case "openstack" :

  }
  return Network_Response
}
