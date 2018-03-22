package main

import (
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
        "context"
        "log"
        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
//        "os"
)

var (
        clientID =        "1b390a83-5255-47e3-bcd3-a5e41061e661"
        subscriptionID =  "0594cd49-9185-425d-9fe2-8d051e4c6054"
        tenantID =        "85c997b9-f494-46b3-a11d-772983cf6f11"
        clientSecret =    "ceZBgXQoryOMGvK6txScc/TruRGaHucs9uayj8d/OtI="
)


func getServicePrincipalToken() (adal.OAuthTokenProvider, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
	code, err := adal.NewServicePrincipalToken(
		*oauthConfig,
		clientID,
		clientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
        if err != nil {
                log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
        }

       return code, err
}


func main(){
  ctx := context.Background()
  token, _ := getServicePrincipalToken()
  vnetClient := network.NewVirtualNetworksClient(subscriptionID)
  vnetClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := vnetClient.CreateOrUpdate(ctx, "Dengine", "vpn",
  network.VirtualNetwork{
    Location: to.StringPtr("CentralIndia"),
    VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
      AddressSpace: &network.AddressSpace{
        AddressPrefixes: &[]string{"10.0.0.0/8"},
        SecurityGroup: &[]network.SecurityGroups{
	  Name: to.StringPtr("mine"),
	  Location: to.StringPtr("CentralIndia"),
	},
      },
      Subnets: &[]network.Subnet{
        {
          Name: to.StringPtr("subnet1"),
          SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
            AddressPrefix: to.StringPtr("10.0.0.0/16"),
          },
        },
        {
          Name: to.StringPtr("subnet2"),
          SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
            AddressPrefix: to.StringPtr("10.1.0.0/16"),
          },
        },
      },
    },
  })
  fmt.Println(response)
/*
  get, _ := vnetClient.Get(ctx, "Dengine", "vpn", "expand")
  fmt.Println(get)


  result, _ := vnetClient.Delete(ctx, "Dengine", "vpn")
  fmt.Println(result.Response)
*/
}

