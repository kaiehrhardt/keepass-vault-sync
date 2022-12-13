package vault

import (
	"log"

	"github.com/kaiehrhardt/keepass-vault-sync/pkg/utils"
	kp "github.com/tobischo/gokeepasslib/v3"
)

func (v *Vault) SearchAndWriteRecursive(enginePath string, groups []kp.Group, syncGroups []string) {
	for _, g := range groups {
		if utils.Contains(syncGroups, g.Name) {
			secret := make(map[string]interface{})
			for _, e := range g.Entries {
				secret[e.GetTitle()] = e.GetPassword()
			}
			err := v.WriteSecrets(enginePath, g.Name, secret)
			if err != nil {
				log.Fatal(err)
			}
		}
		v.SearchAndWriteRecursive(enginePath, g.Groups, syncGroups)
	}
}
