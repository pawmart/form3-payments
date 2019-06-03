.PHONY: help

APP_NAME=form3-payments-srv
MONGO_NAME=form3-mongo
PACKAGE_NAME=github.com/pawmart/form3-payments
VERSION=latest

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

localtest: ## Run godog locally.
	FORM3_DB_DATABASE='form3payments-test' \
	FORM3_DB_HOST='localhost' \
	FORM3_DB_USER='root' \
    FORM3_DB_PASSWORD='example' \
    FORM3_DB_AUTH='admin' \
    godog

localbuild: ## Run build manually.
	dep ensure && go build -o app

localrun: ## Run locally.
	FORM3_DB_DATABASE='form3payments-test' \
	FORM3_DB_HOST='localhost' \
	FORM3_DB_USER='root' \
	FORM3_DB_PASSWORD='example' \
	FORM3_DB_AUTH='admin' \
	./app

test: mongo ## Test.
	docker run --rm \
		-v ${PWD}:/go/src/${PACKAGE_NAME} \
		-v ${HOME}/.ssh/id_rsa:/root/.ssh/id_rsa \
		-w /go/src/${PACKAGE_NAME} \
		-e FORM3_DB_DATABASE='form3payments-test' \
		-e FORM3_DB_HOST='host.docker.internal' \
		-e FORM3_DB_USER='root' \
		-e FORM3_DB_PASSWORD='example' \
		-e FORM3_DB_AUTH='admin' \
		rat4m3n/go-builder:latest /bin/sh -c "dep ensure && godog"

build: test ## Build.
	docker run --rm \
		-v ${PWD}:/go/src/${PACKAGE_NAME} \
		-v ${HOME}/.ssh/id_rsa:/root/.ssh/id_rsa \
		-w /go/src/${PACKAGE_NAME} \
		-e GOOS=linux \
		-e GOARCH=386 \
		rat4m3n/go-builder:latest /bin/sh -c "dep ensure && go build -o app" && \
	docker build -t ${APP_NAME}:${VERSION} .

up: build ## Build and Run.
	docker run -d \
		--name=${APP_NAME} \
		-p 6543:6543 \
		-e FORM3_DB_DATABASE='form3payments' \
		-e FORM3_DB_HOST='host.docker.internal' \
		-e FORM3_DB_USER='root' \
		-e FORM3_DB_PASSWORD='example' \
		-e FORM3_DB_AUTH='admin' \
		${APP_NAME}

mongo: ## Start storage.
	docker run -d \
		--name ${MONGO_NAME} \
		-p 27017:27017 \
		-e MONGO_INITDB_ROOT_USERNAME='root' \
		-e MONGO_INITDB_ROOT_PASSWORD='example' \
		mongo:3.4-jessie

down: ## Clear.
	docker stop ${APP_NAME} && docker rm ${APP_NAME}
	docker stop ${MONGO_NAME} && docker rm ${MONGO_NAME}

validate: ## Validate swagger.
	${GOPATH}/bin/swagger validate ./swagger/swagger.yml

gen: validate ## Generate server from swagger.
	${GOPATH}/bin/swagger generate server \
		--target=./ \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=form3payments

t-build: ## Trigger build
	cd /Users/pawel/Sites/golang/src/github.com/pawmart/terraform-aws-docker && \
	terraform apply -target=null_resource.build


t-deploy: ## Trigger deployment
	cd /Users/pawel/Sites/golang/src/github.com/pawmart/terraform-aws-docker && \
	terraform apply -target=null_resource.deploy


## DEPLOYMENT RELATED COMMANDS
#
#ENV ?= ${USER}
#
#deps:
#	@hash terraform > /dev/null 2>&1 || (echo "Install terraform to continue"; exit 1)
#	@echo Environment: ${ENV}
#
#plan: deps
#	@cd tf; \
#	terraform get; \
#	TF_VAR_environ="${ENV}" terraform plan --state=${ENV}.tfstate
#
#apply: deps
#	@cd tf; \
#	terraform get; \
#	TF_VAR_environ="${ENV}" terraform apply --state=${ENV}.tfstate
#
#destroy: deps
#	@cd tf; \
#	TF_VAR_environ="${ENV}" terraform destroy --state=${ENV}.tfstate
#
#show: deps
#	@cd tf; \
#	TF_VAR_environ="${ENV}" terraform show ${ENV}.tfstate
#
## Force deploy of the code as it presently is
#force_deploy: deps
#	@cd tf; \
#	terraform get; \
#	terraform taint --state=${ENV}.tfstate null_resource.docker; \
#	TF_VAR_environ="${ENV}" terraform apply --state=${ENV}.tfstate