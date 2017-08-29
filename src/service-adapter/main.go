package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Broker struct {
	cfg      Config
	Director Director
}

type Director struct {
	BinaryPath string
	cfg        DirectorConfig
}

type Config struct {
	ServiceInstanceManifest string
	ServiceBindingManifest  string
	Director                DirectorConfig
}

type DirectorConfig struct {
	Host         string
	CACert       string
	Client       string
	ClientSecret string
}

func (c Config) AllowsServiceBinding() bool { return len(c.ServiceBindingManifest) > 0 }

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

	boshBinaryPath := "/var/vcap/packages/odb-service-adapter/bin/bosh"

	respBytes, err := Broker{cfg, Director{boshBinaryPath, cfg.Director}}.BrokerCommand(os.Args)
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

func (s serviceDeploymentJSON) DeploymentNameDNSFriendly() string {
	return strings.Replace(s.DeploymentName, "_", "-", -1)
}

type deploymentManifestYAML struct {
	Name string
}

type CLIErrandResultJSON struct {
	Tables []CLITableJSON
}

type CLITableJSON struct {
	Rows []CLIErrandResultRowJSON
}

type CLIErrandResultRowJSON struct {
	Stdout string
}

var (
	stripNonAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func (b Broker) BrokerCommand(args []string) ([]byte, error) {
	switch args[1] {
	case "generate-manifest":
		// serviceDeploymentJSON(2), planJSON(3), argsJSON(4), previousManifestYAML(5), previousPlanJSON(6)
		var servDep serviceDeploymentJSON

		err := json.Unmarshal([]byte(args[2]), &servDep)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling service deployment arg: %s", err)
		}

		manifest, err := b.adjustNameInManifest(b.cfg.ServiceInstanceManifest, servDep.DeploymentName)
		if err != nil {
			return nil, fmt.Errorf("adjusting instance deployment name: %s", err)
		}

		// todo bosh interpolate?
		manifest = bytes.Replace(manifest,
			[]byte(`((si_deployment_name))`), []byte(servDep.DeploymentName), -1)

		manifest = bytes.Replace(manifest,
			[]byte(`((si_deployment_name_dns_friendly))`), []byte(servDep.DeploymentNameDNSFriendly()), -1)

		return manifest, nil

	case "create-binding":
		// bindingID(2), boshVMsJSON(3), manifestYAML(4), reqParams(5)
		//    Credentials     map[string]interface{} `json:"credentials"`
		//    SyslogDrainURL  string                 `json:"syslog_drain_url,omitempty"`
		//    RouteServiceURL string                 `json:"route_service_url,omitempty"`
		if b.cfg.AllowsServiceBinding() {
			var siDep deploymentManifestYAML

			err := yaml.Unmarshal([]byte(args[4]), &siDep)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling instance deployment manifest: %s", err)
			}

			directorName, err := b.Director.Execute([]string{"env", "--column", "name"}, nil)
			if err != nil {
				return nil, fmt.Errorf("finding director name: %s", err)
			}

			bindingDeploymentName := "service-binding_" + args[2]
			bindingManifest := []byte(b.cfg.ServiceBindingManifest)

			bindingManifest, err = b.adjustNameInManifest(b.cfg.ServiceBindingManifest, bindingDeploymentName)
			if err != nil {
				return nil, fmt.Errorf("adjusting binding deployment name: %s", err)
			}

			_, err = b.Director.Execute([]string{
				"-d", bindingDeploymentName,
				"deploy",
				"-",
				"-v", "si_director_name=" + strings.TrimSpace(string(directorName)),
				"-v", "si_deployment_name=" + siDep.Name,
				"-v", "sb_deployment_name=" + bindingDeploymentName,
				"-v", "sb_deployment_name_alphanum_friendly=" + stripNonAlphaNum.ReplaceAllString(bindingDeploymentName, ""),
			}, bytes.NewReader(bindingManifest))
			if err != nil {
				return nil, fmt.Errorf("deploying service binding deployment: %s", err)
			}

			bindingOutput, err := b.Director.Execute([]string{
				"-d", bindingDeploymentName,
				"run-errand",
				"create-service-binding",
				"--column", "stdout",
				"--json",
			}, nil)
			if err != nil {
				return nil, fmt.Errorf("running service binding errand: %s", err)
			}

			var result CLIErrandResultJSON

			err = json.Unmarshal(bindingOutput, &result)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling errand result: %s", err)
			}

			return []byte(result.Tables[0].Rows[0].Stdout), nil
		}

		return nil, cmdError(10)

	case "delete-binding":
		// bindingID(2), boshVMsJSON(3), manifestYAML(4), unbindingRequestParams(5)
		if b.cfg.AllowsServiceBinding() {
			bindingDeploymentName := "service-binding_" + args[2]

			_, err := b.Director.Execute([]string{
				"-d", bindingDeploymentName,
				"run-errand",
				"delete-service-binding",
			}, nil)
			if err != nil {
				return nil, fmt.Errorf("running service binding errand: %s", err)
			}

			_, err = b.Director.Execute([]string{"-d", bindingDeploymentName, "delete-deployment"}, nil)
			if err != nil {
				return nil, fmt.Errorf("deleting service binding deployment: %s", err)
			}

			return nil, nil
		}

		return nil, nil

	case "dashboard-url":
		// instanceID(2), planJSON(3), manifestYAML(4)
		return nil, cmdError(10)

	default:
		return nil, fmt.Errorf("unknown command '%s'", args[1])
	}
}

func (b Broker) adjustNameInManifest(manifestStr, name string) ([]byte, error) {
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

func (d Director) Execute(args []string, stdin io.Reader) ([]byte, error) {
	cmd := exec.Command(d.BinaryPath, args...)

	cmd.Env = append(os.Environ(),
		"BOSH_ENVIRONMENT="+d.cfg.Host,
		"BOSH_CA_CERT="+d.cfg.CACert,
		"BOSH_CLIENT="+d.cfg.Client,
		"BOSH_CLIENT_SECRET="+d.cfg.ClientSecret,
		"BOSH_NON_INTERACTIVE=true",
		"HOME=/tmp",
	)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Stdin = stdin

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing bosh: %s (stderr: %s)", err, errBuf.String())
	}

	return outBuf.Bytes(), nil
}
