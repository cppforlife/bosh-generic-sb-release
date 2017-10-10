package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pivotal-cf/brokerapi"
)

func (b *BrokerAPIImpl) Bind(ctx context.Context, instanceID, bindingID string,
	details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// todo async please
	return b.broker.Bind(instanceID, bindingID, details)
}

func (b BrokerImpl) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	planConf, err := b.cfg.FindPlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("finding service plan: %s", err)
	}

	if !planConf.AllowsBinding() {
		return brokerapi.Binding{}, fmt.Errorf("binding is not supported")
	}

	depName := bindingDeploymentName{bindingID}

	bindingManifest, err := b.adjustNameInManifest(planConf.ServiceBindingManifest, depName.String())
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("adjusting binding deployment name: %s", err)
	}

	var params map[string]interface{}

	if len(details.RawParameters) > 0 {
		err = json.Unmarshal(details.RawParameters, &params)
		if err != nil {
			return brokerapi.Binding{}, fmt.Errorf("unmarshaling service binding request params: %s", err)
		}
	}

	bindingManifest, err = NewParameters(planConf.ServiceBindingParams, b.director).Apply(bindingManifest, params)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("applying service binding parameters: %s", err)
	}

	directorName, err := b.director.Execute([]string{"env", "--column", "name"}, nil)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("finding director name: %s", err)
	}

	_, err = b.director.Execute([]string{
		"-d", depName.String(),
		"deploy",
		"-",
		"-v", "si_director_name=" + strings.TrimSpace(string(directorName)),
		"-v", "si_deployment_name=" + instanceDeploymentName{instanceID}.String(),
		"-v", "sb_deployment_name=" + depName.String(),
		"-v", "sb_deployment_name_alphanum_friendly=" + depName.AlphaNumFriendly(),
	}, bytes.NewReader(bindingManifest))
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("deploying service binding deployment: %s", err)
	}

	return errandBinding{depName, b.director}.Create()
}
