APP_NAME = finder2d
IMG_NAME = johandry/$(APP_NAME)

GO111MODULE = on

default: test build

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


## Do not use the following rules unless you know what you are doing

mod-redo:
	$(RM) go.mod go.sum
	go clean -modcache
	go mod init github.com/johandry/$(APP_NAME)
	go mod download