package task

// taken from https://github.com/cloudfoundry/bosh-agent/blob/master/agent/task/async_task_service.go

import (
	"sync"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type AsyncTaskService struct {
	uuidGen boshuuid.Generator

	tasks     map[string]Task
	tasksLock sync.RWMutex

	logTag string
	logger boshlog.Logger
}

func NewAsyncTaskService(uuidGen boshuuid.Generator, logger boshlog.Logger) AsyncTaskService {
	return AsyncTaskService{
		uuidGen: uuidGen,
		tasks:   map[string]Task{},

		logTag: "task.AsyncTaskService",
		logger: logger,
	}
}

func (service AsyncTaskService) CreateTaskDefault(taskFunc Func) (Task, error) {
	return service.CreateTask(taskFunc, nil)
}

func (service AsyncTaskService) CreateTask(taskFunc Func, cancelFunc CancelFunc) (Task, error) {
	uuid, err := service.uuidGen.Generate()
	if err != nil {
		return Task{}, err
	}

	return service.CreateTaskWithID(uuid, taskFunc, cancelFunc), nil
}

func (service AsyncTaskService) CreateTaskWithID(
	id string,
	taskFunc Func,
	cancelFunc CancelFunc,
) Task {
	return Task{
		ID:    id,
		State: StateRunning,

		Func:       taskFunc,
		CancelFunc: cancelFunc,
	}
}

func (service AsyncTaskService) StartTask(task Task) {
	service.tasksLock.Lock()
	service.tasks[task.ID] = task
	service.tasksLock.Unlock()

	go func() {
		// todo update task on panic
		defer service.logger.HandlePanic("Task process execution")

		value, err := task.Func()
		if err != nil {
			task.Error = err
			task.State = StateFailed
			service.logger.Error(service.logTag, "Failed task #%s got: %s", task.ID, err.Error())
		} else {
			task.Value = value
			task.State = StateDone
			service.logger.Error(service.logTag, "Completed task #%s got: %s", task.ID)
		}

		// Nil to prevent to memory leaks in case these are closures.
		task.Func = nil
		task.CancelFunc = nil

		service.tasksLock.Lock()
		service.tasks[task.ID] = task
		service.tasksLock.Unlock()
	}()
}

func (service AsyncTaskService) FindTaskWithID(id string) (Task, bool) {
	service.tasksLock.Lock()
	task, found := service.tasks[id]
	service.tasksLock.Unlock()

	return task, found
}
