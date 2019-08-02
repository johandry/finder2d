package finder2d

import (
	"bytes"
	"fmt"
	"io"
)

// DefaultMinMatchPercentage default minimun match percentage. Any match
// percentage lower than this number is not considered
const DefaultMinMatchPercentage float64 = 50.0

// Default values for a one and a zero in a matrix
const (
	DefaultOne  = `+`
	DefaultZero = ` `
)

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

// New create an empty Finder 2D
func New(one, zero byte, percentage float64) *Finder2D {
	if percentage == 0 {
		percentage = DefaultMinMatchPercentage
	}
	// if both are `0`
	if one+zero == 0 {
		one = []byte(DefaultOne)[0]
		zero = []byte(DefaultZero)[0]
	}
	return &Finder2D{
		one:        one,
		zero:       zero,
		Percentage: percentage,
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
