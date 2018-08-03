package azurenetwork

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
//        "dengine/interface/azureinterface/getclient"
)

//var (
//        token, subscription = getclient.GetToken()
//)

func getNicClient() network.InterfacesClient {
	nicClient := network.NewInterfacesClient(subscription)
	nicClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return nicClient
}

// CreateNIC creates a new network interface.

func CreateNIC(ctx context.Context, resourceGroup string, subnetID string, nsgID string, ipID, nicName string, location string) (nic network.Interface, err error) {

	nicClient := getNicClient()
	future, err := nicClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		nicName,
		network.Interface{
			Name:     to.StringPtr(nicName),
			Location: to.StringPtr(location),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: to.StringPtr(nicName + "-ipConfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Subnet: &network.Subnet{
								ID: to.StringPtr(subnetID),
							},
							PrivateIPAllocationMethod: network.Dynamic,
							PublicIPAddress: &network.PublicIPAddress{
								ID: to.StringPtr(ipID),
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

	if err != nil {
		return nic, fmt.Errorf("cannot create nic: %v", err)
	}

	err = future.WaitForCompletion(ctx, nicClient.Client)
	if err != nil {
		return nic, fmt.Errorf("cannot get nic create or update future response: %v", err)
	}

	return future.Result(nicClient)
}
