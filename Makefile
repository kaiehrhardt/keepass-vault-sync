VAULT_TOKEN = root
VAULT_ADDR = http://127.0.0.1:8200
VAULT_NAME = vault-test

vault/start:
	docker run -d --name $(VAULT_NAME) --rm -p 8200:8200 --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=$(VAULT_TOKEN)' vault server -dev

dev/build:
	go build

dev/systemtest:
	VAULT_TOKEN="$(VAULT_TOKEN)" VAULT_ADDR="$(VAULT_ADDR)" ./keepass-vault-sync

dev/show:
	VAULT_TOKEN="$(VAULT_TOKEN)" VAULT_ADDR="$(VAULT_ADDR)" vkv -p kv --show-values

vault/stop:
	docker stop $(VAULT_NAME)

dev/fullsystemtest: vault/start dev/build dev/systemtest dev/show vault/stop
