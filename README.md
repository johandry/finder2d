# Finder 2D

Finder 2D is a package to recognize or find a target string or pattern in a source string. Both, the pattern and the source string, are bidimensional arrays.

This package is created to solve the following exercise:

```text
Couple of years back, Google create a brain simulator in their labs, that learned
how to recognize cats in YouTube videos (http://www.wired.co.uk/news/archive/2012-06/26/google-brain-recognises-cats)

Your task is a stripped down version of the same idea.
Imagine you have a single video frame (image_with_cats.txt) with some cat images.
In addition, you have a perfect cat image (perfect_cat_image.txt). Your goal is
to find the cats in the video frame. You can return the position of the cat in the
image matrix, and the percentage match. You can optionally provide a threshold
match value (like 85% confidence this is the image you are looking for).

You should expose the functionality above as a JSON REST service
that takes the video frame on the input (as a text matrix), and the threshold
match value, and returns the output again as a REST response. You don't have to
pass the image you are looking for (the cat image) in the REST interface - it remains constant.

The above should be implemented in Java or Go. Feel free to use frameworks or libraries as needed to accomplish the task.

Note: The video frame could be "noisy" - you may not find the perfect cat image,
therefore, you have to deal with certain possibilities what you have found is the
correct answer.
```

## How to get `finder2d`

To get the binary use `go get`:

```bash
go get github.com/johandry/finder2d
```

Or build it from source after cloning this repository:

```bash
git clone https://github.com/johandry/finder2d
cd finder2d
make
```

You can also get the Docker image:

```bash
docker pull johandry/finder2d
docker run --rm johandry/finder2d
```

Or, build and run the Docker image from source:

```bash
make docker-build
make docker-run
```

## How to use (CLI mode)

Build and execute:

```bash
make build
./bin/finder2d --frame test_data/image_with_cats.txt --image test_data/perfect_cat_image.txt
```

Or using Docker:

```bash
docker run --rm \
    -v $(pwd)/test_data:/data \
    johandry/finder2d \
    --frame /data/image_with_cats.txt \
    --image /data/perfect_cat_image.txt
```
