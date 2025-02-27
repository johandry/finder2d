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

package finder2d

import (
	"bytes"
	"fmt"
	"io"
)

const (
	uno  = "\033[44m \033[0m" // Blue // `◻️`
	cero = "\033[40m \033[0m" // Black // `◼️️`
)

// Matrix represents a 2D array
type Matrix struct {
	Content    [][]int
	maxX, maxY int
}

// LoadMatrix create a matrix from a reader
func LoadMatrix(r io.Reader, one, zero byte) (*Matrix, error) {
	m := &Matrix{}
	err := m.Load(r, one, zero)
	return m, err
}

// Load loads a matrix from a reader replacing the cell value given in `one`
// for `1` and `zero` for `0`
func (m *Matrix) Load(r io.Reader, one, zero byte) error {
	var x, y int
	NL := []byte("\n")[0]

	defer func() {
		if m.Content == nil || m.maxX+m.maxY == 0 {
			m.Content = nil
			m.maxX, m.maxY = 0, 0
		}
	}()

	m.Content = [][]int{[]int{}}
	buf := make([]byte, 10)
	for {
		n, errRead := r.Read(buf)
		for i := 0; i < n; i++ {
			switch buf[i] {
			case byte(NL):
				if m.maxX == 0 {
					m.maxX = x
				} else if x > m.maxX {
					m.Content = nil
					return fmt.Errorf("source width = %d, especified by the first row, is larger at line #%d (%d)", m.maxX, y, x)
				} else if x < m.maxX {
					for i := 0; i < m.maxX-x; i++ {
						m.Content[y] = append(m.Content[y], 0)
					}
				}
				x = 0
				y = y + 1
				m.Content = append(m.Content, []int{})
			case one:
				x = x + 1
				m.Content[y] = append(m.Content[y], 1)
			case zero:
				x = x + 1
				m.Content[y] = append(m.Content[y], 0)
			default:
				m.Content = nil
				return fmt.Errorf("found invalid value in the source matrix %q", buf[i])
			}
		}
		if errRead == io.EOF {
			break
		}
	}

	if y > 0 && len(m.Content[y]) == 0 {
		m.Content = m.Content[:len(m.Content)-1]
		y = y - 1
	}
	if y != 0 {
		y = y + 1
	}
	m.maxY = y

	return nil
}

// Size returns the size of the matrix
func (m *Matrix) Size() (int, int) {
	return m.maxX, m.maxY
}

func (m *Matrix) String() string {
	return m.Sprintf(cero, uno)
}

// Sprintf returns a string representing the matrix but using the given `one`,
// `zero` strings to represent the one and zero values
func (m *Matrix) Sprintf(zero, one string) string {
	var b bytes.Buffer
	for _, row := range m.Content {
		for _, v := range row {
			switch v {
			case 0:
				b.WriteString(zero)
			case 1:
				b.WriteString(one)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Sample gets a sample of the Matrix from the coordinates (x,y) and the given
// width and height
func (m *Matrix) Sample(x, y, w, h int) *Matrix {
	if x+w > m.maxX || y+h > m.maxY {
		return nil
	}
	if w+h == 0 {
		return &Matrix{}
	}
	sample := make([][]int, h)
	for yi := 0; yi < h; yi++ {
		sample[yi] = make([]int, w)
		for xi := 0; xi < w; xi++ {
			sample[yi][xi] = m.Content[y+yi][x+xi]
		}
	}

	return &Matrix{
		Content: sample,
		maxX:    w,
		maxY:    h,
	}
}

// Compare returns the matching percentage between this and the given matrix
func (m *Matrix) Compare(m1 *Matrix) (float64, error) {
	if m.maxX+m.maxY == 0 || m1.maxX+m1.maxY == 0 {
		return 0, nil
	}
	if m.maxX != m1.maxX || m.maxY != m1.maxY {
		return 0, fmt.Errorf("matrix to compare with is not the same size (%d,%d) != (%d,%d)", m.maxX, m.maxY, m1.maxX, m1.maxY)
	}
	var same float64
	for y := 0; y < m.maxY; y++ {
		for x := 0; x < m.maxY; x++ {
			if m.Content[y][x] == m1.Content[y][x] { // && (m.Content[y][x] == 1) {
				same++
			}
		}
	}

	return same / float64(m.maxX*m.maxY) * 100.0, nil
}
