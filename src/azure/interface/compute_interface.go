package auth

import (
        "github.com/Azure/go-autorest/autorest"
        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
//        "context"
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
        "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
)
/*

In this Program I'm creating a VM, NIC, publicIP and a OS disk.
Storage Account will not created with this Program.

Author:  Ranjith Janagama.

*/

func CreateVm(){
/////////////////////////////// Creating a public IP for the VM  /////////////////////////

  pubIpClient := network.NewPublicIPAddressesClient(SubscriptionID)
  pubIpClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  publicip, err := pubIpClient.CreateOrUpdate(Ctx, "test", "test-ip",
  network.PublicIPAddress{
    Location: to.StringPtr("CentralIndia"),
    PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
      PublicIPAllocationMethod: "Dynamic",
      DNSSettings: &network.PublicIPAddressDNSSettings{
        DomainNameLabel: to.StringPtr("test-ip"),
      },
    },
  })
  err = publicip.WaitForCompletion(Ctx, pubIpClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the ip create or update future response: %v", err)
  }

  ip, _ := publicip.Result(pubIpClient)

////////////////////////// Creating a Network Interface Card for the VM //////////////////

  nicClient := network.NewInterfacesClient(SubscriptionID)
  nicClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  nic, err := nicClient.CreateOrUpdate(Ctx, "test", "test-nic",
  network.Interface{
    Location: to.StringPtr("CentralIndia"),
    InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
      IPConfigurations: &[]network.InterfaceIPConfiguration{
        {
          Name: to.StringPtr("ipconfig1"),
          InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
            PrivateIPAllocationMethod: "Dynamic",
            PublicIPAddress: &network.PublicIPAddress{
              ID: ip.ID,
            },
            Subnet: &network.Subnet{
              ID: to.StringPtr("/subscriptions/0594cd49-9185-425d-9fe2-8d051e4c6054/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/go_vpn/subnets/go_sub1"),
            },
          },
        },
      },
      NetworkSecurityGroup: &network.SecurityGroup{
        ID: to.StringPtr("/subscriptions/0594cd49-9185-425d-9fe2-8d051e4c6054/resourceGroups/test/providers/Microsoft.Network/networkSecurityGroups/go_nsg"),
      },
    },
  })
  err = nic.WaitForCompletion(Ctx, nicClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vm create or update future response: %v", err)
  }

  nicard, _ := nic.Result(nicClient)

///////////////////////////////// Creating a VM starts from here /////////////////////////

  vmClient := compute.NewVirtualMachinesClient(SubscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := vmClient.CreateOrUpdate(Ctx, "test", "test-vm",
  compute.VirtualMachine{
    Location: to.StringPtr("CentralIndia"),
    VirtualMachineProperties: &compute.VirtualMachineProperties{
      NetworkProfile: &compute.NetworkProfile{
        NetworkInterfaces: &[]compute.NetworkInterfaceReference{
          {
            ID: nicard.ID,
            NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
              Primary: to.BoolPtr(true),
            },
          },
        },
      },
      OsProfile: &compute.OSProfile{
        ComputerName: to.StringPtr("test-vm"),
        AdminUsername: to.StringPtr("ubuntu"),
        AdminPassword: to.StringPtr("ubuntu@12345"),
        LinuxConfiguration: &compute.LinuxConfiguration{
          DisablePasswordAuthentication: to.BoolPtr(true),
          SSH: &compute.SSHConfiguration{
            PublicKeys: &[]compute.SSHPublicKey{
              {
                Path: to.StringPtr(fmt.Sprintf("/home/ubuntu/.ssh/authorized_keys")),
                KeyData: to.StringPtr("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCjCX8wh0lnk2KUvoCulBER4TQ+4+repQF5vvQeCVc5eWHNQKIPuSxy4fGcEbar15U4wjEJYDsXUGhW0JIh4peIKFf+dXUtZlMQEo7QvPGGORjVm8Zf+je/cVqGQJOvUP4s1/J8EQ+/n6gidtByBL+4lN/vDp/lgPSZzRgb08zVuW40z6jFrxfwalru10FHzzPmkCEtW54YkdJ2yEnLzk+xZDJXmG7JE4c2yRl+Y35HCzHfeRsUqcF1ErV2KYHcRWqwzD9oDZ5V2uTC4ERHkF102Ve7LOSyYK3cvJ8QSWMoOCOPA/UpdrkJRq9e2eVdpIqvnbu2vp6xazU080ZNu/BB"),
              },
            },
          },
        },
      },
      HardwareProfile: &compute.HardwareProfile{
        VMSize: "Basic_A2",
      },
      StorageProfile: &compute.StorageProfile{
        ImageReference: &compute.ImageReference{
          Publisher: to.StringPtr("Canonical"),
          Offer:     to.StringPtr("UbuntuServer"),
          Sku:       to.StringPtr("16.04-LTS"),
          Version:   to.StringPtr("latest"),
        },
      },
    },
  })
  err = response.WaitForCompletion(Ctx, vmClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vm create or update future response: %v", err)
  }

  result, _ := response.Result(vmClient)
  fmt.Println(*result.ID)
}

