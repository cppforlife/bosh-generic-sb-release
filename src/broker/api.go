package broker

import (
	"broker/task"
)

type BrokerAPIImpl struct {
	broker      BrokerImpl
	taskService task.AsyncTaskService
}

func NewBrokerAPIImpl(broker BrokerImpl, taskService task.AsyncTaskService) *BrokerAPIImpl {
	return &BrokerAPIImpl{broker, taskService}
}
