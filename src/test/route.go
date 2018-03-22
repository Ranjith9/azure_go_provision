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
  routetbClient := network.NewRouteTablesClient(subscriptionID)
  routetbClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := routetbClient.CreateOrUpdate(ctx, "Dengine", "routable1",
  network.RouteTable{
    Location: to.StringPtr("CentralIndia"),
    RouteTablePropertiesFormat: &network.RouteTablePropertiesFormat{
      Routes: &[]network.Route{
        {
          Name: to.StringPtr("route1"),
          RoutePropertiesFormat: &network.RoutePropertiesFormat{
            AddressPrefix: to.StringPtr("10.0.1.0/24"),
            NextHopType: "VirtualNetworkGateway",
          },
        },
      },
    },
  })
  result, _ := response.Result(routetbClient)
  fmt.Println(*result.ID)

}
