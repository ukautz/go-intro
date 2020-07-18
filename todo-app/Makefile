
APP_ADDR ?= 127.0.01:12345
APP_STORE_DIR ?= store
APP_CREATE_TITLE ?= random $(shell echo $$RANDOM)
APP_USER ?= alice
APP_PASSWORD ?= secret1
APP_ID ?= missing_APP_ID
APP_URL_PREFIX = /v1

.PHONY: create-todo
create-todo:
	curl -s -u "$(APP_USER):$(APP_PASSWORD)" -X POST http://$(APP_ADDR)$(APP_URL_PREFIX)/todo \
		-d "{\"title\":\"$(APP_CREATE_TITLE)\", \"description\":\"a longer description\"}" | jq

.PHONY: list-todos
list-todos:
	curl -s -u "$(APP_USER):$(APP_PASSWORD)" http://$(APP_ADDR)$(APP_URL_PREFIX)/todo | jq

.PHONY: get-todo
get-todo:
	curl -s -u "$(APP_USER):$(APP_PASSWORD)" http://$(APP_ADDR)$(APP_URL_PREFIX)/todo/$(APP_ID) | jq

.PHONY: delete-todo
delete-todo:
	curl -s -u "$(APP_USER):$(APP_PASSWORD)" -X DELETE http://$(APP_ADDR)$(APP_URL_PREFIX)/todo/$(APP_ID) | jq

.PHONY: run-server
run-server:
	@echo "hit ctrl+c to stop"
	go run cmd/server/main.go

.PHONY: run
run: run-server
