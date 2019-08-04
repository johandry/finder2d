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

## Installing `finder2d`

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

## Running `finder2d` in CLI mode

After download or build `finder2d` you can execute the binary or run the container. Examples:

```bash
./bin/finder2d \
  --source test_data/image_with_cats.txt \
  --target test_data/perfect_cat_image.txt \
  -p 80
```

```bash
docker run --rm \
    -v $(pwd)/test_data:/data \
    -e FINDER2D_SOURCE="/data/image_with_cats.txt" \
    johandry/finder2d \
    --target /data/perfect_cat_image.txt \
    -p 80
```

Using the docker container requires to mount a volume with the source and target matrix files, to be used with the parameters `--source` and `--target` or the environment variables `FINDER2D_SOURCE` and `FINDER2D_TARGET`.

The `finder2d` has the following parameters in flags or environment variables:

- `--source` or `FINDER2D_SOURCE`: is the source matrix file. The given image or target matrix will be searched into the frame or source matrix. It's required in CLI mode but not in Service mode.
- `--target` or `FINDER2D_TARGET`:  is the target matrix file. If set `finder2d` is executed in CLI mode. 
- `--on` or `FINDER2D_ON`: is the character in the given matrixes to identify a one or on bit of the image. The default value is `+`.
- `--off` or `FINDER2D_OFF`: is the character in the given matrixes to identify a one or on bit of the image. The default value is an space character.
- `-p` or `FINDER2D_PERCENTAGE`: is the matching percentage. The finder will find multiple matches, some of them are noise. The higher the percentage the more the image is equal to the found match. The default value is `50.0`. With the examples matrix the best results are with percentages **61%**
- `-d` or `FINDER2D_DELTA`: is the matches blurry delta. Read below the Delta section. The default delta value is **1**

For more information use `--help`

### Delta

The finder finds multiple matches for the same image/pattern found, all these matches are near by 1, 2, or more bits. Just like a blurry image, all the blurry images are one next to the other in multiple directions.

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

## Running `finder2d` in server mode

Execute `finder2d` either as a binary or in a container without the flag `--target`. The frame or source matrix file with the flag `--source` is optional, if no source file is provided it has to load it before any other action.

```bash
FINDER2D_SOURCE=/test_data/image_with_cats.txt ./bin/finder2d
```

```bash
docker run --rm \
    -v $(pwd)/test_data:/data \
    -p 8080:8080 \
    johandry/finder2d \
    --source /data/image_with_cats.txt
```

In a different terminal access the server either to the gRPC or REST/HTTP API using `grpcurl` or `curl` with `jq` for better JSON formatting:

```bash
curl -s "http://localhost:8080/api/v1/matrixes/0" | jq
```

```bash
grpcurl -plaintext -d '{"name": 0}' localhost:8080 finder2d.v1.Finder2D.GetMatrix
```

To install `grpcurl` or `jq` on MacOS execute `brew install grpcurl jq`.

To list all the available gRPC methods from the Finder2D service, use:

```bash
grpcurl -plaintext localhost:8080 list finder2d.v1.Finder2D
```

To access the Swagger API definition open [localhost:8080/api/v1/swagger/finder2d.json](http://localhost:8080/api/v1/swagger/finder2d.json) or use the Swagger Docker container to see it with Swagger UI:

```bash
docker run --rm \
  -p 80:8080 \
  -e API_URL=http://localhost:8080/api/v1/swagger/finder2d.json \
  swaggerapi/swagger-ui
```

Then navigate to [localhost](http://localhost).

You can also import the Swagger into Postman for a better API interaction and testing.

## Running `finder2d` with Docker Compose

In the project repository there is a `docker-compose.yml` file. This file can be used to start the Finder2D server and the Swagger UI.

Execute the following commands:

```bash
docker-compose up -d
docker-compose ps
```

Navigate to [localhost](http://localhost) to see the Swagger UI or [localhost:8080/api/v1/swagger/finder2d.json](http://localhost:8080/api/v1/swagger/finder2d.json) to read the Swagger JSON API definition.

Execute the commands to access the API using `curl` or `grpcurl` (refer to the [API version 1](#API_version_1) section below). For example:

```bash
# Load the image
img=$(awk '{printf "%s\\n" , $0}' test_data/perfect_cat_image.txt)
grpcurl -plaintext \
  -d '{"api": "v1", "name": 1, "matrix": {"content": "'$img'"}}' \
  localhost:8080 finder2d.v1.Finder2D.LoadMatrix

# Search the image in the frame
grpcurl -plaintext \
  -d '{"api": "v1", "percentage": 70.5, "delta": 1}' \
  localhost:8080 finder2d.v1.Finder2D.Search

# Get all the matches
curl -s "http://localhost:8080/api/v1/matches" | jq

# Get match number 2 and the matrix found
curl -s "http://localhost:8080/api/v1/matches/1" | jq

# Get the 2nd image found
curl -s "http://localhost:8080/api/v1/matches/1" | \
  jq '.matrix.content' | \
  awk '{gsub(/\\n/,"\n")}1'
```

When you are done, execute the following commands to destroy the services.

```bash
docker-compose stop
docker-compose down
```

## API version 1

### GetMatrix

The gRPC method `GetMatrix` is to request the frame or source matrix and the image or target matrix. The frame is identified by a `0` and the image by a `1`. The received object has the matrix content and size.

The REST/HTTP route is `/api/v1/matrixes/{name}` with the HTTP method `GET`.

Using `curl` and `jq`:

```bash
# Requesting the frame
curl -s "http://localhost:8080/api/v1/matrixes/0" | jq
# Requesting the image
curl -s "http://localhost:8080/api/v1/matrixes/1" | jq
```

Using `grpcurl`:

```bash
# Requesting the frame
grpcurl -plaintext -d '{"name": 0}' localhost:8080 finder2d.v1.Finder2D.GetMatrix
# Requesting the image
grpcurl -plaintext -d '{"name": 1}' localhost:8080 finder2d.v1.Finder2D.GetMatrix
```

Sample Output:

```json
{
  "api": "v1",
  "name": "SOURCE",
  "matrix": {
    "width": 15,
    "height": 15,
    "content": "+             +\n+++         +++\n ++++++    .......        \n"
  }
}
```

### LoadMatrix

The gRPC method `LoadMatrix` is to load into the Finder2D the frame or source matrix and the image or target matrix.

The request is a JSON object with the matrix type (`"name"`) and the matrix object only with the content (`"matrix": {"content": "...."}`). The frame is identified by a `0` and the image by a `1`.

The response only contain the API version number, if there was an error it will be in the response.

The REST/HTTP route is `/api/v1/matrixes/{name}` with the HTTP method `POST`.

Using `curl` and `jq`:

```bash
# Load the content into an environment variable, replacing '\n' for '\\n'
img=$(awk '{printf "%s\\n" , $0}' test_data/perfect_cat_image.txt)

# Loading the frame
curl -s \
  -d '{"api": "v1", "name": 0, "matrix": {"content": "'$img'"}}' \
  -H "Content-Type: application/json" \
  -X POST  "http://localhost:8080/api/v1/matrixes/0" | jq

# Loading the image
curl -s \
  -d '{"api": "v1", "name": 1, "matrix": {"content": "'$img'"}}' \
  -H "Content-Type: application/json" \
  -X POST  "http://localhost:8080/api/v1/matrixes/1" | jq
```

Using `grpcurl`:

```bash
# Load the content into an environment variable, replacing '\n' for '\\n'
img=$(awk '{printf "%s\\n" , $0}' test_data/image_with_cats.txt)

# Loading the frame
grpcurl -plaintext \
  -d '{"api": "v1", "name": 0, "matrix": {"content": "'$img'"}}' \
  localhost:8080 finder2d.v1.Finder2D.LoadMatrix

# Loading the image
grpcurl -plaintext \
  -d '{"api": "v1", "name": 1, "matrix": {"content": "'$img'"}}' \
  localhost:8080 finder2d.v1.Finder2D.LoadMatrix
```

Sample Output:

```json
{
  "api": "v1"
}
```

Or an error, in this example the `\n` was not replaced by `\\n`:

```bash
{
  "error": "invalid character '\\n' in string literal",
  "message": "invalid character '\\n' in string literal",
  "code": 3,
  "details": []
}
```

### Search

The gRPC method `Search` is used to search the image or target matrix in the frame or source matrix using the `SearchSimple()` method of the Finder2D.

It's important to remember to load the target matrix before execute a search, otherwise an error will be received.

The request is a JSON object with the percentage (`"percentage"`) and the delta (`"delta"`) values. The response has the total number of matches found (`"total_matches"`).

The REST/HTTP route is `/api/v1/search` with the HTTP method `POST`.

 Using `curl` and `jq`:

```bash
# Load the image
img=$(awk '{printf "%s\\n" , $0}' test_data/perfect_cat_image.txt)
curl -s \
  -d '{"api": "v1", "name": 1, "matrix": {"content": "'$img'"}}' \
  -H "Content-Type: application/json" \
  -X POST  "http://localhost:8080/api/v1/matrixes/1" | jq

# Search the image in the frame
curl -s \
  -d '{"api": "v1", "percentage": 70.5, "delta": 1}' \
  -H "Content-Type: application/json" \
  -X POST  "http://localhost:8080/api/v1/search" | jq
```

Using `grpcurl`:

```bash
# Load the image like in previous example
grpcurl -plaintext \
  -d '{"api": "v1", "name": 1, "matrix": {"content": "'$img'"}}' \
  localhost:8080 finder2d.v1.Finder2D.LoadMatrix

# Search the image in the frame
grpcurl -plaintext \
  -d '{"api": "v1", "percentage": 70.5, "delta": 1}'  \
  localhost:8080 finder2d.v1.Finder2D.Search
```

Sample Output:

```json
{
  "api": "v1",
  "total_matches": 6
}
```

### GetMatches

The gRPC method `GetMatches` is used retrieve all the matches found from a previous search. This method will return an empty list if the search is not done before.

The request is a JSON object only with the API version. The response is a JSON object with an array/list of matches (`"matches"`). Each Match is a JSON object with the coordinates (`"x"`, `"y"`) and the matching percentage (`"percentage"`).

The REST/HTTP route is `/api/v1/matches` with the HTTP method `GET`.

Using `curl` and `jq`:

```bash
# Get the list of found matches
curl -s "http://localhost:8080/api/v1/matches" | jq
```

Using `grpcurl`:

```bash
# Get the list of found matches
grpcurl -plaintext localhost:8080 finder2d.v1.Finder2D.GetMatches
```

Sample Output:

```json
{
  "api": "v1",
  "matches": [
    {
      "x": 80,
      "y": 0,
      "percentage": 99.111115
    },
....
    {
      "x": 84,
      "y": 84,
      "percentage": 98.666664
    }
  ]
}
```

### GetMatch

The gRPC method `GetMatch` return the requested match identified by it's index in the list. This method will return an error if the index is out of range.

The request is a JSON object  with the index or ID (`"id"`) of the required match. The response is a JSON object with the match (`"match"`) and the Matrix (`"matrix"`). The Match is a JSON object with the coordinates (`"x"`, `"y"`) and the matching percentage (`"percentage"`). The Matrix is a JSON object with the width (`"width"`), height (`"height"`) and the content of the matrix (`"content"`) as a string.

The REST/HTTP route is `/api/v1/match/{id}` with the HTTP method `GET`.

Using `curl` and `jq`:

```bash
# Get the list of found matches
curl -s "http://localhost:8080/api/v1/match/0" | jq
```

Using `grpcurl`:

```bash
# Get the list of found matches
grpcurl -plaintext \
  -d '{"api": "v1", "id": 1}' \
  localhost:8080 finder2d.v1.Finder2D.GetMatch
```

Sample Output:

```json
{
  "api": "v1",
  "match": {
    "x": 84,
    "y": 84,
    "percentage": 98.666664
  },
  "matrix": {
    "width": 15,
    "height": 15,
    "content": "+             +\n+++         +++\n +++++++ ......             \n"
  }
}
```

## TODO

- [x] Implement the LoadMatrix gRPC method
- [x] Define and implement the Search gRPC method
- [x] Create the Docker Compose for the services
- [x] Start server mode by default when the image file is not provided
- [x] Accept parameter in environment variables as well as with arguments
- [ ] Create the Kubernetes Manifest for the services
- [ ] Allow to load multiple targets
- [ ] Create a DB service to store the matrixes
- [ ] Make a client `finder2dctl`
- [ ] Make a client in other language (Python? Ruby?)
- [ ] Build a UI
- [ ] Go serverless 
- [ ] Security: Implement TLS
- [ ] Security: Implement JWT
- [ ] CI/CD with Travis and/or others (CircleCI?)

## Other algorithms to improve the search

**2D Convolution**:

- [Convolution Theorem](https://en.wikipedia.org/wiki/Convolution_theorem)
- [Convolution 2D Example](http://www.songho.ca/dsp/convolution/convolution2d_example.html)
- [Two Dimensional Convolution in Image Processing](https://www.allaboutcircuits.com/technical-articles/two-dimensional-convolution-in-image-processing/)
- [Understanding Convolutions for Deep Learning](https://towardsdatascience.com/intuitively-understanding-convolutions-for-deep-learning-1f6f42faee1)
- [Efficient 2D Approximate Matching of Non-rectangular Figures](https://www.cs.rutgers.edu/~farach/pubs/HalfRec.pdf)

**Rabin–Karp algorithm**:

- [Rabin–Karp Algorithm](https://en.wikipedia.org/wiki/Rabin–Karp_algorithm)
- [Rabin Fingerprint](https://en.wikipedia.org/wiki/Rabin_fingerprint)

**Boyer–Moore algorithm**:

- [Boyer–Moore String Search Algorithm](https://en.wikipedia.org/wiki/Boyer–Moore_string-search_algorithm)
