package aws

import (
	"github.com/nathanwilk7/zcloud/storage"
)

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

func (p awsProvider) Rm (params storage.RmParams) (string, error) {
	return "", nil
}

func (p awsProvider) Mv (params storage.MvParams) (string, error) {
	return "", nil
}

func (p awsProvider) Mb (params storage.MbParams) (string, error) {
	return "", nil
}

func (p awsProvider) Rb (params storage.RmParams) (string, error) {
	return "", nil
}

func (p awsProvider) Sync (params storage.SyncParams) (string, error) {
	return "", nil
}
