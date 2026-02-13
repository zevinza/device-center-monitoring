GO ?= go

.PHONY: run-master run-client run-device swag-master

# Adjust these package paths if your services live elsewhere

run-master:
	cd app/master-service && $(GO) run .

run-client:
	cd app/client-service && $(GO) run .

run-device:
	cd app/device-simulator && $(GO) run .

swag-master:
	cd app/master-service && swag init -o ../../docs/master -d ./,../../utils,../../entity



