APP_VERSION:=edge
GOLANG_VERSION:=1.17
DOCKER_IMAGE:=abdollahpour/almaniha-draft

compile:
	for i in darwin linux windows ; do \
		GOOS="$${i}" GOARCH=amd64 go build -ldflags "-X main.Version=$(APP_VERSION)" -o bin/mpg-server-"$${i}"-amd64 cmd/server/main.go; \
	done

archive:
	rm -f bin/*.zip
	for i in darwin linux windows ; do \
		zip -j "bin/mpg-$${i}-amd64.zip" "bin/mpg-server-$${i}-amd64" -x "*.DS_Store"; \
	done

run:
	go run cmd/server/main.go

get:
	go get -d -u ./...

image:
	docker build \
		--cache-from "$(DOCKER_IMAGE):builder" \
		--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
		--tag "$(DOCKER_IMAGE):builder" \
		--tag builder \
		--file docker/Dockerfile.builder .
	docker build \
		--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
		--build-arg APP_VERSION="$(APP_VERSION)" \
		--tag "$(DOCKER_IMAGE):$(APP_VERSION)" \
		--tag server \
		--file docker/Dockerfile.server .

push:
	docker push "$(DOCKER_IMAGE):$(APP_VERSION)"
	# We update latest when a real version change happens
	if [ "$(APP_VERSION)" != "edge" ]; then \
		docker tag "$(DOCKER_IMAGE):$(APP_VERSION)" "$(DOCKER_IMAGE):latest"; \
		docker push "$(DOCKER_IMAGE):latest"; \
	fi

test:
	go test -covermode=count -coverprofile=coverage.out -cover $$(go list ./... | grep -v /smoke_test/)

smoke:
	docker-compose -f smoke_test/search/docker-compose.yml up -d
	go test ./smoke_test/search/...
	# It will stay up in case in error
	docker-compose -f smoke_test/search/docker-compose.yml down

goveralls:
	$$GOPATH/bin/goveralls -service=travis-ci -coverprofile=coverage.out