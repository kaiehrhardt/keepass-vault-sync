package main

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	keepassDbPass = "test"
	keepassDb     = "./tests/testdata/test.kdbx"
	groups        = "test1,test2"
	bin           = "./keepass-vault-sync"
	vaultToken    = "root"
	vaultPort     = "8200"
)

func TestIntegrationWithLatestVault(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "vault:latest",
		ExposedPorts: []string{vaultPort + "/tcp"},
		WaitingFor:   wait.ForLog("Development mode should NOT be used in production installations!"),
		Cmd:          []string{"server", "-dev"},
		Env: map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID": vaultToken,
		},
	}

	vaultC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}

	ip, err := vaultC.Host(ctx)
	if err != nil {
		t.Error(err)
	}

	mappedPort, err := vaultC.MappedPort(ctx, vaultPort)
	if err != nil {
		t.Error(err)
	}

	vaultAddr := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	cmd := exec.Command(bin, "-f", keepassDb, "-g", groups, "-p", keepassDbPass)
	cmd.Env = append(cmd.Env, "VAULT_TOKEN="+vaultToken, "VAULT_ADDR="+vaultAddr)
	output, err := cmd.Output()
	if err != nil {
		t.Log(output, cmd.Environ(), cmd.String())
		t.Error(err)
	}

	config := vault.DefaultConfig()
	config.Address = vaultAddr

	c, err := vault.NewClient(config)
	if err != nil {
		t.Error(err)
	}

	c.SetToken(vaultToken)

	testCases := []struct {
		path  string
		key   string
		match string
	}{
		{
			path:  "test1",
			key:   "asdasda",
			match: "test",
		},
		{
			path:  "test2",
			key:   "test2",
			match: "test",
		},
	}

	for _, tc := range testCases {
		secret, err := c.KVv2("kv").Get(ctx, tc.path)
		if err != nil {
			t.Error(err)
		}

		value, ok := secret.Data[tc.key].(string)
		if !ok {
			t.Error(err)
		}

		assert.Equal(t, tc.match, value)
	}

	defer func() {
		if err := vaultC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminiate container: %s", err.Error())
		}
	}()
}
