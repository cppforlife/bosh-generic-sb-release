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

type cliErrandResultJSON struct{ Tables []cliTableJSON }
type cliTableJSON struct{ Rows []cliErrandResultRowJSON }
type cliErrandResultRowJSON struct{ Stdout string }

func (b BrokerImpl) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	if !b.cfg.AllowsServiceBinding() {
		return brokerapi.Binding{}, fmt.Errorf("binding is not supported")
	}

	directorName, err := b.director.Execute([]string{"env", "--column", "name"}, nil)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("finding director name: %s", err)
	}

	depName := bindingDeploymentName{bindingID}
	bindingManifest := []byte(b.cfg.ServiceBindingManifest)

	bindingManifest, err = b.adjustNameInManifest(b.cfg.ServiceBindingManifest, depName.String())
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

	bindingManifest, err = NewParameters(b.cfg.ServiceBindingParams, b.director).Apply(bindingManifest, params)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("applying service binding parameters: %s", err)
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

	bindingOutput, err := b.director.Execute([]string{
		"-d", depName.String(),
		"run-errand",
		"create-service-binding",
		"--column", "stdout",
		"--json",
	}, nil)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("running service binding errand: %s", err)
	}

	var result cliErrandResultJSON

	err = json.Unmarshal(bindingOutput, &result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("unmarshaling errand result: %s", err)
	}

	creds := json.RawMessage([]byte(result.Tables[0].Rows[0].Stdout))

	return brokerapi.Binding{Credentials: creds}, nil
}
