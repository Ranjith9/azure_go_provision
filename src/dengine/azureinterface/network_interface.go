package DengineAzureInterface

import (
        "github.com/Azure/go-autorest/autorest"
        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
        "context"
        "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
)

var (
       Token, _   = GetServicePrincipalToken()
       Ctx = context.Background()
)

// Author: Ranjith Janagama.



////////////////////////////// Creation of Vnet //////////////////////////////

func CreateVnet(resourceGroup, name, vpn_cidr, location string) network.VirtualNetwork {
  vnetClient := network.NewVirtualNetworksClient(SubscriptionID)
  vnetClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := vnetClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.VirtualNetwork{
    Location: to.StringPtr(location),
    VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
      AddressSpace: &network.AddressSpace{
        AddressPrefixes: &[]string{vpn_cidr},
      },
    },
  })
  err = response.WaitForCompletion(Ctx, vnetClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vnet create or update future response: %v", err)
  }

  result, _ := response.Result(vnetClient)
  return result
}

//////////////////////////// Creation of NSG ////////////////////////////////////

func CreateNsg(resourceGroup,name,location,cidr string) network.SecurityGroup {
  nsgClient := network.NewSecurityGroupsClient(SubscriptionID)
  nsgClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := nsgClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.SecurityGroup{
    Location: to.StringPtr(location),//location
    SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
      SecurityRules: &[]network.SecurityRule{
        {
          Name: to.StringPtr(name+"Rule"),
          SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
            Protocol:                 "Tcp",
            SourcePortRange:          to.StringPtr("*"),
            DestinationPortRange:     to.StringPtr("22"),//port
            SourceAddressPrefix:      to.StringPtr("*"),
            DestinationAddressPrefix: to.StringPtr(cidr),//sub_cidr
            Access:                   "Allow",
            Priority:                 to.Int32Ptr(100), //priority
            Direction:                "Inbound",
          },
        },
      },
    },
  })

  err = response.WaitForCompletion(Ctx, nsgClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nsg create or update future response: %v", err)
  }

  result, _ := response.Result(nsgClient)
  return result
}

func SecurityRuleCreate(resourceGroup,rule_name,nsgname,cidr,port string, priority int32) network.SecurityRule {
  fmt.Println("I'm in..")
  nsgruleClient := network.NewSecurityRulesClient(SubscriptionID)
  nsgruleClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, _ := nsgruleClient.CreateOrUpdate(Ctx, resourceGroup, nsgname, rule_name,
//  nsgruleClient.CreateOrUpdate(Ctx, resourceGroup, nsgname, rule_name,
  network.SecurityRule{
    SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
      Protocol:                 "Tcp",
      SourcePortRange:          to.StringPtr("*"),
      DestinationPortRange:     to.StringPtr(port),//port
      SourceAddressPrefix:      to.StringPtr("*"),
      DestinationAddressPrefix: to.StringPtr(cidr),//cidr
      Access:                   "Allow",
      Priority:                 to.Int32Ptr(priority), //priority
      Direction:                "Inbound",
    },
  })

/*  err = response.WaitForCompletion(Ctx, nsgruleClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nsg rule create or update future response: %v", err)
  }
*/

  result, _ := response.Result(nsgruleClient)
  return result
}

/*
func nsgRule(rule_name,cidr,port string, priority int32) network.SecurityRule {
  return network.SecurityRule{
    Name: to.StringPtr(rule_name),
    SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
      Protocol:                 "Tcp",
      SourcePortRange:          to.StringPtr("*"),
      DestinationPortRange:     to.StringPtr(port),//port
      SourceAddressPrefix:      to.StringPtr("*"),
      DestinationAddressPrefix: to.StringPtr(cidr),//cidr
      Access:                   "Allow",
      Priority:                 to.Int32Ptr(priority), //priority
      Direction:                "Inbound",
    },
  }
}
*/

//////////////////////////////// Creation of routes ////////////////////////////

func CreateRouteTb(resourceGroup,name,location,sub_cidr string) network.RouteTable {
  routetbClient := network.NewRouteTablesClient(SubscriptionID)
  routetbClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := routetbClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.RouteTable{
    Location: to.StringPtr(location),
    RouteTablePropertiesFormat: &network.RouteTablePropertiesFormat{
      Routes: &[]network.Route{
        {
          Name: to.StringPtr("route1"),
          RoutePropertiesFormat: &network.RoutePropertiesFormat{
            AddressPrefix: to.StringPtr(sub_cidr),
            NextHopType:   "VirtualNetworkGateway",
          },
        },
      },
    },
  })
  err = response.WaitForCompletion(Ctx, routetbClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the route table create or update future response: %v", err)
  }

  result, _ := response.Result(routetbClient)
  return result

}

/////////////////////////////////// Creation of Subnet ///////////////////////////////

func CreateSubnet(resourceGroup,name,sub_cidr string, vpn,nsg,routetable *string) (*string,*string){
  subnetClient := network.NewSubnetsClient(SubscriptionID)
  subnetClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := subnetClient.CreateOrUpdate(Ctx, resourceGroup, *vpn, name,
  network.Subnet{
    SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
      AddressPrefix: to.StringPtr(sub_cidr),
      NetworkSecurityGroup: &network.SecurityGroup{
        ID: nsg,// to.StringPtr(nsg)
      },
      RouteTable: &network.RouteTable{
        ID: routetable, // to.StringPtr(routetable),
      },
    },
  })
  err = response.WaitForCompletion(Ctx, subnetClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the subnet create or update future response: %v", err)
  }
  result, _ := response.Result(subnetClient)
  return result.Name, result.ID
}

///////////////////////////////// Creation of Public IP /////////////////////////////////

func CreatePublicIp(resourceGroup,name,location string) *string {
  pubIpClient := network.NewPublicIPAddressesClient(SubscriptionID)
  pubIpClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  publicip, err := pubIpClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.PublicIPAddress{
    Location: to.StringPtr(location),
    PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
      PublicIPAllocationMethod: "Dynamic",
      DNSSettings: &network.PublicIPAddressDNSSettings{
        DomainNameLabel: to.StringPtr(name),
      },
    },
  })
  err = publicip.WaitForCompletion(Ctx, pubIpClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the ip create or update future response: %v", err)
  }

  ip, _ := publicip.Result(pubIpClient)
  return ip.ID
}

//////////////////////////////// Creation of NIC //////////////////////////////////////////

func CreateNic(resourceGroup,name,location,subnetID,nsgID,lbname string, ipID *string) *string {
  nicClient := network.NewInterfacesClient(SubscriptionID)
  nicClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  nic, err := nicClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.Interface{
    Location: to.StringPtr(location),
    InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
      IPConfigurations: &[]network.InterfaceIPConfiguration{
        {
          Name: to.StringPtr(name+"-ipconfig"),
          InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
            PrivateIPAllocationMethod: "Dynamic",
            PublicIPAddress: &network.PublicIPAddress{
//              ID: ip.ID,
              ID: ipID,
            },
            Subnet: &network.Subnet{
              ID: to.StringPtr(subnetID),
//              ID: subnetID,
            },
            if lbname != nil {

              LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
                {
                  ID: to.StringPtr(CreateID(resourceGroup, lbname, "backendAddressPools", "backend-pool")),
                },
              },
              LoadBalancerInboundNatRules: &[]network.InboundNatRule{
                  NatRule(resourceGroup,"nat1",lbname,"lb-fip", 1122),
              },
            }
          },
        },
      },
      NetworkSecurityGroup: &network.SecurityGroup{
        ID: to.StringPtr(nsgID),
//          ID: nsgID,
      },
    },
  })
  err = nic.WaitForCompletion(Ctx, nicClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nic create or update future response: %v", err)
  }

  nicard, _ := nic.Result(nicClient)
  return nicard.ID
}
