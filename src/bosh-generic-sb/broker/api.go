package broker

import (
	btask "bosh-generic-sb/broker/task"
)

type BrokerAPIImpl struct {
	broker      BrokerImpl
	taskService btask.AsyncTaskService
}

func NewBrokerAPIImpl(broker BrokerImpl, taskService btask.AsyncTaskService) *BrokerAPIImpl {
	return &BrokerAPIImpl{broker, taskService}
}
