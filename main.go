package main

import (
	"fmt"
	"log"
	"os"

	kp "github.com/tobischo/gokeepasslib/v3"
)

func main() {
	v, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	v.EnableKV2Engine(k2EnginePath)

	file, err := os.Open(dbPath)
	if err != nil {
		fmt.Println("Failed to open file(%s): %s", file, err)
	}
	defer file.Close()

	db := kp.NewDatabase()
	db.Credentials = kp.NewPasswordCredentials(dbPass)
	_ = kp.NewDecoder(file).Decode(db)

	db.UnlockProtectedEntries()

	v.SearchAndWriteRecursive(db.Content.Root.Groups, syncGroups)

}
