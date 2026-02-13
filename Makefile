GO ?= go

.PHONY: run-master run-client run-device swag-master

# Adjust these package paths if your services live elsewhere

run-master:
	$(GO) run ./app/master-service

run-client:
	$(GO) run ./app/client-service

run-device:
	$(GO) run ./app/device-simulator

swag-master:
	cd app/master-service && swag init -o ../../docs/master -d ./,../../utils,../../entity



