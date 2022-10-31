VAULT_TOKEN = root
VAULT_ADDR = http://127.0.0.1:8200
VAULT_NAME = vault-test
VAULT_ENGINE_PATH = kv
KEEPASS_DB_PASS = test

vault/start:
	@docker run -d --name $(VAULT_NAME) --rm -p 8200:8200 --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=$(VAULT_TOKEN)' vault server -dev

vault/stop:
	@docker stop $(VAULT_NAME)

dev/build:
	@go build

dev/lint:
	@docker pull golangci/golangci-lint:latest
	@docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

test/system:
	@VAULT_TOKEN="$(VAULT_TOKEN)" VAULT_ADDR="$(VAULT_ADDR)" ./keepass-vault-sync -f test-data/test.kdbx -g test1,test2,test3 -p $(KEEPASS_DB_PASS)

test/show:
	@VAULT_TOKEN="$(VAULT_TOKEN)" VAULT_ADDR="$(VAULT_ADDR)" vkv -p $(VAULT_ENGINE_PATH) --show-values

test/full: vault/start dev/build test/system test/show vault/stop
