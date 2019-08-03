APP_NAME = finder2d
DOKR_USR = johandry

IMG_NAME = $(DOKR_USR)/$(APP_NAME)

GO111MODULE = on

default: test build

docker: docker-build docker-push

build:
	go build -o bin/$(APP_NAME) cmd/main.go

test:
	go test ./...

mod:
	go mod tidy
	go mod download

docker-build:
	docker build -t $(IMG_NAME) .

docker-run:
	docker run --rm \
		-v $(PWD)/test_data:/data \
		$(IMG_NAME) \
		--frame /data/image_with_cats.txt \
		--image /data/perfect_cat_image.txt

docker-push:
	docker push $(IMG_NAME)

## Do not use the following rules unless you know what you are doing

mod-redo:
	$(RM) go.mod go.sum
	go clean -modcache
	go mod init github.com/johandry/$(APP_NAME)
	go mod download