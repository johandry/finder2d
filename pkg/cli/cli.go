package cli

import (
	"fmt"
	"os"

	"github.com/johandry/finder2d"
)

// Exec executes the CLI mode, loading the matrixes and printing the matches
func Exec(sourceFileName, targetFileName, zero, one string, percentage float64, delta int, format string) error {
	// Open files
	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		return fmt.Errorf("fail to open the frame file %q. %s", sourceFileName, err)
	}
	targetFile, err := os.Open(targetFileName)
	if err != nil {
		return fmt.Errorf("fail to open the image file %q. %s", targetFileName, err)
	}

	// Load matrixes from files
	f := finder2d.New([]byte(one)[0], []byte(zero)[0], percentage, delta)
	if err := f.LoadSource(sourceFile); err != nil {
		return fmt.Errorf("fail to load the source file %q. %s", sourceFileName, err)
	}
	if err := f.LoadTarget(targetFile); err != nil {
		return fmt.Errorf("fail to load the target file %q. %s", sourceFileName, err)
	}

	// DEBUG:
	// x, y := f.Source.Size()
	// fmt.Printf("Source (%dx%d): \n%s\n", x, y, f.Source)
	// x, y = f.Target.Size()
	// fmt.Printf("Target (%dx%d): \n%s\n", x, y, f.Target)
	// fmt.Println("Finding matches ...")

	f.SearchSimple()

	fmt.Println(f.Stringf(format))

	return nil
}
