package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
)

type Broker struct {
  cfg Config
}

type Config struct {
  Manifest string
}

type cmdError int
func (e cmdError) ExitCode() int { return int(e) }
func (e cmdError) Error() string { return fmt.Sprintf("code '%d'", e) }

type ExitCodeError interface {
  ExitCode() int
}

func main() {
  cfg, err := NewConfigFromPath("/var/vcap/jobs/generic-odb-adapter/config/config.json")
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err)
    os.Exit(1)
  }

  respBytes, err := Broker{cfg}.BrokerCommand(os.Args)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err)
    if typedErr, ok := err.(ExitCodeError); ok {
      os.Exit(typedErr.ExitCode())
    }
    os.Exit(1)
  }

  fmt.Fprintf(os.Stdout, "%s", respBytes)
}

func NewConfigFromPath(path string) (Config, error) {
  var cfg Config

  configBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return cfg, fmt.Errorf("reading job config: %s", err)
  }

  err = json.Unmarshal(configBytes, &cfg)
  if err != nil {
    return cfg, fmt.Errorf("unmarshaling job config: %s", err)
  }

  return cfg, nil
}

type serviceDeploymentJSON struct {
  DeploymentName string `json:"deployment_name"`
}

func (b Broker) BrokerCommand(args []string) ([]byte, error) {
  switch args[1] {
  case "generate-manifest":
    // serviceDeploymentJSON(2), planJSON(3), argsJSON(4), previousManifestYAML(5), previousPlanJSON(6)
    var servDep serviceDeploymentJSON

    err := json.Unmarshal([]byte(args[2]), &servDep)
    if err != nil {
      return nil, fmt.Errorf("unmarshaling service deployment arg: %s", err)
    }

    return b.generatedManifest(servDep)

  case "create-binding":
    // bindingID(2), boshVMsJSON(3), manifestYAML(4), reqParams(5)
    return nil, cmdError(10)

  case "delete-binding":
    // bindingID(2), boshVMsJSON(3), manifestYAML(4), unbindingRequestParams(5)
    return nil, cmdError(10)

  case "dashboard-url":
    // instanceID(2), planJSON(3), manifestYAML(4)
    return nil, cmdError(10)

  default:
    return nil, fmt.Errorf("unknown command '%s'", args[1])
  }
}

func (b Broker) generatedManifest(servDep serviceDeploymentJSON) ([]byte, error) {
  var manifest map[interface{}]interface{}

  err := yaml.Unmarshal([]byte(b.cfg.Manifest), &manifest)
  if err != nil {
    return nil, fmt.Errorf("unmarshaling manifest: %s", err)
  }

  manifest["name"] = servDep.DeploymentName

  manifestBytes, err := yaml.Marshal(manifest)
  if err != nil {
    return nil, fmt.Errorf("marshaling manifest: %s", err)
  }

  return manifestBytes, nil
}
