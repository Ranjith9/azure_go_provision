package DengineAzureInterface

import (
        "github.com/Azure/go-autorest/autorest"
        "github.com/Azure/go-autorest/autorest/to"
        "fmt"
        "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
)
/*

Author:  Ranjith Janagama.

Create a IP and add it to a NIC using Network Interface and pass the NIC ID to CreateVM() function

*/

func CreateVm(resourceGroup, name, location, password string, nicID *string) *string {

///////////////////////////////// Creating a VM starts from here /////////////////////////

  vmClient := compute.NewVirtualMachinesClient(SubscriptionID)
  vmClient.Authorizer = autorest.NewBearerAuthorizer(Token)

  response, err := vmClient.CreateOrUpdate(Ctx, resourceGroup, name,
  compute.VirtualMachine{
    Location: to.StringPtr(location),
    VirtualMachineProperties: &compute.VirtualMachineProperties{
      NetworkProfile: &compute.NetworkProfile{
        NetworkInterfaces: &[]compute.NetworkInterfaceReference{
          {
            ID: nicID,
            NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
              Primary: to.BoolPtr(true),
            },
          },
        },
      },
      OsProfile: &compute.OSProfile{
        ComputerName: to.StringPtr(name),
        AdminUsername: to.StringPtr("ubuntu"),
        AdminPassword: to.StringPtr(password),
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
      DiagnosticsProfile: &compute.DiagnosticsProfile{
        BootDiagnostics: &compute.BootDiagnostics{
          Enabled: to.BoolPtr(true),
          StorageURI: to.StringPtr("https://dengine.blob.core.windows.net/"),
        },
      },
    },
  })
  err = response.WaitForCompletion(Ctx, vmClient.Client)
  if err != nil {
    fmt.Errorf("cannot get the vm create or update future response: %v", err)
  }

  result, _ := response.Result(vmClient)
  return result.ID
}

