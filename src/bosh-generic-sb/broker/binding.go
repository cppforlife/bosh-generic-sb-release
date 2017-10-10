package broker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pivotal-cf/brokerapi"
)

type errandBinding struct {
	depName  bindingDeploymentName
	director DirectorImpl
}

type errand struct{ Name string }

type cliErrandResultsJSON struct{ Tables []cliErrandResultsTableJSON }
type cliErrandResultsTableJSON struct{ Rows []cliErrandResultsRowJSON }
type cliErrandResultsRowJSON struct{ Stdout string }

func (b errandBinding) Create() (brokerapi.Binding, error) {
	errand, err := b.find("create")
	if err != nil {
		return brokerapi.Binding{}, err
	}

	bindingOutput, err := b.director.Execute([]string{
		"-d", b.depName.String(),
		"run-errand",
		errand.Name,
		"--column", "stdout",
		"--json",
	}, nil)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("running errand '%s': %s", errand.Name, err)
	}

	var result cliErrandResultsJSON

	err = json.Unmarshal(bindingOutput, &result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("unmarshaling errand result: %s", err)
	}

	creds := json.RawMessage([]byte(result.Tables[0].Rows[0].Stdout))

	return brokerapi.Binding{Credentials: creds}, nil
}

func (b errandBinding) Delete() error {
	errand, err := b.find("delete")
	if err != nil {
		return err
	}

	_, err = b.director.Execute([]string{"-d", b.depName.String(), "run-errand", errand.Name}, nil)
	if err != nil {
		return fmt.Errorf("running errand '%s': %s", errand.Name, err)
	}

	return nil
}

type cliErrandsJSON struct{ Tables []cliErrandsTableJSON }
type cliErrandsTableJSON struct{ Rows []cliErrandsRowJSON }
type cliErrandsRowJSON struct{ Name string }

func (b errandBinding) find(prefix string) (errand, error) {
	prefix += "-"

	errandsOutput, err := b.director.Execute([]string{
		"-d", b.depName.String(),
		"errands",
		"--column", "name",
		"--json",
	}, nil)
	if err != nil {
		return errand{}, fmt.Errorf("listing errands: %s", err)
	}

	var result cliErrandsJSON

	err = json.Unmarshal(errandsOutput, &result)
	if err != nil {
		return errand{}, fmt.Errorf("unmarshaling errands result: %s", err)
	}

	var potentialErrands []errand
	var allNames []string

	for _, e := range result.Tables[0].Rows {
		if strings.HasPrefix(e.Name, prefix) {
			potentialErrands = append(potentialErrands, errand{e.Name})
		}
		allNames = append(allNames, e.Name)
	}

	if len(potentialErrands) != 1 {
		return errand{}, fmt.Errorf("expected to find an errand matching '%s*' in '%s'", prefix, strings.Join(allNames, "', '"))
	}

	return potentialErrands[0], nil
}
