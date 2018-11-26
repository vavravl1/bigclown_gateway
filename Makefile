
DEPLOY_HOST := "gateway"
DEPLOY_PATH := "/home/vlada/services/bc_gateway/"

all: clean deps build deploy

deps:
	dep ensure

build: clean
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ${GOPATH}/bin/bigclown_gateway github.com/vavravl1/bigclown_gateway/

clean:
	rm -f ${GOPATH}/bin/bigclown_gateway

deploy: build
	ssh $(DEPLOY_HOST) "killall bigclown_gateway"; \
	ssh $(DEPLOY_HOST) "rm $(DEPLOY_PATH)/bigclown_gateway"; \
	scp ${GOPATH}/bin/bigclown_gateway $(DEPLOY_HOST):$(DEPLOY_PATH)/
