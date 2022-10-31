package main

const (
	k2EnginePath         = "kv"
	dbPath               = "./test-data/test.kdbx"
	dbPass               = "test"
	mountEnginePath      = "sys/mounts/%s"
	readWriteSecretsPath = "%s/data/%s"
)

var syncGroups = []string{"test1", "test2", "test3"}
