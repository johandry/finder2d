APP_NAME = finder2d
DOKR_USR = johandry

IMG_NAME = $(DOKR_USR)/$(APP_NAME)

GO111MODULE = on

# first rule, so this will be done when executed `make` without rules
default: mod test build

docker: docker-build docker-push

build:
	go build -o bin/$(APP_NAME) cmd/main.go

test:
	go test ./...

# remove the unused modules and download the missing ones
mod:
	go mod tidy
	go mod download

# build the container with the application
docker-build: test
	docker build -t $(IMG_NAME) .

# push the built image to the docker registry. Make sure you set the correct user
docker-push:
	docker push $(IMG_NAME)

api-build:
	$(MAKE) -C api build

##################################################################### 
# Danger Zone:                                                      #
# Do not use the following rules unless you know what you are doing #
#####################################################################

# install in your local system, at $GOPATH, the Go packages, binaries and proto 
# definitions required to build the API. This may cause problems if the downloaded 
# version are not the correct one for the code.
# TODO (johandry): consider to use `third_party` directory to store all the 
# dependencies
api-dependencies: 
	$(MAKE) -C api dependencies

# clean up all the modules in your system and re-do the modules file, downloading 
# all the needed modules to your system
mod-redo:
	$(RM) go.mod go.sum
	go clean -modcache
	go mod init github.com/johandry/$(APP_NAME)
	go mod download