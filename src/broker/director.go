package broker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type DirectorImpl struct {
	cfg DirectorConfig
}

type DirectorConfig struct {
	Host         string
	CACert       string
	Client       string
	ClientSecret string
}

func NewDirectorImpl(cfg DirectorConfig) DirectorImpl {
	return DirectorImpl{cfg}
}

func (d DirectorImpl) Execute(args []string, stdin io.Reader) ([]byte, error) {
	cmd := exec.Command("bosh", args...)

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

func (d DirectorImpl) ExecuteWithBash(args []string, stdin io.Reader, env []string) ([]byte, error) {
	boshCmd := strings.Join(append([]string{"bosh"}, args...), " ")
	cmd := exec.Command("bash", "-c", boshCmd)

	cmd.Env = append(os.Environ(),
		"BOSH_ENVIRONMENT="+d.cfg.Host,
		"BOSH_CA_CERT="+d.cfg.CACert,
		"BOSH_CLIENT="+d.cfg.Client,
		"BOSH_CLIENT_SECRET="+d.cfg.ClientSecret,
		"BOSH_NON_INTERACTIVE=true",
		"HOME=/tmp",
	)

	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}

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
