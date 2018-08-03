package azurenetwork

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
//        "dengine/interface/azureinterface/getclient"
        "dengine/access/azureaccess"
//        "github.com/Azure/go-autorest/autorest/adal"
)

var (
        token, _, subscription = auth.GetServicePrincipalToken()
)

func getVnetClient() network.VirtualNetworksClient {
	vnetClient := network.NewVirtualNetworksClient(subscription)
	vnetClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return vnetClient
}

// CreateVirtualNetwork creates a virtual network
func CreateVirtualNetwork(ctx context.Context, resourceGroup string, vnetName string, cidr string, location string) (vnet network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		vnetName,
		network.VirtualNetwork{
			Location: to.StringPtr(location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{cidr},
				},
			},
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletion(ctx, vnetClient.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(vnetClient)
}
