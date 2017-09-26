package broker

import (
	"fmt"
)

type BrokerConfig struct {
	Services []ServiceConfig
}

type ServiceConfig struct {
	ID          string
	Name        string
	Description string

	Plans []ServicePlanConfig
}

type ServicePlanConfig struct {
	ID          string
	Name        string
	Description string

	ServiceInstanceManifest string  `json:"si_manifest"`
	ServiceInstanceParams   []Param `json:"si_params"`
	ServiceBindingManifest  string  `json:"sb_manifest"`
	ServiceBindingParams    []Param `json:"sb_params"`
}

type Param struct {
	Name string
	Ops  []interface{}
}

func (c BrokerConfig) FindPlan(serviceID, planID string) (ServicePlanConfig, error) {
	for _, srv := range c.Services {
		for _, plan := range srv.Plans {
			if srv.ID == serviceID && plan.ID == planID {
				return plan, nil
			}
		}
	}
	return ServicePlanConfig{}, fmt.Errorf(
		"expected to find service ID '%s' with plan ID '%s'", serviceID, planID)
}

func (c ServicePlanConfig) AllowsBinding() bool {
	return len(c.ServiceBindingManifest) > 0
}

func (c BrokerConfig) Validate() error {
	for i, srvConf := range c.Services {
		err := srvConf.Validate()
		if err != nil {
			return fmt.Errorf("validating services[%i]: %s", i, err)
		}
	}
	return nil
}

func (c ServiceConfig) Validate() error {
	if len(c.ID) == 0 {
		return fmt.Errorf("expected non-empty 'ID'")
	}
	if len(c.Name) == 0 {
		return fmt.Errorf("expected non-empty 'Name'")
	}
	if len(c.Description) == 0 {
		return fmt.Errorf("expected non-empty 'Description'")
	}
	for i, planConf := range c.Plans {
		err := planConf.Validate()
		if err != nil {
			return fmt.Errorf("validating plans[%i]: %s", i, err)
		}
	}
	return nil
}

func (c ServicePlanConfig) Validate() error {
	if len(c.ID) == 0 {
		return fmt.Errorf("expected non-empty 'ID'")
	}
	if len(c.Name) == 0 {
		return fmt.Errorf("expected non-empty 'Name'")
	}
	if len(c.Description) == 0 {
		return fmt.Errorf("expected non-empty 'Description'")
	}
	if len(c.ServiceInstanceManifest) == 0 {
		return fmt.Errorf("expected non-empty 'ServiceInstanceManifest'")
	}
	return nil
}
