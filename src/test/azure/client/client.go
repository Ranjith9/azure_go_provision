package client

import (
        "github.com/Azure/go-autorest/autorest/adal"
        "github.com/Azure/go-autorest/autorest/azure"
        "log"
)

var (
        clientID =        "1b390a83-5255-47e3-bcd3-a5e41061e661"
        subscriptionID =  "0594cd49-9185-425d-9fe2-8d051e4c6054"
        tenantID =        "85c997b9-f494-46b3-a11d-772983cf6f11"
        clientSecret =    "ceZBgXQoryOMGvK6txScc/TruRGaHucs9uayj8d/OtI="
)


func GetServicePrincipalToken() (adal.OAuthTokenProvider, error) {
        oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
        code, err := adal.NewServicePrincipalToken(
                *oauthConfig,
                clientID,
                clientSecret,
                azure.PublicCloud.ResourceManagerEndpoint)
        if err != nil {
                log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
        }

       return code, err
}
