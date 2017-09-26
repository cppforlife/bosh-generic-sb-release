package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails,
	asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	if !asyncAllowed {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	task, err := b.taskService.CreateTaskDefault(
		func() (interface{}, error) { return nil, b.broker.Provision(instanceID, details) },
	)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("creating provision task")
	}

	b.taskService.StartTask(task)

	// todo DashboardURL:  dashboardURL,
	return brokerapi.ProvisionedServiceSpec{IsAsync: true, OperationData: task.ID}, nil
}

func (b BrokerImpl) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	planConf, err := b.cfg.FindPlan(details.ServiceID, details.PlanID)
	if err != nil {
		return fmt.Errorf("finding service plan: %s", err)
	}

	depName := instanceDeploymentName{instanceID}

	manifest, err := b.adjustNameInManifest(planConf.ServiceInstanceManifest, depName.String())
	if err != nil {
		return fmt.Errorf("adjusting instance deployment name: %s", err)
	}

	var params map[string]interface{}

	if len(details.RawParameters) > 0 {
		err = json.Unmarshal(details.RawParameters, &params)
		if err != nil {
			return fmt.Errorf("unmarshaling service deployment request params: %s", err)
		}
	}

	manifest, err = NewParameters(planConf.ServiceInstanceParams, b.director).Apply(manifest, params)
	if err != nil {
		return fmt.Errorf("applying service instance parameters: %s", err)
	}

	_, err = b.director.Execute([]string{
		"-d", depName.String(),
		"deploy",
		"-",
		"-v", "si_deployment_name=" + depName.String(),
		"-v", "si_deployment_name_dns_friendly=" + depName.DNSFriendly(),
	}, bytes.NewReader(manifest))
	if err != nil {
		return fmt.Errorf("deploying service instance deployment: %s", err)
	}

	return nil
}
