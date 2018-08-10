package azurenetwork

import (
//	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
//        "dengine/interface/azureinterface/getclient"
)

//var (
//        token, subscription = getclient.GetToken()
//)

func getSubnetsClient() network.SubnetsClient {
	subnetsClient := network.NewSubnetsClient(subscription)
	subnetsClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return subnetsClient
}

type SubnetIn struct {
        ResourceGroup string
        VnetName string      `json:"vnetname,omitempty"`
        SubnetName string    `json:"subnetname,omitempty"`
        SubnetCidr string    `json:"cidr,omitempty"`
}

// CreateVirtualNetworkSubnet creates a subnet in an existing vnet

func (sub SubnetIn) CreateVirtualNetworkSubnet() (subnet network.Subnet, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		sub.ResourceGroup,
		sub.VnetName,
		sub.SubnetName,
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr(sub.SubnetCidr),
			},
		})
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletion(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

func (sub SubnetIn) DeleteVirtualNetworkSubnet() (ar autorest.Response, err error) {
        subnetsClient := getSubnetsClient()

        future, err := subnetsClient.Delete(
                ctx,
                sub.ResourceGroup,
                sub.VnetName,
                sub.SubnetName,
                )

        err = future.WaitForCompletion(ctx, subnetsClient.Client)
        if err != nil {
                return ar, fmt.Errorf("cannot get the subnet delete future response: %v", err)
        }

        return future.Result(subnetsClient)
}

func (sub SubnetIn) GetVirtualNetworkSubnet() (subnet network.Subnet, err error) {
        subnetsClient := getSubnetsClient()

        future, err := subnetsClient.Get(
                ctx,
                sub.ResourceGroup,
                sub.VnetName,
                sub.SubnetName,
                "")
        if err != nil {
                return subnet, fmt.Errorf("cannot get the subnet: %v", err)
        }

        return future, err
}

func (sub SubnetIn) ListVirtualNetworkSubnet() (subnet []network.Subnet, err error) {
        subnetsClient := getSubnetsClient()

        future, err := subnetsClient.List(
                ctx,
                sub.ResourceGroup,
                sub.VnetName,
                )

        if err != nil {
                return subnet, fmt.Errorf("cannot list dubnetwork: %v", err)
        }

        return future.Values(), err
}
