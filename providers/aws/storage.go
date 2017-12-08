package aws

import (
	"os/exec"
	
	"github.com/nathanwilk7/zcloud/storage"
)

func (p awsProvider) Cp (params storage.CpParams) (string, error) {
	keyId, secret := p.getEnvCreds()
	cmd := exec.Command("aws")
	args := []string{"s3", "cp", ConvertURL(params.Src), ConvertURL(params.Dest)}
	if params.Recursive {
		args = append(args, "--recursive")
	}
	cmd.Args = GetCmdArgs(cmd, args)
	cmd.Env = []string{keyId, secret}
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "Copy completed successfully", nil
}
