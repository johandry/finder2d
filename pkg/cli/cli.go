/*
Copyright The Finder2D Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"fmt"
	"os"

	"github.com/johandry/finder2d"
)

// Execute executes the CLI mode, loading the matrixes and printing the matches
func Execute(sourceFileName, targetFileName, zero, one string, percentage float64, delta int, format string) error {
	switch format {
	case "", "text", "matrix", "json":
	default:
		return fmt.Errorf("unknown output format %q. Available options are: 'json', 'text' or 'matrix'", format)
	}

	if len(sourceFileName) == 0 {
		return fmt.Errorf("source file is required")
	}

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

	if err := f.SearchSimple(); err != nil {
		return fmt.Errorf("failed to search the target matrix. %s", err)
	}

	fmt.Println(f.Stringf(format))

	return nil
}
