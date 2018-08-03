package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/go-autorest/autorest"
	"strings"
	"test/azure/client"
)

func main() {
	ctx := context.Background()
	token, _, subscriptionID := auth.GetServicePrincipalToken()
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)
	vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

	vm, _ := vmClient.Get(ctx, "M1038273", "test-vm", "")

	osDisk := vm.VirtualMachineProperties.StorageProfile.OsDisk

	var nic string
	net := vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces
	for _, nictemp := range *net {
		nic = *nictemp.ID
	}
	nic_slice := strings.Split(nic, "/")
	nic_name := nic_slice[len(nic_slice)-1]

	nsg, ip := get_nsg("M1038273", nic_name)
	fmt.Println(nsg, ip, nic_name, *vm.Name,*osDisk.Name)

        delete_vm_resources("M1038273","test-vm",*osDisk.Name,nic_name,nsg,ip)
}

func get_nsg(resourcegroup string, nic_name string) (string,string){
        ctx := context.Background()
        token, _, subscriptionID := auth.GetServicePrincipalToken()

	nicClient := network.NewInterfacesClient(subscriptionID)
	nicClient.Authorizer = autorest.NewBearerAuthorizer(token)

	net_nic, _ := nicClient.Get(ctx, resourcegroup, nic_name, "")

	var config string
	ipconfig_val := net_nic.InterfacePropertiesFormat.IPConfigurations

	for _, iptemp := range *ipconfig_val {
		config = *iptemp.ID
	}
	ipconfig_slice := strings.Split(config, "/")
	ipconfig := ipconfig_slice[len(ipconfig_slice)-1]

	nsg_val := net_nic.InterfacePropertiesFormat.NetworkSecurityGroup

	nsg_slice := strings.Split(*nsg_val.ID, "/")

	ip := get_ip(resourcegroup, nic_name, ipconfig)

	return nsg_slice[len(nsg_slice)-1], ip
}

func get_ip(resourcegroup string, nic_name string, ipconfig_name string) string {
        ctx := context.Background()
        token, _, subscriptionID := auth.GetServicePrincipalToken()

	ipconfigClient := network.NewInterfaceIPConfigurationsClient(subscriptionID)
	ipconfigClient.Authorizer = autorest.NewBearerAuthorizer(token)

	pubip, _ := ipconfigClient.Get(ctx, resourcegroup, nic_name, ipconfig_name)

	ip := pubip.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress

	ip_slice := strings.Split(*ip.ID, "/")
	return ip_slice[len(ip_slice)-1]
}

func delete_vm_resources(resourcegroup string, vm_name string,osDisk_Name string, nic_name string ,nsg string ,ip string) {

        fmt.Println("Deleting all VM components")
        ctx := context.Background()
        token, _, subscriptionID := auth.GetServicePrincipalToken()
        vmClient := compute.NewVirtualMachinesClient(subscriptionID)
        vmClient.Authorizer = autorest.NewBearerAuthorizer(token)

        vm_delete, err := vmClient.Delete(ctx,resourcegroup,vm_name)

        if err != nil {
		fmt.Errorf("cannot delete vm: %v", err)
	}

	err = vm_delete.WaitForCompletion(ctx, vmClient.Client)
	if err != nil {
		fmt.Errorf("cannot get the vm delete future response: %v", err)
	}
        fmt.Println("Deleted VM")

//      ************************************@@@@@*********************************************

        nicClient := network.NewInterfacesClient(subscriptionID)
        nicClient.Authorizer = autorest.NewBearerAuthorizer(token)

        nic_delete, err := nicClient.Delete(ctx,resourcegroup,nic_name)
        if err != nil {
		fmt.Errorf("cannot delete nic: %v", err)
	}

	err = nic_delete.WaitForCompletion(ctx, nicClient.Client)
	if err != nil {
		fmt.Errorf("cannot get nic delete future response: %v", err)
	}
        fmt.Println("Deleted NIC")
//      ************************************@@@@@*********************************************

        nsgClient := network.NewSecurityGroupsClient(subscriptionID)
        nsgClient.Authorizer = autorest.NewBearerAuthorizer(token)

        nsg_delete, err := nsgClient.Delete(ctx, resourcegroup, nsg)
        if err != nil {
		fmt.Errorf("cannot delete nsg: %v", err)
	}

	err = nsg_delete.WaitForCompletion(ctx, nsgClient.Client)
	if err != nil {
		fmt.Errorf("cannot get nsg delete future response: %v", err)
	}
        fmt.Println("Deleted NSG")
//      ************************************@@@@@*********************************************

        ipClient := network.NewPublicIPAddressesClient(subscriptionID)
        ipClient.Authorizer = autorest.NewBearerAuthorizer(token)

        ip_delete, err := ipClient.Delete(ctx, resourcegroup, ip)
        if err != nil {
		fmt.Errorf("cannot create public ip address: %v", err)
	}

	err = ip_delete.WaitForCompletion(ctx, ipClient.Client)
	if err != nil {
		fmt.Errorf("cannot get public ip address create or update future response: %v", err)
	}
        fmt.Println("Deleted IP")

//      ************************************@@@@@*********************************************

        diskClient := compute.NewDisksClient(subscriptionID)
        diskClient.Authorizer = autorest.NewBearerAuthorizer(token)

        disk_delete, err := diskClient.Delete(ctx , resourcegroup, osDisk_Name)
        if err != nil {
		fmt.Errorf("cannot delete disk: %v", err)
	}

	err = disk_delete.WaitForCompletion(ctx, vmClient.Client)
	if err != nil {
		fmt.Errorf("cannot get the disk delete future response: %v", err)
	}
        fmt.Println("Deleted Disk")
}
