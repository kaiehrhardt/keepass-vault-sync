package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	kp "github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"
)

const (
	mountEnginePath      = "sys/mounts/%s"
	readWriteSecretsPath = "%s/data/%s"
)

var (
	dbPath      string
	syncGroups  []string
	enginePath  string
	dbPassStdin string
)

var rootCmd = &cobra.Command{
	Use:   "keepass-vault-sync",
	Short: "Sync secrets keepass -> vault",
	Long:  `Sync secrets keepass -> vault`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		if err := v.EnableKV2Engine(enginePath); err != nil {
			log.Fatal(err)
		}

		file, err := os.Open(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if len(dbPassStdin) == 0 {
			log.Println("Datebase Password: ")
			bytepw, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatal(err)
			}
			dbPassStdin = string(bytepw)
		}

		db := kp.NewDatabase()
		db.Credentials = kp.NewPasswordCredentials(dbPassStdin)
		err = kp.NewDecoder(file).Decode(db)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Login OK")
		}
		err = db.UnlockProtectedEntries()
		if err != nil {
			log.Fatal(err)
		}

		v.SearchAndWriteRecursive(db.Content.Root.Groups, syncGroups)
		log.Println("Sync done")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "dbPath", "f", "", "path to .kdbx")
	rootCmd.PersistentFlags().StringSliceVarP(&syncGroups, "syncGroups", "g", []string{}, "comma separated list of groups to sync")
	rootCmd.PersistentFlags().StringVarP(&enginePath, "enginePath", "e", "kv", "vault engine path")
	rootCmd.PersistentFlags().StringVarP(&dbPassStdin, "dbPassStdin", "p", "", ".kdbx password")

	if err := rootCmd.MarkPersistentFlagRequired("dbPath"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("syncGroups"); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

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

func (v *Vault) WriteSecrets(subPath string, secret map[string]interface{}) error {
	options := map[string]interface{}{
		"data": secret,
	}
	_, err := v.Client.Logical().Write(fmt.Sprintf(readWriteSecretsPath, enginePath, subPath), options)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) SearchAndWriteRecursive(groups []kp.Group, syncGroups []string) {
	for _, g := range groups {
		if contains(syncGroups, g.Name) {
			secret := make(map[string]interface{})
			for _, e := range g.Entries {
				secret[e.GetTitle()] = e.GetPassword()
			}
			err := v.WriteSecrets(g.Name, secret)
			if err != nil {
				log.Fatal(err)
			}
		}
		v.SearchAndWriteRecursive(g.Groups, syncGroups)
	}
}
