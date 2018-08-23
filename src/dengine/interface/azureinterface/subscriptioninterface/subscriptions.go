package azuresubscription

import (
        "context"
        "fmt"
        "github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
        "github.com/Azure/go-autorest/autorest"
//        "github.com/Azure/go-autorest/autorest/to"
//        "dengine/interface/azureinterface/getclient"
        "dengine/access/azureaccess"
//        "github.com/Azure/go-autorest/autorest/adal"
)

var (
        token, _, _ = auth.GetServicePrincipalToken()
        ctx = context.Background()
)

type SubcriptionIn struct {
        Subscription string
}

func getSubscriptionClient() subscription.SubscriptionsClient {
        subscriptionClient := subscription.NewSubscriptionsClient()
        subscriptionClient.Authorizer = autorest.NewBearerAuthorizer(token)

        return subscriptionClient
}

func (s SubcriptionIn) GetSubscription() (sub subscription.Model, err error) {
        subscriptionClient := getSubscriptionClient()
        future, err := subscriptionClient.Get(
                ctx,
                s.Subscription,
                )

        if err != nil {
                return sub, fmt.Errorf("cannot get subscription: %v", err)
        }

        return future, err
}

func ListSubscription() (sub []subscription.Model, err error) {
        subscriptionClient := getSubscriptionClient()
        future, err := subscriptionClient.List(
                ctx,
                )

        if err != nil {
                return sub, fmt.Errorf("cannot list subscriptions: %v", err)
        }

        return future.Values(), err
}
