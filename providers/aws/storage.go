package aws

import (
	"github.com/nathanwilk7/zcloud/storage"
)

// TODO: How to avoid duplication between Cp and Ls? Maybe use an interface with GetArgs() or something?
func (p awsProvider) Cp (params storage.CpParams) (string, error) {
	args := []string{}
	if params.Recursive {
		args = append(args, "--recursive")
	}
	cmd := awsStorageCmd(
		"cp",
		[]string{params.Src, params.Dest},
		args,
	)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "Copy completed successfully", nil
}

func (p awsProvider) Ls (params storage.LsParams) (string, error) {
	args := []string{}
	if params.Recursive {
		args = append(args, "--recursive")
	}
	cmd := awsStorageCmd(
		"ls",
		[]string{params.Url},
		args,
	)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "List completed successfully", nil
}
