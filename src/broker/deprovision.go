package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails,
	asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	if !asyncAllowed {
		return brokerapi.DeprovisionServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	task, err := b.taskService.CreateTaskDefault(
		func() (interface{}, error) { return nil, b.broker.Deprovision(instanceID, details) },
	)
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("creating deprovision task")
	}

	b.taskService.StartTask(task)

	return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: task.ID}, nil
}

func (b BrokerImpl) Deprovision(instanceID string, _ brokerapi.DeprovisionDetails) error {
	_, err := b.director.Execute([]string{"-d", instanceDeploymentName{instanceID}.String(), "delete-deployment"}, nil)
	if err != nil {
		return fmt.Errorf("deleting service instance deployment: %s", err)
	}

	return nil
}
