package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails,
	asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	provDetails := brokerapi.ProvisionDetails{
		ServiceID:     details.ServiceID,
		PlanID:        details.PlanID,
		RawParameters: details.RawParameters,
	}

	provSpec, err := b.Provision(ctx, instanceID, provDetails, asyncAllowed)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, err
	}

	return brokerapi.UpdateServiceSpec{IsAsync: provSpec.IsAsync, OperationData: provSpec.OperationData}, nil
}
