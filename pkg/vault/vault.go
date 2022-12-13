package vault

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

const (
	mountEnginePath      = "sys/mounts/%s"
	readWriteSecretsPath = "%s/data/%s"
)

type Vault struct {
	Client *api.Client
}

func NewClient() (*Vault, error) {
	_, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		return nil, fmt.Errorf("VAULT_ADDR required but not set")
	}

	vaultToken, ok := os.LookupEnv("VAULT_TOKEN")
	if !ok {
		return nil, fmt.Errorf("VAULT_TOKEN required but not set")
	}

	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}

	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c.SetToken(vaultToken)

	vaultNamespace, ok := os.LookupEnv("VAULT_NAMESPACE")
	if ok {
		c.SetNamespace(vaultNamespace)
	}

	return &Vault{Client: c}, nil
}

func (v *Vault) EnableKV2Engine(rootPath string) error {
	options := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"path":    rootPath,
			"version": 2,
		},
	}

	_, err := v.Client.Logical().Write(fmt.Sprintf(mountEnginePath, rootPath), options)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) WriteSecrets(enginePath string, subPath string, secret map[string]interface{}) error {
	options := map[string]interface{}{
		"data": secret,
	}
	_, err := v.Client.Logical().Write(fmt.Sprintf(readWriteSecretsPath, enginePath, subPath), options)
	if err != nil {
		return err
	}

	return nil
}
