package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Unbind(ctx context.Context, instanceID, bindingID string,
	details brokerapi.UnbindDetails) error {

	// todo async please
	return b.broker.Unbind(instanceID, bindingID, details)
}

func (b BrokerImpl) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	planConf, err := b.cfg.FindPlan(details.ServiceID, details.PlanID)
	if err != nil {
		return fmt.Errorf("finding service plan: %s", err)
	}

	if !planConf.AllowsBinding() {
		return fmt.Errorf("binding is not supported")
	}

	depName := bindingDeploymentName{bindingID}

	err = errandBinding{depName, b.director}.Delete()
	if err != nil {
		return err
	}

	_, err = b.director.Execute([]string{"-d", depName.String(), "delete-deployment"}, nil)
	if err != nil {
		return fmt.Errorf("deleting service binding deployment: %s", err)
	}

	return nil
}
