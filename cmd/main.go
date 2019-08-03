package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/johandry/finder2d"
)

func main() {
	var frameFileName, imageFileName string
	var zero, one string
	var percentage float64
	var delta int

	flag.StringVar(&frameFileName, "frame", "", "frame or source matrix file (required)")
	flag.StringVar(&imageFileName, "image", "", "image or target matrix file (required)")
	flag.StringVar(&zero, "off", " ", "matrix character that represents a zero or off bit")
	flag.StringVar(&one, "on", "+", "matrix character that represents a one or on bit")
	flag.Float64Var(&percentage, "p", 50.0, "matching percentage")
	flag.IntVar(&delta, "d", 1, "matches blurry delta, the higher it is the less blurry patterns will find")
	flag.Parse()

	if len(frameFileName) == 0 {
		log.Fatalf("frame file is required. Use the flag '--frame'")
	}
	if len(imageFileName) == 0 {
		log.Fatalf("image file is required. Use the flag '--image'")
	}

	frameFile, err := os.Open(frameFileName)
	if err != nil {
		log.Fatalf("fail to open the frame file %q. %s", frameFileName, err)
	}
	imageFile, err := os.Open(imageFileName)
	if err != nil {
		log.Fatalf("fail to open the image file %q. %s", imageFileName, err)
	}

	f := finder2d.New([]byte(one)[0], []byte(zero)[0], percentage, delta)
	if err := f.LoadSource(frameFile); err != nil {
		log.Fatalf("fail to load the frame file %q. %s", frameFileName, err)
	}
	if err := f.LoadTarget(imageFile); err != nil {
		log.Fatalf("fail to load the image file %q. %s", frameFileName, err)
	}

	x, y := f.Source.Size()
	fmt.Printf("Frame (%dx%d): \n%s\n", x, y, f.Source)

	x, y = f.Target.Size()
	fmt.Printf("Image (%dx%d): \n%s\n", x, y, f.Target)

	fmt.Println("Finding matches ...")
	f.SearchSimple()

	n := len(f.Matches)

	fmt.Printf("Matches (%d): \n%s\n", n, strings.Replace(f.String(), "),(", ")\n(", -1))
}
