package main

import (
         "azure/interface/networkinterface"
         "fmt"
         "encoding/json"
)

func main() {
/*         l := azurenetwork.VnetIn{ResourceGroup:"M1038273"}
         list, _ := l.ListVirtualNetwork()
         json_val, _ := json.Marshal(list)

         fmt.Println(string(json_val))

         c := azurenetwork.VnetIn{"M1038273", "test", "192.168.0.0/16", "CentralIndia"}
         create, _ := c.CreateVirtualNetwork()
         json, _ := json.MarshalIndent(create, "", "  ")

         fmt.Println(string(json))
*/
         la,_ := azurenetwork.ListAllVirtualNetwork()
         for _,t := range la {
                  k,_ := json.Marshal(t)
                  fmt.Println(string(k))
              fmt.Println(*t.Name)
         }

//         r := azurenetwork.VnetIn{ResourceGroup: "M1038273", VnetName: "test"}
//        resp, _ := r.DeleteVirtualNetwork()
//         if resp.Response.StatusCode == 200 { fmt.Println("Deleted network "+r.VnetName+" successfully")}

}
