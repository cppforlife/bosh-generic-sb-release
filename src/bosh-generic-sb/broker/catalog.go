package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Services(_ context.Context) []brokerapi.Service {
	return b.broker.Services()
}

func (b BrokerImpl) Services() []brokerapi.Service {
	var srvs []brokerapi.Service

	for _, srvConf := range b.cfg.Services {
		srv := brokerapi.Service{
			ID:            srvConf.ID,
			Name:          srvConf.Name,
			Description:   srvConf.Description,
			Bindable:      true,
			PlanUpdatable: true,

			// Metadata: &brokerapi.ServiceMetadata{
			//  DisplayName:         b.serviceOffering.Metadata.DisplayName,
			//  ImageUrl:            b.serviceOffering.Metadata.ImageURL,
			//  LongDescription:     b.serviceOffering.Metadata.LongDescription,
			//  ProviderDisplayName: b.serviceOffering.Metadata.ProviderDisplayName,
			//  DocumentationUrl:    b.serviceOffering.Metadata.DocumentationURL,
			//  SupportUrl:          b.serviceOffering.Metadata.SupportURL,
			// },
			// DashboardClient: &brokerapi.ServiceDashboardClient{
			//  ID:          b.serviceOffering.DashboardClient.ID,
			//  Secret:      b.serviceOffering.DashboardClient.Secret,
			//  RedirectURI: b.serviceOffering.DashboardClient.RedirectUri,
			// },
			// Requires:        requiredPermissions(b.serviceOffering.Requires),
			// Tags:            b.serviceOffering.Tags,
		}

		var plans []brokerapi.ServicePlan

		for _, planConf := range srvConf.Plans {
			plan := brokerapi.ServicePlan{
				ID:          planConf.ID,
				Name:        planConf.Name,
				Description: planConf.Description,
				Free:        brokerapi.FreeValue(false),
				Bindable:    brokerapi.BindableValue(planConf.AllowsBinding()),

				// Metadata: &brokerapi.ServicePlanMetadata{
				//  DisplayName: plan.Metadata.DisplayName,
				//  Bullets:     plan.Metadata.Bullets,
				//  Costs:       []brokerapi.ServicePlanCost{{Amount: cost.Amount, Unit: cost.Unit}},
				// },
			}

			plans = append(plans, plan)
		}

		srv.Plans = plans

		srvs = append(srvs, srv)
	}

	return srvs
}
