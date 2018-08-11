
all: clean deps build deploy

deps:
	dep ensure

build: clean
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ${GOPATH}/bin/bigclown_gateway github.com/vavravl1/bigclown_gateway/

clean:
	rm -f ${GOPATH}/bin/bigclown_gateway

deploy: build
	ssh gateway "killall bigclown_gateway"; \
	ssh gateway "rm /home/vlada/services/bc_gateway/bigclown_gateway"; \
	scp ${GOPATH}/bin/bigclown_gateway gateway:/home/vlada/services/bc_gateway/
#	scp ${GOPATH}/src/github.com/vavravl1/bigclown_gateway/Dockerfile raspberry2:/home/vlada/docker_images/go_bigclown_gateway/ ;\
#    scp ${GOPATH}/bin/bigclown_gateway raspberry2:/home/vlada/docker_images/go_bigclown_gateway/ ;\
#    ssh raspberry2 "docker build /home/vlada/docker_images/go_bigclown_gateway/ -t go-bigclown-gateway:v1"
