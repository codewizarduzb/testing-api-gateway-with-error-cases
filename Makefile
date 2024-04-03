CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs
