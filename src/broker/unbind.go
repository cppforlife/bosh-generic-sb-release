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
	if !b.cfg.AllowsServiceBinding() {
		return fmt.Errorf("binding is not supported")
	}

	depName := bindingDeploymentName{bindingID}

	_, err := b.director.Execute([]string{
		"-d", depName.String(),
		"run-errand",
		"delete-service-binding",
	}, nil)
	if err != nil {
		return fmt.Errorf("running service binding errand: %s", err)
	}

	_, err = b.director.Execute([]string{"-d", depName.String(), "delete-deployment"}, nil)
	if err != nil {
		return fmt.Errorf("deleting service binding deployment: %s", err)
	}

	return nil
}
