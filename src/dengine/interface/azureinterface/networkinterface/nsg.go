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

func getNsgClient() network.SecurityGroupsClient {
	nsgClient := network.NewSecurityGroupsClient(subscription)
	nsgClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return nsgClient
}

type NsgIn struct {
        ResourceGroup string
        NsgName string      `json:"nsgname,omitempty"`
        Location string      `json:"location,omitempty"`
}

// CreateNetworkSecurityGroup creates a new network security group.

func (ns NsgIn) CreateNetworkSecurityGroup() (nsg network.SecurityGroup, err error) {
	nsgClient := getNsgClient()
	future, err := nsgClient.CreateOrUpdate(
		ctx,
		ns.ResourceGroup,
		ns.NsgName,
		network.SecurityGroup{
			Location: to.StringPtr(ns.Location),
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

func (ns NsgIn) DeleteNetworkSecurityGroup() (ar autorest.Response, err error) {
        nsgClient := getNsgClient()
        future, err := nsgClient.Delete(
                ctx,
                ns.ResourceGroup,
                ns.NsgName,
                )

        if err != nil {
                return ar, fmt.Errorf("cannot delete nsg: %v", err)
        }

        err = future.WaitForCompletion(ctx, nsgClient.Client)
        if err != nil {
                return ar, fmt.Errorf("cannot get nsg delete future response: %v", err)
        }

        return  future.Result(nsgClient)
}

func (ns NsgIn) GetNetworkSecurityGroup() (nsg network.SecurityGroup, err error) {
        nsgClient := getNsgClient()
        future, err := nsgClient.Get(
                ctx,
                ns.ResourceGroup,
                ns.NsgName,
                "")

        if err != nil {
                return nsg, fmt.Errorf("cannot list nsg: %v", err)
        }

        return  future, err
}

func (ns NsgIn) ListNetworkSecurityGroup() (nsg []network.SecurityGroup, err error) {
        nsgClient := getNsgClient()
        future, err := nsgClient.List(
                ctx,
                ns.ResourceGroup,
                )

        if err != nil {
                return nsg, fmt.Errorf("cannot list nsg: %v", err)
        }

        return  future.Values(), err
}

func ListAllNetworkSecurityGroup() (nsg []network.SecurityGroup, err error) {
        nsgClient := getNsgClient()
        future, err := nsgClient.List(
                ctx,
                )

        if err != nil {
                return nsg, fmt.Errorf("cannot list nsg: %v", err)
        }

        return  future.Values(), err
}
