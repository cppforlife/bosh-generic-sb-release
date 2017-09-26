package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Parameters struct {
	params   []Param
	director DirectorImpl
}

func NewParameters(params []Param, director DirectorImpl) Parameters {
	return Parameters{params, director}
}

func (p Parameters) Apply(manifest []byte, givenParams map[string]interface{}) ([]byte, error) {
	for name, value := range givenParams {
		param, err := p.findParam(name)
		if err != nil {
			return nil, err
		}

		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshaling service deployment request param value: %s", err)
		}

		opsBytes, err := yaml.Marshal(param.Ops)
		if err != nil {
			return nil, fmt.Errorf("marshaling service deployment request ops: %s", err)
		}

		interpolatedOps, err := p.director.Execute(
			[]string{"int", "-", "-v", fmt.Sprintf("value=%s", jsonValue)},
			bytes.NewReader(opsBytes),
		)
		if err != nil {
			return nil, fmt.Errorf("interpolating create service param value: %s", err)
		}

		manifest, err = p.director.ExecuteWithBash(
			[]string{"int", "-", "-o", `<(echo "$SB_OPS")`},
			bytes.NewReader(manifest),
			[]string{"SB_OPS=" + string(interpolatedOps)},
		)
		if err != nil {
			return nil, fmt.Errorf("interpolating create service params: %s", err)
		}
	}

	return manifest, nil
}

func (p Parameters) findParam(name string) (Param, error) {
	for _, param := range p.params {
		if param.Name == name {
			return param, nil
		}
	}
	return Param{}, fmt.Errorf("unexpected request param '%s'", name)
}
