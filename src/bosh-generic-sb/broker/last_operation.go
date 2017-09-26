package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"

	btask "bosh-generic-sb/broker/task"
)

func (b *BrokerAPIImpl) LastOperation(ctx context.Context, instanceID, operationDataRaw string,
) (brokerapi.LastOperation, error) {
	taskID := operationDataRaw

	task, found := b.taskService.FindTaskWithID(taskID)
	if !found {
		return brokerapi.LastOperation{}, fmt.Errorf("task '%s' not found", taskID)
	}

	op := brokerapi.LastOperation{}

	switch task.State {
	case btask.StateRunning:
		op.State = brokerapi.InProgress
		op.Description = "running"
	case btask.StateDone:
		op.State = brokerapi.Succeeded
		op.Description = "done"
	case btask.StateFailed:
		op.State = brokerapi.Failed
		op.Description = task.Error.Error()
	default:
		return brokerapi.LastOperation{}, fmt.Errorf("unknown task '%s' state", taskID)
	}

	return op, nil
}
