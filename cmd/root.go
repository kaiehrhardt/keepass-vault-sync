package cmd

import (
	"log"
	"os"
	"syscall"

	"github.com/kaiehrhardt/keepass-vault-sync/pkg/vault"
	"github.com/spf13/cobra"
	kp "github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"
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
		v, err := vault.NewClient()
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

		v.SearchAndWriteRecursive(enginePath, db.Content.Root.Groups, syncGroups)
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
