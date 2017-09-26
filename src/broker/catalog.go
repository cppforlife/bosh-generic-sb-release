package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Services(_ context.Context) []brokerapi.Service {
	return []brokerapi.Service{
		{
			ID:            "service-id",
			Name:          "serice-name",
			Description:   "serice-desc",
			Bindable:      true,
			PlanUpdatable: true,

			Plans: []brokerapi.ServicePlan{{
				ID:          "serice-plan-id",
				Name:        "serice-plan-name",
				Description: "serice-plan-description",
				Free:        brokerapi.FreeValue(false),
				Bindable:    brokerapi.BindableValue(true),
				// Metadata: &brokerapi.ServicePlanMetadata{
				// 	DisplayName: plan.Metadata.DisplayName,
				// 	Bullets:     plan.Metadata.Bullets,
				// 	Costs:       []brokerapi.ServicePlanCost{{Amount: cost.Amount, Unit: cost.Unit}},
				// },
			}},
			// Metadata: &brokerapi.ServiceMetadata{
			// 	DisplayName:         b.serviceOffering.Metadata.DisplayName,
			// 	ImageUrl:            b.serviceOffering.Metadata.ImageURL,
			// 	LongDescription:     b.serviceOffering.Metadata.LongDescription,
			// 	ProviderDisplayName: b.serviceOffering.Metadata.ProviderDisplayName,
			// 	DocumentationUrl:    b.serviceOffering.Metadata.DocumentationURL,
			// 	SupportUrl:          b.serviceOffering.Metadata.SupportURL,
			// },
			// DashboardClient: &brokerapi.ServiceDashboardClient{
			// 	ID:          b.serviceOffering.DashboardClient.ID,
			// 	Secret:      b.serviceOffering.DashboardClient.Secret,
			// 	RedirectURI: b.serviceOffering.DashboardClient.RedirectUri,
			// },
			// Requires:        requiredPermissions(b.serviceOffering.Requires),
			// Tags:            b.serviceOffering.Tags,
		},
	}
}
