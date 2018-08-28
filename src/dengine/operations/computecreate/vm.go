package main

import (
	"dengine/interface/azureinterface/computeinterface"
	"dengine/interface/azureinterface/networkinterface"
	"encoding/json"
	"fmt"
)

func main() {
	v := azurenetwork.VnetIn{"test", "test", "192.168.0.0/16", "CentralIndia"}
	vnet, _ := v.CreateVirtualNetwork()
	jsonMarshall(vnet)

	s := azurenetwork.SubnetIn{"test", "test", "test-sub", "192.168.10.0/24"}
	subnet, _ := s.CreateVirtualNetworkSubnet()
	jsonMarshall(subnet)

	n := azurenetwork.NsgIn{"M1038273", "test-nsg", "CentralIndia"}
	nsg, _ := n.CreateNetworkSecurityGroup()
	jsonMarshall(nsg)

	nr := azurenetwork.SecurityRuleIn{"M1038273", "test-nsg", "ssh", "22", 123}
	nrule, _ := nr.CreateNetworkSecurityRule()
	jsonMarshall(nrule)

	i := azurenetwork.IpIn{"test", "test-ip", "CentralIndia"}
	ip, _ := i.CreatePublicIP()
	jsonMarshall(ip)

	ic := azurenetwork.NicIn{"test", "test-nic", *nsg.ID, *subnet.ID, *ip.ID, "CentralIndia"}
	nic, _ := ic.CreateNIC()
	jsonMarshall(nic)

	m := azurecompute.VMIn{"test", "test-vm", *nic.ID, "ubuntu", "ubuntu@12345", "ubuntu", "", "CentralIndia"}
	vm, _ := m.CreateVM()
	jsonMarshall(vm)

}

func jsonMarshall(i interface{}) {
	json, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(json))
}
