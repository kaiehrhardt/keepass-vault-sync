# keepass-vault-sync

Simple little utility that syncs certain groups from a keepass db to a path in vault.

## install

### golang

```bash
go install github.com/kaiehrhardt/keepass-vault-sync@latest
```

## docker

```bash
docker pull ghcr.io/kaiehrhardt/keepass-vault-sync:latest
```

## getting started

```bash
VAULT_TOKEN="{your_vault_token}" VAULT_ADDR="{your_vault_addr}" keepass-vault-sync -f {name}.kdbx -g {group1_name},{group2_name} -p {keepass_db_pass}
```
