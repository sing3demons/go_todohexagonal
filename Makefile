.PHONY: build maria

build:
	go build \
		-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
		-o app

maria:
	docker run -p 127.0.0.1:3306:3306  --name my-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest

image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p:8080:8080 --env-file ./dev.env --link my-mariadb:db \
	--name goapp todo:test