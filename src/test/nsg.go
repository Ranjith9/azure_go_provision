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
//        "reflect"
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
  in := int32(101)
  ctx := context.Background()
  token, _ := getServicePrincipalToken()
  nsgClient := network.NewSecurityGroupsClient(subscriptionID)
  nsgClient.Authorizer = autorest.NewBearerAuthorizer(token)

  response, _ := nsgClient.CreateOrUpdate(ctx, "Dengine", "test",
  network.SecurityGroup{
    Location: to.StringPtr("CentralIndia"),
    SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
      SecurityRules: &[]network.SecurityRule{
        {
          Name: to.StringPtr("Rule1"),
          SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
            Protocol:                 "Tcp",
            SourcePortRange:          to.StringPtr("*"),
            DestinationPortRange:     to.StringPtr("22"),
            SourceAddressPrefix:      to.StringPtr("*"),
            DestinationAddressPrefix: to.StringPtr("10.0.1.0/24"),
            Access:                   "Allow",
            Priority:                 &in,  // to.Int32Ptr(100)
            Direction:                "Inbound",
          },
        },
      },
    },
  })
  result, _ := response.Result(nsgClient)
  fmt.Println(result)

//  result, _ := nsgClient.Get(ctx, "Dengine", "test", "")
//  fmt.Println(*result.ID)
}
