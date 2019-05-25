APP_NAME=form3-payments-srv
PACKAGE_NAME=github.com/pawmart/form3-payments
VERSION=latest

all: test build run ## Run all

test: ## Test
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

build: ## Build
	docker run --rm \
		-v ${PWD}:/go/src/${PACKAGE_NAME} \
		-v ${HOME}/.ssh/id_rsa:/root/.ssh/id_rsa \
		-w /go/src/${PACKAGE_NAME} \
		-e GOOS=linux \
		-e GOARCH=386 \
		rat4m3n/go-builder:latest /bin/sh -c "dep ensure && go build -o app" && \
	docker build -t ${APP_NAME}:${VERSION} .

run: ## Run
	docker run \
		--name=${APP_NAME} \
		-p 6543:6543 \
		-e FORM3_DB_DATABASE='form3payments' \
		-e FORM3_DB_HOST='host.docker.internal' \
		-e FORM3_DB_USER='root' \
		-e FORM3_DB_PASSWORD='example' \
		-e FORM3_DB_AUTH='admin' \
		${APP_NAME}

clear: ## Clear
	docker rm ${APP_NAME}

validate: ## Validate swagger
	${GOPATH}/bin/swagger validate ./swagger/swagger.yml

gen: validate ## Generate server from swagger
	${GOPATH}/bin/swagger generate server \
		--target=./ \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=form3payments