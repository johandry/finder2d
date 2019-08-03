# Finder 2D

Finder 2D is a package to recognize or find a target string or pattern in a source string. Both, the pattern and the source string, are bi-dimensional arrays.

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

## How to use `finder2d` in CLI mode

After download or build `finder2d` you can execute the binary or run the container. Examples:

```bash
./bin/finder2d --frame test_data/image_with_cats.txt --image test_data/perfect_cat_image.txt
```

```bash
docker run --rm \
    -v $(pwd)/test_data:/data \
    johandry/finder2d \
    --frame /data/image_with_cats.txt \
    --image /data/perfect_cat_image.txt
```

Using the docker container requires to mount a volume with the source and target matrix files, to be used with the parameters `--frame` and `--image`.

The `finder2d` has the following parameters:

- `--frame`: (required) is the source matrix file. The given image or target matrix will be searched into the frame or source matrix.
- `--image`: (required) is the target matrix file.
- `--on`: is the character in the given matrixes to identify a one or on bit of the image. The default value is `+`.
- `--off`: is the character in the given matrixes to identify a one or on bit of the image. The default value is an space ` `.
- `-p`: is the matching percentage. The finder will find multiple matches, some of them are noise. The higher the percentage the more the image is equal to the found match. The default value is `50.0`. With the examples matrix the best results are with percentages **61%**
- `-d`: is the matches blurry delta. Read below the Delta section. The default delta value is **1**

For more information use `--help`

### Delta

The finder finds multiple matches for the same image/patter found, all these matches are near by 1, 2, or more bits. Just like a blurry image, all the blurry images are one next to the other in multiple directions.

The delta value is used to reduce all these blurry images to one. The higher the delta, the less blurry images will found or, in the best case, all these blurry images will be reduced to one.

However, in the source matrix may be multiple images/patters near that may be falsely identified as blurry images. If this is the case, you need to reduce the delta, trying to not identify them as blurry images and identify them as just images one near to the other.

Using the default matching percentage (50%) the best results are found with **delta = 5**. Read the following table with the minimum percentage for possible delta values:

| Delta | Minimum Percentage |
| ----- | ----------- |
| 6     | 48%            |
| 5     | 48%            |
| 4     | 51%            |
| 3     | 51%            |
| 2     | 54%            |
| 1     | 61%            |

The maximum percentage is **96%** for any delta.

## How to use `finder2d` in server mode
