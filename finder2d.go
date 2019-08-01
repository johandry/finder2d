package finder2d

import (
	"bytes"
	"fmt"
	"io"
)

const (
	uno  = `◻️`
	cero = `◼️️`
)

// Matrix represents a 2D array
type Matrix struct {
	Content    [][]int
	maxX, maxY int
}

// Match represents the coordinate the target matrix was found in the source
// matrix and the percentage match
type Match struct {
	X, Y       int
	Percentage float64
}

// Finder2D is the struct used to find a 2D pattern into a 2D source matrix
type Finder2D struct {
	Target     *Matrix
	Source     *Matrix
	one, zero  byte
	Matches    []Match
	Percentage float64
}

func (r *Match) String() string {
	return fmt.Sprintf("(%d,%d,%f)", r.X, r.Y, r.Percentage)
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
	var b bytes.Buffer
	for _, row := range m.Content {
		for _, v := range row {
			switch v {
			case 0:
				b.WriteString(cero)
			case 1:
				b.WriteString(uno)
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
	if m.maxX != m1.maxX || m.maxY != m1.maxY {
		return 0, fmt.Errorf("matrix to compare with is not the same size (%d,%d) != (%d,%d)", m.maxX, m.maxY, m1.maxX, m1.maxY)
	}
	var same int
	for y := 0; y < m.maxY; y++ {
		for x := 0; x < m.maxY; x++ {
			if m.Content[y][x] == m1.Content[y][x] {
				same++
			}
		}
	}

	return float64(same/(m.maxX*m.maxY)) * 100, nil
}

// New create an empty Finder 2D
func New(one, zero byte, percentage float64) *Finder2D {
	return &Finder2D{
		one:        one,
		zero:       zero,
		Percentage: percentage,
	}
}

// LoadSource loads the source from a reader replacing the cell value given in
// `one` for `1` and `zero` for `0`
func (f *Finder2D) LoadSource(r io.Reader) error {
	m, err := LoadMatrix(r, f.one, f.zero)
	if err != nil {
		return err
	}

	f.Source = m
	return nil
}

// LoadTarget loads the target from a reader replacing the cell value given in
// `one` for `1` and `zero` for `0`
func (f *Finder2D) LoadTarget(r io.Reader) error {
	m, err := LoadMatrix(r, f.one, f.zero)
	if err != nil {
		return err
	}

	f.Target = m
	return nil
}

// SearchSimple find the occurences of the target in the source and the percentage
// match in the simplest way which is to iterate thru the entire matrix searching
// for the pattern, storing the match when the match percentage is higher than
// the required
func (f *Finder2D) SearchSimple() {
	maxX, maxY := f.Source.Size()
	width, height := f.Target.Size()

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			sample := f.Source.Sample(x, y, width, height)
			if sample == nil {
				break
			}
			p, err := sample.Compare(f.Target)
			if err != nil {
				return
			}
			if p >= f.Percentage {
				f.Matches = append(f.Matches, Match{
					X:          x,
					Y:          y,
					Percentage: p,
				})
			}
		}
	}
}

func (f *Finder2D) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	first := true
	for _, m := range f.Matches {
		if first {
			b.WriteString(m.String())
			first = false
		} else {
			b.WriteString("," + m.String())
		}
	}
	b.WriteString("]")
	return b.String()
}
