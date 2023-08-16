# keepass-vault-sync

[![Go Report Card](https://goreportcard.com/badge/github.com/kaiehrhardt/keepass-vault-sync)](https://goreportcard.com/report/github.com/kaiehrhardt/keepass-vault-sync)
[![container build](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/container-build.yaml/badge.svg)](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/container-build.yaml)
[![go test](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/go-test.yml/badge.svg)](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaiehrhardt/keepass-vault-sync/actions/workflows/golangci-lint.yml)

Simple little utility that syncs certain groups from a keepass db to a path in vault.

## install

### golang

```bash
go install github.com/kaiehrhardt/keepass-vault-sync@latest
```

## docker

```bash
docker pull ghcr.io/kaiehrhardt/keepass-vault-sync:edge
```

## getting started

```bash
VAULT_TOKEN="{your_vault_token}" VAULT_ADDR="{your_vault_addr}" keepass-vault-sync \
  -f {name}.kdbx \
  -g {group1_name},{group2_name} \
  -p {keepass_db_pass}
```
