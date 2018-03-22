package auth

import (
        "github.com/Azure/go-autorest/autorest"
        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
//        "context"
        "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
)

/*

Author: Ranjith Janagama.

To create a loadbalancer you need to create an IP and you need to pass the IP ID to this CreateLB() function

*/

////////////////////////////// Creation of LB //////////////////////////////


func CreateLB(resourceGroup, name, location string, ipID *string) *string {
  fmt.Println("LB INVOKED")
  lbClient := network.NewLoadBalancersClient(SubscriptionID)
  lbClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := lbClient.CreateOrUpdate(Ctx, resourceGroup, name,
  network.LoadBalancer{
    Sku: &network.LoadBalancerSku{
      Name: "Basic",
    },
    Location: to.StringPtr(location),
    LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
      FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
        {
          Name: to.StringPtr("lb-fip"),
          FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
            PrivateIPAllocationMethod: "Dynamic",
            PublicIPAddress: &network.PublicIPAddress{
              ID: ipID,  // to.StringPtr(ipID),
            },
          },
        },
      },
      BackendAddressPools: &[]network.BackendAddressPool{
        {
          Name: to.StringPtr("backend-pool"),
        },
      },
      Probes: &[]network.Probe{
        {
          Name: to.StringPtr("probe1"),
          ProbePropertiesFormat: &network.ProbePropertiesFormat{
            Protocol:            "Http",
            Port:                to.Int32Ptr(80),
            IntervalInSeconds:   to.Int32Ptr(5),
            NumberOfProbes:      to.Int32Ptr(2),
            RequestPath:         to.StringPtr("/index.html"),
          },
        },
      },
      LoadBalancingRules: &[]network.LoadBalancingRule{
        {
          Name:to.StringPtr("lb-rule1"),
          LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
            Protocol:                "Tcp",
            FrontendPort:            to.Int32Ptr(80),
            BackendPort:             to.Int32Ptr(80),
            IdleTimeoutInMinutes:    to.Int32Ptr(4),
            EnableFloatingIP:        to.BoolPtr(false),
            FrontendIPConfiguration: &network.SubResource {
              ID: to.StringPtr(CreateID(resourceGroup, name, "frontendIPConfigurations", "lb-fip")),
            },
            BackendAddressPool:      &network.SubResource {
              ID: to.StringPtr(CreateID(resourceGroup, name, "backendAddressPools", "backend-pool")),
            },
            Probe:                   &network.SubResource {
              ID: to.StringPtr(CreateID(resourceGroup, name, "probes", "probe1")),
            },
          },
        },
      },
      InboundNatRules: &[]network.InboundNatRule{
        NatRule(resourceGroup,"nat1",name,"lb-fip", 1122),
        NatRule(resourceGroup,"nat2",name,"lb-fip", 1123),
      },
    },
  })
  err = response.WaitForCompletion(Ctx, lbClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vnet create or update future response: %v", err)
  }


  result, _ := response.Result(lbClient)
  return result.ID
}

func CreateID(resourceGroup, lbName, lbSub,lbSubName string) string {
  return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/%s/%s",SubscriptionID, resourceGroup, lbName, lbSub, lbSubName)
}

func NatRule(resourceGroup,rule_name,name,fip_name string, priority int32) network.InboundNatRule {
  return network.InboundNatRule{
    Name: to.StringPtr(rule_name),
    InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
      FrontendIPConfiguration: &network.SubResource {
        ID: to.StringPtr(CreateID(resourceGroup, name, "frontendIPConfigurations", "lb-fip")),
      },
      Protocol:                "Tcp",
      FrontendPort:            to.Int32Ptr(priority),
      BackendPort:             to.Int32Ptr(priority),
    },
  }
}
