package auth

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

func CreateVnet(resourceGroup, name, vpn_cidr, location string) *string {
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
//  if err != nil {
//    return vnet, fmt.Errorf("cannot create virtual network: %v", err)
//  }
  err = response.WaitForCompletion(Ctx, vnetClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vnet create or update future response: %v", err)
  }

  result, _ := response.Result(vnetClient)
  return result.Name
}

func CreateNsg(resourceGroup,name,location,sub_cidr,port string, priority int32) *string {
//  ctx = context.Background()
  nsgClient := network.NewSecurityGroupsClient(SubscriptionID)
  nsgClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := nsgClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.SecurityGroup{
    Location: to.StringPtr(location),//location
    SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
      SecurityRules: &[]network.SecurityRule{
        {
          Name: to.StringPtr("Rule1"),
          SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
            Protocol:                 "Tcp",
            SourcePortRange:          to.StringPtr("*"),
            DestinationPortRange:     to.StringPtr(port),//port
            SourceAddressPrefix:      to.StringPtr("*"),
            DestinationAddressPrefix: to.StringPtr(sub_cidr),//sub_cidr
            Access:                   "Allow",
            Priority:                 to.Int32Ptr(priority), //priority
            Direction:                "Inbound",
          },
        },
      },
    },
  })
//  if err != nil {
//    return vnet, fmt.Errorf("cannot create network security group: %v", err)
//  }
  err = response.WaitForCompletion(Ctx, nsgClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nsg create or update future response: %v", err)
  }

  result, _ := response.Result(nsgClient)
  return result.ID
}

func CreateRouteTb(resourceGroup,name,location,sub_cidr string) *string {
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
            NextHopType: "VirtualNetworkGateway",
          },
        },
      },
    },
  })
//  if err != nil {
//    return vnet, fmt.Errorf("cannot create route table: %v", err)
//  }
  err = response.WaitForCompletion(Ctx, routetbClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the route table create or update future response: %v", err)
  }

  result, _ := response.Result(routetbClient)
  return result.ID

}

func CreateSubnet(resourceGroup,name,sub_cidr string, vpn,nsg,routetable *string) *string {
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
  return result.ID
}

Errorf("cannot get the subnet create or update future response: %v", err)
  }
  result, _ := response.Result(subnetClient)
  return result
}

///////////////////////////////// Creation of Public IP /////////////////////////////////

func CreatePublicIp(resourceGroup,name,location string) *string {
  pubIpClient := network.NewPublicIPAddressesClient(SubscriptionID)
  pubIpClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  publicip, err := pubIpClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.PublicIPAddress{
    Location: to.StringPtr(location),
    PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
//      PublicIPAddressVersion:   "IPv4",
      PublicIPAllocationMethod: "Dynamic",
//      IdleTimeoutInMinutes: to.Int32Ptr(4),
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

//  nic, err := nicClient.CreateOrUpdate(Ctx, resourceGroup, name,

  if len(lbname) != 0 {
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
                ID: ipID,
              },
              Subnet: &network.Subnet{
                ID: to.StringPtr(subnetID),
              },
              LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
                {
                  ID: to.StringPtr(CreateID(resourceGroup, lbname, "backendAddressPools", "backend-pool")),
                },
              },
              LoadBalancerInboundNatRules: &[]network.InboundNatRule{
                  NatRule(resourceGroup,"nat1",lbname,"lb-fip", 1122),
              },
            },
          },
        },
        NetworkSecurityGroup: &network.SecurityGroup{
          ID: to.StringPtr(nsgID),
        },
      },
    })
  err = nic.WaitForCompletion(Ctx, nicClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nic create or update future response: %v", err)
  }

  nicard, _ := nic.Result(nicClient)
  return nicard.ID
  } else {
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
                ID: ipID,
              },
              Subnet: &network.Subnet{
                ID: to.StringPtr(subnetID),
              },
            },
          },
        },
        NetworkSecurityGroup: &network.SecurityGroup{
          ID: to.StringPtr(nsgID),
        },
      },
    },
  )
  err = nic.WaitForCompletion(Ctx, nicClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the nic create or update future response: %v", err)
  }

  nicard, _ := nic.Result(nicClient)
  return nicard.ID
  }
}
