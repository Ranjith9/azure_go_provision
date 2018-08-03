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

func getNsgClient() network.SecurityGroupsClient {
	nsgClient := network.NewSecurityGroupsClient(subscription)
	nsgClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return nsgClient
}

// CreateNetworkSecurityGroup creates a new network security group.

func CreateNetworkSecurityGroup(ctx context.Context, resourceGroup string, nsgName string, location string) (nsg network.SecurityGroup, err error) {
	nsgClient := getNsgClient()
	future, err := nsgClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		nsgName,
		network.SecurityGroup{
			Location: to.StringPtr(location),
		},
	)

	if err != nil {
		return nsg, fmt.Errorf("cannot create nsg: %v", err)
	}

	err = future.WaitForCompletion(ctx, nsgClient.Client)
	if err != nil {
		return nsg, fmt.Errorf("cannot get nsg create or update future response: %v", err)
	}

	return future.Result(nsgClient)
}
