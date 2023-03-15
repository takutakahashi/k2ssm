package output

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Sops struct {
	Raw            Output
	OutputPath     string
	TempFilePath   string
	SopsBinaryPath string
}

type SopsOutput struct {
	Name string
	Tags map[string]string
	Data map[string]string
}

func (o Sops) Execute(ctx context.Context) error {
	f, err := os.CreateTemp("/tmp", "k2ssm-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp file")
	}
	defer os.Remove(f.Name())
	buf, err := yaml.Marshal(o.convertOutput())
	if err != nil {
		return errors.Wrap(err, "failed to marshal to yaml")
	}
	if _, err := f.Write(buf); err != nil {
		return errors.Wrap(err, "failed to write data to temp file")
	}
	c := exec.Command(o.SopsBinaryPath, "-e", f.Name())
	c.Env = os.Environ()
	out, err := c.Output()
	if err != nil {
		return errors.Wrap(err, "failed to encrypt by sops")
	}
	return os.WriteFile(o.OutputPath, out, 0644)
}

func (o Sops) convertOutput() []SopsOutput {
	ret := []SopsOutput{}
	for _, s := range o.Raw.Secrets {
		ret = append(ret, SopsOutput{
			Name: s.Name,
			Tags: s.Labels,
			Data: s.Data,
		})
	}
	return ret
}
