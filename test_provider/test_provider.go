package test_provider

import (
	"fmt"
	"log"
	"os/exec"
)

type TestProvider struct {}

func (p TestProvider) getEnvCreds () (string, string) {
	// os.Expand, ZCLOUD_AWS_KEY_ID, ZCLOUD_AWS_SECRET_KEY
	return "ID", "SECRET"
}

func (p TestProvider) Cp (src, dst string) {
	keyId, secret := p.getEnvCreds()
	cmd := exec.Command("aws", "s3", "cp", src, dst)
	env := cmd.Env
	env = append(env, keyId)
	env = append(env, secret)
	cmd.Env = env
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(string(out))
	}
}

func (p TestProvider) Start () {
	cmd := exec.Command("echo", "STARTED")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(string(out))
	}
}
