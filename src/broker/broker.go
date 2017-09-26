package broker

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

type BrokerImpl struct {
	cfg      Config
	director DirectorImpl
}

func NewBrokerImpl(cfg Config, director DirectorImpl) BrokerImpl {
	return BrokerImpl{cfg, director}
}

func (b BrokerImpl) adjustNameInManifest(manifestStr, name string) ([]byte, error) {
	var manifest map[interface{}]interface{}

	err := yaml.Unmarshal([]byte(manifestStr), &manifest)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling manifest: %s", err)
	}

	manifest["name"] = name

	manifestBytes, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("marshaling manifest: %s", err)
	}

	return manifestBytes, nil
}

var (
	stripNonAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")
)

type bindingDeploymentName struct{ id string }

func (i bindingDeploymentName) String() string { return "service-binding-" + i.id }
func (i bindingDeploymentName) AlphaNumFriendly() string {
	return stripNonAlphaNum.ReplaceAllString(i.String(), "")
}

type instanceDeploymentName struct{ id string }

func (i instanceDeploymentName) String() string      { return "service-instance-" + i.id }
func (i instanceDeploymentName) DNSFriendly() string { return strings.Replace(i.String(), "_", "-", -1) }
